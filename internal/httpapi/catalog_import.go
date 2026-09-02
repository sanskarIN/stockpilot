package httpapi

import (
	"net/http"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvimport"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const maxProductImportBytes = 5 << 20

func (a *API) validateProductImport(w http.ResponseWriter, r *http.Request) {
	result, err := a.parseProductImport(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	errors, valid := a.validateProductRows(r, result)
	writeJSON(w, http.StatusOK, map[string]any{"dryRun": true, "validRows": len(valid), "errorRows": len(errors), "errors": errors, "valid": valid})
}

func (a *API) importProducts(w http.ResponseWriter, r *http.Request) {
	importer, ok := a.catalog.(repository.ProductBatchImporter)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "product import persistence is not available"})
		return
	}

	result, err := a.parseProductImport(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	errors, _ := a.validateProductRows(r, result)
	if len(errors) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "product import validation failed", "errors": errors})
		return
	}
	if len(result.Rows) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "CSV contains no product rows"})
		return
	}

	products := make([]domainProductAlias, 0, len(result.Rows))
	for _, row := range result.Rows {
		product := row.Product
		if strings.TrimSpace(product.ID) == "" || product.ID == "pending" {
			product.ID, err = idgen.New("prd")
			if err != nil {
				writeDomainError(w, err)
				return
			}
		}
		products = append(products, domainProductAlias{row: row.Row, product: product})
	}

	batch := make([]domain.Product, 0, len(products))
	for _, item := range products { batch = append(batch, item.product) }
	if err := importer.ImportProducts(r.Context(), batch); err != nil {
		writeDomainError(w, err)
		return
	}

	a.recordAudit(r.Context(), authenticatedActorID(r), "products.imported", "product_import", requestIDFromContext(r.Context()), map[string]any{"count": len(batch)})
	writeJSON(w, http.StatusCreated, map[string]any{"imported": len(batch), "products": batch})
}

// domainProductAlias keeps row provenance available for future import result
// reporting without exposing CSV parser internals from the HTTP layer.
type domainProductAlias struct { row int; product domain.Product }

func (a *API) parseProductImport(r *http.Request) (csvimport.ValidationResult, error) {
	r.Body = http.MaxBytesReader(nil, r.Body, maxProductImportBytes)
	if err := r.ParseMultipartForm(maxProductImportBytes); err != nil {
		return csvimport.ValidationResult{}, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return csvimport.ValidationResult{}, err
	}
	defer file.Close()
	return csvimport.ParseProducts(file)
}

func (a *API) validateProductRows(r *http.Request, result csvimport.ValidationResult) ([]csvimport.RowError, []map[string]any) {
	categories, err := a.catalog.ListCategories(r.Context())
	if err != nil { return []csvimport.RowError{{Row: 0, Message: err.Error()}}, nil }
	suppliers, err := a.catalog.ListSuppliers(r.Context(), false)
	if err != nil { return []csvimport.RowError{{Row: 0, Message: err.Error()}}, nil }
	categoryIDs := make(map[string]struct{}, len(categories)); for _, item := range categories { categoryIDs[item.ID] = struct{}{} }
	supplierIDs := make(map[string]struct{}, len(suppliers)); for _, item := range suppliers { supplierIDs[item.ID] = struct{}{} }

	errors := append([]csvimport.RowError(nil), result.Errors...)
	valid := make([]map[string]any, 0, len(result.Rows))
	for _, row := range result.Rows {
		product := row.Product
		rowValid := true
		if product.CategoryID != "" { if _, ok := categoryIDs[product.CategoryID]; !ok { errors = append(errors, csvimport.RowError{Row: row.Row, Message: "category_id does not reference an existing category"}); rowValid = false } }
		if product.SupplierID != "" { if _, ok := supplierIDs[product.SupplierID]; !ok { errors = append(errors, csvimport.RowError{Row: row.Row, Message: "supplier_id does not reference an existing supplier"}); rowValid = false } }
		existing, lookupErr := a.catalog.ListProducts(r.Context(), repositoryProductExactFilter(product.SKU))
		if lookupErr != nil { errors = append(errors, csvimport.RowError{Row: row.Row, Message: lookupErr.Error()}); rowValid = false }
		for _, current := range existing { if strings.EqualFold(current.SKU, product.SKU) { errors = append(errors, csvimport.RowError{Row: row.Row, Message: "SKU already exists"}); rowValid = false; break } }
		if rowValid { valid = append(valid, map[string]any{"row": row.Row, "sku": product.SKU, "name": product.Name}) }
	}
	return errors, valid
}

func repositoryProductExactFilter(sku string) repository.ProductFilter { return repository.ProductFilter{Query: strings.TrimSpace(sku), Limit: 20} }

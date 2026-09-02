package httpapi

import (
	"net/http"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvimport"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const maxProductImportBytes = 5 << 20

func (a *API) validateProductImport(w http.ResponseWriter, r *http.Request) {
	result, err := a.parseProductImport(w, r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	errors, valid, err := a.validateProductRows(r, result)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"dryRun": true, "validRows": len(valid), "errorRows": len(errors), "errors": errors, "valid": valid})
}

func (a *API) importProducts(w http.ResponseWriter, r *http.Request) {
	importer, ok := a.catalog.(repository.ProductBatchImporter)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "product import persistence is not available"})
		return
	}

	result, err := a.parseProductImport(w, r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	errors, _, err := a.validateProductRows(r, result)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if len(errors) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "product import validation failed", "errors": errors})
		return
	}
	if len(result.Rows) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "CSV contains no product rows"})
		return
	}

	products := make([]domain.Product, 0, len(result.Rows))
	for _, row := range result.Rows {
		product := row.Product
		if strings.TrimSpace(product.ID) == "" || product.ID == "pending" {
			product.ID, err = idgen.New("prd")
			if err != nil {
				writeDomainError(w, err)
				return
			}
		}
		products = append(products, product)
	}

	if err := importer.ImportProducts(r.Context(), products); err != nil {
		writeDomainError(w, err)
		return
	}

	a.recordAudit(r.Context(), authenticatedActorID(r), "products.imported", "product_import", requestIDFromContext(r.Context()), map[string]any{"count": len(products)})
	writeJSON(w, http.StatusCreated, map[string]any{"imported": len(products), "products": products})
}

func (a *API) parseProductImport(w http.ResponseWriter, r *http.Request) (csvimport.ValidationResult, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxProductImportBytes)
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

func (a *API) validateProductRows(r *http.Request, result csvimport.ValidationResult) ([]csvimport.RowError, []map[string]any, error) {
	categories, err := a.catalog.ListCategories(r.Context())
	if err != nil {
		return nil, nil, err
	}
	suppliers, err := a.catalog.ListSuppliers(r.Context(), false)
	if err != nil {
		return nil, nil, err
	}
	categoryIDs := make(map[string]struct{}, len(categories))
	for _, item := range categories {
		categoryIDs[item.ID] = struct{}{}
	}
	supplierIDs := make(map[string]struct{}, len(suppliers))
	for _, item := range suppliers {
		supplierIDs[item.ID] = struct{}{}
	}

	errors := append([]csvimport.RowError(nil), result.Errors...)
	valid := make([]map[string]any, 0, len(result.Rows))
	for _, row := range result.Rows {
		product := row.Product
		rowValid := true
		if product.CategoryID != "" {
			if _, ok := categoryIDs[product.CategoryID]; !ok {
				errors = append(errors, csvimport.RowError{Row: row.Row, Message: "category_id does not reference an existing category"})
				rowValid = false
			}
		}
		if product.SupplierID != "" {
			if _, ok := supplierIDs[product.SupplierID]; !ok {
				errors = append(errors, csvimport.RowError{Row: row.Row, Message: "supplier_id does not reference an existing supplier"})
				rowValid = false
			}
		}
		existing, lookupErr := a.catalog.ListProducts(r.Context(), repositoryProductExactFilter(product.SKU))
		if lookupErr != nil {
			return nil, nil, lookupErr
		}
		for _, current := range existing {
			if strings.EqualFold(current.SKU, product.SKU) {
				errors = append(errors, csvimport.RowError{Row: row.Row, Message: "SKU already exists"})
				rowValid = false
				break
			}
		}
		if rowValid {
			valid = append(valid, map[string]any{"row": row.Row, "sku": product.SKU, "name": product.Name})
		}
	}
	return errors, valid, nil
}

func repositoryProductExactFilter(sku string) repository.ProductFilter {
	return repository.ProductFilter{Query: strings.TrimSpace(sku), Limit: 20}
}

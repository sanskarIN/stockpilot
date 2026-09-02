package httpapi

import (
	"net/http"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvimport"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (a *API) validateProductImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "expected a multipart CSV upload"})
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "CSV file field is required"})
		return
	}
	defer file.Close()

	result, err := csvimport.ParseProducts(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	categories, err := a.catalog.ListCategories(r.Context())
	if err != nil { writeDomainError(w, err); return }
	suppliers, err := a.catalog.ListSuppliers(r.Context(), false)
	if err != nil { writeDomainError(w, err); return }
	categoryIDs := make(map[string]struct{}, len(categories)); for _, item := range categories { categoryIDs[item.ID] = struct{}{} }
	supplierIDs := make(map[string]struct{}, len(suppliers)); for _, item := range suppliers { supplierIDs[item.ID] = struct{}{} }

	valid := make([]map[string]any, 0, len(result.Rows))
	errors := append([]csvimport.RowError(nil), result.Errors...)
	for _, row := range result.Rows {
		product := row.Product
		rowValid := true
		if product.CategoryID != "" { if _, ok := categoryIDs[product.CategoryID]; !ok { errors = append(errors, csvimport.RowError{Row: row.Row, Message: "category_id does not reference an existing category"}); rowValid = false } }
		if product.SupplierID != "" { if _, ok := supplierIDs[product.SupplierID]; !ok { errors = append(errors, csvimport.RowError{Row: row.Row, Message: "supplier_id does not reference an existing supplier"}); rowValid = false } }
		existing, lookupErr := a.catalog.ListProducts(r.Context(), repositoryProductExactFilter(product.SKU))
		if lookupErr != nil { writeDomainError(w, lookupErr); return }
		for _, current := range existing { if strings.EqualFold(current.SKU, product.SKU) { errors = append(errors, csvimport.RowError{Row: row.Row, Message: "SKU already exists"}); rowValid = false; break } }
		if rowValid { valid = append(valid, map[string]any{"row": row.Row, "sku": product.SKU, "name": product.Name}) }
	}
	writeJSON(w, http.StatusOK, map[string]any{"dryRun": true, "validRows": len(valid), "errorRows": len(errors), "errors": errors, "valid": valid})
}

func repositoryProductExactFilter(sku string) repository.ProductFilter { return repository.ProductFilter{Query: strings.TrimSpace(sku), Limit: 20} }

package httpapi

import (
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (a *API) listCategories(w http.ResponseWriter, r *http.Request) {
	items, err := a.catalog.ListCategories(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
func (a *API) createCategory(w http.ResponseWriter, r *http.Request) {
	var item domain.Category
	if !decodeJSON(w, r, &item) {
		return
	}
	id, err := idgen.New("cat")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	item.ID = id
	if err := a.catalog.CreateCategory(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "category.created", "category", item.ID, map[string]any{"name": item.Name})
	writeJSON(w, http.StatusCreated, item)
}
func (a *API) listSuppliers(w http.ResponseWriter, r *http.Request) {
	items, err := a.catalog.ListSuppliers(r.Context(), queryBool(r, "active", false))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
func (a *API) createSupplier(w http.ResponseWriter, r *http.Request) {
	var item domain.Supplier
	if !decodeJSON(w, r, &item) {
		return
	}
	id, err := idgen.New("sup")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	item.ID = id
	if err := a.catalog.CreateSupplier(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "supplier.created", "supplier", item.ID, map[string]any{"code": item.Code, "name": item.Name})
	writeJSON(w, http.StatusCreated, item)
}
func (a *API) listProducts(w http.ResponseWriter, r *http.Request) {
	filter := repository.ProductFilter{Query: r.URL.Query().Get("q"), CategoryID: r.URL.Query().Get("categoryId"), SupplierID: r.URL.Query().Get("supplierId"), ActiveOnly: queryBool(r, "active", false), Limit: queryInt(r, "limit", 50), Offset: queryInt(r, "offset", 0)}
	items, err := a.catalog.ListProducts(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "limit": filter.Limit, "offset": filter.Offset})
}
func (a *API) createProduct(w http.ResponseWriter, r *http.Request) {
	var item domain.Product
	if !decodeJSON(w, r, &item) {
		return
	}
	id, err := idgen.New("prd")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	item.ID = id
	if item.Currency == "" {
		item.Currency = "INR"
	}
	if err := a.catalog.CreateProduct(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	created, err := a.catalog.GetProduct(r.Context(), item.ID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "product.created", "product", created.ID, map[string]any{"sku": created.SKU, "name": created.Name, "active": created.Active})
	writeJSON(w, http.StatusCreated, created)
}
func (a *API) getProduct(w http.ResponseWriter, r *http.Request) {
	item, err := a.catalog.GetProduct(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (a *API) getProductByBarcode(w http.ResponseWriter, r *http.Request) {
	item, err := a.catalog.GetProductByBarcode(r.Context(), r.PathValue("barcode"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (a *API) updateProduct(w http.ResponseWriter, r *http.Request) {
	var item domain.Product
	if !decodeJSON(w, r, &item) {
		return
	}
	item.ID = r.PathValue("id")
	if item.Currency == "" {
		item.Currency = "INR"
	}
	if err := a.catalog.UpdateProduct(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	updated, err := a.catalog.GetProduct(r.Context(), item.ID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "product.updated", "product", updated.ID, map[string]any{"sku": updated.SKU, "active": updated.Active})
	writeJSON(w, http.StatusOK, updated)
}
func queryInt(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
func queryBool(r *http.Request, key string, fallback bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

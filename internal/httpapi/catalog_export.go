package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const (
	defaultCatalogExportRows = 1000
	maxCatalogExportRows     = 5000
)

func normalizeCatalogExportBounds(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultCatalogExportRows
	}
	if limit > maxCatalogExportRows {
		limit = maxCatalogExportRows
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (a *API) exportProductsCSV(w http.ResponseWriter, r *http.Request) {
	limit, offset := normalizeCatalogExportBounds(
		queryInt(r, "limit", defaultCatalogExportRows),
		queryInt(r, "offset", 0),
	)
	filter := repository.ProductFilter{
		Query:      strings.TrimSpace(r.URL.Query().Get("q")),
		CategoryID: strings.TrimSpace(r.URL.Query().Get("categoryId")),
		SupplierID: strings.TrimSpace(r.URL.Query().Get("supplierId")),
		ActiveOnly: r.URL.Query().Get("activeOnly") != "false",
		Limit:      limit,
		Offset:     offset,
	}

	products, err := a.catalog.ListProducts(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"stockpilot-products.csv\"")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("id", "sku", "name", "description", "categoryId", "supplierId", "barcode", "unit", "unitCostMinor", "currency", "reorderPoint", "reorderQuantity", "trackLots", "trackExpiry", "active", "createdAt", "updatedAt"); err != nil {
		return
	}
	for _, product := range products {
		if err := writer.WriteRow(
			product.ID,
			product.SKU,
			product.Name,
			product.Description,
			product.CategoryID,
			product.SupplierID,
			product.Barcode,
			product.Unit,
			strconv.FormatInt(product.UnitCostMinor, 10),
			product.Currency,
			strconv.FormatInt(product.ReorderPoint, 10),
			strconv.FormatInt(product.ReorderQuantity, 10),
			strconv.FormatBool(product.TrackLots),
			strconv.FormatBool(product.TrackExpiry),
			strconv.FormatBool(product.Active),
			product.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			product.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

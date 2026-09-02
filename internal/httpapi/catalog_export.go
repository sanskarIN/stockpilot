package httpapi

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/repository"
)

const productExportLimit = 10000

func (a *API) exportProductsCSV(w http.ResponseWriter, r *http.Request) {
	filter := repository.ProductFilter{
		Query: r.URL.Query().Get("q"),
		CategoryID: r.URL.Query().Get("categoryId"),
		SupplierID: r.URL.Query().Get("supplierId"),
		ActiveOnly: queryBool(r, "active", false),
		Limit: productExportLimit,
	}
	items, err := a.catalog.ListProducts(r.Context(), filter)
	if err != nil { writeDomainError(w, err); return }
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="stockpilot-products.csv"`)
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "sku", "name", "description", "category_id", "supplier_id", "barcode", "unit", "unit_cost_minor", "currency", "reorder_point", "reorder_quantity", "track_lots", "track_expiry", "active"}); err != nil { return }
	for _, item := range items {
		if err := writer.Write([]string{
			item.ID, item.SKU, item.Name, item.Description, item.CategoryID, item.SupplierID, item.Barcode,
			item.Unit, strconv.FormatInt(item.UnitCostMinor, 10), item.Currency,
			strconv.FormatInt(item.ReorderPoint, 10), strconv.FormatInt(item.ReorderQuantity, 10),
			strconv.FormatBool(item.TrackLots), strconv.FormatBool(item.TrackExpiry), strconv.FormatBool(item.Active),
		}); err != nil { return }
	}
	writer.Flush()
	if err := writer.Error(); err != nil { fmt.Fprintln(w, "") }
}

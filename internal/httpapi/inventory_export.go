package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
)

const (
	defaultInventoryExportRows = 1000
	maxInventoryExportRows     = 5000
)

func (a *API) exportInventoryCSV(w http.ResponseWriter, r *http.Request) {
	limit, offset := normalizeInventoryExportBounds(queryInt(r, "limit", defaultInventoryExportRows), queryInt(r, "offset", 0))
	items, err := a.inventory.ListBalances(r.Context(), limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	setCSVDownloadHeaders(w, "stockpilot-inventory.csv")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("productId", "locationId", "lotId", "quantity", "updatedAt"); err != nil {
		return
	}
	for _, item := range items {
		if err := writer.WriteRow(item.ProductID, item.LocationID, item.LotID, strconv.FormatInt(item.Quantity, 10), formatExportTime(item.UpdatedAt)); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func (a *API) exportLowStockCSV(w http.ResponseWriter, r *http.Request) {
	limit := normalizeInventoryExportLimit(queryInt(r, "limit", defaultInventoryExportRows))
	items, err := a.inventory.ListLowStock(r.Context(), limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	setCSVDownloadHeaders(w, "stockpilot-low-stock.csv")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("productId", "locationId", "lotId", "quantity", "updatedAt"); err != nil {
		return
	}
	for _, item := range items {
		if err := writer.WriteRow(item.ProductID, item.LocationID, item.LotID, strconv.FormatInt(item.Quantity, 10), formatExportTime(item.UpdatedAt)); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func (a *API) exportReorderSuggestionsCSV(w http.ResponseWriter, r *http.Request) {
	limit := normalizeInventoryExportLimit(queryInt(r, "limit", defaultInventoryExportRows))
	items, err := a.inventory.ListReorderSuggestions(r.Context(), limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	setCSVDownloadHeaders(w, "stockpilot-reorder-suggestions.csv")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("productId", "sku", "name", "supplierId", "unit", "onHand", "reorderPoint", "reorderQuantity", "targetStock", "suggestedQuantity"); err != nil {
		return
	}
	for _, item := range items {
		if err := writer.WriteRow(
			item.ProductID,
			item.SKU,
			item.Name,
			item.SupplierID,
			item.Unit,
			strconv.FormatInt(item.OnHand, 10),
			strconv.FormatInt(item.ReorderPoint, 10),
			strconv.FormatInt(item.ReorderQuantity, 10),
			strconv.FormatInt(item.TargetStock, 10),
			strconv.FormatInt(item.SuggestedQuantity, 10),
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func normalizeInventoryExportBounds(limit, offset int) (int, int) {
	return normalizeInventoryExportLimit(limit), normalizeInventoryExportOffset(offset)
}

func normalizeInventoryExportLimit(limit int) int {
	if limit <= 0 {
		return defaultInventoryExportRows
	}
	if limit > maxInventoryExportRows {
		return maxInventoryExportRows
	}
	return limit
}

func normalizeInventoryExportOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func formatExportTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

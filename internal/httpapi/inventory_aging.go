package httpapi

import (
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
)

const (
	defaultInventoryAgingRows = 1000
	maxInventoryAgingRows     = 5000
)

func (a *API) inventoryAging(w http.ResponseWriter, r *http.Request) {
	limit := normalizeInventoryAgingLimit(queryInt(r, "limit", defaultInventoryAgingRows))
	report, err := a.inventory.GetInventoryAging(r.Context(), limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		setCSVDownloadHeaders(w, "stockpilot-inventory-aging.csv")
		writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
		if err := writer.WriteHeader("productId", "sku", "name", "locationId", "lotId", "quantity", "ageDays", "bucket", "asOf", "lastMovementAt"); err != nil {
			return
		}
		for _, item := range report.Items {
			if err := writer.WriteRow(item.ProductID, item.SKU, item.Name, item.LocationID, item.LotID, strconv.FormatInt(item.Quantity, 10), strconv.FormatInt(item.AgeDays, 10), string(item.Bucket), item.AsOf.Format("2006-01-02T15:04:05Z07:00"), item.LastMovementAt.Format("2006-01-02T15:04:05Z07:00")); err != nil {
				return
			}
		}
		_ = writer.Flush()
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func normalizeInventoryAgingLimit(limit int) int {
	if limit <= 0 {
		return defaultInventoryAgingRows
	}
	if limit > maxInventoryAgingRows {
		return maxInventoryAgingRows
	}
	return limit
}

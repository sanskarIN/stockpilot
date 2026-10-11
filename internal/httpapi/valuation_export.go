package httpapi

import (
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
)

const (
	defaultValuationExportRows = 1000
	maxValuationExportRows     = 5000
)

func (a *API) exportInventoryValuationCSV(w http.ResponseWriter, r *http.Request) {
	limit := normalizeValuationExportLimit(queryInt(r, "limit", defaultValuationExportRows))
	report, err := a.inventory.GetInventoryValuation(r.Context(), limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	setCSVDownloadHeaders(w, "stockpilot-inventory-valuation.csv")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("productId", "sku", "name", "unit", "onHand", "unitCostMinor", "currency", "valueMinor"); err != nil {
		return
	}
	for _, item := range report.Items {
		if err := writer.WriteRow(
			item.ProductID,
			item.SKU,
			item.Name,
			item.Unit,
			strconv.FormatInt(item.OnHand, 10),
			strconv.FormatInt(item.UnitCostMinor, 10),
			item.Currency,
			strconv.FormatInt(item.ValueMinor, 10),
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func normalizeValuationExportLimit(limit int) int {
	if limit <= 0 {
		return defaultValuationExportRows
	}
	if limit > maxValuationExportRows {
		return maxValuationExportRows
	}
	return limit
}

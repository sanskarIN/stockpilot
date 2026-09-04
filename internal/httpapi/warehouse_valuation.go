package httpapi

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

const (
	defaultWarehouseValuationLimit = 1000
	maxWarehouseValuationLimit     = 5000
)

func (a *API) warehouseValuation(w http.ResponseWriter, r *http.Request) {
	if a.reports == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "reporting is not available"})
		return
	}
	limit := queryInt(r, "limit", defaultWarehouseValuationLimit)
	if limit <= 0 {
		limit = defaultWarehouseValuationLimit
	}
	if limit > maxWarehouseValuationLimit {
		limit = maxWarehouseValuationLimit
	}
	report, err := a.reports.WarehouseValuation(r.Context(), limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		exportWarehouseValuationCSV(w, report)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func exportWarehouseValuationCSV(w http.ResponseWriter, report domain.WarehouseValuationReport) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=stockpilot-warehouse-valuation.csv")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	writer := csv.NewWriter(w)
	_ = writer.Write([]string{"warehouse_id", "warehouse_code", "warehouse_name", "location_id", "location_code", "location_name", "currency", "on_hand", "valuation_minor", "product_count"})
	for _, item := range report.Items {
		_ = writer.Write([]string{
			sanitizeCSVCell(item.WarehouseID), sanitizeCSVCell(item.WarehouseCode), sanitizeCSVCell(item.WarehouseName),
			sanitizeCSVCell(item.LocationID), sanitizeCSVCell(item.LocationCode), sanitizeCSVCell(item.LocationName),
			sanitizeCSVCell(item.Currency), strconv.FormatInt(item.OnHand, 10), strconv.FormatInt(item.ValuationMinor, 10), strconv.FormatInt(item.ProductCount, 10),
		})
	}
	writer.Flush()
}

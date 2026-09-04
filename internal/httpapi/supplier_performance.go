package httpapi

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

const (
	defaultSupplierPerformanceDays  = 30
	maxSupplierPerformanceDays      = 365
	defaultSupplierPerformanceLimit = 1000
	maxSupplierPerformanceLimit     = 5000
)

func (a *API) supplierPerformance(w http.ResponseWriter, r *http.Request) {
	if a.reports == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "reporting is not available"})
		return
	}
	days := normalizeSupplierPerformanceDays(queryInt(r, "days", defaultSupplierPerformanceDays))
	limit := normalizeSupplierPerformanceLimit(queryInt(r, "limit", defaultSupplierPerformanceLimit))
	report, err := a.reports.SupplierPerformance(r.Context(), days, limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		exportSupplierPerformanceCSV(w, report)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func normalizeSupplierPerformanceDays(value int) int {
	if value <= 0 {
		return defaultSupplierPerformanceDays
	}
	if value > maxSupplierPerformanceDays {
		return maxSupplierPerformanceDays
	}
	return value
}

func normalizeSupplierPerformanceLimit(value int) int {
	if value <= 0 {
		return defaultSupplierPerformanceLimit
	}
	if value > maxSupplierPerformanceLimit {
		return maxSupplierPerformanceLimit
	}
	return value
}

func exportSupplierPerformanceCSV(w http.ResponseWriter, report domain.SupplierPerformanceReport) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=stockpilot-supplier-performance.csv")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	writer := csv.NewWriter(w)
	_ = writer.Write([]string{"supplier_id", "supplier_code", "supplier_name", "order_count", "ordered_units", "received_units", "open_units", "ordered_value_minor", "received_value_minor", "average_lead_time_days", "completed_order_count", "on_time_order_count"})
	for _, item := range report.Items {
		_ = writer.Write([]string{
			sanitizeCSVCell(item.SupplierID), sanitizeCSVCell(item.SupplierCode), sanitizeCSVCell(item.SupplierName),
			strconv.FormatInt(item.OrderCount, 10), strconv.FormatInt(item.OrderedUnits, 10), strconv.FormatInt(item.ReceivedUnits, 10), strconv.FormatInt(item.OpenUnits, 10),
			strconv.FormatInt(item.OrderedValueMinor, 10), strconv.FormatInt(item.ReceivedValueMinor, 10), strconv.FormatFloat(item.AverageLeadTimeDays, 'f', 2, 64),
			strconv.FormatInt(item.CompletedOrderCount, 10), strconv.FormatInt(item.OnTimeOrderCount, 10),
		})
	}
	writer.Flush()
}

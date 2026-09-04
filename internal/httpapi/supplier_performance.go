package httpapi

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
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
	generatedAt := time.Now().UTC()
	period, err := reporting.NewPeriod(generatedAt.AddDate(0, 0, -(days - 1)), generatedAt)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	bounds, err := reporting.NewBounds(limit, 0, defaultSupplierPerformanceLimit, maxSupplierPerformanceLimit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	complete := len(report.Items) < limit
	if r.URL.Query().Get("format") == "csv" {
		exportSupplierPerformanceCSV(w, report, period, bounds, complete, generatedAt)
		return
	}
	writeReportMetadata(w, reportRequest{Period: period, Bounds: bounds}, generatedAt, complete)
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

func exportSupplierPerformanceCSV(w http.ResponseWriter, report domain.SupplierPerformanceReport, period reporting.Period, bounds reporting.Bounds, complete bool, generatedAt time.Time) {
	writeReportExportHeaders(w, "stockpilot-supplier-performance.csv", period, bounds, complete, generatedAt)
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

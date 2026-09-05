package httpapi

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
	"github.com/sanskarIN/stockpilot/internal/repository"
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
	offset, err := parseReportOffsetParameter(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	generatedAt := time.Now().UTC()
	period, err := reporting.NewPeriod(generatedAt, generatedAt)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	bounds, err := reporting.NewBounds(limit, offset, defaultWarehouseValuationLimit, maxWarehouseValuationLimit)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	var report domain.WarehouseValuationReport
	if bounded, ok := a.reports.(repository.BoundedReports); ok {
		report, err = bounded.WarehouseValuationQuery(r.Context(), makeBoundedReportQuery(period, bounds))
	} else {
		if offset != 0 {
			writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "report offset requires bounded reporting capability"})
			return
		}
		report, err = a.reports.WarehouseValuation(r.Context(), limit)
	}
	if err != nil {
		writeDomainError(w, err)
		return
	}

	complete := len(report.Items) < limit
	if r.URL.Query().Get("format") == "csv" {
		exportWarehouseValuationCSV(w, report, period, bounds, complete, generatedAt)
		return
	}
	writeReportMetadata(w, reportRequest{Period: period, Bounds: bounds}, generatedAt, complete)
	writeJSON(w, http.StatusOK, report)
}

func exportWarehouseValuationCSV(w http.ResponseWriter, report domain.WarehouseValuationReport, period reporting.Period, bounds reporting.Bounds, complete bool, generatedAt time.Time) {
	writeReportExportHeaders(w, "stockpilot-warehouse-valuation.csv", period, bounds, complete, generatedAt)
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

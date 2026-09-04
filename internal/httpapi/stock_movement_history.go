package httpapi

import (
	"net/http"
	"strconv"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const (
	defaultMovementHistoryWindowDays = 30
	maxMovementHistoryWindowDays     = 365
	defaultMovementHistoryRows       = 1000
	maxMovementHistoryRows           = 5000
)

func (a *API) stockMovementHistory(w http.ResponseWriter, r *http.Request) {
	reader, ok := a.inventory.(repository.StockMovementHistory)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "stock movement history is not available"})
		return
	}
	windowDays := normalizeMovementHistoryWindow(queryInt(r, "days", defaultMovementHistoryWindowDays))
	limit := normalizeMovementHistoryLimit(queryInt(r, "limit", defaultMovementHistoryRows))
	report, err := reader.GetStockMovementHistory(r.Context(), windowDays, limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		setCSVDownloadHeaders(w, "stockpilot-stock-movement-history.csv")
		writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
		if err := writer.WriteHeader("productId", "sku", "name", "locationId", "lotId", "movementCount", "inboundUnits", "outboundUnits", "netUnits", "averageDailyOutbound", "lastMovementAt", "asOf", "windowDays"); err != nil {
			return
		}
		for _, item := range report.Items {
			if err := writer.WriteRow(item.ProductID, item.SKU, item.Name, item.LocationID, item.LotID, strconv.FormatInt(item.MovementCount, 10), strconv.FormatInt(item.InboundUnits, 10), strconv.FormatInt(item.OutboundUnits, 10), strconv.FormatInt(item.NetUnits, 10), strconv.FormatFloat(item.AverageDailyOut, 'f', 2, 64), item.LastMovementAt.Format("2006-01-02T15:04:05Z07:00"), report.AsOf.Format("2006-01-02T15:04:05Z07:00"), strconv.Itoa(report.WindowDays)); err != nil {
				return
			}
		}
		_ = writer.Flush()
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func normalizeMovementHistoryWindow(days int) int {
	if days <= 0 {
		return defaultMovementHistoryWindowDays
	}
	if days > maxMovementHistoryWindowDays {
		return maxMovementHistoryWindowDays
	}
	return days
}

func normalizeMovementHistoryLimit(limit int) int {
	if limit <= 0 {
		return defaultMovementHistoryRows
	}
	if limit > maxMovementHistoryRows {
		return maxMovementHistoryRows
	}
	return limit
}

var _ domain.StockMovementHistoryReport

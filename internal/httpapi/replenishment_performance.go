package httpapi

import (
	"net/http"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (a *API) replenishmentPerformance(w http.ResponseWriter, r *http.Request) {
	reader, ok := a.reports.(repository.ReplenishmentReports)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "replenishment performance is not available"})
		return
	}
	request, err := parseReportRequest(r, time.Now().UTC())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	query := makeBoundedReportQuery(request.Period, request.Bounds)
	report, err := reader.ReplenishmentPerformance(r.Context(), query)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeReportMetadata(w, request, report.AsOf, len(report.Items) < request.Bounds.Limit)
	writeJSON(w, http.StatusOK, report)
}

var _ = domain.ReplenishmentPerformanceReport{}

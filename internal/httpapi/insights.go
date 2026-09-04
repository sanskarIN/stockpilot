package httpapi

import (
	"net/http"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (a *API) reportOverview(w http.ResponseWriter, r *http.Request) {
	if a.reports == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "reporting is not available"})
		return
	}
	request, err := parseReportRequest(r, time.Now())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	generatedAt := time.Now().UTC()
	inventory, err := a.reports.InventorySummary(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	purchasing, err := a.reports.PurchasingSummary(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeReportMetadata(w, request, generatedAt, true)
	writeJSON(w, http.StatusOK, map[string]any{"inventory": inventory, "purchasing": purchasing})
}

func (a *API) reportInventory(w http.ResponseWriter, r *http.Request) {
	if a.reports == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "reporting is not available"})
		return
	}
	request, err := parseReportRequest(r, time.Now())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	generatedAt := time.Now().UTC()
	report, err := a.reports.InventorySummary(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeReportMetadata(w, request, generatedAt, true)
	writeJSON(w, http.StatusOK, report)
}

func (a *API) reportPurchasing(w http.ResponseWriter, r *http.Request) {
	if a.reports == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "reporting is not available"})
		return
	}
	request, err := parseReportRequest(r, time.Now())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	generatedAt := time.Now().UTC()
	report, err := a.reports.PurchasingSummary(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeReportMetadata(w, request, generatedAt, true)
	writeJSON(w, http.StatusOK, report)
}

func (a *API) listAuditEvents(w http.ResponseWriter, r *http.Request) {
	if a.audit == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "audit is not available"})
		return
	}
	filter := domain.AuditFilter{
		ActorID:    r.URL.Query().Get("actorId"),
		Action:     r.URL.Query().Get("action"),
		EntityType: r.URL.Query().Query().Get("entityType"),
		EntityID:   r.URL.Query().Get("entityId"),
		Limit:      queryInt(r, "limit", 100),
		Offset:     queryInt(r, "offset", 0),
	}
	items, err := a.audit.ListAuditEvents(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "limit": filter.Limit, "offset": filter.Offset})
}

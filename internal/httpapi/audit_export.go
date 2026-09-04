package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/domain"
)

const (
	defaultAuditExportRows = 500
	maxAuditExportRows     = 5000
)

func (a *API) exportAuditCSV(w http.ResponseWriter, r *http.Request) {
	limit, offset := normalizeAuditExportBounds(queryInt(r, "limit", defaultAuditExportRows), queryInt(r, "offset", 0))
	filter := domain.AuditFilter{
		ActorID: strings.TrimSpace(r.URL.Query().Get("actorId")),
		Action: strings.TrimSpace(r.URL.Query().Get("action")),
		EntityType: strings.TrimSpace(r.URL.Query().Get("entityType")),
		EntityID: strings.TrimSpace(r.URL.Query().Get("entityId")),
		Limit: limit,
		Offset: offset,
	}
	items, err := a.audit.ListAuditEvents(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="stockpilot-audit-log.csv"`)
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("id", "occurredAt", "actorId", "action", "entityType", "entityId", "requestId", "metadata"); err != nil {
		return
	}
	for _, item := range items {
		metadata := "{}"
		if item.Metadata != nil {
			encoded, encodeErr := json.Marshal(item.Metadata)
			if encodeErr != nil {
				return
			}
			metadata = string(encoded)
		}
		if err := writer.WriteRow(
			strconv.FormatInt(item.ID, 10),
			formatExportTime(item.OccurredAt),
			item.ActorID,
			item.Action,
			item.EntityType,
			item.EntityID,
			item.RequestID,
			metadata,
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func normalizeAuditExportBounds(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultAuditExportRows
	}
	if limit > maxAuditExportRows {
		limit = maxAuditExportRows
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

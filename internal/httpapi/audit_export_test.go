package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestNormalizeAuditExportBounds(t *testing.T) {
	tests := []struct {
		name                   string
		limit, offset           int
		wantLimit, wantOffset   int
	}{
		{"defaults", 0, 0, 500, 0},
		{"negative", -1, -5, 500, 0},
		{"clamped", 9000, 4, 5000, 4},
		{"valid", 125, 8, 125, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := normalizeAuditExportBounds(tt.limit, tt.offset)
			if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
				t.Fatalf("got (%d,%d), want (%d,%d)", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestAuditExportCSV(t *testing.T) {
	occurred := time.Date(2026, 9, 4, 4, 0, 0, 0, time.FixedZone("IST", 19800))
	store := &fakeStore{auditEvents: []domain.AuditEvent{{
		ID: 42, OccurredAt: occurred, ActorID: "usr_1", Action: "product.updated", EntityType: "product", EntityID: "prd_1", RequestID: "req_1", Metadata: map[string]any{"note": "=safe"},
	}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil }, WithInsights(nil, store))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, authenticatedRequest(http.MethodGet, "/api/v1/audit/export.csv?entityType=product&limit=20&offset=2", ""))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "text/csv; charset=utf-8" {
		t.Fatalf("content type=%q", recorder.Header().Get("Content-Type"))
	}
	if recorder.Header().Get("Content-Disposition") != `attachment; filename="stockpilot-audit-log.csv"` {
		t.Fatalf("content disposition=%q", recorder.Header().Get("Content-Disposition"))
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "id,occurredAt,actorId,action,entityType,entityId,requestId,metadata") {
		t.Fatalf("missing header: %s", body)
	}
	if !strings.Contains(body, "42,2026-09-03T22:30:00Z,usr_1,product.updated") {
		t.Fatalf("missing audit row: %s", body)
	}
	if !strings.Contains(body, `"note":"'=safe"`) {
		t.Fatalf("formula-safe metadata missing: %s", body)
	}
}

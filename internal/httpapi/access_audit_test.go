package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type accessAuditRecorder struct{ events []domain.AuditEvent }

func (a *accessAuditRecorder) AppendAuditEvent(_ context.Context, event domain.AuditEvent) error {
	a.events = append(a.events, event)
	return nil
}

func (a *accessAuditRecorder) ListAuditEvents(context.Context, domain.AuditFilter) ([]domain.AuditEvent, error) {
	return append([]domain.AuditEvent(nil), a.events...), nil
}

func TestAccessAuditRecordsSafeAuthenticationEvent(t *testing.T) {
	recorder := &accessAuditRecorder{}
	handler := &accessHandler{audit: recorder}
	ctx := context.WithValue(context.Background(), requestIDContextKey{}, "req_test")

	handler.recordAudit(ctx, "usr_123", "auth.login.success", "user", "usr_123", map[string]any{"outcome": "accepted"})

	if len(recorder.events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(recorder.events))
	}
	event := recorder.events[0]
	if event.Action != "auth.login.success" || event.ActorID != "usr_123" || event.EntityID != "usr_123" {
		t.Fatalf("unexpected audit event: %#v", event)
	}
	if event.RequestID != "req_test" {
		t.Fatalf("expected request id req_test, got %q", event.RequestID)
	}
	if len(event.Metadata) != 1 || event.Metadata["outcome"] != "accepted" {
		t.Fatalf("unexpected metadata: %#v", event.Metadata)
	}
}

func TestAccessAuditNeverStoresAuthenticationSecrets(t *testing.T) {
	recorder := &accessAuditRecorder{}
	handler := &accessHandler{audit: recorder}
	ctx := context.WithValue(context.Background(), requestIDContextKey{}, "req_secret")

	handler.recordAudit(ctx, "", "auth.login.failure", "session", "", map[string]any{"outcome": "rejected"})

	event := recorder.events[0]
	if _, ok := event.Metadata["password"]; ok {
		t.Fatal("audit metadata must not include passwords")
	}
	if _, ok := event.Metadata["token"]; ok {
		t.Fatal("audit metadata must not include session tokens")
	}
	if event.RequestID != "req_secret" {
		t.Fatalf("expected request id req_secret, got %q", event.RequestID)
	}
}

func TestIsCSVExportRequest(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   bool
	}{
		{name: "get export", method: http.MethodGet, path: "/api/v1/products/export.csv", want: true},
		{name: "uppercase suffix", method: http.MethodGet, path: "/api/v1/inventory/EXPORT.CSV", want: true},
		{name: "head is not export audit", method: http.MethodHead, path: "/api/v1/products/export.csv", want: false},
		{name: "post is not export audit", method: http.MethodPost, path: "/api/v1/products/export.csv", want: false},
		{name: "json route", method: http.MethodGet, path: "/api/v1/products", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, nil)
			if got := isCSVExportRequest(r); got != tt.want {
				t.Fatalf("isCSVExportRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessAuditRecordsCSVExportWithActorAndRequestID(t *testing.T) {
	recorder := &accessAuditRecorder{}
	handler := &accessHandler{audit: recorder}
	ctx := context.WithValue(context.Background(), requestIDContextKey{}, "req_export")
	r := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/export.csv", nil).WithContext(ctx)

	if !isCSVExportRequest(r) {
		t.Fatal("expected CSV export request to be recognized")
	}
	handler.recordAudit(r.Context(), "usr_export", "export.csv.requested", "export", r.URL.Path, map[string]any{"method": r.Method})

	if len(recorder.events) != 1 {
		t.Fatalf("expected one export audit event, got %d", len(recorder.events))
	}
	event := recorder.events[0]
	if event.Action != "export.csv.requested" {
		t.Fatalf("expected export action, got %q", event.Action)
	}
	if event.ActorID != "usr_export" || event.EntityType != "export" || event.EntityID != r.URL.Path {
		t.Fatalf("unexpected export audit identity: %#v", event)
	}
	if event.RequestID != "req_export" {
		t.Fatalf("expected request id req_export, got %q", event.RequestID)
	}
	if event.Metadata["method"] != http.MethodGet {
		t.Fatalf("expected GET metadata, got %#v", event.Metadata)
	}
}

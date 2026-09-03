package httpapi

import (
	"context"
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

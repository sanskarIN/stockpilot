package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestReceiveOrderLineNewLotUsesAtomicRepositoryPath(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(http.MethodPost, "/api/v1/orders/po_1/lines/pol_1/receive", `{"quantity":2,"locationId":"loc_1","newLot":{"lotNumber":"LOT-42","manufacturedAt":"2026-01-01T00:00:00Z","expiresAt":"2027-01-01T00:00:00Z"}}`)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if store.receiptActor != "usr_session" {
		t.Fatalf("receipt actor = %q, want authenticated user", store.receiptActor)
	}
}

func TestReceiveOrderLineNewLotRejectsInvalidDateOrdering(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(http.MethodPost, "/api/v1/orders/po_1/lines/pol_1/receive", `{"quantity":2,"locationId":"loc_1","newLot":{"lotNumber":"LOT-42","manufacturedAt":"2027-01-01T00:00:00Z","expiresAt":"2026-01-01T00:00:00Z"}}`)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
}

package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (f *fakeStore) UpdateOrderStatus(_ context.Context, _ string, status domain.PurchaseOrderStatus, _ string) error {
	if err := domain.ValidatePurchaseOrderTransition(domain.PurchaseOrderDraft, status); err != nil {
		return err
	}
	return nil
}

func TestUpdateOrderStatusEndpointAllowsDraftToOrdered(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recording := httptest.NewRecorder()
	request := authenticatedRequest(http.MethodPatch, "/api/v1/orders/po_1/status", `{"status":"ordered"}`)
	handler.ServeHTTP(recording, request)
	if recording.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recording.Code, http.StatusOK, recording.Body.String())
	}
}

func TestUpdateOrderStatusEndpointRejectsManualReceivedState(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recording := httptest.NewRecorder()
	request := authenticatedRequest(http.MethodPatch, "/api/v1/orders/po_1/status", `{"status":"received"}`)
	handler.ServeHTTP(recording, request)
	if recording.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body=%s", recording.Code, http.StatusConflict, recording.Body.String())
	}
}

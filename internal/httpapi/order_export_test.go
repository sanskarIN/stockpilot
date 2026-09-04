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

func (f *fakeStore) ListOrderExportRows(_ context.Context, _ domain.PurchaseOrderStatus, _, _ int) ([]domain.PurchaseOrderExportRow, error) {
	created := time.Date(2026, 9, 4, 4, 0, 0, 0, time.FixedZone("IST", 5*60*60+30*60))
	expected := time.Date(2026, 9, 12, 0, 0, 0, 0, time.UTC)
	return []domain.PurchaseOrderExportRow{{
		OrderID: "po_1", OrderNumber: "=PO-100", SupplierID: "sup_1", WarehouseID: "wh_1",
		Status: domain.PurchaseOrderPartiallyReceived, Currency: "INR", ExpectedAt: &expected,
		Notes: "deliver before noon", CreatedBy: "usr_1", CreatedAt: created, UpdatedAt: created,
		LineID: "pol_1", ProductID: "prd_1", Quantity: 10, Received: 4, UnitCostMinor: 2500,
	}}, nil
}

func TestNormalizeOrderExportBounds(t *testing.T) {
	tests := []struct {
		name                 string
		limit, offset         int
		wantLimit, wantOffset int
	}{
		{name: "defaults", limit: 0, offset: 0, wantLimit: 500, wantOffset: 0},
		{name: "negative limit", limit: -2, offset: -4, wantLimit: 500, wantOffset: 0},
		{name: "clamped", limit: 9000, offset: 7, wantLimit: 5000, wantOffset: 7},
		{name: "preserves valid", limit: 125, offset: 9, wantLimit: 125, wantOffset: 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := normalizeOrderExportBounds(tt.limit, tt.offset)
			if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset {
				t.Fatalf("got (%d,%d), want (%d,%d)", gotLimit, gotOffset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestPurchaseOrderExportCSV(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(http.MethodGet, "/api/v1/orders/export.csv?status=partially_received&limit=100&offset=2", "")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content type=%q", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="stockpilot-purchase-orders.csv"` {
		t.Fatalf("content disposition=%q", got)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "orderId,orderNumber,supplierId,warehouseId,status,currency") {
		t.Fatalf("missing header: %s", body)
	}
	if !strings.Contains(body, "'=PO-100") {
		t.Fatalf("formula-safe order number missing: %s", body)
	}
	if !strings.Contains(body, "6,2500,25000") {
		t.Fatalf("remaining/line total missing: %s", body)
	}
	if !strings.Contains(body, "2026-09-12T00:00:00Z") || !strings.Contains(body, "2026-09-03T22:30:00Z") {
		t.Fatalf("UTC timestamps missing: %s", body)
	}
}

func TestPurchaseOrderExportRejectsInvalidStatus(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/orders/export.csv?status=unknown", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", recorder.Code)
	}
}

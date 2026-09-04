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

func TestInventoryExportBounds(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{name: "defaults", limit: 0, offset: 0, wantLimit: 1000, wantOffset: 0},
		{name: "negative", limit: -4, offset: -9, wantLimit: 1000, wantOffset: 0},
		{name: "clamps maximum", limit: 9000, offset: 7, wantLimit: 5000, wantOffset: 7},
		{name: "preserves valid", limit: 25, offset: 11, wantLimit: 25, wantOffset: 11},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limit, offset := normalizeInventoryExportBounds(test.limit, test.offset)
			if limit != test.wantLimit || offset != test.wantOffset {
				t.Fatalf("got limit=%d offset=%d", limit, offset)
			}
		})
	}
}

func TestInventoryExportCSV(t *testing.T) {
	now := time.Date(2026, 9, 3, 12, 30, 0, 0, time.FixedZone("IST", 19800))
	store := &fakeStore{balances: []domain.StockBalance{{ProductID: "=PRODUCT", LocationID: "loc_1", LotID: "lot_1", Quantity: 42, UpdatedAt: now}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/export.csv?limit=10&offset=2", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content type=%q", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="stockpilot-inventory.csv"` {
		t.Fatalf("content disposition=%q", got)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "productId,locationId,lotId,quantity,updatedAt") {
		t.Fatalf("header missing: %q", body)
	}
	if !strings.Contains(body, `'=PRODUCT,loc_1,lot_1,42,2026-09-03T07:00:00Z`) {
		t.Fatalf("formula-safe/time formatting missing: %q", body)
	}
}

func TestLowStockExportCSV(t *testing.T) {
	store := &fakeStore{lowStock: []domain.StockBalance{{ProductID: "prd_1", LocationID: "loc_1", Quantity: 2}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/low-stock/export.csv", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "productId,locationId,lotId,quantity,updatedAt") || !strings.Contains(recorder.Body.String(), "prd_1,loc_1,,2,") {
		t.Fatalf("body=%q", recorder.Body.String())
	}
}

func TestReorderSuggestionsExportCSV(t *testing.T) {
	store := &fakeStore{reorderSuggestions: []domain.ReorderSuggestion{{ProductID: "prd_1", SKU: "SKU-1", Name: "Widget", SupplierID: "sup_1", Unit: "piece", OnHand: 3, ReorderPoint: 5, ReorderQuantity: 10, TargetStock: 15, SuggestedQuantity: 12}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/reorder-suggestions/export.csv?limit=25", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d", recorder.Code)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "productId,sku,name,supplierId,unit,onHand,reorderPoint,reorderQuantity,targetStock,suggestedQuantity") {
		t.Fatalf("header missing: %q", body)
	}
	if !strings.Contains(body, "prd_1,SKU-1,Widget,sup_1,piece,3,5,10,15,12") {
		t.Fatalf("row missing: %q", body)
	}
}

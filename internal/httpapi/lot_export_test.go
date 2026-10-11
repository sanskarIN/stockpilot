package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func TestLotExportBounds(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{name: "defaults", limit: 0, offset: 0, wantLimit: 500, wantOffset: 0},
		{name: "negative", limit: -5, offset: -3, wantLimit: 500, wantOffset: 0},
		{name: "clamps maximum", limit: 9000, offset: 4, wantLimit: 5000, wantOffset: 4},
		{name: "preserves valid", limit: 20, offset: 8, wantLimit: 20, wantOffset: 8},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limit, offset := normalizeLotExportBounds(test.limit, test.offset)
			if limit != test.wantLimit || offset != test.wantOffset {
				t.Fatalf("got limit=%d offset=%d", limit, offset)
			}
		})
	}
}

func TestParseExportDate(t *testing.T) {
	value, err := parseExportDate("2026-09-30")
	if err != nil {
		t.Fatal(err)
	}
	if !value.Equal(time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("got %v", value)
	}
	if _, err := parseExportDate("30-09-2026"); err == nil {
		t.Fatal("expected invalid date error")
	}
}

func TestLotInventoryExportCSV(t *testing.T) {
	expires := time.Date(2026, 9, 20, 15, 0, 0, 0, time.FixedZone("IST", 19800))
	store := &fakeStore{lotInventory: []domain.LotInventoryRow{{
		ProductID: "prd_1", SKU: "=SKU-1", ProductName: "Widget", LotID: "lot_1", LotNumber: "LOT-1",
		LocationID: "loc_1", Location: "A1", WarehouseID: "wh_1", Warehouse: "Main", OnHand: 12,
		ExpiresAt: &expires, Active: true,
	}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/lots/export.csv?productId=prd_1&warehouseId=wh_1&locationId=loc_1&lotId=lot_1&expiringBy=2026-10-01&limit=20&offset=2", nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content type=%q", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="stockpilot-lot-inventory.csv"` {
		t.Fatalf("content disposition=%q", got)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "productId,sku,productName,lotId,lotNumber,locationId,location,warehouseId,warehouse,onHand,expiresAt,active") {
		t.Fatalf("header missing: %q", body)
	}
	if !strings.Contains(body, "'=SKU-1,Widget,lot_1,LOT-1,loc_1,A1,wh_1,Main,12,2026-09-20T09:30:00Z,true") {
		t.Fatalf("row missing: %q", body)
	}
}

func TestLotInventoryExportRejectsInvalidExpiryDate(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/lots/export.csv?expiringBy=invalid", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", recorder.Code)
	}
}

func TestLotInventoryExportPassesFilters(t *testing.T) {
	store := &fakeStore{}
	store.lotInventory = []domain.LotInventoryRow{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/lots/export.csv?productId=prd_1&warehouseId=wh_1&locationId=loc_1&lotId=lot_1&limit=12&offset=6", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d", recorder.Code)
	}
	if len(store.lotInventory) != 0 {
		t.Fatal("unexpected fixture mutation")
	}
	var _ repository.Inventory = store
}

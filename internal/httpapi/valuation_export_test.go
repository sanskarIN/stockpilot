package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestValuationExportBounds(t *testing.T) {
	cases := []struct{ input, want int }{{0, 1000}, {-1, 1000}, {42, 42}, {9000, 5000}}
	for _, tc := range cases {
		if got := normalizeValuationExportLimit(tc.input); got != tc.want { t.Fatalf("limit %d: got %d want %d", tc.input, got, tc.want) }
	}
}

func TestValuationExportCSV(t *testing.T) {
	store := &fakeStore{valuation: domain.InventoryValuationReport{Items: []domain.InventoryValuationItem{{ProductID: "prd_1", SKU: "=SKU-1", Name: "Widget", Unit: "piece", OnHand: 12, UnitCostMinor: 1250, Currency: "INR", ValueMinor: 15000}}, Totals: []domain.InventoryValuationTotal{{Currency: "INR", ValueMinor: 15000}}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/inventory-valuation/export.csv?limit=10", nil))
	if recorder.Code != http.StatusOK { t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String()) }
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="stockpilot-inventory-valuation.csv"` { t.Fatalf("content disposition=%q", got) }
	if got := recorder.Header().Get("Cache-Control"); got != "no-store" { t.Fatalf("cache control=%q", got) }
	body := recorder.Body.String()
	if !strings.Contains(body, "productId,sku,name,unit,onHand,unitCostMinor,currency,valueMinor") { t.Fatalf("header missing: %q", body) }
	if !strings.Contains(body, "prd_1,'=SKU-1,Widget,piece,12,1250,INR,15000") { t.Fatalf("row missing: %q", body) }
}

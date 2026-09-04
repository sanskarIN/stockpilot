package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type warehouseValuationReportsStub struct{}

func (warehouseValuationReportsStub) InventorySummary(context.Context) (domain.InventorySummary, error) {
	return domain.InventorySummary{}, nil
}
func (warehouseValuationReportsStub) PurchasingSummary(context.Context) (domain.PurchasingSummary, error) {
	return domain.PurchasingSummary{}, nil
}
func (warehouseValuationReportsStub) SupplierPerformance(context.Context, int, int) (domain.SupplierPerformanceReport, error) {
	return domain.SupplierPerformanceReport{}, nil
}
func (warehouseValuationReportsStub) WarehouseValuation(context.Context, int) (domain.WarehouseValuationReport, error) {
	return domain.WarehouseValuationReport{Items: []domain.WarehouseValuationItem{{
		WarehouseID: "wh_1", WarehouseCode: "MAIN", WarehouseName: "Main",
		LocationID: "loc_1", LocationCode: "A1", LocationName: "A1",
		Currency: "INR", OnHand: 4, ValuationMinor: 10000, ProductCount: 1,
	}}}, nil
}

type warehouseValuationAuditStub struct{}
func (warehouseValuationAuditStub) AppendAuditEvent(context.Context, domain.AuditEvent) error { return nil }
func (warehouseValuationAuditStub) ListAuditEvents(context.Context, domain.AuditFilter) ([]domain.AuditEvent, error) { return nil, nil }

func TestWarehouseValuationDefaultsAndJSON(t *testing.T) {
	handler := NewCore(nil, nil, nil, func(context.Context) error { return nil }, WithInsights(warehouseValuationReportsStub{}, warehouseValuationAuditStub{}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/warehouse-valuation", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"warehouseCode":"MAIN"`) {
		t.Fatalf("body=%q", recorder.Body.String())
	}
}

func TestWarehouseValuationCSVHeadersAndSafety(t *testing.T) {
	report := domain.WarehouseValuationReport{Items: []domain.WarehouseValuationItem{{
		WarehouseID: "wh_1", WarehouseCode: "=BAD", WarehouseName: "Main",
		LocationID: "loc_1", LocationCode: "A1", LocationName: "A1",
		Currency: "INR", OnHand: 4, ValuationMinor: 10000, ProductCount: 1,
	}}}
	recorder := httptest.NewRecorder()
	exportWarehouseValuationCSV(recorder, report)
	if recorder.Header().Get("Cache-Control") != "no-store, no-cache, must-revalidate" {
		t.Fatalf("cache-control=%q", recorder.Header().Get("Cache-Control"))
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "warehouse_id,warehouse_code,warehouse_name") {
		t.Fatalf("missing header: %q", body)
	}
	if strings.Contains(body, "\n=BAD,") {
		t.Fatalf("formula-like CSV cell was not sanitized: %q", body)
	}
}

func TestWarehouseValuationLimitIsBounded(t *testing.T) {
	handler := NewCore(nil, nil, nil, func(context.Context) error { return nil }, WithInsights(warehouseValuationReportsStub{}, warehouseValuationAuditStub{}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/warehouse-valuation?limit=999999", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d", recorder.Code)
	}
}

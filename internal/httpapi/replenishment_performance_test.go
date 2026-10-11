package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
)

type replenishmentReportStub struct {
	query reporting.Query
}

func (s *replenishmentReportStub) InventorySummary(context.Context) (domain.InventorySummary, error) {
	return domain.InventorySummary{}, nil
}
func (s *replenishmentReportStub) PurchasingSummary(context.Context) (domain.PurchasingSummary, error) {
	return domain.PurchasingSummary{}, nil
}
func (s *replenishmentReportStub) SupplierPerformance(context.Context, int, int) (domain.SupplierPerformanceReport, error) {
	return domain.SupplierPerformanceReport{}, nil
}
func (s *replenishmentReportStub) WarehouseValuation(context.Context, int) (domain.WarehouseValuationReport, error) {
	return domain.WarehouseValuationReport{}, nil
}
func (s *replenishmentReportStub) ReplenishmentPerformance(_ context.Context, query reporting.Query) (domain.ReplenishmentPerformanceReport, error) {
	s.query = query
	return domain.ReplenishmentPerformanceReport{
		AsOf: time.Date(2026, 9, 5, 13, 0, 0, 0, time.UTC), WindowDays: 14,
		Items: []domain.ReplenishmentPerformanceItem{{SupplierID: "sup-1", SupplierName: "Acme", OrderCount: 2, OrderedUnits: 20, ReceivedUnits: 18, OutstandingUnits: 2, FillRate: 0.9, OnTimeOrderCount: 1, LateOrderCount: 1, AverageLeadDays: 3.5}},
	}, nil
}

func TestReplenishmentPerformanceEndpointUsesBoundedQuery(t *testing.T) {
	stub := &replenishmentReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/replenishment-performance?from=2026-08-01T00:00:00Z&to=2026-08-14T00:00:00Z&limit=25&offset=5", nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.query.Limit != 25 || stub.query.Offset != 5 {
		t.Fatalf("bounds=%+v", stub.query)
	}
	if stub.query.From == nil || stub.query.To == nil {
		t.Fatal("bounded period was not propagated")
	}
	if got := stub.query.To.Sub(*stub.query.From).Hours()/24 + 1; got != 14 {
		t.Fatalf("period days=%v", got)
	}
	if recorder.Header().Get("X-Report-Limit") != "25" || recorder.Header().Get("X-Report-Offset") != "5" {
		t.Fatalf("pagination metadata missing: limit=%q offset=%q", recorder.Header().Get("X-Report-Limit"), recorder.Header().Get("X-Report-Offset"))
	}
	var report domain.ReplenishmentPerformanceReport
	if err := json.Unmarshal(recorder.Body.Bytes(), &report); err != nil {
		t.Fatalf("decode report: %v", err)
	}
	if len(report.Items) != 1 || report.Items[0].FillRate != 0.9 {
		t.Fatalf("unexpected report=%+v", report)
	}
}

func TestReplenishmentPerformanceEndpointValidatesBounds(t *testing.T) {
	stub := &replenishmentReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/replenishment-performance?limit=0", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.query.Limit != 0 || stub.query.Offset != 0 {
		t.Fatal("invalid request reached repository capability")
	}
}

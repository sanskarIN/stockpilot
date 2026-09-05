package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
)

type reportCapabilityStub struct {
	supplierCalls        int
	warehouseCalls       int
	legacySupplierCalls  int
	legacyWarehouseCalls int
	supplierQuery        reporting.Query
	warehouseQuery       reporting.Query
}

func (s *reportCapabilityStub) InventorySummary(context.Context) (domain.InventorySummary, error) {
	return domain.InventorySummary{}, nil
}

func (s *reportCapabilityStub) PurchasingSummary(context.Context) (domain.PurchasingSummary, error) {
	return domain.PurchasingSummary{}, nil
}

func (s *reportCapabilityStub) SupplierPerformance(context.Context, int, int) (domain.SupplierPerformanceReport, error) {
	s.legacySupplierCalls++
	return domain.SupplierPerformanceReport{}, nil
}

func (s *reportCapabilityStub) WarehouseValuation(context.Context, int) (domain.WarehouseValuationReport, error) {
	s.legacyWarehouseCalls++
	return domain.WarehouseValuationReport{}, nil
}

func (s *reportCapabilityStub) SupplierPerformanceQuery(_ context.Context, query reporting.Query) (domain.SupplierPerformanceReport, error) {
	s.supplierCalls++
	s.supplierQuery = query
	return domain.SupplierPerformanceReport{}, nil
}

func (s *reportCapabilityStub) WarehouseValuationQuery(_ context.Context, query reporting.Query) (domain.WarehouseValuationReport, error) {
	s.warehouseCalls++
	s.warehouseQuery = query
	return domain.WarehouseValuationReport{}, nil
}

type countedReportStub struct {
	reportCapabilityStub
	supplierCountCalls  int
	warehouseCountCalls int
}

func (s *countedReportStub) SupplierPerformanceCount(_ context.Context, query reporting.Query) (int64, error) {
	s.supplierCountCalls++
	s.supplierQuery = query
	return 37, nil
}

func (s *countedReportStub) WarehouseValuationCount(_ context.Context, query reporting.Query) (int64, error) {
	s.warehouseCountCalls++
	s.warehouseQuery = query
	return 23, nil
}

type legacyReportStub struct {
	legacySupplierCalls  int
	legacyWarehouseCalls int
}

func (s *legacyReportStub) InventorySummary(context.Context) (domain.InventorySummary, error) {
	return domain.InventorySummary{}, nil
}

func (s *legacyReportStub) PurchasingSummary(context.Context) (domain.PurchasingSummary, error) {
	return domain.PurchasingSummary{}, nil
}

func (s *legacyReportStub) SupplierPerformance(context.Context, int, int) (domain.SupplierPerformanceReport, error) {
	s.legacySupplierCalls++
	return domain.SupplierPerformanceReport{}, nil
}

func (s *legacyReportStub) WarehouseValuation(context.Context, int) (domain.WarehouseValuationReport, error) {
	s.legacyWarehouseCalls++
	return domain.WarehouseValuationReport{}, nil
}

func TestMakeBoundedReportQueryPreservesPeriodAndBounds(t *testing.T) {
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	period, err := reporting.NewPeriod(from, to)
	if err != nil {
		t.Fatalf("new period: %v", err)
	}
	bounds, err := reporting.NewBounds(25, 50, 100, 5000)
	if err != nil {
		t.Fatalf("new bounds: %v", err)
	}
	query := makeBoundedReportQuery(period, bounds)
	if query.From == nil || query.To == nil || !query.From.Equal(from) || !query.To.Equal(to) {
		t.Fatalf("period not preserved: %+v", query)
	}
	if query.Limit != 25 || query.Offset != 50 {
		t.Fatalf("bounds not preserved: %+v", query)
	}
}

func TestSupplierPerformanceUsesBoundedCapability(t *testing.T) {
	stub := &reportCapabilityStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/supplier-performance?days=14&limit=25&offset=50", nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.supplierCalls != 1 || stub.legacySupplierCalls != 0 {
		t.Fatalf("bounded calls=%d legacy calls=%d", stub.supplierCalls, stub.legacySupplierCalls)
	}
	if stub.supplierQuery.Limit != 25 || stub.supplierQuery.Offset != 50 {
		t.Fatalf("query bounds=%+v", stub.supplierQuery)
	}
	if stub.supplierQuery.From == nil || stub.supplierQuery.To == nil {
		t.Fatal("supplier query period is missing")
	}
	if got := stub.supplierQuery.To.Sub(*stub.supplierQuery.From).Hours()/24 + 1; got != 14 {
		t.Fatalf("period days=%v", got)
	}
	if recorder.Header().Get("X-Report-Limit") != "25" || recorder.Header().Get("X-Report-Offset") != "50" {
		t.Fatalf("metadata bounds: limit=%q offset=%q", recorder.Header().Get("X-Report-Limit"), recorder.Header().Get("X-Report-Offset"))
	}
}

func TestSupplierPerformanceExposesTotalCount(t *testing.T) {
	stub := &countedReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/supplier-performance?days=21&limit=10&offset=20", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.supplierCalls != 1 || stub.supplierCountCalls != 1 {
		t.Fatalf("report calls=%d count calls=%d", stub.supplierCalls, stub.supplierCountCalls)
	}
	if got := recorder.Header().Get("X-Total-Count"); got != "37" {
		t.Fatalf("X-Total-Count=%q", got)
	}
	if stub.supplierQuery.Limit != 10 || stub.supplierQuery.Offset != 20 {
		t.Fatalf("count query bounds=%+v", stub.supplierQuery)
	}
	if stub.supplierQuery.From == nil || stub.supplierQuery.To == nil {
		t.Fatal("count query period is missing")
	}
}

func TestSupplierPerformanceFallsBackToLegacyRepositoryWithoutOffset(t *testing.T) {
	stub := &legacyReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/supplier-performance?days=10&limit=20", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.legacySupplierCalls != 1 {
		t.Fatalf("legacy supplier calls=%d", stub.legacySupplierCalls)
	}
	if recorder.Header().Get("X-Report-Limit") != "20" || recorder.Header().Get("X-Report-Offset") != "0" {
		t.Fatalf("metadata bounds: limit=%q offset=%q", recorder.Header().Get("X-Report-Limit"), recorder.Header().Get("X-Report-Offset"))
	}
}

func TestSupplierPerformanceRejectsOffsetWithoutBoundedCapability(t *testing.T) {
	stub := &legacyReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/supplier-performance?offset=7", nil))
	if recorder.Code != http.StatusNotImplemented {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.legacySupplierCalls != 0 {
		t.Fatal("legacy repository was called despite unsupported offset")
	}
}

func TestWarehouseValuationUsesBoundedCapability(t *testing.T) {
	stub := &reportCapabilityStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/warehouse-valuation?limit=40&offset=12", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.warehouseCalls != 1 || stub.legacyWarehouseCalls != 0 {
		t.Fatalf("bounded calls=%d legacy calls=%d", stub.warehouseCalls, stub.legacyWarehouseCalls)
	}
	if stub.warehouseQuery.Limit != 40 || stub.warehouseQuery.Offset != 12 {
		t.Fatalf("query bounds=%+v", stub.warehouseQuery)
	}
	if stub.warehouseQuery.From == nil || stub.warehouseQuery.To == nil || !stub.warehouseQuery.From.Equal(*stub.warehouseQuery.To) {
		t.Fatalf("warehouse snapshot period=%+v", stub.warehouseQuery)
	}
}

func TestWarehouseValuationExposesTotalCount(t *testing.T) {
	stub := &countedReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/warehouse-valuation?limit=8&offset=16", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.warehouseCalls != 1 || stub.warehouseCountCalls != 1 {
		t.Fatalf("report calls=%d count calls=%d", stub.warehouseCalls, stub.warehouseCountCalls)
	}
	if got := recorder.Header().Get("X-Total-Count"); got != "23" {
		t.Fatalf("X-Total-Count=%q", got)
	}
	if stub.warehouseQuery.Limit != 8 || stub.warehouseQuery.Offset != 16 {
		t.Fatalf("count query bounds=%+v", stub.warehouseQuery)
	}
}

func TestWarehouseValuationFallsBackToLegacyRepositoryWithoutOffset(t *testing.T) {
	stub := &legacyReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/warehouse-valuation?limit=40", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.legacyWarehouseCalls != 1 {
		t.Fatalf("legacy warehouse calls=%d", stub.legacyWarehouseCalls)
	}
}

func TestWarehouseValuationRejectsOffsetWithoutBoundedCapability(t *testing.T) {
	stub := &legacyReportStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/warehouse-valuation?offset=12", nil))
	if recorder.Code != http.StatusNotImplemented {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if stub.legacyWarehouseCalls != 0 {
		t.Fatal("legacy repository was called despite unsupported offset")
	}
}

func TestBoundedReportEndpointsRejectInvalidOffset(t *testing.T) {
	stub := &reportCapabilityStub{}
	handler := NewCore(nil, nil, nil, nil, WithInsights(stub, nil))
	for _, path := range []string{
		"/api/v1/reports/supplier-performance?offset=-1",
		"/api/v1/reports/warehouse-valuation?offset=not-a-number",
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("path=%s status=%d body=%s", path, recorder.Code, recorder.Body.String())
		}
	}
	if stub.supplierCalls != 0 || stub.warehouseCalls != 0 {
		t.Fatal("invalid offsets reached bounded repository")
	}
}

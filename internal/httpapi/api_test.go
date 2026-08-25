package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

type fakeStore struct {
	createdProduct     domain.Product
	reorderSuggestions []domain.ReorderSuggestion
	valuation          domain.InventoryValuationReport
}

func (f *fakeStore) CreateCategory(context.Context, domain.Category) error { return nil }
func (f *fakeStore) ListCategories(context.Context) ([]domain.Category, error) { return nil, nil }
func (f *fakeStore) CreateSupplier(context.Context, domain.Supplier) error { return nil }
func (f *fakeStore) ListSuppliers(context.Context, bool) ([]domain.Supplier, error) { return nil, nil }
func (f *fakeStore) CreateProduct(_ context.Context, product domain.Product) error {
	f.createdProduct = product
	return nil
}
func (f *fakeStore) UpdateProduct(context.Context, domain.Product) error { return nil }
func (f *fakeStore) GetProduct(_ context.Context, id string) (domain.Product, error) {
	if f.createdProduct.ID == id {
		return f.createdProduct, nil
	}
	return domain.Product{}, domain.ErrNotFound
}
func (f *fakeStore) GetProductByBarcode(_ context.Context, barcode string) (domain.Product, error) {
	if f.createdProduct.Barcode != "" && f.createdProduct.Barcode == barcode {
		return f.createdProduct, nil
	}
	return domain.Product{}, domain.ErrNotFound
}
func (f *fakeStore) ListProducts(context.Context, repository.ProductFilter) ([]domain.Product, error) { return nil, nil }
func (f *fakeStore) CreateWarehouse(context.Context, domain.Warehouse) error { return nil }
func (f *fakeStore) ListWarehouses(context.Context, bool) ([]domain.Warehouse, error) { return nil, nil }
func (f *fakeStore) CreateLocation(context.Context, domain.Location) error { return nil }
func (f *fakeStore) ListLocations(context.Context, string, bool) ([]domain.Location, error) { return nil, nil }
func (f *fakeStore) CreateLot(context.Context, domain.Lot) error { return nil }
func (f *fakeStore) ApplyMovement(context.Context, domain.StockMovement) (domain.StockBalance, error) {
	return domain.StockBalance{}, nil
}
func (f *fakeStore) Transfer(context.Context, domain.TransferRequest) error { return nil }
func (f *fakeStore) GetBalance(context.Context, string, string, string) (domain.StockBalance, error) {
	return domain.StockBalance{}, nil
}
func (f *fakeStore) ListLowStock(context.Context, int) ([]domain.StockBalance, error) { return nil, nil }
func (f *fakeStore) ListReorderSuggestions(context.Context, int) ([]domain.ReorderSuggestion, error) {
	return f.reorderSuggestions, nil
}
func (f *fakeStore) GetInventoryValuation(context.Context, int) (domain.InventoryValuationReport, error) {
	return f.valuation, nil
}
func (f *fakeStore) CreateOrder(context.Context, domain.PurchaseOrder) error { return nil }
func (f *fakeStore) GetOrder(context.Context, string) (domain.PurchaseOrder, error) { return domain.PurchaseOrder{}, nil }
func (f *fakeStore) ListOrders(context.Context, domain.PurchaseOrderStatus, int, int) ([]domain.PurchaseOrder, error) {
	return nil, nil
}
func (f *fakeStore) ReceiveLine(context.Context, string, string, int64, string, string) error { return nil }

func TestHealthEndpoint(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `"status":"ok"`) {
		t.Fatalf("body = %q, want status ok", recorder.Body.String())
	}
	if recorder.Header().Get("X-Request-ID") == "" {
		t.Fatal("X-Request-ID header is missing")
	}
	if recorder.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q", recorder.Header().Get("X-Content-Type-Options"))
	}
}

func TestReadyEndpointReportsDependencyFailure(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return errors.New("database unavailable") }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/readyz", nil))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
}

func TestMiddlewareRejectsUnknownOrigin(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, []string{"https://stock.example"}, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	request.Header.Set("Origin", "https://evil.example")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("disallowed origin must not receive Access-Control-Allow-Origin")
	}
}

func TestCreateProductRejectsUnknownJSONFields(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	body := `{"sku":"SKU-1","name":"Widget","unit":"piece","mystery":true}`
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(body)))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestCreateProductGeneratesIDAndDefaultsCurrency(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	body := `{"sku":"SKU-1","name":"Widget","unit":"piece","active":true}`
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(body)))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	if !strings.HasPrefix(store.createdProduct.ID, "prd_") {
		t.Fatalf("created product id = %q, want prd_ prefix", store.createdProduct.ID)
	}
	if store.createdProduct.Currency != "INR" {
		t.Fatalf("currency = %q, want INR", store.createdProduct.Currency)
	}
}

func TestBarcodeLookupReturnsMatchingProduct(t *testing.T) {
	store := &fakeStore{createdProduct: domain.Product{ID: "prd_1", SKU: "SKU-1", Name: "Widget", Barcode: "8901234567890", Unit: "piece", Currency: "INR", Active: true}}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/products/by-barcode/8901234567890", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"sku":"SKU-1"`) {
		t.Fatalf("body = %q, want matching SKU", recorder.Body.String())
	}
}

func TestReorderSuggestionsEndpoint(t *testing.T) {
	store := &fakeStore{reorderSuggestions: []domain.ReorderSuggestion{{ProductID: "prd_1", SKU: "SKU-1", Name: "Widget", OnHand: 3, ReorderPoint: 5, ReorderQuantity: 10, TargetStock: 15, SuggestedQuantity: 12}}}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/reorder-suggestions?limit=10", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"suggestedQuantity":12`) {
		t.Fatalf("body = %q, want suggested quantity", recorder.Body.String())
	}
}

func TestInventoryValuationEndpoint(t *testing.T) {
	store := &fakeStore{valuation: domain.InventoryValuationReport{
		Items: []domain.InventoryValuationItem{{ProductID: "prd_1", SKU: "SKU-1", Name: "Widget", Unit: "piece", OnHand: 4, UnitCostMinor: 2500, Currency: "INR", ValueMinor: 10000}},
		Totals: []domain.InventoryValuationTotal{{Currency: "INR", ValueMinor: 10000}},
	}}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/inventory-valuation", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"valueMinor":10000`) || !strings.Contains(recorder.Body.String(), `"currency":"INR"`) {
		t.Fatalf("body = %q, want INR valuation", recorder.Body.String())
	}
}

package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/auth"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

type fakeStore struct {
	createdProduct domain.Product
	movementActor  string
	transferActor  string
	createdOrder   domain.PurchaseOrder
	receiptActor   string
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
func (f *fakeStore) ListProducts(context.Context, repository.ProductFilter) ([]domain.Product, error) {
	return nil, nil
}
func (f *fakeStore) CreateWarehouse(context.Context, domain.Warehouse) error { return nil }
func (f *fakeStore) ListWarehouses(context.Context, bool) ([]domain.Warehouse, error) { return nil, nil }
func (f *fakeStore) CreateLocation(context.Context, domain.Location) error { return nil }
func (f *fakeStore) ListLocations(context.Context, string, bool) ([]domain.Location, error) {
	return nil, nil
}
func (f *fakeStore) CreateLot(context.Context, domain.Lot) error { return nil }
func (f *fakeStore) ApplyMovement(_ context.Context, movement domain.StockMovement) (domain.StockBalance, error) {
	f.movementActor = movement.ActorID
	return domain.StockBalance{}, nil
}
func (f *fakeStore) Transfer(_ context.Context, request domain.TransferRequest) error {
	f.transferActor = request.ActorID
	return nil
}
func (f *fakeStore) GetBalance(context.Context, string, string, string) (domain.StockBalance, error) {
	return domain.StockBalance{}, nil
}
func (f *fakeStore) ListLowStock(context.Context, int) ([]domain.StockBalance, error) { return nil, nil }
func (f *fakeStore) CreateOrder(_ context.Context, order domain.PurchaseOrder) error {
	f.createdOrder = order
	return nil
}
func (f *fakeStore) GetOrder(context.Context, string) (domain.PurchaseOrder, error) {
	return domain.PurchaseOrder{}, nil
}
func (f *fakeStore) ListOrders(context.Context, domain.PurchaseOrderStatus, int, int) ([]domain.PurchaseOrder, error) {
	return nil, nil
}
func (f *fakeStore) ReceiveLine(_ context.Context, _ string, _ string, _ int64, _ string, _ string, actorID string) error {
	f.receiptActor = actorID
	return nil
}

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

func TestInventoryMutationUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(
		http.MethodPost,
		"/api/v1/inventory/movements",
		`{"productId":"prd_1","locationId":"loc_1","type":"stock_in","quantityDelta":5,"actorId":"usr_spoofed"}`,
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	if store.movementActor != "usr_session" {
		t.Fatalf("movement actor = %q, want authenticated user", store.movementActor)
	}
}

func TestTransferUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(
		http.MethodPost,
		"/api/v1/inventory/transfers",
		`{"productId":"prd_1","fromLocationId":"loc_a","toLocationId":"loc_b","quantity":2,"actorId":"usr_spoofed"}`,
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	if store.transferActor != "usr_session" {
		t.Fatalf("transfer actor = %q, want authenticated user", store.transferActor)
	}
}

func TestCreateOrderUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(
		http.MethodPost,
		"/api/v1/orders",
		`{"number":"PO-1","supplierId":"sup_1","warehouseId":"wh_1","createdBy":"usr_spoofed","lines":[]}`,
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	if store.createdOrder.CreatedBy != "usr_session" {
		t.Fatalf("createdBy = %q, want authenticated user", store.createdOrder.CreatedBy)
	}
}

func TestReceiveOrderLineUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(
		http.MethodPost,
		"/api/v1/orders/po_1/lines/pol_1/receive",
		`{"quantity":1,"locationId":"loc_1"}`,
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if store.receiptActor != "usr_session" {
		t.Fatalf("receipt actor = %q, want authenticated user", store.receiptActor)
	}
}

func authenticatedRequest(method, target, body string) *http.Request {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	principal := auth.Principal{User: domain.User{ID: "usr_session", Role: domain.RoleAdmin, Active: true}}
	return request.WithContext(context.WithValue(request.Context(), principalContextKey{}, principal))
}

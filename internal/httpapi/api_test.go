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
	createdProduct     domain.Product
	movementActor      string
	transferActor      string
	createdOrder       domain.PurchaseOrder
	receiptActor       string
	reorderSuggestions []domain.ReorderSuggestion
	valuation          domain.InventoryValuationReport
}

func (f *fakeStore) CreateCategory(context.Context, domain.Category) error { return nil }
func (f *fakeStore) ListCategories(context.Context) ([]domain.Category, error) { return nil, nil }
func (f *fakeStore) CreateSupplier(context.Context, domain.Supplier) error { return nil }
func (f *fakeStore) ListSuppliers(context.Context, bool) ([]domain.Supplier, error) { return nil, nil }
func (f *fakeStore) CreateProduct(_ context.Context, product domain.Product) error { f.createdProduct = product; return nil }
func (f *fakeStore) UpdateProduct(context.Context, domain.Product) error { return nil }
func (f *fakeStore) GetProduct(_ context.Context, id string) (domain.Product, error) { if f.createdProduct.ID == id { return f.createdProduct, nil }; return domain.Product{}, domain.ErrNotFound }
func (f *fakeStore) GetProductByBarcode(_ context.Context, barcode string) (domain.Product, error) { if f.createdProduct.Barcode != "" && f.createdProduct.Barcode == barcode { return f.createdProduct, nil }; return domain.Product{}, domain.ErrNotFound }
func (f *fakeStore) ListProducts(context.Context, repository.ProductFilter) ([]domain.Product, error) { return nil, nil }
func (f *fakeStore) CreateWarehouse(context.Context, domain.Warehouse) error { return nil }
func (f *fakeStore) ListWarehouses(context.Context, bool) ([]domain.Warehouse, error) { return nil, nil }
func (f *fakeStore) CreateLocation(context.Context, domain.Location) error { return nil }
func (f *fakeStore) ListLocations(context.Context, string, bool) ([]domain.Location, error) { return nil, nil }
func (f *fakeStore) CreateLot(context.Context, domain.Lot) error { return nil }
func (f *fakeStore) ApplyMovement(_ context.Context, movement domain.StockMovement) (domain.StockBalance, error) { f.movementActor = movement.ActorID; return domain.StockBalance{}, nil }
func (f *fakeStore) Transfer(_ context.Context, request domain.TransferRequest) error { f.transferActor = request.ActorID; return nil }
func (f *fakeStore) GetBalance(context.Context, string, string, string) (domain.StockBalance, error) { return domain.StockBalance{}, nil }
func (f *fakeStore) ListLowStock(context.Context, int) ([]domain.StockBalance, error) { return nil, nil }
func (f *fakeStore) ListReorderSuggestions(context.Context, int) ([]domain.ReorderSuggestion, error) { return f.reorderSuggestions, nil }
func (f *fakeStore) GetInventoryValuation(context.Context, int) (domain.InventoryValuationReport, error) { return f.valuation, nil }
func (f *fakeStore) CreateOrder(_ context.Context, order domain.PurchaseOrder) error { f.createdOrder = order; return nil }
func (f *fakeStore) GetOrder(context.Context, string) (domain.PurchaseOrder, error) { return domain.PurchaseOrder{}, nil }
func (f *fakeStore) ListOrders(context.Context, domain.PurchaseOrderStatus, int, int) ([]domain.PurchaseOrder, error) { return nil, nil }
func (f *fakeStore) ReceiveLine(_ context.Context, _ string, _ string, _ int64, _ string, _ string, actorID string) error { f.receiptActor = actorID; return nil }

func TestHealthEndpoint(t *testing.T) {
	store := &fakeStore{}; handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusOK { t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK) }
	if !strings.Contains(recorder.Body.String(), `"status":"ok"`) { t.Fatalf("body = %q, want status ok", recorder.Body.String()) }
	if recorder.Header().Get("X-Request-ID") == "" { t.Fatal("X-Request-ID header is missing") }
	if recorder.Header().Get("X-Content-Type-Options") != "nosniff" { t.Fatalf("X-Content-Type-Options = %q", recorder.Header().Get("X-Content-Type-Options")) }
}

func TestReadyEndpointReportsDependencyFailure(t *testing.T) {
	store := &fakeStore{}; handler := New(store, store, store, func(context.Context) error { return errors.New("database unavailable") }, nil, nil); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if recorder.Code != http.StatusServiceUnavailable { t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable) }
}

func TestMiddlewareRejectsUnknownOrigin(t *testing.T) {
	store := &fakeStore{}; handler := New(store, store, store, func(context.Context) error { return nil }, []string{"https://stock.example"}, nil); request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil); request.Header.Set("Origin", "https://evil.example"); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden { t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden) }
	if recorder.Header().Get("Access-Control-Allow-Origin") != "" { t.Fatal("disallowed origin must not receive Access-Control-Allow-Origin") }
}

func TestCreateProductRejectsUnknownJSONFields(t *testing.T) {
	store := &fakeStore{}; handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil); body := `{"sku":"SKU-1","name":"Widget","unit":"piece","mystery":true}`; recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(body)))
	if recorder.Code != http.StatusBadRequest { t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest) }
}

func TestCreateProductGeneratesIDAndDefaultsCurrency(t *testing.T) {
	store := &fakeStore{}; handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil); body := `{"sku":"SKU-1","name":"Widget","unit":"piece","active":true}`; recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(body)))
	if recorder.Code != http.StatusCreated { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String()) }
	if !strings.HasPrefix(store.createdProduct.ID, "prd_") { t.Fatalf("created product id = %q, want prd_ prefix", store.createdProduct.ID) }
	if store.createdProduct.Currency != "INR" { t.Fatalf("currency = %q, want INR", store.createdProduct.Currency) }
}

func TestInventoryMutationUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}; handler := NewCore(store, store, store, func(context.Context) error { return nil }); request := authenticatedRequest(http.MethodPost, "/api/v1/inventory/movements", `{"productId":"prd_1","locationId":"loc_1","type":"stock_in","quantityDelta":5,"actorId":"usr_spoofed"}`); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String()) }; if store.movementActor != "usr_session" { t.Fatalf("movement actor = %q, want authenticated user", store.movementActor) }
}

func TestTransferUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}; handler := NewCore(store, store, store, func(context.Context) error { return nil }); request := authenticatedRequest(http.MethodPost, "/api/v1/inventory/transfers", `{"productId":"prd_1","fromLocationId":"loc_a","toLocationId":"loc_b","quantity":2,"actorId":"usr_spoofed"}`); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String()) }; if store.transferActor != "usr_session" { t.Fatalf("transfer actor = %q, want authenticated user", store.transferActor) }
}

func TestCreateOrderUsesAuthenticatedActorAndGeneratesIDs(t *testing.T) {
	store := &fakeStore{}; handler := NewCore(store, store, store, func(context.Context) error { return nil }); request := authenticatedRequest(http.MethodPost, "/api/v1/orders", `{"number":"PO-100","supplierId":"sup_1","warehouseId":"wh_1","createdBy":"usr_spoofed","status":"draft","currency":"INR","lines":[{"productId":"prd_1","quantity":5,"received":0,"unitCostMinor":2500}]}`); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusCreated, recorder.Body.String()) }
	if store.createdOrder.CreatedBy != "usr_session" { t.Fatalf("createdBy = %q, want authenticated user", store.createdOrder.CreatedBy) }
	if !strings.HasPrefix(store.createdOrder.ID, "po_") { t.Fatalf("order id = %q, want po_ prefix", store.createdOrder.ID) }
	if len(store.createdOrder.Lines) != 1 { t.Fatalf("lines = %d, want 1", len(store.createdOrder.Lines)) }
	if !strings.HasPrefix(store.createdOrder.Lines[0].ID, "pol_") { t.Fatalf("line id = %q, want pol_ prefix", store.createdOrder.Lines[0].ID) }
	if store.createdOrder.Lines[0].PurchaseOrderID != store.createdOrder.ID { t.Fatalf("line order id = %q, want %q", store.createdOrder.Lines[0].PurchaseOrderID, store.createdOrder.ID) }
}

func TestReceiveOrderLineUsesAuthenticatedActor(t *testing.T) {
	store := &fakeStore{}; handler := NewCore(store, store, store, func(context.Context) error { return nil }); request := authenticatedRequest(http.MethodPost, "/api/v1/orders/po_1/lines/pol_1/receive", `{"quantity":1,"locationId":"loc_1"}`); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String()) }
	if store.receiptActor != "usr_session" { t.Fatalf("receipt actor = %q, want authenticated user", store.receiptActor) }
}

func TestBarcodeLookupReturnsMatchingProduct(t *testing.T) {
	store := &fakeStore{createdProduct: domain.Product{ID: "prd_1", SKU: "SKU-1", Name: "Widget", Barcode: "8901234567890", Unit: "piece", Currency: "INR", Active: true}}; handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/products/by-barcode/8901234567890", nil))
	if recorder.Code != http.StatusOK { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String()) }; if !strings.Contains(recorder.Body.String(), `"sku":"SKU-1"`) { t.Fatalf("body = %q, want matching SKU", recorder.Body.String()) }
}

func TestReorderSuggestionsEndpoint(t *testing.T) {
	store := &fakeStore{reorderSuggestions: []domain.ReorderSuggestion{{ProductID: "prd_1", SKU: "SKU-1", Name: "Widget", OnHand: 3, ReorderPoint: 5, ReorderQuantity: 10, TargetStock: 15, SuggestedQuantity: 12}}}; handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/inventory/reorder-suggestions?limit=10", nil))
	if recorder.Code != http.StatusOK { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String()) }; if !strings.Contains(recorder.Body.String(), `"suggestedQuantity":12`) { t.Fatalf("body = %q, want suggested quantity", recorder.Body.String()) }
}

func TestInventoryValuationEndpoint(t *testing.T) {
	store := &fakeStore{valuation: domain.InventoryValuationReport{Items: []domain.InventoryValuationItem{{ProductID: "prd_1", SKU: "SKU-1", Name: "Widget", Unit: "piece", OnHand: 4, UnitCostMinor: 2500, Currency: "INR", ValueMinor: 10000}}, Totals: []domain.InventoryValuationTotal{{Currency: "INR", ValueMinor: 10000}}}}; handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil); recorder := httptest.NewRecorder(); handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/reports/inventory-valuation", nil))
	if recorder.Code != http.StatusOK { t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String()) }; if !strings.Contains(recorder.Body.String(), `"valueMinor":10000`) || !strings.Contains(recorder.Body.String(), `"currency":"INR"`) { t.Fatalf("body = %q, want INR valuation", recorder.Body.String()) }
}

func authenticatedRequest(method, target, body string) *http.Request {
	request := httptest.NewRequest(method, target, strings.NewReader(body)); principal := auth.Principal{User: domain.User{ID: "usr_session", Role: domain.RoleAdmin, Active: true}}; return request.WithContext(context.WithValue(request.Context(), principalContextKey{}, principal))
}

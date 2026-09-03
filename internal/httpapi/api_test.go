package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	auditEvents        []domain.AuditEvent
	updatedOrder       domain.PurchaseOrder
	updatedWarehouse   domain.Warehouse
	updatedLocation    domain.Location
	lotInventory       []domain.LotInventoryRow
	balances           []domain.StockBalance
}

func (f *fakeStore) CreateCategory(context.Context, domain.Category) error          { return nil }
func (f *fakeStore) ListCategories(context.Context) ([]domain.Category, error)      { return nil, nil }
func (f *fakeStore) CreateSupplier(context.Context, domain.Supplier) error          { return nil }
func (f *fakeStore) ListSuppliers(context.Context, bool) ([]domain.Supplier, error) { return nil, nil }
func (f *fakeStore) CreateProduct(_ context.Context, p domain.Product) error {
	f.createdProduct = p
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
func (f *fakeStore) ListProducts(context.Context, repository.ProductFilter) ([]domain.Product, error) {
	return nil, nil
}
func (f *fakeStore) CreateWarehouse(context.Context, domain.Warehouse) error { return nil }
func (f *fakeStore) UpdateWarehouse(_ context.Context, v domain.Warehouse) error {
	f.updatedWarehouse = v
	return nil
}
func (f *fakeStore) ListWarehouses(context.Context, bool) ([]domain.Warehouse, error) {
	return []domain.Warehouse{{ID: "wh_1", Code: "MAIN", Name: "Main", Timezone: "Asia/Kolkata", Active: true}}, nil
}
func (f *fakeStore) CreateLocation(context.Context, domain.Location) error { return nil }
func (f *fakeStore) UpdateLocation(_ context.Context, v domain.Location) error {
	f.updatedLocation = v
	return nil
}
func (f *fakeStore) ListLocations(context.Context, string, bool) ([]domain.Location, error) {
	return []domain.Location{{ID: "loc_1", WarehouseID: "wh_1", Code: "A1", Name: "A1", Active: true}}, nil
}
func (f *fakeStore) CreateLot(context.Context, domain.Lot) error { return nil }
func (f *fakeStore) ListLots(_ context.Context, productID string, _ int) ([]domain.Lot, error) {
	if productID == "prd_1" {
		return []domain.Lot{{ID: "lot_1", ProductID: "prd_1", LotNumber: "LOT-1"}}, nil
	}
	return []domain.Lot{}, nil
}
func (f *fakeStore) ListLotInventory(_ context.Context, filter repository.LotInventoryFilter) ([]domain.LotInventoryRow, error) {
	if filter.ExpiringBy != nil {
		items := make([]domain.LotInventoryRow, 0)
		for _, item := range f.lotInventory {
			if item.ExpiresAt == nil {
				continue
			}
			if !item.ExpiresAt.After(*filter.ExpiringBy) {
				items = append(items, item)
			}
		}
		return items, nil
	}
	return f.lotInventory, nil
}
func (f *fakeStore) ApplyMovement(_ context.Context, m domain.StockMovement) (domain.StockBalance, error) {
	f.movementActor = m.ActorID
	return domain.StockBalance{}, nil
}
func (f *fakeStore) Transfer(_ context.Context, r domain.TransferRequest) error {
	f.transferActor = r.ActorID
	return nil
}
func (f *fakeStore) GetBalance(context.Context, string, string, string) (domain.StockBalance, error) {
	return domain.StockBalance{}, nil
}
func (f *fakeStore) ListBalances(context.Context, int, int) ([]domain.StockBalance, error) {
	return f.balances, nil
}
func (f *fakeStore) ListLowStock(context.Context, int) ([]domain.StockBalance, error) {
	return nil, nil
}
func (f *fakeStore) ListReorderSuggestions(context.Context, int) ([]domain.ReorderSuggestion, error) {
	return f.reorderSuggestions, nil
}
func (f *fakeStore) GetInventoryValuation(context.Context, int) (domain.InventoryValuationReport, error) {
	return f.valuation, nil
}
func (f *fakeStore) CreateOrder(_ context.Context, o domain.PurchaseOrder) error {
	f.createdOrder = o
	return nil
}
func (f *fakeStore) UpdateOrder(_ context.Context, o domain.PurchaseOrder) error {
	f.updatedOrder = o
	return nil
}
func (f *fakeStore) GetOrder(_ context.Context, id string) (domain.PurchaseOrder, error) {
	if f.updatedOrder.ID == id {
		return f.updatedOrder, nil
	}
	return domain.PurchaseOrder{ID: id, Status: domain.PurchaseOrderDraft, Lines: []domain.PurchaseOrderLine{{ID: "pol_1", PurchaseOrderID: id, ProductID: "prd_1", Quantity: 5}}}, nil
}
func (f *fakeStore) ListOrders(context.Context, domain.PurchaseOrderStatus, int, int) ([]domain.PurchaseOrder, error) {
	return nil, nil
}
func (f *fakeStore) ReceiveLine(_ context.Context, _ string, _ string, _ int64, _ string, _ string, actorID string) error {
	f.receiptActor = actorID
	return nil
}
func (f *fakeStore) ReceiveLineWithNewLot(_ context.Context, _ string, _ string, _ int64, _ string, lot domain.Lot, actorID string) error {
	lot.ID = "lot_test"
	lot.ProductID = "prd_test"
	if err := lot.Validate(); err != nil {
		return err
	}
	f.receiptActor = actorID
	return nil
}
func (f *fakeStore) UpdateOrderStatus(_ context.Context, _ string, status domain.PurchaseOrderStatus, _ string) error {
	if err := domain.ValidatePurchaseOrderTransition(domain.PurchaseOrderDraft, status); err != nil {
		return err
	}
	return nil
}
func (f *fakeStore) AppendAuditEvent(_ context.Context, e domain.AuditEvent) error {
	f.auditEvents = append(f.auditEvents, e)
	return nil
}
func (f *fakeStore) ListAuditEvents(context.Context, domain.AuditFilter) ([]domain.AuditEvent, error) {
	return f.auditEvents, nil
}

func TestHealthEndpoint(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"status":"ok"`) {
		t.Fatalf("body=%q", recorder.Body.String())
	}
	if recorder.Header().Get("X-Request-ID") == "" {
		t.Fatal("X-Request-ID missing")
	}
}

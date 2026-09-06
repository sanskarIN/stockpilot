package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestReplenishmentTraceabilityIntegration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	store, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer store.Close()

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	productID := "prd_trace_" + suffix
	supplierID := "sup_trace_" + suffix
	warehouseID := "wh_trace_" + suffix
	locationID := "loc_trace_" + suffix
	orderID := "po_trace_" + suffix
	lineID := "pol_trace_" + suffix
	reviewID := "rr_trace_" + suffix

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM replenishment_reviews WHERE product_id=$1`, productID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM purchase_order_lines WHERE purchase_order_id=$1`, orderID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM purchase_orders WHERE id=$1`, orderID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM locations WHERE id=$1`, locationID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM warehouses WHERE id=$1`, warehouseID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM products WHERE id=$1`, productID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM suppliers WHERE id=$1`, supplierID)
	}()

	if err := store.CreateSupplier(ctx, domain.Supplier{ID: supplierID, Code: "TRACE-" + suffix, Name: "Trace supplier", Active: true}); err != nil {
		t.Fatalf("CreateSupplier() error = %v", err)
	}
	if err := store.CreateProduct(ctx, domain.Product{ID: productID, SKU: "TRACE-" + suffix, Name: "Trace product", SupplierID: supplierID, Unit: "piece", Currency: "INR", ReorderPoint: 10, ReorderQuantity: 20, Active: true}); err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}
	if err := store.CreateWarehouse(ctx, domain.Warehouse{ID: warehouseID, Code: "TRACE-" + suffix, Name: "Trace warehouse", Timezone: "UTC", Active: true}); err != nil {
		t.Fatalf("CreateWarehouse() error = %v", err)
	}
	if err := store.CreateLocation(ctx, domain.Location{ID: locationID, WarehouseID: warehouseID, Code: "TRACE-" + suffix, Name: "Trace location", Active: true}); err != nil {
		t.Fatalf("CreateLocation() error = %v", err)
	}

	suggestions, err := store.ListReorderSuggestions(ctx, 500)
	if err != nil {
		t.Fatalf("ListReorderSuggestions() error = %v", err)
	}
	var suggestion *domain.ReorderSuggestion
	for i := range suggestions {
		if suggestions[i].ProductID == productID {
			suggestion = &suggestions[i]
			break
		}
	}
	if suggestion == nil {
		t.Fatal("expected reorder suggestion")
	}
	order := domain.PurchaseOrder{
		ID: orderID, Number: "TRACE-" + suffix, SupplierID: supplierID, WarehouseID: warehouseID,
		Status: domain.PurchaseOrderDraft, Currency: "INR", CreatedBy: "user_trace",
		Lines: []domain.PurchaseOrderLine{{ID: lineID, ProductID: productID, Quantity: suggestion.SuggestedQuantity, UnitCostMinor: 100}},
	}
	if err := store.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	review := domain.ReplenishmentReview{
		ID: reviewID, ProductID: suggestion.ProductID, SKU: suggestion.SKU, ProductName: suggestion.Name, Unit: suggestion.Unit,
		SupplierID: suggestion.SupplierID, OnHand: suggestion.OnHand, ReorderPoint: suggestion.ReorderPoint,
		ReorderQuantity: suggestion.ReorderQuantity, TargetStock: suggestion.TargetStock, SuggestedQuantity: suggestion.SuggestedQuantity,
		ReviewedBy: "user_trace", ReviewedAt: time.Now().UTC(), PurchaseOrderID: orderID, Outcome: domain.ReplenishmentReviewAccepted,
	}
	if err := store.CreateReplenishmentReview(ctx, review); err != nil {
		t.Fatalf("CreateReplenishmentReview() error = %v", err)
	}
	items, err := store.ListReplenishmentReviews(ctx, "accepted", 50, 0)
	if err != nil {
		t.Fatalf("ListReplenishmentReviews() error = %v", err)
	}
	if len(items) != 1 || items[0].ID != reviewID || items[0].PurchaseOrderID != orderID || items[0].SuggestedQuantity != suggestion.SuggestedQuantity {
		t.Fatalf("review list = %+v, want persisted immutable snapshot", items)
	}
}

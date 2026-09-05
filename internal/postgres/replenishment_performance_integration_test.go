package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
)

func TestReplenishmentPerformanceIntegration(t *testing.T) {
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
	supplierID := "sup_repl_" + suffix
	warehouseID := "wh_repl_" + suffix
	locationID := "loc_repl_" + suffix
	productID := "prd_repl_" + suffix
	orderID := "po_repl_" + suffix
	orderLineID := "pol_repl_" + suffix

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM stock_movements WHERE reference=$1`, "PO:"+orderID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM purchase_order_lines WHERE id=$1`, orderLineID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM purchase_orders WHERE id=$1`, orderID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM locations WHERE id=$1`, locationID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM warehouses WHERE id=$1`, warehouseID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM products WHERE id=$1`, productID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM suppliers WHERE id=$1`, supplierID)
	}()

	if err := store.CreateSupplier(ctx, domain.Supplier{ID: supplierID, Code: "REPL-" + suffix, Name: "Replenishment integration supplier", Active: true}); err != nil {
		t.Fatalf("CreateSupplier() error = %v", err)
	}
	if err := store.CreateProduct(ctx, domain.Product{ID: productID, SKU: "REPL-" + suffix, Name: "Replenishment integration product", Unit: "piece", UnitCostMinor: 100, Currency: "INR", ReorderPoint: 5, ReorderQuantity: 10, Active: true, SupplierID: supplierID}); err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}
	if err := store.CreateWarehouse(ctx, domain.Warehouse{ID: warehouseID, Code: "REPL-" + suffix, Name: "Replenishment integration warehouse", Timezone: "UTC", Active: true}); err != nil {
		t.Fatalf("CreateWarehouse() error = %v", err)
	}
	if err := store.CreateLocation(ctx, domain.Location{ID: locationID, WarehouseID: warehouseID, Code: "REPL-" + suffix, Name: "Replenishment integration location", Active: true}); err != nil {
		t.Fatalf("CreateLocation() error = %v", err)
	}

	createdAt := time.Now().UTC().Add(-48 * time.Hour)
	expectedAt := createdAt.Add(24 * time.Hour)
	if err := store.CreateOrder(ctx, domain.PurchaseOrder{ID: orderID, Number: "REPL-" + suffix, SupplierID: supplierID, WarehouseID: warehouseID, Status: domain.PurchaseOrderOrdered, Currency: "INR", ExpectedAt: &expectedAt, CreatedAt: createdAt, UpdatedAt: createdAt, Lines: []domain.PurchaseOrderLine{{ID: orderLineID, PurchaseOrderID: orderID, ProductID: productID, Quantity: 20, Received: 0, UnitCostMinor: 100}}}); err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	if _, err := store.pool.Exec(ctx, `UPDATE purchase_orders SET created_at=$2, updated_at=$2 WHERE id=$1`, orderID, createdAt); err != nil {
		t.Fatalf("backdate order: %v", err)
	}
	if _, err := store.pool.Exec(ctx, `INSERT INTO stock_movements (id, product_id, location_id, type, quantity_delta, reference, occurred_at) VALUES ($1,$2,$3,'receive',15,$4,$5)`, "mov_repl_"+suffix, productID, locationID, "PO:"+orderID, createdAt.Add(12*time.Hour)); err != nil {
		t.Fatalf("insert receipt movement: %v", err)
	}

	report, err := store.ReplenishmentPerformance(ctx, reporting.Query{From: timePtr(createdAt.Add(-time.Hour)), To: timePtr(createdAt.Add(72 * time.Hour)), Limit: 10})
	if err != nil {
		t.Fatalf("ReplenishmentPerformance() error = %v", err)
	}
	if len(report.Items) != 1 {
		t.Fatalf("items=%d, want 1", len(report.Items))
	}
	item := report.Items[0]
	if item.OrderCount != 1 || item.OrderedUnits != 20 || item.ReceivedUnits != 15 || item.OutstandingUnits != 5 {
		t.Fatalf("item=%+v", item)
	}
	if item.OnTimeOrderCount != 1 || item.LateOrderCount != 0 || item.FillRate != 0.75 {
		t.Fatalf("timeliness/fill metrics=%+v", item)
	}
}

func timePtr(value time.Time) *time.Time { return &value }

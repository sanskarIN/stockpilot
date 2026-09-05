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

func TestReportingIntegration(t *testing.T) {
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
	productID := "prd_report_" + suffix
	warehouseID := "wh_report_" + suffix
	locationID := "loc_report_" + suffix
	movementID := "mov_report_" + suffix
	barcode := "report-" + suffix

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM stock_movements WHERE product_id=$1`, productID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM inventory_balances WHERE product_id=$1`, productID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM locations WHERE id=$1`, locationID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM warehouses WHERE id=$1`, warehouseID)
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM products WHERE id=$1`, productID)
	}()

	product := domain.Product{
		ID: productID, SKU: "REPORT-" + suffix, Name: "Reporting integration product", Barcode: barcode,
		Unit: "piece", UnitCostMinor: 250, Currency: "INR", ReorderPoint: 10, ReorderQuantity: 20, Active: true,
	}
	if err := store.CreateProduct(ctx, product); err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}
	if err := store.CreateWarehouse(ctx, domain.Warehouse{ID: warehouseID, Code: "REPORT-" + suffix, Name: "Reporting integration warehouse", Timezone: "UTC", Active: true}); err != nil {
		t.Fatalf("CreateWarehouse() error = %v", err)
	}
	if err := store.CreateLocation(ctx, domain.Location{ID: locationID, WarehouseID: warehouseID, Code: "REPORT-" + suffix, Name: "Reporting integration location", Active: true}); err != nil {
		t.Fatalf("CreateLocation() error = %v", err)
	}
	if _, err := store.ApplyMovement(ctx, domain.StockMovement{
		ID: movementID, ProductID: productID, LocationID: locationID, Type: domain.MovementStockIn,
		QuantityDelta: 4, OccurredAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("ApplyMovement() error = %v", err)
	}

	lookup, err := store.GetProductByBarcode(ctx, barcode)
	if err != nil {
		t.Fatalf("GetProductByBarcode() error = %v", err)
	}
	if lookup.ID != productID {
		t.Fatalf("GetProductByBarcode().ID = %q, want %q", lookup.ID, productID)
	}

	suggestions, err := store.ListReorderSuggestions(ctx, 500)
	if err != nil {
		t.Fatalf("ListReorderSuggestions() error = %v", err)
	}
	var suggestion *domain.ReorderSuggestion
	for index := range suggestions {
		if suggestions[index].ProductID == productID {
			suggestion = &suggestions[index]
			break
		}
	}
	if suggestion == nil {
		t.Fatalf("reorder suggestions did not include %s", productID)
	}
	if suggestion.OnHand != 4 || suggestion.TargetStock != 30 || suggestion.SuggestedQuantity != 26 {
		t.Fatalf("suggestion = %+v, want onHand=4 target=30 suggested=26", *suggestion)
	}

	report, err := store.GetInventoryValuation(ctx, 1000)
	if err != nil {
		t.Fatalf("GetInventoryValuation() error = %v", err)
	}
	var valuation *domain.InventoryValuationItem
	for index := range report.Items {
		if report.Items[index].ProductID == productID {
			valuation = &report.Items[index]
			break
		}
	}
	if valuation == nil {
		t.Fatalf("valuation report did not include %s", productID)
	}
	if valuation.OnHand != 4 || valuation.ValueMinor != 1000 || valuation.Currency != "INR" {
		t.Fatalf("valuation = %+v, want onHand=4 valueMinor=1000 currency=INR", *valuation)
	}

	bounded, err := store.WarehouseValuationQuery(ctx, reporting.Query{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("WarehouseValuationQuery() error = %v", err)
	}
	if len(bounded.Items) != 0 {
		t.Fatalf("WarehouseValuationQuery() returned %d items after the only row", len(bounded.Items))
	}
	if len(bounded.Totals) == 0 {
		t.Fatal("WarehouseValuationQuery() dropped warehouse totals when paginating items")
	}

	warehouseCount, err := store.WarehouseValuationCount(ctx, reporting.Query{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("WarehouseValuationCount() error = %v", err)
	}
	if warehouseCount < 1 {
		t.Fatalf("WarehouseValuationCount() = %d, want at least 1", warehouseCount)
	}

	supplierCount, err := store.SupplierPerformanceCount(ctx, reporting.Query{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("SupplierPerformanceCount() error = %v", err)
	}
	if supplierCount < 0 {
		t.Fatalf("SupplierPerformanceCount() = %d, want non-negative", supplierCount)
	}
}

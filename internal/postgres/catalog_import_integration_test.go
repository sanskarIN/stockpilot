package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestProductImportIntegrationIsAtomic(t *testing.T) {
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
	firstID := "prd_import_" + suffix + "_1"
	secondID := "prd_import_" + suffix + "_2"
	sku := "IMPORT-" + suffix
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, _ = store.pool.Exec(cleanupCtx, `DELETE FROM products WHERE id IN ($1, $2)`, firstID, secondID)
	}()

	valid := domain.Product{ID: firstID, SKU: sku, Name: "Atomic import first", Unit: "pcs", Currency: "INR", Active: true}
	duplicate := domain.Product{ID: secondID, SKU: sku, Name: "Atomic import duplicate", Unit: "pcs", Currency: "INR", Active: true}

	if err := store.ImportProducts(ctx, []domain.Product{valid, duplicate}); err == nil {
		t.Fatal("ImportProducts() error = nil, want duplicate conflict")
	}

	if _, err := store.GetProduct(ctx, firstID); err == nil {
		t.Fatal("first product was persisted after failed batch; transaction was not atomic")
	}
	if _, err := store.GetProduct(ctx, secondID); err == nil {
		t.Fatal("second product was persisted after failed batch; transaction was not atomic")
	}
}

package postgres

import (
	"context"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

// ImportProducts inserts the complete batch in one PostgreSQL transaction.
// Validation is intentionally repeated here because a dry-run request and the
// later write request may observe different database state.
func (s *Store) ImportProducts(ctx context.Context, products []domain.Product) error {
	if len(products) == 0 {
		return domain.ErrInvalid
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, product := range products {
		if err := product.Validate(); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO products (
				id, sku, name, description, category_id, supplier_id, barcode, unit,
				unit_cost_minor, currency, reorder_point, reorder_quantity, track_lots, track_expiry, active
			) VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), $8, $9, upper($10), $11, $12, $13, $14, $15)`,
			product.ID,
			strings.TrimSpace(product.SKU),
			strings.TrimSpace(product.Name),
			strings.TrimSpace(product.Description),
			product.CategoryID,
			product.SupplierID,
			strings.TrimSpace(product.Barcode),
			strings.TrimSpace(product.Unit),
			product.UnitCostMinor,
			strings.TrimSpace(product.Currency),
			product.ReorderPoint,
			product.ReorderQuantity,
			product.TrackLots,
			product.TrackExpiry,
			product.Active,
		); err != nil {
			return mapError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return mapError(err)
	}
	return nil
}

var _ repository.ProductBatchImporter = (*Store)(nil)

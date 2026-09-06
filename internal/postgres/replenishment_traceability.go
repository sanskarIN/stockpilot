package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) CreateReplenishmentReview(ctx context.Context, review domain.ReplenishmentReview) error {
	if err := review.Validate(); err != nil {
		return err
	}
	var purchaseOrderID any
	if id := strings.TrimSpace(review.PurchaseOrderID); id != "" {
		purchaseOrderID = id
		var exists bool
		if err := s.pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM purchase_order_lines
				WHERE purchase_order_id=$1 AND product_id=$2
			)`, id, review.ProductID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("%w: purchase order does not contain the reviewed product", domain.ErrInvalid)
		}
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO replenishment_reviews (
			id, product_id, sku, product_name, unit, supplier_id,
			on_hand, reorder_point, reorder_quantity, target_stock, suggested_quantity,
			reviewed_by, reviewed_at, purchase_order_id, outcome
		) VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		review.ID, review.ProductID, strings.TrimSpace(review.SKU), strings.TrimSpace(review.ProductName), strings.TrimSpace(review.Unit), strings.TrimSpace(review.SupplierID),
		review.OnHand, review.ReorderPoint, review.ReorderQuantity, review.TargetStock, review.SuggestedQuantity,
		review.ReviewedBy, review.ReviewedAt, purchaseOrderID, review.Outcome)
	return mapError(err)
}

func (s *Store) ListReplenishmentReviews(ctx context.Context, outcome string, limit, offset int) ([]domain.ReplenishmentReview, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	outcome = strings.TrimSpace(outcome)
	if outcome != "" && !domain.ReplenishmentReviewOutcome(outcome).Valid() {
		return nil, fmt.Errorf("%w: unsupported replenishment review outcome", domain.ErrInvalid)
	}
	args := []any{limit, offset}
	query := `
		SELECT id, product_id, sku, product_name, unit, COALESCE(supplier_id, ''),
			on_hand, reorder_point, reorder_quantity, target_stock, suggested_quantity,
			reviewed_by, reviewed_at, COALESCE(purchase_order_id, ''), outcome, created_at
		FROM replenishment_reviews`
	if outcome != "" {
		query += ` WHERE outcome=$3`
		args = append(args, outcome)
	}
	query += ` ORDER BY reviewed_at DESC, id DESC LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.ReplenishmentReview, 0)
	for rows.Next() {
		var item domain.ReplenishmentReview
		if err := rows.Scan(&item.ID, &item.ProductID, &item.SKU, &item.ProductName, &item.Unit, &item.SupplierID,
			&item.OnHand, &item.ReorderPoint, &item.ReorderQuantity, &item.TargetStock, &item.SuggestedQuantity,
			&item.ReviewedBy, &item.ReviewedAt, &item.PurchaseOrderID, &item.Outcome, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

var _ repository.ReplenishmentTraceability = (*Store)(nil)

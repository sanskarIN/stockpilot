package postgres

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (s *Store) ListBalances(ctx context.Context, limit, offset int) ([]domain.StockBalance, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 5000 {
		limit = 5000
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.pool.Query(ctx, `SELECT product_id, location_id, COALESCE(lot_id, ''), quantity, updated_at FROM inventory_balances ORDER BY product_id, location_id, lot_id NULLS FIRST LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.StockBalance, 0)
	for rows.Next() {
		var item domain.StockBalance
		if err := rows.Scan(&item.ProductID, &item.LocationID, &item.LotID, &item.Quantity, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

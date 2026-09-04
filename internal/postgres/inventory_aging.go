package postgres

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (s *Store) GetInventoryAging(ctx context.Context, limit int) (domain.InventoryAgingReport, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 5000 {
		limit = 5000
	}
	asOf := time.Now().UTC()
	rows, err := s.pool.Query(ctx, `
SELECT b.product_id,p.sku,p.name,b.location_id,COALESCE(b.lot_id,''),b.quantity,
       COALESCE(MAX(m.occurred_at), l.created_at) AS last_movement_at
FROM inventory_balances b
JOIN products p ON p.id=b.product_id
LEFT JOIN lots l ON l.id=b.lot_id
LEFT JOIN stock_movements m ON m.product_id=b.product_id
  AND m.location_id=b.location_id
  AND COALESCE(m.lot_id,'')=COALESCE(b.lot_id,'')
WHERE b.quantity > 0
GROUP BY b.product_id,p.sku,p.name,b.location_id,b.lot_id,b.quantity,l.created_at
ORDER BY last_movement_at ASC,p.name ASC,b.location_id ASC,COALESCE(b.lot_id,'') ASC
LIMIT $1`, limit)
	if err != nil {
		return domain.InventoryAgingReport{}, err
	}
	defer rows.Close()

	report := domain.InventoryAgingReport{Items: make([]domain.InventoryAgingItem, 0)}
	for rows.Next() {
		var item domain.InventoryAgingItem
		if err := rows.Scan(&item.ProductID, &item.SKU, &item.Name, &item.LocationID, &item.LotID, &item.Quantity, &item.LastMovementAt); err != nil {
			return domain.InventoryAgingReport{}, err
		}
		item.AsOf = asOf
		age := asOf.Sub(item.LastMovementAt).Hours() / 24
		item.AgeDays = int64(age)
		if item.AgeDays < 0 {
			item.AgeDays = 0
		}
		item.Bucket = domain.AgingBucket(item.AgeDays)
		report.Items = append(report.Items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.InventoryAgingReport{}, err
	}
	return report, nil
}

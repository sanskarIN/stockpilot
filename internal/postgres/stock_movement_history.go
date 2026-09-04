package postgres

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) GetStockMovementHistory(ctx context.Context, windowDays, limit int) (domain.StockMovementHistoryReport, error) {
	if windowDays <= 0 {
		windowDays = 30
	}
	if windowDays > 365 {
		windowDays = 365
	}
	if limit <= 0 {
		limit = 1000
	}
	if limit > 5000 {
		limit = 5000
	}
	asOf := time.Now().UTC()
	rows, err := s.pool.Query(ctx, `
SELECT m.product_id,p.sku,p.name,m.location_id,COALESCE(m.lot_id,''),
       COUNT(*) AS movement_count,
       COALESCE(SUM(CASE WHEN m.quantity_delta > 0 THEN m.quantity_delta ELSE 0 END),0) AS inbound_units,
       COALESCE(SUM(CASE WHEN m.quantity_delta < 0 THEN -m.quantity_delta ELSE 0 END),0) AS outbound_units,
       COALESCE(SUM(m.quantity_delta),0) AS net_units,
       MAX(m.occurred_at) AS last_movement_at
FROM stock_movements m
JOIN products p ON p.id=m.product_id
WHERE m.occurred_at >= $1
GROUP BY m.product_id,p.sku,p.name,m.location_id,m.lot_id
ORDER BY outbound_units DESC, last_movement_at DESC, p.name ASC, m.location_id ASC, COALESCE(m.lot_id,'') ASC
LIMIT $2`, asOf.Add(-time.Duration(windowDays)*24*time.Hour), limit)
	if err != nil {
		return domain.StockMovementHistoryReport{}, err
	}
	defer rows.Close()

	report := domain.StockMovementHistoryReport{AsOf: asOf, WindowDays: windowDays, Items: make([]domain.StockMovementHistoryItem, 0)}
	for rows.Next() {
		var item domain.StockMovementHistoryItem
		if err := rows.Scan(&item.ProductID, &item.SKU, &item.Name, &item.LocationID, &item.LotID, &item.MovementCount, &item.InboundUnits, &item.OutboundUnits, &item.NetUnits, &item.LastMovementAt); err != nil {
			return domain.StockMovementHistoryReport{}, err
		}
		item.AverageDailyOut = float64(item.OutboundUnits) / float64(windowDays)
		report.Items = append(report.Items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.StockMovementHistoryReport{}, err
	}
	return report, nil
}

var _ repository.StockMovementHistory = (*Store)(nil)

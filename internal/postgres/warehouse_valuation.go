package postgres

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) WarehouseValuation(ctx context.Context, limit int) (domain.WarehouseValuationReport, error) {
	if limit <= 0 {
		limit = 1000
	}
	if limit > 5000 {
		limit = 5000
	}

	rows, err := s.pool.Query(ctx, `
		SELECT w.id, w.code, w.name,
		       l.id, l.code, l.name,
		       p.currency,
		       COALESCE(sum(b.quantity), 0)::bigint AS on_hand,
		       COALESCE(sum(b.quantity * p.unit_cost_minor), 0)::bigint AS valuation_minor,
		       count(DISTINCT p.id)::bigint AS product_count
		FROM inventory_balances b
		JOIN products p ON p.id=b.product_id
		JOIN locations l ON l.id=b.location_id
		JOIN warehouses w ON w.id=l.warehouse_id
		WHERE b.quantity > 0 AND p.active=true
		GROUP BY w.id, w.code, w.name, l.id, l.code, l.name, p.currency
		ORDER BY valuation_minor DESC, w.name ASC, l.name ASC, p.currency ASC
		LIMIT $1`, limit)
	if err != nil {
		return domain.WarehouseValuationReport{}, err
	}
	defer rows.Close()

	items := make([]domain.WarehouseValuationItem, 0)
	for rows.Next() {
		var item domain.WarehouseValuationItem
		if err := rows.Scan(
			&item.WarehouseID, &item.WarehouseCode, &item.WarehouseName,
			&item.LocationID, &item.LocationCode, &item.LocationName,
			&item.Currency, &item.OnHand, &item.ValuationMinor, &item.ProductCount,
		); err != nil {
			return domain.WarehouseValuationReport{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.WarehouseValuationReport{}, err
	}

	totalRows, err := s.pool.Query(ctx, `
		SELECT w.id, w.code, w.name, p.currency,
		       COALESCE(sum(b.quantity), 0)::bigint,
		       COALESCE(sum(b.quantity * p.unit_cost_minor), 0)::bigint
		FROM inventory_balances b
		JOIN products p ON p.id=b.product_id
		JOIN locations l ON l.id=b.location_id
		JOIN warehouses w ON w.id=l.warehouse_id
		WHERE b.quantity > 0 AND p.active=true
		GROUP BY w.id, w.code, w.name, p.currency
		ORDER BY w.name ASC, p.currency ASC`)
	if err != nil {
		return domain.WarehouseValuationReport{}, err
	}
	defer totalRows.Close()

	totals := make([]domain.WarehouseValuationTotal, 0)
	for totalRows.Next() {
		var total domain.WarehouseValuationTotal
		if err := totalRows.Scan(&total.WarehouseID, &total.WarehouseCode, &total.WarehouseName, &total.Currency, &total.OnHand, &total.ValuationMinor); err != nil {
			return domain.WarehouseValuationReport{}, err
		}
		totals = append(totals, total)
	}
	if err := totalRows.Err(); err != nil {
		return domain.WarehouseValuationReport{}, err
	}
	return domain.WarehouseValuationReport{Items: items, Totals: totals}, nil
}

var _ repository.Reports = (*Store)(nil)

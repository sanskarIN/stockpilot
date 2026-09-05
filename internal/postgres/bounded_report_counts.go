package postgres

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/reporting"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) SupplierPerformanceCount(ctx context.Context, query reporting.Query) (int64, error) {
	if err := query.Validate(); err != nil {
		return 0, err
	}

	windowDays := 30
	if query.From != nil && query.To != nil {
		days := int(query.To.Sub(*query.From).Hours()/24) + 1
		if days > 0 {
			windowDays = days
		}
	}

	var count int64
	err := s.pool.QueryRow(ctx, `
		WITH selected_orders AS (
			SELECT o.id, o.supplier_id
			FROM purchase_orders o
			WHERE o.created_at >= now() - ($1::int * interval '1 day')
		), order_totals AS (
			SELECT l.purchase_order_id
			FROM purchase_order_lines l
			JOIN selected_orders o ON o.id=l.purchase_order_id
			GROUP BY l.purchase_order_id
		)
		SELECT count(*)
		FROM (
			SELECT s.id
			FROM suppliers s
			JOIN selected_orders o ON o.supplier_id=s.id
			JOIN order_totals t ON t.purchase_order_id=o.id
			GROUP BY s.id
		) supplier_rows`, windowDays).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) WarehouseValuationCount(ctx context.Context, query reporting.Query) (int64, error) {
	if err := query.Validate(); err != nil {
		return 0, err
	}

	var count int64
	err := s.pool.QueryRow(ctx, `
		SELECT count(*)
		FROM (
			SELECT w.id, l.id, p.currency
			FROM inventory_balances b
			JOIN products p ON p.id=b.product_id
			JOIN locations l ON l.id=b.location_id
			JOIN warehouses w ON w.id=l.warehouse_id
			WHERE b.quantity > 0 AND p.active=true
			GROUP BY w.id, w.code, w.name, l.id, l.code, l.name, p.currency
		) warehouse_rows`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

var _ repository.CountedReports = (*Store)(nil)

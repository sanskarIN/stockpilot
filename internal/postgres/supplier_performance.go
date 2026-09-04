package postgres

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) SupplierPerformance(ctx context.Context, windowDays, limit int) (domain.SupplierPerformanceReport, error) {
	if windowDays <= 0 {
		windowDays = 30
	}
	if windowDays > 365 {
		windowDays = 365
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 5000 {
		limit = 5000
	}

	rows, err := s.pool.Query(ctx, `
		WITH selected_orders AS (
			SELECT o.id, o.supplier_id, o.status, o.expected_at, o.created_at
			FROM purchase_orders o
			WHERE o.created_at >= now() - ($1::int * interval '1 day')
			ORDER BY o.created_at DESC, o.id DESC
		), order_totals AS (
			SELECT l.purchase_order_id,
			       COALESCE(sum(l.quantity), 0)::bigint AS ordered_units,
			       COALESCE(sum(l.received), 0)::bigint AS received_units,
			       COALESCE(sum(l.quantity * l.unit_cost_minor), 0)::bigint AS ordered_value_minor,
			       COALESCE(sum(l.received * l.unit_cost_minor), 0)::bigint AS received_value_minor
			FROM purchase_order_lines l
			JOIN selected_orders o ON o.id=l.purchase_order_id
			GROUP BY l.purchase_order_id
		), receipt_times AS (
			SELECT o.id AS order_id,
			       min(m.occurred_at) AS first_received_at
			FROM selected_orders o
			JOIN stock_movements m ON m.reference='PO:' || o.id
			WHERE m.movement_type='receive'
			GROUP BY o.id
		), supplier_rows AS (
			SELECT s.id, s.code, s.name,
			       count(o.id)::bigint AS order_count,
			       COALESCE(sum(t.ordered_units),0)::bigint AS ordered_units,
			       COALESCE(sum(t.received_units),0)::bigint AS received_units,
			       COALESCE(sum(t.ordered_units-t.received_units),0)::bigint AS open_units,
			       COALESCE(sum(t.ordered_value_minor),0)::bigint AS ordered_value_minor,
			       COALESCE(sum(t.received_value_minor),0)::bigint AS received_value_minor,
			       COALESCE(avg(EXTRACT(epoch FROM (r.first_received_at-o.created_at))/86400.0) FILTER (WHERE r.first_received_at IS NOT NULL),0)::double precision AS average_lead_time_days,
			       count(o.id) FILTER (WHERE r.first_received_at IS NOT NULL)::bigint AS completed_order_count,
			       count(o.id) FILTER (WHERE r.first_received_at IS NOT NULL AND (o.expected_at IS NULL OR r.first_received_at <= o.expected_at))::bigint AS on_time_order_count
			FROM suppliers s
			JOIN selected_orders o ON o.supplier_id=s.id
			JOIN order_totals t ON t.purchase_order_id=o.id
			LEFT JOIN receipt_times r ON r.order_id=o.id
			GROUP BY s.id, s.code, s.name
		)
		SELECT id, code, name, order_count, ordered_units, received_units, open_units,
		       ordered_value_minor, received_value_minor, average_lead_time_days,
		       completed_order_count, on_time_order_count
		FROM supplier_rows
		ORDER BY order_count DESC, name ASC, id ASC
		LIMIT $2`, windowDays, limit)
	if err != nil {
		return domain.SupplierPerformanceReport{}, err
	}
	defer rows.Close()

	items := make([]domain.SupplierPerformanceItem, 0)
	for rows.Next() {
		var item domain.SupplierPerformanceItem
		if err := rows.Scan(
			&item.SupplierID, &item.SupplierCode, &item.SupplierName,
			&item.OrderCount, &item.OrderedUnits, &item.ReceivedUnits, &item.OpenUnits,
			&item.OrderedValueMinor, &item.ReceivedValueMinor, &item.AverageLeadTimeDays,
			&item.CompletedOrderCount, &item.OnTimeOrderCount,
		); err != nil {
			return domain.SupplierPerformanceReport{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.SupplierPerformanceReport{}, err
	}
	return domain.SupplierPerformanceReport{AsOf: time.Now().UTC(), WindowDays: windowDays, Items: items}, nil
}

var _ repository.Reports = (*Store)(nil)

package postgres

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) ReplenishmentPerformance(ctx context.Context, query reporting.Query) (domain.ReplenishmentPerformanceReport, error) {
	if err := query.Validate(); err != nil {
		return domain.ReplenishmentPerformanceReport{}, err
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		WITH selected_orders AS (
			SELECT o.id, o.supplier_id, o.expected_at, o.created_at
			FROM purchase_orders o
			WHERE ($1::timestamptz IS NULL OR o.created_at >= $1::timestamptz)
				AND ($2::timestamptz IS NULL OR o.created_at < ($2::timestamptz + interval '1 day'))
				AND (($1::timestamptz IS NOT NULL AND $2::timestamptz IS NOT NULL)
					OR o.created_at >= now() - ($3::int * interval '1 day'))
				AND o.status <> 'cancelled'
		), line_totals AS (
			SELECT l.purchase_order_id,
				sum(l.quantity)::bigint AS ordered_units,
				sum(l.received)::bigint AS received_units
			FROM purchase_order_lines l
			JOIN selected_orders o ON o.id=l.purchase_order_id
			GROUP BY l.purchase_order_id
		), receipt_stats AS (
			SELECT m.reference, min(m.occurred_at) AS first_received_at
			FROM stock_movements m
			WHERE m.movement_type='receive'
			GROUP BY m.reference
		)
		SELECT s.id, s.name,
			count(*)::bigint AS order_count,
			coalesce(sum(t.ordered_units), 0)::bigint AS ordered_units,
			coalesce(sum(t.received_units), 0)::bigint AS received_units,
			coalesce(sum(t.ordered_units - t.received_units), 0)::bigint AS outstanding_units,
			coalesce(sum(t.received_units)::numeric / nullif(sum(t.ordered_units),0), 0)::double precision AS fill_rate,
			count(*) FILTER (WHERE r.first_received_at IS NOT NULL AND (o.expected_at IS NULL OR r.first_received_at <= o.expected_at))::bigint AS on_time_order_count,
			count(*) FILTER (WHERE r.first_received_at IS NOT NULL AND o.expected_at IS NOT NULL AND r.first_received_at > o.expected_at)::bigint AS late_order_count,
			coalesce(avg(EXTRACT(EPOCH FROM (r.first_received_at - o.created_at)) / 86400.0) FILTER (WHERE r.first_received_at IS NOT NULL), 0)::double precision AS average_lead_days
		FROM suppliers s
		JOIN selected_orders o ON o.supplier_id=s.id
		JOIN line_totals t ON t.purchase_order_id=o.id
		LEFT JOIN receipt_stats r ON r.reference='PO:' || o.id
		GROUP BY s.id, s.name
		ORDER BY outstanding_units DESC, s.name ASC, s.id ASC
		LIMIT $4 OFFSET $5`, query.From, query.To, reporting.DefaultPeriodDays, limit, query.Offset)
	if err != nil {
		return domain.ReplenishmentPerformanceReport{}, err
	}
	defer rows.Close()
	items := make([]domain.ReplenishmentPerformanceItem, 0)
	for rows.Next() {
		var item domain.ReplenishmentPerformanceItem
		if err := rows.Scan(&item.SupplierID, &item.SupplierName, &item.OrderCount, &item.OrderedUnits, &item.ReceivedUnits,
			&item.OutstandingUnits, &item.FillRate, &item.OnTimeOrderCount, &item.LateOrderCount, &item.AverageLeadDays); err != nil {
			return domain.ReplenishmentPerformanceReport{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.ReplenishmentPerformanceReport{}, err
	}
	windowDays := reporting.DefaultPeriodDays
	if query.From != nil && query.To != nil {
		windowDays = int(query.To.Sub(*query.From).Hours()/24) + 1
	}
	return domain.ReplenishmentPerformanceReport{AsOf: time.Now().UTC(), WindowDays: windowDays, Items: items}, nil
}

var _ repository.ReplenishmentReports = (*Store)(nil)

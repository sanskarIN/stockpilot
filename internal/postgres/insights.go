package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) InventorySummary(ctx context.Context) (domain.InventorySummary, error) {
	var summary domain.InventorySummary
	err := s.pool.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM products),
			(SELECT count(*) FROM products WHERE active=true),
			(SELECT count(*) FROM warehouses WHERE active=true),
			(SELECT count(*) FROM locations WHERE active=true),
			COALESCE((SELECT sum(quantity) FROM inventory_balances), 0),
			(
				SELECT count(*)
				FROM inventory_balances b
				JOIN products p ON p.id=b.product_id
				WHERE p.active=true AND b.quantity <= p.reorder_point
			),
			(
				SELECT count(*)
				FROM products p
				WHERE p.active=true
				  AND COALESCE((SELECT sum(b.quantity) FROM inventory_balances b WHERE b.product_id=p.id), 0) = 0
			)`).Scan(
		&summary.ProductCount,
		&summary.ActiveProductCount,
		&summary.ActiveWarehouseCount,
		&summary.ActiveLocationCount,
		&summary.TotalUnits,
		&summary.LowStockBalanceCount,
		&summary.OutOfStockCount,
	)
	if err != nil {
		return domain.InventorySummary{}, err
	}
	return summary, nil
}

func (s *Store) PurchasingSummary(ctx context.Context) (domain.PurchasingSummary, error) {
	var summary domain.PurchasingSummary
	err := s.pool.QueryRow(ctx, `
		SELECT
			count(*),
			count(*) FILTER (WHERE status='draft'),
			count(*) FILTER (WHERE status='ordered'),
			count(*) FILTER (WHERE status='partially_received'),
			count(*) FILTER (WHERE status='received'),
			count(*) FILTER (WHERE status='cancelled'),
			COALESCE((
				SELECT sum(l.quantity-l.received)
				FROM purchase_order_lines l
				JOIN purchase_orders o ON o.id=l.purchase_order_id
				WHERE o.status IN ('ordered', 'partially_received')
			), 0)
		FROM purchase_orders`).Scan(
		&summary.TotalOrders,
		&summary.DraftOrders,
		&summary.OrderedOrders,
		&summary.PartiallyReceivedOrders,
		&summary.ReceivedOrders,
		&summary.CancelledOrders,
		&summary.OutstandingUnits,
	)
	if err != nil {
		return domain.PurchasingSummary{}, err
	}
	return summary, nil
}

func (s *Store) AppendAuditEvent(ctx context.Context, event domain.AuditEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		return fmt.Errorf("encode audit metadata: %w", err)
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO audit_log (occurred_at, actor_id, action, entity_type, entity_id, request_id, metadata)
		VALUES (COALESCE($1, now()), $2, $3, $4, $5, $6, $7::jsonb)`,
		event.OccurredAt,
		strings.TrimSpace(event.ActorID),
		strings.TrimSpace(event.Action),
		strings.TrimSpace(event.EntityType),
		strings.TrimSpace(event.EntityID),
		strings.TrimSpace(event.RequestID),
		string(metadata),
	)
	return err
}

func (s *Store) ListAuditEvents(ctx context.Context, filter domain.AuditFilter) ([]domain.AuditEvent, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	args := make([]any, 0, 6)
	where := make([]string, 0, 4)
	addFilter := func(column, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		args = append(args, value)
		where = append(where, fmt.Sprintf("%s=$%d", column, len(args)))
	}
	addFilter("actor_id", filter.ActorID)
	addFilter("action", filter.Action)
	addFilter("entity_type", filter.EntityType)
	addFilter("entity_id", filter.EntityID)
	query := `SELECT id, occurred_at, actor_id, action, entity_type, entity_id, request_id, metadata FROM audit_log`
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, ` AND `)
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(` ORDER BY occurred_at DESC, id DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args))
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.AuditEvent, 0)
	for rows.Next() {
		var item domain.AuditEvent
		var metadata []byte
		if err := rows.Scan(&item.ID, &item.OccurredAt, &item.ActorID, &item.Action, &item.EntityType, &item.EntityID, &item.RequestID, &metadata); err != nil {
			return nil, err
		}
		item.Metadata = make(map[string]any)
		if len(metadata) > 0 {
			if err := json.Unmarshal(metadata, &item.Metadata); err != nil {
				return nil, fmt.Errorf("decode audit metadata: %w", err)
			}
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

var _ repository.Reports = (*Store)(nil)
var _ repository.Audit = (*Store)(nil)

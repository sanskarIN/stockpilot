package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const receiptHistoryRepositoryMax = 5000

func (s *Store) ListReceiptHistory(ctx context.Context, filter repository.ReceiptHistoryFilter) ([]domain.ReceiptHistoryRow, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 500
	}
	if limit > receiptHistoryRepositoryMax {
		limit = receiptHistoryRepositoryMax
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	where := []string{"m.movement_type = 'receive'"}
	args := make([]any, 0, 10)
	add := func(value any, clause string) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(clause, len(args)))
	}
	if value := strings.TrimSpace(filter.ProductID); value != "" {
		add(value, "m.product_id = $%d")
	}
	if value := strings.TrimSpace(filter.WarehouseID); value != "" {
		add(value, "l.warehouse_id = $%d")
	}
	if value := strings.TrimSpace(filter.LocationID); value != "" {
		add(value, "m.location_id = $%d")
	}
	if value := strings.TrimSpace(filter.LotID); value != "" {
		add(value, "m.lot_id = $%d")
	}
	if value := strings.TrimSpace(filter.ActorID); value != "" {
		add(value, "m.actor_id = $%d")
	}
	if value := strings.TrimSpace(filter.Reference); value != "" {
		add(value, "m.reference = $%d")
	}
	if filter.From != nil {
		add(*filter.From, "m.occurred_at >= $%d")
	}
	if filter.To != nil {
		add(*filter.To, "m.occurred_at < $%d")
	}

	limitArg := len(args) + 1
	offsetArg := len(args) + 2
	args = append(args, limit, filter.Offset)
	query := fmt.Sprintf(`
		SELECT
			m.id,
			m.product_id,
			p.sku,
			p.name,
			m.location_id,
			l.name,
			l.warehouse_id,
			w.name,
			COALESCE(m.lot_id, ''),
			COALESCE(lo.lot_number, ''),
			m.quantity_delta,
			m.reference,
			m.note,
			m.actor_id,
			m.occurred_at,
			m.created_at
		FROM stock_movements m
		JOIN products p ON p.id = m.product_id
		JOIN locations l ON l.id = m.location_id
		JOIN warehouses w ON w.id = l.warehouse_id
		LEFT JOIN lots lo ON lo.id = m.lot_id
		WHERE %s
		ORDER BY m.occurred_at DESC, m.id DESC
		LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), limitArg, offsetArg)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ReceiptHistoryRow, 0)
	for rows.Next() {
		var item domain.ReceiptHistoryRow
		if err := rows.Scan(
			&item.MovementID,
			&item.ProductID,
			&item.SKU,
			&item.ProductName,
			&item.LocationID,
			&item.Location,
			&item.WarehouseID,
			&item.Warehouse,
			&item.LotID,
			&item.LotNumber,
			&item.Quantity,
			&item.Reference,
			&item.Note,
			&item.ActorID,
			&item.OccurredAt,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

package postgres

import (
	"context"
	"fmt"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (s *Store) ListOrderExportRows(ctx context.Context, status domain.PurchaseOrderStatus, limit, offset int) ([]domain.PurchaseOrderExportRow, error) {
	if limit <= 0 {
		limit = 500
	}
	if limit > 5000 {
		limit = 5000
	}
	if offset < 0 {
		offset = 0
	}
	if status != "" && !status.Valid() {
		return nil, fmt.Errorf("%w: unsupported purchase order status", domain.ErrInvalid)
	}

	rows, err := s.pool.Query(ctx, `
		WITH selected_orders AS (
			SELECT id, order_number, supplier_id, warehouse_id, status, currency, expected_at,
			       notes, created_by, created_at, updated_at
			FROM purchase_orders
			WHERE ($1 = '' OR status = $1)
			ORDER BY created_at DESC, id DESC
			LIMIT $2 OFFSET $3
		)
		SELECT o.id, o.order_number, o.supplier_id, o.warehouse_id, o.status, o.currency,
		       o.expected_at, o.notes, o.created_by, o.created_at, o.updated_at,
		       l.id, l.product_id, l.quantity, l.received, l.unit_cost_minor
		FROM selected_orders o
		JOIN purchase_order_lines l ON l.purchase_order_id = o.id
		ORDER BY o.created_at DESC, o.id DESC, l.id`, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.PurchaseOrderExportRow, 0)
	for rows.Next() {
		var item domain.PurchaseOrderExportRow
		if err := rows.Scan(
			&item.OrderID, &item.OrderNumber, &item.SupplierID, &item.WarehouseID, &item.Status, &item.Currency,
			&item.ExpectedAt, &item.Notes, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt,
			&item.LineID, &item.ProductID, &item.Quantity, &item.Received, &item.UnitCostMinor,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

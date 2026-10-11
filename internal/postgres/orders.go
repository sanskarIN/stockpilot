package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) CreateOrder(ctx context.Context, order domain.PurchaseOrder) error {
	if err := order.Validate(); err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `
		INSERT INTO purchase_orders (id, order_number, supplier_id, warehouse_id, status, currency, expected_at, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, upper($6), $7, $8, $9)`,
		order.ID, strings.TrimSpace(order.Number), order.SupplierID, order.WarehouseID, order.Status, strings.TrimSpace(order.Currency),
		order.ExpectedAt, strings.TrimSpace(order.Notes), strings.TrimSpace(order.CreatedBy))
	if err != nil {
		return mapError(err)
	}
	for _, line := range order.Lines {
		_, err = tx.Exec(ctx, `
			INSERT INTO purchase_order_lines (id, purchase_order_id, product_id, quantity, received, unit_cost_minor)
			VALUES ($1, $2, $3, $4, $5, $6)`, line.ID, order.ID, line.ProductID, line.Quantity, line.Received, line.UnitCostMinor)
		if err != nil {
			return mapError(err)
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) GetOrder(ctx context.Context, id string) (domain.PurchaseOrder, error) {
	var order domain.PurchaseOrder
	err := s.pool.QueryRow(ctx, `
		SELECT id, order_number, supplier_id, warehouse_id, status, currency, expected_at, notes, created_by, created_at, updated_at
		FROM purchase_orders WHERE id=$1`, id).Scan(
		&order.ID, &order.Number, &order.SupplierID, &order.WarehouseID, &order.Status, &order.Currency, &order.ExpectedAt,
		&order.Notes, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return domain.PurchaseOrder{}, mapError(err)
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, purchase_order_id, product_id, quantity, received, unit_cost_minor
		FROM purchase_order_lines WHERE purchase_order_id=$1 ORDER BY id`, id)
	if err != nil {
		return domain.PurchaseOrder{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var line domain.PurchaseOrderLine
		if err := rows.Scan(&line.ID, &line.PurchaseOrderID, &line.ProductID, &line.Quantity, &line.Received, &line.UnitCostMinor); err != nil {
			return domain.PurchaseOrder{}, err
		}
		order.Lines = append(order.Lines, line)
	}
	if err := rows.Err(); err != nil {
		return domain.PurchaseOrder{}, err
	}
	return order, nil
}

func (s *Store) ListOrders(ctx context.Context, status domain.PurchaseOrderStatus, limit, offset int) ([]domain.PurchaseOrder, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	args := make([]any, 0, 3)
	query := `SELECT id, order_number, supplier_id, warehouse_id, status, currency, expected_at, notes, created_by, created_at, updated_at FROM purchase_orders`
	if status != "" {
		if !status.Valid() {
			return nil, fmt.Errorf("%w: unsupported purchase order status", domain.ErrInvalid)
		}
		args = append(args, status)
		query += ` WHERE status=$1`
	}
	args = append(args, limit, offset)
	query += fmt.Sprintf(` ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, len(args)-1, len(args))
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.PurchaseOrder, 0)
	for rows.Next() {
		var item domain.PurchaseOrder
		if err := rows.Scan(&item.ID, &item.Number, &item.SupplierID, &item.WarehouseID, &item.Status, &item.Currency, &item.ExpectedAt,
			&item.Notes, &item.CreatedBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) ReceiveLine(ctx context.Context, orderID, lineID string, quantity int64, locationID, lotID, actorID string) error {
	if strings.TrimSpace(orderID) == "" || strings.TrimSpace(lineID) == "" || strings.TrimSpace(locationID) == "" || quantity <= 0 {
		return fmt.Errorf("%w: order, line, location, and positive quantity are required", domain.ErrInvalid)
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var productID, warehouseID string
	var ordered, received int64
	var status domain.PurchaseOrderStatus
	err = tx.QueryRow(ctx, `
		SELECT l.product_id, l.quantity, l.received, o.status, o.warehouse_id
		FROM purchase_order_lines l
		JOIN purchase_orders o ON o.id=l.purchase_order_id
		WHERE l.id=$1 AND l.purchase_order_id=$2
		FOR UPDATE OF l, o`, lineID, orderID).Scan(&productID, &ordered, &received, &status, &warehouseID)
	if err != nil {
		return mapError(err)
	}
	if status == domain.PurchaseOrderDraft || status == domain.PurchaseOrderCancelled || status == domain.PurchaseOrderReceived {
		return fmt.Errorf("%w: order cannot receive stock in status %s", domain.ErrConflict, status)
	}
	if received+quantity > ordered {
		return fmt.Errorf("%w: receipt exceeds remaining order quantity", domain.ErrInvalid)
	}
	var locationWarehouse string
	if err := tx.QueryRow(ctx, `SELECT warehouse_id FROM locations WHERE id=$1 AND active=true`, locationID).Scan(&locationWarehouse); err != nil {
		return mapError(err)
	}
	if locationWarehouse != warehouseID {
		return fmt.Errorf("%w: receiving location belongs to a different warehouse", domain.ErrInvalid)
	}

	movementID, err := idgen.New("mov")
	if err != nil {
		return err
	}
	movement := domain.StockMovement{
		ID: movementID, ProductID: productID, LocationID: locationID, LotID: lotID, Type: domain.MovementReceive,
		QuantityDelta: quantity, Reference: "PO:" + orderID, ActorID: strings.TrimSpace(actorID), OccurredAt: time.Now().UTC(),
	}
	if _, err := applyMovementTx(ctx, tx, movement); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE purchase_order_lines SET received=received+$2 WHERE id=$1`, lineID, quantity); err != nil {
		return err
	}

	var allReceived, anyReceived bool
	if err := tx.QueryRow(ctx, `
		SELECT bool_and(received >= quantity), bool_or(received > 0)
		FROM purchase_order_lines WHERE purchase_order_id=$1`, orderID).Scan(&allReceived, &anyReceived); err != nil {
		return err
	}
	newStatus := domain.PurchaseOrderOrdered
	if allReceived {
		newStatus = domain.PurchaseOrderReceived
	} else if anyReceived {
		newStatus = domain.PurchaseOrderPartiallyReceived
	}
	if _, err := tx.Exec(ctx, `UPDATE purchase_orders SET status=$2, updated_at=now() WHERE id=$1`, orderID, newStatus); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

var _ repository.Orders = (*Store)(nil)

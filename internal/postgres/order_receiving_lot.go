package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
)

// ReceiveLineWithNewLot creates the requested lot and receives the order line
// inside one transaction. A failed receipt therefore rolls back the lot too.
func (s *Store) ReceiveLineWithNewLot(ctx context.Context, orderID, lineID string, quantity int64, locationID string, lot domain.Lot, actorID string) error {
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
	if err := tx.QueryRow(ctx, `
		SELECT l.product_id, l.quantity, l.received, o.status, o.warehouse_id
		FROM purchase_order_lines l
		JOIN purchase_orders o ON o.id=l.purchase_order_id
		WHERE l.id=$1 AND l.purchase_order_id=$2
		FOR UPDATE OF l, o`, lineID, orderID).Scan(&productID, &ordered, &received, &status, &warehouseID); err != nil {
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

	lot.ProductID = productID
	lot.ID, err = idgen.New("lot")
	if err != nil {
		return err
	}
	if err := lot.Validate(); err != nil {
		return err
	}

	var trackLots bool
	if err := tx.QueryRow(ctx, `SELECT track_lots FROM products WHERE id=$1`, productID).Scan(&trackLots); err != nil {
		return mapError(err)
	}
	if !trackLots {
		return fmt.Errorf("%w: product does not use lot tracking", domain.ErrInvalid)
	}

	var existingLotID string
	err = tx.QueryRow(ctx, `SELECT id FROM lots WHERE product_id=$1 AND lot_number=$2`, productID, strings.TrimSpace(lot.LotNumber)).Scan(&existingLotID)
	if err == nil {
		return fmt.Errorf("%w: lot number already exists for this product", domain.ErrConflict)
	}
	if err != pgx.ErrNoRows {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO lots (id, product_id, lot_number, manufactured_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)`, lot.ID, productID, strings.TrimSpace(lot.LotNumber), lot.Manufactured, lot.ExpiresAt); err != nil {
		return mapError(err)
	}

	movementID, err := idgen.New("mov")
	if err != nil {
		return err
	}
	movement := domain.StockMovement{
		ID: movementID, ProductID: productID, LocationID: locationID, LotID: lot.ID,
		Type: domain.MovementReceive, QuantityDelta: quantity, Reference: "PO:" + orderID,
		ActorID: strings.TrimSpace(actorID), OccurredAt: time.Now().UTC(),
	}
	if _, err := applyMovementTx(ctx, tx, movement); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE purchase_order_lines SET received=received+$2 WHERE id=$1`, lineID, quantity); err != nil {
		return err
	}

	var allReceived, anyReceived bool
	if err := tx.QueryRow(ctx, `SELECT bool_and(received >= quantity), bool_or(received > 0) FROM purchase_order_lines WHERE purchase_order_id=$1`, orderID).Scan(&allReceived, &anyReceived); err != nil {
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

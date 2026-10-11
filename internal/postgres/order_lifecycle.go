package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (s *Store) UpdateOrderStatus(ctx context.Context, orderID string, target domain.PurchaseOrderStatus, _ string) error {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" || !target.Valid() {
		return fmt.Errorf("%w: order id and valid target status are required", domain.ErrInvalid)
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var current domain.PurchaseOrderStatus
	if err := tx.QueryRow(ctx, `SELECT status FROM purchase_orders WHERE id=$1 FOR UPDATE`, orderID).Scan(&current); err != nil {
		return mapError(err)
	}
	if err := domain.ValidatePurchaseOrderTransition(current, target); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE purchase_orders SET status=$2, updated_at=now() WHERE id=$1`, orderID, target); err != nil {
		return mapError(err)
	}
	return tx.Commit(ctx)
}

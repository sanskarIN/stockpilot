package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) UpdateOrder(ctx context.Context, order domain.PurchaseOrder) error {
	if err := order.Validate(); err != nil { return err }
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil { return err }
	defer func(){ _ = tx.Rollback(ctx) }()
	var status domain.PurchaseOrderStatus
	var createdBy string
	if err := tx.QueryRow(ctx, `SELECT status, created_by FROM purchase_orders WHERE id=$1 FOR UPDATE`, order.ID).Scan(&status,&createdBy); err != nil { return mapError(err) }
	if status != domain.PurchaseOrderDraft { return fmt.Errorf("%w: only draft purchase orders can be edited", domain.ErrConflict) }
	for _, line := range order.Lines { if line.Received != 0 { return fmt.Errorf("%w: draft order lines must have zero received quantity", domain.ErrInvalid) } }
	if _, err := tx.Exec(ctx, `UPDATE purchase_orders SET order_number=$2,supplier_id=$3,warehouse_id=$4,currency=upper($5),expected_at=$6,notes=$7,updated_at=now() WHERE id=$1`, order.ID, strings.TrimSpace(order.Number), order.SupplierID, order.WarehouseID, strings.TrimSpace(order.Currency), order.ExpectedAt, strings.TrimSpace(order.Notes)); err != nil { return mapError(err) }
	if _, err := tx.Exec(ctx, `DELETE FROM purchase_order_lines WHERE purchase_order_id=$1`, order.ID); err != nil { return err }
	for _, line := range order.Lines { if _, err := tx.Exec(ctx, `INSERT INTO purchase_order_lines (id,purchase_order_id,product_id,quantity,received,unit_cost_minor) VALUES ($1,$2,$3,$4,0,$5)`,line.ID,order.ID,line.ProductID,line.Quantity,line.UnitCostMinor);err!=nil{return mapError(err)} }
	return tx.Commit(ctx)
}

var _ repository.Orders = (*Store)(nil)

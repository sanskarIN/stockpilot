package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) UpdateWarehouse(ctx context.Context, item domain.Warehouse) error {
	if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("%w: warehouse id, code, and name are required", domain.ErrInvalid)
	}
	if strings.TrimSpace(item.Timezone) == "" {
		item.Timezone = "UTC"
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil { return err }
	defer func(){ _ = tx.Rollback(ctx) }()
	var active bool
	if err := tx.QueryRow(ctx, `SELECT active FROM warehouses WHERE id=$1 FOR UPDATE`, item.ID).Scan(&active); err != nil { return mapError(err) }
	if active && !item.Active {
		var count int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM locations WHERE warehouse_id=$1 AND active=true`, item.ID).Scan(&count); err != nil { return err }
		if count > 0 { return fmt.Errorf("%w: warehouse has %d active locations", domain.ErrConflict, count) }
	}
	if _, err := tx.Exec(ctx, `UPDATE warehouses SET code=$2,name=$3,timezone=$4,active=$5 WHERE id=$1`, item.ID, strings.TrimSpace(item.Code), strings.TrimSpace(item.Name), strings.TrimSpace(item.Timezone), item.Active); err != nil { return mapError(err) }
	return tx.Commit(ctx)
}

func (s *Store) UpdateLocation(ctx context.Context, item domain.Location) error {
	if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.WarehouseID) == "" || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("%w: location id, warehouse, code, and name are required", domain.ErrInvalid)
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil { return err }
	defer func(){ _ = tx.Rollback(ctx) }()
	var currentWarehouse string
	var active bool
	if err := tx.QueryRow(ctx, `SELECT warehouse_id, active FROM locations WHERE id=$1 FOR UPDATE`, item.ID).Scan(&currentWarehouse, &active); err != nil { return mapError(err) }
	if !item.Active && active {
		var onHand int64
		if err := tx.QueryRow(ctx, `SELECT COALESCE(sum(quantity),0) FROM inventory_balances WHERE location_id=$1`, item.ID).Scan(&onHand); err != nil { return err }
		if onHand != 0 { return fmt.Errorf("%w: location has non-zero inventory", domain.ErrConflict) }
	}
	var warehouseActive bool
	if err := tx.QueryRow(ctx, `SELECT active FROM warehouses WHERE id=$1`, item.WarehouseID).Scan(&warehouseActive); err != nil { return mapError(err) }
	if !warehouseActive && item.Active { return fmt.Errorf("%w: location cannot be active in an archived warehouse", domain.ErrConflict) }
	if _, err := tx.Exec(ctx, `UPDATE locations SET warehouse_id=$2,code=$3,name=$4,active=$5 WHERE id=$1`, item.ID, item.WarehouseID, strings.TrimSpace(item.Code), strings.TrimSpace(item.Name), item.Active); err != nil { return mapError(err) }
	_ = currentWarehouse
	return tx.Commit(ctx)
}

var _ repository.Inventory = (*Store)(nil)

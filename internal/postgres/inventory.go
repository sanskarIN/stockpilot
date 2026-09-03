package postgres

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) CreateWarehouse(ctx context.Context, warehouse domain.Warehouse) error {
	if err := warehouse.Validate(); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO warehouses (id, code, name, address, timezone, active) VALUES ($1, $2, $3, $4, $5, $6)`, warehouse.ID, strings.TrimSpace(warehouse.Code), strings.TrimSpace(warehouse.Name), strings.TrimSpace(warehouse.Address), strings.TrimSpace(warehouse.Timezone), warehouse.Active)
	return mapError(err)
}
func (s *Store) ListWarehouses(ctx context.Context, activeOnly bool) ([]domain.Warehouse, error) {
	query := `SELECT id, code, name, address, timezone, active, created_at, updated_at FROM warehouses`
	if activeOnly {
		query += ` WHERE active = true`
	}
	query += ` ORDER BY name`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Warehouse, 0)
	for rows.Next() {
		var item domain.Warehouse
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Address, &item.Timezone, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (s *Store) CreateLocation(ctx context.Context, location domain.Location) error {
	if err := location.Validate(); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO locations (id, warehouse_id, code, name, description, active) VALUES ($1, $2, $3, $4, $5, $6)`, location.ID, location.WarehouseID, strings.TrimSpace(location.Code), strings.TrimSpace(location.Name), strings.TrimSpace(location.Description), location.Active)
	return mapError(err)
}
func (s *Store) ListLocations(ctx context.Context, warehouseID string, activeOnly bool) ([]domain.Location, error) {
	args := make([]any, 0, 1)
	where := make([]string, 0, 2)
	if warehouseID != "" {
		args = append(args, warehouseID)
		where = append(where, fmt.Sprintf("warehouse_id=$%d", len(args)))
	}
	if activeOnly {
		where = append(where, "active=true")
	}
	query := `SELECT id, warehouse_id, code, name, description, active, created_at, updated_at FROM locations`
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, ` AND `)
	}
	query += ` ORDER BY name`
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Location, 0)
	for rows.Next() {
		var item domain.Location
		if err := rows.Scan(&item.ID, &item.WarehouseID, &item.Code, &item.Name, &item.Description, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (s *Store) CreateLot(ctx context.Context, lot domain.Lot) error {
	if err := lot.Validate(); err != nil {
		return err
	}
	var trackLots bool
	if err := s.pool.QueryRow(ctx, `SELECT track_lots FROM products WHERE id=$1`, lot.ProductID).Scan(&trackLots); err != nil {
		return mapError(err)
	}
	if !trackLots {
		return fmt.Errorf("%w: product does not use lot tracking", domain.ErrInvalid)
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO lots (id, product_id, lot_number, manufactured_at, expires_at) VALUES ($1, $2, $3, $4, $5)`, lot.ID, lot.ProductID, strings.TrimSpace(lot.LotNumber), lot.Manufactured, lot.ExpiresAt)
	return mapError(err)
}
func (s *Store) ListLots(ctx context.Context, productID string, limit int) ([]domain.Lot, error) {
	if strings.TrimSpace(productID) == "" {
		return nil, fmt.Errorf("%w: product id is required", domain.ErrInvalid)
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	rows, err := s.pool.Query(ctx, `SELECT id, product_id, lot_number, manufactured_at, expires_at, created_at FROM lots WHERE product_id=$1 ORDER BY expires_at NULLS LAST, lot_number LIMIT $2`, strings.TrimSpace(productID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Lot, 0)
	for rows.Next() {
		var item domain.Lot
		if err := rows.Scan(&item.ID, &item.ProductID, &item.LotNumber, &item.Manufactured, &item.ExpiresAt, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (s *Store) ApplyMovement(ctx context.Context, movement domain.StockMovement) (domain.StockBalance, error) {
	if err := movement.Validate(); err != nil {
		return domain.StockBalance{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.StockBalance{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	balance, err := applyMovementTx(ctx, tx, movement)
	if err != nil {
		return domain.StockBalance{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.StockBalance{}, err
	}
	return balance, nil
}
func (s *Store) Transfer(ctx context.Context, request domain.TransferRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	keys := []string{balanceLockKey(request.ProductID, request.FromLocationID, request.LotID), balanceLockKey(request.ProductID, request.ToLocationID, request.LotID)}
	sort.Strings(keys)
	for _, key := range keys {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, key); err != nil {
			return err
		}
	}
	reference := strings.TrimSpace(request.Reference)
	if reference == "" {
		reference = request.ID
	}
	out := domain.StockMovement{ID: request.ID + "_out", ProductID: request.ProductID, LocationID: request.FromLocationID, LotID: request.LotID, Type: domain.MovementTransferOut, QuantityDelta: -request.Quantity, Reference: reference, Note: request.Note, ActorID: request.ActorID, OccurredAt: request.OccurredAt}
	if _, err := applyMovementTx(ctx, tx, out); err != nil {
		return err
	}
	in := domain.StockMovement{ID: request.ID + "_in", ProductID: request.ProductID, LocationID: request.ToLocationID, LotID: request.LotID, Type: domain.MovementTransferIn, QuantityDelta: request.Quantity, Reference: reference, Note: request.Note, ActorID: request.ActorID, OccurredAt: request.OccurredAt}
	if _, err := applyMovementTx(ctx, tx, in); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (s *Store) GetBalance(ctx context.Context, productID, locationID, lotID string) (domain.StockBalance, error) {
	query := `SELECT product_id, location_id, COALESCE(lot_id, ''), quantity, updated_at FROM inventory_balances WHERE product_id=$1 AND location_id=$2 AND lot_id IS NULL`
	args := []any{productID, locationID}
	if lotID != "" {
		query = `SELECT product_id, location_id, COALESCE(lot_id, ''), quantity, updated_at FROM inventory_balances WHERE product_id=$1 AND location_id=$2 AND lot_id=$3`
		args = append(args, lotID)
	}
	var balance domain.StockBalance
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&balance.ProductID, &balance.LocationID, &balance.LotID, &balance.Quantity, &balance.UpdatedAt); err != nil {
		return domain.StockBalance{}, mapError(err)
	}
	return balance, nil
}
func (s *Store) ListLowStock(ctx context.Context, limit int) ([]domain.StockBalance, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	rows, err := s.pool.Query(ctx, `SELECT b.product_id,b.location_id,COALESCE(b.lot_id,''),b.quantity,b.updated_at FROM inventory_balances b JOIN products p ON p.id=b.product_id WHERE p.active=true AND b.quantity <= p.reorder_point ORDER BY b.quantity ASC,p.name ASC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.StockBalance, 0)
	for rows.Next() {
		var item domain.StockBalance
		if err := rows.Scan(&item.ProductID, &item.LocationID, &item.LotID, &item.Quantity, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func applyMovementTx(ctx context.Context, tx pgx.Tx, movement domain.StockMovement) (domain.StockBalance, error) {
	if err := movement.Validate(); err != nil {
		return domain.StockBalance{}, err
	}
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, balanceLockKey(movement.ProductID, movement.LocationID, movement.LotID)); err != nil {
		return domain.StockBalance{}, err
	}
	if err := validateLotPolicy(ctx, tx, movement.ProductID, movement.LotID); err != nil {
		return domain.StockBalance{}, err
	}
	var balanceID string
	var current int64
	query := `SELECT id, quantity FROM inventory_balances WHERE product_id=$1 AND location_id=$2 AND lot_id IS NULL FOR UPDATE`
	args := []any{movement.ProductID, movement.LocationID}
	if movement.LotID != "" {
		query = `SELECT id, quantity FROM inventory_balances WHERE product_id=$1 AND location_id=$2 AND lot_id=$3 FOR UPDATE`
		args = append(args, movement.LotID)
	}
	err := tx.QueryRow(ctx, query, args...).Scan(&balanceID, &current)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domain.StockBalance{}, err
	}
	newQuantity := current + movement.QuantityDelta
	if newQuantity < 0 {
		return domain.StockBalance{}, domain.ErrInsufficientStock
	}
	var updatedAt time.Time
	if errors.Is(err, pgx.ErrNoRows) {
		if movement.QuantityDelta < 0 {
			return domain.StockBalance{}, domain.ErrInsufficientStock
		}
		balanceID, err = idgen.New("bal")
		if err != nil {
			return domain.StockBalance{}, err
		}
		if err := tx.QueryRow(ctx, `INSERT INTO inventory_balances (id,product_id,location_id,lot_id,quantity) VALUES ($1,$2,$3,NULLIF($4,''),$5) RETURNING updated_at`, balanceID, movement.ProductID, movement.LocationID, movement.LotID, newQuantity).Scan(&updatedAt); err != nil {
			return domain.StockBalance{}, mapError(err)
		}
	} else {
		if err := tx.QueryRow(ctx, `UPDATE inventory_balances SET quantity=$2,updated_at=now() WHERE id=$1 RETURNING updated_at`, balanceID, newQuantity).Scan(&updatedAt); err != nil {
			return domain.StockBalance{}, err
		}
	}
	_, err = tx.Exec(ctx, `INSERT INTO stock_movements (id,product_id,location_id,lot_id,movement_type,quantity_delta,reference,note,actor_id,occurred_at) VALUES ($1,$2,$3,NULLIF($4,''),$5,$6,$7,$8,$9,$10)`, movement.ID, movement.ProductID, movement.LocationID, movement.LotID, movement.Type, movement.QuantityDelta, strings.TrimSpace(movement.Reference), strings.TrimSpace(movement.Note), strings.TrimSpace(movement.ActorID), movement.OccurredAt)
	if err != nil {
		return domain.StockBalance{}, mapError(err)
	}
	return domain.StockBalance{ProductID: movement.ProductID, LocationID: movement.LocationID, LotID: movement.LotID, Quantity: newQuantity, UpdatedAt: updatedAt}, nil
}
func validateLotPolicy(ctx context.Context, row queryRower, productID, lotID string) error {
	var trackLots bool
	if err := row.QueryRow(ctx, `SELECT track_lots FROM products WHERE id=$1`, productID).Scan(&trackLots); err != nil {
		return mapError(err)
	}
	if trackLots && lotID == "" {
		return fmt.Errorf("%w: this product requires a lot", domain.ErrInvalid)
	}
	if !trackLots && lotID != "" {
		return fmt.Errorf("%w: this product does not use lots", domain.ErrInvalid)
	}
	if lotID != "" {
		var exists bool
		if err := row.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM lots WHERE id=$1 AND product_id=$2)`, lotID, productID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("%w: lot does not belong to product", domain.ErrInvalid)
		}
	}
	return nil
}

type queryRower interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func balanceLockKey(productID, locationID, lotID string) string {
	return productID + "|" + locationID + "|" + lotID
}

var _ repository.Inventory = (*Store)(nil)

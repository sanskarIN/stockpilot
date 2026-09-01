package repository

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type LotInventoryFilter struct {
	ProductID  string
	WarehouseID string
	LocationID string
	LotID      string
	ExpiringBy *time.Time
	Limit      int
	Offset     int
}

type Inventory interface {
	CreateWarehouse(context.Context, domain.Warehouse) error
	UpdateWarehouse(context.Context, domain.Warehouse) error
	ListWarehouses(context.Context, bool) ([]domain.Warehouse, error)
	CreateLocation(context.Context, domain.Location) error
	UpdateLocation(context.Context, domain.Location) error
	ListLocations(context.Context, string, bool) ([]domain.Location, error)
	CreateLot(context.Context, domain.Lot) error
	ListLots(context.Context, string, int) ([]domain.Lot, error)
	ApplyMovement(context.Context, domain.StockMovement) (domain.StockBalance, error)
	Transfer(context.Context, domain.TransferRequest) error
	GetBalance(context.Context, string, string, string) (domain.StockBalance, error)
	ListLowStock(context.Context, int) ([]domain.StockBalance, error)
	ListReorderSuggestions(context.Context, int) ([]domain.ReorderSuggestion, error)
	GetInventoryValuation(context.Context, int) (domain.InventoryValuationReport, error)
	ListLotInventory(context.Context, LotInventoryFilter) ([]domain.LotInventoryRow, error)
}

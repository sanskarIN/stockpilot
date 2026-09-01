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

type LotInventory interface {
	ListLotInventory(context.Context, LotInventoryFilter) ([]domain.LotInventoryRow, error)
}

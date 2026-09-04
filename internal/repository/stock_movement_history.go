package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

// StockMovementHistory reads aggregated inventory movement activity without
// changing the existing Inventory interface used by mutation-focused callers.
type StockMovementHistory interface {
	GetStockMovementHistory(context.Context, int, int) (domain.StockMovementHistoryReport, error)
}

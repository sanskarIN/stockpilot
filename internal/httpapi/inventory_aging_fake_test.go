package httpapi

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (f *fakeStore) GetInventoryAging(context.Context, int) (domain.InventoryAgingReport, error) {
	return domain.InventoryAgingReport{Items: []domain.InventoryAgingItem{}}, nil
}

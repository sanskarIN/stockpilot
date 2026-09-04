package httpapi

import (
	"net/http"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

const (
	defaultInventoryAgingRows = 1000
	maxInventoryAgingRows     = 5000
)

func (a *API) inventoryAging(w http.ResponseWriter, r *http.Request) {
	limit := normalizeInventoryAgingLimit(queryInt(r, "limit", defaultInventoryAgingRows))
	report, err := a.inventory.GetInventoryAging(r.Context(), limit)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func normalizeInventoryAgingLimit(limit int) int {
	if limit <= 0 {
		return defaultInventoryAgingRows
	}
	if limit > maxInventoryAgingRows {
		return maxInventoryAgingRows
	}
	return limit
}

var _ = domain.Aging0To30

package httpapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (a *API) listLotInventory(w http.ResponseWriter, r *http.Request) {
	filter := repository.LotInventoryFilter{ProductID: strings.TrimSpace(r.URL.Query().Get("productId")), WarehouseID: strings.TrimSpace(r.URL.Query().Get("warehouseId")), LocationID: strings.TrimSpace(r.URL.Query().Get("locationId")), LotID: strings.TrimSpace(r.URL.Query().Get("lotId")), Limit: queryInt(r, "limit", 100), Offset: queryInt(r, "offset", 0)}
	if value := strings.TrimSpace(r.URL.Query().Get("expiringBy")); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "expiringBy must use YYYY-MM-DD"}); return }
		filter.ExpiringBy = &parsed
	}
	items, err := a.inventory.ListLotInventory(r.Context(), filter)
	if err != nil { writeDomainError(w, err); return }
	if items == nil { items = []domain.LotInventoryRow{} }
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

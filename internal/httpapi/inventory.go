package httpapi

import (
	"net/http"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
)

func (a *API) listWarehouses(w http.ResponseWriter, r *http.Request) {
	items, err := a.inventory.ListWarehouses(r.Context(), queryBool(r, "active", false))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (a *API) createWarehouse(w http.ResponseWriter, r *http.Request) {
	var item domain.Warehouse
	if !decodeJSON(w, r, &item) {
		return
	}
	id, err := idgen.New("wh")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	item.ID = id
	if item.Timezone == "" {
		item.Timezone = "UTC"
	}
	if err := a.inventory.CreateWarehouse(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) listLocations(w http.ResponseWriter, r *http.Request) {
	items, err := a.inventory.ListLocations(r.Context(), r.URL.Query().Get("warehouseId"), queryBool(r, "active", false))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (a *API) createLocation(w http.ResponseWriter, r *http.Request) {
	var item domain.Location
	if !decodeJSON(w, r, &item) {
		return
	}
	id, err := idgen.New("loc")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	item.ID = id
	if err := a.inventory.CreateLocation(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) createLot(w http.ResponseWriter, r *http.Request) {
	var item domain.Lot
	if !decodeJSON(w, r, &item) {
		return
	}
	id, err := idgen.New("lot")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	item.ID = id
	if err := a.inventory.CreateLot(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) applyMovement(w http.ResponseWriter, r *http.Request) {
	var movement domain.StockMovement
	if !decodeJSON(w, r, &movement) {
		return
	}
	id, err := idgen.New("mov")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	movement.ID = id
	if movement.OccurredAt.IsZero() {
		movement.OccurredAt = time.Now().UTC()
	}
	balance, err := a.inventory.ApplyMovement(r.Context(), movement)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"movementId": movement.ID, "balance": balance})
}

func (a *API) transfer(w http.ResponseWriter, r *http.Request) {
	var transfer domain.TransferRequest
	if !decodeJSON(w, r, &transfer) {
		return
	}
	id, err := idgen.New("xfer")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	transfer.ID = id
	if transfer.OccurredAt.IsZero() {
		transfer.OccurredAt = time.Now().UTC()
	}
	if err := a.inventory.Transfer(r.Context(), transfer); err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"transferId": transfer.ID, "status": "completed"})
}

func (a *API) listLowStock(w http.ResponseWriter, r *http.Request) {
	items, err := a.inventory.ListLowStock(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

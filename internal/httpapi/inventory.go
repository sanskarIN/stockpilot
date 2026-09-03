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
	a.recordAudit(r.Context(), authenticatedActorID(r), "warehouse.created", "warehouse", item.ID, map[string]any{"code": item.Code, "name": item.Name})
	writeJSON(w, http.StatusCreated, item)
}
func (a *API) updateWarehouse(w http.ResponseWriter, r *http.Request) {
	var item domain.Warehouse
	if !decodeJSON(w, r, &item) {
		return
	}
	item.ID = r.PathValue("id")
	if item.Timezone == "" {
		item.Timezone = "UTC"
	}
	if err := a.inventory.UpdateWarehouse(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "warehouse.updated", "warehouse", item.ID, map[string]any{"code": item.Code, "name": item.Name})
	writeJSON(w, http.StatusOK, item)
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
	a.recordAudit(r.Context(), authenticatedActorID(r), "location.created", "location", item.ID, map[string]any{"warehouseId": item.WarehouseID, "code": item.Code, "name": item.Name})
	writeJSON(w, http.StatusCreated, item)
}
func (a *API) updateLocation(w http.ResponseWriter, r *http.Request) {
	var item domain.Location
	if !decodeJSON(w, r, &item) {
		return
	}
	item.ID = r.PathValue("id")
	if err := a.inventory.UpdateLocation(r.Context(), item); err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "location.updated", "location", item.ID, map[string]any{"warehouseId": item.WarehouseID, "code": item.Code, "name": item.Name, "active": item.Active})
	writeJSON(w, http.StatusOK, item)
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
	a.recordAudit(r.Context(), authenticatedActorID(r), "lot.created", "lot", item.ID, map[string]any{"productId": item.ProductID, "lotNumber": item.LotNumber, "expiresAt": item.ExpiresAt})
	writeJSON(w, http.StatusCreated, item)
}
func (a *API) listLots(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productId")
	if productID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "productId is required"})
		return
	}
	items, err := a.inventory.ListLots(r.Context(), productID, queryInt(r, "limit", 50))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	if items == nil {
		items = []domain.Lot{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
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
	movement.ActorID = authenticatedActorID(r)
	if movement.OccurredAt.IsZero() {
		movement.OccurredAt = time.Now().UTC()
	}
	balance, err := a.inventory.ApplyMovement(r.Context(), movement)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), movement.ActorID, "inventory.movement.applied", "inventory_balance", movement.LocationID, map[string]any{"movementId": movement.ID, "productId": movement.ProductID, "lotId": movement.LotID, "type": movement.Type, "quantityDelta": movement.QuantityDelta})
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
	transfer.ActorID = authenticatedActorID(r)
	if transfer.OccurredAt.IsZero() {
		transfer.OccurredAt = time.Now().UTC()
	}
	if err := a.inventory.Transfer(r.Context(), transfer); err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), transfer.ActorID, "inventory.transfer.completed", "inventory_transfer", transfer.ID, map[string]any{"productId": transfer.ProductID, "fromLocationId": transfer.FromLocationID, "toLocationId": transfer.ToLocationID, "quantity": transfer.Quantity, "lotId": transfer.LotID})
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
func authenticatedActorID(r *http.Request) string {
	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		return ""
	}
	return principal.User.ID
}

package httpapi

import (
	"net/http"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
)

func (a *API) listOrders(w http.ResponseWriter, r *http.Request) {
	status := domain.PurchaseOrderStatus(r.URL.Query().Get("status"))
	items, err := a.orders.ListOrders(r.Context(), status, queryInt(r, "limit", 50), queryInt(r, "offset", 0))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.PurchaseOrder
	if !decodeJSON(w, r, &order) {
		return
	}
	id, err := idgen.New("po")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	order.ID = id
	order.CreatedBy = authenticatedActorID(r)
	if order.Status == "" {
		order.Status = domain.PurchaseOrderDraft
	}
	if order.Currency == "" {
		order.Currency = "INR"
	}
	for i := range order.Lines {
		lineID, err := idgen.New("pol")
		if err != nil {
			writeDomainError(w, err)
			return
		}
		order.Lines[i].ID = lineID
		order.Lines[i].PurchaseOrderID = order.ID
	}
	if err := a.orders.CreateOrder(r.Context(), order); err != nil {
		writeDomainError(w, err)
		return
	}
	created, err := a.orders.GetOrder(r.Context(), order.ID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), order.CreatedBy, "purchase_order.created", "purchase_order", created.ID, map[string]any{"status": created.Status, "supplierId": created.SupplierID, "warehouseId": created.WarehouseID, "lineCount": len(created.Lines)})
	writeJSON(w, http.StatusCreated, created)
}
func (a *API) updateOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.PurchaseOrder
	if !decodeJSON(w, r, &order) {
		return
	}
	order.ID = r.PathValue("id")
	if order.Status == "" {
		order.Status = domain.PurchaseOrderDraft
	}
	if order.Currency == "" {
		order.Currency = "INR"
	}
	for i := range order.Lines {
		if order.Lines[i].ID == "" {
			id, err := idgen.New("pol")
			if err != nil {
				writeDomainError(w, err)
				return
			}
			order.Lines[i].ID = id
		}
		order.Lines[i].PurchaseOrderID = order.ID
		order.Lines[i].Received = 0
	}
	if err := a.orders.UpdateOrder(r.Context(), order); err != nil {
		writeDomainError(w, err)
		return
	}
	updated, err := a.orders.GetOrder(r.Context(), order.ID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), authenticatedActorID(r), "purchase_order.updated", "purchase_order", updated.ID, map[string]any{"lineCount": len(updated.Lines), "status": updated.Status})
	writeJSON(w, http.StatusOK, updated)
}
func (a *API) getOrder(w http.ResponseWriter, r *http.Request) {
	item, err := a.orders.GetOrder(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (a *API) receiveOrderLine(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Quantity   int64  `json:"quantity"`
		LocationID string `json:"locationId"`
		LotID      string `json:"lotId,omitempty"`
		NewLot     *struct {
			LotNumber      string     `json:"lotNumber"`
			ManufacturedAt *time.Time `json:"manufacturedAt,omitempty"`
			ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
		} `json:"newLot,omitempty"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	actor := authenticatedActorID(r)
	var err error
	usingNewLot := body.NewLot != nil
	if usingNewLot {
		lot := domain.Lot{LotNumber: body.NewLot.LotNumber, Manufactured: body.NewLot.ManufacturedAt, ExpiresAt: body.NewLot.ExpiresAt}
		err = a.orders.ReceiveLineWithNewLot(r.Context(), r.PathValue("orderID"), r.PathValue("lineID"), body.Quantity, body.LocationID, lot, actor)
	} else {
		err = a.orders.ReceiveLine(r.Context(), r.PathValue("orderID"), r.PathValue("lineID"), body.Quantity, body.LocationID, body.LotID, actor)
	}
	if err != nil {
		writeDomainError(w, err)
		return
	}
	updated, err := a.orders.GetOrder(r.Context(), r.PathValue("orderID"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), actor, "purchase_order.receipt_recorded", "purchase_order", updated.ID, map[string]any{"lineId": r.PathValue("lineID"), "quantity": body.Quantity, "locationId": body.LocationID, "lotId": body.LotID, "newLot": usingNewLot, "status": updated.Status})
	writeJSON(w, http.StatusOK, updated)
}

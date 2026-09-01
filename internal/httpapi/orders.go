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
	if err != nil { writeDomainError(w, err); return }
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.PurchaseOrder
	if !decodeJSON(w, r, &order) { return }
	id, err := idgen.New("po"); if err != nil { writeDomainError(w, err); return }
	order.ID = id; order.CreatedBy = authenticatedActorID(r)
	if order.Status == "" { order.Status = domain.PurchaseOrderDraft }
	if order.Currency == "" { order.Currency = "INR" }
	for i := range order.Lines { lineID, err := idgen.New("pol"); if err != nil { writeDomainError(w, err); return }; order.Lines[i].ID = lineID; order.Lines[i].PurchaseOrderID = order.ID }
	if err := a.orders.CreateOrder(r.Context(), order); err != nil { writeDomainError(w, err); return }
	created, err := a.orders.GetOrder(r.Context(), order.ID); if err != nil { writeDomainError(w, err); return }
	writeJSON(w, http.StatusCreated, created)
}

func (a *API) getOrder(w http.ResponseWriter, r *http.Request) {
	item, err := a.orders.GetOrder(r.Context(), r.PathValue("id")); if err != nil { writeDomainError(w, err); return }
	writeJSON(w, http.StatusOK, item)
}

func (a *API) receiveOrderLine(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Quantity   int64 `json:"quantity"`
		LocationID string `json:"locationId"`
		LotID      string `json:"lotId,omitempty"`
		NewLot     *struct {
			LotNumber     string     `json:"lotNumber"`
			ManufacturedAt *time.Time `json:"manufacturedAt,omitempty"`
			ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
		} `json:"newLot,omitempty"`
	}
	if !decodeJSON(w, r, &body) { return }

	var err error
	if body.NewLot != nil {
		lot := domain.Lot{LotNumber: body.NewLot.LotNumber, Manufactured: body.NewLot.ManufacturedAt, ExpiresAt: body.NewLot.ExpiresAt}
		err = a.orders.ReceiveLineWithNewLot(r.Context(), r.PathValue("orderID"), r.PathValue("lineID"), body.Quantity, body.LocationID, lot, authenticatedActorID(r))
	} else {
		err = a.orders.ReceiveLine(r.Context(), r.PathValue("orderID"), r.PathValue("lineID"), body.Quantity, body.LocationID, body.LotID, authenticatedActorID(r))
	}
	if err != nil { writeDomainError(w, err); return }
	updated, err := a.orders.GetOrder(r.Context(), r.PathValue("orderID")); if err != nil { writeDomainError(w, err); return }
	writeJSON(w, http.StatusOK, updated)
}

package httpapi

import (
	"net/http"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type orderStatusRequest struct {
	Status domain.PurchaseOrderStatus `json:"status"`
}

func (a *API) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var body orderStatusRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if err := a.orders.UpdateOrderStatus(r.Context(), r.PathValue("id"), body.Status, authenticatedActorID(r)); err != nil {
		writeDomainError(w, err)
		return
	}
	updated, err := a.orders.GetOrder(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

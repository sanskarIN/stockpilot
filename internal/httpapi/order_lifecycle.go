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
	actor := authenticatedActorID(r)
	if err := a.orders.UpdateOrderStatus(r.Context(), r.PathValue("id"), body.Status, actor); err != nil {
		writeDomainError(w, err)
		return
	}
	updated, err := a.orders.GetOrder(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), actor, "purchase_order.status_changed", "purchase_order", updated.ID, map[string]any{"status": updated.Status})
	writeJSON(w, http.StatusOK, updated)
}

package httpapi

import (
	"net/http"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (a *API) listLots(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productId")
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

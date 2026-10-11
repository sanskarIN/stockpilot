package httpapi

import "net/http"

func (a *API) listReorderSuggestions(w http.ResponseWriter, r *http.Request) {
	items, err := a.inventory.ListReorderSuggestions(r.Context(), queryInt(r, "limit", 50))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (a *API) inventoryValuation(w http.ResponseWriter, r *http.Request) {
	report, err := a.inventory.GetInventoryValuation(r.Context(), queryInt(r, "limit", 100))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

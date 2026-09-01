package httpapi

import "net/http"

func (a *API) getProductByBarcode(w http.ResponseWriter, r *http.Request) {
	item, err := a.catalog.GetProductByBarcode(r.Context(), r.PathValue("barcode"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

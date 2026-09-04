package httpapi

import "net/http"

// setCSVDownloadHeaders applies a consistent privacy-oriented download policy
// to every CSV export. Exported datasets can contain operational information,
// so browsers and intermediary caches should not retain them by default.
func setCSVDownloadHeaders(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}

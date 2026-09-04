package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sanskarIN/stockpilot/internal/reporting"
)

func writeReportExportHeaders(w http.ResponseWriter, filename string, period reporting.Period, bounds reporting.Bounds, complete bool, generatedAt time.Time) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	metadata := reporting.NewMetadata(period, bounds, complete, generatedAt)
	w.Header().Set("X-Report-Generated-At", metadata.GeneratedAt.UTC().Format(time.RFC3339))
	w.Header().Set("X-Report-From", metadata.From.UTC().Format(time.RFC3339))
	w.Header().Set("X-Report-To", metadata.To.UTC().Format(time.RFC3339))
	w.Header().Set("X-Report-Limit", strconv.Itoa(metadata.Limit))
	w.Header().Set("X-Report-Offset", strconv.Itoa(metadata.Offset))
	w.Header().Set("X-Report-Complete", strconv.FormatBool(metadata.Complete))
}

package httpapi

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/reporting"
)

func TestWriteReportExportHeaders(t *testing.T) {
	now := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	period, err := reporting.NewPeriod(now.AddDate(0, 0, -29), now)
	if err != nil {
		t.Fatal(err)
	}
	bounds, err := reporting.NewBounds(25, 5, 1, 5000)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	writeReportExportHeaders(recorder, "report.csv", period, bounds, false, now)

	checks := map[string]string{
		"Content-Type":          "text/csv; charset=utf-8",
		"Content-Disposition":   "attachment; filename=report.csv",
		"Cache-Control":         "no-store, no-cache, must-revalidate",
		"Pragma":                "no-cache",
		"X-Report-Generated-At": "2026-09-04T12:00:00Z",
		"X-Report-Limit":        "25",
		"X-Report-Offset":       "5",
		"X-Report-Complete":     "false",
	}
	for header, want := range checks {
		if got := recorder.Header().Get(header); got != want {
			t.Fatalf("%s=%q want %q", header, got, want)
		}
	}
}

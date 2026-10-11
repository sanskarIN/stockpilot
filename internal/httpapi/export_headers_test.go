package httpapi

import (
	"net/http/httptest"
	"testing"
)

func TestSetCSVDownloadHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	setCSVDownloadHeaders(recorder, "stockpilot-test.csv")

	if got := recorder.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content type=%q", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="stockpilot-test.csv"` {
		t.Fatalf("content disposition=%q", got)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache control=%q", got)
	}
	if got := recorder.Header().Get("Pragma"); got != "no-cache" {
		t.Fatalf("pragma=%q", got)
	}
}

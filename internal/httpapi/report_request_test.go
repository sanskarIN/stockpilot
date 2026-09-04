package httpapi

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseReportRequestDefaults(t *testing.T) {
	now := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)
	r := httptest.NewRequest("GET", "/api/v1/reports/inventory", nil)
	request, err := parseReportRequest(r, now)
	if err != nil {
		t.Fatal(err)
	}
	if request.Period.Days() != 30 {
		t.Fatalf("days=%d", request.Period.Days())
	}
	if request.Bounds.Limit != 100 || request.Bounds.Offset != 0 {
		t.Fatalf("bounds=%+v", request.Bounds)
	}
}

func TestParseReportRequestAcceptsPeriodAndBounds(t *testing.T) {
	r := httptest.NewRequest("GET", "/api/v1/reports/inventory?from=2026-09-01T00:00:00Z&to=2026-09-04T00:00:00Z&limit=25&offset=50", nil)
	request, err := parseReportRequest(r, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if request.Period.Days() != 4 || request.Bounds.Limit != 25 || request.Bounds.Offset != 50 {
		t.Fatalf("request=%+v", request)
	}
}

func TestParseReportRequestRejectsIncompletePeriod(t *testing.T) {
	r := httptest.NewRequest("GET", "/api/v1/reports/inventory?from=2026-09-01T00:00:00Z", nil)
	if _, err := parseReportRequest(r, time.Now()); err == nil {
		t.Fatal("expected incomplete period error")
	}
}

func TestParseReportRequestRejectsInvalidBounds(t *testing.T) {
	for _, query := range []string{"?limit=0", "?limit=5001", "?offset=-1", "?offset=abc"} {
		r := httptest.NewRequest("GET", "/api/v1/reports/inventory"+query, nil)
		if _, err := parseReportRequest(r, time.Now()); err == nil {
			t.Fatalf("expected bounds error for %s", query)
		}
	}
}

func TestWriteReportMetadata(t *testing.T) {
	now := time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC)
	r := httptest.NewRequest("GET", "/api/v1/reports/inventory?limit=25&offset=5", nil)
	request, err := parseReportRequest(r, now)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	writeReportMetadata(w, request, now, true)
	for key, want := range map[string]string{
		"X-Report-Generated-At": "2026-09-04T10:00:00Z",
		"X-Report-From":         "2026-08-06T10:00:00Z",
		"X-Report-To":           "2026-09-04T10:00:00Z",
		"X-Report-Limit":        "25",
		"X-Report-Offset":       "5",
		"X-Report-Complete":     "true",
	} {
		if got := w.Header().Get(key); got != want {
			t.Fatalf("%s=%q want %q", key, got, want)
		}
	}
}

package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReportInventoryEmitsMetadataHeaders(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil }, WithInsights(store, nil))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/inventory?limit=25&offset=5", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if recorder.Header().Get("X-Report-Limit") != "25" || recorder.Header().Get("X-Report-Offset") != "5" {
		t.Fatalf("bounds headers: limit=%q offset=%q", recorder.Header().Get("X-Report-Limit"), recorder.Header().Get("X-Report-Offset"))
	}
	if recorder.Header().Get("X-Report-Complete") != "true" {
		t.Fatalf("complete=%q", recorder.Header().Get("X-Report-Complete"))
	}
	if _, err := time.Parse(time.RFC3339, recorder.Header().Get("X-Report-Generated-At")); err != nil {
		t.Fatalf("generated-at=%q: %v", recorder.Header().Get("X-Report-Generated-At"), err)
	}
}

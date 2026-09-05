package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORSExposesReportPaginationHeaders(t *testing.T) {
	handler := WrapCommon(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "37")
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}), []string{"https://app.stockpilot.example"}, nil)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/supplier-performance", nil)
	request.Header.Set("Origin", "https://app.stockpilot.example")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Header().Get("Access-Control-Expose-Headers"), "X-Total-Count") {
		t.Fatalf("expose headers=%q", recorder.Header().Get("Access-Control-Expose-Headers"))
	}
	if recorder.Header().Get("X-Total-Count") != "37" {
		t.Fatalf("X-Total-Count=%q", recorder.Header().Get("X-Total-Count"))
	}
}

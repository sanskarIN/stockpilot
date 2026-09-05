package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetaEndpointReportsReleaseVersion(t *testing.T) {
	handler := NewCore(nil, nil, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/meta", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode meta response: %v", err)
	}
	if response.Name != "StockPilot" {
		t.Fatalf("name=%q", response.Name)
	}
	if response.Version != "0.4.0" {
		t.Fatalf("version=%q", response.Version)
	}
}

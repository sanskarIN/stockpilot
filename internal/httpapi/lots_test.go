package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestListLotsEndpointFiltersByProduct(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/lots?productId=prd_1&limit=20", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"lotNumber":"LOT-1"`) {
		t.Fatalf("body = %q, want lot response", recorder.Body.String())
	}
}

func TestListLotsEndpointRequiresProduct(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, store, store, func(context.Context) error { return nil }, nil, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/lots", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

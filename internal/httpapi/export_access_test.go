package httpapi

import (
	"net/http"
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestExportRoutesUseReadPermissions(t *testing.T) {
	tests := []struct {
		path       string
		permission domain.Permission
	}{
		{"/api/v1/products/export.csv", domain.PermissionCatalogRead},
		{"/api/v1/inventory/export.csv", domain.PermissionInventoryRead},
		{"/api/v1/inventory/low-stock/export.csv", domain.PermissionInventoryRead},
		{"/api/v1/inventory/reorder-suggestions/export.csv", domain.PermissionInventoryRead},
		{"/api/v1/inventory/lots/export.csv", domain.PermissionInventoryRead},
		{"/api/v1/inventory/receipts/export.csv", domain.PermissionInventoryRead},
		{"/api/v1/orders/export.csv", domain.PermissionOrdersRead},
		{"/api/v1/audit/export.csv", domain.PermissionAuditRead},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			permission, ok := permissionFor(httpRequest(http.MethodGet, tt.path))
			if !ok || permission != tt.permission {
				t.Fatalf("permission=%q ok=%v want=%q", permission, ok, tt.permission)
			}
		})
	}
}

func httpRequest(method, path string) *http.Request {
	return &http.Request{Method: method, URL: mustParseURL(path)}
}

func mustParseURL(path string) *url.URL {
	value, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return value
}

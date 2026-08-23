package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRolePermissions(t *testing.T) {
	tests := []struct {
		role       Role
		permission Permission
		want       bool
	}{
		{RoleAdmin, PermissionUsersManage, true},
		{RoleManager, PermissionCatalogWrite, true},
		{RoleManager, PermissionUsersManage, false},
		{RoleOperator, PermissionInventoryWrite, true},
		{RoleOperator, PermissionAuditRead, false},
		{RoleViewer, PermissionReportsRead, true},
		{RoleViewer, PermissionOrdersWrite, false},
	}
	for _, tt := range tests {
		name := strings.Join([]string{string(tt.role), string(tt.permission)}, "/")
		t.Run(name, func(t *testing.T) {
			if got := tt.role.Can(tt.permission); got != tt.want {
				t.Fatalf("Can(%q) = %v, want %v", tt.permission, got, tt.want)
			}
		})
	}
}

func TestUserValidateRejectsInvalidEmailAndRole(t *testing.T) {
	user := User{ID: "usr_1", Email: "not-an-email", DisplayName: "Inventory Admin", Role: Role("owner")}
	if err := user.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid", err)
	}

	user.Email = "admin@example.com"
	if err := user.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid for role", err)
	}
}

func TestSessionValidateRequiresFutureExpiry(t *testing.T) {
	now := time.Now().UTC()
	session := Session{ID: "ses_1", UserID: "usr_1", TokenHash: strings.Repeat("a", 64), ExpiresAt: now.Add(-time.Minute)}
	if err := session.Validate(now); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid", err)
	}
}

package domain

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleManager, RoleOperator, RoleViewer:
		return true
	default:
		return false
	}
}

type Permission string

const (
	PermissionCatalogRead    Permission = "catalog:read"
	PermissionCatalogWrite   Permission = "catalog:write"
	PermissionInventoryRead  Permission = "inventory:read"
	PermissionInventoryWrite Permission = "inventory:write"
	PermissionOrdersRead     Permission = "orders:read"
	PermissionOrdersWrite    Permission = "orders:write"
	PermissionReportsRead    Permission = "reports:read"
	PermissionAuditRead      Permission = "audit:read"
	PermissionUsersManage    Permission = "users:manage"
	PermissionSettingsManage Permission = "settings:manage"
)

func (r Role) Can(permission Permission) bool {
	if r == RoleAdmin {
		return true
	}
	switch r {
	case RoleManager:
		switch permission {
		case PermissionCatalogRead, PermissionCatalogWrite,
			PermissionInventoryRead, PermissionInventoryWrite,
			PermissionOrdersRead, PermissionOrdersWrite,
			PermissionReportsRead, PermissionAuditRead:
			return true
		}
	case RoleOperator:
		switch permission {
		case PermissionCatalogRead, PermissionInventoryRead, PermissionInventoryWrite, PermissionOrdersRead, PermissionOrdersWrite:
			return true
		}
	case RoleViewer:
		switch permission {
		case PermissionCatalogRead, PermissionInventoryRead, PermissionOrdersRead, PermissionReportsRead:
			return true
		}
	}
	return false
}

type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"displayName"`
	Role        Role       `json:"role"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("%w: user id is required", ErrInvalid)
	}
	address, err := mail.ParseAddress(strings.TrimSpace(u.Email))
	if err != nil || !strings.EqualFold(address.Address, strings.TrimSpace(u.Email)) {
		return fmt.Errorf("%w: valid email is required", ErrInvalid)
	}
	if name := strings.TrimSpace(u.DisplayName); len(name) < 2 || len(name) > 120 {
		return fmt.Errorf("%w: display name must be 2-120 characters", ErrInvalid)
	}
	if !u.Role.Valid() {
		return fmt.Errorf("%w: invalid role", ErrInvalid)
	}
	return nil
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	LastSeen  time.Time `json:"lastSeenAt"`
}

func (s Session) Validate(now time.Time) error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.UserID) == "" {
		return fmt.Errorf("%w: session id and user id are required", ErrInvalid)
	}
	if len(s.TokenHash) != 64 {
		return fmt.Errorf("%w: session token hash must be SHA-256 hex", ErrInvalid)
	}
	if !s.ExpiresAt.After(now) {
		return fmt.Errorf("%w: session expiry must be in the future", ErrInvalid)
	}
	return nil
}

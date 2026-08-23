package domain

import (
	"fmt"
	"strings"
	"time"
)

type Warehouse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   string    `json:"address,omitempty"`
	Timezone  string    `json:"timezone"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (w Warehouse) Validate() error {
	if strings.TrimSpace(w.ID) == "" {
		return fmt.Errorf("%w: warehouse id is required", ErrInvalid)
	}
	if c := strings.TrimSpace(w.Code); len(c) < 2 || len(c) > 48 {
		return fmt.Errorf("%w: warehouse code must be 2-48 characters", ErrInvalid)
	}
	if n := strings.TrimSpace(w.Name); len(n) < 2 || len(n) > 160 {
		return fmt.Errorf("%w: warehouse name must be 2-160 characters", ErrInvalid)
	}
	if strings.TrimSpace(w.Timezone) == "" {
		return fmt.Errorf("%w: warehouse timezone is required", ErrInvalid)
	}
	return nil
}

type Location struct {
	ID          string    `json:"id"`
	WarehouseID string    `json:"warehouseId"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (l Location) Validate() error {
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.WarehouseID) == "" {
		return fmt.Errorf("%w: location id and warehouse id are required", ErrInvalid)
	}
	if c := strings.TrimSpace(l.Code); len(c) < 1 || len(c) > 64 {
		return fmt.Errorf("%w: location code must be 1-64 characters", ErrInvalid)
	}
	if n := strings.TrimSpace(l.Name); len(n) < 1 || len(n) > 160 {
		return fmt.Errorf("%w: location name must be 1-160 characters", ErrInvalid)
	}
	return nil
}

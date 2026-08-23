package domain

import (
	"fmt"
	"strings"
	"time"
)

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (c Category) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("%w: category id is required", ErrInvalid)
	}
	if n := strings.TrimSpace(c.Name); len(n) < 2 || len(n) > 120 {
		return fmt.Errorf("%w: category name must be 2-120 characters", ErrInvalid)
	}
	return nil
}

type Supplier struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s Supplier) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("%w: supplier id is required", ErrInvalid)
	}
	if c := strings.TrimSpace(s.Code); len(c) < 2 || len(c) > 48 {
		return fmt.Errorf("%w: supplier code must be 2-48 characters", ErrInvalid)
	}
	if n := strings.TrimSpace(s.Name); len(n) < 2 || len(n) > 160 {
		return fmt.Errorf("%w: supplier name must be 2-160 characters", ErrInvalid)
	}
	return nil
}

type Product struct {
	ID              string    `json:"id"`
	SKU             string    `json:"sku"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	CategoryID      string    `json:"categoryId,omitempty"`
	SupplierID      string    `json:"supplierId,omitempty"`
	Barcode         string    `json:"barcode,omitempty"`
	Unit            string    `json:"unit"`
	UnitCostMinor   int64     `json:"unitCostMinor"`
	Currency        string    `json:"currency"`
	ReorderPoint    int64     `json:"reorderPoint"`
	ReorderQuantity int64     `json:"reorderQuantity"`
	TrackLots       bool      `json:"trackLots"`
	TrackExpiry     bool      `json:"trackExpiry"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (p Product) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("%w: product id is required", ErrInvalid)
	}
	if s := strings.TrimSpace(p.SKU); len(s) < 2 || len(s) > 64 {
		return fmt.Errorf("%w: SKU must be 2-64 characters", ErrInvalid)
	}
	if n := strings.TrimSpace(p.Name); len(n) < 2 || len(n) > 200 {
		return fmt.Errorf("%w: product name must be 2-200 characters", ErrInvalid)
	}
	if u := strings.TrimSpace(p.Unit); len(u) < 1 || len(u) > 32 {
		return fmt.Errorf("%w: unit must be 1-32 characters", ErrInvalid)
	}
	if p.UnitCostMinor < 0 {
		return fmt.Errorf("%w: unit cost cannot be negative", ErrInvalid)
	}
	if len(strings.TrimSpace(p.Currency)) != 3 {
		return fmt.Errorf("%w: currency must be a 3-letter code", ErrInvalid)
	}
	if p.ReorderPoint < 0 || p.ReorderQuantity < 0 {
		return fmt.Errorf("%w: reorder values cannot be negative", ErrInvalid)
	}
	if p.TrackExpiry && !p.TrackLots {
		return fmt.Errorf("%w: expiry tracking requires lot tracking", ErrInvalid)
	}
	return nil
}

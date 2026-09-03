package domain

import (
	"errors"
	"testing"
)

func TestProductValidateAcceptsWellFormedProduct(t *testing.T) {
	product := Product{
		ID:              "prd_1",
		SKU:             "SKU-001",
		Name:            "USB-C Cable",
		Unit:            "piece",
		UnitCostMinor:   129900,
		Currency:        "INR",
		ReorderPoint:    10,
		ReorderQuantity: 25,
		TrackLots:       true,
		TrackExpiry:     true,
		Active:          true,
	}
	if err := product.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestProductValidateRejectsExpiryWithoutLotTracking(t *testing.T) {
	product := Product{
		ID:          "prd_1",
		SKU:         "SKU-001",
		Name:        "Perishable item",
		Unit:        "piece",
		Currency:    "INR",
		TrackExpiry: true,
	}
	if err := product.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid", err)
	}
}

func TestProductValidateRejectsNegativeValues(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Product)
	}{
		{name: "unit cost", mutate: func(p *Product) { p.UnitCostMinor = -1 }},
		{name: "reorder point", mutate: func(p *Product) { p.ReorderPoint = -1 }},
		{name: "reorder quantity", mutate: func(p *Product) { p.ReorderQuantity = -1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := Product{ID: "prd_1", SKU: "SKU-001", Name: "Valid Product", Unit: "piece", Currency: "INR"}
			tt.mutate(&product)
			if err := product.Validate(); !errors.Is(err, ErrInvalid) {
				t.Fatalf("Validate() error = %v, want ErrInvalid", err)
			}
		})
	}
}

package domain

import (
	"errors"
	"testing"
	"time"
)

func TestStockMovementValidateEnforcesDirection(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name  string
		type_ MovementType
		delta int64
		valid bool
	}{
		{name: "stock in positive", type_: MovementStockIn, delta: 5, valid: true},
		{name: "stock in negative", type_: MovementStockIn, delta: -5, valid: false},
		{name: "stock out negative", type_: MovementStockOut, delta: -5, valid: true},
		{name: "stock out positive", type_: MovementStockOut, delta: 5, valid: false},
		{name: "adjustment positive", type_: MovementAdjustment, delta: 2, valid: true},
		{name: "adjustment negative", type_: MovementAdjustment, delta: -2, valid: true},
		{name: "adjustment zero", type_: MovementAdjustment, delta: 0, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movement := StockMovement{ID: "mov_1", ProductID: "prd_1", LocationID: "loc_1", Type: tt.type_, QuantityDelta: tt.delta, OccurredAt: now}
			err := movement.Validate()
			if tt.valid && err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}
			if !tt.valid && !errors.Is(err, ErrInvalid) {
				t.Fatalf("Validate() error = %v, want ErrInvalid", err)
			}
		})
	}
}

func TestTransferRequestValidateRejectsSameLocation(t *testing.T) {
	request := TransferRequest{
		ID: "xfer_1", ProductID: "prd_1", FromLocationID: "loc_1", ToLocationID: "loc_1", Quantity: 3, OccurredAt: time.Now().UTC(),
	}
	if err := request.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid", err)
	}
}

func TestLotValidateRejectsExpiryBeforeManufacture(t *testing.T) {
	manufactured := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	expires := manufactured.Add(-time.Hour)
	lot := Lot{ID: "lot_1", ProductID: "prd_1", LotNumber: "LOT-001", Manufactured: &manufactured, ExpiresAt: &expires}
	if err := lot.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid", err)
	}
}

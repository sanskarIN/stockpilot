package postgres

import (
	"math"
	"testing"
)

func TestReorderTarget(t *testing.T) {
	tests := []struct {
		name            string
		reorderPoint    int64
		reorderQuantity int64
		onHand          int64
		wantTarget      int64
		wantSuggested   int64
	}{
		{name: "stockout", reorderPoint: 10, reorderQuantity: 20, onHand: 0, wantTarget: 30, wantSuggested: 30},
		{name: "below point", reorderPoint: 10, reorderQuantity: 20, onHand: 4, wantTarget: 30, wantSuggested: 26},
		{name: "at point", reorderPoint: 10, reorderQuantity: 20, onHand: 10, wantTarget: 30, wantSuggested: 20},
		{name: "already at target", reorderPoint: 10, reorderQuantity: 20, onHand: 30, wantTarget: 30, wantSuggested: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target, suggested, err := reorderTarget(test.reorderPoint, test.reorderQuantity, test.onHand)
			if err != nil {
				t.Fatalf("reorderTarget() error = %v", err)
			}
			if target != test.wantTarget || suggested != test.wantSuggested {
				t.Fatalf("reorderTarget() = (%d, %d), want (%d, %d)", target, suggested, test.wantTarget, test.wantSuggested)
			}
		})
	}
}

func TestReorderTargetRejectsInvalidOrOverflowingInput(t *testing.T) {
	cases := []struct {
		point    int64
		quantity int64
		onHand   int64
	}{
		{point: -1, quantity: 1, onHand: 0},
		{point: 1, quantity: -1, onHand: 0},
		{point: 1, quantity: 1, onHand: -1},
		{point: math.MaxInt64, quantity: 1, onHand: 0},
	}
	for _, test := range cases {
		if _, _, err := reorderTarget(test.point, test.quantity, test.onHand); err == nil {
			t.Fatalf("reorderTarget(%d, %d, %d) expected error", test.point, test.quantity, test.onHand)
		}
	}
}

func TestParseMinorValue(t *testing.T) {
	value, err := parseMinorValue(" 12345 ")
	if err != nil {
		t.Fatalf("parseMinorValue() error = %v", err)
	}
	if value != 12345 {
		t.Fatalf("parseMinorValue() = %d, want 12345", value)
	}

	if _, err := parseMinorValue("999999999999999999999999"); err == nil {
		t.Fatal("parseMinorValue() expected overflow error")
	}
}

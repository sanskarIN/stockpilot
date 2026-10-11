package httpapi

import "testing"

func TestNormalizeInventoryAgingLimit(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "zero", in: 0, want: defaultInventoryAgingRows},
		{name: "negative", in: -10, want: defaultInventoryAgingRows},
		{name: "within limit", in: 250, want: 250},
		{name: "maximum", in: maxInventoryAgingRows, want: maxInventoryAgingRows},
		{name: "above maximum", in: 9000, want: maxInventoryAgingRows},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeInventoryAgingLimit(tt.in); got != tt.want {
				t.Fatalf("normalizeInventoryAgingLimit(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

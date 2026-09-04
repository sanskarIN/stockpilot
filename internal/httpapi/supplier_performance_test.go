package httpapi

import "testing"

func TestNormalizeSupplierPerformanceDays(t *testing.T) {
	cases := []struct{ in, want int }{{0, 30}, {-1, 30}, {30, 30}, {365, 365}, {999, 365}}
	for _, tc := range cases {
		if got := normalizeSupplierPerformanceDays(tc.in); got != tc.want {
			t.Fatalf("days=%d got %d want %d", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeSupplierPerformanceLimit(t *testing.T) {
	cases := []struct{ in, want int }{{0, 1000}, {-1, 1000}, {100, 100}, {5000, 5000}, {9999, 5000}}
	for _, tc := range cases {
		if got := normalizeSupplierPerformanceLimit(tc.in); got != tc.want {
			t.Fatalf("limit=%d got %d want %d", tc.in, got, tc.want)
		}
	}
}

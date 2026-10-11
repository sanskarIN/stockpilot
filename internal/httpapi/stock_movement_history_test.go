package httpapi

import "testing"

func TestNormalizeMovementHistoryWindow(t *testing.T) {
	cases := []struct {
		name     string
		in, want int
	}{
		{"zero", 0, defaultMovementHistoryWindowDays},
		{"negative", -7, defaultMovementHistoryWindowDays},
		{"within", 90, 90},
		{"maximum", maxMovementHistoryWindowDays, maxMovementHistoryWindowDays},
		{"above", 999, maxMovementHistoryWindowDays},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeMovementHistoryWindow(tc.in); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestNormalizeMovementHistoryLimit(t *testing.T) {
	cases := []struct {
		name     string
		in, want int
	}{
		{"zero", 0, defaultMovementHistoryRows},
		{"negative", -1, defaultMovementHistoryRows},
		{"within", 250, 250},
		{"maximum", maxMovementHistoryRows, maxMovementHistoryRows},
		{"above", 9000, maxMovementHistoryRows},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeMovementHistoryLimit(tc.in); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

package httpapi

import (
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestReplenishmentCursorRoundTrip(t *testing.T) {
	want := replenishmentCursor{RiskRank: 2, SuggestedQuantity: 15, SKU: "SKU-015", ProductID: "product-015"}
	got, err := decodeReplenishmentCursor(encodeReplenishmentCursor(want))
	if err != nil {
		t.Fatalf("decode cursor: %v", err)
	}
	if *got != want {
		t.Fatalf("cursor mismatch: got %+v want %+v", *got, want)
	}
}

func TestReplenishmentCursorRejectsMalformedValues(t *testing.T) {
	for _, value := range []string{"not-base64", "", "eyJza3UiOiIiLCJwcm9kdWN0SWQiOiJwMSJ9"} {
		if value == "" {
			continue
		}
		if _, err := decodeReplenishmentCursor(value); err == nil {
			t.Fatalf("expected malformed cursor %q to fail", value)
		}
	}
}

func TestReplenishmentCursorFollowsReportOrdering(t *testing.T) {
	cursor := replenishmentCursor{RiskRank: 2, SuggestedQuantity: 20, SKU: "SKU-020", ProductID: "product-020"}

	cases := []struct {
		name string
		item domain.ReplenishmentReadinessItem
		want bool
	}{
		{name: "higher risk rank comes after", item: domain.ReplenishmentReadinessItem{Risk: domain.ReplenishmentRiskWatch, SuggestedQuantity: 99, SKU: "SKU-001", ProductID: "product-001"}, want: true},
		{name: "lower suggested quantity comes after", item: domain.ReplenishmentReadinessItem{Risk: domain.ReplenishmentRiskReorder, SuggestedQuantity: 10, SKU: "SKU-001", ProductID: "product-001"}, want: true},
		{name: "higher suggested quantity comes before", item: domain.ReplenishmentReadinessItem{Risk: domain.ReplenishmentRiskReorder, SuggestedQuantity: 30, SKU: "SKU-999", ProductID: "product-999"}, want: false},
		{name: "later SKU tie break", item: domain.ReplenishmentReadinessItem{Risk: domain.ReplenishmentRiskReorder, SuggestedQuantity: 20, SKU: "SKU-021", ProductID: "product-021"}, want: true},
		{name: "same row is not after", item: domain.ReplenishmentReadinessItem{Risk: domain.ReplenishmentRiskReorder, SuggestedQuantity: 20, SKU: "SKU-020", ProductID: "product-020"}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := afterReplenishmentCursor(tc.item, cursor); got != tc.want {
				t.Fatalf("after cursor=%v want %v", got, tc.want)
			}
		})
	}
}

package httpapi

import (
	"testing"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func TestClassifyReplenishmentRisk(t *testing.T) {
	tests := []struct {
		name           string
		onHand         int64
		reorderPoint   int64
		averageDaily   float64
		want           domain.ReplenishmentRisk
	}{
		{"out of stock", 0, 10, 2, domain.ReplenishmentRiskOutOfStock},
		{"critical cover", 10, 5, 2, domain.ReplenishmentRiskCritical},
		{"reorder point", 10, 10, 1, domain.ReplenishmentRiskReorder},
		{"watch cover", 20, 5, 2, domain.ReplenishmentRiskWatch},
		{"healthy", 40, 10, 2, domain.ReplenishmentRiskHealthy},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyReplenishmentRisk(tt.onHand, tt.reorderPoint, tt.averageDaily); got != tt.want {
				t.Fatalf("risk = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeReplenishmentBounds(t *testing.T) {
	if got := normalizeReplenishmentWindow(0); got != 30 { t.Fatalf("default days = %d", got) }
	if got := normalizeReplenishmentWindow(999); got != 365 { t.Fatalf("max days = %d", got) }
	if got := normalizeReplenishmentLimit(0); got != 1000 { t.Fatalf("default limit = %d", got) }
	if got := normalizeReplenishmentLimit(99999); got != 5000 { t.Fatalf("max limit = %d", got) }
}

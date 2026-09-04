package domain

import "testing"

func TestReplenishmentRiskValues(t *testing.T) {
	if ReplenishmentRiskOutOfStock == ReplenishmentRiskHealthy {
		t.Fatal("risk values must remain distinct")
	}
	if ReplenishmentRiskCritical == ReplenishmentRiskReorder {
		t.Fatal("critical and reorder risks must remain distinct")
	}
}

func TestReplenishmentReadinessReportContract(t *testing.T) {
	report := ReplenishmentReadinessReport{
		WindowDays: 30,
		Items: []ReplenishmentReadinessItem{{
			ProductID:         "p1",
			SKU:               "SKU-1",
			OnHand:            12,
			ReorderPoint:      10,
			SuggestedQuantity: 20,
			AverageDailyOut:   2,
			Risk:              ReplenishmentRiskWatch,
		}},
	}
	if report.WindowDays != 30 || len(report.Items) != 1 || report.Items[0].SuggestedQuantity != 20 {
		t.Fatal("unexpected replenishment readiness contract")
	}
}

package domain

import (
	"math"
	"testing"
	"time"
)

func TestStockMovementHistoryAverageDailyOutbound(t *testing.T) {
	item := StockMovementHistoryItem{OutboundUnits: 91}
	windowDays := 30
	average := float64(item.OutboundUnits) / float64(windowDays)
	if math.Abs(average-3.0333333333) > 0.0000001 {
		t.Fatalf("average = %v", average)
	}
}

func TestStockMovementHistoryReportContract(t *testing.T) {
	now := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	report := StockMovementHistoryReport{
		AsOf:       now,
		WindowDays: 30,
		Items: []StockMovementHistoryItem{{
			ProductID:      "p1",
			LocationID:     "l1",
			MovementCount:  2,
			InboundUnits:   10,
			OutboundUnits:  4,
			NetUnits:       6,
			LastMovementAt: now,
		}},
	}
	if report.WindowDays != 30 || len(report.Items) != 1 || report.Items[0].NetUnits != 6 {
		t.Fatal("unexpected movement history contract")
	}
}

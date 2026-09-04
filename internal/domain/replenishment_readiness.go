package domain

import "time"

// ReplenishmentRisk classifies how urgently a product should be reviewed for replenishment.
type ReplenishmentRisk string

const (
	ReplenishmentRiskOutOfStock ReplenishmentRisk = "out_of_stock"
	ReplenishmentRiskCritical   ReplenishmentRisk = "critical"
	ReplenishmentRiskReorder    ReplenishmentRisk = "reorder"
	ReplenishmentRiskWatch      ReplenishmentRisk = "watch"
	ReplenishmentRiskHealthy    ReplenishmentRisk = "healthy"
)

// ReplenishmentReadinessItem combines reorder policy with recent outbound velocity.
type ReplenishmentReadinessItem struct {
	ProductID         string            `json:"productId"`
	SKU               string            `json:"sku"`
	Name              string            `json:"name"`
	SupplierID        string            `json:"supplierId,omitempty"`
	Unit              string            `json:"unit"`
	OnHand            int64             `json:"onHand"`
	ReorderPoint      int64             `json:"reorderPoint"`
	ReorderQuantity   int64             `json:"reorderQuantity"`
	TargetStock       int64             `json:"targetStock"`
	SuggestedQuantity int64             `json:"suggestedQuantity"`
	OutboundUnits     int64             `json:"outboundUnits"`
	AverageDailyOut   float64           `json:"averageDailyOutbound"`
	DaysOfCover       *float64          `json:"daysOfCover,omitempty"`
	Risk              ReplenishmentRisk `json:"risk"`
}

type ReplenishmentReadinessReport struct {
	AsOf       time.Time                    `json:"asOf"`
	WindowDays int                          `json:"windowDays"`
	Items      []ReplenishmentReadinessItem `json:"items"`
}

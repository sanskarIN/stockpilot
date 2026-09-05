package domain

import "time"

// StockMovementHistoryItem aggregates movement activity for a product at a location/lot.
type StockMovementHistoryItem struct {
	ProductID       string    `json:"productId"`
	SKU             string    `json:"sku"`
	Name            string    `json:"name"`
	LocationID      string    `json:"locationId"`
	LotID           string    `json:"lotId,omitempty"`
	MovementCount   int64     `json:"movementCount"`
	InboundUnits    int64     `json:"inboundUnits"`
	OutboundUnits   int64     `json:"outboundUnits"`
	NetUnits        int64     `json:"netUnits"`
	AverageDailyOut float64   `json:"averageDailyOutbound"`
	LastMovementAt  time.Time `json:"lastMovementAt"`
}

type StockMovementHistoryReport struct {
	AsOf       time.Time                 `json:"asOf"`
	WindowDays int                       `json:"windowDays"`
	Items      []StockMovementHistoryItem `json:"items"`
}

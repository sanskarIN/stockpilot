package domain

import "time"

type ReplenishmentPerformanceItem struct {
	SupplierID       string  `json:"supplierId"`
	SupplierName     string  `json:"supplierName"`
	OrderCount       int64   `json:"orderCount"`
	OrderedUnits     int64   `json:"orderedUnits"`
	ReceivedUnits    int64   `json:"receivedUnits"`
	OutstandingUnits int64   `json:"outstandingUnits"`
	FillRate         float64 `json:"fillRate"`
	OnTimeOrderCount int64   `json:"onTimeOrderCount"`
	LateOrderCount   int64   `json:"lateOrderCount"`
	AverageLeadDays  float64 `json:"averageLeadDays"`
}

type ReplenishmentPerformanceReport struct {
	AsOf       time.Time                      `json:"asOf"`
	WindowDays int                            `json:"windowDays"`
	Items      []ReplenishmentPerformanceItem `json:"items"`
}

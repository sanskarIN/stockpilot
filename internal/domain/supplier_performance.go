package domain

import "time"

// SupplierPerformanceItem summarizes purchasing activity and observed receipt lead time.
type SupplierPerformanceItem struct {
	SupplierID           string  `json:"supplierId"`
	SupplierCode         string  `json:"supplierCode"`
	SupplierName         string  `json:"supplierName"`
	OrderCount           int64   `json:"orderCount"`
	OrderedUnits         int64   `json:"orderedUnits"`
	ReceivedUnits        int64   `json:"receivedUnits"`
	OpenUnits            int64   `json:"openUnits"`
	OrderedValueMinor    int64   `json:"orderedValueMinor"`
	ReceivedValueMinor   int64   `json:"receivedValueMinor"`
	AverageLeadTimeDays  float64 `json:"averageLeadTimeDays"`
	CompletedOrderCount  int64   `json:"completedOrderCount"`
	OnTimeOrderCount     int64   `json:"onTimeOrderCount"`
}

type SupplierPerformanceReport struct {
	AsOf       time.Time                  `json:"asOf"`
	WindowDays int                        `json:"windowDays"`
	Items      []SupplierPerformanceItem `json:"items"`
}

package domain

// WarehouseValuationItem represents the on-hand inventory value at a location.
// Values are grouped by currency so no implicit currency conversion is performed.
type WarehouseValuationItem struct {
	WarehouseID    string `json:"warehouseId"`
	WarehouseCode  string `json:"warehouseCode"`
	WarehouseName  string `json:"warehouseName"`
	LocationID     string `json:"locationId"`
	LocationCode   string `json:"locationCode"`
	LocationName   string `json:"locationName"`
	Currency       string `json:"currency"`
	OnHand         int64  `json:"onHand"`
	ValuationMinor int64  `json:"valuationMinor"`
	ProductCount   int64  `json:"productCount"`
}

type WarehouseValuationTotal struct {
	WarehouseID    string `json:"warehouseId"`
	WarehouseCode  string `json:"warehouseCode"`
	WarehouseName  string `json:"warehouseName"`
	Currency       string `json:"currency"`
	OnHand         int64  `json:"onHand"`
	ValuationMinor int64  `json:"valuationMinor"`
}

type WarehouseValuationReport struct {
	Items  []WarehouseValuationItem  `json:"items"`
	Totals []WarehouseValuationTotal `json:"totals"`
}

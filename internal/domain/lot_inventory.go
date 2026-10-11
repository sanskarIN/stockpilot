package domain

import "time"

type LotInventoryRow struct {
	ProductID   string     `json:"productId"`
	SKU         string     `json:"sku"`
	ProductName string     `json:"productName"`
	LotID       string     `json:"lotId"`
	LotNumber   string     `json:"lotNumber"`
	LocationID  string     `json:"locationId"`
	Location    string     `json:"location"`
	WarehouseID string     `json:"warehouseId"`
	Warehouse   string     `json:"warehouse"`
	OnHand      int64      `json:"onHand"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	Active      bool       `json:"active"`
}

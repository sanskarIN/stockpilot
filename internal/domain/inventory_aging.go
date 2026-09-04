package domain

import "time"

// InventoryAgingBucket identifies the age range of inventory held at a location/lot.
type InventoryAgingBucket string

const (
	Aging0To30   InventoryAgingBucket = "0-30"
	Aging31To60  InventoryAgingBucket = "31-60"
	Aging61To90  InventoryAgingBucket = "61-90"
	Aging91To180 InventoryAgingBucket = "91-180"
	Aging181Plus InventoryAgingBucket = "181+"
)

type InventoryAgingItem struct {
	ProductID      string               `json:"productId"`
	SKU            string               `json:"sku"`
	Name           string               `json:"name"`
	LocationID     string               `json:"locationId"`
	LotID          string               `json:"lotId,omitempty"`
	Quantity       int64                `json:"quantity"`
	AgeDays        int64                `json:"ageDays"`
	Bucket         InventoryAgingBucket `json:"bucket"`
	AsOf           time.Time            `json:"asOf"`
	LastMovementAt time.Time            `json:"lastMovementAt"`
}

type InventoryAgingReport struct {
	Items []InventoryAgingItem `json:"items"`
}

func AgingBucket(ageDays int64) InventoryAgingBucket {
	switch {
	case ageDays <= 30:
		return Aging0To30
	case ageDays <= 60:
		return Aging31To60
	case ageDays <= 90:
		return Aging61To90
	case ageDays <= 180:
		return Aging91To180
	default:
		return Aging181Plus
	}
}

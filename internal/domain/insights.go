package domain

import "time"

type InventorySummary struct {
	ProductCount         int64 `json:"productCount"`
	ActiveProductCount   int64 `json:"activeProductCount"`
	ActiveWarehouseCount int64 `json:"activeWarehouseCount"`
	ActiveLocationCount  int64 `json:"activeLocationCount"`
	TotalUnits           int64 `json:"totalUnits"`
	LowStockBalanceCount int64 `json:"lowStockBalanceCount"`
	OutOfStockCount      int64 `json:"outOfStockCount"`
}

type PurchasingSummary struct {
	TotalOrders             int64 `json:"totalOrders"`
	DraftOrders             int64 `json:"draftOrders"`
	OrderedOrders           int64 `json:"orderedOrders"`
	PartiallyReceivedOrders int64 `json:"partiallyReceivedOrders"`
	ReceivedOrders          int64 `json:"receivedOrders"`
	CancelledOrders         int64 `json:"cancelledOrders"`
	OutstandingUnits        int64 `json:"outstandingUnits"`
}

type AuditEvent struct {
	ID         int64          `json:"id"`
	OccurredAt time.Time      `json:"occurredAt"`
	ActorID    string         `json:"actorId,omitempty"`
	Action     string         `json:"action"`
	EntityType string         `json:"entityType"`
	EntityID   string         `json:"entityId"`
	RequestID  string         `json:"requestId,omitempty"`
	Metadata   map[string]any `json:"metadata"`
}

type AuditFilter struct {
	ActorID    string
	Action     string
	EntityType string
	EntityID   string
	Limit      int
	Offset     int
}

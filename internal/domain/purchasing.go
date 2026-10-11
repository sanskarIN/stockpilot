package domain

import (
	"fmt"
	"strings"
	"time"
)

type PurchaseOrderStatus string

const (
	PurchaseOrderDraft             PurchaseOrderStatus = "draft"
	PurchaseOrderOrdered           PurchaseOrderStatus = "ordered"
	PurchaseOrderPartiallyReceived PurchaseOrderStatus = "partially_received"
	PurchaseOrderReceived          PurchaseOrderStatus = "received"
	PurchaseOrderCancelled         PurchaseOrderStatus = "cancelled"
)

func (s PurchaseOrderStatus) Valid() bool {
	switch s {
	case PurchaseOrderDraft, PurchaseOrderOrdered, PurchaseOrderPartiallyReceived, PurchaseOrderReceived, PurchaseOrderCancelled:
		return true
	default:
		return false
	}
}

type PurchaseOrderLine struct {
	ID              string `json:"id"`
	PurchaseOrderID string `json:"purchaseOrderId"`
	ProductID       string `json:"productId"`
	Quantity        int64  `json:"quantity"`
	Received        int64  `json:"received"`
	UnitCostMinor   int64  `json:"unitCostMinor"`
}

func (l PurchaseOrderLine) Validate() error {
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.ProductID) == "" {
		return fmt.Errorf("%w: purchase order line id and product id are required", ErrInvalid)
	}
	if l.Quantity <= 0 {
		return fmt.Errorf("%w: ordered quantity must be positive", ErrInvalid)
	}
	if l.Received < 0 || l.Received > l.Quantity {
		return fmt.Errorf("%w: received quantity must be between zero and ordered quantity", ErrInvalid)
	}
	if l.UnitCostMinor < 0 {
		return fmt.Errorf("%w: unit cost cannot be negative", ErrInvalid)
	}
	return nil
}

type PurchaseOrder struct {
	ID          string              `json:"id"`
	Number      string              `json:"number"`
	SupplierID  string              `json:"supplierId"`
	WarehouseID string              `json:"warehouseId"`
	Status      PurchaseOrderStatus `json:"status"`
	Currency    string              `json:"currency"`
	ExpectedAt  *time.Time          `json:"expectedAt,omitempty"`
	Notes       string              `json:"notes,omitempty"`
	Lines       []PurchaseOrderLine `json:"lines"`
	CreatedBy   string              `json:"createdBy,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

func (p PurchaseOrder) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("%w: purchase order id is required", ErrInvalid)
	}
	if n := strings.TrimSpace(p.Number); len(n) < 2 || len(n) > 64 {
		return fmt.Errorf("%w: purchase order number must be 2-64 characters", ErrInvalid)
	}
	if strings.TrimSpace(p.SupplierID) == "" || strings.TrimSpace(p.WarehouseID) == "" {
		return fmt.Errorf("%w: supplier and warehouse are required", ErrInvalid)
	}
	if !p.Status.Valid() {
		return fmt.Errorf("%w: unsupported purchase order status", ErrInvalid)
	}
	if len(strings.TrimSpace(p.Currency)) != 3 {
		return fmt.Errorf("%w: currency must be a 3-letter code", ErrInvalid)
	}
	if len(p.Lines) == 0 {
		return fmt.Errorf("%w: purchase order requires at least one line", ErrInvalid)
	}
	seen := make(map[string]struct{}, len(p.Lines))
	for _, line := range p.Lines {
		if err := line.Validate(); err != nil {
			return err
		}
		if _, exists := seen[line.ProductID]; exists {
			return fmt.Errorf("%w: product appears more than once in purchase order", ErrInvalid)
		}
		seen[line.ProductID] = struct{}{}
	}
	return nil
}

func (p PurchaseOrder) TotalMinor() int64 {
	var total int64
	for _, line := range p.Lines {
		total += line.Quantity * line.UnitCostMinor
	}
	return total
}

func (p PurchaseOrder) ReceivedStatus() PurchaseOrderStatus {
	var ordered, received int64
	for _, line := range p.Lines {
		ordered += line.Quantity
		received += line.Received
	}
	if received == 0 {
		return p.Status
	}
	if received >= ordered {
		return PurchaseOrderReceived
	}
	return PurchaseOrderPartiallyReceived
}

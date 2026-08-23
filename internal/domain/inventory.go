package domain

import (
	"fmt"
	"strings"
	"time"
)

type MovementType string

const (
	MovementStockIn     MovementType = "stock_in"
	MovementStockOut    MovementType = "stock_out"
	MovementAdjustment  MovementType = "adjustment"
	MovementTransferIn  MovementType = "transfer_in"
	MovementTransferOut MovementType = "transfer_out"
	MovementReceive     MovementType = "receive"
)

func (t MovementType) Valid() bool {
	switch t {
	case MovementStockIn, MovementStockOut, MovementAdjustment, MovementTransferIn, MovementTransferOut, MovementReceive:
		return true
	default:
		return false
	}
}

type Lot struct {
	ID           string     `json:"id"`
	ProductID    string     `json:"productId"`
	LotNumber    string     `json:"lotNumber"`
	Manufactured *time.Time `json:"manufacturedAt,omitempty"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (l Lot) Validate() error {
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.ProductID) == "" {
		return fmt.Errorf("%w: lot id and product id are required", ErrInvalid)
	}
	if n := strings.TrimSpace(l.LotNumber); len(n) < 1 || len(n) > 96 {
		return fmt.Errorf("%w: lot number must be 1-96 characters", ErrInvalid)
	}
	if l.Manufactured != nil && l.ExpiresAt != nil && !l.ExpiresAt.After(*l.Manufactured) {
		return fmt.Errorf("%w: expiry must be after manufacturing time", ErrInvalid)
	}
	return nil
}

type StockBalance struct {
	ProductID  string    `json:"productId"`
	LocationID string    `json:"locationId"`
	LotID      string    `json:"lotId,omitempty"`
	Quantity   int64     `json:"quantity"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type StockMovement struct {
	ID            string       `json:"id"`
	ProductID     string       `json:"productId"`
	LocationID    string       `json:"locationId"`
	LotID         string       `json:"lotId,omitempty"`
	Type          MovementType `json:"type"`
	QuantityDelta int64        `json:"quantityDelta"`
	Reference     string       `json:"reference,omitempty"`
	Note          string       `json:"note,omitempty"`
	ActorID       string       `json:"actorId,omitempty"`
	OccurredAt    time.Time    `json:"occurredAt"`
	CreatedAt     time.Time    `json:"createdAt"`
}

func (m StockMovement) Validate() error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.ProductID) == "" || strings.TrimSpace(m.LocationID) == "" {
		return fmt.Errorf("%w: movement id, product id, and location id are required", ErrInvalid)
	}
	if !m.Type.Valid() {
		return fmt.Errorf("%w: unsupported movement type", ErrInvalid)
	}
	if m.QuantityDelta == 0 {
		return fmt.Errorf("%w: movement quantity delta cannot be zero", ErrInvalid)
	}
	switch m.Type {
	case MovementStockIn, MovementTransferIn, MovementReceive:
		if m.QuantityDelta < 0 {
			return fmt.Errorf("%w: inbound movement quantity must be positive", ErrInvalid)
		}
	case MovementStockOut, MovementTransferOut:
		if m.QuantityDelta > 0 {
			return fmt.Errorf("%w: outbound movement quantity must be negative", ErrInvalid)
		}
	}
	if m.OccurredAt.IsZero() {
		return fmt.Errorf("%w: movement occurrence time is required", ErrInvalid)
	}
	return nil
}

type TransferRequest struct {
	ID             string    `json:"id"`
	ProductID      string    `json:"productId"`
	FromLocationID string    `json:"fromLocationId"`
	ToLocationID   string    `json:"toLocationId"`
	LotID          string    `json:"lotId,omitempty"`
	Quantity       int64     `json:"quantity"`
	Reference      string    `json:"reference,omitempty"`
	Note           string    `json:"note,omitempty"`
	ActorID        string    `json:"actorId,omitempty"`
	OccurredAt     time.Time `json:"occurredAt"`
}

func (r TransferRequest) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.ProductID) == "" {
		return fmt.Errorf("%w: transfer id and product id are required", ErrInvalid)
	}
	if strings.TrimSpace(r.FromLocationID) == "" || strings.TrimSpace(r.ToLocationID) == "" {
		return fmt.Errorf("%w: transfer source and destination are required", ErrInvalid)
	}
	if r.FromLocationID == r.ToLocationID {
		return fmt.Errorf("%w: transfer locations must differ", ErrInvalid)
	}
	if r.Quantity <= 0 {
		return fmt.Errorf("%w: transfer quantity must be positive", ErrInvalid)
	}
	if r.OccurredAt.IsZero() {
		return fmt.Errorf("%w: transfer occurrence time is required", ErrInvalid)
	}
	return nil
}

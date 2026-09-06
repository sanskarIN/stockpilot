package domain

import (
	"fmt"
	"strings"
	"time"
)

type ReplenishmentReviewOutcome string

const (
	ReplenishmentReviewAccepted  ReplenishmentReviewOutcome = "accepted"
	ReplenishmentReviewModified  ReplenishmentReviewOutcome = "modified"
	ReplenishmentReviewDismissed ReplenishmentReviewOutcome = "dismissed"
	ReplenishmentReviewExpired   ReplenishmentReviewOutcome = "expired"
)

func (o ReplenishmentReviewOutcome) Valid() bool {
	switch o {
	case ReplenishmentReviewAccepted, ReplenishmentReviewModified, ReplenishmentReviewDismissed, ReplenishmentReviewExpired:
		return true
	default:
		return false
	}
}

type ReplenishmentReview struct {
	ID                string                     `json:"id"`
	ProductID         string                     `json:"productId"`
	SKU               string                     `json:"sku"`
	ProductName       string                     `json:"productName"`
	Unit              string                     `json:"unit"`
	SupplierID        string                     `json:"supplierId,omitempty"`
	OnHand            int64                      `json:"onHand"`
	ReorderPoint      int64                      `json:"reorderPoint"`
	ReorderQuantity   int64                      `json:"reorderQuantity"`
	TargetStock       int64                      `json:"targetStock"`
	SuggestedQuantity int64                      `json:"suggestedQuantity"`
	ReviewedBy        string                     `json:"reviewedBy"`
	ReviewedAt        time.Time                  `json:"reviewedAt"`
	PurchaseOrderID   string                     `json:"purchaseOrderId,omitempty"`
	Outcome           ReplenishmentReviewOutcome `json:"outcome"`
	CreatedAt         time.Time                  `json:"createdAt"`
}

func (r ReplenishmentReview) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.ProductID) == "" {
		return fmt.Errorf("%w: review id and product id are required", ErrInvalid)
	}
	if strings.TrimSpace(r.SKU) == "" || strings.TrimSpace(r.ProductName) == "" || strings.TrimSpace(r.Unit) == "" {
		return fmt.Errorf("%w: product snapshot fields are required", ErrInvalid)
	}
	if r.OnHand < 0 || r.ReorderPoint < 0 || r.ReorderQuantity < 0 || r.TargetStock < 0 || r.SuggestedQuantity <= 0 {
		return fmt.Errorf("%w: replenishment quantities are invalid", ErrInvalid)
	}
	if strings.TrimSpace(r.ReviewedBy) == "" || r.ReviewedAt.IsZero() {
		return fmt.Errorf("%w: reviewer and review time are required", ErrInvalid)
	}
	if !r.Outcome.Valid() {
		return fmt.Errorf("%w: unsupported replenishment review outcome", ErrInvalid)
	}
	if r.Outcome == ReplenishmentReviewAccepted || r.Outcome == ReplenishmentReviewModified {
		if strings.TrimSpace(r.PurchaseOrderID) == "" {
			return fmt.Errorf("%w: purchase order is required for %s outcome", ErrInvalid, r.Outcome)
		}
	} else if strings.TrimSpace(r.PurchaseOrderID) != "" {
		return fmt.Errorf("%w: purchase order is not allowed for %s outcome", ErrInvalid, r.Outcome)
	}
	return nil
}

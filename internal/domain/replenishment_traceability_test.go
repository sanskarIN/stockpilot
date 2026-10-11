package domain

import (
	"testing"
	"time"
)

func validReplenishmentReview() ReplenishmentReview {
	return ReplenishmentReview{
		ID: "rr_test", ProductID: "prd_test", SKU: "SKU-TEST", ProductName: "Test product", Unit: "piece",
		OnHand: 4, ReorderPoint: 10, ReorderQuantity: 20, TargetStock: 30, SuggestedQuantity: 26,
		ReviewedBy: "user_test", ReviewedAt: time.Unix(100, 0).UTC(), Outcome: ReplenishmentReviewDismissed,
	}
}

func TestReplenishmentReviewValidateOutcomes(t *testing.T) {
	tests := []struct {
		name          string
		outcome       ReplenishmentReviewOutcome
		purchaseOrder string
		wantErr       bool
	}{
		{name: "accepted requires order", outcome: ReplenishmentReviewAccepted, wantErr: true},
		{name: "accepted with order", outcome: ReplenishmentReviewAccepted, purchaseOrder: "po_test"},
		{name: "modified requires order", outcome: ReplenishmentReviewModified, wantErr: true},
		{name: "modified with order", outcome: ReplenishmentReviewModified, purchaseOrder: "po_test"},
		{name: "dismissed without order", outcome: ReplenishmentReviewDismissed},
		{name: "dismissed rejects order", outcome: ReplenishmentReviewDismissed, purchaseOrder: "po_test", wantErr: true},
		{name: "expired without order", outcome: ReplenishmentReviewExpired},
		{name: "invalid outcome", outcome: ReplenishmentReviewOutcome("unknown"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			review := validReplenishmentReview()
			review.Outcome = tt.outcome
			review.PurchaseOrderID = tt.purchaseOrder
			if err := review.Validate(); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package domain

import "testing"

func TestValidatePurchaseOrderTransition(t *testing.T) {
	tests := []struct {
		name string
		from PurchaseOrderStatus
		to   PurchaseOrderStatus
		want bool
	}{
		{"draft to ordered", PurchaseOrderDraft, PurchaseOrderOrdered, true},
		{"draft to cancelled", PurchaseOrderDraft, PurchaseOrderCancelled, true},
		{"ordered to cancelled", PurchaseOrderOrdered, PurchaseOrderCancelled, true},
		{"draft to received", PurchaseOrderDraft, PurchaseOrderReceived, false},
		{"ordered to received", PurchaseOrderOrdered, PurchaseOrderReceived, false},
		{"partial to cancelled", PurchaseOrderPartiallyReceived, PurchaseOrderCancelled, false},
		{"received to cancelled", PurchaseOrderReceived, PurchaseOrderCancelled, false},
		{"same status", PurchaseOrderDraft, PurchaseOrderDraft, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePurchaseOrderTransition(tt.from, tt.to)
			if (err == nil) != tt.want { t.Fatalf("error = %v, want success=%v", err, tt.want) }
		})
	}
}

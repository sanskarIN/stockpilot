package domain

import (
	"errors"
	"testing"
)

func TestPurchaseOrderValidateRejectsDuplicateProducts(t *testing.T) {
	order := PurchaseOrder{
		ID: "po_1", Number: "PO-001", SupplierID: "sup_1", WarehouseID: "wh_1", Status: PurchaseOrderDraft, Currency: "INR",
		Lines: []PurchaseOrderLine{
			{ID: "line_1", ProductID: "prd_1", Quantity: 2, UnitCostMinor: 1000},
			{ID: "line_2", ProductID: "prd_1", Quantity: 3, UnitCostMinor: 1000},
		},
	}
	if err := order.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Validate() error = %v, want ErrInvalid", err)
	}
}

func TestPurchaseOrderTotalMinor(t *testing.T) {
	order := PurchaseOrder{Lines: []PurchaseOrderLine{
		{Quantity: 2, UnitCostMinor: 1500},
		{Quantity: 3, UnitCostMinor: 2500},
	}}
	if got, want := order.TotalMinor(), int64(10500); got != want {
		t.Fatalf("TotalMinor() = %d, want %d", got, want)
	}
}

func TestPurchaseOrderReceivedStatus(t *testing.T) {
	tests := []struct {
		name   string
		lines  []PurchaseOrderLine
		status PurchaseOrderStatus
		want   PurchaseOrderStatus
	}{
		{
			name: "nothing received keeps current status",
			lines: []PurchaseOrderLine{{Quantity: 2, Received: 0}},
			status: PurchaseOrderOrdered,
			want: PurchaseOrderOrdered,
		},
		{
			name: "partial receipt",
			lines: []PurchaseOrderLine{{Quantity: 2, Received: 1}, {Quantity: 2, Received: 0}},
			status: PurchaseOrderOrdered,
			want: PurchaseOrderPartiallyReceived,
		},
		{
			name: "complete receipt",
			lines: []PurchaseOrderLine{{Quantity: 2, Received: 2}, {Quantity: 3, Received: 3}},
			status: PurchaseOrderPartiallyReceived,
			want: PurchaseOrderReceived,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := PurchaseOrder{Lines: tt.lines, Status: tt.status}
			if got := order.ReceivedStatus(); got != tt.want {
				t.Fatalf("ReceivedStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

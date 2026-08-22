package domain

import "testing"

func TestProductInputValidation(t *testing.T) {
	active := true
	valid := ProductInput{SKU: " abc-1 ", Name: "Widget", Unit: "PCS", ReorderPoint: 5, CostCents: 100, PriceCents: 150, Active: &active}
	if err := valid.NormalizeAndValidate(); err != nil {
		t.Fatalf("expected valid product, got %v", err)
	}
	if valid.SKU != "ABC-1" || valid.Unit != "pcs" {
		t.Fatalf("normalization failed: %#v", valid)
	}

	invalid := ProductInput{SKU: "A", Name: "X", TrackExpiry: true}
	if err := invalid.NormalizeAndValidate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestTransferRejectsSameLocation(t *testing.T) {
	in := TransferInput{ProductID: "p1", FromLocationID: "l1", ToLocationID: "l1", Quantity: 1}
	if err := in.NormalizeAndValidate(); err == nil {
		t.Fatal("expected same-location transfer rejection")
	}
}

func TestRoleAllows(t *testing.T) {
	if !RoleAllows(RoleOwner, RoleAdmin) {
		t.Fatal("owner should satisfy admin requirement")
	}
	if RoleAllows(RoleViewer, RoleManager) {
		t.Fatal("viewer should not satisfy manager requirement")
	}
}

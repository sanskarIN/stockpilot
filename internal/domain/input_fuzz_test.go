package domain

import "testing"

func FuzzProductInputNeverPanics(f *testing.F) {
	f.Add("SKU-1", "Widget", "unit", 1.0, int64(100), int64(200), false, false)
	f.Add("x", "y", "", -1.0, int64(-1), int64(0), true, true)
	f.Fuzz(func(t *testing.T, sku, name, unit string, reorder float64, cost, price int64, lots, expiry bool) {
		in := ProductInput{SKU: sku, Name: name, Unit: unit, ReorderPoint: reorder, CostCents: cost, PriceCents: price, TrackLots: lots, TrackExpiry: expiry}
		_ = in.NormalizeAndValidate()
	})
}

func FuzzTransferValidationNeverPanics(f *testing.F) {
	f.Add("p", "from", "to", 1.0, "lot")
	f.Fuzz(func(t *testing.T, product, from, to string, qty float64, lot string) {
		in := TransferInput{ProductID: product, FromLocationID: from, ToLocationID: to, Quantity: qty, LotNumber: lot}
		_ = in.NormalizeAndValidate()
	})
}

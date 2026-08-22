package domain

import "testing"

func BenchmarkProductValidation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		active := true
		in := ProductInput{SKU: "SKU-12345", Name: "Benchmark Widget", Barcode: "8901234567890", Unit: "pcs", ReorderPoint: 12.5, CostCents: 1299, PriceCents: 1999, Active: &active}
		if err := in.NormalizeAndValidate(); err != nil {
			b.Fatal(err)
		}
	}
}

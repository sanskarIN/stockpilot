package csvimport

import (
	"strings"
	"testing"
)

const validHeader = "sku,name,description,unit_cost_minor,currency,reorder_point,reorder_quantity,track_lots,track_expiry,active\n"

func TestParseProductsValidRows(t *testing.T) {
	input := validHeader + "SKU-1,Widget,Demo,1250,INR,10,20,false,false,true\n"
	result, err := ParseProducts(strings.NewReader(input))
	if err != nil { t.Fatalf("ParseProducts() error = %v", err) }
	if len(result.Errors) != 0 { t.Fatalf("unexpected errors: %+v", result.Errors) }
	if len(result.Rows) != 1 { t.Fatalf("rows = %d, want 1", len(result.Rows)) }
	if result.Rows[0].Product.SKU != "SKU-1" || result.Rows[0].Product.UnitCostMinor != 1250 { t.Fatalf("unexpected product: %+v", result.Rows[0].Product) }
}

func TestParseProductsReportsFieldErrors(t *testing.T) {
	input := validHeader + "S,Widget,Demo,-1,INR,nope,20,true,false,true\n"
	result, err := ParseProducts(strings.NewReader(input))
	if err != nil { t.Fatalf("ParseProducts() error = %v", err) }
	if len(result.Rows) != 0 { t.Fatalf("rows = %d, want 0", len(result.Rows)) }
	if len(result.Errors) < 2 { t.Fatalf("errors = %+v, want multiple validation errors", result.Errors) }
}

func TestParseProductsRejectsDuplicateSKUAndBarcode(t *testing.T) {
	input := "sku,name,unit,unit_cost_minor,currency,reorder_point,reorder_quantity,track_lots,track_expiry,active,barcode\nSKU-1,A,each,100,INR,1,1,false,false,true,890\nSKU-1,B,each,200,INR,1,1,false,false,true,890\n"
	result, err := ParseProducts(strings.NewReader(input))
	if err != nil { t.Fatalf("ParseProducts() error = %v", err) }
	if len(result.Rows) != 1 { t.Fatalf("rows = %d, want 1", len(result.Rows)) }
	if len(result.Errors) != 2 { t.Fatalf("errors = %+v, want duplicate SKU and barcode", result.Errors) }
}

func TestParseProductsRequiresHeaders(t *testing.T) {
	_, err := ParseProducts(strings.NewReader("sku,name\nSKU-1,Widget\n"))
	if err == nil { t.Fatal("expected required-header error") }
}

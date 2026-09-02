package csvimport

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

const MaxProductRows = 1000

type ProductRow struct {
	Row     int
	Product domain.Product
}

type RowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Rows    []ProductRow
	Errors  []RowError
	Headers []string
}

var productHeaders = []string{
	"id", "sku", "name", "description", "category_id", "supplier_id", "barcode", "unit",
	"unit_cost_minor", "currency", "reorder_point", "reorder_quantity", "track_lots", "track_expiry", "active",
}

func ParseProducts(r io.Reader) (ValidationResult, error) {
	reader := csv.NewReader(io.LimitReader(r, 4<<20))
	records, err := reader.Read()
	if err == io.EOF {
		return ValidationResult{}, fmt.Errorf("CSV file is empty")
	}
	if err != nil {
		return ValidationResult{}, fmt.Errorf("read CSV header: %w", err)
	}

	index, err := headerIndex(records)
	if err != nil {
		return ValidationResult{}, err
	}
	result := ValidationResult{Rows: make([]ProductRow, 0), Errors: make([]RowError, 0), Headers: append([]string(nil), records...)}
	seenSKU := make(map[string]int)
	seenBarcode := make(map[string]int)
	rowNumber := 1
	for {
		records, err = reader.Read()
		if err == io.EOF {
			break
		}
		rowNumber++
		if err != nil {
			result.Errors = append(result.Errors, RowError{Row: rowNumber, Message: "malformed CSV row"})
			continue
		}
		if rowNumber > MaxProductRows+1 {
			result.Errors = append(result.Errors, RowError{Row: rowNumber, Message: fmt.Sprintf("maximum of %d product rows is allowed", MaxProductRows)})
			for reader.Scan() { rowNumber++ }
			break
		}
		if blankRecord(records) {
			continue
		}
		product, rowErrors := parseProduct(index, records)
		if len(rowErrors) == 0 {
			sku := strings.ToUpper(strings.TrimSpace(product.SKU))
			if previous, ok := seenSKU[sku]; ok { rowErrors = append(rowErrors, fmt.Sprintf("SKU duplicates row %d", previous)) } else { seenSKU[sku] = rowNumber }
			barcode := strings.TrimSpace(product.Barcode)
			if barcode != "" {
				if previous, ok := seenBarcode[barcode]; ok { rowErrors = append(rowErrors, fmt.Sprintf("barcode duplicates row %d", previous)) } else { seenBarcode[barcode] = rowNumber }
			}
		}
		for _, message := range rowErrors { result.Errors = append(result.Errors, RowError{Row: rowNumber, Message: message}) }
		if len(rowErrors) == 0 { result.Rows = append(result.Rows, ProductRow{Row: rowNumber, Product: product}) }
	}
	return result, nil
}

func headerIndex(headers []string) (map[string]int, error) {
	result := make(map[string]int, len(headers))
	for i, raw := range headers {
		name := strings.ToLower(strings.TrimSpace(raw))
		if name == "" { return nil, fmt.Errorf("CSV header %d is empty", i+1) }
		if _, exists := result[name]; exists { return nil, fmt.Errorf("CSV header %q is duplicated", name) }
		result[name] = i
	}
	for _, required := range productHeaders[1:2] {
		if _, ok := result[required]; !ok { return nil, fmt.Errorf("CSV header %q is required", required) }
	}
	for _, required := range []string{"name", "unit", "unit_cost_minor", "currency", "reorder_point", "reorder_quantity", "track_lots", "track_expiry", "active"} {
		if _, ok := result[required]; !ok { return nil, fmt.Errorf("CSV header %q is required", required) }
	}
	return result, nil
}

func parseProduct(index map[string]int, values []string) (domain.Product, []string) {
	get := func(name string) string { if i, ok := index[name]; ok && i < len(values) { return strings.TrimSpace(values[i]) }; return "" }
	product := domain.Product{ID: get("id"), SKU: get("sku"), Name: get("name"), Description: get("description"), CategoryID: get("category_id"), SupplierID: get("supplier_id"), Barcode: get("barcode"), Unit: get("unit"), Currency: get("currency")}
	if product.ID == "" { product.ID = "pending" }
	errors := make([]string, 0, 4)
	parseInt := func(name string, target *int64) { value, err := strconv.ParseInt(get(name), 10, 64); if err != nil { errors = append(errors, fmt.Sprintf("%s must be a whole number", name)); return }; *target = value }
	parseBool := func(name string, target *bool) { value, err := strconv.ParseBool(get(name)); if err != nil { errors = append(errors, fmt.Sprintf("%s must be true or false", name)); return }; *target = value }
	parseInt("unit_cost_minor", &product.UnitCostMinor); parseInt("reorder_point", &product.ReorderPoint); parseInt("reorder_quantity", &product.ReorderQuantity)
	parseBool("track_lots", &product.TrackLots); parseBool("track_expiry", &product.TrackExpiry); parseBool("active", &product.Active)
	if err := product.Validate(); err != nil { errors = append(errors, err.Error()) }
	return product, errors
}

func blankRecord(values []string) bool { for _, value := range values { if strings.TrimSpace(value) != "" { return false } }; return true }

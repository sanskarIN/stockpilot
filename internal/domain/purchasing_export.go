package domain

import "time"

// PurchaseOrderExportRow is the flattened, stable CSV representation of one
// purchase-order line. Keeping the export shape separate from PurchaseOrder
// prevents presentation fields from leaking into the transactional model.
type PurchaseOrderExportRow struct {
	OrderID         string
	OrderNumber     string
	SupplierID      string
	WarehouseID     string
	Status          PurchaseOrderStatus
	Currency        string
	ExpectedAt      *time.Time
	Notes           string
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LineID          string
	ProductID       string
	Quantity        int64
	Received        int64
	UnitCostMinor   int64
}

func (r PurchaseOrderExportRow) LineTotalMinor() int64 {
	return r.Quantity * r.UnitCostMinor
}

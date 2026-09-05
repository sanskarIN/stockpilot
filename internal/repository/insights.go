package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
)

type Reports interface {
	InventorySummary(context.Context) (domain.InventorySummary, error)
	PurchasingSummary(context.Context) (domain.PurchasingSummary, error)
	SupplierPerformance(context.Context, int, int) (domain.SupplierPerformanceReport, error)
	WarehouseValuation(context.Context, int) (domain.WarehouseValuationReport, error)
}

// BoundedReports is an additive capability for repositories that can apply
// validated reporting periods and pagination at the storage boundary.
type BoundedReports interface {
	SupplierPerformanceQuery(context.Context, reporting.Query) (domain.SupplierPerformanceReport, error)
	WarehouseValuationQuery(context.Context, reporting.Query) (domain.WarehouseValuationReport, error)
}

// CountedReports is an additive capability for repositories that can report
// the complete number of rows represented by a bounded report query.
// It leaves existing report response bodies unchanged while enabling clients
// to reason about pagination without guessing from page length alone.
type CountedReports interface {
	SupplierPerformanceCount(context.Context, reporting.Query) (int64, error)
	WarehouseValuationCount(context.Context, reporting.Query) (int64, error)
}

type Audit interface {
	AppendAuditEvent(context.Context, domain.AuditEvent) error
	ListAuditEvents(context.Context, domain.AuditFilter) ([]domain.AuditEvent, error)
}

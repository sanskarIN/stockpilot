package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type Reports interface {
	InventorySummary(context.Context) (domain.InventorySummary, error)
	PurchasingSummary(context.Context) (domain.PurchasingSummary, error)
	SupplierPerformance(context.Context, int, int) (domain.SupplierPerformanceReport, error)
	WarehouseValuation(context.Context, int) (domain.WarehouseValuationReport, error)
}

type Audit interface {
	AppendAuditEvent(context.Context, domain.AuditEvent) error
	ListAuditEvents(context.Context, domain.AuditFilter) ([]domain.AuditEvent, error)
}

package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type Reports interface {
	InventorySummary(context.Context) (domain.InventorySummary, error)
	PurchasingSummary(context.Context) (domain.PurchasingSummary, error)
}

type Audit interface {
	ListAuditEvents(context.Context, domain.AuditFilter) ([]domain.AuditEvent, error)
}

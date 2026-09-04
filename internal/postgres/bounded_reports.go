package postgres

import (
	"context"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/reporting"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

// SupplierPerformanceQuery provides the additive bounded repository capability.
// The legacy query already enforces a storage-side limit; offset is applied to
// the bounded result to preserve the legacy repository interface.
func (s *Store) SupplierPerformanceQuery(ctx context.Context, query reporting.Query) (domain.SupplierPerformanceReport, error) {
	if err := query.Validate(); err != nil {
		return domain.SupplierPerformanceReport{}, err
	}
	windowDays := 30
	if query.From != nil && query.To != nil {
		days := int(query.To.Sub(*query.From).Hours()/24) + 1
		if days > 0 {
			windowDays = days
		}
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 100
	}
	report, err := s.SupplierPerformance(ctx, windowDays, limit+query.Offset)
	if err != nil {
		return domain.SupplierPerformanceReport{}, err
	}
	if query.Offset >= len(report.Items) {
		report.Items = []domain.SupplierPerformanceItem{}
	} else if query.Offset > 0 {
		report.Items = report.Items[query.Offset:]
	}
	if len(report.Items) > limit {
		report.Items = report.Items[:limit]
	}
	return report, nil
}

// WarehouseValuationQuery provides bounded pagination through the additive
// repository capability while preserving the existing report implementation.
func (s *Store) WarehouseValuationQuery(ctx context.Context, query reporting.Query) (domain.WarehouseValuationReport, error) {
	if err := query.Validate(); err != nil {
		return domain.WarehouseValuationReport{}, err
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 1000
	}
	report, err := s.WarehouseValuation(ctx, limit+query.Offset)
	if err != nil {
		return domain.WarehouseValuationReport{}, err
	}
	if query.Offset >= len(report.Items) {
		report.Items = []domain.WarehouseValuationItem{}
	} else if query.Offset > 0 {
		report.Items = report.Items[query.Offset:]
	}
	if len(report.Items) > limit {
		report.Items = report.Items[:limit]
	}
	return report, nil
}

var _ repository.BoundedReports = (*Store)(nil)

var _ = time.UTC

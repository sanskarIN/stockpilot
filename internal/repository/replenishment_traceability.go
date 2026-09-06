package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type ReplenishmentTraceability interface {
	CreateReplenishmentReview(context.Context, domain.ReplenishmentReview) error
	ListReplenishmentReviews(context.Context, string, int, int) ([]domain.ReplenishmentReview, error)
}

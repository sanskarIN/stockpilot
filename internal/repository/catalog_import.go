package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

// ProductBatchImporter persists a complete product import atomically.
// Implementations must revalidate every product and rely on database
// constraints as the final integrity boundary; callers must not treat a
// previous dry-run validation as an authorization to write.
type ProductBatchImporter interface {
	ImportProducts(context.Context, []domain.Product) error
}

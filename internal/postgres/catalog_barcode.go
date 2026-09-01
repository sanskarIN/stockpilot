package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (s *Store) GetProductByBarcode(ctx context.Context, barcode string) (domain.Product, error) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return domain.Product{}, fmt.Errorf("%w: barcode is required", domain.ErrInvalid)
	}
	if len(barcode) > 128 {
		return domain.Product{}, fmt.Errorf("%w: barcode must be at most 128 characters", domain.ErrInvalid)
	}
	return scanProduct(s.pool.QueryRow(ctx, productSelect+` WHERE p.barcode=$1`, barcode))
}

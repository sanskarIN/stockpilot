package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) CreateCategory(ctx context.Context, category domain.Category) error {
	if err := category.Validate(); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO categories (id, name, description)
		VALUES ($1, $2, $3)`,
		category.ID, strings.TrimSpace(category.Name), strings.TrimSpace(category.Description))
	return mapError(err)
}

func (s *Store) ListCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Category, 0)
	for rows.Next() {
		var item domain.Category
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CreateSupplier(ctx context.Context, supplier domain.Supplier) error {
	if err := supplier.Validate(); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO suppliers (id, code, name, email, phone, notes, active)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), $6, $7)`,
		supplier.ID, strings.TrimSpace(supplier.Code), strings.TrimSpace(supplier.Name), strings.TrimSpace(supplier.Email), strings.TrimSpace(supplier.Phone), strings.TrimSpace(supplier.Notes), supplier.Active)
	return mapError(err)
}

func (s *Store) ListSuppliers(ctx context.Context, activeOnly bool) ([]domain.Supplier, error) {
	query := `SELECT id, code, name, COALESCE(email, ''), COALESCE(phone, ''), notes, active, created_at, updated_at FROM suppliers`
	if activeOnly {
		query += ` WHERE active = true`
	}
	query += ` ORDER BY name`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Supplier, 0)
	for rows.Next() {
		var item domain.Supplier
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Email, &item.Phone, &item.Notes, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CreateProduct(ctx context.Context, product domain.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO products (
			id, sku, name, description, category_id, supplier_id, barcode, unit,
			unit_cost_minor, currency, reorder_point, reorder_quantity, track_lots, track_expiry, active
		) VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), $8, $9, upper($10), $11, $12, $13, $14, $15)`,
		product.ID, strings.TrimSpace(product.SKU), strings.TrimSpace(product.Name), strings.TrimSpace(product.Description), product.CategoryID,
		product.SupplierID, strings.TrimSpace(product.Barcode), strings.TrimSpace(product.Unit), product.UnitCostMinor, strings.TrimSpace(product.Currency),
		product.ReorderPoint, product.ReorderQuantity, product.TrackLots, product.TrackExpiry, product.Active)
	return mapError(err)
}

func (s *Store) UpdateProduct(ctx context.Context, product domain.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	command, err := s.pool.Exec(ctx, `
		UPDATE products SET
			sku=$2, name=$3, description=$4, category_id=NULLIF($5, ''), supplier_id=NULLIF($6, ''), barcode=NULLIF($7, ''), unit=$8,
			unit_cost_minor=$9, currency=upper($10), reorder_point=$11, reorder_quantity=$12, track_lots=$13, track_expiry=$14, active=$15, updated_at=now()
		WHERE id=$1`,
		product.ID, strings.TrimSpace(product.SKU), strings.TrimSpace(product.Name), strings.TrimSpace(product.Description), product.CategoryID,
		product.SupplierID, strings.TrimSpace(product.Barcode), strings.TrimSpace(product.Unit), product.UnitCostMinor, strings.TrimSpace(product.Currency),
		product.ReorderPoint, product.ReorderQuantity, product.TrackLots, product.TrackExpiry, product.Active)
	if err != nil {
		return mapError(err)
	}
	if command.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) GetProduct(ctx context.Context, id string) (domain.Product, error) {
	return scanProduct(s.pool.QueryRow(ctx, productSelect+` WHERE p.id=$1`, id))
}

func (s *Store) ListProducts(ctx context.Context, filter repository.ProductFilter) ([]domain.Product, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var where []string
	args := make([]any, 0, 6)
	add := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if q := strings.TrimSpace(filter.Query); q != "" {
		placeholder := add("%" + q + "%")
		where = append(where, `(p.sku ILIKE `+placeholder+` OR p.name ILIKE `+placeholder+` OR COALESCE(p.barcode, '') ILIKE `+placeholder+`)`)
	}
	if filter.CategoryID != "" {
		where = append(where, `p.category_id = `+add(filter.CategoryID))
	}
	if filter.SupplierID != "" {
		where = append(where, `p.supplier_id = `+add(filter.SupplierID))
	}
	if filter.ActiveOnly {
		where = append(where, `p.active = true`)
	}

	query := productSelect
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, ` AND `)
	}
	query += ` ORDER BY p.name, p.id LIMIT ` + add(limit) + ` OFFSET ` + add(filter.Offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Product, 0)
	for rows.Next() {
		item, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

const productSelect = `SELECT p.id, p.sku, p.name, p.description, COALESCE(p.category_id, ''), COALESCE(p.supplier_id, ''),
	COALESCE(p.barcode, ''), p.unit, p.unit_cost_minor, p.currency, p.reorder_point, p.reorder_quantity,
	p.track_lots, p.track_expiry, p.active, p.created_at, p.updated_at FROM products p`

type scanner interface {
	Scan(dest ...any) error
}

func scanProduct(row scanner) (domain.Product, error) {
	var item domain.Product
	err := row.Scan(&item.ID, &item.SKU, &item.Name, &item.Description, &item.CategoryID, &item.SupplierID, &item.Barcode, &item.Unit,
		&item.UnitCostMinor, &item.Currency, &item.ReorderPoint, &item.ReorderQuantity, &item.TrackLots, &item.TrackExpiry, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return domain.Product{}, mapError(err)
	}
	return item, nil
}

var _ repository.Catalog = (*Store)(nil)

package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type ProductFilter struct {
	Query      string
	CategoryID string
	SupplierID string
	ActiveOnly bool
	Limit      int
	Offset     int
}

type Catalog interface {
	CreateCategory(context.Context, domain.Category) error
	ListCategories(context.Context) ([]domain.Category, error)
	CreateSupplier(context.Context, domain.Supplier) error
	ListSuppliers(context.Context, bool) ([]domain.Supplier, error)
	CreateProduct(context.Context, domain.Product) error
	UpdateProduct(context.Context, domain.Product) error
	GetProduct(context.Context, string) (domain.Product, error)
	ListProducts(context.Context, ProductFilter) ([]domain.Product, error)
}

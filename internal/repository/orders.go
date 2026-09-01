package repository

import (
	"context"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

type Orders interface {
	CreateOrder(context.Context, domain.PurchaseOrder) error
	GetOrder(context.Context, string) (domain.PurchaseOrder, error)
	ListOrders(context.Context, domain.PurchaseOrderStatus, int, int) ([]domain.PurchaseOrder, error)
	UpdateOrder(context.Context, domain.PurchaseOrder) error
	ReceiveLine(context.Context, string, string, int64, string, string, string) error
	ReceiveLineWithNewLot(context.Context, string, string, int64, string, domain.Lot, string) error
	UpdateOrderStatus(context.Context, string, domain.PurchaseOrderStatus, string) error
}

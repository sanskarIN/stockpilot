package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (s *Store) ListLotInventory(ctx context.Context, filter repository.LotInventoryFilter) ([]domain.LotInventoryRow, error) {
	if filter.Limit <= 0 { filter.Limit = 100 }
	if filter.Limit > 500 { filter.Limit = 500 }
	if filter.Offset < 0 { filter.Offset = 0 }
	conditions := []string{"b.lot_id IS NOT NULL", "b.quantity > 0", "loc.active = true", "w.active = true"}
	args := make([]any, 0, 8)
	arg := func(value any) string { args = append(args, value); return fmt.Sprintf("$%d", len(args)) }
	if strings.TrimSpace(filter.ProductID) != "" { conditions = append(conditions, "b.product_id="+arg(strings.TrimSpace(filter.ProductID))) }
	if strings.TrimSpace(filter.WarehouseID) != "" { conditions = append(conditions, "w.id="+arg(strings.TrimSpace(filter.WarehouseID))) }
	if strings.TrimSpace(filter.LocationID) != "" { conditions = append(conditions, "loc.id="+arg(strings.TrimSpace(filter.LocationID))) }
	if strings.TrimSpace(filter.LotID) != "" { conditions = append(conditions, "l.id="+arg(strings.TrimSpace(filter.LotID))) }
	if filter.ExpiringBy != nil { conditions = append(conditions, "l.expires_at IS NOT NULL", "l.expires_at <= "+arg(*filter.ExpiringBy)) }
	limitArg := arg(filter.Limit)
	offsetArg := arg(filter.Offset)
	query := `SELECT b.product_id,p.sku,p.name,b.lot_id,l.lot_number,b.location_id,loc.name,w.id,w.name,b.quantity,l.expires_at,loc.active
FROM inventory_balances b
JOIN products p ON p.id=b.product_id
JOIN lots l ON l.id=b.lot_id
JOIN locations loc ON loc.id=b.location_id
JOIN warehouses w ON w.id=loc.warehouse_id
WHERE ` + strings.Join(conditions, " AND ") + `
ORDER BY CASE WHEN l.expires_at IS NULL THEN 1 ELSE 0 END, l.expires_at ASC, p.name ASC, l.lot_number ASC, loc.name ASC
LIMIT ` + limitArg + ` OFFSET ` + offsetArg
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]domain.LotInventoryRow, 0)
	for rows.Next() {
		var item domain.LotInventoryRow
		if err := rows.Scan(&item.ProductID,&item.SKU,&item.ProductName,&item.LotID,&item.LotNumber,&item.LocationID,&item.Location,&item.WarehouseID,&item.Warehouse,&item.OnHand,&item.ExpiresAt,&item.Active); err != nil { return nil, err }
		items = append(items, item)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return items, nil
}

var _ repository.LotInventory = (*Store)(nil)

package postgres

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/domain"
)

func (s *Store) ListReorderSuggestions(ctx context.Context, limit int) ([]domain.ReorderSuggestion, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			p.id,
			p.sku,
			p.name,
			COALESCE(p.supplier_id, ''),
			p.unit,
			COALESCE(SUM(b.quantity), 0)::bigint AS on_hand,
			p.reorder_point,
			p.reorder_quantity
		FROM products p
		LEFT JOIN inventory_balances b ON b.product_id = p.id
		WHERE p.active = true
		GROUP BY p.id, p.sku, p.name, p.supplier_id, p.unit, p.reorder_point, p.reorder_quantity
		HAVING COALESCE(SUM(b.quantity), 0) <= p.reorder_point
		ORDER BY (p.reorder_point - COALESCE(SUM(b.quantity), 0)) DESC, p.name ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ReorderSuggestion, 0)
	for rows.Next() {
		var item domain.ReorderSuggestion
		if err := rows.Scan(
			&item.ProductID,
			&item.SKU,
			&item.Name,
			&item.SupplierID,
			&item.Unit,
			&item.OnHand,
			&item.ReorderPoint,
			&item.ReorderQuantity,
		); err != nil {
			return nil, err
		}

		target, suggested, err := reorderTarget(item.ReorderPoint, item.ReorderQuantity, item.OnHand)
		if err != nil {
			return nil, fmt.Errorf("reorder suggestion for product %s: %w", item.ProductID, err)
		}
		item.TargetStock = target
		item.SuggestedQuantity = suggested
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GetInventoryValuation(ctx context.Context, limit int) (domain.InventoryValuationReport, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.unit,
			COALESCE(SUM(b.quantity), 0)::bigint AS on_hand,
			p.unit_cost_minor,
			p.currency,
			(COALESCE(SUM(b.quantity), 0)::numeric * p.unit_cost_minor::numeric)::text AS value_minor
		FROM products p
		LEFT JOIN inventory_balances b ON b.product_id = p.id
		WHERE p.active = true
		GROUP BY p.id, p.sku, p.name, p.unit, p.unit_cost_minor, p.currency
		ORDER BY p.name ASC
		LIMIT $1`, limit)
	if err != nil {
		return domain.InventoryValuationReport{}, err
	}
	defer rows.Close()

	report := domain.InventoryValuationReport{Items: make([]domain.InventoryValuationItem, 0)}
	for rows.Next() {
		var item domain.InventoryValuationItem
		var valueText string
		if err := rows.Scan(
			&item.ProductID,
			&item.SKU,
			&item.Name,
			&item.Unit,
			&item.OnHand,
			&item.UnitCostMinor,
			&item.Currency,
			&valueText,
		); err != nil {
			return domain.InventoryValuationReport{}, err
		}
		value, err := parseMinorValue(valueText)
		if err != nil {
			return domain.InventoryValuationReport{}, fmt.Errorf("valuation for product %s: %w", item.ProductID, err)
		}
		item.ValueMinor = value
		report.Items = append(report.Items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.InventoryValuationReport{}, err
	}

	totalRows, err := s.pool.Query(ctx, `
		WITH product_totals AS (
			SELECT
				p.currency,
				p.unit_cost_minor,
				COALESCE(SUM(b.quantity), 0)::numeric AS on_hand
			FROM products p
			LEFT JOIN inventory_balances b ON b.product_id = p.id
			WHERE p.active = true
			GROUP BY p.id, p.currency, p.unit_cost_minor
		)
		SELECT currency, SUM(on_hand * unit_cost_minor::numeric)::text AS value_minor
		FROM product_totals
		GROUP BY currency
		ORDER BY currency ASC`)
	if err != nil {
		return domain.InventoryValuationReport{}, err
	}
	defer totalRows.Close()

	report.Totals = make([]domain.InventoryValuationTotal, 0)
	for totalRows.Next() {
		var total domain.InventoryValuationTotal
		var valueText string
		if err := totalRows.Scan(&total.Currency, &valueText); err != nil {
			return domain.InventoryValuationReport{}, err
		}
		value, err := parseMinorValue(valueText)
		if err != nil {
			return domain.InventoryValuationReport{}, fmt.Errorf("valuation total for currency %s: %w", total.Currency, err)
		}
		total.ValueMinor = value
		report.Totals = append(report.Totals, total)
	}
	if err := totalRows.Err(); err != nil {
		return domain.InventoryValuationReport{}, err
	}

	return report, nil
}

func reorderTarget(reorderPoint, reorderQuantity, onHand int64) (int64, int64, error) {
	if reorderPoint < 0 || reorderQuantity < 0 || onHand < 0 {
		return 0, 0, fmt.Errorf("reorder inputs cannot be negative")
	}
	if reorderQuantity > math.MaxInt64-reorderPoint {
		return 0, 0, fmt.Errorf("reorder target exceeds supported range")
	}
	target := reorderPoint + reorderQuantity
	if onHand >= target {
		return target, 0, nil
	}
	return target, target - onHand, nil
}

func parseMinorValue(value string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("minor-unit value is outside the supported range: %w", err)
	}
	return parsed, nil
}

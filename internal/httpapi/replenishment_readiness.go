package httpapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/domain"
)

const (
	defaultReplenishmentWindowDays = 30
	maxReplenishmentWindowDays     = 365
	defaultReplenishmentRows       = 1000
	maxReplenishmentRows           = 5000
)

type movementHistoryReader interface {
	GetStockMovementHistory(context.Context, int, int) (domain.StockMovementHistoryReport, error)
}

type replenishmentCursor struct {
	RiskRank          int    `json:"riskRank"`
	SuggestedQuantity int64  `json:"suggestedQuantity"`
	SKU               string `json:"sku"`
	ProductID         string `json:"productId"`
}

func (a *API) replenishmentReadiness(w http.ResponseWriter, r *http.Request) {
	reader, ok := a.inventory.(movementHistoryReader)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "replenishment readiness is not available"})
		return
	}
	windowDays := normalizeReplenishmentWindow(queryInt(r, "days", defaultReplenishmentWindowDays))
	limit := normalizeReplenishmentLimit(queryInt(r, "limit", defaultReplenishmentRows))

	cursor, err := decodeReplenishmentCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid replenishment readiness cursor"})
		return
	}

	suggestions, err := a.inventory.ListReorderSuggestions(r.Context(), maxReplenishmentRows)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	movementReport, err := reader.GetStockMovementHistory(r.Context(), windowDays, maxReplenishmentRows)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	type velocity struct {
		outbound int64
		average  float64
	}
	velocities := make(map[string]velocity)
	for _, item := range movementReport.Items {
		v := velocities[item.ProductID]
		v.outbound += item.OutboundUnits
		v.average += item.AverageDailyOut
		velocities[item.ProductID] = v
	}

	items := make([]domain.ReplenishmentReadinessItem, 0, len(suggestions))
	for _, suggestion := range suggestions {
		v := velocities[suggestion.ProductID]
		item := domain.ReplenishmentReadinessItem{
			ProductID:         suggestion.ProductID,
			SKU:               suggestion.SKU,
			Name:              suggestion.Name,
			SupplierID:        suggestion.SupplierID,
			Unit:              suggestion.Unit,
			OnHand:            suggestion.OnHand,
			ReorderPoint:      suggestion.ReorderPoint,
			ReorderQuantity:   suggestion.ReorderQuantity,
			TargetStock:       suggestion.TargetStock,
			SuggestedQuantity: suggestion.SuggestedQuantity,
			OutboundUnits:     v.outbound,
			AverageDailyOut:   v.average,
			Risk:              classifyReplenishmentRisk(suggestion.OnHand, suggestion.ReorderPoint, v.average),
		}
		if v.average > 0 {
			days := float64(suggestion.OnHand) / v.average
			if days < 0 {
				days = 0
			}
			item.DaysOfCover = &days
		}
		items = append(items, item)
	}

	sort.SliceStable(items, func(i, j int) bool {
		ri, rj := replenishmentRiskRank(items[i].Risk), replenishmentRiskRank(items[j].Risk)
		if ri != rj {
			return ri < rj
		}
		if items[i].SuggestedQuantity != items[j].SuggestedQuantity {
			return items[i].SuggestedQuantity > items[j].SuggestedQuantity
		}
		if items[i].SKU != items[j].SKU {
			return items[i].SKU < items[j].SKU
		}
		return items[i].ProductID < items[j].ProductID
	})

	if cursor != nil {
		start := 0
		for start < len(items) && !afterReplenishmentCursor(items[start], *cursor) {
			start++
		}
		items = items[start:]
	}

	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	report := domain.ReplenishmentReadinessReport{AsOf: movementReport.AsOf, WindowDays: windowDays, Items: items}
	if hasMore && len(items) > 0 {
		last := items[len(items)-1]
		report.NextCursor = encodeReplenishmentCursor(replenishmentCursor{
			RiskRank:          replenishmentRiskRank(last.Risk),
			SuggestedQuantity: last.SuggestedQuantity,
			SKU:               last.SKU,
			ProductID:         last.ProductID,
		})
	}

	if r.URL.Query().Get("format") == "csv" {
		setCSVDownloadHeaders(w, "stockpilot-replenishment-readiness.csv")
		writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
		if err := writer.WriteHeader("productId", "sku", "name", "supplierId", "unit", "onHand", "reorderPoint", "reorderQuantity", "targetStock", "suggestedQuantity", "outboundUnits", "averageDailyOutbound", "daysOfCover", "risk", "asOf", "windowDays"); err != nil {
			return
		}
		for _, item := range report.Items {
			days := ""
			if item.DaysOfCover != nil {
				days = strconv.FormatFloat(*item.DaysOfCover, 'f', 2, 64)
			}
			if err := writer.WriteRow(item.ProductID, item.SKU, item.Name, item.SupplierID, item.Unit, strconv.FormatInt(item.OnHand, 10), strconv.FormatInt(item.ReorderPoint, 10), strconv.FormatInt(item.ReorderQuantity, 10), strconv.FormatInt(item.TargetStock, 10), strconv.FormatInt(item.SuggestedQuantity, 10), strconv.FormatInt(item.OutboundUnits, 10), strconv.FormatFloat(item.AverageDailyOut, 'f', 2, 64), days, string(item.Risk), report.AsOf.Format(time.RFC3339), strconv.Itoa(report.WindowDays)); err != nil {
				return
			}
		}
		_ = writer.Flush()
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func encodeReplenishmentCursor(cursor replenishmentCursor) string {
	payload, _ := json.Marshal(cursor)
	return base64.RawURLEncoding.EncodeToString(payload)
}

func decodeReplenishmentCursor(value string) (*replenishmentCursor, error) {
	if value == "" {
		return nil, nil
	}
	payload, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	var cursor replenishmentCursor
	if err := json.Unmarshal(payload, &cursor); err != nil {
		return nil, err
	}
	if cursor.RiskRank < 0 || cursor.RiskRank > 4 || cursor.SKU == "" || cursor.ProductID == "" {
		return nil, strconv.ErrSyntax
	}
	return &cursor, nil
}

func afterReplenishmentCursor(item domain.ReplenishmentReadinessItem, cursor replenishmentCursor) bool {
	riskRank := replenishmentRiskRank(item.Risk)
	if riskRank != cursor.RiskRank {
		return riskRank > cursor.RiskRank
	}
	if item.SuggestedQuantity != cursor.SuggestedQuantity {
		return item.SuggestedQuantity < cursor.SuggestedQuantity
	}
	if item.SKU != cursor.SKU {
		return item.SKU > cursor.SKU
	}
	return item.ProductID > cursor.ProductID
}

func classifyReplenishmentRisk(onHand, reorderPoint int64, averageDailyOut float64) domain.ReplenishmentRisk {
	if onHand <= 0 {
		return domain.ReplenishmentRiskOutOfStock
	}
	if averageDailyOut > 0 {
		daysOfCover := float64(onHand) / averageDailyOut
		if daysOfCover < 7 {
			return domain.ReplenishmentRiskCritical
		}
	}
	if onHand <= reorderPoint {
		return domain.ReplenishmentRiskReorder
	}
	if averageDailyOut > 0 && float64(onHand)/averageDailyOut < 14 {
		return domain.ReplenishmentRiskWatch
	}
	return domain.ReplenishmentRiskHealthy
}

func replenishmentRiskRank(risk domain.ReplenishmentRisk) int {
	switch risk {
	case domain.ReplenishmentRiskOutOfStock:
		return 0
	case domain.ReplenishmentRiskCritical:
		return 1
	case domain.ReplenishmentRiskReorder:
		return 2
	case domain.ReplenishmentRiskWatch:
		return 3
	default:
		return 4
	}
}

func normalizeReplenishmentWindow(days int) int {
	if days <= 0 {
		return defaultReplenishmentWindowDays
	}
	if days > maxReplenishmentWindowDays {
		return maxReplenishmentWindowDays
	}
	return days
}

func normalizeReplenishmentLimit(limit int) int {
	if limit <= 0 {
		return defaultReplenishmentRows
	}
	if limit > maxReplenishmentRows {
		return maxReplenishmentRows
	}
	return limit
}

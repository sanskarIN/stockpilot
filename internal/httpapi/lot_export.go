package httpapi

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const (
	defaultLotExportRows = 500
	maxLotExportRows     = 5000
)

func (a *API) exportLotInventoryCSV(w http.ResponseWriter, r *http.Request) {
	limit, offset := normalizeLotExportBounds(queryInt(r, "limit", defaultLotExportRows), queryInt(r, "offset", 0))
	filter := repository.LotInventoryFilter{
		ProductID:   strings.TrimSpace(r.URL.Query().Get("productId")),
		WarehouseID: strings.TrimSpace(r.URL.Query().Get("warehouseId")),
		LocationID:  strings.TrimSpace(r.URL.Query().Get("locationId")),
		LotID:       strings.TrimSpace(r.URL.Query().Get("lotId")),
		Limit:       limit,
		Offset:      offset,
	}
	if raw := strings.TrimSpace(r.URL.Query().Get("expiringBy")); raw != "" {
		value, err := parseExportDate(raw)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "expiringBy must be YYYY-MM-DD"})
			return
		}
		filter.ExpiringBy = &value
	}
	items, err := a.inventory.ListLotInventory(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	setCSVDownloadHeaders(w, "stockpilot-lot-inventory.csv")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader("productId", "sku", "productName", "lotId", "lotNumber", "locationId", "location", "warehouseId", "warehouse", "onHand", "expiresAt", "active"); err != nil {
		return
	}
	for _, item := range items {
		expiresAt := ""
		if item.ExpiresAt != nil {
			expiresAt = formatExportTime(*item.ExpiresAt)
		}
		if err := writer.WriteRow(
			item.ProductID,
			item.SKU,
			item.ProductName,
			item.LotID,
			item.LotNumber,
			item.LocationID,
			item.Location,
			item.WarehouseID,
			item.Warehouse,
			strconv.FormatInt(item.OnHand, 10),
			expiresAt,
			strconv.FormatBool(item.Active),
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func normalizeLotExportBounds(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultLotExportRows
	}
	if limit > maxLotExportRows {
		limit = maxLotExportRows
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func parseExportDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

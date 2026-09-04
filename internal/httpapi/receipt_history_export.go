package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

const (
	defaultReceiptHistoryExportRows = 500
	maxReceiptHistoryExportRows     = 5000
)

func (a *API) exportReceiptHistoryCSV(w http.ResponseWriter, r *http.Request) {
	limit, offset := normalizeReceiptHistoryExportBounds(queryInt(r, "limit", defaultReceiptHistoryExportRows), queryInt(r, "offset", 0))
	filter := repository.ReceiptHistoryFilter{
		ProductID:   strings.TrimSpace(r.URL.Query().Get("productId")),
		WarehouseID: strings.TrimSpace(r.URL.Query().Get("warehouseId")),
		LocationID:  strings.TrimSpace(r.URL.Query().Get("locationId")),
		LotID:       strings.TrimSpace(r.URL.Query().Get("lotId")),
		ActorID:     strings.TrimSpace(r.URL.Query().Get("actorId")),
		Reference:   strings.TrimSpace(r.URL.Query().Get("reference")),
		Limit:       limit,
		Offset:      offset,
	}
	if raw := strings.TrimSpace(r.URL.Query().Get("from")); raw != "" {
		value, err := parseExportDate(raw)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "from must be YYYY-MM-DD"})
			return
		}
		filter.From = &value
	}
	if raw := strings.TrimSpace(r.URL.Query().Get("to")); raw != "" {
		value, err := parseExportDate(raw)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "to must be YYYY-MM-DD"})
			return
		}
		filter.To = &value
	}
	if filter.From != nil && filter.To != nil && !filter.From.Before(*filter.To) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "from must be before to"})
		return
	}

	items, err := a.inventory.ListReceiptHistory(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	setCSVDownloadHeaders(w, "stockpilot-receipt-history.csv")
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader(
		"movementId", "productId", "sku", "productName", "locationId", "location", "warehouseId", "warehouse",
		"lotId", "lotNumber", "quantity", "reference", "note", "actorId", "occurredAt", "createdAt",
	); err != nil {
		return
	}
	for _, item := range items {
		if err := writer.WriteRow(
			item.MovementID,
			item.ProductID,
			item.SKU,
			item.ProductName,
			item.LocationID,
			item.Location,
			item.WarehouseID,
			item.Warehouse,
			item.LotID,
			item.LotNumber,
			strconv.FormatInt(item.Quantity, 10),
			item.Reference,
			item.Note,
			item.ActorID,
			formatExportTime(item.OccurredAt),
			formatExportTime(item.CreatedAt),
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func normalizeReceiptHistoryExportBounds(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultReceiptHistoryExportRows
	}
	if limit > maxReceiptHistoryExportRows {
		limit = maxReceiptHistoryExportRows
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

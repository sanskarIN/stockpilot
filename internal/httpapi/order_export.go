package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/sanskarIN/stockpilot/internal/csvexport"
	"github.com/sanskarIN/stockpilot/internal/domain"
)

const (
	defaultOrderExportRows = 500
	maxOrderExportRows     = 5000
)

func (a *API) exportOrdersCSV(w http.ResponseWriter, r *http.Request) {
	limit, offset := normalizeOrderExportBounds(queryInt(r, "limit", defaultOrderExportRows), queryInt(r, "offset", 0))
	status := domain.PurchaseOrderStatus(strings.TrimSpace(r.URL.Query().Get("status")))
	if status != "" && !status.Valid() {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported purchase order status"})
		return
	}

	items, err := a.orders.ListOrderExportRows(r.Context(), status, limit, offset)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="stockpilot-purchase-orders.csv"`)
	writer := csvexport.New(w, csvexport.Options{FormulaSafe: true})
	if err := writer.WriteHeader(
		"orderId", "orderNumber", "supplierId", "warehouseId", "status", "currency", "expectedAt",
		"notes", "createdBy", "createdAt", "updatedAt", "lineId", "productId", "quantity", "received",
		"remaining", "unitCostMinor", "lineTotalMinor",
	); err != nil {
		return
	}
	for _, item := range items {
		expectedAt := ""
		if item.ExpectedAt != nil {
			expectedAt = formatExportTime(*item.ExpectedAt)
		}
		remaining := item.Quantity - item.Received
		if err := writer.WriteRow(
			item.OrderID,
			item.OrderNumber,
			item.SupplierID,
			item.WarehouseID,
			string(item.Status),
			item.Currency,
			expectedAt,
			item.Notes,
			item.CreatedBy,
			formatExportTime(item.CreatedAt),
			formatExportTime(item.UpdatedAt),
			item.LineID,
			item.ProductID,
			strconv.FormatInt(item.Quantity, 10),
			strconv.FormatInt(item.Received, 10),
			strconv.FormatInt(remaining, 10),
			strconv.FormatInt(item.UnitCostMinor, 10),
			strconv.FormatInt(item.LineTotalMinor(), 10),
		); err != nil {
			return
		}
	}
	_ = writer.Flush()
}

func normalizeOrderExportBounds(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultOrderExportRows
	}
	if limit > maxOrderExportRows {
		limit = maxOrderExportRows
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

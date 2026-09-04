package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

func (f *fakeStore) ListReceiptHistory(_ context.Context, filter repository.ReceiptHistoryFilter) ([]domain.ReceiptHistoryRow, error) {
	items := make([]domain.ReceiptHistoryRow, 0, len(f.receiptHistory))
	for _, item := range f.receiptHistory {
		if filter.ProductID != "" && item.ProductID != filter.ProductID {
			continue
		}
		if filter.WarehouseID != "" && item.WarehouseID != filter.WarehouseID {
			continue
		}
		if filter.LocationID != "" && item.LocationID != filter.LocationID {
			continue
		}
		if filter.LotID != "" && item.LotID != filter.LotID {
			continue
		}
		if filter.ActorID != "" && item.ActorID != filter.ActorID {
			continue
		}
		if filter.Reference != "" && item.Reference != filter.Reference {
			continue
		}
		if filter.From != nil && item.OccurredAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && !item.OccurredAt.Before(*filter.To) {
			continue
		}
		items = append(items, item)
	}
	if filter.Offset >= len(items) {
		return []domain.ReceiptHistoryRow{}, nil
	}
	items = items[filter.Offset:]
	if filter.Limit > 0 && len(items) > filter.Limit {
		items = items[:filter.Limit]
	}
	return items, nil
}

func TestNormalizeReceiptHistoryExportBounds(t *testing.T) {
	if got, offset := normalizeReceiptHistoryExportBounds(0, -10); got != defaultReceiptHistoryExportRows || offset != 0 {
		t.Fatalf("defaults got limit=%d offset=%d", got, offset)
	}
	if got, _ := normalizeReceiptHistoryExportBounds(maxReceiptHistoryExportRows+1, 0); got != maxReceiptHistoryExportRows {
		t.Fatalf("clamped limit=%d", got)
	}
}

func TestReceiptHistoryExportCSV(t *testing.T) {
	occurred := time.Date(2026, 9, 4, 4, 30, 0, 0, time.FixedZone("IST", 19800))
	created := occurred.Add(2 * time.Minute)
	store := &fakeStore{receiptHistory: []domain.ReceiptHistoryRow{{
		MovementID: "mov_1", ProductID: "prd_1", SKU: "SKU-1", ProductName: "Widget",
		LocationID: "loc_1", Location: "A1", WarehouseID: "wh_1", Warehouse: "Main",
		LotID: "lot_1", LotNumber: "LOT-1", Quantity: 7, Reference: "PO-1",
		Note: "=formula", ActorID: "usr_1", OccurredAt: occurred, CreatedAt: created,
	}}}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(http.MethodGet, "/api/v1/inventory/receipts/export.csv?productId=prd_1&warehouseId=wh_1&locationId=loc_1&lotId=lot_1&actorId=usr_1&reference=PO-1&from=2026-09-04&to=2026-09-05&limit=10&offset=0", "")
	recorder := httptest.NewRecorder()	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content type=%q", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="stockpilot-receipt-history.csv"` {
		t.Fatalf("content disposition=%q", got)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "movementId,productId,sku,productName") {
		t.Fatalf("missing header body=%q", body)
	}
	if !strings.Contains(body, "mov_1,prd_1,SKU-1,Widget") || !strings.Contains(body, `'=formula`) {
		t.Fatalf("missing exported row body=%q", body)
	}
	if !strings.Contains(body, "2026-09-04T04:30:00Z") || !strings.Contains(body, "2026-09-04T04:32:00Z") {
		t.Fatalf("timestamps not normalized body=%q", body)
	}
}

func TestReceiptHistoryExportRejectsInvalidDateRange(t *testing.T) {
	store := &fakeStore{}
	handler := NewCore(store, store, store, func(context.Context) error { return nil })
	request := authenticatedRequest(http.MethodGet, "/api/v1/inventory/receipts/export.csv?from=2026-09-05&to=2026-09-04", "")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

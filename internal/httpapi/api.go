package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

type API struct {
	catalog   repository.Catalog
	inventory repository.Inventory
	orders    repository.Orders
	reports   repository.Reports
	audit     repository.Audit
	ping      func(context.Context) error
	origins   map[string]struct{}
	logger    *slog.Logger
}

type CoreOption func(*API)

func WithInsights(reports repository.Reports, audit repository.Audit) CoreOption {
	return func(api *API) { api.reports = reports; api.audit = audit }
}

func New(catalog repository.Catalog, inventory repository.Inventory, orders repository.Orders, ping func(context.Context) error, origins []string, logger *slog.Logger, options ...CoreOption) http.Handler {
	return WrapCommon(NewCore(catalog, inventory, orders, ping, options...), origins, logger)
}

func NewCore(catalog repository.Catalog, inventory repository.Inventory, orders repository.Orders, ping func(context.Context) error, options ...CoreOption) http.Handler {
	a := &API{catalog: catalog, inventory: inventory, orders: orders, ping: ping}
	for _, option := range options {
		if option != nil {
			option(a)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.health)
	mux.HandleFunc("GET /readyz", a.ready)
	mux.HandleFunc("GET /api/v1/meta", a.meta)
	mux.HandleFunc("GET /api/v1/categories", a.listCategories)
	mux.HandleFunc("POST /api/v1/categories", a.createCategory)
	mux.HandleFunc("GET /api/v1/suppliers", a.listSuppliers)
	mux.HandleFunc("POST /api/v1/suppliers", a.createSupplier)
	mux.HandleFunc("GET /api/v1/products", a.listProducts)
	mux.HandleFunc("GET /api/v1/products/export.csv", a.exportProductsCSV)
	mux.HandleFunc("POST /api/v1/products", a.createProduct)
	mux.HandleFunc("POST /api/v1/products/import/validate", a.validateProductImport)
	mux.HandleFunc("POST /api/v1/products/import", a.importProducts)
	mux.HandleFunc("GET /api/v1/products/by-barcode/{barcode}", a.getProductByBarcode)
	mux.HandleFunc("GET /api/v1/products/{id}", a.getProduct)
	mux.HandleFunc("PUT /api/v1/products/{id}", a.updateProduct)
	mux.HandleFunc("GET /api/v1/warehouses", a.listWarehouses)
	mux.HandleFunc("POST /api/v1/warehouses", a.createWarehouse)
	mux.HandleFunc("PUT /api/v1/warehouses/{id}", a.updateWarehouse)
	mux.HandleFunc("GET /api/v1/locations", a.listLocations)
	mux.HandleFunc("POST /api/v1/locations", a.createLocation)
	mux.HandleFunc("PUT /api/v1/locations/{id}", a.updateLocation)
	mux.HandleFunc("GET /api/v1/lots", a.listLots)
	mux.HandleFunc("POST /api/v1/lots", a.createLot)
	mux.HandleFunc("GET /api/v1/inventory/lots", a.listLotInventory)
	mux.HandleFunc("GET /api/v1/inventory/export.csv", a.exportInventoryCSV)
	mux.HandleFunc("GET /api/v1/inventory/low-stock", a.listLowStock)
	mux.HandleFunc("GET /api/v1/inventory/low-stock/export.csv", a.exportLowStockCSV)
	mux.HandleFunc("GET /api/v1/inventory/reorder-suggestions", a.listReorderSuggestions)
	mux.HandleFunc("GET /api/v1/inventory/reorder-suggestions/export.csv", a.exportReorderSuggestionsCSV)
	mux.HandleFunc("POST /api/v1/inventory/movements", a.applyMovement)
	mux.HandleFunc("POST /api/v1/inventory/transfers", a.transfer)
	mux.HandleFunc("GET /api/v1/reports/inventory-valuation", a.inventoryValuation)
	mux.HandleFunc("GET /api/v1/orders", a.listOrders)
	mux.HandleFunc("POST /api/v1/orders", a.createOrder)
	mux.HandleFunc("PUT /api/v1/orders/{id}", a.updateOrder)
	mux.HandleFunc("GET /api/v1/orders/{id}", a.getOrder)
	mux.HandleFunc("PATCH /api/v1/orders/{id}/status", a.updateOrderStatus)
	mux.HandleFunc("POST /api/v1/orders/{orderID}/lines/{lineID}/receive", a.receiveOrderLine)
	if a.reports != nil {
		mux.HandleFunc("GET /api/v1/reports/overview", a.reportOverview)
		mux.HandleFunc("GET /api/v1/reports/inventory", a.reportInventory)
		mux.HandleFunc("GET /api/v1/reports/purchasing", a.reportPurchasing)
	}
	if a.audit != nil {
		mux.HandleFunc("GET /api/v1/audit", a.listAuditEvents)
	}
	return mux
}

func WrapCommon(next http.Handler, origins []string, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	a := &API{logger: logger, origins: make(map[string]struct{}, len(origins))}
	for _, origin := range origins {
		if origin = strings.TrimSpace(origin); origin != "" {
			a.origins[origin] = struct{}{}
		}
	}
	return a.middleware(next)
}

func (a *API) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := idgen.New("req")
		if err != nil {
			requestID = "req_unavailable"
		}
		r = r.WithContext(context.WithValue(r.Context(), requestIDContextKey{}, requestID))
		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'; connect-src 'self'")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		origin := r.Header.Get("Origin")
		if origin != "" && strings.HasPrefix(r.URL.Path, "/api/") {
			if _, ok := a.origins[origin]; !ok {
				writeJSON(w, http.StatusForbidden, map[string]any{"error": "origin is not allowed", "requestId": requestID})
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID, X-StockPilot-CSRF")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		started := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		defer func() {
			if recovered := recover(); recovered != nil {
				a.logger.Error("http panic recovered", "request_id", requestID, "error", fmt.Sprint(recovered))
				if !sw.wroteHeader {
					writeJSON(sw, http.StatusInternalServerError, map[string]any{"error": "internal server error", "requestId": requestID})
				}
			}
			status := sw.status
			if status == 0 {
				status = http.StatusOK
			}
			a.logger.Info("http request", "request_id", requestID, "method", r.Method, "path", r.URL.Path, "status", status, "duration_ms", time.Since(started).Milliseconds())
		}()
		next.ServeHTTP(sw, r)
	})
}

func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if a.ping == nil || a.ping(ctx) != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a *API) meta(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"name": "StockPilot", "version": "0.1.0-dev", "credit": "Made by the Sanskar"})
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request, out *T) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return false
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request body must contain one JSON value"})
		return false
	}
	return true
}

func writeDomainError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"
	switch {
	case errors.Is(err, domain.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, message = http.StatusNotFound, "resource not found"
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrInsufficientStock):
		status, message = http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status, message = http.StatusForbidden, "forbidden"
	}
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

type statusWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *statusWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

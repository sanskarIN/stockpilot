package httpapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/domain"
	"github.com/sanskarIN/stockpilot/internal/idgen"
	"github.com/sanskarIN/stockpilot/internal/repository"
)

type replenishmentTraceabilityStore interface {
	repository.ReplenishmentTraceability
}

func (a *API) createReplenishmentReview(w http.ResponseWriter, r *http.Request) {
	store, ok := a.orders.(replenishmentTraceabilityStore)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "replenishment review traceability is not available"})
		return
	}
	var body struct {
		ProductID       string                             `json:"productId"`
		Outcome         domain.ReplenishmentReviewOutcome `json:"outcome"`
		PurchaseOrderID string                             `json:"purchaseOrderId,omitempty"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	body.ProductID = strings.TrimSpace(body.ProductID)
	body.PurchaseOrderID = strings.TrimSpace(body.PurchaseOrderID)
	if body.ProductID == "" || !body.Outcome.Valid() {
		writeDomainError(w, domain.ErrInvalid)
		return
	}

	suggestions, err := a.inventory.ListReorderSuggestions(r.Context(), maxReplenishmentRows)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	var suggestion *domain.ReorderSuggestion
	for i := range suggestions {
		if suggestions[i].ProductID == body.ProductID {
			suggestion = &suggestions[i]
			break
		}
	}
	if suggestion == nil {
		writeDomainError(w, domain.ErrNotFound)
		return
	}

	if body.Outcome == domain.ReplenishmentReviewAccepted || body.Outcome == domain.ReplenishmentReviewModified {
		if body.PurchaseOrderID == "" {
			writeDomainError(w, domain.ErrInvalid)
			return
		}
		order, err := a.orders.GetOrder(r.Context(), body.PurchaseOrderID)
		if err != nil {
			writeDomainError(w, err)
			return
		}
		if order.Status == domain.PurchaseOrderCancelled {
			writeDomainError(w, domain.ErrConflict)
			return
		}
		var lineQuantity int64
		found := false
		for _, line := range order.Lines {
			if line.ProductID == body.ProductID {
				lineQuantity = line.Quantity
				found = true
				break
			}
		}
		if !found || lineQuantity <= 0 {
			writeDomainError(w, domain.ErrInvalid)
			return
		}
		if body.Outcome == domain.ReplenishmentReviewAccepted && lineQuantity != suggestion.SuggestedQuantity {
			writeDomainError(w, domain.ErrInvalid)
			return
		}
		if body.Outcome == domain.ReplenishmentReviewModified && lineQuantity == suggestion.SuggestedQuantity {
			writeDomainError(w, domain.ErrInvalid)
			return
		}
	} else if body.PurchaseOrderID != "" {
		writeDomainError(w, domain.ErrInvalid)
		return
	}

	id, err := idgen.New("rr")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	now := time.Now().UTC()
	review := domain.ReplenishmentReview{
		ID:                id,
		ProductID:         suggestion.ProductID,
		SKU:               suggestion.SKU,
		ProductName:       suggestion.Name,
		Unit:              suggestion.Unit,
		SupplierID:        suggestion.SupplierID,
		OnHand:            suggestion.OnHand,
		ReorderPoint:      suggestion.ReorderPoint,
		ReorderQuantity:   suggestion.ReorderQuantity,
		TargetStock:       suggestion.TargetStock,
		SuggestedQuantity: suggestion.SuggestedQuantity,
		ReviewedBy:        authenticatedActorID(r),
		ReviewedAt:        now,
		PurchaseOrderID:   body.PurchaseOrderID,
		Outcome:           body.Outcome,
		CreatedAt:         now,
	}
	if err := store.CreateReplenishmentReview(r.Context(), review); err != nil {
		writeDomainError(w, err)
		return
	}
	a.recordAudit(r.Context(), review.ReviewedBy, "replenishment_review.created", "replenishment_review", review.ID, map[string]any{
		"productId": review.ProductID, "outcome": review.Outcome, "purchaseOrderId": review.PurchaseOrderID,
	})
	writeJSON(w, http.StatusCreated, review)
}

func (a *API) listReplenishmentReviews(w http.ResponseWriter, r *http.Request) {
	store, ok := a.orders.(replenishmentTraceabilityStore)
	if !ok {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "replenishment review traceability is not available"})
		return
	}
	items, err := store.ListReplenishmentReviews(r.Context(), strings.TrimSpace(r.URL.Query().Get("outcome")), queryInt(r, "limit", 50), queryInt(r, "offset", 0))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

package domain

import "fmt"

// ValidatePurchaseOrderTransition allows only operator-controlled lifecycle changes.
// Partial and received states are produced by the receiving transaction itself.
func ValidatePurchaseOrderTransition(from, to PurchaseOrderStatus) error {
	if from == to {
		return fmt.Errorf("%w: purchase order is already in status %s", ErrConflict, from)
	}
	switch from {
	case PurchaseOrderDraft:
		if to == PurchaseOrderOrdered || to == PurchaseOrderCancelled {
			return nil
		}
	case PurchaseOrderOrdered:
		if to == PurchaseOrderCancelled {
			return nil
		}
	}
	return fmt.Errorf("%w: cannot transition purchase order from %s to %s", ErrConflict, from, to)
}

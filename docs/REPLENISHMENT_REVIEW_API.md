# Replenishment Review Traceability API

StockPilot records a durable snapshot when a reorder recommendation is reviewed. The snapshot is server-generated from the current recommendation so callers cannot replace the historical product and quantity fields with arbitrary values.

## Create a review

`POST /api/v1/replenishment/reviews`

Request:

```json
{
  "productId": "prd_123",
  "outcome": "accepted",
  "purchaseOrderId": "po_123"
}
```

Valid outcomes:

- `accepted`: the linked purchase-order line quantity exactly matches the recommendation.
- `modified`: the linked purchase-order line quantity differs from the recommendation.
- `dismissed`: no purchase order is linked.
- `expired`: no purchase order is linked.

Creating a review does **not** create, submit, approve, receive, or otherwise mutate a purchase order.

For `accepted` and `modified`, the referenced purchase order must exist, must not be cancelled, and must contain the reviewed product. The server stores the recommendation values as an immutable snapshot.

## List reviews

`GET /api/v1/replenishment/reviews?outcome=accepted&limit=50&offset=0`

The `outcome` filter is optional. Results are ordered by review time descending and then review ID descending for deterministic pagination with the existing offset-based list contract.

## Audit behavior

Successful review creation records the existing append-only audit event `replenishment_review.created` with the product, outcome, and optional purchase-order identifier.

## Compatibility

The capability is additive. Existing purchase-order endpoints and lifecycle rules remain unchanged, and deployments without the v0.4.4 migration must return `501 Not Implemented` for the optional review endpoints rather than changing unrelated order behavior.

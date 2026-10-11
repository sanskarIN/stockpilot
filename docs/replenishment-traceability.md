# Replenishment Recommendation Traceability

## Purpose

StockPilot already calculates reorder recommendations and can seed a reviewed draft purchase order from a recommendation. v0.4.2 makes that relationship durable instead of relying on UI state, free-form notes, or timestamps.

## Traceability model

A traceability record represents the recommendation as it existed when a user reviewed it. The stored snapshot should include:

- product identifier and SKU;
- product name and unit;
- supplier identifier when available;
- on-hand quantity;
- reorder point;
- target stock;
- recommended/reorder quantity;
- suggested quantity presented to the reviewer;
- review actor and review time;
- resulting purchase-order identifier when one is created;
- explicit outcome.

The recommendation snapshot is immutable after creation. Later catalog or inventory changes must not rewrite historical recommendation values.

## Outcome states

The implementation uses explicit outcome semantics rather than inferring a decision from timestamps:

- `accepted` — recommendation was used without a quantity change;
- `modified` — recommendation was used but the reviewed purchase quantity differed;
- `dismissed` — recommendation was reviewed and intentionally not converted into a purchase order;
- `expired` — recommendation was no longer actionable when the review window ended.

## Safety rules

1. Creating a traceability record must not submit a purchase order automatically.
2. Existing purchase-order approval/status rules remain authoritative.
3. A traceability record may point to a purchase order, but a purchase order must remain valid without the optional traceability capability.
4. Existing report contracts remain unchanged.
5. The migration must be additive and have a tested rollback path.
6. Historical records are append-oriented; destructive update/delete APIs are intentionally avoided.

## Validation plan

Before v0.4.2 is tagged:

- run Go formatting, vet, tests, and build;
- run web quality checks;
- run PostgreSQL migration smoke tests from a clean database;
- test migration rollback where supported by the repository migration tooling;
- verify reorder-to-draft behavior remains intact;
- verify purchase-order lifecycle and audit behavior remain intact;
- verify CodeQL and dependency checks pass.

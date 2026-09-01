# Changelog

All notable StockPilot changes are recorded here. The project is pre-1.0, so current work remains under **Unreleased**.

## Unreleased

### Added

- Product, category, supplier, warehouse, location, lot, inventory movement, transfer, and purchase-order foundations.
- Session-based authentication, RBAC, CSRF protection, and administrator bootstrap tooling.
- PostgreSQL persistence with ordered migrations and transactional inventory operations.
- Responsive React + TypeScript web dashboard and installable PWA behavior.
- Native Android client with encrypted session storage, authenticated API access, dark mode, and release TLS enforcement.
- Manifest V3 browser companion with scoped optional host permissions and server health/launcher flow.
- CodeQL and Dependabot automation.
- Aggregate product-level reorder recommendations, including products with zero balance rows.
- Suggested replenishment quantities targeting reorder point plus reorder quantity.
- Currency-safe inventory valuation by product and grouped currency total.
- Exact barcode lookup API for scanner-driven clients.
- PostgreSQL and HTTP regression coverage for reporting capabilities.
- Security policy, contributor workflow, API compatibility policy, architecture guide, restore drill, release checklist, and issue/PR templates.
- Product catalog management, guided inventory operations, warehouse/location administration, and purchase-order creation/receiving workflows.
- Lot-aware receiving with reusable lots, new-lot creation, manufacturing/expiry metadata, and near-expiry warnings.
- Atomic new-lot receipt transactions that roll back lot creation when receipt processing fails.
- Explicit purchase-order lifecycle transitions for deliberate draft submission and cancellation.
- Append-only audit-event persistence for completed sensitive catalog, inventory, warehouse/location/lot, and purchasing mutations.
- Web audit history viewer with actor, action, entity, request ID, metadata, and pagination filters.
- Multi-line purchase-order creation, draft editing, per-line receiving, and draft-order total calculation.
- Reorder-suggestion actions that seed reviewed draft purchase orders without auto-submitting them.
- Warehouse/location edit, archive, and reactivation workflows with server-side integrity guards.

### Changed

- Dashboard stock health now uses product-level reorder recommendations instead of counting every low-stock location independently.
- Dashboard reporting now includes replenishment and inventory valuation insights.
- Repository contracts expose barcode lookup, replenishment, valuation, lot listing, order lifecycle, audit, draft-order update, reorder workflow, and warehouse/location lifecycle capabilities.
- Purchasing supports deliberate `draft → ordered/cancelled` transitions; partial/received states remain receipt-controlled.
- Receiving reuses existing lots for partial receipts and creates new lots atomically when required.
- Sensitive successful mutations now emit auditable actor/action/entity/request records.
- Draft purchase orders can now be edited transactionally until submission.
- Reorder alerts can open the standard purchasing editor with a reviewed product/quantity proposal; persistence still requires explicit user confirmation.
- Warehouse/location administration can edit active records, archive historical records instead of deleting them, and reactivate archived records when safe.

### Security

- HTTP mutation requests retain explicit CSRF confirmation requirements.
- Authorization remains server-side and role-aware.
- Production Android networking requires TLS.
- Browser companion permissions remain scoped to the configured StockPilot origin.
- Valuation arithmetic is calculated with PostgreSQL numeric values before checked conversion to `int64`.
- Lot ownership and product tracking policy remain enforced server-side during inventory mutations.
- Purchase-order status transitions are row-locked and validated by the domain state machine.
- Audit history is read-only from the web client and audit storage exposes append/list operations without an update/delete API.
- Draft purchase-order edits are limited to draft state and protected by a row lock.
- Reorder shortcuts do not bypass purchase-order creation or approval controls.
- Warehouses with active locations cannot be archived, and locations with non-zero inventory cannot be archived.
- Locations cannot be activated under an archived warehouse.

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.

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
- PostgreSQL and HTTP regression coverage for the reporting capabilities.
- Security policy, contributor workflow, API compatibility policy, architecture guide, restore drill, release checklist, and issue/PR templates.
- Product catalog management, guided inventory operations, warehouse/location administration, and purchase-order creation/receiving workflows.
- Lot-aware receiving with reusable existing lots, new-lot creation, manufacturing/expiry metadata, and configurable near-expiry warnings.
- Atomic new-lot receipt transactions that roll back lot creation when inventory receiving fails.

### Changed

- Dashboard stock health now uses product-level reorder recommendations instead of counting every low-stock location independently.
- Dashboard reporting now includes replenishment and inventory valuation insights.
- Repository contracts expose barcode lookup, replenishment, valuation, and lot-listing capabilities.
- Receiving now distinguishes existing lots from new lots, preventing duplicate batches during partial receipts.
- New-lot receipts are committed atomically with the purchase-order receipt and inventory movement.

### Security

- HTTP mutation requests retain explicit CSRF confirmation requirements.
- Authorization remains server-side and role-aware.
- Production Android networking requires TLS.
- Browser companion permissions remain scoped to the configured StockPilot origin.
- Valuation arithmetic is calculated with PostgreSQL numeric values before checked conversion to `int64`.
- Lot ownership and product tracking policy remain enforced server-side during inventory mutations.
- Atomic receiving prevents a failed stock transaction from leaving an uncommitted new lot.

## Release discipline

Dependency upgrades and generated artifacts are reviewed independently from product changes. Stable releases require the checks in `docs/RELEASE_CHECKLIST.md` and the restore procedure in `docs/RESTORE_DRILL.md`.

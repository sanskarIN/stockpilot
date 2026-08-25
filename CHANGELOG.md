# Changelog

All notable StockPilot changes are recorded here. The project is still pre-stable, so the current work remains under **Unreleased**.

## Unreleased

### Added

- Product, category, supplier, warehouse, location, lot, inventory movement, transfer, and purchase-order domain foundations.
- PostgreSQL persistence and ordered schema migrations.
- Secure HTTP API with health/readiness probes, validation, request IDs, CORS allowlisting, and security headers.
- User roles, permission model, password authentication, opaque sessions, CSRF protection, and operator-only administrator bootstrap.
- Responsive React + TypeScript inventory dashboard.
- Installable PWA behavior with production service-worker registration.
- Native Android application with secure session storage, authenticated API access, responsive dashboard UI, dark mode, release TLS enforcement, tests, and CI.
- Manifest V3 browser companion with configurable StockPilot origin, scoped optional permission, health check, launcher, source validation, tests, security documentation, and CI.
- CodeQL security analysis and Dependabot update streams.
- Aggregate reorder recommendations based on total product on-hand quantity, including products with no balance rows yet.
- Suggested replenishment quantities targeting `reorder point + reorder quantity`.
- Inventory valuation report with per-product values and full totals grouped by currency.
- Exact product lookup by unique barcode for scanner-driven clients.
- Root README, delivery roadmap, and refreshed work-continuity documentation.

### Changed

- Dashboard replenishment data now uses product-level reorder recommendations instead of treating every low location/lot balance as a separate product alert.
- Dashboard reporting now includes inventory valuation by currency.
- Repository contracts now expose replenishment, valuation, and barcode lookup capabilities.

### Security

- Session token hashes are peppered before persistence.
- Browser mutation requests require explicit CSRF confirmation.
- Role-based permissions protect catalog, inventory, purchasing, reporting, audit, and user-management surfaces.
- Android cleartext networking is limited to debug development and release builds require TLS.
- Browser-extension host access is optional and scoped to the configured StockPilot origin.
- Valuation arithmetic is performed with PostgreSQL numeric values before checked conversion to application `int64`, avoiding silent database integer overflow.

# StockPilot v0.2.0 — Reports & Analytics

## Release summary

StockPilot v0.2.0 begins the broader pre-1.0 reporting milestone. It turns the existing reporting API foundation into a dedicated web Reports & Analytics workspace and adds a secure inventory-valuation CSV export.

## Highlights

### Reports & Analytics workspace

- Dedicated authenticated Reports screen.
- Inventory operational metrics for active products, units on hand, low-stock balances, and outstanding purchase-order units.
- Purchasing pipeline breakdown for draft, ordered, partially received, received, and cancelled orders.
- Currency-grouped inventory valuation summary.
- Product-level valuation breakdown with on-hand quantity, unit cost, currency, and value.
- Refresh and session-expiry handling consistent with the rest of the web application.

### Inventory valuation export

- New endpoint: `GET /api/v1/reports/inventory-valuation/export.csv`.
- Bounded exports with a default of 1,000 rows and a maximum of 5,000 rows.
- Formula-safe CSV serialization.
- Consistent `no-store` and `no-cache` response headers.
- Deterministic filename: `stockpilot-inventory-valuation.csv`.
- Existing reporting RBAC remains authoritative.
- Existing v0.1.10 export audit instrumentation records authenticated CSV requests.

## Compatibility

- No database migration is required for this milestone.
- Existing reporting endpoints remain compatible.
- Existing CSV exports remain unchanged.
- The project remains pre-1.0 and should not be treated as a stable API contract.

## Security notes

The new export is read-only and does not expose credentials or session values. Formula-safe serialization is enabled to reduce spreadsheet formula-injection risk. Export downloads retain the repository-wide privacy-oriented cache policy.

## Verification

Focused regression coverage has been added for valuation export bounds, response headers, formula-safe cells, and CSV content. Full release verification must be performed through GitHub Actions or a local development environment because the connected workspace does not provide a trustworthy project shell.

## Known limitations

- Inventory aging, configurable expiry-risk reporting, movement velocity, supplier lead-time analytics, and warehouse/location valuation breakdowns remain planned follow-up work.
- Large dataset cursor/streaming exports remain future work.
- Reproducible release artifacts and full end-to-end test coverage remain release-hardening tasks.

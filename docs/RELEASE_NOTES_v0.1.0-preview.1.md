# StockPilot v0.1.0-preview.1

StockPilot's first public preview packages the current inventory and purchasing foundation for early evaluation and feedback.

## Highlights

- Product catalog, categories, suppliers, warehouses, locations, and lot tracking.
- Inventory movements, transfers, low-stock reporting, reorder suggestions, and inventory valuation.
- Purchase-order creation, lifecycle controls, and receiving workflows.
- Lot-aware receiving with atomic new-lot creation.
- Session authentication, role-based authorization, CSRF protection, and administrator bootstrap tooling.
- Append-only business and authentication audit events with request-ID correlation.
- CSV product dry-run validation and transactional CSV product import.
- Server-side import revalidation, generated product IDs, database constraints, and complete-batch rollback on persistence failure.
- React/TypeScript web dashboard and PWA assets.
- Native Kotlin Android client with encrypted session storage and release TLS enforcement.
- Manifest V3 browser companion with scoped host permissions and navigation-only inventory handoff.
- PostgreSQL migrations, Docker deployment, backups, health/readiness endpoints, CI, and CodeQL analysis.

## CSV product import

The preview includes two deliberately separate CSV workflows:

1. **Validate** parses the CSV and reports field, duplicate, and reference errors without writing products.
2. **Import** reparses and revalidates the submitted file immediately before opening one PostgreSQL transaction.

The import transaction relies on application validation plus database uniqueness and foreign-key constraints. If a batch fails during persistence, the complete batch is rolled back. The successful import audit event records request context and batch count but does not store the CSV contents.

## Verification

The release candidate was merged only after the final PR checks passed:

- Go module tidy verification.
- `gofmt` formatting gate.
- `go vet ./...`.
- Race-enabled Go tests with coverage.
- Server build.
- PostgreSQL migration/readiness smoke test.
- Web typecheck and production build.
- CodeQL analysis for Go and JavaScript/TypeScript.

## Preview status

This is a **pre-1.0 preview**, not a production-stability guarantee. The remaining release work includes a real backup/restore drill, manual authentication and UI smoke tests, Android device smoke testing, browser-companion installation testing, reproducible artifact publication, and post-release verification.

## Upgrade and migration notes

- Apply the repository's ordered PostgreSQL migrations before using the server against an existing database.
- Review the environment variables in `.env.example` and never commit production secrets.
- Use HTTPS and a trusted TLS termination path for production deployments.
- Review role permissions before exposing mutation endpoints to operators.

## Feedback

Please report reproducible bugs with the StockPilot issue tracker and include the relevant server, web, Android, or browser-companion component, expected behavior, actual behavior, and safe diagnostic information.
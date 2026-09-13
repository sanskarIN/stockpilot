# StockPilot

StockPilot is an open-source inventory and purchasing platform with a Go/PostgreSQL backend, a React/TypeScript web application, a native Android client, and a Manifest V3 browser companion.

The project is currently in pre-1.0 development. The current development release is `0.6.0`; the Android and browser companion clients remain independently versioned.

## What StockPilot includes

- Product catalog, categories, suppliers, warehouses, and locations.
- Inventory movements, transfers, lot tracking, balances, and low-stock reporting.
- Purchase orders and receiving workflows.
- Session-based authentication and role-based access control.
- Administrator bootstrap tooling and user-management APIs.
- PostgreSQL migrations with automatic migration support.
- Responsive React/TypeScript web dashboard and PWA assets.
- Native Kotlin Android client with encrypted session storage.
- Manifest V3 browser companion with optional per-server host permissions.
- Health/readiness endpoints, structured HTTP hardening, Docker deployment, backups, and CI/CodeQL checks.
- Reports & Analytics for inventory valuation, warehouse/location valuation, inventory aging, expiry risk, movement velocity, and supplier performance.

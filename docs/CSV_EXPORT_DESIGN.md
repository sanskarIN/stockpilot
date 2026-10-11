# CSV Export Design

## Purpose

StockPilot will use CSV exports for operational reporting without turning the export layer into a second data-access API. The export layer is responsible for deterministic serialization, while repositories remain responsible for filtering, authorization, ordering, and pagination.

## v0.1.2 foundation

The reusable `internal/csvexport` package provides:

- RFC 4180-compatible CSV quoting through Go's standard `encoding/csv` package.
- Explicit header writing.
- Buffered writes with error reporting on `Flush`.
- Optional spreadsheet-formula protection for cells beginning with `=`, `+`, `-`, or `@` after leading whitespace.
- No external dependency.

The writer intentionally does not fetch data, enforce permissions, or decide business-level columns.

## Planned export endpoints

Later application-level export work should add authenticated, read-only endpoints for bounded operational datasets, initially:

1. Product catalog.
2. Current inventory balances.
3. Low-stock and reorder suggestions.
4. Lot inventory / expiry-risk rows.
5. Purchase orders and receiving history.
6. Audit/report datasets where the requesting role is authorized.

Each endpoint should define a stable column contract and an explicit maximum result size. Large exports should move to a streaming or asynchronous job model rather than buffering an entire dataset in memory.

## Security requirements

- Apply the same server-side authorization checks as the corresponding JSON/report endpoint.
- Never export session cookies, passwords, CSRF secrets, API credentials, or other secrets.
- Treat CSV as untrusted data when opened by spreadsheet software; enable formula-safe output for user-downloadable exports where appropriate.
- Set `Content-Type: text/csv; charset=utf-8` and a deterministic `Content-Disposition` filename.
- Record an audit event for sensitive or administrator-only exports when the existing audit model supports it.
- Bound filters, page sizes, and total rows to protect database and application resources.

## Determinism

Export column order must be explicit and documented. Database queries should provide deterministic ordering so repeated exports over an unchanged dataset produce stable row order.

## Compatibility

CSV is deliberately chosen as the first portable export format. JSON APIs remain the machine-to-machine interface. XLSX/PDF generation is out of scope for this milestone and should be considered separately after the CSV contracts have stabilized.

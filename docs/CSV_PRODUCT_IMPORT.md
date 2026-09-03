# CSV Product Import

StockPilot supports a two-step CSV product import workflow: **dry-run validation** followed by an explicit **transactional write**.

A ready-to-copy template is available at `examples/products-import.csv`.

## Endpoints

### Dry run

`POST /api/v1/products/import/validate`

### Write

`POST /api/v1/products/import`

Both endpoints require the normal authenticated catalog-write permission and CSRF protection. Upload a multipart form with a field named `file`.

The upload is bounded to 5 MiB and 1,000 product rows. The parser also applies a 4 MiB CSV reader limit so malformed or unexpectedly large input cannot grow unbounded in memory.

## Required columns

```text
sku,name,unit,unit_cost_minor,currency,reorder_point,reorder_quantity,track_lots,track_expiry,active
```

Optional columns are:

```text
id,description,category_id,supplier_id,barcode
```

`unit_cost_minor` is stored as an integer minor-unit amount. `currency` must be a three-letter code. Boolean fields must use `true` or `false`.

## Validation

The dry run checks:

- required headers and duplicate headers;
- product domain constraints;
- duplicate SKUs within the upload;
- duplicate barcodes within the upload;
- existing SKUs in the catalog;
- category and supplier references;
- row-level parsing errors;
- the maximum row count.

The response contains the number of valid rows, the number of errors, row-level error messages, and a small valid-row preview.

## Explicit write step

A successful dry run does **not** automatically write anything. The web UI presents a separate **Import products** action after a clean dry run.

The write endpoint reparses the uploaded CSV and repeats server-side validation against current database state. This closes the time-of-check/time-of-use gap between validation and persistence.

Products without an `id` receive a server-generated product ID. Supplied IDs are retained and remain subject to the database primary-key constraint.

The PostgreSQL persistence layer inserts the complete batch in one transaction. If any product fails domain validation, a foreign-key constraint, the SKU/barcode uniqueness constraint, or another database constraint, the entire batch is rolled back and no partial import remains.

The write response returns the number of imported products and the created product representations. A successful import records a `products.imported` audit event containing the request ID and batch count, without storing the CSV contents in audit metadata.

## Integrity boundary

The dry-run result is never treated as authorization to write. The write request is independently authenticated and authorized, revalidates its own payload, and relies on database constraints as the final integrity boundary.

The web catalog exposes the import workflow only through the existing administrator/manager catalog workspace. The server remains authoritative for permissions and all final product constraints.

## Failure behavior

- malformed or oversized multipart/CSV input returns a client error;
- row-level validation errors prevent the write and are returned with row numbers;
- an unsupported repository implementation returns `501 Not Implemented` rather than silently falling back to non-atomic writes;
- concurrent SKU or barcode conflicts are resolved by the database uniqueness constraints and roll back the complete batch.

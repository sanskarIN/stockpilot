# CSV Product Import

StockPilot's first CSV import milestone provides a **dry-run validation** workflow. It does not write products to PostgreSQL.

A ready-to-copy template is available at `examples/products-import.csv`.

## Endpoint

`POST /api/v1/products/import/validate`

The endpoint requires the normal authenticated catalog-write permission and CSRF protection. Upload a multipart form with a field named `file`.

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

## Integrity boundary

The dry-run endpoint intentionally has no persistence path. A successful dry run means the file is suitable for the next write/import milestone; it does **not** create products.

The web catalog exposes the validation panel only to administrator and manager roles. The server remains authoritative for permissions and all final product constraints.

## Next persistence milestone

The production import endpoint must revalidate the complete payload server-side and persist the complete batch transactionally. It must not trust a previous dry-run result as authorization or as proof that the database has not changed since validation. Unique database constraints and the transaction must remain the final integrity boundary.

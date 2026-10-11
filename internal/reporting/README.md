# Reporting primitives

`internal/reporting` contains small, dependency-free helpers shared by reporting features.

The package deliberately keeps reporting semantics separate from PostgreSQL and HTTP concerns. It provides:

- bounded inclusive date periods (1–365 days);
- deterministic previous-period calculation;
- period-over-period deltas with an explicit undefined state for zero baselines;
- bounded limit/offset validation;
- safe parsing of non-negative integer query values.

These primitives are read-only and do not perform database access or inventory/purchasing mutations.

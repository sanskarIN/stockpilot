BEGIN;

DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS purchase_order_lines;
DROP TABLE IF EXISTS purchase_orders;
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS inventory_balances;
DROP TABLE IF EXISTS lots;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS warehouses;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS suppliers;
DROP TABLE IF EXISTS categories;
DELETE FROM schema_migrations WHERE version = 1;
DROP TABLE IF EXISTS schema_migrations;

COMMIT;

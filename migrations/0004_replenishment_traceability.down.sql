BEGIN;

DROP INDEX IF EXISTS replenishment_reviews_purchase_order_idx;
DROP INDEX IF EXISTS replenishment_reviews_outcome_time_idx;
DROP INDEX IF EXISTS replenishment_reviews_product_time_idx;
DROP TABLE IF EXISTS replenishment_reviews;
DELETE FROM schema_migrations WHERE version = 4;

COMMIT;

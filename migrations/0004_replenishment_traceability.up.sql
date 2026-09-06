BEGIN;

CREATE TABLE replenishment_reviews (
  id text PRIMARY KEY,
  product_id text NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  sku varchar(64) NOT NULL,
  product_name varchar(200) NOT NULL,
  unit varchar(32) NOT NULL,
  supplier_id text REFERENCES suppliers(id) ON DELETE RESTRICT,
  on_hand bigint NOT NULL,
  reorder_point bigint NOT NULL,
  reorder_quantity bigint NOT NULL,
  target_stock bigint NOT NULL,
  suggested_quantity bigint NOT NULL,
  reviewed_by text NOT NULL,
  reviewed_at timestamptz NOT NULL,
  purchase_order_id text REFERENCES purchase_orders(id) ON DELETE RESTRICT,
  outcome varchar(16) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (on_hand >= 0),
  CHECK (reorder_point >= 0),
  CHECK (reorder_quantity >= 0),
  CHECK (target_stock >= 0),
  CHECK (suggested_quantity > 0),
  CHECK (length(btrim(sku)) BETWEEN 2 AND 64),
  CHECK (length(btrim(product_name)) BETWEEN 2 AND 200),
  CHECK (length(btrim(unit)) BETWEEN 1 AND 32),
  CHECK (length(btrim(reviewed_by)) > 0),
  CHECK (outcome IN ('accepted', 'modified', 'dismissed', 'expired')),
  CHECK ((outcome IN ('accepted', 'modified') AND purchase_order_id IS NOT NULL) OR
         (outcome IN ('dismissed', 'expired') AND purchase_order_id IS NULL))
);
CREATE INDEX replenishment_reviews_product_time_idx ON replenishment_reviews (product_id, reviewed_at DESC, id DESC);
CREATE INDEX replenishment_reviews_outcome_time_idx ON replenishment_reviews (outcome, reviewed_at DESC, id DESC);
CREATE INDEX replenishment_reviews_purchase_order_idx ON replenishment_reviews (purchase_order_id) WHERE purchase_order_id IS NOT NULL;

INSERT INTO schema_migrations(version) VALUES (4) ON CONFLICT (version) DO NOTHING;

COMMIT;

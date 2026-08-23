BEGIN;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version bigint PRIMARY KEY,
  applied_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE categories (
  id text PRIMARY KEY,
  name varchar(120) NOT NULL,
  description text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(name)) BETWEEN 2 AND 120)
);
CREATE UNIQUE INDEX categories_name_ci_uq ON categories (lower(name));

CREATE TABLE suppliers (
  id text PRIMARY KEY,
  code varchar(48) NOT NULL UNIQUE,
  name varchar(160) NOT NULL,
  email varchar(254),
  phone varchar(40),
  notes text NOT NULL DEFAULT '',
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(code)) BETWEEN 2 AND 48),
  CHECK (length(btrim(name)) BETWEEN 2 AND 160)
);
CREATE INDEX suppliers_active_name_idx ON suppliers (active, name);

CREATE TABLE products (
  id text PRIMARY KEY,
  sku varchar(64) NOT NULL UNIQUE,
  name varchar(200) NOT NULL,
  description text NOT NULL DEFAULT '',
  category_id text REFERENCES categories(id) ON DELETE SET NULL,
  supplier_id text REFERENCES suppliers(id) ON DELETE SET NULL,
  barcode varchar(128),
  unit varchar(32) NOT NULL,
  unit_cost_minor bigint NOT NULL DEFAULT 0,
  currency char(3) NOT NULL DEFAULT 'INR',
  reorder_point bigint NOT NULL DEFAULT 0,
  reorder_quantity bigint NOT NULL DEFAULT 0,
  track_lots boolean NOT NULL DEFAULT false,
  track_expiry boolean NOT NULL DEFAULT false,
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(sku)) BETWEEN 2 AND 64),
  CHECK (length(btrim(name)) BETWEEN 2 AND 200),
  CHECK (length(btrim(unit)) BETWEEN 1 AND 32),
  CHECK (unit_cost_minor >= 0),
  CHECK (reorder_point >= 0),
  CHECK (reorder_quantity >= 0),
  CHECK (currency ~ '^[A-Z]{3}$'),
  CHECK (NOT track_expiry OR track_lots)
);
CREATE UNIQUE INDEX products_barcode_uq ON products (barcode) WHERE barcode IS NOT NULL AND barcode <> '';
CREATE INDEX products_active_name_idx ON products (active, name);
CREATE INDEX products_category_idx ON products (category_id);
CREATE INDEX products_supplier_idx ON products (supplier_id);

CREATE TABLE warehouses (
  id text PRIMARY KEY,
  code varchar(48) NOT NULL UNIQUE,
  name varchar(160) NOT NULL,
  address text NOT NULL DEFAULT '',
  timezone varchar(64) NOT NULL DEFAULT 'UTC',
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(code)) BETWEEN 2 AND 48),
  CHECK (length(btrim(name)) BETWEEN 2 AND 160),
  CHECK (length(btrim(timezone)) >= 1)
);

CREATE TABLE locations (
  id text PRIMARY KEY,
  warehouse_id text NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
  code varchar(64) NOT NULL,
  name varchar(160) NOT NULL,
  description text NOT NULL DEFAULT '',
  active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (warehouse_id, code),
  CHECK (length(btrim(code)) BETWEEN 1 AND 64),
  CHECK (length(btrim(name)) BETWEEN 1 AND 160)
);
CREATE INDEX locations_warehouse_active_idx ON locations (warehouse_id, active, name);

CREATE TABLE lots (
  id text PRIMARY KEY,
  product_id text NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  lot_number varchar(96) NOT NULL,
  manufactured_at timestamptz,
  expires_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (product_id, lot_number),
  CHECK (length(btrim(lot_number)) BETWEEN 1 AND 96),
  CHECK (manufactured_at IS NULL OR expires_at IS NULL OR expires_at > manufactured_at)
);
CREATE INDEX lots_expiry_idx ON lots (expires_at) WHERE expires_at IS NOT NULL;

CREATE TABLE inventory_balances (
  id text PRIMARY KEY,
  product_id text NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  location_id text NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
  lot_id text REFERENCES lots(id) ON DELETE RESTRICT,
  quantity bigint NOT NULL DEFAULT 0,
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (quantity >= 0)
);
CREATE UNIQUE INDEX inventory_balance_no_lot_uq
  ON inventory_balances (product_id, location_id)
  WHERE lot_id IS NULL;
CREATE UNIQUE INDEX inventory_balance_lot_uq
  ON inventory_balances (product_id, location_id, lot_id)
  WHERE lot_id IS NOT NULL;
CREATE INDEX inventory_balance_product_idx ON inventory_balances (product_id);
CREATE INDEX inventory_balance_location_idx ON inventory_balances (location_id);

CREATE TABLE stock_movements (
  id text PRIMARY KEY,
  product_id text NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  location_id text NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
  lot_id text REFERENCES lots(id) ON DELETE RESTRICT,
  movement_type varchar(24) NOT NULL,
  quantity_delta bigint NOT NULL,
  reference varchar(128) NOT NULL DEFAULT '',
  note text NOT NULL DEFAULT '',
  actor_id text NOT NULL DEFAULT '',
  occurred_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (movement_type IN ('stock_in', 'stock_out', 'adjustment', 'transfer_in', 'transfer_out', 'receive')),
  CHECK (
    (movement_type IN ('stock_in', 'transfer_in', 'receive') AND quantity_delta > 0) OR
    (movement_type IN ('stock_out', 'transfer_out') AND quantity_delta < 0) OR
    (movement_type = 'adjustment' AND quantity_delta <> 0)
  )
);
CREATE INDEX stock_movements_product_time_idx ON stock_movements (product_id, occurred_at DESC);
CREATE INDEX stock_movements_location_time_idx ON stock_movements (location_id, occurred_at DESC);
CREATE INDEX stock_movements_reference_idx ON stock_movements (reference) WHERE reference <> '';

CREATE TABLE purchase_orders (
  id text PRIMARY KEY,
  order_number varchar(64) NOT NULL UNIQUE,
  supplier_id text NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
  warehouse_id text NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
  status varchar(24) NOT NULL DEFAULT 'draft',
  currency char(3) NOT NULL DEFAULT 'INR',
  expected_at timestamptz,
  notes text NOT NULL DEFAULT '',
  created_by text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (status IN ('draft', 'ordered', 'partially_received', 'received', 'cancelled')),
  CHECK (currency ~ '^[A-Z]{3}$')
);
CREATE INDEX purchase_orders_status_time_idx ON purchase_orders (status, created_at DESC);
CREATE INDEX purchase_orders_supplier_idx ON purchase_orders (supplier_id, created_at DESC);

CREATE TABLE purchase_order_lines (
  id text PRIMARY KEY,
  purchase_order_id text NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
  product_id text NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  quantity bigint NOT NULL,
  received bigint NOT NULL DEFAULT 0,
  unit_cost_minor bigint NOT NULL DEFAULT 0,
  UNIQUE (purchase_order_id, product_id),
  CHECK (quantity > 0),
  CHECK (received >= 0 AND received <= quantity),
  CHECK (unit_cost_minor >= 0)
);
CREATE INDEX purchase_order_lines_product_idx ON purchase_order_lines (product_id);

CREATE TABLE audit_log (
  id bigserial PRIMARY KEY,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  actor_id text NOT NULL DEFAULT '',
  action varchar(80) NOT NULL,
  entity_type varchar(80) NOT NULL,
  entity_id text NOT NULL,
  request_id text NOT NULL DEFAULT '',
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX audit_log_entity_idx ON audit_log (entity_type, entity_id, occurred_at DESC);
CREATE INDEX audit_log_actor_idx ON audit_log (actor_id, occurred_at DESC) WHERE actor_id <> '';

INSERT INTO schema_migrations(version) VALUES (1) ON CONFLICT (version) DO NOTHING;

COMMIT;

BEGIN;

CREATE TABLE users (
  id text PRIMARY KEY,
  email varchar(254) NOT NULL,
  display_name varchar(120) NOT NULL,
  password_hash text NOT NULL,
  role varchar(24) NOT NULL,
  active boolean NOT NULL DEFAULT true,
  last_login_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(btrim(email)) BETWEEN 3 AND 254),
  CHECK (length(btrim(display_name)) BETWEEN 2 AND 120),
  CHECK (length(password_hash) >= 20),
  CHECK (role IN ('admin', 'manager', 'operator', 'viewer'))
);
CREATE UNIQUE INDEX users_email_ci_uq ON users (lower(email));
CREATE INDEX users_active_role_idx ON users (active, role, display_name);

CREATE TABLE sessions (
  id text PRIMARY KEY,
  user_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash char(64) NOT NULL UNIQUE,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  last_seen_at timestamptz NOT NULL DEFAULT now(),
  CHECK (length(token_hash) = 64),
  CHECK (expires_at > created_at)
);
CREATE INDEX sessions_user_expiry_idx ON sessions (user_id, expires_at DESC);
CREATE INDEX sessions_expiry_idx ON sessions (expires_at);

INSERT INTO schema_migrations(version) VALUES (2) ON CONFLICT (version) DO NOTHING;

COMMIT;

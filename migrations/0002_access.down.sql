BEGIN;

DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DELETE FROM schema_migrations WHERE version = 2;

COMMIT;

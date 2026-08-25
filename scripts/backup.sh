#!/usr/bin/env sh
set -eu

umask 077

BACKUP_DIR=${BACKUP_DIR:-backups}
DB_SERVICE=${DB_SERVICE:-db}
DB_USER=${POSTGRES_USER:-stockpilot}
DB_NAME=${POSTGRES_DB:-stockpilot}
TIMESTAMP=$(date -u +%Y%m%dT%H%M%SZ)
FINAL_PATH="${BACKUP_DIR}/stockpilot-${TIMESTAMP}.dump"
TEMP_PATH="${FINAL_PATH}.partial"

mkdir -p "$BACKUP_DIR"

cleanup() {
  rm -f "$TEMP_PATH"
}
trap cleanup EXIT HUP INT TERM

echo "Creating StockPilot PostgreSQL backup: ${FINAL_PATH}"
docker compose exec -T "$DB_SERVICE" \
  pg_dump \
  --username="$DB_USER" \
  --dbname="$DB_NAME" \
  --format=custom \
  --compress=6 \
  --no-owner \
  --no-acl > "$TEMP_PATH"

if [ ! -s "$TEMP_PATH" ]; then
  echo "Backup failed: pg_dump produced an empty file." >&2
  exit 1
fi

mv "$TEMP_PATH" "$FINAL_PATH"
trap - EXIT HUP INT TERM

echo "Backup complete: ${FINAL_PATH}"

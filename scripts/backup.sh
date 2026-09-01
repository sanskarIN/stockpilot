#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-${ROOT_DIR}/backups}"
BACKUP_RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-14}"
POSTGRES_USER="${POSTGRES_USER:-stockpilot}"
POSTGRES_DB="${POSTGRES_DB:-stockpilot}"
TIMESTAMP="$(date -u +'%Y%m%dT%H%M%SZ')"
FINAL_PATH="${BACKUP_DIR}/stockpilot-${TIMESTAMP}.dump"
TEMP_PATH="${FINAL_PATH}.tmp"

if ! [[ "${BACKUP_RETENTION_DAYS}" =~ ^[0-9]+$ ]]; then
  echo "BACKUP_RETENTION_DAYS must be a non-negative integer" >&2
  exit 2
fi

umask 077
mkdir -p "${BACKUP_DIR}"
rm -f "${TEMP_PATH}"
trap 'rm -f "${TEMP_PATH}"' EXIT

backup_with_local_pg_dump() {
  if ! command -v pg_dump >/dev/null 2>&1 || [[ -z "${DATABASE_URL:-}" ]]; then
    return 1
  fi

  echo "Creating StockPilot backup with local pg_dump..."
  pg_dump \
    --format=custom \
    --no-owner \
    --no-privileges \
    --file="${TEMP_PATH}" \
    "${DATABASE_URL}"
}

backup_with_docker_compose() {
  if ! command -v docker >/dev/null 2>&1 || ! docker compose version >/dev/null 2>&1; then
    return 1
  fi

  if ! docker compose ps --status running --services 2>/dev/null | grep -qx 'db'; then
    return 1
  fi

  echo "Creating StockPilot backup through the Docker Compose database service..."
  docker compose exec -T db pg_dump \
    -U "${POSTGRES_USER}" \
    -d "${POSTGRES_DB}" \
    --format=custom \
    --no-owner \
    --no-privileges >"${TEMP_PATH}"
}

if ! backup_with_local_pg_dump && ! backup_with_docker_compose; then
  cat >&2 <<'EOF'
Unable to create a StockPilot backup.

Use one of these options:
  1. Install pg_dump and set DATABASE_URL.
  2. Start the repository's Docker Compose db service.
EOF
  exit 1
fi

if [[ ! -s "${TEMP_PATH}" ]]; then
  echo "Backup command completed without producing a non-empty dump" >&2
  exit 1
fi

mv "${TEMP_PATH}" "${FINAL_PATH}"
trap - EXIT

if (( BACKUP_RETENTION_DAYS > 0 )); then
  find "${BACKUP_DIR}" \
    -type f \
    -name 'stockpilot-*.dump' \
    -mtime "+${BACKUP_RETENTION_DAYS}" \
    -delete
fi

printf 'StockPilot backup created: %s\n' "${FINAL_PATH}"

#!/usr/bin/env bash
set -euo pipefail
# Simple helper to open a psql session into the dev Postgres container.
# Usage: ./tools/psql.sh [database] [user]
# Defaults: database=pontoon user=postgres
DB_NAME="${1:-pontoon}"
DB_USER="${2:-postgres}"
CONTAINER="infra-db-dev"

if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
  echo "Container ${CONTAINER} not running. Start stack with ./tools/up.sh" >&2
  exit 1
fi

docker exec -e PGPASSWORD=postgres -it "${CONTAINER}" psql -U "${DB_USER}" -d "${DB_NAME}" "$@"

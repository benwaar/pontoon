#!/usr/bin/env bash
set -euo pipefail

# Simple migration runner applying *.sql in services/game/migrations inside the db container.
# Usage: ./tools/migrate.sh [up]
# (Future: add down or status as needed.)

MODE=${1:-up}
MIGRATIONS_DIR="services/game/migrations"
DB_CONTAINER="pontoon-db-1" # docker compose default naming (verify if different)
DATABASE="pontoon"
USER="postgres"

if [[ "$MODE" != "up" ]]; then
  echo "Only 'up' supported currently" >&2
  exit 1
fi

if ! docker ps --format '{{.Names}}' | grep -q "$DB_CONTAINER"; then
  echo "Database container '$DB_CONTAINER' not running. Run ./tools/up.sh first." >&2
  exit 2
fi

echo "Applying migrations from $MIGRATIONS_DIR..."
for f in $(ls $MIGRATIONS_DIR/*.sql | sort); do
  echo "→ $f"
  docker exec -i "$DB_CONTAINER" psql -U "$USER" -d "$DATABASE" -v ON_ERROR_STOP=1 -f - < "$f"
 done

echo "Migrations applied successfully."

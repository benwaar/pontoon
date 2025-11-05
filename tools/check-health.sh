#!/bin/bash
# Check health endpoints for Phase 1 services

set -euo pipefail

function check_endpoint() {
  local url="$1"
  local name="$2"
  echo -n "Checking $name at $url ... "
  if curl -fsS "$url" > /dev/null; then
    echo "OK"
    return 0
  else
    echo "FAILED"
    return 1
  fi
}


function ensure_docker_up() {
  local running
  running=$(docker ps --format '{{.Names}}' | grep -E 'infra-keycloak-dev|infra-game-service-dev|infra-ai-service-dev|infra-db-dev' | wc -l || true)
  if [ "$running" -lt 4 ]; then
    echo "[health] Some services are not running. Starting stack via up.sh ..."
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    ROOT_DIR="$SCRIPT_DIR/.."
    (cd "$ROOT_DIR" && ./tools/up.sh)
    sleep 5
  fi
}

function check_db() {
  echo -n "Checking Postgres (SELECT 1) ... "
  if PGPASSWORD=game psql -h localhost -p 55432 -U game -d pontoon -c 'SELECT 1;' >/dev/null 2>&1; then
    echo "OK"
  else
    echo "FAILED"
    return 1
  fi
}

ensure_docker_up

check_endpoint "http://localhost:8081" "Keycloak (login page)"
check_endpoint "http://localhost:9000/healthz" "Game service"
check_endpoint "http://localhost:9001/healthz" "AI service"
check_db

echo "All health checks completed."

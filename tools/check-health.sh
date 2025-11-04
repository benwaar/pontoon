#!/bin/bash
# Check health endpoints for Phase 1 services

set -e

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
  local running=$(docker ps --format '{{.Names}}' | grep -E 'keycloak|game|ai' | wc -l)
  if [ "$running" -lt 3 ]; then
    echo "Some services are not running. Starting docker compose..."
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    ROOT_DIR="$SCRIPT_DIR/.."
    (cd "$ROOT_DIR/infra" && docker compose up -d)
    sleep 5
  fi
}

ensure_docker_up

check_endpoint "http://localhost:8080" "Keycloak (login page)"
check_endpoint "http://localhost:9000/healthz" "Game service"
check_endpoint "http://localhost:9001/healthz" "AI service"

echo "Done."

#!/usr/bin/env bash
set -euo pipefail

# test-go.sh — run Go tests for the game service with optional coverage & CGO disable.
# Flags:
#   --verbose    : verbose test output
#   --no-cgo     : set CGO_ENABLED=0 to avoid macOS linker issues
#   --no-cover   : skip coverage collection
# Usage examples:
#   ./tools/test-go.sh
#   ./tools/test-go.sh --no-cgo --verbose

ROOT_DIR="$(dirname "$0")/.."
SERVICE_DIR="$ROOT_DIR/services/game"
EXTRA_FLAGS=""
CGO_FLAG=""
COVERAGE=true

for arg in "$@"; do
	case "$arg" in
		--verbose) EXTRA_FLAGS="-v" ;;
		--no-cgo) CGO_FLAG="CGO_ENABLED=0" ;;
		--no-cover) COVERAGE=false ;;
		*) echo "Unknown flag: $arg" >&2; exit 2 ;;
	esac
done

if [[ ! -d "$SERVICE_DIR" ]]; then
	echo "Game service directory not found: $SERVICE_DIR" >&2
	exit 1
fi

pushd "$SERVICE_DIR" > /dev/null

echo "Running tests (COVERAGE=$COVERAGE CGO_FLAG='${CGO_FLAG:-default}')..."
if $COVERAGE; then
	if [[ -n "$CGO_FLAG" ]]; then export $CGO_FLAG; fi
	go test $EXTRA_FLAGS -cover -coverprofile=coverage.out ./internal/domain || { echo "Tests failed" >&2; popd > /dev/null; exit 1; }
	echo "Coverage summary (domain):"
	go tool cover -func=coverage.out | tail -n 5 || true
else
	if [[ -n "$CGO_FLAG" ]]; then export $CGO_FLAG; fi
	go test $EXTRA_FLAGS ./internal/domain || { echo "Tests failed" >&2; popd > /dev/null; exit 1; }
fi

echo "OK"
popd > /dev/null
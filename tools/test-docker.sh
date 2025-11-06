#!/usr/bin/env bash
set -euo pipefail

# test-docker.sh — run Go tests for the game service inside a Docker golang image.
# This ensures host quirks (CGO/linker) don't affect results and matches CI behavior.
# Flags:
#   --verbose    : pass -v to go test for verbose output
#   --no-cover   : skip coverage collection
#   --with-cgo   : enable CGO (default disabled)
#   --image <img>: specify alternate golang image (default golang:1.25.1)
# Usage examples:
#   ./tools/test-docker.sh
#   ./tools/test-docker.sh --verbose
#   ./tools/test-docker.sh --image golang:1.25 --with-cgo

IMAGE="golang:1.25.1"
EXTRA_FLAGS=""
COVERAGE=true
CGO_ENABLED=0

while (( "$#" )); do
  case "$1" in
    --verbose) EXTRA_FLAGS="-v" ; shift ;;
    --no-cover) COVERAGE=false ; shift ;;
    --with-cgo) CGO_ENABLED=1 ; shift ;;
    --image) IMAGE="$2" ; shift 2 ;;
    *) echo "Unknown flag: $1" >&2; exit 2 ;;
  esac
done

ROOT="$(git rev-parse --show-toplevel)"
SERVICE_DIR="services/game"
WORKDIR="/workspace/${SERVICE_DIR}"

if [[ ! -d "$ROOT/$SERVICE_DIR" ]]; then
  echo "Service directory not found: $ROOT/$SERVICE_DIR" >&2
  exit 1
fi

echo "[docker-test] Image=$IMAGE CGO_ENABLED=$CGO_ENABLED COVERAGE=$COVERAGE FLAGS='$EXTRA_FLAGS'"

TEST_CMD="go test $EXTRA_FLAGS"
if $COVERAGE; then
  TEST_CMD+=" -cover -coverprofile=coverage.out"
fi
TEST_CMD+=" ./internal/domain"

# Run tests inside the container; mount workspace read/write so coverage persists.
docker run --rm \
  -e CGO_ENABLED="$CGO_ENABLED" \
  -v "$ROOT":/workspace \
  -w "$WORKDIR" \
  "$IMAGE" bash -c "$TEST_CMD" || { echo "Tests failed" >&2; exit 1; }

if $COVERAGE; then
  # Use host 'go' if available to summarize, else run inside container.
  if command -v go >/dev/null 2>&1; then
    echo "Coverage summary (host):"
    (cd "$ROOT/$SERVICE_DIR" && go tool cover -func=coverage.out | tail -n 10 || true)
  else
    echo "Coverage summary (docker):"
    docker run --rm -v "$ROOT":/workspace -w "$WORKDIR" "$IMAGE" bash -c "go tool cover -func=coverage.out | tail -n 10" || true
  fi
fi

echo "OK (docker tests)"
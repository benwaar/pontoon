#!/bin/bash
set -euo pipefail

# Ensure goenv shims take precedence over any Homebrew Go to avoid version mismatch.
if [ -d "$HOME/.goenv" ]; then
	export PATH="$HOME/.goenv/shims:$HOME/.goenv/bin:$PATH"
fi

GO_VERSION_EXPECTED=$(cat "$(dirname "$0")/../.go-version" 2>/dev/null || cat "$(dirname "$0")/../../.go-version" 2>/dev/null || true)
if [[ -n "$GO_VERSION_EXPECTED" ]]; then
	ACTUAL=$(go version 2>/dev/null || true)
	if [[ -n "$ACTUAL" && "$ACTUAL" != *"$GO_VERSION_EXPECTED"* ]]; then
		echo "[lint-go] Warning: expected Go $GO_VERSION_EXPECTED but found: $ACTUAL"
	fi
fi

pushd "$(dirname "$0")/../services/game" > /dev/null || exit 1

echo "[lint-go] go vet"
go vet ./...

echo "[lint-go] gofmt check"
UNFORMATTED=$(gofmt -l . | grep -v vendor || true)
if [[ -n "$UNFORMATTED" ]]; then
	echo "The following files are not formatted (gofmt):"
	echo "$UNFORMATTED"
	echo "Run: gofmt -w <files> (or use your editor's format command)."
	exit 1
fi

echo "[lint-go] golangci-lint"
if ! command -v golangci-lint >/dev/null 2>&1; then
	echo "golangci-lint not found. Install with: brew install golangci-lint (macOS) or: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	exit 1
fi
golangci-lint run ./...

popd > /dev/null
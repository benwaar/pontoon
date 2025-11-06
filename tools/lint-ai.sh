#!/bin/bash
set -euo pipefail
pushd "$(dirname "$0")/../services/ai" > /dev/null || exit 1

# Ensure dev tools installed (black, isort, ruff, flake8). Lightweight check.
if ! command -v black >/dev/null 2>&1 || ! command -v ruff >/dev/null 2>&1; then
	if [ -f dev-requirements.txt ]; then
		echo "[lint-ai] Installing Python dev requirements..."
		python -m pip install -r dev-requirements.txt >/dev/null
	else
		echo "[lint-ai] Missing dev-requirements.txt; aborting."; exit 1
	fi
fi

echo "[lint-ai] black (check)"
black --check .

echo "[lint-ai] isort (check)"
isort --check-only .

echo "[lint-ai] ruff"
ruff check .

echo "[lint-ai] flake8"
flake8 .

popd > /dev/null
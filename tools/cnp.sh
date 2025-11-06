#!/bin/bash
# cnp: Commit & push with conventional commit after running linters.
# Modes:
#   Explicit: bash tools/cnp.sh <type> "message"
#   Auto:     bash tools/cnp.sh (no args) -> infers type & message
set -euo pipefail

ALLOWED_TYPES="feat fix docs style refactor perf test chore"

infer_type() {
  # If only docs changed
  if git diff --cached --name-only | grep -q '.'; then
    : # placeholder if staging already done externally
  fi
  CHANGED=$(git diff --name-only)
  if [[ -z "$CHANGED" ]]; then
    echo "chore"
    return
  fi
  if echo "$CHANGED" | grep -Eq '^docs/|\.md$'; then
    echo "docs"; return
  fi
  if echo "$CHANGED" | grep -Eq '\.(go)$'; then
    # If changes appear mostly formatting (no code hunks with keywords), still style
    if git diff | grep -E '^[+-]' | grep -vE '^(---|\+\+\+)' | grep -E 'func |type |return|if |for |switch ' >/dev/null; then
      echo "feat" # default to feat for Go logic changes
    else
      echo "style"
    fi
    return
  fi
  if echo "$CHANGED" | grep -Eq '\.(py)$'; then
    echo "feat"; return
  fi
  if echo "$CHANGED" | grep -Eq '\.(yml|yaml)$'; then
    echo "chore"; return
  fi
  echo "chore"
}

infer_message() {
  SUMMARY=$(git status --short)
  if [[ -z "$SUMMARY" ]]; then
    echo "no changes"
    return
  fi
  COUNT=$(echo "$SUMMARY" | wc -l | tr -d ' ')
  FIRST=$(echo "$SUMMARY" | head -1 | awk '{print $2}')
  if [[ $COUNT -eq 1 ]]; then
    echo "update $FIRST"
  else
    echo "update $COUNT files"
  fi
}

if [ "$#" -ge 1 ]; then
  TYPE="$1"; shift || true
  MESSAGE="$*"
  if [[ -z "$MESSAGE" ]]; then
    echo "Provide a commit message after type."; exit 1
  fi
  if ! echo "$ALLOWED_TYPES" | grep -qw "$TYPE"; then
    echo "Invalid type: $TYPE"; echo "Allowed: $ALLOWED_TYPES"; exit 1
  fi
else
  echo "[cnp] Auto mode: inferring type & message..."
  # Stage first so diff includes all changes
  git add .
  TYPE=$(infer_type)
  MESSAGE=$(infer_message)
fi

echo "[cnp] Running linters (Go/Python/YAML) before commit..."
bash tools/lint-go.sh || { echo "Go lint failed"; exit 1; }
bash tools/lint-ai.sh || { echo "Python lint failed"; exit 1; }
bash tools/lint-yaml.sh || { echo "YAML lint failed"; exit 1; }

if [ "$#" -ge 1 ]; then
  echo "[cnp] Staging changes"
  git add .
fi

GIT_COMMIT_MSG="${TYPE}: ${MESSAGE}"
echo "[cnp] Commit: $GIT_COMMIT_MSG"
git commit -m "$GIT_COMMIT_MSG" || { echo "Commit failed"; exit 1; }

echo "[cnp] Push"
git push || { echo "Push failed"; exit 1; }

echo "[cnp] Done"

#!/bin/bash
set -euo pipefail
# Lint all YAML files (excluding vendored or generated paths if needed later)
yamllint -c .yamllint.yml .

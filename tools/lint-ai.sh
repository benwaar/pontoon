#!/bin/bash
pushd "$(dirname "$0")/../services/ai" > /dev/null || exit 1
flake8 .
popd > /dev/null
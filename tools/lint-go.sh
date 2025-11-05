#!/bin/bash
pushd "$(dirname "$0")/../services/game" > /dev/null || exit 1
golangci-lint run
popd > /dev/null
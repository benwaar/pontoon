#!/bin/bash
pushd "$(dirname "$0")/../services/game" > /dev/null || exit 1
go test ./...
popd > /dev/null
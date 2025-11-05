#!/bin/bash
pushd "$(dirname "$0")/../infra" > /dev/null || exit 1
docker compose down
popd > /dev/null
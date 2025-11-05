#!/bin/bash
pushd "$(dirname "$0")/../infra" > /dev/null || exit 1
docker compose up -d
popd > /dev/null
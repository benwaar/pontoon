#!/bin/bash
pushd "$(dirname "$0")/../services/ai" > /dev/null || exit 1
pytest
popd > /dev/null
#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

go run tools/gen.go
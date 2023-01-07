#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

# Flags:
#   --local=<bool> if generated data should be for local dev.
#   --version=<string> image version to be used

go run worker_generator.go $1 $2

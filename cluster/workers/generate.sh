#!/usr/bin/env sh

# Flags:
#   --local=<bool> if generated data should be for local dev.
#   --version=<string> image version to be used

go run worker_generator.go $1 $2

#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

#if [ -z $1 ]; then
#  echo "Usage $0 <version>"
#  exit 1
#fi

#TAG=$1
TAG=0.0.1

DOCKER_BUILDKIT=1 docker build -f Dockerfile -t "shortest-path/postgres_cleaner:$TAG" ../.. || exit 1

# Only on local and if using kind
kind load docker-image "shortest-path/postgres_cleaner:$TAG"
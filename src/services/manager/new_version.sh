#!/usr/bin/env sh

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

TAG=$1

DOCKER_BUILDKIT=1 docker build -f Dockerfile -t "shortest-path/manager:$TAG" ../.. || exit 1

# Only on local and if using kind
kind load docker-image "shortest-path/manager:$TAG"
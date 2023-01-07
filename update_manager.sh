#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

VERSION=$1

./src/services/manager/new_version.sh "$VERSION"

kubectl set image --namespace shortest-path deployment/manager manager=shortest-path/manager:"$VERSION"
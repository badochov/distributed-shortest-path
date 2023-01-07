#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

VERSION=$1

./src/services/worker/new_version.sh "$VERSION"

./cluster/workers/update.sh "$VERSION"

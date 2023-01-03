#!/usr/bin/env sh

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

VERSION=$1

(
  cd src/services/worker || exit 1
  ./new_version.sh "$VERSION"
)

cd cluster/workers || exit 1
./update.sh "$VERSION"

#!/usr/bin/env sh

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

VERSION=$1

cd src/services/manager || exit 1
./new_version.sh "$VERSION"
kubectl set image --namespace shortest-path deployment/manager manager=shortest-path/manager:"$VERSION"
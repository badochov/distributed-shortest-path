#!/usr/bin/env sh

kubectl apply -f shortest-path-namespace.yaml

(
  cd postgres || exit 1
  ./deploy.sh
)

(
  cd manager || exit 1
  ./deploy.sh
)

(
  cd workers || exit 1
  ./generate --local=true
  ./deploy.sh
)
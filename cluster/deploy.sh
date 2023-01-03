#!/usr/bin/env sh

kubectl apply -f shortest-path-namespace.yaml
kubectl apply -f postgres-config.yaml

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
  ./generate.sh --local=true || exit 1
  ./deploy.sh
)
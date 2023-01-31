#!/usr/bin/env sh

set -e

cd -P -- "$(dirname -- "$0")"

if [ "$1" = "--local" ]; then
  IS_LOCAL=true
fi

kubectl apply -f shortest-path-namespace.yaml
kubectl apply -f shortest-path-commons-config.yaml
kubectl apply -f postgres-config.yaml
kubectl apply -f workers-manager-role.yaml

if [ "$IS_LOCAL" = true ]; then
  kubectl apply -f metrics-server.local.yaml
fi

./postgres/deploy.sh

./manager/deploy.sh


if [ "$IS_LOCAL" = true ]; then
  ./workers/generate.sh --local
else
  ./workers/generate.sh
fi

./workers/deploy.sh

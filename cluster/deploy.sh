#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

if [ "$1" = "--local" ]; then
  IS_LOCAL=true
fi

kubectl apply -f shortest-path-namespace.yaml
kubectl apply -f shortest-path-commons-config.yaml
kubectl apply -f postgres-config.yaml
kubectl apply -f workers-manager-role.yaml

if [ "$IS_LOCAL" = true ]; then
  kubectl apply -f metrics-server.local.yaml
else
  kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
fi

./postgres/deploy.sh

./manager/deploy.sh


if [ "$IS_LOCAL" = true ]; then
  ./workers/generate.sh --local || exit 1
else
  ./workers/generate.sh || exit 1
fi
./workers/deploy.sh

#!/usr/bin/env sh

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
  if [ "$IS_LOCAL" = true ]; then
    ./generate.sh --local || exit 1
  else
    ./generate.sh || exit 1
  fi
  ./deploy.sh
)
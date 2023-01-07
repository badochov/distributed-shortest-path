#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

kubectl apply -f manager-deployment.yaml
kubectl apply -f manager-service.yaml

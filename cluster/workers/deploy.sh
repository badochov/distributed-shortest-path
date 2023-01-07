#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

kubectl apply -f worker-service-account.yaml
kubectl apply -f generated/workers-deployment.yaml
kubectl apply -f generated/workers-hpa.yaml
kubectl apply -f generated/workers-service.yaml

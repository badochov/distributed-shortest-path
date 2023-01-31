#!/usr/bin/env sh

set -e

cd -P -- "$(dirname -- "$0")"

kubectl apply -f postgres-namespace.yaml
kubectl apply -f postgres-secrets.yaml
kubectl apply -f postgres-persistent-volume.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml
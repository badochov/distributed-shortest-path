#!/usr/bin/env sh

kubectl apply -f endpoint-observer-role.yaml
kubectl apply -f generated/workers-deployment.yaml
kubectl apply -f generated/workers-hpa.yaml
kubectl apply -f generated/workers-service.yaml

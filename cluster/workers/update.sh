#!/usr/bin/env bash

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

VERSION=$1

for i in {0..15}; do
  kubectl set image --namespace shortest-path "deployment/workers-region-$i" "worker=shortest-path/worker:$VERSION"
done


for i in {0..15}; do
  kubectl delete --namespace shortest-path "deployment/workers-region-$i"
done


for i in {0..15}; do
  kubectl delete --namespace shortest-path "hpa/workers-region-$i"
done


for i in {0..15}; do
  kubectl scale --replicas=1 --namespace shortest-path "deployment/workers-region-$i"
done
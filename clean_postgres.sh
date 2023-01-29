#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

./src/cmd/postgres_cleaner/new_version.sh || exit 1

kubectl delete -n shortest-path jobs postgres-cleaner
kubectl apply -f cluster/jobs/postgres-cleaner.yaml

#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

./src/libs/db/tools/new_version.sh

kubectl delete pods db-update-schema -n shortest-path

kubectl apply -f cluster/db-update-schema-pod.yaml
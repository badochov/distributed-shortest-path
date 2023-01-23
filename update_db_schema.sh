#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

(cd src/libs/db/ && go generate) || exit 1

kubectl delete pods db-update-schema -n shortest-path

kubectl apply -f cluster/db-update-schema-pod.yaml
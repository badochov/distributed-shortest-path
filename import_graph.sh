#!/usr/bin/env sh

cd -P -- "$(dirname -- "$0")" || exit 1

./src/cmd/osm_map_importer/new_version.sh || exit 1

kubectl delete -n shortest-path jobs osm-map-import
kubectl apply -f cluster/jobs/osm_map_importer.yaml

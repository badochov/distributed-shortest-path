#!/usr/bin/env sh

docker exec --interactive --tty docker-kind-demo \
  curl manager.shortest-path.svc.cluster.local:8080"$@"
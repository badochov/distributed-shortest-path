#!/usr/bin/env sh

docker run \
  --cap-add NET_ADMIN \
  --detach \
  --dns 10.96.0.10 \
  --dns-search svc.cluster.local \
  --dns-search cluster.local \
  --interactive \
  --name docker-kind-demo \
  --net kind \
  --rm \
  --tty \
  curlimages/curl:7.71.0 cat

ADDR="$(docker container inspect kind-control-plane --format '{{ .NetworkSettings.Networks.kind.IPAddress }}')"

docker exec --interactive --tty --user 0 \
  docker-kind-demo ip route add 10.96.0.0/12 via "$ADDR"

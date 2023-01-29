# syntax = docker/dockerfile:1-experimental

FROM golang:1.19 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY libs /go/src/libs
COPY . ./

WORKDIR /go/src/cmd/osm_map_importer

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/osm_map_importer



FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/bin/osm_map_importer /usr/bin/osm_map_importer
COPY  cmd/osm_map_importer/data/monaco.osm /data/monaco.osm

ENTRYPOINT ["/usr/bin/osm_map_importer", "--file", "/data/monaco.osm"]
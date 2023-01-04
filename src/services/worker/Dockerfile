# syntax = docker/dockerfile:1-experimental

FROM golang:1.19 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src

ENV CGO_ENABLED=0

COPY libs /go/src/libs
COPY services/worker/go.mod services/worker/go.sum services/worker/
RUN cd services/worker && go mod download

COPY services/worker /go/src/services/worker

WORKDIR /go/src/services/worker

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/worker



FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/bin/worker /usr/bin/worker

ENTRYPOINT ["/usr/bin/worker"]
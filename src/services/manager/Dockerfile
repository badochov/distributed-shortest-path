# syntax = docker/dockerfile:1-experimental

FROM golang:1.19 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src

ENV CGO_ENABLED=0

COPY libs /go/src/libs
COPY services/manager/go.mod services/manager/go.sum services/manager/
RUN cd services/manager && go mod download

COPY services/manager /go/src/services/manager

WORKDIR /go/src/services/manager

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/manager



FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/bin/manager /usr/bin/manager

ENTRYPOINT ["/usr/bin/manager"]
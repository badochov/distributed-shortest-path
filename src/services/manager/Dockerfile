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

WORKDIR /go/src/services/manager

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/manager



FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/bin/manager /usr/bin/manager

ENTRYPOINT ["/usr/bin/manager"]
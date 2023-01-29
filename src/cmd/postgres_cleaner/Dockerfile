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

WORKDIR /go/src/cmd/postgres_cleaner

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/postgres_cleaner



FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/bin/postgres_cleaner /usr/bin/postgres_cleaner

ENTRYPOINT ["/usr/bin/postgres_cleaner"]
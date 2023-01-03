FROM golang:1.19 as builder

COPY ./ /go/src

WORKDIR /go/src/services/manager

RUN go build -o /go/bin/manager

# Two step compilation reduces size from 1.3 GB to ~20 MB.
FROM busybox

COPY --from=builder /go/bin/manager /usr/bin/manager

ENTRYPOINT ["/usr/bin/manager"]
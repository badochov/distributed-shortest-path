FROM golang:1.19 as builder

COPY ./ /go/src

WORKDIR /go/src/services/worker

RUN go build -o /go/bin/worker

# Two step compilation reduces size from 1.3 GB to ~20 MB.
FROM busybox

COPY --from=builder /go/bin/worker /usr/bin/worker

ENTRYPOINT ["/usr/bin/worker"]
# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /charts
ADD https://charts.bitnami.com/bitnami/redis-12.9.0.tgz .
ADD https://charts.bitnami.com/bitnami/memcached-5.8.0.tgz .
ADD https://charts.bitnami.com/bitnami/rabbitmq-8.11.4.tgz .
ADD https://charts.bitnami.com/bitnami/kafka-12.13.2.tgz .
ADD https://charts.bitnami.com/bitnami/postgresql-10.3.13.tgz .
ADD https://charts.bitnami.com/bitnami/mysql-8.5.1.tgz .

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --chown=65532:65532 --from=builder /workspace/manager .
COPY --chown=65532:65532 --from=builder /charts /charts

USER 65532:65532

ENTRYPOINT ["/manager"]

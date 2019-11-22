# Build the manager binary
FROM registry.cn-hangzhou.aliyuncs.com/larus/golang:1.12.9 as builder

# Copy in the go src
WORKDIR /go/src/github.com/kubeway/appos/

COPY cmd/ cmd/
COPY api/ api/
COPY pkg/ pkg/
COPY ari/ ari/
COPY config/ config/
COPY vendor/ vendor/

# Build
RUN mkdir -p cmd/helloworld/bin
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o cmd/helloworld/bin/helloworld github.com/kubeway/appos/cmd/helloworld

FROM registry.cn-hangzhou.aliyuncs.com/larus/centos:7.6.1810
WORKDIR /app/
RUN mkdir -p /app/bin/
COPY --from=builder /go/src/github.com/kubeway/appos/cmd/helloworld/bin/helloworld /app/bin/helloworld
ENTRYPOINT ["/app/bin/helloworld"]


FROM golang:1.20.1-alpine3.17 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 更新下载软件
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache ca-certificates git openssh-client make bash yarn curl \
    && rm -rf /var/cache/apk/* \
    && git config --global http.version HTTP/1.1 && git config --global http.postBuffer 524288000

# 安装buf插件
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/envoyproxy/protoc-gen-validate@latest \
    && go install github.com/bufbuild/buf/cmd/buf@latest \
    && go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest \
    && go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

FROM alpine:3.17

COPY --from=builder /go/bin /usr/local/bin
#COPY ./buf-Linux-aarch64 /usr/local/bin

#RUN chmod +x "/usr/local/bin/buf-Linux-aarch64"

#ENTRYPOINT ["/usr/local/bin/buf-Linux-aarch64"]
ENTRYPOINT ["/usr/local/bin/buf"]

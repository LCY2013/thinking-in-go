FROM golang:latest
WORKDIR /root
COPY / /root
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o register
EXPOSE 9527
ENTRYPOINT ./register -consul.addr=$consulAddr -service.addr=$serviceAddr
FROM golang:latest
WORKDIR /root/user
COPY / /root/user
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
EXPOSE 10086
ENTRYPOINT ./user -mysql.addr $mysqlAddr -redis.addr $redisAddr

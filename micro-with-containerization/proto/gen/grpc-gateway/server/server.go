package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/proto/gen/grpc-gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	svr := grpc.NewServer()
	v1.RegisterGrpcGatewayServer(svr, &GrpcGatewayServerProc{})
	reflection.Register(svr)
	l, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}
	//port := l.Addr().(*net.TCPAddr).Port
	err = svr.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}

type GrpcGatewayServerProc struct {
	v1.GrpcGatewayServer
}

func (g *GrpcGatewayServerProc) Echo(ctx context.Context, message *v1.StringMessage) (*v1.StringMessage, error) {
	return &v1.StringMessage{
		Value: "hello, " + message.Value,
	}, nil
}

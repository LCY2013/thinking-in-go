// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: helloworld/helloworld.proto

package helloworld

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for HelloWorld service

func NewHelloWorldEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for HelloWorld service

type HelloWorldService interface {
	SayHello(ctx context.Context, in *SayRequest, opts ...client.CallOption) (*SayResponse, error)
}

type helloWorldService struct {
	c    client.Client
	name string
}

func NewHelloWorldService(name string, c client.Client) HelloWorldService {
	return &helloWorldService{
		c:    c,
		name: name,
	}
}

func (c *helloWorldService) SayHello(ctx context.Context, in *SayRequest, opts ...client.CallOption) (*SayResponse, error) {
	req := c.c.NewRequest(c.name, "HelloWorld.SayHello", in)
	out := new(SayResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for HelloWorld service

type HelloWorldHandler interface {
	SayHello(context.Context, *SayRequest, *SayResponse) error
}

func RegisterHelloWorldHandler(s server.Server, hdlr HelloWorldHandler, opts ...server.HandlerOption) error {
	type helloWorld interface {
		SayHello(ctx context.Context, in *SayRequest, out *SayResponse) error
	}
	type HelloWorld struct {
		helloWorld
	}
	h := &helloWorldHandler{hdlr}
	return s.Handle(s.NewHandler(&HelloWorld{h}, opts...))
}

type helloWorldHandler struct {
	HelloWorldHandler
}

func (h *helloWorldHandler) SayHello(ctx context.Context, in *SayRequest, out *SayResponse) error {
	return h.HelloWorldHandler.SayHello(ctx, in, out)
}

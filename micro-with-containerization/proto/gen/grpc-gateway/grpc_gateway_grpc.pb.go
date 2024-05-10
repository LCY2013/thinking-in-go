// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: gen/grpc-gateway/grpc_gateway.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GrpcGateway_Echo_FullMethodName = "/grpc.gateway.service.v1.GrpcGateway/Echo"
)

// GrpcGatewayClient is the client API for GrpcGateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcGatewayClient interface {
	// rpc Echo(StringMessage) returns (StringMessage) {}
	Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error)
}

type grpcGatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcGatewayClient(cc grpc.ClientConnInterface) GrpcGatewayClient {
	return &grpcGatewayClient{cc}
}

func (c *grpcGatewayClient) Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := c.cc.Invoke(ctx, GrpcGateway_Echo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcGatewayServer is the server API for GrpcGateway service.
// All implementations must embed UnimplementedGrpcGatewayServer
// for forward compatibility
type GrpcGatewayServer interface {
	// rpc Echo(StringMessage) returns (StringMessage) {}
	Echo(context.Context, *StringMessage) (*StringMessage, error)
	mustEmbedUnimplementedGrpcGatewayServer()
}

// UnimplementedGrpcGatewayServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcGatewayServer struct {
}

func (UnimplementedGrpcGatewayServer) Echo(context.Context, *StringMessage) (*StringMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedGrpcGatewayServer) mustEmbedUnimplementedGrpcGatewayServer() {}

// UnsafeGrpcGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcGatewayServer will
// result in compilation errors.
type UnsafeGrpcGatewayServer interface {
	mustEmbedUnimplementedGrpcGatewayServer()
}

func RegisterGrpcGatewayServer(s grpc.ServiceRegistrar, srv GrpcGatewayServer) {
	s.RegisterService(&GrpcGateway_ServiceDesc, srv)
}

func _GrpcGateway_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcGatewayServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcGateway_Echo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcGatewayServer).Echo(ctx, req.(*StringMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// GrpcGateway_ServiceDesc is the grpc.ServiceDesc for GrpcGateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GrpcGateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.gateway.service.v1.GrpcGateway",
	HandlerType: (*GrpcGatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _GrpcGateway_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gen/grpc-gateway/grpc_gateway.proto",
}
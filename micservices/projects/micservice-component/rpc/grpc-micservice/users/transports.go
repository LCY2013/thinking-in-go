/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-27
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package users

import (
	"context"
	"errors"
	"github.com/go-kit/kit/transport/grpc"
	"grpc-mic/pb"
)

// 定义坏请求
var ErrBadRequest = errors.New("invalid request parameter")

type grpcServer struct {
	checkPassword grpc.Handler
}

func (grpcServ *grpcServer) CheckPassword(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, response, err := grpcServ.checkPassword.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.LoginResponse), nil
}

// 创建用户服务
func NewUserServer(ctx context.Context, endpoints Endpoints) pb.UserServiceServer {
	return &grpcServer{
		checkPassword: grpc.NewServer(
			endpoints.UserEndpoint,
			decodeLoginRequest,
			encodeLoginResponse,
		),
	}
}

// 定义编码登陆响应函数
func encodeLoginResponse(ctx context.Context, resp interface{}) (response interface{}, err error) {
	// 转换请求响应
	rpcResp := resp.(*RpcResponse)
	retStr := "login fail"
	if rpcResp.Ret {
		retStr = "login success"
	}

	errStr := ""
	if rpcResp.Err != nil {
		errStr = rpcResp.Err.Error()
	}

	return &pb.LoginResponse{
		Ret: retStr,
		Err: errStr,
	}, nil
}

// 定义解码登陆请求函数
func decodeLoginRequest(ctx context.Context, req interface{}) (request interface{}, err error) {
	// 转换登陆请求到pb.LoginRequest
	loginRequest := req.(*pb.LoginRequest)
	return &RpcRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	}, nil
}

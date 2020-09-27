### grpc 特点

```
1、内置流式 RPC 支持: 这意味着你可以使用同一 RPC 框架来处理普通的 RPC 调用和分块进行的数据传输调用，这在很大程度上统一了网络相关的基础代码并简化了逻辑。

2、内置拦截器的支持: gRPC 提供了一种向多个服务端点添加通用功能的强大方法，这使得你可以轻松使用拦截器对所有接口进行共享的运行状况检查和身份验证。

3、内置流量控制和 TLS 支持: gRPC 是基于 HTTP/2 协议构建的，具有很多强大的特性。这使得客户端的实现更简单，并且可以轻松实现更多语言的绑定。

4、基于 ProtoBuf 进行数据序列化: ProtoBuf 是由 Google 开源的数据序列化协议，用于将数据进行序列化，在数据存储和通信协议等方面有较大规模的应用和成熟案例。gRPC 直接使用成熟的 ProtoBuf 来定义服务、接口和数据类型，其序列化性能、稳定性和兼容性得到保障。

5、底层基于 HTTP/2 标准设计: gRPC 正是基于 HTTP/2 才得以实现更多强大功能，如双向流、多复用请求和头部压缩等，从而可以节省带宽、降低 TCP 连接次数和提高 CPU 利用率等。同时，基于 HTTP/2 标准的 gRPC 还提高了云端服务和 Web 应用的性能，使得 gRPC 既能够在客户端应用，也能够在服务器端应用，从而实现客户端和服务器端的通信以及简化通信系统的构建。

6、优秀的社区支持: 作为一个开源项目，gRPC 拥有良好的社区支持和维护，发展迅速，并且 gRPC 的文档也很丰富，这些对用户都很有帮助。

7、提供多种语言支持: gRPC 支持多种语言，如 C、C++、Go 、Python、Ruby、Java 、PHP 、C# 和 Node.js 等，并且能够基于 ProtoBuf 定义自动生成相应的客户端和服务端代码。目前已提供了 Java 语言版本的 gRPC-Java 和 Go 语言版本的 gRPC-Go。
```

### 使用流程

#### 1、定义 proto 文件

[protobuf3官方文档](https://developers.google.com/protocol-buffers/docs/proto3)

定义两个参数结果 LoginRequest And LoginResponse    
定义一个服务结构 UserService 

protobuf 定义文件如下: user.proto
```
syntax = "proto3";
package user;

service UserService{
rpc CheckPassword(LoginRequest) returns (LoginResponse) {}
}

message LoginRequest {
string Username = 1;
string Password = 2;
}

message LoginResponse {
string Ret = 1;
string err = 2;
}
```

#### 2、编译 proto 文件

安装protoc编译器
    https://github.com/protocolbuffers/protobuf/releases
  按不同的平台下载protoc编译器，比如mac平台的: protoc-3.13.0-osx-x86_64.zip。

protoc --go_out=plugins=grpc:. user/user.proto

编译不成功没有go 语言的编译插件，需按如下流程获取protoc-gen-go:
    go get -u github.com/golang/protobuf/protoc-gen-go
  然后去GOPATH中的bin获取protoc-gen-go执行文件copy到protoc的bin下面就行。  


#### 3、编译完成的proto 得到的go源码需要使用下面的方式解决依赖问题

https://github.com/grpc/grpc-go

go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest
go mod tidy
go mod vendor
go build -mod=vendor

#### 4、客户端发送RPC请求

首先调用 grpc.Dial 建立网络连接，然后使用 protoc 编译生成的代码 pb.NewUserServiceClient 函数创建 gRPC 客户端，最后再调用客户端的 CheckPassword 函数进行 RPC 调用。

```
// grpc 服务端地址
serviceAddress := "127.0.0.1:9527"
// 创建目标连接的客户端
clientConn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
if err != nil {
    log.Printf("grpc connction err : %s\n",err)
}
// 最后关闭连接
defer clientConn.Close()

// 创建一个客户端连接服务
userClient := pb.NewUserServiceClient(clientConn)

// 定义请求参数结构体
loginRequest := &pb.LoginRequest{
    Username: "fufeng",
    Password: "123456",
}

// 调用远程服务
loginResponse, err := userClient.CheckPassword(context.Background(), loginRequest)
if err != nil {
    log.Printf("grpc erro : %s\n",err)
}

log.Printf("grpc call method CheckPassword return : %s\n",loginResponse.Ret)
```

#### 5、服务端生成RPC服务相关信息

定义一个UserService结构体和其CheckPassword函数实现

```
// 定义UserService结构体
type UserService struct {
}

// 定义结构体函数
func (userService *UserService) CheckPassword(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)  {
	if request.Username == "fufeng" && request.Password == "123456" {
		return &pb.LoginResponse{Ret: "login success",Err: ""},nil
	}
	return &pb.LoginResponse{Ret: "login fail",Err: "username and password not mach"},nil
}
```

首先需要调用 grpc.NewServer() 来建立 RPC 的服务端，然后将 UserService 注册到 RPC 服务端上，UserService 中实现了 CheckPassword 方法。

```
// 获取命令行参数信息
flag.Parse()

// 开启一个tcp监听
listen, err := net.Listen("tcp", "localhost:9527")
if err != nil {
    log.Fatalf("fail to listen: %v\n",listen)
}

// 创建grpc服务
grpcServer := grpc.NewServer()

// 创建服务
userService := new(users.UserService)

// 注册到grpc中
pb.RegisterUserServiceServer(grpcServer,userService)

// 启动监听
grpcServer.Serve(listen)    
```    





























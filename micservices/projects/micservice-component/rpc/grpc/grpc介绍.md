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
package pb;

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

protoc --go_out=plugins=grpc:. user/user.proto





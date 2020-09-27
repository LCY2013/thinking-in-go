#### Thrift 介绍 [官网](http://thrift.apache.org/)

开源的跨平台、支持多语言的成熟 RPC 框架，它通过定义中间语言（IDL） 自动生成 RPC 客户端与服务端通信代码，从而可以在 C++、Java、Python、PHP 和 Go 等多种编程语言间构建无缝结合的、高效的 RPC 通信服务。

Thrift 通过中间语言来定义 RPC 的接口和数据类型，然后通过编译器生成不同语言的代码并由生成的代码负责 RPC 协议层和传输层的实现。

#### Thrift 使用流程

##### 1、定义 Thrift 文件

定义一个用来检测用户名密码的user服务的thrift文件
```
namespace go user_service

struct LoginRequest {
1: string username;
2: string password;
}

struct LoginResponse {
1: string msg;
}

service User {
LoginResponse checkPassword(1: LoginRequest req);
}    
```

##### 2、编译 Thrift 文件

thrift -r --gen go user.thrift

可能存在不同版本的thrift的函数参数不一致情况，需要自行修改生成代码。

##### 3、客户端发送RPC请求

thrift 客户端调用
```
// 创建thrift的tSocket
tSocket, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
if err != nil {
    log.Panicf(" thrift listen error : %v \n",err)
}

// 创建thrift传输工厂
transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

transport, err := transportFactory.GetTransport(tSocket)
if err != nil{
    log.Fatalln("get transport error : ",err)
}

factoryDefault := thrift.NewTBinaryProtocolFactoryDefault()

userClient := user_service.NewUserClientFactory(transport, factoryDefault)

if err := transport.Open() ; err != nil {
    log.Fatalln("error open : ",HOST,":",PORT)
}

// 关闭传输通道
defer transport.Close()

// 构建请求
loginRequest := &user_service.LoginRequest{
    Username: "fufeng",
    Password: "123456",
}

// 检测密码是否正常
loginResponse, err := userClient.CheckPassword(context.Background(),loginRequest)

if err != nil {
    log.Fatalf("user CheckPassword call error : %s \n",err)
}

fmt.Println(loginResponse)
```

##### 4、服务端建立RPC服务

定义具体实现的结构体
```
// 定义UserService结构体
type UserService struct {
}

// 定义结构体函数
func (userService *UserService) CheckPassword(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	if request.Username == "fufeng" && request.Password == "123456" {
		return &LoginResponse{Msg: "login success"}, nil
	}
	return &LoginResponse{Msg: "username and password not mach"}, nil
}
```

thrift服务启动
```
// 创建结构体方法
userServiceHandler := &user_service.UserService{}

// 创建一个具体的处理器，将自定义实现Handler绑定到处理器上
userProcessor := user_service.NewUserProcessor(userServiceHandler)

// 创建ServerSocket
serverSocket, err := thrift.NewTServerSocket(HOST + ":" + PORT)
if err != nil {
    log.Panicf("create server socket error : %s \n",err)
}

transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
factoryDefault := thrift.NewTBinaryProtocolFactoryDefault()

server4 := thrift.NewTSimpleServer4(userProcessor, serverSocket,
    transportFactory, factoryDefault)

fmt.Println("thrift server running at : ",HOST,":",PORT)

server4.Serve()
```

总结: Thrift 可以让用户选择客户端和服务端之间进行 RPC 网络传输和序列化协议，对于服务端，还提供了不同的网络处理模型，例如本例使用的NewTFramedTransportFactory。

对于通信协议（TProtocol）: Thrift 提供了基于文本和二进制传输协议，可选的协议有：二进制编码协议（TBinaryProtocol）、压缩的二进制编码协议（TCompactProtocol）、JSON 格式的编码协议（TJSONProtocol）和用于调试的可读编码协议（TDebugProtocol）。示例中使用的是默认的二进制协议，也就是 TBinaryProtocol。

对于传输方式（TTransport）: Thrift 提供了丰富的传输方式，可选的传输方式有：最常见的阻塞式 I/O 的 TSocket、HTTP 协议传输的 THttpTransport、以 frame 为单位进行非阻塞传输的 TFramedTransport 和以内存进行传输的 TMemoryTransport 等。

对于服务端模型（TServer）: Thrift 目前提供了：单线程服务器端使用标准的阻塞式 I/O 的 TServer、多线程服务器端使用标准的阻塞式 I/O 的 TThreadedServer 和多线程网络模型使用配有线程池的阻塞式 I/O 的 TThreadPoolServer 等。














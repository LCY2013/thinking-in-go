#### 为什么说 gRPC 离生产级别的微服务还差点意思？
 
缺少了大量微服务场景下的组件功能，比如：连接池、限流熔断、配置中心、服务框架、服务发现、服务治理、分布式链路追踪、埋点和上下文日志等。

#### 利用 Go-kit 和 gRPC 结合 来完善这些实际场景下的问题

Go-kit 框架可以和 gRPC 结合使用，将 RPC 作为传输层的组件，而自身则提供诸如服务注册和发现、断路器等微服务远程交互的通用功能组件。比如，gRPC 缺乏服务治理的功能，就可以通过 Go-kit 结合 gRPC 来完善。

Go-kit 框架抽象的 Endpoint 层设计让开发者可以很容易地封装使用其他微服务组件，如：服务注册与发现、断路器和负载均衡策略等。

##### Go-kit 和 gRPC结合的相关原理

###### Go-kit 提供 Transport 层和 Endpoint 层

Transport 层：主要负责网络传输，例如处理HTTP、gRPC、Thrift等相关的逻辑。

Endpoint 层：主要负责 request/response 格式的转换，以及公用拦截器相关的逻辑。作为 Go-kit 的核心，Endpoint 层采用类似洋葱的模型，提供了对日志、限流、熔断、链路追踪和服务监控等方面的扩展能力。

Go-kit 和 gRPC 结合的关键在于需要将 gRPC 集成到 Go-kit 的 Transport 层。Go-kit 的 Transport 层用于接收用户网络请求并将其转为 Endpoint 可以处理的对象，然后交由 Endpoint 层执行，最后再将处理结果转为响应对象返回给客户端。为了完成这项工作，Transport 层需要具备两个工具方法：

解码器：把用户的请求内容转换为请求对象。

编码器：把处理结果转换为响应对象。

具体流程如下:
```
client request     -(grpc)->  [ decode  encode ]   -(grpc)->   server response
                               request  response
                                  |       |                 
                              [ ...endpoint... ]
                                  |       |
                              [ ...service...]
```

##### Go-kit 集成 gRPC 的项目实践

项目前准备，定义proto文件，生成ProtoBuffer源码:
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

protoc --go_out=plugins=grpc:. user.proto

go mod 文件执行如下命令: 

go mod edit -replace=google.golang.org/grpc=github.com/grpc/grpc-go@latest

go mod tidy

1、定义Service，业务实现
```
// 定义服务接口
type UserService interface {
	CheckPassword(ctx context.Context, userName, password string) (bool, error)
}

type UserServiceImpl struct {
}

// 定义服务实现
func (userService *UserServiceImpl) CheckPassword(ctx context.Context, userName, password string) (bool, error) {
	if userName == "fufeng" && password == "123456" {
		return true,nil
	}
	return false,nil
}
```

2、定义 Endpoint，提供参数转换能力
```
// 定义请求结构体
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 定义响应结构体
type LoginResponse struct {
	Ret bool `json:"ret"`
	Err error `json:"err"`
}

// 定义go kit 的端点
type Endpoints struct {
	UserEndpoint endpoint.Endpoint
}

// 创建User Endpoint
func MakeUserEndpoint(svc UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 强制转换登陆请求结构体
		loginRequest := request.(LoginRequest)
		// 调用业务函数
		isLogin, err := svc.CheckPassword(ctx,
			loginRequest.Username, loginRequest.Password)
		return &LoginResponse{Ret: isLogin,Err: err},err
	}
}
```

3、定义 Middleware，提供限流和日志中间件
```
// 定义限流异常
var ErrLimitExceed = errors.New("rate limit exceed")

// 创建限流中间件(令牌桶限流策略)
func NewTokenBucketLimiterWithBuildIn(limiter *rate.Limiter) endpoint.Middleware {
	return func(endpoint endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// 超过设定的流量阀值，就启动限流
			if !limiter.Allow() {
				return nil,ErrLimitExceed
			}
			// 执行业务逻辑
			return endpoint(ctx,request)
		}
	}
}

// 使用时的代码实例
ratebucket := rate.NewLimiter(rate.Every(time.Second * 1), 1000) 
endpoint = user.NewTokenBucketLimitterWithBuildIn(ratebucket)(endpoint)
```

```
// 定义类型服务中间件为函数类型
type ServiceMiddleware func(service UserService) UserService

// 定义日志中间件结构体
type loggingMiddleware struct {
	UserService
	logger log.Logger
}

// 定义函数含有日志组件的服务中间件
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next UserService) UserService {
		return &loggingMiddleware{
			next, logger,
		}
	}
}

// 实现接口方法
func (logMiddleware *loggingMiddleware) CheckPassword(ctx context.Context, userName, password string) (bool, error) {
	// 最后记录方法调用时间
	defer func(begin time.Time) {
		logMiddleware.logger.Log(
			"function", "CheckPassword",
			"userName", userName,
			"password", password,
			"took", time.Since(begin),
		)
	}(time.Now())
	// 具体服务调用
	isOk, err := logMiddleware.UserService.CheckPassword(ctx, userName, password)
	return isOk, err
}
```

4、定义 Transport，提供网络传输能力
```
// 定义坏请求
var ErrBadRequest = errors.New("invalid request parameter")

type grpcServer struct {
	checkPassword grpc.Handler
}

func (grpcServ *grpcServer) CheckPassword(ctx context.Context,req *pb.LoginRequest) (*pb.LoginResponse,error) {
	_, response, err := grpcServ.checkPassword.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.LoginResponse),nil
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
	},nil
}

// 定义解码登陆请求函数
func decodeLoginRequest(ctx context.Context, req interface{}) (request interface{}, err error) {
	// 转换登陆请求到pb.LoginRequest
	loginRequest := req.(*pb.LoginRequest)
	return &RpcRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	},nil
}
```

5、启动服务端，注册RPC服务
```
// 获取命令行参数信息
flag.Parse()

// 定义日志相关参数信息
var logger log.Logger
{
    logger = log.NewLogfmtLogger(os.Stderr)
    logger = log.With(logger,log.DefaultTimestampUTC)
    logger = log.With(logger,log.DefaultCaller)
}

// 定义应用上下文信息
ctx := context.Background()

// 创建Service
var svc users.UserService
// 创建UserService实现
svc = &users.UserServiceImpl{}

// 构建日志中间件
service := users.LoggingMiddleware(logger)(svc)

// 创建Endpoint
endpoint := users.MakeUserEndpoint(service)
// 构建限流中间件
limiter := rate.NewLimiter(rate.Every(time.Second*1), 200)
endpoint = users.NewTokenBucketLimiterWithBuildIn(limiter)(endpoint)

endpoints := users.Endpoints{
    endpoint,
}

// 构建UserService
userServiceServer := users.NewUserServer(ctx, endpoints)

// grpc 启动，监听端口、注册grpc服务信息
listen, err := net.Listen("tcp", "localhost:9527")
if err != nil {
    logGo.Printf("gpc listen err : %s\n", err)
}
server := grpc.NewServer()
pb.RegisterUserServiceServer(server,userServiceServer)
_ = server.Serve(listen)
```

6、启动客户端调用RPC服务
```
// 定义RPC服务地址
serviceAddress := "localhost:9527"
// 创建客户端连接
conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
if err != nil {
    panic("grpc client connect err")
}
defer conn.Close()

userServiceClient := pb.NewUserServiceClient(conn)
ret, err := userServiceClient.CheckPassword(context.Background(), &pb.LoginRequest{
    Username: "fufeng",
    Password: "123456",
})
if err != nil {
    panic(err)
}
fmt.Println("check password status : ",ret.Ret)
```











#### 什么事GateWay?

在微服务架构中，网关位于接入层之下和业务服务层之上,它是微服务架构中的一个基础服务。

微服务网关就是一个处于应用程序或服务之前的系统，用来处理管理授权、访问控制和流量限制等功能。

#### GateWay职责有哪些?

请求接入: 管理所有接入请求，作为所有 API 接口的请求入口。在生产环境中，为了保护内部系统的安全性，往往内网与外网都是隔离的，服务端应用都是运行在内网环境中，为了安全，一般不允许外部直接访问。网关可以通过校验规则和配置白名单，对外部请求进行初步过滤，这种方式更加动态灵活。

统一管理: 可以提供统一的监控工具、配置管理和接口的 API 文档管理等基础设施。例如，统一配置日志切面，并记录对应的日志文件。

解耦: 可以使得微服务系统的各方能够独立、自由、高效、灵活地调整，而不用担心给其他方面带来影响。软件系统的整个过程中包括不同的角色，有服务的开发提供方、服务的用户、运维人员、安全管理人员等，每个角色的职责和关注点都不同。微服务网关可以很好地解耦各方的相互依赖关系，让各个角色的用户更加专注自己的目标。

拦截插件: 服务网关层除了处理请求的路由转发外，还需要负责认证鉴权、限流熔断、监控和安全防范等，这些功能的实现方式，往往随着业务的变化不断调整。这就要求网关层提供一套机制，可以很好地支持这种动态扩展。拦截策略提供了一个扩展点，方便通过扩展机制对请求进行一系列加工和处理。同时还可以提供统一的安全、路由和流控等公共服务组件。

#### GateWay 如何实现?

##### 目标

API 网关根据客户端 HTTP 请求，动态查询注册中心的服务实例，通过反向代理实现对后台服务的调用。

API 网关将符合规则的请求路由调用对应的后端服务。

##### 思路

1、定义需要实现的网关转发规则

```
网关实现的规则有很多种，如 HTTP 请求的资源路径、方法、头部和参数等。
本例以最简单的请求路径为例，规则为 ：/{serviceInstance}/uri。
路径第一部分为注册中心服务实例名称，后面为服务实例的 REST URI 路径。

例如:
    /user-service/register
    /user-service/login

user-service  服务实例
register , login 服务实例提供的URI
 
```

2、思路分析

客户端向网关发起请求，网关解析请求资源路径中的信息，根据服务实例名称查询注册中心的服务实例；然后使用反向代理技术把客户端请求转发至后端真实的服务实例，请求执行完毕后，再把响应信息返回客户端。

3、本例网关主要功能实现技术概览

HTTP请求规则 /{serviceInstanceName}/URI，否则不识别。

使用 Go 提供的反向代理包 httputil.ReverseProxy 实现一个简单的反向代理，它能够对请求实现负载均衡，随机地把请求发送给服务实例。

使用 Consul 客户端 API 动态查询服务实例。

4、主方法实现
```
// 获取命令行环境参数
var (
    consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
    consulPort = flag.String("consul.port", "127.0.0.1", "consul server port")
)
// 处理环境参数
flag.Parse()

// 创建日志组件，设置日志组件相关内容
var logger log.Logger
{
    logger = log.NewLogfmtLogger(os.Stderr)
    logger = log.With(logger, log.DefaultTimestampUTC)
    logger = log.With(logger, log.DefaultCaller)
}

// 创建consul api 的客户端信息
consulConfig := api.DefaultConfig()
// 设置consul的访问地址
consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
// 利用consulConfig创建一个consul客户端
consulClient, err := api.NewClient(consulConfig)
if err != nil {
    _ = logger.Log("log", err)
    // 定制退出码，用于容器退出捕获
    os.Exit(7)
}

// 创建反向代理，通过传入consul的client端以及日志组件
proxy := NewReverseProxy(consulClient,logger)

errs := make(chan error)
go func() {
    c := make(chan os.Signal)
    signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)
    errs <- fmt.Errorf("%s",<-c)
}()

// 启动协程监听端口网络服务
go func() {
    _ = logger.Log("transport", "HTTP", "addr", "9527")
    errs <- http.ListenAndServe(":9527",proxy)
}()

// 等待退出指令并且打印退出请求
logger.Log("exit",<-errs)
```

5、反向代理方法实现
```
// 创建ReverseProxy需要的Director
director := func(req *http.Request) {
    // 获取请求原始路径
    reqPath := req.URL.Path
    if reqPath == "" {
        return
    }

    // 按照"/"对路径进行分割，获取到实例名称ServiceName
    reqs := strings.Split(reqPath, "/")
    serviceName := reqs[1]

    // 根据服务名称去Consul查询所有的服务实例
    services, _, err := consulClient.Catalog().Service(serviceName, "", nil)
    if err != nil {
        _ = logger.Log("reverse proxy fail", "no such service instance", err.Error())
        return
    }

    // 判断服务实例的数量
    if len(services) == 0 {
        _ = logger.Log("reverse proxy fail","no such service instance",serviceName)
        return
    }

    // 重新组织请求路径，去掉原有的ServiceName
    realRequestPath := strings.Join(reqs[2:], "/")

    // 随机选择一个服务实例
    targetInstance := services[rand.Int()%len(services)]
    // 添加日志打印，输出获取到的目标实例ID
    logger.Log("service id",targetInstance.ServiceID)

    // 设置代理相关的配置信息
    req.URL.Scheme = "http"
    req.URL.Host = fmt.Sprintf("%s:%d",targetInstance.ServiceAddress,targetInstance.ServicePort)
    req.URL.Path = realRequestPath
}

// 返回一个反向代理实现
return &httputil.ReverseProxy{
    Director: director,
}
```















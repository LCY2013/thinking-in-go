Go-kit 是一套强大的微服务开发工具集，用于指导开发人员解决分布式系统开发过程中所遇到的问题，帮助开发人员更专注于业务开发。
Go-kit 推荐使用 transport、endpoint 和 service 3 层结构来组织项目，它们的作用分别为：
    1、transport 层，指定项目提供服务的方式，比如 HTTP 或者 gRPC 等 。
    2、endpoint 层，负责接收请求并返回响应。对于每一个服务接口，endpoint 层都使用一个抽象的 Endpoint 来表示 ，
        我们可以为每一个 Endpoint 装饰 Go-kit 提供的附加功能，如日志记录、限流、熔断等。
    3、service 层，提供具体的业务实现接口，endpoint 层中的 Endpoint 通过调用 service 层的接口方法处理请求。






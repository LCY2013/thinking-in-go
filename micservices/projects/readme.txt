Go-kit 是一套强大的微服务开发工具集，用于指导开发人员解决分布式系统开发过程中所遇到的问题，帮助开发人员更专注于业务开发。
Go-kit 推荐使用 transport、endpoint 和 service 3 层结构来组织项目，它们的作用分别为：
    1、transport 层，指定项目提供服务的方式，比如 HTTP 或者 gRPC 等 。
    2、endpoint 层，负责接收请求并返回响应。对于每一个服务接口，endpoint 层都使用一个抽象的 Endpoint 来表示 ，
        我们可以为每一个 Endpoint 装饰 Go-kit 提供的附加功能，如日志记录、限流、熔断等。
    3、service 层，提供具体的业务实现接口，endpoint 层中的 Endpoint 通过调用 service 层的接口方法处理请求。


DDD:
    1、限界上下文主要有以下几种映射方式。
        合作关系（Partnership）：两个上下文紧密合作的关系，一荣俱荣，一损俱损。

        共享内核（Shared Kernel）：两个上下文依赖部分共享的模型。

        防腐层（Anticorruption Layer）：一个上下文通过一些适配和转换与另一个上下文交互。

        客户方-供应方开发（Customer-Supplier Development）：上下文之间有组织的上下游依赖。

        开放主机服务（Open Host Service）：定义一种协议来让其他上下文对本上下文进行访问。

        遵奉者（Conformist）：下游上下文只能盲目依赖上游上下文。

        发布语言（Published Language）：通常与 OHS 一起使用，用于定义开放主机的协议。

        大泥球（Big Ball of Mud）：混杂在一起的上下文关系，边界不清晰。

        另谋他路（Separate Way）：两个完全没有任何联系的上下文。
    实践中，主要通过防腐层映射不同的限界上下文中相同的领域对象，保证整体领域概念的完整和统一。























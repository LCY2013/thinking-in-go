#### service mesh起源
[service mesh](https://philcalcado.com/2017/08/03/pattern_service_mesh.html)

#### 何为service mesh
service mesh就是sidecar网络拓扑模式
目前的service mesh包含两个主要部分:
    1、数据平面
    2、控制平面

#### service mesh主要功能点
1、流控(蓝绿部署、灰度发布、A/B测试)
    a) 路由
    b) 流量转移
    c) 超时重试
    d) 熔断
    e) 故障注入
    f) 流量镜像
2、策略
    a) 流量限制
    b) 黑白名单
3、网络安全
    a) 授权已经身份认证 
4、可观测型
    a) 指标收集和展示 
    b) 日志收集
    c) 分布式追踪

#### service mesh 与 kubernetes ?
1、service mesh 
    解决服务间网络通信问题
    本质上是管理服务通信(代理)
    对kubernetes网络功能的扩展和延伸
2、kubernetes
    解决容器编排与调度问题
    本质上是管理应用生命周期(调度)     
    给予service mesh
    
#### service mesh 标准
1、UDPA(统一数据平面API)  
2、SMI(服务网格接口)
 
#### service mesh 产品
数据平面
1、Linkerd 
2、Envoy
新增控制平面
3、Istio 
4、Conduit
5、AWS App mesh GA
6、Google Traffic Director beta
7、kong kuma
8、蚂蚁 Mosn
 
 
 
 
 
 
 
 
        





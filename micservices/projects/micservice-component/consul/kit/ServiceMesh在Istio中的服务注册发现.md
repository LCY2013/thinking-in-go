### 第一部分: [什么是Service Mesh?](https://philcalcado.com/2017/08/03/pattern_service_mesh.html)

Buoyant’s CEO William Morgan对Service Mesh作出的定义:
```
    A service mesh is a dedicated infrastructure layer for handling service-to-service communication. 
It’s responsible for the reliable delivery of requests through the complex topology of services that comprise a modern, cloud native application. 
In practice, the service mesh is typically implemented as an array of lightweight network proxies that are deployed alongside application code, without the application needing to be aware.
```

理想中: 
```
90年代，sun公司关于分布式系统的八大谬论
    1、The network is reliable
    2、Latency is zero
    3、Bandwidth is infinite
    4、The network is secure
    5、Topology doesn’t change
    6、There is one administrator
    7、Transport cost is zero
    8、The network is homogeneous
```

现实中:
```
微服务架构需要解决的问题
    1、Rapid provisioning of compute resources
    2、Basic monitoring
    3、Rapid deployment
    4、Easy to provision storage
    5、Easy access to the edge
    6、Authentication/Authorisation
    7、Standardised RPC
```

### 第二部分 Service Mesh 具体能做指什么？
   Service Mesh 作为下一代的微服务架构，它将服务间的通信从基础设施中抽离出来，达到交付更可靠的应用请求、监控和控制流量的目的。
   Service Mesh 一般与应用程序一同部署(SideCar)，作为"数据平面(Data Plane)"代理网络以及"控制平面(Control Plane)"代替应用与其他代理交互。
   Service Mesh 的出现让业务开发人员从基础架构的底层细节中解放出来，从而把更多的精力放在业务开发上，提高需求迭代的效率。

![Service Mesh 架构](https://philcalcado.com/img/service-mesh/6-a.png)

### 第三部分 主流的Service Mesh产品有哪些？
[CNCF Cloud Native Interactive Landscape](https://landscape.cncf.io/)
   Isotio
   Linkerd
   Kuma
   ...

### 第四部分 [Istio](https://istio.io/)    
```
    目前Istio 作为 Service Mesh 的落地产品之一，依托 Kubernetes 快速发展，已经成为最受欢迎的 Service Mesh 之一。
```
Istio 在逻辑上分为数据平面和控制平面
    1、数据平面: 由一组高性能的智能代理（基于 Envoy 改进的 istio-proxy）组成，它们控制和协调了被代理服务的所有网络通信，同时也负责收集和上报相关的监控数据。
    2、控制平面: 负责制定应用策略来控制网络流量的路由。

Istio 由多个组件组成，核心组件及其作用为如下：
    1、Ingressgateway: 控制外部流量访问 Istio 内部的服务。
    2、Egressgateway: 控制 Istio 内部访问外部服务的流量。
    3、Pilot: 负责管理服务网格内部的服务和流量策略。它将服务信息和流量控制的高级路由规则在运行时传播给 Proxy，并将特定平台的服务发现机制抽象为 Proxy 可使用的标准格式。
    4、Citadel: 提供身份认证和凭证管理。
    5、Galley: 负责验证、提取、处理和分发配置。
    6、Proxy: 作为服务代理，调节所有 Service Mesh 单元的入口和出口流量。
    
Proxy 属于数据平面，以 Sidecar 的方式与应用程序一同部署到 Pod 中，而 Pilot、Citadel 和 Galley 属于控制平面。
除此之外，Istio 中还提供一些额外的插件，如 grafana、istio-tracing、kiali 和 prometheus，用于进行可视化的数据查看、流量监控和链路追踪等。
The components marked as X are installed within each profile:
	                    default 	demo	minimal	    remote
Core components				
istio-egressgateway		             X		
istio-ingressgateway	  X          X		
istiod	                  X	         X	      X	

istiod 组件封装了 Pilot、Citadel 和 Galley 等控制平面组件，将它们进行统一打包部署，降低多组件维护和管理的困难性。
从上表可以看出，demo profile 是功能最全的配置清单，适合于学习和功能演示。
preview profile 将可能使用一些开发阶段的测试组件，开启的组件不定。
官方推荐使用 default profile 进行安装，因为它在核心组件和插件上做到了最优的选择，比如组件只开启了 Ingressgateway 和 istiod，插件只开启了 prometheus。
```
根据实践的需求选择合适的 profile 进行安装启动，比如下面的安装命令我们使用的是 demo profile:
    istioctl manifest apply --set profile=demo 
上述命令以 demo profile 部署 Istio，该配置下的 Istio 能够通过可视化界面监控 Istio 中应用的方方面面。
Istio 以 Sidecar 的方式在应用程序运行的 Pod 中注入 Proxy，全面接管应用程序的网络流入流出。
可以通过标记 Kubernetes 命名空间的方式，让 Sidecar 注入器自动将 Proxy 注入在该命名空间下启动的 Pod 中，开启标记的命令如下：
    kubectl label namespace default istio-injection=enabled 
上述命令中，就是将 default 命名空间标记为 istio-injection。
如果不想开启命令空间的标记，也可以通过 istioctl kube-inject 为 Pod 注入 Proxy Sidecar 容器。
接下来，我们就为 register 服务所在的 Pod 注入 Proxy，启动命令如下：
    istioctl kube-inject -f register-service.yaml | kubectl apply -f - 
    istioctl kube-inject -f register-service.yaml | kubectl delete -f -
register-service.yaml 服务的 yaml 配置如下：
apiVersion: apps/v1 
kind: Deployment 
metadata: 
  name: register 
  labels: 
    name: register 
spec: 
  selector: 
    matchLabels: 
      name: register 
  replicas: 1 
  template: 
    metadata: 
      name: register 
      labels: 
        name: register 
        app: register  # 添加 app 标签 
    spec: 
      containers: 
        - name: register 
          image: registry.cn-hangzhou.aliyuncs.com/luochunyun/register 
          ports: 
            - containerPort: 9527 
          imagePullPolicy: IfNotPresent 
          env:
          - name: consulAddr
            valueFrom:
              fieldRef:
                fieldPath: status.hostIP
#                  fieldPath: spec.nodeName
          - name: serviceAddr
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
--- 
# 添加 Service 资源 
apiVersion: v1 
kind: Service 
metadata: 
  name: register-service 
  labels: 
    name: register-service 
spec: 
  selector: 
    name: register 
  ports: 
    - protocol: TCP 
      port: 9527 
      targetPort: 9527 
      name: register-service-http 

主要的改动有：为 register 服务添加 Deployment Controller，添加了新的标签 app，以及为 register 添加相应的 Service 资源。
如果在部署 Istio 时启动了 kiali 插件，即可在 kiali 平台中查看到 register 服务的相关信息，通过以下命令即可打开 kiali 控制面板，默认账户和密码都为 admin：
    istioctl dashboard kiali 
    istioctl dashboard kiali --address 192.168.99.120
卸载kiali
    kubectl delete -f kiali.yaml --ignore-not-found
kiali控制台解析:
    Overview，网格概述，展示 Istio 内具有服务的所有命名空间；
    
    Graph，服务拓扑图；
    
    Applications，应用维度，识别设置了 app 标签的应用；
    
    Workloads，负载维度，检测 Kubernetes 中的资源，包括 Deployment、Job、DaemonSet 等，无论这些资源有没有加入 Istio 中都能检测到；
    
    Services，服务维度，检测 Kubernetes 的 Service；
    
    Istio Config，配置维度，查看 Istio 相关配置类信息。

```

Istio 依托 Kubernetes 的快速发展和推广，对 Kubernetes 有着极强的依赖性，其服务注册与发现的实现也主要依赖于 Kubernetes 的 Service 管理。
![Istio架构图](https://istio.io/latest/docs/ops/deployment/architecture/arch.svg)

由上图可以看出Istio服务注册与发现参与的模块如下:
    1、ConfigController：负责管理配置数据，包括用户配置的流量管理和路由规则。
    2、ServiceController：负责加载各类 ServiceRegistry，从 ServiceRegistry 中同步需要在网格中管理的服务。
        主要包含：
        1、KubeServiceRegistry，从 Kubernetes 同步 Service 和 Endpoint 到 Istio；
        2、ConsulServiceRegistry，从 Consul 中同步服务信息到 Istio；
        3、ExternalServiceRegistry，监听 ConfigController 中的配置变化，获取 ServiceEntry 和 WorkloadEntry 资源并封装成服务数据提供给 ServiceController。
    3、DiscoveryServer：负责将 ConfigController 中的路由配置信息和 ServiceController 中的服务信息封装成 Proxy 可以理解的标准格式，并下发到 Proxy 中。


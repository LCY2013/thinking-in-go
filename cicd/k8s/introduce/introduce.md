#### kubernetes 简介
```
    Kubernetes是由 Google 开源的，目的是管理公司内部运行的成千上万的服务器，降低应用程序部署管理的成本。
Kubernetes 将基础设施抽象，简化了应用开发、部署和运维等工作，提高了硬件资源的利用率，是一款优秀的容器管理和编排系统。
```
#### kubernetes 主要组件
```
    节点分类: Master节点、Node节点
    Master节点(负责管理和控制):
        1、API Server: 对外提供 Kubernetes 的服务接口，供各类客户端使用
        2、Scheduler: 负责对集群内部的资源进行调度，按照预设的策略将 Pod 调度到相应的 Node 节点
        3、Controller Manager: 作为管理控制器，负责维护整个集群的状态
        4、etcd: 保存整个集群的状态数据
    Node节点(工作节点):
        1、Pod: Kubernetes 创建和部署的基本操作单位，它代表了集群中运行的一个进程，
            内部由一个或者多个共享资源的容器组成，我们可以简单将 Pod 理解成一台虚拟主机，
            主机内的容器共享网络、存储等资源
        2、Docker: 是 Pod 中最常见的容器 runtime，Pod 也支持其他容器 runtime
        3、Kubelet: 负责维护调度到它所在 Node 节点的 Pod 的生命周期，包括创建、修改、删除和监控等
        4、Kube-proxy: 负责为 Pod 提供代理，为 Service 提供集群内部的服务发现和负载均衡，Service 可以看作一组提供相同服务的 Pod 的对外访问接口
```









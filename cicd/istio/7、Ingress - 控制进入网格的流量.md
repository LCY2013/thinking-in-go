- [部署 Bookinfo](https://istio.io/latest/docs/examples/bookinfo/)

- [示例地址](https://istio.io/latest/zh/docs/tasks/traffic-management/ingress/ingress-control/)

### Ingress 基本概念
- 服务的访问入口，接收外部请求并转发到后端服务

- Istio 的 Ingress gateway 和 Kubernetes Ingress 的区别
  
  - Kubernetes: 针对L7协议(资源受限)，可定义路由规则
    
  - Istio: 针对 L4-6 协议，只定义接入点，复用 Virtual Service 的 L7 路由定义


### 创建 Ingress 网关
- 任务说明
  - 为 httpbin 服务配置 Ingress 网关
    
- 任务目标
  - 理解 Istio 实现自己的 Ingress 的意义
  - Gateway 的配置方法
  - Virtual Service 的配置方法

### 操作
- 部署 httpbin 服务

> kubectl apply -f samples/httpbin/httpbin.yaml

>  kubectl get pod
```yaml
NAME                                      READY   STATUS    RESTARTS   AGE
httpbin-74fb669cc6-z2rmq                  2/2     Running   0          82s
```

- 确定 Ingress IP 和端口
> kubectl get svc istio-ingressgateway -n istio-system

- 定义 Ingress gateway

```yaml
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "httpbin.example.com"
EOF
```

- 定义对应的 Virtual Service 

```yaml
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - "httpbin.example.com"
  gateways:
  - httpbin-gateway
  http:
  - match:
    - uri:
        prefix: /status
    - uri:
        prefix: /delay
    route:
    - destination:
        port:
          number: 8000
        host: httpbin
EOF
```
  
- 测试

> curl -s -I -HHost:httpbin.example.com "http://$INGRESS_HOST:$INGRESS_PORT/status/200"


### 修改任务中的配置，使其不受域名限制 
- 提示:hosts 使用通配符 *
  
```yaml
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - "*"
  gateways:
  - httpbin-gateway
  http:
  - match:
    - uri:
        prefix: /headers
    route:
    - destination:
        port:
          number: 8000
        host: httpbin
EOF
```

- 清空示例
> kubectl delete gw httpbin-gateway

> kubectl delete vs  httpbin

- 思考为什么 Istio 不使用 Kubernetes 的 Ingress 而要自己实现它?

> Gateway 配置资源允许外部流量进入 Istio 服务网格，并对边界服务实施流量管理和 Istio 可用的策略特性。

### 问题排查
- 检查环境变量 INGRESS_HOST and INGRESS_PORT。确保环境变量的值有效，命令如下：
```shell
$ kubectl get svc -n istio-system
$ echo "INGRESS_HOST=$INGRESS_HOST, INGRESS_PORT=$INGRESS_PORT"
```

- 检查没有在相同的端口上定义其它 Istio Ingress Gateways：
```shell
$ kubectl get gateway --all-namespaces
```

- 检查没有在相同的 IP 和端口上定义 Kubernetes Ingress 资源：
```shell
$ kubectl get ingress --all-namespaces
```

- 如果使用了外部负载均衡器，该外部负载均衡器无法正常工作，尝试 [通过 Node Port 访问 Gateway](https://istio.io/latest/zh/docs/tasks/traffic-management/ingress/ingress-control/#determining-the-ingress-i-p-and-ports) 。

### 清空
```shell
$ kubectl delete gateway httpbin-gateway
$ kubectl delete virtualservice httpbin
$ kubectl delete --ignore-not-found=true -f samples/httpbin/httpbin.yaml
```
- [部署 Bookinfo](https://istio.io/latest/docs/examples/bookinfo/) 

```text
注入 Sidecar
自动:$ kubectl label namespace default istio-injection=enabled
手动: $ kubectl apply -f < (istioctl kube-inject -f samples/bookinfo/platform/kube/bookinfo.yaml)

部署应用
$ kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
确认服务、Pod 已启动

查看部署的资源定义
$ more samples/bookinfo/platform/kube/bookinfo.yaml

查看pod情况,(每一个pod都包含两个服务一个是应用服务，一个是sidecar服务)
$ kubectl get pod 
details-v1-79f774bdb9-g2nc4               0/2     PodInitializing   0          35s
productpage-v1-6b746f74dc-479jc           0/2     PodInitializing   0          35s
ratings-v1-b6994bb9-nxb9z                 0/2     PodInitializing   0          35s
reviews-v1-545db77b95-tpklj               0/2     PodInitializing   0          34s
reviews-v2-7bf8c9648f-jph9h               0/2     PodInitializing   0          35s
reviews-v3-84779c7bbc-ngl86               0/2     PodInitializing   0          35s

查看某个pod的具体情况
$ kubectl describe po details-v1-79f774bdb9-g2nc4 | less 

查看目前服务信息
$ kubectl get services

通过命令行验证
$ kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"

开启网关应用让外部流量进入

查看istio网关配置信息
$ more samples/bookinfo/networking/bookinfo-gateway.yaml

应用istio网关
$ kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml

确认配置没有问题
$ istioctl analyze

确认ingress ip和端口 （Determining the ingress IP and ports）

Execute the following command to determine if your Kubernetes cluster is running in an environment that supports external load balancers:
执行以下命令来确定Kubernetes集群运行在一个支持的环境外部负载平衡器:
$ kubectl get svc istio-ingressgateway -n istio-system
NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)                                                                      AGE
istio-ingressgateway   LoadBalancer   10.97.138.99   <pending>     15021:31387/TCP,80:30371/TCP,443:30829/TCP,31400:30312/TCP,15443:32026/TCP   49m

Set the ingress IP and ports:
设置ingree ip 和 端口:
$ export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
$ export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
$ export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')

In certain environments, the load balancer may be exposed using a host name, instead of an IP address. In this case, the ingress gateway’s EXTERNAL-IP value will not be an IP address, but rather a host name, and the above command will have failed to set the INGRESS_HOST environment variable. Use the following command to correct the INGRESS_HOST value:
如果有的环境中暴露的是主机名称而不是ip地址，可以用下面命令解决
$ export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')

Follow these instructions if your environment does not have an external load balancer and choose a node port instead.
Set the ingress ports:
$ export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
$ export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')

Set GATEWAY_URL:
设置网关地址：
$ export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

输出网关地址
Ensure an IP address and port were successfully assigned to the environment variable:
$ echo "$GATEWAY_URL"

Verify external access
检测外部访问:
$ echo "http://$GATEWAY_URL/productpage"

============== 在dashboard中展示 ================
View the dashboard

********************  dashboard  *******************

1、Install Kiali and the other addons and wait for them to be deployed.
$ kubectl apply -f samples/addons
$ kubectl rollout status deployment/kiali -n istio-system

2、Access the Kiali dashboard.
$ istioctl dashboard kiali
$ istioctl dashboard --address 192.168.0.180 -p 20001 kiali



3、In the left navigation menu, select Graph and in the Namespace drop down, select default.
To see trace data, you must send requests to your service. The number of requests depends on Istio’s sampling rate. You set this rate when you install Istio. The default sampling rate is 1%. You need to send at least 100 requests before the first trace is visible. To send a 100 requests to the productpage service, use the following command:
$ for i in $(seq 1 100); do curl -s -o /dev/null "http://$GATEWAY_URL/productpage"; done

如遇到下面问题：
2021-06-12T18:56:27.060054Z	error	klog	an error occurred forwarding 20001 -> 20001: error forwarding port 20001 to pod a49b498b6a4d6543adda5d8df7bb1af4ac7b164b557ef074bff741170ed7a744, uid : unable to do port forwarding: socat not found
centos8 操作如下:
$ yum -y install socat
```
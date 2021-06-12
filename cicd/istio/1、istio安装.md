### [安装](https://istio.io/latest/docs/setup/getting-started/)

- 获取istio相关信息

> curl -L https://istio.io/downloadIstio | sh -

- 进入istio目录

> cd istio-1.10.1

- 添加istio的命令到系统环境变量中

> export PATH=$PWD/bin:$PATH

- 开始安装istio 到k8s集群

> istioctl install --set profile=demo -y
```text
Detected that your cluster does not support third party JWT authentication. Falling back to less secure first party JWT. See https://istio.io/v1.10/docs/ops/best-practices/security/#configure-third-party-service-account-tokens for details.
! values.global.jwtPolicy is deprecated; use Values.global.jwtPolicy=third-party-jwt. See http://istio.io/latest/docs/ops/best-practices/security/#configure-third-party-service-account-tokens for more information instead
✔ Istio core installed                                                                                         
✔ Istiod installed                                                                                             
✔ Egress gateways installed                                                                                    
✔ Ingress gateways installed                                                                                   
✔ Installation complete                                                                                        
Thank you for installing Istio 1.10.  Please take a few minutes to tell us about your install/upgrade experience!  https://forms.gle/KjkrDnMPByq7akrYA
```

- 安装校验
```text
生成清单通过“kubectl apply -f”安装
• $ istioctl manifest generate > $HOME/generated-manifest.yaml
• $ kubectl apply -f $HOME/generated-manifest.yaml

验证安装
• $ istioctl verify-install -f $HOME/generated-manifest.yaml
• istioctl dashboard 方式
```

- 启动istio dashboard

> istioctl dashboard kiali

删除istio信息
```text
Uninstall
To delete the Bookinfo sample application and its configuration, see Bookinfo cleanup.

The Istio uninstall deletes the RBAC permissions and all resources hierarchically under the istio-system namespace. It is safe to ignore errors for non-existent resources because they may have been deleted hierarchically.

$ kubectl delete -f samples/addons
$ istioctl manifest generate --set profile=demo | kubectl delete --ignore-not-found=true -f -

The istio-system namespace is not removed by default. If no longer needed, use the following command to remove it:

$ kubectl delete namespace istio-system

The label to instruct Istio to automatically inject Envoy sidecar proxies is not removed by default. If no longer needed, use the following command to remove it:

$ kubectl label namespace default istio-injection-
```












### order

- 获取产生的命名空间

> kubectl get ns 

- 获取istio命名空间下的pod信息

> kubectl get po -n istio-system

- 获取istio crd信息

> kubectl get crd | grep istio

- 获取istio 的api资源信息

> kubectl api-resources | grep istio











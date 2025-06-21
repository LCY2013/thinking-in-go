## 添加 Jenkins Helm 仓库
```kubernetes helm
helm repo add jenkinsci https://charts.jenkins.io
helm repo update
```

## 创建命名空间（可选）
```shell
kubectl create namespace jenkins
```

## 自定义安装配置
###  生成默认 values 文件（可选）
```shell
helm show values jenkinsci/jenkins > jenkins-values.yaml
```
### 关键配置项（编辑 jenkins-values.yaml）
```shell
controller:
  image: "jenkins/jenkins"  # 官方镜像
  tag: "2.426.3-jdk17"      # 指定版本（推荐 LTS）
  adminUser: "admin"        # 管理员用户名（默认随机生成）
  adminPassword: "admin123" # 管理员密码（生产环境建议使用 Secret）
  serviceType: LoadBalancer # 或 NodePort/ClusterIP
  servicePort: 8080
  installPlugins:           # 预安装插件
    - kubernetes
    - workflow-aggregator
    - git
    - blueocean
  resources:               # 资源限制
    requests:
      cpu: "1000m"
      memory: "2Gi"
    limits:
      cpu: "2000m"
      memory: "4Gi"

agent:
  enabled: true            # 启用 Jenkins Agent
  image: "jenkins/inbound-agent"
  tag: "4.11-1"

persistence:
  enabled: true            # 启用持久化存储
  storageClass: "nfs"      # 替换为你的 StorageClass
  size: "8Gi"
```

## 安装 Jenkins
```shell
helm install jenkins jenkinsci/jenkins \
  --namespace jenkins \
  --values jenkins-values.yaml \
  --wait
```

## 安装
```shell
root@k8s200:~/cicd# helm install jenkins -n jenkins -f helm-jenkins/jenkins-values.yaml $chart
NAME: jenkins
LAST DEPLOYED: Sat Jun 21 13:06:06 2025
NAMESPACE: jenkins
STATUS: deployed
REVISION: 1
NOTES:
1. Get your 'admin' user password by running:
  kubectl exec --namespace jenkins -it svc/jenkins -c jenkins -- /bin/cat /run/secrets/additional/chart-admin-password && echo
2. Get the Jenkins URL to visit by running these commands in the same shell:
  echo http://127.0.0.1:8080
  kubectl --namespace jenkins port-forward svc/jenkins 8080:8080

3. Login with the password from step 1 and the username: admin
4. Configure security realm and authorization strategy
5. Use Jenkins Configuration as Code by specifying configScripts in your values.yaml file, see documentation: http://127.0.0.1:8080/configuration-as-code and examples: https://github.com/jenkinsci/configuration-as-code-plugin/tree/master/demos

For more information on running Jenkins on Kubernetes, visit:
https://cloud.google.com/solutions/jenkins-on-container-engine

For more information about Jenkins Configuration as Code, visit:
https://jenkins.io/projects/jcasc/


NOTE: Consider using a custom image with pre-installed plugins
```

## 获取访问信息

### (1) 获取管理员密码
```shell
kubectl exec --namespace jenkins -it svc/jenkins -c jenkins -- /bin/cat /run/secrets/additional/chart-admin-password
```

### 或（如果使用自定义密码）：
```shell
echo "admin123"  # 与 values.yaml 中的 adminPassword 一致
```

### (2) 获取访问 URL
#### LoadBalancer 类型：
```shell
kubectl get svc jenkins -n jenkins -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
#访问 http://<EXTERNAL-IP>:8080
```

#### NodePort 类型：
```shell
kubectl get svc jenkins -n jenkins
#访问 http://<NODE-IP>:<NODE-PORT>
```

## 验证安装
```shell
kubectl get pods -n jenkins -w  # 等待 Pod 状态变为 Running
kubectl logs jenkins-0 -n jenkins -c jenkins  # 查看日志
```

## 升级或卸载
###  升级配置
```shell
helm upgrade jenkins jenkinsci/jenkins \
  --namespace jenkins \
  --values jenkins-values.yaml
```

###  卸载
```shell
helm uninstall jenkins -n jenkins
kubectl delete pvc -n jenkins --all  # 可选：删除持久化数据
```



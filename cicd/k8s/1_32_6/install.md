## [安装](https://sealos.run/docs/k8s/quick-start/deploy-kubernetes)
```shell
sudo tee -a /etc/hosts <<EOF
192.168.0.200 k8s200
192.168.0.201 k8s201
192.168.0.202 k8s202
EOF

sealos run docker.io/labring/kubernetes:v1.32.6 docker.io/labring/helm:v3.15.4 docker.io/labring/cilium:1.15.13 \
     --masters 192.168.0.200 \
     --nodes 192.168.0.201,192.168.0.202 -p 123456
```

## [安装k8s]
```shell
https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/

https://v1-32.docs.kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/

```

```shell
# 关闭 Swap
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# 设置主机名解析（所有节点）
sudo tee -a /etc/hosts <<EOF
192.168.0.200 k8s200
192.168.0.201 k8s201
192.168.0.202 k8s202
EOF

# 加载内核模块
sudo modprobe overlay
sudo modprobe br_netfilter

# 设置内核参数
sudo tee /etc/sysctl.d/kubernetes.conf <<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
sudo sysctl --system

# 安装容器运行时（containerd）
sudo apt-get update
sudo apt-get install -y containerd
sudo mkdir -p /etc/containerd
sudo containerd config default | sudo tee /etc/containerd/config.toml
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
sudo systemctl restart containerd
sudo systemctl enable containerd

# 安装 kubeadm/kubelet/kubectl
# 删除旧仓库配置（如果存在）
sudo rm -f /etc/apt/sources.list.d/kubernetes.list

# 添加适用于Ubuntu 24.04的Kubernetes仓库
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.32/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# 更新软件包索引
sudo apt update


# 安装Kubernetes 1.32.6组件
sudo apt install -y kubelet=1.32.6-1.1 kubeadm=1.32.6-1.1 kubectl=1.32.6-1.1

# 锁定版本防止自动升级
sudo apt-mark hold kubelet kubeadm kubectl

# 使用阿里云镜像源替代
sudo sed -i 's|https://pkgs.k8s.io|https://mirrors.aliyun.com/kubernetes|g' /etc/apt/sources.list.d/kubernetes.list
sudo apt update

# 初始化集群
sudo kubeadm init \
  --apiserver-advertise-address=192.168.0.200 \
  --kubernetes-version=v1.32.6 \
  --pod-network-cidr=10.42.0.0/16 \
  --image-repository=registry.k8s.io \
  --cri-socket=unix:///run/containerd/containerd.sock
  
# 配置 kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# 安装 Calico 网络插件
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml

# 等待 Calico 就绪
kubectl wait --for=condition=ready pod -l k8s-app=calico-node -n kube-system --timeout=180s  
  
#Your Kubernetes control-plane has initialized successfully!
#
#To start using your cluster, you need to run the following as a regular user:
#
#  mkdir -p $HOME/.kube
#  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
#  sudo chown $(id -u):$(id -g) $HOME/.kube/config
#
#Alternatively, if you are the root user, you can run:
#
#  export KUBECONFIG=/etc/kubernetes/admin.conf
#
#You should now deploy a pod network to the cluster.
#Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
#  https://kubernetes.io/docs/concepts/cluster-administration/addons/
#
#Then you can join any number of worker nodes by running the following on each as root:
#
#
# 在 worker1 (201) 和 worker2 (202) 上执行
# 使用 master 初始化时生成的 join 命令（替换为你的实际 token）
#sudo kubeadm join master:6443 \
#  --token <your-token> \
#  --discovery-token-ca-cert-hash sha256:<your-hash>
# Worker 节点加入集群 (201, 202)
kubeadm join 192.168.0.200:6443 --token 7ggoib.q21uo13cdguju7u1 \
	--discovery-token-ca-cert-hash sha256:8a1e38e85fb8b980d6bc916e465cdd4318960ab186ea40560218f836c60f3282 
#[preflight] Running pre-flight checks
#[preflight] Reading configuration from the "kubeadm-config" ConfigMap in namespace "kube-system"...
#[preflight] Use 'kubeadm init phase upload-config --config your-config.yaml' to re-upload it.
#[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
#[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
#[kubelet-start] Starting the kubelet
#[kubelet-check] Waiting for a healthy kubelet at http://127.0.0.1:10248/healthz. This can take up to 4m0s
#[kubelet-check] The kubelet is healthy after 500.387935ms
#[kubelet-start] Waiting for the kubelet to perform the TLS Bootstrap
#
#This node has joined the cluster:
#* Certificate signing request was sent to apiserver and a response was received.
#* The Kubelet was informed of the new secure connection details.
#
#Run 'kubectl get nodes' on the control-plane to see this node join the cluster.

# 获取 Token 和 Hash（默认有效期24小时）
kubeadm token create --print-join-command

# 输出示例
kubeadm join 192.168.0.200:6443 \
  --token abcdef.0123456789abcdef \
  --discovery-token-ca-cert-hash sha256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  

# 在每个 Worker 节点上执行从 Master 获取的 kubeadm join 命令：
# 示例（替换为你的实际 Token 和 Hash）
sudo kubeadm join 192.168.0.200:6443 \
  --token abcdef.0123456789abcdef \
  --discovery-token-ca-cert-hash sha256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  
# This node has joined the cluster:
#* Certificate signing request was sent to apiserver and a response was received.
#* The Kubelet was informed of the new secure connection details.

#回到 Master 节点，检查节点是否已加入：
#NAME           STATUS     ROLES           AGE   VERSION
#192.168.0.200  Ready      control-plane   10m   v1.32.6
#192.168.0.201  Ready      <none>          2m    v1.32.6
#192.168.0.202  Ready      <none>          1m    v1.32.6
  
#  ⚠️ 如果 Worker 节点状态为 NotReady：

#检查网络插件（如 Calico/Cilium）是否已安装：kubectl get pods -n kube-system

#确保 Worker 节点的防火墙开放了必要端口（如 6443、10250）。

#安装 Cilium
# 在 Master 节点执行
helm repo add cilium https://helm.cilium.io/
helm repo update

# 安装 Cilium（使用 VXLAN 隧道模式）
helm install cilium cilium/cilium \
  --version 1.14.3 \
  --namespace kube-system \
  --set tunnel=vxlan \
  --set ipv4NativeRoutingCIDR=10.42.0.0/16 \
  --set kubeProxyReplacement=strict

# 等待 Cilium 就绪
kubectl -n kube-system wait --for=condition=ready pod -l k8s-app=cilium --timeout=180s

# 检查节点状态
kubectl get nodes -o wide

# 检查所有 Pod 状态
kubectl get pods -A

# 检查网络插件
kubectl get pods -n kube-system -l k8s-app=calico-node
kubectl get pods -n kube-system -l k8s-app=cilium

# 测试网络连通性
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port=80
kubectl run test --image=busybox --rm -it -- sh
> wget -qO- nginx

#为 Worker 节点添加标签
#如果需要标记节点角色（如 worker）：
kubectl label node k8s201 node-role.kubernetes.io/worker=
kubectl label node k8s202 node-role.kubernetes.io/worker=
#验证标签
kubectl get nodes --show-labels

## Token 过期
# 在 Master 节点重新生成 Token 和 Hash
kubeadm token create --print-join-command

## 在 Worker 节点检查 kubelet 状态：
sudo systemctl restart kubelet
sudo journalctl -u kubelet -f  # 查看日志

```

### containerd重启后问题
1. 清理残留的容器和沙箱
```shell
# 停止 kubelet 服务
sudo systemctl stop kubelet

# 清理 containerd 中的残留容器
sudo crictl rm -fa
sudo crictl rmp -fa

# 清理网络命名空间
sudo ip netns list | xargs -I {} sudo ip netns delete {}

# 清理 containerd 状态
sudo systemctl stop containerd
sudo rm -rf /var/lib/containerd/io.containerd.runtime.v2.task/k8s.io/*
sudo systemctl start containerd

# 启动 kubelet 服务
sudo systemctl start kubelet
```



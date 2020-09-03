#### 查看k8s集群状态
kubectl get cs

#### 查看节点状态
kubectl get nodes

#### kubectl describe命令来查看这个节点(Node)对象的详细信息、状态和事件（Event)
kubectl describe node k8sm-120

#### 可以通过 kubectl 检查这个节点上各个系统 Pod 的状态，其中，kube-system 是 Kubernetes 项目预留的系统 Pod 的工作空间（Namepsace，注意它并不是 Linux Namespace，它只是 Kubernetes 划分不同工作空间的单位）
kubectl get pod -n kube-system -o wide

#### 可以查看所有节点pod状态
kubectl get pod --all-namespaces -o wide

#### 





























### DaemonSet 的主要作用，是在 Kubernetes 集群里，运行一个 DaemonPod
这个 Pod 有如下三个特征：
```text
1. 这个 Pod 运行在 Kubernetes 集群里的每一个节点（Node）上；
2. 每个节点上只有一个这样的 Pod 实例；
3. 当有新的节点加入 Kubernetes 集群后，该 Pod 会自动地在新节点上被创建出来；而当旧节点被删除后，它上面的 Pod 也相应地会被回收掉。
```

这个机制看起来很简单，但 Daemon Pod 的意义确实是非常重要的：
```text
1. 各种网络插件的 Agent 组件，都必须运行在每一个节点上，用来处理这个节点上的容器网络；
2. 各种存储插件的 Agent 组件，也必须运行在每一个节点上，用来在这个节点上挂载远程存储目录，操作容器的 Volume 目录；
3. 各种监控组件和日志组件，也必须运行在每一个节点上，负责这个节点上的监控信息和日志搜集。
```

更重要的是，跟其他编排对象不一样，DaemonSet 开始运行的时机，很多时候比整个Kubernetes 集群出现的时机都要早。

daemonSet示例：
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd-elasticsearch
  namespace: kube-system
  labels:
    k8s-app: fluentd-logging
spec:
  selector:
    matchLabels:
      name: fluentd-elasticsearch
  template:
    metadata:
      labels:
        name: fluentd-elasticsearch
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: fluentd-elasticsearch
        image: k8s.gcr.io/fluentd-elasticsearch:1.20
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      terminationGracePeriodSeconds: 30
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
```
这个 DaemonSet，管理的是一个 fluentd-elasticsearch 镜像的 Pod，这个镜像的功能非常实用：通过 fluentd 将 Docker 容器里的日志转发到 ElasticSearch 中。

DaemonSet 跟 Deployment 其实非常相似，只不过是没有 replicas 字段；它也使用 selector 选择管理所有携带了 name=fluentd-elasticsearch 标签的 Pod。

fluentd 启动之后，它会从这两个目录里搜集日志信息，并转发给 ElasticSearch 保存。这样，通过 ElasticSearch 就可以很方便地检索这些日志了。

需要注意的是，Docker 容器里应用的日志，默认会保存在宿主机的/var/lib/docker/containers/{{. 容器 ID}}/{{. 容器 ID}}-json.log 文件里，所以这个目录正是 fluentd 的搜集目标。

### DaemonSet 又是如何保证每个 Node 上有且只有一个被管理的 Pod 呢？

这是一个典型的“控制器模型”能够处理的问题。

DaemonSet Controller，从 Etcd 里获取所有的 Node 列表，然后遍历所有的Node。这时，就可以很容易地去检查，当前这个 Node 上是不是有一个携带了name=fluentd-elasticsearch 标签的 Pod 在运行。

检查的结果，可能有这么三种情况：
```text
1. 没有这种 Pod，那么就意味着要在这个 Node 上创建这样一个 Pod；
2. 有这种 Pod，但是数量大于 1，那就说明要把多余的 Pod 从这个 Node 上删除掉；
3. 正好只有一个这种 Pod，那说明这个节点是正常的。
```

其中，删除节点（Node）上多余的 Pod 非常简单，直接调用 Kubernetes API 就可以了，选择节点使用 nodeSelector，选择 Node 的名字即可。
```yaml
nodeSelector:    
  name: <Node 名字 >
```

在 Kubernetes 项目里，nodeSelector 其实已经是一个将要被废弃的字段了。因为，现在有了一个新的、功能更完善的字段可以代替它，即：nodeAffinity。举个例子：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: with-node-affinity
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: metadata.name
            operator: In
            values:
            - node-k8s
```
在这个 Pod 里，声明了一个 spec.affinity 字段，然后定义了一个 nodeAffinity。其中，spec.affinity 字段，是 Pod 里跟调度相关的一个字段。

这里定义的 nodeAffinity 的含义是：
```text
1. requiredDuringSchedulingIgnoredDuringExecution：它的意思是说，这个nodeAffinity 必须在每次调度的时候予以考虑。同时，这也意味着你可以设置在某些情况下不考虑这个 nodeAffinity；

2. 这个 Pod，将来只允许运行在“metadata.name”是“node-k8s”的节点上。
```

在这里，你应该注意到 nodeAffinity 的定义，可以支持更加丰富的语法，比如 operator:In（即：部分匹配；如果你定义 operator: Equal，就是完全匹配），这也正是nodeAffinity 会取代 nodeSelector 的原因之一。

DaemonSet Controller 会在创建 Pod 的时候，自动在这个 Pod 的 API对象里，加上这样一个 nodeAffinity 定义。其中，需要绑定的节点名字，正是当前正在遍历的这个 Node。

DaemonSet 并不需要修改用户提交的 YAML 文件里的 Pod 模板，而是在向Kubernetes 发起请求之前，直接修改根据模板生成的 Pod 对象。

DaemonSet 还会给这个 Pod 自动加上另外一个与调度相关的字段，叫作tolerations。这个字段意味着这个 Pod，会“容忍”（Toleration）某些 Node 的“污点”（Taint）。

DaemonSet 自动加上的 tolerations 字段，格式如下所示：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: with-toleration
spec:
  tolerations:
  - key: node.kubernetes.io/unschedulable
    operator: Exists
    effect: NoSchedule
```
这个 Toleration 的含义是：“容忍”所有被标记为 unschedulable“污点”的Node；“容忍”的效果是允许调度。

正常情况下，被标记了 unschedulable“污点”的 Node，是不会有任何 Pod 被调度上去的（effect: NoSchedule）。可是，DaemonSet 自动地给被管理的 Pod 加上了这个特殊的 Toleration，就使得这些 Pod 可以忽略这个限制，继而保证每个节点上都会被调度一个 Pod。当然，如果这个节点有故障的话，这个 Pod 可能会启动失败，而 DaemonSet则会始终尝试下去，直到 Pod 启动成功。

DaemonSet 的“过人之处”，其实就是依靠Toleration 实现的。

假如当前 DaemonSet 管理的，是一个网络插件的 Agent Pod，那么你就必须在这个DaemonSet 的 YAML 文件里，给它的 Pod 模板加上一个能够“容忍”node.kubernetes.io/network-unavailable“污点”的 Toleration，如下：
```yaml
#...
template:    
  metadata:      
    labels:        
      name: network-plugin-agent    
spec:      
  tolerations:      
  - key: node.kubernetes.io/network-unavailable        
    operator: Exists        
    effect: NoSchedule
```
在 Kubernetes 项目中，当一个节点的网络插件尚未安装时，这个节点就会被自动加上名为node.kubernetes.io/network-unavailable的“污点”。

而通过这样一个 Toleration，调度器在调度这个 Pod 的时候，就会忽略当前节点上的“污点”，从而成功地将网络插件的 Agent 组件调度到这台机器上启动起来。

这种机制，正是在部署 Kubernetes 集群的时候，能够先部署 Kubernetes 本身、再部署网络插件的根本原因：因为当时所创建的 Weave 的 YAML，实际上就是一个DaemonSet。

DaemonSet 其实是一个非常简单的控制器。在它的控制循环中，只需要遍历所有节点，然后根据节点上是否有被管理 Pod 的情况，来决定是否要创建或者删除一个 Pod。

在创建每个 Pod 的时候，DaemonSet 会自动给这个 Pod 加上一个nodeAffinity，从而保证这个 Pod 只会在指定节点上启动。同时，它还会自动给这个 Pod加上一个 Toleration，从而忽略节点的 unschedulable“污点”。

也可以在 Pod 模板里加上更多种类的 Toleration，从而利用 DaemonSet 实现自己的目的。比如，在这个 fluentd-elasticsearch DaemonSet 里，就给它加上了这样的 Toleration：
```yaml
tolerations:
- key: node-role.kubernetes.io/master  
  effect: NoSchedule
```
这是因为在默认情况下，Kubernetes 集群不允许用户在 Master 节点部署 Pod。因为，Master 节点默认携带了一个叫作node-role.kubernetes.io/master的“污点”。所以，为了能在 Master 节点上部署 DaemonSet 的 Pod，就必须让这个 Pod“容忍”这个“污点”。

### 示例演示DaemonSet
创建这个 DaemonSet 对象：
```text
$ kubectl create -f fluentd-elasticsearch-daemonset.yaml

需要注意的是，在 DaemonSet 上，一般都应该加上 resources 字段，来限制它的CPU 和内存使用，防止它占用过多的宿主机资源。

$ kubectl get pod -n kube-system -l name=fluentd-elasticsearch
NAME                          READY   STATUS    RESTARTS   AGE
fluentd-elasticsearch-fm84f   1/1     Running   0          60s
fluentd-elasticsearch-xnwb4   1/1     Running   0          60s
fluentd-elasticsearch-ztqgv   1/1     Running   0          60s

通过 kubectl get 查看一下 Kubernetes 集群里的 DaemonSet 对象：
$ kubectl get ds -n kube-system fluentd-elasticsearch
NAME                    DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
fluentd-elasticsearch   3         3         3       3            3           <none>          3m4s

Kubernetes 里比较长的 API 对象都有短名字，比如 DaemonSet 对应的是 ds，Deployment 对应的是 deploy。

发现 DaemonSet 和 Deployment 一样，也有 DESIRED、CURRENT 等多个状态字段。这也就意味着，DaemonSet 可以像 Deployment 那样，进行版本管理。这个版本，可以使用 kubectl rollout history 看到：
$ kubectl rollout history daemonset fluentd-elasticsearch -n kube-system
daemonset.apps/fluentd-elasticsearch 
REVISION  CHANGE-CAUSE
1         <none>

升级这个 DaemonSet 的容器镜像版本到 v2.2.0：
$ kubectl set image ds/fluentd-elasticsearch fluentd-elasticsearch=k8s.gcr.io/fluentd-elasticsearch:v2.2.0 --record -n=kube-system
如果访问不了外网直接使用docker pull gydtc/fluentd-elasticsearch:v2.2.0
$ kubectl set image ds/fluentd-elasticsearch fluentd-elasticsearch=gydtc/fluentd-elasticsearch:v2.2.0 --record -n=kube-system
这个 kubectl set image 命令里，第一个 fluentd-elasticsearch 是 DaemonSet 的名字，第二个 fluentd-elasticsearch 是容器的名字。

使用 kubectl rollout status 命令看到这个“滚动更新”的过程
$ kubectl rollout status ds/fluentd-elasticsearch -n kube-system

由于这一次在升级命令后面加上了–record 参数，所以这次升级使用到的指令就会自动出现在 DaemonSet 的 rollout history 里面，如下所示：
$ kubectl rollout history daemonset fluentd-elasticsearch -n kube-system
有了版本号，也就可以像 Deployment 一样，将 DaemonSet 回滚到某个指定的历史版本了。

Deployment 管理这些版本，靠的是“一个版本对应一个 ReplicaSet 对象”。可是，DaemonSet 控制器操作的直接就是 Pod，不可能有 ReplicaSet 这样的对象参与其中。那么，它的这些版本又是如何维护的呢？

Kubernetes v1.7 之后添加了一个 API 对象，名叫ControllerRevision，专门用来记录某种 Controller 对象的版本。比如，可以通过如下命令查看 fluentd-elasticsearch 对应的ControllerRevision：
$ kubectl get controllerrevision -n kube-system -l name=fluentd-elasticsearc

使用 kubectl describe 查看这个 ControllerRevision 对象：
$ kubectl describe controllerrevision fluentd-elasticsearch-fm84f -n kube-system
这个 ControllerRevision 对象，实际上是在 Data 字段保存了该版本对应的完整的 DaemonSet 的 API 对象。并且，在 Annotation 字段保存了创建这个对象所使用的kubectl 命令。

可以尝试将这个 DaemonSet 回滚到 Revision=1 时的状态：
$ kubectl rollout undo daemonset fluentd-elasticsearch --to-revision=1 -n kube-system

这个 kubectl rollout undo 操作，实际上相当于读取到了 Revision=1 的ControllerRevision 对象保存的 Data 字段。而这个 Data 字段里保存的信息，就是Revision=1 时这个 DaemonSet 的完整 API 对象。

所以，现在 DaemonSet Controller 就可以使用这个历史 API 对象，对现有的DaemonSet 做一次 PATCH 操作（等价于执行一次 kubectl apply -f “旧的 DaemonSet对象”），从而把这个 DaemonSet“更新”到一个旧版本。

这也是为什么，在执行完这次回滚完成后，你会发现，DaemonSet 的 Revision 并不会从Revision=2 退回到 1，而是会增加成 Revision=3。这是因为，一个新的ControllerRevision 被创建了出来。
```

相比于 Deployment，DaemonSet 只管理 Pod 对象，然后通过 nodeAffinity 和Toleration 这两个调度器的小功能，保证了每个节点上有且只有一个 Pod。

与此同时，DaemonSet 使用 ControllerRevision，来保存和管理自己对应的“版本”。这种“面向 API 对象”的设计思路，大大简化了控制器本身的逻辑，也正是 Kubernetes 项目“声明式 API”的优势所在。

StatefulSet 也是直接控制 Pod 对象的，那么它是不是也在使用 ControllerRevision 进行版本管理呢？

没错。在 Kubernetes 项目里，ControllerRevision 其实是一个通用的版本管理对象。这样，Kubernetes 项目就巧妙地避免了每种控制器都要维护一套冗余的代码和逻辑的问题。


































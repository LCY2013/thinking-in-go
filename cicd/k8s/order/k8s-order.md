#### 查看k8s集群 master (ComponentStatus)组件状态
kubectl get cs

#### 查看节点状态
kubectl get nodes

#### kubectl describe命令来查看这个节点(Node)对象的详细信息、状态和事件（Event)
kubectl describe node k8sm-120

#### 可以通过 kubectl 检查这个节点上各个系统 Pod 的状态，其中，kube-system 是 Kubernetes 项目预留的系统 Pod 的工作空间（Namepsace，注意它并不是 Linux Namespace，它只是 Kubernetes 划分不同工作空间的单位）
kubectl get pod -n kube-system -o wide

#### 可以查看所有节点pod状态
kubectl get pod --all-namespaces -o wide

#### 查看某个node详细情况、状态、事件(Event)
kubectl describe node [master]

#### 查看所有POD情况
kubectl get pods

#### 查看命名空间是kube-system下的所有pod(kube-system是kubeadm安装后的默认k8s组件的命名空间)
kubectl get pod -n kube-system -o wide

#### 创建yaml文件部署nginx nginx-deployment.yaml
创建一个nginx-deployment.yaml API对象
$ kubectl create -f nginx-deployment.yaml
deployment.apps/nginx-deployment created

其中nginx-deployment.yaml对应到kubernetes就是一个API Object(API 对象)，里面的每一个字段的值都会被kubernetes解析出来变成特定的容器或者其他类型的API资源。
    kind: 指定API字段的类型(Type),本例中指定的是Deployment类型
    Deployment: 是一个定义多副本应用(多个副本pod)的对象,它还负责在pod定义发生变化时,对每个副本进行滚动升级(Rolling Update),本例中定义了两个副本集
    replicas: 副本集的定义由spec.template内容定义
    pod: pod就是kubernetes世界中的应用,一个应用可以由多个容器组成
kubernetes中使用一种API对象(Deployment)去控制另一种API对象(Pod)的方法叫做<控制器模式(controller pattern)>
每一个API对象都有一个metadata的字段,这个字段的对象叫做API对象的标识,也就是元数据,它是我们在k8s中找API对象的一个重要标识,其中最重要的是labels字段。
    labels字段就是由一组组的key-value 标签组成,他是我们用来筛选对象的主要依据,本例中这个Deployment对象被创建后就会携带一个app:nginx标签,通过这个标签来保证只有两个的pod。
  过滤的字段设置就在spec.selector.matchLabels(Label Selector)。
    另一个metadata中的字段annotations是用来携带key-value格式的内部信息(就是k8s自身感兴趣的东西,而不是用户感兴趣),大多数的annotations都是在k8s运行过程中加在API对象身上的。
API对象主要描述信息:
    1、metadata: 存放API对象的元信息,所有API对象基本统一。
    2、spc: API对象的独有定义,用来描述它所要表达的功能。
    
获取匹配标签app:nginx的所有pod    
$ kubectl get pods -l app=nginx     
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-7cf55fb7bb-ckbgq   1/1     Running   0          31s
nginx-deployment-7cf55fb7bb-lw75n   1/1     Running   0          31s
$ kubectl exec -it nginx-deployment-7cf55fb7bb-ckbgq -n default -- /bin/bash
root@nginx-deployment-7cf55fb7bb-ckbgq:/# curl http://127.0.0.1
<!DOCTYPE html>
<html lang="zh">
<head>
<title>Welcome to nginx!</title>
</head>
</html>

查看其中一个pod的详细信息
$ kubectl describe pod nginx-deployment-7cf55fb7bb-lw75n

Node:         k8s-122/192.168.99.122
Start Time:   Tue, 08 Sep 2020 14:38:10 +0800
Labels:       app=nginx
              pod-template-hash=7cf55fb7bb
Annotations:  <none>
Status:       Running
IP:           10.244.2.33
IPs:
  IP:           10.244.2.33
Controlled By:  ReplicaSet/nginx-deployment-7cf55fb7bb
Containers:
  nginx:
    Container ID:   docker://acfa78c9a72d92cfb741079128f19882de4854ee434f87a95e4d18a23b551367
    Image:          nginx:1.19
    Image ID:       docker-pullable://nginx@sha256:b0ad43f7ee5edbc0effbc14645ae7055e21bc1973aee5150745632a24a752661
    Port:           80/TCP

Events:
  Type    Reason     Age    From              Message
  ----    ------     ----   ----              -------
  Normal  Scheduled  6m39s                    Successfully assigned default/nginx-deployment-7cf55fb7bb-lw75n to k8s-122
  Normal  Pulling    6m38s  kubelet, k8s-122  Pulling image "nginx:1.19"
  Normal  Pulled     6m10s  kubelet, k8s-122  Successfully pulled image "nginx:1.19" in 28.115063072s
  Normal  Created    6m10s  kubelet, k8s-122  Created container nginx
  Normal  Started    6m10s  kubelet, k8s-122  Started container nginx
从上面的详细信息可以看出ip,event等,event是一个重要东西,可以看到pod被调度到了k8s-122节点,拉去镜像的时间等一系列操作记录,如果发生了错误我们就需要查看event相关信息。

接下来尝试升级nginx镜像的版本信息: 将nginx-deployment.yaml image: nginx:1.19 -> image: nginx:1.19.2
$ kubectl replace -f nginx-deployment.yaml
deployment.apps/nginx-deployment replaced
$ kubectl get pods -l app=nginx
NAME                                READY   STATUS        RESTARTS   AGE
nginx-deployment-55df9cfb4b-g8qqk   1/1     Running       0          46s
nginx-deployment-55df9cfb4b-rxcrt   1/1     Running       0          27s
nginx-deployment-7cf55fb7bb-lw75n   0/1     Terminating   0          18m
再次查看nginx-deployment.yaml部署的pod信息
$ kubectl describe pod nginx-deployment-55df9cfb4b-g8qqk 
..
Containers:
  nginx:
    Container ID:   docker://16095b3eee7e740491cbf9954a16cb74259c06b95ee77b558ffdae1683279e18
    Image:          nginx:1.19.2
    Image ID:       docker-pullable://nginx@sha256:b0ad43f7ee5edbc0effbc14645ae7055e21bc1973aee5150745632a24a752661
..
Events:
  Type    Reason     Age   From              Message
  ----    ------     ----  ----              -------
  Normal  Scheduled  101s                    Successfully assigned default/nginx-deployment-55df9cfb4b-g8qqk to k8s-121
  Normal  Pulling    100s  kubelet, k8s-121  Pulling image "nginx:1.19.2"
  Normal  Pulled     83s   kubelet, k8s-121  Successfully pulled image "nginx:1.19.2" in 16.947648883s
  Normal  Created    83s   kubelet, k8s-121  Created container nginx
  Normal  Started    83s   kubelet, k8s-121  Started container nginx
event 可以看到镜像已经替换了

建议使用kubectl apply 统一进行更新和创建的操作,接下来修改nginx-deployment.yaml image: nginx:1.19.2 -> image: nginx:1.19
$ kubectl apply -f nginx-deployment.yaml
$ kubectl get pods -l app=nginx
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-7cf55fb7bb-2qh69   1/1     Running   0          48s
nginx-deployment-7cf55fb7bb-cc5dl   1/1     Running   0          47s
  
接下来就是在pod中添加一个volume,volume作为pod的一部分,可以在nginx-deployment.yaml的spec.template中添加
新建一个nginx-deployment-volume.yaml,创建volume卷:
     Deployment 的 Pod 模板部分添加了一个 volumes 字段，定义了这个 Pod 声明的所有 Volume。它的名字叫作 nginx-vol，类型是 emptyDir。
     emptyDir其实就等同于我们之前讲过的 Docker 的隐式 Volume 参数，即：不显式声明宿主机目录的 Volume。所以，Kubernetes 也会在宿主机上创建一个临时目录，这个目录将来就会被绑定挂载到容器所声明的 Volume 目录上。
     Kubernetes 也提供了显式的 Volume 定义，它叫做 hostPath。比如下面的这个 YAML 文件:
        volumes:
          - name: nginx-vol
            hostPath: 
              path: /var/data
$ kubectl apply -f nginx-deployment-volume.yaml
deployment.apps/nginx-deployment configured
$ kubectl get pods -l app=nginx
NAME                               READY   STATUS    RESTARTS   AGE
nginx-deployment-559d7c9c4-82vgw   1/1     Running   0          38s
nginx-deployment-559d7c9c4-kgwxs   1/1     Running   0          37s
$ kubectl describe pod nginx-deployment-559d7c9c4-82vgw
''''''
Containers:
  nginx:
    Mounts:
      /usr/share/nginx/html from nginx-vol (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-v87cv (ro)
''''''

上面可以看到目录已经被挂载,进入容器查看该目录(/usr/share/nginx/html):
$ kubectl exec -it nginx-deployment-559d7c9c4-82vgw -- /bin/bash
root@nginx-deployment-559d7c9c4-82vgw:/# ls /usr/share/nginx/html/

删除部署
$ kubectl delete -f nginx-deployment-volume.yaml

### pod 
Pod，实际上是在扮演传统基础设施里“虚拟机”的角色；而容器，则是这个虚拟机里运行的用户程序。

显示当前系统运行的进程树
$ pstree -g 

Pod 最重要的一个事实是：它只是一个逻辑概念,Kubernetes 真正处理的，还是宿主机操作系统上 Linux 容器的 Namespace 和 Cgroups，而并不存在一个所谓的 Pod 的边界或者隔离环境。
Pod 如何被创建的？ Pod是一组共享了某些资源(network namespace、volume)的容器。 
docker run --net --volumes-from 这样的命令也能实现:
    $ docker run --net=B --volumes-from=B --name=A image-A ...
但是又一个问题就是B容器需要先于A容器运行,这样依赖就是拓扑图关系。
kubernetes如何解决这样的问题?
   在 Kubernetes 项目里，Pod 的实现需要使用一个中间容器，这个容器叫作 Infra 容器。
    在这个 Pod 中，Infra 容器永远都是第一个被创建的容器，而其他用户定义的容器，
    则通过 Join Network Namespace 的方式，与 Infra 容器关联在一起。
    
   在 Kubernetes 项目里，Infra 容器一定要占用极少的资源，所以它使用的是一个非常特殊的镜像，
    叫作：k8s.gcr.io/pause。这个镜像是一个用汇编语言编写的、永远处于“暂停”状态的容器，解压后的大小也只有 100~200 KB 左右。
    
   这也就意味着，对于 Pod 里的容器 A 和容器 B 来说：
        它们可以直接使用 localhost 进行通信；
        它们看到的网络设备跟 Infra 容器看到的完全一样；
        一个 Pod 只有一个 IP 地址，也就是这个 Pod 的 Network Namespace 对应的 IP 地址；
        当然，其他的所有网络资源，都是一个 Pod 一份，并且被该 Pod 中的所有容器共享；
        Pod 的生命周期只跟 Infra 容器一致，而与容器 A 和 B 无关。   
   同一个 Pod 里面的所有用户容器来说，它们的进出流量，也可以认为都是通过 Infra 容器完成的。
   这一点很重要，因为将来如果你要为 Kubernetes 开发一个网络插件时，
   应该重点考虑的是如何配置这个 Pod 的 Network Namespace，而不是每一个用户容器如何使用你的网络配置，这是没有意义的。                           

例子一(nginx-pod.yaml): 
    接下来创建一个nginx-pod.yaml,里面定义两个容器,一个volume采用hostPath:/data 目录, 
而容器一个挂载到/usr/share/nginx/html,一个挂载到/pod-data中并且在这个容器中写入一个index.html,
这样两个容器就共享了宿主机的/data目录。

例子二(war-tomcat-pod.yaml):
    我们现在有一个 Java Web 应用的 WAR 包，它需要被放在 Tomcat 的 webapps 目录下运行起来。
    docker 方案:
        1、打包镜像的时候将war包直接放到tomcat中的webapps里面，一起打成一个镜像
        2、打包tomcat的时候通过挂载一个宿主机的目录到tomcat的webapps目录中
    这两种方案都存在一定的问题,所以使用kubernetes中的pod(容器设计模式),通过war容器镜像/app里面挂载war包,通过tomcat容器镜像挂载webapps目录达到同享目录。
    initContainers 会比containers先启动,所以会将war包复制到对应的/app目录(也就是两个容器共享的目录),然后就会退出。
    我们就用一种"组合"方式，解决了 WAR 包与 Tomcat 容器之间耦合关系的问题。
    实际上，这个所谓的"组合"操作，正是容器设计模式里最常用的一种模式，它的名字叫：sidecar。
    顾名思义，sidecar 指的就是我们可以在一个 Pod 中，启动一个辅助容器，来完成一些独立于主进程（主容器）之外的工作。
    比如，在我们的这个应用 Pod 中，Tomcat 容器是我们要使用的主容器，而 WAR 包容器的存在，只是为了给它提供一个 WAR 包而已。
    所以，我们用 Init Container 的方式优先运行 WAR 包容器，扮演了一个 sidecar 的角色。

例子三:
    比如，我现在有一个应用，需要不断地把日志文件输出到容器的 /var/log 目录中。
    这时，我就可以把一个 Pod 里的 Volume 挂载到应用容器的 /var/log 目录上。
    然后，我在这个 Pod 里同时运行一个 sidecar 容器，它也声明挂载同一个 Volume 到自己的 /var/log 目录上。
    这样，接下来 sidecar 容器就只需要做一件事儿，那就是不断地从自己的 /var/log 目录里读取日志文件，转发到 MongoDB 或者 Elasticsearch 中存储起来。这样，一个最基本的日志收集工作就完成了。
    跟第二个例子一样，这个例子中的 sidecar 的主要工作也是使用共享的 Volume 来完成对文件的操作。
    Pod 的另一个重要特性是，它的所有容器都共享同一个 Network Namespace。这就使得很多与 Pod 网络相关的配置和管理，
        也都可以交给 sidecar 完成，而完全无须干涉用户容器。这里最典型的例子莫过于 Istio 这个微服务治理项目了。





















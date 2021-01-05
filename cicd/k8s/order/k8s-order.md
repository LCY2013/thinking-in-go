### 所有*.yaml文件位置cicd/k8s/case
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
```text
kind: 指定API字段的类型(Type),本例中指定的是Deployment类型
Deployment: 是一个定义多副本应用(多个副本pod)的对象,它还负责在pod定义发生变化时,对每个副本进行滚动升级(Rolling Update),本例中定义了两个副本集
replicas: 副本集的定义由spec.template内容定义
pod: pod就是kubernetes世界中的应用,一个应用可以由多个容器组成
```

```text
kubernetes中使用一种API对象(Deployment)去控制另一种API对象(Pod)的方法叫做<控制器模式(controller pattern)>
每一个API对象都有一个metadata的字段,这个字段的对象叫做API对象的标识,也就是元数据,它是我们在k8s中找API对象的一个重要标识,其中最重要的是labels字段。
    labels字段就是由一组组的key-value 标签组成,他是我们用来筛选对象的主要依据,本例中这个Deployment对象被创建后就会携带一个app:nginx标签,通过这个标签来保证只有两个的pod。
  过滤的字段设置就在spec.selector.matchLabels(Label Selector)。
    另一个metadata中的字段annotations是用来携带key-value格式的内部信息(就是k8s自身感兴趣的东西,而不是用户感兴趣),大多数的annotations都是在k8s运行过程中加在API对象身上的。
```

API对象主要描述信息:
```text
1、metadata: 存放API对象的元信息,所有API对象基本统一。
2、spc: API对象的独有定义,用来描述它所要表达的功能。
3、status: Pending、Running、Succeeded、Failed、Unknown
```
    
获取匹配标签app:nginx的所有pod  
```text
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
```

查看其中一个pod的详细信息
```text
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
```

接下来尝试升级nginx镜像的版本信息: 将nginx-deployment.yaml image: nginx:1.19 -> image: nginx:1.19.2
```text
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

```  

接下来就是在pod中添加一个volume,volume作为pod的一部分,可以在nginx-deployment.yaml的spec.template中添加
```text
新建一个nginx-deployment-volume.yaml,创建volume卷:

 Deployment 的 Pod 模板部分添加了一个 volumes 字段，定义了这个 Pod 声明的所有 Volume，它的名字叫作 nginx-vol，类型是 emptyDir。
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

```

### pod 
Pod，实际上是在扮演传统基础设施里“虚拟机”的角色；而容器，则是这个虚拟机里运行的用户程序。

显示当前系统运行的进程树

$ pstree -g 

Pod 最重要的一个事实是：它只是一个逻辑概念,Kubernetes 真正处理的，还是宿主机操作系统上 Linux 容器的 Namespace 和 Cgroups，而并不存在一个所谓的 Pod 的边界或者隔离环境。

Pod 如何被创建的？ Pod是一组共享了某些资源(network namespace、volume)的容器。 

docker run --net --volumes-from 这样的命令也能实现:
```text
$ docker run --net=B --volumes-from=B --name=A image-A ...
```

但是又一个问题就是B容器需要先于A容器运行,这样依赖就是拓扑图关系。

kubernetes如何解决这样的问题?
```text
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
```

例子一(nginx-pod.yaml): 

接下来创建一个nginx-pod.yaml,里面定义两个容器,一个volume采用hostPath:/data 目录, 
而容器一个挂载到/usr/share/nginx/html,一个挂载到/pod-data中并且在这个容器中写入一个index.html,
这样两个容器就共享了宿主机的/data目录。

例子二(war-tomcat-pod.yaml):
```text
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
```

例子三:
```text
比如，我现在有一个应用，需要不断地把日志文件输出到容器的 /var/log 目录中。
这时，我就可以把一个 Pod 里的 Volume 挂载到应用容器的 /var/log 目录上。
然后，我在这个 Pod 里同时运行一个 sidecar 容器，它也声明挂载同一个 Volume 到自己的 /var/log 目录上。
这样，接下来 sidecar 容器就只需要做一件事儿，那就是不断地从自己的 /var/log 目录里读取日志文件，转发到 MongoDB 或者 Elasticsearch 中存储起来。这样，一个最基本的日志收集工作就完成了。
跟第二个例子一样，这个例子中的 sidecar 的主要工作也是使用共享的 Volume 来完成对文件的操作。
Pod 的另一个重要特性是，它的所有容器都共享同一个 Network Namespace。这就使得很多与 Pod 网络相关的配置和管理，
    也都可以交给 sidecar 完成，而完全无须干涉用户容器。这里最典型的例子莫过于 Istio 这个微服务治理项目了。
```

### pod 重要概念

```text
理解pod: Pod 看成传统环境里的“机器”、把容器看作是运行在这个"机器"里的"用户程序",
            凡是调度、网络、存储，以及安全相关的属性，基本上是 Pod 级别的。
        这些属性的共同特征是，它们描述的是“机器”这个整体，而不是里面运行的“程序”。
        比如，配置这个“机器”的网卡（即：Pod 的网络定义），配置这个“机器”的磁盘（即：Pod 的存储定义），
        配置这个“机器”的防火墙（即：Pod 的安全定义）。更不用说，这台“机器”运行在哪个服务器之上（即：Pod 的调度）。

 pod主要字段含义和用法: 
   1、NodeSelector：是一个供用户将 Pod 与 Node 进行绑定的字段，示例如下:
      apiVersion: v1
      kind: Pod
      ...
      spec:
       nodeSelector:
         disktype: ssd
      含义: 这个pod只能运行在带有"disktype:ssd"标签(Label)的节点上,否则它的调度就会失败。        
      
   2、NodeName：一旦 Pod 的这个字段被赋值，Kubernetes 项目就会被认为这个 Pod 已经经过了调度，
      调度的结果就是赋值的节点名字。所以，这个字段一般由调度器负责设置，但用户也可以设置它来“骗过”调度器，
      当然这个做法一般是在测试或者调试的时候才会用到。   
   
   3、HostAliases：定义了 Pod 的 hosts 文件（比如 /etc/hosts）里的内容，用法如下：

       apiVersion: v1
       kind: Pod
       ...
       spec:
         hostAliases:
         - ip: "10.1.2.3"
           hostnames:
           - "foo.remote"
           - "bar.remote"
       ...
       上面这个pod中的yaml设置了一组IP和hostname,这样这个pod启动后/etc/hosts内容如下:
       # cat /etc/hosts
       ...
       10.1.2.3 foo.remote
       10.1.2.3 bar.remote
       ...
       注意: 设置pod网络一定要通过这种方式而不是在pod内部添加host信息,否则在pod被删除重建后,kubelet会覆盖掉被修改的内容。 
    举例如下:
        凡是跟容器的 Linux Namespace 相关的属性，也一定是 Pod 级别的。这个原因也很容易理解：Pod 的设计，
        就是要让它里面的容器尽可能多地共享 Linux Namespace，仅保留必要的隔离和限制能力。
        这样，Pod 模拟出的效果，就跟虚拟机里程序间的关系非常类似了。
     例子一: 
      创建一个share-nginx-busybox.yaml文件,这个文件里面"shareProcessNamespace: true"就是说容器共享PID namespace。
      里面还定义了两个容器(nginx,busybox),其中busybox开启了tty(-t)和stdin(-i),tty就是linux提供给用户的常驻小程序,stdin就是一个linux中的标准输入。
       启动Pod
      $ kubectl apply -f share-nginx-busybox.yaml
       使用kubectl attach 连接pod容器
      $ kubectl attach -it nginx -c shell   
      / # ps ax
      PID   USER     TIME  COMMAND
          1 root      0:00 /pause
          6 root      0:00 nginx: master process nginx -g daemon off;
         33 101       0:00 nginx: worker process
         34 root      0:00 sh
         39 root      0:00 ps ax
      可以看到里面有nginx容器的进程信息,以及 Infra 容器的 /pause 进程,所有可以看到整个pod中的容器进程信息(共享同一个PID Namespace)。   
     例子二:
       凡是 Pod 中的容器要共享宿主机的 Namespace，也一定是 Pod 级别的定义，创建一个share-nginx-vm.yaml文件,
       里面 hostNetwork: true,hostIPC: true,hostPID: true,就会直接使用宿主机的网络、直接与宿主机的IPC通信、可以查看宿主机的进程信息。
   
   4、containers

        "containers"字段也是Pod中的一个重要概念，除此之外还有"Init Containers"，这两个字段都属于 Pod 对容器的定义，
     内容也完全相同，只是 Init Containers 的生命周期，会先于所有的 Containers，并且严格按照定义的顺序执行。
        containers: 
            Image（镜像）
            Command（启动命令）
            workingDir（容器的工作目录）
            Ports（容器要开发的端口）
            volumeMounts（容器要挂载的 Volume）
            ImagePullPolicy（拉取策略，老版本默认是 Always, 新版本默认是IfNotPresent ）：https://kubernetes.io/docs/concepts/containers/images/
                always： 每次创建 Pod 都重新拉取一次镜像
                Never或者IfNotPresent： Pod 永远不会主动拉取这个镜像，或者只在宿主机上不存在这个镜像时才拉取。
            Lifecycle：Container Lifecycle Hooks 的作用，是在容器状态发生变化时触发一系列“钩子”。
                示例查看：nginx-lifecycle.yaml
                    postStart 指的是，在容器启动后，立刻执行一个指定的操作。需要明确的是，postStart 定义的操作，
                虽然是在 Docker 容器 ENTRYPOINT 执行之后，但它并不严格保证顺序。也就是说，在 postStart 启动时，ENTRYPOINT 有可能还没有结束。
                当然，如果 postStart 执行超时或者错误，Kubernetes 会在该 Pod 的 Events 中报出该容器启动失败的错误信息，导致 Pod 也处于失败的状态。 
                    preStop 发生的时机，则是容器被杀死之前（比如，收到了 SIGKILL 信号）。而需要明确的是，preStop 操作的执行，是同步的。
                所以，它会阻塞当前的容器杀死流程，直到这个 Hook 定义操作完成之后，才允许容器被杀死，这跟 postStart 不一样。
   
                     
```
```text 
pod中另一个重要字段(status):
    pod.status.phase，就是 Pod 的当前状态，它有如下几种可能的情况：
        1、Pending: 这个状态意味着，Pod 的 YAML 文件已经提交给了 Kubernetes，API 对象已经被创建并保存在 Etcd 当中。
                  但是，这个 Pod 里有些容器因为某种原因而不能被顺利创建。比如，调度不成功。
        2、Running: 这个状态下，Pod 已经调度成功，跟一个具体的节点绑定。它包含的容器都已经创建成功，并且至少有一个正在运行中。
        3、Succeeded: 这个状态意味着，Pod 里的所有容器都正常运行完毕，并且已经退出了。这种情况在运行一次性任务时最为常见。
        4、Failed: 这个状态下，Pod 里至少有一个容器以不正常的状态（非 0 的返回码）退出。这个状态的出现，意味着你得想办法 Debug 这个容器的应用，比如查看 Pod 的 Events 和日志。
        5、Unknown: 这是一个异常状态，意味着 Pod 的状态不能持续地被 kubelet 汇报给 kube-apiserver，这很有可能是主从节点（Master 和 Kubelet）间的通信出现了问题。
    Pod 对象的 Status 字段，还可以再细分出一组 Conditions。这些细分状态的值包括：PodScheduled、Ready、Initialized，以及 Unschedulable。
    它们主要用于描述造成当前 Status 的具体原因是什么。 比如，Pod 当前的 Status 是 Pending，对应的 Condition 是 Unschedulable，这就意味着它的调度出现了问题。   
    $GOPATH/src/k8s.io/kubernetes/vendor/k8s.io/api/core/v1/types.go 里，type Pod struct ，尤其是 PodSpec 部分的内容。争取做到下次看到一个 Pod 的 YAML 文件时，不再需要查阅文档，就能做到把常用字段及其作用信手拈来。
```

#### Pod中重要字段 Projected Volume
[secret](https://kubernetes.io/zh/docs/concepts/configuration/secret/)

[configmap](https://kubernetes.io/zh/docs/concepts/configuration/configmap/)

[downward api](https://kubernetes.io/zh/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/)

```text
特殊的 Volume，叫作 Projected Volume，你可以把它翻译为“投射数据卷”。
    注意：Projected Volume 是 Kubernetes v1.11 之后的新特性
 在 Kubernetes 中，有几种特殊的 Volume，它们存在的意义不是为了存放容器里的数据，也不是用来进行容器和宿主机之间的数据交换。
这些特殊 Volume 的作用，是为容器提供预先定义好的数据。所以，从容器的角度来看，这些 Volume 里的信息就是仿佛是被 Kubernetes“投射”（Project）进入容器当中的。 
目前Kubernetes 支持的 Projected Volume 一共有四种：
   1、Secret: 把 Pod 想要访问的加密数据，存放到 Etcd 中，然后就可以通过在 Pod 的容器里挂载 Volume 的方式，访问到这些 Secret 里保存的信息了。
        secret 最典型的使用场景，莫过于存放数据库的 Credential 信息，如下:
          mysql-secret.yaml中 这个 Pod 中，定义了一个简单的容器,它声明挂载的 Volume，并不是常见的 emptyDir 或者 hostPath 类型，
        而是 projected 类型。而这个 Volume 的数据来源（sources），则是名为 user 和 pass 的 Secret 对象，分别对应的是数据库的用户名和密码。
        这里用到的数据库的用户名、密码，正是以 Secret 对象的方式交给 Kubernetes 保存的。完成这个操作的指令，如下所示:
        $ ctest-projected-volumeat ./username.txt
         root
        $ cat ./password.txt
         123456
         创建两个secret对象user、pass，分别用来存储username.txt，password.txt信息
        $ kubectl create secret generic user --from-file=./username.txt
        $ kubectl create secret generic pass --from-file=./password.txt
         如何获取已经创建的secret信息?
        $ kubectl get secrets  
        $ kubectl get secrets  -o yaml  # 以yaml的形式查看
         查看secret的详细信息
        $ kubectl describe secret user
          Name:         user
          Namespace:    default
          Labels:       <none>
          Annotations:  <none>
          Type:  Opaque
          Data
          ====
          username.txt:  5 bytes 
            通过编写yaml文件来生成secret API对象,具体查看secret-create.yaml。
            通过编写 YAML 文件创建出来的 Secret 对象只有一个。但它的 data 字段，却以 Key-Value 的格式保存了两份 Secret 数据。
        其中，“user”就是第一份数据的 Key，“pass”是第二份数据的 Key。
            需要注意的是，Secret 对象要求这些数据必须是经过 Base64 转码的，以免出现明文密码的安全隐患。
        这个转码操作也很简单，比如： 
          $ echo -n root | base64
            cm9vdA==
          $ echo -n 123456 | base64
            MTIzNDU2
         删除secret
        $ kubectl delete secret user
        $ kubectl delete secret pass
        利用secret-create.yaml创建一个secret 对象
        $ kubectl apply -f secret-create.yaml
        接下来利用mysql-secret.yaml创建这个pod
        $ kubectl create -f mysql-secret.yaml
        登陆容器查看信息
        $ kubectl exec -it test-projected-volume -- /bin/sh
          / # ls -l projected-volume/
          total 0
          lrwxrwxrwx    1 root     root            19 Sep  9 06:39 password.txt -> ..data/password.txt
          lrwxrwxrwx    1 root     root            19 Sep  9 06:39 username.txt -> ..data/username.txt 
          / # cat projected-volume/password.txt 
          123456
          / # cat projected-volume/username.txt 
          root
        通过挂载方式进入到容器里的 Secret，一旦其对应的 Etcd 里的数据被更新，这些 Volume 里的文件内容，同样也会被更新。其实，这是 kubelet 组件在定时维护这些 Volume。  
   2、ConfigMap
            与 Secret 类似的是 ConfigMap，它与 Secret 的区别在于，ConfigMap 保存的是不需要加密的、应用所需的配置信息。
        而 ConfigMap 的用法几乎与 Secret 完全相同：你可以使用 kubectl create configmap 从文件或者目录创建 ConfigMap，也可以直接编写 ConfigMap 对象的 YAML 文件。
        例子: 一个 Java 应用所需的配置文件（.properties 文件），就可以通过下面这样的方式保存在 ConfigMap 里：
        # 编辑一个properties文件信息,内容如下
        $ cat > ui.properties
        color.good=purple
        color.bad=yellow
        allow.textmode=true
        how.nice.to.look=fairlyNice
        # 利用properties文件创建一个Configmap
        $ kubectl create configmap ui-config --from-file=./ui.properties
        # 以yaml的形式输出configmap API对象中的内容
        $ kubectl get configmap ui-config -o yaml
   3、Downward API
            作用: 让 Pod 里的容器能够直接获取到这个 Pod API 对象本身的信息。
        定义一个test-downwardapi-volume.yaml，这个 Pod 的 YAML 文件中，定义了一个简单的容器，声明了一个 projected 类型的 Volume。
        只不过这次 Volume 的数据来源，变成了 Downward API。而这个 Downward API Volume，则声明了要暴露 Pod 的 metadata.labels 信息给容器。    
        这样的声明方式，当前 Pod 的 Labels 字段的值，就会被 Kubernetes 自动挂载成为容器里的 /etc/podinfo/labels 文件。
        $ kubectl apply -f test-downwardapi-volume.yaml
        $ kubectl logs test-downwardapi-volume
         cluster="test-cluster1"
         rack="rack-22"
         zone="us-est-coast"
        Downward API 支持的字段已经非常丰富了，比如：
            1. 使用 fieldRef 可以声明使用:
                spec.nodeName - 宿主机名字
                status.hostIP - 宿主机 IP
                metadata.name - Pod 的名字
                metadata.namespace - Pod 的 Namespace
                status.podIP - Pod 的 IP
                spec.serviceAccountName - Pod 的 Service Account 的名字
                metadata.uid - Pod 的 UID
                metadata.labels['<KEY>'] - 指定 <KEY> 的 Label 值
                metadata.annotations['<KEY>'] - 指定 <KEY> 的 Annotation 值
                metadata.labels - Pod 的所有 Label
                metadata.annotations - Pod 的所有 Annotation
            2. 使用 resourceFieldRef 可以声明使用:
                容器的 CPU limit
                容器的 CPU request
                容器的 memory limit
                容器的 memory request 
   4、ServiceAccountToken
            现在有了一个 Pod，我能不能在这个 Pod 里安装一个 Kubernetes 的 Client，这样就可以从容器里直接访问并且操作这个 Kubernetes 的 API 了呢？
        答: 可以，首先要解决 API Server 的授权问题。
             Service Account 对象的作用，就是 Kubernetes 系统内置的一种“服务账户”，它是 Kubernetes 进行权限分配的对象。
             比如，Service Account A，可以只被允许对 Kubernetes API 进行 GET 操作，而 Service Account B，则可以有 Kubernetes API 的所有操作的权限。
             像这样的 Service Account 的授权信息和文件，实际上保存在它所绑定的一个特殊的 Secret 对象里的。
             这个特殊的 Secret 对象，就叫作ServiceAccountToken。任何运行在 Kubernetes 集群上的应用，都必须使用这个 ServiceAccountToken 里保存的授权信息，也就是 Token，才可以合法地访问 API Server。
            Kubernetes 已经为你提供了一个的默认“服务账户”（default Service Account）。并且，任何一个运行在 Kubernetes 里的 Pod，都可以直接使用这个默认的 Service Account，而无需显示地声明挂载它。     
        意一个运行在 Kubernetes 集群里的 Pod，就会发现，每一个 Pod，都已经自动声明一个类型是 Secret、名为 default-token-xxxx 的 Volume，比如:
            $ kubectl describe pod test-downwardapi-volume
            Mounts:
                  /etc/podinfo from podinfo (rw)
                  /var/run/secrets/kubernetes.io/serviceaccount from default-token-v87cv (ro)
            Volumes:
              podinfo:
                Type:         Projected (a volume that contains injected data from multiple sources)
                DownwardAPI:  true
              default-token-v87cv:
                Type:        Secret (a volume populated by a Secret)
                SecretName:  default-token-v87cv
                Optional:    false     
            这个容器内的路径在 Kubernetes 里是固定的，即：/var/run/secrets/kubernetes.io/serviceaccount ，而这个 Secret 类型的 Volume 里面的内容如下所示：
            / # ls /var/run/secrets/kubernetes.io/serviceaccount
            ca.crt     namespace  token        
Secret、ConfigMap，以及 Downward API 这三种 Projected Volume 定义的信息，大多还可以通过环境变量的方式出现在容器里。
但是，通过环境变量获取这些信息的方式，不具备自动更新的能力。所以，一般情况下，都建议你使用 Volume 文件的方式获取这些信息。 
```

#### [Pod 的重要的配置：容器健康检查和恢复机制](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
````text
在 Kubernetes 中，你可以为 Pod 里的容器定义一个健康检查“探针”（Probe）。这样，kubelet 就会根据这个 Probe 的返回值决定这个容器的状态，
而不是直接以容器进行是否运行（来自 Docker 返回的信息）作为依据。这种机制，是生产环境中保证应用健康存活的重要手段。
例子一(probe-create.yaml) :
        在这个 Pod 中，我们定义了一个有趣的容器。它在启动之后做的第一件事，就是在 /tmp 目录下创建了一个 healthy 文件，以此作为自己已经正常运行的标志。而 30 s 过后，它会把这个文件删除掉。
        与此同时，我们定义了一个这样的 livenessProbe（健康检查）。它的类型是 exec，这意味着，它会在容器启动后，在容器里面执行一句我们指定的命令，比如：“cat /tmp/healthy”。
    这时，如果这个文件存在，这条命令的返回值就是 0，Pod 就会认为这个容器不仅已经启动，而且是健康的。这个健康检查，在容器启动 5 s 后开始执行（initialDelaySeconds: 5），每 5 s 执行一次（periodSeconds: 5）。
    $ kubectl apply -f probe-create.yaml
    $ kubectl describe pod test-liveness-exec
    Events:
      Type     Reason     Age                  From              Message
      ----     ------     ----                 ----              -------
      Warning  Unhealthy  32s (x6 over 2m7s)   kubelet, k8s-121  Liveness probe failed: cat: can't open '/tmp/healthy': No such file or directory
    $  kubectl get pods
    NAME                 READY   STATUS             RESTARTS   AGE
    test-liveness-exec   0/1     CrashLoopBackOff   11         32m
   ### 这时我们发现，Pod 并没有进入 Failed 状态，而是保持了 Running 状态。这是为什么呢？
       其实，如果你注意到 RESTARTS 字段从 0 到 1 的变化，就明白原因了：这个异常的容器已经被 Kubernetes 重启了。在这个过程中，Pod 保持 Running 状态不变。
       需要注意的是：Kubernetes 中并没有 Docker 的 Stop 语义。所以虽然是 Restart（重启），但实际却是重新创建了容器。
       这个功能就是 Kubernetes 里的Pod 恢复机制，也叫 restartPolicy。它是 Pod 的 Spec 部分的一个标准字段（pod.spec.restartPolicy），默认值是 Always，即：任何时候这个容器发生了异常，它一定会被重新创建。
       但一定要强调的是，Pod 的恢复过程，永远都是发生在当前节点上，而不会跑到别的节点上去。事实上，一旦一个 Pod 与一个节点（Node）绑定，除非这个绑定发生了变化（pod.spec.node 字段被修改），否则它永远都不会离开这个节点。这也就意味着，如果这个宿主机宕机了，这个 Pod 也不会主动迁移到其他节点上去。
       而如果你想让 Pod 出现在其他的可用节点上，就必须使用 Deployment 这样的“控制器”来管理 Pod，哪怕你只需要一个 Pod 副本。一个单 Pod 的 Deployment 与一个 Pod 最主要的区别就在这里。 
   ### 作为用户，你还可以通过设置 restartPolicy，改变 Pod 的恢复策略。除了 Always，它还有 OnFailure 和 Never 两种情况：
       Always：在任何情况下，只要容器不在运行状态，就自动重启容器；
       OnFailure: 只在容器 异常时才自动重启容器；
       Never: 从来不重启容器。
   ### restartPolicy 和 Pod 里容器的状态，以及 Pod 状态的对应关系:
       1、只要 Pod 的 restartPolicy 指定的策略允许重启异常的容器（比如：Always），那么这个 Pod 就会保持 Running 状态，并进行容器重启。否则，Pod 就会进入 Failed 状态 。
       2、对于包含多个容器的 Pod，只有它里面所有的容器都进入异常状态后，Pod 才会进入 Failed 状态。在此之前，Pod 都是 Running 状态。此时，Pod 的 READY 字段会显示正常容器的个数，比如：  
         $ kubectl get pod test-liveness-exec
           NAME                 READY   STATUS             RESTARTS   AGE
           test-liveness-exec   0/1     CrashLoopBackOff   11         32m
       在容器中执行命令外，livenessProbe 也可以定义为发起 HTTP 或者 TCP 请求的方式，定义格式如下：
       HTTP：
          ...
          livenessProbe:
               httpGet:
                 path: /healthz
                 port: 8080
                 httpHeaders:
                 - name: X-Custom-Header
                   value: Awesome
                 initialDelaySeconds: 3
                 periodSeconds: 3
       TCP：
          ...
           livenessProbe:
             tcpSocket:
               port: 8080
             initialDelaySeconds: 15
             periodSeconds: 20           
       Pod 可以暴露一个健康检查 URL（比如 /healthz），或者直接让健康检查去检测应用的监听端口。这两种配置方法，在 Web 服务类的应用中非常常用。
       在 Kubernetes 的 Pod 中，还有一个叫 readinessProbe 的字段。虽然它的用法与 livenessProbe 类似，但作用却大不一样。
       readinessProbe 检查结果的成功与否，决定的这个 Pod 是不是能被通过 Service 的方式访问到，而并不影响 Pod 的生命周期。
       Pod 的字段这么多，我又不可能全记住，Kubernetes 能不能自动给 Pod 填充某些字段呢？
        答: 这个叫作 PodPreset（Pod 预设置）的功能 已经出现在了 v1.11 版本的 Kubernetes 中。
            举例: preset.yaml、pod.yaml定义
            $ kubectl apply -f preset.yaml
            $ kubectl apply -f pod.yaml
            查看这个pod的相关信息
            $ kubectl get pod website -o wide
                apiVersion: v1
                kind: Pod
                metadata:
                  name: website
                  labels:
                    app: website
                    role: frontend
                  annotations:
                    podpreset.admission.kubernetes.io/podpreset-allow-database: "resource version"
                spec:
                  containers:
                    - name: website
                      image: nginx
                      volumeMounts:
                        - mountPath: /cache
                          name: cache-volume
                      ports:
                        - containerPort: 80
                      env:
                        - name: DB_PORT
                          value: "6379"
                  volumes:
                    - name: cache-volume
                      emptyDir: {}
             可以看到这个pod这里多了新添加的 labels、env、volumes 和 volumeMount 的定义，它们的配置跟 PodPreset 的内容一样。
             此外，这个 Pod 还被自动加上了一个 annotation 表示这个 Pod 对象被 PodPreset 改动过。 
             注意: PodPreset 里定义的内容，只会在 Pod API 对象被创建之前追加在这个对象本身上，而不会影响任何 Pod 的控制器的定义。
             我们现在提交的是一个 nginx-deployment，那么这个 Deployment 对象本身是永远不会被 PodPreset 改变的，被修改的只是这个 Deployment 创建出来的所有 Pod。
             这里有一个问题：如果你定义了同时作用于一个 Pod 对象的多个 PodPreset，会发生什么呢？
             实际上，Kubernetes 项目会帮你合并（Merge）这两个 PodPreset 要做的修改。而如果它们要做的修改有冲突的话，这些冲突字段就不会被修改。
````

### 编排 
```text
控制器模式之Deployment(遵循滚动升级[rolling update]),案例 - nginx-deployment.yaml
 定义了控制app=nginx 标签的 Pod 的个数，永远等于 spec.replicas 指定的个数，即 2 个。
 这就意味着，如果在这个集群中，携带 app=nginx 标签的 Pod 的个数大于 2 的时候，就会有旧的 Pod 被删除；反之，就会有新的 Pod 被创建。

kubernetes项目控制器模块位置: https://github.com/kubernetes/kubernetes/tree/master/pkg/controller
 下面有很多的控制器模块(Deployment、job、、、)  
 这些控制器之所以被统一放在 pkg/controller 目录下，就是因为它们都遵循 Kubernetes 项目中的一个通用编排模式，即：控制循环（control loop）。
 比如，现在有一种待编排的对象 X，它有一个对应的控制器。那么，我就可以用一段 Go 语言风格的伪代码，为你描述这个控制循环：     
    for {
      实际状态 := 获取集群中对象 X 的实际状态（Actual State）
      期望状态 := 获取集群中对象 X 的期望状态（Desired State）
      if 实际状态 == 期望状态{
        什么都不做
      } else {
        执行编排动作，将实际状态调整为期望状态
      }
    }
    实际状态主要来自于kubernetes项目对Pod对象的监控状态   
    期望状态来自于用户的yaml文件,如: spec.replicas=2 这样的信息，这样的信息往往存在etcd中。
    Deployment控制器的大致逻辑如下:
        1、Deployment 控制器从 Etcd 中获取到所有携带了“app: nginx”标签的 Pod，然后统计它们的数量，这就是实际状态；
        2、Deployment 对象的 Replicas 字段的值就是期望状态；
        3、Deployment 控制器将两个状态做比较，然后根据比较结果，确定是创建 Pod，还是删除已有的 Pod。
     这个操作，通常被叫作调谐（Reconcile）。这个调谐的过程，则被称作“Reconcile Loop”（调谐循环）或者“Sync Loop”（同步循环）。
     在所有 API 对象的 Metadata 里，都有一个字段叫作 ownerReference，用于保存当前这个 API 对象的拥有者（Owner）的信息。  

Kubernetes 项目中的一个非常重要的概念（API 对象）：ReplicaSet - ReplicaSet.yaml
    这个ReplicaSet API对象其实可以理解为Deployment API对象的一个子集，其实Deployment是通过控制ReplicaSet控制Pod的扩缩容。
 ReplicaSet.yaml中定义了一个副本集为三和Pod的模版，而Deployment操作的就是Replicaset API对象，而不是Pod API对象。
 nginx-deployment.yaml 描述的关系如下:
   Deployment --->
                   ReplicaSet(2) --->
                                      Pod(1)                                            
                                      Pod(2)   
   ReplicaSet 负责通过“控制器模式”，保证系统中 Pod 的个数永远等于指定的个数（比如，2 个）。
   这也正是 Deployment 只允许容器的 restartPolicy=Always 的主要原因：只有在容器能保证自己始终是 Running 状态的前提下，ReplicaSet 调整 Pod 的个数才有意义。                                                                            
 Deployment 同样通过“控制器模式”，来操作 ReplicaSet 的个数和属性，进而实现“水平扩展 / 收缩”和“滚动更新”这两个编排动作。
 -----水平扩展
 $ kubectl scale deployment nginx-deployment --replicas=3
    deployment.apps/nginx-deployment scaled
 -----滚动升级   
 $ kubectl create -f nginx-deployment.yaml --record
    额外加了一个–record 参数。它的作用，是记录下你每次操作所执行的命令，以方便后面查看。
 $ kubectl get deployments
     NAME               READY   UP-TO-DATE   AVAILABLE   AGE
     nginx-deployment   2/2     2            2           14s     
   1、DESIRED：用户期望的 Pod 副本个数（spec.replicas 的值）；
   2、CURRENT：当前处于 Running 状态的 Pod 的个数；
   3、UP-TO-DATE：当前处于最新版本的 Pod 的个数，所谓最新版本指的是 Pod 的 Spec 部分与 Deployment 里 Pod 模板里定义的完全一致；
   4、AVAILABLE：当前已经可用的 Pod 的个数，即：既是 Running 状态，又是最新版本，并且已经处于 Ready（健康检查正确）状态的 Pod 的个数。  
 $ kubectl rollout status deployment/nginx-deployment
   deployment "nginx-deployment" successfully rolled out
   实时查看 Deployment 对象的状态变化
 $ kubectl get rs
    查看一下这个 Deployment 所控制的 ReplicaSet
 如果我们修改了 Deployment 的 Pod 模板，“滚动更新”就会被自动触发。
 修改 Deployment 有很多方法。比如，我可以直接使用 kubectl edit 指令编辑 Etcd 里的 API 对象。   
 $ kubectl edit deployment/nginx-deployment
 kubectl edit 指令编辑完成后，保存退出，Kubernetes 就会立刻触发“滚动更新”的过程。你还可以通过 kubectl rollout status 指令查看 nginx-deployment 的状态变化：  
 $ kubectl rollout status deployment/nginx-deployment
     Waiting for deployment "nginx-deployment" rollout to finish: 1 old replicas are pending termination...
     Waiting for deployment "nginx-deployment" rollout to finish: 1 old replicas are pending termination...
     deployment "nginx-deployment" successfully rolled out
 $ kubectl describe deployment nginx-deployment 
   ...
    Events:
      Type    Reason             Age    From                   Message
      ----    ------             ----   ----                   -------
      Normal  ScalingReplicaSet  8m44s  deployment-controller  Scaled up replica set nginx-deployment-7cf55fb7bb to 2
      Normal  ScalingReplicaSet  107s   deployment-controller  Scaled up replica set nginx-deployment-55df9cfb4b to 1
      Normal  ScalingReplicaSet  106s   deployment-controller  Scaled down replica set nginx-deployment-7cf55fb7bb to 1
      Normal  ScalingReplicaSet  106s   deployment-controller  Scaled up replica set nginx-deployment-55df9cfb4b to 2
      Normal  ScalingReplicaSet  87s    deployment-controller  Scaled down replica set nginx-deployment-7cf55fb7bb to 0
    可以看到replica-7cf55fb7bb从2->1的时候，replica-55df9cfb4b从0->1，106s 时 replica-55df9cfb4b从1->2，replica-7cf55fb7bb从1->0，这个就是滚动升级。
 $ kubectl get rs
    NAME                          DESIRED   CURRENT   READY   AGE
    nginx-deployment-55df9cfb4b   2         2         2       6m30s
    nginx-deployment-7cf55fb7bb   0         0         0       13m
 为了进一步保证服务的连续性，Deployment Controller 还会确保，在任何时间窗口内，只有指定比例的 Pod 处于离线状态。
 同时，它也会确保，在任何时间窗口内，只有指定比例的新 Pod 被创建出来。这两个比例的值都是可以配置的，默认都是 DESIRED 值的 25%。      
 如果上面这个 Deployment 的例子中，它有 3 个 Pod 副本，那么控制器在“滚动更新”的过程中永远都会确保至少有 2 个 Pod 处于可用状态，至多只有 4 个 Pod 同时存在于集群中。
 这个策略，是 Deployment 对象的一个字段，名叫 RollingUpdateStrategy，如下所示：
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nginx-deployment
      labels:
        app: nginx
    spec:
    ...
      strategy:
        type: RollingUpdate
        rollingUpdate:
          maxSurge: 1
          maxUnavailable: 1
    在上面这个 RollingUpdateStrategy 的配置中，maxSurge 指定的是除了 DESIRED 数量之外，在一次“滚动”中，Deployment 控制器还可以创建多少个新 Pod；
    而 maxUnavailable 指的是，在一次“滚动”中，Deployment 控制器可以删除多少个旧 Pod。
    同时，这两个配置还可以用前面我们介绍的百分比形式来表示，比如：maxUnavailable=50%，指的是我们最多可以一次删除“50%*DESIRED 数量”个 Pod。 
    Deployment、ReplicaSet、Pod 滚动升级关系图如下:
    Deployment -> 
                  ReplicaSet(V1) 
                                 -> Pod(1)  
                                 -> Pod(2)
                  ReplicaSet(V2) 
                                 -> Pod(1)  
    Deployment -> 
                  ReplicaSet(V1) 
                                 -> Pod(1)  
                  ReplicaSet(V2) 
                                 -> Pod(1)
                                 -> Pod(2)                              
    Deployment -> 
                  ReplicaSet(V1) 
                  ReplicaSet(V2) 
                                 -> Pod(1)
                                 -> Pod(2)  
                                 -> Pod(3)  
   通过这样的多个 ReplicaSet 对象，Kubernetes 项目就实现了对多个“应用版本”的描述。
 例子如下:
 设置一个错误的镜像
 $ kubectl set image deployment/nginx-deployment nginx=nginx:1.19.19
   deployment.apps/nginx-deployment image updated
 $ kubectl get rs
   NAME                          DESIRED   CURRENT   READY   AGE
   nginx-deployment-55df9cfb4b   2         2         2       21m
   nginx-deployment-59dd67fd98   1         1         0       10s
   nginx-deployment-7cf55fb7bb   0         0         0       28m
 如何让这个 Deployment 的 2 个 Pod，都回滚到以前的旧版本呢？                                  
 $ kubectl rollout undo deployment/nginx-deployment
    deployment.apps/nginx-deployment rolled back
  上面的具体操作上，Deployment 的控制器，其实就是让这个旧 ReplicaSet（hash=55df9cfb4b）再次“扩展”成 2 个 Pod，而让新的 ReplicaSet（hash=59dd67fd98）重新“收缩”到 0 个 Pod。   
 如果我想回滚到更早之前的版本，要怎么办呢？
 $ kubectl rollout history deployment/nginx-deployment 
   deployment.apps/nginx-deployment 
   REVISION  CHANGE-CAUSE
   1         kubectl apply --filename=nginx-deployment.yaml --record=true
   3         kubectl apply --filename=nginx-deployment.yaml --record=true
   5         kubectl apply --filename=nginx-deployment.yaml --record=true
 第一种  
 $ kubectl rollout history deployment/nginx-deployment --revision=5 
 第二种
 $ kubectl rollout undo deployment/nginx-deployment --to-revision=5
 Kubernetes 项目还提供了一个指令，使得我们对 Deployment 的多次更新操作，最后 只生成一个 ReplicaSet。
 具体的做法是，在更新 Deployment 前，你要先执行一条 kubectl rollout pause 指令。它的用法如下所示： 
 $ kubectl rollout pause deployment/nginx-deployment
    这个 kubectl rollout pause 的作用，是让这个 Deployment 进入了一个“暂停”状态。
  所以接下来，你就可以随意使用 kubectl edit 或者 kubectl set image 指令，修改这个 Deployment 的内容了。
    由于此时 Deployment 正处于“暂停”状态，所以我们对 Deployment 的所有修改，都不会触发新的“滚动更新”，也不会创建新的 ReplicaSet。
  而等到我们对 Deployment 修改操作都完成之后，只需要再执行一条 kubectl rollout resume 指令，就可以把这个 Deployment“恢复”回来，如下所示：
 $ kubectl rollout resume deployment/nginx-deployment
    在这个 kubectl rollout resume 指令执行之前，在 kubectl rollout pause 指令之后的这段时间里，我们对 Deployment 进行的所有修改，最后只会触发一次“滚动更新”。
 $ kubectl get rs
    通过检查 ReplicaSet 状态的变化，来验证一下 kubectl rollout pause 和 kubectl rollout resume 指令的执行效果    
不过，即使你像上面这样小心翼翼地控制了 ReplicaSet 的生成数量，随着应用版本的不断增加，Kubernetes 中还是会为同一个 Deployment 保存很多很多不同的 ReplicaSet。
那么，我们又该如何控制这些“历史”ReplicaSet 的数量呢？
很简单，Deployment 对象有一个字段，叫作 spec.revisionHistoryLimit，就是 Kubernetes 为 Deployment 保留的“历史版本”个数。所以，如果把它设置为 0，你就再也不能做回滚操作了。
 $ kubectl edit deployment/nginx-deployment
     ...
     revisionHistoryLimit : 0 # 默认是10，将这个数字改成0，那么就不会在存在ReplicaSet的历史记录了。
     ... 
 $ kubectl get rs
    NAME                          DESIRED   CURRENT   READY   AGE
    nginx-deployment-7cf55fb7bb   2         2         2       50m
 $ kubectl rollout history deployment/nginx-deployment
   deployment.apps/nginx-deployment 
   REVISION  CHANGE-CAUSE
   6         kubectl apply --filename=nginx-deployment.yaml --record=true   
总结: Deployment 控制 ReplicaSet（版本），ReplicaSet 控制 Pod（副本数）。这个两层控制关系一定要牢记。
```

### StatefulSet (https://kubernetes.io/zh/docs/tutorials/stateful-application/basic-stateful-set/)[]
```text
StatefulSet 的核心功能，就是通过某种方式记录这些状态，然后在 Pod 被重新创建时，能够为新 Pod 恢复这些状态。    
Headless Service: Service 是 Kubernetes 项目中用来将一组 Pod 暴露给外界访问的一种机制。比如，一个 Deployment 有 3 个 Pod，那么我就可以定义一个 Service。然后，用户只要能访问到这个 Service，它就能访问到某个具体的 Pod。
那么，这个 Service 又是如何被访问的呢？
1、第一种方式，是以 Service 的 VIP（Virtual IP，即：虚拟 IP）方式。比如：当我访问 10.0.23.1 这个 Service 的 IP 地址时，10.0.23.1 其实就是一个 VIP，它会把请求转发到该 Service 所代理的某一个 Pod 上。
2、第二种方式，就是以 Service 的 DNS 方式。比如：这时候，只要我访问“my-svc.my-namespace.svc.cluster.local”这条 DNS 记录，就可以访问到名叫 my-svc 的 Service 所代理的某一个 Pod。
    第一种处理方法，是 Normal Service。这种情况下，你访问“my-svc.my-namespace.svc.cluster.local”解析到的，正是 my-svc 这个 Service 的 VIP，后面的流程就跟 VIP 方式一致了。
    第二种处理方法，正是 Headless Service。这种情况下，你访问“my-svc.my-namespace.svc.cluster.local”解析到的，直接就是 my-svc 代理的某一个 Pod 的 IP 地址。可以看到，这里的区别在于，Headless Service 不需要分配一个 VIP，而是可以直接以 DNS 记录的方式解析出被代理 Pod 的 IP 地址。
Headless Service 对应的 YAML 文件 -> nginx-service.yaml :
    所谓的 Headless Service，其实仍是一个标准 Service 的 YAML 文件。只不过，它的 clusterIP 字段的值是：None，即：这个 Service，没有一个 VIP 作为“头”。这也就是 Headless 的含义。所以，这个 Service 被创建后并不会被分配一个 VIP，而是会以 DNS 记录的方式暴露出它所代理的 Pod。
    它所代理的所有 Pod 的 IP 地址，都会被绑定一个这样格式的 DNS 记录，如下所示：
        <pod-name>.<svc-name>.<namespace>.svc.cluster.local
      有了这个“可解析身份”，只要你知道了一个 Pod 的名字，以及它对应的 Service 的名字，你就可以非常确定地通过这条 DNS 记录访问到 Pod 的 IP 地址。
创建一个StatefulSet.yaml，内容与nginx-deployment.yaml 的唯一区别，就是多了一个 serviceName=nginx 字段。
这个字段的作用，就是告诉 StatefulSet 控制器，在执行控制循环（Control Loop）的时候，请使用 nginx 这个 Headless Service 来保证 Pod 的“可解析身份”。
创建上面的两个API对象 (StatefulSet、Service)
$ kubectl create -f nginx-service.yaml
$ kubectl get service nginx
$ kubectl create -f StatefulSet.yaml
$ kubectl get statefulset web
$ kubectl get pods -w -l app=nginx
    NAME    READY   STATUS              RESTARTS   AGE
    web-0   1/1     Running             0          50s
    web-1   0/1     ContainerCreating   0          18s
    web-1   1/1     Running             0          47s
   通过上面这个 Pod 的创建过程，我们不难看到，StatefulSet 给它所管理的所有 Pod 的名字，进行了编号，编号规则是：-。
   而且这些编号都是从 0 开始累加，与 StatefulSet 的每个 Pod 实例一一对应，绝不重复。
   更重要的是，这些 Pod 的创建，也是严格按照编号顺序进行的。比如，在 web-0 进入到 Running 状态、并且细分状态（Conditions）成为 Ready 之前，web-1 会一直处于 Pending 状态。
   备注：Ready 状态再一次提醒了我们，为 Pod 设置 livenessProbe 和 readinessProbe 的重要性。
   当这两个 Pod 都进入了 Running 状态之后，你就可以查看到它们各自唯一的“网络身份”了。    
$ kubectl exec web-0 -- sh -c 'hostname'     
 web-0
$ kubectl exec web-1 -- sh -c 'hostname'     
 web-1
$ for i in 0 1; do kubectl exec web-$i -- sh -c 'hostname'; done 
这两个 Pod 的 hostname 与 Pod 名字是一致的，都被分配了对应的编号。接下来，我们再试着以 DNS 的方式，访问一下这个 Headless Service： 
$ kubectl run -i --tty --image busybox:1.28.4 dns-test --restart=Never --rm /bin/sh
/ # nslookup web-0.nginx    =========注意这里用1.28.4这个busybox镜像，最新的镜像可能存在问题
如果这里拿不到这个准确的ip地址、可以用以下的流程操作
$ kubectl exec web-0 -- sh -c 'cat /etc/hosts'
    # Kubernetes-managed hosts file.
    127.0.0.1	localhost
    ::1	localhost ip6-localhost ip6-loopback
    fe00::0	ip6-localnet
    fe00::0	ip6-mcastprefix
    fe00::1	ip6-allnodes
    fe00::2	ip6-allrouters
    10.244.1.94	web-0.nginx.default.svc.cluster.local	web-0
/ # nslookup web-0.nginx.default.svc.cluster.local   
/ # nslookup web-0.nginx.default  
/ # ping web-0.nginx
当前Terminal启动pod监听
$ kubectl get pods -w -l app=nginx
另外一个 Terminal 里把这两个“有状态应用”的 Pod 删掉：
$ kubectl delete pods -l app=nginx
可以看到，当我们把这两个 Pod 删除之后，Kubernetes 会按照原先编号的顺序，创建出了两个新的 Pod。并且，Kubernetes 依然为它们分配了与原来相同的“网络身份”：web-0.nginx 和 web-1.nginx。
通过这种严格的对应规则，StatefulSet 就保证了 Pod 网络标识的稳定性。
如果 web-0 是一个需要先启动的主节点，web-1 是一个后启动的从节点，那么只要这个 StatefulSet 不被删除，你访问 web-0.nginx 时始终都会落在主节点上，访问 web-1.nginx 时，则始终都会落在从节点上，这个关系绝对不会发生任何变化。
$ kubectl run -i --tty --image busybox:1.28.4 dns-test --restart=Never --rm /bin/sh 
/ # nslookup web-1.nginx
重点: StatefulSet 这个控制器的主要作用之一，就是使用 Pod 模板创建 Pod 的时候，对它们进行编号，并且按照编号顺序逐一完成创建工作。
    而当 StatefulSet 的“控制循环”发现 Pod 的“实际状态”与“期望状态”不一致，需要新建或者删除 Pod 进行“调谐”的时候，它会严格按照这些 Pod 编号的顺序，逐一完成这些操作。 
```

### Persistent Volume Claim (https://kubernetes.io/docs/concepts/storage/persistent-volumes)[]
```text
产生原因: Volume 的管理和远程持久化存储的知识，不仅超越了开发者的知识储备，还会有暴露公司基础设施秘密的风险。
下面这个例子，就是一个声明了 Ceph RBD 类型 Volume 的 Pod：
    apiVersion: v1
    kind: Pod
    metadata:
      name: rbd
    spec:
      containers:
        - image: kubernetes/pause
          name: rbd-rw
          volumeMounts:
          - name: rbdpd
            mountPath: /mnt/rbd
      volumes:
        - name: rbdpd
          rbd:
            monitors:
            - '10.16.154.78:6789'
            - '10.16.154.82:6789'
            - '10.16.154.83:6789'
            pool: kube
            image: foo
            fsType: ext4
            readOnly: true
            user: admin
            keyring: /etc/ceph/keyring
            imageformat: "2"
            imagefeatures: "layering"
    其一，如果不懂得 Ceph RBD 的使用方法，那么这个 Pod 里 Volumes 字段，你十有八九也完全看不懂。
    其二，这个 Ceph RBD 对应的存储服务器的地址、用户名、授权文件的位置，也都被轻易地暴露给了全公司的所有开发人员。
    
引入新的API: Kubernetes 项目引入了一组叫作 Persistent Volume Claim（PVC）和 Persistent Volume（PV）的 API 对象，大大降低了用户声明和使用持久化 Volume 的门槛。 

举例: 有了 PVC 之后，一个开发人员想要使用一个 Volume，只需要简单的两步即可。
   第一步：定义一个 PVC，声明想要的 Volume 的属性：
        kind: PersistentVolumeClaim
        apiVersion: v1
        metadata:
          name: pv-claim
        spec:
          accessModes:
          - ReadWriteOnce
          resources:
            requests:
              storage: 1Gi
        具体含义: 在这个 PVC 对象里，不需要任何关于 Volume 细节的字段，只有描述性的属性和定义。
                比如，storage: 1Gi，表示我想要的 Volume 大小至少是 1 GiB；accessModes: ReadWriteOnce，
                表示这个 Volume 的挂载方式是可读写，并且只能被挂载在一个节点上而非被多个节点共享。      
   第二步：在应用的 Pod 中，声明使用这个 PVC：
        apiVersion: v1
        kind: Pod
        metadata:
          name: pv-pod
        spec:
          containers:
            - name: pv-container
              image: nginx
              ports:
                - containerPort: 80
                  name: "http-server"
              volumeMounts:
                - mountPath: "/usr/share/nginx/html"
                  name: pv-storage
          volumes:
            - name: pv-storage
              persistentVolumeClaim:
                claimName: pv-claim
        具体含义: 在这个 Pod 的 Volumes 定义中，我们只需要声明它的类型是 persistentVolumeClaim，
                然后指定 PVC 的名字，而完全不必关心 Volume 本身的定义。 
   ==> 疑问？ 这些符合条件的 Volume 又是从哪里来的
    ===> 下面的pv来自于ceph，可以利用k8s搭建一个rook集群，具体见文件<kubeadm安装3node集群.txt>
     官方地址: https://rook.io/docs/rook/v1.4/ceph-quickstart.html
              https://rook.io/docs/rook/v1.4/ceph-toolbox.html
     获取安装后的pod
     $ kubectl get pods -n rook-ceph
     获取安装后的service信息
     $ kubectl get services -n rook-ceph
     它们来自于由运维人员维护的 PV（Persistent Volume）对象。下面是一个常见的 PV 对象的 YAML 定义文件：
         kind: PersistentVolume
         apiVersion: v1
         metadata:
           name: pv-volume
           labels:
             type: local
         spec:
           capacity:
             storage: 10Gi
           rbd:
             monitors:
             - '10.16.154.78:6789'
             - '10.16.154.82:6789'
             - '10.16.154.83:6789'
             pool: kube
             image: foo
             fsType: ext4
             readOnly: true
             user: admin
             keyring: /etc/ceph/keyring
             imageformat: "2"
             imagefeatures: "layering"
         具体含义: 这个 PV 对象的 spec.rbd 字段，正是我们前面介绍过的 Ceph RBD Volume 的详细定义。
                而且，它还声明了这个 PV 的容量是 10 GiB。这样，Kubernetes 就会为我们刚刚创建的 PVC 对象绑定这个 PV。   
                
Kubernetes 中 PVC 和 PV 的设计，实际上类似于“接口”和“实现”的思想。开发者只要知道并会使用“接口”，即：PVC；而运维人员则负责给“接口”绑定具体的实现，即：PV。       

创建一个basic-stateful-set.yaml 声明statefulSet信息:
      为这个 StatefulSet 额外添加了一个 volumeClaimTemplates 字段。从名字就可以看出来，它跟 Deployment 里 Pod 模板（PodTemplate）的作用类似。
   也就是说，凡是被这个 StatefulSet 管理的 Pod，都会声明一个对应的 PVC；而这个 PVC 的定义，就来自于 volumeClaimTemplates 这个模板字段。
   更重要的是，这个 PVC 的名字，会被分配一个与这个 Pod 完全一致的编号。
   这个自动创建的 PVC，与 PV 绑定成功后，就会进入 Bound 状态，这就意味着这个 Pod 可以挂载并使用这个 PV 了。
实际操作:
  $ kubectl apply -f basic-stateful-set.yaml
  $ kubectl get pvc -l app=nginx
    NAME        STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
    www-web-0   Bound     pvc-15c268c7-b507-11e6-932f-42010a800002   1Gi        RWO           48s
    www-web-1   Bound     pvc-15c79307-b507-11e6-932f-42010a800002   1Gi        RWO           48s
这些 PVC，都以“<PVC 名字 >-<StatefulSet 名字 >-< 编号 >”的方式命名，并且处于 Bound 状态。
我们就可以使用如下所示的指令，在 Pod 的 Volume 目录里写入一个文件，来验证一下上述 Volume 的分配情况：
  $ for i in 0 1; do kubectl exec web-$i -- sh -c 'echo hello $(hostname) > /usr/share/nginx/html/index.html'; done
在这个 Pod 容器里访问“http://localhost”，你实际访问到的就是 Pod 里 Nginx 服务器进程，而它会为你返回 /usr/share/nginx/html/index.html 里的内容。这个操作的执行方法如下所示：
  $ for i in 0 1; do kubectl exec -it web-$i -- curl localhost; done
如果你使用 kubectl delete 命令删除这两个 Pod，这些 Volume 里的文件会不会丢失呢？
  $ kubectl delete pod -l app=nginx
    pod "web-0" deleted
    pod "web-1" deleted
在被重新创建出来的 Pod 容器里访问 http://localhost
  $ kubectl exec -it web-0 -- curl localhost     
    hello web-0
  就会发现，这个请求依然会返回：hello web-0。也就是说，原先与名叫 web-0 的 Pod 绑定的 PV，在这个 Pod 被重新创建之后，依然同新的名叫 web-0 的 Pod 绑定在了一起。
  当你把一个 Pod，比如 web-0，删除之后，这个 Pod 对应的 PVC 和 PV，并不会被删除，而这个 Volume 里已经写入的数据，也依然会保存在远程存储服务里（比如，我们在这个例子里用到的 Ceph 服务器）。      
```      
                               















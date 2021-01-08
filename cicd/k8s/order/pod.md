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
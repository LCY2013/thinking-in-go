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
   
   kubectl get pod --selector=role=frontend -o yaml       
  
kubernetes 启动 podPreset 需要在(/etc/kubernetes/manifests/kube-apiserver.yaml)添加下面信息(spec.containers.command):
- --enable-admission-plugins=NodeRestriction,PodPreset
- --runtime-config=settings.k8s.io/v1alpha1=true
修改后kube-apiserver就会自动重启
      
````
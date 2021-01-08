### 编排之Deployment
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
 第一步 
 $ kubectl rollout history deployment/nginx-deployment --revision=5 
 第二步
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
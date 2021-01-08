### Job 和 CronJob

Deployment、StatefulSet、DaemonSet编排主要用在了长时间运行的项目上面，如nginx、mysql等，如果是某一类作业是"离线任务"，那么这三种编排方式就不能够满足条件。

在 Borg 项目中，Google 就已经对作业进行了分类处理，提出了 LRS（LongRunning Service）和 Batch Jobs 两种作业形态，对它们进行“分别管理”和“混合调度”。

在 2015 年 Borg 论文刚刚发布的时候，Kubernetes 项目并不支持对 Batch Job 的管理。直到 v1.4 版本之后，社区才逐步设计出了一个用来描述离线业务的 API 对象，它的名字就是：Job。

Job API 对象的定义示例如下：
```yaml
apiVersion: batch/v1
kind: Job
metadata:  
  name: pi
spec:  
  template:    
    spec:      
      containers:      
      - name: pi        
        image: resouer/ubuntu-bc         
        command: ["sh", "-c", "echo 'scale=10000; 4*a(1)' | bc -l "]      
      restartPolicy: Never  
  backoffLimit: 4
```
其中，bc 命令是 Linux 里的“计算器”；-l 表示，现在要使用标准数学库；而 a(1)，则是调用数学库中的 arctangent 函数，计算 atan(1)。这是什么意思呢？

tan(π/4) = 1，所以，4*atan(1)正好就是π，也就是3.1415926。

所以，这就是一个计算π值的容器。而通过 scale=10000，指定了输出的小数点后的位数是 10000。
```text
创建这个job
$ kubectl create -f pi-job.yaml

查询这个job信息
$ kubectl describe job/pi
Name:           pi
Namespace:      default
Selector:       controller-uid=c08ac26b-7ea8-4491-bb35-2ec10d2b5396
Labels:         controller-uid=c08ac26b-7ea8-4491-bb35-2ec10d2b5396
                job-name=pi
Annotations:    <none>
Parallelism:    1
Completions:    1
Start Time:     Fri, 08 Jan 2021 10:04:12 +0800
Pods Statuses:  1 Running / 0 Succeeded / 0 Failed
Pod Template:
  Labels:  controller-uid=c08ac26b-7ea8-4491-bb35-2ec10d2b5396
           job-name=pi
  Containers:
   pi:
    Image:      resouer/ubuntu-bc
    Port:       <none>
    Host Port:  <none>
    Command:
      sh
      -c
      echo 'scale=10000; 4*a(1)' | bc -l 
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Events:
  Type    Reason            Age   From            Message
  ----    ------            ----  ----            -------
  Normal  SuccessfulCreate  70s   job-controller  Created pod: pi-8dxlw

可以看到，这个 Job 对象在创建后，它的 Pod 模板，被自动加上了一个 controller-uid=< 一个随机字符串 > 这样的 Label。而这个 Job 对象本身，则被自动加上了这个 Label 对应的 Selector，从而 保证了 Job 与它所管理的 Pod 之间的匹配关系。

而 Job Controller 之所以要使用这种携带了 UID 的 Label，就是为了避免不同 Job 对象所管理的 Pod 发生重合。需要注意的是，这种自动生成的 Label 对用户来说并不友好，所以不太适合推广到 Deployment 等长作业编排对象上。

可以看到这个 Job 创建的 Pod 进入了 Running 状态，这意味着它正在计算Pi 的值。
$ kubectl get pods
$ kubectl get pods
NAME                                      READY   STATUS      RESTARTS   AGE
pi-8dxlw                                  1/1     Running   0          1m16s
而几分钟后计算结束，这个 Pod 就会进入 Completed 状态：
$ kubectl get pods
NAME                                      READY   STATUS      RESTARTS   AGE
pi-8dxlw                                  0/1     Completed   0          3m36s
这也是需要在 Pod 模板中定义 restartPolicy=Never 的原因：离线计算的 Pod 永远都不应该被重启，否则它们会再重新计算一遍。
事实上，restartPolicy 在 Job 对象里只允许被设置为 Never 和 OnFailure；而在 Deployment 对象里，restartPolicy 则只允许被设置为 Always。

此时，通过 kubectl logs 查看一下这个 Pod 的日志，就可以看到计算得到的 Pi 值已经被打印了出来：
$ kubectl logs pi-8dxlw
3.141592653589793238462643383279502884197169399375105820974944592307\
81640628620899862803482534211706798214808651328230664709384460955058\
....
```

#### 如果这个离线作业失败了要怎么办？
在这个例子中定义了 restartPolicy=Never，那么离线作业失败后 JobController 就会不断地尝试创建一个新 Pod。

当然，这个尝试肯定不能无限进行下去。所以，就在 Job 对象的 spec.backoffLimit字段里定义了重试次数为 4（即，backoffLimit=4），而这个字段的默认值是 6。

需要注意的是，Job Controller 重新创建 Pod 的间隔是呈指数增加的，即下一次重新创建Pod 的动作会分别发生在 10 s、20 s、40 s ...后。

如果定义的 restartPolicy=OnFailure，那么离线作业失败后，Job Controller 就不会去尝试创建新的 Pod。但是，它会不断地尝试重启 Pod 里的容器。

#### 当一个 Job 的 Pod 运行结束后，它会进入 Completed 状态。但是，如果这个Pod 因为某种原因一直不肯结束呢？
在 Job 的 API 对象里，有一个 spec.activeDeadlineSeconds 字段可以设置最长运行时间，比如：
```yaml
spec: 
  backoffLimit: 5 
  activeDeadlineSeconds: 100
```
一旦运行超过了 100 s，这个 Job 的所有 Pod 都会被终止。并且，可以在 Pod 的状态里看到终止的原因是 reason: DeadlineExceeded。

#### 离线业务之所以被称为Batch Job，是因为可以以“Batch”方式，也就是并行的方式去运行
Job 中控制并行的两个参数信息：
```text
1. spec.parallelism：它定义的是一个 Job 在任意时间最多可以启动多少个 Pod 同时运行；

2. spec.completions：它定义的是 Job 至少要完成的 Pod 数目，即 Job 的最小完成数。
```

在pi-job.yaml示例中添加这两个参数信息：
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  parallelism: 2
  completions: 4
  template:
    spec:
      containers:
        - name: pi
          image: resouer/ubuntu-bc
          command: ["sh", "-c", "echo 'scale=10000; 4*a(1)' | bc -l "]
      restartPolicy: Never
  backoffLimit: 4
```
指定了这个 Job 最大的并行数是 2，而最小的完成数是 4。

下面创建这个job信息:
```text
$ kubectl create -f pi-batch-job.yaml

这个 Job 其实也维护了两个状态字段，即 DESIRED 和 SUCCESSFUL，如下所示：
$ kubectl get job
NAME   COMPLETIONS   DURATION   AGE
pi     0/4           31s        31s
COMPLETIONS 定义的最小完成数。
$ kubectl get job
NAME   COMPLETIONS   DURATION   AGE
pi     4/4           5m43s      13m

Job 首先创建了两个并行运行的 Pod 来计算 Pi：
$ kubectl get pods
NAME                                      READY   STATUS      RESTARTS     AGE
pi-56mbn                                  1/1     Running       0          18s
pi-r5kv5                                  1/1     Running       0          15s

每当有一个 Pod 完成计算进入 Completed 状态时，就会有一个新的Pod 被自动创建出来，并且快速地从 Pending 状态进入到 ContainerCreating 状态：
$ kubectl get pods
NAME                                      READY   STATUS      RESTARTS   AGE
pi-56mbn                                  0/1     Completed   0          21m
pi-r5kv5                                  0/1     Completed   0          18m
pi-rmg5v                                  0/1     Completed   0          21m
pi-w7zwq                                  0/1     Completed   0          18m
```

#### Job Controller 原理分析 
首先，Job Controller 控制的对象，直接就是 Pod。

其次，Job Controller 在控制循环中进行的调谐（Reconcile）操作，是根据实际在Running 状态 Pod 的数目、已经成功退出的 Pod 的数目，以及 parallelism、completions 参数的值共同计算出在这个周期里，应该创建或者删除的 Pod 数目，然后调用 Kubernetes API 来执行这个操作。

以创建 Pod 为例。在上面计算 Pi 值的这个例子中，当 Job 一开始创建出来时，实际处于Running 状态的 Pod 数目 =0，已经成功退出的 Pod 数目 =0，而用户定义的completions，也就是最终用户需要的 Pod 数目 =4。

所以，在这个时刻，需要创建的 Pod 数目 = 最终需要的 Pod 数目 - 实际在 Running 状态Pod 数目 - 已经成功退出的 Pod 数目 = 4 - 0 - 0= 4。也就是说，Job Controller 需要创建 4 个 Pod 来纠正这个不一致状态。

可是，定义了这个 Job 的 parallelism=2。也就是说，规定了每次并发创建的Pod 个数不能超过 2 个。所以，Job Controller 会对前面的计算结果做一个修正，修正后的期望创建的 Pod 数目应该是：2 个。

这时候，Job Controller 就会并发地向 kube-apiserver 发起两个创建 Pod 的请求。

如果在这次调谐周期里，Job Controller 发现实际在 Running 状态的 Pod 数目，比 parallelism 还大，那么它就会删除一些 Pod，使两者相等。

综上所述，Job Controller 实际上控制了，作业执行的并行度，以及总共需要完成的任务数这两个重要参数。而在实际使用时，你需要根据作业的特性，来决定并行度（parallelism）和任务数（completions）的合理取值。

### [三种常用的、使用 Job 对象的方法](https://kubernetes.io/docs/tasks/job/fine-parallel-processing-work-queue/)

#### [第一种用法，也是最简单粗暴的用法：外部管理器 +Job 模板](https://kubernetes.io/docs/tasks/job/parallel-processing-expansion/)
这种模式的特定用法是：把 Job 的 YAML 文件定义为一个“模板”，然后用一个外部工具控制这些“模板”来生成 Job。这时，Job 的定义方式如下所示：
```yaml
apiVersion: batch/v1
kind: Job
metadata:  
  name: process-item-$ITEM  
  labels:    
    jobgroup: jobexample
spec:  
  template:    
    metadata:      
      name: jobexample      
      labels:        
        jobgroup: jobexample    
    spec:      
      containers:      
      - name: c        
        image: busybox        
        command: ["sh", "-c", "echo Processing item $ITEM && sleep 5"]      
      restartPolicy: Never
```
可以看到，我们在这个 Job 的 YAML 里，定义了 $ITEM 这样的“变量”。

所以，在控制这种 Job 时，只要注意如下两个方面即可：
```text
1. 创建 Job 时，替换掉 $ITEM 这样的变量；

2. 所有来自于同一个模板的 Job，都有一个 jobgroup: jobexample 标签，也就是说这一组 Job 使用这样一个相同的标识。
```

而做到第一点非常简单。比如，可以通过这样一句 shell 把 $ITEM 替换掉：
```text
$ mkdir ./jobs
$ for i in apple banana cherry
do  
    cat job-tmpl.yaml | sed "s/\$ITEM/$i/" > ./jobs/job-$i.yaml
done
```
这样，一组来自于同一个模板的不同 Job 的 yaml 就生成了。接下来，就可以通过一句kubectl create 指令创建这些 Job 了：
```text
$ kubectl create -f ./jobs
$ kubectl get pods -l jobgroup=jobexample
NAME                        READY   STATUS      RESTARTS   AGE
process-item-apple-v4l5b    0/1     Completed   0          99s
process-item-banana-kcwfd   0/1     Completed   0          99s
process-item-cherry-w6cp4   0/1     Completed   0          99s
```
这个模式看起来虽然很“傻”，但却是 Kubernetes 社区里使用 Job 的一个很普遍的模式。

原因很简单：大多数用户在需要管理 Batch Job 的时候，都已经有了一套自己的方案，需要做的往往就是集成工作。这时候，Kubernetes 项目对这些方案来说最有价值的，就是Job 这个 API 对象。所以，只需要编写一个外部工具（等同于我们这里的 for 循环）来管理这些 Job 即可。

这种模式最典型的应用，就是 TensorFlow 社区的 KubeFlow 项目。

很容易理解，在这种模式下使用 Job 对象，completions 和 parallelism 这两个字段都应该使用默认值 1，而不应该由我们自行设置。而作业 Pod 的并行控制，应该完全交由外部工具来进行管理（比如，KubeFlow）。

#### [第二种用法：拥有固定任务数目的并行 Job](https://kubernetes.io/docs/tasks/job/coarse-parallel-processing-work-queue/)
这种模式下，只关心最后是否有指定数目（spec.completions）个任务成功退出。至于执行时的并行度是多少，并不关心。

比如，这个计算 Pi 值的例子，就是这样一个典型的、拥有固定任务数目（completions=4）的应用场景。 它的 parallelism 值是 2；或者，你可以干脆不指定parallelism，直接使用默认的并行度（即：1）。

还可以使用一个工作队列（Work Queue）进行任务分发。这时，Job 的 YAML文件定义如下所示：
```yaml
apiVersion: batch/v1
kind: Job
metadata:  
  name: job-wq-1
spec:  
  completions: 8  
  parallelism: 2  
  template:    
    metadata:      
      name: job-wq-1    
    spec:      
      containers:      
      - name: c        
        image: luochunyun/job-wq-1        
        env:        
        - name: BROKER_URL          
          value: amqp://guest:guest@rabbitmq-service:5672        
        - name: QUEUE          
          value: job1      
      restartPolicy: OnFailure
```
可以看到，completions 的值是：8，这意味着总共要处理的任务数目是 8个。也就是说，总共会有 8 个任务会被逐一放入工作队列里（可以运行一个外部小程序作为生产者，来提交任务）。

在这个实例中，选择充当工作队列的是一个运行在 Kubernetes 里的 RabbitMQ。所以，需要在 Pod 模板里定义 BROKER_URL，来作为消费者。

所以，一旦用 kubectl create 创建了这个 Job，它就会以并发度为 2 的方式，每两个Pod 一组，创建出 8 个 Pod。每个 Pod 都会去连接 BROKER_URL，从 RabbitMQ 里读取任务，然后各自进行处理。这个 Pod 里的执行逻辑，可以用这样一段伪代码来表示：
```text
/* job-wq-1 的伪代码 */
queue := newQueue($BROKER_URL, $QUEUE)
task := queue.Pop()
process(task)
exit
```
可以看到，每个 Pod 只需要将任务信息读取出来，处理完成，然后退出即可。而作为用户，只关心最终一共有 8 个计算任务启动并且退出，只要这个目标达到，就认为整个Job 处理完成了。所以说，这种用法，对应的就是“任务总数固定”的场景。

#### [第三种用法，也是很常用的一个用法：指定并行度（parallelism），但不设置固定的completions 的值](https://kubernetes.io/docs/tasks/job/fine-parallel-processing-work-queue/)
此时，就必须自己想办法，来决定什么时候启动新 Pod，什么时候 Job 才算执行完成。在这种情况下，任务的总数是未知的，所以你不仅需要一个工作队列来负责任务分发，还需要能够判断工作队列已经为空（即：所有的工作已经结束了）。

这时候，Job 的定义基本上没变化，只不过是不再需要定义 completions 的值了而已：
```yaml
apiVersion: batch/v1
kind: Job
metadata:  
  name: job-wq-2
spec:  
  parallelism: 2  
  template:    
    metadata:      
      name: job-wq-2  
    spec:      
      containers:      
      - name: c        
        image: luochunyun/job-wq-2        
        env:        
        - name: BROKER_URL          
          value: amqp://guest:guest@rabbitmq-service:5672        
        - name: QUEUE          
          value: job2     
      restartPolicy: OnFailure
```
而对应的 Pod 的逻辑会稍微复杂一些，可以用这样一段伪代码来描述：
```text
/* job-wq-2 的伪代码 */
for !queue.IsEmpty($BROKER_URL, $QUEUE) {  
    task := queue.Pop()  
    process(task)
}
print("Queue empty, exiting")
exit
```
由于任务数目的总数不固定，所以每一个 Pod 必须能够知道，自己什么时候可以退出。比如，在这个例子中，简单地以“队列为空”，作为任务全部完成的标志。所以说，这种用法，对应的是“任务总数不固定”的场景。

不过，在实际的应用中，需要处理的条件往往会非常复杂。比如，任务完成后的输出、每个任务 Pod 之间是不是有资源的竞争和协同等等。

### [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
CronJob 描述的，正是定时任务，它的 API 对象，如下所示：
```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            imagePullPolicy: IfNotPresent
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```
在这个 YAML 文件中，最重要的关键词就是jobTemplate。CronJob 是一个 Job 对象的控制器（Controller）！

CronJob 与 Job 的关系，如同 Deployment 与 Pod 的关系一样，CronJob 是一个专门用来管理 Job 对象的控制器。它是创建和删除 Job 的依据，是 schedule 字段定义的、一个标准的 [Unix Cron](https://en.wikipedia.org/wiki/Cron) 格式的表达式。

Cron schedule syntax
```text
# ┌───────────── minute (0 - 59)
# │ ┌───────────── hour (0 - 23)
# │ │ ┌───────────── day of the month (1 - 31)
# │ │ │ ┌───────────── month (1 - 12)
# │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
# │ │ │ │ │                                   7 is also Sunday on some systems)
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
```

```text
上述例子中 CronJob 对象在创建 1 分钟后，就会有一个 Job 产生了，如下所示：
$ kubectl create -f ./cronjob.yaml

$ kubectl get jobs
NAME               COMPLETIONS   DURATION   AGE
hello-1610088840   1/1           2s         11s

CronJob 对象会记录下这次 Job 执行的时间：
$ kubectl get cronjob hello
NAME    SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
hello   */1 * * * *   False     0        55s             87s
```

需要注意的是，由于定时任务的特殊性，很可能某个 Job 还没有执行完，另外一个新 Job就产生了。这时候，可以通过 spec.concurrencyPolicy 字段来定义具体的处理策略。比如：
```text
1. concurrencyPolicy=Allow，这也是默认情况，这意味着这些 Job 可以同时存在；
2. concurrencyPolicy=Forbid，这意味着不会创建新的 Pod，该创建周期被跳过；
3. concurrencyPolicy=Replace，这意味着新产生的 Job 会替换旧的、没有执行完的Job。

而如果某一次 Job 创建失败，这次创建就会被标记为“miss”。当在指定的时间窗口内，miss 的数目达到 100 时，那么 CronJob 会停止再创建这个 Job。

这个时间窗口，可以由 spec.startingDeadlineSeconds 字段指定。比如startingDeadlineSeconds=200，意味着在过去 200 s 里，如果 miss 的数目达到了 100次，那么这个 Job 就不会被创建执行了。
```













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

```











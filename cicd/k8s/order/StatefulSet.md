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

### 部署一个mysql集群

#### 需求
```text
1. 是一个“主从复制”（Maser-Slave Replication）的 MySQL 集群；
2. 有 1 个主节点（Master）；
3. 有多个从节点（Slave）；
4. 从节点需要能水平扩展；
5. 所有的写操作，只能在主节点上执行；
6. 读操作可以在所有节点上执行。
```

#### 问题分析
```text
在常规环境里，部署这样一个主从模式的 MySQL 集群的主要难点在于：如何让从节点能够拥有主节点的数据，即：如何配置主（Master）从（Slave）节点的复制与同步。

在安装好 MySQL 的 Master 节点之后，你需要做的第一步工作，就是通过XtraBackup 将 Master 节点的数据备份到指定目录。
    备注：XtraBackup 是业界主要使用的开源 MySQL 备份和恢复工具。

#### 第一步
这一步会自动在目标目录里生成一个备份信息文件，名叫：xtrabackup_binlog_info，这个文件一般会包含如下两个信息：
$ cat xtrabackup_binlog_info
TheMaster-bin.000001     481
这两个信息会在接下来配置 Slave 节点的时候用到

#### 第二步
配置 Slave 节点，Slave 节点在第一次启动前，需要先把 Master 节点的备份数据，连同备份信息文件，一起拷贝到自己的数据目录（/var/lib/mysql）下，然后再执行这样一句 SQL：
TheSlave|mysql> CHANGE MASTER TO                
    MASTER_HOST='$masterip',                
    MASTER_USER='xxx',                
    MASTER_PASSWORD='xxx',                
    MASTER_LOG_FILE='TheMaster-bin.000001',                
    MASTER_LOG_POS=481;

其中，MASTER_LOG_FILE 和 MASTER_LOG_POS，就是该备份对应的二进制日志（Binary Log）文件的名称和开始的位置（偏移量），也正是 xtrabackup_binlog_info 文件里的那两部分内容（即：TheMaster-bin.000001 和 481）。

#### 第三步
启动 Slave 节点，在这一步，需要执行这样一句 SQL：
TheSlave|mysql> START SLAVE;

这样，Slave 节点就启动了，它会使用备份信息文件中的二进制日志文件和偏移量，与主节点进行数据同步。

#### 第四步
在这个集群中添加更多的 Slave 节点：
需要注意的是，新添加的 Slave 节点的备份数据，来自于已经存在的 Slave 节点。
所以，在这一步，需要将 Slave 节点的数据备份在指定目录。
而这个备份操作会自动生成另一种备份信息文件，名叫：xtrabackup_slave_info。
同样地，这个文件也包含了MASTER_LOG_FILE 和 MASTER_LOG_POS 两个字段。
然后，就可以执行跟前面一样的“CHANGE MASTER TO”和“START SLAVE” 指令，来初始化并启动这个新的 Slave 节点了。
```

#### 通过问题分析，可以看出部署 MySQL 集群的流程迁移到 Kubernetes 项目上，需要能够“容器化”地解决下面的“三座大山”
```text
1. Master 节点和 Slave 节点需要有不同的配置文件（即：不同的 my.cnf）；
2. Master 节点和 Salve 节点需要能够传输备份信息文件；
3. 在 Slave 节点第一次启动之前，需要执行一些初始化 SQL 操作；
```
由于 MySQL 本身同时拥有拓扑状态（主从节点的区别）和存储状态（MySQL 保存在本地的数据），自然要通过 StatefulSet 来解决这“三座大山”的问题。

##### “第一座大山：Master 节点和 Slave 节点需要有不同的配置文件”，处理：需要给主从节点分别准备两份不同的 MySQL 配置文件，然后根据 Pod 的序号（Index）挂载进去即可。

这里的配置文件信息，应该保存在 ConfigMap 里供 Pod使用。定义如下所示：
```yaml
apiVersion: v1
kind: ConfigMap
metadata:  
  name: mysql  
  labels:    
    app: mysql
data:  
  master.cnf: |    
    # 主节点 MySQL 的配置文件    
    [mysqld]    
    log-bin  
  slave.cnf: |    
    # 从节点 MySQL 的配置文件    
    [mysqld]    
    super-read-only
```

在这里，定义了 master.cnf 和 slave.cnf 两个 MySQL 的配置文件。
```text
master.cnf 开启了 log-bin，即：使用二进制日志文件的方式进行主从复制，这是一个标准的设置。

slave.cnf 的开启了 super-read-only，代表的是从节点会拒绝除了主节点的数据同步操作之外的所有写操作，即：它对用户是只读的。

上述 ConfigMap 定义里的 data 部分，是 Key-Value 格式的。
比如，master.cnf 就是这份配置数据的 Key，而“|”后面的内容，就是这份配置数据的 Value，这份数据将来挂载进 Master 节点对应的 Pod 后，就会在 Volume 目录里生成一个叫作 master.cnf 的文件。
```

接下来，需要创建两个 Service 来供 StatefulSet 以及用户使用，这两个 Service 的定义如下所示：
```yaml
apiVersion: v1
kind: Service
metadata:  
  name: mysql  
  labels:    
    app: mysql
spec:  
  ports:  
  - name: mysql    
    port: 3306  
  clusterIP: None  
  selector:    
   app: mysql
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-read  
  labels:    
    app: mysql
spec:  
  ports:  
  - name: mysql    
    port: 3306  
  selector:    
    app: mysql
```
可以看到，这两个 Service 都代理了所有携带 app=mysql 标签的 Pod，也就是所有的MySQL Pod。端口映射都是用 Service 的 3306 端口对应 Pod 的 3306 端口。

不同的是，第一个名叫“mysql”的 Service 是一个 Headless Service（即：clusterIP=None）。所以它的作用，是通过为 Pod 分配 DNS 记录来固定它的拓扑状态，比如“mysql-0.mysql”和“mysql-1.mysql”这样的 DNS 名字。其中，编号为 0 的节点就是我们的主节点。

而第二个名叫“mysql-read”的 Service，则是一个常规的 Service。

并且我们规定，所有用户的读请求，都必须访问第二个 Service 被自动分配的 DNS 记录，即：“mysql-read”（当然，也可以访问这个 Service 的 VIP）。这样，读请求就可以被转发到任意一个 MySQL 的主节点或者从节点上。
```text
备注：Kubernetes 中的所有 Service、Pod 对象，都会被自动分配同名的DNS 记录。
```

而所有用户的写请求，则必须直接以 DNS 记录的方式访问到 MySQL 的主节点，也就是：“mysql-0.mysql“这条 DNS 记录。

##### 第二座大山：Master 节点和 Salve 节点需要能够传输备份文件”的问题
翻越这座大山的思路，推荐的做法是：先搭建框架，再完善细节。其中，Pod 部分如何定义，是完善细节时的重点。

###### 先为 StatefulSet 对象规划一个大致的框架，如下所示:
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  serviceName: msyql
  replicas: 3
  template:
    metadata:
      labels:
        app: mysql
    spec:
      initContainers:
        - name: init-mysql
        - name: clone-mysql
      containers:
        - name: mysql
        - name: xtrabackup
      volumes:
        - name: conf
          emptyDir: {}
        - name: config-map
          configMap:
            name: mysql
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```
可以先为 StatefulSet 定义一些通用的字段。

比如：selector 表示，这个 StatefulSet 要管理的 Pod 必须携带 app=mysql 标签；它声明要使用的 Headless Service 的名字是：mysql。

这个 StatefulSet 的 replicas 值是 3，表示它定义的 MySQL 集群有三个节点：一个Master 节点，两个 Slave 节点。

StatefulSet 管理的“有状态应用”的多个实例，也都是通过同一份 Pod 模板创建出来的，使用的是同一个 Docker 镜像。这也就意味着：如果你的应用要求不同节点的镜像不一样，那就不能再使用 StatefulSet 了。对于这种情况，应该考虑Operator。

除了这些基本的字段外，作为一个有存储状态的 MySQL 集群，StatefulSet 还需要管理存储状态。所以，需要通过 volumeClaimTemplate（PVC 模板）来为每个 Pod 定义PVC。比如，这个 PVC 模板的 resources.requests.storage 指定了存储的大小为 10GiB；ReadWriteOnce 指定了该存储的属性为可读写，并且一个 PV 只允许挂载在一个宿主机上。将来，这个 PV 对应的的 Volume 就会充当 MySQL Pod 的存储数据目录。

###### 来重点设计一下这个 StatefulSet 的 Pod 模板，也就是 template 字段
由于 StatefulSet 管理的 Pod 都来自于同一个镜像，这就要求在编写 Pod 时，一定要保持清醒，用“人格分裂”的方式进行思考：
```text
1. 如果这个 Pod 是 Master 节点，我们要怎么做；
2. 如果这个 Pod 是 Slave 节点，我们又要怎么做。
```

第一步：从 ConfigMap 中，获取 MySQL 的 Pod 对应的配置文件。

需要进行一个初始化操作，根据节点的角色是 Master 还是 Slave 节点，为Pod 分配对应的配置文件。此外，MySQL 还要求集群里的每个节点都有一个唯一的 ID 文件，名叫 server-id.cnf。

这些初始化操作显然适合通过 InitContainer 来完成。所以，首先定义了一个 InitContainer，如下所示：
```yaml
#...
# template.spec
initContainers:
  - name: init-mysql
    image: mysql:5.7.32
    command:
    - bash 
    - "-c"
    - |
      set -ex
      # 从 Pod 的序号，生成 server-id
      [[ `hostname` =~ -([0-9]+)$ ]] || exit 1
      ordinal=${BASH_REMATCH[1]}
      echo [mysqld] > /mnt/conf.d/server-id.cnf
      # 由于 server-id=0 有特殊含义，我们给 ID 加一个 100 来避开它
      echo server-id=$((100 + $ordinal)) >> /mnt/conf.d/server-id.cnf
      # 如果 Pod 序号是 0，说明它是 Master 节点，从 ConfigMap 里把 Master 的配置文件拷贝
      # 否则，拷贝 Slave 的配置文件
      if [[ $ordinal -eq 0 ]]; then
        cp /mnt/config-map/master.cnf /mnt/conf.d/
      else 
        cp /mnt/config-map/slave.cnf /mnt/conf.d/
      fi
    volumeMounts:
    - name: conf
      mountPath: /mnt/conf.d
    - name: config-map
      mountPath: /mnt/config-map
```
在这个名叫 init-mysql 的 InitContainer 的配置中，它从 Pod 的 hostname 里，读取到了 Pod 的序号，以此作为 MySQL 节点的 server-id。

然后，init-mysql 通过这个序号，判断当前 Pod 到底是 Master 节点（即：序号为 0）还是 Slave 节点（即：序号不为 0），从而把对应的配置文件从 /mnt/config-map 目录拷贝到 /mnt/conf.d/ 目录下。

其中，文件拷贝的源目录 /mnt/config-map，正是 ConfigMap 在这个 Pod 的Volume，如下所示：
```yaml
# template.spec
volumes:
- name: conf
  emptyDir: {}
- name: config-map
  configMap:
    name: mysql
```

通过这个定义，init-mysql 在声明了挂载 config-map 这个 Volume 之后，ConfigMap里保存的内容，就会以文件的方式出现在它的 /mnt/config-map 目录当中。

而文件拷贝的目标目录，即容器里的 /mnt/conf.d/ 目录，对应的则是一个名叫 conf 的、emptyDir 类型的 Volume。基于 Pod Volume 共享的原理，当 InitContainer 复制完配置文件退出后，后面启动的 MySQL 容器只需要直接声明挂载这个名叫 conf 的 Volume，它所需要的.cnf 配置文件已经出现在里面了。

第二步：在 Slave Pod 启动前，从 Master 或者其他 Slave Pod 里拷贝数据库数据到自己的目录下。

为了实现这个操作，我们就需要再定义第二个 InitContainer，如下所示：
```yaml
#...
# template.spec.initContainers
- name: clone-mysql
  #image: gcr.io/google-samples/xtrabackup:1.0
  image: ipunktbs/xtrabackup:1.2.0
  command:
  - bash
  - "-c"
  - |
    set -ex
    # 拷贝操作只需要在第一次启动时进行，所以如果数据已经存在，跳过
    [[ -d /var/lib/mysql/mysql ]] && exit 0
    # Master 节点 (序号为 0) 不需要做这个操作
    [[ `hostname` =~ -([0-9]+)$ ]] || exit 1
    ordinal=${BASH_REMATCH[1]}
    [[ $ordinal -eq 0 ]] && exit 0
    # 使用 ncat 指令，远程地从前一个节点拷贝数据到本地
    ncat --recv-only mysql-$(($ordinal-1)).mysql 3307 | xbstream -x -C /var/lib/mysql
    # 执行 --prepare，这样拷贝来的数据就可以用作恢复了
    xtrabackup --prepare --target-dir=/var/lib/mysql
  volumeMounts:
  - name: data
    mountPath: /var/lib/mysql
    subPath: mysql
  - name: conf
    mountPath: /etc/mysql/conf.d
```
在这个名叫 clone-mysql 的 InitContainer 里，使用的是 xtrabackup 镜像（它里面安装了 xtrabackup 工具）。

而在启动命令里，首先做了一个判断。即：当初始化所需的数据（/var/lib/mysql/mysql 目录）已经存在，或者当前 Pod 是 Master 节点的时候，不需要做拷贝操作。

接下来，clone-mysql 会使用 Linux 自带的 ncat 指令，向 DNS 记录为“mysql-< 当前序号减一 >.mysql”的 Pod，也就是当前 Pod 的前一个 Pod，发起数据传输请求，并且直接用 xbstream 指令将收到的备份数据保存在 /var/lib/mysql 目录下。
```text
备注：3307 是一个特殊端口，运行着一个专门负责备份 MySQL 数据的辅助进程。
```

这一步可以随意选择用自己喜欢的方法来传输数据。比如，用 scp 或者 rsync，都没问题。

可能已经注意到，这个容器里的 /var/lib/mysql 目录，实际上正是一个名为 data 的PVC，即：在前面声明的持久化存储。

这就可以保证，哪怕宿主机宕机了，数据库的数据也不会丢失。更重要的是，由于Pod Volume 是被 Pod 里的容器共享的，所以后面启动的 MySQL 容器，就可以把这个Volume 挂载到自己的 /var/lib/mysql 目录下，直接使用里面的备份数据进行恢复操作。

不过，clone-mysql 容器还要对 /var/lib/mysql 目录，执行一句 xtrabackup --prepare操作，目的是让拷贝来的数据进入一致性状态，这样，这些数据才能被用作数据恢复。

至此，就通过 InitContainer 完成了对“主、从节点间备份文件传输”操作的处理过程，也就是翻越了“第二座大山”。

接下来，可以开始定义 MySQL 容器, 启动 MySQL 服务了。由于 StatefulSet 里的所有 Pod 都来自用同一个 Pod 模板，所以还要“人格分裂”地去思考：这个 MySQL 容器的启动命令，在 Master 和 Slave 两种情况下有什么不同。

有了 Docker 镜像，在 Pod 里声明一个 Master 角色的 MySQL 容器并不是什么困难的事情：直接执行 MySQL 启动命令即可。

但是，如果这个 Pod 是一个第一次启动的 Slave 节点，在执行 MySQL 启动命令之前，它就需要使用前面 InitContainer 拷贝来的备份数据进行初始化。

可是，别忘了，容器是一个单进程模型。

所以，一个 Slave 角色的 MySQL 容器启动之前，谁能负责给它执行初始化的 SQL 语句呢？

这就是需要解决的“第三座大山”的问题，即：如何在 Slave 节点的 MySQL 容器第一次启动之前，执行初始化 SQL。

你可能已经想到了，可以为这个 MySQL 容器额外定义一个 sidecar 容器，来完成这个操作，它的定义如下所示：
```yaml
# template.spec.containers
- name: xtrabackup
  #image: gcr.io/google-samples/xtrabackup:1.0
  image: ipunktbs/xtrabackup:1.2.0
  ports:
  - name: xtrabackup
    containerPort: 3307
  command:
  - bash
  - "-c"
  - |
    set -ex
    cd /var/lib/mysql
    # 从备份信息文件里读取 MASTER_LOG_FILEM 和 MASTER_LOG_POS 这两个字段的值，用来拼装
    if [[ -f xtrabackup_slave_info ]]; the
      # 如果 xtrabackup_slave_info 文件存在，说明这个备份数据来自于另一个 Slave 节点。
      mv xtrabackup_slave_info change_master_to.sql.in
      # 所以，也就用不着 xtrabackup_binlog_info 了
      rm -f xtrabackup_binlog_info
    elif [[ -f xtrabackup_binlog_info ]]; then
      # 如果只存在 xtrabackup_binlog_inf 文件，那说明备份来自于 Master 节点，我们就需备份
      [[ `cat xtrabackup_binlog_info` =~ ^(.*?)[[:space:]]+(.*?)$ ]] || exit 1
      rm xtrabackup_binlog_info
      # 把两个字段的值拼装成 SQL，写入 change_master_to.sql.in 文件
      echo "CHANGE MASTER TO MASTER_LOG_FILE='${BASH_REMATCH[1]}',\
        MASTER_LOG_POS=${BASH_REMATCH[2]}" > change_master_to.sql.in
    fi

    # 如果 change_master_to.sql.in，就意味着需要做集群初始化工作
    if [[ -f change_master_to.sql.in ]]; then
      # 但一定要先等 MySQL 容器启动之后才能进行下一步连接 MySQL 的操作
      echo "Waiting for mysqld to be ready (accepting connections)"
      until mysql -h 127.0.0.1 -e "SELECT 1"; do sleep 1; done
      
      echo "Initializing replication from clone position"
      # 将文件 change_master_to.sql.in 改个名字，防止这个 Container 重启的时候，因为又找到了 change_master_to.sql.in，从而重复执行一遍这个初始化流程
      mv change_master_to.sql.in change_master_to.sql.orig
      # 使用 change_master_to.sql.orig 的内容，也是就是前面拼装的 SQL，组成一个完整的初始化和启动 Slave 的 SQL 语句
      mysql -h 127.0.0.1 <<EOF
    $(<change_master_to.sql.orig),
        MASTER_HOST='mysql-0.mysql',
        MASTER_USER='root',
        MASTER_PASSWORD='',
        MASTER_CONNECT_RETRY=10;
    START SLAVE;
    EOF
    fi

    # 使用 ncat 监听 3307 端口。它的作用是，在收到传输请求的时候，直接执行 "xtrabackup --backup" 命令，备份 MySQL 的数据并发送给请求者
    exec ncat --listen --keep-open --send-only --max-conns=1 3307 -c \
            "xtrabackup --backup --slave-info --stream=xbstream --host=127.0.0.1 --user=root"
  volumeMounts:
  - name: data
    mountPath: /var/lib/mysql
    subPath: mysql
  - name: conf
    mountPath: /etc/mysql/conf.d
```
可以看到，在这个名叫 xtrabackup 的 sidecar 容器的启动命令里，其实实现了两部分工作。

第一部分工作，当然是 MySQL 节点的初始化工作。这个初始化需要使用的 SQL，是sidecar 容器拼装出来、保存在一个名为 change_master_to.sql.in 的文件里的，具体过程如下所示：
```text
sidecar 容器首先会判断当前 Pod 的 /var/lib/mysql 目录下，是否有xtrabackup_slave_info 这个备份信息文件。
    如果有，则说明这个目录下的备份数据是由一个 Slave 节点生成的。这种情况下，XtraBackup 工具在备份的时候，就已经在这个文件里自动生成了 "CHANGE MASTERTO" SQL 语句。所以，我们只需要把这个文件重命名为 change_master_to.sql.in，后面直接使用即可。
    如果没有 xtrabackup_slave_info 文件、但是存在 xtrabackup_binlog_info 文件，那就说明备份数据来自于 Master 节点。这种情况下，sidecar 容器就需要解析这个备份信息文件，读取 MASTER_LOG_FILE 和 MASTER_LOG_POS 这两个字段的值，用它们拼装出初始化 SQL 语句，然后把这句 SQL 写入到 change_master_to.sql.in 文件中。

接下来，sidecar 容器就可以执行初始化了。从上面的叙述中可以看到，只要这个change_master_to.sql.in 文件存在，那就说明接下来需要进行集群初始化操作。

所以，这时候，sidecar 容器只需要读取并执行 change_master_to.sql.in 里面的“CHANGE MASTER TO”指令，再执行一句 START SLAVE 命令，一个 Slave 节点就被成功启动了。
    需要注意的是：Pod 里的容器并没有先后顺序，所以在执行初始化 SQL 之前，必须先执行一句 SQL（select 1）来检查一下 MySQL 服务是否已经可用。
    
当然，上述这些初始化操作完成后，还要删除掉前面用到的这些备份信息文件。否则，下次这个容器重启时，就会发现这些文件存在，所以又会重新执行一次数据恢复和集群初始化的操作，这是不对的。

同理，change_master_to.sql.in 在使用后也要被重命名，以免容器重启时因为发现这个文件存在又执行一遍初始化。

在完成 MySQL 节点的初始化后，这个 sidecar 容器的第二个工作，则是启动一个数据传输服务。

具体做法是：sidecar 容器会使用 ncat 命令启动一个工作在 3307 端口上的网络发送服务。一旦收到数据传输请求时，sidecar 容器就会调用 xtrabackup --backup 指令备份当前 MySQL 的数据，然后把这些备份数据返回给请求者。这就是为什么我们在InitContainer 里定义数据拷贝的时候，访问的是“上一个 MySQL 节点”的 3307 端口。

值得一提的是，由于 sidecar 容器和 MySQL 容器同处于一个 Pod 里，所以它是直接通过Localhost 来访问和备份 MySQL 容器里的数据的，非常方便。

同样地，我在这里举例用的只是一种备份方法而已，你完全可以选择其他自己喜欢的方案。比如，你可以使用 innobackupex 命令做数据备份和准备，它的使用方法几乎与本文的备份方法一样。
```

扳倒了这“三座大山”后，终于可以定义 Pod 里的主角，MySQL 容器了。有了前面这些定义和初始化工作，MySQL 容器本身的定义就非常简单了，如下所示：
```yaml
#...
# template.spec
containers:
- name: mysql
  image: mysql:5.7.32
  env:
  - name: MYSQL_ALLOW_EMPTY_PASSWORD
    value: "1"
  ports:
  - name: mysql
    containerPort: 3306
  volumeMounts:
  - name: data
    mountPath: /var/lib/mysql
    subPath: mysql
  - name: conf
    mountPath: /etc/mysql/conf.d
  resources:
    requests:
      cpu: 500m
      memory: 1Gi
  livenessProbe:
    exec:
      command: ["mysqladmin", "ping"]
    initialDelaySeconds: 30
    periodSeconds: 10
    timeoutSeconds: 5
  readinessProbe:
    exec:
      # 通过 TCP 连接的方式进行健康检查
      command: ["mysql", "-h", "127.0.0.1", "-e", "SELECT 1"]
    initialDelaySeconds: 5
    periodSeconds: 2
    timeoutSeconds: 1
```
在这个容器的定义里，使用了一个标准的 MySQL 5.7.32 的官方镜像。它的数据目录是/var/lib/mysql，配置文件目录是 /etc/mysql/conf.d。

如果 MySQL 容器是 Slave 节点的话，它的数据目录里的数据，就来自于 InitContainer 从其他节点里拷贝而来的备份。它的配置文件目录/etc/mysql/conf.d 里的内容，则来自于 ConfigMap 对应的 Volume。而它的初始化工作，则是由同一个 Pod 里的 sidecar 容器完成的。

为它定义了一个 livenessProbe，通过 mysqladmin ping 命令来检查它是否健康；还定义了一个 readinessProbe，通过查询 SQL（select 1）来检查 MySQL 服务是否可用。当然，凡是 readinessProbe 检查失败的 MySQL Pod，都会从 Service 里被摘除掉。

##### 部署这个mysql集群
注意：用到了 StorageClass 来完成这个操作。它的作用，是自动地为集群里存在的每一个 PVC，调用存储插件（Rook）创建对应的 PV。

```text
$ kubectl create -f mysql-configmap.yaml
configmap/mysql created

$ kubectl create -f mysql-service.yaml
service/mysql created
service/mysql-read created

$ kubectl create -f mysql-statefulset.yaml
statefulset.apps/mysql created

$ kubectl get pod -l app=mysql
NAME                                      READY   STATUS    RESTARTS   AGE
mysql-0                                   2/2     Running   0          96s
mysql-1                                   2/2     Running   1          73s
mysql-2                                   2/2     Running   1          50s

尝试向这个 MySQL 集群发起请求，执行一些 SQL 操作来验证它是否正常：
$ kubectl run mysql-client --image=mysql:5.7 -i --rm --restart=Never -- \
   mysql -h mysql-0.mysql << EOF
CREATE DATABASE test;
CREATE TABLE test.messages (message VARCHAR(250));
INSERT INTO test.messages VALUES ('hello');
EOF
If you don't see a command prompt, try pressing enter.
pod "mysql-client" deleted

通过启动一个容器，使用 MySQL client 执行了创建数据库和表、以及插入数据的操作。需要注意的是，连接的 MySQL 的地址必须是 mysql-0.mysql（即：Master 节点的 DNS 记录）。因为，只有 Master 节点才能处理写操作。

通过连接 mysql-read 这个 Service，我们就可以用 SQL 进行读操作，如下所示：
$ kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-read -e "SELECT * FROM test.messages"
+---------+
| message |
+---------+
| hello   |
+---------+
pod "mysql-client" deleted

有了 StatefulSet 以后，你就可以像 Deployment 那样，非常方便地扩展这个 MySQL集群：
$ kubectl scale statefulset mysql  --replicas=5

直接连接 mysql-3.mysql，即 mysql-3 这个 Pod 的 DNS 名字来进行查询操作：
$ kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
    mysql -h mysql-3.mysql -e "SELECT * FROM test.messages"
  
```

### 对statefulSet进行滚动升级，只需要对statefulSet中的pod模版进行修改
```text
$ kubectl patch statefulset mysql --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"mysql:5.7.23"}]'
statefulset.apps/mysql patched

使用了 kubectl patch 命令，它的意思是，以“补丁”的方式（JSON 格式的）修改一个 API 对象的指定字段，也就是后面指定的“spec/template/spec/containers/0/image”。

这样，StatefulSet Controller 就会按照与 Pod 编号相反的顺序，从最后一个 Pod 开始，逐一更新这个 StatefulSet 管理的每个 Pod。而如果更新发生了错误，这次“滚动更新”就会停止。此外，StatefulSet 的“滚动更新”还允许进行更精细的控制，比如金丝雀发布（Canary Deploy）或者灰度发布，这意味着应用的多个实例中被指定的一部分不会被更新到最新的版本。

这个字段，正是 StatefulSet 的 spec.updateStrategy.rollingUpdate 的 partition 字段。

将前面这个 StatefulSet 的 partition 字段设置为 2：
$ kubectl patch statefulset mysql -p '{"spec":{"updateStrategy":{"type":"RollingUpdate","rollingUpdate":{"partition":2}}}}'
statefulset.apps/mysql patched
kubectl patch 命令后面的参数（JSON 格式的），就是 partition 字段在 API 对象里的路径。所以，上述操作等同于直接使用 kubectl edit 命令，打开这个对象，把partition 字段修改为 2。

指定了当 Pod 模板发生变化的时候，比如 MySQL 镜像更新到 5.7.23，那么只有序号大于或者等于 2 的 Pod 会被更新到这个版本。并且，如果你删除或者重启了序号小于 2 的 Pod，等它再次启动后，也会保持原先的 5.7.2 版本，绝不会被升级到 5.7.23 版本。

StatefulSet 可以说是 Kubernetes 项目中最为复杂的编排对象。
```

### 移除相关资源
```text
移除statefulSet编排相关信息
$ kubectl delete -f mysql-statefulset.yaml 

移除service
$ kubectl delete -f mysql-service.yaml 

移除configmap
$ kubectl delete -f mysql-configmap.yaml 

移除所有pvc
$ kubectl delete -f $(kubectl get pvc -l app=mysql | awk {'print$1'})
```

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
  
Kubernetes 项目引入了一组叫作 Persistent VolumeClaim（PVC）和 Persistent Volume（PV）的 API 对象，大大降低了用户声明和使用持久化 Volume 的门槛。
https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes

PVC 使用流程如下：

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

可以看到，在这个 PVC 对象里，不需要任何关于 Volume 细节的字段，只有描述性的属性和定义。
比如，storage: 1Gi，表示我想要的 Volume 大小至少是 1 GiB；
accessModes:ReadWriteOnce，表示这个 Volume 的挂载方式是可读写，并且只能被挂载在一个节点上而非被多个节点共享。

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
在这个 Pod 的 Volumes 定义中，只需要声明它的类型是persistentVolumeClaim，然后指定 PVC 的名字，而完全不必关心 Volume 本身的定义。

StatefulSet pvc : cicd/k8s/case/statefulset/basic-stateful-set.yaml
$ kubectl create -f basic-stateful-set.yaml.yaml 
$ kubectl get pvc -l app=nginx
NAME        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS         AGE
www-web-0   Bound    pvc-f942fb2e-f6d1-4268-a64a-bc7379e91a5b   100M       RWO            course-nfs-storage   13m
www-web-1   Bound    pvc-65b9309b-c153-4c04-918e-ec5c5e70aa3e   100M       RWO            course-nfs-storage   13m
这些 PVC，都以“<PVC 名字 >-<StatefulSet 名字 >-< 编号 >”的方式命名，并且处于 Bound 状态。  

可以使用如下所示的指令，在 Pod 的 Volume 目录里写入一个文件，来验证一下上述 Volume 的分配情况：
$ for i in 0 1; do kubectl exec web-$i -- sh -c 'echo hello $(hostname) > /usr/share/nginx/html/index.html' ; done

在这个 Pod 容器里访问“http://localhost”，你实际访问到的就是 Pod里 Nginx 服务器进程，而它会为你返回 /usr/share/nginx/html/index.html 里的内容，这个操作的执行方法如下所示
$ for i in 0 1; do kubectl exec -it web-$i -- curl localhost; done  
hello web-0
hello web-1

如果使用 kubectl delete 命令删除这两个 Pod，这些 Volume 里的文件会不会丢失呢？
$ kubectl delete pod -l app=nginx

再次访问nginx相关服务，发现该数据还是存在
$ for i in 0 1; do kubectl exec -it web-$i -- curl localhost; done
hello web-0
hello web-1

会发现，这个请求依然会返回：hello web-0。也就是说，原先与名叫 web-0 的 Pod 绑定的 PV，在这个 Pod 被重新创建之后，依然同新的名叫 web-0 的 Pod 绑定在了一起。

首先，StatefulSet 的控制器直接管理的是 Pod。这是因为，StatefulSet 里的不同 Pod实例，不再像 ReplicaSet 中那样都是完全一样的，而是有了细微区别的。比如，每个 Pod的 hostname、名字等都是不同的、携带了编号的。而 StatefulSet 区分这些实例的方式，就是通过在 Pod 的名字里加上事先约定好的编号。

其次，Kubernetes 通过 Headless Service，为这些有编号的 Pod，在 DNS 服务器中生成带有同样编号的 DNS 记录。只要 StatefulSet 能够保证这些 Pod 名字里的编号不变，那么 Service 里类似于 web-0.nginx.default.svc.cluster.local 这样的 DNS 记录也就不会变，而这条记录解析出来的 Pod 的 IP 地址，则会随着后端 Pod 的删除和再创建而自动更新。这当然是 Service 机制本身的能力，不需要 StatefulSet 操心。

最后，StatefulSet 还为每一个 Pod 分配并创建一个同样编号的 PVC。这样，Kubernetes 就可以通过 Persistent Volume 机制为这个 PVC 绑定上对应的 PV，从而保证了每一个 Pod 都拥有一个独立的 Volume。

在这种情况下，即使 Pod 被删除，它所对应的 PVC 和 PV 依然会保留下来。所以当这个Pod 被重新创建出来之后，Kubernetes 会为它找到同样编号的 PVC，挂载这个 PVC 对应的 Volume，从而获取到以前保存在 Volume 里的数据。

```
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
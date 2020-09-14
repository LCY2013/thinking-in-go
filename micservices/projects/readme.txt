Go-kit 是一套强大的微服务开发工具集，用于指导开发人员解决分布式系统开发过程中所遇到的问题，帮助开发人员更专注于业务开发。
Go-kit 推荐使用 transport、endpoint 和 service 3 层结构来组织项目，它们的作用分别为：
    1、transport 层，指定项目提供服务的方式，比如 HTTP 或者 gRPC 等 。
    2、endpoint 层，负责接收请求并返回响应。对于每一个服务接口，endpoint 层都使用一个抽象的 Endpoint 来表示 ，
        我们可以为每一个 Endpoint 装饰 Go-kit 提供的附加功能，如日志记录、限流、熔断等。
    3、service 层，提供具体的业务实现接口，endpoint 层中的 Endpoint 通过调用 service 层的接口方法处理请求。


DDD:
    1、限界上下文主要有以下几种映射方式。
        合作关系（Partnership）：两个上下文紧密合作的关系，一荣俱荣，一损俱损。

        共享内核（Shared Kernel）：两个上下文依赖部分共享的模型。

        防腐层（Anticorruption Layer）：一个上下文通过一些适配和转换与另一个上下文交互。

        客户方-供应方开发（Customer-Supplier Development）：上下文之间有组织的上下游依赖。

        开放主机服务（Open Host Service）：定义一种协议来让其他上下文对本上下文进行访问。

        遵奉者（Conformist）：下游上下文只能盲目依赖上游上下文。

        发布语言（Published Language）：通常与 OHS 一起使用，用于定义开放主机的协议。

        大泥球（Big Ball of Mud）：混杂在一起的上下文关系，边界不清晰。

        另谋他路（Separate Way）：两个完全没有任何联系的上下文。
    实践中，主要通过防腐层映射不同的限界上下文中相同的领域对象，保证整体领域概念的完整和统一。

Go 本身提供了一套轻量级的测试框架，用于对 Go 程序进行单元测试和基准测试。
    go test 命令是一个按照一定的约定和组织来测试代码的程序，它执行的文件都是以“_test.go” 作为后缀，这部分文件不会包含在 go build 的代码构建中。
    在测试文件中主要存在以下三种函数类型：
        1、以 Test 作为函数名前缀的测试函数，一般用作单元测试，测试函数的逻辑行为是否正确；
        2、以 Benchmark 作为函数名前缀的基准测试函数，一般用来衡量函数的性能；
        3、以 Example 作为函数名前缀的示例函数，主要用于提供示例文档。

通过pipeline部署user服务
        user-redis-deployment.yaml 文件通过 Deployment Controller 管理 Pod，当 Controller 中的 Pod 出现异常被重启时，很可能导致 Pod 的 IP 发生变化。
    如果此时 user 服务通过固定 IP 的方式访问 Redis，很可能会访问失败。为了避免这种情况，我们可以为 user-redis Pod 定义一个 Service(user-redis-service.yaml)。
    Service 定义了一组 Pod 的逻辑集合和一个用于访问它们的策略，Kubernetes 集群会为 Service 分配一个固定的 Cluster IP，用于集群内部的访问。

    通过 Cluster IP 访问 MySQL 和 Redis 等服务，我们就无须担心 Pod IP 的变化。

    通过 Pipeline 部署服务到 Kubernetes 集群，主要有以下步骤：
        1、从 GitHub 中拉取代码。
        2、构建 Docker 镜像。
        3、上传 Docker 镜像到 Docker Hub。
        4、将应用部署 Kubernetes。
        5、接口测试。
    在 Pipeline 中，将上述步骤组织成相应的 Stage，让 Jenkins 为我们完成服务的持续集成和自动化测试，接下来以 user 服务的部署作为例子。
    Pipeline 脚本是由 Groovy 语言实现，支持 Declarative（声明式）和 Scripted（脚本式）语法，下面的演示就基于脚本式语法进行介绍。
    1、拉取代码。
     stage声明如下:
      stage('clone code from github') {
          echo "first stage: clone code"
          git url: "https://github.com/LCY2013/thinking-in-go.git"
          script {
              commit_id = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
          }
      }
      通过 git url 命令从 GitHub 中获取 user 服务的代码，并将本次提交记录的 commit_id 提取出来作为变量使用。
    2、使用 user 服务中的 Dockfile 定义构建相应的 user 镜像。
      Stage 声明如下:
       stage('build image') {
           echo "second stage: build docker image"
           sh "docker build -t luochunyun/user:${commit_id} micservices/projects/user"
       }
      为了方便在排查问题时可以根据对应的代码记录定位代码，采用了 GitHub 的提交记录 commit_id 作为镜像的 tag。
      同时为了将 MySQL 和 Redis 的地址作为参数传入，修改 user 服务的 Dockerfile 为：micservices/projects/部署pipeline服务/Dockerfile。
      mysqlAddr 和 redisAddr 将在 user.yaml 配置文件中以环境变量的方式指定 MySQL 和 Redis 的地址。
    3、为了方便 Kubernetes 拉取服务的镜像，我们将第二步构建好的 Docker 镜像推送到镜像仓库中。
       Stage 声明如下：
        stage('push image') {
            echo "third stage: push docker image to registry"
            sh "docker login -u user -p password"
            sh "docker push luochunyun/user:${commit_id}"
        }
        Docker 中默认的镜像仓库为 Docker Hub，上述声明中就将 user 镜像推送到 Docker Hub 中，当然你也可以选择将镜像推送到私有仓库中。
        往 Docker Hub 中推送镜像需要提交账号密码，这需要我们预先注册申请一个 Docker Hub 账户。
    4、我们使用 kubectl 将 user 服务部署到 Kubernetes 中。
        为了保证部署到正确版本的镜像，我们需要将 commit_id 替换到 user.yaml 文件中，以及将 mysqlAddr 和 redisAddr 作为环境变量输入，
        user-deployment.yaml 的配置为: micservices/projects/部署pipeline服务/user-deployment.yaml
       Stage 声明如下:
        stage('deploy to Kubernetes') {
            echo "forth stage: deploy to Kubernetes"
            sh "sed -i 's/<COMMIT_ID_TAG>/${commit_id}/' user-deployment.yaml"
            sh "sed -i 's/<MYSQL_ADDR_TAG>/${mysql_addr}/' user-deployment.yaml"
            sh "sed -i 's/<REDIS_ADDR_TAG>/${redis_addr}/' user-deployment.yaml"
            sh "kubectl apply -f user-deployment.yaml"
        }
        在上述声明中，首先使用 sed 命令将 yaml 文件中标识替换为对应的变量，再通过 kubectl apply 命令重新部署了 user-service Pod。
        user-service.yaml 配置为: micservices/projects/部署pipeline服务/user-service.yaml
        指定的 Service 的类型为 NodePort，并将 user 服务的接口通过 Node 节点的 19527 暴露出去，对此，我们就可以在集群外部通过 NodeIP:NodePort 的方式访问 user 服务了。
    5、我们通过 go test 对 user 中的 HTTP 接口进行接口测试，验证代码集成的效果
       Stage声明如下:
        stage('http test') {
            echo "fifth stage: http test"
            sh "cd micservices/projects/user/transport && go test -args ${user_addr}"
        }

    在 Jenkins 中创建一个 Pipeline 任务，将上述脚本复制到 Script 区域中，保存后触发构建，不过在这之前需要在 Jenkins 中安装和配置好 Kubernetes Plugin 和 Docker Plugin。
    在实际的开发中，我们可以将上述 Pipeline 脚本放入到 Jenkinsfile 中，与代码一同提交到代码库，将 Pipeline 任务的脚本配置类型修改为 Pipeline Script from SCM，引用代码库中 Pipeline 脚本进行构建。





















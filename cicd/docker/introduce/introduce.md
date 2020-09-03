## docker 三个核心概念
####  1、镜像(Images)
####  2、容器(Container)
####  3、仓库(Repository)
```
    docker 将 配置文件、依赖、环境变量打包成镜像，为应用程序提供运行环境。
镜像是由多个镜像叠加起来的文件系统，底层为UnionFS和AUFS联合文件系统，这是
一种支持分层且可以叠加的轻量级且高性能的文件系统。镜像只是一个可读模版，我们
可以在这个只读模版上叠加其他不同的模版组成不同的镜像。
    容器通过镜像启动，是镜像的运行实例，一个镜像也可以运行多个容器。容器间
是互相隔离的。容器可以被创建、运行、停止、删除、暂停、重启等操作。这个就
类似于容器给我们提供了一个类Linux的沙箱环境。
    仓库是用来存储和管理镜像的地方，分为公有仓库和私有仓库。目前docker仓库
默认的公有仓库是DockerHub，国内可以使用阿里云的仓库地址或者网易云仓库地址:
http://hub-mirror.c.163.com。
    docker镜像可以被运行在任何装有Docker环境的OS上面(Build onece,Run anywhere)。
```
#### Dockerfile基本概念
```
From：Dockerfile 中必须出现的第一个指令，用于指定基础镜像，例如:指定了基础镜像为 golang:latest
WORKDIR：指定工作目录，之后的命令将会在该目录下执行。
COPY：将本地文件添加到容器指定位置中。
RUN：创建镜像执行的命令，一个 Dockerfile 可以有多个 RUN 命令。在上述 RUN 指令中我们指定了 Go 的代理，并通过 go build 构建了 user 服务。
ENTRYPOINT：容器被启动后执行的命令，每个 Dockerfile 只有一个。我们通过该命令在容器启动后，又启动了 user 服务。
ENV：声明环境变量信息。
CMD：容器被启动后执行的命令，但是有个bug，建议使用ENTRYPOINT
```
#### Docker 基础命令
```
    Usage:	docker [OPTIONS] COMMAND
    
    A self-sufficient runtime for containers
    
    Options:
          --config string      Location of client config files (default
                               "/Users/magicLuoMacBook/.docker")
      -c, --context string     Name of the context to use to connect to the
                               daemon (overrides DOCKER_HOST env var and
                               default context set with "docker context use")
      -D, --debug              Enable debug mode
      -H, --host list          Daemon socket(s) to connect to
      -l, --log-level string   Set the logging level
                               ("debug"|"info"|"warn"|"error"|"fatal")
                               (default "info")
          --tls                Use TLS; implied by --tlsverify
          --tlscacert string   Trust certs signed only by this CA (default
                               "/Users/magicLuoMacBook/.docker/ca.pem")
          --tlscert string     Path to TLS certificate file (default
                               "/Users/magicLuoMacBook/.docker/cert.pem")
          --tlskey string      Path to TLS key file (default
                               "/Users/magicLuoMacBook/.docker/key.pem")
          --tlsverify          Use TLS and verify the remote
      -v, --version            Print version information and quit
    
    Management Commands:
      builder     Manage builds
      config      Manage Docker configs
      container   Manage containers
      context     Manage contexts
      image       Manage images
      network     Manage networks
      node        Manage Swarm nodes
      plugin      Manage plugins
      secret      Manage Docker secrets
      service     Manage services
      stack       Manage Docker stacks
      swarm       Manage Swarm
      system      Manage Docker
      trust       Manage trust on Docker images
      volume      Manage volumes
    
    Commands:
      attach      Attach local standard input, output, and error streams to a running container
      build       Build an image from a Dockerfile
      commit      Create a new image from a container's changes
      cp          Copy files/folders between a container and the local filesystem
      create      Create a new container
      diff        Inspect changes to files or directories on a container's filesystem
      events      Get real time events from the server
      exec        Run a command in a running container
      export      Export a container's filesystem as a tar archive
      history     Show the history of an image
      images      List images
      import      Import the contents from a tarball to create a filesystem image
      info        Display system-wide information
      inspect     Return low-level information on Docker objects
      kill        Kill one or more running containers
      load        Load an image from a tar archive or STDIN
      login       Log in to a Docker registry
      logout      Log out from a Docker registry
      logs        Fetch the logs of a container
      pause       Pause all processes within one or more containers
      port        List port mappings or a specific mapping for the container
      ps          List containers
      pull        Pull an image or a repository from a registry
      push        Push an image or a repository to a registry
      rename      Rename a container
      restart     Restart one or more containers
      rm          Remove one or more containers
      rmi         Remove one or more images
      run         Run a command in a new container
      save        Save one or more images to a tar archive (streamed to STDOUT by default)
      search      Search the Docker Hub for images
      start       Start one or more stopped containers
      stats       Display a live stream of container(s) resource usage statistics
      stop        Stop one or more running containers
      tag         Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE
      top         Display the running processes of a container
      unpause     Unpause all processes within one or more containers
      update      Update configuration of one or more containers
      version     Show the Docker version information
      wait        Block until one or more containers stop, then print their exit codes
```



















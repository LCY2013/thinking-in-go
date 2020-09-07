执行步骤:
1、编写Dockerfile

2、构建镜像
-t --tag
$ Docker build -t hello-py .
Dockerfile每执行一个原语，都会生成一个layer(层)
那么我们刚刚构建的镜像就一共有7层,如下:
"RootFS": {
    "Type": "layers",
    "Layers": [
        "sha256:b60e5c3bcef2f42ec42648b3acf7baf6de1fa780ca16d9180f3b4a3f266fe7bc",
        "sha256:568944187d9378b07cf2e2432115605b71c36ef566ec77fbf04516aab0bcdf8e",
        "sha256:7ea2b60b0a086d9faf2ba0a52d4e2f940d9361ed4179642686d1d8b59460667c",
        "sha256:7a287aad297b39792ee705ad5ded9ba839ee3f804fa3fb0b81bb8eb9f9acbf88",
        "sha256:28ebef53a6e9083c68acf9a7ae720ac92dd35abe9cf7f1e630976f7e3e93d9de",
        "sha256:7416a1f292e51722fac211f27e1326245cd570287d6d249823ef210ccd8a3e3a",
        "sha256:ec52c6b515d9336351930b68146ea7e03468394e4354a7e005639bd14d0c763c"
    ]
}

3、运行刚刚构建镜像
$ docker run -p 9527:80 hello-py
访问本机地址http://127.0.0.1:9527

4、查看刚刚容器的详细信息
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                  NAMES
46e3da051f1e        hello-py            "python app.py"     8 seconds ago       Up 7 seconds        0.0.0.0:9527->80/tcp   ecstatic_hellman
$ docker inspect 46e3
查看该容器的详细信息
...
"Networks": {
    "bridge": {
        "IPAMConfig": null,
        "Links": null,
        "Aliases": null,
        "NetworkID": "f4775b0ff3ba0b28d2ff25370536fd5eb877ad0e2c1fe973b532a7d862e6219c",
        "EndpointID": "06011e9ef6032ef5640ee9acd05fa8ea687a427c104e909492beb256a6d07ad7",
        "Gateway": "172.17.0.1",
        "IPAddress": "172.17.0.2",
        "IPPrefixLen": 16,
        "IPv6Gateway": "",
        "GlobalIPv6Address": "",
        "GlobalIPv6PrefixLen": 0,
        "MacAddress": "02:42:ac:11:00:02",
        "DriverOpts": null
    }
}
...

5、给该镜像构建一个tag
$ docker tag hello-py luochunyun/hello-py:1.0.0
$ docker images
REPOSITORY                                      TAG                 IMAGE ID            CREATED             SIZE
hello-py                                        latest              6a99b9d0894d        9 minutes ago       158MB
luochunyun/hello-py                             1.0.0               6a99b9d0894d        9 minutes ago       158MB

6、推送本地镜像到镜像仓库
默认的推送格式是:register-url/repository-id/image:tag
register-url:docker hub 官方地址
repository-id:自定义仓储
image:镜像名称
tag:版本号信息
luochunyun 是我的仓库名称
$ docker push luochunyun/hello-py:1.0.0

7、移除镜像
$ docker rmi -f hello-py:1.0.0

8、利用运行的容器提交保存镜像
docker commit -- 该命令用于创建一个运行中的容器状态的镜像
$ docker run -it busybox /bin/sh
创建一个新的目录/test
/ # mkdir /test
进入test目录
/ # cd test/
在test目录中创建一个name.txt的文件
/test # echo "lcy" > name.txt
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED              STATUS              PORTS               NAMES
b1680268af37        busybox             "/bin/sh"           About a minute ago   Up About a minute                       reverent_greider
$ docker commit b1680268af37 busybox-container:1.0.0
$ docker images
REPOSITORY                                      TAG                 IMAGE ID            CREATED             SIZE
busybox-container                               1.0.0               2496dbfd24db        8 seconds ago       1.22MB
$ docker run -it busybox-container:1.0.0 /bin/sh
可以看到这里的容器就是刚刚创建的目录内容一致
/ # ls /test/
name.txt
$ docker exec -it container-id /bin/sh
获取容器的进程ID号
$ docker inspect --format '{{ .State.Pid }}' container-id
$ docker run -it busybox /bin/sh
$ docker ps                                                                                                                                                                                                        1
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
375ad06ac0ad        busybox             "/bin/sh"           25 seconds ago      Up 24 seconds                           objective_bell
$ docker inspect --format '{{ .State.Pid }}' 375ad06ac0ad
1911
$ ls -l  /proc/1911/ns
总用量 0
lrwxrwxrwx. 1 root root 0 9月   7 01:11 cgroup -> cgroup:[4026531835]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 ipc -> ipc:[4026532199]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 mnt -> mnt:[4026532197]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 net -> net:[4026532202]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 pid -> pid:[4026532200]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 pid_for_children -> pid:[4026532200]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 time -> time:[4026531834]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 time_for_children -> time:[4026531834]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 user -> user:[4026531837]
lrwxrwxrwx. 1 root root 0 9月   7 01:11 uts -> uts:[4026532198]
/proc/[进程号]/ns
这也就意味着：一个进程，可以选择加入到某个进程已有的 Namespace 当中，从而达到"进入"这个进程所在容器的目的，这正是 docker exec 的实现原理。
而这个操作所依赖的，乃是一个名叫 setns() 的 Linux 系统调用。它的调用方法，用< setns.c >来演示这个问题。
它一共接收两个参数，第一个参数是 argv[1]，即当前进程要加入的 Namespace 文件的路径，比如 /proc/1911/ns/net；
而第二个参数，则是你要在这个 Namespace 里运行的进程，比如 /bin/bash。

下面开始编译这个文件
$ gcc -o set_ns setns.c
$ ./set_ns /proc/1911/ns/net /bin/bash
$ ifconfig
$ ip addr
发现网卡少了几个,这里就和容器看到的网卡信息一样了
可以在宿主机上查询/bin/bash信息
$ ps -aux | grep /bin/bash
root      2431  0.0  0.1 115560  3492 pts/0    S+   01:25   0:00 /bin/bash
$ ll /proc/2431/ns/net
lrwxrwxrwx. 1 root root 0 9月   7 01:27 /proc/2431/ns/net -> net:[4026532202]
$ ll /proc/1911/ns/net
lrwxrwxrwx. 1 root root 0 9月   7 01:22 /proc/1911/ns/net -> net:[4026532202]

可以看出宿主机上的bash和容器中的bash共同指向了同一个network namespace
docker 也提供net参数让我们的容器同用一个network namespace，如下:
$ docker run -it --net container:375ad06ac0ad busybox ifconfig
这个容器启动的时候就会使用容器ID为375ad06ac0ad的容器的network namespace，所以可以做到容器的网络互通
如果这里的net换成host(--net=host)，那么就不会使用network namespace隔离，直接使用宿主机的network



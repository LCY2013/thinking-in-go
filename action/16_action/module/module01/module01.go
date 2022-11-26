package module01

/*
1、导入本地 module

Go Module 从 Go 1.11 版本开始引入到 Go 中，现在它已经成为了 Go 语言的依赖管理与构建的标准。
因此，建议你彻底抛弃 Gopath 构建模式，全面拥抱 Go Module 构建模式。

当项目依赖已发布在 GitHub 等代码托管站点的公共 Go Module 时，Go 命令工具可以很好地完成依赖版本选择以及 Go Module 拉取的工作。

不过，如果项目依赖的是本地正在开发、尚未发布到公共站点上的 Go Module，那么应该如何做呢？

假设一个项目，这个项目中的 module a 依赖 module b，而 module b 是你另外一个项目中的 module，它本来是要发布到github.com/user/b上的。

但此时此刻，module b 还没有发布到公共托管站点上，它源码还在你的开发机器上。
也就是说，go 命令无法在github.com/user/b上找到并拉取 module a 的依赖 module b，这时，如果你针对 module a 所在项目使用 go mod tidy 命令，就会收到类似下面这样的报错信息：
$go mod tidy
go: finding module for package github.com/user/b
github.com/user/a imports
    github.com/user/b: cannot find module providing package github.com/user/b: module github.com/user/b: reading https://goproxy.io/github.com/user/b/@v/list: 404 Not Found
    server response:
    not found: github.com/user/b@latest: terminal prompts disabled
    Confirm the import path was entered correctly.
    If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.

这个时候，就可以借助 go.mod 的 replace 指示符，来解决这个问题。解决的步骤是这样的：

首先，需要在 module a 的 go.mod 中的 require 块中，手工加上这一条（这也可以通过 go mod edit 命令实现）：
require github.com/user/b v1.0.0

注意了，这里的 v1.0.0 版本号是一个“假版本号”，目的是满足 go.mod 中 require 块的语法要求。

然后，再在 module a 的 go.mod 中使用 replace，将上面对 module b v1.0.0 的依赖，替换为本地路径上的 module b:
replace github.com/user/b v1.0.0 => module b的本地源码路径

这样修改之后，go 命令就会让 module a 依赖你本地正在开发、尚未发布到代码托管网站的 module b 的源码了。

而且，如果 module b 已经提交到类 GitHub 的站点上，但 module b 的作者正在本地开发新版本，那么上面这种方法，也同样适合 module b 的作者在本地测试验证 module b 的最新版本源码。

虽然“伪造”go.mod 文件内容，可以解决上述这两个场景中的问题，但显然这种方法也是有“瑕疵”的。

首先，这个方法中，require 指示符将github.com/user/b v1.0.0替换为一个本地路径下的 module b 的源码版本，但这个本地路径是因开发者环境而异的。

go.mod 文件通常是要上传到代码服务器上的，这就意味着，另外一个开发人员下载了这份代码后，极大可能是无法成功编译的，想完成 module a 的编译，就得将 replace 后面的本地路径改为适配自己环境下的路径。

于是，每当开发人员 pull 代码后，第一件事就是要修改 module a 的 go.mod 中的 replace 块，每次上传代码前，可能也要将 replace 路径复原，这是一个很糟心的事情。
但即便如此，目前 Go 版本（最新为 Go 1.17.x）也没有一个完美的应对方案。

针对这个问题，Go 核心团队在 Go 社区的帮助下，在预计 2022 年 2 月发布的 go 1.19 版本中加入了 Go 工作区（Go workspace，也译作 Go 工作空间）辅助构建机制。

基于这个机制，可以将多个本地路径放入同一个 workspace 中，这样，在这个 workspace 下各个 module 的构建将优先使用 workspace 下的 module 的源码。
工作区配置数据会放在一个名为 go.work 的文件中，这个文件是开发者环境相关的，因此并不需要提交到源码服务器上，这就解决了上面“伪造 go.mod”方案带来的那些问题。

2、拉取私有 module 的需求与参考方案
Go 1.11 版本引入 Go Module 构建模式后，用 Go 命令拉取项目依赖的公共 Go Module，已不再是“痛点”，只需要在每个开发机上为环境变量 GOPROXY，配置一个高效好用的公共 GOPROXY 服务，就可以轻松拉取所有公共 Go Module 了：
action/16_action/module/module01/goModule机制.png

但随着公司内 Go 使用者和 Go 项目的增多，“重造轮子”的问题就出现了。
抽取公共代码放入一个独立的、可被复用的内部私有仓库成为了必然，这样就有了拉取私有 Go Module 的需求。

一些公司或组织的所有代码，都放在公共 vcs 托管服务商那里（比如 github.com），私有 Go Module 则直接放在对应的公共 vcs 服务的 private repository（私有仓库）中。
如果公司也是这样，那么拉取托管在公共 vcs 私有仓库中的私有 Go Module，也很容易，见下图：
action/16_action/module/module01/privateModule机制.png

也就是说，只要在每个开发机上，配置公共 GOPROXY 服务拉取公共 Go Module，同时再把私有仓库配置到 GOPRIVATE 环境变量，就可以了。
这样，所有私有 module 的拉取，都会直连代码托管服务器，不会走 GOPROXY 代理服务，也不会去 GOSUMDB 服务器做 Go 包的 hash 值校验。

当然，这个方案有一个前提，那就是每个开发人员都需要具有访问公共 vcs 服务上的私有 Go Module 仓库的权限，凭证的形式不限，可以是 basic auth 的 user 和 password，也可以是 personal access token（类似 GitHub 那种），只要按照公共 vcs 的身份认证要求提供就可以了。

不过，更多的公司 / 组织，可能会将私有 Go Module 放在公司 / 组织内部的 vcs（代码版本控制）服务器上，就像下面图中所示：
action/16_action/module/module01/privateModule.png

那么这种情况，该如何让 Go 命令，自动拉取内部服务器上的私有 Go Module 呢？这里给出两个参考方案。
1）第一个方案是通过直连组织公司内部的私有 Go Module 服务器拉取。
action/16_action/module/module01/直连组织公司内部的私有GoModule服务器.png

公司内部会搭建一个内部 goproxy 服务（也就是上图中的 in-house goproxy）。
这样做有两个目的：
一是为那些无法直接访问外网的开发机器，以及 ci 机器提供拉取外部 Go Module 的途径。
二来，由于 in-house goproxy 的 cache 的存在，这样做还可以加速公共 Go Module 的拉取效率。

2）第二种方案是将外部 Go Module 与私有 Go Module 都交给内部统一的 GOPROXY 服务去处理：
action/16_action/module/module01/内部统一的GOPROXY服务.png

在这种方案中，开发者只需要把 GOPROXY 配置为 in-house goproxy，就可以统一拉取外部 Go Module 与私有 Go Module。

但由于 go 命令默认会对所有通过 goproxy 拉取的 Go Module，进行 sum 校验（默认到 sum.golang.org)，而私有 Go Module 在公共 sum 验证 server 中又没有数据记录。
因此，开发者需要将私有 Go Module 填到 GONOSUMDB 环境变量中，这样，go 命令就不会对其进行 sum 校验了。

不过这种方案有一处要注意：in-house goproxy 需要拥有对所有 private module 所在 repo 的访问权限，才能保证每个私有 Go Module 都拉取成功。

可以对比一下上面这两个参考方案，看看你更倾向于哪一个，推荐第二个方案。在第二个方案中，可以将所有复杂性都交给 in-house goproxy 这个节点，开发人员可以无差别地拉取公共 module 与私有 module，心智负担降到最低。

3、统一 Goproxy 方案的实现思路与步骤
后续的方案实现准备一个示例环境，拓扑如下图：
action/16_action/module/module01/统一Goproxy方案的实现思路与步骤.png

1）选择一个 GOPROXY 实现
Go module proxy 协议规范发布后，Go 社区出现了很多成熟的 Goproxy 开源实现，比如有最初的athens，还有国内的两个优秀的开源实现：goproxy.cn和goproxy.io等。
其中，goproxy.io 在官方站点给出了企业内部部署的方法，所以就基于 goproxy.io 来实现我们的方案。

在上图中的 in-house goproxy 节点上执行这几个步骤安装 goproxy：
$mkdir ~/.bin/goproxy
$cd ~/.bin/goproxy
$git clone https://github.com/goproxyio/goproxy.git
$cd goproxy
$make

编译后，我们会在当前的 bin 目录（~/.bin/goproxy/goproxy/bin）下看到名为 goproxy 的可执行文件。

然后，建立 goproxy cache 目录：
$mkdir /root/.bin/goproxy/goproxy/bin/cache

再启动 goproxy：
$./goproxy -listen=0.0.0.0:8081 -cacheDir=/root/.bin/goproxy/goproxy/bin/cache -proxy https://goproxy.io
goproxy.io: ProxyHost https://goproxy.io

启动后，goproxy 会在 8081 端口上监听（即便不指定，goproxy 的默认端口也是 8081），指定的上游 goproxy 服务为 goproxy.io。

不过要注意下：goproxy 的这个启动参数并不是最终版本的，这里仅仅想验证一下 goproxy 是否能按预期工作。现在就来实际验证一下。

首先，在开发机上配置 GOPROXY 环境变量指向 10.10.20.20:8081：
// .bashrc
export GOPROXY=http://10.10.20.20:8081

生效环境变量后，执行下面命令：
$go get github.com/pkg/errors

结果和预期的一致，开发机顺利下载了 github.com/pkg/errors 包。可以在 goproxy 侧，看到了相应的日志：
goproxy.io: ------ --- /github.com/pkg/@v/list [proxy]
goproxy.io: ------ --- /github.com/pkg/errors/@v/list [proxy]
goproxy.io: ------ --- /github.com/@v/list [proxy]
goproxy.io: 0.146s 404 /github.com/@v/list
goproxy.io: 0.156s 404 /github.com/pkg/@v/list
goproxy.io: 0.157s  /github.com/pkg/errors/@v/list

在 goproxy 的 cache 目录下，也看到了下载并缓存的 github.com/pkg/errors 包：
$cd /root/.bin/goproxy/goproxy/bin/cache
$tree
.
└── pkg
    └── mod
        └── cache
            └── download
                └── github.com
                    └── pkg
                        └── errors
                            └── @v
                                └── list
8 directories, 1 file

2）自定义包导入路径并将其映射到内部的 vcs 仓库
一般公司可能没有为 vcs 服务器分配域名，也不能在 Go 私有包的导入路径中放入 ip 地址，因此需要给私有 Go Module 自定义一个路径，比如：mycompany.com/go/module1。
统一将私有 Go Module 放在 mycompany.com/go 下面的代码仓库中。

那么，接下来的问题就是，当 goproxy 去拉取 mycompany.com/go/module1 时，应该得到 mycompany.com/go/module1 对应的内部 vcs 上 module1 仓库的地址，这样，goproxy 才能从内部 vcs 代码服务器上下载 module1 对应的代码，具体的过程如下：
action/16_action/module/module01/自定义包导入路径并将其映射到内部的vcs仓库.png

那么如何实现为私有 module 自定义包导入路径，并将它映射到内部的 vcs 仓库呢？

其实方案不止一种，这里使用了 Google 云开源的一个名为govanityurls的工具，来为私有 module 自定义包导入路径。
然后，结合 govanityurls 和 nginx，就可以将私有 Go Module 的导入路径映射为其在 vcs 上的代码仓库的真实地址。
具体原理可以看一下这张图：
action/16_action/module/module01/vcs上的代码仓库的真实地址.png

首先，goproxy 要想不把收到的拉取私有 Go Module（mycompany.com/go/module1）的请求转发给公共代理，需要在其启动参数上做一些手脚，比如下面这个就是修改后的 goproxy 启动命令：
$./goproxy -listen=0.0.0.0:8081 -cacheDir=/root/.bin/goproxy/goproxy/bin/cache -proxy https://goproxy.io -exclude "mycompany.com/go"

这样，凡是与 -exclude 后面的值匹配的 Go Module 拉取请求，goproxy 都不会转给 goproxy.io，而是直接请求 Go Module 的“源站”。

而上面这张图中要做的，就是将这个“源站”的地址，转换为企业内部 vcs 服务中的一个仓库地址。
然后假设 mycompany.com 这个域名并不存在（很多小公司没有内部域名解析能力），从图中可以看到，我们会在 goproxy 所在节点的 /etc/hosts 中加上这样一条记录：
127.0.0.1 mycompany.com

这样做了后，goproxy 发出的到 mycompany.com 的请求实际上是发向了本机。
而上面这图中显示，监听本机 80 端口的正是 nginx，nginx 关于 mycompany.com 这一主机的配置如下：
// /etc/nginx/conf.d/gomodule.conf
server {
        listen 80;
        server_name mycompany.com;
        location /go {
                proxy_pass http://127.0.0.1:8080;
                proxy_redirect off;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "upgrade";
        }
}

对于路径为 mycompany.com/go/xxx 的请求，nginx 将请求转发给了 127.0.0.1:8080，而这个服务地址恰恰就是 govanityurls 工具监听的地址。

govanityurls 这个工具，是前 Go 核心开发团队成员Jaana B.Dogan开源的一个工具，这个工具可以帮助 Gopher 快速实现自定义 Go 包的 go get 导入路径。

govanityurls 本身，就好比一个“导航”服务器。
当 go 命令向自定义包地址发起请求时，实际上是将请求发送给了 govanityurls 服务，之后，govanityurls 会将请求中的包所在仓库的真实地址（从 vanity.yaml 配置文件中读取）返回给 go 命令，后续 go 命令再从真实的仓库地址获取包数据。

注：govanityurls 的安装方法很简单，直接 go install/go get github.com/GoogleCloudPlatform/govanityurls 就可以了。

在示例中，vanity.yaml 的配置如下：
host: mycompany.com
paths:
  /go/module1:
      repo: ssh://admin@10.10.30.30/module1
      vcs: git

也就是说，当 govanityurls 收到 nginx 转发的请求后，会将请求与 vanity.yaml 中配置的 module 路径相匹配，如果匹配 ok，就会将该 module 的真实 repo 地址，通过 go 命令期望的应答格式返回。
在这里看到，module1 对应的真实 vcs 上的仓库地址为：ssh://admin@10.10.30.30/module1。

所以，goproxy 会收到这个地址，并再次向这个真实地址发起请求，并最终将 module1 缓存到本地 cache 并返回给客户端。

3）开发机 (客户端) 的设置
将开发机的 GOPROXY 环境变量，设置为 goproxy 的服务地址。
但说过，凡是通过 GOPROXY 拉取的 Go Module，go 命令都会默认把它的 sum 值放到公共 GOSUM 服务器上去校验。

但实质上拉取的是私有 Go Module，GOSUM 服务器上并没有我们的 Go Module 的 sum 数据。
这样就会导致 go build 命令报错，无法继续构建过程。

因此，开发机客户端还需要将 mycompany.com/go，作为一个值设置到 GONOSUMDB 环境变量中：
export GONOSUMDB=mycompany.com/go

这个环境变量配置一旦生效，就相当于告诉 go 命令，凡是与 mycompany.com/go 匹配的 Go Module，都不需要在做 sum 校验了。

到这里，就实现了拉取私有 Go Module 的方案。

4）方案的“不足”
第一点：开发者还是需要额外配置 GONOSUMDB 变量。

由于 Go 命令默认会对从 GOPROXY 拉取的 Go Module 进行 sum 校验，因此需要将私有 Go Module 配置到 GONOSUMDB 环境变量中，这就给开发者带来了一个小小的“负担”。

对于这个问题，解决建议是：公司内部可以将私有 go 项目都放在一个特定域名下，这样就不需要为每个 go 私有项目单独增加 GONOSUMDB 配置了，只需要配置一次就可以了。

第二点：新增私有 Go Module，vanity.yaml 需要手工同步更新。
这是这个方案最不灵活的地方了，由于目前 govanityurls 功能有限，针对每个私有 Go Module，可能都需要单独配置它对应的 vcs 仓库地址，以及获取方式（git、svn or hg）。

关于这一点，建议是：在一个 vcs 仓库中管理多个私有 Go Module。
相比于最初 go 官方建议的一个 repo 只管理一个 module，新版本的 go 在一个 repo 下管理多个 Go Module方面，已经有了长足的进步，我们已经可以通过 repo 的 tag 来区别同一个 repo 下的不同 Go Module。

第三点：无法划分权限。
goproxy 所在节点需要具备访问所有私有 Go Module 所在 vcs repo 的权限，但又无法对 go 开发者端做出有差别授权，这样，只要是 goproxy 能拉取到的私有 Go Module，go 开发者都能拉取到。

不过对于多数公司而言，内部所有源码原则上都是企业内部公开的，这个问题似乎也不大。
如果觉得这是个问题，那么只能使用前面提到的第一个方案，也就是直连私有 Go Module 的源码服务器的方案了。








*/

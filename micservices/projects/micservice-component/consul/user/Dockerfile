FROM golang:latest
WORKDIR /temp
COPY . /temp
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
EXPOSE 9527
ENTRYPOINT ["./user"]


#Dockerfile 中出现了五种指令。
#From：Dockerfile 中必须出现的第一个指令，用于指定基础镜像，在上述例子中我们指定了基础镜像为 golang:latest 版本。
#WORKDIR：指定工作目录，之后的命令将会在该目录下执行。
#COPY：将本地文件添加到容器指定位置中。
#RUN：创建镜像执行的命令，一个 Dockerfile 可以有多个 RUN 命令。在上述 RUN 指令中我们指定了 Go 的代理，并通过 go build 构建了 user 服务。
#ENTRYPOINT：容器被启动后执行的命令，每个 Dockerfile 只有一个。我们通过该命令在容器启动后，又启动了 user 服务。
#Dockerfile 放在 user 服务的代码下，即可通过以下命令构建一个带 user 服务的镜像：
#docker build -t user:1.0.0 .
#-t 选项用于指定镜像的名称和标签，不指定标签默认为 latest；命令最后的 . 为 Dockerfile 所在的地址。
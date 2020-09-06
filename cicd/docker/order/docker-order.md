#### 官方安装教程
[https://docs.docker.com/engine/install/centos/]()

#### 查看本机docker相关详细信息
docker info

#### 













-------------------harbor使用教程----------------------
#### 客户端登陆harbor
docker -H IP:2375 info 
如果存在443这种安全服务登陆,解决方案如下:
vim /etc/docker/daemon.json
{
    "insecure-registries": [
       "<HOSTNAME/IP>"
    ]
}
#### 为镜像打tag
docker tag nginx:v1 192.168.99.124/k8s-dev/nginx:v1
#### 将镜像推送到harbor
docker push 192.168.99.124/k8s-dev/nginx:v1
#### 登出registry
docker logout
#### 拉取镜像
docker pull 192.168.99.124/k8s-dev/nginx:v1





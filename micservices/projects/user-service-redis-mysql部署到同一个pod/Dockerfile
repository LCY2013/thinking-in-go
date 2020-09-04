FROM mysql:8.0.20
WORKDIR /docker-entrypoint-initdb.d
ENV LANG=C.UTF-8
COPY init.sql .


# MySQL 的官方镜像支持在容器启动的时候自动执行指定的 sql 脚本或者 shell 脚本，
# 只要在构建容器时将相关 sql 脚本或者 shell 脚本复制到 /docker-entrypoint-initdb.d 目录下即可。
# 比如上述例子中，我们把初始化 init.sql 放到该目录下，让容器在启动时帮我们初始化好 user 数据库

# docker build -t mysql-for-user .
# docker run  -itd --name mysql-for-user -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql-for-user

# 通过以上命令就构建和启动了 mysql-for-user 镜像，并指定 MySQL 的 root 账户密码为 123456，
# 然后我们再通过以下命令启动 user 容器:
# docker run -itd --name user --network host user

#与前面的启动命令不同，我们将使用 host 网络模式启动 user 容器，这意味着 user 容器的网络与宿主机的网络是一样的，
#这样启动的目的是方便我们在 user 容器的代码中直接使用 localhost 访问部署在同一台宿主机的 Redis 和 MySQL 容器。
#user 容器启动成功后，我们就可以通过进入到 user 容器中或者在宿主机中访问 user 服务中的接口。
#非 Linux 的宿主机不支持 host 网络模式，无法将容器暴露的端口直接绑定到宿主机，只能进入到 user 容器内访问接口。
#可通过以下指令访问 user 容器：
# docker exec -it user /bin/bash








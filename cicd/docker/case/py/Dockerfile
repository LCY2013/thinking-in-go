# 使用官方提供的 Python 开发镜像作为基础镜像
FROM python:2.7-slim
# 切换工作目录到/app
WORKDIR /app
# 将当前目录中的文件负责到容器的/app目录下
# add 可以添加本地或者远程的文件到目标目录中
ADD . /app
# 使用pip安装所有的该项目py依赖环境
RUN pip install --trusted-host pypi.python.org -r requirements.txt
# 允许外界访问容器的80端口
EXPOSE 80
# 设置环境变量
ENV NAME lcy
# 设置容器启动命令 python app.py , 等价于docker run python app.py
## ENTRYPOINT 这个指令和CMD是一样的效果都是运行容器进程 ENTRYPOINT CMD
## 默认情况下，Docker会提供一个 /bin/sh -c，在这个例子中即执行的就是/bin/sh -c python app.py, ENTRYPOINT CMD
## 统一称 Docker 容器的启动进程为 ENTRYPOINT，而不是 CMD。
CMD ["python","app.py"]
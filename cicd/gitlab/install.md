# 创建 docker-compose.yml 文件
```shell
version: '3.7'

services:
  gitlab:
    image: gitlab/gitlab-ce:latest
    container_name: gitlab
    hostname: 'gitlab.example.com'  # 替换为你的域名或IP
    restart: always
    environment:
      GITLAB_OMNIBUS_CONFIG: |
        external_url 'http://192.168.0.180'  # 替换为你的域名或IP
        # 配置SMTP邮件服务（可选）
        gitlab_rails['smtp_enable'] = true
        gitlab_rails['smtp_address'] = "smtp.example.com"
        gitlab_rails['smtp_port'] = 587
        gitlab_rails['smtp_user_name'] = "your_email@example.com"
        gitlab_rails['smtp_password'] = "your_password"
        gitlab_rails['smtp_domain'] = "example.com"
        gitlab_rails['smtp_authentication'] = "login"
        gitlab_rails['smtp_enable_starttls_auto'] = true
        gitlab_rails['gitlab_email_from'] = "gitlab@example.com"
    ports:
      - "80:80"          # HTTP访问
#      - "443:443"        # HTTPS访问（配置SSL后使用）
#      - "2222:22"        # SSH访问（避免与宿主机22端口冲突）
    volumes:
      - ./gitlab_config:/etc/gitlab
      - ./gitlab_logs:/var/log/gitlab
      - ./gitlab_data:/var/opt/gitlab
    networks:
      - gitlab_network

volumes:
  gitlab_config:
  gitlab_logs:
  gitlab_data:

networks:
  gitlab_network:
    driver: bridge
```

# 启动 GitLab 容器
```shell
# 创建并启动容器（后台运行）
docker-compose up -d

# 查看启动日志
docker-compose logs -f gitlab
```

# 访问 GitLab
- 等待约2-5分钟初始化完成

- 在浏览器访问 http://your-server-ip（或你配置的域名）

- 首次访问需要设置 root 用户密码（最少8个字符）

## 常用管理命令
## 进入容器执行命令
docker exec -it gitlab gitlab-rake gitlab:check

## 重启GitLab服务
docker-compose restart gitlab

## 停止服务
docker-compose down

## 备份GitLab数据
docker exec -t gitlab gitlab-backup create

## 查看资源使用情况
docker stats gitlab

# 配置优化建议

## 调整资源配置（编辑docker-compose.yml）
```shell
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 4G
    reservations:
      memory: 2G
```

## 启用 HTTPS（需要准备SSL证书）
```shell
environment:
  GITLAB_OMNIBUS_CONFIG: |
    external_url 'https://gitlab.example.com'
    nginx['redirect_http_to_https'] = true
    nginx['ssl_certificate'] = "/etc/gitlab/ssl/gitlab.example.com.crt"
    nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/gitlab.example.com.key"
```

# 备份与恢复

## 创建备份
> docker exec -t gitlab gitlab-backup create STRATEGY=copy

## 恢复备份
```shell
# 停止相关服务
docker exec -it gitlab gitlab-ctl stop unicorn
docker exec -it gitlab gitlab-ctl stop sidekiq

# 恢复备份（替换BACKUP_TIMESTAMP）
docker exec -it gitlab gitlab-backup restore BACKUP=BACKUP_TIMESTAMP

# 重启服务
docker exec -it gitlab gitlab-ctl restart
```


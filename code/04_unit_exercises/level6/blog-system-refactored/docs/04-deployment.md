# 部署与运行指南 🚀

## 概述 📋

本文档详细介绍了博客系统的部署和运行方法，包括开发环境搭建、生产环境部署、配置管理和常见问题解决方案。系统支持多种部署方式，适用于不同的使用场景。

## 系统要求 💻

### 最低系统要求

- **操作系统**: Windows 10/11, macOS 10.15+, Linux (Ubuntu 18.04+)
- **CPU**: 双核 2.0GHz 或更高
- **内存**: 4GB RAM (推荐 8GB+)
- **存储**: 10GB 可用空间
- **网络**: 稳定的互联网连接

### 软件依赖

- **Go**: 1.19+ (推荐 1.21+)
- **Git**: 2.30+
- **数据库**: MySQL 8.0+ 或 SQLite 3.35+
- **可选**: Docker 20.10+, Docker Compose 2.0+

## 开发环境搭建 🛠️

### 1. 安装 Go 语言环境

#### Windows 系统

```powershell
# 方法1: 使用 Chocolatey
choco install golang

# 方法2: 手动下载安装
# 访问 https://golang.org/dl/ 下载 Windows 安装包
# 运行安装程序，按提示完成安装

# 验证安装
go version
```

#### macOS 系统

```bash
# 方法1: 使用 Homebrew
brew install go

# 方法2: 手动下载安装
# 访问 https://golang.org/dl/ 下载 macOS 安装包
# 运行安装程序，按提示完成安装

# 验证安装
go version
```

#### Linux 系统

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang

# 或者手动安装最新版本
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

### 2. 配置 Go 环境变量

```bash
# 设置 GOPATH 和 GOROOT
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# 启用 Go Modules
export GO111MODULE=on

# 配置代理（中国大陆用户推荐）
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
```

### 3. 安装数据库

#### MySQL 安装

**Windows:**
```powershell
# 使用 Chocolatey
choco install mysql

# 或下载 MySQL Installer
# https://dev.mysql.com/downloads/installer/
```

**macOS:**
```bash
# 使用 Homebrew
brew install mysql
brew services start mysql

# 安全配置
mysql_secure_installation
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# CentOS/RHEL
sudo yum install mysql-server
sudo systemctl start mysqld
sudo systemctl enable mysqld

# 安全配置
sudo mysql_secure_installation
```

#### SQLite 安装

```bash
# Ubuntu/Debian
sudo apt install sqlite3

# macOS
brew install sqlite

# Windows
choco install sqlite
```

### 4. 克隆项目代码

```bash
# 克隆仓库
git clone <repository-url>
cd blog-system-refactored

# 查看项目结构
ls -la
```

### 5. 安装项目依赖

```bash
# 初始化 Go Modules
go mod tidy

# 下载依赖包
go mod download

# 验证依赖
go mod verify
```

## 配置管理 ⚙️

### 1. 配置文件结构

```
configs/
├── config.yaml          # 主配置文件
├── config.dev.yaml      # 开发环境配置
├── config.prod.yaml     # 生产环境配置
└── config.test.yaml     # 测试环境配置
```

### 2. 主配置文件 (config.yaml)

```yaml
# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release, test
  read_timeout: 30s
  write_timeout: 30s
  max_header_bytes: 1048576

# 数据库配置
database:
  type: "mysql"  # mysql, sqlite
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "blog_system"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600s
  
# SQLite 配置示例
# database:
#   type: "sqlite"
#   path: "./blog.db"

# JWT 配置
jwt:
  secret: "your-secret-key-here"
  expires_in: 24h
  refresh_expires_in: 168h  # 7 days

# 日志配置
log:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, text
  output: "stdout"  # stdout, file
  file_path: "./logs/app.log"
  max_size: 100  # MB
  max_backups: 5
  max_age: 30  # days
  compress: true

# 限流配置
rate_limit:
  requests_per_second: 100
  burst: 200
  cleanup_interval: 60s

# CORS 配置
cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:8080"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Origin"
    - "Content-Type"
    - "Authorization"
  allow_credentials: true
  max_age: 86400

# 文件上传配置
upload:
  max_size: 10485760  # 10MB
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
  upload_path: "./uploads"
  url_prefix: "/uploads"

# 缓存配置
cache:
  type: "memory"  # memory, redis
  ttl: 3600s
  cleanup_interval: 600s
  
# Redis 配置示例
# cache:
#   type: "redis"
#   host: "localhost"
#   port: 6379
#   password: ""
#   db: 0
#   ttl: 3600s

# 邮件配置
mail:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from_name: "博客系统"
  from_email: "your-email@gmail.com"
```

### 3. 环境变量配置

创建 `.env` 文件：

```bash
# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库配置
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_DATABASE=blog_system

# JWT 密钥
JWT_SECRET=your-very-secure-secret-key-here

# 日志级别
LOG_LEVEL=info

# 环境标识
ENVIRONMENT=development
```

### 4. 数据库初始化

#### MySQL 数据库创建

```sql
-- 连接到 MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE blog_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选）
CREATE USER 'blog_user'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON blog_system.* TO 'blog_user'@'localhost';
FLUSH PRIVILEGES;

-- 退出
EXIT;
```

#### 数据库迁移

```bash
# 运行数据库迁移
go run cmd/migrate/main.go

# 或者启动应用时自动迁移
go run cmd/main.go --migrate
```

## 启动应用 🎯

### 1. 开发环境启动

```bash
# 方法1: 直接运行
go run cmd/main.go

# 方法2: 编译后运行
go build -o blog-system cmd/main.go
./blog-system

# 方法3: 使用 air 热重载（推荐开发时使用）
# 安装 air
go install github.com/cosmtrek/air@latest

# 创建 .air.toml 配置文件
air init

# 启动热重载
air
```

### 2. 指定配置文件启动

```bash
# 使用特定配置文件
go run cmd/main.go --config=configs/config.dev.yaml

# 使用环境变量
ENVIRONMENT=development go run cmd/main.go

# 指定端口
go run cmd/main.go --port=8081
```

### 3. 验证启动

```bash
# 检查健康状态
curl http://localhost:8080/health

# 预期响应
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "database": "connected"
}

# 查看 API 文档
# 浏览器访问: http://localhost:8080/docs
```

## Docker 部署 🐳

### 1. Dockerfile

```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 创建工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# 创建必要的目录
RUN mkdir -p logs uploads

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./main"]
```

### 2. Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  # 博客应用
  blog-app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=mysql
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USERNAME=blog_user
      - DB_PASSWORD=secure_password
      - DB_DATABASE=blog_system
      - JWT_SECRET=your-very-secure-secret-key
      - LOG_LEVEL=info
    volumes:
      - ./logs:/root/logs
      - ./uploads:/root/uploads
      - ./configs:/root/configs
    depends_on:
      - mysql
    restart: unless-stopped
    networks:
      - blog-network

  # MySQL 数据库
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=blog_system
      - MYSQL_USER=blog_user
      - MYSQL_PASSWORD=secure_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
    networks:
      - blog-network
    command: --default-authentication-plugin=mysql_native_password

  # Redis 缓存（可选）
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    networks:
      - blog-network
    command: redis-server --appendonly yes

  # Nginx 反向代理（可选）
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - blog-app
    restart: unless-stopped
    networks:
      - blog-network

volumes:
  mysql_data:
  redis_data:

networks:
  blog-network:
    driver: bridge
```

### 3. Docker 部署命令

```bash
# 构建镜像
docker build -t blog-system .

# 运行容器
docker run -d \
  --name blog-system \
  -p 8080:8080 \
  -e DB_TYPE=sqlite \
  -e DB_PATH=/app/blog.db \
  -v $(pwd)/blog.db:/app/blog.db \
  blog-system

# 使用 Docker Compose
docker-compose up -d

# 查看日志
docker-compose logs -f blog-app

# 停止服务
docker-compose down

# 重新构建并启动
docker-compose up --build -d
```

## 生产环境部署 🏭

### 1. 系统准备

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装必要软件
sudo apt install -y git curl wget unzip

# 创建应用用户
sudo useradd -m -s /bin/bash blog
sudo usermod -aG sudo blog

# 切换到应用用户
su - blog
```

### 2. 应用部署

```bash
# 创建应用目录
mkdir -p /home/blog/app
cd /home/blog/app

# 克隆代码
git clone <repository-url> .

# 编译应用
go build -o blog-system cmd/main.go

# 设置权限
chmod +x blog-system

# 创建必要目录
mkdir -p logs uploads configs

# 复制配置文件
cp configs/config.prod.yaml configs/config.yaml
```

### 3. Systemd 服务配置

创建服务文件 `/etc/systemd/system/blog-system.service`：

```ini
[Unit]
Description=Blog System API Server
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=blog
Group=blog
WorkingDirectory=/home/blog/app
ExecStart=/home/blog/app/blog-system
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=blog-system

# 环境变量
Environment=ENVIRONMENT=production
Environment=LOG_LEVEL=info
Environment=DB_TYPE=mysql
Environment=DB_HOST=localhost
Environment=DB_PORT=3306
Environment=DB_USERNAME=blog_user
Environment=DB_PASSWORD=secure_password
Environment=DB_DATABASE=blog_system
Environment=JWT_SECRET=your-production-secret-key

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/home/blog/app/logs /home/blog/app/uploads

[Install]
WantedBy=multi-user.target
```

### 4. 启动和管理服务

```bash
# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start blog-system

# 设置开机自启
sudo systemctl enable blog-system

# 查看服务状态
sudo systemctl status blog-system

# 查看日志
sudo journalctl -u blog-system -f

# 重启服务
sudo systemctl restart blog-system

# 停止服务
sudo systemctl stop blog-system
```

### 5. Nginx 反向代理配置

创建 Nginx 配置文件 `/etc/nginx/sites-available/blog-system`：

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;
    
    # SSL 证书配置
    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # 安全头
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # 日志配置
    access_log /var/log/nginx/blog-system.access.log;
    error_log /var/log/nginx/blog-system.error.log;
    
    # 客户端上传限制
    client_max_body_size 10M;
    
    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }
    
    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:8080;
        access_log off;
    }
    
    # 文档页面
    location /docs {
        proxy_pass http://127.0.0.1:8080;
    }
    
    # 静态文件
    location /uploads/ {
        alias /home/blog/app/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # 前端应用（如果有）
    location / {
        root /var/www/blog-frontend;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
}
```

启用站点：

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/blog-system /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重新加载 Nginx
sudo systemctl reload nginx
```

## 监控和日志 📊

### 1. 应用监控

```bash
# 查看应用状态
curl http://localhost:8080/health

# 查看系统资源
top -p $(pgrep blog-system)
htop

# 查看端口占用
sudo netstat -tlnp | grep :8080
sudo ss -tlnp | grep :8080
```

### 2. 日志管理

```bash
# 查看应用日志
tail -f /home/blog/app/logs/app.log

# 查看系统日志
sudo journalctl -u blog-system -f

# 查看 Nginx 日志
sudo tail -f /var/log/nginx/blog-system.access.log
sudo tail -f /var/log/nginx/blog-system.error.log

# 日志轮转配置
sudo vim /etc/logrotate.d/blog-system
```

日志轮转配置内容：

```
/home/blog/app/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 blog blog
    postrotate
        systemctl reload blog-system
    endscript
}
```

### 3. 性能监控

```bash
# 安装监控工具
sudo apt install -y htop iotop nethogs

# CPU 和内存监控
htop

# 磁盘 I/O 监控
iotop

# 网络监控
nethogs

# 数据库监控
mysql -u root -p -e "SHOW PROCESSLIST;"
mysql -u root -p -e "SHOW STATUS LIKE 'Threads%';"
```

## 备份策略 💾

### 1. 数据库备份脚本

创建备份脚本 `/home/blog/scripts/backup.sh`：

```bash
#!/bin/bash

# 配置变量
DB_NAME="blog_system"
DB_USER="blog_user"
DB_PASS="secure_password"
BACKUP_DIR="/home/blog/backups"
DATE=$(date +"%Y%m%d_%H%M%S")
RETENTION_DAYS=30

# 创建备份目录
mkdir -p $BACKUP_DIR

# 数据库备份
echo "开始备份数据库..."
mysqldump -u$DB_USER -p$DB_PASS $DB_NAME > $BACKUP_DIR/db_backup_$DATE.sql

# 压缩备份文件
gzip $BACKUP_DIR/db_backup_$DATE.sql

# 备份应用文件
echo "开始备份应用文件..."
tar -czf $BACKUP_DIR/app_backup_$DATE.tar.gz \
    -C /home/blog/app \
    --exclude='logs' \
    --exclude='*.log' \
    .

# 清理旧备份
echo "清理旧备份文件..."
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "备份完成: $DATE"
```

### 2. 定时备份

```bash
# 设置执行权限
chmod +x /home/blog/scripts/backup.sh

# 添加到 crontab
crontab -e

# 每天凌晨 2 点执行备份
0 2 * * * /home/blog/scripts/backup.sh >> /home/blog/logs/backup.log 2>&1
```

## 常见问题解决 🔧

### 1. 端口占用问题

```bash
# 查看端口占用
sudo netstat -tlnp | grep :8080

# 杀死占用进程
sudo kill -9 <PID>

# 或者修改配置文件中的端口
vim configs/config.yaml
```

### 2. 数据库连接问题

```bash
# 检查数据库服务状态
sudo systemctl status mysql

# 测试数据库连接
mysql -u blog_user -p -h localhost blog_system

# 检查防火墙设置
sudo ufw status
sudo ufw allow 3306
```

### 3. 权限问题

```bash
# 检查文件权限
ls -la /home/blog/app/

# 修复权限
sudo chown -R blog:blog /home/blog/app/
sudo chmod -R 755 /home/blog/app/
sudo chmod +x /home/blog/app/blog-system
```

### 4. 内存不足问题

```bash
# 查看内存使用
free -h

# 查看进程内存占用
ps aux --sort=-%mem | head

# 添加交换空间
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# 永久启用交换空间
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

### 5. SSL 证书问题

```bash
# 使用 Let's Encrypt 免费证书
sudo apt install certbot python3-certbot-nginx

# 申请证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 性能优化建议 ⚡

### 1. 应用层优化

- 启用 Gzip 压缩
- 使用连接池
- 实现缓存机制
- 优化数据库查询
- 使用 CDN 加速静态资源

### 2. 系统层优化

```bash
# 调整文件描述符限制
echo "* soft nofile 65535" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65535" | sudo tee -a /etc/security/limits.conf

# 调整内核参数
echo "net.core.somaxconn = 65535" | sudo tee -a /etc/sysctl.conf
echo "net.ipv4.tcp_max_syn_backlog = 65535" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### 3. 数据库优化

```sql
-- 优化 MySQL 配置
-- 编辑 /etc/mysql/mysql.conf.d/mysqld.cnf

[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
max_connections = 200
query_cache_size = 64M
query_cache_type = 1
```

---

**注意**：本文档提供了完整的部署指南，请根据实际环境调整配置参数。生产环境部署前请务必进行充分测试。
# 企业级在线教育平台后端系统

## 项目概述

这是一个基于 Go + GORM + Gin 构建的企业级在线教育平台后端系统，展示了完整的项目架构设计、数据模型设计、业务逻辑实现和 API 接口开发。

## 技术栈

- **语言**: Go 1.19+
- **Web框架**: Gin
- **ORM**: GORM v2
- **数据库**: MySQL 8.0+
- **缓存**: Redis（配置支持）
- **配置管理**: Viper
- **日志**: 结构化日志
- **认证**: JWT（简化实现）

## 项目特性

### 🏗️ 架构设计
- 分层架构：Controller -> Service -> Model
- 依赖注入和接口抽象
- 统一的错误处理和响应格式
- 配置文件管理
- 中间件支持

### 📊 数据模型
- 用户系统（用户、角色、用户资料）
- 课程系统（分类、课程、章节、课时）
- 订单系统（订单、订单项、支付）
- 学习系统（学习进度、课程评价、收藏）
- 营销系统（优惠券、通知）
- 系统日志

### 🔐 权限管理
- 基于角色的访问控制（RBAC）
- JWT 认证机制
- 接口权限验证
- 用户状态管理

### 💰 业务功能
- 用户注册、登录、资料管理
- 课程浏览、搜索、筛选
- 购物车、订单、支付
- 学习进度跟踪
- 课程评价和收藏
- 优惠券系统

### 🚀 性能优化
- 数据库连接池
- 查询优化和索引
- 分页查询
- 预加载关联数据
- 缓存策略（配置支持）

## 项目结构

```
edu-platform/
├── cmd/                    # 应用入口
│   └── server/
│       └── main.go
├── internal/               # 内部代码
│   ├── config/            # 配置管理
│   │   └── config.go
│   ├── models/            # 数据模型
│   │   └── models.go
│   ├── services/          # 业务逻辑
│   │   └── services.go
│   ├── controllers/       # API控制器
│   │   └── controllers.go
│   └── middleware/        # 中间件
├── pkg/                   # 公共包
├── configs/               # 配置文件
│   └── config.yaml
├── docs/                  # 文档
├── scripts/               # 脚本
├── tests/                 # 测试
├── go.mod
├── go.sum
└── README.md
```

## 数据库设计

### 核心表结构

#### 用户相关
- `users` - 用户基本信息
- `roles` - 角色定义
- `user_profiles` - 用户详细资料

#### 课程相关
- `categories` - 课程分类
- `courses` - 课程信息
- `chapters` - 课程章节
- `lessons` - 课程课时
- `course_reviews` - 课程评价
- `course_favorites` - 课程收藏

#### 订单相关
- `orders` - 订单主表
- `order_items` - 订单详情
- `coupons` - 优惠券

#### 学习相关
- `learning_progress` - 学习进度

#### 系统相关
- `notifications` - 系统通知
- `system_logs` - 系统日志

## API 接口

### 用户接口
```
POST   /api/users/register     # 用户注册
POST   /api/users/login        # 用户登录
GET    /api/users/profile      # 获取用户资料
PUT    /api/users/profile      # 更新用户资料
GET    /api/admin/users        # 获取用户列表（管理员）
```

### 课程接口
```
GET    /api/courses            # 获取课程列表
GET    /api/courses/:id        # 获取课程详情
POST   /api/courses            # 创建课程（讲师）
PUT    /api/courses/:id        # 更新课程
POST   /api/courses/:id/publish # 发布课程
```

### 订单接口
```
POST   /api/orders             # 创建订单
GET    /api/orders             # 获取订单列表
POST   /api/orders/:order_no/pay # 支付订单
DELETE /api/orders/:order_no   # 取消订单
```

### 学习接口
```
GET    /api/learning/courses   # 获取学习的课程
POST   /api/learning/progress  # 更新学习进度
GET    /api/learning/courses/:course_id/progress # 获取课程学习进度
```

## 快速开始

### 1. 环境准备

```bash
# 安装 Go 1.19+
# 安装 MySQL 8.0+
# 安装 Redis（可选）
```

### 2. 克隆项目

```bash
git clone <repository-url>
cd edu-platform
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置数据库

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE edu_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 修改配置

编辑 `config.yaml` 文件，配置数据库连接信息：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: edu_platform
```

### 6. 运行项目

```bash
# 开发模式
go run main.go

# 或者构建后运行
go build -o edu-platform main.go
./edu-platform
```

### 7. 测试接口

服务启动后，访问 `http://localhost:8080`

```bash
# 注册用户
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456",
    "nickname": "测试用户"
  }'

# 用户登录
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'

# 获取课程列表
curl -X GET "http://localhost:8080/api/courses?page=1&page_size=10"
```

## 配置说明

### 数据库配置
```yaml
database:
  host: localhost          # 数据库主机
  port: 3306              # 数据库端口
  username: root          # 用户名
  password: password      # 密码
  database: edu_platform  # 数据库名
  charset: utf8mb4        # 字符集
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数
  conn_max_lifetime: 3600 # 连接最大生存时间（秒）
```

### 服务器配置
```yaml
server:
  host: 0.0.0.0          # 监听地址
  port: 8080             # 监听端口
  mode: debug            # 运行模式：debug/release
  read_timeout: 60       # 读取超时（秒）
  write_timeout: 60      # 写入超时（秒）
```

### JWT配置
```yaml
jwt:
  secret: your-secret-key # JWT密钥
  expires_in: 86400      # 过期时间（秒）
  issuer: edu-platform   # 签发者
```

## 开发指南

### 添加新的API接口

1. **定义数据模型**（如果需要）
```go
// internal/models/models.go
type NewModel struct {
    BaseModel
    Name string `json:"name" gorm:"size:100;not null"`
    // 其他字段...
}
```

2. **实现业务逻辑**
```go
// internal/services/services.go
type NewService struct {
    db *gorm.DB
}

func (s *NewService) CreateNew(model *models.NewModel) error {
    return s.db.Create(model).Error
}
```

3. **添加控制器**
```go
// internal/controllers/controllers.go
func (ctrl *NewController) CreateNew(c *gin.Context) {
    // 实现逻辑
}
```

4. **注册路由**
```go
// main.go
api.POST("/news", newController.CreateNew)
```

### 数据库迁移

项目启动时会自动执行数据库迁移，创建所需的表结构。如果需要手动迁移：

```go
db.AutoMigrate(
    &models.User{},
    &models.Role{},
    &models.Course{},
    // 其他模型...
)
```

### 添加中间件

```go
// 自定义中间件
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}

// 使用中间件
r.Use(CustomMiddleware())
```

## 部署指南

### Docker 部署

```dockerfile
# Dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o edu-platform main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/edu-platform .
COPY --from=builder /app/config.yaml .
CMD ["./edu-platform"]
```

```bash
# 构建镜像
docker build -t edu-platform .

# 运行容器
docker run -p 8080:8080 edu-platform
```

### 生产环境配置

1. **修改运行模式**
```yaml
server:
  mode: release
```

2. **配置日志**
```yaml
logger:
  level: info
  format: json
  output: file
  filename: logs/app.log
```

3. **配置HTTPS**
```yaml
server:
  tls:
    enabled: true
    cert_file: cert.pem
    key_file: key.pem
```

## 测试

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/services

# 运行测试并显示覆盖率
go test -cover ./...
```

### 集成测试

```bash
# 运行集成测试
go test -tags=integration ./tests/...
```

## 性能优化

### 数据库优化

1. **索引优化**
```sql
-- 用户表索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- 课程表索引
CREATE INDEX idx_courses_category_id ON courses(category_id);
CREATE INDEX idx_courses_instructor_id ON courses(instructor_id);
CREATE INDEX idx_courses_status ON courses(status);
```

2. **查询优化**
```go
// 使用预加载避免N+1问题
db.Preload("Category").Preload("Instructor").Find(&courses)

// 使用选择字段减少数据传输
db.Select("id, title, price").Find(&courses)

// 使用原生SQL处理复杂查询
db.Raw("SELECT ... FROM ... WHERE ...").Scan(&result)
```

### 缓存策略

```go
// Redis缓存示例
func (s *CourseService) GetCourseFromCache(id uint) (*models.Course, error) {
    key := fmt.Sprintf("course:%d", id)
    
    // 从缓存获取
    if cached := s.redis.Get(key); cached != nil {
        var course models.Course
        json.Unmarshal(cached, &course)
        return &course, nil
    }
    
    // 从数据库获取
    course, err := s.GetCourseByID(id)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    data, _ := json.Marshal(course)
    s.redis.Set(key, data, time.Hour)
    
    return course, nil
}
```

## 监控和日志

### 结构化日志

```go
import "github.com/sirupsen/logrus"

// 配置日志
logrus.SetFormatter(&logrus.JSONFormatter{})
logrus.SetLevel(logrus.InfoLevel)

// 使用日志
logrus.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "create_order",
}).Info("Order created successfully")
```

### 性能监控

```go
// 中间件记录请求时间
func RequestTimeMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        
        logrus.WithFields(logrus.Fields{
            "method":   c.Request.Method,
            "path":     c.Request.URL.Path,
            "status":   c.Writer.Status(),
            "duration": duration.Milliseconds(),
        }).Info("Request completed")
    }
}
```

## 安全考虑

### 输入验证

```go
// 使用binding标签验证输入
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### SQL注入防护

```go
// 使用参数化查询
db.Where("email = ?", email).First(&user)

// 避免直接拼接SQL
// 错误：db.Where(fmt.Sprintf("email = '%s'", email))
```

### 密码安全

```go
import "golang.org/x/crypto/bcrypt"

// 密码加密
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// 密码验证
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## 常见问题

### Q: 数据库连接失败
A: 检查数据库配置信息，确保MySQL服务正在运行，用户名密码正确。

### Q: 端口被占用
A: 修改配置文件中的端口号，或者停止占用端口的进程。

### Q: JWT认证失败
A: 检查JWT密钥配置，确保客户端正确传递Authorization头。

### Q: 性能问题
A: 检查数据库索引，使用查询分析工具，考虑添加缓存。

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目维护者：[Your Name]
- 邮箱：[your.email@example.com]
- 项目链接：[https://github.com/yourusername/edu-platform]

---

**注意**: 这是一个学习项目，用于演示企业级Go应用的开发模式。在生产环境中使用前，请确保进行充分的安全审计和性能测试。
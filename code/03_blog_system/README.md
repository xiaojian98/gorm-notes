# 博客系统实战项目

这是一个基于 Gin + GORM 的博客系统实战项目，对应学习文档：`02_GORM背景示例_博客系统实战.md`

## 项目特色

- **完整的博客功能**：用户管理、文章发布、评论系统、分类标签
- **RESTful API设计**：标准的REST接口，易于前端集成
- **GORM最佳实践**：展示GORM的各种用法和最佳实践
- **数据库关联**：一对一、一对多、多对多关联的完整示例
- **软删除和审计**：完整的数据审计和软删除机制
- **性能优化**：预加载、索引、连接池等性能优化

## 技术栈

- **Web框架**：Gin
- **ORM**：GORM v2
- **数据库**：SQLite（开发）/ MySQL（生产）
- **密码加密**：bcrypt
- **API文档**：内置API文档

## 项目结构

```
blog-system/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块文件
├── README.md              # 项目说明
├── config/                # 配置相关
│   └── database.go        # 数据库配置
├── models/                # 数据模型
│   └── models.go          # 所有数据模型定义
├── services/              # 业务逻辑层
│   └── services.go        # 业务服务实现
├── handlers/              # 控制器层
│   └── handlers.go        # HTTP处理器
├── routes/                # 路由配置
│   └── routes.go          # 路由定义
└── blog.db               # SQLite数据库文件（运行后生成）
```

## 数据模型关系

### 核心实体
- **User（用户）**：系统用户
- **Profile（用户资料）**：用户详细信息
- **Post（文章）**：博客文章
- **Comment（评论）**：文章评论
- **Category（分类）**：文章分类
- **Tag（标签）**：文章标签
- **Like（点赞）**：文章点赞

### 关联关系
- User ↔ Profile：一对一关系
- User → Post：一对多关系（用户可以发布多篇文章）
- User → Comment：一对多关系（用户可以发表多条评论）
- Category → Post：一对多关系（分类包含多篇文章）
- Post ↔ Tag：多对多关系（文章可以有多个标签）
- Post → Comment：一对多关系（文章可以有多条评论）
- Post → Like：一对多关系（文章可以被多次点赞）
- Comment → Comment：自关联（评论回复）

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行项目

```bash
go run main.go
```

### 3. 访问API

服务启动后，访问以下地址：

- **健康检查**：http://localhost:8080/health
- **API文档**：http://localhost:8080/api/docs
- **用户注册**：POST http://localhost:8080/api/users/register
- **用户登录**：POST http://localhost:8080/api/users/login

## API接口说明

### 用户相关
- `POST /api/users/register` - 用户注册
- `POST /api/users/login` - 用户登录
- `GET /api/users/:id` - 获取用户信息
- `PUT /api/users/:id/profile` - 更新用户资料

### 文章相关
- `GET /api/posts` - 获取文章列表（支持分页、搜索、筛选）
- `GET /api/posts/:id` - 获取文章详情
- `GET /api/posts/slug/:slug` - 通过slug获取文章
- `POST /api/posts` - 创建文章
- `PUT /api/posts/:id` - 更新文章
- `DELETE /api/posts/:id` - 删除文章
- `POST /api/posts/:id/publish` - 发布文章
- `POST /api/posts/:id/like` - 点赞文章
- `DELETE /api/posts/:id/like` - 取消点赞

### 评论相关
- `GET /api/comments/post/:post_id` - 获取文章评论
- `POST /api/comments` - 创建评论
- `PUT /api/comments/:id/approve` - 审核通过评论
- `PUT /api/comments/:id/reject` - 拒绝评论

### 分类相关
- `GET /api/categories` - 获取分类列表
- `GET /api/categories/:slug` - 获取分类详情
- `POST /api/categories` - 创建分类

### 标签相关
- `GET /api/tags` - 获取标签列表
- `GET /api/tags/popular` - 获取热门标签
- `POST /api/tags` - 创建标签

### 统计相关
- `GET /api/stats/overview` - 获取统计概览
- `GET /api/stats/posts/popular` - 获取热门文章

## 示例请求

### 用户注册
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "nickname": "测试用户"
  }'
```

### 创建文章
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "我的第一篇文章",
    "content": "这是文章内容...",
    "excerpt": "文章摘要",
    "author_id": 1,
    "category_id": 1,
    "tags": [1, 2],
    "status": "published"
  }'
```

### 获取文章列表
```bash
curl "http://localhost:8080/api/posts?page=1&limit=10&category=tech&search=golang"
```

## 学习要点

### 1. GORM模型定义
- 基础模型和字段标签
- 关联关系定义
- 钩子函数使用
- 索引和约束

### 2. 数据库操作
- CRUD基本操作
- 复杂查询和关联查询
- 事务处理
- 预加载优化

### 3. 业务逻辑
- 服务层设计
- 错误处理
- 数据验证
- 业务规则实现

### 4. API设计
- RESTful接口设计
- 请求参数验证
- 响应格式统一
- 错误处理机制

## 扩展功能

可以基于此项目继续实现：

1. **用户认证**：JWT token认证
2. **权限控制**：基于角色的访问控制
3. **文件上传**：图片上传和管理
4. **搜索功能**：全文搜索
5. **缓存机制**：Redis缓存
6. **消息队列**：异步任务处理
7. **监控日志**：性能监控和日志记录

## 常见问题

### Q: 如何切换到MySQL数据库？
A: 修改 `config/database.go` 中的 `InitDB` 函数，使用 `InitMySQLDB` 替代SQLite配置。

### Q: 如何添加新的数据模型？
A: 在 `models/models.go` 中定义新模型，然后在 `AutoMigrate` 函数中添加迁移。

### Q: 如何自定义API响应格式？
A: 修改 `handlers/handlers.go` 中的响应结构，或创建统一的响应包装器。

## 相关文档

- [GORM官方文档](https://gorm.io/docs/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [学习文档：02_GORM背景示例_博客系统实战.md](../docs/02_GORM背景示例_博客系统实战.md)

## 许可证

MIT License
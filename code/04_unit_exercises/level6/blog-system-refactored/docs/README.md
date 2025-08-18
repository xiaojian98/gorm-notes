# 博客系统技术文档 📚

## 项目概述 🎯

本项目是一个基于 Go 语言和 Gin 框架开发的现代化博客系统，采用分层架构设计，支持用户管理、文章发布、评论互动、数据分析等完整功能。系统设计遵循 RESTful API 规范，具备良好的可扩展性和维护性。

### 主要特性 ✨

- 🔐 **用户系统**：完整的用户注册、登录、权限管理
- 📝 **内容管理**：文章创建、编辑、发布、分类管理
- 💬 **评论系统**：多级评论、点赞、举报功能
- 📊 **数据分析**：实时统计、趋势分析、性能监控
- 🛡️ **安全防护**：JWT认证、CORS跨域、限流保护
- 🗄️ **数据库支持**：支持 SQLite 和 MySQL 双数据库

## 技术栈 🛠️

### 后端技术
- **语言**：Go 1.23.0
- **框架**：Gin Web Framework
- **ORM**：GORM v2
- **数据库**：MySQL 8.0+ / SQLite 3
- **认证**：JWT (JSON Web Token)
- **限流**：golang.org/x/time/rate

### 开发工具
- **包管理**：Go Modules
- **API文档**：Swagger/OpenAPI
- **日志**：Gin内置日志中间件
- **测试**：Go原生测试框架

## 文档导航 📖

### 核心文档
1. [项目目录结构](./01-project-structure.md) - 详细的目录组织说明
2. [架构设计](./02-architecture.md) - 系统架构和技术栈详解
3. [数据库设计](./03-database.md) - 数据模型和表结构设计
4. [核心功能模块](./04-modules.md) - 主要功能模块介绍
5. [部署运行指南](./05-deployment.md) - 环境搭建和部署说明
6. [开发维护指南](./06-development.md) - 开发规范和维护指南

### 快速开始 🚀

#### 环境要求
- Go 1.23.0+
- MySQL 8.0+ 或 SQLite 3
- Git

#### 快速启动
```bash
# 1. 克隆项目
git clone <repository-url>
cd blog-system-refactored

# 2. 安装依赖
go mod download

# 3. 配置数据库（可选，默认使用SQLite）
# 编辑 internal/config/config.go 中的数据库配置

# 4. 启动服务
go run ./cmd

# 5. 访问服务
# API文档: http://localhost:8080/docs
# 健康检查: http://localhost:8080/health
```

## API 接口概览 🔌

### 用户相关
- `POST /api/v1/users` - 用户注册
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

### 文章相关
- `POST /api/v1/posts` - 创建文章
- `GET /api/v1/posts` - 获取文章列表
- `GET /api/v1/posts/:id` - 获取文章详情
- `PUT /api/v1/posts/:id` - 更新文章
- `DELETE /api/v1/posts/:id` - 删除文章

### 评论相关
- `POST /api/v1/comments` - 创建评论
- `GET /api/v1/comments` - 获取评论列表
- `PUT /api/v1/comments/:id` - 更新评论
- `DELETE /api/v1/comments/:id` - 删除评论

### 分析统计
- `GET /api/v1/analytics/dashboard` - 仪表板统计
- `GET /api/v1/analytics/content` - 内容统计
- `GET /api/v1/analytics/users` - 用户统计

## 项目结构概览 📁

```
blog-system-refactored/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序文件
├── internal/              # 内部包（不对外暴露）
│   ├── config/            # 配置管理
│   ├── handlers/          # HTTP处理器
│   ├── middleware/        # 中间件
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   ├── routes/           # 路由配置
│   └── services/         # 业务逻辑层
├── pkg/                  # 公共包
├── configs/              # 配置文件
├── docs/                 # 项目文档
├── scripts/              # 脚本文件
├── tests/                # 测试文件
├── go.mod               # Go模块文件
└── README.md            # 项目说明
```

## 贡献指南 🤝

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证 📄

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。

## 联系方式 📧

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发送邮件
- 项目讨论区

---

**注意**：本文档会随着项目的发展持续更新，请关注最新版本。
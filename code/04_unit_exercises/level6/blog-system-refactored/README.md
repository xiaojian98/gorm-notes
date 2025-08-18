# 博客系统重构版 📝

一个基于Go语言和Gin框架构建的现代化博客系统，采用分层架构设计，具有完整的用户管理、文章发布、评论互动和数据分析功能。

## ✨ 特性

### 🔐 用户管理
- 用户注册、登录、资料管理
- 用户关注系统
- 角色权限控制（普通用户、管理员）
- 用户状态管理（激活、停用）

### 📖 文章系统
- 文章创建、编辑、发布
- 分类和标签管理
- 文章搜索和筛选
- 热门文章推荐
- 文章浏览量统计

### 💬 评论系统
- 多级评论回复
- 评论点赞功能
- 评论审核机制
- 垃圾评论检测
- 评论举报功能

### 📊 数据分析
- 仪表板统计
- 用户行为分析
- 内容统计报告
- 趋势分析
- 实时数据监控

## 项目结构 🏗️

```
blog-system-refactored/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序文件
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   │   ├── config.go     # 配置结构和加载
│   │   └── database.go   # 数据库配置
│   ├── models/           # 数据模型
│   │   ├── user.go       # 用户模型
│   │   ├── post.go       # 文章模型
│   │   ├── comment.go    # 评论模型
│   │   └── analytics.go  # 分析模型
│   ├── repository/       # 数据访问层
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── comment_repository.go
│   │   └── analytics_repository.go
│   ├── services/         # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── comment_service.go
│   │   └── analytics_service.go
│   ├── handlers/         # HTTP处理层
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   ├── comment_handler.go
│   │   ├── analytics_handler.go
│   │   └── common.go     # 通用响应结构
│   ├── middleware/       # 中间件
│   │   └── middleware.go # 认证、权限、限流等
│   └── routes/           # 路由配置
│       └── routes.go     # 路由设置
├── configs/              # 配置文件
│   └── .env.example     # 环境变量示例
├── docs/                 # 文档
├── scripts/              # 脚本文件
├── tests/                # 测试文件
├── go.mod               # Go模块文件
├── go.sum               # 依赖校验文件
└── README.md            # 项目说明
```

## 功能特性 ✨

- 🔐 用户认证和授权
- 👤 用户资料管理
- 📝 文章发布和管理
- 💬 评论系统
- ❤️ 点赞功能
- 👥 关注系统
- 🔔 通知推送
- 📊 数据分析
- 🗃️ 支持SQLite和MySQL

## 快速开始 🚀

1. 克隆项目
```bash
cd blog-system-refactored
```

2. 安装依赖
```bash
go mod tidy
```

3. 运行程序
```bash
go run cmd/blog/main.go
```

## 数据库支持 🗄️

项目支持两种数据库：
- SQLite（默认，适合开发和测试）
- MySQL（适合生产环境）

## 项目架构 🏛️

本项目采用分层架构：

1. **表示层（Handlers）**: 处理HTTP请求和响应
2. **业务逻辑层（Services）**: 实现核心业务逻辑
3. **数据访问层（Repository）**: 封装数据库操作
4. **数据模型层（Models）**: 定义数据结构

## 开发规范 📋

- 遵循Go语言编码规范
- 使用依赖注入模式
- 接口优先设计
- 完善的错误处理
- 详细的代码注释

## 贡献指南 🤝

欢迎提交Issue和Pull Request来改进项目！

## 许可证 📄

MIT License
# 项目目录结构设计 📁

## 概述 🎯

本项目采用标准的 Go 项目布局，遵循 Go 社区最佳实践，采用分层架构设计，确保代码的可维护性、可扩展性和可测试性。

## 完整目录结构 🌳

```
blog-system-refactored/
├── cmd/                           # 应用程序入口目录
│   └── main.go                   # 主程序文件，应用启动入口
├── internal/                      # 内部包目录（不对外暴露）
│   ├── config/                   # 配置管理模块
│   │   ├── config.go             # 应用配置结构和加载逻辑
│   │   └── database.go           # 数据库配置和连接管理
│   ├── handlers/                 # HTTP 处理器层
│   │   ├── analytics_handler.go  # 数据分析相关接口处理
│   │   ├── comment_handler.go    # 评论相关接口处理
│   │   ├── common.go             # 通用响应结构和工具函数
│   │   ├── post_handler.go       # 文章相关接口处理
│   │   └── user_handler.go       # 用户相关接口处理
│   ├── middleware/               # 中间件模块
│   │   └── middleware.go         # 认证、CORS、限流等中间件
│   ├── models/                   # 数据模型层
│   │   ├── base.go               # 基础模型和通用字段定义
│   │   ├── comment.go            # 评论模型定义
│   │   ├── post.go               # 文章模型定义
│   │   └── user.go               # 用户模型定义
│   ├── repository/               # 数据访问层（Repository 模式）
│   │   ├── analytics_repository.go # 数据分析相关数据访问
│   │   ├── comment_repository.go   # 评论数据访问
│   │   ├── post_repository.go      # 文章数据访问
│   │   └── user_repository.go      # 用户数据访问
│   ├── routes/                   # 路由配置模块
│   │   └── routes.go             # API 路由定义和中间件配置
│   └── services/                 # 业务逻辑层
│       ├── analytics_service.go  # 数据分析业务逻辑
│       ├── comment_service.go    # 评论业务逻辑
│       ├── post_service.go       # 文章业务逻辑
│       └── user_service.go       # 用户业务逻辑
├── pkg/                          # 公共包目录（可被外部引用）
├── configs/                      # 配置文件目录
├── docs/                         # 项目文档目录
│   ├── README.md                 # 主文档
│   ├── 01-project-structure.md   # 项目结构说明（本文档）
│   ├── 02-architecture.md        # 架构设计文档
│   ├── 03-database.md            # 数据库设计文档
│   ├── 04-modules.md             # 功能模块文档
│   ├── 05-deployment.md          # 部署指南
│   └── 06-development.md         # 开发指南
├── scripts/                      # 脚本文件目录
├── tests/                        # 测试文件目录
├── go.mod                        # Go 模块依赖文件
├── go.sum                        # Go 模块校验文件
└── README.md                     # 项目根目录说明文件
```

## 目录详细说明 📋

### 1. cmd/ - 应用程序入口 🚀

**用途**：存放应用程序的主入口文件

**命名规范**：
- 主程序文件命名为 `main.go`
- 如有多个可执行程序，可创建子目录如 `cmd/server/`、`cmd/cli/`

**主要职责**：
- 应用程序初始化
- 配置加载
- 服务启动
- 优雅关闭处理

**示例代码结构**：
```go
// cmd/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "your-project/internal/config"
    "your-project/internal/routes"
)

func main() {
    // 加载配置
    cfg := config.Load()
    
    // 初始化路由
    router := routes.SetupRoutes(cfg)
    
    // 启动服务器
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    // 优雅关闭处理
    // ...
}
```

### 2. internal/ - 内部包目录 🔒

**用途**：存放项目内部使用的包，不对外暴露

**Go 语言特性**：`internal` 目录下的包只能被同一模块内的代码导入

**优势**：
- 封装内部实现细节
- 防止外部依赖内部包
- 提高代码安全性

#### 2.1 internal/config/ - 配置管理 ⚙️

**职责**：
- 应用配置结构定义
- 配置文件加载和解析
- 环境变量处理
- 数据库连接配置

**文件说明**：
- `config.go`：主配置结构和加载逻辑
- `database.go`：数据库配置和连接管理

#### 2.2 internal/handlers/ - HTTP 处理器层 🌐

**职责**：
- HTTP 请求处理
- 请求参数验证
- 响应格式化
- 错误处理

**命名规范**：`{模块名}_handler.go`

**文件说明**：
- `common.go`：通用响应结构、错误码定义、工具函数
- `user_handler.go`：用户相关接口（注册、登录、用户管理）
- `post_handler.go`：文章相关接口（CRUD、发布、分类）
- `comment_handler.go`：评论相关接口（评论管理、审核）
- `analytics_handler.go`：数据分析接口（统计、报表）

#### 2.3 internal/middleware/ - 中间件模块 🛡️

**职责**：
- 请求预处理
- 身份认证和授权
- 跨域处理
- 限流保护
- 日志记录

**主要中间件**：
- `CORS`：跨域资源共享处理
- `Logger`：HTTP 请求日志记录
- `Recovery`：panic 恢复处理
- `RateLimit`：请求限流保护
- `AuthRequired`：JWT 身份认证
- `AdminRequired`：管理员权限验证

#### 2.4 internal/models/ - 数据模型层 📊

**职责**：
- 数据库表结构定义
- 模型关联关系
- 数据验证规则
- 模型方法定义

**命名规范**：`{模型名}.go`

**文件说明**：
- `base.go`：基础模型、通用字段、接口定义
- `user.go`：用户模型、用户资料、关注关系
- `post.go`：文章模型、分类、标签、元数据
- `comment.go`：评论模型、评论状态、关联关系

#### 2.5 internal/repository/ - 数据访问层 🗄️

**职责**：
- 数据库操作封装
- SQL 查询实现
- 事务管理
- 数据访问接口定义

**设计模式**：Repository 模式

**命名规范**：`{模块名}_repository.go`

**主要特点**：
- 接口与实现分离
- 支持单元测试
- 数据库无关的业务逻辑

#### 2.6 internal/services/ - 业务逻辑层 💼

**职责**：
- 核心业务逻辑实现
- 业务规则验证
- 跨模块协调
- 事务管理

**命名规范**：`{模块名}_service.go`

**设计原则**：
- 单一职责原则
- 依赖注入
- 接口编程

#### 2.7 internal/routes/ - 路由配置 🛣️

**职责**：
- API 路由定义
- 中间件配置
- 路由分组
- 版本管理

**主要功能**：
- RESTful API 路由设计
- 中间件链配置
- 权限控制
- API 版本管理

### 3. pkg/ - 公共包目录 📦

**用途**：存放可被外部项目引用的公共包

**使用场景**：
- 工具函数库
- 通用组件
- SDK 包
- 第三方集成

**注意事项**：
- 保持 API 稳定性
- 良好的文档说明
- 向后兼容性

### 4. configs/ - 配置文件目录 📝

**用途**：存放各种配置文件

**常见文件类型**：
- `config.yaml`：应用主配置
- `database.yaml`：数据库配置
- `redis.yaml`：缓存配置
- `logging.yaml`：日志配置

### 5. docs/ - 项目文档目录 📚

**用途**：存放项目相关文档

**文档类型**：
- 技术文档
- API 文档
- 部署指南
- 开发规范

### 6. scripts/ - 脚本文件目录 🔧

**用途**：存放项目相关脚本

**常见脚本**：
- 数据库迁移脚本
- 部署脚本
- 构建脚本
- 测试脚本

### 7. tests/ - 测试文件目录 🧪

**用途**：存放测试相关文件

**测试类型**：
- 单元测试
- 集成测试
- 端到端测试
- 性能测试

## 分层架构设计 🏗️

### 架构层次

```
┌─────────────────────────────────────┐
│           Handlers Layer            │  ← HTTP 处理层
│        (HTTP Request/Response)      │
├─────────────────────────────────────┤
│           Services Layer            │  ← 业务逻辑层
│         (Business Logic)            │
├─────────────────────────────────────┤
│          Repository Layer           │  ← 数据访问层
│         (Data Access)               │
├─────────────────────────────────────┤
│            Models Layer             │  ← 数据模型层
│         (Data Models)               │
└─────────────────────────────────────┘
```

### 数据流向

```
HTTP Request → Middleware → Handler → Service → Repository → Database
                    ↓           ↓        ↓         ↓
                 认证/授权    参数验证   业务逻辑   数据操作
                    ↓           ↓        ↓         ↓
HTTP Response ← Response ← Result ← Business ← Data ← Database
```

## 命名规范 📝

### 文件命名
- 使用小写字母和下划线：`user_handler.go`
- 功能模块 + 层次后缀：`{module}_{layer}.go`
- 测试文件：`{filename}_test.go`

### 包命名
- 使用小写字母：`handlers`、`services`
- 简洁明了：`config`、`models`
- 避免复数形式的混淆

### 结构体命名
- 使用大驼峰命名：`UserHandler`、`PostService`
- 接口以 `er` 结尾：`UserRepository`、`PostService`

## 依赖关系 🔗

### 依赖方向
```
Handlers → Services → Repository → Models
    ↓         ↓          ↓
 Middleware  Config   Database
```

### 依赖注入
- 使用接口进行依赖抽象
- 构造函数注入依赖
- 支持单元测试和模拟

## 扩展指南 🚀

### 添加新功能模块

1. **创建模型**：在 `internal/models/` 中定义数据模型
2. **实现 Repository**：在 `internal/repository/` 中实现数据访问
3. **编写 Service**：在 `internal/services/` 中实现业务逻辑
4. **添加 Handler**：在 `internal/handlers/` 中实现 HTTP 处理
5. **配置路由**：在 `internal/routes/` 中添加路由配置
6. **编写测试**：在 `tests/` 中添加相应测试

### 最佳实践

1. **保持层次清晰**：严格按照分层架构组织代码
2. **接口编程**：使用接口定义层间契约
3. **单一职责**：每个文件和函数职责单一
4. **依赖注入**：通过构造函数注入依赖
5. **错误处理**：统一的错误处理机制
6. **文档完善**：及时更新文档和注释

---

**注意**：本文档描述的是当前项目的目录结构，随着项目发展可能会有所调整，请以实际代码结构为准。
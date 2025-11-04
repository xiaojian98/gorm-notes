# GORM 学习代码示例

本目录包含了完整的GORM学习代码示例，对应`docs`目录中的学习文档。所有代码都经过测试，可以直接运行。

## 📁 目录结构

```
code/
├── README.md                     # 本说明文件
├── 01_gorm_basics/              # GORM基础入门代码
│   ├── go.mod
│   ├── main.go                  # 基础使用示例
│   └── README.md
├── 02_gorm_advanced/            # GORM进阶特性代码
│   ├── go.mod
│   ├── main.go                  # 进阶特性示例
│   └── README.md
├── 03_blog_system/              # 博客系统实战项目
│   ├── go.mod
│   ├── main.go                  # 主程序入口
│   ├── models/
│   │   └── models.go            # 数据模型定义
│   ├── routes/
│   │   └── routes.go            # 路由配置
│   ├── handlers/
│   │   └── handlers.go          # 请求处理器
│   └── README.md
└── 04_unit_exercises/           # 单元练习代码
    ├── go.mod
    ├── level1_basic.go          # Level 1: 基础操作
    ├── level2_associations.go   # Level 2: 关联关系
    ├── level3_advanced_queries.go # Level 3: 高级查询
    ├── level4_transactions_hooks.go # Level 4: 事务和钩子
    ├── level5_performance.go    # Level 5: 性能优化
    ├── level6_comprehensive.go  # Level 6: 综合实战
    └── README.md
```

## 🎯 学习路径

### 阶段一：基础入门
📖 **对应文档：** `01_GORM基础入门_核心概念与基本操作.md`  
💻 **代码目录：** `01_gorm_basics/`

**学习内容：**
- GORM安装和配置
- 数据库连接
- 模型定义和迁移
- 基本CRUD操作
- 查询方法

### 阶段二：进阶特性
📖 **对应文档：** `01_GORM基础入门_核心概念与基本操作.md`（进阶部分）  
💻 **代码目录：** `02_gorm_advanced/`

**学习内容：**
- 关联关系（一对一、一对多、多对多）
- 预加载和延迟加载
- 钩子函数
- 事务处理
- 高级查询

### 阶段三：实战项目
📖 **对应文档：** `02_GORM背景示例_博客系统实战.md`  
💻 **代码目录：** `03_blog_system/`

**学习内容：**
- 完整的博客系统设计
- RESTful API实现
- 用户认证和授权
- 内容管理
- 数据关联和查询优化

### 阶段四：技能训练
📖 **对应文档：** `03_GORM单元练习_基础技能训练.md`  
💻 **代码目录：** `04_unit_exercises/`

**学习内容：**
- 6个级别的渐进式练习
- 从基础到高级的全面覆盖
- 性能优化和最佳实践
- 综合应用场景

## 🚀 快速开始

### 环境要求
- Go 1.21 或更高版本
- Git（用于克隆代码）

### 运行步骤

1. **进入对应目录**
   ```bash
   cd f:/Study/GO/Gorm/gorm/学习资料/code/[目录名]
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **运行代码**
   ```bash
   go run main.go  # 对于有main.go的项目
   # 或
   go run [具体文件名].go  # 对于单元练习
   ```

### 示例：运行博客系统

```bash
# 进入博客系统目录
cd f:/Study/GO/Gorm/gorm/学习资料/code/03_blog_system

# 安装依赖
go mod tidy

# 运行项目
go run main.go

# 访问API（在另一个终端）
curl http://localhost:8080/api/health
```

### 示例：运行单元练习

```bash
# 进入练习目录
cd f:/Study/GO/Gorm/gorm/学习资料/code/04_unit_exercises

# 安装依赖
go mod tidy

# 运行Level 1练习
go run level1_basic.go

# 运行Level 6综合练习
go run level6_comprehensive.go
```

## 🛠️ 技术栈

### 核心框架
- **GORM：** Go语言ORM框架
- **Gin：** HTTP Web框架（博客系统）
- **SQLite：** 轻量级数据库（练习用）
- **MySQL：** 生产级数据库（可选）

### 依赖库
- `gorm.io/gorm` - GORM核心库
- `gorm.io/driver/sqlite` - SQLite驱动
- `gorm.io/driver/mysql` - MySQL驱动
- `github.com/gin-gonic/gin` - Gin框架
- `golang.org/x/crypto` - 加密库

## 📚 学习建议

### 1. 循序渐进
- 按照阶段一到阶段四的顺序学习
- 每个阶段都要动手实践
- 理解概念后再进入下一阶段

### 2. 实践为主
- 运行所有代码示例
- 修改代码验证理解
- 尝试添加新功能

### 3. 深入理解
- 查看生成的数据库文件
- 分析SQL查询语句
- 理解性能优化原理

### 4. 扩展学习
- 尝试不同的数据库
- 集成其他Go框架
- 部署到生产环境

## 🔧 常见问题

### Q: 运行时提示找不到模块？
A: 确保在正确的目录下运行`go mod tidy`

### Q: 数据库连接失败？
A: 检查数据库配置，SQLite会自动创建文件

### Q: 端口被占用？
A: 修改代码中的端口号，或停止占用端口的程序

### Q: 依赖下载失败？
A: 配置Go代理：`go env -w GOPROXY=https://goproxy.cn,direct`

## 📊 学习进度跟踪

- [ ] 完成GORM基础入门
- [ ] 完成GORM进阶特性
- [ ] 完成博客系统实战
- [ ] 完成Level 1基础练习
- [ ] 完成Level 2关联关系练习
- [ ] 完成Level 3高级查询练习
- [ ] 完成Level 4事务和钩子练习
- [ ] 完成Level 5性能优化练习
- [ ] 完成Level 6综合实战练习

## 🎉 学习成果

完成所有代码练习后，您将掌握：

✅ **GORM核心概念** - 模型、迁移、关联关系  
✅ **数据库操作** - CRUD、查询、事务  
✅ **性能优化** - 索引、预加载、批量操作  
✅ **实战经验** - 完整项目开发  
✅ **最佳实践** - 代码组织、错误处理  

## 🤝 贡献和反馈

如果您发现代码问题或有改进建议，欢迎：
- 提交Issue报告问题
- 提交Pull Request改进代码
- 分享您的学习心得

祝您学习愉快，成为GORM专家！🚀

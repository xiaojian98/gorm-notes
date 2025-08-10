# GORM 单元练习代码

本目录包含了GORM从基础到高级的6个级别的练习代码，对应`03_GORM单元练习_基础技能训练.md`文档中的内容。

## 📁 文件结构

```
04_unit_exercises/
├── go.mod                    # Go模块依赖管理
├── README.md                 # 本说明文件
├── level1_basic.go           # Level 1: 基础操作练习
├── level2_associations.go    # Level 2: 关联关系练习
├── level3_advanced_queries.go # Level 3: 高级查询练习
├── level4_transactions_hooks.go # Level 4: 事务和钩子练习
├── level5_performance.go     # Level 5: 性能优化练习
└── level6_comprehensive.go   # Level 6: 综合实战练习
```

## 🚀 快速开始

### 1. 安装依赖

```bash
cd f:/Study/GO/Gorm/gorm/学习资料/code/04_unit_exercises
go mod tidy
```

### 2. 运行练习

每个级别的练习都可以独立运行：

```bash
# Level 1: 基础操作
go run level1_basic.go

# Level 2: 关联关系
go run level2_associations.go

# Level 3: 高级查询
go run level3_advanced_queries.go

# Level 4: 事务和钩子
go run level4_transactions_hooks.go

# Level 5: 性能优化
go run level5_performance.go

# Level 6: 综合实战
go run level6_comprehensive.go
```

## 📚 练习内容概览

### Level 1: 基础操作练习
- 数据库连接与配置
- 基本模型定义
- CRUD操作（创建、读取、更新、删除）
- 简单查询操作
- 数据库迁移和索引

**学习目标：** 掌握GORM的基本使用方法

### Level 2: 关联关系练习
- 一对一关系（User-Profile）
- 一对多关系（Category-Posts, User-Posts, Post-Comments）
- 多对多关系（Post-Tags）
- 关联数据的创建和查询
- 预加载（Preload）的使用

**学习目标：** 理解和实现各种数据关联关系

### Level 3: 高级查询练习
- 复杂条件查询和排序
- 聚合查询（COUNT, SUM, AVG等）
- 子查询的使用
- 复杂连接查询
- 窗口函数和排名
- 分页查询优化

**学习目标：** 掌握高级查询技巧和SQL优化

### Level 4: 事务和钩子练习
- 数据库事务的使用
- 钩子函数（BeforeCreate, AfterCreate等）
- 数据验证和自动处理
- 错误处理和回滚
- 批量操作

**学习目标：** 理解事务机制和模型生命周期

### Level 5: 性能优化练习
- 数据库连接池配置
- 索引创建和优化
- 查询性能测试
- N+1查询问题解决
- 批量操作优化
- 预编译语句使用

**学习目标：** 掌握数据库性能优化技巧

### Level 6: 综合实战练习
- 完整的博客系统数据模型
- 复杂业务逻辑实现
- 用户管理、内容管理、互动功能
- 通知系统和统计分析
- 高级查询和性能测试
- 真实场景的综合应用

**学习目标：** 综合运用所有技能构建完整应用

## 🛠️ 技术栈

- **Go版本：** 1.21+
- **ORM框架：** GORM v1.25.5
- **数据库：** SQLite（用于练习，易于部署）
- **加密库：** golang.org/x/crypto（用于密码哈希）

## 📝 学习建议

1. **循序渐进：** 按照Level 1到Level 6的顺序学习
2. **动手实践：** 运行每个练习，观察输出结果
3. **修改实验：** 尝试修改代码，验证自己的理解
4. **查看数据：** 使用SQLite工具查看生成的数据库文件
5. **性能分析：** 关注Level 5和Level 6中的性能测试结果

## 🔍 数据库文件

每个练习会在当前目录生成对应的SQLite数据库文件：
- `level1_basic.db`
- `level2_associations.db`
- `level3_advanced.db`
- `level4_transactions.db`
- `level5_performance.db`
- `level6_comprehensive.db`

可以使用SQLite浏览器工具查看这些数据库的结构和数据。

## 🎯 练习目标

完成所有练习后，您将能够：

✅ 熟练使用GORM进行数据库操作  
✅ 设计复杂的数据模型和关联关系  
✅ 编写高效的数据库查询  
✅ 处理事务和实现业务逻辑  
✅ 优化数据库性能  
✅ 构建完整的数据驱动应用  

## 🤝 问题反馈

如果在练习过程中遇到问题，请：
1. 检查Go版本和依赖是否正确安装
2. 查看控制台输出的错误信息
3. 参考对应的文档说明
4. 尝试重新运行`go mod tidy`

祝您学习愉快！🎉
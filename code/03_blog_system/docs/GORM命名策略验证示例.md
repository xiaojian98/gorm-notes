# GORM命名策略验证示例 🧪

## 验证目的 🎯

通过实际代码验证不同命名策略对数据库字段名的影响，帮助理解`03_blog_system`项目中字段命名的真实原因。

## 测试代码 💻

### 1. 基础测试结构体

```go
package main

import (
    "fmt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

// TestModel 测试模型
type TestModel struct {
    ID          uint   `gorm:"primarykey"`
    UserName    string `gorm:"size:50"`
    EmailAddr   string `gorm:"size:100"`
    PhoneNumber string `gorm:"size:20"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 2. 默认命名策略测试

```go
// testDefaultNaming 测试GORM默认命名策略
func testDefaultNaming() {
    fmt.Println("=== 测试默认命名策略 ===")
    
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        // 使用默认命名策略
    })
    if err != nil {
        panic(err)
    }
    
    // 自动迁移
    db.AutoMigrate(&TestModel{})
    
    // 查看生成的表结构
    var columns []struct {
        Name string `gorm:"column:name"`
        Type string `gorm:"column:type"`
    }
    
    db.Raw("PRAGMA table_info(test_models)").Scan(&columns)
    
    fmt.Println("生成的字段名：")
    for _, col := range columns {
        fmt.Printf("- %s (%s)\n", col.Name, col.Type)
    }
}
```

### 3. 自定义命名策略测试

```go
// CustomNamingStrategy 自定义命名策略（模拟blog_system项目）
type CustomNamingStrategy struct{}

func (ns *CustomNamingStrategy) TableName(table string) string {
    return table
}

func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
    return column // 直接返回原字段名
}

func (ns *CustomNamingStrategy) JoinTableName(str string) string {
    return str
}

func (ns *CustomNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
    return rel.Name + "_id"
}

func (ns *CustomNamingStrategy) CheckerName(table, column string) string {
    return "chk_" + table + "_" + column
}

func (ns *CustomNamingStrategy) IndexName(table, column string) string {
    return "idx_" + table + "_" + column
}

func (ns *CustomNamingStrategy) UniqueName(table, column string) string {
    return "uniq_" + table + "_" + column
}

// testCustomNaming 测试自定义命名策略
func testCustomNaming() {
    fmt.Println("\n=== 测试自定义命名策略 ===")
    
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        NamingStrategy: &CustomNamingStrategy{}, // 使用自定义命名策略
    })
    if err != nil {
        panic(err)
    }
    
    // 自动迁移
    db.AutoMigrate(&TestModel{})
    
    // 查看生成的表结构
    var columns []struct {
        Name string `gorm:"column:name"`
        Type string `gorm:"column:type"`
    }
    
    db.Raw("PRAGMA table_info(TestModel)").Scan(&columns)
    
    fmt.Println("生成的字段名：")
    for _, col := range columns {
        fmt.Printf("- %s (%s)\n", col.Name, col.Type)
    }
}
```

### 4. 完整测试程序

```go
package main

import (
    "fmt"
    "time"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

func main() {
    fmt.Println("🧪 GORM命名策略验证测试")
    
    // 测试默认命名策略
    testDefaultNaming()
    
    // 测试自定义命名策略
    testCustomNaming()
    
    // 对比分析
    fmt.Println("\n=== 对比分析 ===")
    fmt.Println("默认策略：Go字段名 -> 蛇形命名")
    fmt.Println("- UserName -> user_name")
    fmt.Println("- EmailAddr -> email_addr")
    fmt.Println("- PhoneNumber -> phone_number")
    fmt.Println("")
    fmt.Println("自定义策略：Go字段名 -> 原字段名")
    fmt.Println("- UserName -> UserName")
    fmt.Println("- EmailAddr -> EmailAddr")
    fmt.Println("- PhoneNumber -> PhoneNumber")
}
```

## 预期结果 📊

### 默认命名策略输出：
```
=== 测试默认命名策略 ===
生成的字段名：
- id (INTEGER)
- user_name (TEXT)
- email_addr (TEXT)
- phone_number (TEXT)
- created_at (DATETIME)
- updated_at (DATETIME)
```

### 自定义命名策略输出：
```
=== 测试自定义命名策略 ===
生成的字段名：
- ID (INTEGER)
- UserName (TEXT)
- EmailAddr (TEXT)
- PhoneNumber (TEXT)
- CreatedAt (DATETIME)
- UpdatedAt (DATETIME)
```

## 运行测试 🚀

### 1. 创建测试文件

```bash
# 在项目根目录创建测试文件
touch naming_strategy_test.go
```

### 2. 运行测试

```bash
# 运行测试程序
go run naming_strategy_test.go
```

### 3. 查看结果

观察输出结果，验证不同命名策略的实际效果。

## 实际验证blog_system项目 🔍

### 1. 查看当前数据库结构

```sql
-- 连接到blog_system数据库
mysql -h 10.6.2.7 -u root -p blog_system

-- 查看User表结构
DESC User;

-- 查看建表语句
SHOW CREATE TABLE User;
```

### 2. 预期结果

```sql
-- 应该看到字段名为：
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| ID          | bigint(20)   | NO   | PRI | NULL    | auto_increment |
| Username    | varchar(50)  | NO   | UNI | NULL    |                |
| Email       | varchar(100) | NO   | UNI | NULL    |                |
| Password    | varchar(255) | NO   |     | NULL    |                |
| Nickname    | varchar(50)  | YES  |     | NULL    |                |
| Avatar      | varchar(255) | YES  |     | NULL    |                |
| Status      | varchar(20)  | YES  |     | NULL    |                |
| LastLoginAt | datetime(3)  | YES  |     | NULL    |                |
| LoginCount  | bigint(20)   | YES  |     | 0       |                |
| CreatedAt   | datetime(3)  | YES  |     | NULL    |                |
| UpdatedAt   | datetime(3)  | YES  |     | NULL    |                |
| DeletedAt   | datetime(3)  | YES  | MUL | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
```

## 学习要点 📚

### 1. 命名策略的重要性
- 影响数据库表结构的生成
- 决定字段名的最终形式
- 需要在项目初期确定并保持一致

### 2. 自定义vs默认
- **默认策略**：符合数据库命名规范（蛇形命名）
- **自定义策略**：保持Go代码风格的一致性
- **选择依据**：团队约定、项目需求、数据库规范

### 3. 实践建议
- 🔍 **仔细检查**：项目配置中的命名策略设置
- 📝 **文档记录**：明确记录项目使用的命名策略
- 🧪 **验证测试**：通过实际测试验证预期行为
- 🔄 **保持一致**：在整个项目中使用统一的命名策略

## 总结 🎯

通过这个验证示例，我们可以：

1. **理解原理**：不同命名策略如何影响数据库字段名
2. **验证假设**：通过实际代码测试验证理论
3. **解决疑问**：明确`03_blog_system`项目字段命名的真实原因
4. **指导实践**：为将来的项目选择合适的命名策略

这个发现再次证明了**实践验证**的重要性！ 🙏
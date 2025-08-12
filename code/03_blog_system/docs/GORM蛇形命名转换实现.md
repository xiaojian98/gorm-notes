# GORM蛇形命名转换实现指南 🐍

## 问题背景 📋

在`03_blog_system`项目中，当前的`CustomNamingStrategy`直接返回原字段名：

```go
// 当前实现 - 直接返回原字段名
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
    return column  // UserName -> UserName
}
```

如果需要实现蛇形命名转换，需要修改这个方法。

## 蛇形命名转换实现 🔧

### 1. 简单版本实现

```go
import (
    "regexp"
    "strings"
)

// ColumnName 列名命名策略 - 蛇形转换版本
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
    return toSnakeCase(column)
}

// toSnakeCase 将驼峰命名转换为蛇形命名
// 参数: str - 需要转换的字符串
// 返回值: 转换后的蛇形命名字符串
func toSnakeCase(str string) string {
    // 处理连续大写字母的情况，如 "HTTPSProxy" -> "HTTPS_Proxy"
    re1 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
    str = re1.ReplaceAllString(str, "${1}_${2}")
    
    // 处理普通驼峰命名，如 "UserName" -> "User_Name"
    re2 := regexp.MustCompile(`([a-z\d])([A-Z])`)
    str = re2.ReplaceAllString(str, "${1}_${2}")
    
    // 转换为小写
    return strings.ToLower(str)
}
```

### 2. 完整版本实现（推荐）

```go
import (
    "regexp"
    "strings"
    "unicode"
)

// ColumnName 列名命名策略 - 完整蛇形转换版本
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
    return toSnakeCaseAdvanced(column)
}

// toSnakeCaseAdvanced 高级蛇形命名转换
// 支持更复杂的命名场景
// 参数: str - 需要转换的字符串
// 返回值: 转换后的蛇形命名字符串
func toSnakeCaseAdvanced(str string) string {
    if str == "" {
        return ""
    }
    
    // 特殊字段名处理
    specialCases := map[string]string{
        "ID":       "id",
        "URL":      "url",
        "HTTP":     "http",
        "HTTPS":    "https",
        "API":      "api",
        "JSON":     "json",
        "XML":      "xml",
        "UUID":     "uuid",
        "SQL":      "sql",
        "HTML":     "html",
        "CSS":      "css",
        "JS":       "js",
        "OAuth":    "oauth",
        "JWT":      "jwt",
    }
    
    // 检查是否为特殊情况
    if snake, exists := specialCases[str]; exists {
        return snake
    }
    
    var result strings.Builder
    
    for i, r := range str {
        if i > 0 && unicode.IsUpper(r) {
            // 检查前一个字符
            prevRune := rune(str[i-1])
            
            // 如果前一个字符是小写字母或数字，添加下划线
            if unicode.IsLower(prevRune) || unicode.IsDigit(prevRune) {
                result.WriteRune('_')
            }
            
            // 检查连续大写字母的情况
            if i < len(str)-1 {
                nextRune := rune(str[i+1])
                if unicode.IsUpper(prevRune) && unicode.IsLower(nextRune) {
                    result.WriteRune('_')
                }
            }
        }
        
        result.WriteRune(unicode.ToLower(r))
    }
    
    return result.String()
}
```

### 3. 使用GORM内置转换（最简单）

```go
import (
    "gorm.io/gorm/schema"
)

// ColumnName 列名命名策略 - 使用GORM内置转换
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
    // 使用GORM默认的命名策略进行转换
    defaultNaming := schema.NamingStrategy{}
    return defaultNaming.ColumnName(table, column)
}
```

## 转换效果对比 📊

### 测试用例

```go
// 测试不同命名转换的效果
func testNamingConversion() {
    testCases := []string{
        "ID",
        "UserName",
        "EmailAddr",
        "PhoneNumber",
        "CreatedAt",
        "UpdatedAt",
        "LastLoginAt",
        "LoginCount",
        "HTTPSProxy",
        "APIKey",
        "JSONData",
        "XMLParser",
        "UUIDGenerator",
    }
    
    fmt.Println("字段名转换对比：")
    fmt.Println("原字段名 -> 蛇形命名")
    fmt.Println("========================")
    
    for _, testCase := range testCases {
        snake := toSnakeCaseAdvanced(testCase)
        fmt.Printf("%-15s -> %s\n", testCase, snake)
    }
}
```

### 预期输出

```
字段名转换对比：
原字段名 -> 蛇形命名
========================
ID              -> id
UserName        -> user_name
EmailAddr       -> email_addr
PhoneNumber     -> phone_number
CreatedAt       -> created_at
UpdatedAt       -> updated_at
LastLoginAt     -> last_login_at
LoginCount      -> login_count
HTTPSProxy      -> https_proxy
APIKey          -> api_key
JSONData        -> json_data
XMLParser       -> xml_parser
UUIDGenerator   -> uuid_generator
```

## 完整的修改方案 🔄

### 1. 修改database.go文件

```go
// config/database.go

import (
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
    "time"
    "unicode"
    
    _ "github.com/go-sql-driver/mysql"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

// CustomNamingStrategy 自定义命名策略（蛇形版本）
type CustomNamingStrategy struct{}

// SchemaName 数据库模式命名策略
func (ns *CustomNamingStrategy) SchemaName(table string) string {
    return table
}

// TableName 表名命名策略
func (ns *CustomNamingStrategy) TableName(table string) string {
    return toSnakeCase(table)
}

// ColumnName 列名命名策略 - 蛇形转换
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
    return toSnakeCase(column)
}

// JoinTableName 连接表命名策略
func (ns *CustomNamingStrategy) JoinTableName(str string) string {
    return toSnakeCase(str)
}

// RelationshipFKName 外键命名策略
func (ns *CustomNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
    return toSnakeCase(rel.Name) + "_id"
}

// CheckerName 检查器命名策略
func (ns *CustomNamingStrategy) CheckerName(table, column string) string {
    return "chk_" + toSnakeCase(table) + "_" + toSnakeCase(column)
}

// IndexName 索引命名策略
func (ns *CustomNamingStrategy) IndexName(table, column string) string {
    return "idx_" + toSnakeCase(table) + "_" + toSnakeCase(column)
}

// UniqueName 唯一约束命名策略
func (ns *CustomNamingStrategy) UniqueName(table, column string) string {
    return "uniq_" + toSnakeCase(table) + "_" + toSnakeCase(column)
}

// toSnakeCase 将驼峰命名转换为蛇形命名
func toSnakeCase(str string) string {
    if str == "" {
        return ""
    }
    
    // 特殊字段名处理
    if str == "ID" {
        return "id"
    }
    
    // 处理连续大写字母
    re1 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
    str = re1.ReplaceAllString(str, "${1}_${2}")
    
    // 处理普通驼峰命名
    re2 := regexp.MustCompile(`([a-z\d])([A-Z])`)
    str = re2.ReplaceAllString(str, "${1}_${2}")
    
    return strings.ToLower(str)
}
```

### 2. 数据库迁移注意事项

```go
// 如果已有数据，需要创建迁移脚本
func migration004Up(db *gorm.DB) error {
    // 重命名现有字段
    migrations := []struct {
        table    string
        oldName  string
        newName  string
        dataType string
    }{
        {"User", "Username", "username", "VARCHAR(50)"},
        {"User", "Email", "email", "VARCHAR(100)"},
        {"User", "Password", "password", "VARCHAR(255)"},
        {"User", "Nickname", "nickname", "VARCHAR(50)"},
        {"User", "Avatar", "avatar", "VARCHAR(255)"},
        {"User", "Status", "status", "VARCHAR(20)"},
        {"User", "LastLoginAt", "last_login_at", "DATETIME(3)"},
        {"User", "LoginCount", "login_count", "BIGINT"},
        {"User", "CreatedAt", "created_at", "DATETIME(3)"},
        {"User", "UpdatedAt", "updated_at", "DATETIME(3)"},
        {"User", "DeletedAt", "deleted_at", "DATETIME(3)"},
    }
    
    for _, m := range migrations {
        sql := fmt.Sprintf("ALTER TABLE %s CHANGE %s %s %s", 
            m.table, m.oldName, m.newName, m.dataType)
        if err := db.Exec(sql).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

## 选择建议 💡

### 1. 新项目推荐
- 使用**GORM内置转换**（方案3）
- 简单、可靠、与GORM生态一致

### 2. 现有项目迁移
- 评估数据迁移成本
- 考虑API兼容性影响
- 可能需要分阶段迁移

### 3. 混合方案
- 新表使用蛇形命名
- 旧表保持原有命名
- 通过`gorm:"column:xxx"`标签指定

## 实施步骤 🚀

### 1. 备份数据
```bash
mysqldump -h 10.6.2.7 -u root -p blog_system > backup.sql
```

### 2. 修改代码
- 更新`CustomNamingStrategy`
- 添加转换函数
- 测试转换效果

### 3. 创建迁移
- 编写字段重命名迁移
- 测试迁移脚本
- 验证数据完整性

### 4. 更新应用
- 重新部署应用
- 验证功能正常
- 监控错误日志

## 总结 📝

蛇形命名转换的实现方式：

1. **简单正则版本** - 适合基本需求
2. **高级版本** - 处理复杂命名场景
3. **GORM内置** - 最推荐的方案

选择合适的方案需要考虑：
- 项目阶段（新建 vs 迁移）
- 数据量大小
- 团队技术水平
- 维护成本

记住：**一致性比完美更重要！** 🎯✨
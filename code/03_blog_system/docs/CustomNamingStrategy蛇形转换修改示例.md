# CustomNamingStrategy 蛇形转换修改示例 🐍

## 当前代码（第157-160行）📍

```go
// ColumnName 列名命名策略
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
	return column
}
```

## 修改方案 🔧

### 方案一：简单蛇形转换（推荐新手）

```go
// ColumnName 列名命名策略 - 蛇形转换版本
// 参数: table - 表名, column - 列名
// 返回值: 转换后的蛇形命名列名
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
	return toSnakeCase(column)
}

// toSnakeCase 将驼峰命名转换为蛇形命名
// 参数: str - 需要转换的字符串
// 返回值: 转换后的蛇形命名字符串
func toSnakeCase(str string) string {
	if str == "" {
		return ""
	}
	
	// 特殊处理ID字段
	if str == "ID" {
		return "id"
	}
	
	// 使用正则表达式进行转换
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}
```

### 方案二：使用GORM内置转换（最推荐）⭐

```go
import (
	"gorm.io/gorm/schema"
)

// ColumnName 列名命名策略 - 使用GORM默认蛇形转换
// 参数: table - 表名, column - 列名  
// 返回值: 转换后的蛇形命名列名
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
	// 使用GORM默认的命名策略
	defaultNaming := schema.NamingStrategy{}
	return defaultNaming.ColumnName(table, column)
}
```

### 方案三：高级蛇形转换（处理复杂情况）

```go
import (
	"regexp"
	"strings"
)

// ColumnName 列名命名策略 - 高级蛇形转换
// 参数: table - 表名, column - 列名
// 返回值: 转换后的蛇形命名列名
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
	return toSnakeCaseAdvanced(column)
}

// toSnakeCaseAdvanced 高级蛇形命名转换
// 支持复杂的命名场景，如连续大写字母
// 参数: str - 需要转换的字符串
// 返回值: 转换后的蛇形命名字符串
func toSnakeCaseAdvanced(str string) string {
	if str == "" {
		return ""
	}
	
	// 特殊字段名映射
	specialCases := map[string]string{
		"ID":       "id",
		"URL":      "url", 
		"HTTP":     "http",
		"HTTPS":    "https",
		"API":      "api",
		"JSON":     "json",
		"XML":      "xml",
		"UUID":     "uuid",
	}
	
	// 检查特殊情况
	if snake, exists := specialCases[str]; exists {
		return snake
	}
	
	// 处理连续大写字母，如 "HTTPSProxy" -> "HTTPS_Proxy"
	re1 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	str = re1.ReplaceAllString(str, "${1}_${2}")
	
	// 处理普通驼峰命名，如 "UserName" -> "User_Name"
	re2 := regexp.MustCompile(`([a-z\d])([A-Z])`)
	str = re2.ReplaceAllString(str, "${1}_${2}")
	
	// 转换为小写
	return strings.ToLower(str)
}
```

## 完整修改步骤 📝

### 1. 添加必要的import

在文件顶部添加所需的包导入：

```go
import (
	// ... 现有的导入
	"regexp"    // 方案一和三需要
	"strings"   // 方案一和三需要
	"gorm.io/gorm/schema"  // 方案二需要
)
```

### 2. 替换ColumnName方法

将第157-160行的代码替换为上述任一方案。

### 3. 添加辅助函数

如果选择方案一或三，需要在文件末尾添加相应的辅助函数。

## 转换效果对比 📊

| 原字段名 | 当前输出 | 蛇形转换后 |
|---------|---------|----------|
| ID | ID | id |
| UserName | UserName | user_name |
| EmailAddr | EmailAddr | email_addr |
| CreatedAt | CreatedAt | created_at |
| LastLoginAt | LastLoginAt | last_login_at |
| LoginCount | LoginCount | login_count |
| HTTPSProxy | HTTPSProxy | https_proxy |
| APIKey | APIKey | api_key |

## 推荐选择 💡

### 🥇 方案二（GORM内置）
- **优点**: 简单可靠，与GORM生态一致
- **缺点**: 需要额外导入
- **适用**: 新项目或完全迁移

### 🥈 方案一（简单转换）
- **优点**: 代码简洁，易于理解
- **缺点**: 可能无法处理复杂命名
- **适用**: 简单项目，命名规范统一

### 🥉 方案三（高级转换）
- **优点**: 处理复杂命名场景
- **缺点**: 代码较复杂
- **适用**: 有复杂命名需求的项目

## 注意事项 ⚠️

### 1. 数据库迁移
如果数据库中已有数据，修改命名策略后需要：
- 备份现有数据
- 创建字段重命名迁移
- 测试迁移脚本

### 2. API兼容性
字段名变化可能影响：
- JSON序列化/反序列化
- 前端接口调用
- 第三方集成

### 3. 测试验证
修改后需要：
- 运行单元测试
- 验证数据库表结构
- 检查API响应格式

## 快速实施 🚀

### 最简单的修改（推荐）

直接将第157-160行替换为：

```go
// ColumnName 列名命名策略 - 蛇形转换
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
	// 使用GORM默认的蛇形命名策略
	defaultNaming := schema.NamingStrategy{}
	return defaultNaming.ColumnName(table, column)
}
```

然后在import部分添加：
```go
"gorm.io/gorm/schema"
```

这样就完成了蛇形命名的转换！ 🎉✨
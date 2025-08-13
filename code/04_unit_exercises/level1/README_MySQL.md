# GORM Level 1 MySQL 支持说明

本文档说明如何在 GORM Level 1 基础练习中使用 MySQL 数据库。

## 功能特性

- ✅ 支持 SQLite 和 MySQL 双数据库
- ✅ 统一的数据库配置接口
- ✅ 自动选择数据库驱动
- ✅ 独立的连接池配置
- ✅ 完整的演示示例

## 环境要求

### 1. Go 依赖

确保 `go.mod` 文件包含以下依赖：

```go
require (
    gorm.io/driver/mysql v1.5.2
    gorm.io/driver/sqlite v1.5.4
    gorm.io/gorm v1.25.5
)
```

### 2. MySQL 服务

- MySQL 5.7+ 或 MySQL 8.0+
- 创建测试数据库：`gorm_test`
- 确保 MySQL 服务正在运行

## 安装和配置

### 1. 安装 MySQL 依赖

```bash
go mod tidy
```

### 2. 配置 MySQL 数据库

#### 方法一：使用 MySQL 命令行

```sql
-- 连接到 MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE gorm_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（可选）
CREATE USER 'gorm_user'@'localhost' IDENTIFIED BY 'gorm_password';
GRANT ALL PRIVILEGES ON gorm_test.* TO 'gorm_user'@'localhost';
FLUSH PRIVILEGES;
```

#### 方法二：使用现有用户

如果使用 root 用户，确保密码正确，并且 `gorm_test` 数据库已创建。

### 3. 修改连接配置

在 `DemoMySQL()` 函数中修改 DSN 连接字符串：

```go
// 默认配置（使用 root 用户）
mysqlDSN := "root:your_password@tcp(localhost:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

// 或使用自定义用户
mysqlDSN := "gorm_user:gorm_password@tcp(localhost:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
```

## 使用方法

### 1. 基本使用

```go
// 获取 MySQL 配置
mysqlConfig := GetMySQLConfig("root:password@tcp(localhost:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local")

// 初始化数据库连接
db, err := InitDatabase(mysqlConfig)
if err != nil {
    log.Fatal("MySQL连接失败:", err)
}

// 自动迁移
err = AutoMigrate(db)
if err != nil {
    log.Fatal("数据库迁移失败:", err)
}
```

### 2. 配置参数说明

| 参数 | SQLite 默认值 | MySQL 默认值 | 说明 |
|------|---------------|--------------|------|
| MaxOpenConns | 10 | 20 | 最大打开连接数 |
| MaxIdleConns | 5 | 10 | 最大空闲连接数 |
| MaxLifetime | 1小时 | 1小时 | 连接最大生命周期 |
| LogLevel | Info | Info | 日志级别 |

### 3. DSN 连接字符串格式

```
[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
```

常用参数：
- `charset=utf8mb4`: 字符集
- `parseTime=True`: 解析时间类型
- `loc=Local`: 时区设置
- `timeout=30s`: 连接超时
- `readTimeout=30s`: 读取超时
- `writeTimeout=30s`: 写入超时

## 运行示例

### 1. 运行完整示例

```bash
cd f:\Study\GO\Gorm\gorm\gorm-note\code\04_unit_exercises\level1
go run level1_basic.go
```

### 2. 预期输出

```
=== GORM Level 1 基础练习 ===
✓ 数据库连接成功
✓ 数据库迁移完成
✓ 索引创建完成

=== CRUD操作演示 ===
创建用户成功，ID: 1, 影响行数: 1
查询到用户: &{BaseModel:{ID:1 CreatedAt:... UpdatedAt:... DeletedAt:{...}} Username:alice Email:alice@example.com Password: Age:25 IsActive:true}
更新用户成功，影响行数: 1

=== 查询操作演示 ===
...

=== Level 1 SQLite练习完成 ===

=== MySQL 数据库连接演示 ===
✓ MySQL数据库连接成功
✓ MySQL数据库迁移完成
✓ MySQL创建用户成功: &{...}

=== MySQL连接池统计 ===
max_open_connections: 20
open_connections: 1
...

=== Level 1 全部练习完成 ===
```

## 故障排除

### 1. 连接失败

**错误信息：** `MySQL连接失败: dial tcp [::1]:3306: connect: connection refused`

**解决方案：**
- 确保 MySQL 服务正在运行
- 检查端口 3306 是否被占用
- 尝试使用 `127.0.0.1` 替代 `localhost`

### 2. 认证失败

**错误信息：** `Access denied for user 'root'@'localhost'`

**解决方案：**
- 检查用户名和密码是否正确
- 确保用户有访问数据库的权限
- 重置 MySQL root 密码

### 3. 数据库不存在

**错误信息：** `Unknown database 'gorm_test'`

**解决方案：**
- 手动创建数据库：`CREATE DATABASE gorm_test;`
- 检查数据库名称拼写

### 4. 字符集问题

**错误信息：** 中文字符显示异常

**解决方案：**
- 确保数据库字符集为 `utf8mb4`
- DSN 中包含 `charset=utf8mb4`
- 创建数据库时指定字符集

## 扩展功能

### 1. 添加其他数据库支持

可以按照相同模式添加 PostgreSQL、SQL Server 等数据库支持：

```go
const (
    SQLite     DatabaseType = "sqlite"
    MySQL      DatabaseType = "mysql"
    PostgreSQL DatabaseType = "postgresql"
)

// 在 InitDatabase 函数中添加
case PostgreSQL:
    dialector = postgres.Open(config.DSN)
```

### 2. 环境变量配置

```go
import "os"

func GetMySQLConfigFromEnv() *DatabaseConfig {
    dsn := os.Getenv("MYSQL_DSN")
    if dsn == "" {
        dsn = "root:password@tcp(localhost:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
    }
    return GetMySQLConfig(dsn)
}
```

### 3. 配置文件支持

可以使用 YAML、JSON 或 TOML 配置文件来管理数据库配置。

## 最佳实践

1. **生产环境**：使用环境变量或配置文件管理敏感信息
2. **连接池**：根据应用负载调整连接池参数
3. **监控**：定期检查连接池统计信息
4. **安全**：使用专用数据库用户，避免使用 root
5. **备份**：定期备份重要数据

## 参考资料

- [GORM 官方文档](https://gorm.io/docs/)
- [MySQL 驱动文档](https://github.com/go-sql-driver/mysql)
- [GORM MySQL 驱动](https://github.com/go-gorm/mysql)
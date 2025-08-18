# GORM Level 6 故障排查指南

## 🚨 常见问题快速索引

| 问题类型 | 关键词 | 跳转链接 |
|---------|--------|----------|
| 编译错误 | `import`, `undefined` | [编译问题](#编译问题) |
| 数据库连接 | `connection`, `driver` | [连接问题](#数据库连接问题) |
| 模型定义 | `struct`, `tag`, `association` | [模型问题](#模型定义问题) |
| 查询错误 | `query`, `where`, `join` | [查询问题](#查询问题) |
| 事务问题 | `transaction`, `rollback` | [事务问题](#事务问题) |
| 性能问题 | `slow`, `performance`, `index` | [性能问题](#性能问题) |
| 迁移问题 | `migration`, `table`, `column` | [迁移问题](#迁移问题) |

---

## 📋 问题诊断流程

```
遇到问题
    ↓
查看错误信息
    ↓
确定问题类型
    ↓
查找对应解决方案
    ↓
应用解决方案
    ↓
验证问题是否解决
    ↓
记录解决过程
```

---

## 🔧 编译问题

### 问题1：导入包错误

**错误信息**：
```
imported and not used: "encoding/json"
imported and not used: "strconv"
```

**原因分析**：
- Go语言不允许导入未使用的包
- 代码中导入了包但没有实际使用

**解决方案**：
```go
// 错误示例
import (
    "encoding/json"  // 导入但未使用
    "fmt"
    "gorm.io/gorm"
)

// 正确示例
import (
    "fmt"
    "gorm.io/gorm"
    // 只导入实际使用的包
)

// 或者使用空白标识符（如果确实需要包的副作用）
import (
    _ "encoding/json"  // 使用空白标识符
    "fmt"
    "gorm.io/gorm"
)
```

**预防措施**：
- 使用IDE的自动导入功能
- 定期清理未使用的导入
- 使用 `goimports` 工具自动管理导入

### 问题2：未定义的标识符

**错误信息**：
```
undefined: User
undefined: gorm.Model
```

**原因分析**：
- 结构体或变量未定义
- 包导入路径错误
- 作用域问题

**解决方案**：
```go
// 确保正确导入GORM
import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

// 确保结构体定义在正确位置
type User struct {
    gorm.Model  // 确保gorm包已导入
    Username string
    Email    string
}

// 确保在正确的作用域中使用
func main() {
    var user User  // User必须在此作用域中可见
}
```

---

## 🔌 数据库连接问题

### 问题1：SQLite数据库文件权限错误

**错误信息**：
```
unable to open database file: permission denied
```

**原因分析**：
- 数据库文件所在目录没有写权限
- 文件被其他进程占用
- 路径不存在

**解决方案**：
```go
// 检查并创建目录
func ensureDBDir(dbPath string) error {
    dir := filepath.Dir(dbPath)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        return os.MkdirAll(dir, 0755)
    }
    return nil
}

// 使用绝对路径
func initSQLite() *gorm.DB {
    dbPath := "./data/blog.db"
    
    // 确保目录存在
    if err := ensureDBDir(dbPath); err != nil {
        log.Fatal("创建数据库目录失败:", err)
    }
    
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        log.Fatal("连接数据库失败:", err)
    }
    
    return db
}
```

### 问题2：MySQL连接失败

**错误信息**：
```
Error 1045: Access denied for user 'root'@'localhost'
Error 2003: Can't connect to MySQL server
```

**原因分析**：
- 用户名或密码错误
- MySQL服务未启动
- 网络连接问题
- 防火墙阻止连接

**解决方案**：
```go
// 完整的MySQL连接配置
func connectMySQL() *gorm.DB {
    config := DatabaseConfig{
        Type:     "mysql",
        Host:     "localhost",
        Port:     3306,
        Username: "root",
        Password: "your_password",
        Database: "blog_db",
    }
    
    // 构建DSN时添加更多参数
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.Database,
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        log.Printf("MySQL连接失败: %v", err)
        log.Printf("请检查: 1)MySQL服务是否启动 2)用户名密码是否正确 3)数据库是否存在")
        return nil
    }
    
    // 测试连接
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Printf("数据库ping失败: %v", err)
        return nil
    }
    
    log.Println("MySQL连接成功")
    return db
}
```

**诊断步骤**：
```bash
# 1. 检查MySQL服务状态
net start mysql  # Windows
sudo systemctl status mysql  # Linux

# 2. 测试连接
mysql -u root -p -h localhost

# 3. 检查用户权限
SHOW GRANTS FOR 'root'@'localhost';

# 4. 创建数据库
CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

---

## 📊 模型定义问题

### 问题1：关联关系错误

**错误信息**：
```
invalid association
foreign key not found
```

**原因分析**：
- 外键字段名不正确
- 关联标签配置错误
- 结构体字段类型不匹配

**解决方案**：
```go
// 错误示例
type User struct {
    ID    uint
    Posts []Post `gorm:"foreignKey:UserID"`  // 外键字段名错误
}

type Post struct {
    ID     uint
    UserId uint  // 字段名不规范
}

// 正确示例
type User struct {
    ID    uint
    Posts []Post `gorm:"foreignKey:UserID"`  // 正确的外键配置
}

type Post struct {
    ID     uint
    UserID uint  // 正确的字段名（驼峰命名）
    User   User  // 反向关联
}

// 复杂关联示例
type Post struct {
    ID       uint
    UserID   uint
    User     User     `gorm:"foreignKey:UserID"`
    Tags     []Tag    `gorm:"many2many:post_tags;"`
    Comments []Comment `gorm:"foreignKey:PostID"`
}

type Tag struct {
    ID    uint
    Posts []Post `gorm:"many2many:post_tags;"`
}
```

### 问题2：GORM标签配置错误

**错误信息**：
```
invalid tag format
unknown column type
```

**原因分析**：
- 标签语法错误
- 不支持的数据类型
- 标签参数错误

**解决方案**：
```go
// 错误示例
type User struct {
    Username string `gorm:"unique;size:50;not null"`  // 缺少Index关键字
    Email    string `gorm:"unique_index;size:100"`    // 旧版本语法
    Status   string `gorm:"default:'active'"`         // 引号使用错误
}

// 正确示例
type User struct {
    Username string `gorm:"uniqueIndex;size:50;not null;comment:用户名"`
    Email    string `gorm:"uniqueIndex;size:100;not null;comment:邮箱"`
    Status   string `gorm:"size:20;not null;default:active;comment:状态"`
    Age      int    `gorm:"check:age >= 0;comment:年龄"`
    Score    float64 `gorm:"precision:10;scale:2;comment:评分"`
}

// 常用标签参考
type ModelExample struct {
    // 主键
    ID uint `gorm:"primaryKey;autoIncrement;comment:主键ID"`
    
    // 字符串字段
    Name string `gorm:"size:100;not null;uniqueIndex;comment:名称"`
    
    // 数值字段
    Price decimal.Decimal `gorm:"type:decimal(10,2);not null;comment:价格"`
    
    // 时间字段
    CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间"`
    UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间"`
    
    // 软删除
    DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
    
    // JSON字段
    Metadata datatypes.JSON `gorm:"type:json;comment:元数据"`
    
    // 枚举字段
    Status string `gorm:"type:enum('active','inactive','banned');default:active;comment:状态"`
}
```

---

## 🔍 查询问题

### 问题1：查询结果为空

**错误信息**：
```
record not found
no rows in result set
```

**原因分析**：
- 查询条件错误
- 数据不存在
- 软删除记录被过滤
- 表名或字段名错误

**解决方案**：
```go
// 调试查询问题的方法
func debugQuery(db *gorm.DB) {
    // 1. 启用详细日志
    db = db.Debug()  // 打印SQL语句
    
    // 2. 检查查询条件
    var user User
    result := db.Where("username = ?", "testuser").First(&user)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            fmt.Println("记录不存在")
            
            // 检查是否有类似记录
            var count int64
            db.Model(&User{}).Where("username LIKE ?", "%test%").Count(&count)
            fmt.Printf("类似记录数量: %d\n", count)
            
            // 检查是否被软删除
            db.Unscoped().Where("username = ?", "testuser").First(&user)
            if user.ID > 0 {
                fmt.Println("记录已被软删除")
            }
        } else {
            fmt.Printf("查询错误: %v\n", result.Error)
        }
    }
    
    // 3. 验证表结构
    if db.Migrator().HasTable(&User{}) {
        fmt.Println("User表存在")
    } else {
        fmt.Println("User表不存在")
    }
    
    // 4. 检查字段是否存在
    if db.Migrator().HasColumn(&User{}, "username") {
        fmt.Println("username字段存在")
    } else {
        fmt.Println("username字段不存在")
    }
}
```

### 问题2：关联查询失败

**错误信息**：
```
association not found
invalid association
```

**原因分析**：
- 关联关系配置错误
- 外键值为空
- Preload路径错误

**解决方案**：
```go
// 正确的关联查询方式
func correctAssociationQuery(db *gorm.DB) {
    // 1. 简单预加载
    var users []User
    db.Preload("Posts").Find(&users)
    
    // 2. 嵌套预加载
    db.Preload("Posts.Comments").Find(&users)
    
    // 3. 条件预加载
    db.Preload("Posts", "status = ?", "published").Find(&users)
    
    // 4. 自定义预加载
    db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
        return db.Order("created_at DESC").Limit(5)
    }).Find(&users)
    
    // 5. 选择性预加载
    db.Preload("Posts", "id IN ?", []uint{1, 2, 3}).Find(&users)
    
    // 6. 手动关联查询
    var user User
    db.First(&user, 1)
    
    var posts []Post
    db.Model(&user).Association("Posts").Find(&posts)
    
    // 7. 连接查询
    var results []struct {
        User User
        Post Post
    }
    
    db.Table("user").
        Select("user.*, post.*").
        Joins("LEFT JOIN post ON user.id = post.user_id").
        Where("user.status = ?", "active").
        Scan(&results)
}
```

---

## 💾 事务问题

### 问题1：事务死锁

**错误信息**：
```
Deadlock found when trying to get lock
Lock wait timeout exceeded
```

**原因分析**：
- 多个事务相互等待锁
- 事务持有时间过长
- 锁的获取顺序不一致

**解决方案**：
```go
// 避免死锁的最佳实践
func avoidDeadlock(db *gorm.DB) {
    // 1. 保持事务简短
    err := db.Transaction(func(tx *gorm.DB) error {
        // 只在事务中执行必要的操作
        var user User
        if err := tx.First(&user, 1).Error; err != nil {
            return err
        }
        
        // 避免在事务中执行耗时操作
        // time.Sleep(10 * time.Second)  // 不要这样做
        
        return tx.Model(&user).Update("login_count", user.LoginCount+1).Error
    })
    
    // 2. 统一锁的获取顺序
    err = db.Transaction(func(tx *gorm.DB) error {
        // 总是按照相同的顺序获取锁（例如按ID排序）
        var users []User
        if err := tx.Where("id IN ?", []uint{1, 2, 3}).
            Order("id").  // 按ID排序获取锁
            Find(&users).Error; err != nil {
            return err
        }
        
        // 批量更新
        for _, user := range users {
            if err := tx.Model(&user).Update("status", "updated").Error; err != nil {
                return err
            }
        }
        
        return nil
    })
    
    // 3. 使用重试机制
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err = db.Transaction(func(tx *gorm.DB) error {
            // 事务逻辑
            return nil
        })
        
        if err == nil {
            break
        }
        
        // 检查是否是死锁错误
        if strings.Contains(err.Error(), "Deadlock") {
            time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
            continue
        }
        
        break
    }
}
```

### 问题2：事务回滚失败

**错误信息**：
```
transaction has already been committed or rolled back
```

**原因分析**：
- 重复调用Commit或Rollback
- 事务已经自动提交
- 连接已断开

**解决方案**：
```go
// 正确的事务处理模式
func correctTransactionHandling(db *gorm.DB) error {
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)  // 重新抛出panic
        }
    }()
    
    // 检查事务是否成功开始
    if tx.Error != nil {
        return tx.Error
    }
    
    // 执行业务逻辑
    if err := tx.Create(&User{Username: "test"}).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    if err := tx.Create(&Post{Title: "test"}).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return err
    }
    
    return nil
}

// 使用GORM的Transaction方法（推荐）
func recommendedTransactionHandling(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 所有操作都在这个函数中
        if err := tx.Create(&User{Username: "test"}).Error; err != nil {
            return err  // 自动回滚
        }
        
        if err := tx.Create(&Post{Title: "test"}).Error; err != nil {
            return err  // 自动回滚
        }
        
        return nil  // 自动提交
    })
}
```

---

## ⚡ 性能问题

### 问题1：查询速度慢

**症状**：
- 查询响应时间超过1秒
- 数据库CPU使用率高
- 大量慢查询日志

**诊断方法**：
```go
// 性能分析工具
func analyzePerformance(db *gorm.DB) {
    // 1. 启用查询日志
    db = db.Debug()
    
    // 2. 测量查询时间
    start := time.Now()
    var users []User
    db.Where("status = ?", "active").Find(&users)
    duration := time.Since(start)
    fmt.Printf("查询耗时: %v\n", duration)
    
    // 3. 分析执行计划
    var result []map[string]interface{}
    db.Raw("EXPLAIN SELECT * FROM user WHERE status = 'active'").Scan(&result)
    fmt.Printf("执行计划: %+v\n", result)
    
    // 4. 检查索引使用情况
    db.Raw("SHOW INDEX FROM user").Scan(&result)
    fmt.Printf("索引信息: %+v\n", result)
}
```

**解决方案**：
```go
// 性能优化策略
func optimizePerformance(db *gorm.DB) {
    // 1. 添加索引
    db.Exec("CREATE INDEX idx_user_status ON user(status)")
    db.Exec("CREATE INDEX idx_post_user_status ON post(user_id, status)")
    
    // 2. 使用复合索引
    db.Exec("CREATE INDEX idx_post_created_status ON post(created_at, status)")
    
    // 3. 优化查询
    // 避免SELECT *
    var users []User
    db.Select("id, username, email").Where("status = ?", "active").Find(&users)
    
    // 使用LIMIT
    db.Where("status = ?", "active").Limit(100).Find(&users)
    
    // 4. 使用批量操作
    // 避免N+1查询
    db.Preload("Posts").Find(&users)
    
    // 批量插入
    users = make([]User, 1000)
    db.CreateInBatches(users, 100)
    
    // 5. 使用原生SQL优化复杂查询
    var results []struct {
        UserID    uint
        PostCount int
    }
    
    db.Raw(`
        SELECT user_id, COUNT(*) as post_count 
        FROM post 
        WHERE status = 'published' 
        GROUP BY user_id 
        HAVING COUNT(*) > 10
    `).Scan(&results)
}
```

### 问题2：内存使用过高

**症状**：
- 应用内存持续增长
- 出现OOM错误
- 垃圾回收频繁

**解决方案**：
```go
// 内存优化策略
func optimizeMemory(db *gorm.DB) {
    // 1. 使用分页查询大量数据
    pageSize := 1000
    offset := 0
    
    for {
        var users []User
        result := db.Offset(offset).Limit(pageSize).Find(&users)
        
        if result.Error != nil {
            break
        }
        
        if len(users) == 0 {
            break
        }
        
        // 处理数据
        processUsers(users)
        
        // 清理内存
        users = nil
        runtime.GC()
        
        offset += pageSize
    }
    
    // 2. 使用游标查询
    var lastID uint = 0
    
    for {
        var users []User
        result := db.Where("id > ?", lastID).Order("id").Limit(pageSize).Find(&users)
        
        if result.Error != nil || len(users) == 0 {
            break
        }
        
        // 处理数据
        processUsers(users)
        
        lastID = users[len(users)-1].ID
        users = nil
    }
    
    // 3. 使用Rows进行流式处理
    rows, err := db.Model(&User{}).Where("status = ?", "active").Rows()
    if err != nil {
        return
    }
    defer rows.Close()
    
    for rows.Next() {
        var user User
        db.ScanRows(rows, &user)
        
        // 处理单个用户
        processUser(user)
    }
}

func processUsers(users []User) {
    // 处理用户数据
}

func processUser(user User) {
    // 处理单个用户
}
```

---

## 🔄 迁移问题

### 问题1：自动迁移失败

**错误信息**：
```
Error 1071: Specified key was too long
Error 1005: Can't create table
```

**原因分析**：
- 索引键长度超过限制
- 表名或字段名冲突
- 数据类型不兼容

**解决方案**：
```go
// 安全的迁移策略
func safeMigration(db *gorm.DB) error {
    // 1. 检查迁移前的状态
    if !db.Migrator().HasTable(&User{}) {
        fmt.Println("User表不存在，将创建")
    }
    
    // 2. 分步迁移
    models := []interface{}{
        &User{},
        &UserProfile{},
        &Category{},
        &Tag{},
        &Post{},
        &PostMeta{},
        &Comment{},
        &Like{},
        &Follow{},
        &Notification{},
        &Setting{},
    }
    
    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            fmt.Printf("迁移 %T 失败: %v\n", model, err)
            return err
        }
        fmt.Printf("迁移 %T 成功\n", model)
    }
    
    // 3. 手动创建索引（如果自动迁移失败）
    if err := createIndexesManually(db); err != nil {
        return err
    }
    
    return nil
}

func createIndexesManually(db *gorm.DB) error {
    indexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_user_username ON user(username)",
        "CREATE INDEX IF NOT EXISTS idx_user_email ON user(email)",
        "CREATE INDEX IF NOT EXISTS idx_post_user_id ON post(user_id)",
        "CREATE INDEX IF NOT EXISTS idx_post_status ON post(status)",
        "CREATE INDEX IF NOT EXISTS idx_comment_post_id ON comment(post_id)",
    }
    
    for _, indexSQL := range indexes {
        if err := db.Exec(indexSQL).Error; err != nil {
            fmt.Printf("创建索引失败: %s, 错误: %v\n", indexSQL, err)
            // 继续执行其他索引，不要因为一个失败就停止
        }
    }
    
    return nil
}
```

### 问题2：数据迁移冲突

**错误信息**：
```
Duplicate entry 'xxx' for key 'unique_index'
Data truncated for column 'xxx'
```

**解决方案**：
```go
// 数据清理和迁移
func cleanAndMigrate(db *gorm.DB) error {
    // 1. 备份现有数据
    if err := backupData(db); err != nil {
        return err
    }
    
    // 2. 清理重复数据
    if err := cleanDuplicateData(db); err != nil {
        return err
    }
    
    // 3. 执行迁移
    if err := db.AutoMigrate(&User{}, &Post{}).Error; err != nil {
        return err
    }
    
    return nil
}

func backupData(db *gorm.DB) error {
    // 导出数据到文件
    timestamp := time.Now().Format("20060102_150405")
    backupFile := fmt.Sprintf("backup_%s.sql", timestamp)
    
    // 这里可以使用mysqldump或其他工具
    fmt.Printf("数据已备份到: %s\n", backupFile)
    return nil
}

func cleanDuplicateData(db *gorm.DB) error {
    // 清理重复的用户名
    db.Exec(`
        DELETE u1 FROM user u1
        INNER JOIN user u2 
        WHERE u1.id > u2.id 
        AND u1.username = u2.username
    `)
    
    // 清理重复的邮箱
    db.Exec(`
        DELETE u1 FROM user u1
        INNER JOIN user u2 
        WHERE u1.id > u2.id 
        AND u1.email = u2.email
    `)
    
    return nil
}
```

---

## 🛠️ 调试工具和技巧

### 1. 启用详细日志

```go
import (
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

func setupLogger() logger.Interface {
    return logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,   // 慢查询阈值
            LogLevel:                  logger.Info,   // 日志级别
            IgnoreRecordNotFoundError: true,          // 忽略记录未找到错误
            Colorful:                  true,          // 彩色输出
        },
    )
}

func initDBWithLogger() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
        Logger: setupLogger(),
    })
    
    if err != nil {
        panic("failed to connect database")
    }
    
    return db
}
```

### 2. 性能监控

```go
// 查询性能监控中间件
func QueryMonitor() func(*gorm.DB) {
    return func(db *gorm.DB) {
        start := time.Now()
        
        db.Callback().Query().Before("gorm:query").Register("monitor:before", func(db *gorm.DB) {
            start = time.Now()
        })
        
        db.Callback().Query().After("gorm:query").Register("monitor:after", func(db *gorm.DB) {
            duration := time.Since(start)
            if duration > 100*time.Millisecond {
                fmt.Printf("慢查询警告: %v, SQL: %s\n", duration, db.Statement.SQL.String())
            }
        })
    }
}
```

### 3. 错误处理最佳实践

```go
func handleGORMError(err error) {
    if err == nil {
        return
    }
    
    switch {
    case errors.Is(err, gorm.ErrRecordNotFound):
        fmt.Println("记录未找到")
    case errors.Is(err, gorm.ErrInvalidTransaction):
        fmt.Println("无效事务")
    case errors.Is(err, gorm.ErrNotImplemented):
        fmt.Println("功能未实现")
    case errors.Is(err, gorm.ErrMissingWhereClause):
        fmt.Println("缺少WHERE条件")
    case errors.Is(err, gorm.ErrUnsupportedRelation):
        fmt.Println("不支持的关联关系")
    case errors.Is(err, gorm.ErrPrimaryKeyRequired):
        fmt.Println("需要主键")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

---

## 📞 获取帮助

### 社区资源

- **GORM官方文档**: https://gorm.io/docs/
- **GitHub Issues**: https://github.com/go-gorm/gorm/issues
- **Stack Overflow**: 搜索 `[go] [gorm]` 标签
- **Go语言中文网**: https://studygolang.com/

### 问题报告模板

当遇到无法解决的问题时，请按以下格式提供信息：

```
**环境信息**:
- Go版本: 
- GORM版本: 
- 数据库类型和版本: 
- 操作系统: 

**问题描述**:
简要描述遇到的问题

**重现步骤**:
1. 
2. 
3. 

**期望结果**:
描述期望的行为

**实际结果**:
描述实际发生的情况

**错误信息**:
```
完整的错误堆栈信息
```

**相关代码**:
```go
// 最小可重现的代码示例
```
```

---

## 🎯 预防措施

### 开发阶段

1. **代码审查清单**:
   - [ ] 所有查询都有适当的索引
   - [ ] 事务范围最小化
   - [ ] 错误处理完整
   - [ ] 没有N+1查询问题
   - [ ] 使用了连接池配置

2. **测试策略**:
   - 单元测试覆盖所有数据库操作
   - 集成测试验证关联关系
   - 性能测试确保查询效率
   - 压力测试验证并发安全

### 生产环境

1. **监控指标**:
   - 查询响应时间
   - 数据库连接数
   - 慢查询日志
   - 错误率统计

2. **告警设置**:
   - 慢查询超过阈值
   - 连接池耗尽
   - 错误率异常
   - 磁盘空间不足

---

**记住：遇到问题时，先查看错误信息，然后查阅文档，最后寻求社区帮助。大多数问题都有标准的解决方案！**
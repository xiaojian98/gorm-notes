# GORM Level 6 练习指导

## 📚 练习概述

本练习指导将帮助您逐步掌握 GORM Level 6 项目中的各个知识点，从基础的数据模型理解到高级的性能优化技巧。每个练习都包含详细的说明、示例代码和验证方法。

## 🎯 学习路径

```
基础练习 → 进阶练习 → 高级练习 → 实战项目
    ↓           ↓           ↓           ↓
数据模型     关联查询     性能优化     完整功能
基础操作     事务处理     索引设计     系统集成
```

## 📖 练习分级

- 🟢 **基础级**：适合GORM初学者
- 🟡 **进阶级**：需要一定GORM基础
- 🔴 **高级级**：需要深入理解数据库原理
- 🟣 **专家级**：需要丰富的实战经验

---

## 🟢 基础练习

### 练习 1：理解数据模型结构

**目标**：掌握项目中各个数据模型的设计思路和关联关系

**任务**：
1. 阅读 `User`、`Post`、`Comment` 等核心模型的定义
2. 绘制实体关系图（ERD）
3. 理解各种关联关系的实现方式

**练习代码**：
```go
// 1. 分析User模型的字段和约束
type User struct {
    BaseModel
    Username    string      `gorm:"uniqueIndex;size:50;not null;comment:用户名"`
    Email       string      `gorm:"uniqueIndex;size:100;not null;comment:邮箱"`
    // ... 其他字段
}

// 思考题：
// - 为什么Username和Email都设置了uniqueIndex？
// - BaseModel包含哪些字段？它们的作用是什么？
// - 软删除是如何实现的？
```

**验证方法**：
1. 运行程序，观察数据库表结构
2. 尝试插入重复的用户名或邮箱，观察错误信息
3. 删除一条记录，检查是否为软删除

**扩展练习**：
- 设计一个新的模型（如 `Article` 或 `Product`）
- 为新模型添加适当的字段和约束
- 考虑与现有模型的关联关系

---

### 练习 2：基础CRUD操作

**目标**：掌握GORM的基本增删改查操作

**任务**：
1. 创建用户记录
2. 查询用户信息
3. 更新用户数据
4. 删除用户记录

**练习代码**：
```go
func practiceBasicCRUD(db *gorm.DB) {
    // 1. 创建用户
    user := User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashed_password",
        Status:   "active",
    }
    
    // TODO: 使用GORM创建用户记录
    // 提示：使用 db.Create() 方法
    
    // 2. 查询用户
    var foundUser User
    // TODO: 根据用户名查询用户
    // 提示：使用 db.Where().First() 方法
    
    // 3. 更新用户
    // TODO: 更新用户的邮箱地址
    // 提示：使用 db.Model().Where().Update() 方法
    
    // 4. 删除用户
    // TODO: 软删除用户记录
    // 提示：使用 db.Delete() 方法
    
    // 5. 验证软删除
    var deletedUser User
    // TODO: 尝试查询已删除的用户
    // 提示：使用 db.Unscoped().Where().First() 方法
}
```

**参考答案**：
```go
func practiceBasicCRUD(db *gorm.DB) {
    // 1. 创建用户
    user := User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashed_password",
        Status:   "active",
    }
    
    result := db.Create(&user)
    if result.Error != nil {
        fmt.Printf("创建用户失败: %v\n", result.Error)
        return
    }
    fmt.Printf("创建用户成功，ID: %d\n", user.ID)
    
    // 2. 查询用户
    var foundUser User
    err := db.Where("username = ?", "testuser").First(&foundUser).Error
    if err != nil {
        fmt.Printf("查询用户失败: %v\n", err)
        return
    }
    fmt.Printf("查询到用户: %+v\n", foundUser)
    
    // 3. 更新用户
    err = db.Model(&foundUser).Where("id = ?", foundUser.ID).Update("email", "newemail@example.com").Error
    if err != nil {
        fmt.Printf("更新用户失败: %v\n", err)
        return
    }
    fmt.Println("更新用户成功")
    
    // 4. 删除用户
    err = db.Delete(&foundUser).Error
    if err != nil {
        fmt.Printf("删除用户失败: %v\n", err)
        return
    }
    fmt.Println("删除用户成功")
    
    // 5. 验证软删除
    var deletedUser User
    err = db.Unscoped().Where("id = ?", foundUser.ID).First(&deletedUser).Error
    if err != nil {
        fmt.Printf("查询已删除用户失败: %v\n", err)
        return
    }
    fmt.Printf("已删除用户仍存在: DeletedAt = %v\n", deletedUser.DeletedAt)
}
```

---

### 练习 3：关联关系操作

**目标**：理解和操作一对一、一对多、多对多关系

**任务**：
1. 创建用户及其资料（一对一）
2. 为用户创建多篇文章（一对多）
3. 为文章添加标签（多对多）

**练习代码**：
```go
func practiceAssociations(db *gorm.DB) {
    // 1. 创建用户和用户资料（一对一关系）
    user := User{
        Username: "blogger",
        Email:    "blogger@example.com",
        Password: "hashed_password",
        Status:   "active",
    }
    
    // TODO: 创建用户
    
    profile := UserProfile{
        UserID:  user.ID,
        Bio:     "我是一个博客作者",
        Avatar:  "avatar.jpg",
        Website: "https://myblog.com",
    }
    
    // TODO: 创建用户资料
    
    // 2. 创建文章（一对多关系）
    posts := []Post{
        {
            Title:      "我的第一篇文章",
            Content:    "这是我的第一篇文章内容...",
            UserID:     user.ID,
            CategoryID: 1, // 假设分类ID为1
            Status:     "published",
        },
        {
            Title:      "GORM学习笔记",
            Content:    "今天学习了GORM的基础用法...",
            UserID:     user.ID,
            CategoryID: 1,
            Status:     "published",
        },
    }
    
    // TODO: 批量创建文章
    
    // 3. 创建标签并关联到文章（多对多关系）
    tags := []Tag{
        {Name: "Go语言", Color: "#00ADD8"},
        {Name: "数据库", Color: "#336791"},
        {Name: "GORM", Color: "#FF6B6B"},
    }
    
    // TODO: 创建标签
    
    // TODO: 将标签关联到第二篇文章
    
    // 4. 查询关联数据
    // TODO: 查询用户及其资料
    
    // TODO: 查询用户及其所有文章
    
    // TODO: 查询文章及其标签
}
```

**参考答案**：
```go
func practiceAssociations(db *gorm.DB) {
    // 1. 创建用户和用户资料
    user := User{
        Username: "blogger",
        Email:    "blogger@example.com",
        Password: "hashed_password",
        Status:   "active",
    }
    
    db.Create(&user)
    
    profile := UserProfile{
        UserID:  user.ID,
        Bio:     "我是一个博客作者",
        Avatar:  "avatar.jpg",
        Website: "https://myblog.com",
    }
    
    db.Create(&profile)
    
    // 2. 创建文章
    posts := []Post{
        {
            Title:      "我的第一篇文章",
            Content:    "这是我的第一篇文章内容...",
            UserID:     user.ID,
            CategoryID: 1,
            Status:     "published",
        },
        {
            Title:      "GORM学习笔记",
            Content:    "今天学习了GORM的基础用法...",
            UserID:     user.ID,
            CategoryID: 1,
            Status:     "published",
        },
    }
    
    db.Create(&posts)
    
    // 3. 创建标签并关联
    tags := []Tag{
        {Name: "Go语言", Color: "#00ADD8"},
        {Name: "数据库", Color: "#336791"},
        {Name: "GORM", Color: "#FF6B6B"},
    }
    
    db.Create(&tags)
    
    // 关联标签到第二篇文章
    db.Model(&posts[1]).Association("Tags").Append(tags)
    
    // 4. 查询关联数据
    var userWithProfile User
    db.Preload("Profile").First(&userWithProfile, user.ID)
    fmt.Printf("用户及资料: %+v\n", userWithProfile)
    
    var userWithPosts User
    db.Preload("Posts").First(&userWithPosts, user.ID)
    fmt.Printf("用户文章数: %d\n", len(userWithPosts.Posts))
    
    var postWithTags Post
    db.Preload("Tags").First(&postWithTags, posts[1].ID)
    fmt.Printf("文章标签数: %d\n", len(postWithTags.Tags))
}
```

---

## 🟡 进阶练习

### 练习 4：复杂查询操作

**目标**：掌握GORM的高级查询技巧

**任务**：
1. 条件查询和排序
2. 聚合函数使用
3. 子查询操作
4. 原生SQL查询

**练习代码**：
```go
func practiceAdvancedQueries(db *gorm.DB) {
    // 1. 复杂条件查询
    // TODO: 查询最近30天内发布的文章，按点赞数降序排列
    
    // 2. 聚合查询
    // TODO: 统计每个用户的文章数量
    
    // 3. 子查询
    // TODO: 查询点赞数超过平均值的文章
    
    // 4. 分组查询
    // TODO: 按分类统计文章数量和平均点赞数
    
    // 5. 连接查询
    // TODO: 查询用户及其文章的总点赞数
}
```

**参考答案**：
```go
func practiceAdvancedQueries(db *gorm.DB) {
    // 1. 复杂条件查询
    var recentPosts []Post
    thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
    
    db.Preload("User").Preload("Meta").
        Where("status = ? AND created_at > ?", "published", thirtyDaysAgo).
        Joins("LEFT JOIN post_meta ON post.id = post_meta.post_id").
        Order("post_meta.like_count DESC").
        Find(&recentPosts)
    
    fmt.Printf("最近30天文章数: %d\n", len(recentPosts))
    
    // 2. 聚合查询
    type UserPostCount struct {
        UserID    uint
        Username  string
        PostCount int64
    }
    
    var userStats []UserPostCount
    db.Table("user").
        Select("user.id as user_id, user.username, COUNT(post.id) as post_count").
        Joins("LEFT JOIN post ON user.id = post.user_id AND post.status = 'published'").
        Group("user.id, user.username").
        Having("COUNT(post.id) > 0").
        Find(&userStats)
    
    fmt.Printf("用户文章统计: %+v\n", userStats)
    
    // 3. 子查询
    var avgLikes float64
    db.Table("post_meta").Select("AVG(like_count)").Scan(&avgLikes)
    
    var popularPosts []Post
    db.Preload("User").Preload("Meta").
        Joins("JOIN post_meta ON post.id = post_meta.post_id").
        Where("post_meta.like_count > ?", avgLikes).
        Find(&popularPosts)
    
    fmt.Printf("热门文章数: %d (平均点赞数: %.2f)\n", len(popularPosts), avgLikes)
    
    // 4. 分组查询
    type CategoryStats struct {
        CategoryID   uint
        CategoryName string
        PostCount    int64
        AvgLikes     float64
    }
    
    var categoryStats []CategoryStats
    db.Table("category").
        Select("category.id as category_id, category.name as category_name, COUNT(post.id) as post_count, AVG(post_meta.like_count) as avg_likes").
        Joins("LEFT JOIN post ON category.id = post.category_id AND post.status = 'published'").
        Joins("LEFT JOIN post_meta ON post.id = post_meta.post_id").
        Group("category.id, category.name").
        Find(&categoryStats)
    
    fmt.Printf("分类统计: %+v\n", categoryStats)
    
    // 5. 连接查询
    type UserLikeStats struct {
        UserID     uint
        Username   string
        TotalLikes int64
    }
    
    var userLikeStats []UserLikeStats
    db.Table("user").
        Select("user.id as user_id, user.username, COALESCE(SUM(post_meta.like_count), 0) as total_likes").
        Joins("LEFT JOIN post ON user.id = post.user_id AND post.status = 'published'").
        Joins("LEFT JOIN post_meta ON post.id = post_meta.post_id").
        Group("user.id, user.username").
        Order("total_likes DESC").
        Find(&userLikeStats)
    
    fmt.Printf("用户点赞统计: %+v\n", userLikeStats)
}
```

---

### 练习 5：事务处理

**目标**：掌握GORM的事务处理机制

**任务**：
1. 基础事务操作
2. 嵌套事务处理
3. 事务回滚机制
4. 并发事务控制

**练习代码**：
```go
func practiceTransactions(db *gorm.DB) {
    // 1. 基础事务 - 创建文章和更新统计
    // TODO: 在事务中创建文章、文章元数据，并更新用户文章数
    
    // 2. 事务回滚 - 模拟失败场景
    // TODO: 创建一个会失败的事务，观察回滚效果
    
    // 3. 手动事务控制
    // TODO: 使用Begin、Commit、Rollback手动控制事务
    
    // 4. 并发事务测试
    // TODO: 模拟并发更新同一条记录的场景
}
```

**参考答案**：
```go
func practiceTransactions(db *gorm.DB) {
    // 1. 基础事务
    err := db.Transaction(func(tx *gorm.DB) error {
        // 创建文章
        post := Post{
            Title:      "事务测试文章",
            Content:    "这是一篇用于测试事务的文章",
            UserID:     1,
            CategoryID: 1,
            Status:     "published",
        }
        
        if err := tx.Create(&post).Error; err != nil {
            return err
        }
        
        // 创建文章元数据
        meta := PostMeta{
            PostID:    post.ID,
            ViewCount: 0,
            LikeCount: 0,
        }
        
        if err := tx.Create(&meta).Error; err != nil {
            return err
        }
        
        // 更新用户文章数
        if err := tx.Model(&User{}).Where("id = ?", post.UserID).
            Update("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        fmt.Printf("事务执行失败: %v\n", err)
    } else {
        fmt.Println("事务执行成功")
    }
    
    // 2. 事务回滚测试
    err = db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        user := User{
            Username: "transaction_test",
            Email:    "transaction@test.com",
            Password: "password",
        }
        
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        
        fmt.Printf("用户创建成功，ID: %d\n", user.ID)
        
        // 故意制造错误（违反唯一约束）
        duplicateUser := User{
            Username: "transaction_test", // 重复用户名
            Email:    "another@test.com",
            Password: "password",
        }
        
        if err := tx.Create(&duplicateUser).Error; err != nil {
            fmt.Printf("预期的错误发生: %v\n", err)
            return err // 触发回滚
        }
        
        return nil
    })
    
    if err != nil {
        fmt.Println("事务已回滚")
        
        // 验证回滚效果
        var count int64
        db.Model(&User{}).Where("username = ?", "transaction_test").Count(&count)
        fmt.Printf("回滚后用户数量: %d\n", count)
    }
    
    // 3. 手动事务控制
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            fmt.Println("事务因panic回滚")
        }
    }()
    
    user := User{
        Username: "manual_tx_user",
        Email:    "manual@tx.com",
        Password: "password",
    }
    
    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        fmt.Printf("手动事务回滚: %v\n", err)
        return
    }
    
    // 模拟一些业务逻辑
    time.Sleep(100 * time.Millisecond)
    
    if err := tx.Commit().Error; err != nil {
        fmt.Printf("事务提交失败: %v\n", err)
        return
    }
    
    fmt.Println("手动事务提交成功")
    
    // 4. 并发事务测试
    var wg sync.WaitGroup
    userID := uint(1)
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(goroutineID int) {
            defer wg.Done()
            
            err := db.Transaction(func(tx *gorm.DB) error {
                var user User
                if err := tx.First(&user, userID).Error; err != nil {
                    return err
                }
                
                // 模拟一些处理时间
                time.Sleep(10 * time.Millisecond)
                
                // 更新用户登录次数
                return tx.Model(&user).Update("login_count", gorm.Expr("login_count + ?", 1)).Error
            })
            
            if err != nil {
                fmt.Printf("Goroutine %d 事务失败: %v\n", goroutineID, err)
            } else {
                fmt.Printf("Goroutine %d 事务成功\n", goroutineID)
            }
        }(i)
    }
    
    wg.Wait()
    fmt.Println("并发事务测试完成")
}
```

---

## 🔴 高级练习

### 练习 6：性能优化

**目标**：掌握数据库性能优化技巧

**任务**：
1. 索引优化分析
2. 查询性能测试
3. 连接池调优
4. 内存使用优化

**练习代码**：
```go
func practicePerformanceOptimization(db *gorm.DB) {
    // 1. 索引效果对比测试
    // TODO: 对比有索引和无索引的查询性能
    
    // 2. 预加载vs N+1查询对比
    // TODO: 对比使用Preload和不使用的性能差异
    
    // 3. 批量操作优化
    // TODO: 对比单条插入和批量插入的性能
    
    // 4. 分页查询优化
    // TODO: 对比OFFSET分页和游标分页的性能
}
```

**参考答案**：
```go
func practicePerformanceOptimization(db *gorm.DB) {
    // 1. 索引效果对比
    fmt.Println("=== 索引效果对比测试 ===")
    
    // 创建测试数据
    users := make([]User, 1000)
    for i := 0; i < 1000; i++ {
        users[i] = User{
            Username: fmt.Sprintf("user_%d", i),
            Email:    fmt.Sprintf("user_%d@test.com", i),
            Password: "password",
            Status:   []string{"active", "inactive", "banned"}[i%3],
        }
    }
    db.CreateInBatches(users, 100)
    
    // 测试有索引的查询（status字段有索引）
    start := time.Now()
    var activeUsers []User
    db.Where("status = ?", "active").Find(&activeUsers)
    indexedQueryTime := time.Since(start)
    
    // 测试无索引的查询（假设email字段无索引）
    start = time.Now()
    var emailUsers []User
    db.Where("email LIKE ?", "%user_1%").Find(&emailUsers)
    nonIndexedQueryTime := time.Since(start)
    
    fmt.Printf("有索引查询时间: %v\n", indexedQueryTime)
    fmt.Printf("无索引查询时间: %v\n", nonIndexedQueryTime)
    fmt.Printf("性能提升: %.2fx\n", float64(nonIndexedQueryTime)/float64(indexedQueryTime))
    
    // 2. 预加载vs N+1查询对比
    fmt.Println("\n=== 预加载效果对比测试 ===")
    
    // N+1查询（不使用预加载）
    start = time.Now()
    var posts []Post
    db.Limit(10).Find(&posts)
    for i := range posts {
        db.First(&posts[i].User, posts[i].UserID) // N+1查询
    }
    n1QueryTime := time.Since(start)
    
    // 使用预加载
    start = time.Now()
    var postsWithPreload []Post
    db.Preload("User").Limit(10).Find(&postsWithPreload)
    preloadQueryTime := time.Since(start)
    
    fmt.Printf("N+1查询时间: %v\n", n1QueryTime)
    fmt.Printf("预加载查询时间: %v\n", preloadQueryTime)
    fmt.Printf("性能提升: %.2fx\n", float64(n1QueryTime)/float64(preloadQueryTime))
    
    // 3. 批量操作优化
    fmt.Println("\n=== 批量操作对比测试 ===")
    
    // 单条插入
    testUsers := make([]User, 100)
    for i := 0; i < 100; i++ {
        testUsers[i] = User{
            Username: fmt.Sprintf("batch_test_%d", i),
            Email:    fmt.Sprintf("batch_%d@test.com", i),
            Password: "password",
        }
    }
    
    start = time.Now()
    for _, user := range testUsers {
        db.Create(&user)
    }
    singleInsertTime := time.Since(start)
    
    // 批量插入
    start = time.Now()
    db.CreateInBatches(testUsers, 20)
    batchInsertTime := time.Since(start)
    
    fmt.Printf("单条插入时间: %v\n", singleInsertTime)
    fmt.Printf("批量插入时间: %v\n", batchInsertTime)
    fmt.Printf("性能提升: %.2fx\n", float64(singleInsertTime)/float64(batchInsertTime))
    
    // 4. 分页查询优化
    fmt.Println("\n=== 分页查询对比测试 ===")
    
    // OFFSET分页（传统分页）
    page := 50
    pageSize := 20
    offset := (page - 1) * pageSize
    
    start = time.Now()
    var offsetPosts []Post
    db.Offset(offset).Limit(pageSize).Find(&offsetPosts)
    offsetPagingTime := time.Since(start)
    
    // 游标分页
    var lastID uint = 1000 // 假设从ID 1000开始
    start = time.Now()
    var cursorPosts []Post
    db.Where("id < ?", lastID).Order("id DESC").Limit(pageSize).Find(&cursorPosts)
    cursorPagingTime := time.Since(start)
    
    fmt.Printf("OFFSET分页时间: %v\n", offsetPagingTime)
    fmt.Printf("游标分页时间: %v\n", cursorPagingTime)
    fmt.Printf("性能提升: %.2fx\n", float64(offsetPagingTime)/float64(cursorPagingTime))
    
    // 5. 连接池监控
    fmt.Println("\n=== 连接池状态监控 ===")
    sqlDB, _ := db.DB()
    stats := sqlDB.Stats()
    
    fmt.Printf("最大打开连接数: %d\n", stats.MaxOpenConnections)
    fmt.Printf("当前打开连接数: %d\n", stats.OpenConnections)
    fmt.Printf("使用中连接数: %d\n", stats.InUse)
    fmt.Printf("空闲连接数: %d\n", stats.Idle)
    fmt.Printf("等待连接数: %d\n", stats.WaitCount)
    fmt.Printf("等待总时长: %v\n", stats.WaitDuration)
}
```

---

### 练习 7：数据库设计模式

**目标**：学习常见的数据库设计模式和最佳实践

**任务**：
1. 实现审计日志模式
2. 实现软删除模式
3. 实现版本控制模式
4. 实现读写分离模式

**练习代码**：
```go
// 1. 审计日志模式
type AuditLog struct {
    ID        uint      `gorm:"primaryKey"`
    TableName string    `gorm:"size:50;not null"`
    RecordID  uint      `gorm:"not null"`
    Action    string    `gorm:"size:20;not null"` // CREATE, UPDATE, DELETE
    OldData   string    `gorm:"type:json"`
    NewData   string    `gorm:"type:json"`
    UserID    uint      `gorm:"not null"`
    CreatedAt time.Time
}

// TODO: 实现审计日志的钩子函数
func (u *User) AfterCreate(tx *gorm.DB) error {
    // 记录创建日志
    return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
    // 记录更新日志
    return nil
}

// 2. 版本控制模式
type VersionedPost struct {
    BaseModel
    PostID    uint   `gorm:"not null;index"`
    Version   int    `gorm:"not null;index"`
    Title     string `gorm:"size:200;not null"`
    Content   string `gorm:"type:longtext;not null"`
    CreatedBy uint   `gorm:"not null"`
    IsCurrent bool   `gorm:"not null;default:false;index"`
}

// TODO: 实现版本控制逻辑
func createNewVersion(db *gorm.DB, postID uint, title, content string, userID uint) error {
    // 实现版本创建逻辑
    return nil
}
```

---

## 🟣 专家练习

### 练习 8：分布式数据库设计

**目标**：设计支持分布式部署的数据库架构

**任务**：
1. 实现数据分片策略
2. 设计跨库事务处理
3. 实现数据同步机制
4. 设计故障恢复方案

### 练习 9：性能监控系统

**目标**：构建完整的数据库性能监控系统

**任务**：
1. 实现慢查询监控
2. 设计性能指标收集
3. 实现告警机制
4. 构建性能分析报告

### 练习 10：数据迁移工具

**目标**：开发数据库迁移和版本管理工具

**任务**：
1. 设计迁移脚本格式
2. 实现版本控制机制
3. 支持回滚操作
4. 实现数据验证

---

## 📝 练习评估

### 自我评估清单

**基础知识** (🟢)
- [ ] 理解GORM基本概念和用法
- [ ] 掌握数据模型设计原则
- [ ] 熟悉CRUD操作
- [ ] 理解关联关系

**进阶技能** (🟡)
- [ ] 掌握复杂查询技巧
- [ ] 理解事务处理机制
- [ ] 熟悉性能优化方法
- [ ] 掌握错误处理策略

**高级能力** (🔴)
- [ ] 设计高性能数据库架构
- [ ] 实现复杂业务逻辑
- [ ] 掌握分布式数据库概念
- [ ] 具备故障排查能力

**专家水平** (🟣)
- [ ] 设计企业级数据库方案
- [ ] 实现高可用架构
- [ ] 掌握性能调优技巧
- [ ] 具备架构设计能力

### 项目实战建议

1. **个人博客系统**：基于本项目扩展，添加更多功能
2. **电商平台**：设计商品、订单、支付等复杂业务模型
3. **社交网络**：实现用户关系、动态、消息等功能
4. **内容管理系统**：构建企业级CMS平台

### 学习资源推荐

- **官方文档**：[GORM官方文档](https://gorm.io/docs/)
- **进阶教程**：[Go数据库编程实战]()
- **性能优化**：[MySQL性能调优指南]()
- **架构设计**：[分布式数据库设计模式]()

---

**祝您学习顺利！记住：实践是最好的老师，多动手、多思考、多总结！**
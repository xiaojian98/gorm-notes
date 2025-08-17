# GORM Level 6 技术实现指南

## 📋 目录

- [项目结构](#项目结构)
- [数据模型设计](#数据模型设计)
- [数据库配置](#数据库配置)
- [服务层架构](#服务层架构)
- [查询优化](#查询优化)
- [性能监控](#性能监控)
- [最佳实践](#最佳实践)
- [故障排查](#故障排查)

## 🏗️ 项目结构

```
level6/
├── level6_comprehensive.go    # 主程序文件
├── README.md                  # 快速入门指南
├── TECHNICAL_GUIDE.md         # 技术实现指南
└── go.mod                     # Go模块依赖
```

### 代码组织结构

```go
// 1. 包声明和导入
package main
import (...)

// 2. 配置和枚举定义
type DatabaseType int
type DatabaseConfig struct {...}

// 3. 数据模型定义
type BaseModel struct {...}
type User struct {...}
// ... 其他模型

// 4. 数据库初始化
func initDB(config DatabaseConfig) *gorm.DB {...}
func createIndexes(db *gorm.DB) {...}

// 5. 服务层实现
type UserService struct {...}
type PostService struct {...}
// ... 其他服务

// 6. 业务逻辑演示
func demonstrateComprehensiveScenarios(db *gorm.DB) {...}
func demonstrateAdvancedQueries(db *gorm.DB) {...}

// 7. 性能测试
func performanceTest(db *gorm.DB) {...}

// 8. 主函数
func main() {...}
```

## 🗄️ 数据模型设计

### 基础模型 (BaseModel)

```go
type BaseModel struct {
    ID        uint           `gorm:"primaryKey;autoIncrement;comment:主键ID"`
    CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
    UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
    DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间"`
}
```

**设计要点**：
- **主键策略**：使用自增整型主键，性能优于UUID
- **时间戳**：自动管理创建和更新时间
- **软删除**：支持逻辑删除，保留数据历史
- **索引优化**：为DeletedAt字段建立索引，提升软删除查询性能

### 用户模型设计

```go
type User struct {
    BaseModel
    Username    string      `gorm:"uniqueIndex;size:50;not null;comment:用户名"`
    Email       string      `gorm:"uniqueIndex;size:100;not null;comment:邮箱"`
    Password    string      `gorm:"size:255;not null;comment:密码哈希"`
    Status      string      `gorm:"size:20;not null;default:'active';index;comment:用户状态"`
    LastLoginAt *time.Time  `gorm:"comment:最后登录时间"`
    
    // 关联关系
    Profile       UserProfile    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Posts         []Post         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Comments      []Comment      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Likes         []Like         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Followers     []Follow       `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE"`
    Following     []Follow       `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
    Notifications []Notification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
```

**关键特性**：
- **唯一约束**：用户名和邮箱的唯一性保证
- **状态管理**：支持用户状态控制（active/inactive/banned）
- **关联设计**：完整的用户关系网络
- **级联删除**：保证数据一致性

### 文章模型设计

```go
type Post struct {
    BaseModel
    Title      string    `gorm:"size:200;not null;index;comment:文章标题"`
    Slug       string    `gorm:"uniqueIndex;size:200;not null;comment:URL友好标识"`
    Content    string    `gorm:"type:longtext;not null;comment:文章内容"`
    Summary    string    `gorm:"size:500;comment:文章摘要"`
    UserID     uint      `gorm:"not null;index;comment:作者ID"`
    CategoryID uint      `gorm:"not null;index;comment:分类ID"`
    Status     string    `gorm:"size:20;not null;default:'draft';index;comment:文章状态"`
    
    // 关联关系
    User     User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Category Category   `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
    Meta     PostMeta   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
    Comments []Comment  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
    Tags     []Tag      `gorm:"many2many:post_tag;constraint:OnDelete:CASCADE"`
    Likes    []Like     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
```

**设计亮点**：
- **SEO优化**：Slug字段支持友好URL
- **内容分离**：长文本内容与元数据分离
- **状态控制**：支持草稿、发布、归档等状态
- **多对多关系**：灵活的标签系统

### 评论模型设计

```go
type Comment struct {
    BaseModel
    Content  string `gorm:"type:text;not null;comment:评论内容"`
    PostID   uint   `gorm:"not null;index;comment:文章ID"`
    UserID   uint   `gorm:"not null;index;comment:用户ID"`
    ParentID *uint  `gorm:"index;comment:父评论ID"`
    Status   string `gorm:"size:20;not null;default:'approved';index;comment:评论状态"`
    
    // 关联关系
    Post     Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
    User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Parent   *Comment  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
    Children []Comment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
    Likes    []Like    `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}
```

**特色功能**：
- **层级结构**：支持评论回复的树形结构
- **审核机制**：评论状态管理
- **自关联**：Parent-Children关系实现

## ⚙️ 数据库配置

### 配置结构设计

```go
type DatabaseType int

const (
    SQLite DatabaseType = iota
    MySQL
)

type DatabaseConfig struct {
    Type            DatabaseType
    DSN             string
    MaxIdleConns    int
    MaxOpenConns    int
    ConnMaxLifetime time.Duration
    
    // MySQL专用配置
    Host     string
    Port     int
    Username string
    Password string
    Database string
}
```

### 数据库初始化

```go
func initDB(config DatabaseConfig) *gorm.DB {
    var db *gorm.DB
    var err error
    
    // 根据数据库类型选择驱动
    switch config.Type {
    case MySQL:
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            config.Username, config.Password, config.Host, config.Port, config.Database)
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Info),
            NamingStrategy: schema.NamingStrategy{
                SingularTable: true, // 使用单数表名
            },
            PrepareStmt: true, // 预编译语句
        })
    case SQLite:
        db, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Info),
            NamingStrategy: schema.NamingStrategy{
                SingularTable: true,
            },
            PrepareStmt: true,
        })
    }
    
    if err != nil {
        log.Fatalf("数据库连接失败: %v", err)
    }
    
    // 配置连接池
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(config.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    // 自动迁移
    autoMigrate(db)
    
    // 创建索引
    createIndexes(db)
    
    return db
}
```

**配置要点**：
- **连接池优化**：合理设置连接数和生命周期
- **预编译语句**：提升查询性能
- **日志配置**：便于调试和监控
- **命名策略**：统一的表名规范

### 索引创建策略

```go
func createIndexes(db *gorm.DB) {
    // 复合索引 - 用户相关
    db.Exec("CREATE INDEX IF NOT EXISTS idx_user_status_login ON user(status, last_login_at DESC)")
    
    // 复合索引 - 文章相关
    db.Exec("CREATE INDEX IF NOT EXISTS idx_post_user_status ON post(user_id, status)")
    db.Exec("CREATE INDEX IF NOT EXISTS idx_post_category_created ON post(category_id, created_at DESC)")
    db.Exec("CREATE INDEX IF NOT EXISTS idx_post_status_created ON post(status, created_at DESC)")
    
    // 复合索引 - 评论相关
    db.Exec("CREATE INDEX IF NOT EXISTS idx_comment_post_created ON comment(post_id, created_at DESC)")
    db.Exec("CREATE INDEX IF NOT EXISTS idx_comment_user_created ON comment(user_id, created_at DESC)")
    
    // 唯一索引 - 关注关系
    db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_follow_unique ON follow(follower_id, followed_id)")
    
    // 唯一索引 - 点赞关系
    db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_like_unique ON like(user_id, post_id, comment_id)")
}
```

**索引设计原则**：
- **查询频率**：为高频查询字段建立索引
- **复合索引**：按查询条件的选择性排序
- **唯一约束**：防止重复数据
- **性能平衡**：避免过多索引影响写入性能

## 🏢 服务层架构

### 用户服务 (UserService)

```go
type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// CreateUser 创建用户（事务处理）
func (s *UserService) CreateUser(username, email, password string) (*User, error) {
    var user User
    
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // 检查用户名和邮箱唯一性
        var count int64
        if err := tx.Model(&User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error; err != nil {
            return err
        }
        if count > 0 {
            return fmt.Errorf("用户名或邮箱已存在")
        }
        
        // 创建用户
        user = User{
            Username: username,
            Email:    email,
            Password: hashPassword(password), // 密码哈希
            Status:   "active",
        }
        
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        
        // 创建用户资料
        profile := UserProfile{
            UserID: user.ID,
            Bio:    "这个人很懒，什么都没留下",
        }
        
        return tx.Create(&profile).Error
    })
    
    return &user, err
}

// GetUserWithProfile 获取用户及其资料
func (s *UserService) GetUserWithProfile(userID uint) (*User, error) {
    var user User
    err := s.db.Preload("Profile").First(&user, userID).Error
    return &user, err
}

// UpdateLastLogin 更新最后登录时间
func (s *UserService) UpdateLastLogin(userID uint) error {
    now := time.Now()
    return s.db.Model(&User{}).Where("id = ?", userID).Update("last_login_at", &now).Error
}
```

### 文章服务 (PostService)

```go
type PostService struct {
    db *gorm.DB
}

// CreatePost 创建文章（完整业务流程）
func (s *PostService) CreatePost(userID uint, title, content, summary string, categoryID uint, tagNames []string) (*Post, error) {
    var post Post
    
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // 生成Slug
        slug := generateSlug(title)
        
        // 创建文章
        post = Post{
            Title:      title,
            Slug:       slug,
            Content:    content,
            Summary:    summary,
            UserID:     userID,
            CategoryID: categoryID,
            Status:     "published",
        }
        
        if err := tx.Create(&post).Error; err != nil {
            return err
        }
        
        // 创建文章元数据
        meta := PostMeta{
            PostID:     post.ID,
            ViewCount:  0,
            LikeCount:  0,
            ShareCount: 0,
            ReadTime:   calculateReadTime(content),
        }
        
        if err := tx.Create(&meta).Error; err != nil {
            return err
        }
        
        // 处理标签
        if len(tagNames) > 0 {
            var tags []Tag
            for _, name := range tagNames {
                var tag Tag
                // 查找或创建标签
                if err := tx.Where("name = ?", name).FirstOrCreate(&tag, Tag{Name: name}).Error; err != nil {
                    return err
                }
                tags = append(tags, tag)
            }
            
            // 关联标签
            if err := tx.Model(&post).Association("Tags").Append(tags); err != nil {
                return err
            }
        }
        
        // 更新分类文章数
        return tx.Model(&Category{}).Where("id = ?", categoryID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error
    })
    
    return &post, err
}

// GetPostsWithPagination 分页获取文章列表
func (s *PostService) GetPostsWithPagination(page, pageSize int, categoryID *uint, status string) ([]Post, int64, error) {
    var posts []Post
    var total int64
    
    query := s.db.Model(&Post{}).Where("status = ?", status)
    
    if categoryID != nil {
        query = query.Where("category_id = ?", *categoryID)
    }
    
    // 获取总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    offset := (page - 1) * pageSize
    err := query.Preload("User").Preload("Category").Preload("Tags").Preload("Meta").
        Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
    
    return posts, total, err
}
```

### 评论服务 (CommentService)

```go
type CommentService struct {
    db *gorm.DB
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(userID, postID uint, content string, parentID *uint) (*Comment, error) {
    var comment Comment
    
    err := s.db.Transaction(func(tx *gorm.DB) error {
        // 内容过滤
        filteredContent := filterSensitiveWords(content)
        
        comment = Comment{
            Content:  filteredContent,
            PostID:   postID,
            UserID:   userID,
            ParentID: parentID,
            Status:   "approved",
        }
        
        if err := tx.Create(&comment).Error; err != nil {
            return err
        }
        
        // 更新文章评论数
        return tx.Model(&PostMeta{}).Where("post_id = ?", postID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
    })
    
    return &comment, err
}

// GetCommentTree 获取评论树
func (s *CommentService) GetCommentTree(postID uint) ([]Comment, error) {
    var comments []Comment
    
    // 获取所有评论
    err := s.db.Where("post_id = ? AND status = ?", postID, "approved").
        Preload("User").Order("created_at ASC").Find(&comments).Error
    
    if err != nil {
        return nil, err
    }
    
    // 构建树形结构
    return buildCommentTree(comments), nil
}
```

## 🚀 查询优化

### 预加载策略

```go
// 避免N+1查询问题
func GetPostsWithAllRelations() []Post {
    var posts []Post
    
    db.Preload("User").                    // 预加载作者
      Preload("Category").                // 预加载分类
      Preload("Tags").                    // 预加载标签
      Preload("Meta").                    // 预加载元数据
      Preload("Comments", func(db *gorm.DB) *gorm.DB {
          return db.Where("status = ?", "approved").Limit(5) // 只加载前5条评论
      }).Find(&posts)
    
    return posts
}
```

### 选择性字段查询

```go
// 只查询需要的字段
func GetPostSummaries() []Post {
    var posts []Post
    
    db.Select("id, title, summary, user_id, category_id, created_at").
       Where("status = ?", "published").
       Order("created_at DESC").
       Limit(20).Find(&posts)
    
    return posts
}
```

### 复杂查询示例

```go
// 获取热门文章（基于点赞数和评论数）
func GetHotPosts(limit int) []Post {
    var posts []Post
    
    db.Table("post p").
       Select("p.*, pm.like_count, pm.comment_count, (pm.like_count * 2 + pm.comment_count) as hot_score").
       Joins("LEFT JOIN post_meta pm ON p.id = pm.post_id").
       Where("p.status = ? AND p.created_at > ?", "published", time.Now().AddDate(0, 0, -30)).
       Order("hot_score DESC, p.created_at DESC").
       Limit(limit).Find(&posts)
    
    return posts
}

// 获取用户活跃度统计
func GetUserActivityStats(userID uint) map[string]interface{} {
    var result map[string]interface{}
    
    db.Raw(`
        SELECT 
            u.username,
            COUNT(DISTINCT p.id) as post_count,
            COUNT(DISTINCT c.id) as comment_count,
            COUNT(DISTINCT l.id) as like_count,
            COUNT(DISTINCT f.id) as follower_count
        FROM user u
        LEFT JOIN post p ON u.id = p.user_id AND p.status = 'published'
        LEFT JOIN comment c ON u.id = c.user_id AND c.status = 'approved'
        LEFT JOIN like l ON u.id = l.user_id
        LEFT JOIN follow f ON u.id = f.followed_id
        WHERE u.id = ?
        GROUP BY u.id
    `, userID).Scan(&result)
    
    return result
}
```

### 分页优化

```go
// 游标分页（适用于大数据量）
func GetPostsWithCursor(cursor uint, limit int) []Post {
    var posts []Post
    
    query := db.Where("status = ?", "published")
    
    if cursor > 0 {
        query = query.Where("id < ?", cursor)
    }
    
    query.Order("id DESC").Limit(limit).Find(&posts)
    
    return posts
}

// 传统分页（适用于小数据量）
func GetPostsWithOffset(page, pageSize int) ([]Post, int64) {
    var posts []Post
    var total int64
    
    db.Model(&Post{}).Where("status = ?", "published").Count(&total)
    
    offset := (page - 1) * pageSize
    db.Where("status = ?", "published").
       Order("created_at DESC").
       Offset(offset).Limit(pageSize).Find(&posts)
    
    return posts, total
}
```

## 📊 性能监控

### 查询性能分析

```go
// 启用SQL日志
func EnableSQLLogging(db *gorm.DB) {
    db.Logger = logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,   // 慢查询阈值
            LogLevel:                  logger.Info,   // 日志级别
            IgnoreRecordNotFoundError: true,          // 忽略记录未找到错误
            Colorful:                  true,          // 彩色输出
        },
    )
}

// 性能基准测试
func BenchmarkQueries(db *gorm.DB) {
    tests := []struct {
        name string
        fn   func()
    }{
        {
            name: "创建用户",
            fn: func() {
                user := User{
                    Username: fmt.Sprintf("user_%d", time.Now().UnixNano()),
                    Email:    fmt.Sprintf("user_%d@example.com", time.Now().UnixNano()),
                    Password: "hashed_password",
                }
                db.Create(&user)
            },
        },
        {
            name: "查询文章列表",
            fn: func() {
                var posts []Post
                db.Preload("User").Preload("Category").Limit(10).Find(&posts)
            },
        },
        {
            name: "复杂聚合查询",
            fn: func() {
                var result struct {
                    TotalPosts    int64
                    TotalComments int64
                    ActiveUsers   int64
                }
                db.Raw(`
                    SELECT 
                        (SELECT COUNT(*) FROM post WHERE status = 'published') as total_posts,
                        (SELECT COUNT(*) FROM comment WHERE status = 'approved') as total_comments,
                        (SELECT COUNT(*) FROM user WHERE status = 'active') as active_users
                `).Scan(&result)
            },
        },
    }
    
    for _, test := range tests {
        start := time.Now()
        test.fn()
        duration := time.Since(start)
        fmt.Printf("%s: %v\n", test.name, duration)
    }
}
```

### 连接池监控

```go
// 监控连接池状态
func MonitorConnectionPool(db *gorm.DB) {
    sqlDB, _ := db.DB()
    
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := sqlDB.Stats()
        
        fmt.Printf("连接池状态:\n")
        fmt.Printf("  打开连接数: %d\n", stats.OpenConnections)
        fmt.Printf("  使用中连接数: %d\n", stats.InUse)
        fmt.Printf("  空闲连接数: %d\n", stats.Idle)
        fmt.Printf("  等待连接数: %d\n", stats.WaitCount)
        fmt.Printf("  等待时长: %v\n", stats.WaitDuration)
        fmt.Printf("  最大空闲关闭数: %d\n", stats.MaxIdleClosed)
        fmt.Printf("  最大生命周期关闭数: %d\n", stats.MaxLifetimeClosed)
        fmt.Println("---")
    }
}
```

## 💡 最佳实践

### 1. 事务使用原则

```go
// ✅ 正确的事务使用
func CreatePostWithTransaction(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 所有操作都在事务中
        if err := tx.Create(&post).Error; err != nil {
            return err // 自动回滚
        }
        
        if err := tx.Create(&meta).Error; err != nil {
            return err // 自动回滚
        }
        
        return nil // 自动提交
    })
}

// ❌ 错误的事务使用
func CreatePostWithoutTransaction(db *gorm.DB) error {
    if err := db.Create(&post).Error; err != nil {
        return err
    }
    
    // 如果这里失败，post已经创建，数据不一致
    if err := db.Create(&meta).Error; err != nil {
        return err
    }
    
    return nil
}
```

### 2. 查询优化技巧

```go
// ✅ 使用索引的查询
func GetPostsByStatus(status string) []Post {
    var posts []Post
    // status字段有索引
    db.Where("status = ?", status).Find(&posts)
    return posts
}

// ❌ 避免全表扫描
func GetPostsByTitle(title string) []Post {
    var posts []Post
    // LIKE查询可能导致全表扫描
    db.Where("title LIKE ?", "%"+title+"%").Find(&posts)
    return posts
}

// ✅ 优化的模糊查询
func GetPostsByTitleOptimized(title string) []Post {
    var posts []Post
    // 使用全文索引或搜索引擎
    db.Where("MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE)", title).Find(&posts)
    return posts
}
```

### 3. 内存优化

```go
// ✅ 批量处理大数据
func ProcessLargeDataset(db *gorm.DB) {
    batchSize := 1000
    offset := 0
    
    for {
        var posts []Post
        result := db.Offset(offset).Limit(batchSize).Find(&posts)
        
        if result.Error != nil {
            break
        }
        
        if len(posts) == 0 {
            break
        }
        
        // 处理当前批次
        for _, post := range posts {
            processPost(post)
        }
        
        offset += batchSize
    }
}

// ❌ 避免一次性加载大量数据
func ProcessLargeDatasetBad(db *gorm.DB) {
    var posts []Post
    db.Find(&posts) // 可能导致内存溢出
    
    for _, post := range posts {
        processPost(post)
    }
}
```

### 4. 错误处理

```go
// ✅ 完善的错误处理
func CreateUserSafely(db *gorm.DB, user *User) error {
    if err := db.Create(user).Error; err != nil {
        // 检查具体错误类型
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return fmt.Errorf("用户名或邮箱已存在")
        }
        
        // 记录详细错误信息
        log.Printf("创建用户失败: %v", err)
        return fmt.Errorf("创建用户失败")
    }
    
    return nil
}

// ❌ 简单的错误处理
func CreateUserUnsafely(db *gorm.DB, user *User) error {
    return db.Create(user).Error // 直接返回原始错误
}
```

## 🔧 故障排查

### 常见问题及解决方案

#### 1. 连接池耗尽

**症状**：应用响应缓慢，出现连接超时错误

**排查**：
```go
// 检查连接池状态
sqlDB, _ := db.DB()
stats := sqlDB.Stats()
fmt.Printf("连接池状态: %+v\n", stats)
```

**解决**：
- 调整连接池参数
- 检查是否有连接泄露
- 优化长时间运行的查询

#### 2. 慢查询问题

**症状**：某些查询执行时间过长

**排查**：
```go
// 启用慢查询日志
db.Logger = logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold: 200 * time.Millisecond, // 200ms以上的查询
        LogLevel:      logger.Warn,
    },
)
```

**解决**：
- 添加合适的索引
- 优化查询条件
- 使用EXPLAIN分析查询计划

#### 3. 内存泄露

**症状**：应用内存使用持续增长

**排查**：
```go
// 监控内存使用
func MonitorMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("内存使用: %d KB\n", m.Alloc/1024)
    fmt.Printf("系统内存: %d KB\n", m.Sys/1024)
    fmt.Printf("GC次数: %d\n", m.NumGC)
}
```

**解决**：
- 避免一次性加载大量数据
- 及时释放不需要的对象
- 使用分页或流式处理

#### 4. 死锁问题

**症状**：事务执行卡住，最终超时

**排查**：
```sql
-- MySQL死锁检查
SHOW ENGINE INNODB STATUS;

-- 查看当前锁等待
SELECT * FROM information_schema.INNODB_LOCKS;
SELECT * FROM information_schema.INNODB_LOCK_WAITS;
```

**解决**：
- 保持事务简短
- 按相同顺序访问资源
- 使用合适的隔离级别

### 性能调优检查清单

- [ ] 数据库连接池配置合理
- [ ] 关键查询字段已建立索引
- [ ] 避免N+1查询问题
- [ ] 大数据量操作使用分页
- [ ] 事务范围最小化
- [ ] 定期分析慢查询日志
- [ ] 监控数据库性能指标
- [ ] 定期清理无用数据

---

**本指南涵盖了GORM Level 6项目的核心技术实现，建议结合实际代码进行学习和实践。**
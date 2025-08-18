# 数据库设计文档 🗄️

## 概述 📋

本博客系统采用关系型数据库设计，支持 MySQL 和 SQLite 两种数据库。数据库设计遵循第三范式，具备良好的数据完整性和查询性能。系统使用 GORM v2 作为 ORM 框架，支持自动迁移和关联查询。

## 数据库配置 ⚙️

### 支持的数据库类型

#### MySQL 配置
```go
type DatabaseConfig struct {
    Type     string `json:"type" yaml:"type"`         // "mysql"
    Host     string `json:"host" yaml:"host"`         // "localhost"
    Port     int    `json:"port" yaml:"port"`         // 3306
    Username string `json:"username" yaml:"username"` // "root"
    Password string `json:"password" yaml:"password"` // "password"
    Database string `json:"database" yaml:"database"` // "blog_system"
    Charset  string `json:"charset" yaml:"charset"`   // "utf8mb4"
}
```

#### SQLite 配置
```go
type DatabaseConfig struct {
    Type string `json:"type" yaml:"type"` // "sqlite"
    Path string `json:"path" yaml:"path"` // "./blog.db"
}
```

## 数据库 ER 图 📊

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│      Users      │       │      Posts      │       │    Comments     │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ ID (PK)         │◄──────┤ AuthorID (FK)   │◄──────┤ PostID (FK)     │
│ Username        │       │ ID (PK)         │       │ ID (PK)         │
│ Email           │       │ Title           │       │ UserID (FK)     │
│ PasswordHash    │       │ Content         │       │ ParentID (FK)   │
│ Status          │       │ CategoryID (FK) │       │ Content         │
│ CreatedAt       │       │ Status          │       │ Status          │
│ UpdatedAt       │       │ CreatedAt       │       │ Level           │
│ DeletedAt       │       │ UpdatedAt       │       │ LikeCount       │
└─────────────────┘       │ DeletedAt       │       │ CreatedAt       │
         │                │ PublishedAt     │       │ UpdatedAt       │
         │                └─────────────────┘       │ DeletedAt       │
         │                         │                └─────────────────┘
         │                         │                         │
         │                ┌─────────────────┐                │
         │                │   Categories    │                │
         │                ├─────────────────┤                │
         │                │ ID (PK)         │                │
         │                │ Name            │                │
         │                │ Description     │                │
         │                │ CreatedAt       │                │
         │                │ UpdatedAt       │                │
         │                └─────────────────┘                │
         │                                                   │
         │                ┌─────────────────┐                │
         │                │      Tags       │                │
         │                ├─────────────────┤                │
         │                │ ID (PK)         │                │
         │                │ Name            │                │
         │                │ CreatedAt       │                │
         │                │ UpdatedAt       │                │
         │                └─────────────────┘                │
         │                         │                         │
         │                ┌─────────────────┐                │
         │                │   PostTags      │                │
         │                ├─────────────────┤                │
         │                │ PostID (FK)     │                │
         │                │ TagID (FK)      │                │
         │                └─────────────────┘                │
         │                                                   │
         │                ┌─────────────────┐                │
         └────────────────┤     Follows     │                │
                          ├─────────────────┤                │
                          │ FollowerID (FK) │                │
                          │ FollowingID(FK) │                │
                          │ CreatedAt       │                │
                          └─────────────────┘                │
                                                             │
                          ┌─────────────────┐                │
                          │  UserProfiles   │                │
                          ├─────────────────┤                │
                          │ UserID (FK)     │◄───────────────┘
                          │ Avatar          │
                          │ Bio             │
                          │ Website         │
                          │ Location        │
                          │ Birthday        │
                          │ CreatedAt       │
                          │ UpdatedAt       │
                          └─────────────────┘
```

## 数据表详细设计 📋

### 1. 用户表 (users)

**表名**: `users`
**描述**: 存储用户基本信息

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| id | BIGINT | - | PK, AUTO_INCREMENT | - | 用户ID |
| username | VARCHAR | 50 | UNIQUE, NOT NULL | - | 用户名 |
| email | VARCHAR | 100 | UNIQUE, NOT NULL | - | 邮箱地址 |
| password_hash | VARCHAR | 255 | NOT NULL | - | 密码哈希 |
| status | TINYINT | - | NOT NULL | 1 | 用户状态 |
| last_login_at | TIMESTAMP | - | NULL | NULL | 最后登录时间 |
| login_count | INT | - | NOT NULL | 0 | 登录次数 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | - | NULL | NULL | 软删除时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (id)

-- 唯一索引
UNIQUE KEY uk_username (username)
UNIQUE KEY uk_email (email)

-- 普通索引
KEY idx_status (status)
KEY idx_created_at (created_at)
KEY idx_deleted_at (deleted_at)
```

**用户状态枚举**:
```go
type UserStatus int

const (
    UserStatusInactive UserStatus = 0 // 未激活
    UserStatusActive   UserStatus = 1 // 正常
    UserStatusSuspended UserStatus = 2 // 暂停
    UserStatusBanned   UserStatus = 3 // 封禁
)
```

### 2. 用户资料表 (user_profiles)

**表名**: `user_profiles`
**描述**: 存储用户详细资料信息

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| user_id | BIGINT | - | PK, FK | - | 用户ID |
| avatar | VARCHAR | 255 | NULL | NULL | 头像URL |
| bio | TEXT | - | NULL | NULL | 个人简介 |
| website | VARCHAR | 255 | NULL | NULL | 个人网站 |
| location | VARCHAR | 100 | NULL | NULL | 所在地 |
| birthday | DATE | - | NULL | NULL | 生日 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (user_id)

-- 外键约束
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
```

### 3. 关注关系表 (follows)

**表名**: `follows`
**描述**: 存储用户关注关系

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| follower_id | BIGINT | - | PK, FK | - | 关注者ID |
| following_id | BIGINT | - | PK, FK | - | 被关注者ID |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 关注时间 |

**索引设计**:
```sql
-- 复合主键
PRIMARY KEY (follower_id, following_id)

-- 外键约束
FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE

-- 普通索引
KEY idx_following_id (following_id)
KEY idx_created_at (created_at)
```

### 4. 分类表 (categories)

**表名**: `categories`
**描述**: 存储文章分类信息

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| id | BIGINT | - | PK, AUTO_INCREMENT | - | 分类ID |
| name | VARCHAR | 50 | UNIQUE, NOT NULL | - | 分类名称 |
| description | TEXT | - | NULL | NULL | 分类描述 |
| post_count | INT | - | NOT NULL | 0 | 文章数量 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (id)

-- 唯一索引
UNIQUE KEY uk_name (name)

-- 普通索引
KEY idx_post_count (post_count)
```

### 5. 标签表 (tags)

**表名**: `tags`
**描述**: 存储文章标签信息

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| id | BIGINT | - | PK, AUTO_INCREMENT | - | 标签ID |
| name | VARCHAR | 30 | UNIQUE, NOT NULL | - | 标签名称 |
| post_count | INT | - | NOT NULL | 0 | 使用次数 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (id)

-- 唯一索引
UNIQUE KEY uk_name (name)

-- 普通索引
KEY idx_post_count (post_count)
```

### 6. 文章表 (posts)

**表名**: `posts`
**描述**: 存储文章主要信息

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| id | BIGINT | - | PK, AUTO_INCREMENT | - | 文章ID |
| title | VARCHAR | 200 | NOT NULL | - | 文章标题 |
| content | LONGTEXT | - | NOT NULL | - | 文章内容 |
| summary | TEXT | - | NULL | NULL | 文章摘要 |
| author_id | BIGINT | - | FK, NOT NULL | - | 作者ID |
| category_id | BIGINT | - | FK, NULL | NULL | 分类ID |
| status | TINYINT | - | NOT NULL | 1 | 文章状态 |
| priority | TINYINT | - | NOT NULL | 0 | 优先级 |
| view_count | INT | - | NOT NULL | 0 | 浏览次数 |
| like_count | INT | - | NOT NULL | 0 | 点赞次数 |
| comment_count | INT | - | NOT NULL | 0 | 评论次数 |
| published_at | TIMESTAMP | - | NULL | NULL | 发布时间 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | - | NULL | NULL | 软删除时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (id)

-- 外键约束
FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL

-- 复合索引
KEY idx_author_status (author_id, status)
KEY idx_category_status (category_id, status)
KEY idx_status_published (status, published_at)

-- 普通索引
KEY idx_view_count (view_count)
KEY idx_like_count (like_count)
KEY idx_created_at (created_at)
KEY idx_deleted_at (deleted_at)

-- 全文索引
FULLTEXT KEY ft_title_content (title, content)
```

**文章状态枚举**:
```go
type PostStatus int

const (
    PostStatusDraft     PostStatus = 1 // 草稿
    PostStatusPublished PostStatus = 2 // 已发布
    PostStatusArchived  PostStatus = 3 // 已归档
    PostStatusDeleted   PostStatus = 4 // 已删除
)
```

**优先级枚举**:
```go
type Priority int

const (
    PriorityLow    Priority = 0 // 低
    PriorityNormal Priority = 1 // 普通
    PriorityHigh   Priority = 2 // 高
    PriorityTop    Priority = 3 // 置顶
)
```

### 7. 文章标签关联表 (post_tags)

**表名**: `post_tags`
**描述**: 存储文章和标签的多对多关系

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| post_id | BIGINT | - | PK, FK | - | 文章ID |
| tag_id | BIGINT | - | PK, FK | - | 标签ID |

**索引设计**:
```sql
-- 复合主键
PRIMARY KEY (post_id, tag_id)

-- 外键约束
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE

-- 普通索引
KEY idx_tag_id (tag_id)
```

### 8. 文章元数据表 (post_meta)

**表名**: `post_meta`
**描述**: 存储文章的扩展元数据

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| id | BIGINT | - | PK, AUTO_INCREMENT | - | 元数据ID |
| post_id | BIGINT | - | FK, NOT NULL | - | 文章ID |
| meta_key | VARCHAR | 100 | NOT NULL | - | 元数据键 |
| meta_value | TEXT | - | NULL | NULL | 元数据值 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (id)

-- 外键约束
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE

-- 复合索引
UNIQUE KEY uk_post_meta (post_id, meta_key)

-- 普通索引
KEY idx_meta_key (meta_key)
```

### 9. 评论表 (comments)

**表名**: `comments`
**描述**: 存储文章评论信息

| 字段名 | 类型 | 长度 | 约束 | 默认值 | 描述 |
|--------|------|------|------|--------|---------|
| id | BIGINT | - | PK, AUTO_INCREMENT | - | 评论ID |
| post_id | BIGINT | - | FK, NOT NULL | - | 文章ID |
| user_id | BIGINT | - | FK, NOT NULL | - | 用户ID |
| parent_id | BIGINT | - | FK, NULL | NULL | 父评论ID |
| content | TEXT | - | NOT NULL | - | 评论内容 |
| status | TINYINT | - | NOT NULL | 1 | 评论状态 |
| level | TINYINT | - | NOT NULL | 1 | 评论层级 |
| like_count | INT | - | NOT NULL | 0 | 点赞次数 |
| reply_count | INT | - | NOT NULL | 0 | 回复次数 |
| created_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | - | NULL | NULL | 软删除时间 |

**索引设计**:
```sql
-- 主键索引
PRIMARY KEY (id)

-- 外键约束
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE

-- 复合索引
KEY idx_post_status (post_id, status)
KEY idx_parent_status (parent_id, status)
KEY idx_user_status (user_id, status)

-- 普通索引
KEY idx_level (level)
KEY idx_like_count (like_count)
KEY idx_created_at (created_at)
KEY idx_deleted_at (deleted_at)
```

**评论状态枚举**:
```go
type CommentStatus int

const (
    CommentStatusPending  CommentStatus = 1 // 待审核
    CommentStatusApproved CommentStatus = 2 // 已通过
    CommentStatusRejected CommentStatus = 3 // 已拒绝
    CommentStatusSpam     CommentStatus = 4 // 垃圾评论
)
```

## GORM 模型定义 🔧

### 基础模型

```go
// BaseModel 基础模型，包含公共字段
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 实现表名接口
type TableName interface {
    TableName() string
}

// 实现软删除接口
type SoftDeletable interface {
    GetDeletedAt() *time.Time
    SetDeletedAt(time.Time)
}

// 实现时间戳接口
type Timestampable interface {
    GetCreatedAt() time.Time
    GetUpdatedAt() time.Time
    SetUpdatedAt(time.Time)
}
```

### 用户模型

```go
// User 用户模型
type User struct {
    BaseModel
    Username    string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email       string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    Status      UserStatus `gorm:"not null;default:1" json:"status"`
    LastLoginAt *time.Time `json:"last_login_at,omitempty"`
    LoginCount  int        `gorm:"not null;default:0" json:"login_count"`
    
    // 关联关系
    Profile   *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
    Posts     []Post       `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
    Comments  []Comment    `gorm:"foreignKey:UserID" json:"comments,omitempty"`
    Followers []Follow     `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
    Following []Follow     `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
}

func (User) TableName() string {
    return "users"
}

// UserProfile 用户资料模型
type UserProfile struct {
    UserID   uint       `gorm:"primaryKey" json:"user_id"`
    Avatar   string     `gorm:"size:255" json:"avatar,omitempty"`
    Bio      string     `gorm:"type:text" json:"bio,omitempty"`
    Website  string     `gorm:"size:255" json:"website,omitempty"`
    Location string     `gorm:"size:100" json:"location,omitempty"`
    Birthday *time.Time `gorm:"type:date" json:"birthday,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // 关联关系
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserProfile) TableName() string {
    return "user_profiles"
}

// Follow 关注关系模型
type Follow struct {
    FollowerID  uint      `gorm:"primaryKey" json:"follower_id"`
    FollowingID uint      `gorm:"primaryKey" json:"following_id"`
    CreatedAt   time.Time `json:"created_at"`
    
    // 关联关系
    Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
    Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

func (Follow) TableName() string {
    return "follows"
}
```

### 文章模型

```go
// Post 文章模型
type Post struct {
    BaseModel
    Title        string     `gorm:"size:200;not null" json:"title"`
    Content      string     `gorm:"type:longtext;not null" json:"content"`
    Summary      string     `gorm:"type:text" json:"summary,omitempty"`
    AuthorID     uint       `gorm:"not null;index" json:"author_id"`
    CategoryID   *uint      `gorm:"index" json:"category_id,omitempty"`
    Status       PostStatus `gorm:"not null;default:1;index" json:"status"`
    Priority     Priority   `gorm:"not null;default:0" json:"priority"`
    ViewCount    int        `gorm:"not null;default:0;index" json:"view_count"`
    LikeCount    int        `gorm:"not null;default:0;index" json:"like_count"`
    CommentCount int        `gorm:"not null;default:0" json:"comment_count"`
    PublishedAt  *time.Time `gorm:"index" json:"published_at,omitempty"`
    
    // 关联关系
    Author   User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    Category *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    Tags     []Tag      `gorm:"many2many:post_tags" json:"tags,omitempty"`
    Comments []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`
    Meta     []PostMeta `gorm:"foreignKey:PostID" json:"meta,omitempty"`
}

func (Post) TableName() string {
    return "posts"
}

// Category 分类模型
type Category struct {
    BaseModel
    Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"`
    Description string `gorm:"type:text" json:"description,omitempty"`
    PostCount   int    `gorm:"not null;default:0" json:"post_count"`
    
    // 关联关系
    Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

func (Category) TableName() string {
    return "categories"
}

// Tag 标签模型
type Tag struct {
    BaseModel
    Name      string `gorm:"uniqueIndex;size:30;not null" json:"name"`
    PostCount int    `gorm:"not null;default:0" json:"post_count"`
    
    // 关联关系
    Posts []Post `gorm:"many2many:post_tags" json:"posts,omitempty"`
}

func (Tag) TableName() string {
    return "tags"
}

// PostMeta 文章元数据模型
type PostMeta struct {
    BaseModel
    PostID    uint   `gorm:"not null;index" json:"post_id"`
    MetaKey   string `gorm:"size:100;not null" json:"meta_key"`
    MetaValue string `gorm:"type:text" json:"meta_value,omitempty"`
    
    // 关联关系
    Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

func (PostMeta) TableName() string {
    return "post_meta"
}
```

### 评论模型

```go
// Comment 评论模型
type Comment struct {
    BaseModel
    PostID     uint          `gorm:"not null;index" json:"post_id"`
    UserID     uint          `gorm:"not null;index" json:"user_id"`
    ParentID   *uint         `gorm:"index" json:"parent_id,omitempty"`
    Content    string        `gorm:"type:text;not null" json:"content"`
    Status     CommentStatus `gorm:"not null;default:1;index" json:"status"`
    Level      int           `gorm:"not null;default:1" json:"level"`
    LikeCount  int           `gorm:"not null;default:0" json:"like_count"`
    ReplyCount int           `gorm:"not null;default:0" json:"reply_count"`
    
    // 关联关系
    Post    Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
    User    User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

func (Comment) TableName() string {
    return "comments"
}
```

## 数据库迁移 🔄

### 自动迁移

```go
// 数据库自动迁移
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &UserProfile{},
        &Follow{},
        &Category{},
        &Tag{},
        &Post{},
        &PostMeta{},
        &Comment{},
    )
}
```

### 手动迁移脚本

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS blog_system 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE blog_system;

-- 创建用户表
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status TINYINT NOT NULL DEFAULT 1,
    last_login_at TIMESTAMP NULL,
    login_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 其他表的创建脚本...
```

## 查询优化策略 ⚡

### 索引设计原则

1. **主键索引**：每个表都有自增主键
2. **唯一索引**：用户名、邮箱等唯一字段
3. **复合索引**：多字段组合查询
4. **外键索引**：关联查询优化
5. **全文索引**：文章标题和内容搜索

### 常用查询优化

```go
// 1. 预加载关联数据
db.Preload("Author").Preload("Category").Preload("Tags").Find(&posts)

// 2. 选择特定字段
db.Select("id", "title", "summary", "created_at").Find(&posts)

// 3. 分页查询
db.Offset(offset).Limit(limit).Find(&posts)

// 4. 条件查询
db.Where("status = ? AND author_id = ?", PostStatusPublished, authorID).Find(&posts)

// 5. 排序查询
db.Order("created_at DESC").Find(&posts)

// 6. 统计查询
db.Model(&Post{}).Where("status = ?", PostStatusPublished).Count(&count)
```

### 性能监控

```go
// 启用 SQL 日志
db.Logger = logger.Default.LogMode(logger.Info)

// 慢查询监控
db.Logger = logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold: time.Second,   // 慢查询阈值
        LogLevel:      logger.Warn,  // 日志级别
        Colorful:      true,         // 彩色输出
    },
)
```

## 数据完整性 🔒

### 外键约束

```sql
-- 文章作者外键
ALTER TABLE posts 
ADD CONSTRAINT fk_posts_author 
FOREIGN KEY (author_id) REFERENCES users(id) 
ON DELETE CASCADE;

-- 评论文章外键
ALTER TABLE comments 
ADD CONSTRAINT fk_comments_post 
FOREIGN KEY (post_id) REFERENCES posts(id) 
ON DELETE CASCADE;
```

### 数据验证

```go
// GORM 钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 验证邮箱格式
    if !isValidEmail(u.Email) {
        return errors.New("邮箱格式不正确")
    }
    
    // 验证用户名长度
    if len(u.Username) < 3 || len(u.Username) > 50 {
        return errors.New("用户名长度必须在3-50个字符之间")
    }
    
    return nil
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
    // 验证标题长度
    if len(p.Title) == 0 || len(p.Title) > 200 {
        return errors.New("标题长度必须在1-200个字符之间")
    }
    
    // 自动设置发布时间
    if p.Status == PostStatusPublished && p.PublishedAt == nil {
        now := time.Now()
        p.PublishedAt = &now
    }
    
    return nil
}
```

## 备份与恢复 💾

### 数据备份脚本

```bash
#!/bin/bash
# MySQL 数据备份脚本

DB_NAME="blog_system"
DB_USER="root"
DB_PASS="password"
BACKUP_DIR="/backup/mysql"
DATE=$(date +"%Y%m%d_%H%M%S")

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行备份
mysqldump -u$DB_USER -p$DB_PASS $DB_NAME > $BACKUP_DIR/blog_system_$DATE.sql

# 压缩备份文件
gzip $BACKUP_DIR/blog_system_$DATE.sql

# 删除7天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "数据库备份完成: blog_system_$DATE.sql.gz"
```

### 数据恢复

```bash
#!/bin/bash
# MySQL 数据恢复脚本

DB_NAME="blog_system"
DB_USER="root"
DB_PASS="password"
BACKUP_FILE="$1"

if [ -z "$BACKUP_FILE" ]; then
    echo "请指定备份文件路径"
    exit 1
fi

# 解压备份文件（如果是压缩的）
if [[ $BACKUP_FILE == *.gz ]]; then
    gunzip $BACKUP_FILE
    BACKUP_FILE=${BACKUP_FILE%.gz}
fi

# 恢复数据库
mysql -u$DB_USER -p$DB_PASS $DB_NAME < $BACKUP_FILE

echo "数据库恢复完成"
```

---

**注意**：本文档描述的数据库设计支持高并发访问和大数据量存储，建议在生产环境中根据实际需求调整索引和分区策略。
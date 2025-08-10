# GORM背景示例：博客系统实战

## 📖 项目背景

我们将通过构建一个完整的博客系统来学习GORM的各种功能。这个系统包含用户管理、文章发布、评论互动等核心功能，是学习GORM的绝佳实例。

### 🎯 系统功能需求
- 用户注册、登录、个人资料管理
- 文章的创建、编辑、删除、发布
- 文章分类和标签管理
- 评论系统（支持回复）
- 点赞和收藏功能
- 文章搜索和分页

### 🗄️ 数据库设计

```
用户表(users) ←→ 文章表(posts) ←→ 评论表(comments)
     ↓              ↓              ↓
个人资料表(profiles) 标签表(tags)   点赞表(likes)
                    ↓
                分类表(categories)
```

---

## 🏗️ 项目结构设计

```
blog-system/
├── main.go              # 程序入口
├── config/
│   └── database.go      # 数据库配置
├── models/
│   ├── user.go          # 用户模型
│   ├── post.go          # 文章模型
│   ├── comment.go       # 评论模型
│   ├── category.go      # 分类模型
│   ├── tag.go           # 标签模型
│   └── like.go          # 点赞模型
├── services/
│   ├── user_service.go  # 用户服务
│   ├── post_service.go  # 文章服务
│   └── comment_service.go # 评论服务
├── handlers/
│   ├── user_handler.go  # 用户控制器
│   ├── post_handler.go  # 文章控制器
│   └── comment_handler.go # 评论控制器
└── utils/
    ├── response.go      # 响应工具
    └── pagination.go    # 分页工具
```

---

## 📋 第一步：环境搭建和初始化

### 1.1 项目初始化
```bash
# 创建项目目录
mkdir blog-system
cd blog-system

# 初始化Go模块
go mod init blog-system

# 安装依赖
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/gin-gonic/gin
go get -u golang.org/x/crypto/bcrypt
```

### 1.2 数据库配置
```go
// config/database.go
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

// InitDatabase 初始化数据库连接
func InitDatabase() {
	config := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "blog_system"),
		Charset:  "utf8mb4",
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
	)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
		// 禁用外键约束（可选，根据需求决定）
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
```

---

## 🏛️ 第二步：模型定义

### 2.1 用户模型
```go
// models/user.go
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
	
	// 基本信息
	Username string `json:"username" gorm:"size:50;uniqueIndex;not null" validate:"required,min=3,max=50"`
	Email    string `json:"email" gorm:"size:100;uniqueIndex;not null" validate:"required,email"`
	Password string `json:"-" gorm:"size:255;not null" validate:"required,min=6"` // 不在JSON中显示
	
	// 个人信息
	Nickname string `json:"nickname" gorm:"size:50"`
	Avatar   string `json:"avatar" gorm:"size:255"`
	Bio      string `json:"bio" gorm:"size:500"`
	
	// 状态信息
	Status    string `json:"status" gorm:"size:20;default:active"` // active, inactive, banned
	IsAdmin   bool   `json:"is_admin" gorm:"default:false"`
	LastLogin *time.Time `json:"last_login"`
	
	// 关联关系
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:AuthorID"` // 一对多：用户的文章
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"` // 一对多：用户的评论
	Likes    []Like    `json:"likes,omitempty" gorm:"foreignKey:UserID"`   // 一对多：用户的点赞
}

// BeforeCreate 创建前钩子：密码加密
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate 更新前钩子：如果密码被修改，重新加密
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 检查密码是否被修改
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// AfterCreate 创建后钩子：记录日志
func (u *User) AfterCreate(tx *gorm.DB) error {
	// 这里可以添加日志记录、发送欢迎邮件等逻辑
	log.Printf("New user created: %s (ID: %d)", u.Username, u.ID)
	return nil
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
```

### 2.2 文章模型
```go
// models/post.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 基本信息
	Title   string `json:"title" gorm:"size:200;not null;index" validate:"required,max=200"`
	Slug    string `json:"slug" gorm:"size:200;uniqueIndex" validate:"required,max=200"`
	Content string `json:"content" gorm:"type:longtext;not null" validate:"required"`
	Excerpt string `json:"excerpt" gorm:"size:500"` // 摘要
	
	// 状态信息
	Status      string     `json:"status" gorm:"size:20;default:draft;index"` // draft, published, archived
	PublishedAt *time.Time `json:"published_at" gorm:"index"`
	
	// 统计信息
	ViewCount    int `json:"view_count" gorm:"default:0"`
	LikeCount    int `json:"like_count" gorm:"default:0"`
	CommentCount int `json:"comment_count" gorm:"default:0"`
	
	// 外键
	AuthorID   uint `json:"author_id" gorm:"not null;index"`
	CategoryID uint `json:"category_id" gorm:"index"`
	
	// 关联关系
	Author   User     `json:"author" gorm:"foreignKey:AuthorID"`      // 属于：文章的作者
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`   // 属于：文章的分类
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"` // 一对多：文章的评论
	Likes    []Like    `json:"likes,omitempty" gorm:"foreignKey:PostID"`    // 一对多：文章的点赞
	Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:post_tags;"`  // 多对多：文章的标签
}

// BeforeCreate 创建前钩子
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 如果没有设置摘要，自动生成
	if p.Excerpt == "" && len(p.Content) > 200 {
		p.Excerpt = p.Content[:200] + "..."
	}
	
	// 如果状态是发布，设置发布时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	
	return nil
}

// BeforeUpdate 更新前钩子
func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	// 如果状态改为发布，设置发布时间
	if tx.Statement.Changed("Status") && p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}

// AfterUpdate 更新后钩子：更新统计信息
func (p *Post) AfterUpdate(tx *gorm.DB) error {
	// 更新评论数量
	var commentCount int64
	tx.Model(&Comment{}).Where("post_id = ?", p.ID).Count(&commentCount)
	tx.Model(p).Update("comment_count", commentCount)
	
	// 更新点赞数量
	var likeCount int64
	tx.Model(&Like{}).Where("post_id = ? AND type = ?", p.ID, "post").Count(&likeCount)
	tx.Model(p).Update("like_count", likeCount)
	
	return nil
}

// TableName 自定义表名
func (Post) TableName() string {
	return "posts"
}
```

### 2.3 评论模型
```go
// models/comment.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 基本信息
	Content string `json:"content" gorm:"type:text;not null" validate:"required,max=1000"`
	
	// 状态信息
	Status    string `json:"status" gorm:"size:20;default:approved"` // pending, approved, rejected
	LikeCount int    `json:"like_count" gorm:"default:0"`
	
	// 外键
	UserID   uint  `json:"user_id" gorm:"not null;index"`
	PostID   uint  `json:"post_id" gorm:"not null;index"`
	ParentID *uint `json:"parent_id" gorm:"index"` // 父评论ID，用于回复
	
	// 关联关系
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Post     Post      `json:"post" gorm:"foreignKey:PostID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`   // 父评论
	Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"` // 子评论
	Likes    []Like    `json:"likes,omitempty" gorm:"foreignKey:CommentID"`  // 评论的点赞
}

// AfterCreate 创建后钩子：更新文章评论数
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	// 更新文章的评论数量
	tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	return nil
}

// AfterDelete 删除后钩子：更新文章评论数
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 更新文章的评论数量
	tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	return nil
}

// TableName 自定义表名
func (Comment) TableName() string {
	return "comments"
}
```

### 2.4 分类和标签模型
```go
// models/category.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	Name        string `json:"name" gorm:"size:50;uniqueIndex;not null"`
	Slug        string `json:"slug" gorm:"size:50;uniqueIndex;not null"`
	Description string `json:"description" gorm:"size:200"`
	Color       string `json:"color" gorm:"size:7;default:#007bff"` // 十六进制颜色
	PostCount   int    `json:"post_count" gorm:"default:0"`
	
	// 关联关系
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName 自定义表名
func (Category) TableName() string {
	return "categories"
}

// models/tag.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	Name      string `json:"name" gorm:"size:30;uniqueIndex;not null"`
	Slug      string `json:"slug" gorm:"size:30;uniqueIndex;not null"`
	Color     string `json:"color" gorm:"size:7;default:#6c757d"`
	PostCount int    `json:"post_count" gorm:"default:0"`
	
	// 关联关系
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags;"`
}

// TableName 自定义表名
func (Tag) TableName() string {
	return "tags"
}
```

### 2.5 点赞模型
```go
// models/like.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Like 点赞模型
type Like struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 外键
	UserID    uint  `json:"user_id" gorm:"not null;index"`
	PostID    *uint `json:"post_id" gorm:"index"`    // 文章点赞
	CommentID *uint `json:"comment_id" gorm:"index"` // 评论点赞
	
	// 点赞类型
	Type string `json:"type" gorm:"size:20;not null"` // post, comment
	
	// 关联关系
	User    User     `json:"user" gorm:"foreignKey:UserID"`
	Post    *Post    `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Comment *Comment `json:"comment,omitempty" gorm:"foreignKey:CommentID"`
}

// BeforeCreate 创建前钩子：验证数据
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	// 确保点赞类型和对应的ID匹配
	if l.Type == "post" && l.PostID == nil {
		return errors.New("post_id is required for post like")
	}
	if l.Type == "comment" && l.CommentID == nil {
		return errors.New("comment_id is required for comment like")
	}
	return nil
}

// TableName 自定义表名
func (Like) TableName() string {
	return "likes"
}
```

---

## 🔧 第三步：数据库迁移

```go
// main.go
package main

import (
	"blog-system/config"
	"blog-system/models"
	"log"
)

func main() {
	// 初始化数据库
	config.InitDatabase()
	db := config.GetDB()
	
	// 自动迁移
	if err := autoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	// 创建初始数据
	if err := seedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}
	
	log.Println("Database migration completed successfully!")
}

// autoMigrate 自动迁移数据库
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	)
}

// seedData 创建初始数据
func seedData(db *gorm.DB) error {
	// 创建默认分类
	categories := []models.Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", Color: "#007bff"},
		{Name: "生活", Slug: "life", Description: "生活感悟分享", Color: "#28a745"},
		{Name: "随笔", Slug: "essay", Description: "随笔杂谈", Color: "#ffc107"},
	}
	
	for _, category := range categories {
		db.FirstOrCreate(&category, models.Category{Slug: category.Slug})
	}
	
	// 创建默认标签
	tags := []models.Tag{
		{Name: "Go", Slug: "go", Color: "#00ADD8"},
		{Name: "GORM", Slug: "gorm", Color: "#FF6B6B"},
		{Name: "数据库", Slug: "database", Color: "#4ECDC4"},
		{Name: "Web开发", Slug: "web-dev", Color: "#45B7D1"},
	}
	
	for _, tag := range tags {
		db.FirstOrCreate(&tag, models.Tag{Slug: tag.Slug})
	}
	
	// 创建管理员用户
	admin := models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin123",
		Nickname: "管理员",
		IsAdmin:  true,
		Status:   "active",
	}
	
	db.FirstOrCreate(&admin, models.User{Username: "admin"})
	
	return nil
}
```

---

## 🎯 第四步：业务服务层

### 4.1 用户服务
```go
// services/user_service.go
package services

import (
	"blog-system/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *models.User) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		return errors.New("用户名或邮箱已存在")
	}
	
	// 创建用户
	return s.db.Create(user).Error
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateLastLogin 更新最后登录时间
func (s *UserService) UpdateLastLogin(id uint) error {
	now := time.Now()
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", &now).Error
}

// GetUserPosts 获取用户的文章列表
func (s *UserService) GetUserPosts(userID uint, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	
	// 计算总数
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := s.db.Where("author_id = ?", userID).
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error
	
	return posts, total, err
}
```

### 4.2 文章服务
```go
// services/post_service.go
package services

import (
	"blog-system/models"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// CreatePost 创建文章
func (s *PostService) CreatePost(post *models.Post) error {
	// 检查slug是否已存在
	var existingPost models.Post
	if err := s.db.Where("slug = ?", post.Slug).First(&existingPost).Error; err == nil {
		return errors.New("文章slug已存在")
	}
	
	return s.db.Create(post).Error
}

// GetPostByID 根据ID获取文章
func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := s.db.Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments.User").
		Preload("Comments.Replies.User").
		First(&post, id).Error
	
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostBySlug 根据slug获取文章
func (s *PostService) GetPostBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := s.db.Where("slug = ?", slug).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		First(&post, slug).Error
	
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPosts 获取文章列表（支持搜索、分类、标签过滤）
func (s *PostService) GetPosts(params PostQueryParams) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64
	
	// 构建查询
	query := s.db.Model(&models.Post{})
	
	// 只查询已发布的文章
	query = query.Where("status = ?", "published")
	
	// 搜索条件
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", searchTerm, searchTerm)
	}
	
	// 分类过滤
	if params.CategoryID > 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	
	// 标签过滤
	if params.TagID > 0 {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ?", params.TagID)
	}
	
	// 作者过滤
	if params.AuthorID > 0 {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	
	// 计算总数
	query.Count(&total)
	
	// 排序
	orderBy := "created_at DESC"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}
	query = query.Order(orderBy)
	
	// 分页
	offset := (params.Page - 1) * params.PageSize
	query = query.Limit(params.PageSize).Offset(offset)
	
	// 预加载关联数据
	err := query.Preload("Author").
		Preload("Category").
		Preload("Tags").
		Find(&posts).Error
	
	return posts, total, err
}

// PostQueryParams 文章查询参数
type PostQueryParams struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	Search     string `form:"search"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	AuthorID   uint   `form:"author_id"`
	OrderBy    string `form:"order_by"`
}

// UpdatePost 更新文章
func (s *PostService) UpdatePost(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePost 删除文章
func (s *PostService) DeletePost(id uint) error {
	return s.db.Delete(&models.Post{}, id).Error
}

// IncrementViewCount 增加浏览量
func (s *PostService) IncrementViewCount(id uint) error {
	return s.db.Model(&models.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// PublishPost 发布文章
func (s *PostService) PublishPost(id uint) error {
	now := time.Now()
	return s.db.Model(&models.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": &now,
	}).Error
}
```

---

## 🎮 第五步：控制器层

### 5.1 文章控制器
```go
// handlers/post_handler.go
package handlers

import (
	"blog-system/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// GetPosts 获取文章列表
func (h *PostHandler) GetPosts(c *gin.Context) {
	// 解析查询参数
	var params services.PostQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 设置默认值
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	
	// 获取文章列表
	posts, total, err := h.postService.GetPosts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 计算分页信息
	totalPages := (int(total) + params.PageSize - 1) / params.PageSize
	
	c.JSON(http.StatusOK, gin.H{
		"data": posts,
		"pagination": gin.H{
			"page":        params.Page,
			"page_size":   params.PageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetPost 获取单篇文章
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	// 增加浏览量
	h.postService.IncrementViewCount(uint(id))
	
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// CreatePost 创建文章
func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 从JWT中获取用户ID（这里简化处理）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	post.AuthorID = userID.(uint)
	
	if err := h.postService.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": post})
}
```

---

## 📊 第六步：实际运行示例

### 6.1 完整的main.go
```go
// main.go
package main

import (
	"blog-system/config"
	"blog-system/handlers"
	"blog-system/models"
	"blog-system/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDatabase()
	db := config.GetDB()
	
	// 自动迁移
	db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	)
	
	// 初始化服务
	userService := services.NewUserService(db)
	postService := services.NewPostService(db)
	
	// 初始化控制器
	postHandler := handlers.NewPostHandler(postService)
	
	// 初始化路由
	r := gin.Default()
	
	// API路由组
	api := r.Group("/api/v1")
	{
		// 文章相关路由
		posts := api.Group("/posts")
		{
			posts.GET("", postHandler.GetPosts)     // 获取文章列表
			posts.GET("/:id", postHandler.GetPost)  // 获取单篇文章
			posts.POST("", postHandler.CreatePost)  // 创建文章
		}
	}
	
	// 启动服务器
	log.Println("Server starting on :8080")
	r.Run(":8080")
}
```

### 6.2 测试API

```bash
# 1. 获取文章列表
curl "http://localhost:8080/api/v1/posts?page=1&page_size=5"

# 2. 搜索文章
curl "http://localhost:8080/api/v1/posts?search=GORM&page=1"

# 3. 按分类筛选
curl "http://localhost:8080/api/v1/posts?category_id=1&page=1"

# 4. 获取单篇文章
curl "http://localhost:8080/api/v1/posts/1"

# 5. 创建文章（需要认证）
curl -X POST "http://localhost:8080/api/v1/posts" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "GORM入门教程",
    "slug": "gorm-tutorial",
    "content": "这是一篇关于GORM的详细教程...",
    "category_id": 1,
    "status": "published"
  }'
```

---

## 🎯 关键学习点总结

### 1. **模型设计最佳实践**
- 合理使用标签定义字段属性
- 正确设置关联关系
- 利用钩子函数实现业务逻辑
- 使用软删除保护数据

### 2. **查询优化技巧**
- 使用预加载避免N+1问题
- 合理使用索引提高查询性能
- 分页查询处理大量数据
- 链式调用构建复杂查询

### 3. **事务处理**
- 自动事务处理错误回滚
- 手动事务精确控制
- 嵌套事务的使用场景

### 4. **性能优化**
- 批量操作提高效率
- 预编译语句减少解析开销
- 连接池管理数据库连接
- SQL日志监控性能瓶颈

### 5. **架构设计**
- 分层架构清晰职责
- 服务层封装业务逻辑
- 控制器层处理HTTP请求
- 模型层定义数据结构

---

这个博客系统示例涵盖了GORM的核心功能和最佳实践，通过实际的业务场景帮助你深入理解GORM的使用方法。在后续的练习中，我们将基于这个示例进行更深入的学习和实践。
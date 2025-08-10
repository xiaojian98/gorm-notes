# 🚀 GORM实战练习：企业级项目开发

## 📋 项目概述

### 项目背景
本实战练习将带你开发一个完整的**企业级在线教育平台后端系统**，涵盖用户管理、课程管理、订单支付、学习进度跟踪、数据统计等核心功能。通过这个项目，你将掌握GORM在真实企业环境中的应用技巧。

### 技术栈
- **后端框架**：Gin + GORM
- **数据库**：MySQL 8.0
- **缓存**：Redis
- **消息队列**：RabbitMQ
- **文件存储**：MinIO
- **监控**：Prometheus + Grafana
- **部署**：Docker + Docker Compose

### 项目特色
- 🏗️ **微服务架构设计**
- 🔐 **完整的权限控制系统**
- 📊 **实时数据统计和监控**
- 🚀 **高性能和高并发处理**
- 🛡️ **完善的错误处理和日志系统**
- 🧪 **全面的单元测试和集成测试**

---

## 🏗️ 系统架构设计

### 整体架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │    │  Mobile App     │    │  Admin Panel    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   API Gateway   │
                    └─────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                       │                       │
┌─────────────┐        ┌─────────────┐        ┌─────────────┐
│ User Service│        │Course Service│       │Order Service│
└─────────────┘        └─────────────┘        └─────────────┘
        │                       │                       │
        └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Data Layer    │
                    │  MySQL + Redis  │
                    └─────────────────┘
```

### 数据库设计

#### 核心表结构

```sql
-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar VARCHAR(255),
    nickname VARCHAR(50),
    gender TINYINT DEFAULT 0 COMMENT '0:未知,1:男,2:女',
    birthday DATE,
    status TINYINT DEFAULT 1 COMMENT '1:正常,2:禁用',
    role_id BIGINT,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- 角色表
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    permissions JSON,
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_name (name)
);

-- 课程分类表
CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id BIGINT DEFAULT 0,
    sort_order INT DEFAULT 0,
    icon VARCHAR(255),
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_parent_id (parent_id),
    INDEX idx_status (status)
);

-- 课程表
CREATE TABLE courses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    subtitle VARCHAR(500),
    description TEXT,
    cover_image VARCHAR(255),
    category_id BIGINT NOT NULL,
    teacher_id BIGINT NOT NULL,
    price DECIMAL(10,2) DEFAULT 0.00,
    original_price DECIMAL(10,2) DEFAULT 0.00,
    difficulty TINYINT DEFAULT 1 COMMENT '1:初级,2:中级,3:高级',
    duration INT DEFAULT 0 COMMENT '课程时长(分钟)',
    student_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    rating_count INT DEFAULT 0,
    status TINYINT DEFAULT 1 COMMENT '1:草稿,2:发布,3:下架',
    is_free TINYINT DEFAULT 0,
    tags JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_category_id (category_id),
    INDEX idx_teacher_id (teacher_id),
    INDEX idx_status (status),
    INDEX idx_price (price),
    INDEX idx_created_at (created_at),
    FULLTEXT idx_title_desc (title, description)
);

-- 章节表
CREATE TABLE chapters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    duration INT DEFAULT 0,
    is_free TINYINT DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_course_id (course_id),
    INDEX idx_sort_order (sort_order)
);

-- 课时表
CREATE TABLE lessons (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    chapter_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    video_url VARCHAR(255),
    video_duration INT DEFAULT 0,
    sort_order INT DEFAULT 0,
    is_free TINYINT DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_chapter_id (chapter_id),
    INDEX idx_course_id (course_id),
    INDEX idx_sort_order (sort_order)
);

-- 订单表
CREATE TABLE orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no VARCHAR(32) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    pay_amount DECIMAL(10,2) NOT NULL,
    payment_method TINYINT DEFAULT 1 COMMENT '1:支付宝,2:微信,3:银行卡',
    payment_status TINYINT DEFAULT 1 COMMENT '1:待支付,2:已支付,3:已退款',
    order_status TINYINT DEFAULT 1 COMMENT '1:待支付,2:已完成,3:已取消',
    paid_at TIMESTAMP NULL,
    expired_at TIMESTAMP NULL,
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_payment_status (payment_status),
    INDEX idx_created_at (created_at)
);

-- 订单详情表
CREATE TABLE order_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    course_title VARCHAR(200) NOT NULL,
    course_cover VARCHAR(255),
    price DECIMAL(10,2) NOT NULL,
    quantity INT DEFAULT 1,
    total_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_order_id (order_id),
    INDEX idx_course_id (course_id)
);

-- 学习进度表
CREATE TABLE learning_progress (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    lesson_id BIGINT NOT NULL,
    progress_percent DECIMAL(5,2) DEFAULT 0.00,
    watch_duration INT DEFAULT 0,
    is_completed TINYINT DEFAULT 0,
    last_watch_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_user_lesson (user_id, lesson_id),
    INDEX idx_user_course (user_id, course_id),
    INDEX idx_last_watch_at (last_watch_at)
);
```

---

## 💻 核心代码实现

### 1. 项目结构设计

```
edu-platform/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── migrate/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── mysql.go
│   │   └── redis.go
│   ├── models/
│   │   ├── user.go
│   │   ├── course.go
│   │   ├── order.go
│   │   └── learning.go
│   ├── repositories/
│   │   ├── user_repo.go
│   │   ├── course_repo.go
│   │   └── order_repo.go
│   ├── services/
│   │   ├── user_service.go
│   │   ├── course_service.go
│   │   └── order_service.go
│   ├── handlers/
│   │   ├── user_handler.go
│   │   ├── course_handler.go
│   │   └── order_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── utils/
│       ├── response.go
│       ├── jwt.go
│       └── validator.go
├── pkg/
│   ├── logger/
│   ├── cache/
│   └── queue/
├── configs/
│   ├── config.yaml
│   └── docker-compose.yml
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── tests/
│   ├── integration/
│   └── unit/
├── docs/
│   └── api.md
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

### 2. 配置管理

```go
// internal/config/config.go
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Master DatabaseConnection   `mapstructure:"master"`
	Slaves []DatabaseConnection `mapstructure:"slaves"`
	Pool   PoolConfig           `mapstructure:"pool"`
}

type DatabaseConnection struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Charset  string `mapstructure:"charset"`
}

type PoolConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

type UploadConfig struct {
	Path      string   `mapstructure:"path"`
	MaxSize   int64    `mapstructure:"max_size"`
	AllowExts []string `mapstructure:"allow_exts"`
}

// LoadConfig 加载配置
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}
```

### 3. 数据模型定义

```go
// internal/models/user.go
package models

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	Gender       int8           `gorm:"default:0;comment:0未知1男2女" json:"gender"`
	Birthday     *time.Time     `json:"birthday"`
	Status       int8           `gorm:"default:1;index;comment:1正常2禁用" json:"status"`
	RoleID       uint64         `gorm:"index" json:"role_id"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Role     *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Courses  []Course  `gorm:"many2many:user_courses" json:"courses,omitempty"`
	Orders   []Order   `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Progress []LearningProgress `gorm:"foreignKey:UserID" json:"progress,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Nickname == "" {
		u.Nickname = u.Username
	}
	return nil
}

// AfterFind 查询后钩子
func (u *User) AfterFind(tx *gorm.DB) error {
	// 可以在这里添加一些后处理逻辑
	return nil
}

// Role 角色模型
type Role struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:50;not null;index" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Permissions string    `gorm:"type:json" json:"permissions"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

// UserProfile 用户资料扩展
type UserProfile struct {
	UserID      uint64 `gorm:"primaryKey" json:"user_id"`
	RealName    string `gorm:"size:50" json:"real_name"`
	IDCard      string `gorm:"size:18" json:"id_card"`
	Address     string `gorm:"size:255" json:"address"`
	Company     string `gorm:"size:100" json:"company"`
	Position    string `gorm:"size:50" json:"position"`
	Introduction string `gorm:"type:text" json:"introduction"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
```

```go
// internal/models/course.go
package models

import (
	"time"
	"gorm.io/gorm"
)

// Category 课程分类
type Category struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ParentID    uint64    `gorm:"default:0;index" json:"parent_id"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	Icon        string    `gorm:"size:255" json:"icon"`
	Status      int8      `gorm:"default:1;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Parent   *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Courses  []Course  `gorm:"foreignKey:CategoryID" json:"courses,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

// Course 课程模型
type Course struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Subtitle      string         `gorm:"size:500" json:"subtitle"`
	Description   string         `gorm:"type:text" json:"description"`
	CoverImage    string         `gorm:"size:255" json:"cover_image"`
	CategoryID    uint64         `gorm:"not null;index" json:"category_id"`
	TeacherID     uint64         `gorm:"not null;index" json:"teacher_id"`
	Price         float64        `gorm:"type:decimal(10,2);default:0.00;index" json:"price"`
	OriginalPrice float64        `gorm:"type:decimal(10,2);default:0.00" json:"original_price"`
	Difficulty    int8           `gorm:"default:1;comment:1初级2中级3高级" json:"difficulty"`
	Duration      int            `gorm:"default:0;comment:课程时长分钟" json:"duration"`
	StudentCount  int            `gorm:"default:0" json:"student_count"`
	Rating        float64        `gorm:"type:decimal(3,2);default:0.00" json:"rating"`
	RatingCount   int            `gorm:"default:0" json:"rating_count"`
	Status        int8           `gorm:"default:1;index;comment:1草稿2发布3下架" json:"status"`
	IsFree        bool           `gorm:"default:false" json:"is_free"`
	Tags          string         `gorm:"type:json" json:"tags"`
	CreatedAt     time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Teacher  *User     `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Chapters []Chapter `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"chapters,omitempty"`
	Lessons  []Lesson  `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
	Students []User    `gorm:"many2many:user_courses" json:"students,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}

// BeforeCreate 创建前钩子
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.OriginalPrice == 0 {
		c.OriginalPrice = c.Price
	}
	return nil
}

// AfterUpdate 更新后钩子
func (c *Course) AfterUpdate(tx *gorm.DB) error {
	// 可以在这里触发缓存更新等操作
	return nil
}

// Chapter 章节模型
type Chapter struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID    uint64    `gorm:"not null;index" json:"course_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	SortOrder   int       `gorm:"default:0;index" json:"sort_order"`
	Duration    int       `gorm:"default:0" json:"duration"`
	IsFree      bool      `gorm:"default:false" json:"is_free"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Course  *Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
}

func (Chapter) TableName() string {
	return "chapters"
}

// Lesson 课时模型
type Lesson struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ChapterID     uint64    `gorm:"not null;index" json:"chapter_id"`
	CourseID      uint64    `gorm:"not null;index" json:"course_id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	VideoURL      string    `gorm:"size:255" json:"video_url"`
	VideoDuration int       `gorm:"default:0" json:"video_duration"`
	SortOrder     int       `gorm:"default:0;index" json:"sort_order"`
	IsFree        bool      `gorm:"default:false" json:"is_free"`
	Status        int8      `gorm:"default:1" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联关系
	Chapter *Chapter `gorm:"foreignKey:ChapterID" json:"chapter,omitempty"`
	Course  *Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (Lesson) TableName() string {
	return "lessons"
}
```

### 4. 仓储层实现

```go
// internal/repositories/user_repo.go
package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"your-project/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, params ListParams) ([]*models.User, int64, error)
	UpdateLastLogin(ctx context.Context, id uint64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

type ListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Status   *int8
	RoleID   *uint64
	SortBy   string
	SortDesc bool
}

func (r *userRepository) List(ctx context.Context, params ListParams) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	// 添加搜索条件
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", keyword, keyword, keyword)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.RoleID != nil {
		query = query.Where("role_id = ?", *params.RoleID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		orderBy = params.SortBy
		if params.SortDesc {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}
	}
	query = query.Order(orderBy)

	// 添加分页
	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// 预加载关联数据
	err := query.Preload("Role").Find(&users).Error
	return users, total, err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}
```

### 5. 服务层实现

```go
// internal/services/user_service.go
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"your-project/internal/models"
	"your-project/internal/repositories"
	"your-project/pkg/cache"
	"your-project/pkg/logger"
)

type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID uint64) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint64, req UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID uint64, req ChangePasswordRequest) error
	GetUserList(ctx context.Context, req GetUserListRequest) (*GetUserListResponse, error)
	DeleteUser(ctx context.Context, userID uint64) error
}

type userService struct {
	userRepo repositories.UserRepository
	cache    cache.Cache
	logger   logger.Logger
}

func NewUserService(
	userRepo repositories.UserRepository,
	cache cache.Cache,
	logger logger.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
	}
}

// 请求和响应结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

type UpdateProfileRequest struct {
	Nickname string     `json:"nickname"`
	Phone    string     `json:"phone"`
	Avatar   string     `json:"avatar"`
	Gender   int8       `json:"gender"`
	Birthday *time.Time `json:"birthday"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

type GetUserListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
	RoleID   *uint64 `form:"role_id"`
	SortBy   string `form:"sort_by"`
	SortDesc bool   `form:"sort_desc"`
}

type GetUserListResponse struct {
	Users []models.User `json:"users"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	PageSize int        `json:"page_size"`
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("检查用户名失败", "error", err)
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("检查邮箱失败", "error", err)
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("密码加密失败", "error", err)
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
		Phone:        req.Phone,
		Status:       1, // 正常状态
		RoleID:       2, // 默认普通用户角色
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("创建用户失败", "error", err)
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	s.logger.Info("用户注册成功", "user_id", user.ID, "username", user.Username)
	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		s.logger.Error("获取用户信息失败", "error", err)
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := s.generateJWTToken(user)
	if err != nil {
		s.logger.Error("生成JWT令牌失败", "error", err)
		return nil, fmt.Errorf("生成JWT令牌失败: %w", err)
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Warn("更新最后登录时间失败", "error", err)
	}

	// 清除密码字段
	user.PasswordHash = ""

	s.logger.Info("用户登录成功", "user_id", user.ID, "username", user.Username)
	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID uint64) (*models.User, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:profile:%d", userID)
	var user models.User
	if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 从数据库获取
	userPtr, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		s.logger.Error("获取用户资料失败", "error", err)
		return nil, fmt.Errorf("获取用户资料失败: %w", err)
	}

	// 清除敏感信息
	userPtr.PasswordHash = ""

	// 缓存用户信息
	if err := s.cache.Set(ctx, cacheKey, userPtr, 30*time.Minute); err != nil {
		s.logger.Warn("缓存用户信息失败", "error", err)
	}

	return userPtr, nil
}

// generateJWTToken 生成JWT令牌
func (s *userService) generateJWTToken(user *models.User) (string, error) {
	// 这里应该实现JWT令牌生成逻辑
	// 简化示例，实际应该使用jwt-go库
	return "jwt_token_here", nil
}
```

---

## 🧪 测试实现

### 1. 单元测试

```go
// tests/unit/user_service_test.go
package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"your-project/internal/models"
	"your-project/internal/services"
	"your-project/tests/mocks"
)

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name    string
		req     services.RegisterRequest
		mockFn  func(*mocks.UserRepository, *mocks.Cache, *mocks.Logger)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功注册",
			req: services.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Nickname: "Test User",
			},
			mockFn: func(userRepo *mocks.UserRepository, cache *mocks.Cache, logger *mocks.Logger) {
				userRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, gorm.ErrRecordNotFound)
				userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, gorm.ErrRecordNotFound)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
			wantErr: false,
		},
		{
			name: "用户名已存在",
			req: services.RegisterRequest{
				Username: "existuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockFn: func(userRepo *mocks.UserRepository, cache *mocks.Cache, logger *mocks.Logger) {
				existingUser := &models.User{ID: 1, Username: "existuser"}
				userRepo.On("GetByUsername", mock.Anything, "existuser").Return(existingUser, nil)
			},
			wantErr: true,
			errMsg:  "用户名已存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟对象
			userRepo := new(mocks.UserRepository)
			cache := new(mocks.Cache)
			logger := new(mocks.Logger)

			// 设置模拟行为
			tt.mockFn(userRepo, cache, logger)

			// 创建服务
			userService := services.NewUserService(userRepo, cache, logger)

			// 执行测试
			user, err := userService.Register(context.Background(), tt.req)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.req.Username, user.Username)
				assert.Equal(t, tt.req.Email, user.Email)
			}

			// 验证模拟对象的调用
			userRepo.AssertExpectations(t)
			cache.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name    string
		req     services.LoginRequest
		mockFn  func(*mocks.UserRepository, *mocks.Cache, *mocks.Logger)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功登录",
			req: services.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockFn: func(userRepo *mocks.UserRepository, cache *mocks.Cache, logger *mocks.Logger) {
				// 注意：这里的密码哈希是 "password123" 的bcrypt哈希值
				user := &models.User{
					ID:           1,
					Username:     "testuser",
					PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
					Status:       1,
				}
				userRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil)
				userRepo.On("UpdateLastLogin", mock.Anything, uint64(1)).Return(nil)
				logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
			wantErr: false,
		},
		{
			name: "用户不存在",
			req: services.LoginRequest{
				Username: "nonexist",
				Password: "password123",
			},
			mockFn: func(userRepo *mocks.UserRepository, cache *mocks.Cache, logger *mocks.Logger) {
				userRepo.On("GetByUsername", mock.Anything, "nonexist").Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errMsg:  "用户名或密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟对象
			userRepo := new(mocks.UserRepository)
			cache := new(mocks.Cache)
			logger := new(mocks.Logger)

			// 设置模拟行为
			tt.mockFn(userRepo, cache, logger)

			// 创建服务
			userService := services.NewUserService(userRepo, cache, logger)

			// 执行测试
			resp, err := userService.Login(context.Background(), tt.req)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, tt.req.Username, resp.User.Username)
				assert.Empty(t, resp.User.PasswordHash) // 密码应该被清除
			}

			// 验证模拟对象的调用
			userRepo.AssertExpectations(t)
			cache.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
```

### 2. 集成测试

```go
// tests/integration/user_integration_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"your-project/internal/handlers"
	"your-project/internal/services"
	"your-project/tests/testutils"
)

type UserIntegrationTestSuite struct {
	suite.Suite
	router     *gin.Engine
	testDB     *testutils.TestDB
	userHandler *handlers.UserHandler
}

func (suite *UserIntegrationTestSuite) SetupSuite() {
	// 初始化测试数据库
	suite.testDB = testutils.NewTestDB()

	// 初始化服务和处理器
	userRepo := repositories.NewUserRepository(suite.testDB.DB)
	cache := testutils.NewTestCache()
	logger := testutils.NewTestLogger()
	userService := services.NewUserService(userRepo, cache, logger)
	suite.userHandler = handlers.NewUserHandler(userService)

	// 初始化路由
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.setupRoutes()
}

func (suite *UserIntegrationTestSuite) TearDownSuite() {
	suite.testDB.Close()
}

func (suite *UserIntegrationTestSuite) SetupTest() {
	// 每个测试前清理数据
	suite.testDB.CleanUp()
}

func (suite *UserIntegrationTestSuite) setupRoutes() {
	v1 := suite.router.Group("/api/v1")
	{
		v1.POST("/register", suite.userHandler.Register)
		v1.POST("/login", suite.userHandler.Login)
		v1.GET("/profile", suite.userHandler.GetProfile)
		v1.PUT("/profile", suite.userHandler.UpdateProfile)
	}
}

func (suite *UserIntegrationTestSuite) TestUserRegister() {
	reqBody := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
		"nickname": "Test User",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["status"])
}

func (suite *UserIntegrationTestSuite) TestUserLogin() {
	// 先注册用户
	suite.TestUserRegister()

	reqBody := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", response["status"])
	assert.NotEmpty(suite.T(), response["data"].(map[string]interface{})["token"])
}

func TestUserIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserIntegrationTestSuite))
}
```

### 3. 性能测试

```go
// tests/performance/user_performance_test.go
package performance

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"your-project/internal/services"
	"your-project/tests/testutils"
)

func BenchmarkUserService_Register(b *testing.B) {
	// 初始化测试环境
	testDB := testutils.NewTestDB()
	defer testDB.Close()

	userRepo := repositories.NewUserRepository(testDB.DB)
	cache := testutils.NewTestCache()
	logger := testutils.NewTestLogger()
	userService := services.NewUserService(userRepo, cache, logger)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			req := services.RegisterRequest{
				Username: fmt.Sprintf("user%d", i),
				Email:    fmt.Sprintf("user%d@example.com", i),
				Password: "password123",
				Nickname: fmt.Sprintf("User %d", i),
			}
			_, err := userService.Register(context.Background(), req)
			assert.NoError(b, err)
			i++
		}
	})
}

func BenchmarkUserService_Login(b *testing.B) {
	// 初始化测试环境
	testDB := testutils.NewTestDB()
	defer testDB.Close()

	userRepo := repositories.NewUserRepository(testDB.DB)
	cache := testutils.NewTestCache()
	logger := testutils.NewTestLogger()
	userService := services.NewUserService(userRepo, cache, logger)

	// 预先创建用户
	regReq := services.RegisterRequest{
		Username: "benchuser",
		Email:    "bench@example.com",
		Password: "password123",
	}
	_, err := userService.Register(context.Background(), regReq)
	assert.NoError(b, err)

	loginReq := services.LoginRequest{
		Username: "benchuser",
		Password: "password123",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := userService.Login(context.Background(), loginReq)
			assert.NoError(b, err)
		}
	})
}

// 并发测试
func TestConcurrentUserOperations(t *testing.T) {
	testDB := testutils.NewTestDB()
	defer testDB.Close()

	userRepo := repositories.NewUserRepository(testDB.DB)
	cache := testutils.NewTestCache()
	logger := testutils.NewTestLogger()
	userService := services.NewUserService(userRepo, cache, logger)

	concurrency := 100
	var wg sync.WaitGroup
	errorChan := make(chan error, concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req := services.RegisterRequest{
				Username: fmt.Sprintf("concurrent_user_%d", index),
				Email:    fmt.Sprintf("concurrent_%d@example.com", index),
				Password: "password123",
			}

			_, err := userService.Register(context.Background(), req)
			if err != nil {
				errorChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	duration := time.Since(start)
	t.Logf("并发注册%d个用户耗时: %v", concurrency, duration)

	// 检查是否有错误
	for err := range errorChan {
		t.Errorf("并发操作出错: %v", err)
	}
}
```

---

## 🚀 部署配置

### 1. Docker配置

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./main"]
```

### 2. Docker Compose配置

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=edu_platform
      - REDIS_ADDR=redis:6379
    depends_on:
      - mysql
      - redis
    volumes:
      - ./configs:/root/configs
      - ./logs:/root/logs
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: edu_platform
      MYSQL_USER: app_user
      MYSQL_PASSWORD: app_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    command: --default-authentication-plugin=mysql_native_password
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    restart: unless-stopped

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana:/etc/grafana/provisioning
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
  prometheus_data:
  grafana_data:
```

### 3. Nginx配置

```nginx
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    server {
        listen 80;
        server_name localhost;

        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
```

---

## 📊 监控和日志

### 1. Prometheus监控配置

```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'edu-platform'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'mysql'
    static_configs:
      - targets: ['mysql:3306']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']
```

### 2. 应用监控指标

```go
// pkg/metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP请求总数
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP请求持续时间
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// 数据库连接数
	DatabaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"},
	)

	// 用户注册总数
	UserRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	// 课程购买总数
	CoursePurchases = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "course_purchases_total",
			Help: "Total number of course purchases",
		},
	)
)
```

---

## 🎯 实战练习任务

### 任务1：基础功能实现 (⭐⭐)

**目标**：实现用户注册、登录、课程浏览功能

**要求**：
1. 完成用户模型定义和数据库迁移
2. 实现用户注册和登录API
3. 实现课程列表和详情API
4. 添加基本的参数验证和错误处理
5. 编写单元测试

**验收标准**：
- [ ] 用户可以成功注册和登录
- [ ] 可以获取课程列表和详情
- [ ] API返回格式统一
- [ ] 单元测试覆盖率 > 80%

### 任务2：高级功能开发 (⭐⭐⭐)

**目标**：实现订单系统、学习进度跟踪

**要求**：
1. 实现课程购买和订单管理
2. 添加学习进度跟踪功能
3. 实现课程评价系统
4. 添加Redis缓存优化
5. 实现事务处理

**验收标准**：
- [ ] 用户可以购买课程并生成订单
- [ ] 学习进度可以正确记录和更新
- [ ] 缓存命中率 > 90%
- [ ] 事务回滚正常工作

### 任务3：性能优化 (⭐⭐⭐⭐)

**目标**：优化系统性能，支持高并发

**要求**：
1. 实现数据库读写分离
2. 添加连接池优化
3. 实现查询优化和索引设计
4. 添加限流和熔断机制
5. 性能测试和调优

**验收标准**：
- [ ] 支持1000+并发用户
- [ ] API响应时间 < 100ms
- [ ] 数据库查询优化完成
- [ ] 系统稳定性测试通过

### 任务4：企业级特性 (⭐⭐⭐⭐⭐)

**目标**：实现完整的企业级功能

**要求**：
1. 实现完整的权限控制系统
2. 添加数据统计和报表功能
3. 实现消息队列处理
4. 添加分布式锁
5. 完善监控和告警

**验收标准**：
- [ ] 权限控制精确到接口级别
- [ ] 实时数据统计正常
- [ ] 消息队列处理稳定
- [ ] 监控指标完整
- [ ] 告警机制有效

---

## 📚 学习总结

### 核心知识点回顾

1. **GORM高级特性**
   - 复杂关联关系处理
   - 事务和并发控制
   - 性能优化技巧
   - 插件和钩子使用

2. **企业级开发实践**
   - 项目架构设计
   - 代码组织和模块化
   - 错误处理和日志记录
   - 测试驱动开发

3. **性能优化策略**
   - 数据库优化
   - 缓存策略
   - 并发处理
   - 监控和调优

4. **部署和运维**
   - 容器化部署
   - 监控和告警
   - 日志管理
   - 性能分析

### 进阶学习建议

1. **深入学习微服务架构**
   - 服务拆分策略
   - 服务间通信
   - 分布式事务
   - 服务治理

2. **掌握云原生技术**
   - Kubernetes部署
   - 服务网格
   - 云数据库
   - 自动扩缩容

3. **学习大数据处理**
   - 数据仓库设计
   - 实时数据处理
   - 数据分析
   - 机器学习集成

### 实践建议

1. **持续优化**
   - 定期性能测试
   - 代码重构
   - 技术债务管理
   - 新技术调研

2. **团队协作**
   - 代码审查
   - 文档维护
   - 知识分享
   - 最佳实践总结

3. **生产环境实践**
   - 灰度发布
   - 故障处理
   - 容量规划
   - 安全加固

---

## 🔗 相关资源

### 官方文档
- [GORM官方文档](https://gorm.io/docs/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [Go语言官方文档](https://golang.org/doc/)

### 学习资源
- [Go语言高级编程](https://chai2010.cn/advanced-go-programming-book/)
- [微服务设计模式](https://microservices.io/patterns/)
- [数据库性能优化指南](https://use-the-index-luke.com/)

### 开源项目
- [Gin实战项目](https://github.com/gin-gonic/examples)
- [GORM示例代码](https://github.com/go-gorm/gorm/tree/master/examples)
- [Go微服务框架](https://github.com/go-kit/kit)

---

🎉 **恭喜你完成了GORM企业级项目开发实战练习！**

通过这个完整的项目实战，你已经掌握了GORM在企业级应用中的核心技能。继续实践和优化，你将成为Go语言和GORM的专家！
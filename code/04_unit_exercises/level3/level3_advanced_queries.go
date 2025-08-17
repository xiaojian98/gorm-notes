// 04_unit_exercises/level3_advanced_queries.go - Level 3 高级查询练习
// 对应文档：03_GORM单元练习_基础技能训练.md
// 本文件实现了GORM的高级查询练习，包括条件查询、聚合查询、子查询、连接查询等
// 支持SQLite和MySQL两种数据库类型

package main

import (
	"fmt"  // 格式化输出
	"log"  // 日志记录
	"time" // 时间处理

	"gorm.io/driver/mysql"  // MySQL数据库驱动
	"gorm.io/driver/sqlite" // SQLite数据库驱动
	"gorm.io/gorm"          // GORM核心库
	"gorm.io/gorm/logger"   // GORM日志组件
	"gorm.io/gorm/schema"   // GORM模式配置
)

// 数据库配置相关定义

// DatabaseType 数据库类型枚举
// 定义支持的数据库类型，目前支持SQLite和MySQL
type DatabaseType string

const (
	SQLite DatabaseType = "sqlite" // SQLite数据库类型
	MySQL  DatabaseType = "mysql"  // MySQL数据库类型
)

// DatabaseConfig 数据库配置结构体
// 包含数据库连接和连接池的所有配置参数
type DatabaseConfig struct {
	Type         DatabaseType    // 数据库类型(sqlite/mysql)
	DSN          string          // 数据源名称,用于指定数据库连接字符串
	MaxOpenConns int             // 最大打开连接数
	MaxIdleConns int             // 最大空闲连接数
	MaxLifetime  time.Duration   // 连接最大生命周期
	LogLevel     logger.LogLevel // 日志级别
}

// GetDefaultConfig 获取SQLite默认配置
// 返回一个包含默认参数的SQLite数据库配置对象
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:         SQLite,
		DSN:          "level3_advanced_queries.db", // SQLite数据库文件名
		MaxOpenConns: 10,                           // 最大连接数10
		MaxIdleConns: 5,                            // 最大空闲连接5
		MaxLifetime:  time.Hour,                    // 连接生命周期1小时
		LogLevel:     logger.Info,                  // 日志级别为Info
	}
}

// GetMySQLConfig 获取MySQL配置
// 参数dsn: MySQL数据库连接字符串
// 返回一个包含默认参数的MySQL数据库配置对象
func GetMySQLConfig(dsn string) *DatabaseConfig {
	return &DatabaseConfig{
		Type:         MySQL,
		DSN:          dsn,
		MaxOpenConns: 20,        // MySQL建议更高的连接数
		MaxIdleConns: 10,        // 更多的空闲连接
		MaxLifetime:  time.Hour, // 连接生命周期1小时
		LogLevel:     logger.Info,
	}
}

// 基础模型定义

// BaseModel 基础模型结构体
// 包含所有数据库表通用的字段，采用GORM的软删除机制
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`              // 主键ID，自动递增
	CreatedAt time.Time      `json:"created_at"`                        // 创建时间，GORM自动管理
	UpdatedAt time.Time      `json:"updated_at"`                        // 更新时间，GORM自动管理
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 删除时间，用于软删除，建立索引
}

// 数据模型定义

// User 用户模型
// 用户基础信息表，包含用户的基本信息和统计数据
type User struct {
	BaseModel           // 继承基础模型字段
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"` // 用户名，唯一索引，最大长度50，非空
	Email     string    `gorm:"uniqueIndex;size:100;not null" json:"email"`   // 邮箱，唯一索引，最大长度100，非空
	Age       int       `gorm:"check:age >= 0 AND age <= 150" json:"age"`     // 年龄，检查约束：0-150岁
	City      string    `gorm:"size:100;index" json:"city"`                   // 城市，建立索引，最大长度100
	Salary    float64   `gorm:"precision:10;scale:2" json:"salary"`           // 薪资，精度10位，小数点后2位
	JoinDate  time.Time `gorm:"index" json:"join_date"`                       // 加入日期，建立索引
	IsActive  bool      `gorm:"default:true;index" json:"is_active"`          // 是否激活，默认true，建立索引

	// 关联关系定义
	// Has Many关系，一个用户可以发表多篇文章
	Posts []Post `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	// Has Many关系，一个用户可以发表多条评论
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
}

// Category 分类模型
// 文章分类表，用于组织和管理文章
type Category struct {
	BaseModel          // 继承基础模型字段
	Name        string `gorm:"size:100;not null;index" json:"name"`       // 分类名称，建立索引，最大长度100，非空
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"` // URL友好的分类标识，唯一索引，最大长度100，非空
	Description string `gorm:"type:text" json:"description"`              // 分类描述，文本类型，可存储长文本
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`       // 分类是否激活，默认true，建立索引

	// 关联关系定义
	// Has Many关系，一个分类可以包含多篇文章
	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

// Post 文章模型
// 文章表，包含文章的详细信息和统计数据
type Post struct {
	BaseModel              // 继承基础模型字段
	Title       string     `gorm:"size:200;not null;index" json:"title"`        // 文章标题，建立索引，最大长度200，非空
	Content     string     `gorm:"type:text;not null" json:"content"`           // 文章内容，文本类型，非空，可存储长文本
	Status      string     `gorm:"size:20;default:'draft';index" json:"status"` // 文章状态，建立索引，最大长度20，默认为'draft'(草稿)
	ViewCount   int        `gorm:"default:0;index" json:"view_count"`           // 浏览次数，默认0，建立索引用于排序查询
	LikeCount   int        `gorm:"default:0;index" json:"like_count"`           // 点赞次数，默认0，建立索引用于排序查询
	PublishedAt *time.Time `gorm:"index" json:"published_at"`                   // 发布时间，建立索引，指针类型允许为空
	Rating      float64    `gorm:"precision:3;scale:2;default:0" json:"rating"` // 文章评分，精度3位，小数点后2位，默认0

	// 外键关系定义
	AuthorID   uint  `gorm:"not null;index" json:"author_id"` // 作者ID外键，建立索引，非空
	CategoryID *uint `gorm:"index" json:"category_id"`        // 分类ID外键，建立索引，指针类型允许为空

	// 关联关系定义
	// Belongs To关系，通过外键关联到其他表
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`     // 关联的作者对象，通过AuthorID外键
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 关联的分类对象，通过CategoryID外键，指针类型允许为空
	// Has Many关系，一篇文章可以有多个评论
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	// Many To Many关系，一篇文章可以有多个标签，一个标签可以属于多篇文章
	Tags []Tag `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

// Comment 评论模型
// 评论表，用于存储用户对文章的评论
type Comment struct {
	BaseModel        // 继承基础模型字段
	Content   string `gorm:"type:text;not null" json:"content"`             // 评论内容，文本类型，非空，可存储长文本
	Status    string `gorm:"size:20;default:'pending';index" json:"status"` // 评论状态，建立索引，最大长度20，默认为'pending'(待审核)
	LikeCount int    `gorm:"default:0" json:"like_count"`                   // 评论点赞次数，默认0

	// 外键关系定义
	PostID   uint `gorm:"not null;index" json:"post_id"`   // 文章ID外键，建立索引，非空
	AuthorID uint `gorm:"not null;index" json:"author_id"` // 评论作者ID外键，建立索引，非空

	// 关联关系定义
	// Belongs To关系，通过外键关联到其他表
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`     // 关联的文章对象，通过PostID外键
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"` // 关联的评论作者对象，通过AuthorID外键
}

// Tag 标签模型
// 标签表，用于文章的分类标记，支持多对多关系
type Tag struct {
	BaseModel        // 继承基础模型字段
	Name      string `gorm:"uniqueIndex;size:50;not null" json:"name"` // 标签名称，唯一索引，最大长度50，非空
	Slug      string `gorm:"uniqueIndex;size:50;not null" json:"slug"` // URL友好的标签标识，唯一索引，最大长度50，非空
	IsActive  bool   `gorm:"default:true" json:"is_active"`            // 标签是否激活，默认true

	// 关联关系定义
	// Many To Many关系，通过post_tags中间表关联
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// 统计结构定义
// 这些结构体用于存储聚合查询的结果

// UserStats 用户统计信息结构体
// 用于存储用户的各项统计数据
type UserStats struct {
	UserID       uint    `json:"user_id"`       // 用户ID
	Username     string  `json:"username"`      // 用户名
	PostCount    int64   `json:"post_count"`    // 发表文章数量
	CommentCount int64   `json:"comment_count"` // 发表评论数量
	TotalViews   int64   `json:"total_views"`   // 文章总浏览量
	TotalLikes   int64   `json:"total_likes"`   // 文章总点赞数
	AvgRating    float64 `json:"avg_rating"`    // 文章平均评分
}

// CategoryStats 分类统计信息结构体
// 用于存储分类的各项统计数据
type CategoryStats struct {
	CategoryID   uint    `json:"category_id"`   // 分类ID
	CategoryName string  `json:"category_name"` // 分类名称
	PostCount    int64   `json:"post_count"`    // 分类下文章数量
	TotalViews   int64   `json:"total_views"`   // 分类下文章总浏览量
	AvgRating    float64 `json:"avg_rating"`    // 分类下文章平均评分
}

// MonthlyStats 月度统计信息结构体
// 用于存储按月统计的数据
type MonthlyStats struct {
	Year      int   `json:"year"`       // 年份
	Month     int   `json:"month"`      // 月份
	PostCount int64 `json:"post_count"` // 当月发表文章数量
	UserCount int64 `json:"user_count"` // 当月注册用户数量
}

// 数据库初始化相关函数

// initDB 初始化SQLite数据库连接（兼容性保留）
// 使用默认的SQLite配置创建数据库连接
// 返回: *gorm.DB 数据库连接对象
func initDB() *gorm.DB {
	// 使用SQLite驱动打开数据库连接
	// 配置日志级别为Info，显示详细的SQL执行信息
	// 配置命名策略：使用复数表名(SingularTable: false)
	db, err := gorm.Open(sqlite.Open("level3_advanced_queries.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为Info
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数表名，如users、posts等
		},
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移所有模型到数据库
	// 按照依赖关系顺序迁移：User -> Category -> Post -> Comment -> Tag
	err = db.AutoMigrate(&User{}, &Category{}, &Post{}, &Comment{}, &Tag{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db
}

// InitDatabase 根据配置初始化数据库连接
// 支持SQLite和MySQL两种数据库类型
// 参数config: 数据库配置对象
// 返回: *gorm.DB 数据库连接对象
func InitDatabase(config *DatabaseConfig) *gorm.DB {
	var dialector gorm.Dialector

	// 根据数据库类型选择相应的驱动
	switch config.Type {
	case SQLite:
		// 使用SQLite驱动
		dialector = sqlite.Open(config.DSN)
	case MySQL:
		// 使用MySQL驱动
		dialector = mysql.Open(config.DSN)
	default:
		log.Fatalf("不支持的数据库类型: %s", config.Type)
	}

	// 创建数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel), // 使用配置中的日志级别
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数表名
		},
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 获取底层的sql.DB对象，用于配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)   // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)   // 设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(config.MaxLifetime) // 设置连接最大生命周期

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	// 自动迁移所有模型
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Printf("数据库连接成功 - 类型: %s\n", config.Type)
	return db
}

// AutoMigrate 执行数据库迁移
// 按照正确的依赖顺序迁移所有模型
// 参数db: 数据库连接对象
// 返回: error 迁移过程中的错误
func AutoMigrate(db *gorm.DB) error {
	// 按照依赖关系顺序迁移模型
	// 1. 基础模型：User, Category, Tag
	// 2. 依赖模型：Post(依赖User和Category)
	// 3. 关联模型：Comment(依赖Post和User)
	return db.AutoMigrate(
		&User{},     // 用户表
		&Category{}, // 分类表
		&Post{},     // 文章表(依赖User和Category)
		&Comment{},  // 评论表(依赖Post和User)
		&Tag{},      // 标签表
	)
}

// 创建测试数据函数

// createTestData 创建测试数据
// 为数据库创建一套完整的测试数据，包括用户、分类、文章、评论和标签
// 参数db: 数据库连接对象
func createTestData(db *gorm.DB) {
	// 创建用户测试数据
	// 包含不同年龄、城市、薪资和状态的用户，用于测试各种查询条件
	users := []User{
		{Username: "alice", Email: "alice@example.com", Age: 25, City: "北京", Salary: 8000, JoinDate: time.Now().AddDate(-1, 0, 0), IsActive: true},
		{Username: "bob", Email: "bob@example.com", Age: 30, City: "上海", Salary: 12000, JoinDate: time.Now().AddDate(-2, 0, 0), IsActive: true},
		{Username: "charlie", Email: "charlie@example.com", Age: 28, City: "深圳", Salary: 10000, JoinDate: time.Now().AddDate(-1, -6, 0), IsActive: false},
		{Username: "diana", Email: "diana@example.com", Age: 35, City: "北京", Salary: 15000, JoinDate: time.Now().AddDate(-3, 0, 0), IsActive: true},
		{Username: "eve", Email: "eve@example.com", Age: 22, City: "广州", Salary: 6000, JoinDate: time.Now().AddDate(0, -3, 0), IsActive: true},
	}
	// 批量创建用户记录
	db.Create(&users)

	// 创建分类测试数据
	// 包含不同状态的分类，用于测试分类相关查询
	categories := []Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", IsActive: true},
		{Name: "生活", Slug: "life", Description: "生活分享", IsActive: true},
		{Name: "旅游", Slug: "travel", Description: "旅游攻略", IsActive: true},
		{Name: "美食", Slug: "food", Description: "美食推荐", IsActive: false}, // 非激活状态，用于测试过滤
	}
	// 批量创建分类记录
	db.Create(&categories)

	// 创建标签测试数据
	// 用于测试多对多关系和标签相关查询
	tags := []Tag{
		{Name: "Go", Slug: "go", IsActive: true},
		{Name: "数据库", Slug: "database", IsActive: true},
		{Name: "前端", Slug: "frontend", IsActive: true},
		{Name: "后端", Slug: "backend", IsActive: true},
		{Name: "教程", Slug: "tutorial", IsActive: true},
	}
	// 批量创建标签记录
	db.Create(&tags)

	// 创建文章测试数据
	// 包含不同状态、浏览量、点赞数和评分的文章，用于测试各种排序和过滤条件
	posts := []Post{
		{Title: "Go语言入门", Content: "Go语言基础教程...", Status: "published", ViewCount: 1500, LikeCount: 120, Rating: 4.5, AuthorID: 1, CategoryID: &[]uint{1}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, -1, 0)}[0]},
		{Title: "数据库设计原则", Content: "数据库设计的基本原则...", Status: "published", ViewCount: 800, LikeCount: 65, Rating: 4.2, AuthorID: 2, CategoryID: &[]uint{1}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, -2, 0)}[0]},
		{Title: "我的北京生活", Content: "在北京的生活感悟...", Status: "published", ViewCount: 300, LikeCount: 25, Rating: 3.8, AuthorID: 1, CategoryID: &[]uint{2}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, 0, -15)}[0]},
		{Title: "上海旅游攻略", Content: "上海必去景点推荐...", Status: "published", ViewCount: 1200, LikeCount: 95, Rating: 4.7, AuthorID: 3, CategoryID: &[]uint{3}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, 0, -30)}[0]},
		{Title: "React入门教程", Content: "React基础知识...", Status: "draft", ViewCount: 0, LikeCount: 0, Rating: 0, AuthorID: 4, CategoryID: &[]uint{1}[0]}, // 草稿状态，用于测试状态过滤
		{Title: "深圳美食推荐", Content: "深圳好吃的餐厅...", Status: "published", ViewCount: 600, LikeCount: 40, Rating: 4.0, AuthorID: 5, CategoryID: &[]uint{4}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, 0, -7)}[0]},
	}
	// 批量创建文章记录
	db.Create(&posts)

	// 创建评论测试数据
	// 包含不同状态和点赞数的评论，用于测试评论相关查询
	comments := []Comment{
		{Content: "很好的教程！", Status: "approved", LikeCount: 10, PostID: 1, AuthorID: 2},
		{Content: "学到了很多", Status: "approved", LikeCount: 5, PostID: 1, AuthorID: 3},
		{Content: "写得不错", Status: "approved", LikeCount: 3, PostID: 2, AuthorID: 1},
		{Content: "有用的信息", Status: "pending", LikeCount: 0, PostID: 3, AuthorID: 4}, // 待审核状态，用于测试状态过滤
		{Content: "期待更多内容", Status: "approved", LikeCount: 8, PostID: 4, AuthorID: 5},
	}
	// 批量创建评论记录
	db.Create(&comments)

	// 为文章添加标签（多对多关系）
	// 演示如何建立多对多关系，一篇文章可以有多个标签
	var post1, post2 Post
	db.First(&post1, 1) // 获取第一篇文章
	db.First(&post2, 2) // 获取第二篇文章
	var tag1, tag2, tag5 Tag
	db.First(&tag1, 1) // 获取"Go"标签
	db.First(&tag2, 2) // 获取"数据库"标签
	db.First(&tag5, 5) // 获取"教程"标签
	// 为第一篇文章添加"Go"和"教程"标签
	db.Model(&post1).Association("Tags").Append([]Tag{tag1, tag5})
	// 为第二篇文章添加"数据库"和"教程"标签
	db.Model(&post2).Association("Tags").Append([]Tag{tag2, tag5})
}

// 练习1：条件查询和排序函数

// FindUsersByConditions 多条件查询用户
// 演示如何使用多个WHERE条件和ORDER BY进行复杂查询
// 参数:
//
//	minAge, maxAge: 年龄范围
//	cities: 城市列表
//	isActive: 是否激活
//
// 返回: ([]User, error) 符合条件的用户列表和可能的错误
func FindUsersByConditions(db *gorm.DB, minAge, maxAge int, cities []string, isActive bool) ([]User, error) {
	var users []User
	// 构建查询条件：年龄范围查询，使用BETWEEN操作符
	query := db.Where("age BETWEEN ? AND ?", minAge, maxAge)

	// 如果提供了城市列表，添加城市过滤条件，使用IN操作符
	if len(cities) > 0 {
		query = query.Where("city IN ?", cities)
	}

	// 添加激活状态过滤条件
	query = query.Where("is_active = ?", isActive)

	// 执行查询并排序：按薪资降序，年龄升序
	// 这演示了多字段排序的用法
	result := query.Order("salary DESC, age ASC").Find(&users)
	return users, result.Error
}

// FindPostsByDateRange 按日期范围查询文章
// 演示日期范围查询和条件预加载
// 参数:
//
//	startDate, endDate: 日期范围
//	status: 文章状态（可选）
//
// 返回: ([]Post, error) 符合条件的文章列表和可能的错误
func FindPostsByDateRange(db *gorm.DB, startDate, endDate time.Time, status string) ([]Post, error) {
	var posts []Post
	// 构建日期范围查询条件
	query := db.Where("published_at BETWEEN ? AND ?", startDate, endDate)

	// 如果提供了状态参数，添加状态过滤条件
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 执行查询，按浏览量降序排序，并预加载关联的作者和分类信息
	// Preload用于避免N+1查询问题
	result := query.Order("view_count DESC").Preload("Author").Preload("Category").Find(&posts)
	return posts, result.Error
}

// FindTopRatedPosts 查询高评分文章
// 演示多条件查询、排序和限制结果数量
// 参数:
//
//	minRating: 最低评分
//	limit: 返回结果数量限制
//
// 返回: ([]Post, error) 高评分文章列表和可能的错误
func FindTopRatedPosts(db *gorm.DB, minRating float64, limit int) ([]Post, error) {
	var posts []Post
	// 构建查询条件：评分大于等于指定值且状态为已发布
	// 使用多字段排序：先按评分降序，再按点赞数降序
	// 使用Limit限制返回结果数量
	result := db.Where("rating >= ? AND status = ?", minRating, "published").
		Order("rating DESC, like_count DESC").
		Limit(limit).
		Preload("Author").   // 预加载作者信息
		Preload("Category"). // 预加载分类信息
		Find(&posts)
	return posts, result.Error
}

// 练习2：聚合查询函数

// GetUserStatistics 获取用户统计信息
// 演示复杂的聚合查询，包括COUNT、SUM、AVG等聚合函数
// 使用LEFT JOIN连接多个表，GROUP BY分组统计
// 返回: ([]UserStats, error) 用户统计信息列表和可能的错误
func GetUserStatistics(db *gorm.DB) ([]UserStats, error) {
	var stats []UserStats
	// 使用Table方法指定主表，并给表起别名
	// Select中使用聚合函数进行统计计算
	err := db.Table("users u").
		Select(`
			u.id as user_id,
			u.username,
			COUNT(DISTINCT p.id) as post_count,        -- 统计用户发表的文章数量（去重）
			COUNT(DISTINCT c.id) as comment_count,     -- 统计用户发表的评论数量（去重）
			COALESCE(SUM(p.view_count), 0) as total_views,  -- 统计文章总浏览量，使用COALESCE处理NULL值
			COALESCE(SUM(p.like_count), 0) as total_likes,  -- 统计文章总点赞数
			COALESCE(AVG(p.rating), 0) as avg_rating        -- 计算文章平均评分
		`).
		// 使用LEFT JOIN连接posts表，保留没有发表文章的用户
		Joins("LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL").
		// 使用LEFT JOIN连接comments表，保留没有发表评论的用户
		Joins("LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL").
		// 过滤已删除的用户
		Where("u.deleted_at IS NULL").
		// 按用户分组进行聚合计算
		Group("u.id, u.username").
		// 按总浏览量降序排序
		Order("total_views DESC").
		// 将查询结果扫描到UserStats结构体切片中
		Scan(&stats).Error
	return stats, err
}

// GetCategoryStatistics 获取分类统计信息
// 演示分类维度的聚合查询
// 返回: ([]CategoryStats, error) 分类统计信息列表和可能的错误
func GetCategoryStatistics(db *gorm.DB) ([]CategoryStats, error) {
	var stats []CategoryStats
	// 以categories表为主表进行统计
	err := db.Table("categories c").
		Select(`
			c.id as category_id,
			c.name as category_name,
			COUNT(p.id) as post_count,                     -- 统计分类下的文章数量
			COALESCE(SUM(p.view_count), 0) as total_views, -- 统计分类下文章的总浏览量
			COALESCE(AVG(p.rating), 0) as avg_rating       -- 计算分类下文章的平均评分
		`).
		// LEFT JOIN posts表，只统计已发布的文章
		Joins("LEFT JOIN posts p ON c.id = p.category_id AND p.deleted_at IS NULL AND p.status = 'published'").
		// 只统计激活状态的分类
		Where("c.deleted_at IS NULL AND c.is_active = ?", true).
		// 按分类分组
		Group("c.id, c.name").
		// 按文章数量降序排序
		Order("post_count DESC").
		Scan(&stats).Error
	return stats, err
}

// GetMonthlyStatistics 获取月度统计
// 演示时间维度的聚合查询，使用UNION ALL合并多个查询结果
// 参数months: 统计最近几个月的数据
// 返回: ([]MonthlyStats, error) 月度统计信息列表和可能的错误
func GetMonthlyStatistics(db *gorm.DB, months int) ([]MonthlyStats, error) {
	var stats []MonthlyStats
	// 计算起始日期
	startDate := time.Now().AddDate(0, -months, 0)

	// 使用Raw SQL进行复杂的时间统计查询
	// 这个查询演示了：
	// 1. 使用strftime函数提取年月信息
	// 2. 使用UNION ALL合并posts和users的统计
	// 3. 使用CASE WHEN进行条件计数
	err := db.Raw(`
		SELECT 
			strftime('%Y', date) as year,   -- 提取年份
			strftime('%m', date) as month,  -- 提取月份
			COUNT(CASE WHEN type = 'post' THEN 1 END) as post_count,  -- 条件计数：文章数量
			COUNT(CASE WHEN type = 'user' THEN 1 END) as user_count   -- 条件计数：用户数量
		FROM (
			-- 合并posts和users表的创建时间数据
			SELECT created_at as date, 'post' as type FROM posts WHERE created_at >= ? AND deleted_at IS NULL
			UNION ALL
			SELECT created_at as date, 'user' as type FROM users WHERE created_at >= ? AND deleted_at IS NULL
		) combined
		GROUP BY year, month  -- 按年月分组
		ORDER BY year DESC, month DESC  -- 按时间倒序排列
	`, startDate, startDate).Scan(&stats).Error

	return stats, err
}

// 练习3：子查询函数

// FindUsersWithMostPosts 查询发文最多的用户
// 演示子查询的使用，先统计每个用户的文章数量，再连接用户表获取详细信息
// 参数limit: 返回用户数量限制
// 返回: ([]User, error) 发文最多的用户列表和可能的错误
func FindUsersWithMostPosts(db *gorm.DB, limit int) ([]User, error) {
	var users []User
	// 创建子查询：统计每个用户的文章数量
	subQuery := db.Table("posts").Select("author_id, COUNT(*) as post_count").
		Where("deleted_at IS NULL AND status = ?", "published"). // 只统计已发布的文章
		Group("author_id").                                      // 按作者分组
		Order("post_count DESC").                                // 按文章数量降序
		Limit(limit)                                             // 限制结果数量

	// 主查询：使用子查询结果连接用户表
	err := db.Table("users u").
		Joins("JOIN (?) pc ON u.id = pc.author_id", subQuery). // 连接子查询结果
		Where("u.deleted_at IS NULL").                         // 过滤已删除用户
		Order("pc.post_count DESC").                           // 按文章数量排序
		Find(&users).Error

	return users, err
}

// FindPostsAboveAverageViews 查询浏览量高于平均值的文章
// 演示如何在查询中使用平均值作为条件
// 返回: ([]Post, error) 浏览量高于平均值的文章列表和可能的错误
func FindPostsAboveAverageViews(db *gorm.DB) ([]Post, error) {
	var posts []Post

	// 先计算所有已发布文章的平均浏览量
	var avgViews float64
	db.Table("posts").Where("deleted_at IS NULL AND status = ?", "published").Select("AVG(view_count)").Scan(&avgViews)

	// 查询浏览量高于平均值的文章
	err := db.Where("view_count > ? AND status = ?", avgViews, "published").
		Order("view_count DESC"). // 按浏览量降序排序
		Preload("Author").        // 预加载作者信息
		Preload("Category").      // 预加载分类信息
		Find(&posts).Error

	return posts, err
}

// FindUsersWithNoComments 查询没有评论的用户
// 演示NOT IN子查询的使用
// 返回: ([]User, error) 没有发表评论的用户列表和可能的错误
func FindUsersWithNoComments(db *gorm.DB) ([]User, error) {
	var users []User
	// 使用NOT IN子查询，查找不在评论作者列表中的用户
	err := db.Where("id NOT IN (?)",
		// 子查询：获取所有发表过评论的用户ID
		db.Table("comments").Select("DISTINCT author_id").Where("deleted_at IS NULL"),
	).Find(&users).Error
	return users, err
}

// 练习4：复杂连接查询函数

// GetPostsWithAuthorAndCommentCount 获取文章及作者信息和评论数
// 演示多表连接查询，包括INNER JOIN和LEFT JOIN
// 参数limit: 返回结果数量限制
// 返回: ([]map[string]interface{}, error) 文章详细信息列表和可能的错误
func GetPostsWithAuthorAndCommentCount(db *gorm.DB, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 复杂的多表连接查询
	err := db.Table("posts p").
		Select(`
			p.id,
			p.title,
			p.view_count,
			p.like_count,
			p.rating,
			p.published_at,
			u.username as author_name,     -- 作者用户名
			u.city as author_city,         -- 作者城市
			c.name as category_name,       -- 分类名称
			COUNT(cm.id) as comment_count  -- 评论数量统计
		`).
		// INNER JOIN users表获取作者信息（必须有作者）
		Joins("JOIN users u ON p.author_id = u.id").
		// LEFT JOIN categories表获取分类信息（分类可能为空）
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		// LEFT JOIN comments表统计评论数量（文章可能没有评论）
		Joins("LEFT JOIN comments cm ON p.id = cm.post_id AND cm.deleted_at IS NULL").
		// 只查询已发布的文章
		Where("p.deleted_at IS NULL AND p.status = ?", "published").
		// 按文章分组，因为使用了聚合函数COUNT
		Group("p.id, p.title, p.view_count, p.like_count, p.rating, p.published_at, u.username, u.city, c.name").
		// 按浏览量降序排序
		Order("p.view_count DESC").
		// 限制返回结果数量
		Limit(limit).
		// 扫描结果到map切片中，适用于动态字段的查询结果
		Scan(&results).Error

	return results, err
}

// GetUserEngagementReport 获取用户参与度报告
// 演示复杂的用户活跃度统计查询
// 返回: ([]map[string]interface{}, error) 用户参与度报告和可能的错误
func GetUserEngagementReport(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 复杂的用户参与度统计查询
	err := db.Table("users u").
		Select(`
			u.id,
			u.username,
			u.city,
			u.join_date,
			COUNT(DISTINCT p.id) as post_count,                    -- 发表文章数
			COUNT(DISTINCT c.id) as comment_count,                 -- 发表评论数
			COALESCE(SUM(p.view_count), 0) as total_post_views,    -- 文章总浏览量
			COALESCE(SUM(p.like_count), 0) as total_post_likes,    -- 文章总点赞数
			COALESCE(SUM(c.like_count), 0) as total_comment_likes, -- 评论总点赞数
			COALESCE(AVG(p.rating), 0) as avg_post_rating,         -- 文章平均评分
			(COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) as total_activity  -- 总活跃度（文章+评论）
		`).
		// LEFT JOIN posts表统计用户发表的文章
		Joins("LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL").
		// LEFT JOIN comments表统计用户发表的评论
		Joins("LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL").
		// 只统计未删除的用户
		Where("u.deleted_at IS NULL").
		// 按用户分组
		Group("u.id, u.username, u.city, u.join_date").
		// 使用HAVING过滤掉没有任何活动的用户
		Having("total_activity > 0").
		// 按总活跃度降序排序
		Order("total_activity DESC").
		Scan(&results).Error

	return results, err
}

// 练习5：窗口函数和排名函数

// supportsWindowFunctions 检查数据库是否支持窗口函数
// MySQL 8.0+ 和 SQLite 3.25+ 支持窗口函数
// 返回: bool 是否支持窗口函数
func supportsWindowFunctions(db *gorm.DB) bool {
	// 获取数据库方言名称
	dialector := db.Dialector.Name()

	switch dialector {
	case "mysql":
		// 检查MySQL版本
		var version string
		err := db.Raw("SELECT VERSION()").Scan(&version).Error
		if err != nil {
			return false
		}
		// MySQL 8.0+ 支持窗口函数
		// 简单检查版本号是否以8开头
		return len(version) > 0 && version[0] >= '8'
	case "sqlite":
		// SQLite 3.25+ 支持窗口函数
		// 大多数现代SQLite版本都支持，这里返回true
		return true
	default:
		// 其他数据库默认不支持
		return false
	}
}

// GetTopPostsByCategory 获取每个分类的热门文章
// 演示窗口函数ROW_NUMBER()的使用，实现分组排名
// 参数topN: 每个分类返回的文章数量
// 返回: ([]map[string]interface{}, error) 每个分类的热门文章列表和可能的错误
func GetTopPostsByCategory(db *gorm.DB, topN int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 检查数据库类型和版本，决定使用窗口函数还是兼容性查询
	if supportsWindowFunctions(db) {
		// 使用窗口函数进行分组排名的复杂查询（MySQL 8.0+ 或 SQLite 3.25+）
		err := db.Raw(`
			SELECT 
				category_name,
				title,
				view_count,
				like_count,
				rating,
				rank_in_category
			FROM (
				-- 内层查询：使用窗口函数为每个分类的文章排名
				SELECT 
					c.name as category_name,
					p.title,
					p.view_count,
					p.like_count,
					p.rating,
					-- ROW_NUMBER()窗口函数：按分类分组，按浏览量排序，生成排名
					ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY p.view_count DESC) as rank_in_category
				FROM posts p
				JOIN categories c ON p.category_id = c.id
				WHERE p.deleted_at IS NULL AND p.status = 'published'
				  AND c.deleted_at IS NULL AND c.is_active = 1
			) ranked
			-- 外层查询：过滤出每个分类的前N名文章
			WHERE rank_in_category <= ?
			ORDER BY category_name, rank_in_category
		`, topN).Scan(&results).Error
		return results, err
	} else {
		// 兼容性查询：使用变量模拟排名功能（适用于MySQL 5.7及以下版本）
		err := db.Raw(`
			SELECT 
				category_name,
				title,
				view_count,
				like_count,
				rating,
				@rank := CASE 
					WHEN @prev_category = cat_id THEN @rank + 1
					ELSE 1
				END as rank_in_category,
				@prev_category := cat_id
			FROM (
				SELECT 
					p.title,
					p.view_count,
					p.like_count,
					p.rating,
					c.name as category_name,
					c.id as cat_id
				FROM posts p
				JOIN categories c ON p.category_id = c.id
				WHERE p.deleted_at IS NULL AND p.status = 'published'
				  AND c.deleted_at IS NULL AND c.is_active = 1
				ORDER BY c.id, p.view_count DESC
			) ranked_posts
			CROSS JOIN (SELECT @rank := 0, @prev_category := '') r
			HAVING rank_in_category <= ?
			ORDER BY category_name, rank_in_category
		`, topN).Scan(&results).Error
		return results, err
	}
}

// GetUserRankingByActivity 获取用户活跃度排名
// 演示多个窗口函数的使用，包括全局排名和分组排名
// 返回: ([]map[string]interface{}, error) 用户活跃度排名列表和可能的错误
func GetUserRankingByActivity(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 检查数据库类型和版本，决定使用窗口函数还是兼容性查询
	if supportsWindowFunctions(db) {
		// 复杂的用户排名查询，使用多个窗口函数（MySQL 8.0+ 或 SQLite 3.25+）
		err := db.Raw(`
			SELECT 
				username,
				city,
				total_activity,
				activity_rank,
				city_rank,
				-- 使用CASE WHEN根据排名分配徽章
				CASE 
					WHEN activity_rank <= 3 THEN 'Gold'    -- 前3名获得金牌
					WHEN activity_rank <= 10 THEN 'Silver' -- 4-10名获得银牌
					ELSE 'Bronze'                          -- 其他获得铜牌
				END as badge
			FROM (
				-- 内层查询：计算用户活跃度并生成排名
				SELECT 
					u.username,
					u.city,
					(COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) as total_activity,
					-- 全局活跃度排名
					ROW_NUMBER() OVER (ORDER BY (COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) DESC) as activity_rank,
					-- 城市内活跃度排名
					ROW_NUMBER() OVER (PARTITION BY u.city ORDER BY (COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) DESC) as city_rank
				FROM users u
				LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL
				LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL
				WHERE u.deleted_at IS NULL
				GROUP BY u.id, u.username, u.city
				HAVING (COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) > 0  -- 修复：直接使用表达式而不是别名
			) ranked
			ORDER BY activity_rank
		`).Scan(&results).Error
		return results, err
	} else {
		// 兼容性查询：使用变量模拟排名功能（适用于MySQL 5.7及以下版本）
		err := db.Raw(`
			SELECT 
				username,
				city,
				total_activity,
				@activity_rank := @activity_rank + 1 as activity_rank,
				@city_rank := CASE 
					WHEN @prev_city = city THEN @city_rank + 1
					ELSE 1
				END as city_rank,
				@prev_city := city,
				-- 使用CASE WHEN根据排名分配徽章
				CASE 
					WHEN @activity_rank <= 3 THEN 'Gold'    -- 前3名获得金牌
					WHEN @activity_rank <= 10 THEN 'Silver' -- 4-10名获得银牌
					ELSE 'Bronze'                          -- 其他获得铜牌
				END as badge
			FROM (
				SELECT 
					u.username,
					u.city,
					(COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) as total_activity
				FROM users u
				LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL
				LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL
				WHERE u.deleted_at IS NULL
				GROUP BY u.id, u.username, u.city
				HAVING (COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) > 0
				ORDER BY total_activity DESC, city
			) user_activity
			CROSS JOIN (SELECT @activity_rank := 0, @city_rank := 0, @prev_city := '') r
			ORDER BY total_activity DESC
		`).Scan(&results).Error
		return results, err
	}
}

// 练习6：性能优化查询函数

// GetPostsWithOptimizedPreloading 优化预加载的文章查询
// 演示如何使用选择性预加载和字段选择来优化查询性能
// 参数limit: 返回结果数量限制
// 返回: ([]Post, error) 优化后的文章列表和可能的错误
func GetPostsWithOptimizedPreloading(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post

	// 使用选择性预加载，只加载需要的字段，减少数据传输量
	err := db.Select("id, title, content, view_count, like_count, rating, published_at, author_id, category_id").
		// 预加载作者信息，但只选择必要的字段
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, city")
		}).
		// 预加载分类信息，但只选择必要的字段
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug")
		}).
		// 预加载标签信息，只加载激活的标签，并只选择必要字段
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug").Where("is_active = ?", true)
		}).
		// 只查询已发布的文章
		Where("status = ?", "published").
		// 按发布时间降序排序
		Order("published_at DESC").
		// 限制返回结果数量
		Limit(limit).
		Find(&posts).Error

	return posts, err
}

// GetPostsWithPagination 分页查询文章
// 演示分页查询的实现，包括总数统计和条件过滤
// 参数:
//
//	page: 页码（从1开始）
//	pageSize: 每页大小
//	categoryID: 分类ID过滤（可选）
//	search: 搜索关键词（可选）
//
// 返回: ([]Post, int64, error) 文章列表、总数和可能的错误
func GetPostsWithPagination(db *gorm.DB, page, pageSize int, categoryID *uint, search string) ([]Post, int64, error) {
	var posts []Post
	var total int64

	// 构建查询条件：只查询已发布的文章
	query := db.Model(&Post{}).Where("status = ?", "published")

	// 如果提供了分类ID，添加分类过滤条件
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// 如果提供了搜索关键词，在标题和内容中搜索
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 计算符合条件的记录总数
	query.Count(&total)

	// 计算偏移量：(页码-1) * 每页大小
	offset := (page - 1) * pageSize
	// 执行分页查询
	err := query.Offset(offset).Limit(pageSize).
		// 预加载作者信息，只选择必要字段
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		// 预加载分类信息，只选择必要字段
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		// 按发布时间降序排序
		Order("published_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// 主函数演示
// 演示所有高级查询功能的使用方法
func main() {
	fmt.Println("=== GORM Level 3 高级查询练习 ===")

	// 可以选择使用SQLite或MySQL数据库
	// 默认使用SQLite，如需使用MySQL，请取消注释下面的代码并配置正确的DSN

	// // SQLite配置（默认）
	// db := initDB()
	// fmt.Println("✓ SQLite数据库初始化完成")

	// MySQL配置（可选）
	// mysqlDSN := "root:123456@tcp(192.168.100.124:3306)/gorm_note?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDSN := "root:fastbee@tcp(192.168.100.124:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	config := GetMySQLConfig(mysqlDSN)
	db := InitDatabase(config)
	fmt.Println("✓ MySQL数据库初始化完成")

	// 创建测试数据
	createTestData(db)
	fmt.Println("✓ 测试数据创建完成")

	// 练习1：条件查询和排序
	fmt.Println("\n=== 练习1：条件查询和排序 ===")

	// 1.1 多条件查询用户：查询20-35岁，在北京或上海，且激活状态的用户
	users, err := FindUsersByConditions(db, 20, 35, []string{"北京", "上海"}, true)
	if err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个符合条件的用户:\n", len(users))
		for _, user := range users {
			fmt.Printf("  - %s (%d岁, %s, 薪资: %.0f)\n", user.Username, user.Age, user.City, user.Salary)
		}
	}

	// 1.2 按日期范围查询文章：查询最近3个月发布的文章
	startDate := time.Now().AddDate(0, -3, 0)
	endDate := time.Now()
	posts, err := FindPostsByDateRange(db, startDate, endDate, "published")
	if err != nil {
		fmt.Printf("查询文章失败: %v\n", err)
	} else {
		fmt.Printf("\n最近3个月发布的文章 (%d篇):\n", len(posts))
		for _, post := range posts {
			fmt.Printf("  - %s (浏览: %d, 点赞: %d, 作者: %s)\n",
				post.Title, post.ViewCount, post.LikeCount, post.Author.Username)
		}
	}

	// 1.3 查询高评分文章：评分>=4.0的前3篇文章
	topPosts, err := FindTopRatedPosts(db, 4.0, 3)
	if err != nil {
		fmt.Printf("查询高评分文章失败: %v\n", err)
	} else {
		fmt.Printf("\n高评分文章 (评分>=4.0, 前3篇):\n")
		for _, post := range topPosts {
			fmt.Printf("  - %s (评分: %.1f, 浏览: %d)\n",
				post.Title, post.Rating, post.ViewCount)
		}
	}

	// 练习2：聚合查询
	fmt.Println("\n=== 练习2：聚合查询 ===")

	// 2.1 获取用户统计信息
	userStats, err := GetUserStatistics(db)
	if err != nil {
		fmt.Printf("获取用户统计失败: %v\n", err)
	} else {
		fmt.Println("用户统计信息:")
		for _, stat := range userStats {
			fmt.Printf("  - %s: 文章%d篇, 评论%d条, 总浏览%d, 总点赞%d, 平均评分%.1f\n",
				stat.Username, stat.PostCount, stat.CommentCount,
				stat.TotalViews, stat.TotalLikes, stat.AvgRating)
		}
	}

	// 2.2 获取分类统计信息
	categoryStats, err := GetCategoryStatistics(db)
	if err != nil {
		fmt.Printf("获取分类统计失败: %v\n", err)
	} else {
		fmt.Println("\n分类统计信息:")
		for _, stat := range categoryStats {
			fmt.Printf("  - %s: 文章%d篇, 总浏览%d, 平均评分%.1f\n",
				stat.CategoryName, stat.PostCount, stat.TotalViews, stat.AvgRating)
		}
	}

	// 2.3 获取月度统计信息
	monthlyStats, err := GetMonthlyStatistics(db, 6)
	if err != nil {
		fmt.Printf("获取月度统计失败: %v\n", err)
	} else {
		fmt.Println("\n最近6个月统计:")
		for _, stat := range monthlyStats {
			fmt.Printf("  - %d年%02d月: 文章%d篇, 新用户%d人\n",
				stat.Year, stat.Month, stat.PostCount, stat.UserCount)
		}
	}

	// 练习3：子查询
	fmt.Println("\n=== 练习3：子查询 ===")

	// 3.1 查询发文最多的用户
	topUsers, err := FindUsersWithMostPosts(db, 3)
	if err != nil {
		fmt.Printf("查询发文最多用户失败: %v\n", err)
	} else {
		fmt.Printf("发文最多的用户 (前3名):\n")
		for _, user := range topUsers {
			fmt.Printf("  - %s (%s)\n", user.Username, user.City)
		}
	}

	// 3.2 查询浏览量高于平均值的文章
	highViewPosts, err := FindPostsAboveAverageViews(db)
	if err != nil {
		fmt.Printf("查询高浏览量文章失败: %v\n", err)
	} else {
		fmt.Printf("\n浏览量高于平均值的文章 (%d篇):\n", len(highViewPosts))
		for _, post := range highViewPosts {
			fmt.Printf("  - %s (浏览: %d)\n", post.Title, post.ViewCount)
		}
	}

	// 3.3 查询没有发表评论的用户
	noCommentUsers, err := FindUsersWithNoComments(db)
	if err != nil {
		fmt.Printf("查询无评论用户失败: %v\n", err)
	} else {
		fmt.Printf("\n没有发表评论的用户 (%d人):\n", len(noCommentUsers))
		for _, user := range noCommentUsers {
			fmt.Printf("  - %s\n", user.Username)
		}
	}

	// 练习4：复杂连接查询
	fmt.Println("\n=== 练习4：复杂连接查询 ===")

	// 4.1 获取文章详情（包含作者和评论数）
	postDetails, err := GetPostsWithAuthorAndCommentCount(db, 5)
	if err != nil {
		fmt.Printf("获取文章详情失败: %v\n", err)
	} else {
		fmt.Println("文章详情 (包含作者和评论数):")
		for _, detail := range postDetails {
			fmt.Printf("  - %s (作者: %s, 分类: %s, 评论: %v条)\n",
				detail["title"], detail["author_name"],
				detail["category_name"], detail["comment_count"])
		}
	}

	// 4.2 获取用户参与度报告
	engagementReport, err := GetUserEngagementReport(db)
	if err != nil {
		fmt.Printf("获取用户参与度报告失败: %v\n", err)
	} else {
		fmt.Println("\n用户参与度报告:")
		for _, report := range engagementReport {
			fmt.Printf("  - %s: 总活动%v次 (文章%v篇, 评论%v条)\n",
				report["username"], report["total_activity"],
				report["post_count"], report["comment_count"])
		}
	}

	// 练习5：窗口函数和排名
	fmt.Println("\n=== 练习5：窗口函数和排名 ===")

	// 5.1 获取各分类热门文章
	topByCategory, err := GetTopPostsByCategory(db, 2)
	if err != nil {
		fmt.Printf("获取分类热门文章失败: %v\n", err)
	} else {
		fmt.Println("各分类热门文章 (前2名):")
		for _, item := range topByCategory {
			fmt.Printf("  - [%s] %s (浏览: %v, 排名: %v)\n",
				item["category_name"], item["title"],
				item["view_count"], item["rank_in_category"])
		}
	}

	// 5.2 获取用户活跃度排名
	userRanking, err := GetUserRankingByActivity(db)
	if err != nil {
		fmt.Printf("获取用户活跃度排名失败: %v\n", err)
	} else {
		fmt.Println("\n用户活跃度排名:")
		for _, rank := range userRanking {
			fmt.Printf("  - %s (%s): 活动%v次, 全站排名%v, 城市排名%v, 徽章: %s\n",
				rank["username"], rank["city"], rank["total_activity"],
				rank["activity_rank"], rank["city_rank"], rank["badge"])
		}
	}

	// 练习6：性能优化查询
	fmt.Println("\n=== 练习6：性能优化查询 ===")

	// 6.1 优化预加载查询
	optimizedPosts, err := GetPostsWithOptimizedPreloading(db, 3)
	if err != nil {
		fmt.Printf("优化预加载查询失败: %v\n", err)
	} else {
		fmt.Printf("优化预加载查询结果 (%d篇文章):\n", len(optimizedPosts))
		for _, post := range optimizedPosts {
			fmt.Printf("  - %s (作者: %s, 分类: %s, 标签数: %d)\n",
				post.Title, post.Author.Username,
				func() string {
					if post.Category != nil {
						return post.Category.Name
					}
					return "无分类"
				}(), len(post.Tags))
		}
	}

	// 6.2 分页查询
	paginatedPosts, total, err := GetPostsWithPagination(db, 1, 3, nil, "")
	if err != nil {
		fmt.Printf("分页查询失败: %v\n", err)
	} else {
		fmt.Printf("\n分页查询结果 (第1页, 每页3篇, 总共%d篇):\n", total)
		for _, post := range paginatedPosts {
			fmt.Printf("  - %s (作者: %s)\n", post.Title, post.Author.Username)
		}
	}

	// 演示MySQL和SQLite的差异处理
	fmt.Println("\n=== 数据库兼容性说明 ===")
	fmt.Println("本代码同时支持SQLite和MySQL数据库:")
	fmt.Println("- SQLite: 适合开发和测试环境，无需额外配置")
	fmt.Println("- MySQL: 适合生产环境，需要配置正确的DSN连接字符串")
	fmt.Println("- 窗口函数在SQLite 3.25+和MySQL 8.0+中支持")
	fmt.Println("- 如需切换数据库，请修改main函数中的数据库初始化代码")

	fmt.Println("\n=== Level 3 高级查询练习完成 ===")
	fmt.Println("\n练习总结:")
	fmt.Println("1. ✓ 条件查询和排序 - 掌握复杂WHERE条件和ORDER BY")
	fmt.Println("2. ✓ 聚合查询 - 掌握COUNT、SUM、AVG等聚合函数")
	fmt.Println("3. ✓ 子查询 - 掌握嵌套查询和EXISTS/NOT EXISTS")
	fmt.Println("4. ✓ 复杂连接查询 - 掌握多表JOIN和复杂关联")
	fmt.Println("5. ✓ 窗口函数和排名 - 掌握ROW_NUMBER()等窗口函数")
	fmt.Println("6. ✓ 性能优化查询 - 掌握预加载优化和分页查询")
	fmt.Println("\n恭喜！您已完成GORM高级查询的所有练习！")
}

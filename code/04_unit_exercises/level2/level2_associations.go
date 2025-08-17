// 04_unit_exercises/level2_associations.go - Level 2 关联关系练习
// 对应文档：03_GORM单元练习_基础技能训练.md
// 本文件实现了GORM的关联关系练习，包括一对一、一对多、多对多关系的定义和操作
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
		DSN:          "level2_associations.db", // SQLite数据库文件名
		MaxOpenConns: 10,                       // 最大连接数10
		MaxIdleConns: 5,                        // 最大空闲连接5
		MaxLifetime:  time.Hour,                // 连接生命周期1小时
		LogLevel:     logger.Info,              // 日志级别为Info
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

// 练习1：一对一关系模型定义
// 演示用户(User)与用户资料(Profile)之间的一对一关系
// 一个用户对应一个用户资料，一个用户资料属于一个用户

// User 用户模型
// 用户基础信息表，与Profile表建立一对一关系
type User struct {
	BaseModel        // 继承基础模型字段
	Username  string `gorm:"uniqueIndex;size:50;not null" json:"username"` // 用户名，唯一索引，最大长度50，非空
	Email     string `gorm:"uniqueIndex;size:100;not null" json:"email"`   // 邮箱，唯一索引，最大长度100，非空
	Password  string `gorm:"size:255;not null" json:"-"`                   // 密码，最大长度255，非空，JSON序列化时忽略

	// 一对一关系：用户资料
	// 使用Has One关系，当用户更新时级联更新资料，当用户删除时级联删除资料
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile,omitempty"`
}

// Profile 用户资料模型
// 用户详细资料表，与User表建立一对一关系
type Profile struct {
	BaseModel            // 继承基础模型字段
	UserID    uint       `gorm:"uniqueIndex;not null" json:"user_id"` // 用户ID外键，唯一索引，非空，确保一对一关系
	FirstName string     `gorm:"size:50" json:"first_name"`           // 名字，最大长度50
	LastName  string     `gorm:"size:50" json:"last_name"`            // 姓氏，最大长度50
	Bio       string     `gorm:"type:text" json:"bio"`                // 个人简介，文本类型，可存储长文本
	Avatar    string     `gorm:"size:255" json:"avatar"`              // 头像URL，最大长度255
	Phone     string     `gorm:"size:20" json:"phone"`                // 电话号码，最大长度20
	Address   string     `gorm:"size:255" json:"address"`             // 地址，最大长度255
	BirthDate *time.Time `json:"birth_date"`                          // 生日，指针类型允许为空

	// 反向关联 - 使用指针类型避免循环引用
	// Belongs To关系，通过UserID外键关联到User表
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// 练习2：一对多关系模型定义
// 演示分类(Category)与文章(Post)之间的一对多关系
// 一个分类可以包含多篇文章，一篇文章属于一个分类
// 同时演示用户(User)与文章(Post)之间的一对多关系
// 一个用户可以发表多篇文章，一篇文章属于一个用户

// Category 分类模型
// 文章分类表，与Post表建立一对多关系
type Category struct {
	BaseModel          // 继承基础模型字段
	Name        string `gorm:"size:100;not null;index" json:"name"`       // 分类名称，建立索引，最大长度100，非空
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"` // URL友好的分类标识，唯一索引，最大长度100，非空
	Description string `gorm:"type:text" json:"description"`              // 分类描述，文本类型，可存储长文本
	IsActive    bool   `gorm:"default:true" json:"is_active"`             // 分类是否激活，默认为true

	// 一对多关系：分类下的文章
	// Has Many关系，一个分类可以包含多篇文章
	// 当分类更新时级联更新文章，当分类删除时将文章的分类ID设为NULL
	Posts []Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"posts,omitempty"`
}

// Post 文章模型
// 文章表，与Category和User表建立多对一关系，与Comment表建立一对多关系，与Tag表建立多对多关系
type Post struct {
	BaseModel              // 继承基础模型字段
	Title       string     `gorm:"size:200;not null;index" json:"title"`        // 文章标题，建立索引，最大长度200，非空
	Slug        string     `gorm:"uniqueIndex;size:200;not null" json:"slug"`   // URL友好的文章标识，唯一索引，最大长度200，非空
	Content     string     `gorm:"type:text;not null" json:"content"`           // 文章内容，文本类型，非空，可存储长文本
	Excerpt     string     `gorm:"size:500" json:"excerpt"`                     // 文章摘要，最大长度500
	Status      string     `gorm:"size:20;default:'draft';index" json:"status"` // 文章状态，建立索引，最大长度20，默认为'draft'(草稿)
	ViewCount   int        `gorm:"default:0" json:"view_count"`                 // 浏览次数，默认为0
	LikeCount   int        `gorm:"default:0" json:"like_count"`                 // 点赞次数，默认为0
	PublishedAt *time.Time `gorm:"index" json:"published_at"`                   // 发布时间，建立索引，指针类型允许为空

	// 外键关系定义
	AuthorID   uint  `gorm:"not null;index" json:"author_id"` // 作者ID外键，建立索引，非空
	CategoryID *uint `gorm:"index" json:"category_id"`        // 分类ID外键，建立索引，指针类型允许为空

	// 关联关系定义
	// Belongs To关系，通过外键关联到其他表
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`     // 关联的作者对象，通过AuthorID外键
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 关联的分类对象，通过CategoryID外键，指针类型允许为空

	// 一对多关系：文章的评论
	// Has Many关系，一篇文章可以有多个评论
	// 当文章更新时级联更新评论，当文章删除时级联删除评论
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments,omitempty"`

	// 多对多关系：文章的标签
	// Many To Many关系，一篇文章可以有多个标签，一个标签可以属于多篇文章
	// 使用post_tags作为中间表，支持级联更新和删除
	Tags []Tag `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tags,omitempty"`
}

// Comment 评论模型
// 评论表，演示多种关联关系：
// 1. 与Post表的多对一关系(一篇文章可以有多个评论)
// 2. 与User表的多对一关系(一个用户可以发表多个评论)
// 3. 自关联关系(评论可以回复其他评论，形成树状结构)
type Comment struct {
	BaseModel        // 继承基础模型字段
	Content   string `gorm:"type:text;not null" json:"content"`             // 评论内容，文本类型，非空，可存储长文本
	Status    string `gorm:"size:20;default:'pending';index" json:"status"` // 评论状态，建立索引，最大长度20，默认为'pending'(待审核)
	IPAddress string `gorm:"size:45" json:"ip_address"`                     // 评论者IP地址，最大长度45(支持IPv6)

	// 外键关系定义
	PostID   uint  `gorm:"not null;index" json:"post_id"`   // 文章ID外键，建立索引，非空
	AuthorID uint  `gorm:"not null;index" json:"author_id"` // 评论作者ID外键，建立索引，非空
	ParentID *uint `gorm:"index" json:"parent_id"`          // 父评论ID外键，建立索引，指针类型允许为空(顶级评论)

	// 关联关系定义
	// Belongs To关系，通过外键关联到其他表
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`     // 关联的文章对象，通过PostID外键
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"` // 关联的评论作者对象，通过AuthorID外键

	// 自关联：评论回复关系
	// 实现评论的树状结构，支持多级回复
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`  // 父评论对象，通过ParentID外键，指针类型允许为空
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"` // 子评论列表，Has Many关系
}

// 练习3：多对多关系模型定义
// 演示文章(Post)与标签(Tag)之间的多对多关系
// 一篇文章可以有多个标签，一个标签可以属于多篇文章
// 使用中间表post_tags来维护这种多对多关系

// Tag 标签模型
// 标签表，与Post表建立多对多关系
type Tag struct {
	BaseModel          // 继承基础模型字段
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"` // 标签名称，唯一索引，最大长度50，非空
	Slug        string `gorm:"uniqueIndex;size:50;not null" json:"slug"` // URL友好的标签标识，唯一索引，最大长度50，非空
	Description string `gorm:"type:text" json:"description"`             // 标签描述，文本类型，可存储长文本
	Color       string `gorm:"size:7;default:'#007bff'" json:"color"`    // 标签颜色，最大长度7(十六进制颜色码)，默认为蓝色
	IsActive    bool   `gorm:"default:true" json:"is_active"`            // 标签是否激活，默认为true

	// 多对多关系：标签的文章
	// Many To Many关系，通过post_tags中间表关联
	// 一个标签可以属于多篇文章
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// PostTag 文章标签中间表（自定义）
// 用于维护Post和Tag之间的多对多关系
// 可以在中间表中添加额外的字段，如创建时间等
type PostTag struct {
	PostID    uint      `gorm:"primaryKey" json:"post_id"` // 文章ID，复合主键的一部分
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`  // 标签ID，复合主键的一部分
	CreatedAt time.Time `json:"created_at"`                // 关联创建时间，记录文章和标签的关联时间

	// 关联关系定义
	// Belongs To关系，通过外键关联到Post和Tag表
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"` // 关联的文章对象，通过PostID外键
	Tag  Tag  `gorm:"foreignKey:TagID" json:"tag,omitempty"`   // 关联的标签对象，通过TagID外键
}

// 数据库初始化相关函数

// initDB 初始化SQLite数据库连接
// 使用默认的SQLite配置创建数据库连接
// 返回: *gorm.DB 数据库连接对象
func initDB() *gorm.DB {
	// 使用SQLite驱动打开数据库连接
	// 配置日志级别为Info，显示详细的SQL执行信息
	// 配置命名策略：使用复数表名(SingularTable: false)
	db, err := gorm.Open(sqlite.Open("level2_associations.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为Info
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数表名，如users、posts等
		},
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移所有模型到数据库
	// 按照依赖关系顺序迁移：User -> Profile -> Category -> Post -> Comment -> Tag -> PostTag
	err = db.AutoMigrate(&User{}, &Profile{}, &Category{}, &Post{}, &Comment{}, &Tag{}, &PostTag{})
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
	// 2. 依赖模型：Profile(依赖User), Post(依赖User和Category)
	// 3. 关联模型：Comment(依赖Post和User), PostTag(依赖Post和Tag)
	return db.AutoMigrate(
		&User{},     // 用户表
		&Profile{},  // 用户资料表(依赖User)
		&Category{}, // 分类表
		&Post{},     // 文章表(依赖User和Category)
		&Comment{},  // 评论表(依赖Post和User)
		&Tag{},      // 标签表
		&PostTag{},  // 文章标签关联表(依赖Post和Tag)
	)
}

// 练习1：一对一关系操作函数
// 演示User和Profile之间的一对一关系操作

// CreateUserWithProfile 创建用户及其资料
// 演示一对一关系的创建操作，同时创建用户和用户资料
// 参数:
//
//	db: 数据库连接对象
//	username: 用户名
//	email: 邮箱地址
//	password: 密码
//	firstName: 名字
//	lastName: 姓氏
//	bio: 个人简介
//
// 返回: (*User, error) 创建的用户对象和可能的错误
func CreateUserWithProfile(db *gorm.DB, username, email, password, firstName, lastName, bio string) (*User, error) {
	// 创建用户对象，同时包含关联的Profile对象
	// GORM会自动处理一对一关系，先创建User，再创建Profile并设置外键
	user := &User{
		Username: username, // 设置用户名
		Email:    email,    // 设置邮箱
		Password: password, // 设置密码
		Profile: Profile{ // 创建关联的用户资料
			FirstName: firstName, // 设置名字
			LastName:  lastName,  // 设置姓氏
			Bio:       bio,       // 设置个人简介
			// UserID会由GORM自动设置为user.ID
		},
	}

	// 执行创建操作，GORM会自动处理关联关系
	// 1. 首先创建User记录
	// 2. 然后创建Profile记录，并设置UserID外键
	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// 输出创建成功信息，显示生成的ID
	fmt.Printf("创建用户及资料成功，用户ID: %d, 资料ID: %d\n", user.ID, user.Profile.ID)
	return user, nil
}

// GetUserWithProfile 获取用户及其资料
// 演示一对一关系的查询操作，使用Preload预加载关联数据
// 参数:
//
//	db: 数据库连接对象
//	userID: 用户ID
//
// 返回: (*User, error) 用户对象(包含Profile)和可能的错误
func GetUserWithProfile(db *gorm.DB, userID uint) (*User, error) {
	var user User
	// 使用Preload预加载Profile关联数据
	// 这会执行两个SQL查询：
	// 1. SELECT * FROM users WHERE id = ?
	// 2. SELECT * FROM profiles WHERE user_id = ?
	result := db.Preload("Profile").First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUserProfile 更新用户资料
// 演示一对一关系中关联对象的更新操作
// 参数:
//
//	db: 数据库连接对象
//	userID: 用户ID
//	updates: 要更新的字段映射
//
// 返回: error 可能的错误
func UpdateUserProfile(db *gorm.DB, userID uint, updates map[string]interface{}) error {
	// 通过外键user_id定位Profile记录并更新
	// 使用Updates方法可以批量更新多个字段
	result := db.Model(&Profile{}).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	// 输出更新结果信息
	fmt.Printf("更新用户资料成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// 练习2：一对多关系操作函数
// 演示Category和Post之间的一对多关系操作

// CreateCategoryWithPosts 创建分类及其文章
// 演示一对多关系的创建操作，一个分类包含多篇文章
// 参数:
//
//	db: 数据库连接对象
//	categoryName: 分类名称
//	categorySlug: 分类别名(URL友好)
//	authorID: 文章作者ID
//	postTitles: 文章标题列表
//
// 返回: (*Category, error) 创建的分类对象和可能的错误
func CreateCategoryWithPosts(db *gorm.DB, categoryName, categorySlug string, authorID uint, postTitles []string) (*Category, error) {
	// 创建分类对象
	category := &Category{
		Name: categoryName, // 设置分类名称
		Slug: categorySlug, // 设置分类别名
	}

	// 批量创建文章并关联到分类
	// 遍历文章标题列表，为每个标题创建一篇文章
	for i, title := range postTitles {
		// 创建文章对象
		post := Post{
			Title:    title,                                   // 设置文章标题
			Slug:     fmt.Sprintf("%s-%d", categorySlug, i+1), // 生成文章别名
			Content:  fmt.Sprintf("这是%s的内容", title),           // 生成文章内容
			Excerpt:  fmt.Sprintf("这是%s的摘要", title),           // 生成文章摘要
			Status:   "published",                             // 设置文章状态为已发布
			AuthorID: authorID,                                // 设置文章作者ID
			// CategoryID会由GORM自动设置为category.ID
		}
		// 将文章添加到分类的Posts切片中
		category.Posts = append(category.Posts, post)
	}

	// 执行创建操作，GORM会自动处理一对多关系
	// 1. 首先创建Category记录
	// 2. 然后创建所有Post记录，并设置CategoryID外键
	result := db.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}

	// 输出创建成功信息
	fmt.Printf("创建分类及文章成功，分类ID: %d, 文章数量: %d\n", category.ID, len(category.Posts))
	return category, nil
}

// GetCategoryWithPosts 获取分类及其文章
// 演示一对多关系的查询操作，使用嵌套Preload预加载关联数据
// 参数:
//
//	db: 数据库连接对象
//	categoryID: 分类ID
//
// 返回: (*Category, error) 分类对象(包含Posts和Author)和可能的错误
func GetCategoryWithPosts(db *gorm.DB, categoryID uint) (*Category, error) {
	var category Category
	// 使用多层Preload预加载关联数据
	// "Posts": 预加载分类下的所有文章
	// "Posts.Author": 预加载每篇文章的作者信息
	// 这会执行多个SQL查询来获取完整的关联数据
	result := db.Preload("Posts").Preload("Posts.Author").First(&category, categoryID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

// GetPostsWithAuthorAndCategory 获取文章及其作者和分类
// 演示多对一关系的查询操作，从Post角度查询关联的Author和Category
// 参数:
//
//	db: 数据库连接对象
//	limit: 限制返回的文章数量
//
// 返回: ([]Post, error) 文章列表(包含Author和Category)和可能的错误
func GetPostsWithAuthorAndCategory(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post
	// 使用Preload预加载多个关联关系
	// "Author": 预加载文章作者信息
	// "Category": 预加载文章分类信息
	// Limit: 限制查询结果数量
	result := db.Preload("Author").Preload("Category").Limit(limit).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// CreatePostWithComments 创建文章及其评论
// 演示Post和Comment之间的一对多关系操作
// 参数:
//
//	db: 数据库连接对象
//	title: 文章标题
//	content: 文章内容
//	authorID: 文章作者ID
//	categoryID: 文章分类ID
//	commentContents: 评论内容列表
//
// 返回: (*Post, error) 创建的文章对象和可能的错误
func CreatePostWithComments(db *gorm.DB, title, content string, authorID, categoryID uint, commentContents []string) (*Post, error) {
	// 创建文章对象
	post := &Post{
		Title:      title,                                     // 设置文章标题
		Slug:       fmt.Sprintf("post-%d", time.Now().Unix()), // 生成基于时间戳的文章别名
		Content:    content,                                   // 设置文章内容
		Excerpt:    content[:min(len(content), 100)],          // 生成文章摘要(前100字符)
		Status:     "published",                               // 设置文章状态为已发布
		AuthorID:   authorID,                                  // 设置文章作者ID
		CategoryID: &categoryID,                               // 设置文章分类ID(指针类型)
	}

	// 批量创建评论并关联到文章
	// 遍历评论内容列表，为每个内容创建一条评论
	for _, commentContent := range commentContents {
		// 创建评论对象
		comment := Comment{
			Content:  commentContent, // 设置评论内容
			Status:   "approved",     // 设置评论状态为已审核
			AuthorID: authorID,       // 设置评论作者ID
			// PostID会由GORM自动设置为post.ID
		}
		// 将评论添加到文章的Comments切片中
		post.Comments = append(post.Comments, comment)
	}

	// 执行创建操作，GORM会自动处理一对多关系
	// 1. 首先创建Post记录
	// 2. 然后创建所有Comment记录，并设置PostID外键
	result := db.Create(post)
	if result.Error != nil {
		return nil, result.Error
	}

	// 输出创建成功信息
	fmt.Printf("创建文章及评论成功，文章ID: %d, 评论数量: %d\n", post.ID, len(post.Comments))
	return post, nil
}

// GetPostWithComments 获取文章及其评论
// 演示Post和Comment之间一对多关系的查询操作
// 参数:
//
//	db: 数据库连接对象
//	postID: 文章ID
//
// 返回: (*Post, error) 文章对象(包含Comments和Author)和可能的错误
func GetPostWithComments(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	// 使用多层Preload预加载关联数据
	// "Comments": 预加载文章下的所有评论
	// "Comments.Author": 预加载每条评论的作者信息
	result := db.Preload("Comments").Preload("Comments.Author").First(&post, postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// 练习3：多对多关系操作函数
// 演示Post和Tag之间的多对多关系操作，通过PostTag中间表实现

// CreateTagsAndAssignToPosts 创建标签并分配给文章
// 演示多对多关系的创建和关联操作
// 参数:
//
//	db: 数据库连接对象
//	tagNames: 标签名称列表
//	postIDs: 文章ID列表
//
// 返回: error 可能的错误
func CreateTagsAndAssignToPosts(db *gorm.DB, tagNames []string, postIDs []uint) error {
	// 批量创建标签
	// 首先创建所有标签对象
	var tags []Tag
	for _, name := range tagNames {
		// 创建标签对象
		tag := Tag{
			Name: name, // 设置标签名称
			Slug: name, // 设置标签别名(这里简化为与名称相同)
		}
		// 将标签添加到标签切片中
		tags = append(tags, tag)
	}

	// 批量创建所有标签
	result := db.Create(&tags)
	if result.Error != nil {
		return result.Error
	}

	// 为每篇文章分配所有标签
	// 遍历文章ID列表，为每篇文章关联所有标签
	for _, postID := range postIDs {
		// 查找文章对象
		var post Post
		if err := db.First(&post, postID).Error; err != nil {
			// 如果文章不存在，跳过该文章
			continue
		}

		// 使用Association API关联标签
		// 这会在PostTag中间表中创建关联记录
		if err := db.Model(&post).Association("Tags").Append(&tags); err != nil {
			return err
		}
	}

	// 输出创建成功信息
	fmt.Printf("创建标签并分配成功，标签数量: %d\n", len(tags))
	return nil
}

// GetPostWithTags 获取文章及其标签
// 演示多对多关系的查询操作，从Post角度查询关联的Tags
// 参数:
//
//	db: 数据库连接对象
//	postID: 文章ID
//
// 返回: (*Post, error) 文章对象(包含Tags)和可能的错误
func GetPostWithTags(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	// 使用Preload预加载文章的所有标签
	// 这会通过PostTag中间表查询关联的标签
	result := db.Preload("Tags").First(&post, postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// GetTagWithPosts 获取标签及其文章
// 演示多对多关系的查询操作，从Tag角度查询关联的Posts
// 参数:
//
//	db: 数据库连接对象
//	tagID: 标签ID
//
// 返回: (*Tag, error) 标签对象(包含Posts和Author)和可能的错误
func GetTagWithPosts(db *gorm.DB, tagID uint) (*Tag, error) {
	var tag Tag
	// 使用多层Preload预加载关联数据
	// "Posts": 预加载标签下的所有文章
	// "Posts.Author": 预加载每篇文章的作者信息
	result := db.Preload("Posts").Preload("Posts.Author").First(&tag, tagID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

// AddTagsToPost 为文章添加标签
// 演示多对多关系的动态关联操作
// 参数:
//
//	db: 数据库连接对象
//	postID: 文章ID
//	tagIDs: 要添加的标签ID列表
//
// 返回: error 可能的错误
func AddTagsToPost(db *gorm.DB, postID uint, tagIDs []uint) error {
	// 查找指定的文章
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}

	// 查找要添加的标签
	var tags []Tag
	if err := db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	// 使用Association API添加标签关联
	// 这会在PostTag中间表中创建新的关联记录
	return db.Model(&post).Association("Tags").Append(&tags)
}

// RemoveTagsFromPost 从文章中移除标签
// 演示多对多关系的动态解除关联操作
// 参数:
//
//	db: 数据库连接对象
//	postID: 文章ID
//	tagIDs: 要移除的标签ID列表
//
// 返回: error 可能的错误
func RemoveTagsFromPost(db *gorm.DB, postID uint, tagIDs []uint) error {
	// 查找指定的文章
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}

	// 查找要移除的标签
	var tags []Tag
	if err := db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	// 使用Association API删除标签关联
	// 这会从PostTag中间表中删除对应的关联记录
	return db.Model(&post).Association("Tags").Delete(&tags)
}

// 练习4：复杂查询函数
// 演示复杂的多层关联查询和数据统计操作

// GetPostsWithAllAssociations 获取文章及所有关联数据
// 演示复杂的多层关联查询，一次性加载所有相关数据
// 参数:
//
//	db: 数据库连接对象
//	limit: 限制返回的文章数量
//
// 返回: ([]Post, error) 文章列表(包含所有关联数据)和可能的错误
func GetPostsWithAllAssociations(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post
	// 使用链式Preload预加载所有关联数据
	// 这是一个复杂的查询，会执行多个SQL语句来获取完整的数据结构
	result := db.Preload("Author"). // 预加载文章作者
					Preload("Author.Profile").  // 预加载作者的个人资料
					Preload("Category").        // 预加载文章分类
					Preload("Comments").        // 预加载文章评论
					Preload("Comments.Author"). // 预加载评论作者
					Preload("Tags").            // 预加载文章标签
					Limit(limit).Find(&posts)   // 限制结果数量并执行查询

	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// GetUserPostsWithStats 获取用户文章及统计信息
// 演示复杂的数据查询和统计计算
// 参数:
//
//	db: 数据库连接对象
//	userID: 用户ID
//
// 返回: (map[string]interface{}, error) 包含文章和统计信息的映射和可能的错误
func GetUserPostsWithStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
	// 查询用户的所有文章及其关联数据
	var posts []Post
	result := db.Where("author_id = ?", userID).Preload("Category").Preload("Tags").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	// 初始化统计信息结构
	// 创建一个包含各种统计数据的映射
	stats := map[string]interface{}{
		"total_posts": len(posts),           // 文章总数
		"total_views": 0,                    // 总浏览量
		"total_likes": 0,                    // 总点赞数
		"categories":  make(map[string]int), // 分类统计
		"tags":        make(map[string]int), // 标签统计
		"posts":       posts,                // 文章列表
	}

	// 初始化计数器
	categoryCount := make(map[string]int) // 分类计数器
	tagCount := make(map[string]int)      // 标签计数器
	totalViews := 0                       // 总浏览量计数器
	totalLikes := 0                       // 总点赞数计数器

	// 遍历所有文章，计算统计信息
	for _, post := range posts {
		// 累加浏览量和点赞数
		totalViews += post.ViewCount
		totalLikes += post.LikeCount

		// 统计分类分布
		if post.Category != nil {
			categoryCount[post.Category.Name]++
		}

		// 统计标签分布
		for _, tag := range post.Tags {
			tagCount[tag.Name]++
		}
	}

	// 更新统计结果
	stats["total_views"] = totalViews   // 设置总浏览量
	stats["total_likes"] = totalLikes   // 设置总点赞数
	stats["categories"] = categoryCount // 设置分类统计
	stats["tags"] = tagCount            // 设置标签统计

	return stats, nil
}

// ========================================
// 辅助函数
// ========================================

// min 返回两个整数中的较小值
// 这是一个通用的辅助函数，用于处理字符串截取等操作
// 参数:
//
//	a: 第一个整数
//	b: 第二个整数
//
// 返回: int 较小的整数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ========================================
// 主函数演示
// ========================================

// main 主函数，演示GORM关联关系的各种操作
// 包括一对一、一对多、多对多关系的创建、查询、更新等操作
// 同时演示SQLite和MySQL两种数据库的使用
func main() {
	fmt.Println("=== GORM Level 2 关联关系练习 ===")

	// ========================================
	// SQLite 数据库演示
	// ========================================
	fmt.Println("\n--- SQLite 数据库演示 ---")

	// 初始化SQLite数据库
	// 使用默认配置连接SQLite数据库
	db := initDB()
	fmt.Println("✓ SQLite数据库初始化完成")

	// 执行SQLite演示
	runDatabaseDemo(db, "SQLite")

	// ========================================
	// MySQL 数据库演示
	// ========================================
	fmt.Println("\n--- MySQL 数据库演示 ---")

	// 初始化MySQL数据库
	// 使用MySQL配置连接数据库
	// 注意：请根据实际情况修改MySQL连接字符串
	mysqlDSN := "root:fastbee@tcp(192.168.100.124:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig := GetMySQLConfig(mysqlDSN)
	mysqlDB := InitDatabase(mysqlConfig)
	if mysqlDB != nil {
		fmt.Println("✓ MySQL数据库初始化完成")
		// 执行MySQL演示
		runDatabaseDemo(mysqlDB, "MySQL")
	} else {
		fmt.Println("MySQL数据库连接失败，跳过MySQL演示")
	}

	fmt.Println("\n=== Level 2 关联关系练习完成 ===")
}

// runDatabaseDemo 运行数据库演示
// 这个函数包含了所有关联关系的演示代码，可以在不同数据库上运行
// 参数:
//
//	db: 数据库连接实例
//	dbType: 数据库类型名称（用于显示）
func runDatabaseDemo(db *gorm.DB, dbType string) {
	fmt.Printf("\n=== %s 关联关系演示 ===\n", dbType)

	// ========================================
	// 练习1：一对一关系 (User <-> Profile)
	// ========================================
	fmt.Println("\n=== 一对一关系练习 ===")

	// 创建用户及资料
	// 演示如何同时创建用户和关联的个人资料
	user1, err := CreateUserWithProfile(db, "alice", "alice@example.com", "password123", "Alice", "Smith", "我是Alice，一名开发者")
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// 创建第二个用户，用于后续演示
	user2, err := CreateUserWithProfile(db, "bob", "bob@example.com", "password456", "Bob", "Johnson", "我是Bob，喜欢写作")
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// 获取用户及资料
	// 演示如何使用Preload预加载关联数据
	if user1 != nil {
		fetchedUser, err := GetUserWithProfile(db, user1.ID)
		if err != nil {
			fmt.Printf("获取用户失败: %v\n", err)
		} else {
			fmt.Printf("用户: %s, 全名: %s %s\n", fetchedUser.Username, fetchedUser.Profile.FirstName, fetchedUser.Profile.LastName)
		}
	}

	// 更新用户资料
	// 演示如何更新关联表中的数据
	if user1 != nil {
		updates := map[string]interface{}{
			"bio":   "我是Alice，一名全栈开发者",
			"phone": "123-456-7890",
		}
		if err := UpdateUserProfile(db, user1.ID, updates); err != nil {
			fmt.Printf("更新用户资料失败: %v\n", err)
		}
	}

	// ========================================
	// 练习2：一对多关系 (Category -> Posts, Post -> Comments)
	// ========================================
	fmt.Println("\n=== 一对多关系练习 ===")

	// 创建分类及文章
	// 演示如何创建一个分类并同时创建多篇关联文章
	if user1 != nil {
		category1, err := CreateCategoryWithPosts(db, "技术", "tech", user1.ID, []string{"Go语言入门", "GORM使用指南", "数据库设计原则"})
		if err != nil {
			fmt.Printf("创建分类失败: %v\n", err)
		}

		// 获取分类及文章
		// 演示如何查询分类并预加载所有关联的文章
		if category1 != nil {
			fetchedCategory, err := GetCategoryWithPosts(db, category1.ID)
			if err != nil {
				fmt.Printf("获取分类失败: %v\n", err)
			} else {
				fmt.Printf("分类: %s, 文章数量: %d\n", fetchedCategory.Name, len(fetchedCategory.Posts))
				// 遍历显示分类下的所有文章
				for _, post := range fetchedCategory.Posts {
					fmt.Printf("  - %s (作者: %s)\n", post.Title, post.Author.Username)
				}
			}
		}
	}

	// 创建文章及评论
	// 演示如何创建一篇文章并同时创建多条关联评论
	if user1 != nil && user2 != nil {
		post, err := CreatePostWithComments(db, "GORM高级用法", "这是一篇关于GORM高级用法的文章...", user1.ID, 1, []string{"很有用的文章！", "学到了很多", "期待更多内容"})
		if err != nil {
			fmt.Printf("创建文章失败: %v\n", err)
		} else {
			// 获取文章及评论
			// 演示如何查询文章并预加载所有关联的评论
			fetchedPost, err := GetPostWithComments(db, post.ID)
			if err != nil {
				fmt.Printf("获取文章失败: %v\n", err)
			} else {
				fmt.Printf("文章: %s, 评论数量: %d\n", fetchedPost.Title, len(fetchedPost.Comments))
			}
		}
	}

	// ========================================
	// 练习3：多对多关系 (Posts <-> Tags)
	// ========================================
	fmt.Println("\n=== 多对多关系练习 ===")

	// 创建标签并分配给文章
	// 演示如何创建多个标签并将它们关联到多篇文章
	tagNames := []string{"Go", "数据库", "后端", "教程"}
	postIDs := []uint{1, 2, 3, 4}
	if err := CreateTagsAndAssignToPosts(db, tagNames, postIDs); err != nil {
		fmt.Printf("创建标签失败: %v\n", err)
	}

	// 获取文章及标签
	// 演示如何查询文章并预加载所有关联的标签
	post, err := GetPostWithTags(db, 1)
	if err != nil {
		fmt.Printf("获取文章标签失败: %v\n", err)
	} else {
		fmt.Printf("文章: %s, 标签: ", post.Title)
		// 遍历显示文章的所有标签
		for i, tag := range post.Tags {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(tag.Name)
		}
		fmt.Println()
	}

	// 获取标签及文章
	// 演示如何查询标签并预加载所有关联的文章
	tag, err := GetTagWithPosts(db, 1)
	if err != nil {
		fmt.Printf("获取标签文章失败: %v\n", err)
	} else {
		fmt.Printf("标签: %s, 文章数量: %d\n", tag.Name, len(tag.Posts))
	}

	// ========================================
	// 练习4：复杂查询和统计分析
	// ========================================
	fmt.Println("\n=== 复杂查询练习 ===")

	// 获取文章及所有关联数据
	// 演示如何一次性预加载文章的所有关联数据（作者、分类、标签、评论等）
	posts, err := GetPostsWithAllAssociations(db, 5)
	if err != nil {
		fmt.Printf("获取文章失败: %v\n", err)
	} else {
		fmt.Printf("获取到 %d 篇文章（包含所有关联数据）\n", len(posts))
	}

	// 获取用户文章统计
	// 演示如何进行复杂的统计查询，包括计算总数、分组统计等
	if user1 != nil {
		stats, err := GetUserPostsWithStats(db, user1.ID)
		if err != nil {
			fmt.Printf("获取用户统计失败: %v\n", err)
		} else {
			fmt.Printf("用户 %s 的文章统计:\n", user1.Username)
			fmt.Printf("  总文章数: %v\n", stats["total_posts"]) // 显示文章总数
			fmt.Printf("  总浏览数: %v\n", stats["total_views"]) // 显示总浏览量
			fmt.Printf("  总点赞数: %v\n", stats["total_likes"]) // 显示总点赞数
			fmt.Printf("  分类分布: %v\n", stats["categories"])  // 显示分类分布统计
			fmt.Printf("  标签分布: %v\n", stats["tags"])        // 显示标签分布统计
		}
	}

	// 演示完成提示
	fmt.Printf("\n=== %s 关联关系演示完成 ===\n", dbType)
}

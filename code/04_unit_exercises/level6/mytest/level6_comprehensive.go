// 04_unit_exercises/level6_comprehensive.go - Level 6 综合实战练习
// 对应文档：03_GORM单元练习_基础技能训练.md
// 本文件实现了GORM的综合练习，包括复杂的业务模型、高级查询、性能优化、事务处理等
// 支持SQLite和MySQL两种数据库类型
// 涵盖了博客系统的完整功能：用户管理、文章发布、评论系统、点赞关注、通知推送、数据分析等

package main

import (
	"fmt"       // 格式化输出
	"log"       // 日志记录
	"math/rand" // 随机数生成
	"time"      // 时间处理

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
	Type DatabaseType // 数据库类型(sqlite/mysql)
	DSN  string       // 数据源名称,用于指定数据库连接字符串
	// MySQL专用配置字段
	Host     string // MySQL主机地址
	Port     int    // MySQL端口号
	Username string // MySQL用户名
	Password string // MySQL密码
	Database string // MySQL数据库名
	// 连接池配置
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
		Database:     "level6_comprehensive.db", // SQLite数据库文件名
		MaxOpenConns: 10,                        // 最大连接数10
		MaxIdleConns: 5,                         // 最大空闲连接5
		MaxLifetime:  time.Hour,                 // 连接生命周期1小时
		LogLevel:     logger.Info,               // 日志级别为Info
	}
}

// GetMySQLConfig 获取MySQL配置
// 参数host: MySQL主机地址, port: 端口号, username: 用户名, password: 密码, database: 数据库名
// 返回一个包含默认参数的MySQL数据库配置对象
func GetMySQLConfig(host string, port int, username, password, database string) *DatabaseConfig {
	return &DatabaseConfig{
		Type:         MySQL,
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		Database:     database,
		MaxOpenConns: 20,        // MySQL建议更高的连接数
		MaxIdleConns: 10,        // 更多的空闲连接
		MaxLifetime:  time.Hour, // 连接生命周期1小时
		LogLevel:     logger.Info,
	}
}

// GetMySQLConfigFromDSN 从DSN字符串获取MySQL配置
// 参数dsn: MySQL数据库连接字符串
// 返回一个包含默认参数的MySQL数据库配置对象
func GetMySQLConfigFromDSN(dsn string) *DatabaseConfig {
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
// 所有业务模型都应该嵌入此结构体以获得统一的基础字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`              // 主键ID，自动递增
	CreatedAt time.Time      `json:"created_at"`                        // 创建时间，GORM自动管理
	UpdatedAt time.Time      `json:"updated_at"`                        // 更新时间，GORM自动管理
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 删除时间，用于软删除，建立索引
}

// 业务模型定义

// User 用户模型
// 表示博客系统中的用户实体，包含用户的基本信息、认证信息、个人资料和统计数据
// 支持用户注册、登录、个人资料管理、社交功能等完整的用户管理功能
type User struct {
	BaseModel                // 嵌入基础模型，获得ID、时间戳等通用字段
	Username      string     `gorm:"uniqueIndex:idx_username;size:50;not null" json:"username"` // 用户名，唯一索引，最大50字符，不能为空
	Email         string     `gorm:"uniqueIndex:idx_email;size:100;not null" json:"email"`      // 邮箱地址，唯一索引，最大100字符，不能为空
	PasswordHash  string     `gorm:"size:255;not null" json:"-"`                                // 密码哈希值，最大255字符，JSON序列化时忽略(安全考虑)
	FirstName     string     `gorm:"size:50;not null" json:"first_name"`                        // 名字，最大50字符，不能为空
	LastName      string     `gorm:"size:50;not null" json:"last_name"`                         // 姓氏，最大50字符，不能为空
	Avatar        string     `gorm:"size:255" json:"avatar"`                                    // 头像URL，最大255字符，可为空
	Bio           string     `gorm:"type:text" json:"bio"`                                      // 个人简介，文本类型，可为空
	Website       string     `gorm:"size:255" json:"website"`                                   // 个人网站，最大255字符，可为空
	Location      string     `gorm:"size:100" json:"location"`                                  // 所在地，最大100字符，可为空
	BirthDate     *time.Time `json:"birth_date"`                                                // 出生日期，指针类型允许为空
	Gender        string     `gorm:"size:10" json:"gender"`                                     // 性别，最大10字符，可为空
	Phone         string     `gorm:"size:20" json:"phone"`                                      // 电话号码，最大20字符，可为空
	Status        string     `gorm:"size:20;default:'active';index:idx_status" json:"status"`   // 用户状态(active/inactive/banned)，默认active，建立索引
	Role          string     `gorm:"size:20;default:'user';index:idx_role" json:"role"`         // 用户角色(user/admin/moderator)，默认user，建立索引
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`                       // 邮箱是否已验证，默认false
	LastLoginAt   *time.Time `gorm:"index:idx_last_login" json:"last_login_at"`                 // 最后登录时间，指针类型允许为空，建立索引用于查询活跃用户
	LoginCount    int        `gorm:"default:0" json:"login_count"`                              // 登录次数统计，默认0

	// 统计字段 - 用于快速查询用户的内容统计，避免复杂的聚合查询
	PostCount      int `gorm:"default:0" json:"post_count"`      // 发布文章数量，默认0
	CommentCount   int `gorm:"default:0" json:"comment_count"`   // 发布评论数量，默认0
	FollowerCount  int `gorm:"default:0" json:"follower_count"`  // 粉丝数量，默认0
	FollowingCount int `gorm:"default:0" json:"following_count"` // 关注数量，默认0

	// 关联关系 - 定义用户与其他实体的关联，使用外键建立关系
	Profile       *UserProfile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`        // 一对一：用户详细资料
	Posts         []Post         `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`        // 一对多：用户发布的文章
	Comments      []Comment      `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`     // 一对多：用户发布的评论
	Likes         []Like         `gorm:"foreignKey:UserID" json:"likes,omitempty"`          // 一对多：用户的点赞记录
	Followers     []Follow       `gorm:"foreignKey:FollowingID" json:"followers,omitempty"` // 一对多：关注该用户的人(粉丝)
	Following     []Follow       `gorm:"foreignKey:FollowerID" json:"following,omitempty"`  // 一对多：该用户关注的人
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`  // 一对多：用户收到的通知
}

// UserProfile 用户资料扩展模型
// 存储用户的详细个人资料信息，与User模型形成一对一关系
// 包含职业信息、教育背景、技能特长、社交链接等扩展信息
type UserProfile struct {
	BaseModel           // 嵌入基础模型
	UserID       uint   `gorm:"uniqueIndex:idx_user_profile;not null" json:"user_id"` // 用户ID，外键关联User表，唯一索引确保一对一关系
	Company      string `gorm:"size:100" json:"company"`                              // 公司名称，最大100字符
	JobTitle     string `gorm:"size:100" json:"job_title"`                            // 职位名称，最大100字符
	Education    string `gorm:"size:200" json:"education"`                            // 教育背景，最大200字符
	Skills       string `gorm:"type:text" json:"skills"`                              // 技能列表，文本类型，可存储JSON格式
	Experience   int    `gorm:"default:0" json:"experience"`                          // 工作经验年数，默认0
	SalaryRange  string `gorm:"size:50" json:"salary_range"`                          // 薪资范围，最大50字符
	Languages    string `gorm:"type:text" json:"languages"`                           // 语言能力，文本类型，可存储JSON格式
	Interests    string `gorm:"type:text" json:"interests"`                           // 兴趣爱好，文本类型，可存储JSON格式
	SocialLinks  string `gorm:"type:text" json:"social_links"`                        // 社交媒体链接，文本类型，可存储JSON格式
	PrivacyLevel string `gorm:"size:20;default:'public'" json:"privacy_level"`        // 隐私级别(public/private/friends)，默认public

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 反向关联到User模型
}

// Category 分类模型
// 表示博客文章的分类体系，支持层级结构(树形结构)
// 每个分类可以有父分类和子分类，形成完整的分类层次
type Category struct {
	BaseModel          // 嵌入基础模型
	Name        string `gorm:"size:100;not null;index:idx_category_name" json:"name"`       // 分类名称，最大100字符，不能为空，建立索引
	Slug        string `gorm:"uniqueIndex:idx_category_slug;size:100;not null" json:"slug"` // URL友好的标识符，唯一索引，用于SEO
	Description string `gorm:"type:text" json:"description"`                                // 分类描述，文本类型
	Icon        string `gorm:"size:100" json:"icon"`                                        // 分类图标，最大100字符
	Color       string `gorm:"size:7;default:'#007bff'" json:"color"`                       // 分类颜色，7字符十六进制颜色值，默认蓝色
	ParentID    *uint  `gorm:"index:idx_parent" json:"parent_id"`                           // 父分类ID，指针类型允许为空(顶级分类)，建立索引
	Level       int    `gorm:"default:1;index:idx_level" json:"level"`                      // 分类层级，默认1(顶级)，建立索引用于层级查询
	SortOrder   int    `gorm:"default:0;index:idx_sort" json:"sort_order"`                  // 排序顺序，默认0，建立索引用于排序
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`              // 是否激活，默认true，建立索引用于过滤
	PostCount   int    `gorm:"default:0" json:"post_count"`                                 // 该分类下的文章数量，默认0

	// 关联关系 - 实现树形结构的自关联
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // 父分类，一对一关联
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子分类列表，一对多关联
	Posts    []Post     `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`  // 该分类下的文章列表，一对多关联
}

// Tag 标签模型
// 表示博客文章的标签系统，与文章形成多对多关系
// 标签用于对文章进行更细粒度的分类和检索
type Tag struct {
	BaseModel          // 嵌入基础模型
	Name        string `gorm:"uniqueIndex:idx_tag_name;size:50;not null" json:"name"` // 标签名称，唯一索引，最大50字符，不能为空
	Slug        string `gorm:"uniqueIndex:idx_tag_slug;size:50;not null" json:"slug"` // URL友好的标识符，唯一索引，用于SEO
	Description string `gorm:"type:text" json:"description"`                          // 标签描述，文本类型
	Color       string `gorm:"size:7;default:'#007bff'" json:"color"`                 // 标签颜色，7字符十六进制颜色值，默认蓝色
	UsageCount  int    `gorm:"default:0;index:idx_usage" json:"usage_count"`          // 使用次数统计，默认0，建立索引用于热门标签查询
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`        // 是否激活，默认true，建立索引用于过滤

	// 关联关系 - 多对多关系
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"` // 使用该标签的文章列表，多对多关联，中间表为post_tags
}

// Post 文章模型
// 表示博客系统中的文章实体，是整个系统的核心内容模型
// 包含文章的基本信息、内容、状态、统计数据、SEO信息等完整字段
type Post struct {
	BaseModel                  // 嵌入基础模型
	Title           string     `gorm:"size:200;not null;index:idx_title" json:"title"`               // 文章标题，最大200字符，不能为空，建立索引
	Slug            string     `gorm:"uniqueIndex:idx_post_slug;size:200;not null" json:"slug"`      // URL友好的标识符，唯一索引，用于SEO和路由
	Content         string     `gorm:"type:text;not null" json:"content"`                            // 文章内容，文本类型，不能为空
	Excerpt         string     `gorm:"size:500" json:"excerpt"`                                      // 文章摘要，最大500字符，用于列表显示
	FeaturedImage   string     `gorm:"size:255" json:"featured_image"`                               // 特色图片URL，最大255字符
	Status          string     `gorm:"size:20;default:'draft';index:idx_status" json:"status"`       // 文章状态(draft/published/private)，默认draft，建立索引
	Type            string     `gorm:"size:20;default:'post';index:idx_type" json:"type"`            // 文章类型(post/page/custom)，默认post，建立索引
	Format          string     `gorm:"size:20;default:'standard'" json:"format"`                     // 文章格式(standard/gallery/video等)，默认standard
	ViewCount       int        `gorm:"default:0;index:idx_views" json:"view_count"`                  // 浏览次数，默认0，建立索引用于热门文章查询
	LikeCount       int        `gorm:"default:0;index:idx_likes" json:"like_count"`                  // 点赞次数，默认0，建立索引用于热门文章查询
	CommentCount    int        `gorm:"default:0;index:idx_comments" json:"comment_count"`            // 评论次数，默认0，建立索引用于活跃文章查询
	ShareCount      int        `gorm:"default:0" json:"share_count"`                                 // 分享次数，默认0
	PublishedAt     *time.Time `gorm:"index:idx_published" json:"published_at"`                      // 发布时间，指针类型允许为空，建立索引用于时间排序
	Rating          float64    `gorm:"precision:3;scale:2;default:0;index:idx_rating" json:"rating"` // 文章评分，精度3位小数2位，默认0，建立索引
	Featured        bool       `gorm:"default:false;index:idx_featured" json:"featured"`             // 是否为特色文章，默认false，建立索引
	Sticky          bool       `gorm:"default:false;index:idx_sticky" json:"sticky"`                 // 是否置顶，默认false，建立索引
	AllowComments   bool       `gorm:"default:true" json:"allow_comments"`                           // 是否允许评论，默认true
	MetaTitle       string     `gorm:"size:200" json:"meta_title"`                                   // SEO标题，最大200字符
	MetaDescription string     `gorm:"size:500" json:"meta_description"`                             // SEO描述，最大500字符
	MetaKeywords    string     `gorm:"size:255" json:"meta_keywords"`                                // SEO关键词，最大255字符

	// 外键字段 - 建立与其他表的关联
	AuthorID   uint  `gorm:"not null;index:idx_author" json:"author_id"` // 作者ID，外键关联User表，不能为空，建立索引
	CategoryID *uint `gorm:"index:idx_category" json:"category_id"`      // 分类ID，外键关联Category表，指针类型允许为空，建立索引

	// 关联关系 - 定义与其他模型的关联
	Author   User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`     // 文章作者，多对一关联
	Category *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 文章分类，多对一关联
	Comments []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`     // 文章评论列表，一对多关联
	Likes    []Like     `gorm:"foreignKey:PostID" json:"likes,omitempty"`        // 文章点赞列表，一对多关联
	Tags     []Tag      `gorm:"many2many:post_tags;" json:"tags,omitempty"`      // 文章标签列表，多对多关联
	Meta     []PostMeta `gorm:"foreignKey:PostID" json:"meta,omitempty"`         // 文章元数据列表，一对多关联
}

// PostMeta 文章元数据模型
// 用于存储文章的扩展属性和自定义字段，提供灵活的键值对存储机制
// 例如：自定义字段、第三方插件数据、临时配置等
type PostMeta struct {
	BaseModel        // 嵌入基础模型
	PostID    uint   `gorm:"not null;index:idx_post_meta" json:"post_id"`          // 文章ID，外键关联Post表，不能为空，建立索引
	MetaKey   string `gorm:"size:100;not null;index:idx_meta_key" json:"meta_key"` // 元数据键名，最大100字符，不能为空，建立索引用于快速查找
	MetaValue string `gorm:"type:text" json:"meta_value"`                          // 元数据值，文本类型，可存储大量数据

	// 关联关系 - 定义与其他模型的关联
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"` // 所属文章，多对一关联
}

// Comment 评论模型
// 表示博客系统中的评论实体，支持嵌套回复和多种状态管理
// 包含评论内容、作者信息、审核状态、层级关系等完整功能
type Comment struct {
	BaseModel        // 嵌入基础模型
	Content   string `gorm:"type:text;not null" json:"content"`                        // 评论内容，文本类型，不能为空
	Status    string `gorm:"size:20;default:'pending';index:idx_status" json:"status"` // 评论状态(pending/approved/spam/trash)，默认pending，建立索引
	Type      string `gorm:"size:20;default:'comment'" json:"type"`                    // 评论类型(comment/pingback/trackback)，默认comment
	LikeCount int    `gorm:"default:0;index:idx_likes" json:"like_count"`              // 点赞次数，默认0，建立索引
	ParentID  *uint  `gorm:"index:idx_parent" json:"parent_id"`                        // 父评论ID，外键关联Comment表，指针类型允许为空(顶级评论)，建立索引
	Level     int    `gorm:"default:1;index:idx_level" json:"level"`                   // 评论层级，默认1(顶级评论)，建立索引用于层级查询
	UserAgent string `gorm:"size:255" json:"user_agent"`                               // 用户代理字符串，最大255字符
	UserIP    string `gorm:"size:45" json:"user_ip"`                                   // 用户IP地址，最大45字符(支持IPv6)
	IsSpam    bool   `gorm:"default:false;index:idx_spam" json:"is_spam"`              // 是否为垃圾评论，默认false，建立索引

	// 外键字段 - 建立与其他表的关联
	PostID   uint `gorm:"not null;index:idx_post" json:"post_id"`     // 文章ID，外键关联Post表，不能为空，建立索引
	AuthorID uint `gorm:"not null;index:idx_author" json:"author_id"` // 作者ID，外键关联User表，不能为空，建立索引

	// 关联关系 - 定义与其他模型的关联
	Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`       // 所属文章，多对一关联
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`   // 评论作者，多对一关联
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // 父评论，多对一关联(自关联)
	Children []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子评论列表，一对多关联(自关联)
	Likes    []Like    `gorm:"foreignKey:CommentID" json:"likes,omitempty"`   // 评论点赞列表，一对多关联
}

// Like 点赞模型
// 表示用户对文章或评论的点赞行为，支持多种点赞类型
// 通过PostID和CommentID的互斥性实现对不同内容类型的点赞
type Like struct {
	BaseModel        // 嵌入基础模型
	UserID    uint   `gorm:"not null;index:idx_user_like" json:"user_id"`            // 用户ID，外键关联User表，不能为空，建立索引
	PostID    *uint  `gorm:"index:idx_post_like" json:"post_id"`                     // 文章ID，外键关联Post表，指针类型允许为空(评论点赞时为空)，建立索引
	CommentID *uint  `gorm:"index:idx_comment_like" json:"comment_id"`               // 评论ID，外键关联Comment表，指针类型允许为空(文章点赞时为空)，建立索引
	Type      string `gorm:"size:20;default:'like';index:idx_like_type" json:"type"` // 点赞类型(like/dislike/love等)，默认like，建立索引
	UserIP    string `gorm:"size:45" json:"user_ip"`                                 // 用户IP地址，最大45字符(支持IPv6)，用于防刷和统计

	// 关联关系 - 定义与其他模型的关联
	User    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`       // 点赞用户，多对一关联
	Post    *Post    `gorm:"foreignKey:PostID" json:"post,omitempty"`       // 被点赞的文章，多对一关联
	Comment *Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"` // 被点赞的评论，多对一关联
}

// Follow 关注关系模型
// 表示用户之间的关注关系，实现社交网络功能
// 支持关注状态管理，可扩展为好友、黑名单等关系类型
type Follow struct {
	BaseModel          // 嵌入基础模型
	FollowerID  uint   `gorm:"not null;index:idx_follower" json:"follower_id"`   // 关注者ID，外键关联User表，不能为空，建立索引
	FollowingID uint   `gorm:"not null;index:idx_following" json:"following_id"` // 被关注者ID，外键关联User表，不能为空，建立索引
	Status      string `gorm:"size:20;default:'active'" json:"status"`           // 关注状态(active/blocked/pending)，默认active

	// 关联关系 - 定义与其他模型的关联
	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`   // 关注者，多对一关联
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"` // 被关注者，多对一关联
}

// Notification 通知模型
// 表示系统通知实体，支持多种通知类型和状态管理
// 包含通知内容、读取状态、关联对象等完整信息
type Notification struct {
	BaseModel              // 嵌入基础模型
	UserID      uint       `gorm:"not null;index:idx_user_notification" json:"user_id"`      // 接收通知的用户ID，外键关联User表，不能为空，建立索引
	Type        string     `gorm:"size:50;not null;index:idx_notification_type" json:"type"` // 通知类型(comment/like/follow/system等)，最大50字符，不能为空，建立索引
	Title       string     `gorm:"size:200;not null" json:"title"`                           // 通知标题，最大200字符，不能为空
	Content     string     `gorm:"type:text" json:"content"`                                 // 通知内容，文本类型
	Data        string     `gorm:"type:text" json:"data"`                                    // 通知附加数据(JSON格式)，文本类型
	IsRead      bool       `gorm:"default:false;index:idx_read" json:"is_read"`              // 是否已读，默认false，建立索引用于未读通知查询
	ReadAt      *time.Time `json:"read_at"`                                                  // 读取时间，指针类型允许为空
	RelatedID   *uint      `gorm:"index:idx_related" json:"related_id"`                      // 关联对象ID，指针类型允许为空，建立索引
	RelatedType string     `gorm:"size:50" json:"related_type"`                              // 关联对象类型(post/comment/user等)，最大50字符

	// 关联关系 - 定义与其他模型的关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 接收通知的用户，多对一关联
}

// Setting 系统设置模型
// 表示系统配置项，提供灵活的键值对配置管理
// 支持不同数据类型、分组管理和权限控制
type Setting struct {
	BaseModel          // 嵌入基础模型
	Key         string `gorm:"uniqueIndex:idx_setting_key;size:100;not null" json:"key"` // 配置键名，最大100字符，不能为空，唯一索引
	Value       string `gorm:"type:text" json:"value"`                                   // 配置值，文本类型，可存储各种格式数据
	Type        string `gorm:"size:20;default:'string'" json:"type"`                     // 数据类型(string/int/bool/json等)，默认string
	Description string `gorm:"size:255" json:"description"`                              // 配置描述，最大255字符
	Group       string `gorm:"size:50;index:idx_setting_group" json:"group"`             // 配置分组，最大50字符，建立索引用于分组查询
	IsPublic    bool   `gorm:"default:false" json:"is_public"`                           // 是否为公开配置，默认false，用于权限控制
}

// initDB 初始化数据库连接和配置
// 支持SQLite和MySQL两种数据库类型，根据配置自动选择
// 包含连接池配置、自动迁移、索引创建等完整的数据库初始化流程
func initDB(config DatabaseConfig) *gorm.DB {
	var db *gorm.DB
	var err error

	// 根据数据库类型选择不同的驱动和连接字符串
	switch config.Type {
	case MySQL:
		// MySQL数据库连接
		var dsn string
		if config.DSN != "" {
			// 如果提供了完整的DSN字符串，直接使用
			dsn = config.DSN
		} else {
			// 否则根据配置参数构建DSN
			// DSN格式: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.Username, config.Password, config.Host, config.Port, config.Database)
		}
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(config.LogLevel), // 使用配置的日志级别
			PrepareStmt:                              true,                                    // 启用预编译语句，提高性能
			DisableForeignKeyConstraintWhenMigrating: false,                                   // 迁移时保留外键约束
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	case SQLite:
		// SQLite数据库连接（默认选项）
		// 适用于开发环境和小型应用
		dbPath := config.Database
		if dbPath == "" {
			dbPath = "level6_comprehensive.db" // 默认数据库文件名
		}
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                                   logger.Default.LogMode(config.LogLevel), // 使用配置的日志级别
			PrepareStmt:                              true,                                    // 启用预编译语句
			DisableForeignKeyConstraintWhenMigrating: false,                                   // 保留外键约束
		})
	default:
		// 不支持的数据库类型，使用SQLite作为默认选项
		log.Printf("不支持的数据库类型: %v，使用SQLite作为默认选项", config.Type)
		db, err = gorm.Open(sqlite.Open("level6_comprehensive.db"), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: false,
		})
	}

	// 检查数据库连接是否成功
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 获取底层的sql.DB对象，用于配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接失败:", err)
	}

	// 配置数据库连接池参数
	// SetMaxIdleConns: 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	// SetMaxOpenConns: 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	// SetConnMaxLifetime: 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	// 自动迁移数据库表结构
	// 按照依赖关系的顺序进行迁移，确保外键关系正确建立
	err = db.AutoMigrate(
		&User{},         // 用户表（基础表）
		&UserProfile{},  // 用户资料表（依赖User）
		&Category{},     // 分类表（自引用表）
		&Tag{},          // 标签表（独立表）
		&Post{},         // 文章表（依赖User和Category）
		&PostMeta{},     // 文章元数据表（依赖Post）
		&Comment{},      // 评论表（依赖Post和User）
		&Like{},         // 点赞表（依赖User、Post、Comment）
		&Follow{},       // 关注表（依赖User）
		&Notification{}, // 通知表（依赖User）
		&Setting{},      // 设置表（依赖User）
	)
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 创建数据库索引以优化查询性能
	createIndexes(db)

	// 打印初始化成功信息
	fmt.Printf("✓ 数据库初始化完成 (类型: %v)\n", config.Type)
	return db
}

// createIndexes 创建数据库索引以优化查询性能
// 包括复合索引、唯一索引等，根据常见查询模式设计
// 索引的创建遵循"查询优先"原则，针对高频查询字段组合建立索引
func createIndexes(db *gorm.DB) {
	// 复合索引 - 针对多字段组合查询优化
	// 用户状态和角色的复合索引，用于用户管理和权限控制查询
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status_role ON users(status, role)")

	// 文章状态和发布时间的复合索引，用于文章列表查询和排序
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_status_published ON posts(status, published_at)")

	// 文章作者和分类的复合索引，用于按作者或分类筛选文章
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_author_category ON posts(author_id, category_id)")

	// 评论文章和状态的复合索引，用于获取特定文章的已审核评论
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_status ON comments(post_id, status)")

	// 点赞用户和文章的复合索引，用于快速查询用户点赞记录
	db.Exec("CREATE INDEX IF NOT EXISTS idx_likes_user_post ON likes(user_id, post_id)")

	// 关注关系的复合索引，用于查询关注者和被关注者关系
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_follower_following ON follows(follower_id, following_id)")

	// 通知用户和读取状态的复合索引，用于获取用户未读通知
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read)")

	// 唯一索引 - 防止重复数据，确保数据完整性
	// 点赞唯一性约束：同一用户不能对同一文章或评论重复点赞
	// WHERE条件确保至少有一个目标对象（文章或评论）不为空
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_like ON likes(user_id, post_id, comment_id) WHERE post_id IS NOT NULL OR comment_id IS NOT NULL")

	// 关注关系唯一性约束：防止重复关注同一用户
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_follow ON follows(follower_id, following_id)")

	fmt.Println("✓ 索引创建完成")
}

// ==================== 钩子函数实现 ====================
// GORM钩子函数在数据库操作的特定时机自动执行
// 用于实现业务逻辑、数据验证、统计更新等功能

// ==================== 用户模型钩子 ====================

// AfterCreate 用户创建后的钩子函数
// 在用户记录成功插入数据库后自动执行
// 主要用于初始化用户相关的关联数据
func (u *User) AfterCreate(tx *gorm.DB) error {
	// 自动为新用户创建默认的用户资料记录
	// 确保每个用户都有对应的详细资料信息
	profile := UserProfile{
		UserID:       u.ID,     // 关联到当前用户
		PrivacyLevel: "public", // 默认隐私级别为公开
	}
	// 在同一事务中创建用户资料，确保数据一致性
	return tx.Create(&profile).Error
}

// BeforeUpdate 用户更新前的钩子函数
// 在用户记录更新前自动执行，用于实时计算和更新统计信息
// 确保用户的统计数据始终与实际数据保持同步
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 只有当用户ID有效时才进行统计更新
	if u.ID != 0 {
		// 统计并更新用户的已发布文章数量
		// 只计算状态为"published"的文章
		var postCount int64
		tx.Model(&Post{}).Where("author_id = ? AND status = ?", u.ID, "published").Count(&postCount)
		u.PostCount = int(postCount)

		// 统计并更新用户的已审核评论数量
		// 只计算状态为"approved"的评论
		var commentCount int64
		tx.Model(&Comment{}).Where("author_id = ? AND status = ?", u.ID, "approved").Count(&commentCount)
		u.CommentCount = int(commentCount)

		// 统计并更新用户的关注者和关注数量
		var followerCount, followingCount int64
		// 统计有多少人关注了当前用户（粉丝数）
		tx.Model(&Follow{}).Where("following_id = ?", u.ID).Count(&followerCount)
		// 统计当前用户关注了多少人（关注数）
		tx.Model(&Follow{}).Where("follower_id = ?", u.ID).Count(&followingCount)
		u.FollowerCount = int(followerCount)
		u.FollowingCount = int(followingCount)
	}
	return nil
}

// ==================== 文章模型钩子 ====================

// BeforeCreate 文章创建前的钩子函数
// 在文章记录插入数据库前自动执行
// 用于自动生成摘要、设置发布时间等预处理操作
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 自动生成文章摘要
	// 如果没有手动设置摘要且文章内容超过200字符，则自动截取前200字符作为摘要
	if p.Excerpt == "" && len(p.Content) > 200 {
		p.Excerpt = p.Content[:200] + "..." // 添加省略号表示内容被截断
	}

	// 自动设置文章发布时间
	// 如果文章状态为已发布但没有设置发布时间，则使用当前时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now // 使用指针类型，允许为空值
	}

	return nil
}

// AfterCreate 文章创建后的钩子函数
// 在文章记录成功插入数据库后自动执行
// 用于更新相关统计信息，维护数据一致性
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新分类的文章数量统计
	// 只有当文章属于某个分类时才更新（CategoryID不为空）
	if p.CategoryID != nil {
		// 使用原子操作增加分类的文章计数，避免并发问题
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))
	}

	// 更新作者的文章数量统计
	// 使用原子操作增加用户的文章计数
	tx.Model(&User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))

	return nil
}

// AfterUpdate 文章更新后的钩子函数
// 在文章记录更新后自动执行
// 用于处理状态变更后的相关逻辑
func (p *Post) AfterUpdate(tx *gorm.DB) error {
	// 处理文章发布状态变更
	// 如果文章状态改为已发布但还没有发布时间，则设置当前时间为发布时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		// 更新发布时间字段
		tx.Model(p).Update("published_at", now)
	}
	return nil
}

// AfterDelete 文章删除后的钩子函数
// 在文章记录被删除后自动执行
// 用于清理相关数据和更新统计信息
func (p *Post) AfterDelete(tx *gorm.DB) error {
	// 减少分类的文章数量统计
	// 只有当文章属于某个分类时才更新
	if p.CategoryID != nil {
		// 使用原子操作减少分类的文章计数
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("post_count", gorm.Expr("post_count - ?", 1))
	}

	// 减少作者的文章数量统计
	// 使用原子操作减少用户的文章计数
	tx.Model(&User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("post_count - ?", 1))

	return nil
}

// ==================== 评论模型钩子 ====================

// AfterCreate 评论创建后的钩子函数
// 在评论记录成功插入数据库后自动执行
// 用于更新统计信息和发送通知
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	// 更新文章的评论数量统计
	// 使用原子操作增加文章的评论计数，确保并发安全
	tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	// 更新用户的评论数量统计
	// 使用原子操作增加用户的评论计数
	tx.Model(&User{}).Where("id = ?", c.AuthorID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	// 创建评论通知
	// 当有人评论文章时，通知文章作者
	var post Post
	if err := tx.First(&post, c.PostID).Error; err == nil {
		// 只有当评论者不是文章作者时才发送通知（避免自己给自己发通知）
		if post.AuthorID != c.AuthorID {
			// 构建通知内容
			notification := Notification{
				UserID:      post.AuthorID,                             // 通知接收者（文章作者）
				Type:        "comment",                                 // 通知类型
				Title:       "新评论",                                     // 通知标题
				Content:     fmt.Sprintf("您的文章《%s》收到了新评论", post.Title), // 通知内容
				RelatedID:   &c.ID,                                     // 关联的评论ID
				RelatedType: "comment",                                 // 关联类型
			}
			// 在同一事务中创建通知，确保数据一致性
			tx.Create(&notification)
		}
	}

	return nil
}

// AfterDelete 评论删除后的钩子函数
// 在评论记录被删除后自动执行
// 用于更新相关统计信息
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 减少文章的评论数量统计
	// 使用原子操作减少文章的评论计数
	tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))

	// 减少用户的评论数量统计
	// 使用原子操作减少用户的评论计数
	tx.Model(&User{}).Where("id = ?", c.AuthorID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))

	return nil
}

// ==================== 点赞模型钩子 ====================

// AfterCreate 点赞创建后的钩子函数
// 在点赞记录成功插入数据库后自动执行
// 用于更新点赞统计和发送通知
func (l *Like) AfterCreate(tx *gorm.DB) error {
	// 处理文章点赞
	if l.PostID != nil {
		// 更新文章的点赞数量统计
		// 使用原子操作增加文章的点赞计数，确保并发安全
		tx.Model(&Post{}).Where("id = ?", *l.PostID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

		// 创建点赞通知
		// 当有人点赞文章时，通知文章作者
		var post Post
		if err := tx.First(&post, *l.PostID).Error; err == nil {
			// 只有当点赞者不是文章作者时才发送通知（避免自己给自己发通知）
			if post.AuthorID != l.UserID {
				// 构建通知内容
				notification := Notification{
					UserID:      post.AuthorID,                             // 通知接收者（文章作者）
					Type:        "like",                                    // 通知类型
					Title:       "新点赞",                                     // 通知标题
					Content:     fmt.Sprintf("您的文章《%s》收到了新点赞", post.Title), // 通知内容
					RelatedID:   l.PostID,                                  // 关联的文章ID
					RelatedType: "post",                                    // 关联类型
				}
				// 在同一事务中创建通知，确保数据一致性
				tx.Create(&notification)
			}
		}
	}

	// 处理评论点赞
	if l.CommentID != nil {
		// 更新评论的点赞数量统计
		// 使用原子操作增加评论的点赞计数
		tx.Model(&Comment{}).Where("id = ?", *l.CommentID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
	}

	return nil
}

// AfterDelete 点赞删除后的钩子函数
// 在点赞记录被删除后自动执行
// 用于更新相关统计信息
func (l *Like) AfterDelete(tx *gorm.DB) error {
	// 处理文章点赞删除
	if l.PostID != nil {
		// 减少文章的点赞数量统计
		// 使用原子操作减少文章的点赞计数
		tx.Model(&Post{}).Where("id = ?", *l.PostID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	}

	// 处理评论点赞删除
	if l.CommentID != nil {
		// 减少评论的点赞数量统计
		// 使用原子操作减少评论的点赞计数
		tx.Model(&Comment{}).Where("id = ?", *l.CommentID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	}

	return nil
}

// ==================== 关注模型钩子 ====================

// AfterCreate 关注创建后的钩子函数
// 在关注记录成功插入数据库后自动执行
// 用于更新关注统计和发送通知
func (f *Follow) AfterCreate(tx *gorm.DB) error {
	// 更新关注数量统计
	// 增加关注者的"正在关注"数量
	tx.Model(&User{}).Where("id = ?", f.FollowerID).UpdateColumn("following_count", gorm.Expr("following_count + ?", 1))
	// 增加被关注者的"粉丝"数量
	tx.Model(&User{}).Where("id = ?", f.FollowingID).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))

	// 创建关注通知
	// 当有人关注用户时，通知被关注的用户
	notification := Notification{
		UserID:      f.FollowingID, // 通知接收者（被关注的用户）
		Type:        "follow",      // 通知类型
		Title:       "新关注者",        // 通知标题
		Content:     "您有新的关注者",     // 通知内容
		RelatedID:   &f.FollowerID, // 关联的关注者ID
		RelatedType: "user",        // 关联类型
	}
	// 在同一事务中创建通知，确保数据一致性
	tx.Create(&notification)

	return nil
}

// AfterDelete 关注删除后的钩子函数
// 在关注记录被删除后自动执行（取消关注）
// 用于更新相关统计信息
func (f *Follow) AfterDelete(tx *gorm.DB) error {
	// 更新关注数量统计
	// 减少关注者的"正在关注"数量
	tx.Model(&User{}).Where("id = ?", f.FollowerID).UpdateColumn("following_count", gorm.Expr("following_count - ?", 1))
	// 减少被关注者的"粉丝"数量
	tx.Model(&User{}).Where("id = ?", f.FollowingID).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))

	return nil
}

// ==================== 业务逻辑函数 ====================
// 以下是应用程序的核心业务逻辑实现
// 采用服务层模式，将业务逻辑与数据访问层分离
// 每个服务类负责特定领域的业务操作

// ==================== 用户管理服务 ====================

// UserService 用户管理服务
// 提供用户相关的所有业务操作，包括注册、登录、资料管理等
// 封装了用户相关的复杂业务逻辑和数据库操作
type UserService struct {
	db *gorm.DB // 数据库连接实例
}

// NewUserService 创建新的用户服务实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *UserService: 用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser 创建新用户
// 使用数据库事务确保数据一致性
// 参数:
//   - user: 要创建的用户对象
//
// 返回:
//   - error: 创建失败时返回错误信息
func (s *UserService) CreateUser(user *User) error {
	// 使用事务确保用户创建的原子性
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}

// GetUserByID 根据用户ID获取用户信息
// 自动预加载用户资料信息
// 参数:
//   - id: 用户ID
//
// 返回:
//   - *User: 用户对象（包含资料信息）
//   - error: 查询失败时返回错误信息
func (s *UserService) GetUserByID(id uint) (*User, error) {
	var user User
	// 预加载用户资料，避免N+1查询问题
	err := s.db.Preload("Profile").First(&user, id).Error
	return &user, err
}

// GetUserWithStats 获取用户详细信息和统计数据
// 包含用户资料、最新发布的文章和最新的评论
// 参数:
//   - id: 用户ID
//
// 返回:
//   - *User: 用户对象（包含统计信息）
//   - error: 查询失败时返回错误信息
func (s *UserService) GetUserWithStats(id uint) (*User, error) {
	var user User
	// 预加载用户的详细信息和统计数据
	err := s.db.Preload("Profile").
		// 预加载最新的5篇已发布文章
		Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "published").Order("created_at DESC").Limit(5)
		}).
		// 预加载最新的5条已审核评论
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "approved").Order("created_at DESC").Limit(5)
		}).
		First(&user, id).Error
	return &user, err
}

// UpdateUserProfile 更新用户资料
// 参数:
//   - userID: 用户ID
//   - profile: 要更新的用户资料对象
//
// 返回:
//   - error: 更新失败时返回错误信息
func (s *UserService) UpdateUserProfile(userID uint, profile *UserProfile) error {
	// 根据用户ID更新用户资料，只更新非零值字段
	return s.db.Model(&UserProfile{}).Where("user_id = ?", userID).Updates(profile).Error
}

// FollowUser 关注用户
// 创建用户之间的关注关系
// 参数:
//   - followerID: 关注者用户ID
//   - followingID: 被关注者用户ID
//
// 返回:
//   - error: 关注失败时返回错误信息
func (s *UserService) FollowUser(followerID, followingID uint) error {
	// 验证不能关注自己
	if followerID == followingID {
		return fmt.Errorf("不能关注自己")
	}

	// 创建关注关系记录
	follow := Follow{
		FollowerID:  followerID,  // 关注者ID
		FollowingID: followingID, // 被关注者ID
		Status:      "active",    // 关注状态
	}

	// 保存关注关系，触发AfterCreate钩子更新统计和发送通知
	return s.db.Create(&follow).Error
}

// UnfollowUser 取消关注用户
// 删除用户之间的关注关系
// 参数:
//   - followerID: 关注者用户ID
//   - followingID: 被关注者用户ID
//
// 返回:
//   - error: 取消关注失败时返回错误信息
func (s *UserService) UnfollowUser(followerID, followingID uint) error {
	// 删除关注关系记录，触发AfterDelete钩子更新统计
	return s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&Follow{}).Error
}

// GetUserFollowers 获取用户的关注者列表
// 支持分页查询
// 参数:
//   - userID: 用户ID
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//
// 返回:
//   - []User: 关注者用户列表
//   - int64: 总关注者数量
//   - error: 查询失败时返回错误信息
func (s *UserService) GetUserFollowers(userID uint, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 获取关注者总数
	// 通过JOIN查询获取关注该用户的所有用户数量
	s.db.Model(&User{}).Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).Count(&total)

	// 获取分页的关注者数据
	// 通过JOIN查询获取关注该用户的用户列表
	err := s.db.Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// ==================== 文章管理服务 ====================

// PostService 文章管理服务
// 提供文章相关的所有业务操作，包括创建、查询、分类管理等
// 封装了文章相关的复杂业务逻辑和数据库操作
type PostService struct {
	db *gorm.DB // 数据库连接实例
}

// NewPostService 创建新的文章服务实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *PostService: 文章服务实例
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// CreatePost 创建新文章
// 使用数据库事务确保文章创建和标签统计更新的原子性
// 参数:
//   - post: 要创建的文章对象
//
// 返回:
//   - error: 创建失败时返回错误信息
func (s *PostService) CreatePost(post *Post) error {
	// 使用事务确保文章创建和相关操作的原子性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建文章记录
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		// 更新关联标签的使用次数统计
		// 遍历文章的所有标签，增加每个标签的使用计数
		for _, tag := range post.Tags {
			tx.Model(&Tag{}).Where("id = ?", tag.ID).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1))
		}

		return nil
	})
}

// GetPostBySlug 根据文章别名获取文章详情
// 预加载文章的所有相关信息，包括作者、分类、标签、评论等
// 同时自动增加文章浏览量
// 参数:
//   - slug: 文章别名（URL友好的标识符）
//
// 返回:
//   - *Post: 文章对象（包含完整信息）
//   - error: 查询失败时返回错误信息
func (s *PostService) GetPostBySlug(slug string) (*Post, error) {
	var post Post
	// 预加载文章的完整信息
	err := s.db.Preload("Author"). // 预加载文章作者信息
					Preload("Category"). // 预加载文章分类信息
					Preload("Tags").     // 预加载文章标签信息
		// 预加载顶级评论（已审核且无父评论）
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ? AND parent_id IS NULL", "approved").Order("created_at ASC")
		}).
		Preload("Comments.Author"). // 预加载评论作者信息
		// 预加载子评论（回复）
		Preload("Comments.Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "approved").Order("created_at ASC")
		}).
		Preload("Comments.Children.Author"). // 预加载子评论作者信息
		Where("slug = ?", slug).First(&post).Error

	// 如果查询成功，自动增加文章浏览量
	if err == nil {
		// 使用原子操作增加浏览量，确保并发安全
		s.db.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	}

	return &post, err
}

// GetPostsByCategory 根据分类获取文章列表
// 支持分页查询，优先显示置顶文章
// 参数:
//   - categorySlug: 分类别名
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//
// 返回:
//   - []Post: 文章列表
//   - int64: 该分类下的文章总数
//   - error: 查询失败时返回错误信息
func (s *PostService) GetPostsByCategory(categorySlug string, page, pageSize int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 获取该分类下已发布文章的总数
	// 通过JOIN查询关联分类表
	s.db.Model(&Post{}).Joins("JOIN categories ON posts.category_id = categories.id").
		Where("categories.slug = ? AND posts.status = ?", categorySlug, "published").Count(&total)

	// 获取分页的文章数据
	// 预加载作者、分类、标签信息，避免N+1查询
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN categories ON posts.category_id = categories.id").
		Where("categories.slug = ? AND posts.status = ?", categorySlug, "published").
		// 排序：置顶文章优先，然后按发布时间倒序
		Order("posts.sticky DESC, posts.published_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error

	return posts, total, err
}

// GetPostsByTag 根据标签获取文章列表
// 支持分页查询，按发布时间倒序排列
// 参数:
//   - tagSlug: 标签别名
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//
// 返回:
//   - []Post: 文章列表
//   - int64: 该标签下的文章总数
//   - error: 查询失败时返回错误信息
func (s *PostService) GetPostsByTag(tagSlug string, page, pageSize int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 获取该标签下已发布文章的总数
	// 通过多表JOIN查询：文章表 -> 文章标签关联表 -> 标签表
	s.db.Model(&Post{}).Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.slug = ? AND posts.status = ?", tagSlug, "published").Count(&total)

	// 获取分页的文章数据
	// 预加载作者、分类、标签信息，避免N+1查询
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.slug = ? AND posts.status = ?", tagSlug, "published").
		// 按发布时间倒序排列
		Order("posts.published_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error

	return posts, total, err
}

// SearchPosts 搜索文章
// 在文章标题和内容中搜索关键词，支持分页查询
// 按浏览量和发布时间排序，热门文章优先
// 参数:
//   - keyword: 搜索关键词
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//
// 返回:
//   - []Post: 匹配的文章列表
//   - int64: 匹配的文章总数
//   - error: 搜索失败时返回错误信息
func (s *PostService) SearchPosts(keyword string, page, pageSize int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	// 计算分页偏移量
	offset := (page - 1) * pageSize
	// 构建模糊搜索条件
	searchTerm := "%" + keyword + "%"

	// 获取匹配的已发布文章总数
	// 在标题和内容中搜索关键词
	s.db.Model(&Post{}).Where("(title LIKE ? OR content LIKE ?) AND status = ?", searchTerm, searchTerm, "published").Count(&total)

	// 获取分页的搜索结果
	// 预加载作者、分类、标签信息
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("(title LIKE ? OR content LIKE ?) AND status = ?", searchTerm, searchTerm, "published").
		// 排序：浏览量高的优先，然后按发布时间倒序
		Order("view_count DESC, published_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error

	return posts, total, err
}

// LikePost 点赞文章
// 检查用户是否已经点赞，避免重复点赞
// 参数:
//   - userID: 用户ID
//   - postID: 文章ID
//
// 返回:
//   - error: 点赞失败时返回错误信息
func (s *PostService) LikePost(userID, postID uint) error {
	// 检查用户是否已经对该文章点赞
	var existingLike Like
	if err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike).Error; err == nil {
		return fmt.Errorf("已经点赞过了")
	}

	// 创建点赞记录
	like := Like{
		UserID: userID,  // 点赞用户ID
		PostID: &postID, // 被点赞的文章ID
		Type:   "like",  // 点赞类型
	}

	// 保存点赞记录，触发AfterCreate钩子更新统计和发送通知
	return s.db.Create(&like).Error
}

// UnlikePost 取消点赞文章
// 删除用户对文章的点赞记录
// 参数:
//   - userID: 用户ID
//   - postID: 文章ID
//
// 返回:
//   - error: 取消点赞失败时返回错误信息
func (s *PostService) UnlikePost(userID, postID uint) error {
	// 删除点赞记录，触发AfterDelete钩子更新统计
	return s.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Like{}).Error
}

// ==================== 评论管理服务 ====================

// CommentService 评论管理服务
// 提供评论相关的所有业务操作，包括创建、审核、垃圾评论处理等
// 封装了评论相关的复杂业务逻辑和数据库操作
type CommentService struct {
	db *gorm.DB // 数据库连接实例
}

// NewCommentService 创建新的评论服务实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *CommentService: 评论服务实例
func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

// CreateComment 创建评论
// 使用事务确保数据一致性，支持多层级回复评论
// 包含文章评论权限检查和评论层级自动计算
// 参数:
//   - comment: 评论对象指针，包含用户ID、文章ID、内容等信息
//
// 返回:
//   - error: 创建失败时返回错误信息
func (s *CommentService) CreateComment(comment *Comment) error {
	// 使用数据库事务确保操作的原子性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查目标文章是否存在且允许评论
		var post Post
		if err := tx.First(&post, comment.PostID).Error; err != nil {
			return err
		}

		// 验证文章的评论权限设置
		if !post.AllowComments {
			return fmt.Errorf("该文章不允许评论")
		}

		// 如果是回复评论，需要设置正确的层级关系
		if comment.ParentID != nil {
			var parentComment Comment
			if err := tx.First(&parentComment, *comment.ParentID).Error; err != nil {
				return err
			}
			// 计算评论层级：父评论层级 + 1
			comment.Level = parentComment.Level + 1
		}

		// 创建评论记录，触发AfterCreate钩子更新统计和发送通知
		return tx.Create(comment).Error
	})
}

// ApproveComment 审核通过评论
// 将评论状态从待审核改为已通过，使评论对外可见
// 参数:
//   - commentID: 评论ID
//
// 返回:
//   - error: 审核失败时返回错误信息
func (s *CommentService) ApproveComment(commentID uint) error {
	// 更新评论状态为已通过
	return s.db.Model(&Comment{}).Where("id = ?", commentID).Update("status", "approved").Error
}

// RejectComment 拒绝评论
// 将评论状态从待审核改为已拒绝，评论将不会对外显示
// 参数:
//   - commentID: 评论ID
//
// 返回:
//   - error: 拒绝失败时返回错误信息
func (s *CommentService) RejectComment(commentID uint) error {
	// 更新评论状态为已拒绝
	return s.db.Model(&Comment{}).Where("id = ?", commentID).Update("status", "rejected").Error
}

// MarkAsSpam 标记评论为垃圾评论
// 将评论标记为垃圾评论，同时更新状态和垃圾评论标志
// 用于处理恶意、广告或无意义的评论内容
// 参数:
//   - commentID: 评论ID
//
// 返回:
//   - error: 标记失败时返回错误信息
func (s *CommentService) MarkAsSpam(commentID uint) error {
	// 同时更新评论状态和垃圾评论标志
	return s.db.Model(&Comment{}).Where("id = ?", commentID).Updates(map[string]interface{}{
		"status":  "spam", // 设置状态为垃圾评论
		"is_spam": true,   // 标记为垃圾评论
	}).Error
}

// ==================== 通知管理服务 ====================

// NotificationService 通知管理服务
// 提供用户通知相关的所有业务操作，包括获取、标记已读等
// 处理系统内各种事件产生的通知消息
type NotificationService struct {
	db *gorm.DB // 数据库连接实例
}

// NewNotificationService 创建新的通知服务实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *NotificationService: 通知服务实例
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// GetUserNotifications 获取用户通知列表
// 支持分页查询，按创建时间倒序排列
// 参数:
//   - userID: 用户ID
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//
// 返回:
//   - []Notification: 通知列表
//   - int64: 用户通知总数
//   - error: 查询失败时返回错误信息
func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 获取用户通知总数
	s.db.Model(&Notification{}).Where("user_id = ?", userID).Count(&total)

	// 获取分页的通知数据，按创建时间倒序排列
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&notifications).Error

	return notifications, total, err
}

// MarkAsRead 标记单个通知为已读
// 更新通知的已读状态和已读时间
// 参数:
//   - notificationID: 通知ID
//
// 返回:
//   - error: 标记失败时返回错误信息
func (s *NotificationService) MarkAsRead(notificationID uint) error {
	now := time.Now()
	// 更新通知的已读状态和已读时间
	return s.db.Model(&Notification{}).Where("id = ?", notificationID).Updates(map[string]interface{}{
		"is_read": true, // 标记为已读
		"read_at": now,  // 记录已读时间
	}).Error
}

// MarkAllAsRead 标记用户所有未读通知为已读
// 批量更新用户的所有未读通知状态
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - error: 标记失败时返回错误信息
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	now := time.Now()
	// 批量更新用户所有未读通知的状态和已读时间
	return s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Updates(map[string]interface{}{
		"is_read": true, // 标记为已读
		"read_at": now,  // 记录已读时间
	}).Error
}

// GetUnreadCount 获取用户未读通知数量
// 统计用户当前的未读通知总数，用于显示通知徽章
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - int64: 未读通知数量
//   - error: 查询失败时返回错误信息
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	// 统计用户未读通知数量
	err := s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// ==================== 统计分析服务 ====================

// AnalyticsService 统计分析服务
// 提供系统各项数据的统计分析功能，包括用户、文章、评论等维度的数据统计
// 用于生成管理后台的数据报表和仪表板
type AnalyticsService struct {
	db *gorm.DB // 数据库连接实例
}

// NewAnalyticsService 创建新的统计分析服务实例
// 参数:
//   - db: GORM数据库连接实例
//
// 返回:
//   - *AnalyticsService: 统计分析服务实例
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// DashboardStats 仪表板统计数据结构
// 包含系统各项核心指标的统计信息，用于管理后台仪表板展示
type DashboardStats struct {
	TotalUsers       int64 `json:"total_users"`       // 用户总数
	ActiveUsers      int64 `json:"active_users"`      // 活跃用户数
	TotalPosts       int64 `json:"total_posts"`       // 文章总数
	PublishedPosts   int64 `json:"published_posts"`   // 已发布文章数
	TotalComments    int64 `json:"total_comments"`    // 评论总数
	ApprovedComments int64 `json:"approved_comments"` // 已审核评论数
	TotalViews       int64 `json:"total_views"`       // 总浏览量
	TotalLikes       int64 `json:"total_likes"`       // 总点赞数
}

// GetDashboardStats 获取仪表板统计数据
// 统计系统各项核心指标，包括用户、文章、评论、浏览量、点赞数等
// 返回:
//   - *DashboardStats: 仪表板统计数据
//   - error: 统计失败时返回错误信息
func (s *AnalyticsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 用户相关统计
	s.db.Model(&User{}).Count(&stats.TotalUsers)                                // 用户总数
	s.db.Model(&User{}).Where("status = ?", "active").Count(&stats.ActiveUsers) // 活跃用户数

	// 文章相关统计
	s.db.Model(&Post{}).Count(&stats.TotalPosts)                                      // 文章总数
	s.db.Model(&Post{}).Where("status = ?", "published").Count(&stats.PublishedPosts) // 已发布文章数

	// 评论相关统计
	s.db.Model(&Comment{}).Count(&stats.TotalComments)                                    // 评论总数
	s.db.Model(&Comment{}).Where("status = ?", "approved").Count(&stats.ApprovedComments) // 已审核评论数

	// 浏览量和点赞统计（使用COALESCE处理NULL值）
	s.db.Model(&Post{}).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews) // 总浏览量
	s.db.Model(&Post{}).Select("COALESCE(SUM(like_count), 0)").Scan(&stats.TotalLikes) // 总点赞数

	return stats, nil
}

// PopularPost 热门文章数据结构
// 用于展示热门文章的关键信息，包括浏览量、点赞数等指标
type PopularPost struct {
	ID        uint   `json:"id"`         // 文章ID
	Title     string `json:"title"`      // 文章标题
	Slug      string `json:"slug"`       // 文章别名
	ViewCount int    `json:"view_count"` // 浏览量
	LikeCount int    `json:"like_count"` // 点赞数
	Author    string `json:"author"`     // 作者用户名
}

// GetPopularPosts 获取热门文章列表
// 根据浏览量和点赞数排序，返回最受欢迎的文章
// 参数:
//   - limit: 返回文章数量限制
//
// 返回:
//   - []PopularPost: 热门文章列表
//   - error: 查询失败时返回错误信息
func (s *AnalyticsService) GetPopularPosts(limit int) ([]PopularPost, error) {
	var posts []PopularPost

	// 联表查询文章和作者信息，按浏览量和点赞数排序
	err := s.db.Table("posts p").
		Select("p.id, p.title, p.slug, p.view_count, p.like_count, u.username as author").
		Joins("JOIN users u ON p.author_id = u.id").   // 关联用户表获取作者信息
		Where("p.status = ?", "published").            // 只查询已发布的文章
		Order("p.view_count DESC, p.like_count DESC"). // 按浏览量和点赞数倒序排列
		Limit(limit).Scan(&posts).Error

	return posts, err
}

// ActiveUser 活跃用户数据结构
// 用于展示活跃用户的关键指标，包括发文数、评论数、点赞数等
type ActiveUser struct {
	ID           uint   `json:"id"`            // 用户ID
	Username     string `json:"username"`      // 用户名
	PostCount    int    `json:"post_count"`    // 发文数量
	CommentCount int    `json:"comment_count"` // 评论数量
	LikeCount    int    `json:"like_count"`    // 点赞数量
}

// GetActiveUsers 获取活跃用户列表
// 根据用户的发文数、评论数、点赞数综合评分排序
// 参数:
//   - limit: 返回用户数量限制
//
// 返回:
//   - []ActiveUser: 活跃用户列表
//   - error: 查询失败时返回错误信息
func (s *AnalyticsService) GetActiveUsers(limit int) ([]ActiveUser, error) {
	var users []ActiveUser

	// 复杂查询：联表统计用户活跃度指标
	err := s.db.Table("users u").
		Select("u.id, u.username, u.post_count, u.comment_count, COALESCE(l.like_count, 0) as like_count").
		// 左连接点赞统计子查询
		Joins("LEFT JOIN (SELECT user_id, COUNT(*) as like_count FROM likes GROUP BY user_id) l ON u.id = l.user_id").
		Where("u.status = ?", "active").                                            // 只查询活跃用户
		Order("(u.post_count + u.comment_count + COALESCE(l.like_count, 0)) DESC"). // 按综合活跃度排序
		Limit(limit).Scan(&users).Error

	return users, err
}

// ==================== 测试数据生成 ====================

// generateComprehensiveTestData 生成综合测试数据
// 创建完整的测试数据集，包括用户、分类、标签、文章、评论、点赞、关注等
// 用于演示和测试系统的各项功能
// 参数:
//   - db: GORM数据库连接实例
func generateComprehensiveTestData(db *gorm.DB) {
	fmt.Println("开始生成综合测试数据...")
	start := time.Now()

	// ==================== 生成用户数据 ====================
	// 创建50个测试用户，包含完整的用户信息
	users := make([]User, 50)
	for i := 0; i < 50; i++ {
		users[i] = User{
			Username:      fmt.Sprintf("user%d", i+1),                            // 用户名
			Email:         fmt.Sprintf("user%d@example.com", i+1),                // 邮箱
			PasswordHash:  "hashed_password",                                     // 密码哈希
			FirstName:     fmt.Sprintf("First%d", i+1),                           // 名
			LastName:      fmt.Sprintf("Last%d", i+1),                            // 姓
			Bio:           fmt.Sprintf("这是用户%d的个人简介", i+1),                       // 个人简介
			Location:      []string{"北京", "上海", "深圳", "广州", "杭州"}[rand.Intn(5)],  // 随机地理位置
			Gender:        []string{"male", "female", "other"}[rand.Intn(3)],     // 随机性别
			Status:        "active",                                              // 用户状态
			Role:          []string{"user", "author", "moderator"}[rand.Intn(3)], // 随机角色
			EmailVerified: rand.Float32() > 0.2,                                  // 80%概率邮箱已验证
			LoginCount:    rand.Intn(100),                                        // 随机登录次数
		}
		// 70%概率设置最近登录时间
		if rand.Float32() > 0.3 {
			lastLogin := time.Now().AddDate(0, 0, -rand.Intn(30)) // 最近30天内的随机时间
			users[i].LastLoginAt = &lastLogin
		}
	}
	db.Create(&users)
	fmt.Printf("✓ 用户数据: %d条\n", len(users))

	// ==================== 生成分类数据 ====================
	// 创建5个主要分类，每个分类包含完整的元数据
	categories := []Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", Icon: "code", Color: "#007bff", Level: 1, SortOrder: 1, IsActive: true},
		{Name: "生活", Slug: "life", Description: "生活分享", Icon: "heart", Color: "#28a745", Level: 1, SortOrder: 2, IsActive: true},
		{Name: "旅游", Slug: "travel", Description: "旅游攻略", Icon: "map", Color: "#ffc107", Level: 1, SortOrder: 3, IsActive: true},
		{Name: "美食", Slug: "food", Description: "美食推荐", Icon: "utensils", Color: "#fd7e14", Level: 1, SortOrder: 4, IsActive: true},
		{Name: "娱乐", Slug: "entertainment", Description: "娱乐资讯", Icon: "film", Color: "#e83e8c", Level: 1, SortOrder: 5, IsActive: true},
	}
	db.Create(&categories)
	fmt.Printf("✓ 分类数据: %d条\n", len(categories))

	// ==================== 生成标签数据 ====================
	// 创建10个常用标签，涵盖技术和内容类型
	tags := []Tag{
		{Name: "Go", Slug: "go", Description: "Go语言相关", Color: "#00ADD8", IsActive: true},
		{Name: "Python", Slug: "python", Description: "Python编程", Color: "#3776AB", IsActive: true},
		{Name: "JavaScript", Slug: "javascript", Description: "JavaScript开发", Color: "#F7DF1E", IsActive: true},
		{Name: "数据库", Slug: "database", Description: "数据库技术", Color: "#336791", IsActive: true},
		{Name: "前端", Slug: "frontend", Description: "前端开发", Color: "#61DAFB", IsActive: true},
		{Name: "后端", Slug: "backend", Description: "后端开发", Color: "#68217A", IsActive: true},
		{Name: "教程", Slug: "tutorial", Description: "教程文章", Color: "#FF6B6B", IsActive: true},
		{Name: "实战", Slug: "practice", Description: "实战项目", Color: "#4ECDC4", IsActive: true},
		{Name: "经验分享", Slug: "experience", Description: "经验分享", Color: "#95A5A6", IsActive: true},
		{Name: "工具推荐", Slug: "tools", Description: "工具推荐", Color: "#9B59B6", IsActive: true},
	}
	db.Create(&tags)
	fmt.Printf("✓ 标签数据: %d条\n", len(tags))

	// ==================== 生成文章数据 ====================
	// 创建200篇测试文章，包含各种状态和类型
	posts := make([]Post, 200)
	statuses := []string{"published", "draft", "archived"}       // 文章状态
	types := []string{"post", "page"}                            // 文章类型
	formats := []string{"standard", "video", "gallery", "quote"} // 文章格式

	for i := 0; i < 200; i++ {
		authorID := uint(rand.Intn(50) + 1)          // 随机作者ID
		categoryID := uint(rand.Intn(5) + 1)         // 随机分类ID
		status := statuses[rand.Intn(len(statuses))] // 随机状态

		posts[i] = Post{
			Title:           fmt.Sprintf("精彩文章标题 %d - 深度解析技术要点", i+1),
			Slug:            fmt.Sprintf("awesome-post-%d", i+1),
			Content:         fmt.Sprintf("这是文章 %d 的详细内容，包含了丰富的技术信息和实用的代码示例。文章深入浅出地讲解了相关概念，并提供了实际的应用场景和最佳实践。", i+1),
			Excerpt:         fmt.Sprintf("文章 %d 的精彩摘要，概括了主要内容和核心观点", i+1),
			FeaturedImage:   fmt.Sprintf("/images/post-%d.jpg", i+1),
			Status:          status,
			Type:            types[rand.Intn(len(types))],      // 随机文章类型
			Format:          formats[rand.Intn(len(formats))],  // 随机文章格式
			ViewCount:       rand.Intn(5000),                   // 随机浏览量
			LikeCount:       rand.Intn(500),                    // 随机点赞数
			CommentCount:    rand.Intn(50),                     // 随机评论数
			ShareCount:      rand.Intn(100),                    // 随机分享数
			Rating:          float64(rand.Intn(40))/10.0 + 1.0, // 1.0-5.0的随机评分
			Featured:        rand.Float32() > 0.8,              // 20%概率为精选文章
			Sticky:          rand.Float32() > 0.95,             // 5%概率为置顶文章
			AllowComments:   rand.Float32() > 0.1,              // 90%概率允许评论
			MetaTitle:       fmt.Sprintf("文章%d的SEO标题", i+1),
			MetaDescription: fmt.Sprintf("文章%d的SEO描述", i+1),
			MetaKeywords:    fmt.Sprintf("关键词%d,技术,教程", i+1),
			AuthorID:        authorID,
			CategoryID:      &categoryID,
		}

		// 已发布文章设置发布时间
		if status == "published" {
			publishedAt := time.Now().AddDate(0, 0, -rand.Intn(365)) // 过去一年内的随机时间
			posts[i].PublishedAt = &publishedAt
		}
	}
	db.CreateInBatches(posts, 50) // 批量创建，每批50条
	fmt.Printf("✓ 文章数据: %d条\n", len(posts))

	// ==================== 为文章分配标签 ====================
	// 为每篇文章随机分配1-4个标签，避免重复
	for i := 1; i <= 200; i++ {
		tagCount := rand.Intn(4) + 1 // 每篇文章1-4个标签
		selectedTags := make([]Tag, 0, tagCount)
		usedTagIDs := make(map[uint]bool) // 防止重复标签

		for j := 0; j < tagCount; j++ {
			tagID := uint(rand.Intn(10) + 1)
			if !usedTagIDs[tagID] {
				var tag Tag
				db.First(&tag, tagID)
				selectedTags = append(selectedTags, tag)
				usedTagIDs[tagID] = true
			}
		}

		// 建立文章与标签的多对多关联
		var post Post
		db.First(&post, uint(i))
		db.Model(&post).Association("Tags").Append(selectedTags)
	}
	fmt.Println("✓ 文章标签关联完成")

	// ==================== 生成评论数据 ====================
	// 创建500条评论，包含各种状态和完整的元数据
	comments := make([]Comment, 500)
	commentStatuses := []string{"approved", "pending", "spam"} // 评论状态

	for i := 0; i < 500; i++ {
		comments[i] = Comment{
			Content:   fmt.Sprintf("这是评论 %d 的内容，包含了用户的真实想法和建议。", i+1),
			Status:    commentStatuses[rand.Intn(len(commentStatuses))],               // 随机评论状态
			Type:      "comment",                                                      // 评论类型
			LikeCount: rand.Intn(50),                                                  // 随机点赞数
			Level:     1,                                                              // 一级评论
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", // 用户代理
			UserIP:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)),                    // 随机IP地址
			IsSpam:    rand.Float32() > 0.9,                                           // 10%概率为垃圾评论
			PostID:    uint(rand.Intn(200) + 1),                                       // 随机文章ID
			AuthorID:  uint(rand.Intn(50) + 1),                                        // 随机作者ID
		}
	}
	db.CreateInBatches(comments, 100) // 批量创建，每批100条
	fmt.Printf("✓ 评论数据: %d条\n", len(comments))

	// ==================== 生成回复评论 ====================
	// 创建200条回复评论，形成评论层级结构
	replies := make([]Comment, 200)
	for i := 0; i < 200; i++ {
		parentID := uint(rand.Intn(500) + 1) // 随机父评论ID
		replies[i] = Comment{
			Content:   fmt.Sprintf("这是对评论的回复 %d", i+1),
			Status:    "approved",                                                     // 默认已审核
			Type:      "reply",                                                        // 回复类型
			LikeCount: rand.Intn(20),                                                  // 随机点赞数
			ParentID:  &parentID,                                                      // 设置父评论ID
			Level:     2,                                                              // 二级评论
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", // 用户代理
			UserIP:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)),                    // 随机IP地址
			PostID:    uint(rand.Intn(200) + 1),                                       // 随机文章ID
			AuthorID:  uint(rand.Intn(50) + 1),                                        // 随机作者ID
		}
	}
	db.CreateInBatches(replies, 50) // 批量创建，每批50条
	fmt.Printf("✓ 回复数据: %d条\n", len(replies))

	// ==================== 生成点赞数据 ====================
	// 创建1000条点赞记录，包含文章点赞和评论点赞
	likes := make([]Like, 1000)
	for i := 0; i < 1000; i++ {
		userID := uint(rand.Intn(50) + 1) // 随机用户ID

		if rand.Float32() > 0.3 { // 70%概率点赞文章
			postID := uint(rand.Intn(200) + 1)
			likes[i] = Like{
				UserID: userID,
				PostID: &postID,                                     // 文章ID
				Type:   "like",                                      // 点赞类型
				UserIP: fmt.Sprintf("192.168.1.%d", rand.Intn(255)), // 随机IP地址
			}
		} else { // 30%概率点赞评论
			commentID := uint(rand.Intn(700) + 1)
			likes[i] = Like{
				UserID:    userID,
				CommentID: &commentID,                                  // 评论ID
				Type:      "like",                                      // 点赞类型
				UserIP:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)), // 随机IP地址
			}
		}
	}
	db.CreateInBatches(likes, 100) // 批量创建，每批100条
	fmt.Printf("✓ 点赞数据: %d条\n", len(likes))

	// ==================== 生成关注关系 ====================
	// 创建300条关注关系，建立用户之间的社交网络
	follows := make([]Follow, 300)
	for i := 0; i < 300; i++ {
		followerID := uint(rand.Intn(50) + 1)  // 关注者ID
		followingID := uint(rand.Intn(50) + 1) // 被关注者ID

		// 确保不自己关注自己的逻辑检查
		for followerID == followingID {
			followingID = uint(rand.Intn(50) + 1)
		}

		follows[i] = Follow{
			FollowerID:  followerID,
			FollowingID: followingID,
			Status:      "active", // 关注状态为活跃
		}
	}
	db.CreateInBatches(follows, 50) // 批量创建，每批50条
	fmt.Printf("✓ 关注关系: %d条\n", len(follows))

	// ==================== 生成通知数据 ====================
	// 创建200条通知消息，包含各种类型的系统通知
	notifications := make([]Notification, 200)
	notificationTypes := []string{"comment", "like", "follow", "mention", "system"} // 通知类型

	for i := 0; i < 200; i++ {
		notificationType := notificationTypes[rand.Intn(len(notificationTypes))] // 随机通知类型
		notifications[i] = Notification{
			UserID:      uint(rand.Intn(50) + 1), // 随机用户ID
			Type:        notificationType,
			Title:       fmt.Sprintf("%s通知 %d", notificationType, i+1),
			Content:     fmt.Sprintf("这是一条%s类型的通知内容", notificationType),
			Data:        fmt.Sprintf(`{"id": %d, "type": "%s"}`, i+1, notificationType), // JSON格式的额外数据
			IsRead:      rand.Float32() > 0.4,                                           // 60%概率已读
			RelatedID:   func() *uint { id := uint(rand.Intn(100) + 1); return &id }(),  // 关联对象ID
			RelatedType: []string{"post", "comment", "user"}[rand.Intn(3)],              // 随机关联类型
		}

		// 已读通知设置读取时间
		if notifications[i].IsRead {
			readAt := time.Now().AddDate(0, 0, -rand.Intn(7)) // 最近7天内的随机时间
			notifications[i].ReadAt = &readAt
		}
	}
	db.CreateInBatches(notifications, 50) // 批量创建，每批50条
	fmt.Printf("✓ 通知数据: %d条\n", len(notifications))

	// ==================== 生成系统设置 ====================
	// 创建8条系统配置项，涵盖网站的各种设置
	settings := []Setting{
		{Key: "site_title", Value: "我的博客系统", Type: "string", Description: "网站标题", Group: "general", IsPublic: true},
		{Key: "site_description", Value: "一个功能完整的博客系统", Type: "string", Description: "网站描述", Group: "general", IsPublic: true},
		{Key: "posts_per_page", Value: "10", Type: "integer", Description: "每页文章数量", Group: "display", IsPublic: true},
		{Key: "allow_registration", Value: "true", Type: "boolean", Description: "允许用户注册", Group: "user", IsPublic: false},
		{Key: "comment_moderation", Value: "true", Type: "boolean", Description: "评论需要审核", Group: "comment", IsPublic: false},
		{Key: "email_notifications", Value: "true", Type: "boolean", Description: "邮件通知", Group: "notification", IsPublic: false},
		{Key: "max_upload_size", Value: "10485760", Type: "integer", Description: "最大上传文件大小(字节)", Group: "upload", IsPublic: false},
		{Key: "allowed_file_types", Value: "jpg,jpeg,png,gif,pdf,doc,docx", Type: "string", Description: "允许的文件类型", Group: "upload", IsPublic: false},
	}
	db.Create(&settings)
	fmt.Printf("✓ 系统设置: %d条\n", len(settings))

	// ==================== 数据生成完成 ====================
	elapsed := time.Since(start)
	fmt.Printf("\n🎉 综合测试数据生成完成，耗时: %v\n", elapsed)
}

// ==================== 综合业务场景演示 ====================
// demonstrateComprehensiveScenarios 演示完整的博客系统业务场景
// 包括用户注册、文章发布、用户互动、内容搜索、通知管理和数据统计等功能
// 参数:
//   - db: GORM数据库连接实例
func demonstrateComprehensiveScenarios(db *gorm.DB) {
	fmt.Println("\n=== 综合业务场景演示 ===")

	// 初始化各种业务服务层
	userService := NewUserService(db)                 // 用户管理服务
	postService := NewPostService(db)                 // 文章管理服务
	commentService := NewCommentService(db)           // 评论管理服务
	notificationService := NewNotificationService(db) // 通知管理服务
	analyticsService := NewAnalyticsService(db)       // 数据分析服务

	// ==================== 场景1：用户注册和资料完善 ====================
	// 演示新用户注册流程和个人资料管理
	fmt.Println("\n--- 场景1：用户注册和资料完善 ---")
	newUser := &User{
		Username:     "newuser",             // 用户名
		Email:        "newuser@example.com", // 邮箱地址
		PasswordHash: "hashed_password",     // 加密后的密码
		FirstName:    "New",                 // 名字
		LastName:     "User",                // 姓氏
		Bio:          "我是新用户",               // 个人简介
		Location:     "北京",                  // 所在地
		Status:       "active",              // 用户状态
		Role:         "user",                // 用户角色
	}

	// 创建新用户
	if err := userService.CreateUser(newUser); err != nil {
		fmt.Printf("用户创建失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户创建成功，ID: %d\n", newUser.ID)

		// 完善用户详细资料信息
		profile := &UserProfile{
			Company:      "科技公司",                   // 公司名称
			JobTitle:     "软件工程师",                  // 职位
			Education:    "计算机科学学士",                // 教育背景
			Skills:       "Go, Python, JavaScript", // 技能
			Experience:   3,                        // 工作经验年数
			SalaryRange:  "10k-15k",                // 薪资范围
			Languages:    "中文, 英文",                 // 语言能力
			Interests:    "编程, 阅读, 旅游",             // 兴趣爱好
			PrivacyLevel: "public",                 // 隐私级别
		}

		// 更新用户资料
		if err := userService.UpdateUserProfile(newUser.ID, profile); err != nil {
			fmt.Printf("资料更新失败: %v\n", err)
		} else {
			fmt.Println("✓ 用户资料更新成功")
		}
	}

	// ==================== 场景2：发布文章和标签管理 ====================
	// 演示文章创建、标签关联和内容管理功能
	fmt.Println("\n--- 场景2：发布文章和标签管理 ---")
	// 创建新文章对象
	newPost := &Post{
		Title:           "GORM实战教程 - 从入门到精通",                          // 文章标题
		Slug:            "gorm-tutorial-comprehensive",                // URL友好的标识符
		Content:         "这是一篇关于GORM的详细教程，涵盖了从基础使用到高级特性的所有内容...",      // 文章正文内容
		Excerpt:         "GORM实战教程，带你深入了解Go语言最流行的ORM框架",               // 文章摘要
		FeaturedImage:   "/images/gorm-tutorial.jpg",                  // 特色图片
		Status:          "published",                                  // 发布状态
		Type:            "post",                                       // 内容类型
		Format:          "standard",                                   // 文章格式
		Featured:        true,                                         // 是否为精选文章
		AllowComments:   true,                                         // 是否允许评论
		MetaTitle:       "GORM实战教程 - 完整指南",                            // SEO标题
		MetaDescription: "学习GORM的最佳实践和高级技巧",                           // SEO描述
		MetaKeywords:    "GORM,Go,ORM,数据库,教程",                         // SEO关键词
		AuthorID:        1,                                            // 作者ID
		CategoryID:      func() *uint { id := uint(1); return &id }(), // 分类ID
	}

	// 为文章分配相关标签
	var tags []Tag
	db.Where("slug IN ?", []string{"go", "database", "tutorial"}).Find(&tags) // 查找相关标签
	newPost.Tags = tags                                                       // 建立文章与标签的多对多关联

	// 创建文章
	if err := postService.CreatePost(newPost); err != nil {
		fmt.Printf("文章创建失败: %v\n", err)
	} else {
		fmt.Printf("✓ 文章创建成功，ID: %d\n", newPost.ID)
	}

	// ==================== 场景3：用户互动（评论、点赞、关注） ====================
	// 演示用户之间的社交互动功能
	fmt.Println("\n--- 场景3：用户互动 ---")

	// 创建新评论
	newComment := &Comment{
		Content:   "这篇文章写得非常好，对我帮助很大！",                                            // 评论内容
		Status:    "approved",                                                     // 审核状态
		Type:      "comment",                                                      // 评论类型
		Level:     1,                                                              // 评论层级
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", // 用户代理
		UserIP:    "192.168.1.100",                                                // 用户IP
		PostID:    1,                                                              // 文章ID
		AuthorID:  2,                                                              // 评论作者ID
	}

	// 提交评论
	if err := commentService.CreateComment(newComment); err != nil {
		fmt.Printf("评论创建失败: %v\n", err)
	} else {
		fmt.Printf("✓ 评论创建成功，ID: %d\n", newComment.ID)
	}

	// 用户点赞文章（用户ID=2，文章ID=1）
	if err := postService.LikePost(2, 1); err != nil {
		fmt.Printf("点赞失败: %v\n", err)
	} else {
		fmt.Println("✓ 文章点赞成功")
	}

	// 用户关注功能（用户2关注用户1）
	if err := userService.FollowUser(2, 1); err != nil {
		fmt.Printf("关注失败: %v\n", err)
	} else {
		fmt.Println("✓ 用户关注成功")
	}

	// ==================== 场景4：内容搜索和分页 ====================
	// 演示各种内容查询和分页功能
	fmt.Println("\n--- 场景4：内容搜索和分页 ---")

	// 按分类获取文章（分页查询）
	posts, total, err := postService.GetPostsByCategory("tech", 1, 5) // 技术分类，第1页，每页5篇
	if err != nil {
		fmt.Printf("分类查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 技术分类文章: %d/%d篇\n", len(posts), total)
	}

	// 全文搜索文章
	searchPosts, searchTotal, err := postService.SearchPosts("教程", 1, 5) // 搜索关键词，第1页，每页5篇
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
	} else {
		fmt.Printf("✓ 搜索结果: %d/%d篇\n", len(searchPosts), searchTotal)
	}

	// 获取用户关注者列表
	followers, followerTotal, err := userService.GetUserFollowers(1, 1, 10) // 用户1的关注者，第1页，每页10人
	if err != nil {
		fmt.Printf("关注者查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户关注者: %d/%d人\n", len(followers), followerTotal)
	}

	// ==================== 场景5：通知管理 ====================
	// 演示用户通知系统的各种功能
	fmt.Println("\n--- 场景5：通知管理 ---")

	// 获取用户通知列表（分页查询）
	notifications, notificationTotal, err := notificationService.GetUserNotifications(1, 1, 5) // 用户1的通知，第1页，每页5条
	if err != nil {
		fmt.Printf("通知查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户通知: %d/%d条\n", len(notifications), notificationTotal)

		// 标记第一条通知为已读
		if len(notifications) > 0 {
			if err := notificationService.MarkAsRead(notifications[0].ID); err != nil {
				fmt.Printf("标记已读失败: %v\n", err)
			} else {
				fmt.Println("✓ 通知标记为已读")
			}
		}
	}

	// 获取用户未读通知数量
	unreadCount, err := notificationService.GetUnreadCount(1) // 查询用户1的未读通知数
	if err != nil {
		fmt.Printf("未读数量查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 未读通知数量: %d条\n", unreadCount)
	}

	// ==================== 场景6：数据统计和分析 ====================
	// 演示系统数据分析和统计功能
	fmt.Println("\n--- 场景6：数据统计和分析 ---")

	// 获取系统仪表板统计数据
	stats, err := analyticsService.GetDashboardStats()
	if err != nil {
		fmt.Printf("统计查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 系统统计:\n")
		fmt.Printf("  - 总用户数: %d (活跃: %d)\n", stats.TotalUsers, stats.ActiveUsers)          // 用户统计
		fmt.Printf("  - 总文章数: %d (已发布: %d)\n", stats.TotalPosts, stats.PublishedPosts)      // 文章统计
		fmt.Printf("  - 总评论数: %d (已审核: %d)\n", stats.TotalComments, stats.ApprovedComments) // 评论统计
		fmt.Printf("  - 总浏览量: %d\n", stats.TotalViews)                                      // 浏览量统计
		fmt.Printf("  - 总点赞数: %d\n", stats.TotalLikes)                                      // 点赞统计
	}

	// 获取热门文章排行榜
	popularPosts, err := analyticsService.GetPopularPosts(5) // 获取前5篇热门文章
	if err != nil {
		fmt.Printf("热门文章查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 热门文章 (前5篇):\n")
		for i, post := range popularPosts {
			fmt.Printf("  %d. %s (浏览: %d, 点赞: %d)\n", i+1, post.Title, post.ViewCount, post.LikeCount)
		}
	}

	// 获取活跃用户排行榜
	activeUsers, err := analyticsService.GetActiveUsers(5) // 获取前5名活跃用户
	if err != nil {
		fmt.Printf("活跃用户查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 活跃用户 (前5名):\n")
		for i, user := range activeUsers {
			fmt.Printf("  %d. %s (文章: %d, 评论: %d, 点赞: %d)\n",
				i+1, user.Username, user.PostCount, user.CommentCount, user.LikeCount)
		}
	}
}

// demonstrateAdvancedQueries 演示高级查询功能
// 包括复杂多表连接、窗口函数、时间序列分析等高级SQL特性
// 参数:
//   - db: GORM数据库实例
func demonstrateAdvancedQueries(db *gorm.DB) {
	fmt.Println("\n=== 高级查询演示 ===")

	// ==================== 复杂的多表连接查询 ====================
	// 演示JOIN查询、GROUP BY聚合、子查询等复杂SQL操作
	fmt.Println("\n--- 复杂多表连接查询 ---")

	// PostWithStats 文章统计信息结构体
	// 包含文章基本信息和相关统计数据
	type PostWithStats struct {
		ID           uint   `json:"id"`            // 文章ID
		Title        string `json:"title"`         // 文章标题
		AuthorName   string `json:"author_name"`   // 作者姓名
		CategoryName string `json:"category_name"` // 分类名称
		ViewCount    int    `json:"view_count"`    // 浏览次数
		LikeCount    int    `json:"like_count"`    // 点赞次数
		CommentCount int    `json:"comment_count"` // 评论次数
		TagNames     string `json:"tag_names"`     // 标签名称（逗号分隔）
	}

	var postsWithStats []PostWithStats
	// 执行复杂的多表连接查询
	// 连接用户表、分类表、标签表，获取文章的完整统计信息
	err := db.Table("posts p").
		Select(`p.id, p.title, u.username as author_name, c.name as category_name, 
			p.view_count, p.like_count, p.comment_count,
			GROUP_CONCAT(t.name) as tag_names`). // 选择字段，使用GROUP_CONCAT聚合标签
		Joins("JOIN users u ON p.author_id = u.id").                                             // 内连接用户表获取作者信息
		Joins("LEFT JOIN categories c ON p.category_id = c.id").                                 // 左连接分类表
		Joins("LEFT JOIN post_tags pt ON p.id = pt.post_id").                                    // 左连接文章标签关联表
		Joins("LEFT JOIN tags t ON pt.tag_id = t.id").                                           // 左连接标签表
		Where("p.status = ?", "published").                                                      // 只查询已发布的文章
		Group("p.id, p.title, u.username, c.name, p.view_count, p.like_count, p.comment_count"). // 按文章分组
		Order("p.view_count DESC").                                                              // 按浏览量降序排列
		Limit(5).Scan(&postsWithStats).Error                                                     // 限制返回5条记录

	// 处理查询结果并显示
	if err != nil {
		fmt.Printf("复杂查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 文章详细统计 (前5篇):\n")
		for i, post := range postsWithStats {
			fmt.Printf("  %d. %s\n", i+1, post.Title)
			fmt.Printf("     作者: %s | 分类: %s\n", post.AuthorName, post.CategoryName)
			fmt.Printf("     浏览: %d | 点赞: %d | 评论: %d\n", post.ViewCount, post.LikeCount, post.CommentCount)
			fmt.Printf("     标签: %s\n", post.TagNames)
			fmt.Println()
		}
	}

	// ==================== 窗口函数查询（排名） ====================
	// 演示ROW_NUMBER()等窗口函数的使用
	fmt.Println("--- 窗口函数查询 ---")

	// UserRanking 用户排名结构体
	// 包含用户活跃度统计和排名信息
	type UserRanking struct {
		Username     string `json:"username"`      // 用户名
		PostCount    int    `json:"post_count"`    // 文章数量
		CommentCount int    `json:"comment_count"` // 评论数量
		TotalScore   int    `json:"total_score"`   // 总分（文章*3 + 评论*1）
		Rank         int    `json:"rank"`          // 排名
	}

	var userRankings []UserRanking
	// 使用窗口函数计算用户活跃度排名
	// ROW_NUMBER() OVER() 为每行分配一个唯一的排名
	err = db.Raw(`
		SELECT username, post_count, comment_count, 
		       (post_count * 3 + comment_count) as total_score,  -- 计算总分：文章权重3，评论权重1
		       ROW_NUMBER() OVER (ORDER BY (post_count * 3 + comment_count) DESC) as rank  -- 窗口函数排名
		FROM users 
		WHERE status = 'active'   -- 只统计活跃用户
		ORDER BY total_score DESC -- 按总分降序排列
		LIMIT 10                  -- 限制前10名
	`).Scan(&userRankings).Error

	// 显示用户排名结果
	if err != nil {
		fmt.Printf("排名查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户活跃度排名 (前10名):\n")
		for _, user := range userRankings {
			fmt.Printf("  第%d名: %s (文章: %d, 评论: %d, 总分: %d)\n",
				user.Rank, user.Username, user.PostCount, user.CommentCount, user.TotalScore)
		}
	}

	// ==================== 时间序列分析 ====================
	// 演示时间维度的数据统计和分析
	fmt.Println("\n--- 时间序列分析 ---")

	// MonthlyActivity 月度活动统计结构体
	// 用于统计每月的用户活动情况
	type MonthlyActivity struct {
		Month        string `json:"month"`         // 月份（YYYY-MM格式）
		PostCount    int    `json:"post_count"`    // 文章发布数量
		UserCount    int    `json:"user_count"`    // 新注册用户数量
		CommentCount int    `json:"comment_count"` // 评论数量
	}

	var monthlyActivity []MonthlyActivity
	// 执行时间序列分析查询
	// 使用UNION ALL合并多个表的数据，按月份统计活动情况
	err = db.Raw(`
		SELECT 
		    strftime('%Y-%m', created_at) as month,                    -- 提取年月
		    COUNT(CASE WHEN 'posts' THEN 1 END) as post_count,        -- 统计文章数
		    COUNT(CASE WHEN 'users' THEN 1 END) as user_count,        -- 统计用户数
		    COUNT(CASE WHEN 'comments' THEN 1 END) as comment_count   -- 统计评论数
		FROM (
		    SELECT created_at, 'posts' as type FROM posts WHERE status = 'published'      -- 已发布文章
		    UNION ALL
		    SELECT created_at, 'users' as type FROM users WHERE status = 'active'         -- 活跃用户
		    UNION ALL
		    SELECT created_at, 'comments' as type FROM comments WHERE status = 'approved' -- 已审核评论
		) combined
		WHERE created_at >= date('now', '-6 months')  -- 近6个月的数据
		GROUP BY month                                -- 按月分组
		ORDER BY month DESC                           -- 按月份降序
		LIMIT 6                                       -- 限制6条记录
	`).Scan(&monthlyActivity).Error

	// 显示时间序列分析结果
	if err != nil {
		fmt.Printf("时间序列查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 近6个月活动统计:\n")
		for _, activity := range monthlyActivity {
			fmt.Printf("  %s: 文章 %d篇, 用户 %d人, 评论 %d条\n",
				activity.Month, activity.PostCount, activity.UserCount, activity.CommentCount)
		}
	}
}

// performanceTest 性能测试函数
// 测试各种数据库操作的性能表现，包括批量查询、分页查询、聚合查询和事务操作
// 参数:
//   - db: GORM数据库实例
func performanceTest(db *gorm.DB) {
	fmt.Println("\n=== 性能测试 ===")

	// ==================== 批量查询性能测试 ====================
	// 测试预加载关联数据的查询性能
	fmt.Println("\n--- 批量查询性能测试 ---")
	start := time.Now() // 记录开始时间
	var posts []Post
	// 使用Preload预加载关联数据，避免N+1查询问题
	db.Preload("Author").Preload("Category").Preload("Tags").Limit(100).Find(&posts)
	fmt.Printf("✓ 预加载查询100篇文章: %v\n", time.Since(start))

	// ==================== 分页查询性能测试 ====================
	// 测试分页查询的性能，包括总数统计和数据获取
	start = time.Now()
	var paginatedPosts []Post
	var total int64
	// 先统计总数
	db.Model(&Post{}).Where("status = ?", "published").Count(&total)
	// 再获取分页数据
	db.Where("status = ?", "published").Offset(50).Limit(20).Find(&paginatedPosts)
	fmt.Printf("✓ 分页查询(50-70): %v\n", time.Since(start))

	// ==================== 复杂聚合查询性能测试 ====================
	// 测试包含JOIN、GROUP BY、聚合函数的复杂查询性能
	start = time.Now()

	// CategoryStats 分类统计结构体
	// 用于存储分类的统计信息
	type CategoryStats struct {
		CategoryName string  `json:"category_name"` // 分类名称
		PostCount    int64   `json:"post_count"`    // 文章数量
		AvgViews     float64 `json:"avg_views"`     // 平均浏览量
		TotalLikes   int64   `json:"total_likes"`   // 总点赞数
	}

	var categoryStats []CategoryStats
	// 执行复杂的聚合查询
	db.Table("posts p").
		Select("c.name as category_name, COUNT(p.id) as post_count, AVG(p.view_count) as avg_views, SUM(p.like_count) as total_likes").
		Joins("JOIN categories c ON p.category_id = c.id"). // 连接分类表
		Where("p.status = ?", "published").                 // 只统计已发布文章
		Group("c.id, c.name").                              // 按分类分组
		Order("post_count DESC").                           // 按文章数量降序
		Scan(&categoryStats)
	fmt.Printf("✓ 分类统计查询: %v\n", time.Since(start))

	// ==================== 事务性能测试 ====================
	// 测试数据库事务的执行性能
	start = time.Now()
	// 执行10个事务操作，每个事务包含查询和更新
	for i := 0; i < 10; i++ {
		db.Transaction(func(tx *gorm.DB) error {
			var user User
			// 在事务中查询用户
			tx.First(&user, 1)
			// 在事务中更新登录次数
			tx.Model(&user).Update("login_count", gorm.Expr("login_count + ?", 1))
			return nil // 提交事务
		})
	}
	fmt.Printf("✓ 10个事务操作: %v\n", time.Since(start))
}

// main 主函数
// GORM Level 6 综合实战练习的入口函数
// 按顺序执行数据库初始化、测试数据生成、业务场景演示、高级查询演示和性能测试
// main 主函数 - GORM Level 6 综合实战练习入口
// 提供SQLite和MySQL两种数据库的完整演示
// 包括数据库初始化、测试数据生成、业务场景演示、高级查询和性能测试
func main() {
	fmt.Println("=== GORM Level 6 综合实战练习 ===")
	fmt.Println("本练习将演示GORM的高级特性和综合应用场景")
	fmt.Println("支持SQLite和MySQL两种数据库类型")

	// ==================== 数据库类型选择 ====================
	fmt.Println("\n请选择要使用的数据库类型:")
	fmt.Println("1. SQLite (默认，适合开发和测试)")
	fmt.Println("2. MySQL (适合生产环境)")
	// fmt.Print("请输入选择 (1-2，默认为1): ")

	var choice string
	mysqlDSN := "root:fastbee@tcp(192.168.100.124:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	choice = "2"
	// fmt.Scanln(&choice)

	var config *DatabaseConfig
	switch choice {
	case "2":
		// MySQL配置
		fmt.Println("\n=== MySQL数据库配置 ===")
		config = GetMySQLConfigFromDSN(mysqlDSN)
		// config = configureMySQLDatabase()
	default:
		// SQLite配置（默认）
		fmt.Println("\n=== 使用SQLite数据库 ===")
		config = GetDefaultConfig()
	}

	// ==================== 数据库初始化 ====================
	// 根据配置初始化数据库连接和表结构
	db := initDB(*config)
	fmt.Println("✓ 数据库初始化完成")

	// ==================== 测试数据生成 ====================
	// 生成包含用户、文章、评论、标签等的综合测试数据
	// generateComprehensiveTestData(db)

	// ==================== 综合业务场景演示 ====================
	// 演示用户注册、文章发布、社交互动、内容搜索、通知管理、数据统计等完整业务流程
	demonstrateComprehensiveScenarios(db)

	// ==================== 高级查询演示 ====================
	// 演示复杂多表连接、窗口函数、时间序列分析等高级SQL特性
	demonstrateAdvancedQueries(db)

	// ==================== 性能测试 ====================
	// 测试各种数据库操作的性能表现
	performanceTest(db)

	// ==================== 数据库特性对比演示 ====================
	if config.Type == MySQL {
		// 如果使用MySQL，演示MySQL特有的功能
		demonstrateMySQLFeatures(db)
	}

	// ==================== 练习总结 ====================
	fmt.Println("\n=== Level 6 综合实战练习完成 ===")
	fmt.Printf("\n🎉 恭喜！您已经完成了使用 %s 数据库的GORM综合练习！\n",
		map[DatabaseType]string{SQLite: "SQLite", MySQL: "MySQL"}[config.Type])
	fmt.Println("\n现在您应该能够：")
	fmt.Println("1. 熟练使用GORM进行数据库操作") // 基础CRUD操作
	fmt.Println("2. 设计复杂的数据模型和关联关系")  // 数据建模能力
	fmt.Println("3. 实现高效的查询和事务处理")    // 高级查询技巧
	fmt.Println("4. 优化数据库性能和索引")      // 性能优化能力
	fmt.Println("5. 构建完整的业务应用系统")     // 系统架构能力
	fmt.Println("6. 在不同数据库间进行迁移和适配")  // 数据库适配能力
	fmt.Println("\n继续学习，成为GORM专家！💪")
}

// configureMySQLDatabase 配置MySQL数据库连接
// 通过用户输入获取MySQL连接参数，并返回配置对象
// 返回:
//   - *DatabaseConfig: MySQL数据库配置对象
func configureMySQLDatabase() *DatabaseConfig {
	var host, username, password, database string
	var port int

	// 获取MySQL连接参数
	fmt.Print("请输入MySQL主机地址 (默认: localhost): ")
	fmt.Scanln(&host)
	if host == "" {
		host = "localhost"
	}

	fmt.Print("请输入MySQL端口号 (默认: 3306): ")
	fmt.Scanln(&port)
	if port == 0 {
		port = 3306
	}

	fmt.Print("请输入MySQL用户名 (默认: root): ")
	fmt.Scanln(&username)
	if username == "" {
		username = "root"
	}

	fmt.Print("请输入MySQL密码: ")
	fmt.Scanln(&password)

	fmt.Print("请输入数据库名 (默认: gorm_level6): ")
	fmt.Scanln(&database)
	if database == "" {
		database = "gorm_level6"
	}

	// 返回MySQL配置
	return GetMySQLConfig(host, port, username, password, database)
}

// demonstrateMySQLFeatures 演示MySQL特有功能
// 展示MySQL数据库的特殊功能和优化特性
// 参数:
//   - db: GORM数据库连接实例
func demonstrateMySQLFeatures(db *gorm.DB) {
	fmt.Println("\n=== MySQL特有功能演示 ===")

	// 演示MySQL的JSON字段功能
	demonstrateJSONFields(db)

	// 演示MySQL的全文索引功能
	demonstrateFullTextSearch(db)

	// 演示MySQL的分区表功能
	demonstratePartitioning(db)

	// 演示MySQL的存储引擎特性
	demonstrateStorageEngines(db)
}

// demonstrateJSONFields 演示MySQL的JSON字段功能
// MySQL 5.7+支持原生JSON数据类型，提供高效的JSON存储和查询
// 参数:
//   - db: GORM数据库连接实例
func demonstrateJSONFields(db *gorm.DB) {
	fmt.Println("\n--- MySQL JSON字段演示 ---")

	// 创建包含JSON字段的临时表
	type UserSettings struct {
		ID       uint   `gorm:"primaryKey"`
		UserID   uint   `gorm:"not null;index"`
		Settings string `gorm:"type:json"` // MySQL JSON字段
	}

	// 自动迁移
	db.AutoMigrate(&UserSettings{})

	// 插入JSON数据
	settings := UserSettings{
		UserID:   1,
		Settings: `{"theme": "dark", "language": "zh-CN", "notifications": {"email": true, "push": false}}`,
	}
	db.Create(&settings)

	// 使用JSON函数查询
	var result UserSettings
	db.Where("JSON_EXTRACT(settings, '$.theme') = ?", "dark").First(&result)
	fmt.Printf("查询到主题为dark的用户设置: %+v\n", result)

	fmt.Println("✓ JSON字段演示完成")
}

// demonstrateFullTextSearch 演示MySQL的全文索引功能
// MySQL支持对文本字段创建全文索引，提供高效的文本搜索功能
// 参数:
//   - db: GORM数据库连接实例
func demonstrateFullTextSearch(db *gorm.DB) {
	fmt.Println("\n--- MySQL全文索引演示 ---")

	// 为Post表的title和content字段创建全文索引
	db.Exec("ALTER TABLE post ADD FULLTEXT(title, content)")

	// 使用全文搜索查询
	var posts []Post
	db.Where("MATCH(title, content) AGAINST(? IN NATURAL LANGUAGE MODE)", "技术 编程").Find(&posts)
	fmt.Printf("全文搜索找到 %d 篇相关文章\n", len(posts))

	fmt.Println("✓ 全文索引演示完成")
}

// demonstratePartitioning 演示MySQL的分区表功能
// MySQL支持表分区，可以提高大表的查询性能和管理效率
// 参数:
//   - db: GORM数据库连接实例
func demonstratePartitioning(db *gorm.DB) {
	fmt.Println("\n--- MySQL分区表演示 ---")

	// 创建按日期分区的日志表
	db.Exec(`
		CREATE TABLE IF NOT EXISTS access_log (
			id INT AUTO_INCREMENT,
			user_id INT NOT NULL,
			access_time DATETIME NOT NULL,
			ip_address VARCHAR(45),
			PRIMARY KEY (id, access_time)
		) PARTITION BY RANGE (YEAR(access_time)) (
			PARTITION p2023 VALUES LESS THAN (2024),
			PARTITION p2024 VALUES LESS THAN (2025),
			PARTITION p_future VALUES LESS THAN MAXVALUE
		)
	`)

	fmt.Println("✓ 分区表演示完成")
}

// demonstrateStorageEngines 演示MySQL的存储引擎特性
// MySQL支持多种存储引擎，如InnoDB、MyISAM等，各有特点
// 参数:
//   - db: GORM数据库连接实例
func demonstrateStorageEngines(db *gorm.DB) {
	fmt.Println("\n--- MySQL存储引擎演示 ---")

	// 查询当前数据库支持的存储引擎
	var engines []struct {
		Engine  string
		Support string
		Comment string
	}
	db.Raw("SHOW ENGINES").Scan(&engines)

	fmt.Println("支持的存储引擎:")
	for _, engine := range engines {
		if engine.Support == "YES" || engine.Support == "DEFAULT" {
			fmt.Printf("- %s: %s\n", engine.Engine, engine.Comment)
		}
	}

	fmt.Println("✓ 存储引擎演示完成")
}

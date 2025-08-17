// 04_unit_exercises/level5_performance.go - Level 5 性能优化练习
// 对应文档：03_GORM单元练习_基础技能训练.md
// 本文件实现了GORM的性能优化练习，包括索引优化、查询优化、批量操作等
// 支持SQLite和MySQL两种数据库类型，提供完整的性能测试和优化建议

package main

import (
	"fmt"       // 格式化输出，用于打印测试结果和日志信息
	"log"       // 日志记录，用于错误处理和调试信息输出
	"math/rand" // 随机数生成，用于生成测试数据
	"time"      // 时间处理，用于性能测试计时和时间字段

	"gorm.io/driver/mysql"  // MySQL数据库驱动，支持MySQL数据库连接
	"gorm.io/driver/sqlite" // SQLite数据库驱动，支持SQLite数据库连接
	"gorm.io/gorm"          // GORM核心库，提供ORM功能
	"gorm.io/gorm/logger"   // GORM日志组件，用于SQL日志记录和调试
	"gorm.io/gorm/schema"   // GORM模式配置，用于表名和字段名策略配置
)

// 数据库配置相关定义

// DatabaseType 数据库类型枚举
// 定义支持的数据库类型，目前支持SQLite和MySQL
type DatabaseType string

const (
	SQLite DatabaseType = "sqlite" // SQLite数据库类型，轻量级文件数据库
	MySQL  DatabaseType = "mysql"  // MySQL数据库类型，企业级关系数据库
)

// DatabaseConfig 数据库配置结构体
// 包含数据库连接和连接池的所有配置参数，支持不同数据库类型的灵活配置
type DatabaseConfig struct {
	Type         DatabaseType    // 数据库类型(sqlite/mysql)
	DSN          string          // 数据源名称,用于指定数据库连接字符串
	MaxOpenConns int             // 最大打开连接数，控制并发连接数量
	MaxIdleConns int             // 最大空闲连接数，控制连接池中保持的空闲连接
	MaxLifetime  time.Duration   // 连接最大生命周期，防止长时间连接导致的问题
	LogLevel     logger.LogLevel // 日志级别，控制SQL日志的详细程度
}

// GetDefaultConfig 获取SQLite默认配置
// 返回一个包含默认参数的SQLite数据库配置对象，适用于开发和测试环境
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:         SQLite,
		DSN:          "level5_performance.db", // SQLite数据库文件名
		MaxOpenConns: 10,                      // SQLite建议较少的连接数
		MaxIdleConns: 5,                       // 最大空闲连接5个
		MaxLifetime:  time.Hour,               // 连接生命周期1小时
		LogLevel:     logger.Silent,           // 性能测试时关闭日志以提高准确性
	}
}

// GetMySQLConfig 获取MySQL配置
// 参数dsn: MySQL数据库连接字符串，格式如"user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
// 返回一个包含默认参数的MySQL数据库配置对象，适用于生产环境
func GetMySQLConfig(dsn string) *DatabaseConfig {
	return &DatabaseConfig{
		Type:         MySQL,
		DSN:          dsn,
		MaxOpenConns: 50,            // MySQL支持更高的并发连接数
		MaxIdleConns: 25,            // 更多的空闲连接以提高性能
		MaxLifetime:  time.Hour,     // 连接生命周期1小时
		LogLevel:     logger.Silent, // 性能测试时关闭日志
	}
}

// 基础模型定义

// BaseModel 基础模型结构体
// 包含所有数据库表通用的字段，采用GORM的软删除机制
// 所有业务模型都应该嵌入此结构体以获得统一的基础字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`              // 主键ID，自动递增，作为表的唯一标识
	CreatedAt time.Time      `json:"created_at"`                        // 创建时间，GORM自动管理，记录数据创建的时间戳
	UpdatedAt time.Time      `json:"updated_at"`                        // 更新时间，GORM自动管理，记录数据最后更新的时间戳
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 删除时间，用于软删除，建立索引以提高查询性能
}

// 业务模型定义

// User 用户模型
// 代表系统中的用户实体，包含用户的基本信息、工作信息和关联关系
// 支持完整的用户生命周期管理，包括注册、登录、状态管理等功能
// 与Post、Comment、Order等模型建立关联关系，支持复杂的业务查询
type User struct {
	BaseModel              // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	Username    string     `gorm:"uniqueIndex:idx_username;size:50;not null" json:"username"` // 用户名，最大50字符，唯一索引，用于登录和用户标识
	Email       string     `gorm:"uniqueIndex:idx_email;size:100;not null" json:"email"`      // 用户邮箱，最大100字符，唯一索引，用于登录和通信
	FirstName   string     `gorm:"size:50;not null;index:idx_name" json:"first_name"`         // 名字，最大50字符，非空，建立复合索引用于姓名搜索
	LastName    string     `gorm:"size:50;not null;index:idx_name" json:"last_name"`          // 姓氏，最大50字符，非空，建立复合索引用于姓名搜索
	Age         int        `gorm:"check:age >= 0 AND age <= 150;index:idx_age" json:"age"`    // 年龄，范围0-150，建立索引用于年龄统计和筛选
	City        string     `gorm:"size:100;index:idx_location" json:"city"`                   // 所在城市，最大100字符，建立位置索引用于地域统计
	Country     string     `gorm:"size:100;index:idx_location" json:"country"`                // 所在国家，最大100字符，建立位置索引用于地域统计
	Salary      float64    `gorm:"precision:10;scale:2;index:idx_salary" json:"salary"`       // 薪资，精度10位小数2位，建立索引用于薪资统计和筛选
	Department  string     `gorm:"size:100;index:idx_department" json:"department"`           // 部门，最大100字符，建立索引用于部门管理和统计
	Position    string     `gorm:"size:100;index:idx_position" json:"position"`               // 职位，最大100字符，建立索引用于职位管理和统计
	JoinDate    time.Time  `gorm:"index:idx_join_date" json:"join_date"`                      // 入职日期，建立索引用于工龄统计和查询
	IsActive    bool       `gorm:"default:true;index:idx_active" json:"is_active"`            // 是否活跃，默认true，建立索引用于快速筛选活跃用户
	LastLoginAt *time.Time `gorm:"index:idx_last_login" json:"last_login_at"`                 // 最后登录时间，可为空，建立索引用于用户活跃度分析

	// 关联关系定义
	Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`    // 用户发布的文章列表，一对多关系，通过AuthorID关联
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"` // 用户发表的评论列表，一对多关系，通过AuthorID关联
	Orders   []Order   `gorm:"foreignKey:UserID" json:"orders,omitempty"`     // 用户的订单列表，一对多关系，通过UserID关联
}

// Category 分类模型
// 代表内容分类系统，用于组织和管理文章内容，支持层级分类结构
// 与Post模型建立一对多关系，支持分类层次管理和排序功能
// 提供完整的分类管理功能，包括父子关系、排序、状态管理等
type Category struct {
	BaseModel          // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	Name        string `gorm:"size:100;not null;index:idx_category_name" json:"name"`       // 分类名称，最大100字符，非空，建立索引用于分类搜索和查询
	Slug        string `gorm:"uniqueIndex:idx_category_slug;size:100;not null" json:"slug"` // 分类别名，最大100字符，唯一索引，用于URL友好的分类标识
	Description string `gorm:"type:text" json:"description"`                                // 分类描述，文本类型，可存储较长的分类说明和介绍信息
	ParentID    *uint  `gorm:"index:idx_parent" json:"parent_id"`                           // 父分类ID，可为空，建立索引用于层级查询，支持无限级分类
	Level       int    `gorm:"default:1;index:idx_level" json:"level"`                      // 分类层级，默认1，建立索引用于层级筛选和树形结构查询
	SortOrder   int    `gorm:"default:0;index:idx_sort" json:"sort_order"`                  // 排序顺序，默认0，建立索引用于分类排序显示
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`              // 是否启用，默认true，建立索引用于快速筛选启用的分类

	// 关联关系定义
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // 父分类，通过ParentID关联，支持获取上级分类信息
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子分类列表，一对多关系，支持获取下级分类列表
	Posts    []Post     `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`  // 该分类下的文章列表，一对多关系，通过CategoryID关联
}

// Post 文章模型
// 代表系统中的文章实体，包含文章的完整信息和状态管理
// 支持文章发布、编辑、分类、标签、评论等完整的内容管理功能
// 与User、Category、Comment、Tag等模型建立关联关系，支持复杂的内容查询和统计
type Post struct {
	BaseModel               // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	Title        string     `gorm:"size:200;not null;index:idx_title" json:"title"`               // 文章标题，最大200字符，非空，建立索引用于标题搜索和排序
	Slug         string     `gorm:"uniqueIndex:idx_post_slug;size:200;not null" json:"slug"`      // 文章别名，最大200字符，唯一索引，用于SEO友好的URL标识
	Content      string     `gorm:"type:text;not null" json:"content"`                            // 文章正文内容，文本类型，非空，可存储长篇文章内容和富文本
	Excerpt      string     `gorm:"size:500" json:"excerpt"`                                      // 文章摘要，最大500字符，用于文章列表显示和SEO描述
	Status       string     `gorm:"size:20;default:'draft';index:idx_status" json:"status"`       // 文章状态，最大20字符，默认草稿，建立索引用于状态筛选(draft/published/archived)
	ViewCount    int        `gorm:"default:0;index:idx_views" json:"view_count"`                  // 浏览次数，默认0，建立索引用于热门文章排序和统计分析
	LikeCount    int        `gorm:"default:0;index:idx_likes" json:"like_count"`                  // 点赞次数，默认0，建立索引用于热门文章排序和用户互动统计
	CommentCount int        `gorm:"default:0;index:idx_comments" json:"comment_count"`            // 评论数量，默认0，建立索引用于评论统计和热门文章排序
	PublishedAt  *time.Time `gorm:"index:idx_published" json:"published_at"`                      // 发布时间，可为空，建立索引用于发布时间排序和时间范围查询
	Rating       float64    `gorm:"precision:3;scale:2;default:0;index:idx_rating" json:"rating"` // 文章评分，精度3位小数2位，默认0，建立索引用于评分排序和质量统计
	Featured     bool       `gorm:"default:false;index:idx_featured" json:"featured"`             // 是否精选文章，默认false，建立索引用于快速筛选精选内容

	// 外键定义
	AuthorID   uint  `gorm:"not null;index:idx_author" json:"author_id"` // 作者ID，非空，建立索引用于按作者查询文章和作者统计
	CategoryID *uint `gorm:"index:idx_category" json:"category_id"`      // 分类ID，可为空，建立索引用于按分类查询文章和分类统计

	// 关联关系定义
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`     // 文章作者，多对一关系，通过AuthorID关联User表，支持作者信息查询
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 文章分类，多对一关系，通过CategoryID关联Category表，支持分类管理
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`     // 文章评论列表，一对多关系，通过PostID关联Comment表，支持评论管理
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`      // 文章标签列表，多对多关系，通过post_tags中间表关联Tag表，支持标签分类
}

// Comment 评论模型
// 代表文章评论系统，支持多级评论和回复功能，包含评论内容、状态管理和层级结构
// 与User、Post模型建立关联关系，支持评论的层级管理、状态控制和互动统计
// 提供完整的评论管理功能，包括审核、点赞、回复等社交功能
type Comment struct {
	BaseModel        // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	Content   string `gorm:"type:text;not null" json:"content"`                        // 评论内容，文本类型，非空，可存储较长的评论文本和富文本内容
	Status    string `gorm:"size:20;default:'pending';index:idx_status" json:"status"` // 评论状态，最大20字符，默认待审核，建立索引用于状态筛选(pending/approved/spam/rejected)
	LikeCount int    `gorm:"default:0;index:idx_likes" json:"like_count"`              // 点赞数量，默认0，建立索引用于热门评论排序和用户互动统计
	ParentID  *uint  `gorm:"index:idx_parent" json:"parent_id"`                        // 父评论ID，可为空，建立索引用于层级评论查询和回复关系管理
	Level     int    `gorm:"default:1;index:idx_level" json:"level"`                   // 评论层级，默认1级，建立索引用于层级筛选和树形结构展示

	// 外键定义
	PostID   uint `gorm:"not null;index:idx_post" json:"post_id"`     // 所属文章ID，非空，建立索引用于按文章查询评论和评论统计
	AuthorID uint `gorm:"not null;index:idx_author" json:"author_id"` // 评论作者ID，非空，建立索引用于按用户查询评论和用户活跃度统计

	// 关联关系定义
	Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`       // 所属文章，多对一关系，通过PostID关联Post表，支持获取文章信息
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`   // 评论作者，多对一关系，通过AuthorID关联User表，支持获取作者信息
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // 父评论，自关联关系，通过ParentID关联，支持获取上级评论信息
	Children []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子评论列表，一对多自关联，支持获取回复列表和多级评论展示
}

// Tag 标签模型
// 代表内容标签系统，用于文章的分类和标记功能，支持标签的使用统计和状态管理
// 与Post模型建立多对多关系，支持灵活的内容标记和分类，提供完整的标签管理功能
// 包含标签名称、别名、颜色、使用次数等属性，支持标签的层次化管理和统计分析
type Tag struct {
	BaseModel         // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	Name       string `gorm:"uniqueIndex:idx_tag_name;size:50;not null" json:"name"` // 标签名称，最大50字符，非空，唯一索引确保标签名不重复，用于标签识别和显示
	Slug       string `gorm:"uniqueIndex:idx_tag_slug;size:50;not null" json:"slug"` // 标签别名，最大50字符，非空，唯一索引，用于URL友好的标签标识和SEO优化
	Color      string `gorm:"size:7;default:'#007bff'" json:"color"`                 // 标签颜色，7字符十六进制颜色值，默认蓝色，用于前端标签显示和视觉区分
	UsageCount int    `gorm:"default:0;index:idx_usage" json:"usage_count"`          // 使用次数，默认0，建立索引用于热门标签排序和标签使用统计分析
	IsActive   bool   `gorm:"default:true;index:idx_active" json:"is_active"`        // 是否启用，默认true，建立索引用于快速筛选启用的标签和标签状态管理

	// 关联关系定义
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"` // 使用该标签的文章列表，多对多关系，通过post_tags中间表关联，支持标签文章管理
}

// Order 订单模型
// 代表电商系统中的订单实体，包含订单的完整信息和状态管理
// 支持订单创建、支付、发货、完成等完整的订单生命周期管理
// 与User、OrderItem模型建立关联关系，支持复杂的订单查询和统计
type Order struct {
	BaseModel              // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	OrderNumber string     `gorm:"uniqueIndex:idx_order_number;size:50;not null" json:"order_number"`  // 订单号，最大50字符，非空，唯一索引确保订单号不重复，用于订单标识和查询
	UserID      uint       `gorm:"not null;index:idx_user" json:"user_id"`                             // 用户ID，非空，建立索引用于按用户查询订单和用户订单统计
	TotalAmount float64    `gorm:"precision:10;scale:2;not null;index:idx_amount" json:"total_amount"` // 订单总金额，精度10位小数2位，非空，建立索引用于金额统计和财务分析
	Status      string     `gorm:"size:20;default:'pending';index:idx_status" json:"status"`           // 订单状态，最大20字符，默认待处理，建立索引用于状态筛选(pending/paid/shipped/completed/cancelled)
	OrderDate   time.Time  `gorm:"index:idx_order_date" json:"order_date"`                             // 订单日期，建立索引用于时间范围查询、订单统计和报表分析
	ShippedAt   *time.Time `gorm:"index:idx_shipped" json:"shipped_at"`                                // 发货时间，可为空，建立索引用于发货统计和物流跟踪

	// 关联关系定义
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`         // 订单用户，多对一关系，通过UserID关联User表，支持获取用户信息
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"` // 订单项列表，一对多关系，通过OrderID关联OrderItem表，支持订单明细管理
}

// OrderItem 订单项模型
// 代表订单中的具体商品项，包含商品信息、数量、价格和小计金额
// 与Order、Product模型建立关联关系，支持订单明细管理和商品销售统计
// 提供完整的订单项管理功能，包括数量调整、价格记录、库存扣减等
type OrderItem struct {
	BaseModel         // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	OrderID   uint    `gorm:"not null;index:idx_order" json:"order_id"`      // 订单ID，非空，建立索引用于按订单查询订单项和订单明细统计
	ProductID uint    `gorm:"not null;index:idx_product" json:"product_id"`  // 商品ID，非空，建立索引用于按商品查询销售记录和商品销量统计
	Quantity  int     `gorm:"not null" json:"quantity"`                      // 商品数量，非空，用于库存管理、销量统计和订单金额计算
	Price     float64 `gorm:"precision:10;scale:2;not null" json:"price"`    // 商品单价，精度10位小数2位，非空，记录购买时的价格，用于价格历史追踪
	Subtotal  float64 `gorm:"precision:10;scale:2;not null" json:"subtotal"` // 小计金额，精度10位小数2位，非空，等于单价乘以数量，用于订单总额计算

	// 关联关系定义
	Order   Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`     // 所属订单，多对一关系，通过OrderID关联Order表，支持获取订单信息
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"` // 商品信息，多对一关系，通过ProductID关联Product表，支持获取商品详情
}

// Product 商品模型
// 代表电商系统中的商品实体，包含商品的基本信息、价格、库存和状态管理
// 支持商品创建、编辑、分类、库存管理等完整的商品生命周期管理
// 与Category、OrderItem模型建立关联关系，支持复杂的商品查询和销售统计
type Product struct {
	BaseModel           // 嵌入基础模型，获得ID、创建时间、更新时间、删除时间等通用字段
	Name        string  `gorm:"size:200;not null;index:idx_name" json:"name"`               // 商品名称，最大200字符，非空，建立索引用于商品搜索和名称查询
	SKU         string  `gorm:"uniqueIndex:idx_sku;size:50;not null" json:"sku"`            // 商品SKU，最大50字符，非空，唯一索引确保SKU不重复，用于商品标识和库存管理
	Description string  `gorm:"type:text" json:"description"`                               // 商品描述，文本类型，可存储详细的商品介绍和规格说明
	Price       float64 `gorm:"precision:10;scale:2;not null;index:idx_price" json:"price"` // 商品价格，精度10位小数2位，非空，建立索引用于价格排序和价格区间查询
	Stock       int     `gorm:"default:0;index:idx_stock" json:"stock"`                     // 库存数量，默认0，建立索引用于库存查询和库存预警
	CategoryID  *uint   `gorm:"index:idx_category" json:"category_id"`                      // 分类ID，可为空，建立索引用于按分类查询商品和分类统计
	IsActive    bool    `gorm:"default:true;index:idx_active" json:"is_active"`             // 是否启用，默认true，建立索引用于快速筛选上架商品

	// 关联关系定义
	Category   *Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`   // 商品分类，多对一关系，通过CategoryID关联Category表，支持分类管理
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"order_items,omitempty"` // 订单项列表，一对多关系，通过ProductID关联OrderItem表，支持销售记录查询
}

// 数据库初始化相关函数

// initDB 初始化SQLite数据库连接（保持向后兼容）
// 使用默认的SQLite配置，适用于开发和测试环境
// 返回配置好的GORM数据库实例，包含连接池设置和自动迁移
func initDB() *gorm.DB {
	config := GetDefaultConfig() // 获取SQLite默认配置
	return InitDatabase(config)  // 使用通用初始化函数
}

// InitDatabase 通用数据库初始化函数
// 参数config: 数据库配置对象，包含数据库类型、连接字符串、连接池参数等
// 返回配置好的GORM数据库实例，支持SQLite和MySQL两种数据库类型
// 包含完整的连接池配置、日志设置、连接测试和自动迁移功能
func InitDatabase(config *DatabaseConfig) *gorm.DB {
	var dialector gorm.Dialector

	// 根据数据库类型选择相应的驱动
	switch config.Type {
	case SQLite:
		dialector = sqlite.Open(config.DSN) // 使用SQLite驱动打开数据库文件
	case MySQL:
		dialector = mysql.Open(config.DSN) // 使用MySQL驱动连接数据库服务器
	default:
		log.Fatalf("Unsupported database type: %s", config.Type) // 不支持的数据库类型
	}

	// 配置GORM选项
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel), // 设置日志级别，性能测试时建议使用Silent
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // 表名前缀，为空表示不添加前缀
			SingularTable: false, // 使用复数表名，遵循GORM默认约定
		},
		PrepareStmt:                              true, // 启用预编译语句以提高性能
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束以提高迁移性能
	}

	// 打开数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", config.Type, err)
	}

	// 获取底层sql.DB对象用于连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// 配置连接池参数，优化数据库性能
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)   // 设置最大打开连接数，控制并发连接数量
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)   // 设置最大空闲连接数，提高连接复用效率
	sqlDB.SetConnMaxLifetime(config.MaxLifetime) // 设置连接最大生命周期，防止长时间连接问题

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 执行自动迁移，按依赖顺序创建表结构
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建复合索引
	createCompositeIndexes(db)

	fmt.Printf("Successfully connected to %s database\n", config.Type)
	return db
}

// AutoMigrate 执行数据库自动迁移
// 参数db: GORM数据库实例
// 按照模型依赖关系的正确顺序执行迁移，确保外键约束正确建立
// 返回错误信息，如果迁移失败则返回具体的错误原因
func AutoMigrate(db *gorm.DB) error {
	// 按依赖顺序迁移模型，确保外键关系正确建立
	// 1. 首先迁移基础模型（无外键依赖）
	if err := db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("failed to migrate User: %w", err)
	}

	if err := db.AutoMigrate(&Category{}); err != nil {
		return fmt.Errorf("failed to migrate Category: %w", err)
	}

	if err := db.AutoMigrate(&Tag{}); err != nil {
		return fmt.Errorf("failed to migrate Tag: %w", err)
	}

	if err := db.AutoMigrate(&Product{}); err != nil {
		return fmt.Errorf("failed to migrate Product: %w", err)
	}

	// 2. 然后迁移有外键依赖的模型
	if err := db.AutoMigrate(&Post{}); err != nil {
		return fmt.Errorf("failed to migrate Post: %w", err)
	}

	if err := db.AutoMigrate(&Comment{}); err != nil {
		return fmt.Errorf("failed to migrate Comment: %w", err)
	}

	if err := db.AutoMigrate(&Order{}); err != nil {
		return fmt.Errorf("failed to migrate Order: %w", err)
	}

	if err := db.AutoMigrate(&OrderItem{}); err != nil {
		return fmt.Errorf("failed to migrate OrderItem: %w", err)
	}

	fmt.Println("Database migration completed successfully")
	return nil
}

// createCompositeIndexes 创建复合索引以优化查询性能
// 参数db: GORM数据库实例
// 复合索引是包含多个列的索引，能够显著提高多条件查询的性能
// 索引的列顺序很重要：最常用作查询条件的列应该放在前面
func createCompositeIndexes(db *gorm.DB) {
	// 用户复合索引
	// idx_users_active_department: 优化按活跃状态和部门查询用户的性能
	// 适用于获取特定部门的活跃用户列表，常用于人事管理和统计分析
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_active_department ON users(is_active, department)")

	// idx_users_city_salary: 优化按城市和薪资查询用户的性能
	// 适用于地域薪资统计、人才分布分析等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_city_salary ON users(city, salary)")

	// idx_users_join_date_active: 优化按入职时间和活跃状态查询的性能
	// 适用于员工入职统计、活跃用户时间分析等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_join_date_active ON users(join_date, is_active)")

	// 文章复合索引
	// idx_posts_status_published: 优化按状态和发布时间查询文章的性能
	// 适用于获取已发布文章列表、按时间排序等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_status_published ON posts(status, published_at)")

	// idx_posts_author_status: 优化按作者和状态查询文章的性能
	// 适用于获取特定作者的特定状态文章，如草稿、已发布等
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_author_status ON posts(author_id, status)")

	// idx_posts_category_featured: 优化按分类和精选状态查询的性能
	// 适用于获取特定分类的精选文章，提高首页推荐效率
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_category_featured ON posts(category_id, featured)")

	// idx_posts_views_rating: 优化按浏览量和评分排序的性能
	// 适用于热门文章排序、高质量内容筛选等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_views_rating ON posts(view_count, rating)")

	// 评论复合索引
	// idx_comments_post_status: 优化按文章和状态查询评论的性能
	// 适用于获取特定文章的已审核评论、评论管理等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_status ON comments(post_id, status)")

	// idx_comments_author_created: 优化按作者和创建时间查询评论的性能
	// 适用于获取用户评论历史、按时间排序用户活动等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_author_created ON comments(author_id, created_at)")

	// 订单复合索引
	// idx_orders_user_status: 优化按用户和订单状态查询的性能
	// 适用于获取用户的特定状态订单（如待支付、已完成等）
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status)")

	// idx_orders_date_amount: 优化按订单日期和金额查询的性能
	// 适用于财务统计、销售报表、订单金额分析等场景
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_date_amount ON orders(order_date, total_amount)")

	fmt.Println("✓ 复合索引创建完成")
}

// generateTestData 生成测试数据用于性能测试
// 参数db: GORM数据库实例
// 参数userCount: 要生成的用户数量
// 参数postCount: 要生成的文章数量
// 参数commentCount: 要生成的评论数量
// 该函数创建大量模拟数据，包括用户、分类、标签、文章、评论等
// 使用批量插入优化性能，避免逐条插入造成的性能问题
func generateTestData(db *gorm.DB, userCount, postCount, commentCount int) {
	fmt.Printf("开始生成测试数据: %d用户, %d文章, %d评论\n", userCount, postCount, commentCount)
	start := time.Now()

	// 生成用户数据 - 用户是系统的核心实体，其他数据都会关联到用户
	users := make([]User, userCount)
	// 预定义数据数组，用于随机生成真实感的测试数据
	departments := []string{"技术部", "市场部", "销售部", "人事部", "财务部"}         // 5个主要部门
	positions := []string{"工程师", "经理", "主管", "专员", "总监"}               // 5个职位层级
	cities := []string{"北京", "上海", "深圳", "广州", "杭州", "成都", "武汉", "西安"} // 8个主要城市
	countries := []string{"中国", "美国", "日本", "韩国", "新加坡"}               // 5个国家

	// 循环生成指定数量的用户数据
	for i := 0; i < userCount; i++ {
		users[i] = User{
			Username:   fmt.Sprintf("user%d", i+1),                                        // 用户名：user1, user2, ...
			Email:      fmt.Sprintf("user%d@example.com", i+1),                            // 邮箱：user1@example.com, ...
			FirstName:  fmt.Sprintf("First%d", i+1),                                       // 名字：First1, First2, ...
			LastName:   fmt.Sprintf("Last%d", i+1),                                        // 姓氏：Last1, Last2, ...
			Age:        20 + rand.Intn(40),                                                // 年龄：20-59岁随机
			City:       cities[rand.Intn(len(cities))],                                    // 城市：从预定义城市中随机选择
			Country:    countries[rand.Intn(len(countries))],                              // 国家：从预定义国家中随机选择
			Salary:     float64(5000 + rand.Intn(20000)),                                  // 薪资：5000-25000随机
			Department: departments[rand.Intn(len(departments))],                          // 部门：从预定义部门中随机选择
			Position:   positions[rand.Intn(len(positions))],                              // 职位：从预定义职位中随机选择
			JoinDate:   time.Now().AddDate(-rand.Intn(5), -rand.Intn(12), -rand.Intn(30)), // 入职时间：过去5年内随机
			IsActive:   rand.Float32() > 0.1,                                              // 活跃状态：90%概率为活跃用户
		}
		// 70%的用户有最后登录时间，模拟真实的用户活跃度分布
		if rand.Float32() > 0.3 {
			lastLogin := time.Now().AddDate(0, 0, -rand.Intn(30)) // 过去30天内的随机登录时间
			users[i].LastLoginAt = &lastLogin
		}
	}

	// 批量插入用户数据，每批100条记录，提高插入性能
	db.CreateInBatches(users, 100)
	fmt.Printf("✓ 用户数据生成完成: %d条\n", userCount)

	// 生成分类数据
	categories := []Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", Level: 1, SortOrder: 1, IsActive: true},
		{Name: "生活", Slug: "life", Description: "生活分享", Level: 1, SortOrder: 2, IsActive: true},
		{Name: "旅游", Slug: "travel", Description: "旅游攻略", Level: 1, SortOrder: 3, IsActive: true},
		{Name: "美食", Slug: "food", Description: "美食推荐", Level: 1, SortOrder: 4, IsActive: true},
		{Name: "娱乐", Slug: "entertainment", Description: "娱乐资讯", Level: 1, SortOrder: 5, IsActive: true},
	}
	db.Create(&categories)

	// 生成标签数据
	tags := []Tag{
		{Name: "Go", Slug: "go", Color: "#00ADD8", IsActive: true},
		{Name: "Python", Slug: "python", Color: "#3776AB", IsActive: true},
		{Name: "JavaScript", Slug: "javascript", Color: "#F7DF1E", IsActive: true},
		{Name: "数据库", Slug: "database", Color: "#336791", IsActive: true},
		{Name: "前端", Slug: "frontend", Color: "#61DAFB", IsActive: true},
		{Name: "后端", Slug: "backend", Color: "#68217A", IsActive: true},
		{Name: "教程", Slug: "tutorial", Color: "#FF6B6B", IsActive: true},
		{Name: "实战", Slug: "practice", Color: "#4ECDC4", IsActive: true},
	}
	db.Create(&tags)

	// 生成文章数据
	posts := make([]Post, postCount)
	statuses := []string{"published", "draft", "archived"}

	for i := 0; i < postCount; i++ {
		authorID := uint(rand.Intn(userCount) + 1)
		categoryID := uint(rand.Intn(len(categories)) + 1)
		status := statuses[rand.Intn(len(statuses))]

		posts[i] = Post{
			Title:        fmt.Sprintf("文章标题 %d", i+1),
			Slug:         fmt.Sprintf("post-%d", i+1),
			Content:      fmt.Sprintf("这是文章 %d 的内容，包含了丰富的信息和详细的描述...", i+1),
			Excerpt:      fmt.Sprintf("文章 %d 的摘要", i+1),
			Status:       status,
			ViewCount:    rand.Intn(10000),
			LikeCount:    rand.Intn(1000),
			CommentCount: rand.Intn(100),
			Rating:       float64(rand.Intn(50))/10.0 + 1.0, // 1.0-6.0
			Featured:     rand.Float32() > 0.8,              // 20%的文章是精选的
			AuthorID:     authorID,
			CategoryID:   &categoryID,
		}

		if status == "published" {
			publishedAt := time.Now().AddDate(0, 0, -rand.Intn(365))
			posts[i].PublishedAt = &publishedAt
		}
	}

	// 批量插入文章
	db.CreateInBatches(posts, 100)
	fmt.Printf("✓ 文章数据生成完成: %d条\n", postCount)

	// 生成评论数据
	comments := make([]Comment, commentCount)
	commentStatuses := []string{"approved", "pending", "spam"}

	for i := 0; i < commentCount; i++ {
		comments[i] = Comment{
			Content:   fmt.Sprintf("这是评论 %d 的内容", i+1),
			Status:    commentStatuses[rand.Intn(len(commentStatuses))],
			LikeCount: rand.Intn(100),
			Level:     1,
			PostID:    uint(rand.Intn(postCount) + 1),
			AuthorID:  uint(rand.Intn(userCount) + 1),
		}
	}

	// 批量插入评论
	db.CreateInBatches(comments, 100)
	fmt.Printf("✓ 评论数据生成完成: %d条\n", commentCount)

	elapsed := time.Since(start)
	fmt.Printf("✓ 测试数据生成完成，耗时: %v\n", elapsed)
}

// 性能测试函数

// benchmarkBasicQueries 基本查询性能测试
// 参数db: GORM数据库实例
// 测试各种基础查询操作的性能，包括简单查询、索引查询、范围查询和复合条件查询
// 这些是日常开发中最常用的查询类型，性能优化的重点关注对象
func benchmarkBasicQueries(db *gorm.DB) {
	fmt.Println("\n=== 基本查询性能测试 ===")

	// 测试简单查询 - 单一条件查询，使用索引字段
	// 这是最基础的查询方式，性能取决于索引的使用情况
	start := time.Now()
	var users []User
	db.Where("is_active = ?", true).Limit(100).Find(&users) // 查询活跃用户，限制100条结果
	fmt.Printf("简单查询 (100条): %v\n", time.Since(start))

	// 测试带索引的查询 - 使用已建立索引的字段进行查询
	// 索引能显著提高查询性能，特别是在大数据集上
	start = time.Now()
	var usersByCity []User
	db.Where("city = ?", "北京").Limit(100).Find(&usersByCity) // 按城市查询，city字段已建立索引
	fmt.Printf("索引查询 (城市): %v\n", time.Since(start))

	// 测试范围查询 - 使用BETWEEN进行范围条件查询
	// 范围查询在有索引的情况下性能较好，但比等值查询稍慢
	start = time.Now()
	var usersBySalary []User
	db.Where("salary BETWEEN ? AND ?", 8000, 15000).Limit(100).Find(&usersBySalary) // 薪资范围查询
	fmt.Printf("范围查询 (薪资): %v\n", time.Since(start))

	// 测试复合条件查询 - 多个条件组合查询
	// 复合条件查询能利用复合索引提高性能，条件顺序很重要
	start = time.Now()
	var activeUsers []User
	db.Where("is_active = ? AND department = ?", true, "技术部").Limit(100).Find(&activeUsers) // 活跃状态+部门的复合查询
	fmt.Printf("复合条件查询: %v\n", time.Since(start))
}

// benchmarkPreloading 预加载性能对比测试
// 参数db: GORM数据库实例
// 对比N+1查询问题与预加载的性能差异，演示预加载的重要性和最佳实践
// 预加载是解决ORM性能问题的关键技术，能将多次查询合并为少数几次查询
func benchmarkPreloading(db *gorm.DB) {
	fmt.Println("\n=== 预加载性能测试 ===")

	// N+1 查询问题演示 - 经典的ORM性能陷阱
	// 首先执行1次查询获取文章列表，然后为每篇文章执行1次查询获取作者信息
	// 总共执行N+1次查询（N为文章数量），性能极差
	start := time.Now()
	var postsWithoutPreload []Post
	db.Limit(50).Find(&postsWithoutPreload) // 1次查询获取50篇文章
	// 访问关联数据会产生N+1查询问题
	for _, post := range postsWithoutPreload {
		var author User
		db.First(&author, post.AuthorID) // 每篇文章执行1次查询获取作者，共50次查询
		_ = author.Username
	}
	fmt.Printf("N+1查询 (50篇文章): %v\n", time.Since(start))

	// 使用预加载 - 解决N+1查询问题的标准方案
	// 使用JOIN或IN查询一次性获取所有关联数据，大幅提升性能
	start = time.Now()
	var postsWithPreload []Post
	db.Preload("Author").Limit(50).Find(&postsWithPreload) // 2次查询：1次获取文章，1次获取所有作者
	for _, post := range postsWithPreload {
		_ = post.Author.Username // 直接访问已预加载的数据，无需额外查询
	}
	fmt.Printf("预加载查询 (50篇文章): %v\n", time.Since(start))

	// 选择性预加载 - 只加载需要的字段，进一步优化性能
	// 通过Select指定需要的字段，减少数据传输量和内存使用
	start = time.Now()
	var postsWithSelectivePreload []Post
	db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name") // 只选择需要的用户字段
	}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, slug") // 只选择需要的分类字段
	}).Limit(50).Find(&postsWithSelectivePreload)
	fmt.Printf("选择性预加载 (50篇文章): %v\n", time.Since(start))

	// 多层预加载 - 预加载嵌套关联关系
	// 一次性加载文章、作者、分类、评论及评论作者，适用于复杂的数据展示场景
	start = time.Now()
	var postsWithNestedPreload []Post
	db.Preload("Author").Preload("Category").Preload("Comments").Preload("Comments.Author").Limit(20).Find(&postsWithNestedPreload)
	fmt.Printf("多层预加载 (20篇文章): %v\n", time.Since(start))
}

// benchmarkBatchOperations 批量操作性能测试
// 参数db: GORM数据库实例
// 对比单条操作与批量操作的性能差异，演示批量操作的重要性
// 批量操作能显著减少数据库交互次数，大幅提升大数据量操作的性能
func benchmarkBatchOperations(db *gorm.DB) {
	fmt.Println("\n=== 批量操作性能测试 ===")

	// 准备测试数据 - 创建1000条用户记录用于测试
	// 使用统一的测试数据确保测试结果的可比性
	testUsers := make([]User, 1000)
	for i := 0; i < 1000; i++ {
		testUsers[i] = User{
			Username:   fmt.Sprintf("batch_user_%d", i),          // 唯一用户名
			Email:      fmt.Sprintf("batch_user_%d@test.com", i), // 唯一邮箱
			FirstName:  fmt.Sprintf("First%d", i),                // 名字
			LastName:   fmt.Sprintf("Last%d", i),                 // 姓氏
			Age:        25,                                       // 统一年龄
			City:       "测试城市",                                   // 统一城市
			Country:    "测试国家",                                   // 统一国家
			Salary:     10000,                                    // 统一薪资
			Department: "测试部门",                                   // 统一部门
			Position:   "测试职位",                                   // 统一职位
			JoinDate:   time.Now(),                               // 当前时间
			IsActive:   true,                                     // 活跃状态
		}
	}

	// 单条插入测试 - 逐条执行INSERT语句
	// 每次插入都需要一次数据库交互，网络开销和事务开销较大
	start := time.Now()
	for i := 0; i < 100; i++ {
		db.Create(&testUsers[i]) // 每次调用Create都会执行一次INSERT语句
	}
	fmt.Printf("单条插入 (100条): %v\n", time.Since(start))

	// 批量插入测试 - 使用批量插入优化性能
	// 将多条记录合并为一个或少数几个INSERT语句，减少数据库交互次数
	start = time.Now()
	db.CreateInBatches(testUsers[100:200], 50) // 100条记录分2批插入，每批50条
	fmt.Printf("批量插入 (100条, 批次50): %v\n", time.Since(start))

	// 大批量插入测试 - 测试更大批次的插入性能
	// 批次大小需要平衡性能和内存使用，过大可能导致内存问题或SQL语句过长
	start = time.Now()
	db.CreateInBatches(testUsers[200:700], 100) // 500条记录分5批插入，每批100条
	fmt.Printf("大批量插入 (500条, 批次100): %v\n", time.Since(start))

	// 批量更新测试 - 使用单个UPDATE语句更新多条记录
	// 批量更新比逐条更新效率高得多，特别是在大数据量场景下
	start = time.Now()
	db.Model(&User{}).Where("username LIKE ?", "batch_user_%").Update("salary", 12000) // 一次性更新所有测试用户的薪资
	fmt.Printf("批量更新: %v\n", time.Since(start))

	// 清理测试数据 - 删除本次测试创建的所有数据
	// 使用批量删除确保测试环境的清洁
	db.Where("username LIKE ?", "batch_user_%").Delete(&User{})
}

// benchmarkComplexQueries 复杂查询性能测试
// 参数db: GORM数据库实例
// 测试多种复杂查询场景的性能，包括聚合查询、子查询、多表连接和分页查询
// 这些查询在实际业务中经常出现，但性能优化难度较大，需要重点关注
func benchmarkComplexQueries(db *gorm.DB) {
	fmt.Println("\n=== 复杂查询性能测试 ===")

	// 聚合查询性能测试 - 测试GROUP BY和聚合函数的性能
	// 聚合查询用于统计分析，在大数据量下可能较慢
	// 需要在分组字段上建立索引，并考虑使用物化视图等优化手段
	start := time.Now()
	// 定义用户统计结构体，用于接收聚合查询结果
	type UserStats struct {
		Department string  // 部门名称
		UserCount  int64   // 部门用户数量
		AvgSalary  float64 // 部门平均薪资
		MaxSalary  float64 // 部门最高薪资
		MinSalary  float64 // 部门最低薪资
	}
	var stats []UserStats
	// 执行聚合查询：按部门统计用户数量和薪资信息
	// 使用COUNT、AVG、MAX、MIN等聚合函数进行统计计算
	db.Model(&User{}).Select("department, COUNT(*) as user_count, AVG(salary) as avg_salary, MAX(salary) as max_salary, MIN(salary) as min_salary").
		Where("is_active = ?", true). // 只统计活跃用户
		Group("department").          // 按部门分组
		Scan(&stats)                  // 将结果扫描到结构体切片中
	fmt.Printf("聚合查询 (部门统计): %v\n", time.Since(start))

	// 子查询性能测试 - 测试嵌套查询的性能
	// 子查询在某些场景下比JOIN更直观，但性能可能不如JOIN
	// 现代数据库优化器通常会将子查询转换为JOIN操作
	start = time.Now()
	var topAuthors []User
	// 创建子查询：统计每个作者发布的文章数量，并按数量降序排列取前10名
	subQuery := db.Model(&Post{}).Select("author_id, COUNT(*) as post_count").
		Where("status = ?", "published"). // 只统计已发布的文章
		Group("author_id").               // 按作者分组
		Order("post_count DESC").         // 按文章数量降序排列
		Limit(10)                         // 取前10名
	// 主查询：使用子查询结果关联用户表，获取热门作者的详细信息
	db.Table("users u").Joins("JOIN (?) pc ON u.id = pc.author_id", subQuery). // 将子查询作为临时表进行JOIN
											Select("u.*").Find(&topAuthors) // 选择用户的所有字段
	fmt.Printf("子查询 (热门作者): %v\n", time.Since(start))

	// 复杂连接查询性能测试 - 测试多表JOIN操作的性能
	// 多表连接是关系数据库的核心功能，但在大数据量下性能可能成为瓶颈
	// 需要确保连接字段上有适当的索引，并考虑查询计划的优化
	start = time.Now()
	var results []map[string]interface{} // 使用map接收复杂查询结果
	// 执行四表连接查询：文章表 + 用户表 + 分类表 + 评论表
	// 统计最近一个月内热门文章（评论数大于5的文章）
	db.Table("posts p").
		// 选择需要的字段：文章标题、作者用户名、分类名称、评论数量
		Select("p.title, u.username, c.name as category_name, COUNT(cm.id) as comment_count").
		Joins("JOIN users u ON p.author_id = u.id").                                             // 内连接用户表获取作者信息
		Joins("LEFT JOIN categories c ON p.category_id = c.id").                                 // 左连接分类表（允许文章没有分类）
		Joins("LEFT JOIN comments cm ON p.id = cm.post_id AND cm.status = 'approved'").          // 左连接评论表（只统计已审核的评论）
		Where("p.status = ? AND p.published_at > ?", "published", time.Now().AddDate(0, -1, 0)). // 筛选最近一个月的已发布文章
		Group("p.id, p.title, u.username, c.name").                                              // 按文章、标题、用户名、分类名分组
		Having("comment_count > ?", 5).                                                          // 只显示评论数大于5的文章
		Order("comment_count DESC").                                                             // 按评论数降序排列
		Limit(20).                                                                               // 限制返回20条记录
		Scan(&results)                                                                           // 扫描结果到map切片
	fmt.Printf("复杂连接查询 (热门文章): %v\n", time.Since(start))

	// 分页查询性能测试 - 测试大数据量分页的性能
	// 分页查询在Web应用中非常常见，但OFFSET在大偏移量时性能较差
	// 可以考虑使用游标分页等优化方案
	start = time.Now()
	var paginatedPosts []Post
	var total int64
	// 先统计总记录数，用于计算总页数
	db.Model(&Post{}).Where("status = ?", "published").Count(&total)
	// 执行分页查询：跳过前100条记录，取20条记录（第6页，每页20条）
	db.Where("status = ?", "published").
		Offset(100).Limit(20). // OFFSET 100 LIMIT 20 实现分页
		// 预加载作者信息，但只选择必要的字段以优化性能
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username") // 只选择ID和用户名，减少数据传输量
		}).Find(&paginatedPosts)
	fmt.Printf("分页查询 (第6页, 20条/页, 总计%d条): %v\n", total, time.Since(start))

	// 窗口函数查询测试 - 测试高级SQL功能的性能（如果数据库支持）
	// 窗口函数提供了强大的分析功能，但计算复杂度较高
	// 适用于排名、累计计算、移动平均等高级分析场景
	start = time.Now()
	var rankedUsers []struct {
		Username   string  // 用户名
		Department string  // 部门
		Salary     float64 // 薪资
		Rank       int     // 部门内薪资排名
	}
	// 使用窗口函数计算每个部门内的薪资排名
	// ROW_NUMBER() OVER() 为每个部门内的用户按薪资降序排名
	db.Raw(`
		SELECT username, department, salary,
		       ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) as rank
		FROM users 
		WHERE is_active = true
		ORDER BY department, rank
	`).Scan(&rankedUsers)
	fmt.Printf("窗口函数查询 (部门内薪资排名): %v, 结果数量: %d\n", time.Since(start), len(rankedUsers))
}

// benchmarkIndexOptimization 索引优化效果测试
// 参数db: GORM数据库实例
// 测试不同类型索引对查询性能的影响，验证索引优化的效果
// 包括复合索引、范围查询、排序查询和模糊查询等场景
func benchmarkIndexOptimization(db *gorm.DB) {
	fmt.Println("\n=== 索引优化效果测试 ===")

	// 复合索引查询测试 - 测试多字段组合索引的效果
	// 复合索引能够显著提升多条件查询的性能
	// 索引字段的顺序很重要，应该将选择性高的字段放在前面
	start := time.Now()
	var indexedUsers []User
	// 使用复合索引 idx_users_city_active 进行查询
	// 该查询能够充分利用 (city, is_active) 复合索引
	db.Where("city = ? AND is_active = ?", "北京", true).Find(&indexedUsers)
	fmt.Printf("复合索引查询 (城市+状态): %v, 结果: %d条\n", time.Since(start), len(indexedUsers))

	// 范围查询测试 - 测试范围条件与等值条件组合的性能
	// 范围查询通常只能使用索引的前缀部分，后续字段可能无法使用索引
	// 建议将等值条件放在索引前面，范围条件放在后面
	start = time.Now()
	var salaryUsers []User
	// 查询特定部门薪资范围内的用户
	// 如果有 (department, salary) 复合索引，这个查询性能会很好
	db.Where("salary BETWEEN ? AND ? AND department = ?", 8000, 15000, "技术部").Find(&salaryUsers)
	fmt.Printf("范围+等值查询 (薪资+部门): %v, 结果: %d条\n", time.Since(start), len(salaryUsers))

	// 排序查询测试 - 测试ORDER BY子句对索引的利用
	// 如果排序字段有索引，数据库可以避免额外的排序操作
	// 复合索引 (view_count, rating) 可以同时优化WHERE和ORDER BY
	start = time.Now()
	var sortedPosts []Post
	// 查询已发布文章并按浏览量和评分降序排列
	// 利用 idx_posts_views_rating 复合索引优化排序性能
	db.Where("status = ?", "published").Order("view_count DESC, rating DESC").Limit(50).Find(&sortedPosts)
	fmt.Printf("排序查询 (浏览量+评分): %v, 结果: %d条\n", time.Since(start), len(sortedPosts))

	// 模糊查询测试 - 测试LIKE操作的性能
	// 前缀匹配（如 'abc%'）可以使用索引，但中间匹配（如 '%abc%'）通常无法使用索引
	// 对于全文搜索需求，建议使用专门的全文索引或搜索引擎
	start = time.Now()
	var searchPosts []Post
	// 模糊查询文章标题包含"文章"的记录
	// 这种中间匹配查询通常无法使用索引，性能较差
	// 在实际应用中可以考虑使用全文索引或搜索引擎优化
	db.Where("title LIKE ?", "%文章%").Limit(100).Find(&searchPosts)
	fmt.Printf("模糊查询 (标题中间匹配): %v, 结果: %d条\n", time.Since(start), len(searchPosts))

	// 前缀匹配查询测试 - 对比前缀匹配的性能
	// 前缀匹配可以有效利用索引，性能比中间匹配好很多
	start = time.Now()
	var prefixPosts []Post
	// 查询标题以"Go"开头的文章
	// 这种前缀匹配可以利用title字段的索引（如果存在）
	db.Where("title LIKE ?", "Go%").Limit(100).Find(&prefixPosts)
	fmt.Printf("前缀匹配查询 (标题前缀): %v, 结果: %d条\n", time.Since(start), len(prefixPosts))

	// 覆盖索引查询测试 - 测试只查询索引字段的性能
	// 当查询的所有字段都包含在索引中时，数据库可以只访问索引而不访问表数据
	// 这种"覆盖索引"查询性能最佳
	start = time.Now()
	var userCities []struct {
		City     string // 城市
		IsActive bool   // 活跃状态
	}
	// 只查询索引字段，可以实现覆盖索引查询
	// 如果有 (city, is_active) 索引，这个查询只需要访问索引
	db.Model(&User{}).Select("city, is_active").Where("city IN ?", []string{"北京", "上海", "深圳"}).Scan(&userCities)
	fmt.Printf("覆盖索引查询 (只查询索引字段): %v, 结果: %d条\n", time.Since(start), len(userCities))
}

// benchmarkConnectionAndTransaction 测试连接池和事务性能
// 该函数通过对比普通查询、事务查询和并发查询的性能，
// 帮助理解连接池配置和事务使用对数据库性能的影响
// 参数:
//   - db: GORM数据库实例，应该已经配置好连接池参数
func benchmarkConnectionAndTransaction(db *gorm.DB) {
	fmt.Println("\n=== 连接池和事务性能测试 ===")

	// 测试普通操作 - 直接使用连接池中的连接进行查询
	// 这种方式每次查询都会从连接池获取连接，查询完成后归还连接
	start := time.Now()
	for i := 0; i < 100; i++ {
		var user User                            // 声明用户变量
		db.First(&user, uint(rand.Intn(1000)+1)) // 随机查询一个用户，测试单次查询性能
	}
	fmt.Printf("普通查询 (100次): %v\n", time.Since(start))

	// 测试事务操作 - 在事务中执行多个查询
	// 事务会保持连接直到提交或回滚，可能影响连接池的使用效率
	start = time.Now()
	for i := 0; i < 10; i++ {
		// 使用GORM的Transaction方法创建事务
		db.Transaction(func(tx *gorm.DB) error {
			// 在同一个事务中执行10次查询
			for j := 0; j < 10; j++ {
				var user User                            // 声明用户变量
				tx.First(&user, uint(rand.Intn(1000)+1)) // 在事务中查询用户
			}
			return nil // 返回nil表示事务成功，将自动提交
		})
	}
	fmt.Printf("事务查询 (10个事务, 每个10次查询): %v\n", time.Since(start))

	// 测试并发查询 - 多个goroutine同时执行查询
	// 这将测试连接池在高并发情况下的性能表现
	start = time.Now()
	done := make(chan bool, 10) // 创建缓冲通道，用于同步goroutine完成状态
	for i := 0; i < 10; i++ {
		// 启动goroutine执行并发查询
		go func() {
			// 每个goroutine执行10次查询
			for j := 0; j < 10; j++ {
				var user User                            // 声明用户变量
				db.First(&user, uint(rand.Intn(1000)+1)) // 并发查询用户，测试连接池并发性能
			}
			done <- true // 通知该goroutine已完成
		}()
	}
	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done // 从通道接收完成信号
	}
	fmt.Printf("并发查询 (10个goroutine, 每个10次查询): %v\n", time.Since(start))
}

// showOptimizationTips 显示GORM性能优化建议
// 该函数总结了在实际开发中提升GORM性能的最佳实践和优化策略，
// 涵盖索引优化、查询优化、批量操作、连接池配置等多个方面
func showOptimizationTips() {
	fmt.Println("\n=== GORM 性能优化建议 ===")

	// 1. 索引优化建议
	// 索引是提升查询性能最重要的手段之一
	fmt.Println("1. 索引优化:")
	fmt.Println("   - 为经常查询的字段创建索引")  // 如WHERE条件中的字段
	fmt.Println("   - 使用复合索引优化多字段查询") // 多个字段组合查询时使用
	fmt.Println("   - 避免在小表上创建过多索引")  // 索引维护成本可能超过收益

	// 2. 查询优化建议
	// 通过优化查询语句减少数据传输和处理时间
	fmt.Println("\n2. 查询优化:")
	fmt.Println("   - 使用 Select 只查询需要的字段")     // 减少网络传输和内存使用
	fmt.Println("   - 合理使用 Preload 避免 N+1 查询") // 预加载关联数据，避免循环查询
	fmt.Println("   - 使用 Limit 限制查询结果数量")      // 避免一次性加载大量数据
	fmt.Println("   - 避免在循环中执行查询")             // 应该使用批量查询或预加载

	// 3. 批量操作建议
	// 批量操作可以显著提升大数据量处理的性能
	fmt.Println("\n3. 批量操作:")
	fmt.Println("   - 使用 CreateInBatches 进行批量插入") // 比循环插入效率高很多
	fmt.Println("   - 使用批量更新替代循环更新")              // 减少数据库交互次数
	fmt.Println("   - 合理设置批次大小")                  // 平衡内存使用和性能

	// 4. 连接池配置建议
	// 合理的连接池配置对高并发应用至关重要
	fmt.Println("\n4. 连接池配置:")
	fmt.Println("   - 设置合适的最大连接数") // 根据数据库和应用负载调整
	fmt.Println("   - 配置连接生存时间")   // 避免长时间占用连接
	fmt.Println("   - 监控连接池使用情况")  // 及时发现连接泄漏等问题

	// 5. 其他优化建议
	// 涵盖预编译、事务、关联查询等方面的优化
	fmt.Println("\n5. 其他优化:")
	fmt.Println("   - 启用预编译语句")    // 提升重复查询的性能
	fmt.Println("   - 合理使用事务")     // 保证数据一致性的同时避免长事务
	fmt.Println("   - 避免不必要的关联查询") // 只加载真正需要的关联数据
	fmt.Println("   - 使用适当的日志级别")  // 生产环境避免过多的SQL日志输出
}

// runPerformanceTests 运行性能测试的核心函数
// 参数db: 数据库连接实例
// 参数dbType: 数据库类型字符串，用于显示
// 该函数包含完整的性能测试流程，可被不同数据库类型复用
func runPerformanceTests(db *gorm.DB, dbType string) {
	fmt.Printf("\n=== %s 性能测试开始 ===\n", dbType)

	// 生成测试数据 - 创建足够的数据量来测试性能差异
	// 参数说明: 2000个用户, 5000篇文章, 10000条评论
	// 这个数据量足以体现不同查询策略的性能差异
	generateTestData(db, 2000, 5000, 10000)

	// 运行性能测试 - 按照从简单到复杂的顺序进行测试
	benchmarkBasicQueries(db)             // 测试基础查询操作的性能
	benchmarkPreloading(db)               // 测试预加载vs N+1查询的性能差异
	benchmarkBatchOperations(db)          // 测试批量操作vs单条操作的性能差异
	benchmarkComplexQueries(db)           // 测试复杂查询（聚合、子查询、连接）的性能
	benchmarkIndexOptimization(db)        // 测试索引对查询性能的影响
	benchmarkConnectionAndTransaction(db) // 测试连接池和事务的性能特征

	fmt.Printf("\n=== %s 性能测试完成 ===\n", dbType)
}

// demonstrateMySQL MySQL数据库演示函数
// 演示如何使用MySQL数据库进行性能测试，展示MySQL特有的性能特征
// 注意：运行前需要确保MySQL服务已启动，并且数据库连接配置正确
func demonstrateMySQL() {
	fmt.Println("\n=== MySQL 数据库演示 ===")
	fmt.Println("注意：请确保MySQL服务已启动，数据库连接配置正确")

	// 配置MySQL连接字符串
	// 格式：用户名:密码@tcp(主机:端口)/数据库名?参数
	// 重要参数说明：
	// - charset=utf8mb4: 支持完整的UTF-8字符集，包括emoji
	// - parseTime=True: 自动解析时间类型
	// - loc=Local: 使用本地时区
	// mysqlDSN := "root:123456@tcp(localhost:3306)/gorm_performance?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDSN := "root:fastbee@tcp(192.168.100.124:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	// 获取MySQL配置
	mysqlConfig := GetMySQLConfig(mysqlDSN)

	// 尝试连接MySQL数据库
	fmt.Println("正在连接MySQL数据库...")
	db, err := gorm.Open(mysql.Open(mysqlConfig.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlConfig.LogLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		fmt.Printf("MySQL连接失败: %v\n", err)
		fmt.Println("请检查以下配置：")
		fmt.Println("1. MySQL服务是否已启动")
		fmt.Println("2. 用户名和密码是否正确")
		fmt.Println("3. 数据库 'gorm_performance' 是否存在")
		fmt.Println("4. 网络连接是否正常")
		return
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("获取数据库实例失败: %v\n", err)
		return
	}

	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(mysqlConfig.MaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("MySQL连接测试失败: %v\n", err)
		return
	}

	fmt.Println("✓ MySQL数据库连接成功")

	// 执行自动迁移
	if err := AutoMigrate(db); err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		return
	}

	// 创建复合索引
	createCompositeIndexes(db)
	fmt.Println("✓ 数据库结构初始化完成")

	// 运行性能测试
	runPerformanceTests(db, "MySQL")

	// MySQL特有的性能提示
	fmt.Println("\n=== MySQL 性能优化建议 ===")
	fmt.Println("1. 索引优化：")
	fmt.Println("   - 使用 EXPLAIN 分析查询执行计划")
	fmt.Println("   - 合理设计复合索引，注意索引顺序")
	fmt.Println("   - 避免在WHERE子句中使用函数")
	fmt.Println("2. 查询优化：")
	fmt.Println("   - 使用LIMIT限制结果集大小")
	fmt.Println("   - 避免SELECT *，只查询需要的字段")
	fmt.Println("   - 合理使用JOIN，避免笛卡尔积")
	fmt.Println("3. 配置优化：")
	fmt.Println("   - 调整innodb_buffer_pool_size")
	fmt.Println("   - 优化连接池配置")
	fmt.Println("   - 启用查询缓存（适当情况下）")

	// 关闭数据库连接
	sqlDB.Close()
	fmt.Println("✓ MySQL数据库连接已关闭")
}

// main 主函数演示
// 该函数是程序的入口点，提供SQLite和MySQL两种数据库的性能测试选项
// 用户可以选择使用SQLite进行快速测试，或使用MySQL进行生产环境测试
func main() {
	fmt.Println("=== GORM Level 5 性能优化练习 ===")
	fmt.Println("\n请选择要使用的数据库类型：")
	fmt.Println("1. SQLite (默认，快速测试)")
	fmt.Println("2. MySQL (生产环境测试)")
	fmt.Println("3. 两种数据库对比测试")
	choice := "2"
	// // 读取用户输入
	// var choice string
	// fmt.Scanln(&choice)

	// // 默认选择SQLite
	// if choice == "" {
	// 	choice = "1"
	// }

	switch choice {
	case "1":
		// SQLite 演示
		fmt.Println("\n=== 使用 SQLite 数据库 ===")
		// 初始化数据库 - 建立连接、创建表结构、设置索引
		// 使用SQLite作为默认数据库，便于快速测试和学习
		db := initDB()
		fmt.Println("✓ SQLite数据库初始化完成")

		// 运行性能测试
		runPerformanceTests(db, "SQLite")

	case "2":
		// MySQL 演示
		demonstrateMySQL()

	case "3":
		// 对比测试
		fmt.Println("\n=== 数据库性能对比测试 ===")

		// 先测试SQLite
		fmt.Println("\n--- SQLite 测试 ---")
		db1 := initDB()
		fmt.Println("✓ SQLite数据库初始化完成")
		runPerformanceTests(db1, "SQLite")

		// 再测试MySQL
		fmt.Println("\n--- MySQL 测试 ---")
		demonstrateMySQL()

		// 对比总结
		fmt.Println("\n=== 数据库性能对比总结 ===")
		fmt.Println("SQLite 特点：")
		fmt.Println("  ✓ 轻量级，无需安装服务")
		fmt.Println("  ✓ 适合开发和测试环境")
		fmt.Println("  ✓ 单文件存储，便于部署")
		fmt.Println("  ✗ 并发写入能力有限")
		fmt.Println("  ✗ 不适合大型应用")

		fmt.Println("\nMySQL 特点：")
		fmt.Println("  ✓ 高并发处理能力")
		fmt.Println("  ✓ 丰富的存储引擎")
		fmt.Println("  ✓ 完善的事务支持")
		fmt.Println("  ✓ 适合生产环境")
		fmt.Println("  ✗ 需要安装和配置服务")
		fmt.Println("  ✗ 资源消耗相对较大")

	default:
		fmt.Println("无效选择，使用默认SQLite数据库")
		db := initDB()
		fmt.Println("✓ SQLite数据库初始化完成")
		runPerformanceTests(db, "SQLite")
	}

	// 显示优化建议 - 总结性能优化的最佳实践
	// 为开发者提供实用的性能优化指导
	showOptimizationTips()

	fmt.Println("\n=== Level 5 性能优化练习完成 ===")
}

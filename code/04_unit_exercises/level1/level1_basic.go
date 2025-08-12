// 本文件实现了GORM的基础练习示例,包含数据库连接、模型定义、CRUD操作等功能

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 练习1：数据库连接和基本配置

// DatabaseConfig 数据库配置结构
// DSN: 数据源名称,用于指定数据库连接字符串
// MaxOpenConns: 最大打开连接数
// MaxIdleConns: 最大空闲连接数
// MaxLifetime: 连接最大生命周期
// LogLevel: 日志级别
type DatabaseConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	LogLevel     logger.LogLevel
}

// GetDefaultConfig 获取默认配置
// 返回一个包含默认参数的数据库配置对象:
// - 使用SQLite数据库,文件名为test.db
// - 最大连接数10
// - 最大空闲连接5
// - 连接生命周期1小时
// - 日志级别为Info
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DSN:          "test.db",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		MaxLifetime:  time.Hour,
		LogLevel:     logger.Info,
	}
}

// InitDatabase 初始化数据库连接
// 参数config: 数据库配置对象
// 返回:
// - *gorm.DB: GORM数据库连接对象
// - error: 错误信息
func InitDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	// 配置GORM日志
	// - 设置慢查询阈值为1秒
	// - 忽略记录未找到的错误
	// - 启用参数化查询
	// - 关闭彩色输出
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  config.LogLevel,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	// 打开数据库连接
	// 配置:
	// - 禁用外键约束
	// - 表名前缀为t_
	// - 使用复数表名
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}

// TestConnection 测试数据库连接
// 通过Ping()方法测试数据库连接是否正常
func TestConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// GetConnectionStats 获取连接池统计信息
// 返回包含以下统计数据的map:
// - 最大连接数
// - 当前打开连接数
// - 正在使用的连接数
// - 空闲连接数
// - 等待队列长度
// - 等待总时长
// - 因最大空闲连接关闭的连接数
// - 因空闲超时关闭的连接数
// - 因超过生命周期关闭的连接数
func GetConnectionStats(db *gorm.DB) map[string]interface{} {
	sqlDB, err := db.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// 练习2：基本模型定义

// BaseModel 基础模型
// 包含所有模型通用的字段:
// - ID: 主键
// - CreatedAt: 创建时间
// - UpdatedAt: 更新时间
// - DeletedAt: 删除时间(用于软删除)
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// User 用户模型
// 字段说明:
// - Username: 用户名,唯一索引,最大长度50
// - Email: 邮箱,唯一索引,最大长度100
// - Password: 密码,最大长度255,JSON序列化时忽略
// - Age: 年龄,取值范围0-150
// - IsActive: 是否活跃,默认true
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Age      int    `gorm:"check:age >= 0 AND age <= 150" json:"age"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// TableName 自定义表名
// 返回用户表的实际表名
func (User) TableName() string {
	return "users"
}

// Product 商品模型
// 字段说明:
// - Name: 商品名称,建立索引,最大长度200
// - Description: 商品描述,文本类型
// - Price: 价格,精度10位,小数2位
// - Stock: 库存,默认0,非负数
// - SKU: 商品编码,唯一索引,最大长度50
// - Status: 商品状态,默认active,限制可选值
type Product struct {
	BaseModel
	Name        string  `gorm:"size:200;not null;index" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"precision:10;scale:2;not null;check:price >= 0" json:"price"`
	Stock       int     `gorm:"not null;default:0;check:stock >= 0" json:"stock"`
	SKU         string  `gorm:"uniqueIndex;size:50" json:"sku"`
	Status      string  `gorm:"size:20;default:'active';check:status IN ('active','inactive','discontinued')" json:"status"`
}

// Order 订单模型
// 字段说明:
// - OrderNo: 订单号,唯一索引
// - UserID: 用户ID,外键关联用户表
// - TotalPrice: 订单总价
// - Status: 订单状态,默认pending
// - OrderDate: 下单时间
// - User: 关联的用户信息
type Order struct {
	BaseModel
	OrderNo    string    `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	TotalPrice float64   `gorm:"precision:10;scale:2;not null" json:"total_price"`
	Status     string    `gorm:"size:20;default:'pending'" json:"status"`
	OrderDate  time.Time `gorm:"not null;index" json:"order_date"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// 练习3：CRUD基本操作

// CreateUser 创建用户
// 参数:
// - db: 数据库连接
// - username: 用户名
// - email: 邮箱
// - password: 密码
// - age: 年龄
// 返回:
// - 创建的用户对象
// - 错误信息
func CreateUser(db *gorm.DB, username, email, password string, age int) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
		Age:      age,
		IsActive: true,
	}

	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("创建用户成功，ID: %d, 影响行数: %d\n", user.ID, result.RowsAffected)
	return user, nil
}

// GetUserByID 根据ID获取用户
// 参数:
// - db: 数据库连接
// - id: 用户ID
// 返回:
// - 用户对象
// - 错误信息
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
// 参数:
// - db: 数据库连接
// - email: 用户邮箱
// 返回:
// - 用户对象
// - 错误信息
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
// 参数:
// - db: 数据库连接
// - id: 用户ID
// - updates: 需要更新的字段map
func UpdateUser(db *gorm.DB, id uint, updates map[string]interface{}) error {
	result := db.Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("更新用户成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// DeleteUser 删除用户（软删除）
// 参数:
// - db: 数据库连接
// - id: 用户ID
func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("删除用户成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// GetAllUsers 获取所有用户（分页）
// 参数:
// - db: 数据库连接
// - page: 页码
// - pageSize: 每页记录数
// 返回:
// - 用户列表
// - 总记录数
// - 错误信息
func GetAllUsers(db *gorm.DB, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	// 计算总数
	db.Model(&User{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	result := db.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// 练习4：查询操作

// FindUsersByAge 根据年龄范围查找用户
// 参数:
// - db: 数据库连接
// - minAge: 最小年龄
// - maxAge: 最大年龄
func FindUsersByAge(db *gorm.DB, minAge, maxAge int) ([]User, error) {
	var users []User
	result := db.Where("age BETWEEN ? AND ?", minAge, maxAge).Find(&users)
	return users, result.Error
}

// FindActiveUsers 查找活跃用户
// 参数:
// - db: 数据库连接
func FindActiveUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Where("is_active = ?", true).Find(&users)
	return users, result.Error
}

// SearchUsersByUsername 根据用户名搜索（模糊查询）
// 参数:
// - db: 数据库连接
// - keyword: 搜索关键词
func SearchUsersByUsername(db *gorm.DB, keyword string) ([]User, error) {
	var users []User
	result := db.Where("username LIKE ?", "%"+keyword+"%").Find(&users)
	return users, result.Error
}

// CountUsersByAge 统计不同年龄段的用户数量
// 返回各年龄段的用户数量统计:
// - under_18: 18岁以下
// - 18_30: 18-30岁
// - 31_50: 31-50岁
// - over_50: 50岁以上
func CountUsersByAge(db *gorm.DB) (map[string]int64, error) {
	type AgeGroup struct {
		AgeRange string
		Count    int64
	}

	var results []AgeGroup
	err := db.Model(&User{}).Select(
		"CASE WHEN age < 18 THEN 'under_18' WHEN age BETWEEN 18 AND 30 THEN '18_30' WHEN age BETWEEN 31 AND 50 THEN '31_50' ELSE 'over_50' END as age_range, COUNT(*) as count",
	).Group("age_range").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.AgeRange] = result.Count
	}

	return counts, nil
}

// 练习5：数据库迁移

// AutoMigrate 自动迁移
// 自动创建或更新数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Product{}, &Order{})
}

// CreateIndexes 创建索引
// 创建以下索引:
// - 用户表的username和email复合索引
// - 用户表is_active字段的条件索引
func CreateIndexes(db *gorm.DB) error {
	// 创建复合索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username_email ON users(username, email)").Error; err != nil {
		return err
	}

	// 创建条件索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active) WHERE is_active = true").Error; err != nil {
		return err
	}

	return nil
}

// main函数
// 演示所有功能的使用方法
func main() {
	fmt.Println("=== GORM Level 1 基础练习 ===")

	// 1. 初始化数据库
	config := GetDefaultConfig()
	db, err := InitDatabase(config)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 测试连接
	if err := TestConnection(db); err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}
	fmt.Println("✓ 数据库连接成功")

	// 2. 自动迁移
	if err := AutoMigrate(db); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	fmt.Println("✓ 数据库迁移完成")

	// 3. 创建索引
	if err := CreateIndexes(db); err != nil {
		log.Printf("创建索引失败: %v", err)
	} else {
		fmt.Println("✓ 索引创建完成")
	}

	// 4. CRUD操作演示
	fmt.Println("\n=== CRUD操作演示 ===")

	// 创建用户
	user1, err := CreateUser(db, "alice", "alice@example.com", "password123", 25)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// user2, err := CreateUser(db, "bob", "bob@example.com", "password456", 30)
	// if err != nil {
	// 	fmt.Printf("创建用户失败: %v\n", err)
	// }

	// 查询用户
	if user1 != nil {
		fetchedUser, err := GetUserByID(db, user1.ID)
		if err != nil {
			fmt.Printf("查询用户失败: %v\n", err)
		} else {
			fmt.Printf("查询到用户: %+v\n", fetchedUser)
		}
	}

	// 更新用户
	if user1 != nil {
		updates := map[string]interface{}{
			"age":       26,
			"is_active": false,
		}
		if err := UpdateUser(db, user1.ID, updates); err != nil {
			fmt.Printf("更新用户失败: %v\n", err)
		}
	}

	// 5. 查询操作演示
	fmt.Println("\n=== 查询操作演示 ===")

	// 年龄范围查询
	users, err := FindUsersByAge(db, 20, 35)
	if err != nil {
		fmt.Printf("年龄范围查询失败: %v\n", err)
	} else {
		fmt.Printf("年龄在20-35之间的用户数量: %d\n", len(users))
	}

	// 活跃用户查询
	activeUsers, err := FindActiveUsers(db)
	if err != nil {
		fmt.Printf("活跃用户查询失败: %v\n", err)
	} else {
		fmt.Printf("活跃用户数量: %d\n", len(activeUsers))
	}

	// 用户名搜索
	searchResults, err := SearchUsersByUsername(db, "a")
	if err != nil {
		fmt.Printf("用户名搜索失败: %v\n", err)
	} else {
		fmt.Printf("包含字母'a'的用户数量: %d\n", len(searchResults))
	}

	// 年龄统计
	ageCounts, err := CountUsersByAge(db)
	if err != nil {
		fmt.Printf("年龄统计失败: %v\n", err)
	} else {
		fmt.Printf("年龄分布统计: %+v\n", ageCounts)
	}

	// 6. 分页查询
	allUsers, total, err := GetAllUsers(db, 1, 10)
	if err != nil {
		fmt.Printf("分页查询失败: %v\n", err)
	} else {
		fmt.Printf("总用户数: %d, 当前页用户数: %d\n", total, len(allUsers))
	}

	// 7. 连接池统计
	fmt.Println("\n=== 连接池统计 ===")
	stats := GetConnectionStats(db)
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	fmt.Println("\n=== Level 1 练习完成 ===")
}

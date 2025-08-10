// 04_unit_exercises/level1_basic.go - Level 1 基础练习
// 对应文档：03_GORM单元练习_基础技能训练.md

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 练习1：数据库连接和基本配置

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	LogLevel     logger.LogLevel
}

// GetDefaultConfig 获取默认配置
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
func InitDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	// 配置GORM日志
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
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: newLogger,
		// 禁用外键约束（SQLite特定）
		DisableForeignKeyConstraintWhenMigrating: true,
		// 命名策略
		NamingStrategy: gorm.NamingStrategy{
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
func TestConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// GetConnectionStats 获取连接池统计信息
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
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Age      int    `gorm:"check:age >= 0 AND age <= 150" json:"age"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

// Product 商品模型
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
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(db *gorm.DB, id uint, updates map[string]interface{}) error {
	result := db.Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("更新用户成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// DeleteUser 删除用户（软删除）
func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("删除用户成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// GetAllUsers 获取所有用户（分页）
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
func FindUsersByAge(db *gorm.DB, minAge, maxAge int) ([]User, error) {
	var users []User
	result := db.Where("age BETWEEN ? AND ?", minAge, maxAge).Find(&users)
	return users, result.Error
}

// FindActiveUsers 查找活跃用户
func FindActiveUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Where("is_active = ?", true).Find(&users)
	return users, result.Error
}

// SearchUsersByUsername 根据用户名搜索（模糊查询）
func SearchUsersByUsername(db *gorm.DB, keyword string) ([]User, error) {
	var users []User
	result := db.Where("username LIKE ?", "%"+keyword+"%").Find(&users)
	return users, result.Error
}

// CountUsersByAge 统计不同年龄段的用户数量
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
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Product{}, &Order{})
}

// CreateIndexes 创建索引
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

// 主函数演示
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

	user2, err := CreateUser(db, "bob", "bob@example.com", "password456", 30)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

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
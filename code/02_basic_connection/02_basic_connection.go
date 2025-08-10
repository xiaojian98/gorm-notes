// 02_basic_connection.go - GORM基础连接示例
// 对应文档：01_GORM入门指南_5W1H详解.md

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 用户模型 - 演示基本的模型定义
type User struct {
	// ID 主键字段,使用uint类型自增
	ID uint `json:"id" gorm:"primarykey"`
	// Name 用户名字段,限制长度100,不允许为空
	Name string `json:"name" gorm:"size:100;not null"`
	// Email 邮箱字段,唯一索引,限制长度100
	Email string `json:"email" gorm:"uniqueIndex;size:100"`
	// Age 年龄字段
	Age int `json:"age"`
	// Status 状态字段,默认值为active
	Status string `json:"status" gorm:"default:active"`
	// CreatedAt 创建时间,GORM会自动维护
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间,GORM会自动维护
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt 软删除时间,添加了索引
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	// 配置GORM
	db, err := gorm.Open(sqlite.Open("basic_example.db"), &gorm.Config{
		// 设置日志级别
		Logger: logger.Default.LogMode(logger.Info),
		// 启用预编译语句
		PrepareStmt: true,
		// 禁用外键约束（SQLite默认禁用）
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// demonstrateBasicOperations 演示基本的CRUD操作
func demonstrateBasicOperations(db *gorm.DB) {
	fmt.Println("\n=== 基本CRUD操作演示 ===")

	// 1. 创建用户
	fmt.Println("\n1. 创建用户：")
	user := User{
		Name:   "张三",
		Email:  "zhangsan@example.com",
		Age:    25,
		Status: "active",
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.Printf("创建用户失败: %v", result.Error)
	} else {
		fmt.Printf("创建用户成功，ID: %d, 影响行数: %d\n", user.ID, result.RowsAffected)
	}

	// 2. 查询用户
	fmt.Println("\n2. 查询用户：")
	var foundUser User
	db.First(&foundUser, user.ID)
	fmt.Printf("查询到用户: %+v\n", foundUser)

	// 3. 更新用户
	fmt.Println("\n3. 更新用户：")
	db.Model(&foundUser).Update("age", 26)
	fmt.Printf("更新后用户年龄: %d\n", foundUser.Age)

	// 4. 条件查询
	fmt.Println("\n4. 条件查询：")
	var users []User
	db.Where("age > ?", 20).Find(&users)
	fmt.Printf("年龄大于20的用户数量: %d\n", len(users))

	// 5. 软删除
	// 其中 “软删除”（Soft Delete）的实现依赖于一个关键机制：
	// 模型中必须包含一个 gorm.DeletedAt 字段。
	// 最常见的方式是 嵌入 gorm.Model
	fmt.Println("\n5. 软删除：")
	db.Delete(&foundUser)
	fmt.Println("用户已软删除")

	// // 5. 真实删除
	// fmt.Println("\n5. 真实删除：")
	// // 使用Unscoped()方法来执行真实删除
	// db.Unscoped().Delete(&foundUser)
	// fmt.Println("用户已被真实删除")

	// 6. 查询未删除的用户
	var activeUsers []User
	db.Find(&activeUsers)
	fmt.Printf("活跃用户数量: %d\n", len(activeUsers))

	// 7. 查询包含软删除的用户
	var allUsers []User
	db.Unscoped().Find(&allUsers)
	fmt.Printf("所有用户数量（包含软删除）: %d\n", len(allUsers))

}

// demonstrateChainableAPI 演示链式API
func demonstrateChainableAPI(db *gorm.DB) {
	fmt.Println("\n=== 链式API演示 ===")

	// 创建测试数据
	testUsers := []User{
		{Name: "李四", Email: "lisi@example.com", Age: 30, Status: "active"},
		{Name: "王五", Email: "wangwu@example.com", Age: 28, Status: "inactive"},
		{Name: "赵六", Email: "zhaoliu@example.com", Age: 35, Status: "active"},
	}

	db.Create(&testUsers)

	// 链式查询示例
	fmt.Println("\n1. 链式查询 - 活跃用户，按年龄排序：")
	var activeUsers []User
	db.Where("status = ?", "active").
		Order("age desc").
		Limit(10).
		Find(&activeUsers)

	for _, user := range activeUsers {
		fmt.Printf("  %s (年龄: %d)\n", user.Name, user.Age)
	}

	// 聚合查询
	/**
	📊 二、聚合查询（Aggregate Query）
		✅ 什么是聚合查询？
		聚合查询 是指对一组数据进行“统计计算”，而不是获取原始数据行。常见的聚合函数有：
		函数	说明
		COUNT	统计行数
		SUM	求和
		AVG	平均值
		MAX	最大值
		MIN	最小值
	*/
	fmt.Println("\n2. 聚合查询：")
	var count int64
	var avgAge float64

	db.Model(&User{}).Where("status = ?", "active").Count(&count)
	db.Model(&User{}).Where("status = ?", "active").Select("AVG(age)").Scan(&avgAge)

	fmt.Printf("  活跃用户数量: %d\n", count)
	fmt.Printf("  平均年龄: %.2f\n", avgAge)

	// 批量更新
	fmt.Println("\n3. 批量更新：")
	result := db.Model(&User{}).Where("age > ?", 30).Update("status", "senior")
	fmt.Printf("  更新了 %d 个用户的状态\n", result.RowsAffected)
}

// demonstrateErrorHandling 演示错误处理
func demonstrateErrorHandling(db *gorm.DB) {
	fmt.Println("\n=== 错误处理演示 ===")

	// 1. 记录不存在的错误
	fmt.Println("\n1. 查询不存在的记录：")
	var user User
	err := db.First(&user, 99999).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("  记录不存在")
		} else {
			fmt.Printf("  查询错误: %v\n", err)
		}
	}

	// 2. 唯一约束冲突
	fmt.Println("\n2. 唯一约束冲突：")
	duplicateUser := User{
		Name:  "重复邮箱用户",
		Email: "lisi@example.com", // 使用已存在的邮箱
		Age:   25,
	}

	err = db.Create(&duplicateUser).Error
	if err != nil {
		fmt.Printf("  创建失败: %v\n", err)
	}

	// 3. 事务中的错误处理
	fmt.Println("\n3. 事务错误处理：")
	err = db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{Name: "事务用户", Email: "transaction@example.com", Age: 17}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 模拟业务逻辑错误
		if user.Age < 18 {
			return fmt.Errorf("用户年龄不能小于18岁")
		}

		return nil
	})

	if err != nil {
		fmt.Printf("  事务执行失败: %v\n", err)
	} else {
		fmt.Println("  事务执行成功")
	}
}

func main() {
	fmt.Println("🚀 GORM基础连接和操作示例")
	fmt.Println("对应文档：01_GORM入门指南_5W1H详解.md")

	// 初始化数据库
	db, err := InitDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	fmt.Println("✅ 数据库连接成功！")

	// 自动迁移:核心作用是：自动创建或更新数据库表结构，使其与 Go 的结构体定义保持一致。
	// “自动迁移”就是让 GORM 根据你的 Go 结构体（如 User），自动帮你创建或修改数据库表（如 users 表），省去手动写 SQL 的麻烦。
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	fmt.Println("✅ 数据库迁移完成！")

	// 演示基本操作
	demonstrateBasicOperations(db)
	fmt.Println("=============================================================================")
	// 演示链式API
	demonstrateChainableAPI(db)
	fmt.Println("=============================================================================")
	// 演示错误处理
	demonstrateErrorHandling(db)
	fmt.Println("=============================================================================")
	fmt.Println("\n🎉 示例运行完成！")
	fmt.Println("\n💡 学习要点：")
	fmt.Println("1. GORM的基本配置和连接")
	fmt.Println("2. 模型定义和标签使用")
	fmt.Println("3. 基本的CRUD操作")
	fmt.Println("4. 链式API的使用")
	fmt.Println("5. 错误处理的最佳实践")
	fmt.Println("6. 软删除的工作原理")

	// 清理资源
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

package config

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"blog-system-refactored/internal/models"
)

// DatabaseConfig 数据库配置结构体
// 支持SQLite和MySQL两种数据库类型
type DatabaseConfig struct {
	Type     string `json:"type"`     // 数据库类型: "sqlite" 或 "mysql"
	Host     string `json:"host"`     // 数据库主机地址
	Port     int    `json:"port"`     // 数据库端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	DBName   string `json:"dbname"`   // 数据库名称
	Charset  string `json:"charset"`  // 字符集
	FilePath string `json:"filepath"` // SQLite文件路径
}

// GetDefaultSQLiteConfig 获取默认的SQLite配置
// 返回一个预配置的SQLite数据库配置
func GetDefaultSQLiteConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:     "sqlite",
		FilePath: "blog_system.db",
	}
}

// GetDefaultMySQLConfig 获取默认的MySQL配置
// 返回一个预配置的MySQL数据库配置
func GetDefaultMySQLConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
		DBName:   "blog_system",
		Charset:  "utf8mb4",
	}
}

// ConnectDatabase 根据配置连接数据库
// 参数: config - 数据库配置
// 返回: *gorm.DB - GORM数据库实例, error - 错误信息
func (config *DatabaseConfig) ConnectDatabase() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch config.Type {
	case "sqlite":
		// 连接SQLite数据库
		db, err = gorm.Open(sqlite.Open(config.FilePath), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
		}
		log.Printf("✅ 成功连接SQLite数据库: %s", config.FilePath)

	case "mysql":
		// 构建MySQL连接字符串
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.DBName, config.Charset)
		
		// 连接MySQL数据库
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("连接MySQL数据库失败: %v", err)
		}
		log.Printf("✅ 成功连接MySQL数据库: %s@%s:%d/%s", config.Username, config.Host, config.Port, config.DBName)

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}

	return db, nil
}

// TestConnection 测试数据库连接
// 参数: db - GORM数据库实例
// 返回: error - 错误信息
func TestConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	log.Println("🔗 数据库连接测试成功")
	return nil
}

// AutoMigrate 自动迁移数据库表结构
// 参数: db - GORM数据库实例
// 返回: error - 错误信息
func AutoMigrate(db *gorm.DB) error {
	log.Println("🔄 开始数据库表结构迁移...")

	// 定义所有需要迁移的模型
	models := []interface{}{
		// 用户相关表
		&models.User{},
		&models.UserProfile{},
		&models.Follow{},

		// 内容相关表
		&models.Post{},
		&models.Category{},
		&models.Tag{},
		&models.PostMeta{},

		// 评论相关表
		&models.Comment{},
		&models.Like{},

		// 通知相关表
		&models.Notification{},

		// 分析统计表
		&models.Analytics{},
	}

	// 先删除所有表（如果存在）
	log.Println("🗑️ 清理现有表结构...")
	for i := len(models) - 1; i >= 0; i-- {
		if err := db.Migrator().DropTable(models[i]); err != nil {
			log.Printf("⚠️ 删除表失败（可能不存在）: %v", err)
		}
	}

	// 重新创建所有表
	log.Println("🔨 创建新的表结构...")
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("数据库表结构迁移失败: %v", err)
	}

	log.Println("✅ 数据库表结构迁移完成")
	return nil
}
// 03_blog_system/config/database.go - 数据库配置
// 对应文档：02_GORM背景示例_博客系统实战.md

package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB 全局数据库实例
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

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "192.168.100.124"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "123456"),
		DBName:   getEnv("DB_NAME", "blog_system"),
		Charset:  getEnv("DB_CHARSET", "utf8mb4"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	// 连接数据库
	// 为了演示方便，这里使用SQLite
	// 在生产环境中，建议使用MySQL或PostgreSQL
	DB, err = gorm.Open(sqlite.Open("blog_system.db"), &gorm.Config{
		// 配置日志
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // 慢查询阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound错误
				Colorful:                  true,        // 彩色打印
			},
		),
		// 启用预编译语句
		PrepareStmt: true,
		// 禁用外键约束（SQLite特有）
		DisableForeignKeyConstraintWhenMigrating: true,
		// 命名策略
		NamingStrategy: &CustomNamingStrategy{},
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 连接最大空闲时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ 数据库连接成功")
	return nil
}

// InitMySQLDB 初始化MySQL数据库连接（生产环境推荐）
func InitMySQLDB() error {
	// 注意：需要先安装MySQL驱动
	// go get -u gorm.io/driver/mysql

	config := GetDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		PrepareStmt:    true,
		NamingStrategy: &CustomNamingStrategy{},
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// return nil

	return fmt.Errorf("MySQL driver not imported, please uncomment the code above and install mysql driver")
}

// CustomNamingStrategy 自定义命名策略
type CustomNamingStrategy struct{}

// SchemaName 数据库模式命名策略
func (ns *CustomNamingStrategy) SchemaName(table string) string {
	return table
}

// TableName 表名命名策略
func (ns *CustomNamingStrategy) TableName(table string) string {
	return table
}

// ColumnName 列名命名策略
func (ns *CustomNamingStrategy) ColumnName(table, column string) string {
	return column
}

// JoinTableName 连接表命名策略
func (ns *CustomNamingStrategy) JoinTableName(str string) string {
	return str
}

// RelationshipFKName 外键命名策略
func (ns *CustomNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return rel.Name + "_id"
}

// CheckerName 检查器命名策略
func (ns *CustomNamingStrategy) CheckerName(table, column string) string {
	return "chk_" + table + "_" + column
}

// IndexName 索引命名策略
func (ns *CustomNamingStrategy) IndexName(table, column string) string {
	return "idx_" + table + "_" + column
}

// UniqueName 唯一约束命名策略
func (ns *CustomNamingStrategy) UniqueName(table, column string) string {
	return "uniq_" + table + "_" + column
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// IsConnected 检查数据库是否连接
func IsConnected() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	err = sqlDB.Ping()
	return err == nil
}

// GetConnectionStats 获取连接池统计信息
func GetConnectionStats() map[string]interface{} {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil
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

// clean_db.go - 清理数据库
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// cleanDatabase 清理数据库中的所有表
func cleanDatabase() {
	// 连接数据库
	dsn := "root:123456@tcp(192.168.100.124:3306)/blog_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	fmt.Println("开始清理数据库...")

	// 禁用外键检查
	result := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if result.Error != nil {
		log.Fatal("禁用外键检查失败:", result.Error)
	}
	fmt.Println("✅ 已禁用外键检查")

	// 获取所有表名
	type Table struct {
		TableName string `json:"table_name"`
	}

	var tables []Table
	result = db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'blog_system'").Scan(&tables)
	if result.Error != nil {
		log.Fatal("获取表列表失败:", result.Error)
	}

	// 删除所有表
	for _, table := range tables {
		fmt.Printf("删除表: %s\n", table.TableName)
		result = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table.TableName))
		if result.Error != nil {
			log.Printf("删除表 %s 失败: %v", table.TableName, result.Error)
		} else {
			fmt.Printf("✅ 已删除表: %s\n", table.TableName)
		}
	}

	// 重新启用外键检查
	result = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if result.Error != nil {
		log.Fatal("启用外键检查失败:", result.Error)
	}
	fmt.Println("✅ 已启用外键检查")

	fmt.Println("🎉 数据库清理完成！")
}

func main() {
	cleanDatabase()
}
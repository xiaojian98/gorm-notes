// check_db.go - 检查数据库状态
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// checkDatabase 检查数据库状态
func checkDatabase() {
	// 连接数据库
	dsn := "root:123456@tcp(localhost:3306)/blog_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 检查User表是否存在
	var exists bool
	result := db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'blog_system' AND table_name = 'User')").Scan(&exists)
	if result.Error != nil {
		log.Fatal("检查表存在性失败:", result.Error)
	}

	fmt.Printf("User表是否存在: %v\n", exists)

	if exists {
		// 检查User表的约束
		type Constraint struct {
			ConstraintName string `json:"constraint_name"`
			ConstraintType string `json:"constraint_type"`
		}

		var constraints []Constraint
		result = db.Raw(`
			SELECT 
				CONSTRAINT_NAME as constraint_name,
				CONSTRAINT_TYPE as constraint_type
			FROM information_schema.table_constraints 
			WHERE table_schema = 'blog_system' 
			AND table_name = 'User'
		`).Scan(&constraints)

		if result.Error != nil {
			log.Fatal("检查约束失败:", result.Error)
		}

		fmt.Println("\nUser表的约束:")
		for _, constraint := range constraints {
			fmt.Printf("- %s (%s)\n", constraint.ConstraintName, constraint.ConstraintType)
		}

		// 检查索引
		type Index struct {
			IndexName string `json:"index_name"`
			NonUnique int    `json:"non_unique"`
		}

		var indexes []Index
		result = db.Raw(`
			SELECT 
				index_name,
				non_unique
			FROM information_schema.statistics 
			WHERE table_schema = 'blog_system' 
			AND table_name = 'User'
			AND index_name != 'PRIMARY'
		`).Scan(&indexes)

		if result.Error != nil {
			log.Fatal("检查索引失败:", result.Error)
		}

		fmt.Println("\nUser表的索引:")
		for _, index := range indexes {
			uniqueStr := "唯一"
			if index.NonUnique == 1 {
				uniqueStr = "非唯一"
			}
			fmt.Printf("- %s (%s)\n", index.IndexName, uniqueStr)
		}
	}
}

func main() {
	checkDatabase()
}
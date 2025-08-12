package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Post 文章模型（简化版）
// 用于演示钩子函数的工作原理
type Post struct {
	ID          uint       `gorm:"primarykey"`
	Title       string     `gorm:"size:200;not null"`
	Content     string     `gorm:"type:text"`
	Status      string     `gorm:"size:20;default:draft"`
	PublishedAt *time.Time `gorm:"index"`
	CategoryID  *uint      `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Category 分类模型（简化版）
// 用于演示关联数据的更新
type Category struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"size:100;not null"`
	PostCount int       `gorm:"default:0"` // 文章数量统计
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate 文章创建前钩子
// 功能: 在文章插入数据库之前自动执行
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，创建操作会被取消
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("🪝 [BeforeCreate] 钩子被调用！文章标题: %s\n", p.Title)
	
	// 如果是发布状态且没有设置发布时间，则设置为当前时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
		fmt.Printf("✅ [BeforeCreate] 自动设置发布时间: %v\n", now.Format("2006-01-02 15:04:05"))
	}
	
	fmt.Printf("📝 [BeforeCreate] 数据验证通过，准备插入数据库\n")
	return nil
}

// AfterCreate 文章创建后钩子
// 功能: 在文章成功插入数据库后自动执行
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，整个事务会回滚
func (p *Post) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("🎉 [AfterCreate] 钩子被调用！文章ID: %d\n", p.ID)
	
	// 更新分类的文章数量
	if p.CategoryID != nil {
		result := tx.Model(&Category{}).Where("id = ?", *p.CategoryID).
			Update("post_count", gorm.Expr("post_count + ?", 1))
		
		if result.Error != nil {
			fmt.Printf("❌ [AfterCreate] 更新分类计数失败: %v\n", result.Error)
			return result.Error
		}
		
		fmt.Printf("📊 [AfterCreate] 更新分类计数成功，影响行数: %d\n", result.RowsAffected)
	}
	
	fmt.Printf("✨ [AfterCreate] 文章创建完成，所有后续处理已完成\n")
	return nil
}

// BeforeUpdate 文章更新前钩子
// 功能: 演示更新操作的钩子
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，更新操作会被取消
func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	fmt.Printf("🔄 [BeforeUpdate] 钩子被调用！文章ID: %d\n", p.ID)
	
	// 如果状态改为发布且没有发布时间，设置发布时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
		fmt.Printf("📅 [BeforeUpdate] 设置发布时间: %v\n", now.Format("2006-01-02 15:04:05"))
	}
	
	return nil
}

// AfterUpdate 文章更新后钩子
// 功能: 演示更新操作后的处理
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，整个事务会回滚
func (p *Post) AfterUpdate(tx *gorm.DB) error {
	fmt.Printf("🎯 [AfterUpdate] 钩子被调用！文章ID: %d\n", p.ID)
	fmt.Printf("📝 [AfterUpdate] 文章更新完成\n")
	return nil
}

// BeforeDelete 文章删除前钩子
// 功能: 演示删除操作前的处理
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，删除操作会被取消
func (p *Post) BeforeDelete(tx *gorm.DB) error {
	fmt.Printf("🗑️ [BeforeDelete] 钩子被调用！准备删除文章ID: %d\n", p.ID)
	return nil
}

// AfterDelete 文章删除后钩子
// 功能: 演示删除操作后的清理工作
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，整个事务会回滚
func (p *Post) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("🧹 [AfterDelete] 钩子被调用！文章已删除，ID: %d\n", p.ID)
	
	// 更新分类的文章数量（减1）
	if p.CategoryID != nil {
		result := tx.Model(&Category{}).Where("id = ?", *p.CategoryID).
			Update("post_count", gorm.Expr("post_count - ?", 1))
		
		if result.Error != nil {
			fmt.Printf("❌ [AfterDelete] 更新分类计数失败: %v\n", result.Error)
			return result.Error
		}
		
		fmt.Printf("📉 [AfterDelete] 分类计数已减1，影响行数: %d\n", result.RowsAffected)
	}
	
	return nil
}

// AfterFind 查找后钩子
// 功能: 演示查询操作后的处理
// 参数: tx - GORM事务对象
// 返回值: error - 如果返回错误，查询结果会被丢弃
func (p *Post) AfterFind(tx *gorm.DB) error {
	fmt.Printf("🔍 [AfterFind] 钩子被调用！找到文章: %s (ID: %d)\n", p.Title, p.ID)
	return nil
}

// 初始化数据库连接
// 功能: 创建数据库连接并配置日志
// 返回值: *gorm.DB - 数据库连接对象, error - 错误信息
func initDB() (*gorm.DB, error) {
	// 创建SQLite内存数据库（用于测试）
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	
	if err != nil {
		return nil, err
	}
	
	// 自动迁移表结构
	err = db.AutoMigrate(&Post{}, &Category{})
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

// 创建测试分类
// 功能: 创建一个测试用的分类
// 参数: db - 数据库连接
// 返回值: uint - 分类ID, error - 错误信息
func createTestCategory(db *gorm.DB) (uint, error) {
	category := &Category{
		Name:      "技术分享",
		PostCount: 0,
	}
	
	result := db.Create(category)
	if result.Error != nil {
		return 0, result.Error
	}
	
	fmt.Printf("📁 创建测试分类成功，ID: %d\n", category.ID)
	return category.ID, nil
}

// 演示创建操作的钩子
// 功能: 演示BeforeCreate和AfterCreate钩子的调用
// 参数: db - 数据库连接, categoryID - 分类ID
func demonstrateCreateHooks(db *gorm.DB, categoryID uint) {
	fmt.Println("\n" + "="*50)
	fmt.Println("🚀 演示创建操作的钩子函数")
	fmt.Println("="*50)
	
	// 创建一篇草稿文章
	draftPost := &Post{
		Title:      "我的第一篇草稿",
		Content:    "这是一篇草稿文章的内容...",
		Status:     "draft",
		CategoryID: &categoryID,
	}
	
	fmt.Println("📝 创建草稿文章...")
	result := db.Create(draftPost)
	if result.Error != nil {
		log.Printf("创建草稿失败: %v", result.Error)
		return
	}
	
	fmt.Println("\n" + "-"*30)
	
	// 创建一篇发布文章
	publishedPost := &Post{
		Title:      "我的第一篇发布文章",
		Content:    "这是一篇已发布文章的内容...",
		Status:     "published", // 注意：没有设置PublishedAt
		CategoryID: &categoryID,
	}
	
	fmt.Println("📰 创建发布文章...")
	result = db.Create(publishedPost)
	if result.Error != nil {
		log.Printf("创建发布文章失败: %v", result.Error)
		return
	}
	
	// 查看分类的文章数量
	var category Category
	db.First(&category, categoryID)
	fmt.Printf("\n📊 分类 '%s' 的文章数量: %d\n", category.Name, category.PostCount)
}

// 演示更新操作的钩子
// 功能: 演示BeforeUpdate和AfterUpdate钩子的调用
// 参数: db - 数据库连接
func demonstrateUpdateHooks(db *gorm.DB) {
	fmt.Println("\n" + "="*50)
	fmt.Println("🔄 演示更新操作的钩子函数")
	fmt.Println("="*50)
	
	// 查找第一篇草稿文章
	var post Post
	result := db.Where("status = ?", "draft").First(&post)
	if result.Error != nil {
		log.Printf("查找草稿失败: %v", result.Error)
		return
	}
	
	fmt.Printf("📝 找到草稿文章: %s (ID: %d)\n", post.Title, post.ID)
	fmt.Println("🔄 将草稿改为发布状态...")
	
	// 更新文章状态为发布
	post.Status = "published"
	result = db.Save(&post)
	if result.Error != nil {
		log.Printf("更新文章失败: %v", result.Error)
		return
	}
	
	fmt.Printf("✅ 文章发布时间: %v\n", post.PublishedAt.Format("2006-01-02 15:04:05"))
}

// 演示查找操作的钩子
// 功能: 演示AfterFind钩子的调用
// 参数: db - 数据库连接
func demonstrateFindHooks(db *gorm.DB) {
	fmt.Println("\n" + "="*50)
	fmt.Println("🔍 演示查找操作的钩子函数")
	fmt.Println("="*50)
	
	// 查找所有发布的文章
	var posts []Post
	result := db.Where("status = ?", "published").Find(&posts)
	if result.Error != nil {
		log.Printf("查找文章失败: %v", result.Error)
		return
	}
	
	fmt.Printf("📚 找到 %d 篇发布的文章\n", len(posts))
}

// 演示删除操作的钩子
// 功能: 演示BeforeDelete和AfterDelete钩子的调用
// 参数: db - 数据库连接
func demonstrateDeleteHooks(db *gorm.DB) {
	fmt.Println("\n" + "="*50)
	fmt.Println("🗑️ 演示删除操作的钩子函数")
	fmt.Println("="*50)
	
	// 查找第一篇文章
	var post Post
	result := db.First(&post)
	if result.Error != nil {
		log.Printf("查找文章失败: %v", result.Error)
		return
	}
	
	fmt.Printf("📝 准备删除文章: %s (ID: %d)\n", post.Title, post.ID)
	
	// 删除文章
	result = db.Delete(&post)
	if result.Error != nil {
		log.Printf("删除文章失败: %v", result.Error)
		return
	}
	
	// 查看分类的文章数量变化
	if post.CategoryID != nil {
		var category Category
		db.First(&category, *post.CategoryID)
		fmt.Printf("📊 删除后分类 '%s' 的文章数量: %d\n", category.Name, category.PostCount)
	}
}

// 主函数
// 功能: 程序入口，演示所有钩子函数的工作原理
func main() {
	fmt.Println("🎯 GORM钩子函数演示程序")
	fmt.Println("这个程序将演示GORM钩子函数是如何自动被调用的")
	
	// 初始化数据库
	db, err := initDB()
	if err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	
	// 创建测试分类
	categoryID, err := createTestCategory(db)
	if err != nil {
		log.Fatal("创建测试分类失败:", err)
	}
	
	// 演示各种钩子函数
	demonstrateCreateHooks(db, categoryID)  // 创建操作钩子
	demonstrateUpdateHooks(db)              // 更新操作钩子
	demonstrateFindHooks(db)               // 查找操作钩子
	demonstrateDeleteHooks(db)              // 删除操作钩子
	
	fmt.Println("\n" + "="*50)
	fmt.Println("🎉 演示完成！")
	fmt.Println("通过这个演示，你可以看到：")
	fmt.Println("1. 钩子函数是由GORM自动调用的")
	fmt.Println("2. 不需要手动调用这些钩子函数")
	fmt.Println("3. 钩子函数在特定的数据库操作阶段执行")
	fmt.Println("4. 钩子函数可以修改数据、验证数据、处理关联等")
	fmt.Println("="*50)
}

/*
运行这个程序的步骤：

1. 确保已安装GORM和SQLite驱动：
   go mod init hook-demo
   go get gorm.io/gorm
   go get gorm.io/driver/sqlite

2. 运行程序：
   go run 钩子函数测试示例.go

3. 观察输出，你会看到：
   - 每个数据库操作都会触发相应的钩子函数
   - 钩子函数的执行顺序和时机
   - 钩子函数如何处理业务逻辑

预期输出示例：
🎯 GORM钩子函数演示程序
📁 创建测试分类成功，ID: 1
==================================================
🚀 演示创建操作的钩子函数
==================================================
📝 创建草稿文章...
🪝 [BeforeCreate] 钩子被调用！文章标题: 我的第一篇草稿
📝 [BeforeCreate] 数据验证通过，准备插入数据库
🎉 [AfterCreate] 钩子被调用！文章ID: 1
📊 [AfterCreate] 更新分类计数成功，影响行数: 1
✨ [AfterCreate] 文章创建完成，所有后续处理已完成
...
*/
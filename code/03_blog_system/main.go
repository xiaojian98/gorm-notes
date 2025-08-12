// 03_blog_system/main.go - 博客系统主程序
// 对应文档：02_GORM背景示例_博客系统实战.md

package main

import (
	"log"

	"blog-system/config"
	"blog-system/migrations"
	"blog-system/models"
	"blog-system/routes"
	"blog-system/services"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 启动博客系统...")

	// 初始化数据库
	if err := config.InitMySQLDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	log.Println("✅ 数据库连接成功")

	// 运行数据库迁移
	if err := migrations.RunMigrations(config.DB); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 初始化服务
	services.InitServices(config.DB)
	log.Println("✅ 服务初始化完成")

	// // 创建测试数据
	// if err := createTestData(); err != nil {
	// 	log.Printf("⚠️ 创建测试数据失败: %v", err)
	// } else {
	// 	log.Println("✅ 测试数据创建完成")
	// }

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化路由
	r := routes.SetupRoutes()

	log.Println("🌟 博客系统启动成功！")
	log.Println("📖 API文档: http://localhost:8080/api/docs")
	log.Println("🔗 测试接口:")
	log.Println("   GET  /api/posts - 获取文章列表")
	log.Println("   POST /api/users/register - 用户注册")
	log.Println("   POST /api/users/login - 用户登录")

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// createTestData 创建测试数据
func createTestData() error {
	log.Println("📝 创建测试数据...")

	// 创建分类
	categories := []models.Category{
		{Name: "技术分享", Description: "技术相关的文章", Slug: "tech"},
		{Name: "生活随笔", Description: "生活感悟和随笔", Slug: "life"},
		{Name: "学习笔记", Description: "学习过程中的笔记", Slug: "study"},
	}

	for _, category := range categories {
		var existingCategory models.Category
		// 检查分类是否已存在
		if err := config.DB.Where("slug = ?", category.Slug).First(&existingCategory).Error; err != nil {
			if err := config.DB.Create(&category).Error; err != nil {
				return err
			}
		}
	}

	// 创建标签
	tags := []models.Tag{
		{Name: "Go语言", Slug: "golang"},
		{Name: "数据库", Slug: "database"},
		{Name: "Web开发", Slug: "web-dev"},
		{Name: "GORM", Slug: "gorm"},
		{Name: "教程", Slug: "tutorial"},
	}

	for _, tag := range tags {
		var existingTag models.Tag
		if err := config.DB.Where("slug = ?", tag.Slug).First(&existingTag).Error; err != nil {
			if err := config.DB.Create(&tag).Error; err != nil {
				return err
			}
		}
	}

	// 创建测试用户
	testUser := models.User{
		Username: "admin",
		Email:    "admin@blog.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Nickname: "管理员",
		Status:   "active",
	}

	var existingUser models.User
	if err := config.DB.Where("username = ?", testUser.Username).First(&existingUser).Error; err != nil {
		if err := config.DB.Create(&testUser).Error; err != nil {
			return err
		}

		// 创建用户资料
		profile := models.Profile{
			UserID:   testUser.ID,
			Bio:      "这是一个测试用户的个人简介",
			Website:  "https://blog.example.com",
			Location: "北京",
		}
		config.DB.Create(&profile)

		// 创建示例文章
		var techCategory models.Category
		config.DB.Where("slug = ?", "tech").First(&techCategory)

		var gormTag, tutorialTag models.Tag
		config.DB.Where("slug = ?", "gorm").First(&gormTag)
		config.DB.Where("slug = ?", "tutorial").First(&tutorialTag)

		post := models.Post{
			Title:      "GORM入门教程：从零开始学习Go语言ORM",
			Slug:       "gorm-tutorial-for-beginners",
			Content:    "这是一篇关于GORM的详细教程，将带你从零开始学习Go语言中最流行的ORM框架...",
			Excerpt:    "GORM是Go语言中最受欢迎的ORM库，本文将详细介绍其基本用法和高级特性。",
			UserID:     testUser.ID,
			CategoryID: &techCategory.ID,
			Status:     "published",
			ViewCount:  156,
			Tags:       []models.Tag{gormTag, tutorialTag},
		}

		if err := config.DB.Create(&post).Error; err != nil {
			return err
		}

		// 创建示例评论
		comment := models.Comment{
			PostID:  post.ID,
			UserID:  testUser.ID,
			Content: "这篇文章写得很好，对GORM的介绍很详细！",
			Status:  "approved",
		}
		config.DB.Create(&comment)

		// 创建点赞记录
		like := models.Like{
			UserID:     testUser.ID,
			TargetID:   post.ID,
			TargetType: "post",
		}
		config.DB.Create(&like)
	}

	return nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"blog-system-refactored/internal/config"
	"blog-system-refactored/internal/handlers"
	"blog-system-refactored/internal/repository"
	"blog-system-refactored/internal/routes"
	"blog-system-refactored/internal/services"
)

// main 主函数
// 程序入口点，负责初始化各个组件并启动HTTP服务器
func main() {
	// 加载配置
	cfg := config.GetDefaultConfig()

	// 初始化数据库连接
	db, err := cfg.Database.ConnectDatabase()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 测试数据库连接
	if err := config.TestConnection(db); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 自动迁移数据库表结构
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化Repository层
	_ = repository.NewUserRepository(db)
	_ = repository.NewPostRepository(db)
	_ = repository.NewCommentRepository(db)
	_ = repository.NewAnalyticsRepository(db)

	// 初始化Service层
	userService := services.NewUserService(db)
	postService := services.NewPostService(db)
	commentService := services.NewCommentService(db)
	analyticsService := services.NewAnalyticsService(db)

	// 初始化Handler层
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// 设置Gin模式
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin路由器
	r := gin.New()

	// 设置路由
	routes.SetupRoutes(r, userHandler, postHandler, commentHandler, analyticsHandler)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 博客系统启动成功，监听端口: %d", cfg.Server.Port)
		log.Printf("📖 API文档地址: http://localhost:%d/docs", cfg.Server.Port)
		log.Printf("💚 健康检查地址: http://localhost:%d/health", cfg.Server.Port)
		log.Printf("🌍 环境: %s", cfg.App.Environment)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	gracefulShutdown(srv, db)
}

// gracefulShutdown 优雅关闭服务器
// 参数: srv - HTTP服务器, db - 数据库连接
// 功能: 监听系统信号，优雅地关闭服务器和数据库连接
func gracefulShutdown(srv *http.Server, db *gorm.DB) {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	
	// 监听系统信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// 等待信号
	<-quit
	log.Println("🛑 正在关闭服务器...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ 服务器关闭失败: %v", err)
	} else {
		log.Println("✅ HTTP服务器已关闭")
	}

	// 关闭数据库连接
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("❌ 数据库连接关闭失败: %v", err)
		} else {
			log.Println("✅ 数据库连接已关闭")
		}
	}

	log.Println("🎉 服务器已优雅关闭")
}
// 03_blog_system/routes/routes.go - 路由配置
// 对应文档：02_GORM背景示例_博客系统实战.md

package routes

import (
	"fmt"
	"net/http"

	"blog-system/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(CORSMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(ErrorHandlerMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Blog system is running",
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("/register", handlers.RegisterUser)
			users.POST("/login", handlers.LoginUser)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id/profile", handlers.UpdateUserProfile)
		}

		// 文章相关路由
		posts := api.Group("/posts")
		{
			posts.GET("", handlers.GetPosts)
			posts.GET("/:id", handlers.GetPost)
			posts.GET("/slug/:slug", handlers.GetPostBySlug)
			posts.POST("", handlers.CreatePost)
			posts.PUT("/:id", handlers.UpdatePost)
			posts.DELETE("/:id", handlers.DeletePost)
			posts.POST("/:id/publish", handlers.PublishPost)
			posts.POST("/:id/like", handlers.LikePost)
			posts.DELETE("/:id/like", handlers.UnlikePost)
		}

		// 评论相关路由
		comments := api.Group("/comments")
		{
			comments.GET("/post/:post_id", handlers.GetCommentsByPost)
			comments.POST("", handlers.CreateComment)
			comments.PUT("/:id/approve", handlers.ApproveComment)
			comments.PUT("/:id/reject", handlers.RejectComment)
		}

		// 分类相关路由
		categories := api.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
			categories.GET("/:slug", handlers.GetCategoryBySlug)
			categories.POST("", handlers.CreateCategory)
		}

		// 标签相关路由
		tags := api.Group("/tags")
		{
			tags.GET("", handlers.GetTags)
			tags.GET("/popular", handlers.GetPopularTags)
			tags.POST("", handlers.CreateTag)
		}

		// 统计相关路由
		stats := api.Group("/stats")
		{
			stats.GET("/overview", handlers.GetStatsOverview)
			stats.GET("/posts/popular", handlers.GetPopularPosts)
		}

		// API文档
		api.GET("/docs", handlers.GetAPIDocs)
	}

	return r
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有域名访问
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许携带cookie
		c.Header("Access-Control-Allow-Credentials", "true")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"code":  "INTERNAL_ERROR",
			})
		}
	}
}

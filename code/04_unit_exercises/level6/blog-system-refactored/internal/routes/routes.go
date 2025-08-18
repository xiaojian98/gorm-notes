package routes

import (
	"github.com/gin-gonic/gin"
	"blog-system-refactored/internal/handlers"
	"blog-system-refactored/internal/middleware"
)

// SetupRoutes 设置所有路由
// 参数: r - Gin路由器, userHandler - 用户处理器, postHandler - 文章处理器, commentHandler - 评论处理器, analyticsHandler - 分析处理器
// 返回: 无
func SetupRoutes(
	r *gin.Engine,
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
	analyticsHandler *handlers.AnalyticsHandler,
) {
	// 设置全局中间件
	r.Use(middleware.CORS())           // 跨域中间件
	r.Use(middleware.Logger())         // 日志中间件
	r.Use(middleware.Recovery())       // 恢复中间件
	r.Use(middleware.RateLimit())      // 限流中间件

	// API版本1路由组
	v1 := r.Group("/api/v1")
	{
		// 设置用户相关路由
		setupUserRoutes(v1, userHandler)

		// 设置文章相关路由
		setupPostRoutes(v1, postHandler)

		// 设置评论相关路由
		setupCommentRoutes(v1, commentHandler)

		// 设置分析统计相关路由
		setupAnalyticsRoutes(v1, analyticsHandler)
	}

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Blog system is running",
		})
	})

	// API文档路由
	r.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Documentation",
			"version": "v1",
			"endpoints": map[string]interface{}{
				"users": "/api/v1/users",
				"posts": "/api/v1/posts",
				"comments": "/api/v1/comments",
				"analytics": "/api/v1/analytics",
			},
		})
	})
}

// setupUserRoutes 设置用户相关路由
// 参数: rg - 路由组, handler - 用户处理器
// 返回: 无
func setupUserRoutes(rg *gin.RouterGroup, handler *handlers.UserHandler) {
	users := rg.Group("/users")
	{
		// 公开路由 - 不需要认证
		users.POST("/register", handler.CreateUser)           // 用户注册
		// TODO: 实现用户登录功能
		users.GET("/", handler.ListUsers)                     // 获取用户列表
		users.GET("/:id", handler.GetUser)                    // 获取用户信息
		users.GET("/username/:username", handler.GetUserByUsername) // 根据用户名获取用户
		users.GET("/:id/followers", handler.GetUserFollowers) // 获取用户粉丝列表
		users.GET("/:id/following", handler.GetUserFollowing) // 获取用户关注列表

		// 需要认证的路由
		auth := users.Group("/")
		auth.Use(middleware.AuthRequired()) // 认证中间件
		{
			// 用户资料管理
			auth.PUT("/:id", middleware.OwnershipRequired(), handler.UpdateUser)     // 更新用户信息
			auth.DELETE("/:id", middleware.OwnershipRequired(), handler.DeleteUser)  // 删除用户
			auth.PUT("/:id/password", middleware.OwnershipRequired(), handler.UpdatePassword) // 更新密码

			// 用户关注操作
			auth.POST("/:id/follow", handler.FollowUser)   // 关注用户
			auth.DELETE("/:id/follow", handler.UnfollowUser) // 取消关注

			// 管理员操作
			admin := auth.Group("/")
			admin.Use(middleware.AdminRequired()) // 管理员权限中间件
			{
				admin.PUT("/:id/activate", handler.ActivateUser)   // 激活用户
				admin.PUT("/:id/deactivate", handler.DeactivateUser) // 停用用户
			}
		}
	}
}

// setupPostRoutes 设置文章相关路由
// 参数: rg - 路由组, handler - 文章处理器
// 返回: 无
func setupPostRoutes(rg *gin.RouterGroup, handler *handlers.PostHandler) {
	posts := rg.Group("/posts")
	{
		// 公开路由 - 不需要认证
		posts.GET("/", handler.ListPosts)                   // 获取文章列表
		posts.GET("/:id", handler.GetPost)                  // 获取文章详情
		posts.GET("/slug/:slug", handler.GetPostBySlug)     // 根据slug获取文章
		posts.GET("/popular", handler.GetPopularPosts)      // 获取热门文章
		posts.GET("/recent", handler.GetRecentPosts)        // 获取最新文章
		// TODO: 实现按分类获取文章功能
		// TODO: 实现按标签获取文章功能
		// TODO: 实现按作者获取文章功能
		// TODO: 实现搜索文章功能

		// 分类和标签路由
		posts.GET("/categories", handler.GetCategories)     // 获取分类列表
		posts.GET("/tags", handler.GetTags)                 // 获取标签列表

		// 需要认证的路由
		auth := posts.Group("/")
		auth.Use(middleware.AuthRequired()) // 认证中间件
		{
			// 文章管理
			auth.POST("/", handler.CreatePost)                                    // 创建文章
			auth.PUT("/:id", middleware.OwnershipRequired(), handler.UpdatePost)   // 更新文章
			auth.DELETE("/:id", middleware.OwnershipRequired(), handler.DeletePost) // 删除文章

			// 文章状态管理
			auth.PUT("/:id/publish", middleware.OwnershipRequired(), handler.PublishPost)   // 发布文章
			auth.PUT("/:id/unpublish", middleware.OwnershipRequired(), handler.UnpublishPost) // 取消发布

			// TODO: 实现文章点赞功能
			// auth.POST("/:id/like", handler.LikePost)     // 点赞文章
			// auth.DELETE("/:id/like", handler.UnlikePost) // 取消点赞

			// 管理员操作
			admin := auth.Group("/")
			admin.Use(middleware.AdminRequired()) // 管理员权限中间件
			{
				// TODO: 实现分类管理功能
				// admin.POST("/categories", handler.CreateCategory)        // 创建分类
				// admin.PUT("/categories/:id", handler.UpdateCategory)     // 更新分类
				// admin.DELETE("/categories/:id", handler.DeleteCategory)  // 删除分类

				// TODO: 实现标签管理功能
				// admin.POST("/tags", handler.CreateTag)           // 创建标签
				// admin.PUT("/tags/:id", handler.UpdateTag)        // 更新标签
				// admin.DELETE("/tags/:id", handler.DeleteTag)     // 删除标签
			}
		}
	}
}

// setupCommentRoutes 设置评论相关路由
// 参数: rg - 路由组, handler - 评论处理器
// 返回: 无
func setupCommentRoutes(rg *gin.RouterGroup, handler *handlers.CommentHandler) {
	comments := rg.Group("/comments")
	{
		// 公开路由 - 不需要认证
		comments.GET("/", handler.ListComments)                   // 获取评论列表
		comments.GET("/:id", handler.GetComment)                  // 获取评论详情
		comments.GET("/post/:post_id", handler.GetPostComments)   // 获取文章评论
		comments.GET("/:id/replies", handler.GetCommentReplies)   // 获取评论回复

		// 需要认证的路由
		auth := comments.Group("/")
		auth.Use(middleware.AuthRequired()) // 认证中间件
		{
			// 评论管理
			auth.POST("/", handler.CreateComment)                                    // 创建评论
			auth.PUT("/:id", middleware.OwnershipRequired(), handler.UpdateComment)   // 更新评论
			auth.DELETE("/:id", middleware.OwnershipRequired(), handler.DeleteComment) // 删除评论

			// TODO: 实现评论点赞功能
			// auth.POST("/:id/like", handler.LikeComment)     // 点赞评论
			// auth.DELETE("/:id/like", handler.UnlikeComment) // 取消点赞
			auth.POST("/:id/report", handler.ReportComment) // 举报评论

			// 管理员操作
			admin := auth.Group("/")
			admin.Use(middleware.AdminRequired()) // 管理员权限中间件
			{
				admin.PUT("/:id/approve", handler.ApproveComment)   // 审核通过
				admin.PUT("/:id/reject", handler.RejectComment)     // 审核拒绝
				admin.PUT("/:id/spam", handler.MarkAsSpam)          // 标记为垃圾评论
				// TODO: 实现获取待审核和被举报评论功能
				// admin.GET("/pending", handler.GetPendingComments)   // 获取待审核评论
				// admin.GET("/reported", handler.GetReportedComments) // 获取被举报评论
			}
		}
	}
}

// setupAnalyticsRoutes 设置分析统计相关路由
// 参数: rg - 路由组, handler - 分析处理器
// 返回: 无
func setupAnalyticsRoutes(rg *gin.RouterGroup, handler *handlers.AnalyticsHandler) {
	analytics := rg.Group("/analytics")
	{
		// 需要认证的路由
		auth := analytics.Group("/")
		auth.Use(middleware.AuthRequired()) // 认证中间件
		{
			// 基础统计 - 普通用户可访问
			auth.GET("/dashboard", handler.GetDashboardStats)   // 仪表板统计
			auth.GET("/content", handler.GetContentStats)       // 内容统计
			auth.GET("/popular", handler.GetPopularContent)     // 热门内容
			// TODO: 实现热门标签功能
			// auth.GET("/popular/tags", handler.GetPopularTags)   // 热门标签

			// 管理员统计 - 需要管理员权限
			admin := auth.Group("/")
			admin.Use(middleware.AdminRequired()) // 管理员权限中间件
			{
				// 用户统计
				admin.GET("/users", handler.GetUserStats)           // 用户统计
				// TODO: 实现活跃用户和用户增长功能
				// admin.GET("/users/active", handler.GetActiveUsers)  // 活跃用户
				// admin.GET("/users/growth", handler.GetUserGrowth)   // 用户增长

				// 趋势分析
				admin.GET("/trends/users", handler.GetUserTrend)    // 用户趋势
				admin.GET("/trends/posts", handler.GetPostTrend)    // 文章趋势
				admin.GET("/trends/views", handler.GetViewTrend)    // 浏览量趋势

				// TODO: 实现热门分类功能
				// admin.GET("/popular/categories", handler.GetPopularCategories) // 热门分类

				// 性能统计
				admin.GET("/performance", handler.GetPerformanceStats) // 性能统计

				// 实时统计
				admin.GET("/realtime", handler.GetRealTimeStats)   // 实时统计
			}
		}
	}
}
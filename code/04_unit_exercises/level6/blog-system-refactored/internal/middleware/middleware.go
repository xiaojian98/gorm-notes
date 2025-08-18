package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// CORS 跨域中间件
// 用于处理跨域请求，设置相应的响应头
// 返回: gin.HandlerFunc - Gin中间件函数
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Logger 日志中间件
// 用于记录HTTP请求日志
// 返回: gin.HandlerFunc - Gin中间件函数
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
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

// Recovery 恢复中间件
// 用于捕获panic并返回500错误
// 返回: gin.HandlerFunc - Gin中间件函数
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// RateLimit 限流中间件
// 用于限制请求频率，防止恶意攻击
// 返回: gin.HandlerFunc - Gin中间件函数
func RateLimit() gin.HandlerFunc {
	// 创建限流器：每秒允许100个请求，突发容量200
	limiter := rate.NewLimiter(rate.Limit(100), 200)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
				"message": "请求过于频繁，请稍后再试",
				"code": "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthRequired 认证中间件
// 用于验证用户是否已登录
// 返回: gin.HandlerFunc - Gin中间件函数
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization token
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
				"message": "缺少认证信息",
				"code": "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		tokenString := ""
		if strings.HasPrefix(authorization, "Bearer ") {
			tokenString = authorization[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
				"message": "认证格式错误",
				"code": "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 验证token（这里需要实现JWT验证逻辑）
		userID, role, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
				"message": "认证信息无效",
				"code": "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}

// AdminRequired 管理员权限中间件
// 用于验证用户是否具有管理员权限
// 返回: gin.HandlerFunc - Gin中间件函数
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
				"message": "缺少用户信息",
				"code": "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if userRole != "admin" && userRole != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
				"message": "需要管理员权限",
				"code": "ADMIN_REQUIRED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OwnershipRequired 所有权验证中间件
// 用于验证用户是否拥有资源的所有权
// 返回: gin.HandlerFunc - Gin中间件函数
func OwnershipRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		currentUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
				"message": "缺少用户信息",
				"code": "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// 获取资源ID（从URL参数中）
		resourceIDStr := c.Param("id")
		if resourceIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
				"message": "缺少资源ID",
				"code": "VALIDATION_ERROR",
			})
			c.Abort()
			return
		}

		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
				"message": "资源ID格式错误",
				"code": "VALIDATION_ERROR",
			})
			c.Abort()
			return
		}

		// 检查用户角色，管理员可以访问所有资源
		userRole, _ := c.Get("user_role")
		if userRole == "admin" || userRole == "super_admin" {
			c.Next()
			return
		}

		// 对于普通用户，检查是否为资源所有者
		// 这里需要根据具体的资源类型来验证所有权
		// 例如：检查用户是否拥有该文章、评论等
		if !checkOwnership(c, uint(currentUserID.(uint)), uint(resourceID)) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
				"message": "没有权限访问该资源",
				"code": "OWNERSHIP_REQUIRED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ContentType 内容类型验证中间件
// 用于验证请求的Content-Type
// 参数: contentType - 期望的内容类型
// 返回: gin.HandlerFunc - Gin中间件函数
func ContentType(contentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.GetHeader("Content-Type") != contentType {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Unsupported Media Type",
					"message": fmt.Sprintf("期望的Content-Type: %s", contentType),
					"code": "UNSUPPORTED_MEDIA_TYPE",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequestSize 请求大小限制中间件
// 用于限制请求体的大小
// 参数: maxSize - 最大请求体大小（字节）
// 返回: gin.HandlerFunc - Gin中间件函数
func RequestSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request Entity Too Large",
				"message": fmt.Sprintf("请求体大小超过限制: %d bytes", maxSize),
				"code": "REQUEST_TOO_LARGE",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Security 安全头中间件
// 用于设置安全相关的HTTP响应头
// 返回: gin.HandlerFunc - Gin中间件函数
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止XSS攻击
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// 强制HTTPS（在生产环境中启用）
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// 隐藏服务器信息
		c.Header("Server", "")
		
		c.Next()
	}
}

// validateToken 验证JWT token
// 参数: tokenString - JWT token字符串
// 返回: uint, string, error - 用户ID, 用户角色, 错误信息
func validateToken(tokenString string) (uint, string, error) {
	// TODO: 实现JWT token验证逻辑
	// 这里应该:
	// 1. 解析JWT token
	// 2. 验证签名
	// 3. 检查过期时间
	// 4. 提取用户信息
	
	// 临时实现：简单的token验证
	if tokenString == "" {
		return 0, "", fmt.Errorf("empty token")
	}
	
	// 这里应该实现真正的JWT验证
	// 现在返回模拟数据
	return 1, "user", nil
}

// checkOwnership 检查资源所有权
// 参数: c - Gin上下文, userID - 用户ID, resourceID - 资源ID
// 返回: bool - 是否拥有所有权
func checkOwnership(c *gin.Context, userID, resourceID uint) bool {
	// TODO: 实现资源所有权检查逻辑
	// 这里应该:
	// 1. 根据请求路径确定资源类型
	// 2. 查询数据库验证所有权
	// 3. 返回验证结果
	
	// 获取请求路径来确定资源类型
	path := c.Request.URL.Path
	
	if strings.Contains(path, "/users/") {
		// 用户资源：用户只能操作自己的资源
		return userID == resourceID
	}
	
	if strings.Contains(path, "/posts/") {
		// 文章资源：需要查询数据库验证作者
		// TODO: 查询数据库检查文章作者
		return true // 临时返回true
	}
	
	if strings.Contains(path, "/comments/") {
		// 评论资源：需要查询数据库验证评论者
		// TODO: 查询数据库检查评论作者
		return true // 临时返回true
	}
	
	// 默认拒绝访问
	return false
}

// GetCurrentUserID 获取当前用户ID
// 参数: c - Gin上下文
// 返回: uint, bool - 用户ID, 是否存在
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUserRole 获取当前用户角色
// 参数: c - Gin上下文
// 返回: string, bool - 用户角色, 是否存在
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return userRole.(string), true
}

// IsAdmin 检查当前用户是否为管理员
// 参数: c - Gin上下文
// 返回: bool - 是否为管理员
func IsAdmin(c *gin.Context) bool {
	userRole, exists := GetCurrentUserRole(c)
	if !exists {
		return false
	}
	return userRole == "admin" || userRole == "super_admin"
}
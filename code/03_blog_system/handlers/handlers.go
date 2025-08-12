// 03_blog_system/handlers/handlers.go - 处理器实现
// 对应文档：02_GORM背景示例_博客系统实战.md

package handlers

import (
	"net/http"
	"strconv"
	"time"

	"blog-system/models"
	"blog-system/services"

	"github.com/gin-gonic/gin"
)

// 用户相关处理器

// RegisterUser 用户注册
func RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UserService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
		},
	})
}

// LoginUser 用户登录
func LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UserService.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
		},
		"token": "jwt_token_here", // 实际项目中应该生成JWT token
	})
}

// GetUser 获取用户信息
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := services.UserService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"nickname":   user.Nickname,
			"created_at": user.CreatedAt,
			"profile":    user.Profile,
		},
	})
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
		Website  string `json:"website"`
		Location string `json:"location"`
		Gender   string `json:"gender"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"bio":      req.Bio,
		"website":  req.Website,
		"location": req.Location,
		"gender":   req.Gender,
	}
	err = services.UserService.UpdateUserProfile(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户资料更新成功",
	})
}

// 文章相关处理器

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	categorySlug := c.Query("category")
	tag := c.Query("tag")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	filters := map[string]interface{}{
		"status": "published",
	}
	if categorySlug != "" {
		filters["category_slug"] = categorySlug
	}
	if tag != "" {
		filters["tag"] = tag
	}
	if search != "" {
		filters["keyword"] = search
	}
	posts, total, err := services.PostService.GetPosts(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetPost 获取单篇文章
func GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	post, err := services.PostService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// GetPostBySlug 通过slug获取文章
func GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := services.PostService.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required,max=200"`
		Content    string `json:"content" binding:"required"`
		Excerpt    string `json:"excerpt"`
		Slug       string `json:"slug"`
		AuthorID   uint   `json:"author_id" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
		Tags       []uint `json:"tags"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &models.Post{
		Title:      req.Title,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		Slug:       req.Slug,
		UserID:     req.AuthorID,
		CategoryID: &req.CategoryID,
		Status:     req.Status,
	}
	if err := services.PostService.CreatePost(post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post,
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		Excerpt    string `json:"excerpt"`
		Slug       string `json:"slug"`
		CategoryID uint   `json:"category_id"`
		Tags       []uint `json:"tags"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"content":     req.Content,
		"excerpt":     req.Excerpt,
		"slug":        req.Slug,
		"category_id": req.CategoryID,
		"status":      req.Status,
	}
	err = services.PostService.UpdatePost(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	err = services.PostService.DeletePost(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

// PublishPost 发布文章
func PublishPost(c *gin.Context) {
	idStr := c.Param("id")
	// 检查文章是否存在
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	err = services.PostService.PublishPost(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章发布成功",
	})
}

// LikePost 点赞文章
func LikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	_, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现点赞功能
	// err = services.LikeService.LikePost(uint(postID), req.UserID)
	// if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

// UnlikePost 取消点赞
func UnlikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	_, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现取消点赞功能
	// err = services.LikeService.UnlikePost(uint(postID), req.UserID)
	// if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}

// 评论相关处理器

// GetCommentsByPost 获取文章评论
func GetCommentsByPost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	comments, total, err := services.CommentService.GetCommentsByPostID(uint(postID), 1, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
	})
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	var req struct {
		Content  string `json:"content" binding:"required"`
		PostID   uint   `json:"post_id" binding:"required"`
		AuthorID uint   `json:"author_id" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &models.Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  req.AuthorID,
		Status:  "pending",
	}
	if req.ParentID != nil {
		comment.ParentID = req.ParentID
	}
	if err := services.CommentService.CreateComment(comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// ApproveComment 审核通过评论
func ApproveComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	err = services.CommentService.ApproveComment(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论审核通过",
	})
}

// RejectComment 拒绝评论
func RejectComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	err = services.CommentService.RejectComment(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论已拒绝",
	})
}

// 分类相关处理器

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	categories, err := services.CategoryService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetCategoryBySlug 通过slug获取分类
func GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := services.CategoryService.GetCategoryBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,max=100"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}
	if err := services.CategoryService.CreateCategory(category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "分类创建成功",
		"category": category,
	})
}

// 标签相关处理器

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	tags, err := services.TagService.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// GetPopularTags 获取热门标签
func GetPopularTags(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	tags, err := services.TagService.GetPopularTags(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required,max=50"`
		Slug  string `json:"slug"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name:  req.Name,
		Slug:  req.Slug,
		Color: req.Color,
	}
	if err := services.TagService.CreateTag(tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "标签创建成功",
		"tag":     tag,
	})
}

// 统计相关处理器

// GetStatsOverview 获取统计概览
func GetStatsOverview(c *gin.Context) {
	stats := gin.H{
		"total_posts":    0,
		"total_users":    0,
		"total_comments": 0,
		"total_likes":    0,
		"generated_at":   time.Now(),
	}

	// 这里应该调用实际的统计服务
	// stats, err := services.GetStatsOverview()

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetPopularPosts 获取热门文章
func GetPopularPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	// 这里应该调用实际的服务
	// posts, err := services.GetPopularPosts(limit)

	c.JSON(http.StatusOK, gin.H{
		"posts": []gin.H{},
		"limit": limit,
	})
}

// GetAPIDocs 获取API文档
func GetAPIDocs(c *gin.Context) {
	docs := gin.H{
		"title":       "Blog System API",
		"version":     "1.0.0",
		"description": "博客系统API文档",
		"endpoints": gin.H{
			"users": gin.H{
				"POST /api/users/register":   "用户注册",
				"POST /api/users/login":      "用户登录",
				"GET /api/users/:id":         "获取用户信息",
				"PUT /api/users/:id/profile": "更新用户资料",
			},
			"posts": gin.H{
				"GET /api/posts":           "获取文章列表",
				"GET /api/posts/:id":       "获取文章详情",
				"POST /api/posts":          "创建文章",
				"PUT /api/posts/:id":       "更新文章",
				"DELETE /api/posts/:id":    "删除文章",
				"POST /api/posts/:id/like": "点赞文章",
			},
			"comments": gin.H{
				"GET /api/comments/post/:post_id": "获取文章评论",
				"POST /api/comments":              "创建评论",
				"PUT /api/comments/:id/approve":   "审核通过评论",
			},
			"categories": gin.H{
				"GET /api/categories":       "获取分类列表",
				"POST /api/categories":      "创建分类",
				"GET /api/categories/:slug": "获取分类详情",
			},
			"tags": gin.H{
				"GET /api/tags":         "获取标签列表",
				"GET /api/tags/popular": "获取热门标签",
				"POST /api/tags":        "创建标签",
			},
		},
	}

	c.JSON(http.StatusOK, docs)
}

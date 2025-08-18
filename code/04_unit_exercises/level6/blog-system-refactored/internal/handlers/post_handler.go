package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"blog-system-refactored/internal/models"
	"blog-system-refactored/internal/services"
	"github.com/gin-gonic/gin"
)

// stringToPostStatus 将字符串转换为 PostStatus 类型
// 参数: status - 状态字符串
// 返回: models.PostStatus - 文章状态枚举
func stringToPostStatus(status string) models.PostStatus {
	switch strings.ToLower(status) {
	case "published":
		return models.PostStatusPublished
	case "draft":
		return models.PostStatusDraft
	case "private":
		return models.PostStatusPrivate
	case "trash":
		return models.PostStatusTrash
	default:
		return models.PostStatusDraft
	}
}

// PostHandler 文章处理器
// 负责处理文章相关的HTTP请求
type PostHandler struct {
	postService services.PostService
}

// NewPostHandler 创建文章处理器实例
// 参数: postService - 文章服务接口
// 返回: *PostHandler - 文章处理器实例
func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// 请求和响应结构体定义

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=200"`   // 标题
	Content    string   `json:"content" binding:"required,min=10"`        // 内容
	Summary    string   `json:"summary,omitempty"`                       // 摘要
	Slug       string   `json:"slug,omitempty"`                          // URL别名
	CategoryID uint     `json:"category_id,omitempty"`                   // 分类ID
	Tags       []string `json:"tags,omitempty"`                          // 标签列表
	Status     string   `json:"status,omitempty"`                        // 状态
	Featured   bool     `json:"featured,omitempty"`                      // 是否推荐
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string   `json:"title,omitempty"`      // 标题
	Content    string   `json:"content,omitempty"`    // 内容
	Summary    string   `json:"summary,omitempty"`    // 摘要
	Slug       string   `json:"slug,omitempty"`       // URL别名
	CategoryID uint     `json:"category_id,omitempty"` // 分类ID
	Tags       []string `json:"tags,omitempty"`       // 标签列表
	Status     string   `json:"status,omitempty"`     // 状态
	Featured   bool     `json:"featured,omitempty"`   // 是否推荐
}



// PostListResponse 文章列表响应
type PostListResponse struct {
	Posts      []PostResponse `json:"posts"`       // 文章列表
	Total      int64          `json:"total"`       // 总数
	Page       int            `json:"page"`        // 当前页
	PageSize   int            `json:"page_size"`   // 页大小
	TotalPages int            `json:"total_pages"` // 总页数
}



// PostStatsResponse 文章统计响应
type PostStatsResponse struct {
	TotalPosts     int64 `json:"total_posts"`     // 总文章数
	PublishedPosts int64 `json:"published_posts"` // 已发布文章数
	DraftPosts     int64 `json:"draft_posts"`     // 草稿文章数
	TotalViews     int64 `json:"total_views"`     // 总浏览量
	TotalLikes     int64 `json:"total_likes"`     // 总点赞数
}

// 文章基本操作API

// CreatePost 创建文章
// @Summary 创建新文章
// @Description 创建一篇新的文章
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreatePostRequest true "文章信息"
// @Success 201 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// 获取当前用户ID (实际应用中从JWT或session中获取)
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	// 创建文章模型
	var categoryID *uint
	if req.CategoryID != 0 {
		categoryID = &req.CategoryID
	}

	post := &models.Post{
		Title:      req.Title,
		Content:    req.Content,
		Slug:       req.Slug,
		AuthorID:   userID,
		CategoryID: categoryID,
		Status:     stringToPostStatus(req.Status),
	}

	// 设置默认状态
	if req.Status == "" {
		post.Status = stringToPostStatus("draft")
	} else {
		post.Status = stringToPostStatus(req.Status)
	}

	// 调用服务层创建文章
	err := h.postService.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "创建文章失败",
			Message: err.Error(),
		})
		return
	}

	// TODO: 处理标签功能（需要实现标签ID转换）
	if len(req.Tags) > 0 {
		// 暂时跳过标签处理，因为参数类型不匹配
		// err = h.postService.UpdatePostTags(post.ID, req.Tags)
		// 标签功能需要重新设计参数类型
	}

	// 获取完整的文章信息
	fullPost, err := h.postService.GetPostByID(post.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取文章信息失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, h.toPostResponse(fullPost))
}

// GetPost 获取文章详情
// @Summary 获取文章详情
// @Description 根据文章ID获取文章详细信息
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的文章ID",
			Message: "文章ID必须是有效的数字",
		})
		return
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "文章不存在",
			Message: err.Error(),
		})
		return
	}

	// 增加浏览量
	go func() {
		h.postService.IncrementViewCount(uint(id))
	}()

	c.JSON(http.StatusOK, h.toPostResponse(post))
}

// GetPostBySlug 根据slug获取文章
// @Summary 根据slug获取文章
// @Description 根据文章slug获取文章详细信息
// @Tags posts
// @Produce json
// @Param slug path string true "文章slug"
// @Success 200 {object} PostResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/slug/{slug} [get]
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "slug不能为空",
			Message: "请提供有效的文章slug",
		})
		return
	}

	post, err := h.postService.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "文章不存在",
			Message: err.Error(),
		})
		return
	}

	// 增加浏览量
	go func() {
		h.postService.IncrementViewCount(post.ID)
	}()

	c.JSON(http.StatusOK, h.toPostResponse(post))
}

// UpdatePost 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param post body UpdatePostRequest true "更新信息"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的文章ID",
			Message: "文章ID必须是有效的数字",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// 获取现有文章
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "文章不存在",
			Message: err.Error(),
		})
		return
	}

	// 检查权限 (实际应用中需要验证用户是否有权限修改此文章)
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "只能修改自己的文章",
		})
		return
	}

	// 更新文章信息
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	// Summary 字段在模型中不存在，跳过
	// if req.Summary != "" {
	//     post.Summary = req.Summary
	// }
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.CategoryID != 0 {
		categoryID := req.CategoryID
		post.CategoryID = &categoryID
	}
	if req.Status != "" {
		post.Status = stringToPostStatus(req.Status)
	}
	// Featured 字段在模型中不存在，跳过
	// post.Featured = req.Featured

	err = h.postService.UpdatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "更新文章失败",
			Message: err.Error(),
		})
		return
	}

	// TODO: 更新标签功能需要实现字符串到ID的转换
	if req.Tags != nil {
		// 暂时跳过标签更新，因为需要将 []string 转换为 []uint
		// err = h.postService.UpdatePostTags(post.ID, req.Tags)
		// 需要实现标签名称到ID的映射功能
	}

	// 获取更新后的文章信息
	updatedPost, err := h.postService.GetPostByID(post.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取文章信息失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toPostResponse(updatedPost))
}

// DeletePost 删除文章
// @Summary 删除文章
// @Description 软删除指定的文章
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的文章ID",
			Message: "文章ID必须是有效的数字",
		})
		return
	}

	// 检查权限
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	// 获取文章信息以验证权限
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "文章不存在",
			Message: err.Error(),
		})
		return
	}

	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "只能删除自己的文章",
		})
		return
	}

	err = h.postService.DeletePost(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "删除文章失败",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// 文章查询API

// ListPosts 获取文章列表
// @Summary 获取文章列表
// @Description 分页获取文章列表，支持多种筛选条件
// @Tags posts
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "文章状态"
// @Param category_id query int false "分类ID"
// @Param author_id query int false "作者ID"
// @Param tag query string false "标签名称"
// @Param search query string false "搜索关键词"
// @Param featured query bool false "是否推荐"
// @Param sort query string false "排序方式" default("created_at")
// @Param order query string false "排序顺序" default("desc")
// @Success 200 {object} PostListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	authorID, _ := strconv.ParseUint(c.Query("author_id"), 10, 32)
	tag := c.Query("tag")
	search := c.Query("search")
	featured, _ := strconv.ParseBool(c.Query("featured"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建筛选条件
	filters := map[string]interface{}{}
	if status != "" {
		filters["status"] = status
	}
	if categoryID > 0 {
		filters["category_id"] = uint(categoryID)
	}
	if authorID > 0 {
		filters["author_id"] = uint(authorID)
	}
	if tag != "" {
		filters["tag"] = tag
	}
	if search != "" {
		filters["search"] = search
	}
	if c.Query("featured") != "" {
		filters["featured"] = featured
	}

	// 获取文章列表
	// 将 filters map 转换为 PostFilters 结构体
	postFilters := services.PostFilters{}
	
	// 安全的类型转换
	if status, ok := filters["status"].(string); ok {
		postFilters.Status = status
	}
	if categoryID, ok := filters["category_id"].(uint); ok {
		postFilters.CategoryID = categoryID
	}
	if authorID, ok := filters["author_id"].(uint); ok {
		postFilters.AuthorID = authorID
	}
	if search, ok := filters["search"].(string); ok {
		postFilters.Keyword = search
	}
	posts, total, err := h.postService.ListPosts((page-1)*pageSize, pageSize, postFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取文章列表失败",
			Message: err.Error(),
		})
		return
	}

	// 转换为响应格式
	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = h.toPostResponse(&post)
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, PostListResponse{
		Posts:      postResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetPopularPosts 获取热门文章
// @Summary 获取热门文章
// @Description 获取热门文章列表
// @Tags posts
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Param days query int false "时间范围(天)" default(7)
// @Success 200 {object} PostListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/posts/popular [get]
func (h *PostHandler) GetPopularPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	if limit < 1 || limit > 100 {
		limit = 10
	}
	if days < 1 {
		days = 7
	}

	posts, err := h.postService.GetPopularPosts(limit, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取热门文章失败",
			Message: err.Error(),
		})
		return
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = h.toPostResponse(&post)
	}

	c.JSON(http.StatusOK, PostListResponse{
		Posts:      postResponses,
		Total:      int64(len(posts)),
		Page:       1,
		PageSize:   limit,
		TotalPages: 1,
	})
}

// GetRecentPosts 获取最新文章
// @Summary 获取最新文章
// @Description 获取最新发布的文章列表
// @Tags posts
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} PostListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/posts/recent [get]
func (h *PostHandler) GetRecentPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit < 1 || limit > 100 {
		limit = 10
	}

	posts, err := h.postService.GetRecentPosts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取最新文章失败",
			Message: err.Error(),
		})
		return
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = h.toPostResponse(&post)
	}

	c.JSON(http.StatusOK, PostListResponse{
		Posts:      postResponses,
		Total:      int64(len(posts)),
		Page:       1,
		PageSize:   limit,
		TotalPages: 1,
	})
}

// 文章状态管理API

// PublishPost 发布文章
// @Summary 发布文章
// @Description 将草稿文章发布
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id}/publish [post]
func (h *PostHandler) PublishPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的文章ID",
			Message: "文章ID必须是有效的数字",
		})
		return
	}

	// 检查权限
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	err = h.postService.PublishPost(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "发布文章失败",
			Message: err.Error(),
		})
		return
	}

	// 获取发布后的文章信息
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取文章信息失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toPostResponse(post))
}

// UnpublishPost 取消发布文章
// @Summary 取消发布文章
// @Description 将已发布文章转为草稿
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id}/unpublish [post]
func (h *PostHandler) UnpublishPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的文章ID",
			Message: "文章ID必须是有效的数字",
		})
		return
	}

	// 检查权限
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	err = h.postService.UnpublishPost(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "取消发布失败",
			Message: err.Error(),
		})
		return
	}

	// 获取取消发布后的文章信息
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取文章信息失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toPostResponse(post))
}

// 分类和标签管理API

// GetCategories 获取分类列表
// @Summary 获取分类列表
// @Description 获取所有分类列表
// @Tags categories
// @Produce json
// @Success 200 {array} CategoryResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/categories [get]
func (h *PostHandler) GetCategories(c *gin.Context) {
	categories, err := h.postService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取分类列表失败",
			Message: err.Error(),
		})
		return
	}

	categoryResponses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			PostCount:   int64(category.PostCount),
		}
	}

	c.JSON(http.StatusOK, categoryResponses)
}

// GetTags 获取标签列表
// @Summary 获取标签列表
// @Description 获取所有标签列表
// @Tags tags
// @Produce json
// @Param limit query int false "数量限制" default(50)
// @Success 200 {array} TagResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/tags [get]
func (h *PostHandler) GetTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 200 {
		limit = 50
	}

	tags, err := h.postService.GetPopularTags(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取标签列表失败",
			Message: err.Error(),
		})
		return
	}

	tagResponses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			PostCount:   int64(tag.PostCount),
		}
	}

	c.JSON(http.StatusOK, tagResponses)
}

// 文章统计API

// GetPostStats 获取文章统计信息
// @Summary 获取文章统计
// @Description 获取文章的统计信息
// @Tags posts
// @Produce json
// @Param author_id query int false "作者ID"
// @Success 200 {object} PostStatsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/posts/stats [get]
func (h *PostHandler) GetPostStats(c *gin.Context) {
	authorID, _ := strconv.ParseUint(c.Query("author_id"), 10, 32)

	stats, err := h.postService.GetPostStats(uint(authorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取文章统计失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PostStatsResponse{
		TotalPosts:     stats.TotalPosts,
		PublishedPosts: stats.PublishedPosts,
		DraftPosts:     stats.DraftPosts,
		TotalViews:     stats.TotalViews,
		TotalLikes:     stats.TotalLikes,
	})
}

// 辅助方法

// getCurrentUserID 获取当前用户ID
// 参数: c - Gin上下文
// 返回: uint - 用户ID，0表示未登录
func (h *PostHandler) getCurrentUserID(c *gin.Context) uint {
	// 实际应用中应该从JWT token或session中获取用户ID
	// 这里为了演示，从header中获取
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		return 0
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0
	}

	return uint(userID)
}

// toPostResponse 将文章模型转换为响应格式
// 参数: post - 文章模型
// 返回: PostResponse - 文章响应格式
func (h *PostHandler) toPostResponse(post *models.Post) PostResponse {
	response := PostResponse{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		Excerpt:       post.Excerpt,
		Slug:          post.Slug,
		Status:        post.Status.String(),
		ViewCount:     int64(post.ViewCount),
		LikeCount:     int64(post.LikeCount),
		CommentCount:  int64(post.CommentCount),
		FeaturedImage: post.FeaturedImg,
		PublishedAt:   post.PublishedAt,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}

	// 添加作者信息
	if post.Author != nil {
		response.Author = &UserResponse{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Email:    post.Author.Email,
		}
		if post.Author.Profile != nil {
			response.Author.Avatar = post.Author.Profile.Avatar
			response.Author.Bio = post.Author.Profile.Bio
		}
	}

	// 添加分类信息
	if post.Category != nil {
		response.Category = &CategoryResponse{
			ID:          post.Category.ID,
			Name:        post.Category.Name,
			Slug:        post.Category.Slug,
			Description: post.Category.Description,
			PostCount:   int64(post.Category.PostCount),
		}
	}

	// 添加标签信息
	if len(post.Tags) > 0 {
		response.Tags = make([]TagResponse, len(post.Tags))
		for i, tag := range post.Tags {
			response.Tags[i] = TagResponse{
				ID:        tag.ID,
				Name:      tag.Name,
				Slug:      tag.Slug,
				PostCount: int64(tag.PostCount),
			}
		}
	}

	return response
}

// parseTagsFromString 从字符串解析标签列表
// 参数: tagsStr - 标签字符串，用逗号分隔
// 返回: []string - 标签列表
func (h *PostHandler) parseTagsFromString(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}

	tags := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}

	return result
}
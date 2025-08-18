package handlers

import (
	"net/http"
	"strconv"
	"time"

	"blog-system-refactored/internal/models"
	"blog-system-refactored/internal/services"
	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
// 负责处理评论相关的HTTP请求
type CommentHandler struct {
	commentService services.CommentService
}

// NewCommentHandler 创建评论处理器实例
// 参数: commentService - 评论服务接口
// 返回: *CommentHandler - 评论处理器实例
func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// 请求和响应结构体定义

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`   // 文章ID
	ParentID uint   `json:"parent_id,omitempty"`         // 父评论ID
	Content  string `json:"content" binding:"required,min=1,max=1000"` // 评论内容
	Email    string `json:"email,omitempty"`             // 邮箱(游客评论)
	Website  string `json:"website,omitempty"`           // 网站(游客评论)
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"` // 评论内容
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint               `json:"id"`         // 评论ID
	PostID    uint               `json:"post_id"`    // 文章ID
	ParentID  uint               `json:"parent_id"`  // 父评论ID
	Content   string             `json:"content"`    // 评论内容
	Status    string             `json:"status"`     // 状态
	LikeCount int64              `json:"like_count"` // 点赞数
	Author    *UserResponse      `json:"author"`     // 作者信息
	Email     string             `json:"email"`      // 邮箱(游客评论)
	Website   string             `json:"website"`    // 网站(游客评论)
	Replies   []CommentResponse  `json:"replies"`    // 回复列表
	CreatedAt time.Time          `json:"created_at"` // 创建时间
	UpdatedAt time.Time          `json:"updated_at"` // 更新时间
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`    // 评论列表
	Total      int64             `json:"total"`       // 总数
	Page       int               `json:"page"`        // 当前页
	PageSize   int               `json:"page_size"`   // 页大小
	TotalPages int               `json:"total_pages"` // 总页数
}

// CommentStatsResponse 评论统计响应
type CommentStatsResponse struct {
	TotalComments    int64 `json:"total_comments"`    // 总评论数
	ApprovedComments int64 `json:"approved_comments"` // 已审核评论数
	PendingComments  int64 `json:"pending_comments"`  // 待审核评论数
	SpamComments     int64 `json:"spam_comments"`     // 垃圾评论数
	TotalLikes       int64 `json:"total_likes"`       // 总点赞数
}

// LikeRequest 点赞请求
type LikeRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=post comment"` // 目标类型
	TargetID   uint   `json:"target_id" binding:"required"`                     // 目标ID
}

// ReportRequest 举报请求
type ReportRequest struct {
	CommentID uint   `json:"comment_id" binding:"required"` // 评论ID
	Reason    string `json:"reason" binding:"required"`    // 举报原因
	Details   string `json:"details,omitempty"`           // 详细说明
}

// 评论基本操作API

// CreateComment 创建评论
// @Summary 创建新评论
// @Description 为文章创建新评论或回复
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body CreateCommentRequest true "评论信息"
// @Success 201 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// 获取当前用户ID (可能为0，表示游客评论)
	userID := h.getCurrentUserID(c)

	// 创建评论模型
	comment := &models.Comment{
		PostID:  req.PostID,
		Content: req.Content,
		UserID:  userID,
		Status:  models.CommentStatusPending, // 默认待审核
	}

	// 设置父评论ID（如果有）
	if req.ParentID > 0 {
		comment.ParentID = &req.ParentID
	}

	// 如果是注册用户，状态可以直接设为已审核
	if userID > 0 {
		comment.Status = models.CommentStatusApproved
	}

	// 调用服务层创建评论
	err := h.commentService.CreateComment(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "创建评论失败",
			Message: err.Error(),
		})
		return
	}

	// 获取完整的评论信息
	fullComment, err := h.commentService.GetCommentByID(comment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取评论信息失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, h.toCommentResponse(fullComment))
}

// GetComment 获取评论详情
// @Summary 获取评论详情
// @Description 根据评论ID获取评论详细信息
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	comment, err := h.commentService.GetCommentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "评论不存在",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toCommentResponse(comment))
}

// UpdateComment 更新评论
// @Summary 更新评论
// @Description 更新评论内容
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param comment body UpdateCommentRequest true "更新信息"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// 获取现有评论
	comment, err := h.commentService.GetCommentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "评论不存在",
			Message: err.Error(),
		})
		return
	}

	// 检查权限
	userID := h.getCurrentUserID(c)
	if userID == 0 || comment.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "只能修改自己的评论",
		})
		return
	}

	// 更新评论内容
	comment.Content = req.Content

	err = h.commentService.UpdateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "更新评论失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toCommentResponse(comment))
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 软删除指定的评论
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	// 获取评论信息以验证权限
	comment, err := h.commentService.GetCommentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "评论不存在",
			Message: err.Error(),
		})
		return
	}

	// 检查权限
	userID := h.getCurrentUserID(c)
	if userID == 0 || comment.UserID != userID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "只能删除自己的评论",
		})
		return
	}

	err = h.commentService.DeleteComment(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "删除评论失败",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// 评论查询API

// GetPostComments 获取文章评论列表
// @Summary 获取文章评论
// @Description 获取指定文章的评论列表
// @Tags comments
// @Produce json
// @Param post_id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "评论状态"
// @Param sort query string false "排序方式" default("created_at")
// @Param order query string false "排序顺序" default("desc")
// @Success 200 {object} CommentListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/posts/{post_id}/comments [get]
func (h *CommentHandler) GetPostComments(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的文章ID",
			Message: "文章ID必须是有效的数字",
		})
		return
	}

	// TODO: 实现获取文章评论功能
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "功能暂未实现",
		Message: "获取文章评论功能正在开发中",
	})

	// 功能暂未实现，移除未定义的变量引用
	// TODO: 实现评论列表响应格式转换
}

// GetCommentReplies 获取评论回复
// @Summary 获取评论回复
// @Description 获取指定评论的回复列表
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} CommentListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/comments/{id}/replies [get]
func (h *CommentHandler) GetCommentReplies(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	replies, total, err := h.commentService.GetCommentReplies(uint(parentID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取回复列表失败",
			Message: err.Error(),
		})
		return
	}

	replyResponses := make([]CommentResponse, len(replies))
	for i, reply := range replies {
		replyResponses[i] = h.toCommentResponse(&reply)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, CommentListResponse{
		Comments:   replyResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// ListComments 获取评论列表
// @Summary 获取评论列表
// @Description 分页获取评论列表，支持多种筛选条件
// @Tags comments
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "评论状态"
// @Param user_id query int false "用户ID"
// @Param post_id query int false "文章ID"
// @Param search query string false "搜索关键词"
// @Success 200 {object} CommentListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	postID, _ := strconv.ParseUint(c.Query("post_id"), 10, 32)
	search := c.Query("search")

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
	if userID > 0 {
		filters["user_id"] = uint(userID)
	}
	if postID > 0 {
		filters["post_id"] = uint(postID)
	}
	if search != "" {
		filters["search"] = search
	}

	// TODO: 实现评论列表获取功能
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "功能暂未实现",
		Message: "评论列表管理功能正在开发中",
	})
}

// 评论状态管理API

// ApproveComment 审核通过评论
// @Summary 审核通过评论
// @Description 将待审核评论设为已审核
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/comments/{id}/approve [post]
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	// 检查管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "需要管理员权限",
		})
		return
	}

	err = h.commentService.ApproveComment(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "审核评论失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论审核通过",
	})
}

// RejectComment 拒绝评论
// @Summary 拒绝评论
// @Description 拒绝待审核评论
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/comments/{id}/reject [post]
func (h *CommentHandler) RejectComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	// 检查管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "需要管理员权限",
		})
		return
	}

	err = h.commentService.RejectComment(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "拒绝评论失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论已拒绝",
	})
}

// MarkAsSpam 标记为垃圾评论
// @Summary 标记垃圾评论
// @Description 将评论标记为垃圾评论
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} CommentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/comments/{id}/spam [post]
func (h *CommentHandler) MarkAsSpam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的评论ID",
			Message: "评论ID必须是有效的数字",
		})
		return
	}

	// 检查管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "权限不足",
			Message: "需要管理员权限",
		})
		return
	}

	err = h.commentService.MarkAsSpam(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "标记垃圾评论失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论已标记为垃圾评论",
	})
}

// 点赞和举报API

// LikeTarget 点赞
// @Summary 点赞文章或评论
// @Description 为文章或评论点赞
// @Tags likes
// @Accept json
// @Produce json
// @Param like body LikeRequest true "点赞信息"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/likes [post]
func (h *CommentHandler) LikeTarget(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	// TODO: 实现点赞功能
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "功能暂未实现",
		Message: "点赞功能正在开发中",
	})
}

// UnlikeTarget 取消点赞
// @Summary 取消点赞
// @Description 取消对文章或评论的点赞
// @Tags likes
// @Accept json
// @Produce json
// @Param like body LikeRequest true "取消点赞信息"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/likes [delete]
func (h *CommentHandler) UnlikeTarget(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	// TODO: 实现取消点赞功能
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "功能暂未实现",
		Message: "取消点赞功能正在开发中",
	})
}

// ReportComment 举报评论
// @Summary 举报评论
// @Description 举报不当评论
// @Tags comments
// @Accept json
// @Produce json
// @Param report body ReportRequest true "举报信息"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/comments/report [post]
func (h *CommentHandler) ReportComment(c *gin.Context) {
	var req ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	userID := h.getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "未授权",
			Message: "请先登录",
		})
		return
	}

	// TODO: 实现举报评论功能
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "功能暂未实现",
		Message: "举报功能正在开发中",
	})
}

// 评论统计API

// GetCommentStats 获取评论统计信息
// @Summary 获取评论统计
// @Description 获取评论的统计信息
// @Tags comments
// @Produce json
// @Param post_id query int false "文章ID"
// @Param user_id query int false "用户ID"
// @Success 200 {object} CommentStatsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/comments/stats [get]
func (h *CommentHandler) GetCommentStats(c *gin.Context) {
	// 解析查询参数（暂时不使用）
	_, _ = strconv.ParseUint(c.Query("post_id"), 10, 32)
	_, _ = strconv.ParseUint(c.Query("user_id"), 10, 32)

	// TODO: 实现评论统计功能
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error:   "功能暂未实现",
		Message: "评论统计功能正在开发中",
	})
}

// 辅助方法

// getCurrentUserID 获取当前用户ID
// 参数: c - Gin上下文
// 返回: uint - 用户ID，0表示未登录
func (h *CommentHandler) getCurrentUserID(c *gin.Context) uint {
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

// isAdmin 检查是否为管理员
// 参数: c - Gin上下文
// 返回: bool - 是否为管理员
func (h *CommentHandler) isAdmin(c *gin.Context) bool {
	// 实际应用中应该从JWT token或session中获取用户角色
	// 这里为了演示，从header中获取
	role := c.GetHeader("X-User-Role")
	return role == "admin"
}

// toCommentResponse 将评论模型转换为响应格式
// 参数: comment - 评论模型
// 返回: CommentResponse - 评论响应格式
func (h *CommentHandler) toCommentResponse(comment *models.Comment) CommentResponse {
	var parentID uint
	if comment.ParentID != nil {
		parentID = *comment.ParentID
	}

	response := CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  parentID,
		Content:   comment.Content,
		Status:    string(comment.Status),
		LikeCount: int64(comment.LikeCount),
		Email:     "", // 字段不存在，设为空
		Website:   "", // 字段不存在，设为空
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	// 添加作者信息
	if comment.User != nil {
		response.Author = &UserResponse{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Nickname: comment.User.Username, // 使用 Username 代替不存在的 Nickname
			Email:    comment.User.Email,
		}
		if comment.User.Profile != nil {
			response.Author.Avatar = comment.User.Profile.Avatar
			response.Author.Bio = comment.User.Profile.Bio
		}
	}

	// TODO: 添加回复信息（Replies字段暂未在模型中定义）
	// 回复功能需要在模型中添加相应字段后实现

	return response
}
package handlers

import (
	"net/http"
	"strconv"

	"blog-system-refactored/internal/models"
	"blog-system-refactored/internal/services"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
// 负责处理用户相关的HTTP请求
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler 创建用户处理器实例
// 参数: userService - 用户服务接口
// 返回: *UserHandler - 用户处理器实例
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 请求和响应结构体定义

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"required,email"`          // 邮箱
	Password string `json:"password" binding:"required,min=6"`       // 密码
	Nickname string `json:"nickname,omitempty"`                     // 昵称
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname,omitempty"` // 昵称
	Email    string `json:"email,omitempty"`    // 邮箱
	Bio      string `json:"bio,omitempty"`      // 个人简介
	Avatar   string `json:"avatar,omitempty"`   // 头像URL
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6"` // 新密码
}



// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []UserResponse `json:"users"`       // 用户列表
	Total      int64          `json:"total"`       // 总数
	Page       int            `json:"page"`        // 当前页
	PageSize   int            `json:"page_size"`   // 页大小
	TotalPages int            `json:"total_pages"` // 总页数
}



// 用户基本操作API

// CreateUser 创建用户
// @Summary 创建新用户
// @Description 创建一个新的用户账户
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "用户信息"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// 创建用户模型
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // 注意：实际应用中需要加密
		Status:       models.StatusActive,
	}

	// 调用服务层创建用户
	err := h.userService.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "创建用户失败",
			Message: err.Error(),
		})
		return
	}

	// 返回创建的用户信息
	c.JSON(http.StatusCreated, h.toUserResponse(user))
}

// GetUser 获取用户信息
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "用户不存在",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的基本信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body UpdateUserRequest true "更新信息"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// 获取现有用户
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "用户不存在",
			Message: err.Error(),
		})
		return
	}

	// 更新用户信息
	if req.Email != "" {
		user.Email = req.Email
	}

	// 更新用户资料
	if user.Profile == nil {
		user.Profile = &models.UserProfile{UserID: user.ID}
	}
	if req.Bio != "" {
		user.Profile.Bio = req.Bio
	}
	if req.Avatar != "" {
		user.Profile.Avatar = req.Avatar
	}

	err = h.userService.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "更新用户失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 软删除指定的用户
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "删除用户失败",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// 用户查询API

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags users
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "用户状态"
// @Param search query string false "搜索关键词"
// @Success 200 {object} UserListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 计算偏移量和限制
	offset := (page - 1) * pageSize
	limit := pageSize

	// 获取用户列表
	users, total, err := h.userService.ListUsers(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取用户列表失败",
			Message: err.Error(),
		})
		return
	}

	// 转换为响应格式
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.toUserResponse(&user)
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetUserByUsername 根据用户名获取用户
// @Summary 根据用户名获取用户
// @Description 根据用户名获取用户详细信息
// @Tags users
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} UserResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/username/{username} [get]
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "用户名不能为空",
			Message: "请提供有效的用户名",
		})
		return
	}

	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "用户不存在",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, h.toUserResponse(user))
}

// 用户状态管理API

// ActivateUser 激活用户
// @Summary 激活用户
// @Description 激活指定的用户账户
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	err = h.userService.ActivateUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "激活用户失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "用户激活成功",
	})
}

// DeactivateUser 停用用户
// @Summary 停用用户
// @Description 停用指定的用户账户
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	err = h.userService.DeactivateUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "停用用户失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "用户停用成功",
	})
}

// UpdatePassword 更新密码
// @Summary 更新用户密码
// @Description 更新用户的登录密码
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param password body UpdatePasswordRequest true "密码信息"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "参数验证失败",
			Message: err.Error(),
		})
		return
	}

	// TODO: 实现密码更新逻辑
	// 需要先获取用户，验证旧密码，然后更新新密码
	c.JSON(http.StatusNotImplemented, ErrorResponse{
		Error: "功能暂未实现",
	})
}

// 用户关注API

// FollowUser 关注用户
// @Summary 关注用户
// @Description 关注指定的用户
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Param target_id path int true "目标用户ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/follow/{target_id} [post]
func (h *UserHandler) FollowUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的目标用户ID",
			Message: "目标用户ID必须是有效的数字",
		})
		return
	}

	err = h.userService.FollowUser(uint(userID), uint(targetID))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "关注用户失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "关注成功",
	})
}

// UnfollowUser 取消关注用户
// @Summary 取消关注用户
// @Description 取消关注指定的用户
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Param target_id path int true "目标用户ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/unfollow/{target_id} [delete]
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的目标用户ID",
			Message: "目标用户ID必须是有效的数字",
		})
		return
	}

	err = h.userService.UnfollowUser(uint(userID), uint(targetID))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "取消关注失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "取消关注成功",
	})
}

// GetUserFollowers 获取用户粉丝列表
// @Summary 获取用户粉丝
// @Description 获取指定用户的粉丝列表
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} UserListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/followers [get]
func (h *UserHandler) GetUserFollowers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
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

	followers, total, err := h.userService.GetFollowers(uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取粉丝列表失败",
			Message: err.Error(),
		})
		return
	}

	userResponses := make([]UserResponse, len(followers))
	for i, user := range followers {
		userResponses[i] = h.toUserResponse(&user)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetUserFollowing 获取用户关注列表
// @Summary 获取用户关注
// @Description 获取指定用户的关注列表
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} UserListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/following [get]
func (h *UserHandler) GetUserFollowing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
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

	following, total, err := h.userService.GetFollowing(uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "获取关注列表失败",
			Message: err.Error(),
		})
		return
	}

	userResponses := make([]UserResponse, len(following))
	for i, user := range following {
		userResponses[i] = h.toUserResponse(&user)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// 用户统计API

// GetUserStats 获取用户统计信息
// @Summary 获取用户统计
// @Description 获取指定用户的统计信息
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserStatsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "无效的用户ID",
			Message: "用户ID必须是有效的数字",
		})
		return
	}

	stats, err := h.userService.GetUserStats(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "获取用户统计失败",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserStats{
		TotalPosts:     stats.PostsCount,
		TotalComments:  stats.CommentsCount,
		TotalLikes:     stats.LikesReceived,
		TotalViews:     stats.TotalViews,
		FollowerCount:  stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
		JoinDays:       stats.JoinedDays,
	})
}

// 辅助方法

// toUserResponse 将用户模型转换为响应格式
// 参数: user - 用户模型
// 返回: UserResponse - 用户响应格式
func (h *UserHandler) toUserResponse(user *models.User) UserResponse {
	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 添加用户资料信息
	if user.Profile != nil {
		response.Avatar = user.Profile.Avatar
		response.Bio = user.Profile.Bio
	}

	return response
}
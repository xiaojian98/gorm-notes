package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"../models"
	"../services"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 实际项目中需要加密
		Nickname: req.Nickname,
		Phone:    req.Phone,
		RoleID:   2, // 默认学生角色
		Status:   1, // 正常状态
	}

	if err := ctrl.userService.CreateUser(user); err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 返回用户信息（不包含密码）
	user.Password = ""
	Success(c, user)
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 验证用户
	user, err := ctrl.userService.GetUserByEmail(req.Email)
	if err != nil {
		Error(c, 401, "邮箱或密码错误")
		return
	}

	// 实际项目中需要验证加密后的密码
	if user.Password != req.Password {
		Error(c, 401, "邮箱或密码错误")
		return
	}

	if user.Status != 1 {
		Error(c, 401, "账户已被禁用")
		return
	}

	// 更新最后登录时间
	clientIP := c.ClientIP()
	ctrl.userService.UpdateLastLogin(user.ID, clientIP)

	// 生成JWT Token（这里简化处理）
	token := "jwt_token_" + strconv.Itoa(int(user.ID))

	user.Password = ""
	Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

// GetProfile 获取用户资料
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id") // 从中间件获取

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		Error(c, 404, err.Error())
		return
	}

	user.Password = ""
	Success(c, user)
}

// UpdateProfile 更新用户资料
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Gender   int8   `json:"gender"`
		Birthday string `json:"birthday"`
		Bio      string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Gender > 0 {
		updates["gender"] = req.Gender
	}
	if req.Birthday != "" {
		if birthday, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			updates["birthday"] = &birthday
		}
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}

	if err := ctrl.userService.UpdateUser(userID, updates); err != nil {
		Error(c, 500, "更新失败")
		return
	}

	Success(c, nil)
}

// GetUsers 获取用户列表（管理员）
func (ctrl *UserController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			filters["status"] = s
		}
	}
	if roleID := c.Query("role_id"); roleID != "" {
		if r, err := strconv.Atoi(roleID); err == nil {
			filters["role_id"] = r
		}
	}
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}

	users, total, err := ctrl.userService.GetUsers(page, pageSize, filters)
	if err != nil {
		Error(c, 500, "查询失败")
		return
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	Success(c, PageResponse{
		List:     users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// CourseController 课程控制器
type CourseController struct {
	courseService *services.CourseService
}

// NewCourseController 创建课程控制器
func NewCourseController(courseService *services.CourseService) *CourseController {
	return &CourseController{courseService: courseService}
}

// GetCourses 获取课程列表
func (ctrl *CourseController) GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	
	// 状态过滤
	if status := c.Query("status"); status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			filters["status"] = s
		}
	} else {
		// 默认只显示已发布的课程
		filters["status"] = 2
	}

	// 分类过滤
	if categoryID := c.Query("category_id"); categoryID != "" {
		if cid, err := strconv.Atoi(categoryID); err == nil {
			filters["category_id"] = cid
		}
	}

	// 讲师过滤
	if instructorID := c.Query("instructor_id"); instructorID != "" {
		if iid, err := strconv.Atoi(instructorID); err == nil {
			filters["instructor_id"] = iid
		}
	}

	// 难度过滤
	if level := c.Query("level"); level != "" {
		if l, err := strconv.Atoi(level); err == nil {
			filters["level"] = l
		}
	}

	// 免费课程过滤
	if isFree := c.Query("is_free"); isFree != "" {
		filters["is_free"] = isFree == "true" || isFree == "1"
	}

	// 推荐课程过滤
	if isRecommend := c.Query("is_recommend"); isRecommend != "" {
		filters["is_recommend"] = isRecommend == "true" || isRecommend == "1"
	}

	// 关键词搜索
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}

	// 价格范围
	if priceMin := c.Query("price_min"); priceMin != "" {
		if pm, err := strconv.ParseInt(priceMin, 10, 64); err == nil {
			filters["price_min"] = pm * 100 // 转换为分
		}
	}
	if priceMax := c.Query("price_max"); priceMax != "" {
		if pm, err := strconv.ParseInt(priceMax, 10, 64); err == nil {
			filters["price_max"] = pm * 100 // 转换为分
		}
	}

	// 排序
	if sort := c.Query("sort"); sort != "" {
		filters["sort"] = sort
	}

	courses, total, err := ctrl.courseService.GetCourses(page, pageSize, filters)
	if err != nil {
		Error(c, 500, "查询失败")
		return
	}

	Success(c, PageResponse{
		List:     courses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetCourse 获取课程详情
func (ctrl *CourseController) GetCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "参数错误")
		return
	}

	course, err := ctrl.courseService.GetCourseByID(uint(id))
	if err != nil {
		Error(c, 404, err.Error())
		return
	}

	Success(c, course)
}

// CreateCourse 创建课程（讲师/管理员）
func (ctrl *CourseController) CreateCourse(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Title          string `json:"title" binding:"required"`
		Subtitle       string `json:"subtitle"`
		Slug           string `json:"slug" binding:"required"`
		Description    string `json:"description"`
		Cover          string `json:"cover"`
		CategoryID     uint   `json:"category_id" binding:"required"`
		Level          int8   `json:"level" binding:"required,min=1,max=4"`
		Price          int64  `json:"price"`
		OriginalPrice  int64  `json:"original_price"`
		IsFree         bool   `json:"is_free"`
		IsRecommend    bool   `json:"is_recommend"`
		Tags           string `json:"tags"`
		Requirements   string `json:"requirements"`
		LearningGoals  string `json:"learning_goals"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	course := &models.Course{
		Title:         req.Title,
		Subtitle:      req.Subtitle,
		Slug:          req.Slug,
		Description:   req.Description,
		Cover:         req.Cover,
		CategoryID:    req.CategoryID,
		InstructorID:  userID,
		Level:         req.Level,
		Price:         req.Price * 100, // 转换为分
		OriginalPrice: req.OriginalPrice * 100,
		IsFree:        req.IsFree,
		IsRecommend:   req.IsRecommend,
		Tags:          req.Tags,
		Requirements:  req.Requirements,
		LearningGoals: req.LearningGoals,
		Status:        1, // 草稿状态
	}

	if err := ctrl.courseService.CreateCourse(course); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, course)
}

// UpdateCourse 更新课程
func (ctrl *CourseController) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "参数错误")
		return
	}

	var req struct {
		Title          string `json:"title"`
		Subtitle       string `json:"subtitle"`
		Description    string `json:"description"`
		Cover          string `json:"cover"`
		CategoryID     uint   `json:"category_id"`
		Level          int8   `json:"level"`
		Price          int64  `json:"price"`
		OriginalPrice  int64  `json:"original_price"`
		IsFree         *bool  `json:"is_free"`
		IsRecommend    *bool  `json:"is_recommend"`
		Tags           string `json:"tags"`
		Requirements   string `json:"requirements"`
		LearningGoals  string `json:"learning_goals"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Subtitle != "" {
		updates["subtitle"] = req.Subtitle
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.Price >= 0 {
		updates["price"] = req.Price * 100
	}
	if req.OriginalPrice >= 0 {
		updates["original_price"] = req.OriginalPrice * 100
	}
	if req.IsFree != nil {
		updates["is_free"] = *req.IsFree
	}
	if req.IsRecommend != nil {
		updates["is_recommend"] = *req.IsRecommend
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}
	if req.LearningGoals != "" {
		updates["learning_goals"] = req.LearningGoals
	}

	if err := ctrl.courseService.UpdateCourse(uint(id), updates); err != nil {
		Error(c, 500, "更新失败")
		return
	}

	Success(c, nil)
}

// PublishCourse 发布课程
func (ctrl *CourseController) PublishCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.courseService.PublishCourse(uint(id)); err != nil {
		Error(c, 500, "发布失败")
		return
	}

	Success(c, nil)
}

// OrderController 订单控制器
type OrderController struct {
	orderService    *services.OrderService
	learningService *services.LearningService
}

// NewOrderController 创建订单控制器
func NewOrderController(orderService *services.OrderService, learningService *services.LearningService) *OrderController {
	return &OrderController{
		orderService:    orderService,
		learningService: learningService,
	}
}

// CreateOrder 创建订单
func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CourseIDs   []uint `json:"course_ids" binding:"required,min=1"`
		CouponCode  string `json:"coupon_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	order, err := ctrl.orderService.CreateOrder(userID, req.CourseIDs, req.CouponCode)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, order)
}

// PayOrder 支付订单
func (ctrl *OrderController) PayOrder(c *gin.Context) {
	orderNo := c.Param("order_no")

	var req struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
		PaymentNo     string `json:"payment_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.PayOrder(orderNo, req.PaymentMethod, req.PaymentNo); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

// GetOrders 获取订单列表
func (ctrl *OrderController) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var status *int8
	if s := c.Query("status"); s != "" {
		if st, err := strconv.Atoi(s); err == nil {
			statusVal := int8(st)
			status = &statusVal
		}
	}

	orders, total, err := ctrl.orderService.GetOrdersByUserID(userID, page, pageSize, status)
	if err != nil {
		Error(c, 500, "查询失败")
		return
	}

	Success(c, PageResponse{
		List:     orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// CancelOrder 取消订单
func (ctrl *OrderController) CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	if err := ctrl.orderService.CancelOrder(orderNo, userID); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

// GetLearningCourses 获取学习的课程
func (ctrl *OrderController) GetLearningCourses(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	courses, total, err := ctrl.learningService.GetUserLearningCourses(userID, page, pageSize)
	if err != nil {
		Error(c, 500, "查询失败")
		return
	}

	Success(c, PageResponse{
		List:     courses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// UpdateProgress 更新学习进度
func (ctrl *OrderController) UpdateProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		CourseID  uint `json:"course_id" binding:"required"`
		LessonID  uint `json:"lesson_id" binding:"required"`
		Progress  int  `json:"progress" binding:"required,min=0,max=100"`
		WatchTime int  `json:"watch_time" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.learningService.UpdateProgress(userID, req.CourseID, req.LessonID, req.Progress, req.WatchTime); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

// GetCourseProgress 获取课程学习进度
func (ctrl *OrderController) GetCourseProgress(c *gin.Context) {
	userID := c.GetUint("user_id")
	courseID, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
	if err != nil {
		Error(c, 400, "参数错误")
		return
	}

	progress, err := ctrl.learningService.GetUserCourseProgress(userID, uint(courseID))
	if err != nil {
		Error(c, 500, "查询失败")
		return
	}

	Success(c, progress)
}

// AuthMiddleware JWT认证中间件（简化版）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			Error(c, 401, "未登录")
			c.Abort()
			return
		}

		// 简化的token验证，实际项目中需要验证JWT
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		if strings.HasPrefix(token, "jwt_token_") {
			userIDStr := strings.TrimPrefix(token, "jwt_token_")
			if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
				c.Set("user_id", uint(userID))
				c.Next()
				return
			}
		}

		Error(c, 401, "token无效")
		c.Abort()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里应该检查用户角色，简化处理
		c.Next()
	}
}
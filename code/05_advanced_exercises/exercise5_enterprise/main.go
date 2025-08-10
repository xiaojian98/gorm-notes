package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ========== 数据模型定义 ==========

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string       `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string       `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone       string       `gorm:"uniqueIndex;size:20" json:"phone"`
	Password    string       `gorm:"size:255;not null" json:"-"`
	Nickname    string       `gorm:"size:50" json:"nickname"`
	Avatar      string       `gorm:"size:255" json:"avatar"`
	Status      int8         `gorm:"default:1;comment:1-正常,2-禁用" json:"status"`
	RoleID      uint         `gorm:"index;not null" json:"role_id"`
	LastLoginAt *time.Time   `json:"last_login_at"`
	
	// 关联
	Role            Role             `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Profile         UserProfile      `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Orders          []Order          `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	LearningProgress []LearningProgress `gorm:"foreignKey:UserID" json:"learning_progress,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// UserProfile 用户资料模型
type UserProfile struct {
	BaseModel
	UserID   uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	RealName string `gorm:"size:50" json:"real_name"`
	Gender   int8   `gorm:"default:0;comment:0-未知,1-男,2-女" json:"gender"`
	Birthday *time.Time `json:"birthday"`
	Bio      string `gorm:"type:text" json:"bio"`
	Location string `gorm:"size:100" json:"location"`
	Website  string `gorm:"size:255" json:"website"`
	
	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// Category 课程分类模型
type Category struct {
	BaseModel
	Name        string `gorm:"size:50;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `gorm:"size:255" json:"icon"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联
	Parent   *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Courses  []Course   `gorm:"foreignKey:CategoryID" json:"courses,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Course 课程模型
type Course struct {
	BaseModel
	Title       string `gorm:"size:255;not null" json:"title"`
	Slug        string `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Cover       string `gorm:"size:255" json:"cover"`
	CategoryID  uint   `gorm:"index;not null" json:"category_id"`
	InstructorID uint  `gorm:"index;not null" json:"instructor_id"`
	Price       int64  `gorm:"not null;comment:价格(分)" json:"price"`
	OriginalPrice int64 `gorm:"default:0;comment:原价(分)" json:"original_price"`
	Level       int8   `gorm:"default:1;comment:1-初级,2-中级,3-高级" json:"level"`
	Duration    int    `gorm:"default:0;comment:课程时长(分钟)" json:"duration"`
	StudentCount int   `gorm:"default:0;comment:学生数量" json:"student_count"`
	LessonCount  int   `gorm:"default:0;comment:课时数量" json:"lesson_count"`
	Rating      float32 `gorm:"default:0;comment:评分" json:"rating"`
	Status      int8   `gorm:"default:1;comment:1-草稿,2-发布,3-下架" json:"status"`
	PublishedAt *time.Time `json:"published_at"`
	
	// 关联
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Instructor  User      `gorm:"foreignKey:InstructorID" json:"instructor,omitempty"`
	Chapters    []Chapter `gorm:"foreignKey:CourseID" json:"chapters,omitempty"`
	Orders      []Order   `gorm:"many2many:order_items;" json:"orders,omitempty"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// Chapter 章节模型
type Chapter struct {
	BaseModel
	CourseID    uint   `gorm:"index;not null" json:"course_id"`
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联
	Course  Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:ChapterID" json:"lessons,omitempty"`
}

// TableName 指定表名
func (Chapter) TableName() string {
	return "chapters"
}

// Lesson 课时模型
type Lesson struct {
	BaseModel
	ChapterID   uint   `gorm:"index;not null" json:"chapter_id"`
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	VideoURL    string `gorm:"size:500" json:"video_url"`
	Duration    int    `gorm:"default:0;comment:时长(秒)" json:"duration"`
	Sort        int    `gorm:"default:0" json:"sort"`
	IsFree      bool   `gorm:"default:false;comment:是否免费" json:"is_free"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联
	Chapter          Chapter            `gorm:"foreignKey:ChapterID" json:"chapter,omitempty"`
	LearningProgress []LearningProgress `gorm:"foreignKey:LessonID" json:"learning_progress,omitempty"`
}

// TableName 指定表名
func (Lesson) TableName() string {
	return "lessons"
}

// Order 订单模型
type Order struct {
	BaseModel
	OrderNo        string     `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	TotalAmount    int64      `gorm:"not null;comment:总金额(分)" json:"total_amount"`
	PayAmount      int64      `gorm:"not null;comment:实付金额(分)" json:"pay_amount"`
	DiscountAmount int64      `gorm:"default:0;comment:优惠金额(分)" json:"discount_amount"`
	Status         int8       `gorm:"index;default:1;comment:1-待付款,2-已付款,3-已完成,4-已取消" json:"status"`
	PaymentMethod  string     `gorm:"size:50" json:"payment_method"`
	PaymentNo      string     `gorm:"size:100" json:"payment_no"`
	PaidAt         *time.Time `json:"paid_at"`
	ExpiredAt      *time.Time `json:"expired_at"`
	Remark         string     `gorm:"type:text" json:"remark"`
	
	// 关联
	User    User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items   []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Courses []Course  `gorm:"many2many:order_items;" json:"courses,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	BaseModel
	OrderID     uint   `gorm:"index;not null" json:"order_id"`
	CourseID    uint   `gorm:"index;not null" json:"course_id"`
	CourseName  string `gorm:"size:255;not null" json:"course_name"`
	Price       int64  `gorm:"not null;comment:价格(分)" json:"price"`
	OriginalPrice int64 `gorm:"default:0;comment:原价(分)" json:"original_price"`
	
	// 关联
	Order  Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// LearningProgress 学习进度模型
type LearningProgress struct {
	BaseModel
	UserID     uint `gorm:"index;not null" json:"user_id"`
	CourseID   uint `gorm:"index;not null" json:"course_id"`
	LessonID   uint `gorm:"index;not null" json:"lesson_id"`
	Progress   int  `gorm:"default:0;comment:进度百分比" json:"progress"`
	WatchTime  int  `gorm:"default:0;comment:观看时长(秒)" json:"watch_time"`
	IsCompleted bool `gorm:"default:false;comment:是否完成" json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
	
	// 关联
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Lesson Lesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
}

// TableName 指定表名
func (LearningProgress) TableName() string {
	return "learning_progress"
}

// ========== 数据库配置 ==========

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Charset  string
}

// ConnectDatabase 连接数据库
func ConnectDatabase(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return db, nil
}

// ========== 服务层 ==========

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := s.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := s.db.Preload("Role").Preload("Profile").
		Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&users).Error

	return users, total, err
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*User, error) {
	var user User
	err := s.db.Preload("Role").Preload("Profile").First(&user, id).Error
	return &user, err
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *User) error {
	return s.db.Create(user).Error
}

// CourseService 课程服务
type CourseService struct {
	db *gorm.DB
}

// NewCourseService 创建课程服务
func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{db: db}
}

// GetCourses 获取课程列表
func (s *CourseService) GetCourses(page, pageSize int, categoryID *uint) ([]Course, int64, error) {
	var courses []Course
	var total int64

	query := s.db.Model(&Course{}).Where("status = ?", 2) // 只查询已发布的课程
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	offset := (page - 1) * pageSize

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("Category").Preload("Instructor").
		Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&courses).Error

	return courses, total, err
}

// GetCourseByID 根据ID获取课程详情
func (s *CourseService) GetCourseByID(id uint) (*Course, error) {
	var course Course
	err := s.db.Preload("Category").Preload("Instructor").
		Preload("Chapters.Lessons").First(&course, id).Error
	return &course, err
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(course *Course) error {
	return s.db.Create(course).Error
}

// OrderService 订单服务
type OrderService struct {
	db *gorm.DB
}

// NewOrderService 创建订单服务
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID uint, courseIDs []uint) (*Order, error) {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询课程信息
	var courses []Course
	if err := tx.Where("id IN ? AND status = ?", courseIDs, 2).Find(&courses).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(courses) != len(courseIDs) {
		tx.Rollback()
		return nil, fmt.Errorf("部分课程不存在或已下架")
	}

	// 计算总金额
	var totalAmount int64
	for _, course := range courses {
		totalAmount += course.Price
	}

	// 创建订单
	order := &Order{
		OrderNo:     fmt.Sprintf("EDU%d", time.Now().UnixNano()),
		UserID:      userID,
		TotalAmount: totalAmount,
		PayAmount:   totalAmount,
		Status:      1, // 待付款
		ExpiredAt:   &[]time.Time{time.Now().Add(30 * time.Minute)}[0], // 30分钟后过期
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建订单项
	for _, course := range courses {
		orderItem := OrderItem{
			OrderID:       order.ID,
			CourseID:      course.ID,
			CourseName:    course.Title,
			Price:         course.Price,
			OriginalPrice: course.OriginalPrice,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return order, nil
}

// GetOrdersByUserID 获取用户订单列表
func (s *OrderService) GetOrdersByUserID(userID uint, page, pageSize int) ([]Order, int64, error) {
	var orders []Order
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := s.db.Model(&Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := s.db.Preload("Items.Course").Where("user_id = ?", userID).
		Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&orders).Error

	return orders, total, err
}

// LearningService 学习服务
type LearningService struct {
	db *gorm.DB
}

// NewLearningService 创建学习服务
func NewLearningService(db *gorm.DB) *LearningService {
	return &LearningService{db: db}
}

// UpdateProgress 更新学习进度
func (s *LearningService) UpdateProgress(userID, courseID, lessonID uint, progress, watchTime int) error {
	// 查找或创建学习进度记录
	var learningProgress LearningProgress
	err := s.db.Where("user_id = ? AND course_id = ? AND lesson_id = ?", userID, courseID, lessonID).
		First(&learningProgress).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		learningProgress = LearningProgress{
			UserID:   userID,
			CourseID: courseID,
			LessonID: lessonID,
			Progress: progress,
			WatchTime: watchTime,
		}
		if progress >= 100 {
			learningProgress.IsCompleted = true
			now := time.Now()
			learningProgress.CompletedAt = &now
		}
		return s.db.Create(&learningProgress).Error
	} else if err != nil {
		return err
	}

	// 更新现有记录
	learningProgress.Progress = progress
	learningProgress.WatchTime = watchTime
	if progress >= 100 && !learningProgress.IsCompleted {
		learningProgress.IsCompleted = true
		now := time.Now()
		learningProgress.CompletedAt = &now
	}

	return s.db.Save(&learningProgress).Error
}

// GetUserCourseProgress 获取用户课程学习进度
func (s *LearningService) GetUserCourseProgress(userID, courseID uint) ([]LearningProgress, error) {
	var progress []LearningProgress
	err := s.db.Preload("Lesson").Where("user_id = ? AND course_id = ?", userID, courseID).
		Order("lesson_id").Find(&progress).Error
	return progress, err
}

// ========== API控制器 ==========

// APIResponse API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// UserController 用户控制器
type UserController struct {
	userService *UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUsers 获取用户列表
func (c *UserController) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	users, total, err := c.userService.GetUsers(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取用户列表失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data: PaginationResponse{
			List:  users,
			Total: total,
			Page:  page,
			Size:  pageSize,
		},
	})
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    user,
	})
}

// CourseController 课程控制器
type CourseController struct {
	courseService *CourseService
}

// NewCourseController 创建课程控制器
func NewCourseController(courseService *CourseService) *CourseController {
	return &CourseController{courseService: courseService}
}

// GetCourses 获取课程列表
func (c *CourseController) GetCourses(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	categoryIDStr := ctx.Query("category_id")

	var categoryID *uint
	if categoryIDStr != "" {
		id, _ := strconv.ParseUint(categoryIDStr, 10, 32)
		categoryID = (*uint)(&id)
	}

	courses, total, err := c.courseService.GetCourses(page, pageSize, categoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取课程列表失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data: PaginationResponse{
			List:  courses,
			Total: total,
			Page:  page,
			Size:  pageSize,
		},
	})
}

// GetCourse 获取课程详情
func (c *CourseController) GetCourse(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	course, err := c.courseService.GetCourseByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "课程不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    course,
	})
}

// OrderController 订单控制器
type OrderController struct {
	orderService *OrderService
}

// NewOrderController 创建订单控制器
func NewOrderController(orderService *OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CourseIDs []uint `json:"course_ids" binding:"required"`
}

// CreateOrder 创建订单
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	// 这里应该从JWT token中获取用户ID，简化处理直接使用1
	userID := uint(1)

	order, err := c.orderService.CreateOrder(userID, req.CourseIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "订单创建成功",
		Data:    order,
	})
}

// GetOrders 获取订单列表
func (c *OrderController) GetOrders(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// 这里应该从JWT token中获取用户ID，简化处理直接使用1
	userID := uint(1)

	orders, total, err := c.orderService.GetOrdersByUserID(userID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取订单列表失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data: PaginationResponse{
			List:  orders,
			Total: total,
			Page:  page,
			Size:  pageSize,
		},
	})
}

// ========== 数据初始化 ==========

// SeedData 填充测试数据
func SeedData(db *gorm.DB) error {
	fmt.Println("开始填充测试数据...")

	// 创建角色
	roles := []Role{
		{Name: "admin", Description: "管理员"},
		{Name: "instructor", Description: "讲师"},
		{Name: "student", Description: "学生"},
	}
	db.Create(&roles)

	// 创建用户
	users := []User{
		{Username: "admin", Email: "admin@example.com", Phone: "13800138001", Password: "password", Nickname: "管理员", RoleID: roles[0].ID},
		{Username: "instructor1", Email: "instructor1@example.com", Phone: "13800138002", Password: "password", Nickname: "讲师1", RoleID: roles[1].ID},
		{Username: "student1", Email: "student1@example.com", Phone: "13800138003", Password: "password", Nickname: "学生1", RoleID: roles[2].ID},
	}
	db.Create(&users)

	// 创建用户资料
	profiles := []UserProfile{
		{UserID: users[0].ID, RealName: "管理员", Gender: 1},
		{UserID: users[1].ID, RealName: "张老师", Gender: 1},
		{UserID: users[2].ID, RealName: "李同学", Gender: 2},
	}
	db.Create(&profiles)

	// 创建分类
	categories := []Category{
		{Name: "编程开发", Slug: "programming", Description: "编程开发相关课程"},
		{Name: "设计创意", Slug: "design", Description: "设计创意相关课程"},
		{Name: "产品运营", Slug: "product", Description: "产品运营相关课程"},
	}
	db.Create(&categories)

	// 创建课程
	courses := []Course{
		{
			Title:        "Go语言入门到精通",
			Slug:         "golang-tutorial",
			Description:  "从零开始学习Go语言，掌握现代编程技能",
			CategoryID:   categories[0].ID,
			InstructorID: users[1].ID,
			Price:        19900, // 199元
			OriginalPrice: 29900, // 原价299元
			Level:        1,
			Duration:     1200, // 20小时
			Status:       2,    // 已发布
		},
		{
			Title:        "React前端开发实战",
			Slug:         "react-tutorial",
			Description:  "学习React框架，构建现代化前端应用",
			CategoryID:   categories[0].ID,
			InstructorID: users[1].ID,
			Price:        24900, // 249元
			OriginalPrice: 34900, // 原价349元
			Level:        2,
			Duration:     1800, // 30小时
			Status:       2,    // 已发布
		},
	}
	db.Create(&courses)

	// 创建章节
	chapters := []Chapter{
		{CourseID: courses[0].ID, Title: "Go语言基础", Sort: 1},
		{CourseID: courses[0].ID, Title: "Go语言进阶", Sort: 2},
		{CourseID: courses[1].ID, Title: "React基础", Sort: 1},
		{CourseID: courses[1].ID, Title: "React进阶", Sort: 2},
	}
	db.Create(&chapters)

	// 创建课时
	lessons := []Lesson{
		{ChapterID: chapters[0].ID, Title: "Go语言介绍", Duration: 600, Sort: 1, IsFree: true},
		{ChapterID: chapters[0].ID, Title: "变量和数据类型", Duration: 900, Sort: 2},
		{ChapterID: chapters[1].ID, Title: "并发编程", Duration: 1200, Sort: 1},
		{ChapterID: chapters[2].ID, Title: "React介绍", Duration: 600, Sort: 1, IsFree: true},
		{ChapterID: chapters[2].ID, Title: "组件和Props", Duration: 900, Sort: 2},
		{ChapterID: chapters[3].ID, Title: "状态管理", Duration: 1200, Sort: 1},
	}
	db.Create(&lessons)

	fmt.Println("测试数据填充完成")
	return nil
}

// ========== 路由设置 ==========

// SetupRoutes 设置路由
func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 创建服务实例
	userService := NewUserService(db)
	courseService := NewCourseService(db)
	orderService := NewOrderService(db)

	// 创建控制器实例
	userController := NewUserController(userService)
	courseController := NewCourseController(courseService)
	orderController := NewOrderController(orderService)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
		}

		// 课程相关路由
		courses := api.Group("/courses")
		{
			courses.GET("", courseController.GetCourses)
			courses.GET("/:id", courseController.GetCourse)
		}

		// 订单相关路由
		orders := api.Group("/orders")
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("", orderController.GetOrders)
		}
	}

	return r
}

func main() {
	// 数据库配置
	config := DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "gorm_advanced_exercise5",
		Charset:  "utf8mb4",
	}

	// 连接数据库
	fmt.Println("连接数据库...")
	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 迁移数据库
	fmt.Println("迁移数据库...")
	db.AutoMigrate(
		&Role{}, &User{}, &UserProfile{}, &Category{}, &Course{},
		&Chapter{}, &Lesson{}, &Order{}, &OrderItem{}, &LearningProgress{},
	)

	// 检查是否需要填充测试数据
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		if err := SeedData(db); err != nil {
			log.Fatal("填充测试数据失败:", err)
		}
	}

	// 设置路由
	r := SetupRoutes(db)

	// 启动服务器
	fmt.Println("\n=== 在线教育平台后端系统启动 ===")
	fmt.Println("服务器地址: http://localhost:8080")
	fmt.Println("\nAPI接口:")
	fmt.Println("- GET  /api/v1/users        - 获取用户列表")
	fmt.Println("- GET  /api/v1/users/:id    - 获取用户详情")
	fmt.Println("- GET  /api/v1/courses      - 获取课程列表")
	fmt.Println("- GET  /api/v1/courses/:id  - 获取课程详情")
	fmt.Println("- POST /api/v1/orders       - 创建订单")
	fmt.Println("- GET  /api/v1/orders       - 获取订单列表")
	fmt.Println("\n强化练习任务:")
	fmt.Println("1. JWT认证和权限控制")
	fmt.Println("2. Redis缓存集成")
	fmt.Println("3. 文件上传和存储")
	fmt.Println("4. 支付接口集成")
	fmt.Println("5. 实时通知和WebSocket")
	fmt.Println("6. API文档生成（Swagger）")
	fmt.Println("7. 单元测试和集成测试")
	fmt.Println("8. Docker容器化部署")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
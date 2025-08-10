package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"../models"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *models.User) error {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	s.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已存在")
	}

	// 检查手机号是否已存在
	if user.Phone != "" {
		s.db.Model(&models.User{}).Where("phone = ?", user.Phone).Count(&count)
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}

	return s.db.Create(user).Error
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.Preload("Role").Preload("Profile").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateLastLogin 更新最后登录时间
func (s *UserService) UpdateLastLogin(id uint, ip string) error {
	now := time.Now()
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": &now,
		"login_ip":      ip,
	}).Error
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(page, pageSize int, filters map[string]interface{}) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "role_id":
			query = query.Where("role_id = ?", value)
		case "keyword":
			query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", 
				fmt.Sprintf("%%%v%%", value), fmt.Sprintf("%%%v%%", value), fmt.Sprintf("%%%v%%", value))
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Role").Preload("Profile").
		Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&users).Error

	return users, total, err
}

// CourseService 课程服务
type CourseService struct {
	db *gorm.DB
}

// NewCourseService 创建课程服务
func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{db: db}
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(course *models.Course) error {
	// 检查课程标识是否已存在
	var count int64
	s.db.Model(&models.Course{}).Where("slug = ?", course.Slug).Count(&count)
	if count > 0 {
		return errors.New("课程标识已存在")
	}

	return s.db.Create(course).Error
}

// GetCourseByID 根据ID获取课程详情
func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	var course models.Course
	err := s.db.Preload("Category").Preload("Instructor").
		Preload("Chapters.Lessons").First(&course, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("课程不存在")
		}
		return nil, err
	}

	// 增加浏览次数
	s.db.Model(&course).Update("view_count", gorm.Expr("view_count + ?", 1))

	return &course, nil
}

// GetCourses 获取课程列表
func (s *CourseService) GetCourses(page, pageSize int, filters map[string]interface{}) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := s.db.Model(&models.Course{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "category_id":
			query = query.Where("category_id = ?", value)
		case "instructor_id":
			query = query.Where("instructor_id = ?", value)
		case "level":
			query = query.Where("level = ?", value)
		case "is_free":
			query = query.Where("is_free = ?", value)
		case "is_recommend":
			query = query.Where("is_recommend = ?", value)
		case "keyword":
			query = query.Where("title LIKE ? OR subtitle LIKE ?", 
				fmt.Sprintf("%%%v%%", value), fmt.Sprintf("%%%v%%", value))
		case "price_min":
			query = query.Where("price >= ?", value)
		case "price_max":
			query = query.Where("price <= ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize

	// 排序
	orderBy := "created_at DESC"
	if sort, ok := filters["sort"]; ok {
		switch sort {
		case "price_asc":
			orderBy = "price ASC"
		case "price_desc":
			orderBy = "price DESC"
		case "rating":
			orderBy = "rating DESC"
		case "students":
			orderBy = "student_count DESC"
		case "newest":
			orderBy = "created_at DESC"
		}
	}

	err := query.Preload("Category").Preload("Instructor").
		Order(orderBy).Limit(pageSize).Offset(offset).Find(&courses).Error

	return courses, total, err
}

// UpdateCourse 更新课程信息
func (s *CourseService) UpdateCourse(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Course{}).Where("id = ?", id).Updates(updates).Error
}

// PublishCourse 发布课程
func (s *CourseService) PublishCourse(id uint) error {
	now := time.Now()
	return s.db.Model(&models.Course{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       2, // 发布状态
		"published_at": &now,
	}).Error
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
func (s *OrderService) CreateOrder(userID uint, courseIDs []uint, couponCode string) (*models.Order, error) {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询课程信息
	var courses []models.Course
	if err := tx.Where("id IN ? AND status = ?", courseIDs, 2).Find(&courses).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(courses) != len(courseIDs) {
		tx.Rollback()
		return nil, errors.New("部分课程不存在或已下架")
	}

	// 检查用户是否已购买过这些课程
	var existingOrders []models.Order
	tx.Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.course_id IN ? AND orders.status IN ?", 
			userID, courseIDs, []int{2, 3}).Find(&existingOrders)

	if len(existingOrders) > 0 {
		tx.Rollback()
		return nil, errors.New("您已购买过部分课程")
	}

	// 计算总金额
	var totalAmount int64
	for _, course := range courses {
		totalAmount += course.Price
	}

	// 处理优惠券
	var coupon *models.Coupon
	var discountAmount int64
	if couponCode != "" {
		if err := tx.Where("code = ? AND status = ? AND start_time <= ? AND end_time >= ? AND used_count < total_count", 
			couponCode, 1, time.Now(), time.Now()).First(&coupon).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return nil, errors.New("优惠券不存在或已失效")
			}
			tx.Rollback()
			return nil, err
		}

		// 检查最低消费金额
		if totalAmount < coupon.MinAmount {
			tx.Rollback()
			return nil, fmt.Errorf("订单金额不满足优惠券使用条件，最低消费%.2f元", float64(coupon.MinAmount)/100)
		}

		// 计算优惠金额
		if coupon.Type == 1 { // 满减券
			discountAmount = coupon.Value
		} else { // 折扣券
			discountAmount = totalAmount * (100 - coupon.Value) / 100
		}

		// 检查最大优惠金额
		if coupon.MaxAmount > 0 && discountAmount > coupon.MaxAmount {
			discountAmount = coupon.MaxAmount
		}

		// 更新优惠券使用次数
		if err := tx.Model(coupon).Update("used_count", gorm.Expr("used_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	payAmount := totalAmount - discountAmount
	if payAmount < 0 {
		payAmount = 0
	}

	// 创建订单
	order := &models.Order{
		OrderNo:        s.generateOrderNo(),
		UserID:         userID,
		TotalAmount:    totalAmount,
		PayAmount:      payAmount,
		DiscountAmount: discountAmount,
		Status:         1, // 待付款
		ExpiredAt:      &[]time.Time{time.Now().Add(30 * time.Minute)}[0], // 30分钟后过期
	}

	if coupon != nil {
		order.CouponID = &coupon.ID
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建订单项
	for _, course := range courses {
		orderItem := models.OrderItem{
			OrderID:       order.ID,
			CourseID:      course.ID,
			CourseName:    course.Title,
			CourseImage:   course.Cover,
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

// generateOrderNo 生成订单号
func (s *OrderService) generateOrderNo() string {
	return fmt.Sprintf("EDU%d", time.Now().UnixNano())
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(orderNo, paymentMethod, paymentNo string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找订单
	var order models.Order
	if err := tx.Where("order_no = ? AND status = ?", orderNo, 1).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("订单不存在或状态异常")
		}
		tx.Rollback()
		return err
	}

	// 检查订单是否过期
	if order.ExpiredAt != nil && time.Now().After(*order.ExpiredAt) {
		// 自动取消过期订单
		now := time.Now()
		tx.Model(&order).Updates(map[string]interface{}{
			"status":       4, // 已取消
			"cancelled_at": &now,
		})
		tx.Rollback()
		return errors.New("订单已过期")
	}

	// 更新订单状态
	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":         2, // 已付款
		"payment_method": paymentMethod,
		"payment_no":     paymentNo,
		"paid_at":        &now,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新课程学生数量
	var orderItems []models.OrderItem
	tx.Where("order_id = ?", order.ID).Find(&orderItems)
	for _, item := range orderItems {
		tx.Model(&models.Course{}).Where("id = ?", item.CourseID).
			Update("student_count", gorm.Expr("student_count + ?", 1))
	}

	tx.Commit()
	return nil
}

// GetOrdersByUserID 获取用户订单列表
func (s *OrderService) GetOrdersByUserID(userID uint, page, pageSize int, status *int8) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	offset := (page - 1) * pageSize

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("Items.Course").Preload("Coupon").
		Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&orders).Error

	return orders, total, err
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderNo string, userID uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找订单
	var order models.Order
	if err := tx.Where("order_no = ? AND user_id = ? AND status = ?", orderNo, userID, 1).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("订单不存在或无法取消")
		}
		tx.Rollback()
		return err
	}

	// 更新订单状态
	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":       4, // 已取消
		"cancelled_at": &now,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 回滚优惠券使用次数
	if order.CouponID != nil {
		tx.Model(&models.Coupon{}).Where("id = ?", *order.CouponID).
			Update("used_count", gorm.Expr("used_count - ?", 1))
	}

	tx.Commit()
	return nil
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
	// 检查用户是否有权限学习该课程
	var count int64
	s.db.Table("orders").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.course_id = ? AND orders.status IN ?", 
			userID, courseID, []int{2, 3}).Count(&count)

	if count == 0 {
		// 检查是否是免费课程或免费课时
		var lesson models.Lesson
		if err := s.db.Where("id = ? AND (is_free = ? OR EXISTS (SELECT 1 FROM courses WHERE id = ? AND is_free = ?))", 
			lessonID, true, courseID, true).First(&lesson).Error; err != nil {
			return errors.New("您没有权限学习该课程")
		}
	}

	// 查找或创建学习进度记录
	var learningProgress models.LearningProgress
	err := s.db.Where("user_id = ? AND course_id = ? AND lesson_id = ?", userID, courseID, lessonID).
		First(&learningProgress).Error

	now := time.Now()
	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		learningProgress = models.LearningProgress{
			UserID:      userID,
			CourseID:    courseID,
			LessonID:    lessonID,
			Progress:    progress,
			WatchTime:   watchTime,
			LastWatchAt: &now,
		}
		if progress >= 100 {
			learningProgress.IsCompleted = true
			learningProgress.CompletedAt = &now
		}
		return s.db.Create(&learningProgress).Error
	} else if err != nil {
		return err
	}

	// 更新现有记录
	updates := map[string]interface{}{
		"progress":      progress,
		"watch_time":    watchTime,
		"last_watch_at": &now,
	}

	if progress >= 100 && !learningProgress.IsCompleted {
		updates["is_completed"] = true
		updates["completed_at"] = &now
	}

	return s.db.Model(&learningProgress).Updates(updates).Error
}

// GetUserCourseProgress 获取用户课程学习进度
func (s *LearningService) GetUserCourseProgress(userID, courseID uint) ([]models.LearningProgress, error) {
	var progress []models.LearningProgress
	err := s.db.Preload("Lesson").Where("user_id = ? AND course_id = ?", userID, courseID).
		Order("lesson_id").Find(&progress).Error
	return progress, err
}

// GetUserLearningCourses 获取用户学习的课程列表
func (s *LearningService) GetUserLearningCourses(userID uint, page, pageSize int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	// 子查询：获取用户已购买的课程ID
	subQuery := s.db.Table("orders").
		Select("DISTINCT order_items.course_id").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND orders.status IN ?", userID, []int{2, 3})

	query := s.db.Model(&models.Course{}).Where("id IN (?)", subQuery)

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
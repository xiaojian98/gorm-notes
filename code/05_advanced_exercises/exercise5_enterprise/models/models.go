package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string       `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
	Email       string       `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	Phone       string       `gorm:"uniqueIndex;size:20" json:"phone" validate:"omitempty,len=11"`
	Password    string       `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Nickname    string       `gorm:"size:50" json:"nickname" validate:"omitempty,max=50"`
	Avatar      string       `gorm:"size:255" json:"avatar"`
	Status      int8         `gorm:"default:1;comment:1-正常,2-禁用" json:"status"`
	RoleID      uint         `gorm:"index;not null" json:"role_id" validate:"required"`
	LastLoginAt *time.Time   `json:"last_login_at"`
	LoginIP     string       `gorm:"size:45" json:"login_ip"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
	
	// 关联
	Role            Role             `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Profile         UserProfile      `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Orders          []Order          `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	LearningProgress []LearningProgress `gorm:"foreignKey:UserID" json:"learning_progress,omitempty"`
	Courses         []Course         `gorm:"foreignKey:InstructorID" json:"courses,omitempty"`
	Reviews         []CourseReview   `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 这里可以添加密码加密等逻辑
	return nil
}

// AfterCreate GORM钩子：创建后
func (u *User) AfterCreate(tx *gorm.DB) error {
	// 创建用户资料
	profile := UserProfile{
		UserID: u.ID,
	}
	return tx.Create(&profile).Error
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name" validate:"required,max=50"`
	Description string `gorm:"size:255" json:"description" validate:"omitempty,max=255"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	Permissions string `gorm:"type:text" json:"permissions"` // JSON格式存储权限
	
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
	UserID      uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	RealName    string     `gorm:"size:50" json:"real_name" validate:"omitempty,max=50"`
	Gender      int8       `gorm:"default:0;comment:0-未知,1-男,2-女" json:"gender" validate:"omitempty,oneof=0 1 2"`
	Birthday    *time.Time `json:"birthday"`
	Bio         string     `gorm:"type:text" json:"bio" validate:"omitempty,max=500"`
	Location    string     `gorm:"size:100" json:"location" validate:"omitempty,max=100"`
	Website     string     `gorm:"size:255" json:"website" validate:"omitempty,url"`
	Company     string     `gorm:"size:100" json:"company" validate:"omitempty,max=100"`
	Position    string     `gorm:"size:100" json:"position" validate:"omitempty,max=100"`
	Education   string     `gorm:"size:100" json:"education" validate:"omitempty,max=100"`
	Experience  int        `gorm:"default:0;comment:工作经验(年)" json:"experience"`
	
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
	Name        string `gorm:"size:50;not null" json:"name" validate:"required,max=50"`
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug" validate:"required,max=100"`
	Description string `gorm:"type:text" json:"description" validate:"omitempty,max=500"`
	Icon        string `gorm:"size:255" json:"icon"`
	Cover       string `gorm:"size:255" json:"cover"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	CourseCount int    `gorm:"default:0;comment:课程数量" json:"course_count"`
	
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
	Title         string     `gorm:"size:255;not null" json:"title" validate:"required,max=255"`
	Slug          string     `gorm:"uniqueIndex;size:255;not null" json:"slug" validate:"required,max=255"`
	Subtitle      string     `gorm:"size:500" json:"subtitle" validate:"omitempty,max=500"`
	Description   string     `gorm:"type:text" json:"description" validate:"omitempty,max=2000"`
	Content       string     `gorm:"type:longtext" json:"content"` // 详细内容
	Cover         string     `gorm:"size:255" json:"cover"`
	Video         string     `gorm:"size:500" json:"video"` // 课程预览视频
	CategoryID    uint       `gorm:"index;not null" json:"category_id" validate:"required"`
	InstructorID  uint       `gorm:"index;not null" json:"instructor_id" validate:"required"`
	Price         int64      `gorm:"not null;comment:价格(分)" json:"price" validate:"min=0"`
	OriginalPrice int64      `gorm:"default:0;comment:原价(分)" json:"original_price" validate:"min=0"`
	Level         int8       `gorm:"default:1;comment:1-初级,2-中级,3-高级" json:"level" validate:"oneof=1 2 3"`
	Duration      int        `gorm:"default:0;comment:课程时长(分钟)" json:"duration"`
	StudentCount  int        `gorm:"default:0;comment:学生数量" json:"student_count"`
	LessonCount   int        `gorm:"default:0;comment:课时数量" json:"lesson_count"`
	Rating        float32    `gorm:"default:0;comment:评分" json:"rating"`
	ReviewCount   int        `gorm:"default:0;comment:评价数量" json:"review_count"`
	ViewCount     int        `gorm:"default:0;comment:浏览次数" json:"view_count"`
	FavoriteCount int        `gorm:"default:0;comment:收藏次数" json:"favorite_count"`
	Status        int8       `gorm:"default:1;comment:1-草稿,2-发布,3-下架" json:"status" validate:"oneof=1 2 3"`
	IsFree        bool       `gorm:"default:false;comment:是否免费" json:"is_free"`
	IsRecommend   bool       `gorm:"default:false;comment:是否推荐" json:"is_recommend"`
	PublishedAt   *time.Time `json:"published_at"`
	Tags          string     `gorm:"size:500" json:"tags"` // 标签，逗号分隔
	Requirements  string     `gorm:"type:text" json:"requirements"` // 学习要求
	Goals         string     `gorm:"type:text" json:"goals"` // 学习目标
	
	// 关联
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Instructor  User           `gorm:"foreignKey:InstructorID" json:"instructor,omitempty"`
	Chapters    []Chapter      `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"chapters,omitempty"`
	Orders      []Order        `gorm:"many2many:order_items;" json:"orders,omitempty"`
	Reviews     []CourseReview `gorm:"foreignKey:CourseID" json:"reviews,omitempty"`
	Favorites   []CourseFavorite `gorm:"foreignKey:CourseID" json:"favorites,omitempty"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// Chapter 章节模型
type Chapter struct {
	BaseModel
	CourseID    uint   `gorm:"index;not null" json:"course_id" validate:"required"`
	Title       string `gorm:"size:255;not null" json:"title" validate:"required,max=255"`
	Description string `gorm:"type:text" json:"description" validate:"omitempty,max=1000"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status" validate:"oneof=1 2"`
	LessonCount int    `gorm:"default:0;comment:课时数量" json:"lesson_count"`
	Duration    int    `gorm:"default:0;comment:章节时长(分钟)" json:"duration"`
	
	// 关联
	Course  Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
}

// TableName 指定表名
func (Chapter) TableName() string {
	return "chapters"
}

// Lesson 课时模型
type Lesson struct {
	BaseModel
	ChapterID   uint   `gorm:"index;not null" json:"chapter_id" validate:"required"`
	Title       string `gorm:"size:255;not null" json:"title" validate:"required,max=255"`
	Description string `gorm:"type:text" json:"description" validate:"omitempty,max=1000"`
	Content     string `gorm:"type:longtext" json:"content"` // 课时内容
	VideoURL    string `gorm:"size:500" json:"video_url"`
	VideoSize   int64  `gorm:"default:0;comment:视频大小(字节)" json:"video_size"`
	Attachments string `gorm:"type:text" json:"attachments"` // 附件，JSON格式
	Duration    int    `gorm:"default:0;comment:时长(秒)" json:"duration"`
	Sort        int    `gorm:"default:0" json:"sort"`
	IsFree      bool   `gorm:"default:false;comment:是否免费" json:"is_free"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status" validate:"oneof=1 2"`
	ViewCount   int    `gorm:"default:0;comment:观看次数" json:"view_count"`
	
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
	UserID         uint       `gorm:"index;not null" json:"user_id" validate:"required"`
	TotalAmount    int64      `gorm:"not null;comment:总金额(分)" json:"total_amount" validate:"min=0"`
	PayAmount      int64      `gorm:"not null;comment:实付金额(分)" json:"pay_amount" validate:"min=0"`
	DiscountAmount int64      `gorm:"default:0;comment:优惠金额(分)" json:"discount_amount" validate:"min=0"`
	CouponID       *uint      `gorm:"index" json:"coupon_id"`
	Status         int8       `gorm:"index;default:1;comment:1-待付款,2-已付款,3-已完成,4-已取消,5-已退款" json:"status" validate:"oneof=1 2 3 4 5"`
	PaymentMethod  string     `gorm:"size:50" json:"payment_method"`
	PaymentNo      string     `gorm:"size:100" json:"payment_no"`
	PaidAt         *time.Time `json:"paid_at"`
	ExpiredAt      *time.Time `json:"expired_at"`
	CancelledAt    *time.Time `json:"cancelled_at"`
	RefundedAt     *time.Time `json:"refunded_at"`
	Remark         string     `gorm:"type:text" json:"remark" validate:"omitempty,max=500"`
	RefundReason   string     `gorm:"type:text" json:"refund_reason"`
	
	// 关联
	User    User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items   []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	Courses []Course    `gorm:"many2many:order_items;" json:"courses,omitempty"`
	Coupon  *Coupon     `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	BaseModel
	OrderID       uint   `gorm:"index;not null" json:"order_id" validate:"required"`
	CourseID      uint   `gorm:"index;not null" json:"course_id" validate:"required"`
	CourseName    string `gorm:"size:255;not null" json:"course_name" validate:"required,max=255"`
	CourseImage   string `gorm:"size:255" json:"course_image"`
	Price         int64  `gorm:"not null;comment:价格(分)" json:"price" validate:"min=0"`
	OriginalPrice int64  `gorm:"default:0;comment:原价(分)" json:"original_price" validate:"min=0"`
	DiscountAmount int64 `gorm:"default:0;comment:优惠金额(分)" json:"discount_amount" validate:"min=0"`
	
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
	UserID      uint `gorm:"index;not null" json:"user_id" validate:"required"`
	CourseID    uint `gorm:"index;not null" json:"course_id" validate:"required"`
	LessonID    uint `gorm:"index;not null" json:"lesson_id" validate:"required"`
	Progress    int  `gorm:"default:0;comment:进度百分比" json:"progress" validate:"min=0,max=100"`
	WatchTime   int  `gorm:"default:0;comment:观看时长(秒)" json:"watch_time" validate:"min=0"`
	IsCompleted bool `gorm:"default:false;comment:是否完成" json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
	LastWatchAt *time.Time `json:"last_watch_at"`
	
	// 关联
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Lesson Lesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
}

// TableName 指定表名
func (LearningProgress) TableName() string {
	return "learning_progress"
}

// CourseReview 课程评价模型
type CourseReview struct {
	BaseModel
	UserID   uint    `gorm:"index;not null" json:"user_id" validate:"required"`
	CourseID uint    `gorm:"index;not null" json:"course_id" validate:"required"`
	Rating   float32 `gorm:"not null;comment:评分(1-5)" json:"rating" validate:"required,min=1,max=5"`
	Content  string  `gorm:"type:text" json:"content" validate:"omitempty,max=1000"`
	Status   int8    `gorm:"default:1;comment:1-正常,2-隐藏" json:"status" validate:"oneof=1 2"`
	LikeCount int    `gorm:"default:0;comment:点赞数" json:"like_count"`
	
	// 关联
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName 指定表名
func (CourseReview) TableName() string {
	return "course_reviews"
}

// CourseFavorite 课程收藏模型
type CourseFavorite struct {
	BaseModel
	UserID   uint `gorm:"index;not null" json:"user_id" validate:"required"`
	CourseID uint `gorm:"index;not null" json:"course_id" validate:"required"`
	
	// 关联
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName 指定表名
func (CourseFavorite) TableName() string {
	return "course_favorites"
}

// Coupon 优惠券模型
type Coupon struct {
	BaseModel
	Name        string     `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
	Code        string     `gorm:"uniqueIndex;size:50;not null" json:"code" validate:"required,max=50"`
	Type        int8       `gorm:"not null;comment:1-满减券,2-折扣券" json:"type" validate:"required,oneof=1 2"`
	Value       int64      `gorm:"not null;comment:优惠值(分或折扣*100)" json:"value" validate:"required,min=1"`
	MinAmount   int64      `gorm:"default:0;comment:最低消费金额(分)" json:"min_amount" validate:"min=0"`
	MaxAmount   int64      `gorm:"default:0;comment:最大优惠金额(分)" json:"max_amount" validate:"min=0"`
	TotalCount  int        `gorm:"not null;comment:总数量" json:"total_count" validate:"required,min=1"`
	UsedCount   int        `gorm:"default:0;comment:已使用数量" json:"used_count"`
	StartTime   time.Time  `gorm:"not null" json:"start_time" validate:"required"`
	EndTime     time.Time  `gorm:"not null" json:"end_time" validate:"required"`
	Status      int8       `gorm:"default:1;comment:1-启用,2-禁用" json:"status" validate:"oneof=1 2"`
	Description string     `gorm:"type:text" json:"description" validate:"omitempty,max=500"`
	
	// 关联
	Orders []Order `gorm:"foreignKey:CouponID" json:"orders,omitempty"`
}

// TableName 指定表名
func (Coupon) TableName() string {
	return "coupons"
}

// Notification 通知模型
type Notification struct {
	BaseModel
	UserID   uint   `gorm:"index;not null" json:"user_id" validate:"required"`
	Title    string `gorm:"size:255;not null" json:"title" validate:"required,max=255"`
	Content  string `gorm:"type:text" json:"content" validate:"omitempty,max=1000"`
	Type     int8   `gorm:"not null;comment:1-系统通知,2-课程通知,3-订单通知" json:"type" validate:"required,oneof=1 2 3"`
	IsRead   bool   `gorm:"default:false;comment:是否已读" json:"is_read"`
	ReadAt   *time.Time `json:"read_at"`
	Data     string `gorm:"type:text" json:"data"` // 额外数据，JSON格式
	
	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// SystemLog 系统日志模型
type SystemLog struct {
	BaseModel
	UserID    *uint  `gorm:"index" json:"user_id"`
	Action    string `gorm:"size:100;not null" json:"action" validate:"required,max=100"`
	Module    string `gorm:"size:50;not null" json:"module" validate:"required,max=50"`
	Method    string `gorm:"size:10;not null" json:"method" validate:"required,max=10"`
	URL       string `gorm:"size:500" json:"url" validate:"omitempty,max=500"`
	IP        string `gorm:"size:45" json:"ip"`
	UserAgent string `gorm:"size:500" json:"user_agent"`
	Request   string `gorm:"type:text" json:"request"`
	Response  string `gorm:"type:text" json:"response"`
	Status    int    `gorm:"not null" json:"status"`
	Duration  int64  `gorm:"not null;comment:耗时(毫秒)" json:"duration"`
	
	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (SystemLog) TableName() string {
	return "system_logs"
}
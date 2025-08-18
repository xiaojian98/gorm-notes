package handlers

import (
	"time"
)

// 通用响应结构体

// ErrorResponse 错误响应结构体
// 用于统一的错误信息返回格式
type ErrorResponse struct {
	Error     string    `json:"error"`     // 错误类型
	Message   string    `json:"message"`   // 错误详细信息
	Code      string    `json:"code,omitempty"` // 错误代码
	Timestamp time.Time `json:"timestamp"` // 时间戳
}

// SuccessResponse 成功响应结构体
// 用于统一的成功信息返回格式
type SuccessResponse struct {
	Message   string      `json:"message"`   // 成功信息
	Data      interface{} `json:"data,omitempty"` // 返回数据
	Timestamp time.Time   `json:"timestamp"` // 时间戳
}

// PaginationResponse 分页响应结构体
// 用于统一的分页信息返回格式
type PaginationResponse struct {
	Page       int `json:"page"`        // 当前页码
	PageSize   int `json:"page_size"`   // 每页数量
	Total      int64 `json:"total"`      // 总记录数
	TotalPages int `json:"total_pages"` // 总页数
}

// ListResponse 列表响应结构体
// 用于统一的列表数据返回格式
type ListResponse struct {
	Data       interface{}         `json:"data"`       // 列表数据
	Pagination PaginationResponse  `json:"pagination"` // 分页信息
	Timestamp  time.Time           `json:"timestamp"`  // 时间戳
}

// 用户相关响应结构体

// UserResponse 用户响应结构体
// 用于返回用户基本信息
type UserResponse struct {
	ID          uint      `json:"id"`          // 用户ID
	Username    string    `json:"username"`    // 用户名
	Nickname    string    `json:"nickname"`    // 昵称
	Email       string    `json:"email"`       // 邮箱
	Avatar      string    `json:"avatar"`      // 头像URL
	Bio         string    `json:"bio"`         // 个人简介
	Website     string    `json:"website"`     // 个人网站
	Location    string    `json:"location"`    // 所在地
	Status      string    `json:"status"`      // 状态
	Role        string    `json:"role"`        // 角色
	IsVerified  bool      `json:"is_verified"` // 是否已验证
	PostCount   int64     `json:"post_count"`  // 文章数量
	FollowerCount int64   `json:"follower_count"` // 粉丝数量
	FollowingCount int64  `json:"following_count"` // 关注数量
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

// UserProfileResponse 用户详细资料响应结构体
// 用于返回用户详细信息
type UserProfileResponse struct {
	User      UserResponse `json:"user"`      // 用户基本信息
	Stats     UserStats    `json:"stats"`     // 用户统计信息
	IsFollowing bool       `json:"is_following,omitempty"` // 是否已关注(仅在查看他人资料时返回)
}

// UserStats 用户统计信息
type UserStats struct {
	TotalPosts     int64 `json:"total_posts"`     // 总文章数
	TotalComments  int64 `json:"total_comments"`  // 总评论数
	TotalLikes     int64 `json:"total_likes"`     // 获得的总点赞数
	TotalViews     int64 `json:"total_views"`     // 文章总浏览量
	FollowerCount  int64 `json:"follower_count"`  // 粉丝数
	FollowingCount int64 `json:"following_count"` // 关注数
	JoinDays       int   `json:"join_days"`       // 加入天数
}

// 文章相关响应结构体

// PostResponse 文章响应结构体
// 用于返回文章信息
type PostResponse struct {
	ID          uint              `json:"id"`          // 文章ID
	Title       string            `json:"title"`       // 标题
	Slug        string            `json:"slug"`        // URL别名
	Content     string            `json:"content"`     // 内容
	Excerpt     string            `json:"excerpt"`     // 摘要
	Status      string            `json:"status"`      // 状态
	ViewCount   int64             `json:"view_count"`  // 浏览量
	LikeCount   int64             `json:"like_count"`  // 点赞数
	CommentCount int64            `json:"comment_count"` // 评论数
	WordCount   int               `json:"word_count"`  // 字数
	ReadingTime int               `json:"reading_time"` // 阅读时间(分钟)
	FeaturedImage string          `json:"featured_image"` // 特色图片
	Author      *UserResponse     `json:"author"`      // 作者信息
	Category    *CategoryResponse  `json:"category"`    // 分类信息
	Tags        []TagResponse     `json:"tags"`        // 标签列表
	Meta        map[string]string `json:"meta,omitempty"` // 元数据
	IsLiked     bool              `json:"is_liked,omitempty"` // 是否已点赞(仅在用户登录时返回)
	CreatedAt   time.Time         `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time         `json:"updated_at"`  // 更新时间
	PublishedAt *time.Time        `json:"published_at,omitempty"` // 发布时间
}

// PostSummaryResponse 文章摘要响应结构体
// 用于返回文章列表中的简要信息
type PostSummaryResponse struct {
	ID          uint             `json:"id"`          // 文章ID
	Title       string           `json:"title"`       // 标题
	Slug        string           `json:"slug"`        // URL别名
	Excerpt     string           `json:"excerpt"`     // 摘要
	Status      string           `json:"status"`      // 状态
	ViewCount   int64            `json:"view_count"`  // 浏览量
	LikeCount   int64            `json:"like_count"`  // 点赞数
	CommentCount int64           `json:"comment_count"` // 评论数
	ReadingTime int              `json:"reading_time"` // 阅读时间(分钟)
	FeaturedImage string         `json:"featured_image"` // 特色图片
	Author      *UserResponse    `json:"author"`      // 作者信息
	Category    *CategoryResponse `json:"category"`    // 分类信息
	Tags        []TagResponse    `json:"tags"`        // 标签列表
	CreatedAt   time.Time        `json:"created_at"`  // 创建时间
	PublishedAt *time.Time       `json:"published_at,omitempty"` // 发布时间
}

// 分类和标签响应结构体

// CategoryResponse 分类响应结构体
// 用于返回分类信息
type CategoryResponse struct {
	ID          uint   `json:"id"`          // 分类ID
	Name        string `json:"name"`        // 分类名称
	Slug        string `json:"slug"`        // URL别名
	Description string `json:"description"` // 描述
	Color       string `json:"color"`       // 颜色
	Icon        string `json:"icon"`        // 图标
	PostCount   int64  `json:"post_count"`  // 文章数量
	ParentID    uint   `json:"parent_id,omitempty"` // 父分类ID
	Children    []CategoryResponse `json:"children,omitempty"` // 子分类
	CreatedAt   time.Time `json:"created_at"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at"` // 更新时间
}

// TagResponse 标签响应结构体
// 用于返回标签信息
type TagResponse struct {
	ID          uint      `json:"id"`          // 标签ID
	Name        string    `json:"name"`        // 标签名称
	Slug        string    `json:"slug"`        // URL别名
	Description string    `json:"description"` // 描述
	Color       string    `json:"color"`       // 颜色
	PostCount   int64     `json:"post_count"`  // 文章数量
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

// 通用查询参数结构体

// PaginationParams 分页参数
// 用于解析分页相关的查询参数
type PaginationParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`      // 页码
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页数量
}

// SortParams 排序参数
// 用于解析排序相关的查询参数
type SortParams struct {
	Sort  string `form:"sort" binding:"omitempty"`  // 排序字段
	Order string `form:"order" binding:"omitempty,oneof=asc desc"` // 排序顺序
}

// SearchParams 搜索参数
// 用于解析搜索相关的查询参数
type SearchParams struct {
	Keyword    string `form:"keyword" binding:"omitempty,max=100"`    // 搜索关键词
	Category   string `form:"category" binding:"omitempty"`           // 分类筛选
	Tag        string `form:"tag" binding:"omitempty"`               // 标签筛选
	Author     string `form:"author" binding:"omitempty"`             // 作者筛选
	Status     string `form:"status" binding:"omitempty"`             // 状态筛选
	StartDate  string `form:"start_date" binding:"omitempty"`         // 开始日期
	EndDate    string `form:"end_date" binding:"omitempty"`           // 结束日期
}

// 常用的HTTP状态码常量
const (
	// 成功状态码
	StatusOK      = 200 // 请求成功
	StatusCreated = 201 // 创建成功
	StatusNoContent = 204 // 无内容

	// 客户端错误状态码
	StatusBadRequest   = 400 // 请求参数错误
	StatusUnauthorized = 401 // 未授权
	StatusForbidden    = 403 // 禁止访问
	StatusNotFound     = 404 // 资源不存在
	StatusConflict     = 409 // 资源冲突
	StatusUnprocessableEntity = 422 // 无法处理的实体
	StatusTooManyRequests = 429 // 请求过多

	// 服务器错误状态码
	StatusInternalServerError = 500 // 内部服务器错误
	StatusBadGateway         = 502 // 网关错误
	StatusServiceUnavailable = 503 // 服务不可用
	StatusGatewayTimeout     = 504 // 网关超时
)

// 常用的错误代码常量
const (
	// 通用错误代码
	ErrCodeValidation    = "VALIDATION_ERROR"    // 参数验证错误
	ErrCodeUnauthorized  = "UNAUTHORIZED"        // 未授权
	ErrCodeForbidden     = "FORBIDDEN"           // 禁止访问
	ErrCodeNotFound      = "NOT_FOUND"           // 资源不存在
	ErrCodeConflict      = "CONFLICT"            // 资源冲突
	ErrCodeInternalError = "INTERNAL_ERROR"      // 内部错误
	ErrCodeRateLimit     = "RATE_LIMIT_EXCEEDED" // 请求频率限制

	// 用户相关错误代码
	ErrCodeUserNotFound     = "USER_NOT_FOUND"     // 用户不存在
	ErrCodeUserExists       = "USER_EXISTS"        // 用户已存在
	ErrCodeInvalidPassword  = "INVALID_PASSWORD"   // 密码错误
	ErrCodeUserBanned       = "USER_BANNED"        // 用户被封禁
	ErrCodeUserNotVerified  = "USER_NOT_VERIFIED"  // 用户未验证

	// 文章相关错误代码
	ErrCodePostNotFound     = "POST_NOT_FOUND"     // 文章不存在
	ErrCodePostExists       = "POST_EXISTS"        // 文章已存在
	ErrCodePostNotPublished = "POST_NOT_PUBLISHED" // 文章未发布
	ErrCodePostDraft        = "POST_DRAFT"         // 文章为草稿

	// 评论相关错误代码
	ErrCodeCommentNotFound   = "COMMENT_NOT_FOUND"   // 评论不存在
	ErrCodeCommentNotApproved = "COMMENT_NOT_APPROVED" // 评论未审核
	ErrCodeCommentSpam       = "COMMENT_SPAM"        // 垃圾评论

	// 权限相关错误代码
	ErrCodeInsufficientPermission = "INSUFFICIENT_PERMISSION" // 权限不足
	ErrCodeOwnershipRequired      = "OWNERSHIP_REQUIRED"      // 需要所有权
	ErrCodeAdminRequired          = "ADMIN_REQUIRED"          // 需要管理员权限
)

// NewErrorResponse 创建错误响应
// 参数: error - 错误类型, message - 错误信息, code - 错误代码(可选)
// 返回: ErrorResponse - 错误响应结构体
func NewErrorResponse(error, message string, code ...string) ErrorResponse {
	response := ErrorResponse{
		Error:     error,
		Message:   message,
		Timestamp: time.Now(),
	}

	if len(code) > 0 {
		response.Code = code[0]
	}

	return response
}

// NewSuccessResponse 创建成功响应
// 参数: message - 成功信息, data - 返回数据(可选)
// 返回: SuccessResponse - 成功响应结构体
func NewSuccessResponse(message string, data ...interface{}) SuccessResponse {
	response := SuccessResponse{
		Message:   message,
		Timestamp: time.Now(),
	}

	if len(data) > 0 {
		response.Data = data[0]
	}

	return response
}

// NewListResponse 创建列表响应
// 参数: data - 列表数据, page - 当前页, pageSize - 每页数量, total - 总记录数
// 返回: ListResponse - 列表响应结构体
func NewListResponse(data interface{}, page, pageSize int, total int64) ListResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	return ListResponse{
		Data: data,
		Pagination: PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		Timestamp: time.Now(),
	}
}

// GetDefaultPagination 获取默认分页参数
// 参数: params - 分页参数
// 返回: int, int - 页码, 每页数量
func GetDefaultPagination(params PaginationParams) (int, int) {
	page := params.Page
	if page < 1 {
		page = 1
	}

	pageSize := params.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}

// GetDefaultSort 获取默认排序参数
// 参数: params - 排序参数, defaultSort - 默认排序字段, defaultOrder - 默认排序顺序
// 返回: string, string - 排序字段, 排序顺序
func GetDefaultSort(params SortParams, defaultSort, defaultOrder string) (string, string) {
	sort := params.Sort
	if sort == "" {
		sort = defaultSort
	}

	order := params.Order
	if order == "" {
		order = defaultOrder
	}

	return sort, order
}
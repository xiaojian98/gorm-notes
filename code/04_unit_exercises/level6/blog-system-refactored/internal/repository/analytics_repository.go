package repository

import (
	"database/sql"
	"errors"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// AnalyticsRepository 分析数据访问层接口
// 定义分析统计相关的数据库操作方法
type AnalyticsRepository interface {
	// 基本CRUD操作
	Create(analytics *models.Analytics) error                  // 创建分析记录
	GetByID(id uint) (*models.Analytics, error)               // 根据ID获取分析记录
	Update(analytics *models.Analytics) error                 // 更新分析记录
	Delete(id uint) error                                     // 删除分析记录
	
	// 仪表板统计
	GetDashboardStats() (*models.DashboardStats, error)      // 获取仪表板统计数据
	GetDashboardStatsForPeriod(startDate, endDate time.Time) (*models.DashboardStats, error) // 获取指定时期的仪表板统计
	
	// 内容统计
	GetContentStats() (*models.ContentStats, error)          // 获取内容统计
	GetPopularPosts(limit int, days int) ([]models.PopularPost, error) // 获取热门文章
	GetCategoryStats() ([]models.CategoryStats, error)       // 获取分类统计
	GetTagStats(limit int) ([]models.TagStats, error)        // 获取标签统计
	
	// 用户统计
	GetActiveUsers(limit int, days int) ([]models.ActiveUser, error) // 获取活跃用户
	GetUserGrowthStats(days int) ([]models.UserGrowthStats, error) // 获取用户增长统计
	GetUserEngagementStats(userID uint) (*UserEngagement, error) // 获取用户参与度统计
	
	// 趋势分析
	GetPostTrends(days int) ([]TrendData, error)             // 获取文章发布趋势
	GetCommentTrends(days int) ([]TrendData, error)          // 获取评论趋势
	GetViewTrends(days int) ([]TrendData, error)             // 获取浏览量趋势
	GetUserRegistrationTrends(days int) ([]TrendData, error) // 获取用户注册趋势
	
	// 性能分析
	GetTopPerformingContent(contentType string, metric string, limit int, days int) ([]PerformanceData, error) // 获取表现最佳的内容
	GetEngagementMetrics(days int) (*EngagementMetrics, error) // 获取参与度指标
	GetRetentionAnalysis(days int) (*RetentionAnalysis, error) // 获取留存分析
	
	// 实时统计
	GetRealTimeStats() (*RealTimeStats, error)               // 获取实时统计
	UpdateRealTimeStats(stats *RealTimeStats) error          // 更新实时统计
	
	// 自定义查询
	ExecuteCustomQuery(query string, params ...interface{}) ([]map[string]interface{}, error) // 执行自定义查询
	GetMetricsByDateRange(metric string, startDate, endDate time.Time) ([]MetricData, error) // 获取指定日期范围的指标数据
}

// analyticsRepository 分析数据访问层实现
type analyticsRepository struct {
	db *gorm.DB
}

// NewAnalyticsRepository 创建分析数据访问层实例
// 参数: db - 数据库连接
// 返回: AnalyticsRepository - 分析数据访问层接口实例
func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}

// 辅助数据结构

// TrendData 趋势数据
type TrendData struct {
	Date  time.Time `json:"date"`  // 日期
	Count int64     `json:"count"` // 数量
	Value float64   `json:"value"` // 值
}

// PerformanceData 性能数据
type PerformanceData struct {
	ID          uint    `json:"id"`          // ID
	Title       string  `json:"title"`       // 标题
	Type        string  `json:"type"`        // 类型
	MetricValue float64 `json:"metric_value"` // 指标值
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
}

// EngagementMetrics 参与度指标
type EngagementMetrics struct {
	TotalViews        int64   `json:"total_views"`        // 总浏览量
	TotalComments     int64   `json:"total_comments"`     // 总评论数
	TotalLikes        int64   `json:"total_likes"`        // 总点赞数
	TotalShares       int64   `json:"total_shares"`       // 总分享数
	EngagementRate    float64 `json:"engagement_rate"`    // 参与度
	AverageTimeOnPage float64 `json:"average_time_on_page"` // 平均页面停留时间
	BounceRate        float64 `json:"bounce_rate"`        // 跳出率
}

// RetentionAnalysis 留存分析
type RetentionAnalysis struct {
	Day1Retention  float64 `json:"day1_retention"`  // 1天留存率
	Day7Retention  float64 `json:"day7_retention"`  // 7天留存率
	Day30Retention float64 `json:"day30_retention"` // 30天留存率
	CohortData     []CohortData `json:"cohort_data"`    // 队列数据
}

// CohortData 队列数据
type CohortData struct {
	Cohort         string    `json:"cohort"`          // 队列标识
	RegisteredDate time.Time `json:"registered_date"` // 注册日期
	UserCount      int64     `json:"user_count"`      // 用户数量
	RetentionRates []float64 `json:"retention_rates"` // 留存率数组
}

// RealTimeStats 实时统计
type RealTimeStats struct {
	OnlineUsers     int64     `json:"online_users"`     // 在线用户数
	ActiveUsers     int64     `json:"active_users"`     // 活跃用户数
	CurrentViews    int64     `json:"current_views"`    // 当前浏览量
	TodayPosts      int64     `json:"today_posts"`      // 今日文章数
	TodayComments   int64     `json:"today_comments"`   // 今日评论数
	TodayUsers      int64     `json:"today_users"`      // 今日新用户数
	LastUpdated     time.Time `json:"last_updated"`     // 最后更新时间
}

// UserEngagement 用户参与度
type UserEngagement struct {
	UserID           uint    `json:"user_id"`           // 用户ID
	TotalPosts       int64   `json:"total_posts"`       // 总文章数
	TotalComments    int64   `json:"total_comments"`    // 总评论数
	TotalLikes       int64   `json:"total_likes"`       // 总点赞数
	TotalViews       int64   `json:"total_views"`       // 总浏览量
	EngagementScore  float64 `json:"engagement_score"`  // 参与度评分
	LastActiveDate   time.Time `json:"last_active_date"`  // 最后活跃日期
	RegistrationDate time.Time `json:"registration_date"` // 注册日期
}

// MetricData 指标数据
type MetricData struct {
	Date   time.Time `json:"date"`   // 日期
	Metric string    `json:"metric"` // 指标名称
	Value  float64   `json:"value"`  // 指标值
}

// 基本CRUD操作实现

// Create 创建分析记录
// 参数: analytics - 分析记录对象
// 返回: error - 错误信息
func (r *analyticsRepository) Create(analytics *models.Analytics) error {
	if analytics == nil {
		return errors.New("分析记录对象不能为空")
	}
	
	return r.db.Create(analytics).Error
}

// GetByID 根据ID获取分析记录
// 参数: id - 分析记录ID
// 返回: *models.Analytics - 分析记录对象, error - 错误信息
func (r *analyticsRepository) GetByID(id uint) (*models.Analytics, error) {
	if id == 0 {
		return nil, errors.New("分析记录ID不能为空")
	}
	
	analytics := &models.Analytics{}
	err := r.db.First(analytics, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分析记录不存在")
		}
		return nil, err
	}
	
	return analytics, nil
}

// Update 更新分析记录
// 参数: analytics - 分析记录对象
// 返回: error - 错误信息
func (r *analyticsRepository) Update(analytics *models.Analytics) error {
	if analytics == nil || analytics.ID == 0 {
		return errors.New("分析记录对象或ID不能为空")
	}
	
	return r.db.Save(analytics).Error
}

// Delete 删除分析记录
// 参数: id - 分析记录ID
// 返回: error - 错误信息
func (r *analyticsRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("分析记录ID不能为空")
	}
	
	return r.db.Delete(&models.Analytics{}, id).Error
}

// 仪表板统计实现

// GetDashboardStats 获取仪表板统计数据
// 返回: *models.DashboardStats - 仪表板统计数据, error - 错误信息
func (r *analyticsRepository) GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	
	// 总用户数
	r.db.Model(&models.User{}).Count(&stats.TotalUsers)
	
	// 总文章数
	r.db.Model(&models.Post{}).Count(&stats.TotalPosts)
	
	// 总评论数
	r.db.Model(&models.Comment{}).Count(&stats.TotalComments)
	
	// 总浏览量
	var totalViews sql.NullInt64
	r.db.Model(&models.Post{}).Select("SUM(view_count)").Scan(&totalViews)
	if totalViews.Valid {
		stats.TotalViews = totalViews.Int64
	}
	
	// 今日新用户
	today := time.Now().Truncate(24 * time.Hour)
	r.db.Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.TodayUsers)
	
	// 今日新文章
	r.db.Model(&models.Post{}).Where("created_at >= ?", today).Count(&stats.TodayPosts)
	
	// 今日新评论
	r.db.Model(&models.Comment{}).Where("created_at >= ?", today).Count(&stats.TodayComments)
	
	// 今日浏览量（简化实现，实际应该从访问日志获取）
	stats.TodayViews = 0
	
	// 活跃用户数（最近7天有活动的用户）
	weekAgo := time.Now().AddDate(0, 0, -7)
	r.db.Model(&models.User{}).Where("last_login_at >= ?", weekAgo).Count(&stats.ActiveUsers)
	
	// 已发布文章数
	r.db.Model(&models.Post{}).Where("status = ?", "published").Count(&stats.PublishedPosts)
	
	// 待审核评论数
	r.db.Model(&models.Comment{}).Where("status = ?", "pending").Count(&stats.PendingComments)
	
	// 计算增长率（与昨天比较）
	yesterday := today.AddDate(0, 0, -1)
	
	var yesterdayUsers, yesterdayPosts, yesterdayComments int64
	r.db.Model(&models.User{}).Where("created_at >= ? AND created_at < ?", yesterday, today).Count(&yesterdayUsers)
	r.db.Model(&models.Post{}).Where("created_at >= ? AND created_at < ?", yesterday, today).Count(&yesterdayPosts)
	r.db.Model(&models.Comment{}).Where("created_at >= ? AND created_at < ?", yesterday, today).Count(&yesterdayComments)
	
	if yesterdayUsers > 0 {
		stats.UserGrowthRate = float64(stats.TodayUsers-yesterdayUsers) / float64(yesterdayUsers) * 100
	}
	if yesterdayPosts > 0 {
		stats.PostGrowthRate = float64(stats.TodayPosts-yesterdayPosts) / float64(yesterdayPosts) * 100
	}
	if yesterdayComments > 0 {
		stats.CommentGrowthRate = float64(stats.TodayComments-yesterdayComments) / float64(yesterdayComments) * 100
	}
	
	return stats, nil
}

// GetDashboardStatsForPeriod 获取指定时期的仪表板统计
// 参数: startDate - 开始日期, endDate - 结束日期
// 返回: *models.DashboardStats - 仪表板统计数据, error - 错误信息
func (r *analyticsRepository) GetDashboardStatsForPeriod(startDate, endDate time.Time) (*models.DashboardStats, error) {
	if startDate.After(endDate) {
		return nil, errors.New("开始日期不能晚于结束日期")
	}
	
	stats := &models.DashboardStats{}
	
	// 指定时期内的统计
	r.db.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&stats.TodayUsers)
	r.db.Model(&models.Post{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&stats.TodayPosts)
	r.db.Model(&models.Comment{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&stats.TodayComments)
	
	// 总数统计
	r.db.Model(&models.User{}).Count(&stats.TotalUsers)
	r.db.Model(&models.Post{}).Count(&stats.TotalPosts)
	r.db.Model(&models.Comment{}).Count(&stats.TotalComments)
	
	// 总浏览量
	var totalViews sql.NullInt64
	r.db.Model(&models.Post{}).Select("SUM(view_count)").Scan(&totalViews)
	if totalViews.Valid {
		stats.TotalViews = totalViews.Int64
	}
	
	return stats, nil
}

// 内容统计实现

// GetContentStats 获取内容统计
// 返回: *models.ContentStats - 内容统计, error - 错误信息
func (r *analyticsRepository) GetContentStats() (*models.ContentStats, error) {
	stats := &models.ContentStats{}
	
	// 文章统计
	r.db.Model(&models.Post{}).Count(&stats.TotalPosts)
	r.db.Model(&models.Post{}).Where("status = ?", "published").Count(&stats.PublishedPosts)
	r.db.Model(&models.Post{}).Where("status = ?", "draft").Count(&stats.DraftPosts)
	
	// 评论统计
	r.db.Model(&models.Comment{}).Count(&stats.TotalComments)
	r.db.Model(&models.Comment{}).Where("status = ?", "approved").Count(&stats.ApprovedComments)
	r.db.Model(&models.Comment{}).Where("status = ?", "pending").Count(&stats.PendingComments)
	
	// 分类和标签统计
	r.db.Model(&models.Category{}).Count(&stats.TotalCategories)
	r.db.Model(&models.Tag{}).Count(&stats.TotalTags)
	
	// 平均文章长度
	var avgLength sql.NullFloat64
	r.db.Model(&models.Post{}).Select("AVG(LENGTH(content))").Scan(&avgLength)
	if avgLength.Valid {
		stats.AveragePostLength = avgLength.Float64
	}
	
	// 平均评论长度
	r.db.Model(&models.Comment{}).Select("AVG(LENGTH(content))").Scan(&avgLength)
	if avgLength.Valid {
		stats.AverageCommentLength = avgLength.Float64
	}
	
	// 最受欢迎的分类
	var popularCategory struct {
		Name  string
		Count int64
	}
	r.db.Table("categories").
		Select("categories.name, COUNT(posts.id) as count").
		Joins("LEFT JOIN posts ON categories.id = posts.category_id").
		Group("categories.id").Order("count DESC").Limit(1).Scan(&popularCategory)
	stats.MostPopularCategory = popularCategory.Name
	
	return stats, nil
}

// GetPopularPosts 获取热门文章
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.PopularPost - 热门文章列表, error - 错误信息
func (r *analyticsRepository) GetPopularPosts(limit int, days int) ([]models.PopularPost, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 7
	}
	
	var popularPosts []models.PopularPost
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("posts").
		Select(`
			posts.id,
			posts.title,
			posts.slug,
			posts.view_count,
			COALESCE(comment_counts.comment_count, 0) as comment_count,
			COALESCE(like_counts.like_count, 0) as like_count,
			posts.published_at,
			(posts.view_count * 1.0 + COALESCE(comment_counts.comment_count, 0) * 2.0 + COALESCE(like_counts.like_count, 0) * 3.0) as popularity_score
		`).
		Joins(`LEFT JOIN (
			SELECT post_id, COUNT(*) as comment_count 
			FROM comments 
			WHERE created_at >= ? 
			GROUP BY post_id
		) comment_counts ON posts.id = comment_counts.post_id`, startDate).
		Joins(`LEFT JOIN (
			SELECT target_id, COUNT(*) as like_count 
			FROM likes 
			WHERE target_type = 'post' AND created_at >= ?
			GROUP BY target_id
		) like_counts ON posts.id = like_counts.target_id`, startDate).
		Where("posts.status = ? AND posts.published_at >= ?", "published", startDate).
		Order("popularity_score DESC").Limit(limit).Scan(&popularPosts).Error
	
	return popularPosts, err
}

// GetCategoryStats 获取分类统计
// 返回: []models.CategoryStats - 分类统计列表, error - 错误信息
func (r *analyticsRepository) GetCategoryStats() ([]models.CategoryStats, error) {
	var categoryStats []models.CategoryStats
	
	err := r.db.Table("categories").
		Select(`
			categories.id,
			categories.name,
			categories.slug,
			COUNT(posts.id) as post_count,
			COALESCE(SUM(posts.view_count), 0) as total_views,
			COALESCE(AVG(posts.view_count), 0) as average_views
		`).
		Joins("LEFT JOIN posts ON categories.id = posts.category_id AND posts.status = 'published'").
		Group("categories.id").Order("post_count DESC").Scan(&categoryStats).Error
	
	return categoryStats, err
}

// GetTagStats 获取标签统计
// 参数: limit - 限制数量
// 返回: []models.TagStats - 标签统计列表, error - 错误信息
func (r *analyticsRepository) GetTagStats(limit int) ([]models.TagStats, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var tagStats []models.TagStats
	
	err := r.db.Table("tags").
		Select(`
			tags.id,
			tags.name,
			tags.slug,
			COUNT(post_tags.post_id) as post_count,
			COALESCE(SUM(posts.view_count), 0) as total_views
		`).
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN posts ON post_tags.post_id = posts.id AND posts.status = 'published'").
		Group("tags.id").Order("post_count DESC").Limit(limit).Scan(&tagStats).Error
	
	return tagStats, err
}

// 用户统计实现

// GetActiveUsers 获取活跃用户
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.ActiveUser - 活跃用户列表, error - 错误信息
func (r *analyticsRepository) GetActiveUsers(limit int, days int) ([]models.ActiveUser, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 7
	}
	
	var activeUsers []models.ActiveUser
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("users").
		Select(`
			users.id,
			users.username,
			users.email,
			COALESCE(post_counts.post_count, 0) as post_count,
			COALESCE(comment_counts.comment_count, 0) as comment_count,
			users.last_login_at,
			(COALESCE(post_counts.post_count, 0) * 2.0 + COALESCE(comment_counts.comment_count, 0) * 1.0) as activity_score
		`).
		Joins(`LEFT JOIN (
			SELECT author_id, COUNT(*) as post_count 
			FROM posts 
			WHERE created_at >= ? 
			GROUP BY author_id
		) post_counts ON users.id = post_counts.author_id`, startDate).
		Joins(`LEFT JOIN (
			SELECT user_id, COUNT(*) as comment_count 
			FROM comments 
			WHERE created_at >= ? 
			GROUP BY user_id
		) comment_counts ON users.id = comment_counts.user_id`, startDate).
		Where("users.last_login_at >= ?", startDate).
		Order("activity_score DESC").Limit(limit).Scan(&activeUsers).Error
	
	return activeUsers, err
}

// GetUserGrowthStats 获取用户增长统计
// 参数: days - 统计天数
// 返回: []models.UserGrowthStats - 用户增长统计列表, error - 错误信息
func (r *analyticsRepository) GetUserGrowthStats(days int) ([]models.UserGrowthStats, error) {
	if days <= 0 {
		days = 30
	}
	
	var growthStats []models.UserGrowthStats
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("users").
		Select(`
			DATE(created_at) as date,
			COUNT(*) as new_users,
			SUM(COUNT(*)) OVER (ORDER BY DATE(created_at)) as total_users
		`).
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").Order("date").Scan(&growthStats).Error
	
	return growthStats, err
}

// GetUserEngagementStats 获取用户参与度统计
// 参数: userID - 用户ID
// 返回: *UserEngagement - 用户参与度统计, error - 错误信息
func (r *analyticsRepository) GetUserEngagementStats(userID uint) (*UserEngagement, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	engagement := &UserEngagement{UserID: userID}
	
	// 获取用户基本信息
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	engagement.RegistrationDate = user.CreatedAt
	if user.LastLoginAt != nil {
		engagement.LastActiveDate = *user.LastLoginAt
	}
	
	// 统计文章数
	r.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&engagement.TotalPosts)
	
	// 统计评论数
	r.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&engagement.TotalComments)
	
	// 统计获得的点赞数
	r.db.Table("likes").
		Joins("JOIN posts ON likes.target_id = posts.id AND likes.target_type = 'post'").
		Where("posts.author_id = ?", userID).Count(&engagement.TotalLikes)
	
	// 统计文章总浏览量
	var totalViews sql.NullInt64
	r.db.Model(&models.Post{}).Where("author_id = ?", userID).Select("SUM(view_count)").Scan(&totalViews)
	if totalViews.Valid {
		engagement.TotalViews = totalViews.Int64
	}
	
	// 计算参与度评分
	engagement.EngagementScore = float64(engagement.TotalPosts)*2.0 + 
		float64(engagement.TotalComments)*1.0 + 
		float64(engagement.TotalLikes)*0.5 + 
		float64(engagement.TotalViews)*0.1
	
	return engagement, nil
}

// 趋势分析实现

// GetPostTrends 获取文章发布趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (r *analyticsRepository) GetPostTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("posts").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").Order("date").Scan(&trends).Error
	
	return trends, err
}

// GetCommentTrends 获取评论趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (r *analyticsRepository) GetCommentTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("comments").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").Order("date").Scan(&trends).Error
	
	return trends, err
}

// GetViewTrends 获取浏览量趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (r *analyticsRepository) GetViewTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	// 简化实现：基于文章发布日期统计浏览量
	// 实际应用中应该有专门的访问日志表
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("posts").
		Select("DATE(published_at) as date, SUM(view_count) as value").
		Where("published_at >= ? AND status = ?", startDate, "published").
		Group("DATE(published_at)").Order("date").Scan(&trends).Error
	
	return trends, err
}

// GetUserRegistrationTrends 获取用户注册趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (r *analyticsRepository) GetUserRegistrationTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").Order("date").Scan(&trends).Error
	
	return trends, err
}

// 性能分析实现

// GetTopPerformingContent 获取表现最佳的内容
// 参数: contentType - 内容类型, metric - 指标, limit - 限制数量, days - 统计天数
// 返回: []PerformanceData - 性能数据列表, error - 错误信息
func (r *analyticsRepository) GetTopPerformingContent(contentType string, metric string, limit int, days int) ([]PerformanceData, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 30
	}
	
	var performanceData []PerformanceData
	startDate := time.Now().AddDate(0, 0, -days)
	
	switch contentType {
	case "post":
		var orderBy string
		switch metric {
		case "views":
			orderBy = "view_count DESC"
		case "comments":
			orderBy = "(SELECT COUNT(*) FROM comments WHERE post_id = posts.id) DESC"
		case "likes":
			orderBy = "(SELECT COUNT(*) FROM likes WHERE target_id = posts.id AND target_type = 'post') DESC"
		default:
			orderBy = "view_count DESC"
		}
		
		err := r.db.Table("posts").
			Select("id, title, 'post' as type, view_count as metric_value, created_at").
			Where("created_at >= ? AND status = ?", startDate, "published").
			Order(orderBy).Limit(limit).Scan(&performanceData).Error
		return performanceData, err
		
	case "comment":
		err := r.db.Table("comments").
			Select(`
				id, 
				SUBSTR(content, 1, 50) as title, 
				'comment' as type, 
				(SELECT COUNT(*) FROM likes WHERE target_id = comments.id AND target_type = 'comment') as metric_value,
				created_at
			`).
			Where("created_at >= ? AND status = ?", startDate, "approved").
			Order("metric_value DESC").Limit(limit).Scan(&performanceData).Error
		return performanceData, err
		
	default:
		return nil, errors.New("不支持的内容类型")
	}
}

// GetEngagementMetrics 获取参与度指标
// 参数: days - 统计天数
// 返回: *EngagementMetrics - 参与度指标, error - 错误信息
func (r *analyticsRepository) GetEngagementMetrics(days int) (*EngagementMetrics, error) {
	if days <= 0 {
		days = 30
	}
	
	metrics := &EngagementMetrics{}
	startDate := time.Now().AddDate(0, 0, -days)
	
	// 总浏览量
	var totalViews sql.NullInt64
	r.db.Model(&models.Post{}).Where("created_at >= ?", startDate).Select("SUM(view_count)").Scan(&totalViews)
	if totalViews.Valid {
		metrics.TotalViews = totalViews.Int64
	}
	
	// 总评论数
	r.db.Model(&models.Comment{}).Where("created_at >= ?", startDate).Count(&metrics.TotalComments)
	
	// 总点赞数
	r.db.Model(&models.Like{}).Where("created_at >= ?", startDate).Count(&metrics.TotalLikes)
	
	// 计算参与度（评论数+点赞数）/浏览量
	if metrics.TotalViews > 0 {
		metrics.EngagementRate = float64(metrics.TotalComments+metrics.TotalLikes) / float64(metrics.TotalViews) * 100
	}
	
	// 简化的跳出率计算（实际需要访问日志数据）
	metrics.BounceRate = 65.0 // 假设值
	metrics.AverageTimeOnPage = 180.0 // 假设值（秒）
	
	return metrics, nil
}

// GetRetentionAnalysis 获取留存分析
// 参数: days - 统计天数
// 返回: *RetentionAnalysis - 留存分析, error - 错误信息
func (r *analyticsRepository) GetRetentionAnalysis(days int) (*RetentionAnalysis, error) {
	if days <= 0 {
		days = 30
	}
	
	analysis := &RetentionAnalysis{}
	startDate := time.Now().AddDate(0, 0, -days)
	
	// 简化的留存率计算
	// 实际应用中需要更复杂的用户行为追踪
	
	// 获取注册用户数
	var totalUsers int64
	r.db.Model(&models.User{}).Where("created_at >= ?", startDate).Count(&totalUsers)
	
	if totalUsers == 0 {
		return analysis, nil
	}
	
	// 1天留存（简化：第二天有登录的用户）
	var day1Retained int64
	r.db.Model(&models.User{}).
		Where("created_at >= ? AND last_login_at > DATE_ADD(created_at, INTERVAL 1 DAY)", startDate).
		Count(&day1Retained)
	analysis.Day1Retention = float64(day1Retained) / float64(totalUsers) * 100
	
	// 7天留存
	var day7Retained int64
	r.db.Model(&models.User{}).
		Where("created_at >= ? AND last_login_at > DATE_ADD(created_at, INTERVAL 7 DAY)", startDate).
		Count(&day7Retained)
	analysis.Day7Retention = float64(day7Retained) / float64(totalUsers) * 100
	
	// 30天留存
	var day30Retained int64
	r.db.Model(&models.User{}).
		Where("created_at >= ? AND last_login_at > DATE_ADD(created_at, INTERVAL 30 DAY)", startDate).
		Count(&day30Retained)
	analysis.Day30Retention = float64(day30Retained) / float64(totalUsers) * 100
	
	return analysis, nil
}

// 实时统计实现

// GetRealTimeStats 获取实时统计
// 返回: *RealTimeStats - 实时统计, error - 错误信息
func (r *analyticsRepository) GetRealTimeStats() (*RealTimeStats, error) {
	stats := &RealTimeStats{}
	today := time.Now().Truncate(24 * time.Hour)
	
	// 今日统计
	r.db.Model(&models.Post{}).Where("created_at >= ?", today).Count(&stats.TodayPosts)
	r.db.Model(&models.Comment{}).Where("created_at >= ?", today).Count(&stats.TodayComments)
	r.db.Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.TodayUsers)
	
	// 活跃用户（最近1小时有活动）
	hourAgo := time.Now().Add(-time.Hour)
	r.db.Model(&models.User{}).Where("last_login_at >= ?", hourAgo).Count(&stats.ActiveUsers)
	
	// 在线用户（简化实现，实际需要实时会话管理）
	stats.OnlineUsers = stats.ActiveUsers / 2 // 假设值
	
	// 当前浏览量（简化实现）
	stats.CurrentViews = stats.OnlineUsers * 3 // 假设值
	
	stats.LastUpdated = time.Now()
	
	return stats, nil
}

// UpdateRealTimeStats 更新实时统计
// 参数: stats - 实时统计数据
// 返回: error - 错误信息
func (r *analyticsRepository) UpdateRealTimeStats(stats *RealTimeStats) error {
	if stats == nil {
		return errors.New("实时统计数据不能为空")
	}
	
	// 这里可以将实时统计数据存储到缓存或专门的表中
	// 简化实现：直接返回成功
	return nil
}

// 自定义查询实现

// ExecuteCustomQuery 执行自定义查询
// 参数: query - SQL查询语句, params - 查询参数
// 返回: []map[string]interface{} - 查询结果, error - 错误信息
func (r *analyticsRepository) ExecuteCustomQuery(query string, params ...interface{}) ([]map[string]interface{}, error) {
	if query == "" {
		return nil, errors.New("查询语句不能为空")
	}
	
	var results []map[string]interface{}
	err := r.db.Raw(query, params...).Scan(&results).Error
	return results, err
}

// GetMetricsByDateRange 获取指定日期范围的指标数据
// 参数: metric - 指标名称, startDate - 开始日期, endDate - 结束日期
// 返回: []MetricData - 指标数据列表, error - 错误信息
func (r *analyticsRepository) GetMetricsByDateRange(metric string, startDate, endDate time.Time) ([]MetricData, error) {
	if metric == "" {
		return nil, errors.New("指标名称不能为空")
	}
	if startDate.After(endDate) {
		return nil, errors.New("开始日期不能晚于结束日期")
	}
	
	var metricData []MetricData
	
	switch metric {
	case "posts":
		err := r.db.Table("posts").
			Select("DATE(created_at) as date, 'posts' as metric, COUNT(*) as value").
			Where("created_at BETWEEN ? AND ?", startDate, endDate).
			Group("DATE(created_at)").Order("date").Scan(&metricData).Error
		return metricData, err
		
	case "comments":
		err := r.db.Table("comments").
			Select("DATE(created_at) as date, 'comments' as metric, COUNT(*) as value").
			Where("created_at BETWEEN ? AND ?", startDate, endDate).
			Group("DATE(created_at)").Order("date").Scan(&metricData).Error
		return metricData, err
		
	case "users":
		err := r.db.Table("users").
			Select("DATE(created_at) as date, 'users' as metric, COUNT(*) as value").
			Where("created_at BETWEEN ? AND ?", startDate, endDate).
			Group("DATE(created_at)").Order("date").Scan(&metricData).Error
		return metricData, err
		
	case "views":
		err := r.db.Table("posts").
			Select("DATE(published_at) as date, 'views' as metric, SUM(view_count) as value").
			Where("published_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "published").
			Group("DATE(published_at)").Order("date").Scan(&metricData).Error
		return metricData, err
		
	default:
		return nil, errors.New("不支持的指标类型")
	}
}
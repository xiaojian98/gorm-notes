package services

import (
	"errors"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// AnalyticsService 分析服务接口
// 定义数据统计和分析相关的业务操作
type AnalyticsService interface {
	// 仪表板统计
	GetDashboardStats() (*models.DashboardStats, error)           // 获取仪表板统计数据
	GetDashboardStatsForPeriod(days int) (*models.DashboardStats, error) // 获取指定时间段的仪表板统计
	
	// 内容统计
	GetContentStats() (*models.ContentStats, error)              // 获取内容统计
	GetPopularPosts(limit int, days int) ([]models.PopularPost, error) // 获取热门文章
	GetCategoryStats() ([]models.CategoryStats, error)           // 获取分类统计
	GetTagStats(limit int) ([]models.TagStats, error)           // 获取标签统计
	
	// 用户统计
	GetActiveUsers(limit int, days int) ([]models.ActiveUser, error) // 获取活跃用户
	GetUserGrowthStats(days int) ([]models.UserGrowthStats, error) // 获取用户增长统计
	GetUserEngagementStats(userID uint) (*UserEngagementStats, error) // 获取用户参与度统计
	
	// 趋势分析
	GetPostTrends(days int) ([]TrendData, error)                 // 获取文章发布趋势
	GetCommentTrends(days int) ([]TrendData, error)              // 获取评论趋势
	GetUserRegistrationTrends(days int) ([]TrendData, error)     // 获取用户注册趋势
	GetViewTrends(days int) ([]TrendData, error)                 // 获取浏览量趋势
	
	// 性能分析
	GetTopPerformingContent(limit int, metric string) ([]ContentPerformance, error) // 获取表现最佳的内容
	GetEngagementMetrics(startDate, endDate time.Time) (*EngagementMetrics, error) // 获取参与度指标
	
	// 实时统计
	GetRealTimeStats() (*RealTimeStats, error)                   // 获取实时统计
	UpdateRealTimeStats(event string, data map[string]interface{}) error // 更新实时统计
}

// analyticsService 分析服务实现
type analyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService 创建分析服务实例
// 参数: db - 数据库连接
// 返回: AnalyticsService - 分析服务接口实例
func NewAnalyticsService(db *gorm.DB) AnalyticsService {
	return &analyticsService{
		db: db,
	}
}

// 辅助数据结构

// UserEngagementStats 用户参与度统计
type UserEngagementStats struct {
	UserID           uint    `json:"user_id"`
	TotalPosts       int     `json:"total_posts"`       // 总文章数
	TotalComments    int     `json:"total_comments"`    // 总评论数
	TotalLikes       int     `json:"total_likes"`       // 总点赞数
	TotalViews       int     `json:"total_views"`       // 总浏览数
	EngagementRate   float64 `json:"engagement_rate"`   // 参与度
	LastActiveAt     *time.Time `json:"last_active_at,omitempty"` // 最后活跃时间
	ActiveDays       int     `json:"active_days"`       // 活跃天数
	AvgPostsPerDay   float64 `json:"avg_posts_per_day"` // 平均每日文章数
}

// TrendData 趋势数据
type TrendData struct {
	Date  time.Time `json:"date"`  // 日期
	Count int       `json:"count"` // 数量
	Value float64   `json:"value"` // 数值（可选）
}

// ContentPerformance 内容表现数据
type ContentPerformance struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	Type         string  `json:"type"` // post, comment等
	Score        float64 `json:"score"` // 综合评分
	Views        int     `json:"views"`
	Likes        int     `json:"likes"`
	Comments     int     `json:"comments"`
	Shares       int     `json:"shares"`
	Engagement   float64 `json:"engagement"`
	CreatedAt    time.Time `json:"created_at"`
}

// EngagementMetrics 参与度指标
type EngagementMetrics struct {
	TotalViews       int     `json:"total_views"`
	TotalLikes       int     `json:"total_likes"`
	TotalComments    int     `json:"total_comments"`
	TotalShares      int     `json:"total_shares"`
	EngagementRate   float64 `json:"engagement_rate"`
	AvgTimeOnSite    float64 `json:"avg_time_on_site"`
	BounceRate       float64 `json:"bounce_rate"`
	ReturnVisitorRate float64 `json:"return_visitor_rate"`
}

// RealTimeStats 实时统计
type RealTimeStats struct {
	OnlineUsers      int       `json:"online_users"`
	ActiveUsers      int       `json:"active_users"`
	TodayViews       int       `json:"today_views"`
	TodayPosts       int       `json:"today_posts"`
	TodayComments    int       `json:"today_comments"`
	TodayRegistrations int     `json:"today_registrations"`
	LastUpdated      time.Time `json:"last_updated"`
}

// 仪表板统计实现

// GetDashboardStats 获取仪表板统计数据
// 返回: *models.DashboardStats - 仪表板统计数据, error - 错误信息
func (s *analyticsService) GetDashboardStats() (*models.DashboardStats, error) {
	return s.GetDashboardStatsForPeriod(30) // 默认30天
}

// GetDashboardStatsForPeriod 获取指定时间段的仪表板统计
// 参数: days - 统计天数
// 返回: *models.DashboardStats - 仪表板统计数据, error - 错误信息
func (s *analyticsService) GetDashboardStatsForPeriod(days int) (*models.DashboardStats, error) {
	if days <= 0 {
		days = 30
	}
	
	stats := &models.DashboardStats{}
	startDate := time.Now().AddDate(0, 0, -days)
	
	// 总用户数
	var totalUsers int64
	s.db.Model(&models.User{}).Count(&totalUsers)
	stats.TotalUsers = totalUsers
	
	// 新用户数（指定时间段内）
	var newUsers int64
	s.db.Model(&models.User{}).Where("created_at >= ?", startDate).Count(&newUsers)
	stats.TodayUsers = newUsers
	
	// 总文章数
	var totalPosts int64
	s.db.Model(&models.Post{}).Count(&totalPosts)
	stats.TotalPosts = totalPosts
	
	// 新文章数（指定时间段内）
	var newPosts int64
	s.db.Model(&models.Post{}).Where("created_at >= ?", startDate).Count(&newPosts)
	stats.TodayPosts = newPosts
	
	// 总评论数
	var totalComments int64
	s.db.Model(&models.Comment{}).Count(&totalComments)
	stats.TotalComments = totalComments
	
	// 新评论数（指定时间段内）
	var newComments int64
	s.db.Model(&models.Comment{}).Where("created_at >= ?", startDate).Count(&newComments)
	stats.TodayComments = newComments
	
	// 总浏览量
	var totalViews int64
	s.db.Model(&models.Post{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
	stats.TotalViews = totalViews
	
	// 计算增长率
	prevStartDate := startDate.AddDate(0, 0, -days)
	stats.UserGrowthRate = s.calculateGrowthRate("users", prevStartDate, startDate, startDate, time.Now())
	stats.PostGrowthRate = s.calculateGrowthRate("posts", prevStartDate, startDate, startDate, time.Now())
	stats.CommentGrowthRate = s.calculateGrowthRate("comments", prevStartDate, startDate, startDate, time.Now())
	
	return stats, nil
}

// 内容统计实现

// GetContentStats 获取内容统计
// 返回: *models.ContentStats - 内容统计数据, error - 错误信息
func (s *analyticsService) GetContentStats() (*models.ContentStats, error) {
	stats := &models.ContentStats{}
	
	// 已发布文章数
	var publishedPosts int64
	s.db.Model(&models.Post{}).Where("status = ?", models.PostStatusPublished).Count(&publishedPosts)
	stats.PublishedPosts = publishedPosts
	
	// 草稿文章数
	var draftPosts int64
	s.db.Model(&models.Post{}).Where("status = ?", models.PostStatusDraft).Count(&draftPosts)
	stats.DraftPosts = draftPosts
	
	// 总分类数
	var totalCategories int64
	s.db.Model(&models.Category{}).Count(&totalCategories)
	stats.TotalCategories = totalCategories
	
	// 总标签数
	var totalTags int64
	s.db.Model(&models.Tag{}).Count(&totalTags)
	stats.TotalTags = totalTags
	
	// 平均文章长度
	var avgWordCount float64
	s.db.Model(&models.Post{}).Where("status = ?", "published").Select("AVG(word_count)").Scan(&avgWordCount)
	stats.AvgPostLength = avgWordCount
	
	// 平均阅读时间
	var avgReadTime float64
	s.db.Model(&models.Post{}).Where("status = ?", "published").Select("AVG(read_time)").Scan(&avgReadTime)
	stats.AvgReadTime = avgReadTime
	
	return stats, nil
}

// GetPopularPosts 获取热门文章
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.PopularPost - 热门文章列表, error - 错误信息
func (s *analyticsService) GetPopularPosts(limit int, days int) ([]models.PopularPost, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 7
	}
	
	var posts []models.PopularPost
	startDate := time.Now().AddDate(0, 0, -days)
	
	// 查询热门文章（基于浏览量、点赞数、评论数综合排序）
	err := s.db.Table("posts").
		Select(`
			posts.id,
			posts.title,
			posts.view_count,
			COALESCE(like_counts.like_count, 0) as like_count,
			COALESCE(comment_counts.comment_count, 0) as comment_count,
			(posts.view_count * 1 + COALESCE(like_counts.like_count, 0) * 5 + COALESCE(comment_counts.comment_count, 0) * 10) as popularity_score
		`).
		Joins(`LEFT JOIN (
			SELECT target_id, COUNT(*) as like_count 
			FROM likes 
			WHERE target_type = 'post' AND created_at >= ?
			GROUP BY target_id
		) like_counts ON posts.id = like_counts.target_id`, startDate).
		Joins(`LEFT JOIN (
			SELECT post_id, COUNT(*) as comment_count 
			FROM comments 
			WHERE created_at >= ?
			GROUP BY post_id
		) comment_counts ON posts.id = comment_counts.post_id`, startDate).
		Where("posts.status = ? AND posts.published_at >= ?", "published", startDate).
		Order("popularity_score DESC").
		Limit(limit).
		Scan(&posts).Error
	
	if err != nil {
		return nil, err
	}
	
	return posts, nil
}

// GetCategoryStats 获取分类统计
// 返回: []models.CategoryStats - 分类统计列表, error - 错误信息
func (s *analyticsService) GetCategoryStats() ([]models.CategoryStats, error) {
	var stats []models.CategoryStats
	
	err := s.db.Table("categories").
		Select(`
			categories.id,
			categories.name,
			COUNT(posts.id) as post_count,
			COALESCE(SUM(posts.view_count), 0) as total_views
		`).
		Joins("LEFT JOIN posts ON categories.id = posts.category_id AND posts.status = 'published'").
		Group("categories.id, categories.name").
		Order("post_count DESC").
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}

// GetTagStats 获取标签统计
// 参数: limit - 限制数量
// 返回: []models.TagStats - 标签统计列表, error - 错误信息
func (s *analyticsService) GetTagStats(limit int) ([]models.TagStats, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var stats []models.TagStats
	
	err := s.db.Table("tags").
		Select(`
			tags.id,
			tags.name,
			COUNT(post_tags.post_id) as usage_count
		`).
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN posts ON post_tags.post_id = posts.id AND posts.status = 'published'").
		Group("tags.id, tags.name").
		Order("usage_count DESC").
		Limit(limit).
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}

// 用户统计实现

// GetActiveUsers 获取活跃用户
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.ActiveUser - 活跃用户列表, error - 错误信息
func (s *analyticsService) GetActiveUsers(limit int, days int) ([]models.ActiveUser, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 7
	}
	
	var users []models.ActiveUser
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := s.db.Table("users").
		Select(`
			users.id,
			users.username,
			COALESCE(post_counts.post_count, 0) as post_count,
			COALESCE(comment_counts.comment_count, 0) as comment_count,
			(COALESCE(post_counts.post_count, 0) * 5 + COALESCE(comment_counts.comment_count, 0) * 2) as activity_score
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
		Where("users.status = ?", "active").
		Having("activity_score > 0").
		Order("activity_score DESC").
		Limit(limit).
		Scan(&users).Error
	
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

// GetUserGrowthStats 获取用户增长统计
// 参数: days - 统计天数
// 返回: []models.UserGrowthStats - 用户增长统计列表, error - 错误信息
func (s *analyticsService) GetUserGrowthStats(days int) ([]models.UserGrowthStats, error) {
	if days <= 0 {
		days = 30
	}
	
	var stats []models.UserGrowthStats
	startDate := time.Now().AddDate(0, 0, -days)
	
	// 按日统计用户注册数
	err := s.db.Table("users").
		Select(`
			DATE(created_at) as date,
			COUNT(*) as new_users,
			0 as total_users
		`).
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	// 计算累计用户数
	var totalUsers int
	for i := range stats {
		totalUsers += stats[i].NewUsers
		stats[i].TotalUsers = totalUsers
	}
	
	return stats, nil
}

// GetUserEngagementStats 获取用户参与度统计
// 参数: userID - 用户ID
// 返回: *UserEngagementStats - 用户参与度统计, error - 错误信息
func (s *analyticsService) GetUserEngagementStats(userID uint) (*UserEngagementStats, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	stats := &UserEngagementStats{
		UserID: userID,
	}
	
	// 获取用户基本信息
	user := &models.User{}
	if err := s.db.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	// 总文章数
	var totalPosts int64
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&totalPosts)
	stats.TotalPosts = int(totalPosts)
	
	// 总评论数
	var totalComments int64
	s.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&totalComments)
	stats.TotalComments = int(totalComments)
	
	// 获得的总点赞数
	var totalLikes int64
	s.db.Table("likes").
		Joins("JOIN posts ON likes.target_id = posts.id AND likes.target_type = 'post'").
		Where("posts.author_id = ?", userID).Count(&totalLikes)
	stats.TotalLikes = int(totalLikes)
	
	// 文章总浏览数
	var totalViews int64
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
	stats.TotalViews = int(totalViews)
	
	// 计算参与度（基于发布内容和互动）
	if stats.TotalPosts > 0 {
		stats.EngagementRate = float64(stats.TotalLikes+stats.TotalComments) / float64(stats.TotalPosts)
	}
	
	// 计算活跃天数和平均每日文章数
	daysSinceRegistration := int(time.Since(user.CreatedAt).Hours() / 24)
	if daysSinceRegistration > 0 {
		stats.ActiveDays = daysSinceRegistration
		stats.AvgPostsPerDay = float64(stats.TotalPosts) / float64(daysSinceRegistration)
	}
	
	// 最后活跃时间（最后发布文章或评论的时间）
	var lastPostTime, lastCommentTime time.Time
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Select("MAX(created_at)").Scan(&lastPostTime)
	s.db.Model(&models.Comment{}).Where("user_id = ?", userID).Select("MAX(created_at)").Scan(&lastCommentTime)
	
	if lastPostTime.After(lastCommentTime) {
		stats.LastActiveAt = &lastPostTime
	} else if !lastCommentTime.IsZero() {
		stats.LastActiveAt = &lastCommentTime
	}
	
	return stats, nil
}

// 趋势分析实现

// GetPostTrends 获取文章发布趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (s *analyticsService) GetPostTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := s.db.Table("posts").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trends).Error
	
	return trends, err
}

// GetCommentTrends 获取评论趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (s *analyticsService) GetCommentTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := s.db.Table("comments").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trends).Error
	
	return trends, err
}

// GetUserRegistrationTrends 获取用户注册趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (s *analyticsService) GetUserRegistrationTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := s.db.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trends).Error
	
	return trends, err
}

// GetViewTrends 获取浏览量趋势
// 参数: days - 统计天数
// 返回: []TrendData - 趋势数据列表, error - 错误信息
func (s *analyticsService) GetViewTrends(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 30
	}
	
	// 注意：这里简化处理，实际应该有专门的浏览记录表
	var trends []TrendData
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := s.db.Table("posts").
		Select("DATE(updated_at) as date, SUM(view_count) as count").
		Where("updated_at >= ?", startDate).
		Group("DATE(updated_at)").
		Order("date ASC").
		Scan(&trends).Error
	
	return trends, err
}

// 性能分析实现

// GetTopPerformingContent 获取表现最佳的内容
// 参数: limit - 限制数量, metric - 评估指标（views, likes, comments, engagement）
// 返回: []ContentPerformance - 内容表现列表, error - 错误信息
func (s *analyticsService) GetTopPerformingContent(limit int, metric string) ([]ContentPerformance, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	var content []ContentPerformance
	var orderBy string
	
	switch metric {
	case "views":
		orderBy = "posts.view_count DESC"
	case "likes":
		orderBy = "like_count DESC"
	case "comments":
		orderBy = "comment_count DESC"
	case "engagement":
		orderBy = "(like_count + comment_count * 2) DESC"
	default:
		orderBy = "(posts.view_count + like_count * 5 + comment_count * 10) DESC"
	}
	
	err := s.db.Table("posts").
		Select(`
			posts.id,
			posts.title,
			'post' as type,
			posts.view_count as views,
			COALESCE(like_counts.like_count, 0) as likes,
			COALESCE(comment_counts.comment_count, 0) as comments,
			0 as shares,
			posts.created_at,
			(posts.view_count + COALESCE(like_counts.like_count, 0) * 5 + COALESCE(comment_counts.comment_count, 0) * 10) as score
		`).
		Joins(`LEFT JOIN (
			SELECT target_id, COUNT(*) as like_count 
			FROM likes 
			WHERE target_type = 'post'
			GROUP BY target_id
		) like_counts ON posts.id = like_counts.target_id`).
		Joins(`LEFT JOIN (
			SELECT post_id, COUNT(*) as comment_count 
			FROM comments 
			GROUP BY post_id
		) comment_counts ON posts.id = comment_counts.post_id`).
		Where("posts.status = ?", "published").
		Order(orderBy).
		Limit(limit).
		Scan(&content).Error
	
	if err != nil {
		return nil, err
	}
	
	// 计算参与度
	for i := range content {
		if content[i].Views > 0 {
			content[i].Engagement = float64(content[i].Likes+content[i].Comments) / float64(content[i].Views) * 100
		}
	}
	
	return content, nil
}

// GetEngagementMetrics 获取参与度指标
// 参数: startDate - 开始日期, endDate - 结束日期
// 返回: *EngagementMetrics - 参与度指标, error - 错误信息
func (s *analyticsService) GetEngagementMetrics(startDate, endDate time.Time) (*EngagementMetrics, error) {
	if startDate.After(endDate) {
		return nil, errors.New("开始日期不能晚于结束日期")
	}
	
	metrics := &EngagementMetrics{}
	
	// 总浏览量
	var totalViews int64
	s.db.Model(&models.Post{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
	metrics.TotalViews = int(totalViews)
	
	// 总点赞数
	var totalLikes int64
	s.db.Model(&models.Like{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&totalLikes)
	metrics.TotalLikes = int(totalLikes)
	
	// 总评论数
	var totalComments int64
	s.db.Model(&models.Comment{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&totalComments)
	metrics.TotalComments = int(totalComments)
	
	// 计算参与度（简化计算）
	if metrics.TotalViews > 0 {
		metrics.EngagementRate = float64(metrics.TotalLikes+metrics.TotalComments) / float64(metrics.TotalViews) * 100
	}
	
	// TODO: 实现更复杂的指标计算
	// 平均停留时间、跳出率、回访率等需要额外的数据收集
	metrics.AvgTimeOnSite = 0
	metrics.BounceRate = 0
	metrics.ReturnVisitorRate = 0
	
	return metrics, nil
}

// 实时统计实现

// GetRealTimeStats 获取实时统计
// 返回: *RealTimeStats - 实时统计数据, error - 错误信息
func (s *analyticsService) GetRealTimeStats() (*RealTimeStats, error) {
	stats := &RealTimeStats{
		LastUpdated: time.Now(),
	}
	
	today := time.Now().Truncate(24 * time.Hour)
	
	// 今日浏览量（简化处理）
	var todayViews int64
	s.db.Model(&models.Post{}).Where("updated_at >= ?", today).
		Select("COALESCE(SUM(view_count), 0)").Scan(&todayViews)
	stats.TodayViews = int(todayViews)
	
	// 今日文章数
	var todayPosts int64
	s.db.Model(&models.Post{}).Where("created_at >= ?", today).Count(&todayPosts)
	stats.TodayPosts = int(todayPosts)
	
	// 今日评论数
	var todayComments int64
	s.db.Model(&models.Comment{}).Where("created_at >= ?", today).Count(&todayComments)
	stats.TodayComments = int(todayComments)
	
	// 今日注册数
	var todayRegistrations int64
	s.db.Model(&models.User{}).Where("created_at >= ?", today).Count(&todayRegistrations)
	stats.TodayRegistrations = int(todayRegistrations)
	
	// TODO: 实现在线用户和活跃用户统计
	// 需要额外的会话管理和用户活动跟踪
	stats.OnlineUsers = 0
	stats.ActiveUsers = 0
	
	return stats, nil
}

// UpdateRealTimeStats 更新实时统计
// 参数: event - 事件类型, data - 事件数据
// 返回: error - 错误信息
func (s *analyticsService) UpdateRealTimeStats(event string, data map[string]interface{}) error {
	// 这里可以实现实时统计的更新逻辑
	// 例如：用户登录、文章浏览、评论发布等事件的处理
	
	switch event {
	case "user_login":
		// 处理用户登录事件
	case "post_view":
		// 处理文章浏览事件
	case "comment_created":
		// 处理评论创建事件
	case "user_register":
		// 处理用户注册事件
	default:
		// 未知事件类型
	}
	
	// TODO: 实现具体的统计更新逻辑
	return nil
}

// 辅助方法

// calculateGrowthRate 计算增长率
// 参数: table - 表名, prevStart - 上期开始时间, prevEnd - 上期结束时间, currStart - 当期开始时间, currEnd - 当期结束时间
// 返回: float64 - 增长率
func (s *analyticsService) calculateGrowthRate(table string, prevStart, prevEnd, currStart, currEnd time.Time) float64 {
	var prevCount, currCount int64
	
	// 获取上期数量
	s.db.Table(table).Where("created_at BETWEEN ? AND ?", prevStart, prevEnd).Count(&prevCount)
	
	// 获取当期数量
	s.db.Table(table).Where("created_at BETWEEN ? AND ?", currStart, currEnd).Count(&currCount)
	
	if prevCount == 0 {
		if currCount > 0 {
			return 100.0 // 从0增长到有数据，视为100%增长
		}
		return 0.0
	}
	
	return float64(currCount-prevCount) / float64(prevCount) * 100
}
package models

import (
	"time"
)

// Analytics 分析统计模型
// 存储系统的各种统计数据
type Analytics struct {
	BaseModel
	Date        time.Time `gorm:"uniqueIndex:idx_date_type;not null" json:"date"`         // 统计日期
	Type        string    `gorm:"uniqueIndex:idx_date_type;size:50;not null" json:"type"` // 统计类型
	Value       int64     `gorm:"not null" json:"value"`                                  // 统计值
	Metadata    string    `gorm:"type:json" json:"metadata,omitempty"`                    // 额外元数据
	Description string    `gorm:"size:255" json:"description"`                            // 描述
}

// TableName 自定义表名
func (Analytics) TableName() string {
	return "analytics"
}

// DashboardStats 仪表板统计结构体
// 用于展示系统概览数据
type DashboardStats struct {
	TotalUsers         int64     `json:"total_users"`         // 总用户数
	TotalPosts         int64     `json:"total_posts"`         // 总文章数
	TotalComments      int64     `json:"total_comments"`      // 总评论数
	TotalLikes         int64     `json:"total_likes"`         // 总点赞数
	TotalViews         int64     `json:"total_views"`         // 总浏览数
	ActiveUsers        int64     `json:"active_users"`        // 活跃用户数
	PublishedPosts     int64     `json:"published_posts"`     // 已发布文章数
	PendingComments    int64     `json:"pending_comments"`    // 待审核评论数
	NewUsersToday      int64     `json:"new_users_today"`     // 今日新用户
	NewPostsToday      int64     `json:"new_posts_today"`     // 今日新文章
	NewCommentsToday   int64     `json:"new_comments_today"`  // 今日新评论
	// 新增字段以修复编译错误
	TodayUsers         int64     `json:"today_users"`         // 今日用户数
	TodayPosts         int64     `json:"today_posts"`         // 今日文章数
	TodayComments      int64     `json:"today_comments"`      // 今日评论数
	TodayViews         int64     `json:"today_views"`         // 今日浏览数
	UserGrowthRate     float64   `json:"user_growth_rate"`    // 用户增长率
	PostGrowthRate     float64   `json:"post_growth_rate"`    // 文章增长率
	CommentGrowthRate  float64   `json:"comment_growth_rate"` // 评论增长率
	LastUpdated        time.Time `json:"last_updated"`        // 最后更新时间
}

// PopularPost 热门文章结构体
// 用于展示热门文章数据
type PopularPost struct {
	ID           uint      `json:"id"`            // 文章ID
	Title        string    `json:"title"`         // 文章标题
	Slug         string    `json:"slug"`          // URL别名
	AuthorName   string    `json:"author_name"`   // 作者名称
	ViewCount    int       `json:"view_count"`    // 浏览次数
	LikeCount    int       `json:"like_count"`    // 点赞次数
	CommentCount int       `json:"comment_count"` // 评论次数
	PublishedAt  time.Time `json:"published_at"`  // 发布时间
	Score        float64   `json:"score"`         // 热度评分
}

// ActiveUser 活跃用户结构体
// 用于展示活跃用户数据
type ActiveUser struct {
	ID             uint      `json:"id"`              // 用户ID
	Username       string    `json:"username"`        // 用户名
	Nickname       string    `json:"nickname"`        // 昵称
	Avatar         string    `json:"avatar"`          // 头像
	PostsCount     int       `json:"posts_count"`     // 文章数量
	CommentsCount  int       `json:"comments_count"`  // 评论数量
	LikesReceived  int       `json:"likes_received"`  // 获得点赞数
	FollowersCount int       `json:"followers_count"` // 粉丝数量
	LastActiveAt   time.Time `json:"last_active_at"`  // 最后活跃时间
	ActivityScore  float64   `json:"activity_score"`  // 活跃度评分
}

// CategoryStats 分类统计结构体
// 用于展示分类相关统计数据
type CategoryStats struct {
	ID          uint    `json:"id"`           // 分类ID
	Name        string  `json:"name"`         // 分类名称
	Slug        string  `json:"slug"`         // URL别名
	PostsCount  int     `json:"posts_count"`  // 文章数量
	ViewsCount  int64   `json:"views_count"`  // 浏览总数
	LikesCount  int64   `json:"likes_count"`  // 点赞总数
	Percentage  float64 `json:"percentage"`   // 占比
	GrowthRate  float64 `json:"growth_rate"`  // 增长率
	LastPostAt  *time.Time `json:"last_post_at,omitempty"` // 最后发文时间
}

// TagStats 标签统计结构体
// 用于展示标签相关统计数据
type TagStats struct {
	ID         uint    `json:"id"`          // 标签ID
	Name       string  `json:"name"`        // 标签名称
	Slug       string  `json:"slug"`        // URL别名
	Color      string  `json:"color"`       // 标签颜色
	PostsCount int     `json:"posts_count"` // 文章数量
	ViewsCount int64   `json:"views_count"` // 浏览总数
	Popularity float64 `json:"popularity"`  // 热门度
	Trending   bool    `json:"trending"`    // 是否趋势
}

// UserGrowthStats 用户增长统计结构体
// 用于展示用户增长趋势数据
type UserGrowthStats struct {
	Date           time.Time `json:"date"`            // 日期
	NewUsers       int       `json:"new_users"`       // 新增用户
	ActiveUsers    int       `json:"active_users"`    // 活跃用户
	RetainedUsers  int       `json:"retained_users"`  // 留存用户
	TotalUsers     int       `json:"total_users"`     // 总用户数
	GrowthRate     float64   `json:"growth_rate"`     // 增长率
	RetentionRate  float64   `json:"retention_rate"`  // 留存率
}

// ContentStats 内容统计结构体
// 用于展示内容相关统计数据
type ContentStats struct {
	Date                  time.Time `json:"date"`                    // 日期
	NewPosts              int       `json:"new_posts"`              // 新增文章
	NewComments           int       `json:"new_comments"`           // 新增评论
	TotalPosts            int64     `json:"total_posts"`            // 总文章数
	PublishedPosts        int64     `json:"published_posts"`        // 已发布文章数
	DraftPosts            int64     `json:"draft_posts"`            // 草稿文章数
	TotalComments         int64     `json:"total_comments"`         // 总评论数
	ApprovedComments      int64     `json:"approved_comments"`      // 已审核评论数
	PendingComments       int64     `json:"pending_comments"`       // 待审核评论数
	TotalViews            int64     `json:"total_views"`            // 总浏览量
	TotalLikes            int64     `json:"total_likes"`            // 总点赞数
	TotalCategories       int64     `json:"total_categories"`       // 总分类数
	TotalTags             int64     `json:"total_tags"`             // 总标签数
	MostPopularCategory   string    `json:"most_popular_category"`  // 最受欢迎的分类
	AvgPostLength         float64   `json:"avg_post_length"`        // 平均文章长度
	AvgReadTime           float64   `json:"avg_read_time"`          // 平均阅读时间
	AveragePostLength     float64   `json:"average_post_length"`    // 平均文章长度
	AverageCommentLength  float64   `json:"average_comment_length"` // 平均评论长度
	EngagementRate        float64   `json:"engagement_rate"`        // 参与率
	AverageReadTime       float64   `json:"average_read_time"`      // 平均阅读时间
}

// AnalyticsMethods 分析统计模型的方法

// GetAnalyticsByDateRange 根据日期范围获取分析数据
// 这是一个示例方法，实际实现应该在repository层
func (a *Analytics) GetAnalyticsByDateRange(startDate, endDate time.Time, analyticsType string) []Analytics {
	// 这里应该是数据库查询逻辑
	// 实际实现会在repository层完成
	return []Analytics{}
}

// CalculateGrowthRate 计算增长率
// 参数: currentValue - 当前值, previousValue - 之前值
// 返回: float64 - 增长率百分比
func CalculateGrowthRate(currentValue, previousValue int64) float64 {
	if previousValue == 0 {
		if currentValue > 0 {
			return 100.0 // 从0开始增长视为100%
		}
		return 0.0
	}
	return float64(currentValue-previousValue) / float64(previousValue) * 100.0
}

// CalculateEngagementRate 计算参与率
// 参数: interactions - 互动数, views - 浏览数
// 返回: float64 - 参与率百分比
func CalculateEngagementRate(interactions, views int64) float64 {
	if views == 0 {
		return 0.0
	}
	return float64(interactions) / float64(views) * 100.0
}

// CalculateRetentionRate 计算留存率
// 参数: retainedUsers - 留存用户数, totalUsers - 总用户数
// 返回: float64 - 留存率百分比
func CalculateRetentionRate(retainedUsers, totalUsers int64) float64 {
	if totalUsers == 0 {
		return 0.0
	}
	return float64(retainedUsers) / float64(totalUsers) * 100.0
}

// CalculatePopularityScore 计算热门度评分
// 参数: views - 浏览数, likes - 点赞数, comments - 评论数, days - 天数
// 返回: float64 - 热门度评分
func CalculatePopularityScore(views, likes, comments int64, days int) float64 {
	if days <= 0 {
		days = 1
	}
	
	// 权重: 浏览数 1分，点赞数 3分，评论数 5分
	score := float64(views)*1.0 + float64(likes)*3.0 + float64(comments)*5.0
	
	// 按天数平均
	return score / float64(days)
}

// CalculateActivityScore 计算活跃度评分
// 参数: posts - 文章数, comments - 评论数, likes - 点赞数, days - 天数
// 返回: float64 - 活跃度评分
func CalculateActivityScore(posts, comments, likes int64, days int) float64 {
	if days <= 0 {
		days = 1
	}
	
	// 权重: 文章 10分，评论 3分，点赞 1分
	score := float64(posts)*10.0 + float64(comments)*3.0 + float64(likes)*1.0
	
	// 按天数平均
	return score / float64(days)
}

// DashboardStatsMethods 仪表板统计的方法

// UpdateLastUpdated 更新最后更新时间
func (ds *DashboardStats) UpdateLastUpdated() {
	ds.LastUpdated = time.Now()
}

// GetTotalEngagement 获取总参与度
// 返回: int64 - 总参与度（评论数 + 点赞数）
func (ds *DashboardStats) GetTotalEngagement() int64 {
	return ds.TotalComments + ds.TotalLikes
}

// GetEngagementRate 获取参与率
// 返回: float64 - 参与率百分比
func (ds *DashboardStats) GetEngagementRate() float64 {
	return CalculateEngagementRate(ds.GetTotalEngagement(), ds.TotalViews)
}

// PopularPostMethods 热门文章的方法

// UpdateScore 更新热门度评分
// 根据当前数据重新计算热门度评分
func (pp *PopularPost) UpdateScore() {
	days := int(time.Since(pp.PublishedAt).Hours() / 24)
	if days < 1 {
		days = 1
	}
	pp.Score = CalculatePopularityScore(int64(pp.ViewCount), int64(pp.LikeCount), int64(pp.CommentCount), days)
}

// GetEngagementCount 获取参与数
// 返回: int - 总参与数（点赞数 + 评论数）
func (pp *PopularPost) GetEngagementCount() int {
	return pp.LikeCount + pp.CommentCount
}

// ActiveUserMethods 活跃用户的方法

// UpdateActivityScore 更新活跃度评分
// 根据当前数据重新计算活跃度评分
func (au *ActiveUser) UpdateActivityScore() {
	days := int(time.Since(au.LastActiveAt).Hours() / 24)
	if days < 1 {
		days = 1
	}
	au.ActivityScore = CalculateActivityScore(int64(au.PostsCount), int64(au.CommentsCount), int64(au.LikesReceived), days)
}

// IsRecentlyActive 检查是否最近活跃
// 参数: hours - 小时数阈值
// 返回: bool - 是否在指定小时内活跃
func (au *ActiveUser) IsRecentlyActive(hours int) bool {
	return time.Since(au.LastActiveAt).Hours() <= float64(hours)
}
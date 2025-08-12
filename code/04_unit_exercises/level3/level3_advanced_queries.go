// 04_unit_exercises/level3_advanced_queries.go - Level 3 高级查询练习
// 对应文档：03_GORM单元练习_基础技能训练.md

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 数据模型定义
type User struct {
	BaseModel
	Username string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Age      int       `gorm:"check:age >= 0 AND age <= 150" json:"age"`
	City     string    `gorm:"size:100;index" json:"city"`
	Salary   float64   `gorm:"precision:10;scale:2" json:"salary"`
	JoinDate time.Time `gorm:"index" json:"join_date"`
	IsActive bool      `gorm:"default:true;index" json:"is_active"`
	
	// 关联关系
	Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
}

type Category struct {
	BaseModel
	Name        string `gorm:"size:100;not null;index" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
	
	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

type Post struct {
	BaseModel
	Title       string     `gorm:"size:200;not null;index" json:"title"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Status      string     `gorm:"size:20;default:'draft';index" json:"status"`
	ViewCount   int        `gorm:"default:0;index" json:"view_count"`
	LikeCount   int        `gorm:"default:0;index" json:"like_count"`
	PublishedAt *time.Time `gorm:"index" json:"published_at"`
	Rating      float64    `gorm:"precision:3;scale:2;default:0" json:"rating"`
	
	// 外键
	AuthorID   uint  `gorm:"not null;index" json:"author_id"`
	CategoryID *uint `gorm:"index" json:"category_id"`
	
	// 关联关系
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

type Comment struct {
	BaseModel
	Content   string `gorm:"type:text;not null" json:"content"`
	Status    string `gorm:"size:20;default:'pending';index" json:"status"`
	LikeCount int    `gorm:"default:0" json:"like_count"`
	
	// 外键
	PostID   uint `gorm:"not null;index" json:"post_id"`
	AuthorID uint `gorm:"not null;index" json:"author_id"`
	
	// 关联关系
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

type Tag struct {
	BaseModel
	Name     string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug     string `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// 统计结构
type UserStats struct {
	UserID      uint    `json:"user_id"`
	Username    string  `json:"username"`
	PostCount   int64   `json:"post_count"`
	CommentCount int64  `json:"comment_count"`
	TotalViews  int64   `json:"total_views"`
	TotalLikes  int64   `json:"total_likes"`
	AvgRating   float64 `json:"avg_rating"`
}

type CategoryStats struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	PostCount    int64   `json:"post_count"`
	TotalViews   int64   `json:"total_views"`
	AvgRating    float64 `json:"avg_rating"`
}

type MonthlyStats struct {
	Year      int   `json:"year"`
	Month     int   `json:"month"`
	PostCount int64 `json:"post_count"`
	UserCount int64 `json:"user_count"`
}

// 数据库初始化
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("level3_advanced_queries.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Category{}, &Post{}, &Comment{}, &Tag{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db
}

// 创建测试数据
func createTestData(db *gorm.DB) {
	// 创建用户
	users := []User{
		{Username: "alice", Email: "alice@example.com", Age: 25, City: "北京", Salary: 8000, JoinDate: time.Now().AddDate(-1, 0, 0), IsActive: true},
		{Username: "bob", Email: "bob@example.com", Age: 30, City: "上海", Salary: 12000, JoinDate: time.Now().AddDate(-2, 0, 0), IsActive: true},
		{Username: "charlie", Email: "charlie@example.com", Age: 28, City: "深圳", Salary: 10000, JoinDate: time.Now().AddDate(-1, -6, 0), IsActive: false},
		{Username: "diana", Email: "diana@example.com", Age: 35, City: "北京", Salary: 15000, JoinDate: time.Now().AddDate(-3, 0, 0), IsActive: true},
		{Username: "eve", Email: "eve@example.com", Age: 22, City: "广州", Salary: 6000, JoinDate: time.Now().AddDate(0, -3, 0), IsActive: true},
	}
	db.Create(&users)

	// 创建分类
	categories := []Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", IsActive: true},
		{Name: "生活", Slug: "life", Description: "生活分享", IsActive: true},
		{Name: "旅游", Slug: "travel", Description: "旅游攻略", IsActive: true},
		{Name: "美食", Slug: "food", Description: "美食推荐", IsActive: false},
	}
	db.Create(&categories)

	// 创建标签
	tags := []Tag{
		{Name: "Go", Slug: "go", IsActive: true},
		{Name: "数据库", Slug: "database", IsActive: true},
		{Name: "前端", Slug: "frontend", IsActive: true},
		{Name: "后端", Slug: "backend", IsActive: true},
		{Name: "教程", Slug: "tutorial", IsActive: true},
	}
	db.Create(&tags)

	// 创建文章
	posts := []Post{
		{Title: "Go语言入门", Content: "Go语言基础教程...", Status: "published", ViewCount: 1500, LikeCount: 120, Rating: 4.5, AuthorID: 1, CategoryID: &[]uint{1}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, -1, 0)}[0]},
		{Title: "数据库设计原则", Content: "数据库设计的基本原则...", Status: "published", ViewCount: 800, LikeCount: 65, Rating: 4.2, AuthorID: 2, CategoryID: &[]uint{1}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, -2, 0)}[0]},
		{Title: "我的北京生活", Content: "在北京的生活感悟...", Status: "published", ViewCount: 300, LikeCount: 25, Rating: 3.8, AuthorID: 1, CategoryID: &[]uint{2}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, 0, -15)}[0]},
		{Title: "上海旅游攻略", Content: "上海必去景点推荐...", Status: "published", ViewCount: 1200, LikeCount: 95, Rating: 4.7, AuthorID: 3, CategoryID: &[]uint{3}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, 0, -30)}[0]},
		{Title: "React入门教程", Content: "React基础知识...", Status: "draft", ViewCount: 0, LikeCount: 0, Rating: 0, AuthorID: 4, CategoryID: &[]uint{1}[0]},
		{Title: "深圳美食推荐", Content: "深圳好吃的餐厅...", Status: "published", ViewCount: 600, LikeCount: 40, Rating: 4.0, AuthorID: 5, CategoryID: &[]uint{4}[0], PublishedAt: &[]time.Time{time.Now().AddDate(0, 0, -7)}[0]},
	}
	db.Create(&posts)

	// 创建评论
	comments := []Comment{
		{Content: "很好的教程！", Status: "approved", LikeCount: 10, PostID: 1, AuthorID: 2},
		{Content: "学到了很多", Status: "approved", LikeCount: 5, PostID: 1, AuthorID: 3},
		{Content: "写得不错", Status: "approved", LikeCount: 3, PostID: 2, AuthorID: 1},
		{Content: "有用的信息", Status: "pending", LikeCount: 0, PostID: 3, AuthorID: 4},
		{Content: "期待更多内容", Status: "approved", LikeCount: 8, PostID: 4, AuthorID: 5},
	}
	db.Create(&comments)

	// 为文章添加标签
	var post1, post2 Post
	db.First(&post1, 1)
	db.First(&post2, 2)
	var tag1, tag2, tag5 Tag
	db.First(&tag1, 1)
	db.First(&tag2, 2)
	db.First(&tag5, 5)
	db.Model(&post1).Association("Tags").Append([]Tag{tag1, tag5})
	db.Model(&post2).Association("Tags").Append([]Tag{tag2, tag5})
}

// 练习1：条件查询和排序

// FindUsersByConditions 多条件查询用户
func FindUsersByConditions(db *gorm.DB, minAge, maxAge int, cities []string, isActive bool) ([]User, error) {
	var users []User
	query := db.Where("age BETWEEN ? AND ?", minAge, maxAge)
	
	if len(cities) > 0 {
		query = query.Where("city IN ?", cities)
	}
	
	query = query.Where("is_active = ?", isActive)
	
	result := query.Order("salary DESC, age ASC").Find(&users)
	return users, result.Error
}

// FindPostsByDateRange 按日期范围查询文章
func FindPostsByDateRange(db *gorm.DB, startDate, endDate time.Time, status string) ([]Post, error) {
	var posts []Post
	query := db.Where("published_at BETWEEN ? AND ?", startDate, endDate)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	result := query.Order("view_count DESC").Preload("Author").Preload("Category").Find(&posts)
	return posts, result.Error
}

// FindTopRatedPosts 查询高评分文章
func FindTopRatedPosts(db *gorm.DB, minRating float64, limit int) ([]Post, error) {
	var posts []Post
	result := db.Where("rating >= ? AND status = ?", minRating, "published").
		Order("rating DESC, like_count DESC").
		Limit(limit).
		Preload("Author").
		Preload("Category").
		Find(&posts)
	return posts, result.Error
}

// 练习2：聚合查询

// GetUserStatistics 获取用户统计信息
func GetUserStatistics(db *gorm.DB) ([]UserStats, error) {
	var stats []UserStats
	err := db.Table("users u").
		Select(`
			u.id as user_id,
			u.username,
			COUNT(DISTINCT p.id) as post_count,
			COUNT(DISTINCT c.id) as comment_count,
			COALESCE(SUM(p.view_count), 0) as total_views,
			COALESCE(SUM(p.like_count), 0) as total_likes,
			COALESCE(AVG(p.rating), 0) as avg_rating
		`).
		Joins("LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Group("u.id, u.username").
		Order("total_views DESC").
		Scan(&stats).Error
	return stats, err
}

// GetCategoryStatistics 获取分类统计信息
func GetCategoryStatistics(db *gorm.DB) ([]CategoryStats, error) {
	var stats []CategoryStats
	err := db.Table("categories c").
		Select(`
			c.id as category_id,
			c.name as category_name,
			COUNT(p.id) as post_count,
			COALESCE(SUM(p.view_count), 0) as total_views,
			COALESCE(AVG(p.rating), 0) as avg_rating
		`).
		Joins("LEFT JOIN posts p ON c.id = p.category_id AND p.deleted_at IS NULL AND p.status = 'published'").
		Where("c.deleted_at IS NULL AND c.is_active = ?", true).
		Group("c.id, c.name").
		Order("post_count DESC").
		Scan(&stats).Error
	return stats, err
}

// GetMonthlyStatistics 获取月度统计
func GetMonthlyStatistics(db *gorm.DB, months int) ([]MonthlyStats, error) {
	var stats []MonthlyStats
	startDate := time.Now().AddDate(0, -months, 0)
	
	err := db.Raw(`
		SELECT 
			strftime('%Y', date) as year,
			strftime('%m', date) as month,
			COUNT(CASE WHEN type = 'post' THEN 1 END) as post_count,
			COUNT(CASE WHEN type = 'user' THEN 1 END) as user_count
		FROM (
			SELECT created_at as date, 'post' as type FROM posts WHERE created_at >= ? AND deleted_at IS NULL
			UNION ALL
			SELECT created_at as date, 'user' as type FROM users WHERE created_at >= ? AND deleted_at IS NULL
		) combined
		GROUP BY year, month
		ORDER BY year DESC, month DESC
	`, startDate, startDate).Scan(&stats).Error
	
	return stats, err
}

// 练习3：子查询

// FindUsersWithMostPosts 查询发文最多的用户
func FindUsersWithMostPosts(db *gorm.DB, limit int) ([]User, error) {
	var users []User
	subQuery := db.Table("posts").Select("author_id, COUNT(*) as post_count").
		Where("deleted_at IS NULL AND status = ?", "published").
		Group("author_id").
		Order("post_count DESC").
		Limit(limit)
	
	err := db.Table("users u").
		Joins("JOIN (?) pc ON u.id = pc.author_id", subQuery).
		Where("u.deleted_at IS NULL").
		Order("pc.post_count DESC").
		Find(&users).Error
	
	return users, err
}

// FindPostsAboveAverageViews 查询浏览量高于平均值的文章
func FindPostsAboveAverageViews(db *gorm.DB) ([]Post, error) {
	var posts []Post
	
	// 先获取平均浏览量
	var avgViews float64
	db.Table("posts").Where("deleted_at IS NULL AND status = ?", "published").Select("AVG(view_count)").Scan(&avgViews)
	
	// 查询高于平均值的文章
	err := db.Where("view_count > ? AND status = ?", avgViews, "published").
		Order("view_count DESC").
		Preload("Author").
		Preload("Category").
		Find(&posts).Error
	
	return posts, err
}

// FindUsersWithNoComments 查询没有评论的用户
func FindUsersWithNoComments(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Where("id NOT IN (?)", 
		db.Table("comments").Select("DISTINCT author_id").Where("deleted_at IS NULL"),
	).Find(&users).Error
	return users, err
}

// 练习4：复杂连接查询

// GetPostsWithAuthorAndCommentCount 获取文章及作者信息和评论数
func GetPostsWithAuthorAndCommentCount(db *gorm.DB, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := db.Table("posts p").
		Select(`
			p.id,
			p.title,
			p.view_count,
			p.like_count,
			p.rating,
			p.published_at,
			u.username as author_name,
			u.city as author_city,
			c.name as category_name,
			COUNT(cm.id) as comment_count
		`).
		Joins("JOIN users u ON p.author_id = u.id").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Joins("LEFT JOIN comments cm ON p.id = cm.post_id AND cm.deleted_at IS NULL").
		Where("p.deleted_at IS NULL AND p.status = ?", "published").
		Group("p.id, p.title, p.view_count, p.like_count, p.rating, p.published_at, u.username, u.city, c.name").
		Order("p.view_count DESC").
		Limit(limit).
		Scan(&results).Error
	
	return results, err
}

// GetUserEngagementReport 获取用户参与度报告
func GetUserEngagementReport(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := db.Table("users u").
		Select(`
			u.id,
			u.username,
			u.city,
			u.join_date,
			COUNT(DISTINCT p.id) as post_count,
			COUNT(DISTINCT c.id) as comment_count,
			COALESCE(SUM(p.view_count), 0) as total_post_views,
			COALESCE(SUM(p.like_count), 0) as total_post_likes,
			COALESCE(SUM(c.like_count), 0) as total_comment_likes,
			COALESCE(AVG(p.rating), 0) as avg_post_rating,
			(COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) as total_activity
		`).
		Joins("LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Group("u.id, u.username, u.city, u.join_date").
		Having("total_activity > 0").
		Order("total_activity DESC").
		Scan(&results).Error
	
	return results, err
}

// 练习5：窗口函数和排名

// GetTopPostsByCategory 获取每个分类的热门文章
func GetTopPostsByCategory(db *gorm.DB, topN int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := db.Raw(`
		SELECT 
			category_name,
			title,
			view_count,
			like_count,
			rating,
			rank_in_category
		FROM (
			SELECT 
				c.name as category_name,
				p.title,
				p.view_count,
				p.like_count,
				p.rating,
				ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY p.view_count DESC) as rank_in_category
			FROM posts p
			JOIN categories c ON p.category_id = c.id
			WHERE p.deleted_at IS NULL AND p.status = 'published'
			  AND c.deleted_at IS NULL AND c.is_active = 1
		) ranked
		WHERE rank_in_category <= ?
		ORDER BY category_name, rank_in_category
	`, topN).Scan(&results).Error
	
	return results, err
}

// GetUserRankingByActivity 获取用户活跃度排名
func GetUserRankingByActivity(db *gorm.DB) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := db.Raw(`
		SELECT 
			username,
			city,
			total_activity,
			activity_rank,
			city_rank,
			CASE 
				WHEN activity_rank <= 3 THEN 'Gold'
				WHEN activity_rank <= 10 THEN 'Silver'
				ELSE 'Bronze'
			END as badge
		FROM (
			SELECT 
				u.username,
				u.city,
				(COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) as total_activity,
				ROW_NUMBER() OVER (ORDER BY (COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) DESC) as activity_rank,
				ROW_NUMBER() OVER (PARTITION BY u.city ORDER BY (COUNT(DISTINCT p.id) + COUNT(DISTINCT c.id)) DESC) as city_rank
			FROM users u
			LEFT JOIN posts p ON u.id = p.author_id AND p.deleted_at IS NULL
			LEFT JOIN comments c ON u.id = c.author_id AND c.deleted_at IS NULL
			WHERE u.deleted_at IS NULL
			GROUP BY u.id, u.username, u.city
			HAVING total_activity > 0
		) ranked
		ORDER BY activity_rank
	`).Scan(&results).Error
	
	return results, err
}

// 练习6：性能优化查询

// GetPostsWithOptimizedPreloading 优化预加载的文章查询
func GetPostsWithOptimizedPreloading(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post
	
	// 使用选择性预加载，只加载需要的字段
	err := db.Select("id, title, content, view_count, like_count, rating, published_at, author_id, category_id").
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, city")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, slug").Where("is_active = ?", true)
		}).
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	
	return posts, err
}

// GetPostsWithPagination 分页查询文章
func GetPostsWithPagination(db *gorm.DB, page, pageSize int, categoryID *uint, search string) ([]Post, int64, error) {
	var posts []Post
	var total int64
	
	// 构建查询条件
	query := db.Model(&Post{}).Where("status = ?", "published")
	
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	// 计算总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Order("published_at DESC").
		Find(&posts).Error
	
	return posts, total, err
}

// 主函数演示
func main() {
	fmt.Println("=== GORM Level 3 高级查询练习 ===")

	// 初始化数据库
	db := initDB()
	fmt.Println("✓ 数据库初始化完成")

	// 创建测试数据
	createTestData(db)
	fmt.Println("✓ 测试数据创建完成")

	// 练习1：条件查询和排序
	fmt.Println("\n=== 条件查询和排序 ===")
	
	users, err := FindUsersByConditions(db, 20, 35, []string{"北京", "上海"}, true)
	if err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个符合条件的用户\n", len(users))
		for _, user := range users {
			fmt.Printf("  - %s (%d岁, %s, 薪资: %.0f)\n", user.Username, user.Age, user.City, user.Salary)
		}
	}

	startDate := time.Now().AddDate(0, -3, 0)
	endDate := time.Now()
	posts, err := FindPostsByDateRange(db, startDate, endDate, "published")
	if err != nil {
		fmt.Printf("查询文章失败: %v\n", err)
	} else {
		fmt.Printf("\n最近3个月发布的文章 (%d篇):\n", len(posts))
		for _, post := range posts {
			fmt.Printf("  - %s (浏览: %d, 点赞: %d, 作者: %s)\n", 
				post.Title, post.ViewCount, post.LikeCount, post.Author.Username)
		}
	}

	topPosts, err := FindTopRatedPosts(db, 4.0, 3)
	if err != nil {
		fmt.Printf("查询高评分文章失败: %v\n", err)
	} else {
		fmt.Printf("\n高评分文章 (评分>=4.0, 前3篇):\n")
		for _, post := range topPosts {
			fmt.Printf("  - %s (评分: %.1f, 浏览: %d)\n", 
				post.Title, post.Rating, post.ViewCount)
		}
	}

	// 练习2：聚合查询
	fmt.Println("\n=== 聚合查询 ===")
	
	userStats, err := GetUserStatistics(db)
	if err != nil {
		fmt.Printf("获取用户统计失败: %v\n", err)
	} else {
		fmt.Println("用户统计信息:")
		for _, stat := range userStats {
			fmt.Printf("  - %s: 文章%d篇, 评论%d条, 总浏览%d, 总点赞%d, 平均评分%.1f\n",
				stat.Username, stat.PostCount, stat.CommentCount, 
				stat.TotalViews, stat.TotalLikes, stat.AvgRating)
		}
	}

	categoryStats, err := GetCategoryStatistics(db)
	if err != nil {
		fmt.Printf("获取分类统计失败: %v\n", err)
	} else {
		fmt.Println("\n分类统计信息:")
		for _, stat := range categoryStats {
			fmt.Printf("  - %s: 文章%d篇, 总浏览%d, 平均评分%.1f\n",
				stat.CategoryName, stat.PostCount, stat.TotalViews, stat.AvgRating)
		}
	}

	monthlyStats, err := GetMonthlyStatistics(db, 6)
	if err != nil {
		fmt.Printf("获取月度统计失败: %v\n", err)
	} else {
		fmt.Println("\n最近6个月统计:")
		for _, stat := range monthlyStats {
			fmt.Printf("  - %d年%02d月: 文章%d篇, 新用户%d人\n",
				stat.Year, stat.Month, stat.PostCount, stat.UserCount)
		}
	}

	// 练习3：子查询
	fmt.Println("\n=== 子查询 ===")
	
	topUsers, err := FindUsersWithMostPosts(db, 3)
	if err != nil {
		fmt.Printf("查询发文最多用户失败: %v\n", err)
	} else {
		fmt.Printf("发文最多的用户 (前3名):\n")
		for _, user := range topUsers {
			fmt.Printf("  - %s (%s)\n", user.Username, user.City)
		}
	}

	highViewPosts, err := FindPostsAboveAverageViews(db)
	if err != nil {
		fmt.Printf("查询高浏览量文章失败: %v\n", err)
	} else {
		fmt.Printf("\n浏览量高于平均值的文章 (%d篇):\n", len(highViewPosts))
		for _, post := range highViewPosts {
			fmt.Printf("  - %s (浏览: %d)\n", post.Title, post.ViewCount)
		}
	}

	noCommentUsers, err := FindUsersWithNoComments(db)
	if err != nil {
		fmt.Printf("查询无评论用户失败: %v\n", err)
	} else {
		fmt.Printf("\n没有发表评论的用户 (%d人):\n", len(noCommentUsers))
		for _, user := range noCommentUsers {
			fmt.Printf("  - %s\n", user.Username)
		}
	}

	// 练习4：复杂连接查询
	fmt.Println("\n=== 复杂连接查询 ===")
	
	postDetails, err := GetPostsWithAuthorAndCommentCount(db, 5)
	if err != nil {
		fmt.Printf("获取文章详情失败: %v\n", err)
	} else {
		fmt.Println("文章详情 (包含作者和评论数):")
		for _, detail := range postDetails {
			fmt.Printf("  - %s (作者: %s, 分类: %s, 评论: %v条)\n",
				detail["title"], detail["author_name"], 
				detail["category_name"], detail["comment_count"])
		}
	}

	engagementReport, err := GetUserEngagementReport(db)
	if err != nil {
		fmt.Printf("获取用户参与度报告失败: %v\n", err)
	} else {
		fmt.Println("\n用户参与度报告:")
		for _, report := range engagementReport {
			fmt.Printf("  - %s: 总活动%v次 (文章%v篇, 评论%v条)\n",
				report["username"], report["total_activity"],
				report["post_count"], report["comment_count"])
		}
	}

	// 练习5：窗口函数和排名
	fmt.Println("\n=== 窗口函数和排名 ===")
	
	topByCategory, err := GetTopPostsByCategory(db, 2)
	if err != nil {
		fmt.Printf("获取分类热门文章失败: %v\n", err)
	} else {
		fmt.Println("各分类热门文章 (前2名):")
		for _, item := range topByCategory {
			fmt.Printf("  - [%s] %s (浏览: %v, 排名: %v)\n",
				item["category_name"], item["title"],
				item["view_count"], item["rank_in_category"])
		}
	}

	userRanking, err := GetUserRankingByActivity(db)
	if err != nil {
		fmt.Printf("获取用户活跃度排名失败: %v\n", err)
	} else {
		fmt.Println("\n用户活跃度排名:")
		for _, rank := range userRanking {
			fmt.Printf("  - %s (%s): 活动%v次, 全站排名%v, 城市排名%v, 徽章: %s\n",
				rank["username"], rank["city"], rank["total_activity"],
				rank["activity_rank"], rank["city_rank"], rank["badge"])
		}
	}

	// 练习6：性能优化查询
	fmt.Println("\n=== 性能优化查询 ===")
	
	optimizedPosts, err := GetPostsWithOptimizedPreloading(db, 3)
	if err != nil {
		fmt.Printf("优化预加载查询失败: %v\n", err)
	} else {
		fmt.Printf("优化预加载查询结果 (%d篇文章):\n", len(optimizedPosts))
		for _, post := range optimizedPosts {
			fmt.Printf("  - %s (作者: %s, 分类: %s, 标签数: %d)\n",
				post.Title, post.Author.Username, 
				func() string {
					if post.Category != nil {
						return post.Category.Name
					}
					return "无分类"
				}(), len(post.Tags))
		}
	}

	paginatedPosts, total, err := GetPostsWithPagination(db, 1, 3, nil, "")
	if err != nil {
		fmt.Printf("分页查询失败: %v\n", err)
	} else {
		fmt.Printf("\n分页查询结果 (第1页, 每页3篇, 总共%d篇):\n", total)
		for _, post := range paginatedPosts {
			fmt.Printf("  - %s (作者: %s)\n", post.Title, post.Author.Username)
		}
	}

	fmt.Println("\n=== Level 3 高级查询练习完成 ===")
}
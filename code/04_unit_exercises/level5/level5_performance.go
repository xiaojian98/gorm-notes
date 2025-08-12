// 04_unit_exercises/level5_performance.go - Level 5 性能优化练习
// 对应文档：03_GORM单元练习_基础技能训练.md

package main

import (
	"fmt"
	"log"
	"math/rand"
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
	Username    string `gorm:"uniqueIndex:idx_username;size:50;not null" json:"username"`
	Email       string `gorm:"uniqueIndex:idx_email;size:100;not null" json:"email"`
	FirstName   string `gorm:"size:50;not null;index:idx_name" json:"first_name"`
	LastName    string `gorm:"size:50;not null;index:idx_name" json:"last_name"`
	Age         int    `gorm:"check:age >= 0 AND age <= 150;index:idx_age" json:"age"`
	City        string `gorm:"size:100;index:idx_location" json:"city"`
	Country     string `gorm:"size:100;index:idx_location" json:"country"`
	Salary      float64 `gorm:"precision:10;scale:2;index:idx_salary" json:"salary"`
	Department  string `gorm:"size:100;index:idx_department" json:"department"`
	Position    string `gorm:"size:100;index:idx_position" json:"position"`
	JoinDate    time.Time `gorm:"index:idx_join_date" json:"join_date"`
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`
	LastLoginAt *time.Time `gorm:"index:idx_last_login" json:"last_login_at"`
	
	// 关联关系
	Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
	Orders   []Order   `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

type Category struct {
	BaseModel
	Name        string `gorm:"size:100;not null;index:idx_category_name" json:"name"`
	Slug        string `gorm:"uniqueIndex:idx_category_slug;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	ParentID    *uint  `gorm:"index:idx_parent" json:"parent_id"`
	Level       int    `gorm:"default:1;index:idx_level" json:"level"`
	SortOrder   int    `gorm:"default:0;index:idx_sort" json:"sort_order"`
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`
	
	// 关联关系
	Parent   *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Posts    []Post    `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

type Post struct {
	BaseModel
	Title       string     `gorm:"size:200;not null;index:idx_title" json:"title"`
	Slug        string     `gorm:"uniqueIndex:idx_post_slug;size:200;not null" json:"slug"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Excerpt     string     `gorm:"size:500" json:"excerpt"`
	Status      string     `gorm:"size:20;default:'draft';index:idx_status" json:"status"`
	ViewCount   int        `gorm:"default:0;index:idx_views" json:"view_count"`
	LikeCount   int        `gorm:"default:0;index:idx_likes" json:"like_count"`
	CommentCount int       `gorm:"default:0;index:idx_comments" json:"comment_count"`
	PublishedAt *time.Time `gorm:"index:idx_published" json:"published_at"`
	Rating      float64    `gorm:"precision:3;scale:2;default:0;index:idx_rating" json:"rating"`
	Featured    bool       `gorm:"default:false;index:idx_featured" json:"featured"`
	
	// 外键
	AuthorID   uint  `gorm:"not null;index:idx_author" json:"author_id"`
	CategoryID *uint `gorm:"index:idx_category" json:"category_id"`
	
	// 关联关系
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

type Comment struct {
	BaseModel
	Content   string `gorm:"type:text;not null" json:"content"`
	Status    string `gorm:"size:20;default:'pending';index:idx_status" json:"status"`
	LikeCount int    `gorm:"default:0;index:idx_likes" json:"like_count"`
	ParentID  *uint  `gorm:"index:idx_parent" json:"parent_id"`
	Level     int    `gorm:"default:1;index:idx_level" json:"level"`
	
	// 外键
	PostID   uint `gorm:"not null;index:idx_post" json:"post_id"`
	AuthorID uint `gorm:"not null;index:idx_author" json:"author_id"`
	
	// 关联关系
	Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type Tag struct {
	BaseModel
	Name      string `gorm:"uniqueIndex:idx_tag_name;size:50;not null" json:"name"`
	Slug      string `gorm:"uniqueIndex:idx_tag_slug;size:50;not null" json:"slug"`
	Color     string `gorm:"size:7;default:'#007bff'" json:"color"`
	UsageCount int   `gorm:"default:0;index:idx_usage" json:"usage_count"`
	IsActive  bool   `gorm:"default:true;index:idx_active" json:"is_active"`
	
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

type Order struct {
	BaseModel
	OrderNumber string    `gorm:"uniqueIndex:idx_order_number;size:50;not null" json:"order_number"`
	UserID      uint      `gorm:"not null;index:idx_user" json:"user_id"`
	TotalAmount float64   `gorm:"precision:10;scale:2;not null;index:idx_amount" json:"total_amount"`
	Status      string    `gorm:"size:20;default:'pending';index:idx_status" json:"status"`
	OrderDate   time.Time `gorm:"index:idx_order_date" json:"order_date"`
	ShippedAt   *time.Time `gorm:"index:idx_shipped" json:"shipped_at"`
	
	// 关联关系
	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
	BaseModel
	OrderID   uint    `gorm:"not null;index:idx_order" json:"order_id"`
	ProductID uint    `gorm:"not null;index:idx_product" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"precision:10;scale:2;not null" json:"price"`
	Subtotal  float64 `gorm:"precision:10;scale:2;not null" json:"subtotal"`
	
	// 关联关系
	Order   Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type Product struct {
	BaseModel
	Name        string  `gorm:"size:200;not null;index:idx_name" json:"name"`
	SKU         string  `gorm:"uniqueIndex:idx_sku;size:50;not null" json:"sku"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"precision:10;scale:2;not null;index:idx_price" json:"price"`
	Stock       int     `gorm:"default:0;index:idx_stock" json:"stock"`
	CategoryID  *uint   `gorm:"index:idx_category" json:"category_id"`
	IsActive    bool    `gorm:"default:true;index:idx_active" json:"is_active"`
	
	// 关联关系
	Category   *Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`
}

// 数据库初始化
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("level5_performance.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 关闭日志以提高性能测试准确性
		PrepareStmt: true, // 启用预编译语句
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束以提高性能
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接失败:", err)
	}
	
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Category{}, &Post{}, &Comment{}, &Tag{}, &Order{}, &OrderItem{}, &Product{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 创建复合索引
	createCompositeIndexes(db)

	return db
}

// 创建复合索引
func createCompositeIndexes(db *gorm.DB) {
	// 用户复合索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_active_department ON users(is_active, department)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_city_salary ON users(city, salary)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_join_date_active ON users(join_date, is_active)")
	
	// 文章复合索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_status_published ON posts(status, published_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_author_status ON posts(author_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_category_featured ON posts(category_id, featured)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_views_rating ON posts(view_count, rating)")
	
	// 评论复合索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_status ON comments(post_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_author_created ON comments(author_id, created_at)")
	
	// 订单复合索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_date_amount ON orders(order_date, total_amount)")
	
	fmt.Println("✓ 复合索引创建完成")
}

// 生成测试数据
func generateTestData(db *gorm.DB, userCount, postCount, commentCount int) {
	fmt.Printf("开始生成测试数据: %d用户, %d文章, %d评论\n", userCount, postCount, commentCount)
	start := time.Now()
	
	// 生成用户数据
	users := make([]User, userCount)
	departments := []string{"技术部", "市场部", "销售部", "人事部", "财务部"}
	positions := []string{"工程师", "经理", "主管", "专员", "总监"}
	cities := []string{"北京", "上海", "深圳", "广州", "杭州", "成都", "武汉", "西安"}
	countries := []string{"中国", "美国", "日本", "韩国", "新加坡"}
	
	for i := 0; i < userCount; i++ {
		users[i] = User{
			Username:    fmt.Sprintf("user%d", i+1),
			Email:       fmt.Sprintf("user%d@example.com", i+1),
			FirstName:   fmt.Sprintf("First%d", i+1),
			LastName:    fmt.Sprintf("Last%d", i+1),
			Age:         20 + rand.Intn(40),
			City:        cities[rand.Intn(len(cities))],
			Country:     countries[rand.Intn(len(countries))],
			Salary:      float64(5000 + rand.Intn(20000)),
			Department:  departments[rand.Intn(len(departments))],
			Position:    positions[rand.Intn(len(positions))],
			JoinDate:    time.Now().AddDate(-rand.Intn(5), -rand.Intn(12), -rand.Intn(30)),
			IsActive:    rand.Float32() > 0.1, // 90%的用户是活跃的
		}
		if rand.Float32() > 0.3 { // 70%的用户有最后登录时间
			lastLogin := time.Now().AddDate(0, 0, -rand.Intn(30))
			users[i].LastLoginAt = &lastLogin
		}
	}
	
	// 批量插入用户
	db.CreateInBatches(users, 100)
	fmt.Printf("✓ 用户数据生成完成: %d条\n", userCount)
	
	// 生成分类数据
	categories := []Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", Level: 1, SortOrder: 1, IsActive: true},
		{Name: "生活", Slug: "life", Description: "生活分享", Level: 1, SortOrder: 2, IsActive: true},
		{Name: "旅游", Slug: "travel", Description: "旅游攻略", Level: 1, SortOrder: 3, IsActive: true},
		{Name: "美食", Slug: "food", Description: "美食推荐", Level: 1, SortOrder: 4, IsActive: true},
		{Name: "娱乐", Slug: "entertainment", Description: "娱乐资讯", Level: 1, SortOrder: 5, IsActive: true},
	}
	db.Create(&categories)
	
	// 生成标签数据
	tags := []Tag{
		{Name: "Go", Slug: "go", Color: "#00ADD8", IsActive: true},
		{Name: "Python", Slug: "python", Color: "#3776AB", IsActive: true},
		{Name: "JavaScript", Slug: "javascript", Color: "#F7DF1E", IsActive: true},
		{Name: "数据库", Slug: "database", Color: "#336791", IsActive: true},
		{Name: "前端", Slug: "frontend", Color: "#61DAFB", IsActive: true},
		{Name: "后端", Slug: "backend", Color: "#68217A", IsActive: true},
		{Name: "教程", Slug: "tutorial", Color: "#FF6B6B", IsActive: true},
		{Name: "实战", Slug: "practice", Color: "#4ECDC4", IsActive: true},
	}
	db.Create(&tags)
	
	// 生成文章数据
	posts := make([]Post, postCount)
	statuses := []string{"published", "draft", "archived"}
	
	for i := 0; i < postCount; i++ {
		authorID := uint(rand.Intn(userCount) + 1)
		categoryID := uint(rand.Intn(len(categories)) + 1)
		status := statuses[rand.Intn(len(statuses))]
		
		posts[i] = Post{
			Title:        fmt.Sprintf("文章标题 %d", i+1),
			Slug:         fmt.Sprintf("post-%d", i+1),
			Content:      fmt.Sprintf("这是文章 %d 的内容，包含了丰富的信息和详细的描述...", i+1),
			Excerpt:      fmt.Sprintf("文章 %d 的摘要", i+1),
			Status:       status,
			ViewCount:    rand.Intn(10000),
			LikeCount:    rand.Intn(1000),
			CommentCount: rand.Intn(100),
			Rating:       float64(rand.Intn(50))/10.0 + 1.0, // 1.0-6.0
			Featured:     rand.Float32() > 0.8, // 20%的文章是精选的
			AuthorID:     authorID,
			CategoryID:   &categoryID,
		}
		
		if status == "published" {
			publishedAt := time.Now().AddDate(0, 0, -rand.Intn(365))
			posts[i].PublishedAt = &publishedAt
		}
	}
	
	// 批量插入文章
	db.CreateInBatches(posts, 100)
	fmt.Printf("✓ 文章数据生成完成: %d条\n", postCount)
	
	// 生成评论数据
	comments := make([]Comment, commentCount)
	commentStatuses := []string{"approved", "pending", "spam"}
	
	for i := 0; i < commentCount; i++ {
		comments[i] = Comment{
			Content:   fmt.Sprintf("这是评论 %d 的内容", i+1),
			Status:    commentStatuses[rand.Intn(len(commentStatuses))],
			LikeCount: rand.Intn(100),
			Level:     1,
			PostID:    uint(rand.Intn(postCount) + 1),
			AuthorID:  uint(rand.Intn(userCount) + 1),
		}
	}
	
	// 批量插入评论
	db.CreateInBatches(comments, 100)
	fmt.Printf("✓ 评论数据生成完成: %d条\n", commentCount)
	
	elapsed := time.Since(start)
	fmt.Printf("✓ 测试数据生成完成，耗时: %v\n", elapsed)
}

// 性能测试函数

// 测试1：基本查询性能
func benchmarkBasicQueries(db *gorm.DB) {
	fmt.Println("\n=== 基本查询性能测试 ===")
	
	// 测试简单查询
	start := time.Now()
	var users []User
	db.Where("is_active = ?", true).Limit(100).Find(&users)
	fmt.Printf("简单查询 (100条): %v\n", time.Since(start))
	
	// 测试带索引的查询
	start = time.Now()
	var usersByCity []User
	db.Where("city = ?", "北京").Limit(100).Find(&usersByCity)
	fmt.Printf("索引查询 (城市): %v\n", time.Since(start))
	
	// 测试范围查询
	start = time.Now()
	var usersBySalary []User
	db.Where("salary BETWEEN ? AND ?", 8000, 15000).Limit(100).Find(&usersBySalary)
	fmt.Printf("范围查询 (薪资): %v\n", time.Since(start))
	
	// 测试复合条件查询
	start = time.Now()
	var activeUsers []User
	db.Where("is_active = ? AND department = ?", true, "技术部").Limit(100).Find(&activeUsers)
	fmt.Printf("复合条件查询: %v\n", time.Since(start))
}

// 测试2：预加载性能对比
func benchmarkPreloading(db *gorm.DB) {
	fmt.Println("\n=== 预加载性能测试 ===")
	
	// N+1 查询问题演示
	start := time.Now()
	var postsWithoutPreload []Post
	db.Limit(50).Find(&postsWithoutPreload)
	// 访问关联数据会产生N+1查询
	for _, post := range postsWithoutPreload {
		var author User
		db.First(&author, post.AuthorID)
		_ = author.Username
	}
	fmt.Printf("N+1查询 (50篇文章): %v\n", time.Since(start))
	
	// 使用预加载
	start = time.Now()
	var postsWithPreload []Post
	db.Preload("Author").Limit(50).Find(&postsWithPreload)
	for _, post := range postsWithPreload {
		_ = post.Author.Username
	}
	fmt.Printf("预加载查询 (50篇文章): %v\n", time.Since(start))
	
	// 选择性预加载
	start = time.Now()
	var postsWithSelectivePreload []Post
	db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, slug")
	}).Limit(50).Find(&postsWithSelectivePreload)
	fmt.Printf("选择性预加载 (50篇文章): %v\n", time.Since(start))
	
	// 多层预加载
	start = time.Now()
	var postsWithNestedPreload []Post
	db.Preload("Author").Preload("Category").Preload("Comments").Preload("Comments.Author").Limit(20).Find(&postsWithNestedPreload)
	fmt.Printf("多层预加载 (20篇文章): %v\n", time.Since(start))
}

// 测试3：批量操作性能
func benchmarkBatchOperations(db *gorm.DB) {
	fmt.Println("\n=== 批量操作性能测试 ===")
	
	// 单条插入 vs 批量插入
	testUsers := make([]User, 1000)
	for i := 0; i < 1000; i++ {
		testUsers[i] = User{
			Username:   fmt.Sprintf("batch_user_%d", i),
			Email:      fmt.Sprintf("batch_user_%d@test.com", i),
			FirstName:  fmt.Sprintf("First%d", i),
			LastName:   fmt.Sprintf("Last%d", i),
			Age:        25,
			City:       "测试城市",
			Country:    "测试国家",
			Salary:     10000,
			Department: "测试部门",
			Position:   "测试职位",
			JoinDate:   time.Now(),
			IsActive:   true,
		}
	}
	
	// 单条插入测试
	start := time.Now()
	for i := 0; i < 100; i++ {
		db.Create(&testUsers[i])
	}
	fmt.Printf("单条插入 (100条): %v\n", time.Since(start))
	
	// 批量插入测试
	start = time.Now()
	db.CreateInBatches(testUsers[100:200], 50)
	fmt.Printf("批量插入 (100条, 批次50): %v\n", time.Since(start))
	
	// 大批量插入测试
	start = time.Now()
	db.CreateInBatches(testUsers[200:700], 100)
	fmt.Printf("大批量插入 (500条, 批次100): %v\n", time.Since(start))
	
	// 批量更新测试
	start = time.Now()
	db.Model(&User{}).Where("username LIKE ?", "batch_user_%").Update("salary", 12000)
	fmt.Printf("批量更新: %v\n", time.Since(start))
	
	// 清理测试数据
	db.Where("username LIKE ?", "batch_user_%").Delete(&User{})
}

// 测试4：复杂查询性能
func benchmarkComplexQueries(db *gorm.DB) {
	fmt.Println("\n=== 复杂查询性能测试 ===")
	
	// 聚合查询
	start := time.Now()
	type UserStats struct {
		Department string
		UserCount  int64
		AvgSalary  float64
		MaxSalary  float64
		MinSalary  float64
	}
	var stats []UserStats
	db.Model(&User{}).Select("department, COUNT(*) as user_count, AVG(salary) as avg_salary, MAX(salary) as max_salary, MIN(salary) as min_salary").
		Where("is_active = ?", true).Group("department").Scan(&stats)
	fmt.Printf("聚合查询 (部门统计): %v\n", time.Since(start))
	
	// 子查询
	start = time.Now()
	var topAuthors []User
	subQuery := db.Model(&Post{}).Select("author_id, COUNT(*) as post_count").
		Where("status = ?", "published").Group("author_id").Order("post_count DESC").Limit(10)
	db.Table("users u").Joins("JOIN (?) pc ON u.id = pc.author_id", subQuery).
		Select("u.*").Find(&topAuthors)
	fmt.Printf("子查询 (热门作者): %v\n", time.Since(start))
	
	// 复杂连接查询
	start = time.Now()
	var results []map[string]interface{}
	db.Table("posts p").
		Select("p.title, u.username, c.name as category_name, COUNT(cm.id) as comment_count").
		Joins("JOIN users u ON p.author_id = u.id").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Joins("LEFT JOIN comments cm ON p.id = cm.post_id AND cm.status = 'approved'").
		Where("p.status = ? AND p.published_at > ?", "published", time.Now().AddDate(0, -1, 0)).
		Group("p.id, p.title, u.username, c.name").
		Having("comment_count > ?", 5).
		Order("comment_count DESC").
		Limit(20).
		Scan(&results)
	fmt.Printf("复杂连接查询 (热门文章): %v\n", time.Since(start))
	
	// 分页查询
	start = time.Now()
	var paginatedPosts []Post
	var total int64
	db.Model(&Post{}).Where("status = ?", "published").Count(&total)
	db.Where("status = ?", "published").Offset(100).Limit(20).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).Find(&paginatedPosts)
	fmt.Printf("分页查询 (第6页, 20条/页): %v\n", time.Since(start))
}

// 测试5：索引优化效果
func benchmarkIndexOptimization(db *gorm.DB) {
	fmt.Println("\n=== 索引优化效果测试 ===")
	
	// 测试有索引的查询
	start := time.Now()
	var indexedUsers []User
	db.Where("city = ? AND is_active = ?", "北京", true).Find(&indexedUsers)
	fmt.Printf("复合索引查询 (城市+状态): %v, 结果: %d条\n", time.Since(start), len(indexedUsers))
	
	// 测试范围查询
	start = time.Now()
	var salaryUsers []User
	db.Where("salary BETWEEN ? AND ? AND department = ?", 8000, 15000, "技术部").Find(&salaryUsers)
	fmt.Printf("范围+等值查询 (薪资+部门): %v, 结果: %d条\n", time.Since(start), len(salaryUsers))
	
	// 测试排序查询
	start = time.Now()
	var sortedPosts []Post
	db.Where("status = ?", "published").Order("view_count DESC, rating DESC").Limit(50).Find(&sortedPosts)
	fmt.Printf("排序查询 (浏览量+评分): %v\n", time.Since(start))
	
	// 测试模糊查询
	start = time.Now()
	var searchPosts []Post
	db.Where("title LIKE ?", "%文章%").Limit(100).Find(&searchPosts)
	fmt.Printf("模糊查询 (标题): %v, 结果: %d条\n", time.Since(start), len(searchPosts))
}

// 测试6：连接池和事务性能
func benchmarkConnectionAndTransaction(db *gorm.DB) {
	fmt.Println("\n=== 连接池和事务性能测试 ===")
	
	// 测试普通操作
	start := time.Now()
	for i := 0; i < 100; i++ {
		var user User
		db.First(&user, uint(rand.Intn(1000)+1))
	}
	fmt.Printf("普通查询 (100次): %v\n", time.Since(start))
	
	// 测试事务操作
	start = time.Now()
	for i := 0; i < 10; i++ {
		db.Transaction(func(tx *gorm.DB) error {
			for j := 0; j < 10; j++ {
				var user User
				tx.First(&user, uint(rand.Intn(1000)+1))
			}
			return nil
		})
	}
	fmt.Printf("事务查询 (10个事务, 每个10次查询): %v\n", time.Since(start))
	
	// 测试并发查询
	start = time.Now()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				var user User
				db.First(&user, uint(rand.Intn(1000)+1))
			}
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	fmt.Printf("并发查询 (10个goroutine, 每个10次查询): %v\n", time.Since(start))
}

// 性能优化建议
func showOptimizationTips() {
	fmt.Println("\n=== GORM 性能优化建议 ===")
	fmt.Println("1. 索引优化:")
	fmt.Println("   - 为经常查询的字段创建索引")
	fmt.Println("   - 使用复合索引优化多字段查询")
	fmt.Println("   - 避免在小表上创建过多索引")
	
	fmt.Println("\n2. 查询优化:")
	fmt.Println("   - 使用 Select 只查询需要的字段")
	fmt.Println("   - 合理使用 Preload 避免 N+1 查询")
	fmt.Println("   - 使用 Limit 限制查询结果数量")
	fmt.Println("   - 避免在循环中执行查询")
	
	fmt.Println("\n3. 批量操作:")
	fmt.Println("   - 使用 CreateInBatches 进行批量插入")
	fmt.Println("   - 使用批量更新替代循环更新")
	fmt.Println("   - 合理设置批次大小")
	
	fmt.Println("\n4. 连接池配置:")
	fmt.Println("   - 设置合适的最大连接数")
	fmt.Println("   - 配置连接生存时间")
	fmt.Println("   - 监控连接池使用情况")
	
	fmt.Println("\n5. 其他优化:")
	fmt.Println("   - 启用预编译语句")
	fmt.Println("   - 合理使用事务")
	fmt.Println("   - 避免不必要的关联查询")
	fmt.Println("   - 使用适当的日志级别")
}

// 主函数演示
func main() {
	fmt.Println("=== GORM Level 5 性能优化练习 ===")

	// 初始化数据库
	db := initDB()
	fmt.Println("✓ 数据库初始化完成")

	// 生成测试数据
	generateTestData(db, 2000, 5000, 10000)

	// 运行性能测试
	benchmarkBasicQueries(db)
	benchmarkPreloading(db)
	benchmarkBatchOperations(db)
	benchmarkComplexQueries(db)
	benchmarkIndexOptimization(db)
	benchmarkConnectionAndTransaction(db)

	// 显示优化建议
	showOptimizationTips()

	fmt.Println("\n=== Level 5 性能优化练习完成 ===")
}
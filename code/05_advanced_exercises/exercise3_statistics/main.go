package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 使用exercise2的模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	BaseModel
	Username    string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string     `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone       string     `gorm:"uniqueIndex;size:20" json:"phone"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	Nickname    string     `gorm:"size:50" json:"nickname"`
	Status      int8       `gorm:"default:1;comment:1-正常,2-禁用" json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

type Category struct {
	BaseModel
	Name     string `gorm:"size:50;not null" json:"name"`
	Slug     string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	ParentID *uint  `gorm:"index" json:"parent_id"`
	Status   int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
}

type Brand struct {
	BaseModel
	Name   string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug   string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Status int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
}

type Product struct {
	BaseModel
	Name       string `gorm:"size:255;not null" json:"name"`
	SKU        string `gorm:"uniqueIndex;size:100;not null" json:"sku"`
	CategoryID uint   `gorm:"index;not null" json:"category_id"`
	BrandID    *uint  `gorm:"index" json:"brand_id"`
	Price      int64  `gorm:"not null;comment:价格(分)" json:"price"`
	Stock      int    `gorm:"default:0" json:"stock"`
	Sales      int    `gorm:"default:0" json:"sales"`
	Views      int    `gorm:"default:0" json:"views"`
	Status     int8   `gorm:"default:1;comment:1-上架,2-下架" json:"status"`
}

type Order struct {
	BaseModel
	OrderNo        string     `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	Status         int8       `gorm:"index;default:1;comment:1-待付款,2-待发货,3-待收货,4-已完成,5-已取消" json:"status"`
	TotalAmount    int64      `gorm:"not null;comment:商品总金额(分)" json:"total_amount"`
	PayAmount      int64      `gorm:"not null;comment:实付金额(分)" json:"pay_amount"`
	FreightAmount  int64      `gorm:"default:0;comment:运费(分)" json:"freight_amount"`
	DiscountAmount int64      `gorm:"default:0;comment:优惠金额(分)" json:"discount_amount"`
	PaidAt         *time.Time `json:"paid_at"`
	FinishedAt     *time.Time `json:"finished_at"`
}

type OrderItem struct {
	BaseModel
	OrderID     uint   `gorm:"index;not null" json:"order_id"`
	ProductID   uint   `gorm:"index;not null" json:"product_id"`
	Quantity    int    `gorm:"not null" json:"quantity"`
	Price       int64  `gorm:"not null;comment:单价(分)" json:"price"`
	TotalPrice  int64  `gorm:"not null;comment:总价(分)" json:"total_price"`
	ProductName string `gorm:"size:255;not null" json:"product_name"`
}

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

// StatisticsService 统计服务
type StatisticsService struct {
	db *gorm.DB
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{db: db}
}

// SalesStatistics 销售统计数据
type SalesStatistics struct {
	Date          string  `json:"date"`
	OrderCount    int64   `json:"order_count"`
	SalesAmount   int64   `json:"sales_amount"`
	UserCount     int64   `json:"user_count"`
	AvgOrderValue float64 `json:"avg_order_value"`
}

// ProductSalesRank 商品销量排行
type ProductSalesRank struct {
	ProductID    uint   `json:"product_id"`
	ProductName  string `json:"product_name"`
	SalesCount   int64  `json:"sales_count"`
	SalesAmount  int64  `json:"sales_amount"`
	CategoryName string `json:"category_name"`
	BrandName    string `json:"brand_name"`
}

// UserBehaviorAnalysis 用户行为分析
type UserBehaviorAnalysis struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	OrderCount   int64     `json:"order_count"`
	TotalAmount  int64     `json:"total_amount"`
	AvgAmount    float64   `json:"avg_amount"`
	LastOrderAt  time.Time `json:"last_order_at"`
	RegisterDays int       `json:"register_days"`
}

// DashboardData 数据大屏数据
type DashboardData struct {
	TodayOrders     int64   `json:"today_orders"`
	TodaySales      int64   `json:"today_sales"`
	TodayUsers      int64   `json:"today_users"`
	TotalOrders     int64   `json:"total_orders"`
	TotalSales      int64   `json:"total_sales"`
	TotalUsers      int64   `json:"total_users"`
	TotalProducts   int64   `json:"total_products"`
	AvgOrderValue   float64 `json:"avg_order_value"`
	OrderGrowthRate float64 `json:"order_growth_rate"`
	SalesGrowthRate float64 `json:"sales_growth_rate"`
}

// GetSalesStatistics 获取销售统计数据
func (s *StatisticsService) GetSalesStatistics(startDate, endDate time.Time) ([]SalesStatistics, error) {
	var results []SalesStatistics

	sql := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as order_count,
			SUM(pay_amount) as sales_amount,
			COUNT(DISTINCT user_id) as user_count,
			AVG(pay_amount) as avg_order_value
		FROM orders 
		WHERE created_at >= ? AND created_at <= ? AND status >= 2
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	err := s.db.Raw(sql, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetProductSalesRank 获取商品销量排行
func (s *StatisticsService) GetProductSalesRank(startDate, endDate time.Time, limit int) ([]ProductSalesRank, error) {
	var results []ProductSalesRank

	sql := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			SUM(oi.quantity) as sales_count,
			SUM(oi.total_price) as sales_amount,
			c.name as category_name,
			b.name as brand_name
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN brands b ON p.brand_id = b.id
		WHERE o.created_at >= ? AND o.created_at <= ? AND o.status >= 2
		GROUP BY p.id, p.name, c.name, b.name
		ORDER BY sales_count DESC
		LIMIT ?
	`

	err := s.db.Raw(sql, startDate, endDate, limit).Scan(&results).Error
	return results, err
}

// GetUserBehaviorAnalysis 获取用户行为分析
func (s *StatisticsService) GetUserBehaviorAnalysis(startDate, endDate time.Time, limit int) ([]UserBehaviorAnalysis, error) {
	var results []UserBehaviorAnalysis

	sql := `
		SELECT 
			u.id as user_id,
			u.username,
			COUNT(o.id) as order_count,
			COALESCE(SUM(o.pay_amount), 0) as total_amount,
			COALESCE(AVG(o.pay_amount), 0) as avg_amount,
			MAX(o.created_at) as last_order_at,
			DATEDIFF(NOW(), u.created_at) as register_days
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id 
			AND o.created_at >= ? AND o.created_at <= ? 
			AND o.status >= 2
		WHERE u.created_at <= ?
		GROUP BY u.id, u.username, u.created_at
		HAVING order_count > 0
		ORDER BY total_amount DESC
		LIMIT ?
	`

	err := s.db.Raw(sql, startDate, endDate, endDate, limit).Scan(&results).Error
	return results, err
}

// GetDashboardData 获取数据大屏数据
func (s *StatisticsService) GetDashboardData() (*DashboardData, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	data := &DashboardData{}

	// 今日订单数
	s.db.Model(&Order{}).Where("created_at >= ? AND status >= 2", today).Count(&data.TodayOrders)

	// 今日销售额
	var todaySales struct{ Total int64 }
	s.db.Model(&Order{}).Select("COALESCE(SUM(pay_amount), 0) as total").
		Where("created_at >= ? AND status >= 2", today).Scan(&todaySales)
	data.TodaySales = todaySales.Total

	// 今日新增用户
	s.db.Model(&User{}).Where("created_at >= ?", today).Count(&data.TodayUsers)

	// 总订单数
	s.db.Model(&Order{}).Where("status >= 2").Count(&data.TotalOrders)

	// 总销售额
	var totalSales struct{ Total int64 }
	s.db.Model(&Order{}).Select("COALESCE(SUM(pay_amount), 0) as total").
		Where("status >= 2").Scan(&totalSales)
	data.TotalSales = totalSales.Total

	// 总用户数
	s.db.Model(&User{}).Count(&data.TotalUsers)

	// 总商品数
	s.db.Model(&Product{}).Where("status = 1").Count(&data.TotalProducts)

	// 平均订单价值
	if data.TotalOrders > 0 {
		data.AvgOrderValue = float64(data.TotalSales) / float64(data.TotalOrders)
	}

	// 计算增长率
	var yesterdayOrders int64
	s.db.Model(&Order{}).Where("created_at >= ? AND created_at < ? AND status >= 2", yesterday, today).Count(&yesterdayOrders)

	var yesterdaySales struct{ Total int64 }
	s.db.Model(&Order{}).Select("COALESCE(SUM(pay_amount), 0) as total").
		Where("created_at >= ? AND created_at < ? AND status >= 2", yesterday, today).Scan(&yesterdaySales)

	if yesterdayOrders > 0 {
		data.OrderGrowthRate = float64(data.TodayOrders-yesterdayOrders) / float64(yesterdayOrders) * 100
	}
	if yesterdaySales.Total > 0 {
		data.SalesGrowthRate = float64(data.TodaySales-yesterdaySales.Total) / float64(yesterdaySales.Total) * 100
	}

	return data, nil
}

// GetSalesStatisticsByCategory 按分类获取销售统计
func (s *StatisticsService) GetSalesStatisticsByCategory(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	sql := `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(DISTINCT o.id) as order_count,
			SUM(oi.quantity) as sales_count,
			SUM(oi.total_price) as sales_amount
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id
		LEFT JOIN order_items oi ON p.id = oi.product_id
		LEFT JOIN orders o ON oi.order_id = o.id 
			AND o.created_at >= ? AND o.created_at <= ? 
			AND o.status >= 2
		GROUP BY c.id, c.name
		ORDER BY sales_amount DESC
	`

	err := s.db.Raw(sql, startDate, endDate).Scan(&results).Error
	return results, err
}

// GetHourlyOrderStatistics 获取小时级订单统计
func (s *StatisticsService) GetHourlyOrderStatistics(date time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	sql := `
		SELECT 
			HOUR(created_at) as hour,
			COUNT(*) as order_count,
			SUM(pay_amount) as sales_amount,
			COUNT(DISTINCT user_id) as user_count
		FROM orders 
		WHERE created_at >= ? AND created_at < ? AND status >= 2
		GROUP BY HOUR(created_at)
		ORDER BY hour
	`

	err := s.db.Raw(sql, startOfDay, endOfDay).Scan(&results).Error
	return results, err
}

// GetUserRetentionAnalysis 获取用户留存分析
func (s *StatisticsService) GetUserRetentionAnalysis(startDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	sql := `
		SELECT 
			DATE(u.created_at) as register_date,
			COUNT(u.id) as register_count,
			COUNT(CASE WHEN o1.user_id IS NOT NULL THEN 1 END) as day1_retention,
			COUNT(CASE WHEN o7.user_id IS NOT NULL THEN 1 END) as day7_retention,
			COUNT(CASE WHEN o30.user_id IS NOT NULL THEN 1 END) as day30_retention
		FROM users u
		LEFT JOIN (
			SELECT DISTINCT user_id 
			FROM orders 
			WHERE created_at >= DATE_ADD(?, INTERVAL 1 DAY) 
				AND created_at < DATE_ADD(?, INTERVAL 2 DAY)
				AND status >= 2
		) o1 ON u.id = o1.user_id
		LEFT JOIN (
			SELECT DISTINCT user_id 
			FROM orders 
			WHERE created_at >= DATE_ADD(?, INTERVAL 7 DAY) 
				AND created_at < DATE_ADD(?, INTERVAL 8 DAY)
				AND status >= 2
		) o7 ON u.id = o7.user_id
		LEFT JOIN (
			SELECT DISTINCT user_id 
			FROM orders 
			WHERE created_at >= DATE_ADD(?, INTERVAL 30 DAY) 
				AND created_at < DATE_ADD(?, INTERVAL 31 DAY)
				AND status >= 2
		) o30 ON u.id = o30.user_id
		WHERE u.created_at >= ? AND u.created_at < DATE_ADD(?, INTERVAL 1 DAY)
		GROUP BY DATE(u.created_at)
		ORDER BY register_date
	`

	err := s.db.Raw(sql, startDate, startDate, startDate, startDate, startDate, startDate, startDate, startDate).Scan(&results).Error
	return results, err
}

// GetCohortAnalysis 获取队列分析
func (s *StatisticsService) GetCohortAnalysis(startDate time.Time, months int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 队列分析：按注册月份分组，分析每个月份用户在后续月份的购买行为
	sql := `
		SELECT 
			DATE_FORMAT(u.created_at, '%Y-%m') as cohort_month,
			COUNT(DISTINCT u.id) as total_users,
			COUNT(DISTINCT CASE WHEN PERIOD_DIFF(DATE_FORMAT(o.created_at, '%Y%m'), DATE_FORMAT(u.created_at, '%Y%m')) = 0 THEN u.id END) as month_0,
			COUNT(DISTINCT CASE WHEN PERIOD_DIFF(DATE_FORMAT(o.created_at, '%Y%m'), DATE_FORMAT(u.created_at, '%Y%m')) = 1 THEN u.id END) as month_1,
			COUNT(DISTINCT CASE WHEN PERIOD_DIFF(DATE_FORMAT(o.created_at, '%Y%m'), DATE_FORMAT(u.created_at, '%Y%m')) = 2 THEN u.id END) as month_2,
			COUNT(DISTINCT CASE WHEN PERIOD_DIFF(DATE_FORMAT(o.created_at, '%Y%m'), DATE_FORMAT(u.created_at, '%Y%m')) = 3 THEN u.id END) as month_3
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id AND o.status >= 2
		WHERE u.created_at >= ?
		GROUP BY DATE_FORMAT(u.created_at, '%Y-%m')
		ORDER BY cohort_month
	`

	err := s.db.Raw(sql, startDate).Scan(&results).Error
	return results, err
}

// GetRFMAnalysis 获取RFM分析（最近购买时间、购买频率、购买金额）
func (s *StatisticsService) GetRFMAnalysis() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	sql := `
		SELECT 
			u.id as user_id,
			u.username,
			DATEDIFF(NOW(), MAX(o.created_at)) as recency,
			COUNT(o.id) as frequency,
			SUM(o.pay_amount) as monetary,
			CASE 
				WHEN DATEDIFF(NOW(), MAX(o.created_at)) <= 30 THEN 5
				WHEN DATEDIFF(NOW(), MAX(o.created_at)) <= 60 THEN 4
				WHEN DATEDIFF(NOW(), MAX(o.created_at)) <= 90 THEN 3
				WHEN DATEDIFF(NOW(), MAX(o.created_at)) <= 180 THEN 2
				ELSE 1
			END as r_score,
			CASE 
				WHEN COUNT(o.id) >= 10 THEN 5
				WHEN COUNT(o.id) >= 5 THEN 4
				WHEN COUNT(o.id) >= 3 THEN 3
				WHEN COUNT(o.id) >= 2 THEN 2
				ELSE 1
			END as f_score,
			CASE 
				WHEN SUM(o.pay_amount) >= 100000 THEN 5
				WHEN SUM(o.pay_amount) >= 50000 THEN 4
				WHEN SUM(o.pay_amount) >= 20000 THEN 3
				WHEN SUM(o.pay_amount) >= 10000 THEN 2
				ELSE 1
			END as m_score
		FROM users u
		JOIN orders o ON u.id = o.user_id AND o.status >= 2
		GROUP BY u.id, u.username
		ORDER BY monetary DESC
	`

	err := s.db.Raw(sql).Scan(&results).Error
	return results, err
}

// SeedTestData 填充测试数据
func SeedTestData(db *gorm.DB) error {
	fmt.Println("开始填充测试数据...")

	// 创建用户
	users := []User{
		{Username: "user1", Email: "user1@example.com", Phone: "13800138001", Password: "password", Nickname: "用户1"},
		{Username: "user2", Email: "user2@example.com", Phone: "13800138002", Password: "password", Nickname: "用户2"},
		{Username: "user3", Email: "user3@example.com", Phone: "13800138003", Password: "password", Nickname: "用户3"},
	}
	db.Create(&users)

	// 创建分类
	categories := []Category{
		{Name: "电子产品", Slug: "electronics"},
		{Name: "服装", Slug: "clothing"},
		{Name: "图书", Slug: "books"},
	}
	db.Create(&categories)

	// 创建品牌
	brands := []Brand{
		{Name: "苹果", Slug: "apple"},
		{Name: "华为", Slug: "huawei"},
		{Name: "小米", Slug: "xiaomi"},
	}
	db.Create(&brands)

	// 创建商品
	products := []Product{
		{Name: "iPhone 15", SKU: "IPHONE15", CategoryID: categories[0].ID, BrandID: &brands[0].ID, Price: 599900, Stock: 100, Sales: 50},
		{Name: "华为P60", SKU: "HUAWEIP60", CategoryID: categories[0].ID, BrandID: &brands[1].ID, Price: 499900, Stock: 80, Sales: 30},
		{Name: "小米14", SKU: "XIAOMI14", CategoryID: categories[0].ID, BrandID: &brands[2].ID, Price: 399900, Stock: 120, Sales: 80},
		{Name: "T恤", SKU: "TSHIRT001", CategoryID: categories[1].ID, Price: 9900, Stock: 200, Sales: 150},
		{Name: "编程书籍", SKU: "BOOK001", CategoryID: categories[2].ID, Price: 5900, Stock: 50, Sales: 25},
	}
	db.Create(&products)

	// 创建订单和订单项
	for i := 0; i < 30; i++ {
		userID := uint(i%3 + 1)
		productID := uint(i%5 + 1)
		quantity := i%3 + 1
		price := products[productID-1].Price
		totalPrice := price * int64(quantity)

		order := Order{
			OrderNo:     fmt.Sprintf("ORD%d", time.Now().UnixNano()+int64(i)),
			UserID:      userID,
			Status:      4, // 已完成
			TotalAmount: totalPrice,
			PayAmount:   totalPrice,
			CreatedAt:   time.Now().AddDate(0, 0, -i), // 不同日期
		}
		db.Create(&order)

		orderItem := OrderItem{
			OrderID:     order.ID,
			ProductID:   productID,
			Quantity:    quantity,
			Price:       price,
			TotalPrice:  totalPrice,
			ProductName: products[productID-1].Name,
		}
		db.Create(&orderItem)
	}

	fmt.Println("测试数据填充完成")
	return nil
}

// demonstrateStatistics 演示统计功能
func demonstrateStatistics(db *gorm.DB) {
	fmt.Println("\n=== 演示统计功能 ===")

	statisticsService := NewStatisticsService(db)

	// 1. 销售统计
	fmt.Println("\n1. 销售统计:")
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()
	salesStats, err := statisticsService.GetSalesStatistics(startDate, endDate)
	if err != nil {
		fmt.Printf("获取销售统计失败: %v\n", err)
	} else {
		for _, stat := range salesStats {
			fmt.Printf("日期: %s, 订单数: %d, 销售额: %.2f元, 用户数: %d, 平均订单价值: %.2f元\n",
				stat.Date, stat.OrderCount, float64(stat.SalesAmount)/100, stat.UserCount, stat.AvgOrderValue/100)
		}
	}

	// 2. 商品销量排行
	fmt.Println("\n2. 商品销量排行:")
	productRank, err := statisticsService.GetProductSalesRank(startDate, endDate, 10)
	if err != nil {
		fmt.Printf("获取商品销量排行失败: %v\n", err)
	} else {
		for i, rank := range productRank {
			fmt.Printf("排名%d: %s (分类: %s, 品牌: %s), 销量: %d, 销售额: %.2f元\n",
				i+1, rank.ProductName, rank.CategoryName, rank.BrandName, rank.SalesCount, float64(rank.SalesAmount)/100)
		}
	}

	// 3. 用户行为分析
	fmt.Println("\n3. 用户行为分析:")
	userBehavior, err := statisticsService.GetUserBehaviorAnalysis(startDate, endDate, 10)
	if err != nil {
		fmt.Printf("获取用户行为分析失败: %v\n", err)
	} else {
		for _, behavior := range userBehavior {
			fmt.Printf("用户: %s, 订单数: %d, 总金额: %.2f元, 平均金额: %.2f元, 注册天数: %d\n",
				behavior.Username, behavior.OrderCount, float64(behavior.TotalAmount)/100, behavior.AvgAmount/100, behavior.RegisterDays)
		}
	}

	// 4. 数据大屏
	fmt.Println("\n4. 数据大屏:")
	dashboard, err := statisticsService.GetDashboardData()
	if err != nil {
		fmt.Printf("获取数据大屏数据失败: %v\n", err)
	} else {
		fmt.Printf("今日订单: %d, 今日销售额: %.2f元, 今日新增用户: %d\n",
			dashboard.TodayOrders, float64(dashboard.TodaySales)/100, dashboard.TodayUsers)
		fmt.Printf("总订单: %d, 总销售额: %.2f元, 总用户: %d, 总商品: %d\n",
			dashboard.TotalOrders, float64(dashboard.TotalSales)/100, dashboard.TotalUsers, dashboard.TotalProducts)
		fmt.Printf("平均订单价值: %.2f元, 订单增长率: %.2f%%, 销售额增长率: %.2f%%\n",
			dashboard.AvgOrderValue/100, dashboard.OrderGrowthRate, dashboard.SalesGrowthRate)
	}

	// 5. 按分类统计
	fmt.Println("\n5. 按分类销售统计:")
	categoryStats, err := statisticsService.GetSalesStatisticsByCategory(startDate, endDate)
	if err != nil {
		fmt.Printf("获取分类统计失败: %v\n", err)
	} else {
		for _, stat := range categoryStats {
			fmt.Printf("分类: %v, 订单数: %v, 销量: %v, 销售额: %.2f元\n",
				stat["category_name"], stat["order_count"], stat["sales_count"], 
				float64(stat["sales_amount"].(int64))/100)
		}
	}

	// 6. 小时级统计
	fmt.Println("\n6. 今日小时级订单统计:")
	hourlyStats, err := statisticsService.GetHourlyOrderStatistics(time.Now())
	if err != nil {
		fmt.Printf("获取小时级统计失败: %v\n", err)
	} else {
		for _, stat := range hourlyStats {
			fmt.Printf("%v点: 订单数 %v, 销售额 %.2f元, 用户数 %v\n",
				stat["hour"], stat["order_count"], 
				float64(stat["sales_amount"].(int64))/100, stat["user_count"])
		}
	}

	// 7. RFM分析
	fmt.Println("\n7. RFM分析:")
	rfmAnalysis, err := statisticsService.GetRFMAnalysis()
	if err != nil {
		fmt.Printf("获取RFM分析失败: %v\n", err)
	} else {
		for _, rfm := range rfmAnalysis {
			fmt.Printf("用户: %v, 最近购买: %v天前, 购买频率: %v次, 购买金额: %.2f元, RFM评分: %v-%v-%v\n",
				rfm["username"], rfm["recency"], rfm["frequency"], 
				float64(rfm["monetary"].(int64))/100, rfm["r_score"], rfm["f_score"], rfm["m_score"])
		}
	}
}

func main() {
	// 数据库配置
	config := DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "gorm_advanced_exercise3",
		Charset:  "utf8mb4",
	}

	// 连接数据库
	fmt.Println("连接数据库...")
	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 迁移数据库
	db.AutoMigrate(&User{}, &Category{}, &Brand{}, &Product{}, &Order{}, &OrderItem{})

	// 检查是否需要填充测试数据
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		if err := SeedTestData(db); err != nil {
			log.Fatal("填充测试数据失败:", err)
		}
	}

	// 演示统计功能
	demonstrateStatistics(db)

	fmt.Println("\n=== 练习3：数据统计和报表 演示完成 ===")
	fmt.Println("\n强化练习任务:")
	fmt.Println("1. 复杂查询优化（索引优化、查询重写）")
	fmt.Println("2. 数据可视化（图表生成、导出功能）")
	fmt.Println("3. 实时更新（WebSocket推送、缓存更新）")
	fmt.Println("4. 缓存优化（Redis缓存、查询结果缓存）")
	fmt.Println("5. 导出功能（Excel、PDF、CSV格式）")
}
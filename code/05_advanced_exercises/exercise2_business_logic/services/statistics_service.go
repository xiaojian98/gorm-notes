package services

import (
	"time"

	"gorm.io/gorm"
)

// SalesStatistics 销售统计数据
type SalesStatistics struct {
	Date         string  `json:"date"`
	OrderCount   int64   `json:"order_count"`
	SalesAmount  int64   `json:"sales_amount"`
	UserCount    int64   `json:"user_count"`
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

// StatisticsService 统计服务
type StatisticsService struct {
	db *gorm.DB
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{
		db: db,
	}
}

// GetSalesStatistics 获取销售统计数据
func (s *StatisticsService) GetSalesStatistics(startDate, endDate time.Time) ([]SalesStatistics, error) {
	var results []SalesStatistics

	// 使用原生SQL进行复杂统计查询
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
	if err != nil {
		return nil, err
	}

	return results, nil
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
	if err != nil {
		return nil, err
	}

	return results, nil
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
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetDashboardData 获取数据大屏数据
func (s *StatisticsService) GetDashboardData() (*DashboardData, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	data := &DashboardData{}

	// 今日订单数
	err := s.db.Model(&Order{}).Where("created_at >= ? AND status >= 2", today).Count(&data.TodayOrders).Error
	if err != nil {
		return nil, err
	}

	// 今日销售额
	var todaySales struct {
		Total int64
	}
	err = s.db.Model(&Order{}).Select("COALESCE(SUM(pay_amount), 0) as total").
		Where("created_at >= ? AND status >= 2", today).Scan(&todaySales).Error
	if err != nil {
		return nil, err
	}
	data.TodaySales = todaySales.Total

	// 今日新增用户
	err = s.db.Model(&User{}).Where("created_at >= ?", today).Count(&data.TodayUsers).Error
	if err != nil {
		return nil, err
	}

	// 总订单数
	err = s.db.Model(&Order{}).Where("status >= 2").Count(&data.TotalOrders).Error
	if err != nil {
		return nil, err
	}

	// 总销售额
	var totalSales struct {
		Total int64
	}
	err = s.db.Model(&Order{}).Select("COALESCE(SUM(pay_amount), 0) as total").
		Where("status >= 2").Scan(&totalSales).Error
	if err != nil {
		return nil, err
	}
	data.TotalSales = totalSales.Total

	// 总用户数
	err = s.db.Model(&User{}).Count(&data.TotalUsers).Error
	if err != nil {
		return nil, err
	}

	// 总商品数
	err = s.db.Model(&Product{}).Where("status = 1").Count(&data.TotalProducts).Error
	if err != nil {
		return nil, err
	}

	// 平均订单价值
	if data.TotalOrders > 0 {
		data.AvgOrderValue = float64(data.TotalSales) / float64(data.TotalOrders)
	}

	// 计算增长率
	// 昨日订单数
	var yesterdayOrders int64
	err = s.db.Model(&Order{}).Where("created_at >= ? AND created_at < ? AND status >= 2", yesterday, today).Count(&yesterdayOrders).Error
	if err != nil {
		return nil, err
	}

	// 昨日销售额
	var yesterdaySales struct {
		Total int64
	}
	err = s.db.Model(&Order{}).Select("COALESCE(SUM(pay_amount), 0) as total").
		Where("created_at >= ? AND created_at < ? AND status >= 2", yesterday, today).Scan(&yesterdaySales).Error
	if err != nil {
		return nil, err
	}

	// 计算增长率
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
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetSalesStatisticsByBrand 按品牌获取销售统计
func (s *StatisticsService) GetSalesStatisticsByBrand(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	sql := `
		SELECT 
			b.id as brand_id,
			b.name as brand_name,
			COUNT(DISTINCT o.id) as order_count,
			SUM(oi.quantity) as sales_count,
			SUM(oi.total_price) as sales_amount
		FROM brands b
		LEFT JOIN products p ON b.id = p.brand_id
		LEFT JOIN order_items oi ON p.id = oi.product_id
		LEFT JOIN orders o ON oi.order_id = o.id 
			AND o.created_at >= ? AND o.created_at <= ? 
			AND o.status >= 2
		GROUP BY b.id, b.name
		ORDER BY sales_amount DESC
	`

	err := s.db.Raw(sql, startDate, endDate).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetUserRetentionAnalysis 获取用户留存分析
func (s *StatisticsService) GetUserRetentionAnalysis(startDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 计算用户留存率（注册后1天、7天、30天的留存情况）
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
	if err != nil {
		return nil, err
	}

	return results, nil
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
	if err != nil {
		return nil, err
	}

	return results, nil
}
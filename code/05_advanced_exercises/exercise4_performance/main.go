package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
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
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// ConnectDatabase 连接数据库（优化版）
func ConnectDatabase(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用外键约束检查以提高性能
		DisableForeignKeyConstraintWhenMigrating: true,
		// 预编译语句缓存
		PrepareStmt: true,
	})

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取sql.DB失败: %w", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)       // 最大空闲连接数
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)       // 最大打开连接数
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime) // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime) // 连接最大空闲时间

	return db, nil
}

// PerformanceMonitor 性能监控器
type PerformanceMonitor struct {
	db        *gorm.DB
	queryLogs []QueryLog
	mu        sync.RWMutex
}

// QueryLog 查询日志
type QueryLog struct {
	SQL      string        `json:"sql"`
	Duration time.Duration `json:"duration"`
	Rows     int64         `json:"rows"`
	Time     time.Time     `json:"time"`
}

// NewPerformanceMonitor 创建性能监控器
func NewPerformanceMonitor(db *gorm.DB) *PerformanceMonitor {
	return &PerformanceMonitor{
		db:        db,
		queryLogs: make([]QueryLog, 0),
	}
}

// LogQuery 记录查询
func (pm *PerformanceMonitor) LogQuery(sql string, duration time.Duration, rows int64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.queryLogs = append(pm.queryLogs, QueryLog{
		SQL:      sql,
		Duration: duration,
		Rows:     rows,
		Time:     time.Now(),
	})

	// 保持最近1000条记录
	if len(pm.queryLogs) > 1000 {
		pm.queryLogs = pm.queryLogs[len(pm.queryLogs)-1000:]
	}
}

// GetSlowQueries 获取慢查询
func (pm *PerformanceMonitor) GetSlowQueries(threshold time.Duration) []QueryLog {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var slowQueries []QueryLog
	for _, log := range pm.queryLogs {
		if log.Duration > threshold {
			slowQueries = append(slowQueries, log)
		}
	}
	return slowQueries
}

// GetQueryStats 获取查询统计
func (pm *PerformanceMonitor) GetQueryStats() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if len(pm.queryLogs) == 0 {
		return map[string]interface{}{
			"total_queries": 0,
			"avg_duration":  0,
			"max_duration":  0,
			"min_duration":  0,
		}
	}

	var totalDuration time.Duration
	maxDuration := pm.queryLogs[0].Duration
	minDuration := pm.queryLogs[0].Duration

	for _, log := range pm.queryLogs {
		totalDuration += log.Duration
		if log.Duration > maxDuration {
			maxDuration = log.Duration
		}
		if log.Duration < minDuration {
			minDuration = log.Duration
		}
	}

	return map[string]interface{}{
		"total_queries": len(pm.queryLogs),
		"avg_duration":  totalDuration / time.Duration(len(pm.queryLogs)),
		"max_duration":  maxDuration,
		"min_duration":  minDuration,
	}
}

// OptimizedQueryService 优化查询服务
type OptimizedQueryService struct {
	db      *gorm.DB
	monitor *PerformanceMonitor
}

// NewOptimizedQueryService 创建优化查询服务
func NewOptimizedQueryService(db *gorm.DB, monitor *PerformanceMonitor) *OptimizedQueryService {
	return &OptimizedQueryService{
		db:      db,
		monitor: monitor,
	}
}

// GetProductsWithPagination 分页查询商品（优化版）
func (s *OptimizedQueryService) GetProductsWithPagination(page, pageSize int, categoryID *uint) ([]Product, int64, error) {
	start := time.Now()
	defer func() {
		s.monitor.LogQuery("GetProductsWithPagination", time.Since(start), 0)
	}()

	var products []Product
	var total int64

	query := s.db.Model(&Product{}).Where("status = ?", 1)
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// 先获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，使用索引优化
	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&products).Error

	return products, total, err
}

// GetOrdersWithJoin 关联查询订单（优化版）
func (s *OptimizedQueryService) GetOrdersWithJoin(userID uint, limit int) ([]map[string]interface{}, error) {
	start := time.Now()
	defer func() {
		s.monitor.LogQuery("GetOrdersWithJoin", time.Since(start), 0)
	}()

	var results []map[string]interface{}

	// 使用原生SQL进行优化查询
	sql := `
		SELECT 
			o.id,
			o.order_no,
			o.status,
			o.pay_amount,
			o.created_at,
			u.username,
			COUNT(oi.id) as item_count,
			SUM(oi.quantity) as total_quantity
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = ? AND o.deleted_at IS NULL
		GROUP BY o.id, o.order_no, o.status, o.pay_amount, o.created_at, u.username
		ORDER BY o.created_at DESC
		LIMIT ?
	`

	err := s.db.Raw(sql, userID, limit).Scan(&results).Error
	return results, err
}

// GetSalesStatisticsOptimized 优化的销售统计
func (s *OptimizedQueryService) GetSalesStatisticsOptimized(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	start := time.Now()
	defer func() {
		s.monitor.LogQuery("GetSalesStatisticsOptimized", time.Since(start), 0)
	}()

	var results []map[string]interface{}

	// 使用索引优化的查询
	sql := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as order_count,
			SUM(pay_amount) as sales_amount,
			COUNT(DISTINCT user_id) as user_count,
			AVG(pay_amount) as avg_order_value
		FROM orders 
		WHERE created_at >= ? AND created_at <= ? 
			AND status >= 2 
			AND deleted_at IS NULL
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	err := s.db.Raw(sql, startDate, endDate).Scan(&results).Error
	return results, err
}

// BatchInsertProducts 批量插入商品
func (s *OptimizedQueryService) BatchInsertProducts(products []Product, batchSize int) error {
	start := time.Now()
	defer func() {
		s.monitor.LogQuery("BatchInsertProducts", time.Since(start), int64(len(products)))
	}()

	// 分批插入
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		batch := products[i:end]
		if err := s.db.CreateInBatches(batch, batchSize).Error; err != nil {
			return err
		}
	}

	return nil
}

// UpdateProductStockOptimized 优化的库存更新
func (s *OptimizedQueryService) UpdateProductStockOptimized(productID uint, quantity int) error {
	start := time.Now()
	defer func() {
		s.monitor.LogQuery("UpdateProductStockOptimized", time.Since(start), 1)
	}()

	// 使用原子操作更新库存
	result := s.db.Model(&Product{}).Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("库存不足或商品不存在")
	}

	return nil
}

// GetConnectionStats 获取连接池统计
func GetConnectionStats(db *gorm.DB) (map[string]interface{}, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections":     stats.MaxOpenConnections,
		"open_connections":         stats.OpenConnections,
		"in_use":                   stats.InUse,
		"idle":                     stats.Idle,
		"wait_count":               stats.WaitCount,
		"wait_duration":            stats.WaitDuration,
		"max_idle_closed":          stats.MaxIdleClosed,
		"max_idle_time_closed":     stats.MaxIdleTimeClosed,
		"max_lifetime_closed":      stats.MaxLifetimeClosed,
	}, nil
}

// CreateOptimizedIndexes 创建优化索引
func CreateOptimizedIndexes(db *gorm.DB) error {
	fmt.Println("创建优化索引...")

	// 复合索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_orders_user_status_created ON orders(user_id, status, created_at)",
		"CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders(status, created_at)",
		"CREATE INDEX IF NOT EXISTS idx_products_category_status ON products(category_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_products_brand_status ON products(brand_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_order_items_order_product ON order_items(order_id, product_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_status_created ON users(status, created_at)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			fmt.Printf("创建索引失败: %s, 错误: %v\n", indexSQL, err)
		} else {
			fmt.Printf("索引创建成功: %s\n", indexSQL)
		}
	}

	return nil
}

// BenchmarkTest 性能基准测试
type BenchmarkTest struct {
	db      *gorm.DB
	monitor *PerformanceMonitor
}

// NewBenchmarkTest 创建基准测试
func NewBenchmarkTest(db *gorm.DB, monitor *PerformanceMonitor) *BenchmarkTest {
	return &BenchmarkTest{
		db:      db,
		monitor: monitor,
	}
}

// RunConcurrentQueries 并发查询测试
func (bt *BenchmarkTest) RunConcurrentQueries(concurrency int, iterations int) {
	fmt.Printf("\n开始并发查询测试: %d个并发, 每个执行%d次查询\n", concurrency, iterations)

	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				// 查询商品
				var products []Product
				bt.db.Where("status = ?", 1).Limit(10).Find(&products)

				// 查询订单
				var orders []Order
				bt.db.Where("status >= ?", 2).Limit(5).Find(&orders)

				// 统计查询
				var count int64
				bt.db.Model(&User{}).Where("status = ?", 1).Count(&count)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	totalQueries := concurrency * iterations * 3 // 每次循环3个查询
	qps := float64(totalQueries) / duration.Seconds()

	fmt.Printf("并发测试完成: 总耗时 %v, 总查询数 %d, QPS: %.2f\n", duration, totalQueries, qps)
}

// RunBatchInsertTest 批量插入测试
func (bt *BenchmarkTest) RunBatchInsertTest(totalRecords int, batchSize int) {
	fmt.Printf("\n开始批量插入测试: 总记录数 %d, 批次大小 %d\n", totalRecords, batchSize)

	// 生成测试数据
	products := make([]Product, totalRecords)
	for i := 0; i < totalRecords; i++ {
		products[i] = Product{
			Name:       fmt.Sprintf("测试商品%d", i),
			SKU:        fmt.Sprintf("TEST%d", i),
			CategoryID: 1,
			Price:      int64(1000 + i),
			Stock:      100,
			Status:     1,
		}
	}

	start := time.Now()

	// 批量插入
	service := NewOptimizedQueryService(bt.db, bt.monitor)
	err := service.BatchInsertProducts(products, batchSize)

	duration := time.Since(start)

	if err != nil {
		fmt.Printf("批量插入失败: %v\n", err)
	} else {
		rps := float64(totalRecords) / duration.Seconds()
		fmt.Printf("批量插入完成: 总耗时 %v, 插入速度: %.2f records/s\n", duration, rps)
	}

	// 清理测试数据
	bt.db.Where("sku LIKE 'TEST%'").Delete(&Product{})
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

	// 创建大量订单用于性能测试
	for i := 0; i < 1000; i++ {
		userID := uint(i%3 + 1)
		productID := uint(i%5 + 1)
		quantity := i%3 + 1
		price := products[productID-1].Price
		totalPrice := price * int64(quantity)

		order := Order{
			OrderNo:     fmt.Sprintf("ORD%d", time.Now().UnixNano()+int64(i)),
			UserID:      userID,
			Status:      int8(i%5 + 1), // 随机状态
			TotalAmount: totalPrice,
			PayAmount:   totalPrice,
			CreatedAt:   time.Now().AddDate(0, 0, -i%365), // 随机日期
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

// demonstratePerformanceOptimization 演示性能优化
func demonstratePerformanceOptimization(db *gorm.DB) {
	fmt.Println("\n=== 演示性能优化功能 ===")

	// 创建性能监控器
	monitor := NewPerformanceMonitor(db)
	service := NewOptimizedQueryService(db, monitor)

	// 1. 分页查询测试
	fmt.Println("\n1. 分页查询测试:")
	categoryID := uint(1)
	products, total, err := service.GetProductsWithPagination(1, 10, &categoryID)
	if err != nil {
		fmt.Printf("分页查询失败: %v\n", err)
	} else {
		fmt.Printf("查询到 %d 个商品，总数: %d\n", len(products), total)
	}

	// 2. 关联查询测试
	fmt.Println("\n2. 关联查询测试:")
	orders, err := service.GetOrdersWithJoin(1, 5)
	if err != nil {
		fmt.Printf("关联查询失败: %v\n", err)
	} else {
		fmt.Printf("查询到 %d 个订单\n", len(orders))
		for _, order := range orders {
			fmt.Printf("订单号: %v, 用户: %v, 商品数: %v\n", 
				order["order_no"], order["username"], order["item_count"])
		}
	}

	// 3. 销售统计查询测试
	fmt.Println("\n3. 销售统计查询测试:")
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()
	stats, err := service.GetSalesStatisticsOptimized(startDate, endDate)
	if err != nil {
		fmt.Printf("统计查询失败: %v\n", err)
	} else {
		fmt.Printf("统计数据: %d 条记录\n", len(stats))
	}

	// 4. 库存更新测试
	fmt.Println("\n4. 库存更新测试:")
	err = service.UpdateProductStockOptimized(1, 1)
	if err != nil {
		fmt.Printf("库存更新失败: %v\n", err)
	} else {
		fmt.Println("库存更新成功")
	}

	// 5. 连接池统计
	fmt.Println("\n5. 连接池统计:")
	connStats, err := GetConnectionStats(db)
	if err != nil {
		fmt.Printf("获取连接池统计失败: %v\n", err)
	} else {
		fmt.Printf("最大连接数: %v, 当前连接数: %v, 使用中: %v, 空闲: %v\n",
			connStats["max_open_connections"], connStats["open_connections"],
			connStats["in_use"], connStats["idle"])
	}

	// 6. 查询性能统计
	fmt.Println("\n6. 查询性能统计:")
	queryStats := monitor.GetQueryStats()
	fmt.Printf("总查询数: %v, 平均耗时: %v, 最大耗时: %v, 最小耗时: %v\n",
		queryStats["total_queries"], queryStats["avg_duration"],
		queryStats["max_duration"], queryStats["min_duration"])

	// 7. 慢查询分析
	slowQueries := monitor.GetSlowQueries(100 * time.Millisecond)
	if len(slowQueries) > 0 {
		fmt.Printf("\n发现 %d 个慢查询:\n", len(slowQueries))
		for _, query := range slowQueries {
			fmt.Printf("SQL: %s, 耗时: %v\n", query.SQL, query.Duration)
		}
	} else {
		fmt.Println("\n未发现慢查询")
	}

	// 8. 基准测试
	benchmark := NewBenchmarkTest(db, monitor)
	benchmark.RunConcurrentQueries(10, 100)
	benchmark.RunBatchInsertTest(1000, 100)
}

func main() {
	// 数据库配置（优化版）
	config := DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		User:            "root",
		Password:        "123456",
		DBName:          "gorm_advanced_exercise4",
		Charset:         "utf8mb4",
		MaxIdleConns:    10,                // 最大空闲连接数
		MaxOpenConns:    100,               // 最大打开连接数
		ConnMaxLifetime: time.Hour,         // 连接最大生存时间
		ConnMaxIdleTime: 10 * time.Minute,  // 连接最大空闲时间
	}

	// 连接数据库
	fmt.Println("连接数据库...")
	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 迁移数据库
	db.AutoMigrate(&User{}, &Category{}, &Brand{}, &Product{}, &Order{}, &OrderItem{})

	// 创建优化索引
	CreateOptimizedIndexes(db)

	// 检查是否需要填充测试数据
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		if err := SeedTestData(db); err != nil {
			log.Fatal("填充测试数据失败:", err)
		}
	}

	// 演示性能优化功能
	demonstratePerformanceOptimization(db)

	fmt.Println("\n=== 练习4：性能优化和监控 演示完成 ===")
	fmt.Println("\n强化练习任务:")
	fmt.Println("1. 查询计划分析（EXPLAIN优化）")
	fmt.Println("2. 缓存策略（Redis集成、查询缓存）")
	fmt.Println("3. 读写分离（主从配置、负载均衡）")
	fmt.Println("4. 分库分表（水平分片、垂直分片）")
	fmt.Println("5. 监控告警（Prometheus集成、性能指标）")
}
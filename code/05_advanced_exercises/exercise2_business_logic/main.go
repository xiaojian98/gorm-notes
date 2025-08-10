package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm-advanced-exercises/exercise2_business_logic/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
		// 禁用外键约束检查（开发环境）
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接的最大生存时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// MigrateDatabase 迁移数据库
func MigrateDatabase(db *gorm.DB) error {
	// 自动迁移所有模型
	err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&Address{},
		&Category{},
		&Brand{},
		&Product{},
		&ProductImage{},
		&ProductSKU{},
		&ProductReview{},
		&Cart{},
		&Order{},
		&OrderItem{},
		&Payment{},
		&Coupon{},
		&UserCoupon{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	fmt.Println("数据库迁移完成")
	return nil
}

// SeedTestData 填充测试数据
func SeedTestData(db *gorm.DB) error {
	fmt.Println("开始填充测试数据...")

	// 创建用户
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Password: "password123",
		Nickname: "测试用户",
		Status:   1,
	}
	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 创建用户资料
	profile := &UserProfile{
		UserID:   user.ID,
		RealName: "张三",
		Company:  "测试公司",
		Position: "软件工程师",
	}
	if err := db.Create(profile).Error; err != nil {
		return fmt.Errorf("创建用户资料失败: %w", err)
	}

	// 创建收货地址
	address := &Address{
		UserID:    user.ID,
		Name:      "张三",
		Phone:     "13800138000",
		Province:  "北京市",
		City:      "北京市",
		District:  "朝阳区",
		Detail:    "某某街道某某小区",
		IsDefault: true,
	}
	if err := db.Create(address).Error; err != nil {
		return fmt.Errorf("创建收货地址失败: %w", err)
	}

	// 创建分类
	category := &Category{
		Name:        "电子产品",
		Slug:        "electronics",
		Description: "各种电子产品",
		Status:      1,
	}
	if err := db.Create(category).Error; err != nil {
		return fmt.Errorf("创建分类失败: %w", err)
	}

	// 创建品牌
	brand := &Brand{
		Name:        "苹果",
		Slug:        "apple",
		Description: "苹果公司",
		Website:     "https://www.apple.com",
		Status:      1,
	}
	if err := db.Create(brand).Error; err != nil {
		return fmt.Errorf("创建品牌失败: %w", err)
	}

	// 创建商品
	product := &Product{
		Name:        "iPhone 15 Pro",
		SKU:         "IPHONE15PRO",
		Description: "最新款iPhone",
		Content:     "详细的商品描述...",
		CategoryID:  category.ID,
		BrandID:     &brand.ID,
		Price:       899900, // 8999.00元
		MarketPrice: 999900, // 9999.00元
		Stock:       100,
		Status:      1,
	}
	if err := db.Create(product).Error; err != nil {
		return fmt.Errorf("创建商品失败: %w", err)
	}

	// 创建商品SKU
	sku1 := &ProductSKU{
		ProductID: product.ID,
		SKU:       "IPHONE15PRO-128GB-BLACK",
		Name:      "iPhone 15 Pro 128GB 深空黑色",
		Price:     899900,
		Stock:     50,
		Specs:     json.RawMessage(`{"storage":"128GB","color":"深空黑色"}`),
		Status:    1,
	}
	sku2 := &ProductSKU{
		ProductID: product.ID,
		SKU:       "IPHONE15PRO-256GB-BLACK",
		Name:      "iPhone 15 Pro 256GB 深空黑色",
		Price:     999900,
		Stock:     30,
		Specs:     json.RawMessage(`{"storage":"256GB","color":"深空黑色"}`),
		Status:    1,
	}
	if err := db.Create([]*ProductSKU{sku1, sku2}).Error; err != nil {
		return fmt.Errorf("创建商品SKU失败: %w", err)
	}

	// 创建优惠券
	coupon := &Coupon{
		Name:          "新用户专享",
		Code:          "NEWUSER100",
		Type:          1, // 满减
		Value:         10000, // 100元
		MinAmount:     50000, // 满500元
		TotalQuantity: 1000,
		PerUserLimit:  1,
		StartTime:     time.Now(),
		EndTime:       time.Now().AddDate(0, 1, 0), // 一个月后过期
		Description:   "新用户专享优惠券，满500减100",
		Status:        1,
	}
	if err := db.Create(coupon).Error; err != nil {
		return fmt.Errorf("创建优惠券失败: %w", err)
	}

	// 给用户发放优惠券
	userCoupon := &UserCoupon{
		UserID:   user.ID,
		CouponID: coupon.ID,
		Status:   1, // 未使用
	}
	if err := db.Create(userCoupon).Error; err != nil {
		return fmt.Errorf("发放优惠券失败: %w", err)
	}

	// 添加到购物车
	cart := &Cart{
		UserID:    user.ID,
		ProductID: product.ID,
		SKUID:     &sku1.ID,
		Quantity:  2,
	}
	if err := db.Create(cart).Error; err != nil {
		return fmt.Errorf("添加购物车失败: %w", err)
	}

	fmt.Println("测试数据填充完成")
	return nil
}

// demonstrateOrderService 演示订单服务
func demonstrateOrderService(db *gorm.DB) {
	fmt.Println("\n=== 演示订单服务 ===")

	orderService := services.NewOrderService(db)

	// 获取测试用户和地址
	var user User
	db.First(&user, "username = ?", "testuser")

	var address Address
	db.First(&address, "user_id = ?", user.ID)

	var sku ProductSKU
	db.First(&sku, "sku = ?", "IPHONE15PRO-128GB-BLACK")

	var coupon Coupon
	db.First(&coupon, "code = ?", "NEWUSER100")

	// 创建订单请求
	createOrderReq := &services.CreateOrderRequest{
		UserID:    user.ID,
		AddressID: address.ID,
		Items: []services.CreateOrderItemRequest{
			{
				ProductID: sku.ProductID,
				SKUID:     &sku.ID,
				Quantity:  1,
			},
		},
		CouponID: &coupon.ID,
		Remark:   "测试订单",
	}

	// 创建订单
	fmt.Println("创建订单...")
	order, err := orderService.CreateOrder(createOrderReq)
	if err != nil {
		fmt.Printf("创建订单失败: %v\n", err)
		return
	}

	fmt.Printf("订单创建成功: %s, 订单金额: %.2f元\n", order.OrderNo, float64(order.PayAmount)/100)

	// 查询订单详情
	var orderDetail Order
	db.Preload("Items").Preload("User").Preload("Coupon").First(&orderDetail, order.ID)
	fmt.Printf("订单详情: %+v\n", orderDetail)

	// 取消订单
	fmt.Println("\n取消订单...")
	err = orderService.CancelOrder(order.ID, user.ID, "用户主动取消")
	if err != nil {
		fmt.Printf("取消订单失败: %v\n", err)
	} else {
		fmt.Println("订单取消成功")
	}
}

// demonstrateStatisticsService 演示统计服务
func demonstrateStatisticsService(db *gorm.DB) {
	fmt.Println("\n=== 演示统计服务 ===")

	statisticsService := services.NewStatisticsService(db)

	// 创建一些测试订单数据
	createTestOrders(db)

	// 获取销售统计
	startDate := time.Now().AddDate(0, 0, -7) // 最近7天
	endDate := time.Now()

	fmt.Println("获取销售统计...")
	salesStats, err := statisticsService.GetSalesStatistics(startDate, endDate)
	if err != nil {
		fmt.Printf("获取销售统计失败: %v\n", err)
	} else {
		fmt.Printf("销售统计数据: %+v\n", salesStats)
	}

	// 获取商品销量排行
	fmt.Println("\n获取商品销量排行...")
	productRank, err := statisticsService.GetProductSalesRank(startDate, endDate, 10)
	if err != nil {
		fmt.Printf("获取商品销量排行失败: %v\n", err)
	} else {
		fmt.Printf("商品销量排行: %+v\n", productRank)
	}

	// 获取用户行为分析
	fmt.Println("\n获取用户行为分析...")
	userBehavior, err := statisticsService.GetUserBehaviorAnalysis(startDate, endDate, 10)
	if err != nil {
		fmt.Printf("获取用户行为分析失败: %v\n", err)
	} else {
		fmt.Printf("用户行为分析: %+v\n", userBehavior)
	}

	// 获取数据大屏数据
	fmt.Println("\n获取数据大屏数据...")
	dashboardData, err := statisticsService.GetDashboardData()
	if err != nil {
		fmt.Printf("获取数据大屏数据失败: %v\n", err)
	} else {
		fmt.Printf("数据大屏数据: %+v\n", dashboardData)
	}
}

// createTestOrders 创建测试订单数据
func createTestOrders(db *gorm.DB) {
	// 获取测试数据
	var user User
	db.First(&user, "username = ?", "testuser")

	var product Product
	db.First(&product)

	// 创建几个测试订单
	for i := 0; i < 5; i++ {
		order := &Order{
			OrderNo:         fmt.Sprintf("TEST%d%d", time.Now().Unix(), i),
			UserID:          user.ID,
			Status:          4, // 已完成
			TotalAmount:     int64((i + 1) * 10000), // 100元 * (i+1)
			PayAmount:       int64((i + 1) * 10000),
			ReceiverName:    "测试用户",
			ReceiverPhone:   "13800138000",
			ReceiverAddress: "测试地址",
			CreatedAt:       time.Now().AddDate(0, 0, -i), // 不同日期
		}
		db.Create(order)

		// 创建订单项
		orderItem := &OrderItem{
			OrderID:     order.ID,
			ProductID:   product.ID,
			Quantity:    i + 1,
			Price:       10000,
			TotalPrice:  int64((i + 1) * 10000),
			ProductName: product.Name,
			ProductSKU:  product.SKU,
		}
		db.Create(orderItem)
	}
}

// demonstrateComplexQueries 演示复杂查询
func demonstrateComplexQueries(db *gorm.DB) {
	fmt.Println("\n=== 演示复杂查询 ===")

	// 1. 子查询：查找购买过商品的用户
	fmt.Println("1. 查找购买过商品的用户:")
	var users []User
	db.Where("id IN (?)", db.Table("orders").Select("DISTINCT user_id").Where("status >= ?", 2)).Find(&users)
	fmt.Printf("购买过商品的用户数量: %d\n", len(users))

	// 2. 连接查询：查询用户及其订单统计
	fmt.Println("\n2. 查询用户及其订单统计:")
	type UserOrderStats struct {
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		OrderCount  int64  `json:"order_count"`
		TotalAmount int64  `json:"total_amount"`
	}
	var userStats []UserOrderStats
	db.Table("users u").
		Select("u.id as user_id, u.username, COUNT(o.id) as order_count, COALESCE(SUM(o.pay_amount), 0) as total_amount").
		Joins("LEFT JOIN orders o ON u.id = o.user_id AND o.status >= ?", 2).
		Group("u.id, u.username").
		Find(&userStats)
	for _, stat := range userStats {
		fmt.Printf("用户: %s, 订单数: %d, 总金额: %.2f元\n", stat.Username, stat.OrderCount, float64(stat.TotalAmount)/100)
	}

	// 3. 聚合查询：按分类统计商品数量和平均价格
	fmt.Println("\n3. 按分类统计商品数量和平均价格:")
	type CategoryStats struct {
		CategoryName string  `json:"category_name"`
		ProductCount int64   `json:"product_count"`
		AvgPrice     float64 `json:"avg_price"`
		MinPrice     int64   `json:"min_price"`
		MaxPrice     int64   `json:"max_price"`
	}
	var categoryStats []CategoryStats
	db.Table("categories c").
		Select("c.name as category_name, COUNT(p.id) as product_count, AVG(p.price) as avg_price, MIN(p.price) as min_price, MAX(p.price) as max_price").
		Joins("LEFT JOIN products p ON c.id = p.category_id AND p.status = ?", 1).
		Group("c.id, c.name").
		Find(&categoryStats)
	for _, stat := range categoryStats {
		fmt.Printf("分类: %s, 商品数: %d, 平均价格: %.2f元, 最低价格: %.2f元, 最高价格: %.2f元\n",
			stat.CategoryName, stat.ProductCount, stat.AvgPrice/100, float64(stat.MinPrice)/100, float64(stat.MaxPrice)/100)
	}

	// 4. 窗口函数：商品销量排名
	fmt.Println("\n4. 商品销量排名:")
	type ProductSalesRank struct {
		ProductName string `json:"product_name"`
		SalesCount  int64  `json:"sales_count"`
		Rank        int    `json:"rank"`
	}
	var productRanks []ProductSalesRank
	db.Raw(`
		SELECT 
			p.name as product_name,
			COALESCE(SUM(oi.quantity), 0) as sales_count,
			ROW_NUMBER() OVER (ORDER BY COALESCE(SUM(oi.quantity), 0) DESC) as rank
		FROM products p
		LEFT JOIN order_items oi ON p.id = oi.product_id
		LEFT JOIN orders o ON oi.order_id = o.id AND o.status >= 2
		WHERE p.status = 1
		GROUP BY p.id, p.name
		ORDER BY sales_count DESC
		LIMIT 10
	`).Scan(&productRanks)
	for _, rank := range productRanks {
		fmt.Printf("排名: %d, 商品: %s, 销量: %d\n", rank.Rank, rank.ProductName, rank.SalesCount)
	}
}

func main() {
	// 数据库配置
	config := DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "gorm_advanced_exercise2",
		Charset:  "utf8mb4",
	}

	// 连接数据库
	fmt.Println("连接数据库...")
	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 迁移数据库
	if err := MigrateDatabase(db); err != nil {
		log.Fatal("迁移数据库失败:", err)
	}

	// 检查是否需要填充测试数据
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		if err := SeedTestData(db); err != nil {
			log.Fatal("填充测试数据失败:", err)
		}
	}

	// 演示订单服务
	demonstrateOrderService(db)

	// 演示统计服务
	demonstrateStatisticsService(db)

	// 演示复杂查询
	demonstrateComplexQueries(db)

	fmt.Println("\n=== 练习2：复杂业务逻辑实现 演示完成 ===")
	fmt.Println("\n强化练习任务:")
	fmt.Println("1. 完善支付、发货、确认收货等业务逻辑")
	fmt.Println("2. 进行高并发下单测试")
	fmt.Println("3. 完善异常处理和事务回滚")
	fmt.Println("4. 优化数据库查询性能")
	fmt.Println("5. 为核心业务逻辑编写单元测试")
}
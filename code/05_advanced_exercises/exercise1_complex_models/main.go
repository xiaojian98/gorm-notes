package main

import (
	"fmt"
	"log"
	"time"

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

// GetDefaultConfig 获取默认数据库配置
func GetDefaultConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "gorm_ecommerce",
		Charset:  "utf8mb4",
	}
}

// ConnectDatabase 连接数据库
func ConnectDatabase(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层的sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                   // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)                  // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour)         // 设置了连接可复用的最大时间

	return db, nil
}

// demonstrateModels 演示模型功能
func demonstrateModels(db *gorm.DB) {
	fmt.Println("\n=== 演示复杂数据模型功能 ===")

	// 1. 演示用户创建和关联查询
	demonstateUserOperations(db)

	// 2. 演示商品和分类操作
	demonstrateProductOperations(db)

	// 3. 演示订单操作
	demonstrateOrderOperations(db)

	// 4. 演示优惠券操作
	demonstrateCouponOperations(db)
}

// demonstateUserOperations 演示用户操作
func demonstateUserOperations(db *gorm.DB) {
	fmt.Println("\n--- 用户操作演示 ---")

	// 查询用户及其关联数据
	var user User
	result := db.Preload("Profile").Preload("Addresses").First(&user, 1)
	if result.Error != nil {
		log.Printf("查询用户失败: %v", result.Error)
		return
	}

	fmt.Printf("用户信息: %s (%s)\n", user.Nickname, user.Email)
	if user.Profile != nil {
		fmt.Printf("用户详情: %s\n", user.Profile.RealName)
	}
	fmt.Printf("地址数量: %d\n", len(user.Addresses))

	// 创建新地址
	address := Address{
		UserID:    user.ID,
		Name:      "张三",
		Phone:     "13800138000",
		Province:  "北京市",
		City:      "北京市",
		District:  "朝阳区",
		Detail:    "某某街道123号",
		Postcode:  "100000",
		IsDefault: false,
	}

	if err := db.Create(&address).Error; err != nil {
		log.Printf("创建地址失败: %v", err)
	} else {
		fmt.Printf("创建地址成功: ID=%d\n", address.ID)
	}
}

// demonstrateProductOperations 演示商品操作
func demonstrateProductOperations(db *gorm.DB) {
	fmt.Println("\n--- 商品操作演示 ---")

	// 查询商品及其关联数据
	var products []Product
	result := db.Preload("Category").Preload("Brand").Preload("Images").Limit(5).Find(&products)
	if result.Error != nil {
		log.Printf("查询商品失败: %v", result.Error)
		return
	}

	fmt.Printf("商品列表 (共%d个):\n", len(products))
	for _, product := range products {
		fmt.Printf("- %s (SKU: %s) - 价格: %.2f元\n", 
			product.Name, 
			product.SKU, 
			float64(product.Price)/100)
		if product.Category.ID > 0 {
			fmt.Printf("  分类: %s\n", product.Category.Name)
		}
		if product.Brand.ID > 0 {
			fmt.Printf("  品牌: %s\n", product.Brand.Name)
		}
	}

	// 创建商品图片
	if len(products) > 0 {
		productImage := ProductImage{
			ProductID: products[0].ID,
			URL:       "https://example.com/image1.jpg",
			Alt:       "商品主图",
			Sort:      1,
			IsMain:    true,
		}

		if err := db.Create(&productImage).Error; err != nil {
			log.Printf("创建商品图片失败: %v", err)
		} else {
			fmt.Printf("创建商品图片成功: ID=%d\n", productImage.ID)
		}
	}
}

// demonstrateOrderOperations 演示订单操作
func demonstrateOrderOperations(db *gorm.DB) {
	fmt.Println("\n--- 订单操作演示 ---")

	// 创建测试订单
	order := Order{
		OrderNo:         fmt.Sprintf("ORD%d", time.Now().Unix()),
		UserID:          1,
		Status:          1, // 待付款
		TotalAmount:     99900,
		PayAmount:       99900,
		FreightAmount:   0,
		DiscountAmount:  0,
		ReceiverName:    "张三",
		ReceiverPhone:   "13800138000",
		ReceiverAddress: "北京市朝阳区某某街道123号",
		Remark:          "测试订单",
	}

	if err := db.Create(&order).Error; err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}

	fmt.Printf("创建订单成功: %s\n", order.OrderNo)

	// 创建订单项
	orderItem := OrderItem{
		OrderID:      order.ID,
		ProductID:    1,
		Quantity:     1,
		Price:        99900,
		TotalPrice:   99900,
		ProductName:  "iPhone 15 Pro",
		ProductSKU:   "IPHONE15PRO001",
		ProductImage: "https://example.com/iphone.jpg",
	}

	if err := db.Create(&orderItem).Error; err != nil {
		log.Printf("创建订单项失败: %v", err)
	} else {
		fmt.Printf("创建订单项成功: ID=%d\n", orderItem.ID)
	}

	// 查询订单及其关联数据
	var orderWithItems Order
	result := db.Preload("User").Preload("Items").First(&orderWithItems, order.ID)
	if result.Error != nil {
		log.Printf("查询订单失败: %v", result.Error)
		return
	}

	fmt.Printf("订单详情: %s - 用户: %s - 商品数量: %d\n", 
		orderWithItems.OrderNo, 
		orderWithItems.User.Nickname,
		len(orderWithItems.Items))
}

// demonstrateCouponOperations 演示优惠券操作
func demonstrateCouponOperations(db *gorm.DB) {
	fmt.Println("\n--- 优惠券操作演示 ---")

	// 创建优惠券
	coupon := Coupon{
		Name:          "新用户专享",
		Code:          "NEWUSER100",
		Type:          1, // 满减
		Value:         10000, // 100元
		MinAmount:     50000, // 满500元
		MaxDiscount:   10000, // 最大优惠100元
		TotalQuantity: 1000,
		UsedQuantity:  0,
		UserLimit:     1,
		StartTime:     time.Now(),
		EndTime:       time.Now().AddDate(0, 1, 0), // 一个月后过期
		Status:        1,
		Description:   "新用户专享优惠券，满500减100",
	}

	if err := db.Create(&coupon).Error; err != nil {
		log.Printf("创建优惠券失败: %v", err)
		return
	}

	fmt.Printf("创建优惠券成功: %s\n", coupon.Code)

	// 用户领取优惠券
	userCoupon := UserCoupon{
		UserID:   1,
		CouponID: coupon.ID,
		Status:   1, // 未使用
	}

	if err := db.Create(&userCoupon).Error; err != nil {
		log.Printf("用户领取优惠券失败: %v", err)
	} else {
		fmt.Printf("用户领取优惠券成功: ID=%d\n", userCoupon.ID)
	}

	// 查询用户的优惠券
	var userCoupons []UserCoupon
	result := db.Preload("Coupon").Where("user_id = ? AND status = ?", 1, 1).Find(&userCoupons)
	if result.Error != nil {
		log.Printf("查询用户优惠券失败: %v", result.Error)
		return
	}

	fmt.Printf("用户可用优惠券数量: %d\n", len(userCoupons))
	for _, uc := range userCoupons {
		fmt.Printf("- %s: %s\n", uc.Coupon.Name, uc.Coupon.Description)
	}
}

// demonstrateComplexQueries 演示复杂查询
func demonstrateComplexQueries(db *gorm.DB) {
	fmt.Println("\n=== 演示复杂查询功能 ===")

	// 1. 统计查询
	demonstateStatistics(db)

	// 2. 关联查询
	demonstrateJoinQueries(db)

	// 3. 聚合查询
	demonstrateAggregateQueries(db)
}

// demonstateStatistics 演示统计查询
func demonstateStatistics(db *gorm.DB) {
	fmt.Println("\n--- 统计查询演示 ---")

	// 统计用户数量
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	fmt.Printf("用户总数: %d\n", userCount)

	// 统计商品数量
	var productCount int64
	db.Model(&Product{}).Count(&productCount)
	fmt.Printf("商品总数: %d\n", productCount)

	// 统计订单数量和总金额
	var orderCount int64
	var totalAmount int64
	db.Model(&Order{}).Count(&orderCount)
	db.Model(&Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)
	fmt.Printf("订单总数: %d, 总金额: %.2f元\n", orderCount, float64(totalAmount)/100)
}

// demonstrateJoinQueries 演示关联查询
func demonstrateJoinQueries(db *gorm.DB) {
	fmt.Println("\n--- 关联查询演示 ---")

	// 查询商品及其分类和品牌信息
	type ProductInfo struct {
		ProductName  string
		ProductSKU   string
		Price        int64
		CategoryName string
		BrandName    string
	}

	var productInfos []ProductInfo
	result := db.Table("products p").
		Select("p.name as product_name, p.sku as product_sku, p.price, c.name as category_name, b.name as brand_name").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Joins("LEFT JOIN brands b ON p.brand_id = b.id").
		Where("p.status = ?", 1).
		Limit(5).
		Find(&productInfos)

	if result.Error != nil {
		log.Printf("关联查询失败: %v", result.Error)
		return
	}

	fmt.Println("商品信息列表:")
	for _, info := range productInfos {
		fmt.Printf("- %s (%s) - %.2f元 - %s/%s\n", 
			info.ProductName, 
			info.ProductSKU, 
			float64(info.Price)/100,
			info.CategoryName,
			info.BrandName)
	}
}

// demonstrateAggregateQueries 演示聚合查询
func demonstrateAggregateQueries(db *gorm.DB) {
	fmt.Println("\n--- 聚合查询演示 ---")

	// 按分类统计商品数量
	type CategoryStats struct {
		CategoryName string
		ProductCount int64
		AvgPrice     float64
	}

	var categoryStats []CategoryStats
	result := db.Table("products p").
		Select("c.name as category_name, COUNT(p.id) as product_count, AVG(p.price) as avg_price").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Where("p.status = ?", 1).
		Group("c.id, c.name").
		Find(&categoryStats)

	if result.Error != nil {
		log.Printf("聚合查询失败: %v", result.Error)
		return
	}

	fmt.Println("分类统计:")
	for _, stats := range categoryStats {
		fmt.Printf("- %s: %d个商品, 平均价格: %.2f元\n", 
			stats.CategoryName, 
			stats.ProductCount, 
			stats.AvgPrice/100)
	}
}

func main() {
	fmt.Println("GORM强化练习 - 练习1：复杂数据模型设计")
	fmt.Println("========================================")

	// 1. 连接数据库
	config := GetDefaultConfig()
	fmt.Printf("连接数据库: %s@%s:%d/%s\n", config.User, config.Host, config.Port, config.DBName)

	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	fmt.Println("数据库连接成功!")

	// 2. 执行数据库迁移
	if err := MigrateDatabase(db); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 3. 插入测试数据
	if err := SeedData(db); err != nil {
		log.Printf("插入测试数据失败: %v", err)
	}

	// 4. 演示模型功能
	demonstrateModels(db)

	// 5. 演示复杂查询
	demonstrateComplexQueries(db)

	fmt.Println("\n=== 练习1完成 ===")
	fmt.Println("\n练习要点总结:")
	fmt.Println("1. 设计了完整的电商数据模型")
	fmt.Println("2. 实现了复杂的关联关系")
	fmt.Println("3. 添加了适当的索引和约束")
	fmt.Println("4. 实现了软删除和审计字段")
	fmt.Println("5. 演示了各种查询操作")
	fmt.Println("\n下一步可以尝试:")
	fmt.Println("- 优化查询性能")
	fmt.Println("- 添加更多业务逻辑")
	fmt.Println("- 实现数据验证")
	fmt.Println("- 编写单元测试")
}
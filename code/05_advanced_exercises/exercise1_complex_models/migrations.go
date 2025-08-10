package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// MigrateDatabase 执行数据库迁移
func MigrateDatabase(db *gorm.DB) error {
	fmt.Println("开始执行数据库迁移...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&Address{},
		&Category{},
		&Brand{},
		&Product{},
		&ProductImage{},
		&ProductAttr{},
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
		return fmt.Errorf("自动迁移失败: %v", err)
	}

	// 创建索引
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("创建索引失败: %v", err)
	}

	// 添加外键约束
	if err := addForeignKeys(db); err != nil {
		return fmt.Errorf("添加外键约束失败: %v", err)
	}

	fmt.Println("数据库迁移完成!")
	return nil
}

// createIndexes 创建索引
func createIndexes(db *gorm.DB) error {
	fmt.Println("创建索引...")

	// 用户表索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone)",
		"CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)",
		
		// 地址表索引
		"CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_addresses_is_default ON addresses(is_default)",
		
		// 分类表索引
		"CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id)",
		"CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug)",
		"CREATE INDEX IF NOT EXISTS idx_categories_status ON categories(status)",
		"CREATE INDEX IF NOT EXISTS idx_categories_sort ON categories(sort)",
		
		// 品牌表索引
		"CREATE INDEX IF NOT EXISTS idx_brands_slug ON brands(slug)",
		"CREATE INDEX IF NOT EXISTS idx_brands_status ON brands(status)",
		
		// 商品表索引
		"CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_brand_id ON products(brand_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku)",
		"CREATE INDEX IF NOT EXISTS idx_products_status ON products(status)",
		"CREATE INDEX IF NOT EXISTS idx_products_price ON products(price)",
		"CREATE INDEX IF NOT EXISTS idx_products_sales ON products(sales)",
		"CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at)",
		
		// 商品图片表索引
		"CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_product_images_is_main ON product_images(is_main)",
		
		// 商品属性表索引
		"CREATE INDEX IF NOT EXISTS idx_product_attrs_product_id ON product_attrs(product_id)",
		
		// 商品SKU表索引
		"CREATE INDEX IF NOT EXISTS idx_product_skus_product_id ON product_skus(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_product_skus_sku ON product_skus(sku)",
		"CREATE INDEX IF NOT EXISTS idx_product_skus_status ON product_skus(status)",
		
		// 商品评价表索引
		"CREATE INDEX IF NOT EXISTS idx_product_reviews_product_id ON product_reviews(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_product_reviews_user_id ON product_reviews(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_product_reviews_order_id ON product_reviews(order_id)",
		"CREATE INDEX IF NOT EXISTS idx_product_reviews_rating ON product_reviews(rating)",
		"CREATE INDEX IF NOT EXISTS idx_product_reviews_status ON product_reviews(status)",
		
		// 购物车表索引
		"CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_carts_product_id ON carts(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_carts_sku_id ON carts(sku_id)",
		
		// 订单表索引
		"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_orders_order_no ON orders(order_no)",
		"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)",
		"CREATE INDEX IF NOT EXISTS idx_orders_coupon_id ON orders(coupon_id)",
		"CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_orders_payment_time ON orders(payment_time)",
		
		// 订单项表索引
		"CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)",
		"CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_order_items_sku_id ON order_items(sku_id)",
		
		// 支付记录表索引
		"CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id)",
		"CREATE INDEX IF NOT EXISTS idx_payments_payment_no ON payments(payment_no)",
		"CREATE INDEX IF NOT EXISTS idx_payments_method ON payments(method)",
		"CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)",
		"CREATE INDEX IF NOT EXISTS idx_payments_third_party_no ON payments(third_party_no)",
		
		// 优惠券表索引
		"CREATE INDEX IF NOT EXISTS idx_coupons_code ON coupons(code)",
		"CREATE INDEX IF NOT EXISTS idx_coupons_type ON coupons(type)",
		"CREATE INDEX IF NOT EXISTS idx_coupons_status ON coupons(status)",
		"CREATE INDEX IF NOT EXISTS idx_coupons_start_time ON coupons(start_time)",
		"CREATE INDEX IF NOT EXISTS idx_coupons_end_time ON coupons(end_time)",
		
		// 用户优惠券表索引
		"CREATE INDEX IF NOT EXISTS idx_user_coupons_user_id ON user_coupons(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_coupons_coupon_id ON user_coupons(coupon_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_coupons_status ON user_coupons(status)",
		"CREATE INDEX IF NOT EXISTS idx_user_coupons_used_at ON user_coupons(used_at)",
	}

	for _, indexSQL := range indexes {
		if err := db.Exec(indexSQL).Error; err != nil {
			log.Printf("创建索引失败: %s, 错误: %v", indexSQL, err)
			// 继续执行其他索引，不中断
		}
	}

	fmt.Println("索引创建完成!")
	return nil
}

// addForeignKeys 添加外键约束
func addForeignKeys(db *gorm.DB) error {
	fmt.Println("添加外键约束...")

	// 注意：在生产环境中，外键约束可能会影响性能，需要根据实际情况决定是否使用
	foreignKeys := []string{
		// 用户详细资料外键
		"ALTER TABLE user_profiles ADD CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		
		// 地址外键
		"ALTER TABLE addresses ADD CONSTRAINT fk_addresses_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		
		// 分类自关联外键
		"ALTER TABLE categories ADD CONSTRAINT fk_categories_parent_id FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL",
		
		// 商品外键
		"ALTER TABLE products ADD CONSTRAINT fk_products_category_id FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT",
		"ALTER TABLE products ADD CONSTRAINT fk_products_brand_id FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE RESTRICT",
		
		// 商品图片外键
		"ALTER TABLE product_images ADD CONSTRAINT fk_product_images_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE",
		
		// 商品属性外键
		"ALTER TABLE product_attrs ADD CONSTRAINT fk_product_attrs_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE",
		
		// 商品SKU外键
		"ALTER TABLE product_skus ADD CONSTRAINT fk_product_skus_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE",
		
		// 商品评价外键
		"ALTER TABLE product_reviews ADD CONSTRAINT fk_product_reviews_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE",
		"ALTER TABLE product_reviews ADD CONSTRAINT fk_product_reviews_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		"ALTER TABLE product_reviews ADD CONSTRAINT fk_product_reviews_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE",
		
		// 购物车外键
		"ALTER TABLE carts ADD CONSTRAINT fk_carts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		"ALTER TABLE carts ADD CONSTRAINT fk_carts_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE",
		"ALTER TABLE carts ADD CONSTRAINT fk_carts_sku_id FOREIGN KEY (sku_id) REFERENCES product_skus(id) ON DELETE CASCADE",
		
		// 订单外键
		"ALTER TABLE orders ADD CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT",
		"ALTER TABLE orders ADD CONSTRAINT fk_orders_coupon_id FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE SET NULL",
		
		// 订单项外键
		"ALTER TABLE order_items ADD CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE",
		"ALTER TABLE order_items ADD CONSTRAINT fk_order_items_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT",
		"ALTER TABLE order_items ADD CONSTRAINT fk_order_items_sku_id FOREIGN KEY (sku_id) REFERENCES product_skus(id) ON DELETE RESTRICT",
		
		// 支付记录外键
		"ALTER TABLE payments ADD CONSTRAINT fk_payments_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE",
		
		// 用户优惠券外键
		"ALTER TABLE user_coupons ADD CONSTRAINT fk_user_coupons_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE",
		"ALTER TABLE user_coupons ADD CONSTRAINT fk_user_coupons_coupon_id FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE",
	}

	for _, fkSQL := range foreignKeys {
		if err := db.Exec(fkSQL).Error; err != nil {
			log.Printf("添加外键约束失败: %s, 错误: %v", fkSQL, err)
			// 继续执行其他外键，不中断
		}
	}

	fmt.Println("外键约束添加完成!")
	return nil
}

// SeedData 插入测试数据
func SeedData(db *gorm.DB) error {
	fmt.Println("开始插入测试数据...")

	// 创建测试用户
	users := []User{
		{
			Email:    "admin@example.com",
			Phone:    "13800138000",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Nickname: "管理员",
			Gender:   1,
			Status:   1,
			Points:   1000,
		},
		{
			Email:    "user1@example.com",
			Phone:    "13800138001",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Nickname: "用户1",
			Gender:   2,
			Status:   1,
			Points:   500,
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("创建用户失败: %v", err)
		}
	}

	// 创建测试分类
	categories := []Category{
		{
			Name:        "电子产品",
			Slug:        "electronics",
			Description: "各类电子产品",
			Sort:        1,
			Status:      1,
		},
		{
			Name:        "服装",
			Slug:        "clothing",
			Description: "时尚服装",
			Sort:        2,
			Status:      1,
		},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			log.Printf("创建分类失败: %v", err)
		}
	}

	// 创建测试品牌
	brands := []Brand{
		{
			Name:        "苹果",
			Slug:        "apple",
			Description: "苹果公司",
			Website:     "https://www.apple.com",
			Sort:        1,
			Status:      1,
		},
		{
			Name:        "耐克",
			Slug:        "nike",
			Description: "耐克公司",
			Website:     "https://www.nike.com",
			Sort:        2,
			Status:      1,
		},
	}

	for _, brand := range brands {
		if err := db.Create(&brand).Error; err != nil {
			log.Printf("创建品牌失败: %v", err)
		}
	}

	// 创建测试商品
	products := []Product{
		{
			Name:        "iPhone 15 Pro",
			SKU:         "IPHONE15PRO001",
			Description: "最新款iPhone",
			Price:       999900, // 9999.00元
			MarketPrice: 1099900,
			CostPrice:   799900,
			Stock:       100,
			Weight:      0.2,
			Volume:      0.001,
			Status:      1,
			CategoryID:  1,
			BrandID:     1,
		},
		{
			Name:        "Nike Air Max",
			SKU:         "NIKEAIRMAX001",
			Description: "经典运动鞋",
			Price:       89900, // 899.00元
			MarketPrice: 99900,
			CostPrice:   59900,
			Stock:       200,
			Weight:      0.5,
			Volume:      0.005,
			Status:      1,
			CategoryID:  2,
			BrandID:     2,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			log.Printf("创建商品失败: %v", err)
		}
	}

	fmt.Println("测试数据插入完成!")
	return nil
}
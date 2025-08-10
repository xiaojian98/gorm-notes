# GORM强化练习：综合能力提升

## 📚 练习说明

强化练习是在掌握基础技能后的综合能力提升训练。这些练习将多个GORM特性结合起来，模拟真实的业务场景，帮助你建立系统性的思维和解决复杂问题的能力。

### 🎯 练习目标
- 综合运用GORM的各种特性
- 解决复杂的业务场景问题
- 培养系统设计和架构思维
- 掌握性能优化和错误处理
- 学会编写可维护的代码

### 📋 环境准备
```bash
# 创建强化练习项目
mkdir gorm-advanced-exercises
cd gorm-advanced-exercises
go mod init gorm-advanced-exercises

# 安装依赖
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/go-redis/redis/v8
go get -u github.com/gin-gonic/gin
```

---

## 🏢 综合项目：企业级电商系统

### 项目背景

我们将构建一个企业级电商系统的核心数据层，包含用户管理、商品管理、订单处理、库存管理、优惠券系统等模块。这个项目将综合运用GORM的所有重要特性。

### 系统架构

```
电商系统
├── 用户模块 (User Management)
│   ├── 用户注册/登录
│   ├── 用户资料管理
│   └── 地址管理
├── 商品模块 (Product Management)
│   ├── 商品分类
│   ├── 商品信息
│   ├── 商品规格/SKU
│   └── 库存管理
├── 订单模块 (Order Management)
│   ├── 购物车
│   ├── 订单创建
│   ├── 订单状态流转
│   └── 支付处理
├── 营销模块 (Marketing)
│   ├── 优惠券系统
│   ├── 促销活动
│   └── 积分系统
└── 系统模块 (System)
    ├── 操作日志
    ├── 数据统计
    └── 缓存管理
```

---

## 🚀 强化练习1：复杂数据模型设计

### 目标
设计一个完整的电商数据模型，包含所有必要的关联关系、索引和约束。

### 要求
1. 设计用户、商品、订单等核心模型
2. 实现复杂的关联关系（一对一、一对多、多对多）
3. 添加适当的索引和约束
4. 实现软删除和审计字段
5. 使用GORM标签优化数据库结构

### 代码实现

```go
package models

import (
	"time"
	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string    `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Email       string    `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Phone       string    `json:"phone" gorm:"size:20;index"`
	Password    string    `json:"-" gorm:"size:255;not null"` // 不在JSON中显示
	Nickname    string    `json:"nickname" gorm:"size:50"`
	Avatar      string    `json:"avatar" gorm:"size:255"`
	Gender      int8      `json:"gender" gorm:"default:0"` // 0:未知 1:男 2:女
	Birthday    *time.Time `json:"birthday"`
	Status      int8      `json:"status" gorm:"default:1;index"` // 1:正常 2:禁用
	LastLoginAt *time.Time `json:"last_login_at"`
	Points      int       `json:"points" gorm:"default:0"` // 积分
	
	// 关联关系
	Profile   *UserProfile `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Addresses []Address    `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	Orders    []Order      `json:"orders,omitempty" gorm:"foreignKey:UserID"`
	Carts     []Cart       `json:"carts,omitempty" gorm:"foreignKey:UserID"`
	Coupons   []Coupon     `json:"coupons,omitempty" gorm:"many2many:user_coupons;"`
}

// UserProfile 用户详细资料
type UserProfile struct {
	ID           uint   `json:"id" gorm:"primarykey"`
	UserID       uint   `json:"user_id" gorm:"uniqueIndex;not null"`
	RealName     string `json:"real_name" gorm:"size:50"`
	IDCard       string `json:"id_card" gorm:"size:18;index"`
	Company      string `json:"company" gorm:"size:100"`
	Occupation   string `json:"occupation" gorm:"size:50"`
	Introduction string `json:"introduction" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// 关联关系
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Address 地址模型
type Address struct {
	BaseModel
	UserID      uint   `json:"user_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"size:50;not null"`
	Phone       string `json:"phone" gorm:"size:20;not null"`
	Province    string `json:"province" gorm:"size:50;not null"`
	City        string `json:"city" gorm:"size:50;not null"`
	District    string `json:"district" gorm:"size:50;not null"`
	Detail      string `json:"detail" gorm:"size:255;not null"`
	PostalCode  string `json:"postal_code" gorm:"size:10"`
	IsDefault   bool   `json:"is_default" gorm:"default:false;index"`
	
	// 关联关系
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Category 商品分类
type Category struct {
	BaseModel
	Name        string `json:"name" gorm:"size:50;not null;index"`
	Slug        string `json:"slug" gorm:"size:50;uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:text"`
	Image       string `json:"image" gorm:"size:255"`
	ParentID    *uint  `json:"parent_id" gorm:"index"`
	Sort        int    `json:"sort" gorm:"default:0;index"`
	Status      int8   `json:"status" gorm:"default:1;index"`
	
	// 自关联
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	
	// 关联关系
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// Brand 品牌
type Brand struct {
	BaseModel
	Name        string `json:"name" gorm:"size:50;not null;index"`
	Slug        string `json:"slug" gorm:"size:50;uniqueIndex;not null"`
	Logo        string `json:"logo" gorm:"size:255"`
	Description string `json:"description" gorm:"type:text"`
	Website     string `json:"website" gorm:"size:255"`
	Status      int8   `json:"status" gorm:"default:1;index"`
	
	// 关联关系
	Products []Product `json:"products,omitempty" gorm:"foreignKey:BrandID"`
}

// Product 商品模型
type Product struct {
	BaseModel
	Name         string  `json:"name" gorm:"size:255;not null;index"`
	Slug         string  `json:"slug" gorm:"size:255;uniqueIndex;not null"`
	SKU          string  `json:"sku" gorm:"size:50;uniqueIndex;not null"`
	Description  string  `json:"description" gorm:"type:text"`
	Content      string  `json:"content" gorm:"type:longtext"`
	Price        float64 `json:"price" gorm:"type:decimal(10,2);not null;index"`
	OriginalPrice float64 `json:"original_price" gorm:"type:decimal(10,2)"`
	CostPrice    float64 `json:"cost_price" gorm:"type:decimal(10,2)"`
	Stock        int     `json:"stock" gorm:"default:0;index"`
	Sales        int     `json:"sales" gorm:"default:0;index"`
	Views        int     `json:"views" gorm:"default:0;index"`
	Weight       float64 `json:"weight" gorm:"type:decimal(8,2)"`
	Volume       float64 `json:"volume" gorm:"type:decimal(8,2)"`
	Status       int8    `json:"status" gorm:"default:1;index"` // 1:上架 2:下架
	CategoryID   uint    `json:"category_id" gorm:"not null;index"`
	BrandID      *uint   `json:"brand_id" gorm:"index"`
	
	// 关联关系
	Category   Category        `json:"category" gorm:"foreignKey:CategoryID"`
	Brand      *Brand          `json:"brand,omitempty" gorm:"foreignKey:BrandID"`
	Images     []ProductImage  `json:"images,omitempty" gorm:"foreignKey:ProductID"`
	Attributes []ProductAttr   `json:"attributes,omitempty" gorm:"foreignKey:ProductID"`
	SKUs       []ProductSKU    `json:"skus,omitempty" gorm:"foreignKey:ProductID"`
	Reviews    []ProductReview `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
}

// ProductImage 商品图片
type ProductImage struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	ProductID uint   `json:"product_id" gorm:"not null;index"`
	URL       string `json:"url" gorm:"size:255;not null"`
	Alt       string `json:"alt" gorm:"size:255"`
	Sort      int    `json:"sort" gorm:"default:0"`
	IsMain    bool   `json:"is_main" gorm:"default:false"`
	
	// 关联关系
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// ProductAttr 商品属性
type ProductAttr struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	ProductID uint   `json:"product_id" gorm:"not null;index"`
	Name      string `json:"name" gorm:"size:50;not null"`
	Value     string `json:"value" gorm:"size:255;not null"`
	Sort      int    `json:"sort" gorm:"default:0"`
	
	// 关联关系
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// ProductSKU 商品SKU
type ProductSKU struct {
	BaseModel
	ProductID     uint    `json:"product_id" gorm:"not null;index"`
	SKU           string  `json:"sku" gorm:"size:50;uniqueIndex;not null"`
	Name          string  `json:"name" gorm:"size:255;not null"`
	Price         float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	OriginalPrice float64 `json:"original_price" gorm:"type:decimal(10,2)"`
	Stock         int     `json:"stock" gorm:"default:0"`
	Sales         int     `json:"sales" gorm:"default:0"`
	Image         string  `json:"image" gorm:"size:255"`
	Specs         string  `json:"specs" gorm:"type:json"` // JSON格式存储规格信息
	Status        int8    `json:"status" gorm:"default:1"`
	
	// 关联关系
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// ProductReview 商品评价
type ProductReview struct {
	BaseModel
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	UserID    uint    `json:"user_id" gorm:"not null;index"`
	OrderID   uint    `json:"order_id" gorm:"not null;index"`
	Rating    int8    `json:"rating" gorm:"not null;index"` // 1-5星
	Content   string  `json:"content" gorm:"type:text"`
	Images    string  `json:"images" gorm:"type:json"` // JSON格式存储图片
	Reply     string  `json:"reply" gorm:"type:text"`
	Status    int8    `json:"status" gorm:"default:1;index"`
	
	// 关联关系
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
	Order   Order   `json:"order" gorm:"foreignKey:OrderID"`
}

// Cart 购物车
type Cart struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	SKUID     *uint     `json:"sku_id"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Selected  bool      `json:"selected" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 关联关系
	User    User        `json:"user" gorm:"foreignKey:UserID"`
	Product Product     `json:"product" gorm:"foreignKey:ProductID"`
	SKU     *ProductSKU `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Order 订单模型
type Order struct {
	BaseModel
	OrderNo       string    `json:"order_no" gorm:"size:32;uniqueIndex;not null"`
	UserID        uint      `json:"user_id" gorm:"not null;index"`
	Status        int8      `json:"status" gorm:"not null;index"` // 订单状态
	PaymentStatus int8      `json:"payment_status" gorm:"default:1;index"` // 支付状态
	ShipStatus    int8      `json:"ship_status" gorm:"default:1;index"` // 发货状态
	TotalAmount   float64   `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	PayAmount     float64   `json:"pay_amount" gorm:"type:decimal(10,2);not null"`
	Freight       float64   `json:"freight" gorm:"type:decimal(10,2);default:0"`
	Discount      float64   `json:"discount" gorm:"type:decimal(10,2);default:0"`
	CouponID      *uint     `json:"coupon_id"`
	Remark        string    `json:"remark" gorm:"size:500"`
	PaidAt        *time.Time `json:"paid_at"`
	ShippedAt     *time.Time `json:"shipped_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	
	// 收货地址信息（快照）
	ReceiverName     string `json:"receiver_name" gorm:"size:50;not null"`
	ReceiverPhone    string `json:"receiver_phone" gorm:"size:20;not null"`
	ReceiverProvince string `json:"receiver_province" gorm:"size:50;not null"`
	ReceiverCity     string `json:"receiver_city" gorm:"size:50;not null"`
	ReceiverDistrict string `json:"receiver_district" gorm:"size:50;not null"`
	ReceiverDetail   string `json:"receiver_detail" gorm:"size:255;not null"`
	ReceiverPostal   string `json:"receiver_postal" gorm:"size:10"`
	
	// 关联关系
	User       User        `json:"user" gorm:"foreignKey:UserID"`
	Coupon     *Coupon     `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	Payments   []Payment   `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderItem 订单项
type OrderItem struct {
	ID            uint    `json:"id" gorm:"primarykey"`
	OrderID       uint    `json:"order_id" gorm:"not null;index"`
	ProductID     uint    `json:"product_id" gorm:"not null"`
	SKUID         *uint   `json:"sku_id"`
	Quantity      int     `json:"quantity" gorm:"not null"`
	Price         float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	TotalAmount   float64 `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	
	// 商品信息快照
	ProductName  string `json:"product_name" gorm:"size:255;not null"`
	ProductImage string `json:"product_image" gorm:"size:255"`
	ProductSKU   string `json:"product_sku" gorm:"size:50"`
	ProductSpecs string `json:"product_specs" gorm:"type:json"`
	
	// 关联关系
	Order   Order       `json:"order" gorm:"foreignKey:OrderID"`
	Product Product     `json:"product" gorm:"foreignKey:ProductID"`
	SKU     *ProductSKU `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Payment 支付记录
type Payment struct {
	BaseModel
	OrderID       uint      `json:"order_id" gorm:"not null;index"`
	PaymentNo     string    `json:"payment_no" gorm:"size:64;uniqueIndex;not null"`
	Method        string    `json:"method" gorm:"size:20;not null;index"` // alipay, wechat, bank
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status        int8      `json:"status" gorm:"default:1;index"` // 1:待支付 2:已支付 3:已退款
	TransactionID string    `json:"transaction_id" gorm:"size:64;index"`
	PaidAt        *time.Time `json:"paid_at"`
	Remark        string    `json:"remark" gorm:"size:500"`
	
	// 关联关系
	Order Order `json:"order" gorm:"foreignKey:OrderID"`
}

// Coupon 优惠券
type Coupon struct {
	BaseModel
	Name        string     `json:"name" gorm:"size:100;not null"`
	Code        string     `json:"code" gorm:"size:20;uniqueIndex;not null"`
	Type        int8       `json:"type" gorm:"not null;index"` // 1:满减 2:折扣 3:固定金额
	Value       float64    `json:"value" gorm:"type:decimal(10,2);not null"`
	MinAmount   float64    `json:"min_amount" gorm:"type:decimal(10,2);default:0"`
	MaxDiscount float64    `json:"max_discount" gorm:"type:decimal(10,2);default:0"`
	Total       int        `json:"total" gorm:"not null"`
	Used        int        `json:"used" gorm:"default:0"`
	StartAt     time.Time  `json:"start_at" gorm:"not null;index"`
	EndAt       time.Time  `json:"end_at" gorm:"not null;index"`
	Status      int8       `json:"status" gorm:"default:1;index"`
	
	// 关联关系
	Users  []User  `json:"users,omitempty" gorm:"many2many:user_coupons;"`
	Orders []Order `json:"orders,omitempty" gorm:"foreignKey:CouponID"`
}

// UserCoupon 用户优惠券中间表
type UserCoupon struct {
	ID       uint       `json:"id" gorm:"primarykey"`
	UserID   uint       `json:"user_id" gorm:"not null;index"`
	CouponID uint       `json:"coupon_id" gorm:"not null;index"`
	UsedAt   *time.Time `json:"used_at"`
	OrderID  *uint      `json:"order_id"`
	Status   int8       `json:"status" gorm:"default:1;index"` // 1:未使用 2:已使用 3:已过期
	CreatedAt time.Time `json:"created_at"`
	
	// 关联关系
	User   User    `json:"user" gorm:"foreignKey:UserID"`
	Coupon Coupon  `json:"coupon" gorm:"foreignKey:CouponID"`
	Order  *Order  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// 表名定义
func (User) TableName() string { return "users" }
func (UserProfile) TableName() string { return "user_profiles" }
func (Address) TableName() string { return "addresses" }
func (Category) TableName() string { return "categories" }
func (Brand) TableName() string { return "brands" }
func (Product) TableName() string { return "products" }
func (ProductImage) TableName() string { return "product_images" }
func (ProductAttr) TableName() string { return "product_attrs" }
func (ProductSKU) TableName() string { return "product_skus" }
func (ProductReview) TableName() string { return "product_reviews" }
func (Cart) TableName() string { return "carts" }
func (Order) TableName() string { return "orders" }
func (OrderItem) TableName() string { return "order_items" }
func (Payment) TableName() string { return "payments" }
func (Coupon) TableName() string { return "coupons" }
func (UserCoupon) TableName() string { return "user_coupons" }
```

### 练习任务

1. **分析模型关系**：画出ER图，理解各模型之间的关联关系
2. **优化索引设计**：为高频查询字段添加合适的索引
3. **实现数据迁移**：编写迁移脚本，创建所有表和索引
4. **数据验证**：为关键字段添加验证规则
5. **性能测试**：创建测试数据，验证查询性能

---

## 🔄 强化练习2：复杂业务逻辑实现

### 目标
实现电商系统的核心业务逻辑，包括下单流程、库存管理、优惠券使用等。

### 要求
1. 实现完整的下单流程（库存检查、价格计算、订单创建）
2. 实现库存管理（扣减、回滚、预占）
3. 实现优惠券系统（发放、使用、验证）
4. 使用事务保证数据一致性
5. 实现并发安全的库存操作

### 代码实现

```go
package services

import (
	"errors"
	"fmt"
	"time"
	"gorm.io/gorm"
	"context"
	"sync"
)

// OrderService 订单服务
type OrderService struct {
	db    *gorm.DB
	mutex sync.RWMutex
}

// NewOrderService 创建订单服务
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID    uint                `json:"user_id"`
	Items     []OrderItemRequest  `json:"items"`
	AddressID uint                `json:"address_id"`
	CouponID  *uint               `json:"coupon_id,omitempty"`
	Remark    string              `json:"remark"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	SKUID     *uint `json:"sku_id,omitempty"`
	Quantity  int  `json:"quantity"`
}

// CreateOrder 创建订单（核心业务逻辑）
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	// 参数验证
	if err := s.validateCreateOrderRequest(req); err != nil {
		return nil, err
	}
	
	// 使用事务确保数据一致性
	var order *Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 验证用户
		var user User
		if err := tx.First(&user, req.UserID).Error; err != nil {
			return errors.New("用户不存在")
		}
		
		// 2. 验证收货地址
		var address Address
		if err := tx.Where("id = ? AND user_id = ?", req.AddressID, req.UserID).First(&address).Error; err != nil {
			return errors.New("收货地址不存在")
		}
		
		// 3. 验证商品和库存
		orderItems, totalAmount, err := s.validateAndCalculateItems(tx, req.Items)
		if err != nil {
			return err
		}
		
		// 4. 验证和使用优惠券
		discount, err := s.validateAndUseCoupon(tx, req.CouponID, req.UserID, totalAmount)
		if err != nil {
			return err
		}
		
		// 5. 计算最终金额
		freight := s.calculateFreight(totalAmount)
		payAmount := totalAmount + freight - discount
		
		// 6. 创建订单
		order = &Order{
			OrderNo:          s.generateOrderNo(),
			UserID:           req.UserID,
			Status:           1, // 待支付
			PaymentStatus:    1, // 待支付
			ShipStatus:       1, // 待发货
			TotalAmount:      totalAmount,
			PayAmount:        payAmount,
			Freight:          freight,
			Discount:         discount,
			CouponID:         req.CouponID,
			Remark:           req.Remark,
			ReceiverName:     address.Name,
			ReceiverPhone:    address.Phone,
			ReceiverProvince: address.Province,
			ReceiverCity:     address.City,
			ReceiverDistrict: address.District,
			ReceiverDetail:   address.Detail,
			ReceiverPostal:   address.PostalCode,
		}
		
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		
		// 7. 创建订单项
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}
		
		// 8. 扣减库存
		if err := s.deductStock(tx, req.Items); err != nil {
			return err
		}
		
		// 9. 清空购物车（如果是从购物车下单）
		if err := s.clearCart(tx, req.UserID, req.Items); err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 加载完整的订单信息
	s.db.Preload("OrderItems").Preload("User").First(order, order.ID)
	
	return order, nil
}

// validateCreateOrderRequest 验证创建订单请求
func (s *OrderService) validateCreateOrderRequest(req CreateOrderRequest) error {
	if req.UserID == 0 {
		return errors.New("用户ID不能为空")
	}
	if len(req.Items) == 0 {
		return errors.New("订单项不能为空")
	}
	if req.AddressID == 0 {
		return errors.New("收货地址不能为空")
	}
	
	for _, item := range req.Items {
		if item.ProductID == 0 {
			return errors.New("商品ID不能为空")
		}
		if item.Quantity <= 0 {
			return errors.New("商品数量必须大于0")
		}
	}
	
	return nil
}

// validateAndCalculateItems 验证商品并计算金额
func (s *OrderService) validateAndCalculateItems(tx *gorm.DB, items []OrderItemRequest) ([]OrderItem, float64, error) {
	var orderItems []OrderItem
	var totalAmount float64
	
	for _, item := range items {
		// 查询商品信息
		var product Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			return nil, 0, fmt.Errorf("商品ID %d 不存在", item.ProductID)
		}
		
		// 检查商品状态
		if product.Status != 1 {
			return nil, 0, fmt.Errorf("商品 %s 已下架", product.Name)
		}
		
		var price float64
		var skuInfo *ProductSKU
		var stock int
		
		// 如果指定了SKU
		if item.SKUID != nil {
			var sku ProductSKU
			if err := tx.Where("id = ? AND product_id = ?", *item.SKUID, item.ProductID).First(&sku).Error; err != nil {
				return nil, 0, fmt.Errorf("商品SKU不存在")
			}
			if sku.Status != 1 {
				return nil, 0, fmt.Errorf("商品SKU已下架")
			}
			price = sku.Price
			stock = sku.Stock
			skuInfo = &sku
		} else {
			price = product.Price
			stock = product.Stock
		}
		
		// 检查库存
		if stock < item.Quantity {
			return nil, 0, fmt.Errorf("商品 %s 库存不足，当前库存：%d", product.Name, stock)
		}
		
		// 计算金额
		itemTotal := price * float64(item.Quantity)
		totalAmount += itemTotal
		
		// 构建订单项
		orderItem := OrderItem{
			ProductID:    item.ProductID,
			SKUID:        item.SKUID,
			Quantity:     item.Quantity,
			Price:        price,
			TotalAmount:  itemTotal,
			ProductName:  product.Name,
			ProductSKU:   product.SKU,
		}
		
		if skuInfo != nil {
			orderItem.ProductSKU = skuInfo.SKU
			orderItem.ProductSpecs = skuInfo.Specs
		}
		
		// 获取商品主图
		var mainImage ProductImage
		if err := tx.Where("product_id = ? AND is_main = ?", item.ProductID, true).First(&mainImage).Error; err == nil {
			orderItem.ProductImage = mainImage.URL
		}
		
		orderItems = append(orderItems, orderItem)
	}
	
	return orderItems, totalAmount, nil
}

// validateAndUseCoupon 验证并使用优惠券
func (s *OrderService) validateAndUseCoupon(tx *gorm.DB, couponID *uint, userID uint, totalAmount float64) (float64, error) {
	if couponID == nil {
		return 0, nil
	}
	
	// 查询用户优惠券
	var userCoupon UserCoupon
	if err := tx.Preload("Coupon").Where("user_id = ? AND coupon_id = ? AND status = ?", userID, *couponID, 1).First(&userCoupon).Error; err != nil {
		return 0, errors.New("优惠券不存在或已使用")
	}
	
	coupon := userCoupon.Coupon
	
	// 检查优惠券有效期
	now := time.Now()
	if now.Before(coupon.StartAt) || now.After(coupon.EndAt) {
		return 0, errors.New("优惠券不在有效期内")
	}
	
	// 检查最低消费金额
	if totalAmount < coupon.MinAmount {
		return 0, fmt.Errorf("订单金额不满足优惠券使用条件，最低消费：%.2f", coupon.MinAmount)
	}
	
	// 计算折扣金额
	var discount float64
	switch coupon.Type {
	case 1: // 满减
		discount = coupon.Value
	case 2: // 折扣
		discount = totalAmount * (1 - coupon.Value/100)
		if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
			discount = coupon.MaxDiscount
		}
	case 3: // 固定金额
		discount = coupon.Value
	default:
		return 0, errors.New("无效的优惠券类型")
	}
	
	// 确保折扣不超过订单金额
	if discount > totalAmount {
		discount = totalAmount
	}
	
	// 标记优惠券为已使用
	now = time.Now()
	if err := tx.Model(&userCoupon).Updates(map[string]interface{}{
		"status":  2,
		"used_at": &now,
	}).Error; err != nil {
		return 0, err
	}
	
	// 更新优惠券使用次数
	if err := tx.Model(&coupon).Update("used", gorm.Expr("used + ?", 1)).Error; err != nil {
		return 0, err
	}
	
	return discount, nil
}

// calculateFreight 计算运费
func (s *OrderService) calculateFreight(totalAmount float64) float64 {
	// 简单的运费计算逻辑：满99免运费，否则10元
	if totalAmount >= 99 {
		return 0
	}
	return 10
}

// deductStock 扣减库存
func (s *OrderService) deductStock(tx *gorm.DB, items []OrderItemRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	for _, item := range items {
		if item.SKUID != nil {
			// 扣减SKU库存
			result := tx.Model(&ProductSKU{}).Where("id = ? AND stock >= ?", *item.SKUID, item.Quantity).Update("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("SKU库存不足")
			}
			
			// 同时扣减商品总库存
			tx.Model(&Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity))
		} else {
			// 扣减商品库存
			result := tx.Model(&Product{}).Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).Update("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("商品库存不足")
			}
		}
	}
	
	return nil
}

// clearCart 清空购物车
func (s *OrderService) clearCart(tx *gorm.DB, userID uint, items []OrderItemRequest) error {
	for _, item := range items {
		query := tx.Where("user_id = ? AND product_id = ?", userID, item.ProductID)
		if item.SKUID != nil {
			query = query.Where("sku_id = ?", *item.SKUID)
		}
		query.Delete(&Cart{})
	}
	return nil
}

// generateOrderNo 生成订单号
func (s *OrderService) generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 查询订单
		var order Order
		if err := tx.Preload("OrderItems").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
			return errors.New("订单不存在")
		}
		
		// 检查订单状态
		if order.Status != 1 {
			return errors.New("订单状态不允许取消")
		}
		
		// 更新订单状态
		if err := tx.Model(&order).Update("status", 5).Error; err != nil { // 5: 已取消
			return err
		}
		
		// 回滚库存
		for _, item := range order.OrderItems {
			if item.SKUID != nil {
				tx.Model(&ProductSKU{}).Where("id = ?", *item.SKUID).Update("stock", gorm.Expr("stock + ?", item.Quantity))
			}
			tx.Model(&Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity))
		}
		
		// 回滚优惠券
		if order.CouponID != nil {
			tx.Model(&UserCoupon{}).Where("user_id = ? AND coupon_id = ? AND status = ?", userID, *order.CouponID, 2).Updates(map[string]interface{}{
				"status":  1,
				"used_at": nil,
				"order_id": nil,
			})
			tx.Model(&Coupon{}).Where("id = ?", *order.CouponID).Update("used", gorm.Expr("used - ?", 1))
		}
		
		return nil
	})
}
```

### 练习任务

1. **完善业务逻辑**：补充支付、发货、确认收货等流程
2. **并发测试**：模拟高并发下单，测试库存扣减的准确性
3. **异常处理**：完善各种异常情况的处理逻辑
4. **性能优化**：优化数据库查询，减少事务时间
5. **单元测试**：为核心业务逻辑编写单元测试

---

## 📊 强化练习3：数据统计和报表

### 目标
实现复杂的数据统计和报表功能，掌握GORM的高级查询技巧。

### 要求
1. 实现销售数据统计（日、周、月、年）
2. 实现商品销量排行榜
3. 实现用户行为分析
4. 实现实时数据大屏
5. 优化复杂查询性能

### 代码实现

```go
package services

import (
	"time"
	"gorm.io/gorm"
)

// StatisticsService 统计服务
type StatisticsService struct {
	db *gorm.DB
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{db: db}
}

// SalesStatistics 销售统计
type SalesStatistics struct {
	Date        string  `json:"date"`
	OrderCount  int64   `json:"order_count"`
	SalesAmount float64 `json:"sales_amount"`
	UserCount   int64   `json:"user_count"`
}

// GetSalesStatistics 获取销售统计
func (s *StatisticsService) GetSalesStatistics(startDate, endDate time.Time, groupBy string) ([]SalesStatistics, error) {
	var stats []SalesStatistics
	
	// 根据groupBy确定日期格式
	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "DATE(created_at)"
	case "week":
		dateFormat = "DATE_FORMAT(created_at, '%Y-%u')"
	case "month":
		dateFormat = "DATE_FORMAT(created_at, '%Y-%m')"
	case "year":
		dateFormat = "DATE_FORMAT(created_at, '%Y')"
	default:
		dateFormat = "DATE(created_at)"
	}
	
	err := s.db.Model(&Order{}).
		Select(fmt.Sprintf(`
			%s as date,
			COUNT(*) as order_count,
			SUM(pay_amount) as sales_amount,
			COUNT(DISTINCT user_id) as user_count
		`, dateFormat)).
		Where("created_at BETWEEN ? AND ? AND status IN ?", startDate, endDate, []int{2, 3, 4}). // 已支付、已发货、已完成
		Group(dateFormat).
		Order("date ASC").
		Scan(&stats).Error
	
	return stats, err
}

// ProductSalesRank 商品销量排行
type ProductSalesRank struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CategoryName string  `json:"category_name"`
	SalesCount   int64   `json:"sales_count"`
	SalesAmount  float64 `json:"sales_amount"`
	Rank         int     `json:"rank"`
}

// GetProductSalesRank 获取商品销量排行
func (s *StatisticsService) GetProductSalesRank(startDate, endDate time.Time, limit int) ([]ProductSalesRank, error) {
	var ranks []ProductSalesRank
	
	err := s.db.Model(&OrderItem{}).
		Select(`
			order_items.product_id,
			products.name as product_name,
			categories.name as category_name,
			SUM(order_items.quantity) as sales_count,
			SUM(order_items.total_amount) as sales_amount
		`).
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN categories ON products.category_id = categories.id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status IN ?", startDate, endDate, []int{2, 3, 4}).
		Group("order_items.product_id, products.name, categories.name").
		Order("sales_count DESC").
		Limit(limit).
		Scan(&ranks).Error
	
	// 添加排名
	for i := range ranks {
		ranks[i].Rank = i + 1
	}
	
	return ranks, err
}

// UserBehaviorAnalysis 用户行为分析
type UserBehaviorAnalysis struct {
	UserID          uint    `json:"user_id"`
	Username        string  `json:"username"`
	OrderCount      int64   `json:"order_count"`
	TotalAmount     float64 `json:"total_amount"`
	AvgOrderAmount  float64 `json:"avg_order_amount"`
	LastOrderDate   *time.Time `json:"last_order_date"`
	DaysSinceLastOrder int  `json:"days_since_last_order"`
	CustomerLevel   string  `json:"customer_level"`
}

// GetUserBehaviorAnalysis 获取用户行为分析
func (s *StatisticsService) GetUserBehaviorAnalysis(startDate, endDate time.Time) ([]UserBehaviorAnalysis, error) {
	var analysis []UserBehaviorAnalysis
	
	err := s.db.Model(&Order{}).
		Select(`
			users.id as user_id,
			users.username,
			COUNT(orders.id) as order_count,
			SUM(orders.pay_amount) as total_amount,
			AVG(orders.pay_amount) as avg_order_amount,
			MAX(orders.created_at) as last_order_date,
			DATEDIFF(NOW(), MAX(orders.created_at)) as days_since_last_order
		`).
		Joins("JOIN users ON orders.user_id = users.id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status IN ?", startDate, endDate, []int{2, 3, 4}).
		Group("users.id, users.username").
		Having("order_count > 0").
		Order("total_amount DESC").
		Scan(&analysis).Error
	
	// 计算客户等级
	for i := range analysis {
		analysis[i].CustomerLevel = s.calculateCustomerLevel(analysis[i].TotalAmount, analysis[i].OrderCount)
	}
	
	return analysis, err
}

// calculateCustomerLevel 计算客户等级
func (s *StatisticsService) calculateCustomerLevel(totalAmount float64, orderCount int64) string {
	if totalAmount >= 10000 && orderCount >= 20 {
		return "钻石客户"
	} else if totalAmount >= 5000 && orderCount >= 10 {
		return "黄金客户"
	} else if totalAmount >= 1000 && orderCount >= 5 {
		return "银牌客户"
	} else {
		return "普通客户"
	}
}

// DashboardData 数据大屏
type DashboardData struct {
	TodayOrders     int64   `json:"today_orders"`
	TodaySales      float64 `json:"today_sales"`
	TodayUsers      int64   `json:"today_users"`
	TotalProducts   int64   `json:"total_products"`
	LowStockCount   int64   `json:"low_stock_count"`
	PendingOrders   int64   `json:"pending_orders"`
	MonthlyGrowth   float64 `json:"monthly_growth"`
	HourlyOrders    []HourlyOrderData `json:"hourly_orders"`
	CategoryStats   []CategorySalesData `json:"category_stats"`
	RecentOrders    []RecentOrderData `json:"recent_orders"`
}

// HourlyOrderData 每小时订单数据
type HourlyOrderData struct {
	Hour       int   `json:"hour"`
	OrderCount int64 `json:"order_count"`
}

// CategorySalesData 分类销售数据
type CategorySalesData struct {
	CategoryName string  `json:"category_name"`
	SalesAmount  float64 `json:"sales_amount"`
	OrderCount   int64   `json:"order_count"`
}

// RecentOrderData 最近订单数据
type RecentOrderData struct {
	OrderNo     string    `json:"order_no"`
	Username    string    `json:"username"`
	Amount      float64   `json:"amount"`
	Status      int8      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetDashboardData 获取数据大屏数据
func (s *StatisticsService) GetDashboardData() (*DashboardData, error) {
	var data DashboardData
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	
	// 今日订单数
	s.db.Model(&Order{}).Where("created_at BETWEEN ? AND ?", todayStart, todayEnd).Count(&data.TodayOrders)
	
	// 今日销售额
	s.db.Model(&Order{}).Where("created_at BETWEEN ? AND ? AND status IN ?", todayStart, todayEnd, []int{2, 3, 4}).Select("COALESCE(SUM(pay_amount), 0)").Scan(&data.TodaySales)
	
	// 今日新用户
	s.db.Model(&User{}).Where("created_at BETWEEN ? AND ?", todayStart, todayEnd).Count(&data.TodayUsers)
	
	// 商品总数
	s.db.Model(&Product{}).Where("status = ?", 1).Count(&data.TotalProducts)
	
	// 低库存商品数
	s.db.Model(&Product{}).Where("stock < ? AND status = ?", 10, 1).Count(&data.LowStockCount)
	
	// 待处理订单
	s.db.Model(&Order{}).Where("status = ?", 1).Count(&data.PendingOrders)
	
	// 月度增长率
	data.MonthlyGrowth = s.calculateMonthlyGrowth()
	
	// 每小时订单统计
	data.HourlyOrders = s.getHourlyOrders(todayStart, todayEnd)
	
	// 分类销售统计
	data.CategoryStats = s.getCategoryStats(todayStart, todayEnd)
	
	// 最近订单
	data.RecentOrders = s.getRecentOrders(10)
	
	return &data, nil
}

// calculateMonthlyGrowth 计算月度增长率
func (s *StatisticsService) calculateMonthlyGrowth() float64 {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := currentMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := currentMonthStart.Add(-time.Second)
	
	var currentMonthSales, lastMonthSales float64
	
	s.db.Model(&Order{}).Where("created_at >= ? AND status IN ?", currentMonthStart, []int{2, 3, 4}).Select("COALESCE(SUM(pay_amount), 0)").Scan(&currentMonthSales)
	s.db.Model(&Order{}).Where("created_at BETWEEN ? AND ? AND status IN ?", lastMonthStart, lastMonthEnd, []int{2, 3, 4}).Select("COALESCE(SUM(pay_amount), 0)").Scan(&lastMonthSales)
	
	if lastMonthSales == 0 {
		return 0
	}
	
	return ((currentMonthSales - lastMonthSales) / lastMonthSales) * 100
}

// getHourlyOrders 获取每小时订单统计
func (s *StatisticsService) getHourlyOrders(start, end time.Time) []HourlyOrderData {
	var hourlyData []HourlyOrderData
	
	s.db.Model(&Order{}).
		Select("HOUR(created_at) as hour, COUNT(*) as order_count").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("HOUR(created_at)").
		Order("hour ASC").
		Scan(&hourlyData)
	
	// 补充缺失的小时数据
	hourlyMap := make(map[int]int64)
	for _, data := range hourlyData {
		hourlyMap[data.Hour] = data.OrderCount
	}
	
	result := make([]HourlyOrderData, 24)
	for i := 0; i < 24; i++ {
		result[i] = HourlyOrderData{
			Hour:       i,
			OrderCount: hourlyMap[i],
		}
	}
	
	return result
}

// getCategoryStats 获取分类销售统计
func (s *StatisticsService) getCategoryStats(start, end time.Time) []CategorySalesData {
	var categoryStats []CategorySalesData
	
	s.db.Model(&OrderItem{}).
		Select(`
			categories.name as category_name,
			SUM(order_items.total_amount) as sales_amount,
			COUNT(DISTINCT order_items.order_id) as order_count
		`).
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN categories ON products.category_id = categories.id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status IN ?", start, end, []int{2, 3, 4}).
		Group("categories.id, categories.name").
		Order("sales_amount DESC").
		Limit(10).
		Scan(&categoryStats)
	
	return categoryStats
}

// getRecentOrders 获取最近订单
func (s *StatisticsService) getRecentOrders(limit int) []RecentOrderData {
	var recentOrders []RecentOrderData
	
	s.db.Model(&Order{}).
		Select(`
			orders.order_no,
			users.username,
			orders.pay_amount as amount,
			orders.status,
			orders.created_at
		`).
		Joins("JOIN users ON orders.user_id = users.id").
		Order("orders.created_at DESC").
		Limit(limit).
		Scan(&recentOrders)
	
	return recentOrders
}
```

### 练习任务

1. **复杂查询优化**：分析慢查询，添加合适的索引
2. **数据可视化**：将统计数据用图表展示
3. **实时更新**：实现数据的实时刷新机制
4. **缓存优化**：对频繁查询的数据进行缓存
5. **导出功能**：实现数据导出为Excel功能

---

## 🔧 强化练习4：性能优化和监控

### 目标
掌握GORM的性能优化技巧，实现数据库监控和调优。

### 要求
1. 实现数据库连接池优化
2. 实现查询性能监控
3. 实现慢查询日志
4. 实现数据库读写分离
5. 实现查询缓存机制

### 代码实现

```go
package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/mysql"
	"github.com/go-redis/redis/v8"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Master DatabaseConnection `json:"master"`
	Slaves []DatabaseConnection `json:"slaves"`
	Pool   PoolConfig `json:"pool"`
	Cache  CacheConfig `json:"cache"`
}

type DatabaseConnection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

type PoolConfig struct {
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

type CacheConfig struct {
	Enabled bool          `json:"enabled"`
	TTL     time.Duration `json:"ttl"`
	Prefix  string        `json:"prefix"`
}

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	master *gorm.DB
	slaves []*gorm.DB
	redis  *redis.Client
	config DatabaseConfig
	metrics *DatabaseMetrics
}

// DatabaseMetrics 数据库指标
type DatabaseMetrics struct {
	QueryCount    int64         `json:"query_count"`
	SlowQueryCount int64        `json:"slow_query_count"`
	AvgQueryTime  time.Duration `json:"avg_query_time"`
	ErrorCount    int64         `json:"error_count"`
	CacheHitRate  float64       `json:"cache_hit_rate"`
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(config DatabaseConfig) (*DatabaseManager, error) {
	manager := &DatabaseManager{
		config:  config,
		metrics: &DatabaseMetrics{},
	}
	
	// 初始化主库
	master, err := manager.initDatabase(config.Master, "master")
	if err != nil {
		return nil, fmt.Errorf("初始化主库失败: %v", err)
	}
	manager.master = master
	
	// 初始化从库
	for i, slaveConfig := range config.Slaves {
		slave, err := manager.initDatabase(slaveConfig, fmt.Sprintf("slave-%d", i))
		if err != nil {
			log.Printf("初始化从库 %d 失败: %v", i, err)
			continue
		}
		manager.slaves = append(manager.slaves, slave)
	}
	
	// 初始化Redis缓存
	if config.Cache.Enabled {
		manager.redis = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		})
	}
	
	return manager, nil
}

// initDatabase 初始化数据库连接
func (dm *DatabaseManager) initDatabase(config DatabaseConnection, role string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset)
	
	// 自定义日志记录器
	customLogger := logger.New(
		log.New(log.Writer(), fmt.Sprintf("[%s] ", role), log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, err
	}
	
	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	sqlDB.SetMaxOpenConns(dm.config.Pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dm.config.Pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dm.config.Pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dm.config.Pool.ConnMaxIdleTime)
	
	// 添加性能监控插件
	db.Use(&PerformancePlugin{
		metrics: dm.metrics,
	})
	
	return db, nil
}

// GetMaster 获取主库连接（用于写操作）
func (dm *DatabaseManager) GetMaster() *gorm.DB {
	return dm.master
}

// GetSlave 获取从库连接（用于读操作）
func (dm *DatabaseManager) GetSlave() *gorm.DB {
	if len(dm.slaves) == 0 {
		return dm.master
	}
	
	// 简单的轮询负载均衡
	index := time.Now().UnixNano() % int64(len(dm.slaves))
	return dm.slaves[index]
}

// PerformancePlugin 性能监控插件
type PerformancePlugin struct {
	metrics *DatabaseMetrics
}

func (p *PerformancePlugin) Name() string {
	return "performance"
}

func (p *PerformancePlugin) Initialize(db *gorm.DB) error {
	// 注册回调函数
	db.Callback().Query().Before("gorm:query").Register("performance:before_query", p.beforeQuery)
	db.Callback().Query().After("gorm:query").Register("performance:after_query", p.afterQuery)
	
	db.Callback().Create().Before("gorm:create").Register("performance:before_create", p.beforeQuery)
	db.Callback().Create().After("gorm:create").Register("performance:after_create", p.afterQuery)
	
	db.Callback().Update().Before("gorm:update").Register("performance:before_update", p.beforeQuery)
	db.Callback().Update().After("gorm:update").Register("performance:after_update", p.afterQuery)
	
	db.Callback().Delete().Before("gorm:delete").Register("performance:before_delete", p.beforeQuery)
	db.Callback().Delete().After("gorm:delete").Register("performance:after_delete", p.afterQuery)
	
	return nil
}

func (p *PerformancePlugin) beforeQuery(db *gorm.DB) {
	db.Set("start_time", time.Now())
}

func (p *PerformancePlugin) afterQuery(db *gorm.DB) {
	startTime, exists := db.Get("start_time")
	if !exists {
		return
	}
	
	duration := time.Since(startTime.(time.Time))
	p.metrics.QueryCount++
	
	// 记录慢查询
	if duration > 200*time.Millisecond {
		p.metrics.SlowQueryCount++
		log.Printf("慢查询检测: SQL=%s, 耗时=%v", db.Statement.SQL.String(), duration)
	}
	
	// 更新平均查询时间
	p.updateAvgQueryTime(duration)
	
	// 记录错误
	if db.Error != nil {
		p.metrics.ErrorCount++
		log.Printf("数据库错误: %v, SQL=%s", db.Error, db.Statement.SQL.String())
	}
}

func (p *PerformancePlugin) updateAvgQueryTime(duration time.Duration) {
	// 简单的移动平均算法
	if p.metrics.AvgQueryTime == 0 {
		p.metrics.AvgQueryTime = duration
	} else {
		p.metrics.AvgQueryTime = (p.metrics.AvgQueryTime + duration) / 2
	}
}

// CacheService 缓存服务
type CacheService struct {
	redis  *redis.Client
	config CacheConfig
}

// NewCacheService 创建缓存服务
func NewCacheService(redis *redis.Client, config CacheConfig) *CacheService {
	return &CacheService{
		redis:  redis,
		config: config,
	}
}

// Get 获取缓存
func (cs *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := cs.config.Prefix + key
	result, err := cs.redis.Get(ctx, fullKey).Result()
	if err != nil {
		return err
	}
	
	// 这里应该使用JSON反序列化
	// 简化示例，实际应该使用json.Unmarshal
	return nil
}

// Set 设置缓存
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}) error {
	fullKey := cs.config.Prefix + key
	// 这里应该使用JSON序列化
	// 简化示例，实际应该使用json.Marshal
	return cs.redis.Set(ctx, fullKey, value, cs.config.TTL).Err()
}

// Delete 删除缓存
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	fullKey := cs.config.Prefix + key
	return cs.redis.Del(ctx, fullKey).Err()
}
```

### 练习任务

1. **连接池调优**：根据业务负载调整连接池参数
2. **监控告警**：实现数据库性能监控和告警机制
3. **查询优化**：分析并优化慢查询
4. **缓存策略**：设计合理的缓存策略和失效机制
5. **压力测试**：进行数据库压力测试，验证优化效果

---

## 🎯 综合实战项目

### 项目要求

基于前面的练习，完成一个完整的电商系统后台管理功能：

1. **用户管理模块**
   - 用户列表、详情、编辑
   - 用户行为分析
   - 用户等级管理

2. **商品管理模块**
   - 商品CRUD操作
   - 批量导入/导出
   - 库存管理
   - 价格管理

3. **订单管理模块**
   - 订单列表、详情
   - 订单状态流转
   - 退款处理
   - 物流跟踪

4. **营销管理模块**
   - 优惠券管理
   - 促销活动
   - 积分系统

5. **数据统计模块**
   - 销售报表
   - 用户分析
   - 商品分析
   - 实时大屏

### 技术要求

1. 使用GORM实现所有数据库操作
2. 实现读写分离和缓存机制
3. 添加完整的错误处理和日志记录
4. 实现数据库事务和并发控制
5. 添加性能监控和优化
6. 编写完整的单元测试
7. 提供API文档和部署说明

### 评估标准

1. **功能完整性**（30%）：所有功能模块是否完整实现
2. **代码质量**（25%）：代码结构、命名规范、注释完整性
3. **性能优化**（20%）：查询优化、缓存使用、并发处理
4. **错误处理**（15%）：异常处理、数据验证、日志记录
5. **测试覆盖**（10%）：单元测试、集成测试覆盖率

---

## 📝 学习总结

### 核心知识点回顾

1. **复杂数据模型设计**
   - 多表关联关系设计
   - 索引和约束优化
   - 软删除和审计字段

2. **业务逻辑实现**
   - 事务处理和数据一致性
   - 并发控制和库存管理
   - 复杂业务流程设计

3. **性能优化技巧**
   - 查询优化和索引使用
   - 连接池配置和监控
   - 缓存策略和读写分离

4. **数据统计分析**
   - 复杂聚合查询
   - 报表生成和可视化
   - 实时数据处理

### 进阶学习建议

1. **深入学习数据库原理**
   - MySQL存储引擎
   - 索引原理和优化
   - 事务隔离级别

2. **微服务架构**
   - 数据库拆分策略
   - 分布式事务处理
   - 服务间数据同步

3. **大数据处理**
   - 数据分片和分库分表
   - 数据仓库设计
   - 实时数据流处理

4. **云原生技术**
   - 容器化部署
   - 自动扩缩容
   - 监控和运维

### 实践建议

1. **多做项目实战**：通过实际项目加深理解
2. **关注性能优化**：始终考虑性能和扩展性
3. **学习最佳实践**：参考开源项目和行业标准
4. **持续学习新技术**：跟上技术发展趋势

---

## 🔗 相关资源

- [GORM官方文档](https://gorm.io/docs/)
- [MySQL性能优化指南](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)
- [Go并发编程实战](https://golang.org/doc/effective_go.html#concurrency)
- [微服务设计模式](https://microservices.io/patterns/)

---

**恭喜你完成了GORM强化练习！** 🎉

通过这些综合性的练习，你应该已经掌握了GORM的高级特性和企业级应用开发技能。继续保持学习的热情，在实际项目中不断实践和优化！
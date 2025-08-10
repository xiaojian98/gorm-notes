package main

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone       string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Nickname    string         `gorm:"size:50" json:"nickname"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Gender      int8           `gorm:"default:0;comment:0-未知,1-男,2-女" json:"gender"`
	Birthday    *time.Time     `json:"birthday"`
	Status      int8           `gorm:"default:1;comment:1-正常,2-禁用" json:"status"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	
	// 关联关系
	Profile   *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Addresses []Address    `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Orders    []Order      `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Carts     []Cart       `gorm:"foreignKey:UserID" json:"carts,omitempty"`
	Coupons   []UserCoupon `gorm:"foreignKey:UserID" json:"coupons,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserProfile 用户资料
type UserProfile struct {
	BaseModel
	UserID      uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	RealName    string `gorm:"size:50" json:"real_name"`
	IDCard      string `gorm:"size:20" json:"id_card"`
	Company     string `gorm:"size:100" json:"company"`
	Position    string `gorm:"size:50" json:"position"`
	Address     string `gorm:"size:255" json:"address"`
	Introduction string `gorm:"type:text" json:"introduction"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// Address 收货地址
type Address struct {
	BaseModel
	UserID     uint   `gorm:"index;not null" json:"user_id"`
	Name       string `gorm:"size:50;not null" json:"name"`
	Phone      string `gorm:"size:20;not null" json:"phone"`
	Province   string `gorm:"size:50;not null" json:"province"`
	City       string `gorm:"size:50;not null" json:"city"`
	District   string `gorm:"size:50;not null" json:"district"`
	Detail     string `gorm:"size:255;not null" json:"detail"`
	PostalCode string `gorm:"size:10" json:"postal_code"`
	IsDefault  bool   `gorm:"default:false" json:"is_default"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Address) TableName() string {
	return "addresses"
}

// Category 商品分类
type Category struct {
	BaseModel
	Name        string `gorm:"size:50;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Image       string `gorm:"size:255" json:"image"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联关系
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product  `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Brand 品牌
type Brand struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Logo        string `gorm:"size:255" json:"logo"`
	Website     string `gorm:"size:255" json:"website"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联关系
	Products []Product `gorm:"foreignKey:BrandID" json:"products,omitempty"`
}

// TableName 指定表名
func (Brand) TableName() string {
	return "brands"
}

// Product 商品
type Product struct {
	BaseModel
	Name         string          `gorm:"size:255;not null" json:"name"`
	SKU          string          `gorm:"uniqueIndex;size:100;not null" json:"sku"`
	Description  string          `gorm:"type:text" json:"description"`
	Content      string          `gorm:"type:longtext" json:"content"`
	CategoryID   uint            `gorm:"index;not null" json:"category_id"`
	BrandID      *uint           `gorm:"index" json:"brand_id"`
	Price        int64           `gorm:"not null;comment:价格(分)" json:"price"`
	MarketPrice  int64           `gorm:"comment:市场价(分)" json:"market_price"`
	CostPrice    int64           `gorm:"comment:成本价(分)" json:"cost_price"`
	Stock        int             `gorm:"default:0" json:"stock"`
	Sales        int             `gorm:"default:0" json:"sales"`
	Views        int             `gorm:"default:0" json:"views"`
	Weight       float64         `gorm:"comment:重量(kg)" json:"weight"`
	Volume       float64         `gorm:"comment:体积(立方米)" json:"volume"`
	Keywords     string          `gorm:"size:255" json:"keywords"`
	Tags         json.RawMessage `gorm:"type:json" json:"tags"`
	Attributes   json.RawMessage `gorm:"type:json" json:"attributes"`
	Status       int8            `gorm:"default:1;comment:1-上架,2-下架" json:"status"`
	Sort         int             `gorm:"default:0" json:"sort"`
	
	// 关联关系
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Brand        *Brand         `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Images       []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	SKUs         []ProductSKU   `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	Reviews      []ProductReview `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
	OrderItems   []OrderItem    `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`
	CartItems    []Cart         `gorm:"foreignKey:ProductID" json:"cart_items,omitempty"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// ProductImage 商品图片
type ProductImage struct {
	BaseModel
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	URL       string `gorm:"size:255;not null" json:"url"`
	Alt       string `gorm:"size:255" json:"alt"`
	Sort      int    `gorm:"default:0" json:"sort"`
	IsMain    bool   `gorm:"default:false" json:"is_main"`
	
	// 关联关系
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (ProductImage) TableName() string {
	return "product_images"
}

// ProductSKU 商品SKU
type ProductSKU struct {
	BaseModel
	ProductID uint            `gorm:"index;not null" json:"product_id"`
	SKU       string          `gorm:"uniqueIndex;size:100;not null" json:"sku"`
	Name      string          `gorm:"size:255" json:"name"`
	Image     string          `gorm:"size:255" json:"image"`
	Price     int64           `gorm:"not null;comment:价格(分)" json:"price"`
	Stock     int             `gorm:"default:0" json:"stock"`
	Sales     int             `gorm:"default:0" json:"sales"`
	Weight    float64         `gorm:"comment:重量(kg)" json:"weight"`
	Specs     json.RawMessage `gorm:"type:json;comment:规格参数" json:"specs"`
	Status    int8            `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联关系
	Product    Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:SKUID" json:"order_items,omitempty"`
	CartItems  []Cart      `gorm:"foreignKey:SKUID" json:"cart_items,omitempty"`
}

// TableName 指定表名
func (ProductSKU) TableName() string {
	return "product_skus"
}

// ProductReview 商品评价
type ProductReview struct {
	BaseModel
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	UserID    uint   `gorm:"index;not null" json:"user_id"`
	OrderID   uint   `gorm:"index;not null" json:"order_id"`
	Rating    int8   `gorm:"not null;comment:评分1-5" json:"rating"`
	Content   string `gorm:"type:text" json:"content"`
	Images    json.RawMessage `gorm:"type:json" json:"images"`
	Reply     string `gorm:"type:text" json:"reply"`
	Status    int8   `gorm:"default:1;comment:1-显示,2-隐藏" json:"status"`
	
	// 关联关系
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order   Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (ProductReview) TableName() string {
	return "product_reviews"
}

// Cart 购物车
type Cart struct {
	BaseModel
	UserID    uint  `gorm:"index;not null" json:"user_id"`
	ProductID uint  `gorm:"index;not null" json:"product_id"`
	SKUID     *uint `gorm:"index" json:"sku_id"`
	Quantity  int   `gorm:"not null" json:"quantity"`
	
	// 关联关系
	User    User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU     *ProductSKU `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// TableName 指定表名
func (Cart) TableName() string {
	return "carts"
}

// Order 订单
type Order struct {
	BaseModel
	OrderNo         string     `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	Status          int8       `gorm:"index;default:1;comment:1-待付款,2-待发货,3-待收货,4-已完成,5-已取消" json:"status"`
	TotalAmount     int64      `gorm:"not null;comment:商品总金额(分)" json:"total_amount"`
	PayAmount       int64      `gorm:"not null;comment:实付金额(分)" json:"pay_amount"`
	FreightAmount   int64      `gorm:"default:0;comment:运费(分)" json:"freight_amount"`
	DiscountAmount  int64      `gorm:"default:0;comment:优惠金额(分)" json:"discount_amount"`
	CouponID        *uint      `gorm:"index" json:"coupon_id"`
	PaymentMethod   string     `gorm:"size:50" json:"payment_method"`
	PaymentNo       string     `gorm:"size:100" json:"payment_no"`
	ReceiverName    string     `gorm:"size:50;not null" json:"receiver_name"`
	ReceiverPhone   string     `gorm:"size:20;not null" json:"receiver_phone"`
	ReceiverAddress string     `gorm:"size:255;not null" json:"receiver_address"`
	Remark          string     `gorm:"type:text" json:"remark"`
	PaidAt          *time.Time `json:"paid_at"`
	ShippedAt       *time.Time `json:"shipped_at"`
	DeliveredAt     *time.Time `json:"delivered_at"`
	FinishedAt      *time.Time `json:"finished_at"`
	CancelTime      *time.Time `json:"cancel_time"`
	CancelReason    string     `gorm:"type:text" json:"cancel_reason"`
	
	// 关联关系
	User     User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coupon   *Coupon     `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payments []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
	Reviews  []ProductReview `gorm:"foreignKey:OrderID" json:"reviews,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项
type OrderItem struct {
	BaseModel
	OrderID      uint            `gorm:"index;not null" json:"order_id"`
	ProductID    uint            `gorm:"index;not null" json:"product_id"`
	SKUID        *uint           `gorm:"index" json:"sku_id"`
	Quantity     int             `gorm:"not null" json:"quantity"`
	Price        int64           `gorm:"not null;comment:单价(分)" json:"price"`
	TotalPrice   int64           `gorm:"not null;comment:总价(分)" json:"total_price"`
	ProductName  string          `gorm:"size:255;not null" json:"product_name"`
	ProductSKU   string          `gorm:"size:100" json:"product_sku"`
	ProductImage string          `gorm:"size:255" json:"product_image"`
	ProductSpecs json.RawMessage `gorm:"type:json" json:"product_specs"`
	
	// 关联关系
	Order   Order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU     *ProductSKU `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// Payment 支付记录
type Payment struct {
	BaseModel
	OrderID       uint       `gorm:"index;not null" json:"order_id"`
	PaymentNo     string     `gorm:"uniqueIndex;size:100;not null" json:"payment_no"`
	Method        string     `gorm:"size:50;not null" json:"method"`
	Amount        int64      `gorm:"not null;comment:支付金额(分)" json:"amount"`
	Status        int8       `gorm:"default:1;comment:1-待支付,2-支付成功,3-支付失败" json:"status"`
	ThirdPartyNo  string     `gorm:"size:100" json:"third_party_no"`
	ThirdPartyData json.RawMessage `gorm:"type:json" json:"third_party_data"`
	PaidAt        *time.Time `json:"paid_at"`
	FailedReason  string     `gorm:"type:text" json:"failed_reason"`
	
	// 关联关系
	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

// Coupon 优惠券
type Coupon struct {
	BaseModel
	Name         string    `gorm:"size:100;not null" json:"name"`
	Code         string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Type         int8      `gorm:"not null;comment:1-满减,2-折扣,3-固定金额" json:"type"`
	Value        int64     `gorm:"not null;comment:优惠值" json:"value"`
	MinAmount    int64     `gorm:"default:0;comment:最低消费金额(分)" json:"min_amount"`
	MaxDiscount  int64     `gorm:"default:0;comment:最大优惠金额(分)" json:"max_discount"`
	TotalQuantity int      `gorm:"not null;comment:总数量" json:"total_quantity"`
	UsedQuantity  int      `gorm:"default:0;comment:已使用数量" json:"used_quantity"`
	PerUserLimit  int      `gorm:"default:1;comment:每人限领数量" json:"per_user_limit"`
	StartTime     time.Time `gorm:"not null" json:"start_time"`
	EndTime       time.Time `gorm:"not null" json:"end_time"`
	Description   string    `gorm:"type:text" json:"description"`
	Status        int8      `gorm:"default:1;comment:1-启用,2-禁用" json:"status"`
	
	// 关联关系
	UserCoupons []UserCoupon `gorm:"foreignKey:CouponID" json:"user_coupons,omitempty"`
	Orders      []Order      `gorm:"foreignKey:CouponID" json:"orders,omitempty"`
}

// TableName 指定表名
func (Coupon) TableName() string {
	return "coupons"
}

// UserCoupon 用户优惠券
type UserCoupon struct {
	BaseModel
	UserID   uint       `gorm:"index;not null" json:"user_id"`
	CouponID uint       `gorm:"index;not null" json:"coupon_id"`
	Status   int8       `gorm:"default:1;comment:1-未使用,2-已使用,3-已过期" json:"status"`
	UsedAt   *time.Time `json:"used_at"`
	
	// 关联关系
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coupon Coupon `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

// TableName 指定表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}
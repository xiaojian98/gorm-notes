package main

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy uint           `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy uint           `gorm:"index" json:"updated_by,omitempty"`
}

// User 用户模型
type User struct {
	BaseModel
	Email        string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone        string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Gender       int8           `gorm:"default:0;comment:0-未知,1-男,2-女" json:"gender"`
	Birthday     *time.Time     `json:"birthday"`
	Status       int8           `gorm:"default:1;comment:0-禁用,1-正常" json:"status"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	Points       int            `gorm:"default:0" json:"points"`
	
	// 关联关系
	Profile      *UserProfile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Addresses    []Address      `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Orders       []Order        `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	Carts        []Cart         `gorm:"foreignKey:UserID" json:"carts,omitempty"`
	Coupons      []Coupon       `gorm:"many2many:user_coupons" json:"coupons,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserProfile 用户详细资料
type UserProfile struct {
	BaseModel
	UserID      uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	RealName    string `gorm:"size:50" json:"real_name"`
	IDCard      string `gorm:"size:18" json:"id_card"`
	Company     string `gorm:"size:100" json:"company"`
	Position    string `gorm:"size:50" json:"position"`
	Bio         string `gorm:"type:text" json:"bio"`
	Website     string `gorm:"size:255" json:"website"`
	Location    string `gorm:"size:100" json:"location"`
	
	// 关联关系
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// Address 地址模型
type Address struct {
	BaseModel
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	Name        string `gorm:"size:50;not null" json:"name"`
	Phone       string `gorm:"size:20;not null" json:"phone"`
	Province    string `gorm:"size:50;not null" json:"province"`
	City        string `gorm:"size:50;not null" json:"city"`
	District    string `gorm:"size:50;not null" json:"district"`
	Detail      string `gorm:"size:255;not null" json:"detail"`
	Postcode    string `gorm:"size:10" json:"postcode"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	
	// 关联关系
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Address) TableName() string {
	return "addresses"
}

// Category 商品分类模型
type Category struct {
	BaseModel
	Name        string      `gorm:"size:50;not null" json:"name"`
	Slug        string      `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string      `gorm:"type:text" json:"description"`
	Image       string      `gorm:"size:255" json:"image"`
	ParentID    *uint       `gorm:"index" json:"parent_id"`
	Sort        int         `gorm:"default:0" json:"sort"`
	Status      int8        `gorm:"default:1;comment:0-禁用,1-正常" json:"status"`
	
	// 自关联
	Parent      *Category   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	
	// 关联关系
	Products    []Product   `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Brand 品牌模型
type Brand struct {
	BaseModel
	Name        string    `gorm:"size:50;not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Logo        string    `gorm:"size:255" json:"logo"`
	Website     string    `gorm:"size:255" json:"website"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Status      int8      `gorm:"default:1;comment:0-禁用,1-正常" json:"status"`
	
	// 关联关系
	Products    []Product `gorm:"foreignKey:BrandID" json:"products,omitempty"`
}

// TableName 指定表名
func (Brand) TableName() string {
	return "brands"
}

// Product 商品模型
type Product struct {
	BaseModel
	Name         string         `gorm:"size:255;not null" json:"name"`
	SKU          string         `gorm:"uniqueIndex;size:100;not null" json:"sku"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        int64          `gorm:"not null;comment:价格(分)" json:"price"`
	MarketPrice  int64          `gorm:"comment:市场价(分)" json:"market_price"`
	CostPrice    int64          `gorm:"comment:成本价(分)" json:"cost_price"`
	Stock        int            `gorm:"default:0" json:"stock"`
	Sales        int            `gorm:"default:0" json:"sales"`
	Views        int            `gorm:"default:0" json:"views"`
	Weight       float64        `gorm:"comment:重量(kg)" json:"weight"`
	Volume       float64        `gorm:"comment:体积(m³)" json:"volume"`
	Status       int8           `gorm:"default:1;comment:0-下架,1-上架" json:"status"`
	CategoryID   uint           `gorm:"index;not null" json:"category_id"`
	BrandID      uint           `gorm:"index;not null" json:"brand_id"`
	
	// 关联关系
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Brand        Brand          `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Images       []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Attrs        []ProductAttr  `gorm:"foreignKey:ProductID" json:"attrs,omitempty"`
	SKUs         []ProductSKU   `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	Reviews      []ProductReview `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// ProductImage 商品图片模型
type ProductImage struct {
	BaseModel
	ProductID   uint   `gorm:"index;not null" json:"product_id"`
	URL         string `gorm:"size:255;not null" json:"url"`
	Alt         string `gorm:"size:255" json:"alt"`
	Sort        int    `gorm:"default:0" json:"sort"`
	IsMain      bool   `gorm:"default:false" json:"is_main"`
	
	// 关联关系
	Product     Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (ProductImage) TableName() string {
	return "product_images"
}

// ProductAttr 商品属性模型
type ProductAttr struct {
	BaseModel
	ProductID   uint   `gorm:"index;not null" json:"product_id"`
	Name        string `gorm:"size:50;not null" json:"name"`
	Value       string `gorm:"size:255;not null" json:"value"`
	Sort        int    `gorm:"default:0" json:"sort"`
	
	// 关联关系
	Product     Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (ProductAttr) TableName() string {
	return "product_attrs"
}

// ProductSKU 商品SKU模型
type ProductSKU struct {
	BaseModel
	ProductID   uint            `gorm:"index;not null" json:"product_id"`
	SKU         string          `gorm:"uniqueIndex;size:100;not null" json:"sku"`
	Specs       json.RawMessage `gorm:"type:json;comment:规格信息" json:"specs"`
	Price       int64           `gorm:"not null;comment:价格(分)" json:"price"`
	Stock       int             `gorm:"default:0" json:"stock"`
	Sales       int             `gorm:"default:0" json:"sales"`
	Image       string          `gorm:"size:255" json:"image"`
	Status      int8            `gorm:"default:1;comment:0-禁用,1-正常" json:"status"`
	
	// 关联关系
	Product     Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (ProductSKU) TableName() string {
	return "product_skus"
}

// ProductReview 商品评价模型
type ProductReview struct {
	BaseModel
	ProductID   uint   `gorm:"index;not null" json:"product_id"`
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	OrderID     uint   `gorm:"index;not null" json:"order_id"`
	Rating      int8   `gorm:"not null;comment:评分1-5" json:"rating"`
	Content     string `gorm:"type:text" json:"content"`
	Images      string `gorm:"type:json;comment:评价图片" json:"images"`
	IsAnonymous bool   `gorm:"default:false" json:"is_anonymous"`
	Status      int8   `gorm:"default:1;comment:0-隐藏,1-显示" json:"status"`
	
	// 关联关系
	Product     Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User        User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Order       Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (ProductReview) TableName() string {
	return "product_reviews"
}

// Cart 购物车模型
type Cart struct {
	BaseModel
	UserID      uint        `gorm:"index;not null" json:"user_id"`
	ProductID   uint        `gorm:"index;not null" json:"product_id"`
	SKUID       *uint       `gorm:"index" json:"sku_id"`
	Quantity    int         `gorm:"not null" json:"quantity"`
	
	// 关联关系
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product     Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU         *ProductSKU `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// TableName 指定表名
func (Cart) TableName() string {
	return "carts"
}

// Order 订单模型
type Order struct {
	BaseModel
	OrderNo         string      `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	UserID          uint        `gorm:"index;not null" json:"user_id"`
	Status          int8        `gorm:"index;default:1;comment:1-待付款,2-待发货,3-待收货,4-已完成,5-已取消" json:"status"`
	TotalAmount     int64       `gorm:"not null;comment:总金额(分)" json:"total_amount"`
	PayAmount       int64       `gorm:"not null;comment:实付金额(分)" json:"pay_amount"`
	FreightAmount   int64       `gorm:"default:0;comment:运费(分)" json:"freight_amount"`
	DiscountAmount  int64       `gorm:"default:0;comment:优惠金额(分)" json:"discount_amount"`
	CouponID        *uint       `gorm:"index" json:"coupon_id"`
	PaymentMethod   string      `gorm:"size:20" json:"payment_method"`
	PaymentTime     *time.Time  `json:"payment_time"`
	ShippingTime    *time.Time  `json:"shipping_time"`
	DeliveryTime    *time.Time  `json:"delivery_time"`
	FinishTime      *time.Time  `json:"finish_time"`
	CancelTime      *time.Time  `json:"cancel_time"`
	CancelReason    string      `gorm:"size:255" json:"cancel_reason"`
	Remark          string      `gorm:"type:text" json:"remark"`
	
	// 收货地址信息（冗余存储）
	ReceiverName    string      `gorm:"size:50;not null" json:"receiver_name"`
	ReceiverPhone   string      `gorm:"size:20;not null" json:"receiver_phone"`
	ReceiverAddress string      `gorm:"size:500;not null" json:"receiver_address"`
	
	// 关联关系
	User            User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coupon          *Coupon     `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payments        []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	BaseModel
	OrderID     uint        `gorm:"index;not null" json:"order_id"`
	ProductID   uint        `gorm:"index;not null" json:"product_id"`
	SKUID       *uint       `gorm:"index" json:"sku_id"`
	Quantity    int         `gorm:"not null" json:"quantity"`
	Price       int64       `gorm:"not null;comment:单价(分)" json:"price"`
	TotalPrice  int64       `gorm:"not null;comment:总价(分)" json:"total_price"`
	
	// 商品信息快照（冗余存储）
	ProductName string      `gorm:"size:255;not null" json:"product_name"`
	ProductSKU  string      `gorm:"size:100;not null" json:"product_sku"`
	ProductImage string     `gorm:"size:255" json:"product_image"`
	ProductSpecs json.RawMessage `gorm:"type:json;comment:商品规格" json:"product_specs"`
	
	// 关联关系
	Order       Order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product     Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU         *ProductSKU `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}

// Payment 支付记录模型
type Payment struct {
	BaseModel
	OrderID         uint       `gorm:"index;not null" json:"order_id"`
	PaymentNo       string     `gorm:"uniqueIndex;size:32;not null" json:"payment_no"`
	Method          string     `gorm:"size:20;not null" json:"method"`
	Amount          int64      `gorm:"not null;comment:支付金额(分)" json:"amount"`
	Status          int8       `gorm:"default:1;comment:1-待支付,2-支付成功,3-支付失败" json:"status"`
	ThirdPartyNo    string     `gorm:"size:64" json:"third_party_no"`
	ThirdPartyData  string     `gorm:"type:json" json:"third_party_data"`
	PaidAt          *time.Time `json:"paid_at"`
	FailReason      string     `gorm:"size:255" json:"fail_reason"`
	
	// 关联关系
	Order           Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

// Coupon 优惠券模型
type Coupon struct {
	BaseModel
	Name            string     `gorm:"size:100;not null" json:"name"`
	Code            string     `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Type            int8       `gorm:"not null;comment:1-满减,2-折扣,3-固定金额" json:"type"`
	Value           int64      `gorm:"not null;comment:优惠值" json:"value"`
	MinAmount       int64      `gorm:"default:0;comment:最低消费金额(分)" json:"min_amount"`
	MaxDiscount     int64      `gorm:"default:0;comment:最大优惠金额(分)" json:"max_discount"`
	TotalQuantity   int        `gorm:"not null" json:"total_quantity"`
	UsedQuantity    int        `gorm:"default:0" json:"used_quantity"`
	UserLimit       int        `gorm:"default:1;comment:每人限领数量" json:"user_limit"`
	StartTime       time.Time  `gorm:"not null" json:"start_time"`
	EndTime         time.Time  `gorm:"not null" json:"end_time"`
	Status          int8       `gorm:"default:1;comment:0-禁用,1-正常" json:"status"`
	Description     string     `gorm:"type:text" json:"description"`
	
	// 关联关系
	Users           []User     `gorm:"many2many:user_coupons" json:"users,omitempty"`
	Orders          []Order    `gorm:"foreignKey:CouponID" json:"orders,omitempty"`
}

// TableName 指定表名
func (Coupon) TableName() string {
	return "coupons"
}

// UserCoupon 用户优惠券中间表
type UserCoupon struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	CouponID    uint       `gorm:"index;not null" json:"coupon_id"`
	Status      int8       `gorm:"default:1;comment:1-未使用,2-已使用,3-已过期" json:"status"`
	UsedAt      *time.Time `json:"used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// 关联关系
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coupon      Coupon     `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

// TableName 指定表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}
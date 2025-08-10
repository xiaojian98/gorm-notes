package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID      uint                    `json:"user_id" binding:"required"`
	AddressID   uint                    `json:"address_id" binding:"required"`
	Items       []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	CouponID    *uint                   `json:"coupon_id"`
	Remark      string                  `json:"remark"`
}

// CreateOrderItemRequest 创建订单项请求
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	SKUID     *uint `json:"sku_id"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// OrderService 订单服务
type OrderService struct {
	db *gorm.DB
	mu sync.RWMutex // 用于并发控制
}

// NewOrderService 创建订单服务实例
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*Order, error) {
	// 参数验证
	if err := s.validateCreateOrderRequest(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}

	// 开始事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("开始事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 验证收货地址
	address, err := s.validateAddress(tx, req.UserID, req.AddressID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("验证收货地址失败: %w", err)
	}

	// 验证商品和库存，计算总金额
	validatedItems, totalAmount, err := s.validateAndCalculateItems(tx, req.Items)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("验证商品失败: %w", err)
	}

	// 验证和使用优惠券
	discountAmount := int64(0)
	if req.CouponID != nil {
		discount, err := s.validateAndUseCoupon(tx, req.UserID, *req.CouponID, totalAmount)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("验证优惠券失败: %w", err)
		}
		discountAmount = discount
	}

	// 计算运费
	freightAmount := s.calculateFreight(address, validatedItems)

	// 计算最终金额
	finalAmount := totalAmount + freightAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	// 创建订单
	order := &Order{
		OrderNo:         s.generateOrderNo(),
		UserID:          req.UserID,
		Status:          1, // 待付款
		TotalAmount:     totalAmount,
		PayAmount:       finalAmount,
		FreightAmount:   freightAmount,
		DiscountAmount:  discountAmount,
		CouponID:        req.CouponID,
		ReceiverName:    address.Name,
		ReceiverPhone:   address.Phone,
		ReceiverAddress: fmt.Sprintf("%s%s%s%s", address.Province, address.City, address.District, address.Detail),
		Remark:          req.Remark,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 创建订单项
	for _, item := range validatedItems {
		orderItem := &OrderItem{
			OrderID:      order.ID,
			ProductID:    item.ProductID,
			SKUID:        item.SKUID,
			Quantity:     item.Quantity,
			Price:        item.Price,
			TotalPrice:   item.Price * int64(item.Quantity),
			ProductName:  item.ProductName,
			ProductSKU:   item.ProductSKU,
			ProductImage: item.ProductImage,
			ProductSpecs: item.ProductSpecs,
		}

		if err := tx.Create(orderItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建订单项失败: %w", err)
		}
	}

	// 扣减库存
	for _, item := range validatedItems {
		if err := s.deductStock(tx, item.ProductID, item.SKUID, item.Quantity); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("扣减库存失败: %w", err)
		}
	}

	// 清空购物车
	if err := s.clearCart(tx, req.UserID, req.Items); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("清空购物车失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return order, nil
}

// validateCreateOrderRequest 验证创建订单请求
func (s *OrderService) validateCreateOrderRequest(req *CreateOrderRequest) error {
	if req == nil {
		return errors.New("请求不能为空")
	}

	if req.UserID == 0 {
		return errors.New("用户ID不能为空")
	}

	if req.AddressID == 0 {
		return errors.New("收货地址ID不能为空")
	}

	if len(req.Items) == 0 {
		return errors.New("订单项不能为空")
	}

	for i, item := range req.Items {
		if item.ProductID == 0 {
			return fmt.Errorf("第%d个商品ID不能为空", i+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("第%d个商品数量必须大于0", i+1)
		}
	}

	return nil
}

// validateAddress 验证收货地址
func (s *OrderService) validateAddress(tx *gorm.DB, userID, addressID uint) (*Address, error) {
	var address Address
	err := tx.Where("id = ? AND user_id = ?", addressID, userID).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收货地址不存在")
		}
		return nil, err
	}
	return &address, nil
}

// ValidatedOrderItem 验证后的订单项
type ValidatedOrderItem struct {
	ProductID    uint
	SKUID        *uint
	Quantity     int
	Price        int64
	ProductName  string
	ProductSKU   string
	ProductImage string
	ProductSpecs json.RawMessage
}

// validateAndCalculateItems 验证商品信息、检查库存并计算订单总金额
func (s *OrderService) validateAndCalculateItems(tx *gorm.DB, items []CreateOrderItemRequest) ([]ValidatedOrderItem, int64, error) {
	validatedItems := make([]ValidatedOrderItem, 0, len(items))
	totalAmount := int64(0)

	for _, item := range items {
		if item.SKUID != nil {
			// 验证SKU
			var sku ProductSKU
			err := tx.Preload("Product").Where("id = ? AND product_id = ? AND status = 1", *item.SKUID, item.ProductID).First(&sku).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, 0, fmt.Errorf("商品SKU不存在或已下架")
				}
				return nil, 0, err
			}

			// 检查库存
			if sku.Stock < item.Quantity {
				return nil, 0, fmt.Errorf("商品 %s 库存不足，当前库存：%d", sku.Product.Name, sku.Stock)
			}

			validatedItem := ValidatedOrderItem{
				ProductID:    item.ProductID,
				SKUID:        item.SKUID,
				Quantity:     item.Quantity,
				Price:        sku.Price,
				ProductName:  sku.Product.Name,
				ProductSKU:   sku.SKU,
				ProductImage: sku.Image,
				ProductSpecs: sku.Specs,
			}
			validatedItems = append(validatedItems, validatedItem)
			totalAmount += sku.Price * int64(item.Quantity)
		} else {
			// 验证商品
			var product Product
			err := tx.Where("id = ? AND status = 1", item.ProductID).First(&product).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, 0, fmt.Errorf("商品不存在或已下架")
				}
				return nil, 0, err
			}

			// 检查库存
			if product.Stock < item.Quantity {
				return nil, 0, fmt.Errorf("商品 %s 库存不足，当前库存：%d", product.Name, product.Stock)
			}

			validatedItem := ValidatedOrderItem{
				ProductID:    item.ProductID,
				SKUID:        nil,
				Quantity:     item.Quantity,
				Price:        product.Price,
				ProductName:  product.Name,
				ProductSKU:   product.SKU,
				ProductImage: "", // 可以从商品图片中获取主图
				ProductSpecs: nil,
			}
			validatedItems = append(validatedItems, validatedItem)
			totalAmount += product.Price * int64(item.Quantity)
		}
	}

	return validatedItems, totalAmount, nil
}

// validateAndUseCoupon 验证优惠券的有效性、计算折扣并更新优惠券使用状态
func (s *OrderService) validateAndUseCoupon(tx *gorm.DB, userID, couponID uint, orderAmount int64) (int64, error) {
	// 检查用户是否拥有该优惠券
	var userCoupon UserCoupon
	err := tx.Preload("Coupon").Where("user_id = ? AND coupon_id = ? AND status = 1", userID, couponID).First(&userCoupon).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("优惠券不存在或已使用")
		}
		return 0, err
	}

	coupon := userCoupon.Coupon

	// 检查优惠券是否在有效期内
	now := time.Now()
	if now.Before(coupon.StartTime) || now.After(coupon.EndTime) {
		return 0, errors.New("优惠券不在有效期内")
	}

	// 检查最低消费金额
	if orderAmount < coupon.MinAmount {
		return 0, fmt.Errorf("订单金额不满足优惠券使用条件，最低消费：%.2f元", float64(coupon.MinAmount)/100)
	}

	// 计算折扣金额
	var discountAmount int64
	switch coupon.Type {
	case 1: // 满减
		discountAmount = coupon.Value
	case 2: // 折扣
		discountAmount = orderAmount * (100 - coupon.Value) / 100
	case 3: // 固定金额
		discountAmount = coupon.Value
	default:
		return 0, errors.New("不支持的优惠券类型")
	}

	// 检查最大优惠金额限制
	if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
		discountAmount = coupon.MaxDiscount
	}

	// 更新用户优惠券状态为已使用
	now = time.Now()
	err = tx.Model(&userCoupon).Updates(map[string]interface{}{
		"status":  2, // 已使用
		"used_at": &now,
	}).Error
	if err != nil {
		return 0, fmt.Errorf("更新优惠券状态失败: %w", err)
	}

	// 更新优惠券使用数量
	err = tx.Model(&coupon).UpdateColumn("used_quantity", gorm.Expr("used_quantity + ?", 1)).Error
	if err != nil {
		return 0, fmt.Errorf("更新优惠券使用数量失败: %w", err)
	}

	return discountAmount, nil
}

// calculateFreight 计算运费
func (s *OrderService) calculateFreight(address *Address, items []ValidatedOrderItem) int64 {
	// 简单的运费计算逻辑，实际项目中可能需要更复杂的计算
	// 这里假设：
	// 1. 订单金额超过100元免运费
	// 2. 否则根据地区收取不同运费

	totalAmount := int64(0)
	for _, item := range items {
		totalAmount += item.Price * int64(item.Quantity)
	}

	// 满100元免运费
	if totalAmount >= 10000 { // 100元
		return 0
	}

	// 根据省份计算运费
	switch address.Province {
	case "北京市", "上海市", "天津市", "重庆市":
		return 800 // 8元
	case "广东省", "江苏省", "浙江省":
		return 1000 // 10元
	default:
		return 1500 // 15元
	}
}

// deductStock 并发安全地扣减商品或SKU库存
func (s *OrderService) deductStock(tx *gorm.DB, productID uint, skuID *uint, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if skuID != nil {
		// 扣减SKU库存
		result := tx.Model(&ProductSKU{}).Where("id = ? AND stock >= ?", *skuID, quantity).
			UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("SKU库存不足")
		}
	} else {
		// 扣减商品库存
		result := tx.Model(&Product{}).Where("id = ? AND stock >= ?", productID, quantity).
			UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("商品库存不足")
		}
	}

	return nil
}

// clearCart 清空购物车中对应的商品
func (s *OrderService) clearCart(tx *gorm.DB, userID uint, items []CreateOrderItemRequest) error {
	for _, item := range items {
		query := tx.Where("user_id = ? AND product_id = ?", userID, item.ProductID)
		if item.SKUID != nil {
			query = query.Where("sku_id = ?", *item.SKUID)
		} else {
			query = query.Where("sku_id IS NULL")
		}

		if err := query.Delete(&Cart{}).Error; err != nil {
			return err
		}
	}
	return nil
}

// generateOrderNo 生成订单号
func (s *OrderService) generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderID uint, userID uint, reason string) error {
	// 开始事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 查询订单
	var order Order
	err := tx.Preload("Items").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.New("订单不存在")
		}
		tx.Rollback()
		return err
	}

	// 检查订单状态
	if order.Status != 1 { // 只有待付款状态的订单可以取消
		tx.Rollback()
		return errors.New("订单状态不允许取消")
	}

	// 更新订单状态
	now := time.Now()
	err = tx.Model(&order).Updates(map[string]interface{}{
		"status":        5, // 已取消
		"cancel_time":   &now,
		"cancel_reason": reason,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 回滚库存
	for _, item := range order.Items {
		if err := s.rollbackStock(tx, item.ProductID, item.SKUID, item.Quantity); err != nil {
			tx.Rollback()
			return fmt.Errorf("回滚库存失败: %w", err)
		}
	}

	// 回滚优惠券
	if order.CouponID != nil {
		if err := s.rollbackCoupon(tx, userID, *order.CouponID); err != nil {
			tx.Rollback()
			return fmt.Errorf("回滚优惠券失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// rollbackStock 回滚库存
func (s *OrderService) rollbackStock(tx *gorm.DB, productID uint, skuID *uint, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if skuID != nil {
		// 回滚SKU库存
		err := tx.Model(&ProductSKU{}).Where("id = ?", *skuID).
			UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
		if err != nil {
			return err
		}
	} else {
		// 回滚商品库存
		err := tx.Model(&Product{}).Where("id = ?", productID).
			UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// rollbackCoupon 回滚优惠券
func (s *OrderService) rollbackCoupon(tx *gorm.DB, userID, couponID uint) error {
	// 查找用户优惠券记录
	var userCoupon UserCoupon
	err := tx.Where("user_id = ? AND coupon_id = ? AND status = 2", userID, couponID).First(&userCoupon).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 如果找不到记录，说明优惠券没有被使用，不需要回滚
		}
		return err
	}

	// 恢复用户优惠券状态
	err = tx.Model(&userCoupon).Updates(map[string]interface{}{
		"status":  1, // 未使用
		"used_at": nil,
	}).Error
	if err != nil {
		return err
	}

	// 减少优惠券使用数量
	err = tx.Model(&Coupon{}).Where("id = ?", couponID).
		UpdateColumn("used_quantity", gorm.Expr("used_quantity - ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}
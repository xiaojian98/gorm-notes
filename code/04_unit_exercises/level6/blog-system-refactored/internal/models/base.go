package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型结构体
// 包含所有数据表的公共字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`                    // 主键ID
	CreatedAt time.Time      `json:"created_at"`                              // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                              // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`       // 软删除时间
}

// TableName 接口定义
// 用于自定义表名的模型需要实现此接口
type TableNamer interface {
	TableName() string
}

// SoftDeletable 软删除接口
// 支持软删除的模型需要实现此接口
type SoftDeletable interface {
	IsDeleted() bool
	Delete() error
	Restore() error
}

// Timestampable 时间戳接口
// 需要自动管理时间戳的模型可以实现此接口
type Timestampable interface {
	UpdateTimestamp()
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// IsDeleted 检查记录是否被软删除
// 返回: bool - 是否被删除
func (base *BaseModel) IsDeleted() bool {
	return base.DeletedAt.Valid
}

// UpdateTimestamp 更新时间戳
// 手动更新UpdatedAt字段为当前时间
func (base *BaseModel) UpdateTimestamp() {
	base.UpdatedAt = time.Now()
}

// GetCreatedAt 获取创建时间
// 返回: time.Time - 创建时间
func (base *BaseModel) GetCreatedAt() time.Time {
	return base.CreatedAt
}

// GetUpdatedAt 获取更新时间
// 返回: time.Time - 更新时间
func (base *BaseModel) GetUpdatedAt() time.Time {
	return base.UpdatedAt
}

// BeforeCreate GORM钩子函数 - 创建前
// 在创建记录前自动设置时间戳
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	base.CreatedAt = now
	base.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM钩子函数 - 更新前
// 在更新记录前自动设置更新时间戳
func (base *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	base.UpdatedAt = time.Now()
	return nil
}

// ModelStatus 模型状态枚举
type ModelStatus int

const (
	StatusInactive ModelStatus = iota // 0 - 未激活
	StatusActive                      // 1 - 激活
	StatusSuspended                   // 2 - 暂停
	StatusDeleted                     // 3 - 已删除
)

// String 返回状态的字符串表示
func (s ModelStatus) String() string {
	switch s {
	case StatusInactive:
		return "inactive"
	case StatusActive:
		return "active"
	case StatusSuspended:
		return "suspended"
	case StatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
// 返回: bool - 状态是否有效
func (s ModelStatus) IsValid() bool {
	return s >= StatusInactive && s <= StatusDeleted
}

// Priority 优先级枚举
type Priority int

const (
	PriorityLow    Priority = iota // 0 - 低优先级
	PriorityNormal                 // 1 - 普通优先级
	PriorityHigh                   // 2 - 高优先级
	PriorityUrgent                 // 3 - 紧急优先级
)

// String 返回优先级的字符串表示
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}

// IsValid 检查优先级是否有效
// 返回: bool - 优先级是否有效
func (p Priority) IsValid() bool {
	return p >= PriorityLow && p <= PriorityUrgent
}
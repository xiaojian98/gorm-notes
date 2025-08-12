// 03_blog_system/models/models.go - 数据模型定义
// 对应文档：02_GORM背景示例_博客系统实战.md

package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含公共字段
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string     `json:"username" gorm:"size:50;uniqueIndex:idx_user_username;not null" validate:"required,min=3,max=50"`
	Email       string     `json:"email" gorm:"size:100;uniqueIndex:idx_user_email;not null" validate:"required,email"`
	Password    string     `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
	Nickname    string     `json:"nickname" gorm:"size:50"`
	Avatar      string     `json:"avatar" gorm:"size:255"`
	Status      string     `json:"status" gorm:"size:20;default:active;index" validate:"oneof=active inactive banned"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LoginCount  int        `json:"login_count" gorm:"default:0"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
	// 一个用户只能有一个个人资料
	Profile *Profile `json:"profile,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_profiles_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一个用户可以发布多篇文章
	Posts []*Post `json:"posts,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_posts_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一个用户可以评论多篇文章
	Comments []*Comment `json:"comments,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_comments_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一个用户可以点赞多篇文章
	Likes []*Like `json:"likes,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_likes_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Profile 用户资料模型
type Profile struct {
	BaseModel
	UserID   uint       `json:"user_id" gorm:"uniqueIndex;not null"`
	Bio      string     `json:"bio" gorm:"type:text"`
	Website  string     `json:"website" gorm:"size:255"`
	Location string     `json:"location" gorm:"size:100"`
	Birthday *time.Time `json:"birthday"`
	Gender   string     `json:"gender" gorm:"size:10;default:unknown" validate:"oneof=male female unknown"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
	// 一个用户只能有一个个人资料
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_profiles_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Category 分类模型
type Category struct {
	BaseModel
	Name string `json:"name" gorm:"size:100;not null;index" validate:"required,max=100"`
	// 分类的URL slug，用于生成文章URL
	Slug        string `json:"slug" gorm:"size:100;uniqueIndex;not null" validate:"required,max=100"`
	Description string `json:"description" gorm:"type:text"`
	Color       string `json:"color" gorm:"size:7;default:#007bff"`
	Icon        string `json:"icon" gorm:"size:50"`
	SortOrder   int    `json:"sort_order" gorm:"default:0;index"`
	PostCount   int    `json:"post_count" gorm:"default:0"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
	// 一个分类可以有多篇文章
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:CategoryID;references:ID;constraint:fk_posts_category_id,OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Tag 标签模型
type Tag struct {
	BaseModel
	Name      string `json:"name" gorm:"size:50;not null;index" validate:"required,max=50"`
	Slug      string `json:"slug" gorm:"size:50;uniqueIndex;not null" validate:"required,max=50"`
	Color     string `json:"color" gorm:"size:7;default:#6c757d"`
	PostCount int    `json:"post_count" gorm:"default:0"`

	// 多对多关联关系
	// 一个标签可以对应多篇文章
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Post 文章模型
type Post struct {
	BaseModel
	Title        string     `json:"title" gorm:"size:200;not null;index" validate:"required,max=200"`
	Slug         string     `json:"slug" gorm:"size:200;uniqueIndex;not null" validate:"required,max=200"`
	Content      string     `json:"content" gorm:"type:longtext;not null" validate:"required"`
	Excerpt      string     `json:"excerpt" gorm:"type:text"`
	FeaturedImg  string     `json:"featured_img" gorm:"size:255"`
	Status       string     `json:"status" gorm:"size:20;default:draft;index" validate:"oneof=draft published archived"`
	ViewCount    int        `json:"view_count" gorm:"default:0;index"`
	LikeCount    int        `json:"like_count" gorm:"default:0;index"`
	CommentCount int        `json:"comment_count" gorm:"default:0"`
	PublishedAt  *time.Time `json:"published_at" gorm:"index"`
	UserID       uint       `json:"user_id" gorm:"not null;index"`
	CategoryID   *uint      `json:"category_id" gorm:"index"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
	// 一篇文章只能有一个作者
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_posts_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一篇文章可以有多个标签
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID;constraint:fk_posts_category_id,OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// 一篇文章可以有多个标签
	Tags []Tag `json:"tags,omitempty" gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一篇文章可以有多个评论
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID;references:ID;constraint:fk_comments_post_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一篇文章可以有多个点赞
	Likes []Like `json:"likes,omitempty" gorm:"foreignKey:TargetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	Content   string `json:"content" gorm:"type:text;not null" validate:"required"`
	Status    string `json:"status" gorm:"size:20;default:pending;index" validate:"oneof=pending approved rejected"`
	IP        string `json:"ip" gorm:"size:45"`
	UserAgent string `json:"user_agent" gorm:"size:255"`
	PostID    uint   `json:"post_id" gorm:"not null;index"`
	UserID    uint   `json:"user_id" gorm:"not null;index"`
	ParentID  *uint  `json:"parent_id" gorm:"index"` // 支持回复评论

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
	// 一个评论只能属于一个文章
	Post Post `json:"post,omitempty" gorm:"foreignKey:PostID;references:ID;constraint:fk_comments_post_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一个评论只能属于一个用户
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_comments_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一个评论可以有多个回复
	Parent *Comment `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:fk_comments_parent_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// 一个评论可以有多个回复
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:fk_comments_parent_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Like 点赞模型
type Like struct {
	BaseModel
	UserID     uint   `json:"user_id" gorm:"not null;index"`
	TargetID   uint   `json:"target_id" gorm:"not null;index"`
	TargetType string `json:"target_type" gorm:"size:20;not null;index" validate:"oneof=post comment"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
	// 一个点赞只能属于一个用户
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_likes_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// 复合唯一索引，防止重复点赞
	// 在AutoMigrate中会自动创建
}

// PostTag 文章标签关联表（GORM会自动创建，这里定义是为了自定义字段）
type PostTag struct {
	PostID    uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

// TableName 自定义表名
func (PostTag) TableName() string {
	return "post_tags"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑，比如密码加密
	return nil
}

// AfterCreate 创建后钩子
func (u *User) AfterCreate(tx *gorm.DB) error {
	// 创建用户后自动创建用户资料
	profile := Profile{
		UserID: u.ID,
		Bio:    "这个人很懒，什么都没有留下...",
	}
	return tx.Create(&profile).Error
}

// BeforeCreate 文章创建前钩子
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 如果是发布状态且没有设置发布时间，则设置为当前时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}

// AfterCreate 文章创建后钩子
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新分类的文章数量
	if p.CategoryID != nil {
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))
	}
	return nil
}

// AfterDelete 文章删除后钩子
func (p *Post) AfterDelete(tx *gorm.DB) error {
	// 更新分类的文章数量
	if p.CategoryID != nil {
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("PostCount", gorm.Expr("PostCount - ?", 1))
	}
	return nil
}

// BeforeCreate 点赞创建前钩子
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	// 检查是否已经点赞过
	var count int64
	tx.Model(&Like{}).Where("user_id = ? AND target_id = ? AND target_type = ?", l.UserID, l.TargetID, l.TargetType).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

// AfterCreate 点赞创建后钩子
func (l *Like) AfterCreate(tx *gorm.DB) error {
	// 更新目标对象的点赞数
	if l.TargetType == "post" {
		tx.Model(&Post{}).Where("id = ?", l.TargetID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
	}
	return nil
}

// AfterDelete 点赞删除后钩子
func (l *Like) AfterDelete(tx *gorm.DB) error {
	// 更新目标对象的点赞数
	if l.TargetType == "post" {
		tx.Model(&Post{}).Where("id = ?", l.TargetID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	}
	return nil
}

// AutoMigrate 自动迁移数据库 - 已弃用，请使用迁移系统
// 保留此函数是为了向后兼容，但建议使用 migrations 包
func AutoMigrate(db *gorm.DB) error {
	// 这个函数现在只是一个简单的包装器
	// 实际的迁移逻辑已移至 migrations 包
	return db.AutoMigrate(
		&User{},
		&Profile{},
		&Category{},
		&Tag{},
		&Post{},
		&Comment{},
		&Like{},
	)
}

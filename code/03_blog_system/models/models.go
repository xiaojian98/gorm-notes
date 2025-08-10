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
	Profile  *Profile   `json:"profile,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_profiles_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Posts    []*Post    `json:"posts,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_posts_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments []*Comment `json:"comments,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_comments_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Likes    []*Like    `json:"likes,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_likes_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_profiles_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Category 分类模型
type Category struct {
	BaseModel
	Name        string `json:"name" gorm:"size:100;not null;index" validate:"required,max=100"`
	Slug        string `json:"slug" gorm:"size:100;uniqueIndex;not null" validate:"required,max=100"`
	Description string `json:"description" gorm:"type:text"`
	Color       string `json:"color" gorm:"size:7;default:#007bff"`
	Icon        string `json:"icon" gorm:"size:50"`
	SortOrder   int    `json:"sort_order" gorm:"default:0;index"`
	PostCount   int    `json:"post_count" gorm:"default:0"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
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
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_posts_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID;constraint:fk_posts_category_id,OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID;references:ID;constraint:fk_comments_post_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Likes    []Like    `json:"likes,omitempty" gorm:"foreignKey:TargetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	Post    Post      `json:"post,omitempty" gorm:"foreignKey:PostID;references:ID;constraint:fk_comments_post_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User    User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_comments_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:fk_comments_parent_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:fk_comments_parent_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Like 点赞模型
type Like struct {
	BaseModel
	UserID     uint   `json:"user_id" gorm:"not null;index"`
	TargetID   uint   `json:"target_id" gorm:"not null;index"`
	TargetType string `json:"target_type" gorm:"size:20;not null;index" validate:"oneof=post comment"`

	// 关联关系 - 修复外键约束名称重复问题，为每个外键指定唯一名称
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
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("post_count", gorm.Expr("post_count - ?", 1))
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

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	// 按照依赖关系顺序迁移
	err := db.AutoMigrate(
		&User{},
		&Profile{},
		&Category{},
		&Tag{},
		&Post{},
		&Comment{},
		&Like{},
	)

	if err != nil {
		return err
	}

	// 创建复合索引
	if err := createIndexes(db); err != nil {
		return err
	}

	return nil
}

// createIndexes 创建复合索引
func createIndexes(db *gorm.DB) error {
	// 为Like表创建复合唯一索引
	// 先检查索引是否存在，如果不存在则创建
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'Like' AND index_name = 'idx_likes_user_target'").Scan(&count)
	if count == 0 {
		if err := db.Exec("CREATE UNIQUE INDEX idx_likes_user_target ON `Like`(UserID, TargetID, TargetType)").Error; err != nil {
			return err
		}
	}

	// 为Post表创建复合索引
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'Post' AND index_name = 'idx_posts_status_published'").Scan(&count)
	if count == 0 {
		if err := db.Exec("CREATE INDEX idx_posts_status_published ON `Post`(Status, PublishedAt)").Error; err != nil {
			return err
		}
	}

	// 为Comment表创建复合索引
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'Comment' AND index_name = 'idx_comments_post_status'").Scan(&count)
	if count == 0 {
		if err := db.Exec("CREATE INDEX idx_comments_post_status ON `Comment`(PostID, Status)").Error; err != nil {
			return err
		}
	}

	return nil
}

package models

import (
	"strings"
	"time"
)

// Post 文章模型
// 存储博客文章的主要信息
type Post struct {
	BaseModel
	Title       string      `gorm:"size:200;not null" json:"title"`                    // 文章标题
	Slug        string      `gorm:"uniqueIndex;size:200;not null" json:"slug"`         // URL别名，唯一索引
	Content     string      `gorm:"type:longtext;not null" json:"content"`             // 文章内容
	Excerpt     string      `gorm:"type:text" json:"excerpt"`                         // 文章摘要
	Status      PostStatus  `gorm:"default:0" json:"status"`                          // 文章状态
	AuthorID    uint        `gorm:"not null;index" json:"author_id"`                  // 作者ID
	CategoryID  *uint       `gorm:"index" json:"category_id,omitempty"`               // 分类ID
	FeaturedImg string      `gorm:"size:255" json:"featured_img"`                     // 特色图片
	ViewCount   int         `gorm:"default:0" json:"view_count"`                      // 浏览次数
	LikeCount   int         `gorm:"default:0" json:"like_count"`                      // 点赞次数
	CommentCount int        `gorm:"default:0" json:"comment_count"`                   // 评论次数
	PublishedAt *time.Time  `json:"published_at,omitempty"`                          // 发布时间
	IsFeatured  bool        `gorm:"default:false" json:"is_featured"`                 // 是否为精选文章
	IsTop       bool        `gorm:"default:false" json:"is_top"`                      // 是否置顶
	
	// 关联关系
	Author   *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`   // 作者
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 分类
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`     // 标签（多对多）
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`   // 评论
	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes,omitempty"`      // 点赞
	Meta     []PostMeta `gorm:"foreignKey:PostID" json:"meta,omitempty"`      // 元数据
}

// TableName 自定义表名
func (Post) TableName() string {
	return "posts"
}

// PostStatus 文章状态枚举
type PostStatus int

const (
	PostStatusDraft     PostStatus = iota // 0 - 草稿
	PostStatusPublished                   // 1 - 已发布
	PostStatusPrivate                     // 2 - 私有
	PostStatusTrash                       // 3 - 回收站
)

// String 返回状态的字符串表示
func (s PostStatus) String() string {
	switch s {
	case PostStatusDraft:
		return "draft"
	case PostStatusPublished:
		return "published"
	case PostStatusPrivate:
		return "private"
	case PostStatusTrash:
		return "trash"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
func (s PostStatus) IsValid() bool {
	return s >= PostStatusDraft && s <= PostStatusTrash
}

// Category 分类模型
// 存储文章分类信息
type Category struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`        // 分类名称，唯一索引
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"`        // URL别名，唯一索引
	Description string `gorm:"type:text" json:"description"`                    // 分类描述
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"`                // 父分类ID
	SortOrder   int    `gorm:"default:0" json:"sort_order"`                     // 排序
	PostCount   int    `gorm:"default:0" json:"post_count"`                     // 文章数量
	IsActive    bool   `gorm:"default:true" json:"is_active"`                   // 是否激活
	
	// 关联关系
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`   // 父分类
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子分类
	Posts    []Post     `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`  // 文章
}

// TableName 自定义表名
func (Category) TableName() string {
	return "categories"
}

// Tag 标签模型
// 存储文章标签信息
type Tag struct {
	BaseModel
	Name      string `gorm:"uniqueIndex;size:50;not null" json:"name"`     // 标签名称，唯一索引
	Slug      string `gorm:"uniqueIndex;size:50;not null" json:"slug"`     // URL别名，唯一索引
	Color     string `gorm:"size:7;default:'#007bff'" json:"color"`        // 标签颜色
	PostCount int    `gorm:"default:0" json:"post_count"`                 // 文章数量
	IsActive  bool   `gorm:"default:true" json:"is_active"`               // 是否激活
	
	// 关联关系
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"` // 文章（多对多）
}

// TableName 自定义表名
func (Tag) TableName() string {
	return "tags"
}

// PostMeta 文章元数据模型
// 存储文章的额外元数据信息
type PostMeta struct {
	BaseModel
	PostID uint   `gorm:"not null;index" json:"post_id"` // 文章ID
	Key    string `gorm:"size:100;not null" json:"key"`   // 元数据键
	Value  string `gorm:"type:text" json:"value"`         // 元数据值
	
	// 关联关系
	Post *Post `gorm:"foreignKey:PostID" json:"post,omitempty"` // 文章
}

// TableName 自定义表名
func (PostMeta) TableName() string {
	return "post_meta"
}

// PostMethods 文章模型的方法

// IsPublished 检查文章是否已发布
// 返回: bool - 文章是否已发布
func (p *Post) IsPublished() bool {
	return p.Status == PostStatusPublished && p.PublishedAt != nil
}

// IsDraft 检查文章是否为草稿
// 返回: bool - 文章是否为草稿
func (p *Post) IsDraft() bool {
	return p.Status == PostStatusDraft
}

// IsPrivate 检查文章是否为私有
// 返回: bool - 文章是否为私有
func (p *Post) IsPrivate() bool {
	return p.Status == PostStatusPrivate
}

// Publish 发布文章
// 将文章状态设置为已发布，并设置发布时间
func (p *Post) Publish() {
	p.Status = PostStatusPublished
	now := time.Now()
	p.PublishedAt = &now
	p.UpdateTimestamp()
}

// UnPublish 取消发布文章
// 将文章状态设置为草稿
func (p *Post) UnPublish() {
	p.Status = PostStatusDraft
	p.PublishedAt = nil
	p.UpdateTimestamp()
}

// IncrementViewCount 增加浏览次数
// 将文章的浏览次数加1
func (p *Post) IncrementViewCount() {
	p.ViewCount++
	p.UpdateTimestamp()
}

// IncrementLikeCount 增加点赞次数
// 将文章的点赞次数加1
func (p *Post) IncrementLikeCount() {
	p.LikeCount++
	p.UpdateTimestamp()
}

// DecrementLikeCount 减少点赞次数
// 将文章的点赞次数减1
func (p *Post) DecrementLikeCount() {
	if p.LikeCount > 0 {
		p.LikeCount--
		p.UpdateTimestamp()
	}
}

// IncrementCommentCount 增加评论次数
// 将文章的评论次数加1
func (p *Post) IncrementCommentCount() {
	p.CommentCount++
	p.UpdateTimestamp()
}

// DecrementCommentCount 减少评论次数
// 将文章的评论次数减1
func (p *Post) DecrementCommentCount() {
	if p.CommentCount > 0 {
		p.CommentCount--
		p.UpdateTimestamp()
	}
}

// GetTagNames 获取标签名称列表
// 返回: []string - 标签名称数组
func (p *Post) GetTagNames() []string {
	var names []string
	for _, tag := range p.Tags {
		names = append(names, tag.Name)
	}
	return names
}

// GetTagsString 获取标签字符串
// 返回: string - 逗号分隔的标签名称
func (p *Post) GetTagsString() string {
	return strings.Join(p.GetTagNames(), ", ")
}

// CategoryMethods 分类模型的方法

// IsRootCategory 检查是否为根分类
// 返回: bool - 是否为根分类（没有父分类）
func (c *Category) IsRootCategory() bool {
	return c.ParentID == nil
}

// HasChildren 检查是否有子分类
// 返回: bool - 是否有子分类
func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}

// IncrementPostCount 增加文章数量
// 将分类的文章数量加1
func (c *Category) IncrementPostCount() {
	c.PostCount++
	c.UpdateTimestamp()
}

// DecrementPostCount 减少文章数量
// 将分类的文章数量减1
func (c *Category) DecrementPostCount() {
	if c.PostCount > 0 {
		c.PostCount--
		c.UpdateTimestamp()
	}
}

// TagMethods 标签模型的方法

// IncrementPostCount 增加文章数量
// 将标签的文章数量加1
func (t *Tag) IncrementPostCount() {
	t.PostCount++
	t.UpdateTimestamp()
}

// DecrementPostCount 减少文章数量
// 将标签的文章数量减1
func (t *Tag) DecrementPostCount() {
	if t.PostCount > 0 {
		t.PostCount--
		t.UpdateTimestamp()
	}
}
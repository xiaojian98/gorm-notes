package models

import (
	"time"
)

// Comment 评论模型
// 存储用户对文章的评论信息
type Comment struct {
	BaseModel
	PostID    uint          `gorm:"not null;index" json:"post_id"`               // 文章ID
	UserID    uint          `gorm:"not null;index" json:"user_id"`               // 用户ID
	ParentID  *uint         `gorm:"index" json:"parent_id,omitempty"`            // 父评论ID（用于回复）
	Content   string        `gorm:"type:text;not null" json:"content"`           // 评论内容
	Status    CommentStatus `gorm:"default:0" json:"status"`                    // 评论状态
	Level     int           `gorm:"default:1" json:"level"`                     // 评论层级
	LikeCount int           `gorm:"default:0" json:"like_count"`                // 点赞数
	IPAddress string        `gorm:"size:45" json:"ip_address,omitempty"`        // IP地址
	UserAgent string        `gorm:"size:255" json:"user_agent,omitempty"`       // 用户代理
	IsSpam    bool          `gorm:"default:false" json:"is_spam"`               // 是否为垃圾评论
	
	// 关联关系
	Post     *Post     `gorm:"foreignKey:PostID" json:"post,omitempty"`     // 文章
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`     // 用户
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"` // 父评论
	Children []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子评论
	Likes    []Like    `gorm:"foreignKey:CommentID" json:"likes,omitempty"` // 点赞
}

// TableName 自定义表名
func (Comment) TableName() string {
	return "comments"
}

// CommentStatus 评论状态枚举
type CommentStatus int

const (
	CommentStatusPending  CommentStatus = iota // 0 - 待审核
	CommentStatusApproved                      // 1 - 已通过
	CommentStatusRejected                      // 2 - 已拒绝
	CommentStatusSpam                          // 3 - 垃圾评论
	CommentStatusDeleted                       // 4 - 已删除
)

// String 返回状态的字符串表示
func (s CommentStatus) String() string {
	switch s {
	case CommentStatusPending:
		return "pending"
	case CommentStatusApproved:
		return "approved"
	case CommentStatusRejected:
		return "rejected"
	case CommentStatusSpam:
		return "spam"
	case CommentStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// IsValid 检查状态是否有效
func (s CommentStatus) IsValid() bool {
	return s >= CommentStatusPending && s <= CommentStatusDeleted
}

// Like 点赞模型
// 存储用户对文章或评论的点赞信息
type Like struct {
	BaseModel
	UserID     uint   `gorm:"not null;index" json:"user_id"`               // 用户ID
	TargetID   uint   `gorm:"not null;index" json:"target_id"`             // 目标ID（文章或评论ID）
	TargetType string `gorm:"size:20;not null;index" json:"target_type"`   // 目标类型（post或comment）
	PostID     *uint  `gorm:"index" json:"post_id,omitempty"`              // 文章ID（可选）
	CommentID  *uint  `gorm:"index" json:"comment_id,omitempty"`           // 评论ID（可选）
	Type       LikeType `gorm:"not null" json:"type"`                     // 点赞类型
	
	// 关联关系
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`    // 用户
	Post    *Post    `gorm:"foreignKey:PostID" json:"post,omitempty"`    // 文章
	Comment *Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"` // 评论
}

// TableName 自定义表名
func (Like) TableName() string {
	return "likes"
}

// LikeType 点赞类型枚举
type LikeType int

const (
	LikeTypePost    LikeType = iota // 0 - 文章点赞
	LikeTypeComment                 // 1 - 评论点赞
)

// String 返回类型的字符串表示
func (t LikeType) String() string {
	switch t {
	case LikeTypePost:
		return "post"
	case LikeTypeComment:
		return "comment"
	default:
		return "unknown"
	}
}

// IsValid 检查类型是否有效
func (t LikeType) IsValid() bool {
	return t >= LikeTypePost && t <= LikeTypeComment
}

// Notification 通知模型
// 存储系统通知信息
type Notification struct {
	BaseModel
	UserID   uint             `gorm:"not null;index" json:"user_id"`   // 接收通知的用户ID
	Type     NotificationType `gorm:"not null" json:"type"`            // 通知类型
	Title    string           `gorm:"size:200;not null" json:"title"`  // 通知标题
	Content  string           `gorm:"type:text" json:"content"`        // 通知内容
	Data     string           `gorm:"type:json" json:"data,omitempty"` // 额外数据（JSON格式）
	IsRead   bool             `gorm:"default:false" json:"is_read"`    // 是否已读
	ReadAt   *time.Time       `json:"read_at,omitempty"`              // 阅读时间
	Priority Priority         `gorm:"default:1" json:"priority"`       // 优先级
	
	// 关联关系
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 用户
}

// TableName 自定义表名
func (Notification) TableName() string {
	return "notifications"
}

// NotificationType 通知类型枚举
type NotificationType int

const (
	NotificationTypeSystem  NotificationType = iota // 0 - 系统通知
	NotificationTypeComment                         // 1 - 评论通知
	NotificationTypeLike                            // 2 - 点赞通知
	NotificationTypeFollow                          // 3 - 关注通知
	NotificationTypePost                            // 4 - 文章通知
	NotificationTypeReply                           // 5 - 回复通知
)

// String 返回类型的字符串表示
func (t NotificationType) String() string {
	switch t {
	case NotificationTypeSystem:
		return "system"
	case NotificationTypeComment:
		return "comment"
	case NotificationTypeLike:
		return "like"
	case NotificationTypeFollow:
		return "follow"
	case NotificationTypePost:
		return "post"
	case NotificationTypeReply:
		return "reply"
	default:
		return "unknown"
	}
}

// IsValid 检查类型是否有效
func (t NotificationType) IsValid() bool {
	return t >= NotificationTypeSystem && t <= NotificationTypeReply
}

// CommentMethods 评论模型的方法

// IsApproved 检查评论是否已通过审核
// 返回: bool - 评论是否已通过审核
func (c *Comment) IsApproved() bool {
	return c.Status == CommentStatusApproved
}

// IsPending 检查评论是否待审核
// 返回: bool - 评论是否待审核
func (c *Comment) IsPending() bool {
	return c.Status == CommentStatusPending
}

// IsRejected 检查评论是否被拒绝
// 返回: bool - 评论是否被拒绝
func (c *Comment) IsRejected() bool {
	return c.Status == CommentStatusRejected
}

// IsReply 检查是否为回复评论
// 返回: bool - 是否为回复评论
func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}

// IsRootComment 检查是否为根评论
// 返回: bool - 是否为根评论（非回复）
func (c *Comment) IsRootComment() bool {
	return c.ParentID == nil
}

// Approve 通过评论审核
// 将评论状态设置为已通过
func (c *Comment) Approve() {
	c.Status = CommentStatusApproved
	c.UpdateTimestamp()
}

// Reject 拒绝评论
// 将评论状态设置为已拒绝
func (c *Comment) Reject() {
	c.Status = CommentStatusRejected
	c.UpdateTimestamp()
}

// MarkAsSpam 标记为垃圾评论
// 将评论标记为垃圾评论
func (c *Comment) MarkAsSpam() {
	c.Status = CommentStatusSpam
	c.IsSpam = true
	c.UpdateTimestamp()
}

// IncrementLikeCount 增加点赞数
// 将评论的点赞数加1
func (c *Comment) IncrementLikeCount() {
	c.LikeCount++
	c.UpdateTimestamp()
}

// DecrementLikeCount 减少点赞数
// 将评论的点赞数减1
func (c *Comment) DecrementLikeCount() {
	if c.LikeCount > 0 {
		c.LikeCount--
		c.UpdateTimestamp()
	}
}

// LikeMethods 点赞模型的方法

// IsPostLike 检查是否为文章点赞
// 返回: bool - 是否为文章点赞
func (l *Like) IsPostLike() bool {
	return l.Type == LikeTypePost && l.PostID != nil
}

// IsCommentLike 检查是否为评论点赞
// 返回: bool - 是否为评论点赞
func (l *Like) IsCommentLike() bool {
	return l.Type == LikeTypeComment && l.CommentID != nil
}

// IsValid 检查点赞记录是否有效
// 返回: bool - 点赞记录是否有效
func (l *Like) IsValid() bool {
	if l.Type == LikeTypePost {
		return l.PostID != nil && l.CommentID == nil
	}
	if l.Type == LikeTypeComment {
		return l.CommentID != nil && l.PostID == nil
	}
	return false
}

// NotificationMethods 通知模型的方法

// MarkAsRead 标记通知为已读
// 将通知标记为已读并设置阅读时间
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
	n.UpdateTimestamp()
}

// MarkAsUnread 标记通知为未读
// 将通知标记为未读
func (n *Notification) MarkAsUnread() {
	n.IsRead = false
	n.ReadAt = nil
	n.UpdateTimestamp()
}

// IsHighPriority 检查是否为高优先级通知
// 返回: bool - 是否为高优先级通知
func (n *Notification) IsHighPriority() bool {
	return n.Priority == PriorityHigh || n.Priority == PriorityUrgent
}

// IsUrgent 检查是否为紧急通知
// 返回: bool - 是否为紧急通知
func (n *Notification) IsUrgent() bool {
	return n.Priority == PriorityUrgent
}
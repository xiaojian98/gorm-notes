package models

import (
	"time"
)

// User 用户模型
// 存储用户的基本信息和认证数据
type User struct {
	BaseModel
	Username    string      `gorm:"uniqueIndex;size:50;not null" json:"username"`           // 用户名，唯一索引
	Email       string      `gorm:"uniqueIndex;size:100;not null" json:"email"`             // 邮箱，唯一索引
	PasswordHash string     `gorm:"size:255;not null" json:"-"`                            // 密码哈希，不返回给前端
	Status      ModelStatus `gorm:"default:1" json:"status"`                               // 用户状态
	LastLoginAt *time.Time  `json:"last_login_at,omitempty"`                              // 最后登录时间
	
	// 关联关系
	Profile       *UserProfile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`       // 用户资料
	Posts         []Post         `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`       // 发布的文章
	Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`      // 发表的评论
	Likes         []Like         `gorm:"foreignKey:UserID" json:"likes,omitempty"`         // 点赞记录
	Following     []Follow       `gorm:"foreignKey:FollowerID" json:"following,omitempty"` // 关注的用户
	Followers     []Follow       `gorm:"foreignKey:FollowedID" json:"followers,omitempty"` // 粉丝
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"` // 通知
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

// UserProfile 用户资料模型
// 存储用户的详细个人信息
type UserProfile struct {
	BaseModel
	UserID      uint    `gorm:"uniqueIndex;not null" json:"user_id"`              // 用户ID，外键
	Nickname    string  `gorm:"size:100" json:"nickname"`                         // 昵称
	Avatar      string  `gorm:"size:255" json:"avatar"`                           // 头像URL
	Bio         string  `gorm:"type:text" json:"bio"`                             // 个人简介
	Location    string  `gorm:"size:100" json:"location"`                        // 所在地
	Website     string  `gorm:"size:255" json:"website"`                         // 个人网站
	Birthday    *time.Time `json:"birthday,omitempty"`                           // 生日
	Gender      string  `gorm:"size:10" json:"gender"`                           // 性别
	Phone       string  `gorm:"size:20" json:"phone"`                            // 电话
	IsPublic    bool    `gorm:"default:true" json:"is_public"`                   // 资料是否公开
	
	// 统计字段
	PostsCount     int `gorm:"default:0" json:"posts_count"`     // 文章数量
	FollowersCount int `gorm:"default:0" json:"followers_count"` // 粉丝数量
	FollowingCount int `gorm:"default:0" json:"following_count"` // 关注数量
	LikesCount     int `gorm:"default:0" json:"likes_count"`     // 获赞数量
	
	// 关联关系
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联的用户
}

// TableName 自定义表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// Follow 关注关系模型
// 存储用户之间的关注关系
type Follow struct {
	BaseModel
	FollowerID uint `gorm:"not null;index" json:"follower_id"` // 关注者ID
	FollowedID uint `gorm:"not null;index" json:"followed_id"` // 被关注者ID
	
	// 关联关系
	Follower *User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"` // 关注者
	Followed *User `gorm:"foreignKey:FollowedID" json:"followed,omitempty"` // 被关注者
}

// TableName 自定义表名
func (Follow) TableName() string {
	return "follows"
}

// UserMethods 用户模型的方法

// IsActive 检查用户是否激活
// 返回: bool - 用户是否处于激活状态
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsSuspended 检查用户是否被暂停
// 返回: bool - 用户是否被暂停
func (u *User) IsSuspended() bool {
	return u.Status == StatusSuspended
}

// Activate 激活用户
// 将用户状态设置为激活
func (u *User) Activate() {
	u.Status = StatusActive
	u.UpdateTimestamp()
}

// Suspend 暂停用户
// 将用户状态设置为暂停
func (u *User) Suspend() {
	u.Status = StatusSuspended
	u.UpdateTimestamp()
}

// UpdateLastLogin 更新最后登录时间
// 将最后登录时间设置为当前时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdateTimestamp()
}

// GetDisplayName 获取显示名称
// 优先返回昵称，如果没有昵称则返回用户名
// 返回: string - 显示名称
func (u *User) GetDisplayName() string {
	if u.Profile != nil && u.Profile.Nickname != "" {
		return u.Profile.Nickname
	}
	return u.Username
}

// UserProfileMethods 用户资料模型的方法

// GetAge 计算用户年龄
// 根据生日计算当前年龄
// 返回: int - 年龄，如果没有设置生日则返回0
func (up *UserProfile) GetAge() int {
	if up.Birthday == nil {
		return 0
	}
	return int(time.Since(*up.Birthday).Hours() / 24 / 365)
}

// IsProfileComplete 检查资料是否完整
// 检查必要的资料字段是否已填写
// 返回: bool - 资料是否完整
func (up *UserProfile) IsProfileComplete() bool {
	return up.Nickname != "" && up.Bio != "" && up.Avatar != ""
}

// UpdateStats 更新统计数据
// 参数: postsCount, followersCount, followingCount, likesCount - 各项统计数据
func (up *UserProfile) UpdateStats(postsCount, followersCount, followingCount, likesCount int) {
	up.PostsCount = postsCount
	up.FollowersCount = followersCount
	up.FollowingCount = followingCount
	up.LikesCount = likesCount
	up.UpdateTimestamp()
}

// FollowMethods 关注关系模型的方法

// IsValid 检查关注关系是否有效
// 检查关注者和被关注者是否不同
// 返回: bool - 关注关系是否有效
func (f *Follow) IsValid() bool {
	return f.FollowerID != f.FollowedID && f.FollowerID > 0 && f.FollowedID > 0
}
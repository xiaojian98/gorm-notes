package services

import (
	"errors"
	"fmt"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// UserService 用户服务接口
// 定义用户相关的业务操作
type UserService interface {
	// 用户基本操作
	CreateUser(user *models.User) error                    // 创建用户
	GetUserByID(id uint) (*models.User, error)             // 根据ID获取用户
	GetUserByUsername(username string) (*models.User, error) // 根据用户名获取用户
	GetUserByEmail(email string) (*models.User, error)     // 根据邮箱获取用户
	UpdateUser(user *models.User) error                    // 更新用户信息
	DeleteUser(id uint) error                              // 删除用户
	ListUsers(offset, limit int) ([]models.User, int64, error) // 分页获取用户列表
	
	// 用户资料操作
	CreateUserProfile(profile *models.UserProfile) error   // 创建用户资料
	GetUserProfile(userID uint) (*models.UserProfile, error) // 获取用户资料
	UpdateUserProfile(profile *models.UserProfile) error   // 更新用户资料
	
	// 用户关注操作
	FollowUser(followerID, followingID uint) error         // 关注用户
	UnfollowUser(followerID, followingID uint) error       // 取消关注
	IsFollowing(followerID, followingID uint) (bool, error) // 检查是否关注
	GetFollowers(userID uint, offset, limit int) ([]models.User, int64, error) // 获取粉丝列表
	GetFollowing(userID uint, offset, limit int) ([]models.User, int64, error) // 获取关注列表
	
	// 用户状态操作
	ActivateUser(id uint) error                            // 激活用户
	DeactivateUser(id uint) error                          // 停用用户
	BanUser(id uint, reason string) error                 // 封禁用户
	UnbanUser(id uint) error                               // 解封用户
	
	// 用户统计
	GetUserStats(userID uint) (*UserStats, error)         // 获取用户统计信息
	GetActiveUsers(limit int) ([]models.User, error)      // 获取活跃用户
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
// 参数: db - 数据库连接
// 返回: UserService - 用户服务接口实例
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

// UserStats 用户统计信息结构体
type UserStats struct {
	PostsCount     int64 `json:"posts_count"`     // 文章数量
	CommentsCount  int64 `json:"comments_count"`  // 评论数量
	LikesReceived  int64 `json:"likes_received"`  // 获得点赞数
	FollowersCount int64 `json:"followers_count"` // 粉丝数量
	FollowingCount int64 `json:"following_count"` // 关注数量
	TotalViews     int64 `json:"total_views"`     // 总浏览量
	JoinedDays     int   `json:"joined_days"`     // 加入天数
}

// 用户基本操作实现

// CreateUser 创建用户
// 参数: user - 用户模型
// 返回: error - 错误信息
func (s *userService) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("用户信息不能为空")
	}
	
	// 验证必填字段
	if user.Username == "" {
		return errors.New("用户名不能为空")
	}
	if user.Email == "" {
		return errors.New("邮箱不能为空")
	}
	
	// 检查用户名是否已存在
	existingUser := &models.User{}
	if err := s.db.Where("username = ?", user.Username).First(existingUser).Error; err == nil {
		return errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", user.Email).First(existingUser).Error; err == nil {
		return errors.New("邮箱已存在")
	}
	
	// 设置默认值
	if user.Status == models.ModelStatus(0) {
		user.Status = models.StatusActive
	}
	
	return s.db.Create(user).Error
}

// GetUserByID 根据ID获取用户
// 参数: id - 用户ID
// 返回: *models.User - 用户模型, error - 错误信息
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	user := &models.User{}
	err := s.db.Preload("Profile").First(user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserByUsername 根据用户名获取用户
// 参数: username - 用户名
// 返回: *models.User - 用户模型, error - 错误信息
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	
	user := &models.User{}
	err := s.db.Preload("Profile").Where("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserByEmail 根据邮箱获取用户
// 参数: email - 邮箱地址
// 返回: *models.User - 用户模型, error - 错误信息
func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	
	user := &models.User{}
	err := s.db.Preload("Profile").Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// UpdateUser 更新用户信息
// 参数: user - 用户模型
// 返回: error - 错误信息
func (s *userService) UpdateUser(user *models.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("用户信息不完整")
	}
	
	// 检查用户是否存在
	existingUser := &models.User{}
	if err := s.db.First(existingUser, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	
	// 如果更新用户名，检查是否重复
	if user.Username != "" && user.Username != existingUser.Username {
		tempUser := &models.User{}
		if err := s.db.Where("username = ? AND id != ?", user.Username, user.ID).First(tempUser).Error; err == nil {
			return errors.New("用户名已存在")
		}
	}
	
	// 如果更新邮箱，检查是否重复
	if user.Email != "" && user.Email != existingUser.Email {
		tempUser := &models.User{}
		if err := s.db.Where("email = ? AND id != ?", user.Email, user.ID).First(tempUser).Error; err == nil {
			return errors.New("邮箱已存在")
		}
	}
	
	return s.db.Save(user).Error
}

// DeleteUser 删除用户（软删除）
// 参数: id - 用户ID
// 返回: error - 错误信息
func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	// 检查用户是否存在
	user := &models.User{}
	if err := s.db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	
	return s.db.Delete(user).Error
}

// ListUsers 分页获取用户列表
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.User - 用户列表, int64 - 总数量, error - 错误信息
func (s *userService) ListUsers(offset, limit int) ([]models.User, int64, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	var total int64
	
	// 获取总数
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取用户列表
	err := s.db.Preload("Profile").Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

// 用户资料操作实现

// CreateUserProfile 创建用户资料
// 参数: profile - 用户资料模型
// 返回: error - 错误信息
func (s *userService) CreateUserProfile(profile *models.UserProfile) error {
	if profile == nil || profile.UserID == 0 {
		return errors.New("用户资料信息不完整")
	}
	
	// 检查用户是否存在
	user := &models.User{}
	if err := s.db.First(user, profile.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	
	// 检查是否已有资料
	existingProfile := &models.UserProfile{}
	if err := s.db.Where("user_id = ?", profile.UserID).First(existingProfile).Error; err == nil {
		return errors.New("用户资料已存在")
	}
	
	return s.db.Create(profile).Error
}

// GetUserProfile 获取用户资料
// 参数: userID - 用户ID
// 返回: *models.UserProfile - 用户资料模型, error - 错误信息
func (s *userService) GetUserProfile(userID uint) (*models.UserProfile, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	profile := &models.UserProfile{}
	err := s.db.Where("user_id = ?", userID).First(profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户资料不存在")
		}
		return nil, err
	}
	
	return profile, nil
}

// UpdateUserProfile 更新用户资料
// 参数: profile - 用户资料模型
// 返回: error - 错误信息
func (s *userService) UpdateUserProfile(profile *models.UserProfile) error {
	if profile == nil || profile.ID == 0 {
		return errors.New("用户资料信息不完整")
	}
	
	// 检查资料是否存在
	existingProfile := &models.UserProfile{}
	if err := s.db.First(existingProfile, profile.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户资料不存在")
		}
		return err
	}
	
	return s.db.Save(profile).Error
}

// 用户关注操作实现

// FollowUser 关注用户
// 参数: followerID - 关注者ID, followingID - 被关注者ID
// 返回: error - 错误信息
func (s *userService) FollowUser(followerID, followingID uint) error {
	if followerID == 0 || followingID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	if followerID == followingID {
		return errors.New("不能关注自己")
	}
	
	// 检查用户是否存在
	var count int64
	s.db.Model(&models.User{}).Where("id IN (?, ?)", followerID, followingID).Count(&count)
	if count != 2 {
		return errors.New("用户不存在")
	}
	
	// 检查是否已关注
	existingFollow := &models.Follow{}
	if err := s.db.Where("follower_id = ? AND followed_id = ?", followerID, followingID).First(existingFollow).Error; err == nil {
		return errors.New("已经关注该用户")
	}

	follow := &models.Follow{
		FollowerID: followerID,
		FollowedID: followingID,
	}
	
	return s.db.Create(follow).Error
}

// UnfollowUser 取消关注
// 参数: followerID - 关注者ID, followingID - 被关注者ID
// 返回: error - 错误信息
func (s *userService) UnfollowUser(followerID, followingID uint) error {
	if followerID == 0 || followingID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	follow := &models.Follow{}
	err := s.db.Where("follower_id = ? AND followed_id = ?", followerID, followingID).First(follow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未关注该用户")
		}
		return err
	}
	
	return s.db.Delete(follow).Error
}

// IsFollowing 检查是否关注
// 参数: followerID - 关注者ID, followingID - 被关注者ID
// 返回: bool - 是否关注, error - 错误信息
func (s *userService) IsFollowing(followerID, followingID uint) (bool, error) {
	if followerID == 0 || followingID == 0 {
		return false, errors.New("用户ID不能为空")
	}
	
	var count int64
	err := s.db.Model(&models.Follow{}).Where("follower_id = ? AND followed_id = ?", followerID, followingID).Count(&count).Error
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// GetFollowers 获取粉丝列表
// 参数: userID - 用户ID, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 粉丝列表, int64 - 总数量, error - 错误信息
func (s *userService) GetFollowers(userID uint, offset, limit int) ([]models.User, int64, error) {
	if userID == 0 {
		return nil, 0, errors.New("用户ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	var total int64
	
	// 获取总数
	if err := s.db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取粉丝列表
	err := s.db.Table("users").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("follows.created_at DESC").
		Find(&users).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

// GetFollowing 获取关注列表
// 参数: userID - 用户ID, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 关注列表, int64 - 总数量, error - 错误信息
func (s *userService) GetFollowing(userID uint, offset, limit int) ([]models.User, int64, error) {
	if userID == 0 {
		return nil, 0, errors.New("用户ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	var total int64
	
	// 获取总数
	if err := s.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取关注列表
	err := s.db.Table("users").
		Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("follows.created_at DESC").
		Find(&users).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

// 用户状态操作实现

// ActivateUser 激活用户
// 参数: id - 用户ID
// 返回: error - 错误信息
func (s *userService) ActivateUser(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("status", "active").Error
}

// DeactivateUser 停用用户
// 参数: id - 用户ID
// 返回: error - 错误信息
func (s *userService) DeactivateUser(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("status", "inactive").Error
}

// BanUser 封禁用户
// 参数: id - 用户ID, reason - 封禁原因
// 返回: error - 错误信息
func (s *userService) BanUser(id uint, reason string) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	updates := map[string]interface{}{
		"status": "banned",
	}
	if reason != "" {
		updates["ban_reason"] = reason
	}
	
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// UnbanUser 解封用户
// 参数: id - 用户ID
// 返回: error - 错误信息
func (s *userService) UnbanUser(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	updates := map[string]interface{}{
		"status":     "active",
		"ban_reason": nil,
	}
	
	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// 用户统计实现

// GetUserStats 获取用户统计信息
// 参数: userID - 用户ID
// 返回: *UserStats - 用户统计信息, error - 错误信息
func (s *userService) GetUserStats(userID uint) (*UserStats, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	// 检查用户是否存在
	user := &models.User{}
	if err := s.db.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	stats := &UserStats{}
	
	// 获取文章数量
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&stats.PostsCount)
	
	// 获取评论数量
	s.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&stats.CommentsCount)
	
	// 获取获得的点赞数（文章点赞 + 评论点赞）
	var postLikes, commentLikes int64
	s.db.Table("likes").Joins("JOIN posts ON likes.target_id = posts.id").
		Where("likes.target_type = 'post' AND posts.author_id = ?", userID).Count(&postLikes)
	s.db.Table("likes").Joins("JOIN comments ON likes.target_id = comments.id").
		Where("likes.target_type = 'comment' AND comments.user_id = ?", userID).Count(&commentLikes)
	stats.LikesReceived = postLikes + commentLikes
	
	// 获取粉丝数量
	s.db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&stats.FollowersCount)
	
	// 获取关注数量
	s.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&stats.FollowingCount)
	
	// 获取总浏览量
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews)
	
	// 计算加入天数
	stats.JoinedDays = int(time.Since(user.CreatedAt).Hours() / 24)
	if stats.JoinedDays < 1 {
		stats.JoinedDays = 1
	}
	
	return stats, nil
}

// GetActiveUsers 获取活跃用户
// 参数: limit - 限制数量
// 返回: []models.User - 活跃用户列表, error - 错误信息
func (s *userService) GetActiveUsers(limit int) ([]models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	var users []models.User
	
	// 根据最近活动时间排序获取活跃用户
	err := s.db.Where("status = ? AND last_login_at > ?", "active", time.Now().AddDate(0, 0, -30)).
		Order("last_login_at DESC").
		Limit(limit).
		Find(&users).Error
	
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

// 辅助方法

// ValidateUserData 验证用户数据
// 参数: user - 用户模型
// 返回: error - 验证错误信息
func ValidateUserData(user *models.User) error {
	if user == nil {
		return errors.New("用户信息不能为空")
	}
	
	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}
	
	if len(user.Email) < 5 || len(user.Email) > 100 {
		return errors.New("邮箱长度必须在5-100个字符之间")
	}
	
	// 简单的邮箱格式验证
	if !isValidEmail(user.Email) {
		return errors.New("邮箱格式不正确")
	}
	
	if user.PasswordHash != "" && len(user.PasswordHash) < 6 {
		return errors.New("密码长度不能少于6个字符")
	}
	
	return nil
}

// isValidEmail 简单的邮箱格式验证
// 参数: email - 邮箱地址
// 返回: bool - 是否有效
func isValidEmail(email string) bool {
	// 这里可以使用更复杂的正则表达式或第三方库进行验证
	// 简单验证：包含@符号且@前后都有字符
	if len(email) < 5 {
		return false
	}
	
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 { // 多个@符号
				return false
			}
			atIndex = i
		}
	}
	
	// 必须有@符号，且@前后都有字符
	return atIndex > 0 && atIndex < len(email)-1
}

// GetUserDisplayName 获取用户显示名称
// 参数: user - 用户模型
// 返回: string - 显示名称
func GetUserDisplayName(user *models.User) string {
	if user == nil {
		return "未知用户"
	}
	
	if user.Profile != nil && user.Profile.Nickname != "" {
		return user.Profile.Nickname
	}
	
	return user.Username
}

// FormatUserInfo 格式化用户信息用于显示
// 参数: user - 用户模型
// 返回: string - 格式化的用户信息
func FormatUserInfo(user *models.User) string {
	if user == nil {
		return "用户信息不可用"
	}
	
	displayName := GetUserDisplayName(user)
	joinedDays := int(time.Since(user.CreatedAt).Hours() / 24)
	
	return fmt.Sprintf("%s (ID: %d, 加入: %d天前, 状态: %s)", 
		displayName, user.ID, joinedDays, user.Status)
}
package repository

import (
	"errors"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问层接口
// 定义用户相关的数据库操作方法
type UserRepository interface {
	// 基本CRUD操作
	Create(user *models.User) error                              // 创建用户
	GetByID(id uint) (*models.User, error)                      // 根据ID获取用户
	GetByUsername(username string) (*models.User, error)        // 根据用户名获取用户
	GetByEmail(email string) (*models.User, error)              // 根据邮箱获取用户
	Update(user *models.User) error                             // 更新用户信息
	Delete(id uint) error                                       // 删除用户（软删除）
	HardDelete(id uint) error                                   // 硬删除用户
	
	// 查询操作
	List(offset, limit int) ([]models.User, error)             // 分页获取用户列表
	Search(keyword string, offset, limit int) ([]models.User, error) // 搜索用户
	GetActiveUsers(offset, limit int) ([]models.User, error)   // 获取活跃用户
	GetUsersByStatus(status string, offset, limit int) ([]models.User, error) // 根据状态获取用户
	GetUsersByRole(role string, offset, limit int) ([]models.User, error) // 根据角色获取用户
	
	// 统计操作
	Count() (int64, error)                                      // 获取用户总数
	CountByStatus(status string) (int64, error)                // 根据状态统计用户数
	CountByRole(role string) (int64, error)                    // 根据角色统计用户数
	CountActiveUsers(days int) (int64, error)                  // 统计活跃用户数
	CountNewUsers(days int) (int64, error)                     // 统计新用户数
	
	// 用户资料操作
	CreateProfile(profile *models.UserProfile) error           // 创建用户资料
	GetProfile(userID uint) (*models.UserProfile, error)       // 获取用户资料
	UpdateProfile(profile *models.UserProfile) error           // 更新用户资料
	DeleteProfile(userID uint) error                           // 删除用户资料
	
	// 关注关系操作
	Follow(followerID, followingID uint) error                  // 关注用户
	Unfollow(followerID, followingID uint) error                // 取消关注
	IsFollowing(followerID, followingID uint) (bool, error)     // 检查是否关注
	GetFollowers(userID uint, offset, limit int) ([]models.User, error) // 获取粉丝列表
	GetFollowing(userID uint, offset, limit int) ([]models.User, error) // 获取关注列表
	CountFollowers(userID uint) (int64, error)                  // 统计粉丝数
	CountFollowing(userID uint) (int64, error)                  // 统计关注数
	
	// 批量操作
	BatchCreate(users []models.User) error                      // 批量创建用户
	BatchUpdateStatus(userIDs []uint, status string) error     // 批量更新状态
	BatchDelete(userIDs []uint) error                          // 批量删除用户
	
	// 高级查询
	GetUserWithProfile(userID uint) (*models.User, error)      // 获取用户及其资料
	GetUserWithStats(userID uint) (*UserWithStats, error)      // 获取用户及其统计信息
	GetTopActiveUsers(limit int, days int) ([]models.User, error) // 获取最活跃用户
	GetRecentUsers(limit int) ([]models.User, error)           // 获取最近注册用户
}

// userRepository 用户数据访问层实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问层实例
// 参数: db - 数据库连接
// 返回: UserRepository - 用户数据访问层接口实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// 辅助数据结构

// UserWithStats 用户及其统计信息
type UserWithStats struct {
	models.User
	PostCount     int `json:"post_count"`     // 文章数量
	CommentCount  int `json:"comment_count"`  // 评论数量
	FollowerCount int `json:"follower_count"` // 粉丝数量
	FollowingCount int `json:"following_count"` // 关注数量
	LikeCount     int `json:"like_count"`     // 获得点赞数
	ViewCount     int `json:"view_count"`     // 文章总浏览数
}

// 基本CRUD操作实现

// Create 创建用户
// 参数: user - 用户对象
// 返回: error - 错误信息
func (r *userRepository) Create(user *models.User) error {
	if user == nil {
		return errors.New("用户对象不能为空")
	}
	
	// 检查用户名是否已存在
	var count int64
	r.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	r.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已存在")
	}
	
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
// 参数: id - 用户ID
// 返回: *models.User - 用户对象, error - 错误信息
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	user := &models.User{}
	err := r.db.First(user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// GetByUsername 根据用户名获取用户
// 参数: username - 用户名
// 返回: *models.User - 用户对象, error - 错误信息
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	
	user := &models.User{}
	err := r.db.Where("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// GetByEmail 根据邮箱获取用户
// 参数: email - 邮箱地址
// 返回: *models.User - 用户对象, error - 错误信息
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	
	user := &models.User{}
	err := r.db.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// Update 更新用户信息
// 参数: user - 用户对象
// 返回: error - 错误信息
func (r *userRepository) Update(user *models.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("用户对象或ID不能为空")
	}
	
	// 检查用户是否存在
	var count int64
	r.db.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	
	// 检查用户名是否被其他用户使用
	r.db.Model(&models.User{}).Where("username = ? AND id != ?", user.Username, user.ID).Count(&count)
	if count > 0 {
		return errors.New("用户名已被其他用户使用")
	}
	
	// 检查邮箱是否被其他用户使用
	r.db.Model(&models.User{}).Where("email = ? AND id != ?", user.Email, user.ID).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被其他用户使用")
	}
	
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
// 参数: id - 用户ID
// 返回: error - 错误信息
func (r *userRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	// 检查用户是否存在
	var count int64
	r.db.Model(&models.User{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	
	return r.db.Delete(&models.User{}, id).Error
}

// HardDelete 硬删除用户
// 参数: id - 用户ID
// 返回: error - 错误信息
func (r *userRepository) HardDelete(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}
	
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}

// 查询操作实现

// List 分页获取用户列表
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.User - 用户列表, error - 错误信息
func (r *userRepository) List(offset, limit int) ([]models.User, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

// Search 搜索用户
// 参数: keyword - 搜索关键词, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 用户列表, error - 错误信息
func (r *userRepository) Search(keyword string, offset, limit int) ([]models.User, error) {
	if keyword == "" {
		return r.List(offset, limit)
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	keyword = "%" + keyword + "%"
	err := r.db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", keyword, keyword, keyword).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

// GetActiveUsers 获取活跃用户
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.User - 用户列表, error - 错误信息
func (r *userRepository) GetActiveUsers(offset, limit int) ([]models.User, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	err := r.db.Where("status = ?", "active").
		Offset(offset).Limit(limit).Order("last_login_at DESC").Find(&users).Error
	return users, err
}

// GetUsersByStatus 根据状态获取用户
// 参数: status - 用户状态, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 用户列表, error - 错误信息
func (r *userRepository) GetUsersByStatus(status string, offset, limit int) ([]models.User, error) {
	if status == "" {
		return nil, errors.New("状态不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	err := r.db.Where("status = ?", status).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

// GetUsersByRole 根据角色获取用户
// 参数: role - 用户角色, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 用户列表, error - 错误信息
func (r *userRepository) GetUsersByRole(role string, offset, limit int) ([]models.User, error) {
	if role == "" {
		return nil, errors.New("角色不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	err := r.db.Where("role = ?", role).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

// 统计操作实现

// Count 获取用户总数
// 返回: int64 - 用户总数, error - 错误信息
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计用户数
// 参数: status - 用户状态
// 返回: int64 - 用户数量, error - 错误信息
func (r *userRepository) CountByStatus(status string) (int64, error) {
	if status == "" {
		return 0, errors.New("状态不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountByRole 根据角色统计用户数
// 参数: role - 用户角色
// 返回: int64 - 用户数量, error - 错误信息
func (r *userRepository) CountByRole(role string) (int64, error) {
	if role == "" {
		return 0, errors.New("角色不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

// CountActiveUsers 统计活跃用户数
// 参数: days - 活跃天数（最近N天内有登录）
// 返回: int64 - 活跃用户数, error - 错误信息
func (r *userRepository) CountActiveUsers(days int) (int64, error) {
	if days <= 0 {
		days = 30
	}
	
	var count int64
	startDate := time.Now().AddDate(0, 0, -days)
	err := r.db.Model(&models.User{}).Where("last_login_at >= ? AND status = ?", startDate, "active").Count(&count).Error
	return count, err
}

// CountNewUsers 统计新用户数
// 参数: days - 统计天数（最近N天内注册）
// 返回: int64 - 新用户数, error - 错误信息
func (r *userRepository) CountNewUsers(days int) (int64, error) {
	if days <= 0 {
		days = 7
	}
	
	var count int64
	startDate := time.Now().AddDate(0, 0, -days)
	err := r.db.Model(&models.User{}).Where("created_at >= ?", startDate).Count(&count).Error
	return count, err
}

// 用户资料操作实现

// CreateProfile 创建用户资料
// 参数: profile - 用户资料对象
// 返回: error - 错误信息
func (r *userRepository) CreateProfile(profile *models.UserProfile) error {
	if profile == nil || profile.UserID == 0 {
		return errors.New("用户资料对象或用户ID不能为空")
	}
	
	// 检查用户是否存在
	var count int64
	r.db.Model(&models.User{}).Where("id = ?", profile.UserID).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	
	// 检查是否已有资料
	r.db.Model(&models.UserProfile{}).Where("user_id = ?", profile.UserID).Count(&count)
	if count > 0 {
		return errors.New("用户资料已存在")
	}
	
	return r.db.Create(profile).Error
}

// GetProfile 获取用户资料
// 参数: userID - 用户ID
// 返回: *models.UserProfile - 用户资料对象, error - 错误信息
func (r *userRepository) GetProfile(userID uint) (*models.UserProfile, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	profile := &models.UserProfile{}
	err := r.db.Where("user_id = ?", userID).First(profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户资料不存在")
		}
		return nil, err
	}
	
	return profile, nil
}

// UpdateProfile 更新用户资料
// 参数: profile - 用户资料对象
// 返回: error - 错误信息
func (r *userRepository) UpdateProfile(profile *models.UserProfile) error {
	if profile == nil || profile.UserID == 0 {
		return errors.New("用户资料对象或用户ID不能为空")
	}
	
	// 检查资料是否存在
	var count int64
	r.db.Model(&models.UserProfile{}).Where("user_id = ?", profile.UserID).Count(&count)
	if count == 0 {
		return errors.New("用户资料不存在")
	}
	
	return r.db.Where("user_id = ?", profile.UserID).Save(profile).Error
}

// DeleteProfile 删除用户资料
// 参数: userID - 用户ID
// 返回: error - 错误信息
func (r *userRepository) DeleteProfile(userID uint) error {
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	return r.db.Where("user_id = ?", userID).Delete(&models.UserProfile{}).Error
}

// 关注关系操作实现

// Follow 关注用户
// 参数: followerID - 关注者ID, followingID - 被关注者ID
// 返回: error - 错误信息
func (r *userRepository) Follow(followerID, followingID uint) error {
	if followerID == 0 || followingID == 0 {
		return errors.New("关注者ID和被关注者ID不能为空")
	}
	
	if followerID == followingID {
		return errors.New("不能关注自己")
	}
	
	// 检查用户是否存在
	var count int64
	r.db.Model(&models.User{}).Where("id IN (?, ?)", followerID, followingID).Count(&count)
	if count != 2 {
		return errors.New("用户不存在")
	}
	
	// 检查是否已关注
	r.db.Model(&models.Follow{}).Where("follower_id = ? AND followed_id = ?", followerID, followingID).Count(&count)
	if count > 0 {
		return errors.New("已经关注该用户")
	}
	
	follow := &models.Follow{
		FollowerID: followerID,
		FollowedID: followingID,
	}
	
	return r.db.Create(follow).Error
}

// Unfollow 取消关注
// 参数: followerID - 关注者ID, followingID - 被关注者ID
// 返回: error - 错误信息
func (r *userRepository) Unfollow(followerID, followingID uint) error {
	if followerID == 0 || followingID == 0 {
		return errors.New("关注者ID和被关注者ID不能为空")
	}
	
	return r.db.Where("follower_id = ? AND followed_id = ?", followerID, followingID).Delete(&models.Follow{}).Error
}

// IsFollowing 检查是否关注
// 参数: followerID - 关注者ID, followingID - 被关注者ID
// 返回: bool - 是否关注, error - 错误信息
func (r *userRepository) IsFollowing(followerID, followingID uint) (bool, error) {
	if followerID == 0 || followingID == 0 {
		return false, errors.New("关注者ID和被关注者ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Follow{}).Where("follower_id = ? AND followed_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

// GetFollowers 获取粉丝列表
// 参数: userID - 用户ID, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 粉丝列表, error - 错误信息
func (r *userRepository) GetFollowers(userID uint, offset, limit int) ([]models.User, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	err := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Offset(offset).Limit(limit).Order("follows.created_at DESC").
		Find(&users).Error
	
	return users, err
}

// GetFollowing 获取关注列表
// 参数: userID - 用户ID, offset - 偏移量, limit - 限制数量
// 返回: []models.User - 关注列表, error - 错误信息
func (r *userRepository) GetFollowing(userID uint, offset, limit int) ([]models.User, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var users []models.User
	err := r.db.Table("users").
		Joins("JOIN follows ON users.id = follows.followed_id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).Limit(limit).Order("follows.created_at DESC").
		Find(&users).Error
	
	return users, err
}

// CountFollowers 统计粉丝数
// 参数: userID - 用户ID
// 返回: int64 - 粉丝数, error - 错误信息
func (r *userRepository) CountFollowers(userID uint) (int64, error) {
	if userID == 0 {
		return 0, errors.New("用户ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

// CountFollowing 统计关注数
// 参数: userID - 用户ID
// 返回: int64 - 关注数, error - 错误信息
func (r *userRepository) CountFollowing(userID uint) (int64, error) {
	if userID == 0 {
		return 0, errors.New("用户ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

// 批量操作实现

// BatchCreate 批量创建用户
// 参数: users - 用户列表
// 返回: error - 错误信息
func (r *userRepository) BatchCreate(users []models.User) error {
	if len(users) == 0 {
		return errors.New("用户列表不能为空")
	}
	
	// 检查用户名和邮箱唯一性
	usernames := make([]string, len(users))
	emails := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
		emails[i] = user.Email
	}
	
	var count int64
	r.db.Model(&models.User{}).Where("username IN (?) OR email IN (?)", usernames, emails).Count(&count)
	if count > 0 {
		return errors.New("存在重复的用户名或邮箱")
	}
	
	return r.db.CreateInBatches(users, 100).Error
}

// BatchUpdateStatus 批量更新状态
// 参数: userIDs - 用户ID列表, status - 新状态
// 返回: error - 错误信息
func (r *userRepository) BatchUpdateStatus(userIDs []uint, status string) error {
	if len(userIDs) == 0 {
		return errors.New("用户ID列表不能为空")
	}
	if status == "" {
		return errors.New("状态不能为空")
	}
	
	return r.db.Model(&models.User{}).Where("id IN (?)", userIDs).Update("status", status).Error
}

// BatchDelete 批量删除用户
// 参数: userIDs - 用户ID列表
// 返回: error - 错误信息
func (r *userRepository) BatchDelete(userIDs []uint) error {
	if len(userIDs) == 0 {
		return errors.New("用户ID列表不能为空")
	}
	
	return r.db.Delete(&models.User{}, userIDs).Error
}

// 高级查询实现

// GetUserWithProfile 获取用户及其资料
// 参数: userID - 用户ID
// 返回: *models.User - 用户对象（包含资料）, error - 错误信息
func (r *userRepository) GetUserWithProfile(userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	user := &models.User{}
	err := r.db.Preload("Profile").First(user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	
	return user, nil
}

// GetUserWithStats 获取用户及其统计信息
// 参数: userID - 用户ID
// 返回: *UserWithStats - 用户及统计信息, error - 错误信息
func (r *userRepository) GetUserWithStats(userID uint) (*UserWithStats, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	// 获取用户基本信息
	user, err := r.GetByID(userID)
	if err != nil {
		return nil, err
	}
	
	userWithStats := &UserWithStats{
		User: *user,
	}
	
	// 统计文章数
	var postCount int64
	r.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&postCount)
	userWithStats.PostCount = int(postCount)
	
	// 统计评论数
	var commentCount int64
	r.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount)
	userWithStats.CommentCount = int(commentCount)
	
	// 统计粉丝数
	followerCount, _ := r.CountFollowers(userID)
	userWithStats.FollowerCount = int(followerCount)
	
	// 统计关注数
	followingCount, _ := r.CountFollowing(userID)
	userWithStats.FollowingCount = int(followingCount)
	
	// 统计获得点赞数
	var likeCount int64
	r.db.Table("likes").
		Joins("JOIN posts ON likes.target_id = posts.id AND likes.target_type = 'post'").
		Where("posts.author_id = ?", userID).Count(&likeCount)
	userWithStats.LikeCount = int(likeCount)
	
	// 统计文章总浏览数
	var viewCount int64
	r.db.Model(&models.Post{}).Where("author_id = ?", userID).
		Select("COALESCE(SUM(view_count), 0)").Scan(&viewCount)
	userWithStats.ViewCount = int(viewCount)
	
	return userWithStats, nil
}

// GetTopActiveUsers 获取最活跃用户
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.User - 活跃用户列表, error - 错误信息
func (r *userRepository) GetTopActiveUsers(limit int, days int) ([]models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 30
	}
	
	var users []models.User
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Table("users").
		Select(`
			users.*,
			(COALESCE(post_counts.post_count, 0) * 5 + COALESCE(comment_counts.comment_count, 0) * 2) as activity_score
		`).
		Joins(`LEFT JOIN (
			SELECT author_id, COUNT(*) as post_count 
			FROM posts 
			WHERE created_at >= ?
			GROUP BY author_id
		) post_counts ON users.id = post_counts.author_id`, startDate).
		Joins(`LEFT JOIN (
			SELECT user_id, COUNT(*) as comment_count 
			FROM comments 
			WHERE created_at >= ?
			GROUP BY user_id
		) comment_counts ON users.id = comment_counts.user_id`, startDate).
		Where("users.status = ?", "active").
		Having("activity_score > 0").
		Order("activity_score DESC").
		Limit(limit).
		Find(&users).Error
	
	return users, err
}

// GetRecentUsers 获取最近注册用户
// 参数: limit - 限制数量
// 返回: []models.User - 最近注册用户列表, error - 错误信息
func (r *userRepository) GetRecentUsers(limit int) ([]models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	var users []models.User
	err := r.db.Order("created_at DESC").Limit(limit).Find(&users).Error
	return users, err
}
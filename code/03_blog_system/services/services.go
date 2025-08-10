// 03_blog_system/services/services.go - 服务层实现
// 对应文档：02_GORM背景示例_博客系统实战.md

package services

import (
	"blog-system/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 全局服务实例
var (
	UserService    *userService
	PostService    *postService
	CommentService *commentService
	CategoryService *categoryService
	TagService     *tagService
)

// InitServices 初始化所有服务
func InitServices(db *gorm.DB) {
	UserService = &userService{db: db}
	PostService = &postService{db: db}
	CommentService = &commentService{db: db}
	CategoryService = &categoryService{db: db}
	TagService = &tagService{db: db}
}

// ===== 用户服务 =====

type userService struct {
	db *gorm.DB
}

// RegisterUser 用户注册
func (s *userService) RegisterUser(username, email, password string) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名或邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Nickname: username,
		Status:   "active",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &user, nil
}

// LoginUser 用户登录
func (s *userService) LoginUser(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("用户账号已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 更新登录信息
	now := time.Now()
	s.db.Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
		"login_count":   gorm.Expr("login_count + ?", 1),
	})

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Profile").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// UpdateUserProfile 更新用户资料
func (s *userService) UpdateUserProfile(userID uint, profile map[string]interface{}) error {
	return s.db.Model(&models.Profile{}).Where("user_id = ?", userID).Updates(profile).Error
}

// ===== 文章服务 =====

type postService struct {
	db *gorm.DB
}

// CreatePost 创建文章
func (s *postService) CreatePost(post *models.Post) error {
	return s.db.Create(post).Error
}

// GetPostByID 根据ID获取文章
func (s *postService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := s.db.Preload("User").Preload("Category").Preload("Tags").Preload("Comments.User").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	// 增加浏览量
	s.db.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return &post, nil
}

// GetPostBySlug 根据Slug获取文章
func (s *postService) GetPostBySlug(slug string) (*models.Post, error) {
	var post models.Post
	if err := s.db.Preload("User").Preload("Category").Preload("Tags").Where("slug = ?", slug).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	// 增加浏览量
	s.db.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return &post, nil
}

// GetPosts 获取文章列表
func (s *postService) GetPosts(page, pageSize int, filters map[string]interface{}) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := s.db.Model(&models.Post{})

	// 应用过滤条件
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if categoryID, ok := filters["category_id"]; ok {
		query = query.Where("category_id = ?", categoryID)
	}
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if keyword, ok := filters["keyword"]; ok {
		keyword = fmt.Sprintf("%%%s%%", keyword)
		query = query.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计文章数量失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Category").Preload("Tags").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("查询文章列表失败: %w", err)
	}

	return posts, total, nil
}

// UpdatePost 更新文章
func (s *postService) UpdatePost(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePost 删除文章
func (s *postService) DeletePost(id uint) error {
	return s.db.Delete(&models.Post{}, id).Error
}

// PublishPost 发布文章
func (s *postService) PublishPost(id uint) error {
	now := time.Now()
	return s.db.Model(&models.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	}).Error
}

// ===== 评论服务 =====

type commentService struct {
	db *gorm.DB
}

// CreateComment 创建评论
func (s *commentService) CreateComment(comment *models.Comment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 创建评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 更新文章评论数
		return tx.Model(&models.Post{}).Where("id = ?", comment.PostID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	})
}

// GetCommentsByPostID 获取文章的评论列表
func (s *commentService) GetCommentsByPostID(postID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := s.db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, "approved")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计评论数量失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Replies.User").
		Where("parent_id IS NULL"). // 只获取顶级评论
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, fmt.Errorf("查询评论列表失败: %w", err)
	}

	return comments, total, nil
}

// ApproveComment 审核通过评论
func (s *commentService) ApproveComment(id uint) error {
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "approved").Error
}

// RejectComment 拒绝评论
func (s *commentService) RejectComment(id uint) error {
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "rejected").Error
}

// ===== 分类服务 =====

type categoryService struct {
	db *gorm.DB
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(category *models.Category) error {
	return s.db.Create(category).Error
}

// GetCategories 获取分类列表
func (s *categoryService) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Order("sort_order ASC, created_at DESC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("查询分类列表失败: %w", err)
	}
	return categories, nil
}

// GetCategoryBySlug 根据Slug获取分类
func (s *categoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	if err := s.db.Where("slug = ?", slug).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("分类不存在")
		}
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}
	return &category, nil
}

// ===== 标签服务 =====

type tagService struct {
	db *gorm.DB
}

// CreateTag 创建标签
func (s *tagService) CreateTag(tag *models.Tag) error {
	return s.db.Create(tag).Error
}

// GetTags 获取标签列表
func (s *tagService) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Order("post_count DESC, created_at DESC").Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %w", err)
	}
	return tags, nil
}

// GetTagsByNames 根据名称获取标签（如果不存在则创建）
func (s *tagService) GetTagsByNames(names []string) ([]models.Tag, error) {
	var tags []models.Tag

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		var tag models.Tag
		if err := s.db.Where("name = ?", name).First(&tag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 标签不存在，创建新标签
				tag = models.Tag{
					Name: name,
					Slug: strings.ToLower(strings.ReplaceAll(name, " ", "-")),
				}
				if err := s.db.Create(&tag).Error; err != nil {
					return nil, fmt.Errorf("创建标签失败: %w", err)
				}
			} else {
				return nil, fmt.Errorf("查询标签失败: %w", err)
			}
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetPopularTags 获取热门标签
func (s *tagService) GetPopularTags(limit int) ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Where("post_count > 0").Order("post_count DESC").Limit(limit).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("查询热门标签失败: %w", err)
	}
	return tags, nil
}
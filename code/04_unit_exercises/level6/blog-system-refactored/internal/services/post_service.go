package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// PostService 文章服务接口
// 定义文章相关的业务操作
type PostService interface {
	// 文章基本操作
	CreatePost(post *models.Post) error                    // 创建文章
	GetPostByID(id uint) (*models.Post, error)             // 根据ID获取文章
	GetPostBySlug(slug string) (*models.Post, error)       // 根据slug获取文章
	UpdatePost(post *models.Post) error                    // 更新文章
	DeletePost(id uint) error                              // 删除文章
	ListPosts(offset, limit int, filters PostFilters) ([]models.Post, int64, error) // 分页获取文章列表
	
	// 文章状态操作
	PublishPost(id uint) error                             // 发布文章
	UnpublishPost(id uint) error                           // 取消发布
	ArchivePost(id uint) error                             // 归档文章
	PinPost(id uint) error                                 // 置顶文章
	UnpinPost(id uint) error                               // 取消置顶
	
	// 文章统计操作
	IncrementViewCount(id uint) error                      // 增加浏览次数
	GetPostStats(id uint) (*PostStats, error)             // 获取文章统计
	
	// 文章搜索和筛选
	SearchPosts(keyword string, offset, limit int) ([]models.Post, int64, error) // 搜索文章
	GetPostsByCategory(categoryID uint, offset, limit int) ([]models.Post, int64, error) // 按分类获取文章
	GetPostsByTag(tagID uint, offset, limit int) ([]models.Post, int64, error) // 按标签获取文章
	GetPostsByAuthor(authorID uint, offset, limit int) ([]models.Post, int64, error) // 按作者获取文章
	
	// 热门和推荐
	GetPopularPosts(limit int, days int) ([]models.Post, error) // 获取热门文章
	GetRecentPosts(limit int) ([]models.Post, error)       // 获取最新文章
	GetRecommendedPosts(userID uint, limit int) ([]models.Post, error) // 获取推荐文章
	
	// 文章标签管理
	AddTagsToPost(postID uint, tagIDs []uint) error        // 为文章添加标签
	RemoveTagsFromPost(postID uint, tagIDs []uint) error   // 从文章移除标签
	UpdatePostTags(postID uint, tagIDs []uint) error       // 更新文章标签
	
	// 分类管理
	GetAllCategories() ([]models.Category, error)          // 获取所有分类
	
	// 标签管理
	GetPopularTags(limit int) ([]models.Tag, error)        // 获取热门标签
}

// postService 文章服务实现
type postService struct {
	db *gorm.DB
}

// NewPostService 创建文章服务实例
// 参数: db - 数据库连接
// 返回: PostService - 文章服务接口实例
func NewPostService(db *gorm.DB) PostService {
	return &postService{
		db: db,
	}
}

// PostFilters 文章筛选条件
type PostFilters struct {
	Status     string `json:"status"`      // 状态筛选
	CategoryID uint   `json:"category_id"` // 分类筛选
	AuthorID   uint   `json:"author_id"`   // 作者筛选
	TagID      uint   `json:"tag_id"`      // 标签筛选
	Keyword    string `json:"keyword"`     // 关键词搜索
	StartDate  *time.Time `json:"start_date"` // 开始日期
	EndDate    *time.Time `json:"end_date"`   // 结束日期
	OrderBy    string `json:"order_by"`    // 排序字段
	OrderDir   string `json:"order_dir"`   // 排序方向
}

// PostStats 文章统计信息
type PostStats struct {
	TotalPosts     int64 `json:"total_posts"`     // 总文章数
	PublishedPosts int64 `json:"published_posts"` // 已发布文章数
	DraftPosts     int64 `json:"draft_posts"`     // 草稿文章数
	TotalViews     int64 `json:"total_views"`     // 总浏览量
	TotalLikes     int64 `json:"total_likes"`     // 总点赞数
	ViewCount      int   `json:"view_count"`      // 浏览次数
	LikeCount      int   `json:"like_count"`      // 点赞次数
	CommentCount   int   `json:"comment_count"`   // 评论次数
	ShareCount     int   `json:"share_count"`     // 分享次数
	ReadTime       int   `json:"read_time"`       // 预估阅读时间（分钟）
	WordCount      int   `json:"word_count"`      // 字数统计
	PublishedAt    *time.Time `json:"published_at,omitempty"` // 发布时间
	LastViewAt     *time.Time `json:"last_view_at,omitempty"` // 最后浏览时间
}

// 文章基本操作实现

// CreatePost 创建文章
// 参数: post - 文章模型
// 返回: error - 错误信息
func (s *postService) CreatePost(post *models.Post) error {
	if post == nil {
		return errors.New("文章信息不能为空")
	}
	
	// 验证必填字段
	if err := s.validatePostData(post); err != nil {
		return err
	}
	
	// 生成slug（如果没有提供）
	if post.Slug == "" {
		post.Slug = s.generateSlug(post.Title)
	}
	
	// 检查slug是否重复
	if err := s.checkSlugUnique(post.Slug, 0); err != nil {
		return err
	}
	
	// 设置默认值
	if post.Status == models.PostStatus(0) {
		post.Status = models.PostStatusDraft
	}
	if post.ViewCount == 0 {
		post.ViewCount = 0
	}
	
	// 如果状态为已发布，设置发布时间
	if post.Status == models.PostStatusPublished && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}
	
	return s.db.Create(post).Error
}

// GetPostByID 根据ID获取文章
// 参数: id - 文章ID
// 返回: *models.Post - 文章模型, error - 错误信息
func (s *postService) GetPostByID(id uint) (*models.Post, error) {
	if id == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	post := &models.Post{}
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").Preload("Meta").First(post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	return post, nil
}

// GetPostBySlug 根据slug获取文章
// 参数: slug - 文章slug
// 返回: *models.Post - 文章模型, error - 错误信息
func (s *postService) GetPostBySlug(slug string) (*models.Post, error) {
	if slug == "" {
		return nil, errors.New("文章slug不能为空")
	}
	
	post := &models.Post{}
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").Preload("Meta").Where("slug = ?", slug).First(post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	return post, nil
}

// UpdatePost 更新文章
// 参数: post - 文章模型
// 返回: error - 错误信息
func (s *postService) UpdatePost(post *models.Post) error {
	if post == nil || post.ID == 0 {
		return errors.New("文章信息不完整")
	}
	
	// 验证数据
	if err := s.validatePostData(post); err != nil {
		return err
	}
	
	// 检查文章是否存在
	existingPost := &models.Post{}
	if err := s.db.First(existingPost, post.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	// 如果更新slug，检查是否重复
	if post.Slug != "" && post.Slug != existingPost.Slug {
		if err := s.checkSlugUnique(post.Slug, post.ID); err != nil {
			return err
		}
	}
	
	// 如果状态从非发布改为发布，设置发布时间
	if post.Status == models.PostStatusPublished && existingPost.Status != models.PostStatusPublished && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}
	
	return s.db.Save(post).Error
}

// DeletePost 删除文章（软删除）
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) DeletePost(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	// 检查文章是否存在
	post := &models.Post{}
	if err := s.db.First(post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	return s.db.Delete(post).Error
}

// ListPosts 分页获取文章列表
// 参数: offset - 偏移量, limit - 限制数量, filters - 筛选条件
// 返回: []models.Post - 文章列表, int64 - 总数量, error - 错误信息
func (s *postService) ListPosts(offset, limit int, filters PostFilters) ([]models.Post, int64, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	var total int64
	
	// 构建查询
	query := s.db.Model(&models.Post{})
	
	// 应用筛选条件
	query = s.applyPostFilters(query, filters)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 应用排序
	orderBy := "created_at"
	orderDir := "DESC"
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}
	if filters.OrderDir != "" {
		orderDir = strings.ToUpper(filters.OrderDir)
	}
	
	// 获取文章列表
	err := query.Preload("Author").Preload("Category").Preload("Tags").
		Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", orderBy, orderDir)).
		Find(&posts).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	return posts, total, nil
}

// 文章状态操作实现

// PublishPost 发布文章
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) PublishPost(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	now := time.Now()
	updates := map[string]interface{}{
		"status":       "published",
		"published_at": &now,
	}
	
	return s.db.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

// UnpublishPost 取消发布
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) UnpublishPost(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	updates := map[string]interface{}{
		"status":       "draft",
		"published_at": nil,
	}
	
	return s.db.Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

// ArchivePost 归档文章
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) ArchivePost(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	return s.db.Model(&models.Post{}).Where("id = ?", id).Update("status", "archived").Error
}

// PinPost 置顶文章
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) PinPost(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	return s.db.Model(&models.Post{}).Where("id = ?", id).Update("is_pinned", true).Error
}

// UnpinPost 取消置顶
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) UnpinPost(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	return s.db.Model(&models.Post{}).Where("id = ?", id).Update("is_pinned", false).Error
}

// 文章统计操作实现

// IncrementViewCount 增加浏览次数
// 参数: id - 文章ID
// 返回: error - 错误信息
func (s *postService) IncrementViewCount(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	return s.db.Model(&models.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// GetPostStats 获取文章统计
// 参数: id - 文章ID
// 返回: *PostStats - 文章统计信息, error - 错误信息
func (s *postService) GetPostStats(id uint) (*PostStats, error) {
	if id == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	// 获取文章基本信息
	post := &models.Post{}
	if err := s.db.First(post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	stats := &PostStats{
		ViewCount:   post.ViewCount,

		PublishedAt: post.PublishedAt,
	}
	
	// 获取点赞数
	var likeCount int64
	s.db.Model(&models.Like{}).Where("target_type = ? AND target_id = ?", "post", id).Count(&likeCount)
	stats.LikeCount = int(likeCount)
	
	// 获取评论数
	var commentCount int64
	s.db.Model(&models.Comment{}).Where("post_id = ?", id).Count(&commentCount)
	stats.CommentCount = int(commentCount)
	
	// TODO: 获取分享数（需要实现分享功能）
	stats.ShareCount = 0
	
	// TODO: 获取最后浏览时间（需要实现浏览记录功能）
	
	return stats, nil
}

// 文章搜索和筛选实现

// SearchPosts 搜索文章
// 参数: keyword - 搜索关键词, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, int64 - 总数量, error - 错误信息
func (s *postService) SearchPosts(keyword string, offset, limit int) ([]models.Post, int64, error) {
	if keyword == "" {
		return nil, 0, errors.New("搜索关键词不能为空")
	}
	
	filters := PostFilters{
		Keyword: keyword,
		Status:  "published", // 只搜索已发布的文章
	}
	
	return s.ListPosts(offset, limit, filters)
}

// GetPostsByCategory 按分类获取文章
// 参数: categoryID - 分类ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, int64 - 总数量, error - 错误信息
func (s *postService) GetPostsByCategory(categoryID uint, offset, limit int) ([]models.Post, int64, error) {
	if categoryID == 0 {
		return nil, 0, errors.New("分类ID不能为空")
	}
	
	filters := PostFilters{
		CategoryID: categoryID,
		Status:     "published",
	}
	
	return s.ListPosts(offset, limit, filters)
}

// GetPostsByTag 按标签获取文章
// 参数: tagID - 标签ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, int64 - 总数量, error - 错误信息
func (s *postService) GetPostsByTag(tagID uint, offset, limit int) ([]models.Post, int64, error) {
	if tagID == 0 {
		return nil, 0, errors.New("标签ID不能为空")
	}
	
	filters := PostFilters{
		TagID:  tagID,
		Status: "published",
	}
	
	return s.ListPosts(offset, limit, filters)
}

// GetPostsByAuthor 按作者获取文章
// 参数: authorID - 作者ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, int64 - 总数量, error - 错误信息
func (s *postService) GetPostsByAuthor(authorID uint, offset, limit int) ([]models.Post, int64, error) {
	if authorID == 0 {
		return nil, 0, errors.New("作者ID不能为空")
	}
	
	filters := PostFilters{
		AuthorID: authorID,
		Status:   "published",
	}
	
	return s.ListPosts(offset, limit, filters)
}

// 热门和推荐实现

// GetPopularPosts 获取热门文章
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.Post - 热门文章列表, error - 错误信息
func (s *postService) GetPopularPosts(limit int, days int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 7 // 默认7天
	}
	
	var posts []models.Post
	
	// 根据浏览量、点赞数等综合排序
	startDate := time.Now().AddDate(0, 0, -days)
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("status = ? AND published_at >= ?", "published", startDate).
		Order("view_count DESC, (SELECT COUNT(*) FROM likes WHERE target_type = 'post' AND target_id = posts.id) DESC").
		Limit(limit).
		Find(&posts).Error
	
	if err != nil {
		return nil, err
	}
	
	return posts, nil
}

// GetRecentPosts 获取最新文章
// 参数: limit - 限制数量
// 返回: []models.Post - 最新文章列表, error - 错误信息
func (s *postService) GetRecentPosts(limit int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	var posts []models.Post
	
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	
	if err != nil {
		return nil, err
	}
	
	return posts, nil
}

// GetRecommendedPosts 获取推荐文章
// 参数: userID - 用户ID, limit - 限制数量
// 返回: []models.Post - 推荐文章列表, error - 错误信息
func (s *postService) GetRecommendedPosts(userID uint, limit int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	// 简单的推荐算法：基于用户关注的作者和喜欢的分类
	var posts []models.Post
	
	if userID == 0 {
		// 未登录用户，返回热门文章
		return s.GetPopularPosts(limit, 30)
	}
	
	// 获取用户关注的作者的文章
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN follows ON posts.author_id = follows.following_id").
		Where("follows.follower_id = ? AND posts.status = ?", userID, "published").
		Order("posts.published_at DESC").
		Limit(limit).
		Find(&posts).Error
	
	if err != nil {
		return nil, err
	}
	
	// 如果关注的作者文章不够，补充热门文章
	if len(posts) < limit {
		remainingLimit := limit - len(posts)
		popularPosts, err := s.GetPopularPosts(remainingLimit, 30)
		if err == nil {
			posts = append(posts, popularPosts...)
		}
	}
	
	return posts, nil
}

// 文章标签管理实现

// AddTagsToPost 为文章添加标签
// 参数: postID - 文章ID, tagIDs - 标签ID列表
// 返回: error - 错误信息
func (s *postService) AddTagsToPost(postID uint, tagIDs []uint) error {
	if postID == 0 {
		return errors.New("文章ID不能为空")
	}
	if len(tagIDs) == 0 {
		return errors.New("标签ID列表不能为空")
	}
	
	// 检查文章是否存在
	post := &models.Post{}
	if err := s.db.First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	// 检查标签是否存在
	var existingTagCount int64
	s.db.Model(&models.Tag{}).Where("id IN ?", tagIDs).Count(&existingTagCount)
	if int(existingTagCount) != len(tagIDs) {
		return errors.New("部分标签不存在")
	}
	
	// 获取标签
	var tags []models.Tag
	if err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}
	
	// 添加关联
	return s.db.Model(post).Association("Tags").Append(&tags)
}

// RemoveTagsFromPost 从文章移除标签
// 参数: postID - 文章ID, tagIDs - 标签ID列表
// 返回: error - 错误信息
func (s *postService) RemoveTagsFromPost(postID uint, tagIDs []uint) error {
	if postID == 0 {
		return errors.New("文章ID不能为空")
	}
	if len(tagIDs) == 0 {
		return errors.New("标签ID列表不能为空")
	}
	
	// 检查文章是否存在
	post := &models.Post{}
	if err := s.db.First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	// 获取要移除的标签
	var tags []models.Tag
	if err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}
	
	// 移除关联
	return s.db.Model(post).Association("Tags").Delete(&tags)
}

// UpdatePostTags 更新文章标签
// 参数: postID - 文章ID, tagIDs - 新的标签ID列表
// 返回: error - 错误信息
func (s *postService) UpdatePostTags(postID uint, tagIDs []uint) error {
	if postID == 0 {
		return errors.New("文章ID不能为空")
	}
	
	// 检查文章是否存在
	post := &models.Post{}
	if err := s.db.First(post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}
	
	if len(tagIDs) == 0 {
		// 清空所有标签
		return s.db.Model(post).Association("Tags").Clear()
	}
	
	// 检查标签是否存在
	var existingTagCount int64
	s.db.Model(&models.Tag{}).Where("id IN ?", tagIDs).Count(&existingTagCount)
	if int(existingTagCount) != len(tagIDs) {
		return errors.New("部分标签不存在")
	}
	
	// 获取新标签
	var tags []models.Tag
	if err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}
	
	// 替换所有标签
	return s.db.Model(post).Association("Tags").Replace(&tags)
}

// 辅助方法

// validatePostData 验证文章数据
// 参数: post - 文章模型
// 返回: error - 验证错误信息
func (s *postService) validatePostData(post *models.Post) error {
	if post.Title == "" {
		return errors.New("文章标题不能为空")
	}
	if len(post.Title) > 200 {
		return errors.New("文章标题不能超过200个字符")
	}
	if post.Content == "" {
		return errors.New("文章内容不能为空")
	}
	if post.AuthorID == 0 {
		return errors.New("文章作者不能为空")
	}
	if post.CategoryID == nil || *post.CategoryID == 0 {
		return errors.New("文章分类不能为空")
	}
	
	return nil
}

// generateSlug 生成文章slug
// 参数: title - 文章标题
// 返回: string - 生成的slug
func (s *postService) generateSlug(title string) string {
	// 简单的slug生成：转小写，替换空格为连字符
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	
	// 移除特殊字符（这里简化处理）
	allowedChars := "abcdefghijklmnopqrstuvwxyz0123456789-"
	var result strings.Builder
	for _, char := range slug {
		if strings.ContainsRune(allowedChars, char) {
			result.WriteRune(char)
		}
	}
	
	slug = result.String()
	
	// 限制长度
	if len(slug) > 100 {
		slug = slug[:100]
	}
	
	// 移除首尾的连字符
	slug = strings.Trim(slug, "-")
	
	if slug == "" {
		slug = fmt.Sprintf("post-%d", time.Now().Unix())
	}
	
	return slug
}

// checkSlugUnique 检查slug是否唯一
// 参数: slug - 要检查的slug, excludeID - 排除的文章ID（用于更新时）
// 返回: error - 错误信息
func (s *postService) checkSlugUnique(slug string, excludeID uint) error {
	var count int64
	query := s.db.Model(&models.Post{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	
	if count > 0 {
		return errors.New("文章slug已存在")
	}
	
	return nil
}

// calculateWordCount 计算字数
// 参数: content - 文章内容
// 返回: int - 字数
func (s *postService) calculateWordCount(content string) int {
	// 简单的字数统计：按空格分割
	words := strings.Fields(content)
	return len(words)
}

// calculateReadTime 计算预估阅读时间
// 参数: wordCount - 字数
// 返回: int - 阅读时间（分钟）
func (s *postService) calculateReadTime(wordCount int) int {
	// 假设每分钟阅读200个单词
	readTime := wordCount / 200
	if readTime < 1 {
		readTime = 1
	}
	return readTime
}

// applyPostFilters 应用文章筛选条件
// 参数: query - GORM查询对象, filters - 筛选条件
// 返回: *gorm.DB - 应用筛选后的查询对象
func (s *postService) applyPostFilters(query *gorm.DB, filters PostFilters) *gorm.DB {
	// 状态筛选
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	
	// 分类筛选
	if filters.CategoryID > 0 {
		query = query.Where("category_id = ?", filters.CategoryID)
	}
	
	// 作者筛选
	if filters.AuthorID > 0 {
		query = query.Where("author_id = ?", filters.AuthorID)
	}
	
	// 标签筛选
	if filters.TagID > 0 {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").Where("post_tags.tag_id = ?", filters.TagID)
	}
	
	// 关键词搜索
	if filters.Keyword != "" {
		keyword := "%" + filters.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR excerpt LIKE ?", keyword, keyword, keyword)
	}
	
	// 日期范围筛选
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	
	return query
}

// GetAllCategories 获取所有分类
// 返回: []models.Category - 分类列表, error - 错误信息
func (s *postService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	
	err := s.db.Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %w", err)
	}
	
	return categories, nil
}

// GetPopularTags 获取热门标签
// 参数: limit - 限制数量
// 返回: []models.Tag - 标签列表, error - 错误信息
func (s *postService) GetPopularTags(limit int) ([]models.Tag, error) {
	var tags []models.Tag
	
	err := s.db.Order("post_count DESC, name ASC").Limit(limit).Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("获取热门标签失败: %w", err)
	}
	
	return tags, nil
}
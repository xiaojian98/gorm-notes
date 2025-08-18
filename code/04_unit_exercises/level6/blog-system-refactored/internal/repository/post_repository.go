package repository

import (
	"errors"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// PostRepository 文章数据访问层接口
// 定义文章相关的数据库操作方法
type PostRepository interface {
	// 基本CRUD操作
	Create(post *models.Post) error                             // 创建文章
	GetByID(id uint) (*models.Post, error)                     // 根据ID获取文章
	GetBySlug(slug string) (*models.Post, error)               // 根据slug获取文章
	Update(post *models.Post) error                            // 更新文章
	Delete(id uint) error                                      // 删除文章（软删除）
	HardDelete(id uint) error                                  // 硬删除文章
	
	// 查询操作
	List(offset, limit int) ([]models.Post, error)            // 分页获取文章列表
	ListPublished(offset, limit int) ([]models.Post, error)   // 分页获取已发布文章
	ListByAuthor(authorID uint, offset, limit int) ([]models.Post, error) // 根据作者获取文章
	ListByCategory(categoryID uint, offset, limit int) ([]models.Post, error) // 根据分类获取文章
	ListByTag(tagID uint, offset, limit int) ([]models.Post, error) // 根据标签获取文章
	ListByStatus(status string, offset, limit int) ([]models.Post, error) // 根据状态获取文章
	
	// 搜索操作
	Search(keyword string, offset, limit int) ([]models.Post, error) // 搜索文章
	SearchPublished(keyword string, offset, limit int) ([]models.Post, error) // 搜索已发布文章
	AdvancedSearch(params SearchParams) ([]models.Post, error) // 高级搜索
	
	// 统计操作
	Count() (int64, error)                                     // 获取文章总数
	CountByStatus(status string) (int64, error)               // 根据状态统计文章数
	CountByAuthor(authorID uint) (int64, error)               // 根据作者统计文章数
	CountByCategory(categoryID uint) (int64, error)           // 根据分类统计文章数
	CountByTag(tagID uint) (int64, error)                     // 根据标签统计文章数
	
	// 热门和推荐
	GetPopularPosts(limit int, days int) ([]models.Post, error) // 获取热门文章
	GetRecentPosts(limit int) ([]models.Post, error)          // 获取最新文章
	GetRecommendedPosts(userID uint, limit int) ([]models.Post, error) // 获取推荐文章
	GetRelatedPosts(postID uint, limit int) ([]models.Post, error) // 获取相关文章
	
	// 浏览量操作
	IncrementViewCount(id uint) error                          // 增加浏览量
	GetMostViewedPosts(limit int, days int) ([]models.Post, error) // 获取最多浏览文章
	
	// 标签操作
	AddTags(postID uint, tagIDs []uint) error                  // 为文章添加标签
	RemoveTags(postID uint, tagIDs []uint) error               // 移除文章标签
	UpdateTags(postID uint, tagIDs []uint) error               // 更新文章标签
	GetPostTags(postID uint) ([]models.Tag, error)            // 获取文章标签
	
	// 批量操作
	BatchCreate(posts []models.Post) error                     // 批量创建文章
	BatchUpdateStatus(postIDs []uint, status string) error    // 批量更新状态
	BatchDelete(postIDs []uint) error                         // 批量删除文章
	
	// 高级查询
	GetPostWithDetails(id uint) (*PostWithDetails, error)     // 获取文章详细信息
	GetPostsWithStats(offset, limit int) ([]PostWithStats, error) // 获取文章及统计信息
	GetArchive(year, month int) ([]models.Post, error)        // 获取归档文章
	GetSitemap() ([]PostSitemap, error)                       // 获取站点地图数据
}

// postRepository 文章数据访问层实现
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建文章数据访问层实例
// 参数: db - 数据库连接
// 返回: PostRepository - 文章数据访问层接口实例
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

// 辅助数据结构

// SearchParams 搜索参数
type SearchParams struct {
	Keyword    string    `json:"keyword"`     // 搜索关键词
	AuthorID   uint      `json:"author_id"`   // 作者ID
	CategoryID uint      `json:"category_id"` // 分类ID
	TagIDs     []uint    `json:"tag_ids"`     // 标签ID列表
	Status     string    `json:"status"`      // 状态
	StartDate  time.Time `json:"start_date"`  // 开始日期
	EndDate    time.Time `json:"end_date"`    // 结束日期
	Offset     int       `json:"offset"`      // 偏移量
	Limit      int       `json:"limit"`       // 限制数量
	OrderBy    string    `json:"order_by"`    // 排序字段
	OrderDir   string    `json:"order_dir"`   // 排序方向
}

// PostWithDetails 文章详细信息
type PostWithDetails struct {
	models.Post
	Author   *models.User     `json:"author"`   // 作者信息
	Category *models.Category `json:"category"` // 分类信息
	Tags     []models.Tag     `json:"tags"`     // 标签列表
	CommentCount int          `json:"comment_count"` // 评论数
	LikeCount    int          `json:"like_count"`    // 点赞数
}

// PostWithStats 文章及统计信息
type PostWithStats struct {
	models.Post
	CommentCount int `json:"comment_count"` // 评论数
	LikeCount    int `json:"like_count"`    // 点赞数
	ShareCount   int `json:"share_count"`   // 分享数
}

// PostSitemap 站点地图文章信息
type PostSitemap struct {
	ID          uint      `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
}

// 基本CRUD操作实现

// Create 创建文章
// 参数: post - 文章对象
// 返回: error - 错误信息
func (r *postRepository) Create(post *models.Post) error {
	if post == nil {
		return errors.New("文章对象不能为空")
	}
	
	// 检查作者是否存在
	var count int64
	r.db.Model(&models.User{}).Where("id = ?", post.AuthorID).Count(&count)
	if count == 0 {
		return errors.New("作者不存在")
	}
	
	// 检查分类是否存在
	if post.CategoryID != nil && *post.CategoryID != 0 {
		r.db.Model(&models.Category{}).Where("id = ?", *post.CategoryID).Count(&count)
		if count == 0 {
			return errors.New("分类不存在")
		}
	}
	
	// 检查slug是否已存在
	if post.Slug != "" {
		r.db.Model(&models.Post{}).Where("slug = ?", post.Slug).Count(&count)
		if count > 0 {
			return errors.New("slug已存在")
		}
	}
	
	return r.db.Create(post).Error
}

// GetByID 根据ID获取文章
// 参数: id - 文章ID
// 返回: *models.Post - 文章对象, error - 错误信息
func (r *postRepository) GetByID(id uint) (*models.Post, error) {
	if id == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	post := &models.Post{}
	err := r.db.First(post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	return post, nil
}

// GetBySlug 根据slug获取文章
// 参数: slug - 文章slug
// 返回: *models.Post - 文章对象, error - 错误信息
func (r *postRepository) GetBySlug(slug string) (*models.Post, error) {
	if slug == "" {
		return nil, errors.New("slug不能为空")
	}
	
	post := &models.Post{}
	err := r.db.Where("slug = ?", slug).First(post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	return post, nil
}

// Update 更新文章
// 参数: post - 文章对象
// 返回: error - 错误信息
func (r *postRepository) Update(post *models.Post) error {
	if post == nil || post.ID == 0 {
		return errors.New("文章对象或ID不能为空")
	}
	
	// 检查文章是否存在
	var count int64
	r.db.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count)
	if count == 0 {
		return errors.New("文章不存在")
	}
	
	// 检查分类是否存在
	if post.CategoryID != nil && *post.CategoryID != 0 {
		r.db.Model(&models.Category{}).Where("id = ?", *post.CategoryID).Count(&count)
		if count == 0 {
			return errors.New("分类不存在")
		}
	}
	
	// 检查slug是否被其他文章使用
	if post.Slug != "" {
		r.db.Model(&models.Post{}).Where("slug = ? AND id != ?", post.Slug, post.ID).Count(&count)
		if count > 0 {
			return errors.New("slug已被其他文章使用")
		}
	}
	
	return r.db.Save(post).Error
}

// Delete 删除文章（软删除）
// 参数: id - 文章ID
// 返回: error - 错误信息
func (r *postRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	// 检查文章是否存在
	var count int64
	r.db.Model(&models.Post{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return errors.New("文章不存在")
	}
	
	return r.db.Delete(&models.Post{}, id).Error
}

// HardDelete 硬删除文章
// 参数: id - 文章ID
// 返回: error - 错误信息
func (r *postRepository) HardDelete(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	return r.db.Unscoped().Delete(&models.Post{}, id).Error
}

// 查询操作实现

// List 分页获取文章列表
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) List(offset, limit int) ([]models.Post, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// ListPublished 分页获取已发布文章
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) ListPublished(offset, limit int) ([]models.Post, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Where("status = ? AND published_at <= ?", "published", time.Now()).
		Offset(offset).Limit(limit).Order("published_at DESC").Find(&posts).Error
	return posts, err
}

// ListByAuthor 根据作者获取文章
// 参数: authorID - 作者ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) ListByAuthor(authorID uint, offset, limit int) ([]models.Post, error) {
	if authorID == 0 {
		return nil, errors.New("作者ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Where("author_id = ?", authorID).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// ListByCategory 根据分类获取文章
// 参数: categoryID - 分类ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) ListByCategory(categoryID uint, offset, limit int) ([]models.Post, error) {
	if categoryID == 0 {
		return nil, errors.New("分类ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Where("category_id = ? AND status = ?", categoryID, "published").
		Offset(offset).Limit(limit).Order("published_at DESC").Find(&posts).Error
	return posts, err
}

// ListByTag 根据标签获取文章
// 参数: tagID - 标签ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) ListByTag(tagID uint, offset, limit int) ([]models.Post, error) {
	if tagID == 0 {
		return nil, errors.New("标签ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = ?", tagID, "published").
		Offset(offset).Limit(limit).Order("posts.published_at DESC").Find(&posts).Error
	return posts, err
}

// ListByStatus 根据状态获取文章
// 参数: status - 文章状态, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) ListByStatus(status string, offset, limit int) ([]models.Post, error) {
	if status == "" {
		return nil, errors.New("状态不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Where("status = ?", status).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// 搜索操作实现

// Search 搜索文章
// 参数: keyword - 搜索关键词, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) Search(keyword string, offset, limit int) ([]models.Post, error) {
	if keyword == "" {
		return r.List(offset, limit)
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	keyword = "%" + keyword + "%"
	err := r.db.Preload("Author").Preload("Category").
		Where("title LIKE ? OR content LIKE ? OR excerpt LIKE ?", keyword, keyword, keyword).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// SearchPublished 搜索已发布文章
// 参数: keyword - 搜索关键词, offset - 偏移量, limit - 限制数量
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) SearchPublished(keyword string, offset, limit int) ([]models.Post, error) {
	if keyword == "" {
		return r.ListPublished(offset, limit)
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var posts []models.Post
	keyword = "%" + keyword + "%"
	err := r.db.Preload("Author").Preload("Category").
		Where("(title LIKE ? OR content LIKE ? OR excerpt LIKE ?) AND status = ? AND published_at <= ?", 
			keyword, keyword, keyword, "published", time.Now()).
		Offset(offset).Limit(limit).Order("published_at DESC").Find(&posts).Error
	return posts, err
}

// AdvancedSearch 高级搜索
// 参数: params - 搜索参数
// 返回: []models.Post - 文章列表, error - 错误信息
func (r *postRepository) AdvancedSearch(params SearchParams) ([]models.Post, error) {
	if params.Offset < 0 {
		params.Offset = 0
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.OrderBy == "" {
		params.OrderBy = "created_at"
	}
	if params.OrderDir == "" {
		params.OrderDir = "DESC"
	}
	
	query := r.db.Preload("Author").Preload("Category")
	
	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR excerpt LIKE ?", keyword, keyword, keyword)
	}
	
	// 作者筛选
	if params.AuthorID != 0 {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	
	// 分类筛选
	if params.CategoryID != 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	
	// 状态筛选
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	
	// 日期范围筛选
	if !params.StartDate.IsZero() {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if !params.EndDate.IsZero() {
		query = query.Where("created_at <= ?", params.EndDate)
	}
	
	// 标签筛选
	if len(params.TagIDs) > 0 {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id IN (?)", params.TagIDs)
	}
	
	var posts []models.Post
	err := query.Offset(params.Offset).Limit(params.Limit).
		Order(params.OrderBy + " " + params.OrderDir).Find(&posts).Error
	return posts, err
}

// 统计操作实现

// Count 获取文章总数
// 返回: int64 - 文章总数, error - 错误信息
func (r *postRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计文章数
// 参数: status - 文章状态
// 返回: int64 - 文章数量, error - 错误信息
func (r *postRepository) CountByStatus(status string) (int64, error) {
	if status == "" {
		return 0, errors.New("状态不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Post{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountByAuthor 根据作者统计文章数
// 参数: authorID - 作者ID
// 返回: int64 - 文章数量, error - 错误信息
func (r *postRepository) CountByAuthor(authorID uint) (int64, error) {
	if authorID == 0 {
		return 0, errors.New("作者ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Post{}).Where("author_id = ?", authorID).Count(&count).Error
	return count, err
}

// CountByCategory 根据分类统计文章数
// 参数: categoryID - 分类ID
// 返回: int64 - 文章数量, error - 错误信息
func (r *postRepository) CountByCategory(categoryID uint) (int64, error) {
	if categoryID == 0 {
		return 0, errors.New("分类ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Post{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// CountByTag 根据标签统计文章数
// 参数: tagID - 标签ID
// 返回: int64 - 文章数量, error - 错误信息
func (r *postRepository) CountByTag(tagID uint) (int64, error) {
	if tagID == 0 {
		return 0, errors.New("标签ID不能为空")
	}
	
	var count int64
	err := r.db.Table("posts").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ?", tagID).Count(&count).Error
	return count, err
}

// 热门和推荐实现

// GetPopularPosts 获取热门文章
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.Post - 热门文章列表, error - 错误信息
func (r *postRepository) GetPopularPosts(limit int, days int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if days <= 0 {
		days = 7
	}
	
	var posts []models.Post
	startDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Preload("Author").Preload("Category").
		Where("status = ? AND published_at >= ? AND published_at <= ?", "published", startDate, time.Now()).
		Order("view_count DESC, (SELECT COUNT(*) FROM comments WHERE post_id = posts.id) DESC").
		Limit(limit).Find(&posts).Error
	
	return posts, err
}

// GetRecentPosts 获取最新文章
// 参数: limit - 限制数量
// 返回: []models.Post - 最新文章列表, error - 错误信息
func (r *postRepository) GetRecentPosts(limit int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	var posts []models.Post
	err := r.db.Preload("Author").Preload("Category").
		Where("status = ? AND published_at <= ?", "published", time.Now()).
		Order("published_at DESC").Limit(limit).Find(&posts).Error
	
	return posts, err
}

// GetRecommendedPosts 获取推荐文章
// 参数: userID - 用户ID, limit - 限制数量
// 返回: []models.Post - 推荐文章列表, error - 错误信息
func (r *postRepository) GetRecommendedPosts(userID uint, limit int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	// 简化的推荐算法：基于用户关注的作者和热门文章
	var posts []models.Post
	
	if userID != 0 {
		// 获取用户关注的作者的文章
		err := r.db.Preload("Author").Preload("Category").
			Joins("JOIN follows ON posts.author_id = follows.following_id").
			Where("follows.follower_id = ? AND posts.status = ? AND posts.published_at <= ?", 
				userID, "published", time.Now()).
			Order("posts.published_at DESC").Limit(limit).Find(&posts).Error
		
		if err == nil && len(posts) > 0 {
			return posts, nil
		}
	}
	
	// 如果没有关注的作者或获取失败，返回热门文章
	return r.GetPopularPosts(limit, 30)
}

// GetRelatedPosts 获取相关文章
// 参数: postID - 文章ID, limit - 限制数量
// 返回: []models.Post - 相关文章列表, error - 错误信息
func (r *postRepository) GetRelatedPosts(postID uint, limit int) ([]models.Post, error) {
	if postID == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	
	// 获取当前文章信息
	post, err := r.GetByID(postID)
	if err != nil {
		return nil, err
	}
	
	var relatedPosts []models.Post
	
	// 首先尝试获取同分类的文章
	if post.CategoryID != nil && *post.CategoryID != 0 {
		err = r.db.Preload("Author").Preload("Category").
			Where("category_id = ? AND id != ? AND status = ? AND published_at <= ?", 
				*post.CategoryID, postID, "published", time.Now()).
			Order("published_at DESC").Limit(limit).Find(&relatedPosts).Error
		
		if err == nil && len(relatedPosts) >= limit {
			return relatedPosts, nil
		}
	}
	
	// 如果同分类文章不够，获取同作者的其他文章
	remaining := limit - len(relatedPosts)
	if remaining > 0 {
		var authorPosts []models.Post
		excludeIDs := []uint{postID}
		for _, p := range relatedPosts {
			excludeIDs = append(excludeIDs, p.ID)
		}
		
		r.db.Preload("Author").Preload("Category").
			Where("author_id = ? AND id NOT IN (?) AND status = ? AND published_at <= ?", 
				post.AuthorID, excludeIDs, "published", time.Now()).
			Order("published_at DESC").Limit(remaining).Find(&authorPosts)
		
		relatedPosts = append(relatedPosts, authorPosts...)
	}
	
	return relatedPosts, nil
}

// 浏览量操作实现

// IncrementViewCount 增加浏览量
// 参数: id - 文章ID
// 返回: error - 错误信息
func (r *postRepository) IncrementViewCount(id uint) error {
	if id == 0 {
		return errors.New("文章ID不能为空")
	}
	
	return r.db.Model(&models.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// GetMostViewedPosts 获取最多浏览文章
// 参数: limit - 限制数量, days - 统计天数
// 返回: []models.Post - 最多浏览文章列表, error - 错误信息
func (r *postRepository) GetMostViewedPosts(limit int, days int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	query := r.db.Preload("Author").Preload("Category").
		Where("status = ?", "published")
	
	if days > 0 {
		startDate := time.Now().AddDate(0, 0, -days)
		query = query.Where("published_at >= ?", startDate)
	}
	
	var posts []models.Post
	err := query.Order("view_count DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

// 标签操作实现

// AddTags 为文章添加标签
// 参数: postID - 文章ID, tagIDs - 标签ID列表
// 返回: error - 错误信息
func (r *postRepository) AddTags(postID uint, tagIDs []uint) error {
	if postID == 0 {
		return errors.New("文章ID不能为空")
	}
	if len(tagIDs) == 0 {
		return errors.New("标签ID列表不能为空")
	}
	
	// 检查文章是否存在
	var count int64
	r.db.Model(&models.Post{}).Where("id = ?", postID).Count(&count)
	if count == 0 {
		return errors.New("文章不存在")
	}
	
	// 检查标签是否存在
	r.db.Model(&models.Tag{}).Where("id IN (?)", tagIDs).Count(&count)
	if int(count) != len(tagIDs) {
		return errors.New("部分标签不存在")
	}
	
	// 添加标签关联（忽略已存在的关联）
	for _, tagID := range tagIDs {
		r.db.Exec("INSERT IGNORE INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tagID)
	}
	
	return nil
}

// RemoveTags 移除文章标签
// 参数: postID - 文章ID, tagIDs - 标签ID列表
// 返回: error - 错误信息
func (r *postRepository) RemoveTags(postID uint, tagIDs []uint) error {
	if postID == 0 {
		return errors.New("文章ID不能为空")
	}
	if len(tagIDs) == 0 {
		return errors.New("标签ID列表不能为空")
	}
	
	return r.db.Exec("DELETE FROM post_tags WHERE post_id = ? AND tag_id IN (?)", postID, tagIDs).Error
}

// UpdateTags 更新文章标签
// 参数: postID - 文章ID, tagIDs - 标签ID列表
// 返回: error - 错误信息
func (r *postRepository) UpdateTags(postID uint, tagIDs []uint) error {
	if postID == 0 {
		return errors.New("文章ID不能为空")
	}
	
	// 开启事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 删除现有标签关联
	if err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", postID).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 添加新的标签关联
	if len(tagIDs) > 0 {
		// 检查标签是否存在
		var count int64
		tx.Model(&models.Tag{}).Where("id IN (?)", tagIDs).Count(&count)
		if int(count) != len(tagIDs) {
			tx.Rollback()
			return errors.New("部分标签不存在")
		}
		
		for _, tagID := range tagIDs {
			if err := tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tagID).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	
	return tx.Commit().Error
}

// GetPostTags 获取文章标签
// 参数: postID - 文章ID
// 返回: []models.Tag - 标签列表, error - 错误信息
func (r *postRepository) GetPostTags(postID uint) ([]models.Tag, error) {
	if postID == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	var tags []models.Tag
	err := r.db.Table("tags").
		Joins("JOIN post_tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error
	
	return tags, err
}

// 批量操作实现

// BatchCreate 批量创建文章
// 参数: posts - 文章列表
// 返回: error - 错误信息
func (r *postRepository) BatchCreate(posts []models.Post) error {
	if len(posts) == 0 {
		return errors.New("文章列表不能为空")
	}
	
	// 检查slug唯一性
	slugs := make([]string, 0)
	for _, post := range posts {
		if post.Slug != "" {
			slugs = append(slugs, post.Slug)
		}
	}
	
	if len(slugs) > 0 {
		var count int64
		r.db.Model(&models.Post{}).Where("slug IN (?)", slugs).Count(&count)
		if count > 0 {
			return errors.New("存在重复的slug")
		}
	}
	
	return r.db.CreateInBatches(posts, 100).Error
}

// BatchUpdateStatus 批量更新状态
// 参数: postIDs - 文章ID列表, status - 新状态
// 返回: error - 错误信息
func (r *postRepository) BatchUpdateStatus(postIDs []uint, status string) error {
	if len(postIDs) == 0 {
		return errors.New("文章ID列表不能为空")
	}
	if status == "" {
		return errors.New("状态不能为空")
	}
	
	updateData := map[string]interface{}{"status": status}
	if status == "published" {
		updateData["published_at"] = time.Now()
	}
	
	return r.db.Model(&models.Post{}).Where("id IN (?)", postIDs).Updates(updateData).Error
}

// BatchDelete 批量删除文章
// 参数: postIDs - 文章ID列表
// 返回: error - 错误信息
func (r *postRepository) BatchDelete(postIDs []uint) error {
	if len(postIDs) == 0 {
		return errors.New("文章ID列表不能为空")
	}
	
	return r.db.Delete(&models.Post{}, postIDs).Error
}

// 高级查询实现

// GetPostWithDetails 获取文章详细信息
// 参数: id - 文章ID
// 返回: *PostWithDetails - 文章详细信息, error - 错误信息
func (r *postRepository) GetPostWithDetails(id uint) (*PostWithDetails, error) {
	if id == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	// 获取文章基本信息
	post := &models.Post{}
	err := r.db.Preload("Author").Preload("Category").First(post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	
	postDetails := &PostWithDetails{
		Post:     *post,
		Author:   post.Author,
		Category: post.Category,
	}
	
	// 获取标签
	tags, _ := r.GetPostTags(id)
	postDetails.Tags = tags
	
	// 统计评论数
	var commentCount int64
	r.db.Model(&models.Comment{}).Where("post_id = ?", id).Count(&commentCount)
	postDetails.CommentCount = int(commentCount)
	
	// 统计点赞数
	var likeCount int64
	r.db.Model(&models.Like{}).Where("target_id = ? AND target_type = ?", id, "post").Count(&likeCount)
	postDetails.LikeCount = int(likeCount)
	
	return postDetails, nil
}

// GetPostsWithStats 获取文章及统计信息
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []PostWithStats - 文章及统计信息列表, error - 错误信息
func (r *postRepository) GetPostsWithStats(offset, limit int) ([]PostWithStats, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var postsWithStats []PostWithStats
	
	err := r.db.Table("posts").
		Select(`
			posts.*,
			COALESCE(comment_counts.comment_count, 0) as comment_count,
			COALESCE(like_counts.like_count, 0) as like_count,
			0 as share_count
		`).
		Joins(`LEFT JOIN (
			SELECT post_id, COUNT(*) as comment_count 
			FROM comments 
			GROUP BY post_id
		) comment_counts ON posts.id = comment_counts.post_id`).
		Joins(`LEFT JOIN (
			SELECT target_id, COUNT(*) as like_count 
			FROM likes 
			WHERE target_type = 'post'
			GROUP BY target_id
		) like_counts ON posts.id = like_counts.target_id`).
		Offset(offset).Limit(limit).Order("posts.created_at DESC").
		Scan(&postsWithStats).Error
	
	return postsWithStats, err
}

// GetArchive 获取归档文章
// 参数: year - 年份, month - 月份（0表示整年）
// 返回: []models.Post - 归档文章列表, error - 错误信息
func (r *postRepository) GetArchive(year, month int) ([]models.Post, error) {
	if year <= 0 {
		return nil, errors.New("年份不能为空")
	}
	
	query := r.db.Preload("Author").Preload("Category").
		Where("status = ? AND YEAR(published_at) = ?", "published", year)
	
	if month > 0 && month <= 12 {
		query = query.Where("MONTH(published_at) = ?", month)
	}
	
	var posts []models.Post
	err := query.Order("published_at DESC").Find(&posts).Error
	return posts, err
}

// GetSitemap 获取站点地图数据
// 返回: []PostSitemap - 站点地图文章列表, error - 错误信息
func (r *postRepository) GetSitemap() ([]PostSitemap, error) {
	var sitemap []PostSitemap
	
	err := r.db.Table("posts").
		Select("id, slug, title, updated_at, published_at").
		Where("status = ? AND published_at <= ?", "published", time.Now()).
		Order("published_at DESC").
		Scan(&sitemap).Error
	
	return sitemap, err
}
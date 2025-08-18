package repository

import (
	"database/sql"
	"errors"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层接口
// 定义评论相关的数据库操作方法
type CommentRepository interface {
	// 基本CRUD操作
	Create(comment *models.Comment) error                      // 创建评论
	GetByID(id uint) (*models.Comment, error)                 // 根据ID获取评论
	Update(comment *models.Comment) error                     // 更新评论
	Delete(id uint) error                                     // 删除评论（软删除）
	HardDelete(id uint) error                                 // 硬删除评论
	
	// 查询操作
	List(offset, limit int) ([]models.Comment, error)        // 分页获取评论列表
	ListByPost(postID uint, offset, limit int) ([]models.Comment, error) // 根据文章获取评论
	ListByUser(userID uint, offset, limit int) ([]models.Comment, error) // 根据用户获取评论
	ListByStatus(status string, offset, limit int) ([]models.Comment, error) // 根据状态获取评论
	ListReplies(parentID uint, offset, limit int) ([]models.Comment, error) // 获取回复评论
	ListPending(offset, limit int) ([]models.Comment, error) // 获取待审核评论
	
	// 树形结构操作
	GetCommentTree(postID uint) ([]CommentNode, error)       // 获取评论树
	GetCommentWithReplies(id uint) (*CommentWithReplies, error) // 获取评论及其回复
	
	// 统计操作
	Count() (int64, error)                                   // 获取评论总数
	CountByPost(postID uint) (int64, error)                  // 根据文章统计评论数
	CountByUser(userID uint) (int64, error)                  // 根据用户统计评论数
	CountByStatus(status string) (int64, error)              // 根据状态统计评论数
	CountReplies(parentID uint) (int64, error)               // 统计回复数
	
	// 状态操作
	Approve(id uint) error                                   // 批准评论
	Reject(id uint) error                                    // 拒绝评论
	MarkAsSpam(id uint) error                                // 标记为垃圾评论
	BatchUpdateStatus(commentIDs []uint, status string) error // 批量更新状态
	
	// 点赞操作
	AddLike(commentID, userID uint) error                    // 添加点赞
	RemoveLike(commentID, userID uint) error                 // 移除点赞
	GetLikeCount(commentID uint) (int64, error)              // 获取点赞数
	IsLikedByUser(commentID, userID uint) (bool, error)      // 检查用户是否点赞
	
	// 举报操作
	ReportComment(commentID, userID uint, reason string) error // 举报评论
	GetReports(commentID uint) ([]CommentReport, error)      // 获取举报记录
	
	// 批量操作
	BatchCreate(comments []models.Comment) error             // 批量创建评论
	BatchDelete(commentIDs []uint) error                     // 批量删除评论
	
	// 高级查询
	GetRecentComments(limit int) ([]models.Comment, error)   // 获取最新评论
	GetPopularComments(postID uint, limit int) ([]models.Comment, error) // 获取热门评论
	GetUserCommentStats(userID uint) (*UserCommentStats, error) // 获取用户评论统计
	SearchComments(keyword string, offset, limit int) ([]models.Comment, error) // 搜索评论
}

// commentRepository 评论数据访问层实现
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论数据访问层实例
// 参数: db - 数据库连接
// 返回: CommentRepository - 评论数据访问层接口实例
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

// 辅助数据结构

// CommentNode 评论树节点
type CommentNode struct {
	models.Comment
	Replies []CommentNode `json:"replies"` // 子评论
}

// CommentWithReplies 评论及其回复
type CommentWithReplies struct {
	models.Comment
	Replies []models.Comment `json:"replies"` // 回复列表
	ReplyCount int           `json:"reply_count"` // 回复数量
}

// CommentReport 评论举报
type CommentReport struct {
	ID        uint      `json:"id"`
	CommentID uint      `json:"comment_id"`
	UserID    uint      `json:"user_id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// UserCommentStats 用户评论统计
type UserCommentStats struct {
	TotalComments    int64 `json:"total_comments"`    // 总评论数
	ApprovedComments int64 `json:"approved_comments"` // 已批准评论数
	PendingComments  int64 `json:"pending_comments"`  // 待审核评论数
	RejectedComments int64 `json:"rejected_comments"` // 被拒绝评论数
	SpamComments     int64 `json:"spam_comments"`     // 垃圾评论数
	TotalLikes       int64 `json:"total_likes"`       // 总获赞数
	AverageLength    float64 `json:"average_length"`   // 平均评论长度
}

// 基本CRUD操作实现

// Create 创建评论
// 参数: comment - 评论对象
// 返回: error - 错误信息
func (r *commentRepository) Create(comment *models.Comment) error {
	if comment == nil {
		return errors.New("评论对象不能为空")
	}
	
	// 检查文章是否存在
	var count int64
	r.db.Model(&models.Post{}).Where("id = ?", comment.PostID).Count(&count)
	if count == 0 {
		return errors.New("文章不存在")
	}
	
	// 检查用户是否存在
	r.db.Model(&models.User{}).Where("id = ?", comment.UserID).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	
	// 如果是回复评论，检查父评论是否存在
	if comment.ParentID != nil && *comment.ParentID != 0 {
		r.db.Model(&models.Comment{}).Where("id = ? AND post_id = ?", *comment.ParentID, comment.PostID).Count(&count)
		if count == 0 {
			return errors.New("父评论不存在或不属于同一文章")
		}
	}
	
	return r.db.Create(comment).Error
}

// GetByID 根据ID获取评论
// 参数: id - 评论ID
// 返回: *models.Comment - 评论对象, error - 错误信息
func (r *commentRepository) GetByID(id uint) (*models.Comment, error) {
	if id == 0 {
		return nil, errors.New("评论ID不能为空")
	}
	
	comment := &models.Comment{}
	err := r.db.Preload("User").Preload("Post").First(comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}
	
	return comment, nil
}

// Update 更新评论
// 参数: comment - 评论对象
// 返回: error - 错误信息
func (r *commentRepository) Update(comment *models.Comment) error {
	if comment == nil || comment.ID == 0 {
		return errors.New("评论对象或ID不能为空")
	}
	
	// 检查评论是否存在
	var count int64
	r.db.Model(&models.Comment{}).Where("id = ?", comment.ID).Count(&count)
	if count == 0 {
		return errors.New("评论不存在")
	}
	
	return r.db.Save(comment).Error
}

// Delete 删除评论（软删除）
// 参数: id - 评论ID
// 返回: error - 错误信息
func (r *commentRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	// 检查评论是否存在
	var count int64
	r.db.Model(&models.Comment{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return errors.New("评论不存在")
	}
	
	return r.db.Delete(&models.Comment{}, id).Error
}

// HardDelete 硬删除评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (r *commentRepository) HardDelete(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	return r.db.Unscoped().Delete(&models.Comment{}, id).Error
}

// 查询操作实现

// List 分页获取评论列表
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, error - 错误信息
func (r *commentRepository) List(offset, limit int) ([]models.Comment, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	err := r.db.Preload("User").Preload("Post").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

// ListByPost 根据文章获取评论
// 参数: postID - 文章ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, error - 错误信息
func (r *commentRepository) ListByPost(postID uint, offset, limit int) ([]models.Comment, error) {
	if postID == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	err := r.db.Preload("User").
		Where("post_id = ? AND status = ?", postID, "approved").
		Offset(offset).Limit(limit).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// ListByUser 根据用户获取评论
// 参数: userID - 用户ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, error - 错误信息
func (r *commentRepository) ListByUser(userID uint, offset, limit int) ([]models.Comment, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	err := r.db.Preload("Post").
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

// ListByStatus 根据状态获取评论
// 参数: status - 评论状态, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, error - 错误信息
func (r *commentRepository) ListByStatus(status string, offset, limit int) ([]models.Comment, error) {
	if status == "" {
		return nil, errors.New("状态不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	err := r.db.Preload("User").Preload("Post").
		Where("status = ?", status).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

// ListReplies 获取回复评论
// 参数: parentID - 父评论ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 回复评论列表, error - 错误信息
func (r *commentRepository) ListReplies(parentID uint, offset, limit int) ([]models.Comment, error) {
	if parentID == 0 {
		return nil, errors.New("父评论ID不能为空")
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	err := r.db.Preload("User").
		Where("parent_id = ? AND status = ?", parentID, "approved").
		Offset(offset).Limit(limit).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// ListPending 获取待审核评论
// 参数: offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 待审核评论列表, error - 错误信息
func (r *commentRepository) ListPending(offset, limit int) ([]models.Comment, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	err := r.db.Preload("User").Preload("Post").
		Where("status = ?", "pending").
		Offset(offset).Limit(limit).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// 树形结构操作实现

// GetCommentTree 获取评论树
// 参数: postID - 文章ID
// 返回: []CommentNode - 评论树, error - 错误信息
func (r *commentRepository) GetCommentTree(postID uint) ([]CommentNode, error) {
	if postID == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	// 获取所有评论
	var comments []models.Comment
	err := r.db.Preload("User").
		Where("post_id = ? AND status = ?", postID, "approved").
		Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	
	// 构建评论映射
	commentMap := make(map[uint]*CommentNode)
	var rootComments []CommentNode
	
	// 第一遍：创建所有节点
	for _, comment := range comments {
		node := CommentNode{
			Comment: comment,
			Replies: make([]CommentNode, 0),
		}
		commentMap[comment.ID] = &node
	}
	
	// 第二遍：构建树形结构
	for _, comment := range comments {
		node := commentMap[comment.ID]
		if comment.ParentID == nil || *comment.ParentID == 0 {
			// 根评论
			rootComments = append(rootComments, *node)
		} else {
			// 子评论
			if parentNode, exists := commentMap[*comment.ParentID]; exists {
				parentNode.Replies = append(parentNode.Replies, *node)
			}
		}
	}
	
	return rootComments, nil
}

// GetCommentWithReplies 获取评论及其回复
// 参数: id - 评论ID
// 返回: *CommentWithReplies - 评论及其回复, error - 错误信息
func (r *commentRepository) GetCommentWithReplies(id uint) (*CommentWithReplies, error) {
	if id == 0 {
		return nil, errors.New("评论ID不能为空")
	}
	
	// 获取主评论
	comment, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 获取回复
	replies, err := r.ListReplies(id, 0, 100)
	if err != nil {
		return nil, err
	}
	
	result := &CommentWithReplies{
		Comment:    *comment,
		Replies:    replies,
		ReplyCount: len(replies),
	}
	
	return result, nil
}

// 统计操作实现

// Count 获取评论总数
// 返回: int64 - 评论总数, error - 错误信息
func (r *commentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Count(&count).Error
	return count, err
}

// CountByPost 根据文章统计评论数
// 参数: postID - 文章ID
// 返回: int64 - 评论数量, error - 错误信息
func (r *commentRepository) CountByPost(postID uint) (int64, error) {
	if postID == 0 {
		return 0, errors.New("文章ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, "approved").Count(&count).Error
	return count, err
}

// CountByUser 根据用户统计评论数
// 参数: userID - 用户ID
// 返回: int64 - 评论数量, error - 错误信息
func (r *commentRepository) CountByUser(userID uint) (int64, error) {
	if userID == 0 {
		return 0, errors.New("用户ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计评论数
// 参数: status - 评论状态
// 返回: int64 - 评论数量, error - 错误信息
func (r *commentRepository) CountByStatus(status string) (int64, error) {
	if status == "" {
		return 0, errors.New("状态不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Comment{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountReplies 统计回复数
// 参数: parentID - 父评论ID
// 返回: int64 - 回复数量, error - 错误信息
func (r *commentRepository) CountReplies(parentID uint) (int64, error) {
	if parentID == 0 {
		return 0, errors.New("父评论ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Comment{}).Where("parent_id = ? AND status = ?", parentID, "approved").Count(&count).Error
	return count, err
}

// 状态操作实现

// Approve 批准评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (r *commentRepository) Approve(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "approved").Error
}

// Reject 拒绝评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (r *commentRepository) Reject(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "rejected").Error
}

// MarkAsSpam 标记为垃圾评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (r *commentRepository) MarkAsSpam(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "spam").Error
}

// BatchUpdateStatus 批量更新状态
// 参数: commentIDs - 评论ID列表, status - 新状态
// 返回: error - 错误信息
func (r *commentRepository) BatchUpdateStatus(commentIDs []uint, status string) error {
	if len(commentIDs) == 0 {
		return errors.New("评论ID列表不能为空")
	}
	if status == "" {
		return errors.New("状态不能为空")
	}
	
	return r.db.Model(&models.Comment{}).Where("id IN (?)", commentIDs).Update("status", status).Error
}

// 点赞操作实现

// AddLike 添加点赞
// 参数: commentID - 评论ID, userID - 用户ID
// 返回: error - 错误信息
func (r *commentRepository) AddLike(commentID, userID uint) error {
	if commentID == 0 {
		return errors.New("评论ID不能为空")
	}
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	// 检查是否已经点赞
	var count int64
	r.db.Model(&models.Like{}).Where("target_id = ? AND target_type = ? AND user_id = ?", 
		commentID, "comment", userID).Count(&count)
	if count > 0 {
		return errors.New("已经点赞过了")
	}
	
	like := &models.Like{
		TargetID:   commentID,
		TargetType: "comment",
		UserID:     userID,
	}
	
	return r.db.Create(like).Error
}

// RemoveLike 移除点赞
// 参数: commentID - 评论ID, userID - 用户ID
// 返回: error - 错误信息
func (r *commentRepository) RemoveLike(commentID, userID uint) error {
	if commentID == 0 {
		return errors.New("评论ID不能为空")
	}
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	return r.db.Where("target_id = ? AND target_type = ? AND user_id = ?", 
		commentID, "comment", userID).Delete(&models.Like{}).Error
}

// GetLikeCount 获取点赞数
// 参数: commentID - 评论ID
// 返回: int64 - 点赞数, error - 错误信息
func (r *commentRepository) GetLikeCount(commentID uint) (int64, error) {
	if commentID == 0 {
		return 0, errors.New("评论ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Like{}).Where("target_id = ? AND target_type = ?", 
		commentID, "comment").Count(&count).Error
	return count, err
}

// IsLikedByUser 检查用户是否点赞
// 参数: commentID - 评论ID, userID - 用户ID
// 返回: bool - 是否点赞, error - 错误信息
func (r *commentRepository) IsLikedByUser(commentID, userID uint) (bool, error) {
	if commentID == 0 || userID == 0 {
		return false, errors.New("评论ID和用户ID不能为空")
	}
	
	var count int64
	err := r.db.Model(&models.Like{}).Where("target_id = ? AND target_type = ? AND user_id = ?", 
		commentID, "comment", userID).Count(&count).Error
	return count > 0, err
}

// 举报操作实现

// ReportComment 举报评论
// 参数: commentID - 评论ID, userID - 用户ID, reason - 举报原因
// 返回: error - 错误信息
func (r *commentRepository) ReportComment(commentID, userID uint, reason string) error {
	if commentID == 0 {
		return errors.New("评论ID不能为空")
	}
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	if reason == "" {
		return errors.New("举报原因不能为空")
	}
	
	// 检查是否已经举报过
	var count int64
	r.db.Table("comment_reports").Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count)
	if count > 0 {
		return errors.New("已经举报过了")
	}
	
	// 创建举报记录
	return r.db.Exec("INSERT INTO comment_reports (comment_id, user_id, reason, created_at) VALUES (?, ?, ?, ?)", 
		commentID, userID, reason, time.Now()).Error
}

// GetReports 获取举报记录
// 参数: commentID - 评论ID
// 返回: []CommentReport - 举报记录列表, error - 错误信息
func (r *commentRepository) GetReports(commentID uint) ([]CommentReport, error) {
	if commentID == 0 {
		return nil, errors.New("评论ID不能为空")
	}
	
	var reports []CommentReport
	err := r.db.Table("comment_reports").Where("comment_id = ?", commentID).
		Order("created_at DESC").Find(&reports).Error
	return reports, err
}

// 批量操作实现

// BatchCreate 批量创建评论
// 参数: comments - 评论列表
// 返回: error - 错误信息
func (r *commentRepository) BatchCreate(comments []models.Comment) error {
	if len(comments) == 0 {
		return errors.New("评论列表不能为空")
	}
	
	return r.db.CreateInBatches(comments, 100).Error
}

// BatchDelete 批量删除评论
// 参数: commentIDs - 评论ID列表
// 返回: error - 错误信息
func (r *commentRepository) BatchDelete(commentIDs []uint) error {
	if len(commentIDs) == 0 {
		return errors.New("评论ID列表不能为空")
	}
	
	return r.db.Delete(&models.Comment{}, commentIDs).Error
}

// 高级查询实现

// GetRecentComments 获取最新评论
// 参数: limit - 限制数量
// 返回: []models.Comment - 最新评论列表, error - 错误信息
func (r *commentRepository) GetRecentComments(limit int) ([]models.Comment, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	var comments []models.Comment
	err := r.db.Preload("User").Preload("Post").
		Where("status = ?", "approved").
		Order("created_at DESC").Limit(limit).Find(&comments).Error
	return comments, err
}

// GetPopularComments 获取热门评论
// 参数: postID - 文章ID, limit - 限制数量
// 返回: []models.Comment - 热门评论列表, error - 错误信息
func (r *commentRepository) GetPopularComments(postID uint, limit int) ([]models.Comment, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	query := r.db.Preload("User").Where("status = ?", "approved")
	
	if postID != 0 {
		query = query.Where("post_id = ?", postID)
	}
	
	var comments []models.Comment
	err := query.Select(`
		comments.*,
		(SELECT COUNT(*) FROM likes WHERE target_id = comments.id AND target_type = 'comment') as like_count
	`).Order("like_count DESC, created_at DESC").Limit(limit).Find(&comments).Error
	
	return comments, err
}

// GetUserCommentStats 获取用户评论统计
// 参数: userID - 用户ID
// 返回: *UserCommentStats - 用户评论统计, error - 错误信息
func (r *commentRepository) GetUserCommentStats(userID uint) (*UserCommentStats, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	stats := &UserCommentStats{}
	
	// 总评论数
	r.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&stats.TotalComments)
	
	// 各状态评论数
	r.db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", userID, "approved").Count(&stats.ApprovedComments)
	r.db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&stats.PendingComments)
	r.db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", userID, "rejected").Count(&stats.RejectedComments)
	r.db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", userID, "spam").Count(&stats.SpamComments)
	
	// 总获赞数
	r.db.Table("likes").
		Joins("JOIN comments ON likes.target_id = comments.id").
		Where("likes.target_type = ? AND comments.user_id = ?", "comment", userID).
		Count(&stats.TotalLikes)
	
	// 平均评论长度
	var avgLength sql.NullFloat64
	r.db.Model(&models.Comment{}).Where("user_id = ?", userID).
		Select("AVG(LENGTH(content))").Scan(&avgLength)
	if avgLength.Valid {
		stats.AverageLength = avgLength.Float64
	}
	
	return stats, nil
}

// SearchComments 搜索评论
// 参数: keyword - 搜索关键词, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, error - 错误信息
func (r *commentRepository) SearchComments(keyword string, offset, limit int) ([]models.Comment, error) {
	if keyword == "" {
		return r.List(offset, limit)
	}
	
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	keyword = "%" + keyword + "%"
	err := r.db.Preload("User").Preload("Post").
		Where("content LIKE ? AND status = ?", keyword, "approved").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, err
}
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"blog-system-refactored/internal/models"
	"gorm.io/gorm"
)

// CommentService 评论服务接口
// 定义评论相关的业务操作
type CommentService interface {
	// 评论基本操作
	CreateComment(comment *models.Comment) error                    // 创建评论
	GetCommentByID(id uint) (*models.Comment, error)               // 根据ID获取评论
	UpdateComment(comment *models.Comment) error                   // 更新评论
	DeleteComment(id uint) error                                   // 删除评论
	ListComments(offset, limit int, filters CommentFilters) ([]models.Comment, int64, error) // 分页获取评论列表
	
	// 评论状态操作
	ApproveComment(id uint) error                                  // 审核通过评论
	RejectComment(id uint) error                                   // 拒绝评论
	MarkAsSpam(id uint) error                                      // 标记为垃圾评论
	UnmarkSpam(id uint) error                                      // 取消垃圾标记
	
	// 评论查询
	GetCommentsByPost(postID uint, offset, limit int) ([]models.Comment, int64, error) // 获取文章评论
	GetCommentsByUser(userID uint, offset, limit int) ([]models.Comment, int64, error) // 获取用户评论
	GetCommentReplies(parentID uint, offset, limit int) ([]models.Comment, int64, error) // 获取评论回复
	GetCommentTree(postID uint) ([]CommentNode, error)             // 获取评论树结构
	
	// 评论统计
	GetCommentStats(postID uint) (*CommentStats, error)           // 获取评论统计
	GetUserCommentStats(userID uint) (*UserCommentStats, error)   // 获取用户评论统计
	
	// 评论点赞
	LikeComment(commentID, userID uint) error                      // 点赞评论
	UnlikeComment(commentID, userID uint) error                    // 取消点赞
	IsCommentLiked(commentID, userID uint) (bool, error)          // 检查是否已点赞
	
	// 评论举报
	ReportComment(commentID, userID uint, reason string) error    // 举报评论
	GetCommentReports(commentID uint) ([]CommentReport, error)    // 获取评论举报
}

// commentService 评论服务实现
type commentService struct {
	db *gorm.DB
}

// NewCommentService 创建评论服务实例
// 参数: db - 数据库连接
// 返回: CommentService - 评论服务接口实例
func NewCommentService(db *gorm.DB) CommentService {
	return &commentService{
		db: db,
	}
}

// CommentFilters 评论筛选条件
type CommentFilters struct {
	PostID     uint       `json:"post_id"`     // 文章ID筛选
	UserID     uint       `json:"user_id"`     // 用户ID筛选
	ParentID   *uint      `json:"parent_id"`   // 父评论ID筛选
	Status     string     `json:"status"`      // 状态筛选
	Keyword    string     `json:"keyword"`     // 关键词搜索
	StartDate  *time.Time `json:"start_date"`  // 开始日期
	EndDate    *time.Time `json:"end_date"`    // 结束日期
	OrderBy    string     `json:"order_by"`    // 排序字段
	OrderDir   string     `json:"order_dir"`   // 排序方向
	IsSpam     *bool      `json:"is_spam"`     // 是否垃圾评论
}

// CommentStats 评论统计信息
type CommentStats struct {
	TotalCount    int `json:"total_count"`    // 总评论数
	ApprovedCount int `json:"approved_count"` // 已审核评论数
	PendingCount  int `json:"pending_count"`  // 待审核评论数
	SpamCount     int `json:"spam_count"`     // 垃圾评论数
	ReplyCount    int `json:"reply_count"`    // 回复数
	LikeCount     int `json:"like_count"`     // 总点赞数
}

// UserCommentStats 用户评论统计
type UserCommentStats struct {
	TotalComments    int `json:"total_comments"`    // 总评论数
	ApprovedComments int `json:"approved_comments"` // 已审核评论数
	TotalLikes       int `json:"total_likes"`       // 获得的总点赞数
	TotalReplies     int `json:"total_replies"`     // 获得的总回复数
	LastCommentAt    *time.Time `json:"last_comment_at,omitempty"` // 最后评论时间
}

// CommentNode 评论树节点
type CommentNode struct {
	Comment  *models.Comment `json:"comment"`  // 评论信息
	Children []CommentNode   `json:"children"` // 子评论
	Depth    int             `json:"depth"`    // 层级深度
}

// CommentReport 评论举报信息
type CommentReport struct {
	ID        uint      `json:"id"`
	CommentID uint      `json:"comment_id"`
	UserID    uint      `json:"user_id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	User      *models.User `json:"user,omitempty"`
}

// 评论基本操作实现

// CreateComment 创建评论
// 参数: comment - 评论模型
// 返回: error - 错误信息
func (s *commentService) CreateComment(comment *models.Comment) error {
	if comment == nil {
		return errors.New("评论信息不能为空")
	}
	
	// 验证必填字段
	if err := s.validateCommentData(comment); err != nil {
		return err
	}
	
	// 检查文章是否存在
	if err := s.checkPostExists(comment.PostID); err != nil {
		return err
	}
	
	// 检查用户是否存在
	if err := s.checkUserExists(comment.UserID); err != nil {
		return err
	}
	
	// 如果是回复，检查父评论是否存在
	if comment.ParentID != nil {
		if err := s.checkParentCommentExists(*comment.ParentID, comment.PostID); err != nil {
			return err
		}
	}
	
	// 设置默认状态
	if comment.Status == 0 {
		comment.Status = models.CommentStatusPending // 默认待审核
	}
	
	// 内容过滤和处理
	comment.Content = s.sanitizeContent(comment.Content)
	
	// 自动垃圾评论检测
	if s.isSpamContent(comment.Content) {
		comment.IsSpam = true
		comment.Status = models.CommentStatusRejected
	}
	
	return s.db.Create(comment).Error
}

// GetCommentByID 根据ID获取评论
// 参数: id - 评论ID
// 返回: *models.Comment - 评论模型, error - 错误信息
func (s *commentService) GetCommentByID(id uint) (*models.Comment, error) {
	if id == 0 {
		return nil, errors.New("评论ID不能为空")
	}
	
	comment := &models.Comment{}
	err := s.db.Preload("User").Preload("Post").Preload("Parent").First(comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}
	
	return comment, nil
}

// UpdateComment 更新评论
// 参数: comment - 评论模型
// 返回: error - 错误信息
func (s *commentService) UpdateComment(comment *models.Comment) error {
	if comment == nil || comment.ID == 0 {
		return errors.New("评论信息不完整")
	}
	
	// 检查评论是否存在
	existingComment := &models.Comment{}
	if err := s.db.First(existingComment, comment.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}
	
	// 验证数据
	if comment.Content != "" {
		if len(comment.Content) > 1000 {
			return errors.New("评论内容不能超过1000个字符")
		}
		comment.Content = s.sanitizeContent(comment.Content)
	}
	
	// 只允许更新特定字段
	updateFields := map[string]interface{}{}
	if comment.Content != "" {
		updateFields["content"] = comment.Content
	}
	if comment.Status != 0 {
		updateFields["status"] = comment.Status
	}
	if comment.IsSpam != existingComment.IsSpam {
		updateFields["is_spam"] = comment.IsSpam
	}
	
	if len(updateFields) == 0 {
		return errors.New("没有需要更新的字段")
	}
	
	return s.db.Model(comment).Updates(updateFields).Error
}

// DeleteComment 删除评论（软删除）
// 参数: id - 评论ID
// 返回: error - 错误信息
func (s *commentService) DeleteComment(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	// 检查评论是否存在
	comment := &models.Comment{}
	if err := s.db.First(comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}
	
	// 软删除评论及其所有回复
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除所有子评论
		if err := tx.Where("parent_id = ?", id).Delete(&models.Comment{}).Error; err != nil {
			return err
		}
		
		// 删除评论本身
		return tx.Delete(comment).Error
	})
}

// ListComments 分页获取评论列表
// 参数: offset - 偏移量, limit - 限制数量, filters - 筛选条件
// 返回: []models.Comment - 评论列表, int64 - 总数量, error - 错误信息
func (s *commentService) ListComments(offset, limit int, filters CommentFilters) ([]models.Comment, int64, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	var comments []models.Comment
	var total int64
	
	// 构建查询
	query := s.db.Model(&models.Comment{})
	
	// 应用筛选条件
	query = s.applyCommentFilters(query, filters)
	
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
	
	// 获取评论列表
	err := query.Preload("User").Preload("Post").Preload("Parent").
		Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", orderBy, orderDir)).
		Find(&comments).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	return comments, total, nil
}

// 评论状态操作实现

// ApproveComment 审核通过评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (s *commentService) ApproveComment(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "approved").Error
}

// RejectComment 拒绝评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (s *commentService) RejectComment(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Update("status", "rejected").Error
}

// MarkAsSpam 标记为垃圾评论
// 参数: id - 评论ID
// 返回: error - 错误信息
func (s *commentService) MarkAsSpam(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	updates := map[string]interface{}{
		"is_spam": true,
		"status":  "rejected",
	}
	
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Updates(updates).Error
}

// UnmarkSpam 取消垃圾标记
// 参数: id - 评论ID
// 返回: error - 错误信息
func (s *commentService) UnmarkSpam(id uint) error {
	if id == 0 {
		return errors.New("评论ID不能为空")
	}
	
	updates := map[string]interface{}{
		"is_spam": false,
		"status":  "pending",
	}
	
	return s.db.Model(&models.Comment{}).Where("id = ?", id).Updates(updates).Error
}

// 评论查询实现

// GetCommentsByPost 获取文章评论
// 参数: postID - 文章ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, int64 - 总数量, error - 错误信息
func (s *commentService) GetCommentsByPost(postID uint, offset, limit int) ([]models.Comment, int64, error) {
	if postID == 0 {
		return nil, 0, errors.New("文章ID不能为空")
	}
	
	filters := CommentFilters{
		PostID: postID,
		Status: "approved", // 只显示已审核的评论
	}
	
	return s.ListComments(offset, limit, filters)
}

// GetCommentsByUser 获取用户评论
// 参数: userID - 用户ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 评论列表, int64 - 总数量, error - 错误信息
func (s *commentService) GetCommentsByUser(userID uint, offset, limit int) ([]models.Comment, int64, error) {
	if userID == 0 {
		return nil, 0, errors.New("用户ID不能为空")
	}
	
	filters := CommentFilters{
		UserID: userID,
	}
	
	return s.ListComments(offset, limit, filters)
}

// GetCommentReplies 获取评论回复
// 参数: parentID - 父评论ID, offset - 偏移量, limit - 限制数量
// 返回: []models.Comment - 回复列表, int64 - 总数量, error - 错误信息
func (s *commentService) GetCommentReplies(parentID uint, offset, limit int) ([]models.Comment, int64, error) {
	if parentID == 0 {
		return nil, 0, errors.New("父评论ID不能为空")
	}
	
	filters := CommentFilters{
		ParentID: &parentID,
		Status:   "approved",
	}
	
	return s.ListComments(offset, limit, filters)
}

// GetCommentTree 获取评论树结构
// 参数: postID - 文章ID
// 返回: []CommentNode - 评论树, error - 错误信息
func (s *commentService) GetCommentTree(postID uint) ([]CommentNode, error) {
	if postID == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	// 获取所有已审核的评论
	var comments []models.Comment
	err := s.db.Preload("User").Where("post_id = ? AND status = ?", postID, "approved").
		Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	
	// 构建评论树
	return s.buildCommentTree(comments, nil, 0), nil
}

// 评论统计实现

// GetCommentStats 获取评论统计
// 参数: postID - 文章ID
// 返回: *CommentStats - 评论统计信息, error - 错误信息
func (s *commentService) GetCommentStats(postID uint) (*CommentStats, error) {
	if postID == 0 {
		return nil, errors.New("文章ID不能为空")
	}
	
	stats := &CommentStats{}
	
	// 总评论数
	var totalCount int64
	s.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&totalCount)
	stats.TotalCount = int(totalCount)
	
	// 已审核评论数
	var approvedCount int64
	s.db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, "approved").Count(&approvedCount)
	stats.ApprovedCount = int(approvedCount)
	
	// 待审核评论数
	var pendingCount int64
	s.db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, "pending").Count(&pendingCount)
	stats.PendingCount = int(pendingCount)
	
	// 垃圾评论数
	var spamCount int64
	s.db.Model(&models.Comment{}).Where("post_id = ? AND is_spam = ?", postID, true).Count(&spamCount)
	stats.SpamCount = int(spamCount)
	
	// 回复数（有父评论的评论）
	var replyCount int64
	s.db.Model(&models.Comment{}).Where("post_id = ? AND parent_id IS NOT NULL", postID).Count(&replyCount)
	stats.ReplyCount = int(replyCount)
	
	// 总点赞数
	var likeCount int64
	s.db.Model(&models.Like{}).Joins("JOIN comments ON likes.target_id = comments.id").
		Where("likes.target_type = ? AND comments.post_id = ?", "comment", postID).Count(&likeCount)
	stats.LikeCount = int(likeCount)
	
	return stats, nil
}

// GetUserCommentStats 获取用户评论统计
// 参数: userID - 用户ID
// 返回: *UserCommentStats - 用户评论统计, error - 错误信息
func (s *commentService) GetUserCommentStats(userID uint) (*UserCommentStats, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	
	stats := &UserCommentStats{}
	
	// 总评论数
	var totalComments int64
	s.db.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&totalComments)
	stats.TotalComments = int(totalComments)
	
	// 已审核评论数
	var approvedComments int64
	s.db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", userID, "approved").Count(&approvedComments)
	stats.ApprovedComments = int(approvedComments)
	
	// 获得的总点赞数
	var totalLikes int64
	s.db.Model(&models.Like{}).Joins("JOIN comments ON likes.target_id = comments.id").
		Where("likes.target_type = ? AND comments.user_id = ?", "comment", userID).Count(&totalLikes)
	stats.TotalLikes = int(totalLikes)
	
	// 获得的总回复数
	var totalReplies int64
	s.db.Model(&models.Comment{}).Joins("JOIN comments AS parent ON comments.parent_id = parent.id").
		Where("parent.user_id = ?", userID).Count(&totalReplies)
	stats.TotalReplies = int(totalReplies)
	
	// 最后评论时间
	var lastComment models.Comment
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").First(&lastComment).Error; err == nil {
		stats.LastCommentAt = &lastComment.CreatedAt
	}
	
	return stats, nil
}

// 评论点赞实现

// LikeComment 点赞评论
// 参数: commentID - 评论ID, userID - 用户ID
// 返回: error - 错误信息
func (s *commentService) LikeComment(commentID, userID uint) error {
	if commentID == 0 {
		return errors.New("评论ID不能为空")
	}
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	// 检查是否已经点赞
	isLiked, err := s.IsCommentLiked(commentID, userID)
	if err != nil {
		return err
	}
	if isLiked {
		return errors.New("已经点赞过该评论")
	}
	
	// 检查评论是否存在
	if err := s.checkCommentExists(commentID); err != nil {
		return err
	}
	
	// 创建点赞记录
	like := &models.Like{
		UserID:     userID,
		TargetType: "comment",
		TargetID:   commentID,
	}
	
	return s.db.Create(like).Error
}

// UnlikeComment 取消点赞
// 参数: commentID - 评论ID, userID - 用户ID
// 返回: error - 错误信息
func (s *commentService) UnlikeComment(commentID, userID uint) error {
	if commentID == 0 {
		return errors.New("评论ID不能为空")
	}
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	// 删除点赞记录
	result := s.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, "comment", commentID).
		Delete(&models.Like{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("未找到点赞记录")
	}
	
	return nil
}

// IsCommentLiked 检查是否已点赞
// 参数: commentID - 评论ID, userID - 用户ID
// 返回: bool - 是否已点赞, error - 错误信息
func (s *commentService) IsCommentLiked(commentID, userID uint) (bool, error) {
	if commentID == 0 || userID == 0 {
		return false, nil
	}
	
	var count int64
	err := s.db.Model(&models.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, "comment", commentID).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// 评论举报实现

// ReportComment 举报评论
// 参数: commentID - 评论ID, userID - 用户ID, reason - 举报原因
// 返回: error - 错误信息
func (s *commentService) ReportComment(commentID, userID uint, reason string) error {
	if commentID == 0 {
		return errors.New("评论ID不能为空")
	}
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	if reason == "" {
		return errors.New("举报原因不能为空")
	}
	
	// 检查评论是否存在
	if err := s.checkCommentExists(commentID); err != nil {
		return err
	}
	
	// 检查是否已经举报过
	var count int64
	s.db.Table("comment_reports").Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count)
	if count > 0 {
		return errors.New("已经举报过该评论")
	}
	
	// 创建举报记录（这里简化处理，实际应该有专门的举报表）
	report := map[string]interface{}{
		"comment_id": commentID,
		"user_id":    userID,
		"reason":     reason,
		"created_at": time.Now(),
	}
	
	return s.db.Table("comment_reports").Create(report).Error
}

// GetCommentReports 获取评论举报
// 参数: commentID - 评论ID
// 返回: []CommentReport - 举报列表, error - 错误信息
func (s *commentService) GetCommentReports(commentID uint) ([]CommentReport, error) {
	if commentID == 0 {
		return nil, errors.New("评论ID不能为空")
	}
	
	var reports []CommentReport
	err := s.db.Table("comment_reports").Preload("User").
		Where("comment_id = ?", commentID).Find(&reports).Error
	
	if err != nil {
		return nil, err
	}
	
	return reports, nil
}

// 辅助方法

// validateCommentData 验证评论数据
// 参数: comment - 评论模型
// 返回: error - 验证错误信息
func (s *commentService) validateCommentData(comment *models.Comment) error {
	if comment.Content == "" {
		return errors.New("评论内容不能为空")
	}
	if len(comment.Content) > 1000 {
		return errors.New("评论内容不能超过1000个字符")
	}
	if comment.PostID == 0 {
		return errors.New("文章ID不能为空")
	}
	if comment.UserID == 0 {
		return errors.New("用户ID不能为空")
	}
	
	return nil
}

// checkPostExists 检查文章是否存在
// 参数: postID - 文章ID
// 返回: error - 错误信息
func (s *commentService) checkPostExists(postID uint) error {
	var count int64
	s.db.Model(&models.Post{}).Where("id = ?", postID).Count(&count)
	if count == 0 {
		return errors.New("文章不存在")
	}
	return nil
}

// checkUserExists 检查用户是否存在
// 参数: userID - 用户ID
// 返回: error - 错误信息
func (s *commentService) checkUserExists(userID uint) error {
	var count int64
	s.db.Model(&models.User{}).Where("id = ?", userID).Count(&count)
	if count == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

// checkCommentExists 检查评论是否存在
// 参数: commentID - 评论ID
// 返回: error - 错误信息
func (s *commentService) checkCommentExists(commentID uint) error {
	var count int64
	s.db.Model(&models.Comment{}).Where("id = ?", commentID).Count(&count)
	if count == 0 {
		return errors.New("评论不存在")
	}
	return nil
}

// checkParentCommentExists 检查父评论是否存在且属于同一文章
// 参数: parentID - 父评论ID, postID - 文章ID
// 返回: error - 错误信息
func (s *commentService) checkParentCommentExists(parentID, postID uint) error {
	var count int64
	s.db.Model(&models.Comment{}).Where("id = ? AND post_id = ?", parentID, postID).Count(&count)
	if count == 0 {
		return errors.New("父评论不存在或不属于该文章")
	}
	return nil
}

// sanitizeContent 清理评论内容
// 参数: content - 原始内容
// 返回: string - 清理后的内容
func (s *commentService) sanitizeContent(content string) string {
	// 简单的内容清理：去除首尾空格
	content = strings.TrimSpace(content)
	
	// TODO: 实现更复杂的内容过滤，如HTML标签过滤、敏感词过滤等
	
	return content
}

// isSpamContent 检测是否为垃圾内容
// 参数: content - 评论内容
// 返回: bool - 是否为垃圾内容
func (s *commentService) isSpamContent(content string) bool {
	// 简单的垃圾内容检测
	spamKeywords := []string{"广告", "推广", "加微信", "QQ群", "免费领取"}
	
	contentLower := strings.ToLower(content)
	for _, keyword := range spamKeywords {
		if strings.Contains(contentLower, keyword) {
			return true
		}
	}
	
	// 检测重复字符
	if s.hasRepeatedChars(content, 5) {
		return true
	}
	
	return false
}

// hasRepeatedChars 检测是否有重复字符
// 参数: content - 内容, threshold - 重复阈值
// 返回: bool - 是否有重复字符
func (s *commentService) hasRepeatedChars(content string, threshold int) bool {
	if len(content) < threshold {
		return false
	}
	
	for i := 0; i <= len(content)-threshold; i++ {
		char := content[i]
		count := 1
		for j := i + 1; j < len(content) && j < i+threshold; j++ {
			if content[j] == char {
				count++
			} else {
				break
			}
		}
		if count >= threshold {
			return true
		}
	}
	
	return false
}

// buildCommentTree 构建评论树
// 参数: comments - 评论列表, parentID - 父评论ID, depth - 当前深度
// 返回: []CommentNode - 评论树节点
func (s *commentService) buildCommentTree(comments []models.Comment, parentID *uint, depth int) []CommentNode {
	var nodes []CommentNode
	
	for _, comment := range comments {
		// 检查是否为当前层级的评论
		if (parentID == nil && comment.ParentID == nil) || 
		   (parentID != nil && comment.ParentID != nil && *comment.ParentID == *parentID) {
			
			node := CommentNode{
				Comment: &comment,
				Depth:   depth,
			}
			
			// 递归构建子评论（限制深度避免无限递归）
			if depth < 5 {
				node.Children = s.buildCommentTree(comments, &comment.ID, depth+1)
			}
			
			nodes = append(nodes, node)
		}
	}
	
	return nodes
}

// applyCommentFilters 应用评论筛选条件
// 参数: query - GORM查询对象, filters - 筛选条件
// 返回: *gorm.DB - 应用筛选后的查询对象
func (s *commentService) applyCommentFilters(query *gorm.DB, filters CommentFilters) *gorm.DB {
	// 文章筛选
	if filters.PostID > 0 {
		query = query.Where("post_id = ?", filters.PostID)
	}
	
	// 用户筛选
	if filters.UserID > 0 {
		query = query.Where("user_id = ?", filters.UserID)
	}
	
	// 父评论筛选
	if filters.ParentID != nil {
		query = query.Where("parent_id = ?", *filters.ParentID)
	}
	
	// 状态筛选
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	
	// 垃圾评论筛选
	if filters.IsSpam != nil {
		query = query.Where("is_spam = ?", *filters.IsSpam)
	}
	
	// 关键词搜索
	if filters.Keyword != "" {
		keyword := "%" + filters.Keyword + "%"
		query = query.Where("content LIKE ?", keyword)
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
// 04_unit_exercises/level6_comprehensive.go - Level 6 综合实战练习
// 对应文档：03_GORM单元练习_基础技能训练.md

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 用户模型
type User struct {
	BaseModel
	Username     string     `gorm:"uniqueIndex:idx_username;size:50;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex:idx_email;size:100;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	FirstName    string     `gorm:"size:50;not null" json:"first_name"`
	LastName     string     `gorm:"size:50;not null" json:"last_name"`
	Avatar       string     `gorm:"size:255" json:"avatar"`
	Bio          string     `gorm:"type:text" json:"bio"`
	Website      string     `gorm:"size:255" json:"website"`
	Location     string     `gorm:"size:100" json:"location"`
	BirthDate    *time.Time `json:"birth_date"`
	Gender       string     `gorm:"size:10" json:"gender"`
	Phone        string     `gorm:"size:20" json:"phone"`
	Status       string     `gorm:"size:20;default:'active';index:idx_status" json:"status"`
	Role         string     `gorm:"size:20;default:'user';index:idx_role" json:"role"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	LastLoginAt  *time.Time `gorm:"index:idx_last_login" json:"last_login_at"`
	LoginCount   int        `gorm:"default:0" json:"login_count"`
	
	// 统计字段
	PostCount     int `gorm:"default:0" json:"post_count"`
	CommentCount  int `gorm:"default:0" json:"comment_count"`
	FollowerCount int `gorm:"default:0" json:"follower_count"`
	FollowingCount int `gorm:"default:0" json:"following_count"`
	
	// 关联关系
	Profile   *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Posts     []Post       `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments  []Comment    `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
	Likes     []Like       `gorm:"foreignKey:UserID" json:"likes,omitempty"`
	Followers []Follow     `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
	Following []Follow     `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

// 用户资料扩展
type UserProfile struct {
	BaseModel
	UserID       uint   `gorm:"uniqueIndex:idx_user_profile;not null" json:"user_id"`
	Company      string `gorm:"size:100" json:"company"`
	JobTitle     string `gorm:"size:100" json:"job_title"`
	Education    string `gorm:"size:200" json:"education"`
	Skills       string `gorm:"type:text" json:"skills"`
	Experience   int    `gorm:"default:0" json:"experience"`
	SalaryRange  string `gorm:"size:50" json:"salary_range"`
	Languages    string `gorm:"type:text" json:"languages"`
	Interests    string `gorm:"type:text" json:"interests"`
	SocialLinks  string `gorm:"type:text" json:"social_links"`
	PrivacyLevel string `gorm:"size:20;default:'public'" json:"privacy_level"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// 分类模型
type Category struct {
	BaseModel
	Name        string `gorm:"size:100;not null;index:idx_category_name" json:"name"`
	Slug        string `gorm:"uniqueIndex:idx_category_slug;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `gorm:"size:100" json:"icon"`
	Color       string `gorm:"size:7;default:'#007bff'" json:"color"`
	ParentID    *uint  `gorm:"index:idx_parent" json:"parent_id"`
	Level       int    `gorm:"default:1;index:idx_level" json:"level"`
	SortOrder   int    `gorm:"default:0;index:idx_sort" json:"sort_order"`
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`
	PostCount   int    `gorm:"default:0" json:"post_count"`
	
	// 关联关系
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Posts    []Post     `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

// 标签模型
type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex:idx_tag_name;size:50;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex:idx_tag_slug;size:50;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"size:7;default:'#007bff'" json:"color"`
	UsageCount  int    `gorm:"default:0;index:idx_usage" json:"usage_count"`
	IsActive    bool   `gorm:"default:true;index:idx_active" json:"is_active"`
	
	// 关联关系
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// 文章模型
type Post struct {
	BaseModel
	Title         string     `gorm:"size:200;not null;index:idx_title" json:"title"`
	Slug          string     `gorm:"uniqueIndex:idx_post_slug;size:200;not null" json:"slug"`
	Content       string     `gorm:"type:text;not null" json:"content"`
	Excerpt       string     `gorm:"size:500" json:"excerpt"`
	FeaturedImage string     `gorm:"size:255" json:"featured_image"`
	Status        string     `gorm:"size:20;default:'draft';index:idx_status" json:"status"`
	Type          string     `gorm:"size:20;default:'post';index:idx_type" json:"type"`
	Format        string     `gorm:"size:20;default:'standard'" json:"format"`
	ViewCount     int        `gorm:"default:0;index:idx_views" json:"view_count"`
	LikeCount     int        `gorm:"default:0;index:idx_likes" json:"like_count"`
	CommentCount  int        `gorm:"default:0;index:idx_comments" json:"comment_count"`
	ShareCount    int        `gorm:"default:0" json:"share_count"`
	PublishedAt   *time.Time `gorm:"index:idx_published" json:"published_at"`
	Rating        float64    `gorm:"precision:3;scale:2;default:0;index:idx_rating" json:"rating"`
	Featured      bool       `gorm:"default:false;index:idx_featured" json:"featured"`
	Sticky        bool       `gorm:"default:false;index:idx_sticky" json:"sticky"`
	AllowComments bool       `gorm:"default:true" json:"allow_comments"`
	MetaTitle     string     `gorm:"size:200" json:"meta_title"`
	MetaDescription string   `gorm:"size:500" json:"meta_description"`
	MetaKeywords  string     `gorm:"size:255" json:"meta_keywords"`
	
	// 外键
	AuthorID   uint  `gorm:"not null;index:idx_author" json:"author_id"`
	CategoryID *uint `gorm:"index:idx_category" json:"category_id"`
	
	// 关联关系
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	Meta     []PostMeta `gorm:"foreignKey:PostID" json:"meta,omitempty"`
}

// 文章元数据
type PostMeta struct {
	BaseModel
	PostID    uint   `gorm:"not null;index:idx_post_meta" json:"post_id"`
	MetaKey   string `gorm:"size:100;not null;index:idx_meta_key" json:"meta_key"`
	MetaValue string `gorm:"type:text" json:"meta_value"`
	
	// 关联关系
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// 评论模型
type Comment struct {
	BaseModel
	Content    string `gorm:"type:text;not null" json:"content"`
	Status     string `gorm:"size:20;default:'pending';index:idx_status" json:"status"`
	Type       string `gorm:"size:20;default:'comment'" json:"type"`
	LikeCount  int    `gorm:"default:0;index:idx_likes" json:"like_count"`
	ParentID   *uint  `gorm:"index:idx_parent" json:"parent_id"`
	Level      int    `gorm:"default:1;index:idx_level" json:"level"`
	UserAgent  string `gorm:"size:255" json:"user_agent"`
	UserIP     string `gorm:"size:45" json:"user_ip"`
	IsSpam     bool   `gorm:"default:false;index:idx_spam" json:"is_spam"`
	
	// 外键
	PostID   uint `gorm:"not null;index:idx_post" json:"post_id"`
	AuthorID uint `gorm:"not null;index:idx_author" json:"author_id"`
	
	// 关联关系
	Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Likes    []Like    `gorm:"foreignKey:CommentID" json:"likes,omitempty"`
}

// 点赞模型
type Like struct {
	BaseModel
	UserID    uint   `gorm:"not null;index:idx_user_like" json:"user_id"`
	PostID    *uint  `gorm:"index:idx_post_like" json:"post_id"`
	CommentID *uint  `gorm:"index:idx_comment_like" json:"comment_id"`
	Type      string `gorm:"size:20;default:'like';index:idx_like_type" json:"type"`
	UserIP    string `gorm:"size:45" json:"user_ip"`
	
	// 关联关系
	User    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post    *Post    `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Comment *Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
}

// 关注关系模型
type Follow struct {
	BaseModel
	FollowerID  uint   `gorm:"not null;index:idx_follower" json:"follower_id"`
	FollowingID uint   `gorm:"not null;index:idx_following" json:"following_id"`
	Status      string `gorm:"size:20;default:'active'" json:"status"`
	
	// 关联关系
	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

// 通知模型
type Notification struct {
	BaseModel
	UserID      uint   `gorm:"not null;index:idx_user_notification" json:"user_id"`
	Type        string `gorm:"size:50;not null;index:idx_notification_type" json:"type"`
	Title       string `gorm:"size:200;not null" json:"title"`
	Content     string `gorm:"type:text" json:"content"`
	Data        string `gorm:"type:text" json:"data"`
	IsRead      bool   `gorm:"default:false;index:idx_read" json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	RelatedID   *uint  `gorm:"index:idx_related" json:"related_id"`
	RelatedType string `gorm:"size:50" json:"related_type"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// 系统设置模型
type Setting struct {
	BaseModel
	Key         string `gorm:"uniqueIndex:idx_setting_key;size:100;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Type        string `gorm:"size:20;default:'string'" json:"type"`
	Description string `gorm:"size:255" json:"description"`
	Group       string `gorm:"size:50;index:idx_setting_group" json:"group"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
}

// 数据库初始化
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("level6_comprehensive.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取数据库连接失败:", err)
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	err = db.AutoMigrate(
		&User{}, &UserProfile{}, &Category{}, &Tag{}, &Post{}, &PostMeta{},
		&Comment{}, &Like{}, &Follow{}, &Notification{}, &Setting{},
	)
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 创建索引
	createIndexes(db)

	return db
}

// 创建索引
func createIndexes(db *gorm.DB) {
	// 复合索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status_role ON users(status, role)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_status_published ON posts(status, published_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_author_category ON posts(author_id, category_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_status ON comments(post_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_likes_user_post ON likes(user_id, post_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_follower_following ON follows(follower_id, following_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read)")
	
	// 唯一索引
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_like ON likes(user_id, post_id, comment_id) WHERE post_id IS NOT NULL OR comment_id IS NOT NULL")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_follow ON follows(follower_id, following_id)")
	
	fmt.Println("✓ 索引创建完成")
}

// 钩子函数实现

// 用户钩子
func (u *User) AfterCreate(tx *gorm.DB) error {
	// 创建用户资料
	profile := UserProfile{
		UserID:       u.ID,
		PrivacyLevel: "public",
	}
	return tx.Create(&profile).Error
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 更新统计信息
	if u.ID != 0 {
		// 更新文章数量
		var postCount int64
		tx.Model(&Post{}).Where("author_id = ? AND status = ?", u.ID, "published").Count(&postCount)
		u.PostCount = int(postCount)
		
		// 更新评论数量
		var commentCount int64
		tx.Model(&Comment{}).Where("author_id = ? AND status = ?", u.ID, "approved").Count(&commentCount)
		u.CommentCount = int(commentCount)
		
		// 更新关注数量
		var followerCount, followingCount int64
		tx.Model(&Follow{}).Where("following_id = ?", u.ID).Count(&followerCount)
		tx.Model(&Follow{}).Where("follower_id = ?", u.ID).Count(&followingCount)
		u.FollowerCount = int(followerCount)
		u.FollowingCount = int(followingCount)
	}
	return nil
}

// 文章钩子
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 生成摘要
	if p.Excerpt == "" && len(p.Content) > 200 {
		p.Excerpt = p.Content[:200] + "..."
	}
	
	// 设置发布时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	
	return nil
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新分类文章数量
	if p.CategoryID != nil {
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))
	}
	
	// 更新用户文章数量
	tx.Model(&User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("post_count + ?", 1))
	
	return nil
}

func (p *Post) AfterUpdate(tx *gorm.DB) error {
	// 如果状态改为已发布，设置发布时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		tx.Model(p).Update("published_at", now)
	}
	return nil
}

func (p *Post) AfterDelete(tx *gorm.DB) error {
	// 更新分类文章数量
	if p.CategoryID != nil {
		tx.Model(&Category{}).Where("id = ?", *p.CategoryID).UpdateColumn("post_count", gorm.Expr("post_count - ?", 1))
	}
	
	// 更新用户文章数量
	tx.Model(&User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("post_count - ?", 1))
	
	return nil
}

// 评论钩子
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	// 更新文章评论数量
	tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	
	// 更新用户评论数量
	tx.Model(&User{}).Where("id = ?", c.AuthorID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	
	// 创建通知
	var post Post
	if err := tx.First(&post, c.PostID).Error; err == nil {
		if post.AuthorID != c.AuthorID { // 不给自己发通知
			notification := Notification{
				UserID:      post.AuthorID,
				Type:        "comment",
				Title:       "新评论",
				Content:     fmt.Sprintf("您的文章《%s》收到了新评论", post.Title),
				RelatedID:   &c.ID,
				RelatedType: "comment",
			}
			tx.Create(&notification)
		}
	}
	
	return nil
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 更新文章评论数量
	tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	
	// 更新用户评论数量
	tx.Model(&User{}).Where("id = ?", c.AuthorID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	
	return nil
}

// 点赞钩子
func (l *Like) AfterCreate(tx *gorm.DB) error {
	if l.PostID != nil {
		// 更新文章点赞数量
		tx.Model(&Post{}).Where("id = ?", *l.PostID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
		
		// 创建通知
		var post Post
		if err := tx.First(&post, *l.PostID).Error; err == nil {
			if post.AuthorID != l.UserID {
				notification := Notification{
					UserID:      post.AuthorID,
					Type:        "like",
					Title:       "新点赞",
					Content:     fmt.Sprintf("您的文章《%s》收到了新点赞", post.Title),
					RelatedID:   l.PostID,
					RelatedType: "post",
				}
				tx.Create(&notification)
			}
		}
	}
	
	if l.CommentID != nil {
		// 更新评论点赞数量
		tx.Model(&Comment{}).Where("id = ?", *l.CommentID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
	}
	
	return nil
}

func (l *Like) AfterDelete(tx *gorm.DB) error {
	if l.PostID != nil {
		tx.Model(&Post{}).Where("id = ?", *l.PostID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	}
	
	if l.CommentID != nil {
		tx.Model(&Comment{}).Where("id = ?", *l.CommentID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	}
	
	return nil
}

// 关注钩子
func (f *Follow) AfterCreate(tx *gorm.DB) error {
	// 更新关注数量
	tx.Model(&User{}).Where("id = ?", f.FollowerID).UpdateColumn("following_count", gorm.Expr("following_count + ?", 1))
	tx.Model(&User{}).Where("id = ?", f.FollowingID).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
	
	// 创建通知
	notification := Notification{
		UserID:      f.FollowingID,
		Type:        "follow",
		Title:       "新关注者",
		Content:     "您有新的关注者",
		RelatedID:   &f.FollowerID,
		RelatedType: "user",
	}
	tx.Create(&notification)
	
	return nil
}

func (f *Follow) AfterDelete(tx *gorm.DB) error {
	// 更新关注数量
	tx.Model(&User{}).Where("id = ?", f.FollowerID).UpdateColumn("following_count", gorm.Expr("following_count - ?", 1))
	tx.Model(&User{}).Where("id = ?", f.FollowingID).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
	
	return nil
}

// 业务逻辑函数

// 用户管理
type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user *User) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
	var user User
	err := s.db.Preload("Profile").First(&user, id).Error
	return &user, err
}

func (s *UserService) GetUserWithStats(id uint) (*User, error) {
	var user User
	err := s.db.Preload("Profile").
		Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "published").Order("created_at DESC").Limit(5)
		}).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "approved").Order("created_at DESC").Limit(5)
		}).
		First(&user, id).Error
	return &user, err
}

func (s *UserService) UpdateUserProfile(userID uint, profile *UserProfile) error {
	return s.db.Model(&UserProfile{}).Where("user_id = ?", userID).Updates(profile).Error
}

func (s *UserService) FollowUser(followerID, followingID uint) error {
	if followerID == followingID {
		return fmt.Errorf("不能关注自己")
	}
	
	follow := Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		Status:      "active",
	}
	
	return s.db.Create(&follow).Error
}

func (s *UserService) UnfollowUser(followerID, followingID uint) error {
	return s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&Follow{}).Error
}

func (s *UserService) GetUserFollowers(userID uint, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64
	
	offset := (page - 1) * pageSize
	
	// 获取总数
	s.db.Model(&User{}).Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).Count(&total)
	
	// 获取数据
	err := s.db.Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Offset(offset).Limit(pageSize).Find(&users).Error
	
	return users, total, err
}

// 文章管理
type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(post *Post) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		
		// 更新标签使用次数
		for _, tag := range post.Tags {
			tx.Model(&Tag{}).Where("id = ?", tag.ID).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1))
		}
		
		return nil
	})
}

func (s *PostService) GetPostBySlug(slug string) (*Post, error) {
	var post Post
	err := s.db.Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ? AND parent_id IS NULL", "approved").Order("created_at ASC")
		}).
		Preload("Comments.Author").
		Preload("Comments.Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "approved").Order("created_at ASC")
		}).
		Preload("Comments.Children.Author").
		Where("slug = ?", slug).First(&post).Error
	
	if err == nil {
		// 增加浏览量
		s.db.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	}
	
	return &post, err
}

func (s *PostService) GetPostsByCategory(categorySlug string, page, pageSize int) ([]Post, int64, error) {
	var posts []Post
	var total int64
	
	offset := (page - 1) * pageSize
	
	// 获取总数
	s.db.Model(&Post{}).Joins("JOIN categories ON posts.category_id = categories.id").
		Where("categories.slug = ? AND posts.status = ?", categorySlug, "published").Count(&total)
	
	// 获取数据
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN categories ON posts.category_id = categories.id").
		Where("categories.slug = ? AND posts.status = ?", categorySlug, "published").
		Order("posts.sticky DESC, posts.published_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error
	
	return posts, total, err
}

func (s *PostService) GetPostsByTag(tagSlug string, page, pageSize int) ([]Post, int64, error) {
	var posts []Post
	var total int64
	
	offset := (page - 1) * pageSize
	
	// 获取总数
	s.db.Model(&Post{}).Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.slug = ? AND posts.status = ?", tagSlug, "published").Count(&total)
	
	// 获取数据
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.slug = ? AND posts.status = ?", tagSlug, "published").
		Order("posts.published_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error
	
	return posts, total, err
}

func (s *PostService) SearchPosts(keyword string, page, pageSize int) ([]Post, int64, error) {
	var posts []Post
	var total int64
	
	offset := (page - 1) * pageSize
	searchTerm := "%" + keyword + "%"
	
	// 获取总数
	s.db.Model(&Post{}).Where("(title LIKE ? OR content LIKE ?) AND status = ?", searchTerm, searchTerm, "published").Count(&total)
	
	// 获取数据
	err := s.db.Preload("Author").Preload("Category").Preload("Tags").
		Where("(title LIKE ? OR content LIKE ?) AND status = ?", searchTerm, searchTerm, "published").
		Order("view_count DESC, published_at DESC").
		Offset(offset).Limit(pageSize).Find(&posts).Error
	
	return posts, total, err
}

func (s *PostService) LikePost(userID, postID uint) error {
	// 检查是否已经点赞
	var existingLike Like
	if err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike).Error; err == nil {
		return fmt.Errorf("已经点赞过了")
	}
	
	like := Like{
		UserID: userID,
		PostID: &postID,
		Type:   "like",
	}
	
	return s.db.Create(&like).Error
}

func (s *PostService) UnlikePost(userID, postID uint) error {
	return s.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Like{}).Error
}

// 评论管理
type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) CreateComment(comment *Comment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查文章是否允许评论
		var post Post
		if err := tx.First(&post, comment.PostID).Error; err != nil {
			return err
		}
		
		if !post.AllowComments {
			return fmt.Errorf("该文章不允许评论")
		}
		
		// 设置评论层级
		if comment.ParentID != nil {
			var parentComment Comment
			if err := tx.First(&parentComment, *comment.ParentID).Error; err != nil {
				return err
			}
			comment.Level = parentComment.Level + 1
		}
		
		return tx.Create(comment).Error
	})
}

func (s *CommentService) ApproveComment(commentID uint) error {
	return s.db.Model(&Comment{}).Where("id = ?", commentID).Update("status", "approved").Error
}

func (s *CommentService) RejectComment(commentID uint) error {
	return s.db.Model(&Comment{}).Where("id = ?", commentID).Update("status", "rejected").Error
}

func (s *CommentService) MarkAsSpam(commentID uint) error {
	return s.db.Model(&Comment{}).Where("id = ?", commentID).Updates(map[string]interface{}{
		"status":  "spam",
		"is_spam": true,
	}).Error
}

// 通知管理
type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64
	
	offset := (page - 1) * pageSize
	
	// 获取总数
	s.db.Model(&Notification{}).Where("user_id = ?", userID).Count(&total)
	
	// 获取数据
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&notifications).Error
	
	return notifications, total, err
}

func (s *NotificationService) MarkAsRead(notificationID uint) error {
	now := time.Now()
	return s.db.Model(&Notification{}).Where("id = ?", notificationID).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}).Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}).Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// 统计分析
type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

type DashboardStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	TotalPosts    int64 `json:"total_posts"`
	PublishedPosts int64 `json:"published_posts"`
	TotalComments int64 `json:"total_comments"`
	ApprovedComments int64 `json:"approved_comments"`
	TotalViews    int64 `json:"total_views"`
	TotalLikes    int64 `json:"total_likes"`
}

func (s *AnalyticsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}
	
	// 用户统计
	s.db.Model(&User{}).Count(&stats.TotalUsers)
	s.db.Model(&User{}).Where("status = ?", "active").Count(&stats.ActiveUsers)
	
	// 文章统计
	s.db.Model(&Post{}).Count(&stats.TotalPosts)
	s.db.Model(&Post{}).Where("status = ?", "published").Count(&stats.PublishedPosts)
	
	// 评论统计
	s.db.Model(&Comment{}).Count(&stats.TotalComments)
	s.db.Model(&Comment{}).Where("status = ?", "approved").Count(&stats.ApprovedComments)
	
	// 浏览量和点赞统计
	s.db.Model(&Post{}).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews)
	s.db.Model(&Post{}).Select("COALESCE(SUM(like_count), 0)").Scan(&stats.TotalLikes)
	
	return stats, nil
}

type PopularPost struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ViewCount int    `json:"view_count"`
	LikeCount int    `json:"like_count"`
	Author    string `json:"author"`
}

func (s *AnalyticsService) GetPopularPosts(limit int) ([]PopularPost, error) {
	var posts []PopularPost
	
	err := s.db.Table("posts p").
		Select("p.id, p.title, p.slug, p.view_count, p.like_count, u.username as author").
		Joins("JOIN users u ON p.author_id = u.id").
		Where("p.status = ?", "published").
		Order("p.view_count DESC, p.like_count DESC").
		Limit(limit).Scan(&posts).Error
	
	return posts, err
}

type ActiveUser struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	PostCount    int    `json:"post_count"`
	CommentCount int    `json:"comment_count"`
	LikeCount    int    `json:"like_count"`
}

func (s *AnalyticsService) GetActiveUsers(limit int) ([]ActiveUser, error) {
	var users []ActiveUser
	
	err := s.db.Table("users u").
		Select("u.id, u.username, u.post_count, u.comment_count, COALESCE(l.like_count, 0) as like_count").
		Joins("LEFT JOIN (SELECT user_id, COUNT(*) as like_count FROM likes GROUP BY user_id) l ON u.id = l.user_id").
		Where("u.status = ?", "active").
		Order("(u.post_count + u.comment_count + COALESCE(l.like_count, 0)) DESC").
		Limit(limit).Scan(&users).Error
	
	return users, err
}

// 生成测试数据
func generateComprehensiveTestData(db *gorm.DB) {
	fmt.Println("开始生成综合测试数据...")
	start := time.Now()
	
	// 生成用户数据
	users := make([]User, 50)
	for i := 0; i < 50; i++ {
		users[i] = User{
			Username:     fmt.Sprintf("user%d", i+1),
			Email:        fmt.Sprintf("user%d@example.com", i+1),
			PasswordHash: "hashed_password",
			FirstName:    fmt.Sprintf("First%d", i+1),
			LastName:     fmt.Sprintf("Last%d", i+1),
			Bio:          fmt.Sprintf("这是用户%d的个人简介", i+1),
			Location:     []string{"北京", "上海", "深圳", "广州", "杭州"}[rand.Intn(5)],
			Gender:       []string{"male", "female", "other"}[rand.Intn(3)],
			Status:       "active",
			Role:         []string{"user", "author", "moderator"}[rand.Intn(3)],
			EmailVerified: rand.Float32() > 0.2,
			LoginCount:   rand.Intn(100),
		}
		if rand.Float32() > 0.3 {
			lastLogin := time.Now().AddDate(0, 0, -rand.Intn(30))
			users[i].LastLoginAt = &lastLogin
		}
	}
	db.Create(&users)
	fmt.Printf("✓ 用户数据: %d条\n", len(users))
	
	// 生成分类数据
	categories := []Category{
		{Name: "技术", Slug: "tech", Description: "技术相关文章", Icon: "code", Color: "#007bff", Level: 1, SortOrder: 1, IsActive: true},
		{Name: "生活", Slug: "life", Description: "生活分享", Icon: "heart", Color: "#28a745", Level: 1, SortOrder: 2, IsActive: true},
		{Name: "旅游", Slug: "travel", Description: "旅游攻略", Icon: "map", Color: "#ffc107", Level: 1, SortOrder: 3, IsActive: true},
		{Name: "美食", Slug: "food", Description: "美食推荐", Icon: "utensils", Color: "#fd7e14", Level: 1, SortOrder: 4, IsActive: true},
		{Name: "娱乐", Slug: "entertainment", Description: "娱乐资讯", Icon: "film", Color: "#e83e8c", Level: 1, SortOrder: 5, IsActive: true},
	}
	db.Create(&categories)
	fmt.Printf("✓ 分类数据: %d条\n", len(categories))
	
	// 生成标签数据
	tags := []Tag{
		{Name: "Go", Slug: "go", Description: "Go语言相关", Color: "#00ADD8", IsActive: true},
		{Name: "Python", Slug: "python", Description: "Python编程", Color: "#3776AB", IsActive: true},
		{Name: "JavaScript", Slug: "javascript", Description: "JavaScript开发", Color: "#F7DF1E", IsActive: true},
		{Name: "数据库", Slug: "database", Description: "数据库技术", Color: "#336791", IsActive: true},
		{Name: "前端", Slug: "frontend", Description: "前端开发", Color: "#61DAFB", IsActive: true},
		{Name: "后端", Slug: "backend", Description: "后端开发", Color: "#68217A", IsActive: true},
		{Name: "教程", Slug: "tutorial", Description: "教程文章", Color: "#FF6B6B", IsActive: true},
		{Name: "实战", Slug: "practice", Description: "实战项目", Color: "#4ECDC4", IsActive: true},
		{Name: "经验分享", Slug: "experience", Description: "经验分享", Color: "#95A5A6", IsActive: true},
		{Name: "工具推荐", Slug: "tools", Description: "工具推荐", Color: "#9B59B6", IsActive: true},
	}
	db.Create(&tags)
	fmt.Printf("✓ 标签数据: %d条\n", len(tags))
	
	// 生成文章数据
	posts := make([]Post, 200)
	statuses := []string{"published", "draft", "archived"}
	types := []string{"post", "page"}
	formats := []string{"standard", "video", "gallery", "quote"}
	
	for i := 0; i < 200; i++ {
		authorID := uint(rand.Intn(50) + 1)
		categoryID := uint(rand.Intn(5) + 1)
		status := statuses[rand.Intn(len(statuses))]
		
		posts[i] = Post{
			Title:         fmt.Sprintf("精彩文章标题 %d - 深度解析技术要点", i+1),
			Slug:          fmt.Sprintf("awesome-post-%d", i+1),
			Content:       fmt.Sprintf("这是文章 %d 的详细内容，包含了丰富的技术信息和实用的代码示例。文章深入浅出地讲解了相关概念，并提供了实际的应用场景和最佳实践。", i+1),
			Excerpt:       fmt.Sprintf("文章 %d 的精彩摘要，概括了主要内容和核心观点", i+1),
			FeaturedImage: fmt.Sprintf("/images/post-%d.jpg", i+1),
			Status:        status,
			Type:          types[rand.Intn(len(types))],
			Format:        formats[rand.Intn(len(formats))],
			ViewCount:     rand.Intn(5000),
			LikeCount:     rand.Intn(500),
			CommentCount:  rand.Intn(50),
			ShareCount:    rand.Intn(100),
			Rating:        float64(rand.Intn(40))/10.0 + 1.0,
			Featured:      rand.Float32() > 0.8,
			Sticky:        rand.Float32() > 0.95,
			AllowComments: rand.Float32() > 0.1,
			MetaTitle:     fmt.Sprintf("文章%d的SEO标题", i+1),
			MetaDescription: fmt.Sprintf("文章%d的SEO描述", i+1),
			MetaKeywords:    fmt.Sprintf("关键词%d,技术,教程", i+1),
			AuthorID:       authorID,
			CategoryID:     &categoryID,
		}
		
		if status == "published" {
			publishedAt := time.Now().AddDate(0, 0, -rand.Intn(365))
			posts[i].PublishedAt = &publishedAt
		}
	}
	db.CreateInBatches(posts, 50)
	fmt.Printf("✓ 文章数据: %d条\n", len(posts))
	
	// 为文章分配标签
	for i := 1; i <= 200; i++ {
		tagCount := rand.Intn(4) + 1 // 每篇文章1-4个标签
		selectedTags := make([]Tag, 0, tagCount)
		usedTagIDs := make(map[uint]bool)
		
		for j := 0; j < tagCount; j++ {
			tagID := uint(rand.Intn(10) + 1)
			if !usedTagIDs[tagID] {
				var tag Tag
				db.First(&tag, tagID)
				selectedTags = append(selectedTags, tag)
				usedTagIDs[tagID] = true
			}
		}
		
		var post Post
		db.First(&post, uint(i))
		db.Model(&post).Association("Tags").Append(selectedTags)
	}
	fmt.Println("✓ 文章标签关联完成")
	
	// 生成评论数据
	comments := make([]Comment, 500)
	commentStatuses := []string{"approved", "pending", "spam"}
	
	for i := 0; i < 500; i++ {
		comments[i] = Comment{
			Content:   fmt.Sprintf("这是评论 %d 的内容，包含了用户的真实想法和建议。", i+1),
			Status:    commentStatuses[rand.Intn(len(commentStatuses))],
			Type:      "comment",
			LikeCount: rand.Intn(50),
			Level:     1,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			UserIP:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			IsSpam:    rand.Float32() > 0.9,
			PostID:    uint(rand.Intn(200) + 1),
			AuthorID:  uint(rand.Intn(50) + 1),
		}
	}
	db.CreateInBatches(comments, 100)
	fmt.Printf("✓ 评论数据: %d条\n", len(comments))
	
	// 生成回复评论
	replies := make([]Comment, 200)
	for i := 0; i < 200; i++ {
		parentID := uint(rand.Intn(500) + 1)
		replies[i] = Comment{
			Content:   fmt.Sprintf("这是对评论的回复 %d", i+1),
			Status:    "approved",
			Type:      "reply",
			LikeCount: rand.Intn(20),
			ParentID:  &parentID,
			Level:     2,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			UserIP:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			PostID:    uint(rand.Intn(200) + 1),
			AuthorID:  uint(rand.Intn(50) + 1),
		}
	}
	db.CreateInBatches(replies, 50)
	fmt.Printf("✓ 回复数据: %d条\n", len(replies))
	
	// 生成点赞数据
	likes := make([]Like, 1000)
	for i := 0; i < 1000; i++ {
		userID := uint(rand.Intn(50) + 1)
		
		if rand.Float32() > 0.3 { // 70%点赞文章
			postID := uint(rand.Intn(200) + 1)
			likes[i] = Like{
				UserID: userID,
				PostID: &postID,
				Type:   "like",
				UserIP: fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			}
		} else { // 30%点赞评论
			commentID := uint(rand.Intn(700) + 1)
			likes[i] = Like{
				UserID:    userID,
				CommentID: &commentID,
				Type:      "like",
				UserIP:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			}
		}
	}
	db.CreateInBatches(likes, 100)
	fmt.Printf("✓ 点赞数据: %d条\n", len(likes))
	
	// 生成关注关系
	follows := make([]Follow, 300)
	for i := 0; i < 300; i++ {
		followerID := uint(rand.Intn(50) + 1)
		followingID := uint(rand.Intn(50) + 1)
		
		// 确保不自己关注自己
		for followerID == followingID {
			followingID = uint(rand.Intn(50) + 1)
		}
		
		follows[i] = Follow{
			FollowerID:  followerID,
			FollowingID: followingID,
			Status:      "active",
		}
	}
	db.CreateInBatches(follows, 50)
	fmt.Printf("✓ 关注关系: %d条\n", len(follows))
	
	// 生成通知数据
	notifications := make([]Notification, 200)
	notificationTypes := []string{"comment", "like", "follow", "mention", "system"}
	
	for i := 0; i < 200; i++ {
		notificationType := notificationTypes[rand.Intn(len(notificationTypes))]
		notifications[i] = Notification{
			UserID:      uint(rand.Intn(50) + 1),
			Type:        notificationType,
			Title:       fmt.Sprintf("%s通知 %d", notificationType, i+1),
			Content:     fmt.Sprintf("这是一条%s类型的通知内容", notificationType),
			Data:        fmt.Sprintf(`{"id": %d, "type": "%s"}`, i+1, notificationType),
			IsRead:      rand.Float32() > 0.4,
			RelatedID:   func() *uint { id := uint(rand.Intn(100) + 1); return &id }(),
			RelatedType: []string{"post", "comment", "user"}[rand.Intn(3)],
		}
		
		if notifications[i].IsRead {
			readAt := time.Now().AddDate(0, 0, -rand.Intn(7))
			notifications[i].ReadAt = &readAt
		}
	}
	db.CreateInBatches(notifications, 50)
	fmt.Printf("✓ 通知数据: %d条\n", len(notifications))
	
	// 生成系统设置
	settings := []Setting{
		{Key: "site_title", Value: "我的博客系统", Type: "string", Description: "网站标题", Group: "general", IsPublic: true},
		{Key: "site_description", Value: "一个功能完整的博客系统", Type: "string", Description: "网站描述", Group: "general", IsPublic: true},
		{Key: "posts_per_page", Value: "10", Type: "integer", Description: "每页文章数量", Group: "display", IsPublic: true},
		{Key: "allow_registration", Value: "true", Type: "boolean", Description: "允许用户注册", Group: "user", IsPublic: false},
		{Key: "comment_moderation", Value: "true", Type: "boolean", Description: "评论需要审核", Group: "comment", IsPublic: false},
		{Key: "email_notifications", Value: "true", Type: "boolean", Description: "邮件通知", Group: "notification", IsPublic: false},
		{Key: "max_upload_size", Value: "10485760", Type: "integer", Description: "最大上传文件大小(字节)", Group: "upload", IsPublic: false},
		{Key: "allowed_file_types", Value: "jpg,jpeg,png,gif,pdf,doc,docx", Type: "string", Description: "允许的文件类型", Group: "upload", IsPublic: false},
	}
	db.Create(&settings)
	fmt.Printf("✓ 系统设置: %d条\n", len(settings))
	
	elapsed := time.Since(start)
	fmt.Printf("\n✓ 综合测试数据生成完成，耗时: %v\n", elapsed)
}

// 综合业务场景演示
func demonstrateComprehensiveScenarios(db *gorm.DB) {
	fmt.Println("\n=== 综合业务场景演示 ===")
	
	// 初始化服务
	userService := NewUserService(db)
	postService := NewPostService(db)
	commentService := NewCommentService(db)
	notificationService := NewNotificationService(db)
	analyticsService := NewAnalyticsService(db)
	
	// 场景1：用户注册和资料完善
	fmt.Println("\n--- 场景1：用户注册和资料完善 ---")
	newUser := &User{
		Username:     "newuser",
		Email:        "newuser@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "New",
		LastName:     "User",
		Bio:          "我是新用户",
		Location:     "北京",
		Status:       "active",
		Role:         "user",
	}
	
	if err := userService.CreateUser(newUser); err != nil {
		fmt.Printf("用户创建失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户创建成功，ID: %d\n", newUser.ID)
		
		// 完善用户资料
		profile := &UserProfile{
			Company:      "科技公司",
			JobTitle:     "软件工程师",
			Education:    "计算机科学学士",
			Skills:       "Go, Python, JavaScript",
			Experience:   3,
			SalaryRange:  "10k-15k",
			Languages:    "中文, 英文",
			Interests:    "编程, 阅读, 旅游",
			PrivacyLevel: "public",
		}
		
		if err := userService.UpdateUserProfile(newUser.ID, profile); err != nil {
			fmt.Printf("资料更新失败: %v\n", err)
		} else {
			fmt.Println("✓ 用户资料更新成功")
		}
	}
	
	// 场景2：发布文章和标签管理
	fmt.Println("\n--- 场景2：发布文章和标签管理 ---")
	newPost := &Post{
		Title:         "GORM实战教程 - 从入门到精通",
		Slug:          "gorm-tutorial-comprehensive",
		Content:       "这是一篇关于GORM的详细教程，涵盖了从基础使用到高级特性的所有内容...",
		Excerpt:       "GORM实战教程，带你深入了解Go语言最流行的ORM框架",
		FeaturedImage: "/images/gorm-tutorial.jpg",
		Status:        "published",
		Type:          "post",
		Format:        "standard",
		Featured:      true,
		AllowComments: true,
		MetaTitle:     "GORM实战教程 - 完整指南",
		MetaDescription: "学习GORM的最佳实践和高级技巧",
		MetaKeywords:  "GORM,Go,ORM,数据库,教程",
		AuthorID:      1,
		CategoryID:    func() *uint { id := uint(1); return &id }(),
	}
	
	// 为文章分配标签
	var tags []Tag
	db.Where("slug IN ?", []string{"go", "database", "tutorial"}).Find(&tags)
	newPost.Tags = tags
	
	if err := postService.CreatePost(newPost); err != nil {
		fmt.Printf("文章创建失败: %v\n", err)
	} else {
		fmt.Printf("✓ 文章创建成功，ID: %d\n", newPost.ID)
	}
	
	// 场景3：用户互动（评论、点赞、关注）
	fmt.Println("\n--- 场景3：用户互动 ---")
	
	// 添加评论
	newComment := &Comment{
		Content:   "这篇文章写得非常好，对我帮助很大！",
		Status:    "approved",
		Type:      "comment",
		Level:     1,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		UserIP:    "192.168.1.100",
		PostID:    1,
		AuthorID:  2,
	}
	
	if err := commentService.CreateComment(newComment); err != nil {
		fmt.Printf("评论创建失败: %v\n", err)
	} else {
		fmt.Printf("✓ 评论创建成功，ID: %d\n", newComment.ID)
	}
	
	// 点赞文章
	if err := postService.LikePost(2, 1); err != nil {
		fmt.Printf("点赞失败: %v\n", err)
	} else {
		fmt.Println("✓ 文章点赞成功")
	}
	
	// 关注用户
	if err := userService.FollowUser(2, 1); err != nil {
		fmt.Printf("关注失败: %v\n", err)
	} else {
		fmt.Println("✓ 用户关注成功")
	}
	
	// 场景4：内容搜索和分页
	fmt.Println("\n--- 场景4：内容搜索和分页 ---")
	
	// 按分类获取文章
	posts, total, err := postService.GetPostsByCategory("tech", 1, 5)
	if err != nil {
		fmt.Printf("分类查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 技术分类文章: %d/%d篇\n", len(posts), total)
	}
	
	// 搜索文章
	searchPosts, searchTotal, err := postService.SearchPosts("教程", 1, 5)
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)
	} else {
		fmt.Printf("✓ 搜索结果: %d/%d篇\n", len(searchPosts), searchTotal)
	}
	
	// 获取用户关注者
	followers, followerTotal, err := userService.GetUserFollowers(1, 1, 10)
	if err != nil {
		fmt.Printf("关注者查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户关注者: %d/%d人\n", len(followers), followerTotal)
	}
	
	// 场景5：通知管理
	fmt.Println("\n--- 场景5：通知管理 ---")
	
	// 获取用户通知
	notifications, notificationTotal, err := notificationService.GetUserNotifications(1, 1, 5)
	if err != nil {
		fmt.Printf("通知查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户通知: %d/%d条\n", len(notifications), notificationTotal)
		
		// 标记通知为已读
		if len(notifications) > 0 {
			if err := notificationService.MarkAsRead(notifications[0].ID); err != nil {
				fmt.Printf("标记已读失败: %v\n", err)
			} else {
				fmt.Println("✓ 通知标记为已读")
			}
		}
	}
	
	// 获取未读通知数量
	unreadCount, err := notificationService.GetUnreadCount(1)
	if err != nil {
		fmt.Printf("未读数量查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 未读通知数量: %d条\n", unreadCount)
	}
	
	// 场景6：数据统计和分析
	fmt.Println("\n--- 场景6：数据统计和分析 ---")
	
	// 获取仪表板统计
	stats, err := analyticsService.GetDashboardStats()
	if err != nil {
		fmt.Printf("统计查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 系统统计:\n")
		fmt.Printf("  - 总用户数: %d (活跃: %d)\n", stats.TotalUsers, stats.ActiveUsers)
		fmt.Printf("  - 总文章数: %d (已发布: %d)\n", stats.TotalPosts, stats.PublishedPosts)
		fmt.Printf("  - 总评论数: %d (已审核: %d)\n", stats.TotalComments, stats.ApprovedComments)
		fmt.Printf("  - 总浏览量: %d\n", stats.TotalViews)
		fmt.Printf("  - 总点赞数: %d\n", stats.TotalLikes)
	}
	
	// 获取热门文章
	popularPosts, err := analyticsService.GetPopularPosts(5)
	if err != nil {
		fmt.Printf("热门文章查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 热门文章 (前5篇):\n")
		for i, post := range popularPosts {
			fmt.Printf("  %d. %s (浏览: %d, 点赞: %d)\n", i+1, post.Title, post.ViewCount, post.LikeCount)
		}
	}
	
	// 获取活跃用户
	activeUsers, err := analyticsService.GetActiveUsers(5)
	if err != nil {
		fmt.Printf("活跃用户查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 活跃用户 (前5名):\n")
		for i, user := range activeUsers {
			fmt.Printf("  %d. %s (文章: %d, 评论: %d, 点赞: %d)\n", 
				i+1, user.Username, user.PostCount, user.CommentCount, user.LikeCount)
		}
	}
}

// 高级查询演示
func demonstrateAdvancedQueries(db *gorm.DB) {
	fmt.Println("\n=== 高级查询演示 ===")
	
	// 复杂的多表连接查询
	fmt.Println("\n--- 复杂多表连接查询 ---")
	type PostWithStats struct {
		ID           uint   `json:"id"`
		Title        string `json:"title"`
		AuthorName   string `json:"author_name"`
		CategoryName string `json:"category_name"`
		ViewCount    int    `json:"view_count"`
		LikeCount    int    `json:"like_count"`
		CommentCount int    `json:"comment_count"`
		TagNames     string `json:"tag_names"`
	}
	
	var postsWithStats []PostWithStats
	err := db.Table("posts p").
		Select(`p.id, p.title, u.username as author_name, c.name as category_name, 
			p.view_count, p.like_count, p.comment_count,
			GROUP_CONCAT(t.name) as tag_names`).
		Joins("JOIN users u ON p.author_id = u.id").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Joins("LEFT JOIN post_tags pt ON p.id = pt.post_id").
		Joins("LEFT JOIN tags t ON pt.tag_id = t.id").
		Where("p.status = ?", "published").
		Group("p.id, p.title, u.username, c.name, p.view_count, p.like_count, p.comment_count").
		Order("p.view_count DESC").
		Limit(5).Scan(&postsWithStats).Error
	
	if err != nil {
		fmt.Printf("复杂查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 文章详细统计 (前5篇):\n")
		for i, post := range postsWithStats {
			fmt.Printf("  %d. %s\n", i+1, post.Title)
			fmt.Printf("     作者: %s | 分类: %s\n", post.AuthorName, post.CategoryName)
			fmt.Printf("     浏览: %d | 点赞: %d | 评论: %d\n", post.ViewCount, post.LikeCount, post.CommentCount)
			fmt.Printf("     标签: %s\n", post.TagNames)
			fmt.Println()
		}
	}
	
	// 窗口函数查询（排名）
	fmt.Println("--- 窗口函数查询 ---")
	type UserRanking struct {
		Username     string `json:"username"`
		PostCount    int    `json:"post_count"`
		CommentCount int    `json:"comment_count"`
		TotalScore   int    `json:"total_score"`
		Rank         int    `json:"rank"`
	}
	
	var userRankings []UserRanking
	err = db.Raw(`
		SELECT username, post_count, comment_count, 
		       (post_count * 3 + comment_count) as total_score,
		       ROW_NUMBER() OVER (ORDER BY (post_count * 3 + comment_count) DESC) as rank
		FROM users 
		WHERE status = 'active' 
		ORDER BY total_score DESC 
		LIMIT 10
	`).Scan(&userRankings).Error
	
	if err != nil {
		fmt.Printf("排名查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 用户活跃度排名 (前10名):\n")
		for _, user := range userRankings {
			fmt.Printf("  第%d名: %s (文章: %d, 评论: %d, 总分: %d)\n", 
				user.Rank, user.Username, user.PostCount, user.CommentCount, user.TotalScore)
		}
	}
	
	// 时间序列分析
	fmt.Println("\n--- 时间序列分析 ---")
	type MonthlyActivity struct {
		Month      string `json:"month"`
		PostCount  int    `json:"post_count"`
		UserCount  int    `json:"user_count"`
		CommentCount int  `json:"comment_count"`
	}
	
	var monthlyActivity []MonthlyActivity
	err = db.Raw(`
		SELECT 
		    strftime('%Y-%m', created_at) as month,
		    COUNT(CASE WHEN 'posts' THEN 1 END) as post_count,
		    COUNT(CASE WHEN 'users' THEN 1 END) as user_count,
		    COUNT(CASE WHEN 'comments' THEN 1 END) as comment_count
		FROM (
		    SELECT created_at, 'posts' as type FROM posts WHERE status = 'published'
		    UNION ALL
		    SELECT created_at, 'users' as type FROM users WHERE status = 'active'
		    UNION ALL
		    SELECT created_at, 'comments' as type FROM comments WHERE status = 'approved'
		) combined
		WHERE created_at >= date('now', '-6 months')
		GROUP BY month
		ORDER BY month DESC
		LIMIT 6
	`).Scan(&monthlyActivity).Error
	
	if err != nil {
		fmt.Printf("时间序列查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 近6个月活动统计:\n")
		for _, activity := range monthlyActivity {
			fmt.Printf("  %s: 文章 %d篇, 用户 %d人, 评论 %d条\n", 
				activity.Month, activity.PostCount, activity.UserCount, activity.CommentCount)
		}
	}
}

// 性能测试
func performanceTest(db *gorm.DB) {
	fmt.Println("\n=== 性能测试 ===")
	
	// 测试批量查询性能
	fmt.Println("\n--- 批量查询性能测试 ---")
	start := time.Now()
	var posts []Post
	db.Preload("Author").Preload("Category").Preload("Tags").Limit(100).Find(&posts)
	fmt.Printf("✓ 预加载查询100篇文章: %v\n", time.Since(start))
	
	// 测试分页查询性能
	start = time.Now()
	var paginatedPosts []Post
	var total int64
	db.Model(&Post{}).Where("status = ?", "published").Count(&total)
	db.Where("status = ?", "published").Offset(50).Limit(20).Find(&paginatedPosts)
	fmt.Printf("✓ 分页查询(50-70): %v\n", time.Since(start))
	
	// 测试复杂聚合查询性能
	start = time.Now()
	type CategoryStats struct {
		CategoryName string `json:"category_name"`
		PostCount    int64  `json:"post_count"`
		AvgViews     float64 `json:"avg_views"`
		TotalLikes   int64  `json:"total_likes"`
	}
	
	var categoryStats []CategoryStats
	db.Table("posts p").
		Select("c.name as category_name, COUNT(p.id) as post_count, AVG(p.view_count) as avg_views, SUM(p.like_count) as total_likes").
		Joins("JOIN categories c ON p.category_id = c.id").
		Where("p.status = ?", "published").
		Group("c.id, c.name").
		Order("post_count DESC").
		Scan(&categoryStats)
	fmt.Printf("✓ 分类统计查询: %v\n", time.Since(start))
	
	// 测试事务性能
	start = time.Now()
	for i := 0; i < 10; i++ {
		db.Transaction(func(tx *gorm.DB) error {
			var user User
			tx.First(&user, 1)
			tx.Model(&user).Update("login_count", gorm.Expr("login_count + ?", 1))
			return nil
		})
	}
	fmt.Printf("✓ 10个事务操作: %v\n", time.Since(start))
}

// 主函数
func main() {
	fmt.Println("=== GORM Level 6 综合实战练习 ===")

	// 初始化数据库
	db := initDB()
	fmt.Println("✓ 数据库初始化完成")

	// 生成测试数据
	generateComprehensiveTestData(db)

	// 演示综合业务场景
	demonstrateComprehensiveScenarios(db)

	// 演示高级查询
	demonstrateAdvancedQueries(db)

	// 性能测试
	performanceTest(db)

	fmt.Println("\n=== Level 6 综合实战练习完成 ===")
	fmt.Println("\n🎉 恭喜！您已经完成了GORM的全部练习，现在您应该能够：")
	fmt.Println("1. 熟练使用GORM进行数据库操作")
	fmt.Println("2. 设计复杂的数据模型和关联关系")
	fmt.Println("3. 实现高效的查询和事务处理")
	fmt.Println("4. 优化数据库性能和索引")
	fmt.Println("5. 构建完整的业务应用系统")
	fmt.Println("\n继续学习，成为GORM专家！💪")
}
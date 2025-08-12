// 04_unit_exercises/level2_associations.go - Level 2 关联关系练习
// 对应文档：03_GORM单元练习_基础技能训练.md

package main

import (
	"fmt"
	"log"
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

// 练习1：一对一关系

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	
	// 一对一关系：用户资料
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile,omitempty"`
}

// Profile 用户资料模型
type Profile struct {
	BaseModel
	UserID    uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName string `gorm:"size:50" json:"first_name"`
	LastName  string `gorm:"size:50" json:"last_name"`
	Bio       string `gorm:"type:text" json:"bio"`
	Avatar    string `gorm:"size:255" json:"avatar"`
	Phone     string `gorm:"size:20" json:"phone"`
	Address   string `gorm:"size:255" json:"address"`
	BirthDate *time.Time `json:"birth_date"`
	
	// 反向关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// 练习2：一对多关系

// Category 分类模型
type Category struct {
	BaseModel
	Name        string `gorm:"size:100;not null;index" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// 一对多关系：分类下的文章
	Posts []Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	BaseModel
	Title      string `gorm:"size:200;not null;index" json:"title"`
	Slug       string `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Excerpt    string `gorm:"size:500" json:"excerpt"`
	Status     string `gorm:"size:20;default:'draft';index" json:"status"`
	ViewCount  int    `gorm:"default:0" json:"view_count"`
	LikeCount  int    `gorm:"default:0" json:"like_count"`
	PublishedAt *time.Time `gorm:"index" json:"published_at"`
	
	// 外键关系
	AuthorID   uint `gorm:"not null;index" json:"author_id"`
	CategoryID *uint `gorm:"index" json:"category_id"`
	
	// 关联关系
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	
	// 一对多关系：文章的评论
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments,omitempty"`
	
	// 多对多关系：文章的标签
	Tags []Tag `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tags,omitempty"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	Content   string `gorm:"type:text;not null" json:"content"`
	Status    string `gorm:"size:20;default:'pending';index" json:"status"`
	IPAddress string `gorm:"size:45" json:"ip_address"`
	
	// 外键关系
	PostID   uint `gorm:"not null;index" json:"post_id"`
	AuthorID uint `gorm:"not null;index" json:"author_id"`
	ParentID *uint `gorm:"index" json:"parent_id"`
	
	// 关联关系
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	
	// 自关联：评论回复
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// 练习3：多对多关系

// Tag 标签模型
type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"size:7;default:'#007bff'" json:"color"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// 多对多关系：标签的文章
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// PostTag 文章标签中间表（自定义）
type PostTag struct {
	PostID    uint      `gorm:"primaryKey" json:"post_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// 关联关系
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Tag  Tag  `gorm:"foreignKey:TagID" json:"tag,omitempty"`
}

// 数据库初始化
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("level2_associations.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: gorm.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Profile{}, &Category{}, &Post{}, &Comment{}, &Tag{}, &PostTag{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db
}

// 练习1：一对一关系操作

// CreateUserWithProfile 创建用户及其资料
func CreateUserWithProfile(db *gorm.DB, username, email, password, firstName, lastName, bio string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
		Profile: Profile{
			FirstName: firstName,
			LastName:  lastName,
			Bio:       bio,
		},
	}

	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("创建用户及资料成功，用户ID: %d, 资料ID: %d\n", user.ID, user.Profile.ID)
	return user, nil
}

// GetUserWithProfile 获取用户及其资料
func GetUserWithProfile(db *gorm.DB, userID uint) (*User, error) {
	var user User
	result := db.Preload("Profile").First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(db *gorm.DB, userID uint, updates map[string]interface{}) error {
	result := db.Model(&Profile{}).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("更新用户资料成功，影响行数: %d\n", result.RowsAffected)
	return nil
}

// 练习2：一对多关系操作

// CreateCategoryWithPosts 创建分类及其文章
func CreateCategoryWithPosts(db *gorm.DB, categoryName, categorySlug string, authorID uint, postTitles []string) (*Category, error) {
	category := &Category{
		Name: categoryName,
		Slug: categorySlug,
	}

	// 创建文章
	for i, title := range postTitles {
		post := Post{
			Title:    title,
			Slug:     fmt.Sprintf("%s-%d", categorySlug, i+1),
			Content:  fmt.Sprintf("这是%s的内容", title),
			Excerpt:  fmt.Sprintf("这是%s的摘要", title),
			Status:   "published",
			AuthorID: authorID,
		}
		category.Posts = append(category.Posts, post)
	}

	result := db.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("创建分类及文章成功，分类ID: %d, 文章数量: %d\n", category.ID, len(category.Posts))
	return category, nil
}

// GetCategoryWithPosts 获取分类及其文章
func GetCategoryWithPosts(db *gorm.DB, categoryID uint) (*Category, error) {
	var category Category
	result := db.Preload("Posts").Preload("Posts.Author").First(&category, categoryID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

// GetPostsWithAuthorAndCategory 获取文章及其作者和分类
func GetPostsWithAuthorAndCategory(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post
	result := db.Preload("Author").Preload("Category").Limit(limit).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// CreatePostWithComments 创建文章及其评论
func CreatePostWithComments(db *gorm.DB, title, content string, authorID, categoryID uint, commentContents []string) (*Post, error) {
	post := &Post{
		Title:      title,
		Slug:       fmt.Sprintf("post-%d", time.Now().Unix()),
		Content:    content,
		Excerpt:    content[:min(len(content), 100)],
		Status:     "published",
		AuthorID:   authorID,
		CategoryID: &categoryID,
	}

	// 创建评论
	for _, commentContent := range commentContents {
		comment := Comment{
			Content:  commentContent,
			Status:   "approved",
			AuthorID: authorID,
		}
		post.Comments = append(post.Comments, comment)
	}

	result := db.Create(post)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("创建文章及评论成功，文章ID: %d, 评论数量: %d\n", post.ID, len(post.Comments))
	return post, nil
}

// GetPostWithComments 获取文章及其评论
func GetPostWithComments(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	result := db.Preload("Comments").Preload("Comments.Author").First(&post, postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// 练习3：多对多关系操作

// CreateTagsAndAssignToPosts 创建标签并分配给文章
func CreateTagsAndAssignToPosts(db *gorm.DB, tagNames []string, postIDs []uint) error {
	// 创建标签
	var tags []Tag
	for _, name := range tagNames {
		tag := Tag{
			Name: name,
			Slug: name,
		}
		tags = append(tags, tag)
	}

	result := db.Create(&tags)
	if result.Error != nil {
		return result.Error
	}

	// 为文章分配标签
	for _, postID := range postIDs {
		var post Post
		if err := db.First(&post, postID).Error; err != nil {
			continue
		}

		// 关联标签
		if err := db.Model(&post).Association("Tags").Append(&tags); err != nil {
			return err
		}
	}

	fmt.Printf("创建标签并分配成功，标签数量: %d\n", len(tags))
	return nil
}

// GetPostWithTags 获取文章及其标签
func GetPostWithTags(db *gorm.DB, postID uint) (*Post, error) {
	var post Post
	result := db.Preload("Tags").First(&post, postID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// GetTagWithPosts 获取标签及其文章
func GetTagWithPosts(db *gorm.DB, tagID uint) (*Tag, error) {
	var tag Tag
	result := db.Preload("Posts").Preload("Posts.Author").First(&tag, tagID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tag, nil
}

// AddTagsToPost 为文章添加标签
func AddTagsToPost(db *gorm.DB, postID uint, tagIDs []uint) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []Tag
	if err := db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return db.Model(&post).Association("Tags").Append(&tags)
}

// RemoveTagsFromPost 从文章中移除标签
func RemoveTagsFromPost(db *gorm.DB, postID uint, tagIDs []uint) error {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return err
	}

	var tags []Tag
	if err := db.Find(&tags, tagIDs).Error; err != nil {
		return err
	}

	return db.Model(&post).Association("Tags").Delete(&tags)
}

// 练习4：复杂查询

// GetPostsWithAllAssociations 获取文章及所有关联数据
func GetPostsWithAllAssociations(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post
	result := db.Preload("Author").
		Preload("Author.Profile").
		Preload("Category").
		Preload("Comments").
		Preload("Comments.Author").
		Preload("Tags").
		Limit(limit).Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// GetUserPostsWithStats 获取用户文章及统计信息
func GetUserPostsWithStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
	var posts []Post
	result := db.Where("author_id = ?", userID).Preload("Category").Preload("Tags").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	// 统计信息
	stats := map[string]interface{}{
		"total_posts":    len(posts),
		"total_views":    0,
		"total_likes":    0,
		"categories":     make(map[string]int),
		"tags":           make(map[string]int),
		"posts":          posts,
	}

	categoryCount := make(map[string]int)
	tagCount := make(map[string]int)
	totalViews := 0
	totalLikes := 0

	for _, post := range posts {
		totalViews += post.ViewCount
		totalLikes += post.LikeCount

		if post.Category != nil {
			categoryCount[post.Category.Name]++
		}

		for _, tag := range post.Tags {
			tagCount[tag.Name]++
		}
	}

	stats["total_views"] = totalViews
	stats["total_likes"] = totalLikes
	stats["categories"] = categoryCount
	stats["tags"] = tagCount

	return stats, nil
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 主函数演示
func main() {
	fmt.Println("=== GORM Level 2 关联关系练习 ===")

	// 初始化数据库
	db := initDB()
	fmt.Println("✓ 数据库初始化完成")

	// 练习1：一对一关系
	fmt.Println("\n=== 一对一关系练习 ===")

	// 创建用户及资料
	user1, err := CreateUserWithProfile(db, "alice", "alice@example.com", "password123", "Alice", "Smith", "我是Alice，一名开发者")
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	user2, err := CreateUserWithProfile(db, "bob", "bob@example.com", "password456", "Bob", "Johnson", "我是Bob，喜欢写作")
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// 获取用户及资料
	if user1 != nil {
		fetchedUser, err := GetUserWithProfile(db, user1.ID)
		if err != nil {
			fmt.Printf("获取用户失败: %v\n", err)
		} else {
			fmt.Printf("用户: %s, 全名: %s %s\n", fetchedUser.Username, fetchedUser.Profile.FirstName, fetchedUser.Profile.LastName)
		}
	}

	// 更新用户资料
	if user1 != nil {
		updates := map[string]interface{}{
			"bio":   "我是Alice，一名全栈开发者",
			"phone": "123-456-7890",
		}
		if err := UpdateUserProfile(db, user1.ID, updates); err != nil {
			fmt.Printf("更新用户资料失败: %v\n", err)
		}
	}

	// 练习2：一对多关系
	fmt.Println("\n=== 一对多关系练习 ===")

	// 创建分类及文章
	if user1 != nil {
		category1, err := CreateCategoryWithPosts(db, "技术", "tech", user1.ID, []string{"Go语言入门", "GORM使用指南", "数据库设计原则"})
		if err != nil {
			fmt.Printf("创建分类失败: %v\n", err)
		}

		// 获取分类及文章
		if category1 != nil {
			fetchedCategory, err := GetCategoryWithPosts(db, category1.ID)
			if err != nil {
				fmt.Printf("获取分类失败: %v\n", err)
			} else {
				fmt.Printf("分类: %s, 文章数量: %d\n", fetchedCategory.Name, len(fetchedCategory.Posts))
				for _, post := range fetchedCategory.Posts {
					fmt.Printf("  - %s (作者: %s)\n", post.Title, post.Author.Username)
				}
			}
		}
	}

	// 创建文章及评论
	if user1 != nil && user2 != nil {
		post, err := CreatePostWithComments(db, "GORM高级用法", "这是一篇关于GORM高级用法的文章...", user1.ID, 1, []string{"很有用的文章！", "学到了很多", "期待更多内容"})
		if err != nil {
			fmt.Printf("创建文章失败: %v\n", err)
		} else {
			// 获取文章及评论
			fetchedPost, err := GetPostWithComments(db, post.ID)
			if err != nil {
				fmt.Printf("获取文章失败: %v\n", err)
			} else {
				fmt.Printf("文章: %s, 评论数量: %d\n", fetchedPost.Title, len(fetchedPost.Comments))
			}
		}
	}

	// 练习3：多对多关系
	fmt.Println("\n=== 多对多关系练习 ===")

	// 创建标签并分配给文章
	tagNames := []string{"Go", "数据库", "后端", "教程"}
	postIDs := []uint{1, 2, 3, 4}
	if err := CreateTagsAndAssignToPosts(db, tagNames, postIDs); err != nil {
		fmt.Printf("创建标签失败: %v\n", err)
	}

	// 获取文章及标签
	post, err := GetPostWithTags(db, 1)
	if err != nil {
		fmt.Printf("获取文章标签失败: %v\n", err)
	} else {
		fmt.Printf("文章: %s, 标签: ", post.Title)
		for i, tag := range post.Tags {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(tag.Name)
		}
		fmt.Println()
	}

	// 获取标签及文章
	tag, err := GetTagWithPosts(db, 1)
	if err != nil {
		fmt.Printf("获取标签文章失败: %v\n", err)
	} else {
		fmt.Printf("标签: %s, 文章数量: %d\n", tag.Name, len(tag.Posts))
	}

	// 练习4：复杂查询
	fmt.Println("\n=== 复杂查询练习 ===")

	// 获取文章及所有关联数据
	posts, err := GetPostsWithAllAssociations(db, 5)
	if err != nil {
		fmt.Printf("获取文章失败: %v\n", err)
	} else {
		fmt.Printf("获取到 %d 篇文章（包含所有关联数据）\n", len(posts))
	}

	// 获取用户文章统计
	if user1 != nil {
		stats, err := GetUserPostsWithStats(db, user1.ID)
		if err != nil {
			fmt.Printf("获取用户统计失败: %v\n", err)
		} else {
			fmt.Printf("用户 %s 的文章统计:\n", user1.Username)
			fmt.Printf("  总文章数: %v\n", stats["total_posts"])
			fmt.Printf("  总浏览数: %v\n", stats["total_views"])
			fmt.Printf("  总点赞数: %v\n", stats["total_likes"])
			fmt.Printf("  分类分布: %v\n", stats["categories"])
			fmt.Printf("  标签分布: %v\n", stats["tags"])
		}
	}

	fmt.Println("\n=== Level 2 关联关系练习完成 ===")
}
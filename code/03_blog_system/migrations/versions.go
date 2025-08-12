// migrations/versions.go - 具体的迁移版本定义
// 每个版本对应一个具体的数据库变更

package migrations

import (
	"blog-system/models"

	"gorm.io/gorm"
)

// GetAllMigrations 获取所有迁移版本
func GetAllMigrations() []MigrationItem {
	return []MigrationItem{
		{
			Version: "001_initial_schema",
			Name:    "创建初始数据库结构",
			Up:      migration001Up,
			Down:    migration001Down,
		},
		{
			Version: "002_create_indexes",
			Name:    "创建数据库索引",
			Up:      migration002Up,
			Down:    migration002Down,
		},
		{
			Version: "003_fix_constraints",
			Name:    "修复外键约束问题",
			Up:      migration003Up,
			Down:    migration003Down,
		},
	}
}

// migration001Up 创建初始数据库结构
func migration001Up(db *gorm.DB) error {
	// 按照依赖关系顺序创建表
	return db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	)
}

// migration001Down 回滚初始数据库结构
func migration001Down(db *gorm.DB) error {
	// 按照相反的依赖关系顺序删除表
	tables := []string{"Like", "Comment", "Post", "post_tags", "Tag", "Category", "Profile", "User"}
	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			return err
		}
	}
	return nil
}

// migration002Up 创建数据库索引
func migration002Up(db *gorm.DB) error {
	// 为Like表创建复合唯一索引
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'Like' AND index_name = 'idx_likes_user_target'").Scan(&count)
	if count == 0 {
		if err := db.Exec("CREATE UNIQUE INDEX idx_likes_user_target ON `Like`(user_id, target_id, target_type)").Error; err != nil {
			return err
		}
	}

	// 为Post表创建复合索引
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'Post' AND index_name = 'idx_posts_status_published'").Scan(&count)
	if count == 0 {
		if err := db.Exec("CREATE INDEX idx_posts_status_published ON `Post`(status, published_at)").Error; err != nil {
			return err
		}
	}

	// 为Comment表创建复合索引
	db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'Comment' AND index_name = 'idx_comments_post_status'").Scan(&count)
	if count == 0 {
		if err := db.Exec("CREATE INDEX idx_comments_post_status ON `Comment`(post_id, status)").Error; err != nil {
			return err
		}
	}

	return nil
}

// migration002Down 删除数据库索引
func migration002Down(db *gorm.DB) error {
	indexes := []struct {
		table string
		index string
	}{
		{"Like", "idx_likes_user_target"},
		{"Post", "idx_posts_status_published"},
		{"Comment", "idx_comments_post_status"},
	}

	for _, idx := range indexes {
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", idx.table, idx.index).Scan(&count)
		if count > 0 {
			if err := db.Exec("DROP INDEX " + idx.index + " ON " + idx.table).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// migration003Up 修复外键约束问题
func migration003Up(db *gorm.DB) error {
	// 检查并修复可能存在的约束问题
	// 这个迁移主要是为了处理之前可能存在的约束冲突

	// 检查是否存在问题约束
	var constraintCount int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.table_constraints 
		WHERE constraint_schema = DATABASE() 
		AND constraint_name LIKE 'uniq_User_%'
	`).Scan(&constraintCount)

	// 如果存在问题约束，尝试删除
	if constraintCount > 0 {
		// 获取所有问题约束名称
		var constraints []string
		db.Raw(`
			SELECT constraint_name 
			FROM information_schema.table_constraints 
			WHERE constraint_schema = DATABASE() 
			AND constraint_name LIKE 'uniq_User_%'
		`).Scan(&constraints)

		for _, constraint := range constraints {
			// 安全地删除约束
			db.Exec("ALTER TABLE User DROP INDEX IF EXISTS " + constraint)
		}
	}

	return nil
}

// migration003Down 回滚外键约束修复
func migration003Down(db *gorm.DB) error {
	// 这个迁移的回滚不需要做任何操作
	// 因为我们只是清理了问题约束
	return nil
}

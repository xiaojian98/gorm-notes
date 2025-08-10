// migrations/migration.go - 数据库迁移版本控制
// 解决每次重启都要清理数据库的问题

package migrations

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Migration 迁移记录结构
type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Version   string    `gorm:"size:50;unique;not null"`
	Name      string    `gorm:"size:255;not null"`
	Executed  bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MigrationFunc 迁移函数类型
type MigrationFunc func(*gorm.DB) error

// MigrationItem 迁移项
type MigrationItem struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	db         *gorm.DB
	migrations []MigrationItem
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: make([]MigrationItem, 0),
	}
}

// AddMigration 添加迁移
func (m *MigrationManager) AddMigration(version, name string, up, down MigrationFunc) {
	m.migrations = append(m.migrations, MigrationItem{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

// InitMigrationTable 初始化迁移表
func (m *MigrationManager) InitMigrationTable() error {
	return m.db.AutoMigrate(&Migration{})
}

// RunMigrations 执行迁移
func (m *MigrationManager) RunMigrations() error {
	// 初始化迁移表
	if err := m.InitMigrationTable(); err != nil {
		return fmt.Errorf("初始化迁移表失败: %v", err)
	}

	for _, migration := range m.migrations {
		// 检查迁移是否已执行
		var existingMigration Migration
		err := m.db.Where("version = ?", migration.Version).First(&existingMigration).Error
		
		if err == gorm.ErrRecordNotFound {
			// 迁移未执行，开始执行
			log.Printf("执行迁移: %s - %s", migration.Version, migration.Name)
			
			// 开始事务
			tx := m.db.Begin()
			if tx.Error != nil {
				return fmt.Errorf("开始事务失败: %v", tx.Error)
			}

			// 执行迁移
			if err := migration.Up(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("执行迁移 %s 失败: %v", migration.Version, err)
			}

			// 记录迁移
			migrationRecord := Migration{
				Version:  migration.Version,
				Name:     migration.Name,
				Executed: true,
			}
			if err := tx.Create(&migrationRecord).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("记录迁移失败: %v", err)
			}

			// 提交事务
			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("提交事务失败: %v", err)
			}

			log.Printf("✅ 迁移 %s 执行成功", migration.Version)
		} else if err != nil {
			return fmt.Errorf("检查迁移状态失败: %v", err)
		} else {
			log.Printf("⏭️ 迁移 %s 已执行，跳过", migration.Version)
		}
	}

	return nil
}

// RollbackMigration 回滚迁移
func (m *MigrationManager) RollbackMigration(version string) error {
	// 查找迁移
	var migration *MigrationItem
	for _, m := range m.migrations {
		if m.Version == version {
			migration = &m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("未找到版本 %s 的迁移", version)
	}

	// 检查迁移是否已执行
	var existingMigration Migration
	err := m.db.Where("version = ? AND executed = ?", version, true).First(&existingMigration).Error
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("迁移 %s 未执行或已回滚", version)
	} else if err != nil {
		return fmt.Errorf("检查迁移状态失败: %v", err)
	}

	log.Printf("回滚迁移: %s - %s", migration.Version, migration.Name)

	// 开始事务
	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %v", tx.Error)
	}

	// 执行回滚
	if migration.Down != nil {
		if err := migration.Down(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("回滚迁移 %s 失败: %v", version, err)
		}
	}

	// 更新迁移记录
	if err := tx.Model(&existingMigration).Update("executed", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新迁移记录失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	log.Printf("✅ 迁移 %s 回滚成功", version)
	return nil
}

// GetMigrationStatus 获取迁移状态
func (m *MigrationManager) GetMigrationStatus() ([]Migration, error) {
	var migrations []Migration
	err := m.db.Order("version").Find(&migrations).Error
	return migrations, err
}

// RunMigrations 运行所有迁移的便捷函数
func RunMigrations(db *gorm.DB) error {
	manager := NewMigrationManager(db)
	
	// 添加所有迁移
	for _, migration := range GetAllMigrations() {
		manager.AddMigration(migration.Version, migration.Name, migration.Up, migration.Down)
	}
	
	// 执行迁移
	return manager.RunMigrations()
}
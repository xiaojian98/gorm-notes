// 04_unit_exercises/level4_transactions_hooks.go - Level 4 事务和钩子练习
// 对应文档：03_GORM单元练习_基础技能训练.md
// 本文件实现了GORM的事务管理和钩子函数练习，包括：
// 1. 数据库事务的使用和管理
// 2. GORM钩子函数的实现和应用
// 3. 银行账户系统的业务逻辑实现
// 4. 支持SQLite和MySQL两种数据库类型

package main

import (
	"errors" // 错误处理
	"fmt"    // 格式化输出
	"log"    // 日志记录
	"time"   // 时间处理

	"gorm.io/driver/mysql"  // MySQL数据库驱动
	"gorm.io/driver/sqlite" // SQLite数据库驱动
	"gorm.io/gorm"          // GORM核心库
	"gorm.io/gorm/logger"   // GORM日志组件
	"gorm.io/gorm/schema"   // GORM模式配置
)

// 数据库配置相关定义

// DatabaseType 数据库类型枚举
// 定义支持的数据库类型，目前支持SQLite和MySQL
type DatabaseType string

const (
	SQLite DatabaseType = "sqlite" // SQLite数据库类型
	MySQL  DatabaseType = "mysql"  // MySQL数据库类型
)

// DatabaseConfig 数据库配置结构体
// 包含数据库连接和连接池的所有配置参数
type DatabaseConfig struct {
	Type         DatabaseType    // 数据库类型(sqlite/mysql)
	DSN          string          // 数据源名称,用于指定数据库连接字符串
	MaxOpenConns int             // 最大打开连接数
	MaxIdleConns int             // 最大空闲连接数
	MaxLifetime  time.Duration   // 连接最大生命周期
	LogLevel     logger.LogLevel // 日志级别
}

// GetDefaultConfig 获取SQLite默认配置
// 返回一个包含默认参数的SQLite数据库配置对象
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:         SQLite,
		DSN:          "level4_transactions_hooks.db", // SQLite数据库文件名
		MaxOpenConns: 10,                             // 最大连接数10
		MaxIdleConns: 5,                              // 最大空闲连接5
		MaxLifetime:  time.Hour,                      // 连接生命周期1小时
		LogLevel:     logger.Info,                    // 日志级别为Info
	}
}

// GetMySQLConfig 获取MySQL配置
// 参数dsn: MySQL数据库连接字符串
// 返回一个包含默认参数的MySQL数据库配置对象
func GetMySQLConfig(dsn string) *DatabaseConfig {
	return &DatabaseConfig{
		Type:         MySQL,
		DSN:          dsn,
		MaxOpenConns: 20,        // MySQL建议更高的连接数
		MaxIdleConns: 10,        // 更多的空闲连接
		MaxLifetime:  time.Hour, // 连接生命周期1小时
		LogLevel:     logger.Info,
	}
}

// 基础模型定义

// BaseModel 基础模型结构体
// 包含所有数据库表通用的字段，采用GORM的软删除机制
// 所有业务模型都应该嵌入此结构体以获得统一的基础字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`              // 主键ID，自动递增，GORM自动管理
	CreatedAt time.Time      `json:"created_at"`                        // 创建时间，GORM自动设置
	UpdatedAt time.Time      `json:"updated_at"`                        // 更新时间，GORM自动维护
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 删除时间，用于软删除，建立索引
}

// 银行账户系统数据模型定义
// 本系统模拟一个简单的银行账户管理系统，包含用户、账户、交易、审计日志和通知等功能

// Account 银行账户模型
// 表示银行账户的基本信息，包括余额、类型、限额等
// 与User表建立多对一关系（一个用户可以有多个账户）
// 与Transaction表建立一对多关系（一个账户可以有多笔交易）
type Account struct {
	BaseModel           // 继承基础模型字段
	UserID      uint    `gorm:"not null;index" json:"user_id"`                          // 用户ID外键，建立索引，非空
	AccountType string  `gorm:"size:20;not null;index" json:"account_type"`             // 账户类型：savings(储蓄), checking(支票), credit(信用卡)
	Balance     float64 `gorm:"precision:15;scale:2;not null;default:0" json:"balance"` // 账户余额，精度15位，小数点后2位，默认为0
	Currency    string  `gorm:"size:3;not null;default:'CNY'" json:"currency"`          // 货币类型，3位货币代码，默认人民币
	IsActive    bool    `gorm:"default:true;index" json:"is_active"`                    // 账户是否激活，默认激活，建立索引
	DailyLimit  float64 `gorm:"precision:15;scale:2;default:10000" json:"daily_limit"`  // 日交易限额，默认10000

	// 关联关系定义
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`            // 所属用户，通过UserID外键关联
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"transactions,omitempty"` // 账户的所有交易记录
}

// User 用户模型
// 表示银行系统中的用户基本信息
// 与Account表建立一对多关系（一个用户可以拥有多个银行账户）
// 与Transaction表建立一对多关系（一个用户可以有多笔交易记录）
// 与AuditLog表建立一对多关系（一个用户可以有多条审计日志）
// 支持用户名和邮箱的唯一性约束，确保数据完整性
type User struct {
	BaseModel              // 继承基础模型字段
	Username    string     `gorm:"uniqueIndex;size:50;not null" json:"username"` // 用户名，最大50字符，唯一索引，非空
	Email       string     `gorm:"uniqueIndex;size:100;not null" json:"email"`   // 邮箱地址，最大100字符，唯一索引，非空
	FullName    string     `gorm:"size:100;not null" json:"full_name"`           // 用户全名，最大100字符，非空
	Phone       string     `gorm:"size:20;index" json:"phone"`                   // 手机号码，最大20字符，建立索引
	IsActive    bool       `gorm:"default:true;index" json:"is_active"`          // 用户是否激活，默认激活，建立索引
	LastLoginAt *time.Time `json:"last_login_at"`                                // 最后登录时间，可为空

	// 关联关系定义
	Accounts     []Account     `gorm:"foreignKey:UserID" json:"accounts,omitempty"`     // 用户拥有的所有银行账户
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"` // 用户的所有交易记录
	AuditLogs    []AuditLog    `gorm:"foreignKey:UserID" json:"audit_logs,omitempty"`   // 用户的所有审计日志
}

// Transaction 交易记录模型
// 表示银行系统中的所有交易记录，包括存款、取款、转账等操作
// 与Account表建立多对一关系（一个账户可以有多笔交易）
// 与User表建立多对一关系（一个用户可以有多笔交易）
// 支持转账操作，包含源账户和目标账户的关联
// 记录交易前后的余额变化，确保数据一致性
type Transaction struct {
	BaseModel               // 继承基础模型字段
	AccountID       uint    `gorm:"not null;index" json:"account_id"`                       // 账户ID外键，建立索引，非空
	UserID          uint    `gorm:"not null;index" json:"user_id"`                          // 用户ID外键，建立索引，非空
	TransactionType string  `gorm:"size:20;not null;index" json:"transaction_type"`         // 交易类型：deposit(存款), withdraw(取款), transfer(转账)
	Amount          float64 `gorm:"precision:15;scale:2;not null" json:"amount"`            // 交易金额，精度15位，小数点后2位，非空
	BalanceBefore   float64 `gorm:"precision:15;scale:2;not null" json:"balance_before"`    // 交易前账户余额，用于审计和对账
	BalanceAfter    float64 `gorm:"precision:15;scale:2;not null" json:"balance_after"`     // 交易后账户余额，用于审计和对账
	Description     string  `gorm:"size:500" json:"description"`                            // 交易描述，最大500字符
	Reference       string  `gorm:"size:100;index" json:"reference"`                        // 交易参考号，用于交易追踪，建立索引
	Status          string  `gorm:"size:20;not null;default:'pending';index" json:"status"` // 交易状态：pending(待处理), completed(已完成), failed(失败), cancelled(已取消)

	// 转账相关字段
	ToAccountID  *uint   `gorm:"index" json:"to_account_id,omitempty"`                // 转账目标账户ID，仅转账交易使用，建立索引
	TransferFee  float64 `gorm:"precision:15;scale:2;default:0" json:"transfer_fee"`  // 转账手续费，默认为0
	ExchangeRate float64 `gorm:"precision:10;scale:6;default:1" json:"exchange_rate"` // 汇率，用于跨币种转账，默认为1

	// 关联关系定义
	Account   Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`      // 交易所属账户
	User      User     `gorm:"foreignKey:UserID" json:"user,omitempty"`            // 交易发起用户
	ToAccount *Account `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"` // 转账目标账户（仅转账交易）
}

// AuditLog 审计日志模型
// 记录系统中所有重要操作的审计信息，用于安全监控和合规要求
// 与User表建立多对一关系（一个用户可以有多条审计日志）
// 记录数据变更的前后值，支持操作追溯和数据恢复
// 包含IP地址和用户代理信息，用于安全分析
type AuditLog struct {
	BaseModel          // 继承基础模型字段
	UserID      uint   `gorm:"not null;index" json:"user_id"`            // 操作用户ID外键，建立索引，非空
	Action      string `gorm:"size:50;not null;index" json:"action"`     // 操作类型：CREATE(创建), UPDATE(更新), DELETE(删除)
	TableName   string `gorm:"size:50;not null;index" json:"table_name"` // 操作的数据表名，建立索引
	RecordID    uint   `gorm:"not null;index" json:"record_id"`          // 操作记录的ID，建立索引
	OldValues   string `gorm:"type:text" json:"old_values"`              // 操作前的数据值，JSON格式存储
	NewValues   string `gorm:"type:text" json:"new_values"`              // 操作后的数据值，JSON格式存储
	IPAddress   string `gorm:"size:45" json:"ip_address"`                // 操作者IP地址，支持IPv4和IPv6
	UserAgent   string `gorm:"size:500" json:"user_agent"`               // 用户代理字符串，用于识别客户端
	Description string `gorm:"size:500" json:"description"`              // 操作描述，最大500字符

	// 关联关系定义
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 执行操作的用户
}

// NotificationLog 通知日志模型
// 记录系统发送给用户的所有通知信息，包括邮件、短信、推送等
// 与User表建立多对一关系（一个用户可以收到多条通知）
// 支持通知状态跟踪和重试机制，确保重要通知的送达
// 记录发送时间和错误信息，便于故障排查
type NotificationLog struct {
	BaseModel             // 继承基础模型字段
	UserID     uint       `gorm:"not null;index" json:"user_id"`                          // 接收用户ID外键，建立索引，非空
	Type       string     `gorm:"size:20;not null;index" json:"type"`                     // 通知类型：email(邮件), sms(短信), push(推送)
	Title      string     `gorm:"size:200;not null" json:"title"`                         // 通知标题，最大200字符，非空
	Content    string     `gorm:"type:text;not null" json:"content"`                      // 通知内容，文本类型，非空
	Status     string     `gorm:"size:20;not null;default:'pending';index" json:"status"` // 通知状态：pending(待发送), sent(已发送), failed(发送失败)
	SentAt     *time.Time `json:"sent_at"`                                                // 实际发送时间，可为空
	RetryCount int        `gorm:"default:0" json:"retry_count"`                           // 重试次数，默认为0
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`                             // 错误信息，发送失败时记录具体原因

	// 关联关系定义
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 通知接收用户
}

// GORM钩子函数实现
// 钩子函数是GORM提供的生命周期回调机制，允许在数据库操作的特定时点执行自定义逻辑
// 主要钩子类型：BeforeCreate, AfterCreate, BeforeUpdate, AfterUpdate, BeforeDelete, AfterDelete
// 钩子函数可以用于数据验证、默认值设置、审计日志记录、缓存更新等场景

// User模型的钩子函数
// 实现用户数据的自动验证、默认值设置和审计日志记录

// BeforeCreate 用户创建前钩子
// 在用户记录插入数据库之前执行，用于数据验证和默认值设置
// 参数 tx: GORM数据库事务对象，可用于执行额外的数据库操作
// 返回 error: 如果返回错误，将阻止创建操作并回滚事务
func (u *User) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户创建前: %s\n", u.Username)

	// 验证用户名格式
	// 用户名长度必须至少3个字符，确保用户名的可读性和唯一性
	if len(u.Username) < 3 {
		return errors.New("用户名长度不能少于3个字符")
	}

	// 验证邮箱格式（简单验证）
	// 检查邮箱长度和是否包含@符号，实际项目中应使用更严格的正则表达式验证
	if len(u.Email) < 5 || !contains(u.Email, "@") {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// AfterCreate 用户创建后钩子
// 在用户记录成功插入数据库之后执行，用于执行后续的业务逻辑
// 包括创建默认账户、记录审计日志、发送欢迎通知等操作
// 参数 tx: GORM数据库事务对象，确保所有操作在同一事务中执行
// 返回 error: 如果返回错误，将回滚整个事务，包括用户创建操作
func (u *User) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户创建后: %s (ID: %d)\n", u.Username, u.ID)

	// 自动创建默认储蓄账户
	// 每个新用户都会自动获得一个默认的储蓄账户，便于后续的银行业务操作
	defaultAccount := Account{
		UserID:      u.ID,      // 关联到新创建的用户
		AccountType: "savings", // 默认账户类型为储蓄账户
		Balance:     0,         // 初始余额为0
		Currency:    "CNY",     // 默认货币为人民币
		IsActive:    true,      // 账户默认为激活状态
		DailyLimit:  10000,     // 设置日交易限额为10000元
	}

	// 在同一事务中创建默认账户，确保数据一致性
	if err := tx.Create(&defaultAccount).Error; err != nil {
		fmt.Printf("[Hook] 创建默认账户失败: %v\n", err)
		return err
	}

	// 记录审计日志
	// 为新用户创建操作记录审计日志，用于安全监控和合规要求
	auditLog := AuditLog{
		UserID:      u.ID,                                                        // 操作用户ID
		Action:      "CREATE",                                                    // 操作类型为创建
		TableName:   "users",                                                     // 操作的表名
		RecordID:    u.ID,                                                        // 操作记录的ID
		NewValues:   fmt.Sprintf("username: %s, email: %s", u.Username, u.Email), // 记录新创建的数据值
		Description: "新用户注册",                                                     // 操作描述
	}

	// 在同一事务中创建审计日志
	if err := tx.Create(&auditLog).Error; err != nil {
		fmt.Printf("[Hook] 创建审计日志失败: %v\n", err)
		return err
	}

	// 发送欢迎通知
	// 为新用户创建欢迎通知，提升用户体验
	notification := NotificationLog{
		UserID:  u.ID,                                        // 通知接收用户
		Type:    "email",                                     // 通知类型为邮件
		Title:   "欢迎注册",                                      // 通知标题
		Content: fmt.Sprintf("欢迎 %s 注册我们的银行系统！", u.FullName), // 个性化通知内容
		Status:  "pending",                                   // 通知状态为待发送
	}

	// 在同一事务中创建通知记录
	if err := tx.Create(&notification).Error; err != nil {
		fmt.Printf("[Hook] 创建欢迎通知失败: %v\n", err)
		return err
	}

	fmt.Printf("[Hook] 用户 %s 的默认账户和通知已创建\n", u.Username)
	return nil
}

// BeforeUpdate 用户更新前钩子
// 在用户记录更新到数据库之前执行，用于数据验证和预处理
// 参数 tx: GORM数据库事务对象
// 返回 error: 如果返回错误，将阻止更新操作
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户更新前: %s\n", u.Username)

	// 手动更新时间戳
	// 虽然GORM会自动管理UpdatedAt字段，但在某些情况下可能需要手动设置
	u.UpdatedAt = time.Now()

	return nil
}

// AfterUpdate 用户更新后钩子
// 在用户记录成功更新到数据库之后执行，用于记录审计日志和后续处理
// 参数 tx: GORM数据库事务对象
// 返回 error: 如果返回错误，将回滚更新操作
func (u *User) AfterUpdate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户更新后: %s\n", u.Username)

	// 记录审计日志
	// 用户信息更新是敏感操作，需要记录审计日志用于安全监控
	auditLog := AuditLog{
		UserID:      u.ID,                                                        // 操作用户ID
		Action:      "UPDATE",                                                    // 操作类型为更新
		TableName:   "users",                                                     // 操作的表名
		RecordID:    u.ID,                                                        // 操作记录的ID
		NewValues:   fmt.Sprintf("username: %s, email: %s", u.Username, u.Email), // 记录更新后的数据值
		Description: "用户信息更新",                                                    // 操作描述
	}

	// 在同一事务中创建审计日志，确保数据一致性
	return tx.Create(&auditLog).Error
}

// BeforeDelete 用户删除前钩子
// 在用户记录从数据库删除之前执行，用于业务规则验证和数据完整性检查
// 参数 tx: GORM数据库事务对象
// 返回 error: 如果返回错误，将阻止删除操作
func (u *User) BeforeDelete(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户删除前: %s\n", u.Username)

	// 检查是否有活跃账户
	// 银行业务规则：用户如果还有活跃账户，则不能删除用户记录
	// 这是为了保护用户资产和维护数据完整性
	var activeAccountCount int64
	tx.Model(&Account{}).Where("user_id = ? AND is_active = ?", u.ID, true).Count(&activeAccountCount)

	// 如果用户还有活跃账户，阻止删除操作
	if activeAccountCount > 0 {
		return errors.New("用户还有活跃账户，无法删除")
	}

	return nil
}

// Account模型的钩子函数
// 实现账户数据的验证、默认值设置和业务规则检查

// BeforeCreate 账户创建前钩子
// 在账户记录插入数据库之前执行，用于数据验证和业务规则检查
// 确保账户类型的有效性和用户账户的唯一性约束
// 参数 tx: GORM数据库事务对象
// 返回 error: 如果返回错误，将阻止创建操作并回滚事务
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 账户创建前: 用户ID %d, 类型 %s\n", a.UserID, a.AccountType)

	// 验证账户类型
	// 银行系统只支持三种账户类型：储蓄账户、支票账户和信用卡账户
	validTypes := []string{"savings", "checking", "credit"}
	if !containsString(validTypes, a.AccountType) {
		return errors.New("无效的账户类型")
	}

	// 检查用户是否已有相同类型的账户
	// 业务规则：每个用户每种类型只能有一个活跃账户
	// 这是为了简化账户管理和避免业务逻辑复杂化
	var existingCount int64
	tx.Model(&Account{}).Where("user_id = ? AND account_type = ? AND is_active = ?",
		a.UserID, a.AccountType, true).Count(&existingCount)

	// 如果用户已有相同类型的活跃账户，阻止创建新账户
	if existingCount > 0 {
		return fmt.Errorf("用户已有 %s 类型的活跃账户", a.AccountType)
	}

	return nil
}

// AfterCreate 账户创建后钩子
// 在账户记录成功插入数据库之后执行，用于记录审计日志和后续处理
// 确保账户创建操作的可追溯性和合规性
// 参数 tx: GORM数据库事务对象
// 返回 error: 如果返回错误，将回滚整个事务，包括账户创建操作
func (a *Account) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 账户创建后: ID %d, 类型 %s\n", a.ID, a.AccountType)

	// 记录审计日志
	// 账户创建是重要的业务操作，需要记录审计日志用于合规和安全监控
	auditLog := AuditLog{
		UserID:      a.UserID,                                                                 // 操作用户ID
		Action:      "CREATE",                                                                 // 操作类型为创建
		TableName:   "accounts",                                                               // 操作的表名
		RecordID:    a.ID,                                                                     // 新创建账户的ID
		NewValues:   fmt.Sprintf("account_type: %s, balance: %.2f", a.AccountType, a.Balance), // 记录新账户的关键信息
		Description: "新账户创建",                                                                  // 操作描述
	}

	// 在同一事务中创建审计日志，确保数据一致性
	return tx.Create(&auditLog).Error
}

// Transaction模型的钩子函数
// 实现交易数据的验证、业务规则检查和自动字段生成
// 确保交易的合法性和数据完整性

// BeforeCreate 交易创建前钩子
// 在交易记录插入数据库之前执行，用于数据验证和预处理
// 包括交易类型验证、金额检查、余额验证、限额控制等关键业务逻辑
// 确保所有交易都符合银行业务规则和风险控制要求
// 参数 tx: GORM数据库事务对象，用于查询相关数据
// 返回 error: 如果返回错误，将阻止创建操作并回滚事务
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 交易创建前: 账户ID %d, 类型 %s, 金额 %.2f\n",
		t.AccountID, t.TransactionType, t.Amount)

	// 验证交易类型
	// 银行系统只支持三种基本交易类型：存款、取款和转账
	// 这是核心业务规则，确保系统只处理合法的交易类型
	validTypes := []string{"deposit", "withdraw", "transfer"}
	if !containsString(validTypes, t.TransactionType) {
		return errors.New("无效的交易类型")
	}

	// 验证交易金额
	// 银行业务规则：所有交易金额必须为正数，不允许零金额或负金额交易
	// 这是基本的数据完整性检查，防止异常交易
	if t.Amount <= 0 {
		return errors.New("交易金额必须大于0")
	}

	// 获取并验证账户信息
	// 确保交易的源账户存在且处于可用状态
	var account Account
	if err := tx.First(&account, t.AccountID).Error; err != nil {
		return fmt.Errorf("账户不存在: %v", err)
	}

	// 检查账户状态
	// 冻结或非激活账户不能进行任何交易操作
	// 这是重要的风险控制措施
	if !account.IsActive {
		return errors.New("账户已被冻结")
	}

	// 记录交易前余额
	// 用于审计追踪和数据一致性验证
	t.BalanceBefore = account.Balance

	// 对取款和转账交易进行额外验证
	// 这些交易会减少账户余额，需要进行余额和限额检查
	if t.TransactionType == "withdraw" || t.TransactionType == "transfer" {
		// 验证账户余额是否充足
		// 防止透支，确保账户资金安全
		if account.Balance < t.Amount {
			return errors.New("账户余额不足")
		}

		// 检查日交易限额
		// 计算当日已完成的取款和转账总额，防止超出日限额
		// 这是重要的风险控制和反洗钱措施
		var todayWithdrawTotal float64
		today := time.Now().Format("2006-01-02")
		tx.Model(&Transaction{}).
			Where("account_id = ? AND transaction_type IN ? AND DATE(created_at) = ? AND status = ?",
				t.AccountID, []string{"withdraw", "transfer"}, today, "completed").
			Select("COALESCE(SUM(amount), 0)").Scan(&todayWithdrawTotal)

		// 验证是否超出日限额
		if todayWithdrawTotal+t.Amount > account.DailyLimit {
			return fmt.Errorf("超出日限额 %.2f，今日已使用 %.2f",
				account.DailyLimit, todayWithdrawTotal)
		}

		// 计算交易后余额（减少）
		t.BalanceAfter = account.Balance - t.Amount
	} else {
		// 存款交易：计算交易后余额（增加）
		t.BalanceAfter = account.Balance + t.Amount
	}

	// 生成唯一交易参考号
	// 如果没有提供参考号，系统自动生成一个基于交易类型、账户ID和时间戳的唯一标识
	// 格式：交易类型_账户ID_时间戳，确保交易的可追溯性和唯一性
	if t.Reference == "" {
		t.Reference = fmt.Sprintf("%s_%d_%d",
			t.TransactionType, t.AccountID, time.Now().Unix())
	}

	return nil
}

// AfterCreate 交易创建后钩子
// 在交易记录成功插入数据库之后执行，用于执行后续的业务逻辑
// 包括更新账户余额、修改交易状态、记录审计日志、发送通知等关键操作
// 确保交易的完整性和业务流程的正确执行
// 参数 tx: GORM数据库事务对象，确保所有操作在同一事务中执行
// 返回 error: 如果返回错误，将回滚整个事务，包括交易创建操作
func (t *Transaction) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 交易创建后: ID %d, 参考号 %s\n", t.ID, t.Reference)

	// 更新账户余额
	// 根据交易类型计算余额变化：存款为正，取款和转账为负
	// 使用数据库级别的原子操作确保并发安全
	var balanceChange float64
	if t.TransactionType == "deposit" {
		balanceChange = t.Amount // 存款增加余额
	} else {
		balanceChange = -t.Amount // 取款和转账减少余额
	}

	// 使用GORM的Expr进行原子更新，避免并发问题
	// 直接在数据库层面进行余额计算，确保数据一致性
	if err := tx.Model(&Account{}).Where("id = ?", t.AccountID).
		Update("balance", gorm.Expr("balance + ?", balanceChange)).Error; err != nil {
		return fmt.Errorf("更新账户余额失败: %v", err)
	}

	// 更新交易状态为已完成
	// 交易创建时状态为pending，成功处理后更新为completed
	// 这是交易生命周期管理的重要环节
	if err := tx.Model(t).Update("status", "completed").Error; err != nil {
		return fmt.Errorf("更新交易状态失败: %v", err)
	}

	// 记录审计日志
	// 所有交易操作都需要记录审计日志，用于合规监管和安全审计
	// 记录交易的关键信息，便于后续追溯和分析
	auditLog := AuditLog{
		UserID:    t.UserID,       // 交易用户
		Action:    "CREATE",       // 操作类型
		TableName: "transactions", // 操作表名
		RecordID:  t.ID,           // 交易记录ID
		NewValues: fmt.Sprintf("type: %s, amount: %.2f, reference: %s", // 记录交易详情
			t.TransactionType, t.Amount, t.Reference),
		Description: fmt.Sprintf("%s 交易", t.TransactionType), // 操作描述
	}

	// 创建审计日志，如果失败只记录错误但不中断交易
	// 审计日志的失败不应该影响核心业务流程
	if err := tx.Create(&auditLog).Error; err != nil {
		fmt.Printf("[Hook] 创建审计日志失败: %v\n", err)
	}

	// 发送交易通知
	// 为用户发送交易完成通知，提升用户体验和安全感知
	// 通知包含交易类型、金额和余额等关键信息
	var user User
	if err := tx.First(&user, t.UserID).Error; err == nil {
		// 创建短信通知记录
		notification := NotificationLog{
			UserID: t.UserID, // 通知接收用户
			Type:   "sms",    // 通知类型为短信
			Title:  "交易通知",   // 通知标题
			Content: fmt.Sprintf("您的账户发生 %s 交易，金额: %.2f，余额: %.2f", // 个性化通知内容
				t.TransactionType, t.Amount, t.BalanceAfter),
			Status: "pending", // 通知状态为待发送
		}

		// 创建通知记录，如果失败只记录错误但不中断交易
		// 通知发送的失败不应该影响核心交易流程
		if err := tx.Create(&notification).Error; err != nil {
			fmt.Printf("[Hook] 创建交易通知失败: %v\n", err)
		}
	}

	return nil
}

// 辅助函数
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 数据库初始化函数
// 支持SQLite和MySQL两种数据库类型，提供灵活的配置选项

// initDB 初始化SQLite数据库连接（向后兼容）
// 使用SQLite作为默认数据库，适用于开发和测试环境
// 返回配置好的GORM数据库实例
func initDB() *gorm.DB {
	// 使用SQLite数据库，文件存储在当前目录
	// SQLite是轻量级数据库，无需额外安装，适合开发和小型应用
	db, err := gorm.Open(sqlite.Open("level4_transactions_hooks.db"), &gorm.Config{
		// 启用详细日志模式，便于开发调试
		// 在生产环境中应该使用logger.Silent或logger.Warn级别
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移数据库表结构
	// GORM会根据结构体定义自动创建或更新表结构
	// 包括所有业务模型：用户、账户、交易、审计日志、通知日志
	err = db.AutoMigrate(&User{}, &Account{}, &Transaction{}, &AuditLog{}, &NotificationLog{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db
}

// InitDatabase 通用数据库初始化函数
// 支持SQLite和MySQL两种数据库类型，根据配置自动选择
// 参数 config: 数据库配置信息，包含数据库类型和连接参数
// 返回 *gorm.DB: 配置好的GORM数据库实例
// 返回 error: 初始化过程中的错误信息
func InitDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 根据数据库类型选择相应的驱动和连接方式
	switch config.Type {
	case SQLite:
		// SQLite数据库连接
		// 适用于开发、测试和小型应用场景
		db, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
			// 禁用外键约束检查（SQLite特有配置）
			DisableForeignKeyConstraintWhenMigrating: true,
		})

	case MySQL:
		// MySQL数据库连接
		// 适用于生产环境和大型应用场景
		db, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
			// MySQL特有配置
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false, // 使用复数表名
			},
		})

		// 配置MySQL连接池参数
		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				// 设置连接池参数，优化数据库性能
				sqlDB.SetMaxIdleConns(config.MaxIdleConns)   // 最大空闲连接数
				sqlDB.SetMaxOpenConns(config.MaxOpenConns)   // 最大打开连接数
				sqlDB.SetConnMaxLifetime(config.MaxLifetime) // 连接最大生存时间
			}
		}

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %v", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 执行数据库表结构迁移
	// 自动创建或更新所有业务模型对应的数据库表
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	return db, nil
}

// AutoMigrate 执行数据库表结构迁移
// 根据Go结构体定义自动创建或更新数据库表结构
// 参数 db: GORM数据库实例
// 返回 error: 迁移过程中的错误信息
func AutoMigrate(db *gorm.DB) error {
	// 按照依赖关系顺序迁移表结构
	// 先创建基础表（User），再创建依赖表（Account, Transaction等）
	return db.AutoMigrate(
		&User{},            // 用户表
		&Account{},         // 账户表（依赖User）
		&Transaction{},     // 交易表（依赖User和Account）
		&AuditLog{},        // 审计日志表（依赖User）
		&NotificationLog{}, // 通知日志表（依赖User）
	)
}

// 事务操作示例
// 以下函数展示了银行系统中的核心业务操作，包括用户管理、账户操作、转账等
// 所有操作都使用GORM事务确保数据一致性和完整性

// CreateUserWithAccount 创建用户和账户（事务操作）
// 这是一个完整的用户注册流程，包括创建用户、自动创建默认账户、处理初始存款
// 使用GORM事务确保所有操作的原子性
// 参数 db: GORM数据库实例
// 参数 username: 用户名（唯一标识）
// 参数 email: 用户邮箱地址
// 参数 fullName: 用户全名
// 参数 initialDeposit: 初始存款金额（可选，0表示无初始存款）
// 返回 error: 操作过程中的错误信息
func CreateUserWithAccount(db *gorm.DB, username, email, fullName string, initialDeposit float64) error {
	// 使用GORM事务包装器，自动处理事务的开始、提交和回滚
	// 如果函数返回error，事务会自动回滚；如果返回nil，事务会自动提交
	return db.Transaction(func(tx *gorm.DB) error {
		// 创建用户对象
		// 会触发User模型的BeforeCreate钩子（设置默认值、验证数据）
		// 和AfterCreate钩子（记录审计日志、自动创建默认账户）
		user := User{
			Username: username, // 用户名，必须唯一
			Email:    email,    // 邮箱地址，用于通知和登录
			FullName: fullName, // 用户全名，用于显示和身份验证
			IsActive: true,     // 账户状态，新用户默认为活跃状态
		}

		// 在事务中创建用户记录
		// 如果用户名或邮箱已存在，会返回唯一约束错误
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %v", err)
		}

		// 处理初始存款（如果有）
		// 这是一个可选步骤，允许用户在注册时进行首次存款
		if initialDeposit > 0 {
			// 获取用户的默认储蓄账户
			// 该账户应该在User的AfterCreate钩子中自动创建
			var account Account
			if err := tx.Where("user_id = ? AND account_type = ?", user.ID, "savings").First(&account).Error; err != nil {
				return fmt.Errorf("获取默认账户失败: %v", err)
			}

			// 创建初始存款交易记录
			// 会触发Transaction模型的BeforeCreate钩子（验证金额、生成参考号）
			// 和AfterCreate钩子（更新账户余额、记录审计日志、发送通知）
			transaction := Transaction{
				AccountID:       account.ID,     // 关联到用户的默认账户
				UserID:          user.ID,        // 关联到刚创建的用户
				TransactionType: "deposit",      // 交易类型：存款
				Amount:          initialDeposit, // 存款金额
				Description:     "初始存款",         // 交易描述
				Status:          "pending",      // 交易状态：待处理
			}

			// 在事务中创建存款交易记录
			// 钩子函数会自动更新账户余额并记录相关日志
			if err := tx.Create(&transaction).Error; err != nil {
				return fmt.Errorf("创建初始存款交易失败: %v", err)
			}
		}

		fmt.Printf("✓ 用户 %s 创建成功，ID: %d\n", username, user.ID)
		return nil
	})
}

// TransferMoney 转账操作（事务）
// 实现两个账户间的资金转移，确保转账的原子性和一致性
// 包括余额验证、账户状态检查、交易记录创建等完整流程
// 参数 db: GORM数据库实例
// 参数 fromAccountID: 转出账户ID
// 参数 toAccountID: 转入账户ID
// 参数 amount: 转账金额（必须大于0）
// 参数 description: 转账描述信息
// 返回 error: 操作过程中的错误信息
func TransferMoney(db *gorm.DB, fromAccountID, toAccountID uint, amount float64, description string) error {
	// 使用GORM事务确保转账操作的原子性
	// 转账涉及多个数据库操作，必须保证要么全部成功，要么全部失败
	return db.Transaction(func(tx *gorm.DB) error {
		// 验证转出和转入账户的存在性和活跃状态
		// 只有活跃的账户才能参与转账操作
		var fromAccount, toAccount Account

		// 查询并验证转出账户
		// 检查账户是否存在且处于活跃状态
		if err := tx.Where("id = ? AND is_active = ?", fromAccountID, true).First(&fromAccount).Error; err != nil {
			return fmt.Errorf("源账户不存在或已冻结: %v", err)
		}

		// 查询并验证转入账户
		// 同样检查账户的存在性和活跃状态
		if err := tx.Where("id = ? AND is_active = ?", toAccountID, true).First(&toAccount).Error; err != nil {
			return fmt.Errorf("目标账户不存在或已冻结: %v", err)
		}

		// 防止自转账操作
		// 同一账户间的转账是无意义的，应该被禁止
		if fromAccountID == toAccountID {
			return errors.New("不能向同一账户转账")
		}

		// 创建转出交易记录
		// 记录资金从源账户转出的操作
		// 会触发Transaction的BeforeCreate钩子进行余额验证
		withdrawTx := Transaction{
			AccountID:       fromAccountID,                                         // 转出账户ID
			UserID:          fromAccount.UserID,                                    // 转出账户所属用户ID
			TransactionType: "transfer",                                            // 交易类型：转账
			Amount:          amount,                                                // 转账金额
			Description:     fmt.Sprintf("转账至账户 %d: %s", toAccountID, description), // 交易描述
			ToAccountID:     &toAccountID,                                          // 目标账户ID（用于关联转账记录）
			Status:          "pending",                                             // 交易状态：待处理
		}

		// 在事务中创建转出交易记录
		// BeforeCreate钩子会验证余额是否充足
		// AfterCreate钩子会更新账户余额并记录审计日志
		if err := tx.Create(&withdrawTx).Error; err != nil {
			return fmt.Errorf("创建转出交易失败: %v", err)
		}

		// 创建转入交易记录
		// 记录资金转入目标账户的操作，与转出交易形成完整的转账记录
		depositTx := Transaction{
			AccountID:       toAccountID,                                                // 转入账户ID
			UserID:          toAccount.UserID,                                           // 转入账户所属用户ID
			TransactionType: "deposit",                                                  // 交易类型：存款
			Amount:          amount,                                                     // 转账金额
			Description:     fmt.Sprintf("来自账户 %d 的转账: %s", fromAccountID, description), // 交易描述
			Reference:       withdrawTx.Reference,                                       // 使用相同的参考号关联转账记录
			Status:          "pending",                                                  // 交易状态：待处理
		}

		// 手动设置余额变化信息
		// 虽然钩子函数会自动处理余额更新，但这里预设值有助于数据一致性检查
		depositTx.BalanceBefore = toAccount.Balance         // 转账前余额
		depositTx.BalanceAfter = toAccount.Balance + amount // 转账后余额

		// 在事务中创建转入交易记录
		// AfterCreate钩子会更新目标账户余额并发送通知
		if err := tx.Create(&depositTx).Error; err != nil {
			return fmt.Errorf("创建转入交易失败: %v", err)
		}

		fmt.Printf("✓ 转账成功: 从账户 %d 向账户 %d 转账 %.2f\n", fromAccountID, toAccountID, amount)
		return nil
	})
}

// BatchCreateTransactions 批量创建交易（事务）
// 在单个事务中批量创建多个交易记录，确保数据一致性
// 适用于批量导入、批量处理等场景
// 参数 db: GORM数据库实例
// 参数 transactions: 待创建的交易记录切片
// 返回 error: 操作过程中的错误信息
func BatchCreateTransactions(db *gorm.DB, transactions []Transaction) error {
	// 使用事务确保批量操作的原子性
	// 如果任何一个交易创建失败，整个批次都会回滚
	return db.Transaction(func(tx *gorm.DB) error {
		// 遍历所有交易记录，逐个创建
		// 每个交易的创建都会触发相应的钩子函数
		for i, transaction := range transactions {
			// 创建单个交易记录
			// 钩子函数会处理余额验证、更新和日志记录
			if err := tx.Create(&transaction).Error; err != nil {
				return fmt.Errorf("创建第 %d 个交易失败: %v", i+1, err)
			}
		}

		// 批量操作成功完成的日志输出
		fmt.Printf("✓ 批量创建 %d 个交易成功\n", len(transactions))
		return nil
	})
}

// ProcessBatchTransactions 处理批量交易（事务）
// 批量处理交易记录，包括状态更新和余额计算
// 适用于批量交易处理、定时任务等场景
// 参数 db: GORM数据库实例
// 参数 transactions: 待处理的交易记录切片
// 返回 error: 操作过程中的错误信息
func ProcessBatchTransactions(db *gorm.DB, transactions []Transaction) error {
	// 使用事务确保批量处理的原子性
	// 如果任何一个交易处理失败，整个批次都会回滚
	return db.Transaction(func(tx *gorm.DB) error {
		// 遍历所有交易记录，逐个处理
		// 每个交易的处理都会触发相应的钩子函数
		for i, transaction := range transactions {
			// 创建交易记录
			// 钩子函数会处理余额验证、更新和日志记录
			if err := tx.Create(&transaction).Error; err != nil {
				return fmt.Errorf("处理第 %d 个交易失败: %v", i+1, err)
			}

			// 更新交易状态为已完成
			// 模拟交易处理完成后的状态更新
			if err := tx.Model(&transaction).Update("status", "completed").Error; err != nil {
				return fmt.Errorf("更新第 %d 个交易状态失败: %v", i+1, err)
			}
		}

		// 批量处理成功完成的日志输出
		fmt.Printf("✓ 批量处理 %d 个交易成功\n", len(transactions))
		return nil
	})
}

// UpdateAccountStatus 更新账户状态（事务）
// 用于冻结或激活账户，包括状态更新、审计日志记录和用户通知
// 这是一个重要的风控操作，需要完整的审计追踪
// 参数 db: GORM数据库实例
// 参数 accountID: 要更新的账户ID
// 参数 isActive: 新的账户状态（true=激活，false=冻结）
// 参数 reason: 状态变更原因（用于审计和通知）
// 返回 error: 操作过程中的错误信息
func UpdateAccountStatus(db *gorm.DB, accountID uint, isActive bool, reason string) error {
	// 使用事务确保状态更新、审计记录和通知的一致性
	return db.Transaction(func(tx *gorm.DB) error {
		// 查询要更新的账户信息
		// 需要获取当前状态以便记录变更历史
		var account Account
		if err := tx.First(&account, accountID).Error; err != nil {
			return fmt.Errorf("账户不存在: %v", err)
		}

		// 更新账户的活跃状态
		// 会触发Account模型的BeforeUpdate和AfterUpdate钩子
		if err := tx.Model(&account).Update("is_active", isActive).Error; err != nil {
			return fmt.Errorf("更新账户状态失败: %v", err)
		}

		// 创建审计日志记录
		// 记录账户状态变更的详细信息，用于合规性和安全审计
		auditLog := AuditLog{
			UserID:      account.UserID,                                 // 账户所属用户ID
			Action:      "UPDATE",                                       // 操作类型：更新
			TableName:   "accounts",                                     // 操作的表名
			RecordID:    accountID,                                      // 被操作记录的ID
			OldValues:   fmt.Sprintf("is_active: %t", account.IsActive), // 变更前的值
			NewValues:   fmt.Sprintf("is_active: %t", isActive),         // 变更后的值
			Description: fmt.Sprintf("账户状态变更: %s", reason),              // 变更描述和原因
		}

		// 在事务中创建审计日志
		// 确保审计记录与业务操作的一致性
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("创建审计日志失败: %v", err)
		}

		// 创建用户通知记录
		// 及时通知用户账户状态变更，提升用户体验和透明度
		notification := NotificationLog{
			UserID: account.UserID, // 通知目标用户ID
			Type:   "email",        // 通知类型：邮件
			Title:  "账户状态变更通知",     // 通知标题
			Content: fmt.Sprintf("您的账户状态已变更为: %s，原因: %s", // 通知内容
				func() string {
					// 根据状态值返回中文描述
					if isActive {
						return "激活" // 账户激活
					}
					return "冻结" // 账户冻结
				}(), reason),
			Status: "pending", // 通知状态：待发送
		}

		// 在事务中创建通知记录
		// 确保通知记录与业务操作的一致性
		if err := tx.Create(&notification).Error; err != nil {
			return fmt.Errorf("创建状态变更通知失败: %v", err)
		}

		// 输出操作成功信息
		// 用于调试和日志记录
		fmt.Printf("✓ 账户 %d 状态更新为: %t\n", accountID, isActive)

		// 事务成功完成
		// 所有操作（状态更新、审计记录、通知创建）都已成功
		return nil
	})
}

// ==================== 查询函数 ====================
// 以下函数用于查询和获取数据库中的信息
// 这些函数不涉及数据修改，主要用于数据展示和业务查询

// GetAccountBalance 获取账户余额
// 根据账户ID查询账户的当前余额信息
// 参数 db: GORM数据库实例
// 参数 accountID: 要查询的账户ID
// 返回 float64: 账户余额
// 返回 error: 查询过程中的错误信息
func GetAccountBalance(db *gorm.DB, accountID uint) (float64, error) {
	// 定义账户变量用于存储查询结果
	var account Account

	// 根据账户ID查询账户信息
	// 使用First方法查询第一条匹配记录
	if err := db.First(&account, accountID).Error; err != nil {
		// 如果查询失败（如账户不存在），返回0余额和错误信息
		return 0, err
	}

	// 返回账户余额
	return account.Balance, nil
}

// GetUserTransactionHistory 获取用户交易历史
// 查询指定用户的交易记录，按时间倒序排列
// 支持分页限制和关联账户信息预加载
// 参数 db: GORM数据库实例
// 参数 userID: 要查询的用户ID
// 参数 limit: 返回记录数量限制
// 返回 []Transaction: 交易记录切片
// 返回 error: 查询过程中的错误信息
func GetUserTransactionHistory(db *gorm.DB, userID uint, limit int) ([]Transaction, error) {
	// 定义交易记录切片用于存储查询结果
	var transactions []Transaction

	// 构建查询条件并执行查询
	err := db.Where("user_id = ?", userID). // 按用户ID筛选
						Order("created_at DESC"). // 按创建时间倒序排列（最新的在前）
						Limit(limit).             // 限制返回记录数量
						Preload("Account").       // 预加载关联的账户信息
						Find(&transactions).Error // 执行查询并获取错误信息

	// 返回查询结果和错误信息
	return transactions, err
}

// GetAuditLogs 获取审计日志
// 查询指定用户的审计日志记录，用于合规性检查和操作追踪
// 按时间倒序排列，支持分页限制
// 参数 db: GORM数据库实例
// 参数 userID: 要查询的用户ID
// 参数 limit: 返回记录数量限制
// 返回 []AuditLog: 审计日志记录切片
// 返回 error: 查询过程中的错误信息
func GetAuditLogs(db *gorm.DB, userID uint, limit int) ([]AuditLog, error) {
	// 定义审计日志切片用于存储查询结果
	var logs []AuditLog

	// 构建查询条件并执行查询
	err := db.Where("user_id = ?", userID). // 按用户ID筛选
						Order("created_at DESC"). // 按创建时间倒序排列（最新的在前）
						Limit(limit).             // 限制返回记录数量
						Find(&logs).Error         // 执行查询并获取错误信息

	// 返回查询结果和错误信息
	return logs, err
}

// ==================== 主函数演示 ====================
// 演示GORM事务管理和钩子函数的完整使用流程
// 包括用户创建、转账操作、批量交易、状态更新等核心业务场景
func main() {
	// 声明数据库连接变量
	var db *gorm.DB
	var err error

	// 输出程序标题
	fmt.Println("=== GORM Level 4 事务和钩子练习 ===")

	// 数据库初始化演示
	// 支持SQLite和MySQL两种数据库类型
	fmt.Println("\n=== 数据库初始化演示 ===")

	// 演示1：使用SQLite数据库（默认方式）
	// fmt.Println("\n--- 使用SQLite数据库 ---")
	// db := initDB()
	// fmt.Println("✓ SQLite数据库初始化完成")

	// // 演示2：使用通用初始化函数（SQLite配置）
	// fmt.Println("\n--- 使用通用初始化函数（SQLite配置） ---")
	// sqliteConfig := GetDefaultConfig()
	// sqliteDB, err := InitDatabase(sqliteConfig)
	// if err != nil {
	// 	fmt.Printf("SQLite数据库初始化失败: %v\n", err)
	// } else {
	// 	fmt.Println("✓ 通用SQLite数据库初始化完成")
	// 	// 关闭额外的数据库连接
	// 	sqlDB, _ := sqliteDB.DB()
	// 	if sqlDB != nil {
	// 		sqlDB.Close()
	// 	}
	// }

	// 演示3：MySQL数据库配置示例（注释掉，需要实际MySQL服务器）

	// MySQL数据库初始化示例（需要实际MySQL服务器）
	// 请根据实际环境修改连接参数
	mysqlDSN := "root:fastbee@tcp(192.168.100.124:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig := GetMySQLConfig(mysqlDSN)
	mysqlDB, err := InitDatabase(mysqlConfig)
	if err != nil {
		fmt.Printf("MySQL数据库初始化失败: %v\n", err)
		// MySQL连接失败，使用SQLite作为备选
		fmt.Println("\n=== 使用SQLite数据库进行业务演示 ===")
		sqliteConfig := GetDefaultConfig()
		db, err = InitDatabase(sqliteConfig)
		if err != nil {
			fmt.Printf("SQLite数据库初始化失败: %v\n", err)
			return
		}
		fmt.Println("✓ SQLite数据库初始化完成")
	} else {
		fmt.Println("✓ MySQL数据库初始化完成")
		// 使用MySQL数据库进行后续业务演示
		fmt.Println("\n=== 使用MySQL数据库进行业务演示 ===")
		db = mysqlDB
		// 注意：不要在这里关闭数据库连接，因为后续还要使用
		// 在程序结束时再关闭连接
		defer func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}()
	}

	// ==================== 演示1：创建用户（触发钩子函数） ====================
	// 演示用户注册流程，包括用户创建、账户创建和初始存款
	// 每个操作都会触发相应的GORM钩子函数进行数据验证和业务处理
	fmt.Println("\n=== 演示1：创建用户和账户 ===")

	// 创建用户Alice，初始存款1000元
	// 会触发User.BeforeCreate、User.AfterCreate、Account.BeforeCreate、Account.AfterCreate
	// 以及Transaction.BeforeCreate、Transaction.AfterCreate等钩子函数
	err = CreateUserWithAccount(db, "alice", "alice@example.com", "Alice Wang", 1000.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// 创建用户Bob，初始存款2000元
	// 同样会触发完整的钩子函数链
	err = CreateUserWithAccount(db, "bob", "bob@example.com", "Bob Chen", 2000.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// 创建用户Charlie，初始存款500元
	// 演示不同金额的初始存款处理
	err = CreateUserWithAccount(db, "charlie", "charlie@example.com", "Charlie Li", 500.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// ==================== 演示2：转账操作（事务） ====================
	// 演示账户间转账的完整流程，包括事务管理、余额更新和钩子函数触发
	// 转账操作需要确保原子性，要么全部成功，要么全部失败
	fmt.Println("\n=== 演示2：转账操作 ===")

	// 通过用户名查询对应的账户信息
	// 使用子查询关联users表获取账户ID
	var aliceAccount, bobAccount Account
	db.Where("user_id = (SELECT id FROM users WHERE username = ?)", "alice").First(&aliceAccount)
	db.Where("user_id = (SELECT id FROM users WHERE username = ?)", "bob").First(&bobAccount)

	// 执行转账操作：Alice向Bob转账300元
	// 会在事务中创建两条交易记录（转出和转入）
	// 触发Transaction.BeforeCreate和Transaction.AfterCreate钩子进行余额更新
	err = TransferMoney(db, aliceAccount.ID, bobAccount.ID, 300.0, "还款")
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	}

	// 查询转账后的账户余额
	// 验证转账操作是否正确执行
	aliceBalance, _ := GetAccountBalance(db, aliceAccount.ID)
	bobBalance, _ := GetAccountBalance(db, bobAccount.ID)
	fmt.Printf("转账后余额 - Alice: %.2f, Bob: %.2f\n", aliceBalance, bobBalance)

	// ==================== 演示3：批量交易（事务） ====================
	// 演示批量创建交易记录的事务处理
	// 确保所有交易要么全部成功，要么全部失败，维护数据一致性
	fmt.Println("\n=== 演示3：批量交易 ===")

	// 准备批量交易数据
	// 包括Alice的存款交易和Bob的取款交易
	batchTransactions := []Transaction{
		{
			AccountID:       aliceAccount.ID,     // Alice的账户ID
			UserID:          aliceAccount.UserID, // Alice的用户ID
			TransactionType: "deposit",           // 交易类型：存款
			Amount:          100.0,               // 交易金额：100元
			Description:     "工资",                // 交易描述
			Status:          "pending",           // 交易状态：待处理
		},
		{
			AccountID:       bobAccount.ID,     // Bob的账户ID
			UserID:          bobAccount.UserID, // Bob的用户ID
			TransactionType: "withdraw",        // 交易类型：取款
			Amount:          50.0,              // 交易金额：50元
			Description:     "ATM取款",           // 交易描述
			Status:          "pending",         // 交易状态：待处理
		},
	}

	// 执行批量交易操作
	// 在单个事务中处理所有交易，确保数据一致性
	// 每个交易都会触发相应的钩子函数进行余额更新和审计记录
	err = BatchCreateTransactions(db, batchTransactions)
	if err != nil {
		fmt.Printf("批量交易失败: %v\n", err)
	}

	// ==================== 演示4：账户状态管理（事务） ====================
	// 演示账户冻结和解冻操作，包括状态更新、审计记录和用户通知
	// 这是重要的风控功能，用于处理可疑交易和账户安全
	fmt.Println("\n=== 演示4：账户状态管理 ===")

	// 冻结Alice的账户
	// 会触发Account.BeforeUpdate和Account.AfterUpdate钩子
	// 同时创建审计日志和用户通知
	err = UpdateAccountStatus(db, aliceAccount.ID, false, "可疑交易，临时冻结")
	if err != nil {
		fmt.Printf("冻结账户失败: %v\n", err)
	}

	// 测试冻结账户的交易限制功能
	// 尝试在冻结账户上创建交易，应该被钩子函数拒绝
	fmt.Println("\n尝试在冻结账户上进行交易:")
	frozenTransaction := Transaction{
		AccountID:       aliceAccount.ID,     // 已冻结的Alice账户ID
		UserID:          aliceAccount.UserID, // Alice的用户ID
		TransactionType: "withdraw",          // 交易类型：取款
		Amount:          100.0,               // 交易金额：100元
		Description:     "测试冻结账户交易",          // 交易描述
		Status:          "pending",           // 交易状态：待处理
	}

	// 尝试创建交易记录
	// Transaction.BeforeCreate钩子会检查账户状态并拒绝冻结账户的交易
	err = db.Create(&frozenTransaction).Error
	if err != nil {
		fmt.Printf("✓ 冻结账户交易被正确拒绝: %v\n", err)
	}

	// 解冻Alice的账户
	// 恢复账户的正常使用功能
	err = UpdateAccountStatus(db, aliceAccount.ID, true, "调查完成，恢复正常")
	if err != nil {
		fmt.Printf("解冻账户失败: %v\n", err)
	}

	// ==================== 演示5：交易历史和审计日志 ====================
	// 演示查询功能，展示如何获取用户的交易记录和操作审计
	// 这些功能对于用户查询和合规性检查非常重要
	fmt.Println("\n=== 演示5：交易历史和审计日志 ===")

	// 获取Alice用户信息
	// 通过用户名查找用户记录，获取用户ID用于后续查询
	var aliceUser User
	db.Where("username = ?", "alice").First(&aliceUser)

	// 查询Alice的交易历史记录
	// 获取最近10条交易记录，按时间倒序排列
	// 包含预加载的账户信息，便于显示完整的交易详情
	transactions, err := GetUserTransactionHistory(db, aliceUser.ID, 10)
	if err != nil {
		fmt.Printf("获取交易历史失败: %v\n", err)
	} else {
		fmt.Printf("Alice的交易历史 (%d条):\n", len(transactions))
		// 遍历并显示每条交易记录
		// 包括交易类型、金额、状态和描述信息
		for _, tx := range transactions {
			fmt.Printf("  - %s: %.2f (%s) - %s\n",
				tx.TransactionType, tx.Amount, tx.Status, tx.Description)
		}
	}

	// 查询Alice的审计日志记录
	// 获取最近10条审计记录，用于合规性检查和操作追踪
	// 审计日志记录了所有重要的数据库操作
	auditLogs, err := GetAuditLogs(db, aliceUser.ID, 10)
	if err != nil {
		fmt.Printf("获取审计日志失败: %v\n", err)
	} else {
		fmt.Printf("\nAlice的审计日志 (%d条):\n", len(auditLogs))
		// 遍历并显示每条审计记录
		// 包括操作类型、表名和操作描述
		for _, log := range auditLogs {
			fmt.Printf("  - %s %s: %s\n", log.Action, log.TableName, log.Description)
		}
	}

	// ==================== 演示6：错误处理和回滚 ====================
	// 演示事务的错误处理和自动回滚机制
	// 展示数据验证、业务规则检查和事务完整性保护
	fmt.Println("\n=== 演示6：错误处理和回滚 ===")

	// 测试用户数据验证
	// 尝试创建不符合验证规则的用户，应该被钩子函数拒绝
	fmt.Println("尝试创建无效用户:")
	// 用户名过短（少于2个字符）且邮箱格式无效
	// User.BeforeCreate钩子会进行数据验证并拒绝创建
	err = CreateUserWithAccount(db, "x", "invalid-email", "Invalid User", 100.0)
	if err != nil {
		fmt.Printf("✓ 无效用户创建被正确拒绝: %v\n", err)
	}

	// 测试业务规则验证
	// 尝试进行余额不足的转账，应该被业务逻辑拒绝
	fmt.Println("\n尝试余额不足的转账:")
	// Bob账户余额不足以支付10000元的转账
	// TransferMoney函数会检查余额并拒绝交易
	err = TransferMoney(db, bobAccount.ID, aliceAccount.ID, 10000.0, "大额转账测试")
	if err != nil {
		fmt.Printf("✓ 余额不足转账被正确拒绝: %v\n", err)
	}

	// ==================== 最终余额检查 ====================
	// 验证所有操作完成后的账户余额状态
	// 确保所有事务操作的正确性和数据一致性
	fmt.Println("\n=== 最终余额检查 ===")
	// 查询Alice和Bob的最终账户余额
	// 验证经过所有演示操作后的余额变化是否正确
	aliceBalance, _ = GetAccountBalance(db, aliceAccount.ID)
	bobBalance, _ = GetAccountBalance(db, bobAccount.ID)
	fmt.Printf("最终余额 - Alice: %.2f, Bob: %.2f\n", aliceBalance, bobBalance)

	// 程序演示结束
	// 本演示展示了GORM钩子函数在银行系统中的完整应用
	// 包括用户管理、账户操作、事务处理、状态管理和审计功能

	fmt.Println("\n=== Level 4 事务和钩子练习完成 ===")
}

// ==================== MySQL演示函数 ====================
// DemoWithMySQL 使用MySQL数据库进行事务和钩子演示
// 展示如何在MySQL环境中使用GORM的事务管理和钩子函数
// 注意：此函数需要实际的MySQL服务器环境
// 参数 host: MySQL服务器地址
// 参数 port: MySQL服务器端口
// 参数 database: 数据库名称
// 参数 username: 数据库用户名
// 参数 password: 数据库密码
func DemoWithMySQL(host string, port int, database, username, password string) {
	fmt.Println("\n=== MySQL数据库事务和钩子演示 ===")

	// 创建MySQL DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	// 创建MySQL数据库配置
	mysqlConfig := GetMySQLConfig(dsn)
	fmt.Printf("连接MySQL数据库: %s:%d/%s\n", host, port, database)

	// 初始化MySQL数据库连接
	db, err := InitDatabase(mysqlConfig)
	if err != nil {
		fmt.Printf("MySQL数据库连接失败: %v\n", err)
		return
	}
	fmt.Println("✓ MySQL数据库连接成功")

	// 确保在函数结束时关闭数据库连接
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
			fmt.Println("✓ MySQL数据库连接已关闭")
		}
	}()

	// 清理测试数据（可选）
	fmt.Println("\n--- 清理测试数据 ---")
	db.Exec("DELETE FROM notification_logs")
	db.Exec("DELETE FROM audit_logs")
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM users")
	fmt.Println("✓ 测试数据清理完成")

	// 演示1：MySQL环境下的用户创建
	fmt.Println("\n=== MySQL演示1：用户创建 ===")
	err = CreateUserWithAccount(db, "mysql_alice", "alice@mysql.com", "Alice MySQL", 1500.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}

	err = CreateUserWithAccount(db, "mysql_bob", "bob@mysql.com", "Bob MySQL", 2500.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}
	fmt.Println("✓ MySQL用户创建完成")

	// 演示2：MySQL环境下的转账操作
	fmt.Println("\n=== MySQL演示2：转账操作 ===")
	var aliceAccount, bobAccount Account
	db.Where("user_id = (SELECT id FROM users WHERE username = ?)", "mysql_alice").First(&aliceAccount)
	db.Where("user_id = (SELECT id FROM users WHERE username = ?)", "mysql_bob").First(&bobAccount)

	err = TransferMoney(db, aliceAccount.ID, bobAccount.ID, 500.0, "MySQL转账测试")
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		aliceBalance, _ := GetAccountBalance(db, aliceAccount.ID)
		bobBalance, _ := GetAccountBalance(db, bobAccount.ID)
		fmt.Printf("转账后余额 - Alice: %.2f, Bob: %.2f\n", aliceBalance, bobBalance)
	}

	// 演示3：MySQL事务性能测试
	fmt.Println("\n=== MySQL演示3：事务性能测试 ===")
	start := time.Now()

	// 批量创建交易记录测试MySQL事务性能
	batchTransactions := []Transaction{
		{
			AccountID:       aliceAccount.ID,
			UserID:          aliceAccount.UserID,
			TransactionType: "deposit",
			Amount:          200.0,
			Description:     "MySQL批量存款1",
			Status:          "pending",
		},
		{
			AccountID:       bobAccount.ID,
			UserID:          bobAccount.UserID,
			TransactionType: "deposit",
			Amount:          300.0,
			Description:     "MySQL批量存款2",
			Status:          "pending",
		},
	}

	err = ProcessBatchTransactions(db, batchTransactions)
	if err != nil {
		fmt.Printf("批量交易失败: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("✓ MySQL批量交易完成，耗时: %v\n", duration)
	}

	// 演示4：查询MySQL中的审计日志
	fmt.Println("\n=== MySQL演示4：审计日志查询 ===")
	var aliceUser User
	db.Where("username = ?", "mysql_alice").First(&aliceUser)

	auditLogs, err := GetAuditLogs(db, aliceUser.ID, 5)
	if err != nil {
		fmt.Printf("获取审计日志失败: %v\n", err)
	} else {
		fmt.Printf("MySQL审计日志 (%d条):\n", len(auditLogs))
		for i, log := range auditLogs {
			fmt.Printf("%d. %s - %s (表: %s)\n", i+1, log.Action, log.Description, log.TableName)
		}
	}

	// 最终余额检查
	fmt.Println("\n=== MySQL最终余额检查 ===")
	aliceBalance, _ := GetAccountBalance(db, aliceAccount.ID)
	bobBalance, _ := GetAccountBalance(db, bobAccount.ID)
	fmt.Printf("MySQL最终余额 - Alice: %.2f, Bob: %.2f\n", aliceBalance, bobBalance)

	fmt.Println("\n=== MySQL演示完成 ===")
}

// ==================== 使用示例 ====================
// 以下是如何使用MySQL演示函数的示例代码
// 请根据实际MySQL服务器配置修改参数
/*
func ExampleUsage() {
	// 示例1：本地MySQL服务器
	DemoWithMySQL("localhost", 3306, "gorm_demo", "root", "password")

	// 示例2：远程MySQL服务器
	DemoWithMySQL("192.168.1.100", 3306, "banking_system", "admin", "secure_password")

	// 示例3：云数据库
	DemoWithMySQL("mysql.example.com", 3306, "production_db", "app_user", "app_password")
}
*/

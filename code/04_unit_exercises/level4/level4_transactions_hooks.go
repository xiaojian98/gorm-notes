// 04_unit_exercises/level4_transactions_hooks.go - Level 4 事务和钩子练习
// 对应文档：03_GORM单元练习_基础技能训练.md

package main

import (
	"errors"
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

// 数据模型定义
type Account struct {
	BaseModel
	UserID      uint    `gorm:"not null;index" json:"user_id"`
	AccountType string  `gorm:"size:20;not null;index" json:"account_type"` // savings, checking, credit
	Balance     float64 `gorm:"precision:15;scale:2;not null;default:0" json:"balance"`
	Currency    string  `gorm:"size:3;not null;default:'CNY'" json:"currency"`
	IsActive    bool    `gorm:"default:true;index" json:"is_active"`
	DailyLimit  float64 `gorm:"precision:15;scale:2;default:10000" json:"daily_limit"`
	
	// 关联关系
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"transactions,omitempty"`
}

type User struct {
	BaseModel
	Username    string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	FullName    string `gorm:"size:100;not null" json:"full_name"`
	Phone       string `gorm:"size:20;index" json:"phone"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	
	// 关联关系
	Accounts     []Account     `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	AuditLogs    []AuditLog    `gorm:"foreignKey:UserID" json:"audit_logs,omitempty"`
}

type Transaction struct {
	BaseModel
	AccountID       uint    `gorm:"not null;index" json:"account_id"`
	UserID          uint    `gorm:"not null;index" json:"user_id"`
	TransactionType string  `gorm:"size:20;not null;index" json:"transaction_type"` // deposit, withdraw, transfer
	Amount          float64 `gorm:"precision:15;scale:2;not null" json:"amount"`
	BalanceBefore   float64 `gorm:"precision:15;scale:2;not null" json:"balance_before"`
	BalanceAfter    float64 `gorm:"precision:15;scale:2;not null" json:"balance_after"`
	Description     string  `gorm:"size:500" json:"description"`
	Reference       string  `gorm:"size:100;index" json:"reference"`
	Status          string  `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, completed, failed, cancelled
	
	// 转账相关字段
	ToAccountID   *uint   `gorm:"index" json:"to_account_id,omitempty"`
	TransferFee   float64 `gorm:"precision:15;scale:2;default:0" json:"transfer_fee"`
	ExchangeRate  float64 `gorm:"precision:10;scale:6;default:1" json:"exchange_rate"`
	
	// 关联关系
	Account   Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	User      User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ToAccount *Account `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
}

type AuditLog struct {
	BaseModel
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Action      string `gorm:"size:50;not null;index" json:"action"`
	TableName   string `gorm:"size:50;not null;index" json:"table_name"`
	RecordID    uint   `gorm:"not null;index" json:"record_id"`
	OldValues   string `gorm:"type:text" json:"old_values"`
	NewValues   string `gorm:"type:text" json:"new_values"`
	IPAddress   string `gorm:"size:45" json:"ip_address"`
	UserAgent   string `gorm:"size:500" json:"user_agent"`
	Description string `gorm:"size:500" json:"description"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type NotificationLog struct {
	BaseModel
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Type        string `gorm:"size:20;not null;index" json:"type"` // email, sms, push
	Title       string `gorm:"size:200;not null" json:"title"`
	Content     string `gorm:"type:text;not null" json:"content"`
	Status      string `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, sent, failed
	SentAt      *time.Time `json:"sent_at"`
	RetryCount  int    `gorm:"default:0" json:"retry_count"`
	ErrorMsg    string `gorm:"type:text" json:"error_msg"`
	
	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// 钩子函数实现

// User 钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户创建前: %s\n", u.Username)
	
	// 验证用户名格式
	if len(u.Username) < 3 {
		return errors.New("用户名长度不能少于3个字符")
	}
	
	// 验证邮箱格式（简单验证）
	if len(u.Email) < 5 || !contains(u.Email, "@") {
		return errors.New("邮箱格式不正确")
	}
	
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户创建后: %s (ID: %d)\n", u.Username, u.ID)
	
	// 自动创建默认储蓄账户
	defaultAccount := Account{
		UserID:      u.ID,
		AccountType: "savings",
		Balance:     0,
		Currency:    "CNY",
		IsActive:    true,
		DailyLimit:  10000,
	}
	
	if err := tx.Create(&defaultAccount).Error; err != nil {
		fmt.Printf("[Hook] 创建默认账户失败: %v\n", err)
		return err
	}
	
	// 记录审计日志
	auditLog := AuditLog{
		UserID:      u.ID,
		Action:      "CREATE",
		TableName:   "users",
		RecordID:    u.ID,
		NewValues:   fmt.Sprintf("username: %s, email: %s", u.Username, u.Email),
		Description: "新用户注册",
	}
	
	if err := tx.Create(&auditLog).Error; err != nil {
		fmt.Printf("[Hook] 创建审计日志失败: %v\n", err)
		return err
	}
	
	// 发送欢迎通知
	notification := NotificationLog{
		UserID:  u.ID,
		Type:    "email",
		Title:   "欢迎注册",
		Content: fmt.Sprintf("欢迎 %s 注册我们的银行系统！", u.FullName),
		Status:  "pending",
	}
	
	if err := tx.Create(&notification).Error; err != nil {
		fmt.Printf("[Hook] 创建欢迎通知失败: %v\n", err)
		return err
	}
	
	fmt.Printf("[Hook] 用户 %s 的默认账户和通知已创建\n", u.Username)
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户更新前: %s\n", u.Username)
	
	// 记录更新时间
	u.UpdatedAt = time.Now()
	
	return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户更新后: %s\n", u.Username)
	
	// 记录审计日志
	auditLog := AuditLog{
		UserID:      u.ID,
		Action:      "UPDATE",
		TableName:   "users",
		RecordID:    u.ID,
		NewValues:   fmt.Sprintf("username: %s, email: %s", u.Username, u.Email),
		Description: "用户信息更新",
	}
	
	return tx.Create(&auditLog).Error
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	fmt.Printf("[Hook] 用户删除前: %s\n", u.Username)
	
	// 检查是否有活跃账户
	var activeAccountCount int64
	tx.Model(&Account{}).Where("user_id = ? AND is_active = ?", u.ID, true).Count(&activeAccountCount)
	
	if activeAccountCount > 0 {
		return errors.New("用户还有活跃账户，无法删除")
	}
	
	return nil
}

// Account 钩子函数
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 账户创建前: 用户ID %d, 类型 %s\n", a.UserID, a.AccountType)
	
	// 验证账户类型
	validTypes := []string{"savings", "checking", "credit"}
	if !containsString(validTypes, a.AccountType) {
		return errors.New("无效的账户类型")
	}
	
	// 检查用户是否已有相同类型的账户
	var existingCount int64
	tx.Model(&Account{}).Where("user_id = ? AND account_type = ? AND is_active = ?", 
		a.UserID, a.AccountType, true).Count(&existingCount)
	
	if existingCount > 0 {
		return fmt.Errorf("用户已有 %s 类型的活跃账户", a.AccountType)
	}
	
	return nil
}

func (a *Account) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 账户创建后: ID %d, 类型 %s\n", a.ID, a.AccountType)
	
	// 记录审计日志
	auditLog := AuditLog{
		UserID:      a.UserID,
		Action:      "CREATE",
		TableName:   "accounts",
		RecordID:    a.ID,
		NewValues:   fmt.Sprintf("account_type: %s, balance: %.2f", a.AccountType, a.Balance),
		Description: "新账户创建",
	}
	
	return tx.Create(&auditLog).Error
}

// Transaction 钩子函数
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 交易创建前: 账户ID %d, 类型 %s, 金额 %.2f\n", 
		t.AccountID, t.TransactionType, t.Amount)
	
	// 验证交易类型
	validTypes := []string{"deposit", "withdraw", "transfer"}
	if !containsString(validTypes, t.TransactionType) {
		return errors.New("无效的交易类型")
	}
	
	// 验证金额
	if t.Amount <= 0 {
		return errors.New("交易金额必须大于0")
	}
	
	// 获取账户当前余额
	var account Account
	if err := tx.First(&account, t.AccountID).Error; err != nil {
		return fmt.Errorf("账户不存在: %v", err)
	}
	
	if !account.IsActive {
		return errors.New("账户已被冻结")
	}
	
	// 记录交易前余额
	t.BalanceBefore = account.Balance
	
	// 验证余额（对于取款和转账）
	if t.TransactionType == "withdraw" || t.TransactionType == "transfer" {
		if account.Balance < t.Amount {
			return errors.New("账户余额不足")
		}
		
		// 检查日限额
		var todayWithdrawTotal float64
		today := time.Now().Format("2006-01-02")
		tx.Model(&Transaction{}).
			Where("account_id = ? AND transaction_type IN ? AND DATE(created_at) = ? AND status = ?",
				t.AccountID, []string{"withdraw", "transfer"}, today, "completed").
			Select("COALESCE(SUM(amount), 0)").Scan(&todayWithdrawTotal)
		
		if todayWithdrawTotal+t.Amount > account.DailyLimit {
			return fmt.Errorf("超出日限额 %.2f，今日已使用 %.2f", 
				account.DailyLimit, todayWithdrawTotal)
		}
		
		t.BalanceAfter = account.Balance - t.Amount
	} else {
		t.BalanceAfter = account.Balance + t.Amount
	}
	
	// 生成交易参考号
	if t.Reference == "" {
		t.Reference = fmt.Sprintf("%s_%d_%d", 
			t.TransactionType, t.AccountID, time.Now().Unix())
	}
	
	return nil
}

func (t *Transaction) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[Hook] 交易创建后: ID %d, 参考号 %s\n", t.ID, t.Reference)
	
	// 更新账户余额
	var balanceChange float64
	if t.TransactionType == "deposit" {
		balanceChange = t.Amount
	} else {
		balanceChange = -t.Amount
	}
	
	if err := tx.Model(&Account{}).Where("id = ?", t.AccountID).
		Update("balance", gorm.Expr("balance + ?", balanceChange)).Error; err != nil {
		return fmt.Errorf("更新账户余额失败: %v", err)
	}
	
	// 更新交易状态为完成
	if err := tx.Model(t).Update("status", "completed").Error; err != nil {
		return fmt.Errorf("更新交易状态失败: %v", err)
	}
	
	// 记录审计日志
	auditLog := AuditLog{
		UserID:      t.UserID,
		Action:      "CREATE",
		TableName:   "transactions",
		RecordID:    t.ID,
		NewValues:   fmt.Sprintf("type: %s, amount: %.2f, reference: %s", 
			t.TransactionType, t.Amount, t.Reference),
		Description: fmt.Sprintf("%s 交易", t.TransactionType),
	}
	
	if err := tx.Create(&auditLog).Error; err != nil {
		fmt.Printf("[Hook] 创建审计日志失败: %v\n", err)
	}
	
	// 发送交易通知
	var user User
	if err := tx.First(&user, t.UserID).Error; err == nil {
		notification := NotificationLog{
			UserID:  t.UserID,
			Type:    "sms",
			Title:   "交易通知",
			Content: fmt.Sprintf("您的账户发生 %s 交易，金额: %.2f，余额: %.2f", 
				t.TransactionType, t.Amount, t.BalanceAfter),
			Status:  "pending",
		}
		
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

// 数据库初始化
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("level4_transactions_hooks.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{}, &Account{}, &Transaction{}, &AuditLog{}, &NotificationLog{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db
}

// 事务操作示例

// CreateUserWithAccount 创建用户和账户（事务）
func CreateUserWithAccount(db *gorm.DB, username, email, fullName string, initialDeposit float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{
			Username: username,
			Email:    email,
			FullName: fullName,
			IsActive: true,
		}
		
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %v", err)
		}
		
		// 如果有初始存款，进行存款交易
		if initialDeposit > 0 {
			// 获取默认账户
			var account Account
			if err := tx.Where("user_id = ? AND account_type = ?", user.ID, "savings").First(&account).Error; err != nil {
				return fmt.Errorf("获取默认账户失败: %v", err)
			}
			
			// 创建存款交易
			transaction := Transaction{
				AccountID:       account.ID,
				UserID:          user.ID,
				TransactionType: "deposit",
				Amount:          initialDeposit,
				Description:     "初始存款",
				Status:          "pending",
			}
			
			if err := tx.Create(&transaction).Error; err != nil {
				return fmt.Errorf("创建初始存款交易失败: %v", err)
			}
		}
		
		fmt.Printf("✓ 用户 %s 创建成功，ID: %d\n", username, user.ID)
		return nil
	})
}

// TransferMoney 转账操作（事务）
func TransferMoney(db *gorm.DB, fromAccountID, toAccountID uint, amount float64, description string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 验证账户存在且活跃
		var fromAccount, toAccount Account
		
		if err := tx.Where("id = ? AND is_active = ?", fromAccountID, true).First(&fromAccount).Error; err != nil {
			return fmt.Errorf("源账户不存在或已冻结: %v", err)
		}
		
		if err := tx.Where("id = ? AND is_active = ?", toAccountID, true).First(&toAccount).Error; err != nil {
			return fmt.Errorf("目标账户不存在或已冻结: %v", err)
		}
		
		if fromAccountID == toAccountID {
			return errors.New("不能向同一账户转账")
		}
		
		// 创建转出交易
		withdrawTx := Transaction{
			AccountID:       fromAccountID,
			UserID:          fromAccount.UserID,
			TransactionType: "transfer",
			Amount:          amount,
			Description:     fmt.Sprintf("转账至账户 %d: %s", toAccountID, description),
			ToAccountID:     &toAccountID,
			Status:          "pending",
		}
		
		if err := tx.Create(&withdrawTx).Error; err != nil {
			return fmt.Errorf("创建转出交易失败: %v", err)
		}
		
		// 创建转入交易
		depositTx := Transaction{
			AccountID:       toAccountID,
			UserID:          toAccount.UserID,
			TransactionType: "deposit",
			Amount:          amount,
			Description:     fmt.Sprintf("来自账户 %d 的转账: %s", fromAccountID, description),
			Reference:       withdrawTx.Reference, // 使用相同的参考号
			Status:          "pending",
		}
		
		// 手动设置余额，因为钩子函数会自动处理
		depositTx.BalanceBefore = toAccount.Balance
		depositTx.BalanceAfter = toAccount.Balance + amount
		
		if err := tx.Create(&depositTx).Error; err != nil {
			return fmt.Errorf("创建转入交易失败: %v", err)
		}
		
		fmt.Printf("✓ 转账成功: 从账户 %d 向账户 %d 转账 %.2f\n", fromAccountID, toAccountID, amount)
		return nil
	})
}

// BatchCreateTransactions 批量创建交易（事务）
func BatchCreateTransactions(db *gorm.DB, transactions []Transaction) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for i, transaction := range transactions {
			if err := tx.Create(&transaction).Error; err != nil {
				return fmt.Errorf("创建第 %d 个交易失败: %v", i+1, err)
			}
		}
		
		fmt.Printf("✓ 批量创建 %d 个交易成功\n", len(transactions))
		return nil
	})
}

// UpdateAccountStatus 更新账户状态（事务）
func UpdateAccountStatus(db *gorm.DB, accountID uint, isActive bool, reason string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var account Account
		if err := tx.First(&account, accountID).Error; err != nil {
			return fmt.Errorf("账户不存在: %v", err)
		}
		
		// 更新账户状态
		if err := tx.Model(&account).Update("is_active", isActive).Error; err != nil {
			return fmt.Errorf("更新账户状态失败: %v", err)
		}
		
		// 记录审计日志
		auditLog := AuditLog{
			UserID:      account.UserID,
			Action:      "UPDATE",
			TableName:   "accounts",
			RecordID:    accountID,
			OldValues:   fmt.Sprintf("is_active: %t", account.IsActive),
			NewValues:   fmt.Sprintf("is_active: %t", isActive),
			Description: fmt.Sprintf("账户状态变更: %s", reason),
		}
		
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("创建审计日志失败: %v", err)
		}
		
		// 发送状态变更通知
		notification := NotificationLog{
			UserID:  account.UserID,
			Type:    "email",
			Title:   "账户状态变更通知",
			Content: fmt.Sprintf("您的账户状态已变更为: %s，原因: %s", 
				func() string {
					if isActive {
						return "激活"
					}
					return "冻结"
				}(), reason),
			Status: "pending",
		}
		
		if err := tx.Create(&notification).Error; err != nil {
			return fmt.Errorf("创建状态变更通知失败: %v", err)
		}
		
		fmt.Printf("✓ 账户 %d 状态更新为: %t\n", accountID, isActive)
		return nil
	})
}

// 查询函数

// GetAccountBalance 获取账户余额
func GetAccountBalance(db *gorm.DB, accountID uint) (float64, error) {
	var account Account
	if err := db.First(&account, accountID).Error; err != nil {
		return 0, err
	}
	return account.Balance, nil
}

// GetUserTransactionHistory 获取用户交易历史
func GetUserTransactionHistory(db *gorm.DB, userID uint, limit int) ([]Transaction, error) {
	var transactions []Transaction
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Preload("Account").
		Find(&transactions).Error
	return transactions, err
}

// GetAuditLogs 获取审计日志
func GetAuditLogs(db *gorm.DB, userID uint, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// 主函数演示
func main() {
	fmt.Println("=== GORM Level 4 事务和钩子练习 ===")

	// 初始化数据库
	db := initDB()
	fmt.Println("✓ 数据库初始化完成")

	// 演示1：创建用户（触发钩子函数）
	fmt.Println("\n=== 演示1：创建用户和账户 ===")
	
	err := CreateUserWithAccount(db, "alice", "alice@example.com", "Alice Wang", 1000.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}
	
	err = CreateUserWithAccount(db, "bob", "bob@example.com", "Bob Chen", 2000.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}
	
	err = CreateUserWithAccount(db, "charlie", "charlie@example.com", "Charlie Li", 500.0)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
	}

	// 演示2：转账操作（事务）
	fmt.Println("\n=== 演示2：转账操作 ===")
	
	// 获取账户ID
	var aliceAccount, bobAccount Account
	db.Where("user_id = (SELECT id FROM users WHERE username = ?)", "alice").First(&aliceAccount)
	db.Where("user_id = (SELECT id FROM users WHERE username = ?)", "bob").First(&bobAccount)
	
	// 转账
	err = TransferMoney(db, aliceAccount.ID, bobAccount.ID, 300.0, "还款")
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	}
	
	// 查看余额
	aliceBalance, _ := GetAccountBalance(db, aliceAccount.ID)
	bobBalance, _ := GetAccountBalance(db, bobAccount.ID)
	fmt.Printf("转账后余额 - Alice: %.2f, Bob: %.2f\n", aliceBalance, bobBalance)

	// 演示3：批量交易（事务）
	fmt.Println("\n=== 演示3：批量交易 ===")
	
	batchTransactions := []Transaction{
		{
			AccountID:       aliceAccount.ID,
			UserID:          aliceAccount.UserID,
			TransactionType: "deposit",
			Amount:          100.0,
			Description:     "工资",
			Status:          "pending",
		},
		{
			AccountID:       bobAccount.ID,
			UserID:          bobAccount.UserID,
			TransactionType: "withdraw",
			Amount:          50.0,
			Description:     "ATM取款",
			Status:          "pending",
		},
	}
	
	err = BatchCreateTransactions(db, batchTransactions)
	if err != nil {
		fmt.Printf("批量交易失败: %v\n", err)
	}

	// 演示4：账户状态管理（事务）
	fmt.Println("\n=== 演示4：账户状态管理 ===")
	
	err = UpdateAccountStatus(db, aliceAccount.ID, false, "可疑交易，临时冻结")
	if err != nil {
		fmt.Printf("冻结账户失败: %v\n", err)
	}
	
	// 尝试在冻结账户上进行交易（应该失败）
	fmt.Println("\n尝试在冻结账户上进行交易:")
	frozenTransaction := Transaction{
		AccountID:       aliceAccount.ID,
		UserID:          aliceAccount.UserID,
		TransactionType: "withdraw",
		Amount:          100.0,
		Description:     "测试冻结账户交易",
		Status:          "pending",
	}
	
	err = db.Create(&frozenTransaction).Error
	if err != nil {
		fmt.Printf("✓ 冻结账户交易被正确拒绝: %v\n", err)
	}
	
	// 解冻账户
	err = UpdateAccountStatus(db, aliceAccount.ID, true, "调查完成，恢复正常")
	if err != nil {
		fmt.Printf("解冻账户失败: %v\n", err)
	}

	// 演示5：查看交易历史和审计日志
	fmt.Println("\n=== 演示5：交易历史和审计日志 ===")
	
	// 获取Alice的交易历史
	var aliceUser User
	db.Where("username = ?", "alice").First(&aliceUser)
	
	transactions, err := GetUserTransactionHistory(db, aliceUser.ID, 10)
	if err != nil {
		fmt.Printf("获取交易历史失败: %v\n", err)
	} else {
		fmt.Printf("Alice的交易历史 (%d条):\n", len(transactions))
		for _, tx := range transactions {
			fmt.Printf("  - %s: %.2f (%s) - %s\n", 
				tx.TransactionType, tx.Amount, tx.Status, tx.Description)
		}
	}
	
	// 获取审计日志
	auditLogs, err := GetAuditLogs(db, aliceUser.ID, 10)
	if err != nil {
		fmt.Printf("获取审计日志失败: %v\n", err)
	} else {
		fmt.Printf("\nAlice的审计日志 (%d条):\n", len(auditLogs))
		for _, log := range auditLogs {
			fmt.Printf("  - %s %s: %s\n", log.Action, log.TableName, log.Description)
		}
	}

	// 演示6：错误处理和回滚
	fmt.Println("\n=== 演示6：错误处理和回滚 ===")
	
	// 尝试创建无效用户（应该失败并回滚）
	fmt.Println("尝试创建无效用户:")
	err = CreateUserWithAccount(db, "x", "invalid-email", "Invalid User", 100.0)
	if err != nil {
		fmt.Printf("✓ 无效用户创建被正确拒绝: %v\n", err)
	}
	
	// 尝试余额不足的转账（应该失败）
	fmt.Println("\n尝试余额不足的转账:")
	err = TransferMoney(db, bobAccount.ID, aliceAccount.ID, 10000.0, "大额转账测试")
	if err != nil {
		fmt.Printf("✓ 余额不足转账被正确拒绝: %v\n", err)
	}

	// 最终余额检查
	fmt.Println("\n=== 最终余额检查 ===")
	aliceBalance, _ = GetAccountBalance(db, aliceAccount.ID)
	bobBalance, _ = GetAccountBalance(db, bobAccount.ID)
	fmt.Printf("最终余额 - Alice: %.2f, Bob: %.2f\n", aliceBalance, bobBalance)

	fmt.Println("\n=== Level 4 事务和钩子练习完成 ===")
}
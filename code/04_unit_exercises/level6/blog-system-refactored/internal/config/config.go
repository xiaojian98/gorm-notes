package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// AppConfig 应用程序配置结构体
// 包含应用程序运行所需的所有配置信息
type AppConfig struct {
	App      AppSettings     `json:"app"`      // 应用程序设置
	Database *DatabaseConfig `json:"database"` // 数据库配置
	Server   ServerConfig    `json:"server"`   // 服务器配置
	Logging  LoggingConfig   `json:"logging"`  // 日志配置
}

// AppSettings 应用程序基本设置
type AppSettings struct {
	Name        string `json:"name"`        // 应用程序名称
	Version     string `json:"version"`     // 版本号
	Environment string `json:"environment"` // 运行环境: development, production, test
	Debug       bool   `json:"debug"`       // 是否开启调试模式
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `json:"host"`          // 服务器主机地址
	Port         int    `json:"port"`          // 服务器端口
	ReadTimeout  int    `json:"read_timeout"`  // 读取超时时间(秒)
	WriteTimeout int    `json:"write_timeout"` // 写入超时时间(秒)
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `json:"level"`       // 日志级别: debug, info, warn, error
	Format     string `json:"format"`      // 日志格式: json, text
	Output     string `json:"output"`      // 输出目标: stdout, file
	FilePath   string `json:"file_path"`   // 日志文件路径
	MaxSize    int    `json:"max_size"`    // 日志文件最大大小(MB)
	MaxBackups int    `json:"max_backups"` // 保留的日志文件数量
	MaxAge     int    `json:"max_age"`     // 日志文件保留天数
}

// GetDefaultConfig 获取默认配置
// 返回一个包含默认设置的应用程序配置
func GetDefaultConfig() *AppConfig {
	return &AppConfig{
		App: AppSettings{
			Name:        "Blog System",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Database: &DatabaseConfig{
			Type:     "mysql",
			Host:     "192.168.100.124",
			Port:     3306,
			Username: "root",
			Password: "fastbee",
			DBName:   "gorm_test",
			Charset:  "utf8mb4",
		},
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			FilePath:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
		},
	}
}

// LoadConfigFromFile 从文件加载配置
// 参数: filePath - 配置文件路径
// 返回: *AppConfig - 应用程序配置, error - 错误信息
func LoadConfigFromFile(filePath string) (*AppConfig, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", filePath)
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON配置
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// SaveConfigToFile 保存配置到文件
// 参数: config - 应用程序配置, filePath - 配置文件路径
// 返回: error - 错误信息
func SaveConfigToFile(config *AppConfig, filePath string) error {
	// 序列化配置为JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// Validate 验证配置的有效性
// 返回: error - 验证错误信息
func (c *AppConfig) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("应用程序名称不能为空")
	}

	if c.App.Version == "" {
		return fmt.Errorf("应用程序版本不能为空")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535范围内")
	}

	if c.Database == nil {
		return fmt.Errorf("数据库配置不能为空")
	}

	return nil
}

// IsDevelopment 判断是否为开发环境
// 返回: bool - 是否为开发环境
func (c *AppConfig) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction 判断是否为生产环境
// 返回: bool - 是否为生产环境
func (c *AppConfig) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetServerAddress 获取服务器地址
// 返回: string - 完整的服务器地址
func (c *AppConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

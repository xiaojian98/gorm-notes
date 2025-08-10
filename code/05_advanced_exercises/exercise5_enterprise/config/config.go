package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Email    EmailConfig    `mapstructure:"email"`
	Payment  PaymentConfig  `mapstructure:"payment"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	MaxHeaderMB  int           `mapstructure:"max_header_mb"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	LogLevel        string        `mapstructure:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	ExpireDuration   time.Duration `mapstructure:"expire_duration"`
	RefreshDuration  time.Duration `mapstructure:"refresh_duration"`
	Issuer           string        `mapstructure:"issuer"`
	SigningMethod    string        `mapstructure:"signing_method"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path         string   `mapstructure:"path"`
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	ImageResize  bool     `mapstructure:"image_resize"`
	ImageWidth   int      `mapstructure:"image_width"`
	ImageHeight  int      `mapstructure:"image_height"`
	ImageQuality int      `mapstructure:"image_quality"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	Alipay AlipayConfig `mapstructure:"alipay"`
	Wechat WechatConfig `mapstructure:"wechat"`
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID        string `mapstructure:"app_id"`
	PrivateKey   string `mapstructure:"private_key"`
	PublicKey    string `mapstructure:"public_key"`
	NotifyURL    string `mapstructure:"notify_url"`
	ReturnURL    string `mapstructure:"return_url"`
	SignType     string `mapstructure:"sign_type"`
	IsSandbox    bool   `mapstructure:"is_sandbox"`
}

// WechatConfig 微信支付配置
type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	MchID     string `mapstructure:"mch_id"`
	APIKey    string `mapstructure:"api_key"`
	NotifyURL string `mapstructure:"notify_url"`
	IsSandbox bool   `mapstructure:"is_sandbox"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.max_header_mb", 1)

	// 数据库默认配置
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "edu_platform")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "1h")
	viper.SetDefault("database.log_level", "info")

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")
	viper.SetDefault("redis.idle_timeout", "5m")

	// JWT默认配置
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_duration", "24h")
	viper.SetDefault("jwt.refresh_duration", "168h")
	viper.SetDefault("jwt.issuer", "edu-platform")
	viper.SetDefault("jwt.signing_method", "HS256")

	// 日志默认配置
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.output", "stdout")
	viper.SetDefault("logger.filename", "logs/app.log")
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_age", 30)
	viper.SetDefault("logger.max_backups", 10)
	viper.SetDefault("logger.compress", true)

	// 文件上传默认配置
	viper.SetDefault("upload.path", "uploads")
	viper.SetDefault("upload.max_size", 10485760) // 10MB
	viper.SetDefault("upload.allowed_types", []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx"})
	viper.SetDefault("upload.image_resize", true)
	viper.SetDefault("upload.image_width", 800)
	viper.SetDefault("upload.image_height", 600)
	viper.SetDefault("upload.image_quality", 80)

	// 邮件默认配置
	viper.SetDefault("email.smtp_host", "smtp.gmail.com")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.from_name", "在线教育平台")

	// 支付默认配置
	viper.SetDefault("payment.alipay.sign_type", "RSA2")
	viper.SetDefault("payment.alipay.is_sandbox", true)
	viper.SetDefault("payment.wechat.is_sandbox", true)
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetServerAddr 获取服务器地址
func (c *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key" {
		return fmt.Errorf("请设置JWT密钥")
	}

	if c.Database.Username == "" {
		return fmt.Errorf("请设置数据库用户名")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("请设置数据库名称")
	}

	return nil
}
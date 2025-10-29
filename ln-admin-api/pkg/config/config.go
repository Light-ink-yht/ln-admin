package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Cfg *Config

// Config 应用配置
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	Middleware MiddlewareConfig `mapstructure:"middleware"`
	Captcha    CaptchaConfig    `mapstructure:"captcha"`
	Password   PasswordConfig   `mapstructure:"password"`
	Redis      RedisConfig      `mapstructure:"redis"`
	SMS        SMSConfig        `mapstructure:"sms"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type   string       `mapstructure:"type"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Oracle OracleConfig `mapstructure:"oracle"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
}

// OracleConfig Oracle配置
type OracleConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	ServiceName     string `mapstructure:"service_name"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	Expire        int    `mapstructure:"expire"`
	RefreshExpire int    `mapstructure:"refresh_expire"`
	Issuer        string `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	Path       string `mapstructure:"path"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	CORS CORSConfig `mapstructure:"cors"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Length int `mapstructure:"length"`
	Width  int `mapstructure:"width"`
	Height int `mapstructure:"height"`
	Expiry int `mapstructure:"expiry"`
}

// PasswordConfig 密码配置
type PasswordConfig struct {
	Cost int `mapstructure:"cost"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	DialTimeout  int    `mapstructure:"dial_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// SMSConfig 短信配置
type SMSConfig struct {
	Enabled bool `mapstructure:"enabled"` // 是否启用短信功能（开发环境可关闭）
}

// Load 加载配置
func Load(configPath string) error {
	// 获取配置文件的绝对路径
	if !filepath.IsAbs(configPath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("获取工作目录失败: %w", err)
		}
		configPath = filepath.Join(wd, configPath)
	}

	// 设置配置文件路径和名称
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	if c.Type == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.MySQL.User,
			c.MySQL.Password,
			c.MySQL.Host,
			c.MySQL.Port,
			c.MySQL.DBName,
			c.MySQL.Charset,
		)
	} else if c.Type == "oracle" {
		// Oracle DSN 格式（GORM oracle 驱动使用 godror）:
		// 格式1: user/password@host:port/service_name (推荐)
		// 格式2: user/password@(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=host)(PORT=port))(CONNECT_DATA=(SERVICE_NAME=service_name)))
		return fmt.Sprintf("%s/%s@%s:%d/%s",
			c.Oracle.User,
			c.Oracle.Password,
			c.Oracle.Host,
			c.Oracle.Port,
			c.Oracle.ServiceName,
		)
	}
	return ""
}

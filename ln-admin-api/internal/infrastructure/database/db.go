package database

import (
	"fmt"
	"time"

	infraLogger "github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	var dialector gorm.Dialector
	dsn := config.Cfg.Database.GetDSN()

	switch config.Cfg.Database.Type {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "oracle":
		// Oracle 驱动需要特殊的配置方式
		//dialector = oracle.New(oracle.Config{
		//	DSN: dsn,
		//})
	default:
		return fmt.Errorf("不支持的数据库类型: %s", config.Cfg.Database.Type)
	}

	// GORM配置
	cfg := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层sql.DB设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if config.Cfg.Database.Type == "mysql" {
		// MySQL连接池配置
		sqlDB.SetMaxOpenConns(config.Cfg.Database.MySQL.MaxOpenConns)
		sqlDB.SetMaxIdleConns(config.Cfg.Database.MySQL.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(config.Cfg.Database.MySQL.ConnMaxLifetime) * time.Second)
		sqlDB.SetConnMaxIdleTime(time.Duration(config.Cfg.Database.MySQL.ConnMaxIdleTime) * time.Second)
	} else if config.Cfg.Database.Type == "oracle" {
		// Oracle连接池配置
		sqlDB.SetMaxOpenConns(config.Cfg.Database.Oracle.MaxOpenConns)
		sqlDB.SetMaxIdleConns(config.Cfg.Database.Oracle.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(config.Cfg.Database.Oracle.ConnMaxLifetime) * time.Second)
		sqlDB.SetConnMaxIdleTime(time.Duration(config.Cfg.Database.Oracle.ConnMaxIdleTime) * time.Second)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping数据库失败: %w", err)
	}

	DB = db

	infraLogger.Info("数据库连接成功", zap.String("type", config.Cfg.Database.Type))
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

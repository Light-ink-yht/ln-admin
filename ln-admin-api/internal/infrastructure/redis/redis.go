package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Client *redis.Client

// Init 初始化Redis连接
func Init() error {
	cfg := config.Cfg.Redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	Client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis连接失败: %w", err)
	}

	logger.Info("Redis连接成功", zap.String("addr", addr))
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

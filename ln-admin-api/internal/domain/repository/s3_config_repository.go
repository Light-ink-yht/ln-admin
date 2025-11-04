package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// S3ConfigRepository S3配置仓库接口
type S3ConfigRepository interface {
	// GetActiveConfig 获取当前启用的S3配置
	GetActiveConfig(ctx context.Context) (*entity.S3Config, error)

	// Create 创建S3配置
	Create(ctx context.Context, config *entity.S3Config) error

	// Update 更新S3配置
	Update(ctx context.Context, config *entity.S3Config) error

	// Delete 删除S3配置（软删除）
	Delete(ctx context.Context, configID string) error

	// FindByID 根据配置ID查找
	FindByID(ctx context.Context, configID string) (*entity.S3Config, error)
}

package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// SystemConfigRepository 系统配置仓库接口
type SystemConfigRepository interface {
	// GetByKey 根据配置键获取配置
	GetByKey(ctx context.Context, key string) (*entity.SystemConfig, error)

	// GetByGroup 根据分组获取配置列表
	GetByGroup(ctx context.Context, group string) ([]*entity.SystemConfig, error)

	// Create 创建配置
	Create(ctx context.Context, config *entity.SystemConfig) error

	// Update 更新配置
	Update(ctx context.Context, config *entity.SystemConfig) error

	// GetAll 获取所有配置
	GetAll(ctx context.Context) ([]*entity.SystemConfig, error)
}

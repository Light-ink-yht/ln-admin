package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// WorkbenchRepository 工作台配置仓库接口
type WorkbenchRepository interface {
	// FindByRoleKey 根据角色标识查找工作台配置
	FindByRoleKey(ctx context.Context, roleKey string) (*entity.WorkbenchConfig, error)

	// FindByConfigID 根据配置ID查找工作台配置
	FindByConfigID(ctx context.Context, configID string) (*entity.WorkbenchConfig, error)

	// Create 创建工作台配置
	Create(ctx context.Context, config *entity.WorkbenchConfig) error

	// Update 更新工作台配置
	Update(ctx context.Context, config *entity.WorkbenchConfig) error

	// Delete 删除工作台配置
	Delete(ctx context.Context, configID string) error

	// ListAll 列出所有工作台配置
	ListAll(ctx context.Context) ([]*entity.WorkbenchConfig, error)

	// ExistsByRoleKey 检查角色标识是否存在配置
	ExistsByRoleKey(ctx context.Context, roleKey string) (bool, error)
}

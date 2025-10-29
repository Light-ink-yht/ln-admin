package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// RoleRepository 角色仓库接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) error
	// FindByID 根据ID查找角色
	FindByID(ctx context.Context, roleID string) (*entity.Role, error)
	// FindByKey 根据角色标识查找角色
	FindByKey(ctx context.Context, roleKey string) (*entity.Role, error)
	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) error
	// Delete 删除角色
	Delete(ctx context.Context, roleID string) error
	// List 查询角色列表
	List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.Role, int64, error)
	// FindRolesByUserID 根据用户ID查找角色列表
	FindRolesByUserID(ctx context.Context, userID string) ([]*entity.Role, error)
}

package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// PermissionRepository 权限仓库接口
type PermissionRepository interface {
	// Create 创建权限
	Create(ctx context.Context, permission *entity.Permission) error
	// FindByID 根据ID查找权限
	FindByID(ctx context.Context, permissionID string) (*entity.Permission, error)
	// FindByKey 根据权限标识查找权限
	FindByKey(ctx context.Context, permissionKey string) (*entity.Permission, error)
	// Update 更新权限
	Update(ctx context.Context, permission *entity.Permission) error
	// Delete 删除权限
	Delete(ctx context.Context, permissionID string) error
	// List 查询权限列表
	List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.Permission, int64, error)
	// FindByResource 根据资源路径和方法查找权限
	FindByResource(ctx context.Context, resourcePath, method string) (*entity.Permission, error)
}

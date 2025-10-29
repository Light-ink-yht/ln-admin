package repository

import (
	"context"
)

// UserRoleRepository 用户角色关联仓库接口
type UserRoleRepository interface {
	// AssignRole 为用户分配角色
	AssignRole(ctx context.Context, userID, roleID string) error
	// RemoveRole 移除用户角色
	RemoveRole(ctx context.Context, userID, roleID string) error
	// FindRolesByUserID 根据用户ID查找角色ID列表
	FindRolesByUserID(ctx context.Context, userID string) ([]string, error)
	// FindUsersByRoleID 根据角色ID查找用户ID列表
	FindUsersByRoleID(ctx context.Context, roleID string) ([]string, error)
	// RemoveAllRoles 移除用户的所有角色
	RemoveAllRoles(ctx context.Context, userID string) error
}

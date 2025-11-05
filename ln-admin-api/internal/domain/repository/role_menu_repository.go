package repository

import (
	"context"
)

// RoleMenuRepository 角色菜单关联仓库接口
type RoleMenuRepository interface {
	// AssignMenu 为角色分配菜单
	AssignMenu(ctx context.Context, roleID, menuID string) error
	// RemoveMenu 移除角色菜单
	RemoveMenu(ctx context.Context, roleID, menuID string) error
	// FindMenusByRoleID 根据角色ID查找菜单ID列表
	FindMenusByRoleID(ctx context.Context, roleID string) ([]string, error)
	// FindRolesByMenuID 根据菜单ID查找角色ID列表
	FindRolesByMenuID(ctx context.Context, menuID string) ([]string, error)
	// DeleteByRoleID 删除角色的所有菜单关联
	DeleteByRoleID(ctx context.Context, roleID string) error
	// DeleteByMenuID 删除菜单的所有角色关联
	DeleteByMenuID(ctx context.Context, menuID string) error
}

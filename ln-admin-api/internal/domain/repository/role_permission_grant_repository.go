package repository

import (
	"context"
)

// RolePermissionGrantRepository 角色权限授予关系仓库接口
type RolePermissionGrantRepository interface {
	// GrantPermission 授予权限（记录授予关系）
	GrantPermission(ctx context.Context, grantorRoleID, granteeRoleID, permissionID string) error
	// RevokePermission 撤销权限（删除授予关系）
	RevokePermission(ctx context.Context, grantorRoleID, granteeRoleID, permissionID string) error
	// FindPermissionsByGranteeRoleID 根据被授予者角色ID查找权限ID列表（查找该角色被授予的所有权限）
	FindPermissionsByGranteeRoleID(ctx context.Context, granteeRoleID string) ([]string, error)
	// FindGrantorRolesByGranteeRoleID 根据被授予者角色ID查找授予者角色ID列表（查找谁授予了该角色权限）
	FindGrantorRolesByGranteeRoleID(ctx context.Context, granteeRoleID string) ([]string, error)
	// FindGranteeRolesByGrantorRoleID 根据授予者角色ID查找被授予者角色ID列表（查找该角色授予了哪些角色）
	FindGranteeRolesByGrantorRoleID(ctx context.Context, grantorRoleID string) ([]string, error)
	// CheckGrantExists 检查授予关系是否存在
	CheckGrantExists(ctx context.Context, grantorRoleID, granteeRoleID, permissionID string) (bool, error)
	// DeleteByGranteeRoleID 删除角色的所有被授予权限关系
	DeleteByGranteeRoleID(ctx context.Context, granteeRoleID string) error
	// DeleteByGrantorRoleID 删除角色的所有授予权限关系
	DeleteByGrantorRoleID(ctx context.Context, grantorRoleID string) error
	// DeleteByPermissionID 删除权限的所有授予关系
	DeleteByPermissionID(ctx context.Context, permissionID string) error
}

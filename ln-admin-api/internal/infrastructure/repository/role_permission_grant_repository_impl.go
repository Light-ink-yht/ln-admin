package repository

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.RolePermissionGrantRepository = (*rolePermissionGrantRepositoryImpl)(nil)

type rolePermissionGrantRepositoryImpl struct {
	db *gorm.DB
}

// NewRolePermissionGrantRepository 创建角色权限授予关系仓库实例
func NewRolePermissionGrantRepository() repository.RolePermissionGrantRepository {
	return &rolePermissionGrantRepositoryImpl{
		db: database.GetDB(),
	}
}

// GrantPermission 授予权限（记录授予关系）
func (r *rolePermissionGrantRepositoryImpl) GrantPermission(ctx context.Context, grantorRoleID, granteeRoleID, permissionID string) error {
	grant := &entity.RolePermissionGrant{
		GrantorRoleID: grantorRoleID,
		GranteeRoleID: granteeRoleID,
		PermissionID:  permissionID,
	}
	aa13 := grant.ToAA13()
	if err := r.db.WithContext(ctx).Create(aa13).Error; err != nil {
		// 如果已经存在，忽略错误
		if isDuplicateKeyError(err) {
			return nil
		}
		return fmt.Errorf("授予权限失败: %w", err)
	}
	return nil
}

// RevokePermission 撤销权限（删除授予关系）
func (r *rolePermissionGrantRepositoryImpl) RevokePermission(ctx context.Context, grantorRoleID, granteeRoleID, permissionID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAM001 = ? AND AAM002 = ? AND AAM003 = ?", grantorRoleID, granteeRoleID, permissionID).
		Delete(&entity.AA13{}).Error; err != nil {
		return fmt.Errorf("撤销权限失败: %w", err)
	}
	return nil
}

// FindPermissionsByGranteeRoleID 根据被授予者角色ID查找权限ID列表（查找该角色被授予的所有权限）
func (r *rolePermissionGrantRepositoryImpl) FindPermissionsByGranteeRoleID(ctx context.Context, granteeRoleID string) ([]string, error) {
	var permissionIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA13{}).
		Where("AAM002 = ?", granteeRoleID).
		Distinct().
		Pluck("AAM003", &permissionIDs).Error; err != nil {
		return nil, fmt.Errorf("查询角色被授予的权限失败: %w", err)
	}
	return permissionIDs, nil
}

// FindGrantorRolesByGranteeRoleID 根据被授予者角色ID查找授予者角色ID列表（查找谁授予了该角色权限）
func (r *rolePermissionGrantRepositoryImpl) FindGrantorRolesByGranteeRoleID(ctx context.Context, granteeRoleID string) ([]string, error) {
	var grantorRoleIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA13{}).
		Where("AAM002 = ?", granteeRoleID).
		Distinct().
		Pluck("AAM001", &grantorRoleIDs).Error; err != nil {
		return nil, fmt.Errorf("查询授予者角色失败: %w", err)
	}
	return grantorRoleIDs, nil
}

// FindGranteeRolesByGrantorRoleID 根据授予者角色ID查找被授予者角色ID列表（查找该角色授予了哪些角色）
func (r *rolePermissionGrantRepositoryImpl) FindGranteeRolesByGrantorRoleID(ctx context.Context, grantorRoleID string) ([]string, error) {
	var granteeRoleIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA13{}).
		Where("AAM001 = ?", grantorRoleID).
		Distinct().
		Pluck("AAM002", &granteeRoleIDs).Error; err != nil {
		return nil, fmt.Errorf("查询被授予者角色失败: %w", err)
	}
	return granteeRoleIDs, nil
}

// CheckGrantExists 检查授予关系是否存在
func (r *rolePermissionGrantRepositoryImpl) CheckGrantExists(ctx context.Context, grantorRoleID, granteeRoleID, permissionID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.AA13{}).
		Where("AAM001 = ? AND AAM002 = ? AND AAM003 = ?", grantorRoleID, granteeRoleID, permissionID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查授予关系失败: %w", err)
	}
	return count > 0, nil
}

// DeleteByGranteeRoleID 删除角色的所有被授予权限关系
func (r *rolePermissionGrantRepositoryImpl) DeleteByGranteeRoleID(ctx context.Context, granteeRoleID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAM002 = ?", granteeRoleID).
		Delete(&entity.AA13{}).Error; err != nil {
		return fmt.Errorf("删除角色被授予权限关系失败: %w", err)
	}
	return nil
}

// DeleteByGrantorRoleID 删除角色的所有授予权限关系
func (r *rolePermissionGrantRepositoryImpl) DeleteByGrantorRoleID(ctx context.Context, grantorRoleID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAM001 = ?", grantorRoleID).
		Delete(&entity.AA13{}).Error; err != nil {
		return fmt.Errorf("删除角色授予权限关系失败: %w", err)
	}
	return nil
}

// DeleteByPermissionID 删除权限的所有授予关系
func (r *rolePermissionGrantRepositoryImpl) DeleteByPermissionID(ctx context.Context, permissionID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAM003 = ?", permissionID).
		Delete(&entity.AA13{}).Error; err != nil {
		return fmt.Errorf("删除权限授予关系失败: %w", err)
	}
	return nil
}

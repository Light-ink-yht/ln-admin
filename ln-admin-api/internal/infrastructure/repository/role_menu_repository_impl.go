package repository

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.RoleMenuRepository = (*roleMenuRepositoryImpl)(nil)

type roleMenuRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleMenuRepository 创建角色菜单关联仓库实例
func NewRoleMenuRepository() repository.RoleMenuRepository {
	return &roleMenuRepositoryImpl{
		db: database.GetDB(),
	}
}

// AssignMenu 为角色分配菜单
func (r *roleMenuRepositoryImpl) AssignMenu(ctx context.Context, roleID, menuID string) error {
	roleMenu := &entity.RoleMenu{
		RoleID: roleID,
		MenuID: menuID,
	}
	aa12 := roleMenu.ToAA12()
	if err := r.db.WithContext(ctx).Create(aa12).Error; err != nil {
		// 如果已经存在，忽略错误
		if isDuplicateKeyError(err) {
			return nil
		}
		return fmt.Errorf("分配菜单失败: %w", err)
	}
	return nil
}

// RemoveMenu 移除角色菜单
func (r *roleMenuRepositoryImpl) RemoveMenu(ctx context.Context, roleID, menuID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAL001 = ? AND AAL002 = ?", roleID, menuID).
		Delete(&entity.AA12{}).Error; err != nil {
		return fmt.Errorf("移除菜单失败: %w", err)
	}
	return nil
}

// FindMenusByRoleID 根据角色ID查找菜单ID列表
func (r *roleMenuRepositoryImpl) FindMenusByRoleID(ctx context.Context, roleID string) ([]string, error) {
	var menuIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA12{}).
		Where("AAL001 = ?", roleID).
		Pluck("AAL002", &menuIDs).Error; err != nil {
		return nil, fmt.Errorf("查询角色菜单失败: %w", err)
	}
	return menuIDs, nil
}

// FindRolesByMenuID 根据菜单ID查找角色ID列表
func (r *roleMenuRepositoryImpl) FindRolesByMenuID(ctx context.Context, menuID string) ([]string, error) {
	var roleIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA12{}).
		Where("AAL002 = ?", menuID).
		Pluck("AAL001", &roleIDs).Error; err != nil {
		return nil, fmt.Errorf("查询菜单角色失败: %w", err)
	}
	return roleIDs, nil
}

// DeleteByRoleID 删除角色的所有菜单关联
func (r *roleMenuRepositoryImpl) DeleteByRoleID(ctx context.Context, roleID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAL001 = ?", roleID).
		Delete(&entity.AA12{}).Error; err != nil {
		return fmt.Errorf("删除角色菜单关联失败: %w", err)
	}
	return nil
}

// DeleteByMenuID 删除菜单的所有角色关联
func (r *roleMenuRepositoryImpl) DeleteByMenuID(ctx context.Context, menuID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAL002 = ?", menuID).
		Delete(&entity.AA12{}).Error; err != nil {
		return fmt.Errorf("删除菜单角色关联失败: %w", err)
	}
	return nil
}

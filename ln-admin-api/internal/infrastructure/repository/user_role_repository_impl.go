package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.UserRoleRepository = (*userRoleRepositoryImpl)(nil)

type userRoleRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRoleRepository 创建用户角色关联仓库实例
func NewUserRoleRepository() repository.UserRoleRepository {
	return &userRoleRepositoryImpl{
		db: database.GetDB(),
	}
}

// AssignRole 为用户分配角色
func (r *userRoleRepositoryImpl) AssignRole(ctx context.Context, userID, roleID string) error {
	userRole := &entity.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	aa07 := userRole.ToAA07()
	if err := r.db.WithContext(ctx).Create(aa07).Error; err != nil {
		// 如果已经存在，忽略错误
		if isDuplicateKeyError(err) {
			return nil
		}
		return fmt.Errorf("分配角色失败: %w", err)
	}
	return nil
}

// RemoveRole 移除用户角色
func (r *userRoleRepositoryImpl) RemoveRole(ctx context.Context, userID, roleID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAG001 = ? AND AAG002 = ?", userID, roleID).
		Delete(&entity.AA07{}).Error; err != nil {
		return fmt.Errorf("移除角色失败: %w", err)
	}
	return nil
}

// FindRolesByUserID 根据用户ID查找角色ID列表
func (r *userRoleRepositoryImpl) FindRolesByUserID(ctx context.Context, userID string) ([]string, error) {
	var roleIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA07{}).
		Where("AAG001 = ?", userID).
		Pluck("AAG002", &roleIDs).Error; err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}
	return roleIDs, nil
}

// FindUsersByRoleID 根据角色ID查找用户ID列表
func (r *userRoleRepositoryImpl) FindUsersByRoleID(ctx context.Context, roleID string) ([]string, error) {
	var userIDs []string
	if err := r.db.WithContext(ctx).
		Model(&entity.AA07{}).
		Where("AAG002 = ?", roleID).
		Pluck("AAG001", &userIDs).Error; err != nil {
		return nil, fmt.Errorf("查询角色用户失败: %w", err)
	}
	return userIDs, nil
}

// RemoveAllRoles 移除用户的所有角色
func (r *userRoleRepositoryImpl) RemoveAllRoles(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAG001 = ?", userID).
		Delete(&entity.AA07{}).Error; err != nil {
		return fmt.Errorf("移除所有角色失败: %w", err)
	}
	return nil
}

// isDuplicateKeyError 检查是否是重复键错误
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MySQL duplicate key error format: "Error 1062: Duplicate entry..."
	return strings.Contains(errStr, "Duplicate entry") ||
		strings.Contains(errStr, "duplicated key") ||
		strings.Contains(errStr, "1062")
}

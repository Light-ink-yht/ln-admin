package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.RoleRepository = (*roleRepositoryImpl)(nil)

type roleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库实例
func NewRoleRepository() repository.RoleRepository {
	return &roleRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建角色
func (r *roleRepositoryImpl) Create(ctx context.Context, role *entity.Role) error {
	aa05 := role.ToAA05()
	if err := r.db.WithContext(ctx).Create(aa05).Error; err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}
	role.FromAA05(aa05)
	return nil
}

// FindByID 根据ID查找角色
func (r *roleRepositoryImpl) FindByID(ctx context.Context, roleID string) (*entity.Role, error) {
	var aa05 entity.AA05
	if err := r.db.WithContext(ctx).Where("AAE001 = ?", roleID).First(&aa05).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}
	role := &entity.Role{}
	role.FromAA05(&aa05)
	return role, nil
}

// FindByKey 根据角色标识查找角色
func (r *roleRepositoryImpl) FindByKey(ctx context.Context, roleKey string) (*entity.Role, error) {
	var aa05 entity.AA05
	if err := r.db.WithContext(ctx).Where("AAE002 = ?", roleKey).First(&aa05).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}
	role := &entity.Role{}
	role.FromAA05(&aa05)
	return role, nil
}

// Update 更新角色
func (r *roleRepositoryImpl) Update(ctx context.Context, role *entity.Role) error {
	aa05 := role.ToAA05()
	if err := r.db.WithContext(ctx).Where("AAE001 = ?", aa05.AAE001).Updates(aa05).Error; err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}
	return nil
}

// Delete 删除角色
func (r *roleRepositoryImpl) Delete(ctx context.Context, roleID string) error {
	if err := r.db.WithContext(ctx).Where("AAE001 = ?", roleID).Delete(&entity.AA05{}).Error; err != nil {
		return fmt.Errorf("删除角色失败: %w", err)
	}
	return nil
}

// List 查询角色列表
func (r *roleRepositoryImpl) List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.Role, int64, error) {
	var aa05List []entity.AA05
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AA05{})

	// 应用查询条件
	for key, value := range conditions {
		switch key {
		case "role_key":
			// 支持数组格式（多个搜索词，使用 OR 查询）
			if values, ok := value.([]interface{}); ok && len(values) > 0 {
				var orConditions []string
				var args []interface{}
				for _, v := range values {
					orConditions = append(orConditions, "AAE002 LIKE ?")
					args = append(args, "%"+fmt.Sprintf("%v", v)+"%")
				}
				query = query.Where("("+strings.Join(orConditions, " OR ")+")", args...)
			} else {
				query = query.Where("AAE002 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
			}
		case "role_name":
			if values, ok := value.([]interface{}); ok && len(values) > 0 {
				var orConditions []string
				var args []interface{}
				for _, v := range values {
					orConditions = append(orConditions, "AAE003 LIKE ?")
					args = append(args, "%"+fmt.Sprintf("%v", v)+"%")
				}
				query = query.Where("("+strings.Join(orConditions, " OR ")+")", args...)
			} else {
				query = query.Where("AAE003 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
			}
		case "status":
			query = query.Where("AAE005 = ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询角色总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&aa05List).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询角色失败: %w", err)
	}

	roles := make([]*entity.Role, len(aa05List))
	for i, aa05 := range aa05List {
		role := &entity.Role{}
		role.FromAA05(&aa05)
		roles[i] = role
	}

	return roles, total, nil
}

// FindRolesByUserID 根据用户ID查找角色列表
func (r *roleRepositoryImpl) FindRolesByUserID(ctx context.Context, userID string) ([]*entity.Role, error) {
	var aa05List []entity.AA05
	if err := r.db.WithContext(ctx).
		Table("AA05").
		Joins("INNER JOIN AA07 ON AA05.AAE001 = AA07.AAG002").
		Where("AA07.AAG001 = ?", userID).
		Find(&aa05List).Error; err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}

	roles := make([]*entity.Role, len(aa05List))
	for i, aa05 := range aa05List {
		role := &entity.Role{}
		role.FromAA05(&aa05)
		roles[i] = role
	}

	return roles, nil
}

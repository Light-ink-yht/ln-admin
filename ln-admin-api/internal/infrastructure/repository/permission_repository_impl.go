package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.PermissionRepository = (*permissionRepositoryImpl)(nil)

type permissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库实例
func NewPermissionRepository() repository.PermissionRepository {
	return &permissionRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建权限
func (r *permissionRepositoryImpl) Create(ctx context.Context, permission *entity.Permission) error {
	aa06 := permission.ToAA06()
	if err := r.db.WithContext(ctx).Create(aa06).Error; err != nil {
		return fmt.Errorf("创建权限失败: %w", err)
	}
	permission.FromAA06(aa06)
	return nil
}

// FindByID 根据ID查找权限
func (r *permissionRepositoryImpl) FindByID(ctx context.Context, permissionID string) (*entity.Permission, error) {
	var aa06 entity.AA06
	if err := r.db.WithContext(ctx).Where("AAF001 = ?", permissionID).First(&aa06).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询权限失败: %w", err)
	}
	permission := &entity.Permission{}
	permission.FromAA06(&aa06)
	return permission, nil
}

// FindByKey 根据权限标识查找权限
func (r *permissionRepositoryImpl) FindByKey(ctx context.Context, permissionKey string) (*entity.Permission, error) {
	var aa06 entity.AA06
	if err := r.db.WithContext(ctx).Where("AAF002 = ?", permissionKey).First(&aa06).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询权限失败: %w", err)
	}
	permission := &entity.Permission{}
	permission.FromAA06(&aa06)
	return permission, nil
}

// Update 更新权限
func (r *permissionRepositoryImpl) Update(ctx context.Context, permission *entity.Permission) error {
	aa06 := permission.ToAA06()
	if err := r.db.WithContext(ctx).Where("AAF001 = ?", aa06.AAF001).Updates(aa06).Error; err != nil {
		return fmt.Errorf("更新权限失败: %w", err)
	}
	return nil
}

// Delete 删除权限
func (r *permissionRepositoryImpl) Delete(ctx context.Context, permissionID string) error {
	if err := r.db.WithContext(ctx).Where("AAF001 = ?", permissionID).Delete(&entity.AA06{}).Error; err != nil {
		return fmt.Errorf("删除权限失败: %w", err)
	}
	return nil
}

// List 查询权限列表
func (r *permissionRepositoryImpl) List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.Permission, int64, error) {
	var aa06List []entity.AA06
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AA06{})

	// 应用查询条件
	for key, value := range conditions {
		switch key {
		case "permission_key":
			query = query.Where("AAF002 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "permission_name":
			query = query.Where("AAF003 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "resource_path":
			query = query.Where("AAF004 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "method":
			query = query.Where("AAF005 = ?", value)
		case "status":
			query = query.Where("AAF007 = ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询权限总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&aa06List).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询权限失败: %w", err)
	}

	permissions := make([]*entity.Permission, len(aa06List))
	for i, aa06 := range aa06List {
		permission := &entity.Permission{}
		permission.FromAA06(&aa06)
		permissions[i] = permission
	}

	return permissions, total, nil
}

// FindByResource 根据资源路径和方法查找权限
func (r *permissionRepositoryImpl) FindByResource(ctx context.Context, resourcePath, method string) (*entity.Permission, error) {
	var aa06 entity.AA06
	if err := r.db.WithContext(ctx).
		Where("AAF004 = ? AND AAF005 = ?", resourcePath, method).
		First(&aa06).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询权限失败: %w", err)
	}
	permission := &entity.Permission{}
	permission.FromAA06(&aa06)
	return permission, nil
}

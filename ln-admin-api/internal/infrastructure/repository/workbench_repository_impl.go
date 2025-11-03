package repository

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.WorkbenchRepository = (*workbenchRepositoryImpl)(nil)

type workbenchRepositoryImpl struct {
	db *gorm.DB
}

// NewWorkbenchRepository 创建工作台配置仓库实例
func NewWorkbenchRepository() repository.WorkbenchRepository {
	return &workbenchRepositoryImpl{
		db: database.GetDB(),
	}
}

// FindByRoleKey 根据角色标识查找工作台配置
func (r *workbenchRepositoryImpl) FindByRoleKey(ctx context.Context, roleKey string) (*entity.WorkbenchConfig, error) {
	var aa10 entity.AA10
	if err := r.db.WithContext(ctx).
		Where("AAJ002 = ? AND AAJ005 = ?", roleKey, "1"). // 只查询启用的配置
		First(&aa10).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询工作台配置失败: %w", err)
	}

	config := &entity.WorkbenchConfig{}
	config.FromAA10(&aa10)
	return config, nil
}

// FindByConfigID 根据配置ID查找工作台配置
func (r *workbenchRepositoryImpl) FindByConfigID(ctx context.Context, configID string) (*entity.WorkbenchConfig, error) {
	var aa10 entity.AA10
	if err := r.db.WithContext(ctx).
		Where("AAJ001 = ?", configID).
		First(&aa10).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询工作台配置失败: %w", err)
	}

	config := &entity.WorkbenchConfig{}
	config.FromAA10(&aa10)
	return config, nil
}

// Create 创建工作台配置
func (r *workbenchRepositoryImpl) Create(ctx context.Context, config *entity.WorkbenchConfig) error {
	aa10 := config.ToAA10()
	if err := r.db.WithContext(ctx).Create(aa10).Error; err != nil {
		return fmt.Errorf("创建工作台配置失败: %w", err)
	}
	return nil
}

// Update 更新工作台配置
func (r *workbenchRepositoryImpl) Update(ctx context.Context, config *entity.WorkbenchConfig) error {
	aa10 := config.ToAA10()
	if err := r.db.WithContext(ctx).
		Where("AAJ001 = ?", config.ConfigID).
		Updates(aa10).Error; err != nil {
		return fmt.Errorf("更新工作台配置失败: %w", err)
	}
	return nil
}

// Delete 删除工作台配置
func (r *workbenchRepositoryImpl) Delete(ctx context.Context, configID string) error {
	if err := r.db.WithContext(ctx).
		Where("AAJ001 = ?", configID).
		Delete(&entity.AA10{}).Error; err != nil {
		return fmt.Errorf("删除工作台配置失败: %w", err)
	}
	return nil
}

// ListAll 列出所有工作台配置
func (r *workbenchRepositoryImpl) ListAll(ctx context.Context) ([]*entity.WorkbenchConfig, error) {
	var aa10List []entity.AA10
	if err := r.db.WithContext(ctx).Find(&aa10List).Error; err != nil {
		return nil, fmt.Errorf("查询工作台配置列表失败: %w", err)
	}

	configs := make([]*entity.WorkbenchConfig, len(aa10List))
	for i, aa10 := range aa10List {
		config := &entity.WorkbenchConfig{}
		config.FromAA10(&aa10)
		configs[i] = config
	}

	return configs, nil
}

// ExistsByRoleKey 检查角色标识是否存在配置
func (r *workbenchRepositoryImpl) ExistsByRoleKey(ctx context.Context, roleKey string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.AA10{}).
		Where("AAJ002 = ?", roleKey).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查工作台配置是否存在失败: %w", err)
	}
	return count > 0, nil
}

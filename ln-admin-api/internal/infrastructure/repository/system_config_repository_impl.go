package repository

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"

	"gorm.io/gorm"
)

var _ repository.SystemConfigRepository = (*systemConfigRepositoryImpl)(nil)

type systemConfigRepositoryImpl struct {
	db *gorm.DB
}

// NewSystemConfigRepository 创建系统配置仓库实例
func NewSystemConfigRepository() repository.SystemConfigRepository {
	return &systemConfigRepositoryImpl{
		db: database.GetDB(),
	}
}

// GetByKey 根据配置键获取配置
func (r *systemConfigRepositoryImpl) GetByKey(ctx context.Context, key string) (*entity.SystemConfig, error) {
	var aa02 entity.AA02
	if err := r.db.WithContext(ctx).Where("AAB001 = ? AND AAB006 = ?", key, "1").First(&aa02).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询配置失败: %w", err)
	}
	config := &entity.SystemConfig{}
	config.FromAA02(&aa02)
	return config, nil
}

// GetByGroup 根据分组获取配置列表
func (r *systemConfigRepositoryImpl) GetByGroup(ctx context.Context, group string) ([]*entity.SystemConfig, error) {
	var aa02List []entity.AA02
	if err := r.db.WithContext(ctx).Where("AAB004 = ? AND AAB006 = ?", group, "1").Find(&aa02List).Error; err != nil {
		return nil, fmt.Errorf("查询配置列表失败: %w", err)
	}
	configs := make([]*entity.SystemConfig, len(aa02List))
	for i, aa02 := range aa02List {
		config := &entity.SystemConfig{}
		config.FromAA02(&aa02)
		configs[i] = config
	}
	return configs, nil
}

// Create 创建配置
func (r *systemConfigRepositoryImpl) Create(ctx context.Context, config *entity.SystemConfig) error {
	aa02 := config.ToAA02()
	if err := r.db.WithContext(ctx).Create(aa02).Error; err != nil {
		return fmt.Errorf("创建配置失败: %w", err)
	}
	config.FromAA02(aa02)
	return nil
}

// Update 更新配置
func (r *systemConfigRepositoryImpl) Update(ctx context.Context, config *entity.SystemConfig) error {
	aa02 := config.ToAA02()
	if err := r.db.WithContext(ctx).Model(&entity.AA02{}).
		Where("AAB001 = ?", config.ConfigKey).
		Updates(aa02).Error; err != nil {
		return fmt.Errorf("更新配置失败: %w", err)
	}
	return nil
}

// GetAll 获取所有配置
func (r *systemConfigRepositoryImpl) GetAll(ctx context.Context) ([]*entity.SystemConfig, error) {
	var aa02List []entity.AA02
	if err := r.db.WithContext(ctx).Where("AAB006 = ?", "1").Find(&aa02List).Error; err != nil {
		return nil, fmt.Errorf("查询配置列表失败: %w", err)
	}
	configs := make([]*entity.SystemConfig, len(aa02List))
	for i, aa02 := range aa02List {
		config := &entity.SystemConfig{}
		config.FromAA02(&aa02)
		configs[i] = config
	}
	return configs, nil
}

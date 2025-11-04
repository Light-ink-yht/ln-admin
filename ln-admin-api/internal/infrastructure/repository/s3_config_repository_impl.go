package repository

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"

	"gorm.io/gorm"
)

var _ repository.S3ConfigRepository = (*s3ConfigRepositoryImpl)(nil)

type s3ConfigRepositoryImpl struct {
	db *gorm.DB
}

// NewS3ConfigRepository 创建S3配置仓库实例
func NewS3ConfigRepository() repository.S3ConfigRepository {
	return &s3ConfigRepositoryImpl{
		db: database.GetDB(),
	}
}

// GetActiveConfig 获取当前启用的S3配置
func (r *s3ConfigRepositoryImpl) GetActiveConfig(ctx context.Context) (*entity.S3Config, error) {
	var aa11 entity.AA11
	if err := r.db.WithContext(ctx).
		Where("AAK009 = ?", "1").  // 状态为启用
		Where("AAK002 = ?", "s3"). // 存储类型为s3
		First(&aa11).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询S3配置失败: %w", err)
	}
	config := &entity.S3Config{}
	config.FromAA11(&aa11)
	return config, nil
}

// Create 创建S3配置
func (r *s3ConfigRepositoryImpl) Create(ctx context.Context, config *entity.S3Config) error {
	aa11 := config.ToAA11()
	if err := r.db.WithContext(ctx).Create(aa11).Error; err != nil {
		return fmt.Errorf("创建S3配置失败: %w", err)
	}
	config.FromAA11(aa11)
	return nil
}

// Update 更新S3配置
func (r *s3ConfigRepositoryImpl) Update(ctx context.Context, config *entity.S3Config) error {
	aa11 := config.ToAA11()
	if err := r.db.WithContext(ctx).Model(&entity.AA11{}).
		Where("AAK001 = ?", config.ConfigID).
		Updates(aa11).Error; err != nil {
		return fmt.Errorf("更新S3配置失败: %w", err)
	}
	return nil
}

// Delete 删除S3配置（软删除）
func (r *s3ConfigRepositoryImpl) Delete(ctx context.Context, configID string) error {
	if err := r.db.WithContext(ctx).Where("AAK001 = ?", configID).
		Delete(&entity.AA11{}).Error; err != nil {
		return fmt.Errorf("删除S3配置失败: %w", err)
	}
	return nil
}

// FindByID 根据配置ID查找
func (r *s3ConfigRepositoryImpl) FindByID(ctx context.Context, configID string) (*entity.S3Config, error) {
	var aa11 entity.AA11
	if err := r.db.WithContext(ctx).Where("AAK001 = ?", configID).First(&aa11).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询S3配置失败: %w", err)
	}
	config := &entity.S3Config{}
	config.FromAA11(&aa11)
	return config, nil
}

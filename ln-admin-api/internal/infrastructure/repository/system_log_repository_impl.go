package repository

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"

	"gorm.io/gorm"
)

var _ repository.SystemLogRepository = (*systemLogRepositoryImpl)(nil)

type systemLogRepositoryImpl struct {
	db *gorm.DB
}

// NewSystemLogRepository 创建系统日志仓库实例
func NewSystemLogRepository() repository.SystemLogRepository {
	return &systemLogRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建日志
func (r *systemLogRepositoryImpl) Create(ctx context.Context, log *entity.SystemLog) error {
	aa03 := log.ToAA03()
	if err := r.db.WithContext(ctx).Create(aa03).Error; err != nil {
		return fmt.Errorf("创建日志失败: %w", err)
	}
	log.FromAA03(aa03)
	return nil
}

// List 分页查询日志
func (r *systemLogRepositoryImpl) List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.SystemLog, int64, error) {
	var aa03List []entity.AA03
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AA03{})

	// 应用查询条件
	for key, value := range conditions {
		switch key {
		case "level":
			query = query.Where("AAC002 = ?", value)
		case "module":
			query = query.Where("AAC003 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "action":
			query = query.Where("AAC004 = ?", value)
		case "user_id":
			query = query.Where("AAC006 = ?", value)
		case "ip":
			query = query.Where("AAC007 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "start_time":
			query = query.Where("AAC013 >= ?", value)
		case "end_time":
			query = query.Where("AAC013 <= ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询日志总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("AAC013 DESC").Find(&aa03List).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询日志失败: %w", err)
	}

	logs := make([]*entity.SystemLog, len(aa03List))
	for i, aa03 := range aa03List {
		log := &entity.SystemLog{}
		log.FromAA03(&aa03)
		logs[i] = log
	}

	return logs, total, nil
}

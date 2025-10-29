package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// SystemLogRepository 系统日志仓库接口
type SystemLogRepository interface {
	// Create 创建日志
	Create(ctx context.Context, log *entity.SystemLog) error

	// List 分页查询日志
	List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.SystemLog, int64, error)
}

package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// SMSTemplateRepository 短信模板仓库接口
type SMSTemplateRepository interface {
	// Create 创建模板
	Create(ctx context.Context, template *entity.SMSTemplate) error
	// FindByID 根据ID查找模板
	FindByID(ctx context.Context, templateID string) (*entity.SMSTemplate, error)
	// FindByType 根据类型查找模板
	FindByType(ctx context.Context, templateType string) (*entity.SMSTemplate, error)
	// Update 更新模板
	Update(ctx context.Context, template *entity.SMSTemplate) error
	// Delete 删除模板
	Delete(ctx context.Context, templateID string) error
	// List 查询模板列表
	List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.SMSTemplate, int64, error)
}

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

var _ repository.SMSTemplateRepository = (*smsTemplateRepositoryImpl)(nil)

type smsTemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewSMSTemplateRepository 创建短信模板仓库实例
func NewSMSTemplateRepository() repository.SMSTemplateRepository {
	return &smsTemplateRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建模板
func (r *smsTemplateRepositoryImpl) Create(ctx context.Context, template *entity.SMSTemplate) error {
	aa08 := template.ToAA08()
	if err := r.db.WithContext(ctx).Create(aa08).Error; err != nil {
		return fmt.Errorf("创建短信模板失败: %w", err)
	}
	template.FromAA08(aa08)
	return nil
}

// FindByID 根据ID查找模板
func (r *smsTemplateRepositoryImpl) FindByID(ctx context.Context, templateID string) (*entity.SMSTemplate, error) {
	var aa08 entity.AA08
	if err := r.db.WithContext(ctx).Where("AAH001 = ?", templateID).First(&aa08).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询短信模板失败: %w", err)
	}
	template := &entity.SMSTemplate{}
	template.FromAA08(&aa08)
	return template, nil
}

// FindByType 根据类型查找模板
func (r *smsTemplateRepositoryImpl) FindByType(ctx context.Context, templateType string) (*entity.SMSTemplate, error) {
	var aa08 entity.AA08
	if err := r.db.WithContext(ctx).
		Where("AAH002 = ? AND AAH008 = ?", templateType, "1"). // 查找启用状态的模板
		First(&aa08).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询短信模板失败: %w", err)
	}
	template := &entity.SMSTemplate{}
	template.FromAA08(&aa08)
	return template, nil
}

// Update 更新模板
func (r *smsTemplateRepositoryImpl) Update(ctx context.Context, template *entity.SMSTemplate) error {
	aa08 := template.ToAA08()
	if err := r.db.WithContext(ctx).Where("AAH001 = ?", aa08.AAH001).Updates(aa08).Error; err != nil {
		return fmt.Errorf("更新短信模板失败: %w", err)
	}
	return nil
}

// Delete 删除模板
func (r *smsTemplateRepositoryImpl) Delete(ctx context.Context, templateID string) error {
	if err := r.db.WithContext(ctx).Where("AAH001 = ?", templateID).Delete(&entity.AA08{}).Error; err != nil {
		return fmt.Errorf("删除短信模板失败: %w", err)
	}
	return nil
}

// List 查询模板列表
func (r *smsTemplateRepositoryImpl) List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.SMSTemplate, int64, error) {
	var aa08List []entity.AA08
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AA08{})

	// 应用查询条件
	for key, value := range conditions {
		switch key {
		case "type":
			query = query.Where("AAH002 = ?", value)
		case "template_name":
			query = query.Where("AAH003 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "status":
			query = query.Where("AAH008 = ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询模板总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&aa08List).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询模板失败: %w", err)
	}

	templates := make([]*entity.SMSTemplate, len(aa08List))
	for i, aa08 := range aa08List {
		template := &entity.SMSTemplate{}
		template.FromAA08(&aa08)
		templates[i] = template
	}

	return templates, total, nil
}

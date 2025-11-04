package service

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/google/uuid"
)

// SMSTemplateService 短信模板服务
type SMSTemplateService struct {
	templateRepo repository.SMSTemplateRepository
}

// NewSMSTemplateService 创建短信模板服务
func NewSMSTemplateService(templateRepo repository.SMSTemplateRepository) *SMSTemplateService {
	return &SMSTemplateService{
		templateRepo: templateRepo,
	}
}

// GetTemplateList 获取模板列表
func (s *SMSTemplateService) GetTemplateList(ctx context.Context, req *dto.SMSTemplateListRequest) ([]*dto.SMSTemplateResponse, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	conditions := make(map[string]interface{})
	if req.Type != "" {
		conditions["type"] = req.Type
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}

	templates, total, err := s.templateRepo.List(ctx, page, pageSize, conditions)
	if err != nil {
		return nil, 0, fmt.Errorf("查询模板列表失败: %w", err)
	}

	responses := make([]*dto.SMSTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = s.toTemplateResponse(template)
	}

	return responses, total, nil
}

// GetTemplateByID 根据ID获取模板
func (s *SMSTemplateService) GetTemplateByID(ctx context.Context, templateID string) (*dto.SMSTemplateResponse, error) {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}
	if template == nil {
		return nil, fmt.Errorf("模板不存在")
	}

	return s.toTemplateResponse(template), nil
}

// CreateTemplate 创建模板
func (s *SMSTemplateService) CreateTemplate(ctx context.Context, req *dto.CreateSMSTemplateRequest, userID string) error {
	// 检查类型是否已存在
	existing, err := s.templateRepo.FindByType(ctx, req.Type)
	if err != nil {
		return fmt.Errorf("查询模板失败: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("模板类型 %s 已存在", req.Type)
	}

	template := &entity.SMSTemplate{
		TemplateID:        fmt.Sprintf("template_%s_%s", req.Type, uuid.New().String()[:8]),
		Type:              req.Type,
		TemplateName:      req.TemplateName,
		TencentTemplateID: req.TencentTemplateID,
		Title:             req.Title,
		Content:           req.Content,
		Description:       req.Description,
		Status:            req.Status,
	}
	if template.Status == "" {
		template.Status = "1"
	}
	template.CreatorID = userID
	template.ModifierID = userID

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return fmt.Errorf("创建模板失败: %w", err)
	}

	return nil
}

// UpdateTemplate 更新模板
func (s *SMSTemplateService) UpdateTemplate(ctx context.Context, templateID string, req *dto.UpdateSMSTemplateRequest, userID string) error {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("查询模板失败: %w", err)
	}
	if template == nil {
		return fmt.Errorf("模板不存在")
	}

	if req.TemplateName != "" {
		template.TemplateName = req.TemplateName
	}
	if req.TencentTemplateID != "" {
		template.TencentTemplateID = req.TencentTemplateID
	}
	if req.Title != "" {
		template.Title = req.Title
	}
	if req.Content != "" {
		template.Content = req.Content
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Status != "" {
		template.Status = req.Status
	}
	template.ModifierID = userID

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return fmt.Errorf("更新模板失败: %w", err)
	}

	return nil
}

// DeleteTemplate 删除模板
func (s *SMSTemplateService) DeleteTemplate(ctx context.Context, templateID string) error {
	if err := s.templateRepo.Delete(ctx, templateID); err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}
	return nil
}

// toTemplateResponse 转换为响应DTO
func (s *SMSTemplateService) toTemplateResponse(template *entity.SMSTemplate) *dto.SMSTemplateResponse {
	return &dto.SMSTemplateResponse{
		TemplateID:        template.TemplateID,
		Type:              template.Type,
		TemplateName:      template.TemplateName,
		TencentTemplateID: template.TencentTemplateID,
		Title:             template.Title,
		Content:           template.Content,
		Description:       template.Description,
		Status:            template.Status,
		CreatorID:         template.CreatorID,
		ModifierID:        template.ModifierID,
		CreatedAt:         dto.FormatTime(template.CreatedAt),
		UpdatedAt:         dto.FormatTime(template.UpdatedAt),
	}
}

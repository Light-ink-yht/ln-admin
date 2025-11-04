package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
)

var (
	ErrSystemConfigNotFound = errors.New("系统配置不存在")
	ErrSystemConfigExists   = errors.New("系统配置已存在")
)

// SystemConfigService 系统配置服务
type SystemConfigService struct {
	configRepo repository.SystemConfigRepository
}

// NewSystemConfigService 创建系统配置服务
func NewSystemConfigService(configRepo repository.SystemConfigRepository) *SystemConfigService {
	return &SystemConfigService{
		configRepo: configRepo,
	}
}

// GetConfigByKey 根据配置键获取配置
func (s *SystemConfigService) GetConfigByKey(ctx context.Context, key string) (*dto.SystemConfigResponse, error) {
	config, err := s.configRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("查询配置失败: %w", err)
	}
	if config == nil {
		return nil, ErrSystemConfigNotFound
	}

	return s.toConfigResponse(config), nil
}

// GetConfigsByGroup 根据分组获取配置列表
func (s *SystemConfigService) GetConfigsByGroup(ctx context.Context, group string) ([]*dto.SystemConfigResponse, error) {
	configs, err := s.configRepo.GetByGroup(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("查询配置列表失败: %w", err)
	}

	responses := make([]*dto.SystemConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = s.toConfigResponse(config)
	}

	return responses, nil
}

// GetAllConfigs 获取所有配置
func (s *SystemConfigService) GetAllConfigs(ctx context.Context) ([]*dto.SystemConfigResponse, error) {
	configs, err := s.configRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询所有配置失败: %w", err)
	}

	responses := make([]*dto.SystemConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = s.toConfigResponse(config)
	}

	return responses, nil
}

// CreateConfig 创建配置
func (s *SystemConfigService) CreateConfig(ctx context.Context, req *dto.CreateSystemConfigRequest, creatorID string) error {
	// 检查配置是否已存在
	existing, err := s.configRepo.GetByKey(ctx, req.ConfigKey)
	if err != nil {
		return fmt.Errorf("检查配置是否存在失败: %w", err)
	}
	if existing != nil {
		return ErrSystemConfigExists
	}

	config := &entity.SystemConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigName:  req.ConfigName,
		ConfigGroup: req.ConfigGroup,
		Description: req.Description,
		Status:      req.Status,
		CreatorID:   creatorID,
		ModifierID:  creatorID,
	}

	if err := s.configRepo.Create(ctx, config); err != nil {
		return fmt.Errorf("创建配置失败: %w", err)
	}

	return nil
}

// UpdateConfig 更新配置
func (s *SystemConfigService) UpdateConfig(ctx context.Context, key string, req *dto.UpdateSystemConfigRequest, modifierID string) error {
	// 检查配置是否存在
	existing, err := s.configRepo.GetByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("查询配置失败: %w", err)
	}
	if existing == nil {
		return ErrSystemConfigNotFound
	}

	// 更新配置
	if req.ConfigValue != "" {
		existing.ConfigValue = req.ConfigValue
	}
	if req.ConfigName != "" {
		existing.ConfigName = req.ConfigName
	}
	if req.ConfigGroup != "" {
		existing.ConfigGroup = req.ConfigGroup
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.ModifierID = modifierID

	if err := s.configRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("更新配置失败: %w", err)
	}

	return nil
}

// toConfigResponse 转换为响应DTO
func (s *SystemConfigService) toConfigResponse(config *entity.SystemConfig) *dto.SystemConfigResponse {
	return &dto.SystemConfigResponse{
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigName:  config.ConfigName,
		ConfigGroup: config.ConfigGroup,
		Description: config.Description,
		Status:      config.Status,
		CreatedAt:   dto.FormatTime(config.CreatedAt),
		UpdatedAt:   dto.FormatTime(config.UpdatedAt),
	}
}

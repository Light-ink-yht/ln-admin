package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrWorkbenchConfigNotFound = errors.New("工作台配置不存在")
	ErrWorkbenchConfigExists   = errors.New("工作台配置已存在")
)

// WorkbenchService 工作台配置服务
type WorkbenchService struct {
	workbenchRepo     repository.WorkbenchRepository
	permissionService *PermissionService
}

// NewWorkbenchService 创建工作台配置服务
func NewWorkbenchService(workbenchRepo repository.WorkbenchRepository, permissionService *PermissionService) *WorkbenchService {
	return &WorkbenchService{
		workbenchRepo:     workbenchRepo,
		permissionService: permissionService,
	}
}

// GetWorkbenchConfig 获取用户的工作台配置
// 根据用户的角色返回对应的工作台配置
func (s *WorkbenchService) GetWorkbenchConfig(ctx context.Context, userID string) (*dto.WorkbenchConfigResponse, error) {
	// 获取用户角色
	roles, err := s.permissionService.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 如果用户没有角色，返回默认配置
	if len(roles) == 0 {
		return s.getDefaultConfig("guest"), nil
	}

	// 使用第一个角色（主要角色）查找配置
	primaryRole := roles[0]
	config, err := s.workbenchRepo.FindByRoleKey(ctx, primaryRole.RoleKey)
	if err != nil {
		return nil, fmt.Errorf("查询工作台配置失败: %w", err)
	}

	// 如果数据库中没有配置，返回默认配置
	if config == nil {
		return s.getDefaultConfig(primaryRole.RoleKey), nil
	}

	// 解析JSON配置数据
	var configData dto.WorkbenchConfigData
	if err := json.Unmarshal([]byte(config.ConfigData), &configData); err != nil {
		// 如果解析失败，返回默认配置
		return s.getDefaultConfig(primaryRole.RoleKey), nil
	}

	return &dto.WorkbenchConfigResponse{
		RoleKey:    config.RoleKey,
		ConfigName: config.ConfigName,
		Config:     configData,
	}, nil
}

// getDefaultConfig 获取默认配置（当数据库中不存在时使用）
func (s *WorkbenchService) getDefaultConfig(roleKey string) *dto.WorkbenchConfigResponse {
	// 默认配置（统一配置，实际配置应从数据库读取）
	defaultConfig := dto.WorkbenchConfigData{
		Stats: []dto.WorkbenchStatItem{
			{Key: "users", Title: "用户总数", Value: 0, Icon: "TeamOutlined", Color: "#1890ff", Trend: "up", Change: "12%", Desc: "较上月", Type: "primary", Suffix: "人"},
			{Key: "roles", Title: "角色数量", Value: 0, Icon: "SafetyOutlined", Color: "#52c41a", Trend: "up", Change: "5%", Desc: "较上月", Type: "success", Suffix: "个"},
			{Key: "permissions", Title: "权限数量", Value: 0, Icon: "UnlockOutlined", Color: "#faad14", Trend: "down", Change: "2%", Desc: "较上月", Type: "warning", Suffix: "个"},
			{Key: "logs", Title: "今日操作", Value: 0, Icon: "FileTextOutlined", Color: "#f5222d", Trend: "up", Change: "23%", Desc: "较昨日", Type: "danger", Suffix: "次"},
		},
		QuickActions: []dto.WorkbenchQuickAction{
			{Key: "user", Title: "用户管理", Desc: "管理系统用户", Icon: "UserAddOutlined", Color: "var(--color-primary)", Path: "/user/list"},
			{Key: "role", Title: "角色管理", Desc: "配置用户角色", Icon: "SafetyOutlined", Color: "var(--color-primary)", Path: "/role/list"},
			{Key: "permission", Title: "权限管理", Desc: "管理权限配置", Icon: "UnlockOutlined", Color: "var(--color-primary)", Path: "/permission/list"},
			{Key: "log", Title: "操作日志", Desc: "查看系统日志", Icon: "AuditOutlined", Color: "var(--color-primary)", Path: "/system/log"},
			{Key: "config", Title: "系统配置", Desc: "系统参数设置", Icon: "SettingOutlined", Color: "var(--color-primary)", Path: "/system/config"},
		},
		Announcements: []dto.WorkbenchAnnouncement{
			{Key: "1", Title: "欢迎使用", Content: "欢迎使用 LN Admin 管理系统", Time: "2024-01-01 08:00", Author: "系统管理员", Type: "normal"},
		},
		ShowSystemInfo: true,
	}

	return &dto.WorkbenchConfigResponse{
		RoleKey:    roleKey,
		ConfigName: "默认配置",
		Config:     defaultConfig,
	}
}

// CreateWorkbenchConfig 创建工作台配置
func (s *WorkbenchService) CreateWorkbenchConfig(ctx context.Context, req *dto.WorkbenchConfigCreateRequest, creatorID string) error {
	// 检查是否已存在
	exists, err := s.workbenchRepo.ExistsByRoleKey(ctx, req.RoleKey)
	if err != nil {
		return fmt.Errorf("检查工作台配置是否存在失败: %w", err)
	}
	if exists {
		return ErrWorkbenchConfigExists
	}

	// 将配置数据转为JSON
	configDataJSON, err := json.Marshal(req.Config)
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %w", err)
	}

	// 创建配置实体
	config := &entity.WorkbenchConfig{
		ConfigID:    uuid.New().String(),
		RoleKey:     req.RoleKey,
		ConfigName:  req.ConfigName,
		ConfigData:  string(configDataJSON),
		Status:      req.Status,
		Description: req.Description,
		CreatorID:   creatorID,
		ModifierID:  creatorID,
	}

	if config.Status == "" {
		config.Status = "1" // 默认启用
	}

	return s.workbenchRepo.Create(ctx, config)
}

// UpdateWorkbenchConfig 更新工作台配置
func (s *WorkbenchService) UpdateWorkbenchConfig(ctx context.Context, req *dto.WorkbenchConfigUpdateRequest, modifierID string) error {
	// 查找现有配置
	config, err := s.workbenchRepo.FindByConfigID(ctx, req.ConfigID)
	if err != nil {
		return fmt.Errorf("查询工作台配置失败: %w", err)
	}
	if config == nil {
		return ErrWorkbenchConfigNotFound
	}

	// 更新配置
	if req.ConfigName != "" {
		config.ConfigName = req.ConfigName
	}
	if req.Description != "" {
		config.Description = req.Description
	}
	if req.Status != "" {
		config.Status = req.Status
	}

	// 如果提供了新的配置数据，更新它
	if req.Config.Stats != nil || req.Config.QuickActions != nil || req.Config.Announcements != nil {
		configDataJSON, err := json.Marshal(req.Config)
		if err != nil {
			return fmt.Errorf("序列化配置数据失败: %w", err)
		}
		config.ConfigData = string(configDataJSON)
	}

	config.ModifierID = modifierID

	return s.workbenchRepo.Update(ctx, config)
}

// DeleteWorkbenchConfig 删除工作台配置
func (s *WorkbenchService) DeleteWorkbenchConfig(ctx context.Context, configID string) error {
	config, err := s.workbenchRepo.FindByConfigID(ctx, configID)
	if err != nil {
		return fmt.Errorf("查询工作台配置失败: %w", err)
	}
	if config == nil {
		return ErrWorkbenchConfigNotFound
	}

	return s.workbenchRepo.Delete(ctx, configID)
}

// GetWorkbenchConfigByID 根据配置ID获取工作台配置
func (s *WorkbenchService) GetWorkbenchConfigByID(ctx context.Context, configID string) (*dto.WorkbenchConfigResponse, error) {
	config, err := s.workbenchRepo.FindByConfigID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("查询工作台配置失败: %w", err)
	}
	if config == nil {
		return nil, ErrWorkbenchConfigNotFound
	}

	// 解析JSON配置数据
	var configData dto.WorkbenchConfigData
	if err := json.Unmarshal([]byte(config.ConfigData), &configData); err != nil {
		return nil, fmt.Errorf("解析配置数据失败: %w", err)
	}

	return &dto.WorkbenchConfigResponse{
		RoleKey:    config.RoleKey,
		ConfigName: config.ConfigName,
		Config:     configData,
	}, nil
}

// ListWorkbenchConfigs 列出所有工作台配置
func (s *WorkbenchService) ListWorkbenchConfigs(ctx context.Context) ([]dto.WorkbenchConfigListResponse, error) {
	configs, err := s.workbenchRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询工作台配置列表失败: %w", err)
	}

	result := make([]dto.WorkbenchConfigListResponse, len(configs))
	for i, config := range configs {
		result[i] = dto.WorkbenchConfigListResponse{
			ConfigID:    config.ConfigID,
			RoleKey:     config.RoleKey,
			ConfigName:  config.ConfigName,
			Description: config.Description,
			Status:      config.Status,
			CreatedAt:   config.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   config.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return result, nil
}

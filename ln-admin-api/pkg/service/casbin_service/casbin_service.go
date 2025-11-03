package casbin_service

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
)

// CasbinService Casbin服务封装
type CasbinService struct{}

// NewCasbinService 创建Casbin服务
func NewCasbinService() *CasbinService {
	return &CasbinService{}
}

// AddRoleForUser 为用户添加角色（Casbin策略：g, user_id, role_key）
func (s *CasbinService) AddRoleForUser(ctx context.Context, userID, roleKey string) error {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("Casbin enforcer未初始化")
	}

	_, err := enforcer.AddGroupingPolicy(userID, roleKey)
	if err != nil {
		return fmt.Errorf("添加用户角色失败: %w", err)
	}
	return nil
}

// DeleteRoleForUser 删除用户角色
func (s *CasbinService) DeleteRoleForUser(ctx context.Context, userID, roleKey string) error {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("Casbin enforcer未初始化")
	}

	_, err := enforcer.RemoveGroupingPolicy(userID, roleKey)
	if err != nil {
		return fmt.Errorf("删除用户角色失败: %w", err)
	}
	return nil
}

// AddPolicy 添加权限策略（Casbin策略：p, role_key, resource_path, method）
func (s *CasbinService) AddPolicy(ctx context.Context, roleKey, resourcePath, method string) error {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("Casbin enforcer未初始化")
	}

	_, err := enforcer.AddPolicy(roleKey, resourcePath, method)
	if err != nil {
		return fmt.Errorf("添加权限策略失败: %w", err)
	}
	return nil
}

// RemovePolicy 移除权限策略
func (s *CasbinService) RemovePolicy(ctx context.Context, roleKey, resourcePath, method string) error {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("Casbin enforcer未初始化")
	}

	_, err := enforcer.RemovePolicy(roleKey, resourcePath, method)
	if err != nil {
		return fmt.Errorf("移除权限策略失败: %w", err)
	}
	return nil
}

// AddPolicyForUser 直接为用户添加权限策略（Casbin策略：p, user_id, resource_path, method）
func (s *CasbinService) AddPolicyForUser(ctx context.Context, userID, resourcePath, method string) error {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("Casbin enforcer未初始化")
	}

	_, err := enforcer.AddPolicy(userID, resourcePath, method)
	if err != nil {
		return fmt.Errorf("添加用户权限策略失败: %w", err)
	}
	return nil
}

// RemovePolicyForUser 移除用户直接权限策略
func (s *CasbinService) RemovePolicyForUser(ctx context.Context, userID, resourcePath, method string) error {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return fmt.Errorf("Casbin enforcer未初始化")
	}

	_, err := enforcer.RemovePolicy(userID, resourcePath, method)
	if err != nil {
		return fmt.Errorf("移除用户权限策略失败: %w", err)
	}
	return nil
}

// Enforce 检查权限
func (s *CasbinService) Enforce(userID, resourcePath, method string) (bool, error) {
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return false, fmt.Errorf("Casbin enforcer未初始化")
	}

	return enforcer.Enforce(userID, resourcePath, method)
}

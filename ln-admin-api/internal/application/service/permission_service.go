package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/pkg/service/casbin_service"

	"github.com/google/uuid"
)

var (
	ErrRoleNotFound       = errors.New("角色不存在")
	ErrRoleExists         = errors.New("角色已存在")
	ErrPermissionNotFound = errors.New("权限不存在")
	ErrPermissionExists   = errors.New("权限已存在")
)

// PermissionService 权限管理服务
type PermissionService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	userRoleRepo   repository.UserRoleRepository
	casbinSvc      *casbin_service.CasbinService
}

// NewPermissionService 创建权限管理服务
func NewPermissionService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
) *PermissionService {
	casbinSvc := casbin_service.NewCasbinService()
	return &PermissionService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRoleRepo:   userRoleRepo,
		casbinSvc:      casbinSvc,
	}
}

// AssignRoleToUser 为用户分配角色
func (s *PermissionService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 分配角色
	if err := s.userRoleRepo.AssignRole(ctx, userID, roleID); err != nil {
		return err
	}

	// 更新Casbin策略：用户拥有角色
	if err := s.casbinSvc.AddRoleForUser(ctx, userID, role.RoleKey); err != nil {
		return fmt.Errorf("更新Casbin策略失败: %w", err)
	}

	return nil
}

// RemoveRoleFromUser 移除用户角色
func (s *PermissionService) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 移除角色
	if err := s.userRoleRepo.RemoveRole(ctx, userID, roleID); err != nil {
		return err
	}

	// 更新Casbin策略：移除用户角色
	if err := s.casbinSvc.DeleteRoleForUser(ctx, userID, role.RoleKey); err != nil {
		return fmt.Errorf("更新Casbin策略失败: %w", err)
	}

	return nil
}

// AssignPermissionToRole 为角色分配权限
func (s *PermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 检查权限是否存在
	permission, err := s.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	// 更新Casbin策略：角色拥有权限
	// 策略格式：p, role_key, resource_path, method
	if err := s.casbinSvc.AddPolicy(ctx, role.RoleKey, permission.ResourcePath, permission.Method); err != nil {
		return fmt.Errorf("更新Casbin策略失败: %w", err)
	}

	return nil
}

// RemovePermissionFromRole 移除角色权限
func (s *PermissionService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 检查权限是否存在
	permission, err := s.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	// 更新Casbin策略：移除角色权限
	if err := s.casbinSvc.RemovePolicy(ctx, role.RoleKey, permission.ResourcePath, permission.Method); err != nil {
		return fmt.Errorf("更新Casbin策略失败: %w", err)
	}

	return nil
}

// CreateRole 创建角色
func (s *PermissionService) CreateRole(ctx context.Context, role *entity.Role) error {
	// 检查角色标识是否已存在
	existing, err := s.roleRepo.FindByKey(ctx, role.RoleKey)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if existing != nil {
		return ErrRoleExists
	}

	// 生成角色ID
	if role.RoleID == "" {
		role.RoleID = uuid.New().String()
	}

	// 创建角色
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return err
	}

	return nil
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	// 检查权限标识是否已存在
	existing, err := s.permissionRepo.FindByKey(ctx, permission.PermissionKey)
	if err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}
	if existing != nil {
		return ErrPermissionExists
	}

	// 生成权限ID
	if permission.PermissionID == "" {
		permission.PermissionID = uuid.New().String()
	}

	// 创建权限
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		return err
	}

	return nil
}

// GetUserRoles 获取用户角色列表
func (s *PermissionService) GetUserRoles(ctx context.Context, userID string) ([]*entity.Role, error) {
	return s.roleRepo.FindRolesByUserID(ctx, userID)
}

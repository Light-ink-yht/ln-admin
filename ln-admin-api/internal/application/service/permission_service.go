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
	roleRepo                repository.RoleRepository
	permissionRepo          repository.PermissionRepository
	userRoleRepo            repository.UserRoleRepository
	rolePermissionGrantRepo repository.RolePermissionGrantRepository
	casbinSvc               *casbin_service.CasbinService
}

// NewPermissionService 创建权限管理服务
func NewPermissionService(
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	rolePermissionGrantRepo repository.RolePermissionGrantRepository,
) *PermissionService {
	casbinSvc := casbin_service.NewCasbinService()
	return &PermissionService{
		roleRepo:                roleRepo,
		permissionRepo:          permissionRepo,
		userRoleRepo:            userRoleRepo,
		rolePermissionGrantRepo: rolePermissionGrantRepo,
		casbinSvc:               casbinSvc,
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

// AssignPermissionToRole 为角色分配权限（支持层级授权，需要传入授予者用户ID）
func (s *PermissionService) AssignPermissionToRole(ctx context.Context, grantorUserID, roleID, permissionID string) error {
	// 使用新的层级授权方法
	return s.GrantPermissionToRole(ctx, grantorUserID, roleID, permissionID)
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

// GetAllRoles 获取所有角色（不分页）
func (s *PermissionService) GetAllRoles(ctx context.Context) ([]*entity.Role, int64, error) {
	return s.roleRepo.List(ctx, 1, 10000, map[string]interface{}{})
}

// GetRoleList 获取角色列表（分页）
func (s *PermissionService) GetRoleList(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.Role, int64, error) {
	return s.roleRepo.List(ctx, page, pageSize, conditions)
}

// GetRoleByID 根据ID获取角色
func (s *PermissionService) GetRoleByID(ctx context.Context, roleID string) (*entity.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// GetRoleByKey 根据角色标识获取角色
func (s *PermissionService) GetRoleByKey(ctx context.Context, roleKey string) (*entity.Role, error) {
	role, err := s.roleRepo.FindByKey(ctx, roleKey)
	if err != nil {
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// UpdateRole 更新角色
func (s *PermissionService) UpdateRole(ctx context.Context, role *entity.Role) error {
	// 检查角色是否存在
	existing, err := s.roleRepo.FindByID(ctx, role.RoleID)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	// 更新角色
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return err
	}

	return nil
}

// DeleteRole 删除角色
func (s *PermissionService) DeleteRole(ctx context.Context, roleID string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// 删除角色
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return err
	}

	return nil
}

// GetPermissionList 获取权限列表（分页）
func (s *PermissionService) GetPermissionList(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.Permission, int64, error) {
	return s.permissionRepo.List(ctx, page, pageSize, conditions)
}

// GetPermissionByID 根据ID获取权限
func (s *PermissionService) GetPermissionByID(ctx context.Context, permissionID string) (*entity.Permission, error) {
	permission, err := s.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return nil, fmt.Errorf("查询权限失败: %w", err)
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}
	return permission, nil
}

// GetPermissionByResourceAndMethod 根据资源路径和方法获取权限
func (s *PermissionService) GetPermissionByResourceAndMethod(ctx context.Context, resourcePath, method string) (*entity.Permission, error) {
	permission, err := s.permissionRepo.FindByResource(ctx, resourcePath, method)
	if err != nil {
		return nil, fmt.Errorf("查询权限失败: %w", err)
	}
	// 找不到权限时返回nil，不返回错误
	return permission, nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(ctx context.Context, permission *entity.Permission) error {
	// 检查权限是否存在
	existing, err := s.permissionRepo.FindByID(ctx, permission.PermissionID)
	if err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}
	if existing == nil {
		return ErrPermissionNotFound
	}

	// 更新权限
	if err := s.permissionRepo.Update(ctx, permission); err != nil {
		return err
	}

	return nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(ctx context.Context, permissionID string) error {
	// 检查权限是否存在
	permission, err := s.permissionRepo.FindByID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}
	if permission == nil {
		return ErrPermissionNotFound
	}

	// 删除权限
	if err := s.permissionRepo.Delete(ctx, permissionID); err != nil {
		return err
	}

	return nil
}

// AssignPermissionsToUser 直接为用户分配权限（不通过角色）
func (s *PermissionService) AssignPermissionsToUser(ctx context.Context, userID string, permissionIDs []string) error {
	// 简化方案：先获取所有权限，移除所有用户直接权限，然后添加新的权限
	// 使用List方法获取所有权限（设置大pageSize）
	allPermissions, _, err := s.permissionRepo.List(ctx, 1, 10000, make(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("查询所有权限失败: %w", err)
	}

	// 移除用户的所有直接权限（遍历所有权限，尝试移除）
	// 这样可以确保清理所有旧的直接权限
	for _, permission := range allPermissions {
		_ = s.casbinSvc.RemovePolicyForUser(ctx, userID, permission.ResourcePath, permission.Method)
	}

	// 添加新的权限
	for _, permissionID := range permissionIDs {
		permission, err := s.permissionRepo.FindByID(ctx, permissionID)
		if err != nil {
			return fmt.Errorf("查询权限失败: %w", err)
		}
		if permission == nil {
			return fmt.Errorf("权限不存在: %s", permissionID)
		}
		// 添加用户直接权限策略：p, user_id, resource_path, method
		if err := s.casbinSvc.AddPolicyForUser(ctx, userID, permission.ResourcePath, permission.Method); err != nil {
			return fmt.Errorf("添加用户权限失败: %w", err)
		}
	}

	return nil
}

// IsSuperAdmin 检查用户是否是超级管理员
func (s *PermissionService) IsSuperAdmin(ctx context.Context, userID string) (bool, error) {
	roles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.RoleKey == "super_admin" {
			return true, nil
		}
	}
	return false, nil
}

// GetVisiblePermissions 获取用户可见的权限列表（根据层级授权）
func (s *PermissionService) GetVisiblePermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	// 检查是否是超级管理员
	isSuperAdmin, err := s.IsSuperAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("检查超级管理员失败: %w", err)
	}

	// 超级管理员可以看到所有权限
	if isSuperAdmin {
		allPermissions, _, err := s.permissionRepo.List(ctx, 1, 10000, map[string]interface{}{})
		if err != nil {
			return nil, fmt.Errorf("查询所有权限失败: %w", err)
		}
		return allPermissions, nil
	}

	// 获取用户的角色
	roles, err := s.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 收集所有被授予的权限ID
	permissionIDSet := make(map[string]bool)
	for _, role := range roles {
		// 查询该角色被授予的所有权限
		grantedPermissionIDs, err := s.rolePermissionGrantRepo.FindPermissionsByGranteeRoleID(ctx, role.RoleID)
		if err != nil {
			return nil, fmt.Errorf("查询角色被授予的权限失败: %w", err)
		}
		for _, permissionID := range grantedPermissionIDs {
			permissionIDSet[permissionID] = true
		}
	}

	// 如果没有被授予的权限，返回空列表
	if len(permissionIDSet) == 0 {
		return []*entity.Permission{}, nil
	}

	// 查询权限详情
	permissionIDs := make([]string, 0, len(permissionIDSet))
	for permissionID := range permissionIDSet {
		permissionIDs = append(permissionIDs, permissionID)
	}

	// 批量查询权限
	permissions := make([]*entity.Permission, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		permission, err := s.permissionRepo.FindByID(ctx, permissionID)
		if err != nil {
			continue // 忽略查询失败的权限
		}
		if permission != nil {
			permissions = append(permissions, permission)
		}
	}

	return permissions, nil
}

// CanGrantPermission 检查用户是否有权限授予某个权限
func (s *PermissionService) CanGrantPermission(ctx context.Context, grantorUserID, permissionID string) (bool, error) {
	// 检查是否是超级管理员
	isSuperAdmin, err := s.IsSuperAdmin(ctx, grantorUserID)
	if err != nil {
		return false, err
	}
	if isSuperAdmin {
		return true, nil
	}

	// 获取用户的角色
	roles, err := s.GetUserRoles(ctx, grantorUserID)
	if err != nil {
		return false, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 检查用户的角色是否被授予了该权限
	for _, role := range roles {
		// 查询该角色被授予的所有权限
		grantedPermissionIDs, err := s.rolePermissionGrantRepo.FindPermissionsByGranteeRoleID(ctx, role.RoleID)
		if err != nil {
			continue
		}
		// 检查是否包含该权限
		for _, grantedPermissionID := range grantedPermissionIDs {
			if grantedPermissionID == permissionID {
				return true, nil
			}
		}
	}

	return false, nil
}

// GrantPermissionToRole 层级授权（为角色分配权限，记录授予关系）
func (s *PermissionService) GrantPermissionToRole(ctx context.Context, grantorUserID, roleID, permissionID string) error {
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

	// 获取授予者的角色
	grantorRoles, err := s.GetUserRoles(ctx, grantorUserID)
	if err != nil {
		return fmt.Errorf("获取授予者角色失败: %w", err)
	}
	if len(grantorRoles) == 0 {
		return fmt.Errorf("授予者没有角色")
	}

	// 检查是否是超级管理员
	isSuperAdmin, err := s.IsSuperAdmin(ctx, grantorUserID)
	if err != nil {
		return fmt.Errorf("检查超级管理员失败: %w", err)
	}

	// 如果不是超级管理员，检查是否有权限授予该权限
	if !isSuperAdmin {
		canGrant, err := s.CanGrantPermission(ctx, grantorUserID, permissionID)
		if err != nil {
			return fmt.Errorf("检查授权权限失败: %w", err)
		}
		if !canGrant {
			return fmt.Errorf("没有权限授予该权限")
		}
	}

	// 使用第一个角色作为授予者角色（通常用户只有一个主要角色）
	grantorRoleID := grantorRoles[0].RoleID
	// 如果是超级管理员，grantorRoleID 使用超级管理员角色ID
	if isSuperAdmin {
		// 查找超级管理员角色
		superAdminRole, err := s.roleRepo.FindByKey(ctx, "super_admin")
		if err == nil && superAdminRole != nil {
			grantorRoleID = superAdminRole.RoleID
		}
	}

	// 记录授予关系
	if err := s.rolePermissionGrantRepo.GrantPermission(ctx, grantorRoleID, roleID, permissionID); err != nil {
		return fmt.Errorf("记录授予关系失败: %w", err)
	}

	// 更新Casbin策略：角色拥有权限
	resourcePath := permission.ResourcePath
	// 移除路径前缀 /api（因为中间件会移除这个前缀）
	if len(resourcePath) > 4 && resourcePath[:4] == "/api" {
		resourcePath = resourcePath[4:]
	}
	if err := s.casbinSvc.AddPolicy(ctx, role.RoleKey, resourcePath, permission.Method); err != nil {
		return fmt.Errorf("更新Casbin策略失败: %w", err)
	}

	return nil
}

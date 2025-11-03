package init

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/config"
	"github.com/Light-ink-yht/ln-admin/pkg/service/casbin_service"
	"github.com/Light-ink-yht/ln-admin/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InitDefaultData 初始化默认数据（默认用户和角色）
func InitDefaultData(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	permissionService *service.PermissionService,
) error {
	ctx := context.Background()

	// 初始化默认角色
	if err := initDefaultRoles(ctx, roleRepo); err != nil {
		return fmt.Errorf("初始化默认角色失败: %w", err)
	}

	// 初始化默认权限
	if err := initDefaultPermissions(ctx, permissionRepo); err != nil {
		return fmt.Errorf("初始化默认权限失败: %w", err)
	}

	// 初始化默认用户
	if err := initDefaultUser(ctx, userRepo, roleRepo, userRoleRepo, permissionService); err != nil {
		return fmt.Errorf("初始化默认用户失败: %w", err)
	}

	// 为超级管理员分配所有权限
	if err := assignAllPermissionsToAdmin(ctx, roleRepo, permissionRepo, permissionService); err != nil {
		return fmt.Errorf("为超级管理员分配权限失败: %w", err)
	}

	logger.Info("默认数据初始化成功")
	return nil
}

// initDefaultRoles 初始化默认角色
func initDefaultRoles(ctx context.Context, roleRepo repository.RoleRepository) error {
	defaultRoles := []*entity.Role{
		{
			RoleID:      "role_super_admin",
			RoleKey:     "super_admin",
			RoleName:    "超级管理员",
			Description: "系统超级管理员，拥有所有权限",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			RoleID:      "role_admin",
			RoleKey:     "admin",
			RoleName:    "管理员",
			Description: "系统管理员，拥有管理权限",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			RoleID:      "role_user",
			RoleKey:     "user",
			RoleName:    "普通用户",
			Description: "普通用户角色",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
	}

	for _, role := range defaultRoles {
		// 检查角色是否已存在
		existing, err := roleRepo.FindByKey(ctx, role.RoleKey)
		if err != nil {
			return fmt.Errorf("查询角色失败: %w", err)
		}
		if existing == nil {
			// 角色不存在，创建新角色
			if err := roleRepo.Create(ctx, role); err != nil {
				// 如果是重复键错误，忽略
				if !isDuplicateError(err) {
					return fmt.Errorf("创建角色 %s 失败: %w", role.RoleName, err)
				}
			} else {
				logger.Info("创建默认角色", zap.String("role", role.RoleName))
			}
		}
	}

	return nil
}

// initDefaultPermissions 初始化默认权限
func initDefaultPermissions(ctx context.Context, permissionRepo repository.PermissionRepository) error {
	defaultPermissions := []*entity.Permission{
		// 用户管理权限
		{
			PermissionID:   "perm_user_list",
			PermissionKey:  "user:list",
			PermissionName: "用户列表",
			ResourcePath:   "/api/user/list",
			Method:         "GET",
			Description:    "查看用户列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_user_info",
			PermissionKey:  "user:info",
			PermissionName: "用户信息",
			ResourcePath:   "/api/user/userinfo",
			Method:         "GET",
			Description:    "查看用户信息",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_user_create",
			PermissionKey:  "user:create",
			PermissionName: "创建用户",
			ResourcePath:   "/api/user/create",
			Method:         "POST",
			Description:    "创建新用户",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_user_update",
			PermissionKey:  "user:update",
			PermissionName: "更新用户",
			ResourcePath:   "/api/user/update",
			Method:         "PUT",
			Description:    "更新用户信息",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_user_delete",
			PermissionKey:  "user:delete",
			PermissionName: "删除用户",
			ResourcePath:   "/api/user/delete",
			Method:         "DELETE",
			Description:    "删除用户",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
	}

	for _, permission := range defaultPermissions {
		// 检查权限是否已存在
		existing, err := permissionRepo.FindByKey(ctx, permission.PermissionKey)
		if err != nil {
			return fmt.Errorf("查询权限失败: %w", err)
		}
		if existing == nil {
			// 权限不存在，创建新权限
			if err := permissionRepo.Create(ctx, permission); err != nil {
				// 如果是重复键错误，忽略
				if !isDuplicateError(err) {
					return fmt.Errorf("创建权限 %s 失败: %w", permission.PermissionName, err)
				}
			} else {
				logger.Info("创建默认权限", zap.String("permission", permission.PermissionName))
			}
		}
	}

	return nil
}

// initDefaultUser 初始化默认用户
func initDefaultUser(
	ctx context.Context,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
	permissionService *service.PermissionService,
) error {
	phone := "18797131041"
	defaultPassword := "hello#123world"

	// 检查用户是否已存在
	existing, err := userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	var userID string

	// 如果用户已存在，确保分配角色和权限
	if existing != nil {
		userID = existing.UserID
		logger.Info("默认用户已存在", zap.String("phone", phone), zap.String("user_id", userID))
	} else {
		// 加密密码
		hashedPassword, err := utils.HashPassword(defaultPassword, config.Cfg.Password.Cost)
		if err != nil {
			return fmt.Errorf("加密密码失败: %w", err)
		}

		// 生成用户ID
		userID = uuid.New().String()

		// 创建默认用户
		user := &entity.User{
			UserID:     userID,
			Phone:      &phone,
			Password:   hashedPassword,
			Nickname:   "超级管理员",
			FullName:   "系统管理员",
			Status:     "1", // 启用
			Gender:     "3", // 未知
			CreatorID:  "system",
			ModifierID: "system",
		}

		if err := userRepo.Create(ctx, user); err != nil {
			if isDuplicateError(err) {
				logger.Info("默认用户已存在", zap.String("phone", phone))
				// 如果创建时出现重复错误，重新查询用户
				existing, err = userRepo.FindByPhone(ctx, phone)
				if err != nil || existing == nil {
					return fmt.Errorf("创建默认用户失败: %w", err)
				}
				userID = existing.UserID
			} else {
				return fmt.Errorf("创建默认用户失败: %w", err)
			}
		}
	}

	// 为默认用户分配超级管理员角色（无论用户是新创建还是已存在）
	superAdminRole, err := roleRepo.FindByKey(ctx, "super_admin")
	if err != nil {
		return fmt.Errorf("查询超级管理员角色失败: %w", err)
	}
	if superAdminRole == nil {
		return fmt.Errorf("超级管理员角色不存在")
	}

	// 使用权限服务分配角色（会自动更新Casbin策略）
	// 注意：即使角色已分配，也会确保Casbin策略正确
	if err := permissionService.AssignRoleToUser(ctx, userID, superAdminRole.RoleID); err != nil {
		// 如果是重复键错误，说明角色已经分配，继续执行以确保权限正确
		if !isDuplicateError(err) {
			logger.Warn("分配角色失败", zap.Error(err))
		} else {
			logger.Info("角色已分配，确保Casbin策略正确", zap.String("user_id", userID))
		}
	}

	// 确保Casbin策略中存在用户角色关系（即使已存在，重新添加也不会出错）
	casbinSvc := casbin_service.NewCasbinService()

	// 先检查用户是否已有该角色
	enforcer := casbin.GetEnforcer()
	if enforcer != nil {
		roles, err := enforcer.GetRolesForUser(userID)
		hasRole := false
		if err == nil {
			for _, role := range roles {
				if role == superAdminRole.RoleKey {
					hasRole = true
					break
				}
			}
		}

		if !hasRole {
			// 用户还没有该角色，添加角色关系
			if err := casbinSvc.AddRoleForUser(ctx, userID, superAdminRole.RoleKey); err != nil {
				logger.Warn("添加Casbin用户角色策略失败", zap.Error(err))
			} else {
				logger.Info("已添加Casbin用户角色策略", zap.String("user_id", userID), zap.String("role", superAdminRole.RoleKey))
			}
		} else {
			logger.Info("用户已拥有角色", zap.String("user_id", userID), zap.String("role", superAdminRole.RoleKey))
		}

		// 重新加载策略确保生效
		if err := enforcer.LoadPolicy(); err != nil {
			logger.Warn("重新加载Casbin策略失败", zap.Error(err))
		}
	}

	logger.Info("默认用户初始化完成",
		zap.String("phone", phone),
		zap.String("user_id", userID),
		zap.String("role", "超级管理员"),
		zap.Bool("is_existing", existing != nil))

	return nil
}

// assignAllPermissionsToAdmin 为超级管理员分配所有权限
func assignAllPermissionsToAdmin(
	ctx context.Context,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	permissionService *service.PermissionService,
) error {
	// 获取超级管理员角色
	superAdminRole, err := roleRepo.FindByKey(ctx, "super_admin")
	if err != nil {
		return fmt.Errorf("查询超级管理员角色失败: %w", err)
	}
	if superAdminRole == nil {
		return fmt.Errorf("超级管理员角色不存在")
	}

	// 获取所有权限
	permissions, _, err := permissionRepo.List(ctx, 1, 1000, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("查询权限列表失败: %w", err)
	}

	// 为超级管理员角色分配所有权限
	casbinSvc := casbin_service.NewCasbinService()
	assignedCount := 0
	for _, permission := range permissions {
		// 使用Casbin直接添加策略（p, role_key, resource_path, method）
		if err := casbinSvc.AddPolicy(ctx, superAdminRole.RoleKey, permission.ResourcePath, permission.Method); err != nil {
			logger.Warn("为超级管理员分配权限失败",
				zap.String("permission", permission.PermissionName),
				zap.Error(err))
			continue
		}
		assignedCount++
	}

	logger.Info("为超级管理员分配权限完成",
		zap.String("role", superAdminRole.RoleName),
		zap.Int("total_permissions", len(permissions)),
		zap.Int("assigned_count", assignedCount))

	return nil
}

// isDuplicateError 检查是否是重复键错误
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "Duplicate entry") ||
		contains(errStr, "duplicated key") ||
		contains(errStr, "1062") ||
		contains(errStr, "UNIQUE constraint")
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 || indexOf(s, substr) >= 0)
}

// indexOf 查找子串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

package defaultdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	infraRepo "github.com/Light-ink-yht/ln-admin/internal/infrastructure/repository"
	"github.com/Light-ink-yht/ln-admin/pkg/config"
	"github.com/Light-ink-yht/ln-admin/pkg/service/casbin_service"
	"github.com/Light-ink-yht/ln-admin/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InitDefaultData 初始化默认数据（默认用户、角色、权限、菜单、短信模板、系统配置等）
func InitDefaultData(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	templateRepo repository.SMSTemplateRepository,
	menuRepo repository.MenuRepository,
	permissionService *service.PermissionService,
) error {
	ctx := context.Background()

	// 1. 初始化默认角色
	if err := initDefaultRoles(ctx, roleRepo); err != nil {
		return fmt.Errorf("初始化默认角色失败: %w", err)
	}

	// 2. 初始化默认权限
	if err := initDefaultPermissions(ctx, permissionRepo); err != nil {
		return fmt.Errorf("初始化默认权限失败: %w", err)
	}

	// 3. 初始化默认菜单（必须在权限之后，因为菜单可能关联权限）
	if err := initDefaultMenus(ctx, menuRepo); err != nil {
		return fmt.Errorf("初始化默认菜单失败: %w", err)
	}

	// 4. 初始化默认短信模板
	if err := initDefaultSMSTemplates(ctx, templateRepo); err != nil {
		return fmt.Errorf("初始化默认短信模板失败: %w", err)
	}

	// 5. 初始化默认系统配置
	configRepo := infraRepo.NewSystemConfigRepository()
	if err := initDefaultSystemConfigs(ctx, configRepo); err != nil {
		return fmt.Errorf("初始化默认系统配置失败: %w", err)
	}

	// 6. 初始化默认用户（超级管理员账号：18797131041）
	if err := initDefaultUser(ctx, userRepo, roleRepo, userRoleRepo, permissionService); err != nil {
		return fmt.Errorf("初始化默认用户失败: %w", err)
	}

	// 7. 为超级管理员分配所有权限
	if err := assignAllPermissionsToAdmin(ctx, roleRepo, permissionRepo, permissionService); err != nil {
		return fmt.Errorf("为超级管理员分配权限失败: %w", err)
	}

	// 8. 为超级管理员角色分配所有菜单
	if err := assignAllMenusToAdmin(ctx, roleRepo, menuRepo); err != nil {
		return fmt.Errorf("为超级管理员分配菜单失败: %w", err)
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

// getPermissionCategory 根据权限标识提取分类
func getPermissionCategory(permissionKey string) string {
	if !strings.Contains(permissionKey, ":") {
		return "其他"
	}
	category := strings.Split(permissionKey, ":")[0]
	categoryMap := map[string]string{
		"user":       "用户管理",
		"role":       "角色管理",
		"permission": "权限管理",
		"system":     "系统配置",
		"sms":        "短信管理",
		"menu":       "菜单管理",
		"ops":        "系统运维",
		"file":       "文件管理",
	}
	if cat, ok := categoryMap[category]; ok {
		return cat
	}
	return "其他"
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
			Category:       getPermissionCategory("user:list"),
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
			PermissionID:   "perm_user_detail",
			PermissionKey:  "user:detail",
			PermissionName: "用户详情",
			ResourcePath:   "/api/user/*",
			Method:         "GET",
			Description:    "查看用户详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_user_create",
			PermissionKey:  "user:create",
			PermissionName: "创建用户",
			ResourcePath:   "/api/user",
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
			ResourcePath:   "/api/user/*",
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
			ResourcePath:   "/api/user/*",
			Method:         "DELETE",
			Description:    "删除用户",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_user_grant",
			PermissionKey:  "user:grant",
			PermissionName: "用户授权",
			ResourcePath:   "/api/user/*",
			Method:         "POST",
			Description:    "给用户直接分配权限",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 角色管理权限
		{
			PermissionID:   "perm_role_list",
			PermissionKey:  "role:list",
			PermissionName: "角色列表",
			ResourcePath:   "/api/role/list",
			Method:         "GET",
			Description:    "查看角色列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_role_detail",
			PermissionKey:  "role:detail",
			PermissionName: "角色详情",
			ResourcePath:   "/api/role/*",
			Method:         "GET",
			Description:    "查看角色详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_role_create",
			PermissionKey:  "role:create",
			PermissionName: "创建角色",
			ResourcePath:   "/api/role",
			Method:         "POST",
			Description:    "创建新角色",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_role_update",
			PermissionKey:  "role:update",
			PermissionName: "更新角色",
			ResourcePath:   "/api/role/*",
			Method:         "PUT",
			Description:    "更新角色信息",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_role_delete",
			PermissionKey:  "role:delete",
			PermissionName: "删除角色",
			ResourcePath:   "/api/role/*",
			Method:         "DELETE",
			Description:    "删除角色",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 权限管理权限
		{
			PermissionID:   "perm_permission_list",
			PermissionKey:  "permission:list",
			PermissionName: "权限列表",
			ResourcePath:   "/api/permission/list",
			Method:         "GET",
			Description:    "查看权限列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_permission_detail",
			PermissionKey:  "permission:detail",
			PermissionName: "权限详情",
			ResourcePath:   "/api/permission/*",
			Method:         "GET",
			Description:    "查看权限详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_permission_create",
			PermissionKey:  "permission:create",
			PermissionName: "创建权限",
			ResourcePath:   "/api/permission",
			Method:         "POST",
			Description:    "创建新权限",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_permission_update",
			PermissionKey:  "permission:update",
			PermissionName: "更新权限",
			ResourcePath:   "/api/permission/*",
			Method:         "PUT",
			Description:    "更新权限信息",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_permission_delete",
			PermissionKey:  "permission:delete",
			PermissionName: "删除权限",
			ResourcePath:   "/api/permission/*",
			Method:         "DELETE",
			Description:    "删除权限",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 系统配置管理权限
		{
			PermissionID:   "perm_system_config_list",
			PermissionKey:  "system:config:list",
			PermissionName: "系统配置列表",
			ResourcePath:   "/api/system/config/list",
			Method:         "GET",
			Description:    "查看系统配置列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_system_config_group",
			PermissionKey:  "system:config:group",
			PermissionName: "系统配置分组",
			ResourcePath:   "/api/system/config/group/*",
			Method:         "GET",
			Description:    "根据分组查看系统配置",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_system_config_detail",
			PermissionKey:  "system:config:detail",
			PermissionName: "系统配置详情",
			ResourcePath:   "/api/system/config/*",
			Method:         "GET",
			Description:    "查看系统配置详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_system_config_create",
			PermissionKey:  "system:config:create",
			PermissionName: "创建系统配置",
			ResourcePath:   "/api/system/config",
			Method:         "POST",
			Description:    "创建系统配置",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_system_config_update",
			PermissionKey:  "system:config:update",
			PermissionName: "更新系统配置",
			ResourcePath:   "/api/system/config/*",
			Method:         "PUT",
			Description:    "更新系统配置",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 系统日志管理权限
		{
			PermissionID:   "perm_system_log_list",
			PermissionKey:  "system:log:list",
			PermissionName: "系统日志列表",
			ResourcePath:   "/api/system/log/list",
			Method:         "GET",
			Description:    "查看系统操作日志列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 短信管理权限
		{
			PermissionID:   "perm_sms_template_list",
			PermissionKey:  "sms:template:list",
			PermissionName: "短信模板列表",
			ResourcePath:   "/api/sms/template/list",
			Method:         "GET",
			Description:    "查看短信模板列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_sms_template_detail",
			PermissionKey:  "sms:template:detail",
			PermissionName: "短信模板详情",
			ResourcePath:   "/api/sms/template/*",
			Method:         "GET",
			Description:    "查看短信模板详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_sms_template_create",
			PermissionKey:  "sms:template:create",
			PermissionName: "创建短信模板",
			ResourcePath:   "/api/sms/template",
			Method:         "POST",
			Description:    "创建短信模板",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_sms_template_update",
			PermissionKey:  "sms:template:update",
			PermissionName: "更新短信模板",
			ResourcePath:   "/api/sms/template/*",
			Method:         "PUT",
			Description:    "更新短信模板",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_sms_template_delete",
			PermissionKey:  "sms:template:delete",
			PermissionName: "删除短信模板",
			ResourcePath:   "/api/sms/template/*",
			Method:         "DELETE",
			Description:    "删除短信模板",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_sms_code_list",
			PermissionKey:  "sms:code:list",
			PermissionName: "短信验证码列表",
			ResourcePath:   "/api/sms/code/list",
			Method:         "GET",
			Description:    "查看短信验证码列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 菜单管理权限
		{
			PermissionID:   "perm_menu_list",
			PermissionKey:  "menu:list",
			PermissionName: "菜单列表",
			ResourcePath:   "/api/menu/list",
			Method:         "GET",
			Description:    "查看菜单列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_menu_detail",
			PermissionKey:  "menu:detail",
			PermissionName: "菜单详情",
			ResourcePath:   "/api/menu/*",
			Method:         "GET",
			Description:    "查看菜单详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_menu_create",
			PermissionKey:  "menu:create",
			PermissionName: "创建菜单",
			ResourcePath:   "/api/menu",
			Method:         "POST",
			Description:    "创建菜单",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_menu_update",
			PermissionKey:  "menu:update",
			PermissionName: "更新菜单",
			ResourcePath:   "/api/menu/*",
			Method:         "PUT",
			Description:    "更新菜单",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_menu_delete",
			PermissionKey:  "menu:delete",
			PermissionName: "删除菜单",
			ResourcePath:   "/api/menu/*",
			Method:         "DELETE",
			Description:    "删除菜单",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 系统运维权限
		{
			PermissionID:   "perm_ops_monitor",
			PermissionKey:  "ops:monitor",
			PermissionName: "系统监控",
			ResourcePath:   "/api/ops/monitor",
			Method:         "GET",
			Description:    "查看系统监控信息（CPU、内存、磁盘等）",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_ops_api_doc",
			PermissionKey:  "ops:api-doc",
			PermissionName: "接口文档",
			ResourcePath:   "/api/swagger/*",
			Method:         "GET",
			Description:    "查看API接口文档",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		// 文件管理权限
		{
			PermissionID:   "perm_file_list",
			PermissionKey:  "file:list",
			PermissionName: "文件列表",
			ResourcePath:   "/api/file/list",
			Method:         "GET",
			Description:    "查看文件列表",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_file_detail",
			PermissionKey:  "file:detail",
			PermissionName: "文件详情",
			ResourcePath:   "/api/file/*",
			Method:         "GET",
			Description:    "查看文件详情",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_file_upload",
			PermissionKey:  "file:upload",
			PermissionName: "上传文件",
			ResourcePath:   "/api/file/upload",
			Method:         "POST",
			Description:    "上传文件",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_file_download",
			PermissionKey:  "file:download",
			PermissionName: "下载文件",
			ResourcePath:   "/api/file/*",
			Method:         "GET",
			Description:    "下载文件",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_file_delete",
			PermissionKey:  "file:delete",
			PermissionName: "删除文件",
			ResourcePath:   "/api/file/*",
			Method:         "DELETE",
			Description:    "删除文件",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_file_storage_config_get",
			PermissionKey:  "file:storage:config:get",
			PermissionName: "查看存储配置",
			ResourcePath:   "/api/file/storage/config",
			Method:         "GET",
			Description:    "查看文件存储配置",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
		{
			PermissionID:   "perm_file_storage_config_save",
			PermissionKey:  "file:storage:config:save",
			PermissionName: "保存存储配置",
			ResourcePath:   "/api/file/storage/config",
			Method:         "POST",
			Description:    "保存文件存储配置",
			Status:         "1",
			CreatorID:      "system",
			ModifierID:     "system",
		},
	}

	for _, permission := range defaultPermissions {
		// 如果没有设置分类，根据 permissionKey 自动设置
		if permission.Category == "" {
			permission.Category = getPermissionCategory(permission.PermissionKey)
		}
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
		} else {
			// 如果权限已存在但分类为空，更新分类
			if existing.Category == "" {
				existing.Category = getPermissionCategory(permission.PermissionKey)
				if err := permissionRepo.Update(ctx, existing); err != nil {
					logger.Warn("更新权限分类失败", zap.String("permission", permission.PermissionName), zap.Error(err))
				}
			}
		}
	}

	return nil
}

// initDefaultSystemConfigs 初始化默认系统配置
func initDefaultSystemConfigs(ctx context.Context, configRepo repository.SystemConfigRepository) error {
	defaultConfigs := []*entity.SystemConfig{
		{
			ConfigKey:   "site_name",
			ConfigValue: "LN Admin",
			ConfigName:  "网站名称",
			ConfigGroup: "system",
			Description: "网站/系统的名称，显示在页面标题、LOGO等位置",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_logo",
			ConfigValue: "",
			ConfigName:  "网站LOGO",
			ConfigGroup: "system",
			Description: "网站LOGO图片URL地址",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_favicon",
			ConfigValue: "",
			ConfigName:  "网站图标",
			ConfigGroup: "system",
			Description: "网站Favicon图标URL地址",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_copyright",
			ConfigValue: "© 2024 LN Admin 基于 Vue 3 + Ant Design Vue",
			ConfigName:  "版权信息",
			ConfigGroup: "system",
			Description: "网站底部版权信息",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_description",
			ConfigValue: "LN Admin 管理系统",
			ConfigName:  "网站描述",
			ConfigGroup: "system",
			Description: "网站/系统的描述信息",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_keywords",
			ConfigValue: "LN Admin,管理系统,后台管理",
			ConfigName:  "网站关键词",
			ConfigGroup: "system",
			Description: "网站SEO关键词，多个关键词用逗号分隔",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_beian",
			ConfigValue: "",
			ConfigName:  "备案号",
			ConfigGroup: "system",
			Description: "网站ICP备案号",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_contact_email",
			ConfigValue: "",
			ConfigName:  "联系邮箱",
			ConfigGroup: "system",
			Description: "系统联系邮箱地址",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_contact_phone",
			ConfigValue: "",
			ConfigName:  "联系电话",
			ConfigGroup: "system",
			Description: "系统联系电话",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "site_address",
			ConfigValue: "",
			ConfigName:  "公司地址",
			ConfigGroup: "system",
			Description: "公司/组织地址",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
		{
			ConfigKey:   "api_doc_url",
			ConfigValue: "",
			ConfigName:  "接口文档地址",
			ConfigGroup: "system",
			Description: "API接口文档地址，用于在接口文档页面进行跳转",
			Status:      "1",
			CreatorID:   "system",
			ModifierID:  "system",
		},
	}

	for _, config := range defaultConfigs {
		// 检查配置是否已存在
		existing, err := configRepo.GetByKey(ctx, config.ConfigKey)
		if err != nil {
			return fmt.Errorf("查询配置失败: %w", err)
		}
		if existing == nil {
			// 配置不存在，创建新配置
			if err := configRepo.Create(ctx, config); err != nil {
				// 如果是重复键错误，忽略
				if !isDuplicateError(err) {
					return fmt.Errorf("创建配置 %s 失败: %w", config.ConfigName, err)
				}
			} else {
				logger.Info("创建默认系统配置", zap.String("config", config.ConfigName))
			}
		}
	}

	return nil
}

// initDefaultSMSTemplates 初始化默认短信模板
func initDefaultSMSTemplates(ctx context.Context, templateRepo repository.SMSTemplateRepository) error {
	defaultTemplates := []*entity.SMSTemplate{
		{
			TemplateID:        "template_register_001",
			Type:              "register",
			TemplateName:      "用户注册模板",
			TencentTemplateID: "your_register_template_id",
			Title:             "您的应用名称",
			Content:           "您的注册验证码是：{code}，10分钟内有效，请勿泄露给他人。",
			Description:       "用户注册场景使用的短信模板",
			Status:            "1",
			CreatorID:         "system",
			ModifierID:        "system",
		},
		{
			TemplateID:        "template_forgot_001",
			Type:              "forgot",
			TemplateName:      "忘记密码模板",
			TencentTemplateID: "your_forgot_template_id",
			Title:             "您的应用名称",
			Content:           "您的密码重置验证码是：{code}，10分钟内有效，请勿泄露给他人。",
			Description:       "忘记密码场景使用的短信模板",
			Status:            "1",
			CreatorID:         "system",
			ModifierID:        "system",
		},
	}

	for _, template := range defaultTemplates {
		// 检查模板是否已存在
		existing, err := templateRepo.FindByType(ctx, template.Type)
		if err != nil {
			return fmt.Errorf("查询短信模板失败: %w", err)
		}
		if existing == nil {
			// 模板不存在，创建新模板
			if err := templateRepo.Create(ctx, template); err != nil {
				// 如果是重复键错误，忽略
				if !isDuplicateError(err) {
					return fmt.Errorf("创建短信模板 %s 失败: %w", template.TemplateName, err)
				}
			} else {
				logger.Info("创建默认短信模板", zap.String("template", template.TemplateName))
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
	// 如果用户已存在，检查并分配角色
	if existing != nil {
		userID = existing.UserID
		logger.Info("默认用户已存在", zap.String("phone", phone), zap.String("user_id", userID))

		// 检查用户是否已有超级管理员角色
		userRoles, err := userRoleRepo.FindRolesByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("查询用户角色失败: %w", err)
		}

		// 获取超级管理员角色
		superAdminRole, err := roleRepo.FindByKey(ctx, "super_admin")
		if err != nil {
			return fmt.Errorf("查询超级管理员角色失败: %w", err)
		}
		if superAdminRole == nil {
			return fmt.Errorf("超级管理员角色不存在")
		}

		// 检查用户是否已有超级管理员角色
		hasSuperAdminRole := false
		for _, roleID := range userRoles {
			if roleID == superAdminRole.RoleID {
				hasSuperAdminRole = true
				break
			}
		}

		// 如果没有超级管理员角色，则分配
		if !hasSuperAdminRole {
			logger.Info("为已存在的默认用户分配超级管理员角色",
				zap.String("phone", phone),
				zap.String("user_id", userID))
			if err := permissionService.AssignRoleToUser(ctx, userID, superAdminRole.RoleID); err != nil {
				logger.Warn("分配角色失败", zap.Error(err))
				// 不返回错误，继续执行
			}
		}

		// 确保Casbin策略中存在用户角色关系
		enforcer := casbin.GetEnforcer()
		if enforcer != nil {
			roles, err := enforcer.GetRolesForUser(userID)
			hasRoleInCasbin := false
			if err == nil {
				for _, role := range roles {
					if role == superAdminRole.RoleKey {
						hasRoleInCasbin = true
						break
					}
				}
			}

			if !hasRoleInCasbin {
				casbinSvc := casbin_service.NewCasbinService()
				if err := casbinSvc.AddRoleForUser(ctx, userID, superAdminRole.RoleKey); err != nil {
					logger.Warn("添加Casbin用户角色策略失败", zap.Error(err))
				} else {
					logger.Info("已添加Casbin用户角色策略", zap.String("user_id", userID), zap.String("role", superAdminRole.RoleKey))
					// 重新加载策略
					if err := enforcer.LoadPolicy(); err != nil {
						logger.Warn("重新加载Casbin策略失败", zap.Error(err))
					}
				}
			}
		}

		return nil
	}

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
			// 用户已存在，重新查询并分配角色
			existing, err := userRepo.FindByPhone(ctx, phone)
			if err != nil {
				return fmt.Errorf("查询用户失败: %w", err)
			}
			if existing == nil {
				return fmt.Errorf("用户创建失败且查询不到用户")
			}
			userID = existing.UserID

			// 检查并分配角色
			userRoles, err := userRoleRepo.FindRolesByUserID(ctx, userID)
			if err == nil {
				superAdminRole, err := roleRepo.FindByKey(ctx, "super_admin")
				if err == nil && superAdminRole != nil {
					hasSuperAdminRole := false
					for _, roleID := range userRoles {
						if roleID == superAdminRole.RoleID {
							hasSuperAdminRole = true
							break
						}
					}
					if !hasSuperAdminRole {
						if err := permissionService.AssignRoleToUser(ctx, userID, superAdminRole.RoleID); err != nil {
							logger.Warn("分配角色失败", zap.Error(err))
						}
					}
				}
			}
			logger.Info("默认用户已存在", zap.String("phone", phone), zap.String("user_id", userID))
			return nil
		}
		return fmt.Errorf("创建默认用户失败: %w", err)
	}

	// 为默认用户分配超级管理员角色
	superAdminRole, err := roleRepo.FindByKey(ctx, "super_admin")
	if err != nil {
		return fmt.Errorf("查询超级管理员角色失败: %w", err)
	}
	if superAdminRole == nil {
		return fmt.Errorf("超级管理员角色不存在")
	}

	// 使用权限服务分配角色（会自动更新Casbin策略）
	if err := permissionService.AssignRoleToUser(ctx, userID, superAdminRole.RoleID); err != nil {
		// 如果是重复键错误，忽略
		if !isDuplicateError(err) {
			logger.Warn("分配角色失败", zap.Error(err))
		}
	}

	logger.Info("创建默认用户成功",
		zap.String("phone", phone),
		zap.String("user_id", userID),
		zap.String("role", "超级管理员"))

	return nil
}

// assignAllPermissionsToAdmin 为超级管理员分配所有权限（同时记录授予关系）
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

	// 获取所有权限（分页获取，确保获取所有权限）
	allPermissions := make([]*entity.Permission, 0)
	page := 1
	pageSize := 1000
	for {
		permissions, total, err := permissionRepo.List(ctx, page, pageSize, map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("查询权限列表失败: %w", err)
		}
		allPermissions = append(allPermissions, permissions...)
		if int64(len(allPermissions)) >= total {
			break
		}
		page++
	}

	logger.Info("获取所有权限", zap.Int("total", len(allPermissions)))

	// 获取角色权限授予关系仓库
	rolePermissionGrantRepo := infraRepo.NewRolePermissionGrantRepository()

	// 为超级管理员角色分配所有权限
	casbinSvc := casbin_service.NewCasbinService()
	assignedCount := 0
	grantedCount := 0
	for _, permission := range allPermissions {
		// 移除路径前缀 /api（因为中间件会移除这个前缀）
		resourcePath := strings.TrimPrefix(permission.ResourcePath, "/api")

		// 使用Casbin直接添加策略（p, role_key, resource_path, method）
		// 注意：resourcePath 已经移除了 /api 前缀，与中间件验证时的路径格式一致
		if err := casbinSvc.AddPolicy(ctx, superAdminRole.RoleKey, resourcePath, permission.Method); err != nil {
			logger.Warn("为超级管理员分配权限失败",
				zap.String("permission", permission.PermissionName),
				zap.String("resource_path", resourcePath),
				zap.Error(err))
			continue
		}

		// 记录授予关系（超级管理员自己授予给自己，grantorRoleID 和 granteeRoleID 都是超级管理员角色ID）
		// 这样超级管理员就拥有了所有权限，并且可以授予给其他角色
		if err := rolePermissionGrantRepo.GrantPermission(ctx, superAdminRole.RoleID, superAdminRole.RoleID, permission.PermissionID); err != nil {
			logger.Warn("记录超级管理员权限授予关系失败",
				zap.String("permission", permission.PermissionName),
				zap.String("permission_id", permission.PermissionID),
				zap.Error(err))
			// 继续处理，不中断
		} else {
			grantedCount++
		}

		assignedCount++
		logger.Debug("为超级管理员添加权限策略",
			zap.String("role", superAdminRole.RoleKey),
			zap.String("resource", resourcePath),
			zap.String("method", permission.Method))
	}

	logger.Info("为超级管理员分配权限完成",
		zap.String("role", superAdminRole.RoleName),
		zap.Int("total_permissions", len(allPermissions)),
		zap.Int("assigned_count", assignedCount),
		zap.Int("granted_count", grantedCount))

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

// initDefaultMenus 初始化默认菜单（从menu_service.go中的getAllMenus提取）
func initDefaultMenus(ctx context.Context, menuRepo repository.MenuRepository) error {
	// 定义默认菜单（侧边栏菜单，类型为1）
	defaultMenus := []*struct {
		MenuKey     string
		Title       string
		Path        string
		Icon        string
		ParentID    string
		Sort        int
		Permission  string
		MenuType    string
		Description string
		Children    []*struct {
			MenuKey     string
			Title       string
			Path        string
			Icon        string
			Sort        int
			Permission  string
			Description string
		}
	}{
		{
			MenuKey:     "dashboard",
			Title:       "仪表盘",
			Path:        "",
			Icon:        "DashboardOutlined",
			ParentID:    "0",
			Sort:        1,
			Permission:  "",
			MenuType:    "1",
			Description: "系统仪表盘",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "dashboard-workbench",
					Title:       "工作台",
					Path:        "/",
					Icon:        "",
					Sort:        1,
					Permission:  "/dashboard/workbench:GET",
					Description: "工作台页面",
				},
				{
					MenuKey:     "dashboard-analysis",
					Title:       "分析页",
					Path:        "/dashboard/analysis",
					Icon:        "",
					Sort:        2,
					Permission:  "/dashboard/analysis:GET",
					Description: "数据分析页面",
				},
			},
		},
		{
			MenuKey:     "user-management",
			Title:       "用户管理",
			Path:        "",
			Icon:        "TeamOutlined",
			ParentID:    "0",
			Sort:        2,
			Permission:  "",
			MenuType:    "1",
			Description: "用户管理模块",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "user-list",
					Title:       "用户列表",
					Path:        "/user/list",
					Icon:        "",
					Sort:        1,
					Permission:  "/user/list:GET",
					Description: "用户列表页面",
				},
				{
					MenuKey:     "role-list",
					Title:       "角色管理",
					Path:        "/role/list",
					Icon:        "",
					Sort:        2,
					Permission:  "/role/list:GET",
					Description: "角色管理页面",
				},
				{
					MenuKey:     "permission-list",
					Title:       "权限管理",
					Path:        "/permission/list",
					Icon:        "",
					Sort:        3,
					Permission:  "/permission/list:GET",
					Description: "权限管理页面",
				},
			},
		},
		{
			MenuKey:     "sms-management",
			Title:       "短信管理",
			Path:        "",
			Icon:        "MessageOutlined",
			ParentID:    "0",
			Sort:        3,
			Permission:  "",
			MenuType:    "1",
			Description: "短信管理模块",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "sms-template",
					Title:       "短信模板",
					Path:        "/sms/template",
					Icon:        "",
					Sort:        1,
					Permission:  "/sms/template/list:GET",
					Description: "短信模板管理页面",
				},
				{
					MenuKey:     "sms-code",
					Title:       "短信验证码",
					Path:        "/sms/code",
					Icon:        "",
					Sort:        2,
					Permission:  "/sms/code/list:GET",
					Description: "短信验证码管理页面",
				},
			},
		},
		{
			MenuKey:     "menu-management",
			Title:       "菜单管理",
			Path:        "",
			Icon:        "MenuOutlined",
			ParentID:    "0",
			Sort:        4,
			Permission:  "",
			MenuType:    "1",
			Description: "菜单管理模块",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "menu-list",
					Title:       "菜单列表",
					Path:        "/menu/list",
					Icon:        "",
					Sort:        1,
					Permission:  "/menu/list:GET",
					Description: "菜单列表页面",
				},
			},
		},
		{
			MenuKey:     "file-management",
			Title:       "文件管理",
			Path:        "",
			Icon:        "FileOutlined",
			ParentID:    "0",
			Sort:        5,
			Permission:  "",
			MenuType:    "1",
			Description: "文件管理模块",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "file-list",
					Title:       "文件列表",
					Path:        "/file/list",
					Icon:        "",
					Sort:        1,
					Permission:  "/file/list:GET",
					Description: "文件列表页面",
				},
				{
					MenuKey:     "file-storage-config",
					Title:       "存储配置",
					Path:        "/file/storage/config",
					Icon:        "",
					Sort:        2,
					Permission:  "/file/storage/config:GET",
					Description: "文件存储配置页面",
				},
			},
		},
		{
			MenuKey:     "system-ops",
			Title:       "系统运维",
			Path:        "",
			Icon:        "ToolOutlined",
			ParentID:    "0",
			Sort:        6,
			Permission:  "",
			MenuType:    "1",
			Description: "系统运维模块",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "api-doc",
					Title:       "接口文档",
					Path:        "/ops/api-doc",
					Icon:        "",
					Sort:        1,
					Permission:  "/swagger:GET",
					Description: "API接口文档页面",
				},
				{
					MenuKey:     "system-monitor",
					Title:       "系统监控",
					Path:        "/ops/monitor",
					Icon:        "",
					Sort:        2,
					Permission:  "/ops/monitor:GET",
					Description: "系统监控页面",
				},
			},
		},
		{
			MenuKey:     "system-settings",
			Title:       "系统设置",
			Path:        "",
			Icon:        "SettingOutlined",
			ParentID:    "0",
			Sort:        7,
			Permission:  "",
			MenuType:    "1",
			Description: "系统设置模块",
			Children: []*struct {
				MenuKey     string
				Title       string
				Path        string
				Icon        string
				Sort        int
				Permission  string
				Description string
			}{
				{
					MenuKey:     "system-config",
					Title:       "系统配置",
					Path:        "/system/config",
					Icon:        "",
					Sort:        1,
					Permission:  "/system/config:GET",
					Description: "系统配置页面",
				},
				{
					MenuKey:     "system-log",
					Title:       "操作日志",
					Path:        "/system/log",
					Icon:        "",
					Sort:        2,
					Permission:  "/system/log:GET",
					Description: "操作日志页面",
				},
			},
		},
	}

	// 创建父菜单和子菜单
	for _, menuDef := range defaultMenus {
		// 检查父菜单是否已存在
		existing, err := menuRepo.FindByKey(ctx, menuDef.MenuKey)
		if err != nil {
			return fmt.Errorf("查询菜单失败: %w", err)
		}

		var parentMenuID string
		if existing == nil {
			// 创建父菜单
			parentMenu := &entity.Menu{
				MenuID:      uuid.New().String(),
				MenuKey:     menuDef.MenuKey,
				Title:       menuDef.Title,
				Path:        menuDef.Path,
				Icon:        menuDef.Icon,
				ParentID:    menuDef.ParentID,
				Sort:        menuDef.Sort,
				Permission:  "", // 菜单显示不再依赖Permission字段，只基于角色授权
				MenuType:    menuDef.MenuType,
				Status:      "1",
				Description: menuDef.Description,
				CreatorID:   "system",
				ModifierID:  "system",
			}

			if err := menuRepo.Create(ctx, parentMenu); err != nil {
				if !isDuplicateError(err) {
					return fmt.Errorf("创建菜单 %s 失败: %w", menuDef.Title, err)
				}
				// 如果已存在，重新查询
				existing, err = menuRepo.FindByKey(ctx, menuDef.MenuKey)
				if err != nil {
					return fmt.Errorf("查询菜单失败: %w", err)
				}
				if existing != nil {
					parentMenuID = existing.MenuID
				}
			} else {
				parentMenuID = parentMenu.MenuID
				logger.Info("创建默认菜单", zap.String("menu", menuDef.Title))
			}
		} else {
			parentMenuID = existing.MenuID
		}

		// 创建子菜单
		for _, childDef := range menuDef.Children {
			existingChild, err := menuRepo.FindByKey(ctx, childDef.MenuKey)
			if err != nil {
				return fmt.Errorf("查询子菜单失败: %w", err)
			}
			if existingChild == nil {
				childMenu := &entity.Menu{
					MenuID:      uuid.New().String(),
					MenuKey:     childDef.MenuKey,
					Title:       childDef.Title,
					Path:        childDef.Path,
					Icon:        childDef.Icon,
					ParentID:    parentMenuID,
					Sort:        childDef.Sort,
					Permission:  "", // 菜单显示不再依赖Permission字段，只基于角色授权
					MenuType:    menuDef.MenuType,
					Status:      "1",
					Description: childDef.Description,
					CreatorID:   "system",
					ModifierID:  "system",
				}

				if err := menuRepo.Create(ctx, childMenu); err != nil {
					if !isDuplicateError(err) {
						return fmt.Errorf("创建子菜单 %s 失败: %w", childDef.Title, err)
					}
				} else {
					logger.Info("创建默认子菜单", zap.String("menu", childDef.Title), zap.String("parent", menuDef.Title))
				}
			}
		}
	}

	// 初始化用户下拉菜单（类型为2）
	defaultUserMenus := []*struct {
		MenuKey     string
		Title       string
		Path        string
		Icon        string
		ParentID    string
		Sort        int
		Permission  string
		MenuType    string
		Description string
	}{
		{
			MenuKey:     "profile",
			Title:       "个人资料",
			Path:        "/profile",
			Icon:        "UserOutlined",
			ParentID:    "0",
			Sort:        1,
			Permission:  "",
			MenuType:    "2",
			Description: "个人资料菜单",
		},
		{
			MenuKey:     "logout",
			Title:       "退出登录",
			Path:        "",
			Icon:        "LogoutOutlined",
			ParentID:    "0",
			Sort:        2,
			Permission:  "divider", // 使用permission字段标记分隔线
			MenuType:    "2",
			Description: "退出登录菜单",
		},
	}

	// 创建用户下拉菜单
	for _, menuDef := range defaultUserMenus {
		existing, err := menuRepo.FindByKey(ctx, menuDef.MenuKey)
		if err != nil {
			return fmt.Errorf("查询用户菜单失败: %w", err)
		}
		if existing == nil {
			userMenu := &entity.Menu{
				MenuID:      uuid.New().String(),
				MenuKey:     menuDef.MenuKey,
				Title:       menuDef.Title,
				Path:        menuDef.Path,
				Icon:        menuDef.Icon,
				ParentID:    menuDef.ParentID,
				Sort:        menuDef.Sort,
				Permission:  menuDef.Permission, // 保留Permission字段，用于特殊标记（如"divider"用于分隔线）
				MenuType:    menuDef.MenuType,
				Status:      "1",
				Description: menuDef.Description,
				CreatorID:   "system",
				ModifierID:  "system",
			}

			if err := menuRepo.Create(ctx, userMenu); err != nil {
				if !isDuplicateError(err) {
					return fmt.Errorf("创建用户菜单 %s 失败: %w", menuDef.Title, err)
				}
			} else {
				logger.Info("创建默认用户菜单", zap.String("menu", menuDef.Title))
			}
		}
	}

	return nil
}

// assignAllMenusToAdmin 为超级管理员角色分配所有菜单
func assignAllMenusToAdmin(
	ctx context.Context,
	roleRepo repository.RoleRepository,
	menuRepo repository.MenuRepository,
) error {
	// 获取超级管理员角色
	superAdminRole, err := roleRepo.FindByKey(ctx, "super_admin")
	if err != nil {
		return fmt.Errorf("查询超级管理员角色失败: %w", err)
	}
	if superAdminRole == nil {
		return fmt.Errorf("超级管理员角色不存在")
	}

	// 获取所有菜单
	allMenus, err := menuRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("查询菜单列表失败: %w", err)
	}

	// 获取角色菜单关联仓库
	roleMenuRepo := infraRepo.NewRoleMenuRepository()

	// 为超级管理员角色分配所有菜单
	assignedCount := 0
	for _, menu := range allMenus {
		// 检查是否已分配
		menuIDs, err := roleMenuRepo.FindMenusByRoleID(ctx, superAdminRole.RoleID)
		if err == nil {
			alreadyAssigned := false
			for _, menuID := range menuIDs {
				if menuID == menu.MenuID {
					alreadyAssigned = true
					break
				}
			}
			if alreadyAssigned {
				continue
			}
		}

		// 分配菜单
		if err := roleMenuRepo.AssignMenu(ctx, superAdminRole.RoleID, menu.MenuID); err != nil {
			logger.Warn("为超级管理员分配菜单失败",
				zap.String("menu", menu.Title),
				zap.String("menu_id", menu.MenuID),
				zap.Error(err))
			continue
		}
		assignedCount++
		logger.Debug("为超级管理员添加菜单",
			zap.String("role", superAdminRole.RoleKey),
			zap.String("menu", menu.Title))
	}

	logger.Info("为超级管理员分配菜单完成",
		zap.String("role", superAdminRole.RoleName),
		zap.Int("total_menus", len(allMenus)),
		zap.Int("assigned_count", assignedCount))

	return nil
}

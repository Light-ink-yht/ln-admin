package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrMenuNotFound = errors.New("菜单不存在")
	ErrMenuExists   = errors.New("菜单已存在")
)

// MenuService 菜单服务
type MenuService struct {
	permissionService *PermissionService
	menuRepo          repository.MenuRepository
	permissionRepo    repository.PermissionRepository
	roleMenuRepo      repository.RoleMenuRepository
}

// NewMenuService 创建菜单服务
func NewMenuService(
	permissionService *PermissionService,
	menuRepo repository.MenuRepository,
	permissionRepo repository.PermissionRepository,
	roleMenuRepo repository.RoleMenuRepository,
) *MenuService {
	return &MenuService{
		permissionService: permissionService,
		menuRepo:          menuRepo,
		permissionRepo:    permissionRepo,
		roleMenuRepo:      roleMenuRepo,
	}
}

// getAllMenus 获取所有菜单（不根据角色过滤）
func (s *MenuService) getAllMenus() []dto.MenuItem {
	return []dto.MenuItem{
		{
			Key:   "dashboard",
			Title: "仪表盘",
			Icon:  "DashboardOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "dashboard-workbench",
					Title:      "工作台",
					Path:       "/",
					Permission: "/dashboard/workbench:GET",
				},
				{
					Key:        "dashboard-analysis",
					Title:      "分析页",
					Path:       "/dashboard/analysis",
					Permission: "/dashboard/analysis:GET",
				},
			},
		},
		{
			Key:   "user-management",
			Title: "用户管理",
			Icon:  "TeamOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "user-list",
					Title:      "用户列表",
					Path:       "/user/list",
					Permission: "/user/list:GET",
				},
				{
					Key:        "role-list",
					Title:      "角色管理",
					Path:       "/role/list",
					Permission: "/role/list:GET",
				},
				{
					Key:        "permission-list",
					Title:      "权限管理",
					Path:       "/permission/list",
					Permission: "/permission/list:GET",
				},
			},
		},
		{
			Key:   "sms-management",
			Title: "短信管理",
			Icon:  "MessageOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "sms-template",
					Title:      "短信模板",
					Path:       "/sms/template",
					Permission: "/sms/template/list:GET",
				},
				{
					Key:        "sms-code",
					Title:      "短信验证码",
					Path:       "/sms/code",
					Permission: "/sms/code/list:GET",
				},
			},
		},
		{
			Key:   "menu-management",
			Title: "菜单管理",
			Icon:  "MenuOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "menu-list",
					Title:      "菜单列表",
					Path:       "/menu/list",
					Permission: "/menu/list:GET",
				},
			},
		},
		{
			Key:   "file-management",
			Title: "文件管理",
			Icon:  "FileOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "file-list",
					Title:      "文件列表",
					Path:       "/file/list",
					Permission: "/file/list:GET",
				},
				{
					Key:        "file-storage-config",
					Title:      "存储配置",
					Path:       "/file/storage/config",
					Permission: "/file/storage/config:GET",
				},
			},
		},
		{
			Key:   "system-ops",
			Title: "系统运维",
			Icon:  "ToolOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "api-doc",
					Title:      "接口文档",
					Path:       "/ops/api-doc",
					Permission: "/swagger:GET",
				},
				{
					Key:        "system-monitor",
					Title:      "系统监控",
					Path:       "/ops/monitor",
					Permission: "/ops/monitor:GET",
				},
			},
		},
		{
			Key:   "system-settings",
			Title: "系统设置",
			Icon:  "SettingOutlined",
			Children: []dto.MenuItem{
				{
					Key:        "system-config",
					Title:      "系统配置",
					Path:       "/system/config",
					Permission: "/system/config:GET",
				},
				{
					Key:        "system-log",
					Title:      "操作日志",
					Path:       "/system/log",
					Permission: "/system/log:GET",
				},
			},
		},
	}
}

// GetSidebarMenus 获取侧边栏菜单（根据用户角色过滤）
// roleKeys: 可选的角色标识列表，如果提供则直接使用，否则从数据库查询
func (s *MenuService) GetSidebarMenus(ctx context.Context, userID string, roleKeys ...[]string) ([]dto.MenuItem, error) {
	var roles []*entity.Role

	// 如果提供了角色列表，直接使用；否则从数据库查询
	if len(roleKeys) > 0 && len(roleKeys[0]) > 0 {
		// 从角色标识列表查询角色实体
		roles = make([]*entity.Role, 0, len(roleKeys[0]))
		for _, roleKey := range roleKeys[0] {
			role, err := s.permissionService.GetRoleByKey(ctx, roleKey)
			if err == nil && role != nil {
				roles = append(roles, role)
			}
		}
	} else {
		// 从数据库查询用户角色
		var err error
		roles, err = s.permissionService.GetUserRoles(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("获取用户角色失败: %w", err)
		}
	}

	// 优先从数据库获取侧边栏菜单（类型为1）
	dbMenus, err := s.menuRepo.FindTreeByType(ctx, "1")
	if err != nil {
		// 如果查询失败，使用默认菜单
		dbMenus = nil
	}

	// 如果数据库中没有菜单，使用默认菜单；否则使用数据库中的菜单
	var allMenus []dto.MenuItem
	var menuIDMap map[string]string = make(map[string]string)
	if len(dbMenus) == 0 {
		// 数据库查询不到，使用默认菜单
		allMenus = s.getAllMenus()
	} else {
		// 数据库中有菜单，优先使用数据库菜单
		allMenus = s.convertMenusToDTOWithID(dbMenus, menuIDMap)
	}

	// 获取Casbin Enforcer
	enforcer := casbin.GetEnforcer()
	if enforcer == nil {
		return nil, fmt.Errorf("权限系统未初始化")
	}

	// 如果是超级管理员，返回所有菜单
	isSuperAdmin := false
	for _, role := range roles {
		if role.RoleKey == "super_admin" {
			isSuperAdmin = true
			break
		}
	}

	if isSuperAdmin {
		return allMenus, nil
	}

	// 根据角色过滤菜单
	filteredMenus := s.filterMenusByRole(ctx, allMenus, roles)

	// 调试日志：记录过滤结果
	logger.Debug("菜单过滤结果",
		zap.String("用户ID", userID),
		zap.Int("原始菜单数量", len(allMenus)),
		zap.Int("过滤后菜单数量", len(filteredMenus)),
		zap.Int("角色数量", len(roles)),
	)

	return filteredMenus, nil
}

// convertMenusToDTO 将菜单实体转换为DTO
func (s *MenuService) convertMenusToDTO(menus []*entity.Menu) []dto.MenuItem {
	return s.convertMenusToDTOWithID(menus, nil)
}

// convertMenusToDTOWithID 将菜单实体转换为DTO，并建立MenuID到MenuKey的映射
func (s *MenuService) convertMenusToDTOWithID(menus []*entity.Menu, menuIDMap map[string]string) []dto.MenuItem {
	if menuIDMap == nil {
		menuIDMap = make(map[string]string)
	}

	result := make([]dto.MenuItem, len(menus))
	for i, menu := range menus {
		// 建立MenuID到MenuKey的映射
		menuIDMap[menu.MenuID] = menu.MenuKey

		children := make([]dto.MenuItem, 0)
		if menu.Children != nil && len(menu.Children) > 0 {
			children = s.convertMenusToDTOWithID(menu.Children, menuIDMap)
		}
		result[i] = dto.MenuItem{
			Key:        menu.MenuKey,
			Title:      menu.Title,
			Path:       menu.Path,
			Icon:       menu.Icon,
			Permission: menu.Permission,
			Hidden:     menu.Status == "2", // 禁用状态视为隐藏
			Children:   children,
		}
	}
	return result
}

// filterMenusByRole 根据角色过滤菜单
func (s *MenuService) filterMenusByRole(
	ctx context.Context,
	menus []dto.MenuItem,
	roles []*entity.Role,
) []dto.MenuItem {
	// 获取所有角色拥有的菜单ID集合
	roleMenuIDs := make(map[string]bool)
	for _, role := range roles {
		menuIDs, err := s.roleMenuRepo.FindMenusByRoleID(ctx, role.RoleID)
		if err == nil {
			for _, menuID := range menuIDs {
				roleMenuIDs[menuID] = true
			}
			logger.Debug("角色菜单关联",
				zap.String("角色ID", role.RoleID),
				zap.String("角色Key", role.RoleKey),
				zap.Int("菜单数量", len(menuIDs)),
				zap.Strings("菜单ID列表", menuIDs),
			)
		} else {
			logger.Warn("查询角色菜单失败",
				zap.String("角色ID", role.RoleID),
				zap.String("角色Key", role.RoleKey),
				zap.Error(err),
			)
		}
	}

	// 如果没有分配任何菜单，返回空列表
	if len(roleMenuIDs) == 0 {
		logger.Debug("角色没有分配任何菜单",
			zap.Int("角色数量", len(roles)),
		)
		return []dto.MenuItem{}
	}

	// 建立菜单Key到MenuID的反向映射（通过查找菜单实体）
	menuKeyToID := make(map[string]string)
	allDBMenus, _ := s.menuRepo.FindAll(ctx)
	for _, menu := range allDBMenus {
		menuKeyToID[menu.MenuKey] = menu.MenuID
	}

	filtered := make([]dto.MenuItem, 0)

	for _, menu := range menus {
		// 通过菜单Key查找MenuID
		menuID, exists := menuKeyToID[menu.Key]

		// 检查当前菜单是否被授权
		isMenuAuthorized := false
		if exists {
			isMenuAuthorized = roleMenuIDs[menuID]
		}

		// 如果有子菜单，先递归过滤子菜单
		var filteredChildren []dto.MenuItem
		if len(menu.Children) > 0 {
			filteredChildren = s.filterMenusByRole(ctx, menu.Children, roles)
		}

		// 如果当前菜单被授权，或者有被授权的子菜单，则显示该菜单
		// 注意：即使菜单Key在映射中找不到（可能是默认菜单），如果有被授权的子菜单，也要显示
		if isMenuAuthorized || len(filteredChildren) > 0 {
			menuCopy := menu
			if len(filteredChildren) > 0 {
				menuCopy.Children = filteredChildren
			} else {
				menuCopy.Children = nil
			}
			filtered = append(filtered, menuCopy)
		}
	}

	return filtered
}

// GetUserMenus 获取用户下拉菜单
// roleKeys: 可选的角色标识列表，如果提供则直接使用，否则从数据库查询
func (s *MenuService) GetUserMenus(ctx context.Context, userID string, roleKeys ...[]string) ([]dto.UserMenuItem, error) {
	var roles []*entity.Role

	// 如果提供了角色列表，直接使用；否则从数据库查询
	if len(roleKeys) > 0 && len(roleKeys[0]) > 0 {
		// 从角色标识列表查询角色实体
		roles = make([]*entity.Role, 0, len(roleKeys[0]))
		for _, roleKey := range roleKeys[0] {
			role, err := s.permissionService.GetRoleByKey(ctx, roleKey)
			if err == nil && role != nil {
				roles = append(roles, role)
			}
		}
	} else {
		// 从数据库查询用户角色
		var err error
		roles, err = s.permissionService.GetUserRoles(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("获取用户角色失败: %w", err)
		}
	}

	// 从数据库获取用户菜单（类型为2）
	dbMenus, err := s.menuRepo.FindByType(ctx, "2")
	if err != nil {
		// 如果查询失败，使用默认菜单
		dbMenus = nil
	}

	// 如果数据库中没有菜单，使用默认菜单
	var menus []dto.UserMenuItem
	if len(dbMenus) == 0 {
		menus = s.getDefaultUserMenus(roles)
	} else {
		// 将数据库菜单转换为DTO，并过滤已禁用的菜单
		menus = make([]dto.UserMenuItem, 0)
		for _, menu := range dbMenus {
			if menu.Status == "1" { // 只返回启用的菜单
				menus = append(menus, dto.UserMenuItem{
					Key:     menu.MenuKey,
					Label:   menu.Title,
					Icon:    menu.Icon,
					Divider: menu.Permission == "divider", // 可以通过permission字段标记分隔线
				})
			}
		}
		// 如果数据库菜单为空，使用默认菜单
		if len(menus) == 0 {
			menus = s.getDefaultUserMenus(roles)
		}
	}

	return menus, nil
}

// getDefaultUserMenus 获取默认用户菜单
func (s *MenuService) getDefaultUserMenus(roles []*entity.Role) []dto.UserMenuItem {
	// 基础菜单项（所有用户都有）
	menus := []dto.UserMenuItem{
		{
			Key:   "profile",
			Label: "个人资料",
			Icon:  "UserOutlined",
		},
		{
			Key:   "settings",
			Label: "设置",
			Icon:  "SettingOutlined",
		},
	}

	// 如果是超级管理员或管理员，添加管理相关菜单
	isAdmin := false
	for _, role := range roles {
		if role.RoleKey == "super_admin" || role.RoleKey == "admin" {
			isAdmin = true
			break
		}
	}

	if isAdmin {
		menus = append(menus, dto.UserMenuItem{
			Key:   "admin",
			Label: "管理中心",
			Icon:  "LaptopOutlined",
		})
	}

	// 添加退出登录（所有用户都有，带分隔线）
	menus = append(menus, dto.UserMenuItem{
		Key:     "logout",
		Label:   "退出登录",
		Icon:    "LogoutOutlined",
		Divider: true,
	})

	return menus
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(ctx context.Context, req *dto.MenuCreateRequest, creatorID string) (*entity.Menu, error) {
	// 检查菜单标识是否已存在
	existing, err := s.menuRepo.FindByKey(ctx, req.MenuKey)
	if err != nil {
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	if existing != nil {
		return nil, ErrMenuExists
	}

	// 如果指定了父菜单，验证父菜单是否存在
	if req.ParentID != "" && req.ParentID != "0" {
		parent, err := s.menuRepo.FindByID(ctx, req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("查询父菜单失败: %w", err)
		}
		if parent == nil {
			return nil, errors.New("父菜单不存在")
		}
	}

	// 创建菜单实体
	menu := &entity.Menu{
		MenuID:      uuid.New().String(),
		MenuKey:     req.MenuKey,
		Title:       req.Title,
		Path:        req.Path,
		Icon:        req.Icon,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
		Permission:  req.Permission,
		MenuType:    req.MenuType,
		Status:      "1", // 默认启用
		Description: req.Description,
		CreatorID:   creatorID,
		ModifierID:  creatorID,
	}

	// 保存到数据库
	if err := s.menuRepo.Create(ctx, menu); err != nil {
		return nil, fmt.Errorf("创建菜单失败: %w", err)
	}

	return menu, nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(ctx context.Context, menuID string, req *dto.MenuUpdateRequest, modifierID string) (*entity.Menu, error) {
	// 查找菜单
	menu, err := s.menuRepo.FindByID(ctx, menuID)
	if err != nil {
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	if menu == nil {
		return nil, ErrMenuNotFound
	}

	// 如果指定了父菜单，验证父菜单是否存在且不是自己
	if req.ParentID != "" && req.ParentID != "0" {
		if req.ParentID == menuID {
			return nil, errors.New("不能将自己设置为父菜单")
		}
		parent, err := s.menuRepo.FindByID(ctx, req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("查询父菜单失败: %w", err)
		}
		if parent == nil {
			return nil, errors.New("父菜单不存在")
		}
	}

	// 更新菜单字段
	menu.Title = req.Title
	menu.Path = req.Path
	menu.Icon = req.Icon
	menu.ParentID = req.ParentID
	menu.Sort = req.Sort
	menu.Permission = req.Permission
	menu.Description = req.Description
	menu.ModifierID = modifierID
	if req.Status != "" {
		menu.Status = req.Status
	}

	// 保存到数据库
	if err := s.menuRepo.Update(ctx, menu); err != nil {
		return nil, fmt.Errorf("更新菜单失败: %w", err)
	}

	return menu, nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(ctx context.Context, menuID string) error {
	// 检查菜单是否存在
	menu, err := s.menuRepo.FindByID(ctx, menuID)
	if err != nil {
		return fmt.Errorf("查询菜单失败: %w", err)
	}
	if menu == nil {
		return ErrMenuNotFound
	}

	// 检查是否有子菜单
	children, err := s.menuRepo.FindByParentID(ctx, menuID)
	if err != nil {
		return fmt.Errorf("查询子菜单失败: %w", err)
	}
	if len(children) > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	// 删除菜单
	if err := s.menuRepo.Delete(ctx, menuID); err != nil {
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	return nil
}

// GetMenu 获取菜单详情
func (s *MenuService) GetMenu(ctx context.Context, menuID string) (*entity.Menu, error) {
	menu, err := s.menuRepo.FindByID(ctx, menuID)
	if err != nil {
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	if menu == nil {
		return nil, ErrMenuNotFound
	}
	return menu, nil
}

// ListMenus 获取菜单列表
func (s *MenuService) ListMenus(ctx context.Context, menuType string) ([]*entity.Menu, error) {
	if menuType != "" {
		return s.menuRepo.FindByType(ctx, menuType)
	}
	return s.menuRepo.FindAll(ctx)
}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree(ctx context.Context, menuType string) ([]*entity.Menu, error) {
	return s.menuRepo.FindTreeByType(ctx, menuType)
}

// ConvertMenuToDTO 将菜单实体转换为响应DTO
func (s *MenuService) ConvertMenuToDTO(menu *entity.Menu) dto.MenuResponse {
	children := make([]dto.MenuResponse, len(menu.Children))
	for i, child := range menu.Children {
		children[i] = s.ConvertMenuToDTO(child)
	}

	return dto.MenuResponse{
		ID:          menu.ID,
		CreatedAt:   menu.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   menu.UpdatedAt.Format(time.RFC3339),
		MenuID:      menu.MenuID,
		MenuKey:     menu.MenuKey,
		Title:       menu.Title,
		Path:        menu.Path,
		Icon:        menu.Icon,
		ParentID:    menu.ParentID,
		Sort:        menu.Sort,
		Permission:  menu.Permission,
		MenuType:    menu.MenuType,
		Status:      menu.Status,
		Description: menu.Description,
		CreatorID:   menu.CreatorID,
		ModifierID:  menu.ModifierID,
		Children:    children,
	}
}

// AssignMenuToRole 为角色分配菜单
func (s *MenuService) AssignMenuToRole(ctx context.Context, roleID, menuID string) error {
	return s.roleMenuRepo.AssignMenu(ctx, roleID, menuID)
}

// DeleteRoleMenus 删除角色的所有菜单关联
func (s *MenuService) DeleteRoleMenus(ctx context.Context, roleID string) error {
	return s.roleMenuRepo.DeleteByRoleID(ctx, roleID)
}

// GetRoleMenus 获取角色的菜单ID列表
func (s *MenuService) GetRoleMenus(ctx context.Context, roleID string) ([]string, error) {
	return s.roleMenuRepo.FindMenusByRoleID(ctx, roleID)
}

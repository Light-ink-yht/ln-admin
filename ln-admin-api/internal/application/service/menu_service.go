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
	"github.com/google/uuid"

	cb "github.com/casbin/casbin/v2"
)

var (
	ErrMenuNotFound = errors.New("菜单不存在")
	ErrMenuExists   = errors.New("菜单已存在")
)

// MenuService 菜单服务
type MenuService struct {
	permissionService *PermissionService
	menuRepo          repository.MenuRepository
}

// NewMenuService 创建菜单服务
func NewMenuService(permissionService *PermissionService, menuRepo repository.MenuRepository) *MenuService {
	return &MenuService{
		permissionService: permissionService,
		menuRepo:          menuRepo,
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
func (s *MenuService) GetSidebarMenus(ctx context.Context, userID string) ([]dto.MenuItem, error) {
	// 获取用户角色
	roles, err := s.permissionService.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 优先从数据库获取侧边栏菜单（类型为1）
	dbMenus, err := s.menuRepo.FindTreeByType(ctx, "1")
	if err != nil {
		// 如果查询失败，使用默认菜单
		dbMenus = nil
	}

	// 如果数据库中没有菜单，使用默认菜单；否则使用数据库中的菜单
	var allMenus []dto.MenuItem
	if len(dbMenus) == 0 {
		// 数据库查询不到，使用默认菜单
		allMenus = s.getAllMenus()
	} else {
		// 数据库中有菜单，优先使用数据库菜单
		allMenus = s.convertMenusToDTO(dbMenus)
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

	// 根据权限过滤菜单
	filteredMenus := s.filterMenusByPermission(ctx, allMenus, userID, enforcer)

	return filteredMenus, nil
}

// convertMenusToDTO 将菜单实体转换为DTO
func (s *MenuService) convertMenusToDTO(menus []*entity.Menu) []dto.MenuItem {
	result := make([]dto.MenuItem, len(menus))
	for i, menu := range menus {
		children := make([]dto.MenuItem, 0)
		if menu.Children != nil && len(menu.Children) > 0 {
			children = s.convertMenusToDTO(menu.Children)
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

// filterMenusByPermission 根据权限过滤菜单
func (s *MenuService) filterMenusByPermission(
	ctx context.Context,
	menus []dto.MenuItem,
	userID string,
	enforcer *cb.Enforcer,
) []dto.MenuItem {
	filtered := make([]dto.MenuItem, 0)

	for _, menu := range menus {
		// 如果有子菜单，递归过滤子菜单
		if len(menu.Children) > 0 {
			filteredChildren := s.filterMenusByPermission(ctx, menu.Children, userID, enforcer)
			if len(filteredChildren) > 0 {
				menu.Children = filteredChildren
				filtered = append(filtered, menu)
			}
		} else {
			// 如果没有子菜单，检查当前菜单项的权限
			if menu.Permission != "" {
				// 解析权限：格式为 "路径:方法"
				parts := splitPermission(menu.Permission)
				if len(parts) == 2 {
					resourcePath := parts[0]
					method := parts[1]

					// 检查用户是否有权限
					allowed, err := enforcer.Enforce(userID, resourcePath, method)
					if err == nil && allowed {
						filtered = append(filtered, menu)
					}
				}
			} else {
				// 没有权限要求，直接添加
				filtered = append(filtered, menu)
			}
		}
	}

	return filtered
}

// splitPermission 分割权限字符串 "路径:方法" -> ["路径", "方法"]
func splitPermission(permission string) []string {
	for i := 0; i < len(permission); i++ {
		if permission[i] == ':' {
			return []string{permission[:i], permission[i+1:]}
		}
	}
	return []string{permission}
}

// GetUserMenus 获取用户下拉菜单
func (s *MenuService) GetUserMenus(ctx context.Context, userID string) ([]dto.UserMenuItem, error) {
	// 获取用户角色
	roles, err := s.permissionService.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
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

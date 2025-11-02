package handler

import (
	"context"
	"strconv"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MenuHandler struct {
	menuService *service.MenuService
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// GetSidebarMenus 获取侧边栏菜单列表
// @Summary      获取侧边栏菜单列表
// @Description  根据当前登录用户的角色获取侧边栏菜单列表，不同角色看到不同的菜单
// @Tags         菜单相关
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.MenuItem}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/sidebar [get]
func (h *MenuHandler) GetSidebarMenus(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("获取侧边栏菜单失败: 未找到用户ID")
		response.Unauthorized(c, "未授权")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Warn("获取侧边栏菜单失败: 用户ID格式错误")
		response.Unauthorized(c, "用户ID格式错误")
		return
	}

	ctx := context.Background()

	logger.Info("开始获取侧边栏菜单",
		zap.String("操作", "获取侧边栏菜单"),
		zap.String("用户ID", userIDStr),
		zap.String("ip", c.ClientIP()))

	// 获取菜单
	menus, err := h.menuService.GetSidebarMenus(ctx, userIDStr)
	if err != nil {
		logger.Error("获取侧边栏菜单失败",
			zap.String("操作", "获取侧边栏菜单"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("用户ID", userIDStr),
			zap.String("ip", c.ClientIP()),
			zap.Error(err))
		response.InternalError(c, "获取菜单失败: "+err.Error())
		return
	}

	logger.Info("获取侧边栏菜单成功",
		zap.String("操作", "获取侧边栏菜单"),
		zap.String("结果", "成功"),
		zap.String("用户ID", userIDStr),
		zap.Int("菜单数量", len(menus)),
		zap.String("ip", c.ClientIP()))

	response.SuccessWithMessage(c, "获取菜单成功", menus)
}

// GetUserMenus 获取用户下拉菜单列表
// @Summary      获取用户下拉菜单列表
// @Description  根据当前登录用户的角色获取用户下拉菜单列表，不同角色看到不同的菜单项
// @Tags         菜单相关
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.UserMenuItem}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/user [get]
func (h *MenuHandler) GetUserMenus(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("获取用户菜单失败: 未找到用户ID")
		response.Unauthorized(c, "未授权")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Warn("获取用户菜单失败: 用户ID格式错误")
		response.Unauthorized(c, "用户ID格式错误")
		return
	}

	ctx := context.Background()

	logger.Info("开始获取用户菜单",
		zap.String("操作", "获取用户菜单"),
		zap.String("用户ID", userIDStr),
		zap.String("ip", c.ClientIP()))

	// 获取菜单
	menus, err := h.menuService.GetUserMenus(ctx, userIDStr)
	if err != nil {
		logger.Error("获取用户菜单失败",
			zap.String("操作", "获取用户菜单"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("用户ID", userIDStr),
			zap.String("ip", c.ClientIP()),
			zap.Error(err))
		response.InternalError(c, "获取菜单失败: "+err.Error())
		return
	}

	logger.Info("获取用户菜单成功",
		zap.String("操作", "获取用户菜单"),
		zap.String("结果", "成功"),
		zap.String("用户ID", userIDStr),
		zap.Int("菜单数量", len(menus)),
		zap.String("ip", c.ClientIP()))

	response.SuccessWithMessage(c, "获取菜单成功", menus)
}

// CreateMenu 创建菜单
// @Summary      创建菜单
// @Description  创建新的菜单项，支持侧边栏菜单和用户菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        menu  body      dto.MenuCreateRequest  true  "菜单信息"
// @Success      200   {object}  response.Response{data=dto.MenuResponse}  "成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      401   {object}  response.Response  "未授权"
// @Failure      500   {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userIDStr := userID.(string)

	var req dto.MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建菜单失败: 参数验证失败", zap.Error(err))
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctx := context.Background()
	logger.Info("开始创建菜单",
		zap.String("操作", "创建菜单"),
		zap.String("菜单标识", req.MenuKey),
		zap.String("用户ID", userIDStr),
		zap.String("ip", c.ClientIP()))

	menu, err := h.menuService.CreateMenu(ctx, &req, userIDStr)
	if err != nil {
		if err == service.ErrMenuExists {
			logger.Warn("创建菜单失败: 菜单已存在",
				zap.String("操作", "创建菜单"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "菜单标识已存在"),
				zap.String("菜单标识", req.MenuKey),
				zap.String("用户ID", userIDStr))
			response.Error(c, 400, "菜单标识已存在")
			return
		}
		logger.Error("创建菜单失败",
			zap.String("操作", "创建菜单"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("用户ID", userIDStr),
			zap.Error(err))
		response.InternalError(c, "创建菜单失败: "+err.Error())
		return
	}

	menuResp := h.menuService.ConvertMenuToDTO(menu)
	logger.Info("创建菜单成功",
		zap.String("操作", "创建菜单"),
		zap.String("结果", "成功"),
		zap.String("菜单ID", menu.MenuID),
		zap.String("用户ID", userIDStr))

	response.SuccessWithMessage(c, "创建菜单成功", menuResp)
}

// UpdateMenu 更新菜单
// @Summary      更新菜单
// @Description  更新菜单信息
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "菜单ID"
// @Param        menu  body      dto.MenuUpdateRequest true  "菜单信息"
// @Success      200   {object}  response.Response{data=dto.MenuResponse}  "成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      401   {object}  response.Response  "未授权"
// @Failure      404   {object}  response.Response  "菜单不存在"
// @Failure      500   {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userIDStr := userID.(string)

	menuID := c.Param("id")
	if menuID == "" {
		response.Error(c, 400, "菜单ID不能为空")
		return
	}

	var req dto.MenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新菜单失败: 参数验证失败", zap.Error(err))
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	ctx := context.Background()
	logger.Info("开始更新菜单",
		zap.String("操作", "更新菜单"),
		zap.String("菜单ID", menuID),
		zap.String("用户ID", userIDStr),
		zap.String("ip", c.ClientIP()))

	menu, err := h.menuService.UpdateMenu(ctx, menuID, &req, userIDStr)
	if err != nil {
		if err == service.ErrMenuNotFound {
			logger.Warn("更新菜单失败: 菜单不存在",
				zap.String("操作", "更新菜单"),
				zap.String("菜单ID", menuID))
			response.Error(c, 404, "菜单不存在")
			return
		}
		logger.Error("更新菜单失败",
			zap.String("操作", "更新菜单"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("菜单ID", menuID),
			zap.Error(err))
		response.InternalError(c, "更新菜单失败: "+err.Error())
		return
	}

	menuResp := h.menuService.ConvertMenuToDTO(menu)
	logger.Info("更新菜单成功",
		zap.String("操作", "更新菜单"),
		zap.String("结果", "成功"),
		zap.String("菜单ID", menuID))

	response.SuccessWithMessage(c, "更新菜单成功", menuResp)
}

// DeleteMenu 删除菜单
// @Summary      删除菜单
// @Description  删除菜单（如果有子菜单则无法删除）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "菜单ID"
// @Success      200   {object}  response.Response  "成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      401   {object}  response.Response  "未授权"
// @Failure      404   {object}  response.Response  "菜单不存在"
// @Failure      500   {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	menuID := c.Param("id")
	if menuID == "" {
		response.Error(c, 400, "菜单ID不能为空")
		return
	}

	ctx := context.Background()
	logger.Info("开始删除菜单",
		zap.String("操作", "删除菜单"),
		zap.String("菜单ID", menuID),
		zap.String("ip", c.ClientIP()))

	err := h.menuService.DeleteMenu(ctx, menuID)
	if err != nil {
		if err == service.ErrMenuNotFound {
			logger.Warn("删除菜单失败: 菜单不存在",
				zap.String("操作", "删除菜单"),
				zap.String("菜单ID", menuID))
			response.Error(c, 404, "菜单不存在")
			return
		}
		logger.Error("删除菜单失败",
			zap.String("操作", "删除菜单"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("菜单ID", menuID),
			zap.Error(err))
		response.InternalError(c, "删除菜单失败: "+err.Error())
		return
	}

	logger.Info("删除菜单成功",
		zap.String("操作", "删除菜单"),
		zap.String("结果", "成功"),
		zap.String("菜单ID", menuID))

	response.Success(c, "删除菜单成功")
}

// GetMenu 获取菜单详情
// @Summary      获取菜单详情
// @Description  根据菜单ID获取菜单详细信息
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "菜单ID"
// @Success      200   {object}  response.Response{data=dto.MenuResponse}  "成功"
// @Failure      401   {object}  response.Response  "未授权"
// @Failure      404   {object}  response.Response  "菜单不存在"
// @Failure      500   {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/{id} [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	menuID := c.Param("id")
	if menuID == "" {
		response.Error(c, 400, "菜单ID不能为空")
		return
	}

	ctx := context.Background()
	logger.Info("开始获取菜单详情",
		zap.String("操作", "获取菜单详情"),
		zap.String("菜单ID", menuID),
		zap.String("ip", c.ClientIP()))

	menu, err := h.menuService.GetMenu(ctx, menuID)
	if err != nil {
		if err == service.ErrMenuNotFound {
			logger.Warn("获取菜单详情失败: 菜单不存在",
				zap.String("操作", "获取菜单详情"),
				zap.String("菜单ID", menuID))
			response.Error(c, 404, "菜单不存在")
			return
		}
		logger.Error("获取菜单详情失败",
			zap.String("操作", "获取菜单详情"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("菜单ID", menuID),
			zap.Error(err))
		response.InternalError(c, "获取菜单详情失败: "+err.Error())
		return
	}

	menuResp := h.menuService.ConvertMenuToDTO(menu)
	logger.Info("获取菜单详情成功",
		zap.String("操作", "获取菜单详情"),
		zap.String("结果", "成功"),
		zap.String("菜单ID", menuID))

	response.SuccessWithMessage(c, "获取菜单详情成功", menuResp)
}

// ListMenus 获取菜单列表
// @Summary      获取菜单列表
// @Description  获取菜单列表，可按菜单类型过滤
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        menu_type  query     string  false  "菜单类型：1 侧边栏菜单 2 用户菜单"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200        {object}  response.Response{data=dto.MenuListResponse}  "成功"
// @Failure      401        {object}  response.Response  "未授权"
// @Failure      500        {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/list [get]
func (h *MenuHandler) ListMenus(c *gin.Context) {
	menuType := c.Query("menu_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	ctx := context.Background()
	logger.Info("开始获取菜单列表",
		zap.String("操作", "获取菜单列表"),
		zap.String("菜单类型", menuType),
		zap.Int("页码", page),
		zap.Int("每页数量", pageSize),
		zap.String("ip", c.ClientIP()))

	menus, err := h.menuService.ListMenus(ctx, menuType)
	if err != nil {
		logger.Error("获取菜单列表失败",
			zap.String("操作", "获取菜单列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.Error(err))
		response.InternalError(c, "获取菜单列表失败: "+err.Error())
		return
	}

	// 转换为DTO
	menuResps := make([]dto.MenuResponse, len(menus))
	for i, menu := range menus {
		menuResps[i] = h.menuService.ConvertMenuToDTO(menu)
	}

	logger.Info("获取菜单列表成功",
		zap.String("操作", "获取菜单列表"),
		zap.String("结果", "成功"),
		zap.Int("菜单数量", len(menuResps)))

	response.SuccessWithMessage(c, "获取菜单列表成功", dto.MenuListResponse{
		List:  menuResps,
		Total: int64(len(menuResps)),
	})
}

// GetMenuTree 获取菜单树
// @Summary      获取菜单树
// @Description  根据菜单类型获取菜单树（包含层级关系）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        menu_type  query     string  true  "菜单类型：1 侧边栏菜单 2 用户菜单"
// @Success      200        {object}  response.Response{data=[]dto.MenuResponse}  "成功"
// @Failure      401        {object}  response.Response  "未授权"
// @Failure      500        {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /menu/tree [get]
func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	menuType := c.Query("menu_type")
	if menuType == "" {
		response.Error(c, 400, "菜单类型不能为空")
		return
	}

	ctx := context.Background()
	logger.Info("开始获取菜单树",
		zap.String("操作", "获取菜单树"),
		zap.String("菜单类型", menuType),
		zap.String("ip", c.ClientIP()))

	menus, err := h.menuService.GetMenuTree(ctx, menuType)
	if err != nil {
		logger.Error("获取菜单树失败",
			zap.String("操作", "获取菜单树"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.Error(err))
		response.InternalError(c, "获取菜单树失败: "+err.Error())
		return
	}

	// 转换为DTO
	menuResps := make([]dto.MenuResponse, len(menus))
	for i, menu := range menus {
		menuResps[i] = h.menuService.ConvertMenuToDTO(menu)
	}

	logger.Info("获取菜单树成功",
		zap.String("操作", "获取菜单树"),
		zap.String("结果", "成功"),
		zap.Int("菜单数量", len(menuResps)))

	response.SuccessWithMessage(c, "获取菜单树成功", menuResps)
}

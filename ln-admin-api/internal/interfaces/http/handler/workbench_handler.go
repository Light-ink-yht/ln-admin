package handler

import (
	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WorkbenchHandler struct {
	workbenchService *service.WorkbenchService
}

// NewWorkbenchHandler 创建工作台配置处理器
func NewWorkbenchHandler(workbenchService *service.WorkbenchService) *WorkbenchHandler {
	return &WorkbenchHandler{
		workbenchService: workbenchService,
	}
}

// GetWorkbenchConfig 获取当前用户的工作台配置
// @Summary      获取工作台配置
// @Description  根据当前登录用户的角色获取对应的工作台配置
// @Tags         工作台相关
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.WorkbenchConfigResponse}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /workbench/config [get]
func (h *WorkbenchHandler) GetWorkbenchConfig(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("获取工作台配置失败: 未找到用户ID")
		response.Unauthorized(c, "未授权")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Warn("获取工作台配置失败: 用户ID格式错误")
		response.Unauthorized(c, "用户ID格式错误")
		return
	}

	// 调用服务获取配置
	config, err := h.workbenchService.GetWorkbenchConfig(c.Request.Context(), userIDStr)
	if err != nil {
		logger.Error("获取工作台配置失败", zap.Error(err))
		response.InternalError(c, "获取工作台配置失败")
		return
	}

	response.SuccessWithMessage(c, "获取工作台配置成功", config)
}

// CreateWorkbenchConfig 创建工作台配置
// @Summary      创建工作台配置
// @Description  为指定角色创建工作台配置
// @Tags         工作台相关
// @Accept       json
// @Produce      json
// @Param        request  body  dto.WorkbenchConfigCreateRequest  true  "创建工作台配置请求"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请求参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /workbench/config [post]
func (h *WorkbenchHandler) CreateWorkbenchConfig(c *gin.Context) {
	var req dto.WorkbenchConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建工作台配置失败: 参数验证失败", zap.Error(err))
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 从上下文获取用户ID作为创建人
	creatorID, exists := c.Get("user_id")
	if !exists {
		creatorID = "system"
	}

	creatorIDStr, _ := creatorID.(string)

	// 调用服务创建配置
	if err := h.workbenchService.CreateWorkbenchConfig(c.Request.Context(), &req, creatorIDStr); err != nil {
		if err == service.ErrWorkbenchConfigExists {
			response.Error(c, 400, "该角色的工作台配置已存在")
			return
		}
		logger.Error("创建工作台配置失败", zap.Error(err))
		response.InternalError(c, "创建工作台配置失败")
		return
	}

	response.SuccessWithMessage(c, "创建工作台配置成功", nil)
}

// UpdateWorkbenchConfig 更新工作台配置
// @Summary      更新工作台配置
// @Description  更新指定工作台配置
// @Tags         工作台相关
// @Accept       json
// @Produce      json
// @Param        request  body  dto.WorkbenchConfigUpdateRequest  true  "更新工作台配置请求"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请求参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /workbench/config [put]
func (h *WorkbenchHandler) UpdateWorkbenchConfig(c *gin.Context) {
	var req dto.WorkbenchConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新工作台配置失败: 参数验证失败", zap.Error(err))
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 从上下文获取用户ID作为修改人
	modifierID, exists := c.Get("user_id")
	if !exists {
		modifierID = "system"
	}

	modifierIDStr, _ := modifierID.(string)

	// 调用服务更新配置
	if err := h.workbenchService.UpdateWorkbenchConfig(c.Request.Context(), &req, modifierIDStr); err != nil {
		if err == service.ErrWorkbenchConfigNotFound {
			response.Error(c, 404, "工作台配置不存在")
			return
		}
		logger.Error("更新工作台配置失败", zap.Error(err))
		response.InternalError(c, "更新工作台配置失败")
		return
	}

	response.SuccessWithMessage(c, "更新工作台配置成功", nil)
}

// DeleteWorkbenchConfig 删除工作台配置
// @Summary      删除工作台配置
// @Description  删除指定工作台配置
// @Tags         工作台相关
// @Accept       json
// @Produce      json
// @Param        config_id  path  string  true  "配置ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      404  {object}  response.Response  "配置不存在"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /workbench/config/{config_id} [delete]
func (h *WorkbenchHandler) DeleteWorkbenchConfig(c *gin.Context) {
	configID := c.Param("config_id")
	if configID == "" {
		response.Error(c, 400, "配置ID不能为空")
		return
	}

	// 调用服务删除配置
	if err := h.workbenchService.DeleteWorkbenchConfig(c.Request.Context(), configID); err != nil {
		if err == service.ErrWorkbenchConfigNotFound {
			response.Error(c, 404, "工作台配置不存在")
			return
		}
		logger.Error("删除工作台配置失败", zap.Error(err))
		response.InternalError(c, "删除工作台配置失败")
		return
	}

	response.SuccessWithMessage(c, "删除工作台配置成功", nil)
}

// GetWorkbenchConfigByID 根据配置ID获取工作台配置
// @Summary      获取工作台配置详情
// @Description  根据配置ID获取工作台配置详情
// @Tags         工作台相关
// @Accept       json
// @Produce      json
// @Param        config_id  path  string  true  "配置ID"
// @Success      200  {object}  response.Response{data=dto.WorkbenchConfigResponse}  "成功"
// @Failure      404  {object}  response.Response  "配置不存在"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /workbench/config/{config_id} [get]
func (h *WorkbenchHandler) GetWorkbenchConfigByID(c *gin.Context) {
	configID := c.Param("config_id")
	if configID == "" {
		response.Error(c, 400, "配置ID不能为空")
		return
	}

	// 调用服务获取配置
	config, err := h.workbenchService.GetWorkbenchConfigByID(c.Request.Context(), configID)
	if err != nil {
		if err == service.ErrWorkbenchConfigNotFound {
			response.Error(c, 404, "工作台配置不存在")
			return
		}
		logger.Error("获取工作台配置失败", zap.Error(err))
		response.InternalError(c, "获取工作台配置失败")
		return
	}

	response.SuccessWithMessage(c, "获取工作台配置成功", config)
}

// ListWorkbenchConfigs 列出所有工作台配置
// @Summary      获取工作台配置列表
// @Description  获取所有工作台配置列表
// @Tags         工作台相关
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.WorkbenchConfigListResponse}  "成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /workbench/config/list [get]
func (h *WorkbenchHandler) ListWorkbenchConfigs(c *gin.Context) {
	// 调用服务获取配置列表
	configs, err := h.workbenchService.ListWorkbenchConfigs(c.Request.Context())
	if err != nil {
		logger.Error("获取工作台配置列表失败", zap.Error(err))
		response.InternalError(c, "获取工作台配置列表失败")
		return
	}

	response.SuccessWithMessage(c, "获取工作台配置列表成功", configs)
}

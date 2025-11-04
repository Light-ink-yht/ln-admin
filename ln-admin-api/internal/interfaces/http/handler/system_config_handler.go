package handler

import (
	"net/http"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemConfigHandler struct {
	configService *service.SystemConfigService
}

// NewSystemConfigHandler 创建系统配置处理器
func NewSystemConfigHandler(configService *service.SystemConfigService) *SystemConfigHandler {
	return &SystemConfigHandler{
		configService: configService,
	}
}

// GetAllConfigs 获取所有配置
// @Summary      获取所有配置
// @Description  获取所有系统配置列表
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  response.Response{data=[]dto.SystemConfigResponse}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /system/config/list [get]
func (h *SystemConfigHandler) GetAllConfigs(c *gin.Context) {
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	logger.Info("开始获取所有系统配置",
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	configs, err := h.configService.GetAllConfigs(c.Request.Context())
	if err != nil {
		logger.Error("获取所有系统配置失败",
			zap.String("操作", "获取所有系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取系统配置失败")
		return
	}

	logger.Info("获取所有系统配置成功",
		zap.String("操作", "获取所有系统配置"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.Int("count", len(configs)),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取系统配置成功", configs)
}

// GetConfigsByGroup 根据分组获取配置
// @Summary      根据分组获取配置
// @Description  根据配置分组获取配置列表
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        group  path      string  true  "配置分组"
// @Success      200    {object}  response.Response{data=[]dto.SystemConfigResponse}  "获取成功"
// @Failure      401    {object}  response.Response  "未授权"
// @Failure      500    {object}  response.Response  "服务器错误"
// @Router       /system/config/group/{group} [get]
func (h *SystemConfigHandler) GetConfigsByGroup(c *gin.Context) {
	group := c.Param("group")
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	logger.Info("开始根据分组获取系统配置",
		zap.String("group", group),
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	configs, err := h.configService.GetConfigsByGroup(c.Request.Context(), group)
	if err != nil {
		logger.Error("根据分组获取系统配置失败",
			zap.String("操作", "根据分组获取系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("group", group),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取系统配置失败")
		return
	}

	logger.Info("根据分组获取系统配置成功",
		zap.String("操作", "根据分组获取系统配置"),
		zap.String("结果", "成功"),
		zap.String("group", group),
		zap.String("user_id", userID),
		zap.Int("count", len(configs)),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取系统配置成功", configs)
}

// GetConfigByKey 根据配置键获取配置
// @Summary      根据配置键获取配置
// @Description  根据配置键获取单个配置详情
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        key  path      string  true  "配置键"
// @Success      200  {object}  response.Response{data=dto.SystemConfigResponse}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "配置不存在"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /system/config/{key} [get]
func (h *SystemConfigHandler) GetConfigByKey(c *gin.Context) {
	key := c.Param("key")
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	logger.Info("开始根据配置键获取系统配置",
		zap.String("key", key),
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	config, err := h.configService.GetConfigByKey(c.Request.Context(), key)
	if err != nil {
		if err == service.ErrSystemConfigNotFound {
			logger.Warn("根据配置键获取系统配置失败",
				zap.String("操作", "根据配置键获取系统配置"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "配置不存在"),
				zap.String("key", key),
				zap.String("user_id", userID),
				zap.String("ip", clientIP))
			response.Error(c, http.StatusNotFound, "配置不存在")
			return
		}
		logger.Error("根据配置键获取系统配置失败",
			zap.String("操作", "根据配置键获取系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("key", key),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取系统配置失败")
		return
	}

	logger.Info("根据配置键获取系统配置成功",
		zap.String("操作", "根据配置键获取系统配置"),
		zap.String("结果", "成功"),
		zap.String("key", key),
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取系统配置成功", config)
}

// CreateConfig 创建配置
// @Summary      创建系统配置
// @Description  创建新的系统配置
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.CreateSystemConfigRequest  true  "创建配置请求"
// @Success      200      {object}  response.Response  "创建成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /system/config [post]
func (h *SystemConfigHandler) CreateConfig(c *gin.Context) {
	var req dto.CreateSystemConfigRequest
	clientIP := c.ClientIP()

	creatorID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			creatorID = uidStr
		}
	}

	logger.Info("开始创建系统配置",
		zap.String("creator_id", creatorID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建系统配置失败",
			zap.String("操作", "创建系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("creator_id", creatorID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.configService.CreateConfig(c.Request.Context(), &req, creatorID); err != nil {
		if err == service.ErrSystemConfigExists {
			logger.Warn("创建系统配置失败",
				zap.String("操作", "创建系统配置"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "配置已存在"),
				zap.String("config_key", req.ConfigKey),
				zap.String("creator_id", creatorID),
				zap.String("ip", clientIP))
			response.Error(c, http.StatusBadRequest, "配置已存在")
			return
		}
		logger.Error("创建系统配置失败",
			zap.String("操作", "创建系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("config_key", req.ConfigKey),
			zap.String("creator_id", creatorID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "创建系统配置失败")
		return
	}

	logger.Info("创建系统配置成功",
		zap.String("操作", "创建系统配置"),
		zap.String("结果", "成功"),
		zap.String("config_key", req.ConfigKey),
		zap.String("creator_id", creatorID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "创建系统配置成功", nil)
}

// UpdateConfig 更新配置
// @Summary      更新系统配置
// @Description  更新系统配置信息
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        key      path      string                        true  "配置键"
// @Param        request  body      dto.UpdateSystemConfigRequest  true  "更新配置请求"
// @Success      200      {object}  response.Response  "更新成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      404      {object}  response.Response  "配置不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /system/config/{key} [put]
func (h *SystemConfigHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")
	var req dto.UpdateSystemConfigRequest
	clientIP := c.ClientIP()

	modifierID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			modifierID = uidStr
		}
	}

	logger.Info("开始更新系统配置",
		zap.String("key", key),
		zap.String("modifier_id", modifierID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新系统配置失败",
			zap.String("操作", "更新系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("key", key),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.configService.UpdateConfig(c.Request.Context(), key, &req, modifierID); err != nil {
		if err == service.ErrSystemConfigNotFound {
			logger.Warn("更新系统配置失败",
				zap.String("操作", "更新系统配置"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "配置不存在"),
				zap.String("key", key),
				zap.String("modifier_id", modifierID),
				zap.String("ip", clientIP))
			response.Error(c, http.StatusNotFound, "配置不存在")
			return
		}
		logger.Error("更新系统配置失败",
			zap.String("操作", "更新系统配置"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("key", key),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "更新系统配置失败")
		return
	}

	logger.Info("更新系统配置成功",
		zap.String("操作", "更新系统配置"),
		zap.String("结果", "成功"),
		zap.String("key", key),
		zap.String("modifier_id", modifierID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "更新系统配置成功", nil)
}

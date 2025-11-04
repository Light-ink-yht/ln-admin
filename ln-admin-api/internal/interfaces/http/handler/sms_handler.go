package handler

import (
	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SMSHandler struct {
	templateService *service.SMSTemplateService
	codeService     *service.SMSCodeService
}

// NewSMSHandler 创建短信处理器
func NewSMSHandler(templateService *service.SMSTemplateService, codeService *service.SMSCodeService) *SMSHandler {
	return &SMSHandler{
		templateService: templateService,
		codeService:     codeService,
	}
}

// GetTemplateList 获取模板列表
// @Summary      获取短信模板列表
// @Description  分页查询短信模板，支持多条件筛选
// @Tags         短信管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page      query     int     false  "页码" default(1)
// @Param        pageSize  query     int     false  "每页数量" default(10)
// @Param        type      query     string  false  "模板类型"
// @Param        status    query     string  false  "状态 (1 启用, 2 禁用)"
// @Success      200  {object}  response.Response{data=object{list=[]dto.SMSTemplateResponse,total=int64}}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /sms/template/list [get]
func (h *SMSHandler) GetTemplateList(c *gin.Context) {
	var req dto.SMSTemplateListRequest
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取短信模板列表失败",
			zap.String("操作", "获取短信模板列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	templates, total, err := h.templateService.GetTemplateList(c.Request.Context(), &req)
	if err != nil {
		logger.Error("获取短信模板列表失败",
			zap.String("操作", "获取短信模板列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取短信模板列表失败")
		return
	}

	response.SuccessWithMessage(c, "获取短信模板列表成功", gin.H{
		"list":  templates,
		"total": total,
	})
}

// GetTemplateDetail 获取模板详情
// @Summary      获取短信模板详情
// @Description  根据模板ID获取短信模板详细信息
// @Tags         短信管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        templateId  path      string  true  "模板ID"
// @Success      200  {object}  response.Response{data=dto.SMSTemplateResponse}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "模板不存在"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /sms/template/{templateId} [get]
func (h *SMSHandler) GetTemplateDetail(c *gin.Context) {
	templateID := c.Param("templateId")
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	template, err := h.templateService.GetTemplateByID(c.Request.Context(), templateID)
	if err != nil {
		logger.Error("获取短信模板详情失败",
			zap.String("操作", "获取短信模板详情"),
			zap.String("结果", "失败"),
			zap.String("template_id", templateID),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取短信模板详情失败")
		return
	}

	if template == nil {
		response.NotFound(c, "模板不存在")
		return
	}

	response.SuccessWithMessage(c, "获取短信模板详情成功", template)
}

// CreateTemplate 创建模板
// @Summary      创建短信模板
// @Description  创建新的短信模板
// @Tags         短信管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.CreateSMSTemplateRequest  true  "创建模板请求"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "请求参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /sms/template [post]
func (h *SMSHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateSMSTemplateRequest
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建短信模板失败",
			zap.String("操作", "创建短信模板"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.templateService.CreateTemplate(c.Request.Context(), &req, userID); err != nil {
		logger.Error("创建短信模板失败",
			zap.String("操作", "创建短信模板"),
			zap.String("结果", "失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	logger.Info("创建短信模板成功",
		zap.String("操作", "创建短信模板"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	response.Success(c, "创建短信模板成功")
}

// UpdateTemplate 更新模板
// @Summary      更新短信模板
// @Description  更新短信模板信息
// @Tags         短信管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        templateId  path      string  true  "模板ID"
// @Param        request     body      dto.UpdateSMSTemplateRequest  true  "更新模板请求"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "请求参数错误"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "模板不存在"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /sms/template/{templateId} [put]
func (h *SMSHandler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	var req dto.UpdateSMSTemplateRequest
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新短信模板失败",
			zap.String("操作", "更新短信模板"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("template_id", templateID),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.templateService.UpdateTemplate(c.Request.Context(), templateID, &req, userID); err != nil {
		logger.Error("更新短信模板失败",
			zap.String("操作", "更新短信模板"),
			zap.String("结果", "失败"),
			zap.String("template_id", templateID),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	logger.Info("更新短信模板成功",
		zap.String("操作", "更新短信模板"),
		zap.String("结果", "成功"),
		zap.String("template_id", templateID),
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	response.Success(c, "更新短信模板成功")
}

// DeleteTemplate 删除模板
// @Summary      删除短信模板
// @Description  删除指定的短信模板
// @Tags         短信管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        templateId  path      string  true  "模板ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      404  {object}  response.Response  "模板不存在"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /sms/template/{templateId} [delete]
func (h *SMSHandler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if err := h.templateService.DeleteTemplate(c.Request.Context(), templateID); err != nil {
		logger.Error("删除短信模板失败",
			zap.String("操作", "删除短信模板"),
			zap.String("结果", "失败"),
			zap.String("template_id", templateID),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "删除短信模板失败")
		return
	}

	logger.Info("删除短信模板成功",
		zap.String("操作", "删除短信模板"),
		zap.String("结果", "成功"),
		zap.String("template_id", templateID),
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	response.Success(c, "删除短信模板成功")
}

// GetCodeList 获取验证码列表
// @Summary      获取短信验证码列表
// @Description  分页查询短信验证码，支持多条件筛选
// @Tags         短信管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page      query     int     false  "页码" default(1)
// @Param        pageSize  query     int     false  "每页数量" default(10)
// @Param        phone     query     string  false  "手机号"
// @Param        type      query     string  false  "验证码类型"
// @Param        status    query     string  false  "状态 (1 未使用, 2 已使用, 3 已过期)"
// @Success      200  {object}  response.Response{data=object{list=[]dto.SMSCodeResponse,total=int64}}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /sms/code/list [get]
func (h *SMSHandler) GetCodeList(c *gin.Context) {
	var req dto.SMSCodeListRequest
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取短信验证码列表失败",
			zap.String("操作", "获取短信验证码列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	codes, total, err := h.codeService.GetCodeList(c.Request.Context(), &req)
	if err != nil {
		logger.Error("获取短信验证码列表失败",
			zap.String("操作", "获取短信验证码列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取短信验证码列表失败")
		return
	}

	response.SuccessWithMessage(c, "获取短信验证码列表成功", gin.H{
		"list":  codes,
		"total": total,
	})
}


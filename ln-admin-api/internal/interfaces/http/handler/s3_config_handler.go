package handler

import (
	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type S3ConfigHandler struct {
	s3ConfigService *service.S3ConfigService
}

// NewS3ConfigHandler 创建S3配置处理器
func NewS3ConfigHandler(s3ConfigService *service.S3ConfigService) *S3ConfigHandler {
	return &S3ConfigHandler{
		s3ConfigService: s3ConfigService,
	}
}

// GetS3Config 获取S3配置
// @Summary      获取S3配置
// @Description  获取当前的文件存储配置（本地或S3）
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.S3ConfigResponse}  "成功"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /file/storage/config [get]
func (h *S3ConfigHandler) GetS3Config(c *gin.Context) {
	config, err := h.s3ConfigService.GetS3Config(c.Request.Context())
	if err != nil {
		logger.Error("获取S3配置失败",
			zap.String("操作", "获取S3配置"),
			zap.String("结果", "失败"),
			zap.Error(err))
		response.InternalError(c, "获取S3配置失败: "+err.Error())
		return
	}

	// 如果没有配置，返回默认本地存储配置
	if config == nil {
		config = &dto.S3ConfigResponse{
			StorageType: "local",
			Status:      "1",
		}
	}

	response.Success(c, config)
}

// CreateOrUpdateS3Config 创建或更新S3配置
// @Summary      创建或更新S3配置
// @Description  创建或更新文件存储配置（本地或S3）。如果配置已存在则更新，否则创建新配置。
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        request  body      dto.S3ConfigRequest  true  "S3配置请求"
// @Success      200      {object}  response.Response{data=dto.S3ConfigResponse}  "成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /file/storage/config [post]
func (h *S3ConfigHandler) CreateOrUpdateS3Config(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.S3ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证S3配置必填字段
	if req.StorageType == "s3" {
		if req.Bucket == "" {
			response.BadRequest(c, "存储桶名称不能为空")
			return
		}
		if req.Region == "" {
			response.BadRequest(c, "区域不能为空")
			return
		}
		if req.AccessKeyID == "" {
			response.BadRequest(c, "访问密钥ID不能为空")
			return
		}
		if req.SecretKey == "" {
			response.BadRequest(c, "访问密钥不能为空")
			return
		}
	}

	config, err := h.s3ConfigService.CreateOrUpdateS3Config(c.Request.Context(), &req, userID.(string))
	if err != nil {
		logger.Error("保存S3配置失败",
			zap.String("操作", "保存S3配置"),
			zap.String("结果", "失败"),
			zap.String("user_id", userID.(string)),
			zap.Error(err))
		response.InternalError(c, "保存S3配置失败: "+err.Error())
		return
	}

	logger.Info("保存S3配置成功",
		zap.String("操作", "保存S3配置"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID.(string)),
		zap.String("storage_type", config.StorageType))

	response.SuccessWithMessage(c, "保存配置成功", config)
}

// DeleteS3Config 删除S3配置
// @Summary      删除S3配置
// @Description  删除S3配置（软删除）
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        config_id  path  string  true  "配置ID"
// @Success      200        {object}  response.Response  "删除成功"
// @Failure      400        {object}  response.Response  "请求参数错误"
// @Failure      500        {object}  response.Response  "服务器错误"
// @Router       /file/storage/config/:config_id [delete]
func (h *S3ConfigHandler) DeleteS3Config(c *gin.Context) {
	configID := c.Param("config_id")
	if configID == "" {
		response.BadRequest(c, "配置ID不能为空")
		return
	}

	err := h.s3ConfigService.DeleteS3Config(c.Request.Context(), configID)
	if err != nil {
		logger.Error("删除S3配置失败",
			zap.String("操作", "删除S3配置"),
			zap.String("结果", "失败"),
			zap.String("config_id", configID),
			zap.Error(err))
		response.InternalError(c, "删除S3配置失败: "+err.Error())
		return
	}

	logger.Info("删除S3配置成功",
		zap.String("操作", "删除S3配置"),
		zap.String("结果", "成功"),
		zap.String("config_id", configID))

	response.SuccessWithMessage(c, "删除配置成功", nil)
}

package handler

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemMonitorHandler struct {
	monitorService *service.SystemMonitorService
}

// NewSystemMonitorHandler 创建系统监控处理器
func NewSystemMonitorHandler(monitorService *service.SystemMonitorService) *SystemMonitorHandler {
	return &SystemMonitorHandler{
		monitorService: monitorService,
	}
}

// GetSystemMonitor 获取系统监控信息
// @Summary      获取系统监控信息
// @Description  获取服务器的CPU、内存、磁盘等监控信息
// @Tags         系统运维
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.SystemMonitorResponse}  "成功"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Security     Bearer
// @Router       /ops/monitor [get]
func (h *SystemMonitorHandler) GetSystemMonitor(c *gin.Context) {
	ctx := context.Background()
	logger.Info("开始获取系统监控信息",
		zap.String("操作", "获取系统监控信息"),
		zap.String("ip", c.ClientIP()))

	monitorInfo, err := h.monitorService.GetSystemMonitor(ctx)
	if err != nil {
		logger.Error("获取系统监控信息失败",
			zap.String("操作", "获取系统监控信息"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.Error(err))
		response.InternalError(c, "获取系统监控信息失败: "+err.Error())
		return
	}

	logger.Info("获取系统监控信息成功",
		zap.String("操作", "获取系统监控信息"),
		zap.String("结果", "成功"))

	response.SuccessWithMessage(c, "获取系统监控信息成功", monitorInfo)
}

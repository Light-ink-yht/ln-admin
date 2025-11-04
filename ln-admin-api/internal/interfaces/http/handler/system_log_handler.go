package handler

import (
	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemLogHandler struct {
	logService *service.SystemLogService
}

// NewSystemLogHandler 创建系统日志处理器
func NewSystemLogHandler(logService *service.SystemLogService) *SystemLogHandler {
	return &SystemLogHandler{
		logService: logService,
	}
}

// GetLogList 获取日志列表
// @Summary      获取系统日志列表
// @Description  分页查询系统操作日志，支持多条件筛选
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page       query     int     false  "页码" default(1)
// @Param        pageSize   query     int     false  "每页数量" default(10)
// @Param        level      query     string  false  "日志级别 (info, warn, error, debug)"
// @Param        module     query     string  false  "模块名称"
// @Param        action     query     string  false  "操作类型"
// @Param        userId     query     string  false  "用户ID"
// @Param        ip         query     string  false  "IP地址"
// @Param        startTime  query     string  false  "开始时间 (格式: 2006-01-02 15:04:05)"
// @Param        endTime    query     string  false  "结束时间 (格式: 2006-01-02 15:04:05)"
// @Success      200  {object}  response.Response{data=object{list=[]dto.SystemLogResponse,total=int64}}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /system/log/list [get]
func (h *SystemLogHandler) GetLogList(c *gin.Context) {
	var req dto.SystemLogListRequest
	clientIP := c.ClientIP()

	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	logger.Info("开始获取系统日志列表",
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取系统日志列表失败",
			zap.String("操作", "获取系统日志列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	logs, total, err := h.logService.GetLogList(c.Request.Context(), &req)
	if err != nil {
		logger.Error("获取系统日志列表失败",
			zap.String("操作", "获取系统日志列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取系统日志列表失败")
		return
	}

	logger.Info("获取系统日志列表成功",
		zap.String("操作", "获取系统日志列表"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.Int64("total", total),
		zap.Int("count", len(logs)),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取系统日志列表成功", gin.H{
		"list":  logs,
		"total": total,
	})
}

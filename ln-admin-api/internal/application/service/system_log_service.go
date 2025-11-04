package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
)

// SystemLogService 系统日志服务
type SystemLogService struct {
	logRepo repository.SystemLogRepository
}

// NewSystemLogService 创建系统日志服务
func NewSystemLogService(logRepo repository.SystemLogRepository) *SystemLogService {
	return &SystemLogService{
		logRepo: logRepo,
	}
}

// GetLogList 获取日志列表
func (s *SystemLogService) GetLogList(ctx context.Context, req *dto.SystemLogListRequest) ([]*dto.SystemLogResponse, int64, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.Level != "" {
		conditions["level"] = req.Level
	}
	if req.Module != "" {
		conditions["module"] = req.Module
	}
	if req.Action != "" {
		conditions["action"] = req.Action
	}
	if req.UserID != "" {
		conditions["user_id"] = req.UserID
	}
	if req.IP != "" {
		conditions["ip"] = req.IP
	}

	// 处理时间范围
	if req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err == nil {
			conditions["start_time"] = startTime
		}
	}
	if req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err == nil {
			conditions["end_time"] = endTime
		}
	}

	// 查询日志
	logs, total, err := s.logRepo.List(ctx, page, pageSize, conditions)
	if err != nil {
		return nil, 0, fmt.Errorf("查询日志列表失败: %w", err)
	}

	// 转换为响应DTO
	responses := make([]*dto.SystemLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = s.toLogResponse(log)
	}

	return responses, total, nil
}

// toLogResponse 转换为响应DTO
func (s *SystemLogService) toLogResponse(log *entity.SystemLog) *dto.SystemLogResponse {
	logTime := ""
	if log.LogTime != nil {
		logTime = log.LogTime.Format("2006-01-02 15:04:05")
	}

	return &dto.SystemLogResponse{
		LogID:      log.LogID,
		Level:      log.Level,
		Module:     log.Module,
		Action:     log.Action,
		Content:    log.Content,
		UserID:     log.UserID,
		IP:         log.IP,
		Path:       log.Path,
		Method:     log.Method,
		StatusCode: log.StatusCode,
		UserAgent:  log.UserAgent,
		ErrorMsg:   log.ErrorMsg,
		LogTime:    logTime,
		CreatedAt:  dto.FormatTime(log.CreatedAt),
	}
}

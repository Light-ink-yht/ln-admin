package log

import (
	"context"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"

	"github.com/google/uuid"
)

type DBLogger struct {
	logRepo repository.SystemLogRepository
}

// NewDBLogger 创建数据库日志服务
func NewDBLogger(logRepo repository.SystemLogRepository) *DBLogger {
	return &DBLogger{
		logRepo: logRepo,
	}
}

// Log 记录日志到数据库
func (l *DBLogger) Log(ctx context.Context, level, module, action, content, userID, ip, path, method string, statusCode int, userAgent, errorMsg string) error {
	logID := uuid.New().String()
	now := time.Now()

	log := &entity.SystemLog{
		LogID:      logID,
		Level:      level,
		Module:     module,
		Action:     action,
		Content:    content,
		UserID:     userID,
		IP:         ip,
		Path:       path,
		Method:     method,
		StatusCode: statusCode,
		UserAgent:  userAgent,
		ErrorMsg:   errorMsg,
		LogTime:    &now,
	}

	if err := l.logRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("记录日志失败: %w", err)
	}

	return nil
}

// LogInfo 记录Info级别日志
func (l *DBLogger) LogInfo(ctx context.Context, module, action, content, userID, ip, path, method string, statusCode int, userAgent string) error {
	return l.Log(ctx, "info", module, action, content, userID, ip, path, method, statusCode, userAgent, "")
}

// LogError 记录Error级别日志
func (l *DBLogger) LogError(ctx context.Context, module, action, content, userID, ip, path, method string, statusCode int, userAgent, errorMsg string) error {
	return l.Log(ctx, "error", module, action, content, userID, ip, path, method, statusCode, userAgent, errorMsg)
}

// LogWarn 记录Warn级别日志
func (l *DBLogger) LogWarn(ctx context.Context, module, action, content, userID, ip, path, method string, statusCode int, userAgent string) error {
	return l.Log(ctx, "warn", module, action, content, userID, ip, path, method, statusCode, userAgent, "")
}

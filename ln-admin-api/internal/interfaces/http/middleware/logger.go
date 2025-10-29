package middleware

import (
	"context"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	logService "github.com/Light-ink-yht/ln-admin/pkg/service/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var dbLogger *logService.DBLogger

// SetDBLogger 设置数据库日志服务（在启动时调用）
func SetDBLogger(dbLog *logService.DBLogger) {
	dbLogger = dbLog
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ip),
			zap.Duration("latency", latency),
			zap.String("user_agent", userAgent),
		}

		// 写入文件日志
		if status >= 500 {
			logger.Error("HTTP请求错误", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP请求警告", fields...)
		} else {
			logger.Info("HTTP请求", fields...)
		}

		// 写入数据库日志（异步，不阻塞请求）
		if dbLogger != nil {
			go func() {
				ctx := context.Background()
				// 获取用户ID（如果有）
				userID := ""
				if uid, exists := c.Get("user_id"); exists {
					if uidStr, ok := uid.(string); ok {
						userID = uidStr
					}
				}

				// 根据状态码确定日志级别和操作类型
				var level, action string
				if status >= 500 {
					level = "error"
					action = "error"
				} else if status >= 400 {
					level = "warn"
					action = "warning"
				} else {
					level = "info"
					action = "request"
				}

				// 构建日志内容
				content := method + " " + path
				if query != "" {
					content += "?" + query
				}
				if latency > 0 {
					content += " - " + latency.String()
				}

				errorMsg := ""
				if status >= 500 {
					errorMsg = "HTTP请求错误"
				}

				// 记录到数据库
				if err := dbLogger.Log(ctx, level, "http", action, content, userID, ip, path, method, status, userAgent, errorMsg); err != nil {
					// 数据库日志写入失败，只记录到文件日志，不中断流程
					logger.Warn("数据库日志写入失败", zap.Error(err))
				}
			}()
		}
	}
}

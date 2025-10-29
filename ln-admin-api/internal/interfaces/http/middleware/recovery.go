package middleware

import (
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 恢复中间件（捕获panic）
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", zap.Any("error", err), zap.String("path", c.Request.URL.Path))
				response.InternalError(c, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

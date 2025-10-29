package middleware

import (
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"
	"github.com/Light-ink-yht/ln-admin/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证Token")
			c.Abort()
			return
		}

		// 检查Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		// 解析Token
		token := parts[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			if err == utils.ErrExpiredToken {
				response.Unauthorized(c, "Token已过期")
			} else {
				response.Unauthorized(c, "Token无效")
			}
			logger.Warn("Token验证失败", zap.Error(err))
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

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
		var token string

		// 优先从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 检查Token格式
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果Header中没有token，尝试从query参数获取（用于Swagger等场景）
		if token == "" {
			token = c.Query("token")
		}

		// 如果还是没有token，尝试从cookie获取
		if token == "" {
			cookieToken, err := c.Cookie("token")
			if err == nil && cookieToken != "" {
				token = cookieToken
			}
		}

		// 如果仍然没有token，返回错误
		if token == "" {
			response.Unauthorized(c, "未提供认证Token")
			c.Abort()
			return
		}

		// 解析Token
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

		// 将用户ID和角色信息存储到上下文中
		c.Set("user_id", claims.UserID)
		if len(claims.Roles) > 0 {
			c.Set("user_roles", claims.Roles)
		}
		c.Next()
	}
}

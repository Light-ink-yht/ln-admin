package middleware

import (
	"net/http"
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// CasbinMiddleware Casbin权限验证中间件
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID（从JWT中间件中获取）
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未授权")
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			response.Unauthorized(c, "用户ID格式错误")
			c.Abort()
			return
		}

		// 获取请求路径和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 移除路径前缀（如果有）
		obj = strings.TrimPrefix(obj, "/api")

		// 使用Casbin检查权限
		enforcer := casbin.GetEnforcer()
		if enforcer == nil {
			response.InternalError(c, "权限系统未初始化")
			c.Abort()
			return
		}

		// 检查权限：用户是否有权限访问该资源
		// r.sub = 用户ID, r.obj = 资源路径, r.act = 请求方法
		allowed, err := enforcer.Enforce(userIDStr, obj, act)
		if err != nil {
			response.InternalError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !allowed {
			response.Error(c, http.StatusForbidden, "没有权限访问该资源")
			c.Abort()
			return
		}

		// 权限验证通过，继续处理
		c.Next()
	}
}

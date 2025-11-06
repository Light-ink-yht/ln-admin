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

		// 优先从context中获取角色信息（从JWT token中解析）
		var userRoles []string
		if rolesFromContext, exists := c.Get("user_roles"); exists {
			if rolesList, ok := rolesFromContext.([]string); ok {
				userRoles = rolesList
			}
		}

		// 如果context中没有角色信息，从Casbin中获取（兼容旧token）
		if len(userRoles) == 0 {
			rolesFromCasbin, err := enforcer.GetRolesForUser(userIDStr)
			if err == nil {
				userRoles = rolesFromCasbin
			}
		}

		// 检查用户是否是超级管理员，如果是则直接放行
		for _, role := range userRoles {
			if role == "super_admin" {
				// 超级管理员拥有所有权限，直接放行
				c.Next()
				return
			}
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

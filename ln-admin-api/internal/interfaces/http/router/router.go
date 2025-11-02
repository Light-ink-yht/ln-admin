package router

import (
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/repository"
	"github.com/Light-ink-yht/ln-admin/internal/interfaces/http/handler"
	"github.com/Light-ink-yht/ln-admin/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 注意：需要先执行 swag init -g cmd/server/main.go -o ./docs/swagger 生成文档
	// swagger 生成的文档必须导入以触发 init 函数
	_ "github.com/Light-ink-yht/ln-admin/docs/swagger"
)

// SetupRouter 设置路由
func SetupRouter(userAppService *service.UserAppService, smsAppService *service.SMSAppService, permissionService *service.PermissionService) *gin.Engine {
	// 创建菜单仓库和服务
	menuRepo := repository.NewMenuRepository()
	menuService := service.NewMenuService(permissionService, menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)
	// 创建Gin引擎
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 创建用户处理器
	userHandler := handler.NewUserHandler(userAppService, smsAppService)

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关路由（不需要认证）
		user := api.Group("/user")
		{
			user.GET("/captcha", userHandler.GetCaptcha)          // 获取图片验证码
			user.POST("/login", userHandler.Login)                // 登录
			user.POST("/signup", userHandler.Register)            // 注册
			user.POST("/password", userHandler.ForgotPassword)    // 忘记密码
			user.POST("/signup/code", userHandler.SendSms)        // 发送短信验证码
			user.POST("/refresh-token", userHandler.RefreshToken) // 刷新Token
		}

		// 需要认证但不需权限验证的路由（基础用户信息接口）
		auth := api.Group("/user")
		auth.Use(middleware.Auth()) // 只需JWT认证
		{
			auth.GET("/userinfo", userHandler.GetUserInfo) // 获取当前用户信息（所有已登录用户都可访问）
			auth.POST("/logout", func(c *gin.Context) {    // 退出登录
				c.JSON(200, gin.H{"code": 200, "message": "success"})
			})
		}

		// 菜单相关路由（需要认证）
		menu := api.Group("/menu")
		menu.Use(middleware.Auth()) // 只需JWT认证
		{
			menu.GET("/sidebar", menuHandler.GetSidebarMenus) // 获取侧边栏菜单（动态）
			menu.GET("/user", menuHandler.GetUserMenus)       // 获取用户下拉菜单（动态）
			menu.GET("/list", menuHandler.ListMenus)          // 获取菜单列表
			menu.GET("/tree", menuHandler.GetMenuTree)        // 获取菜单树
			menu.POST("", menuHandler.CreateMenu)             // 创建菜单
			menu.GET("/:id", menuHandler.GetMenu)             // 获取菜单详情（必须放在最后，避免与/list等冲突）
			menu.PUT("/:id", menuHandler.UpdateMenu)          // 更新菜单
			menu.DELETE("/:id", menuHandler.DeleteMenu)       // 删除菜单
		}

		// 需要认证和权限验证的路由（管理功能）
		admin := api.Group("/user")
		admin.Use(middleware.Auth())             // JWT认证
		admin.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			admin.GET("/list", userHandler.ListUsers) // 获取用户列表（需要权限）
		}
	}

	// Swagger文档路由
	// 配置 Swagger UI，指定文档路径为 /swagger/doc.json
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

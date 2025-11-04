package router

import (
	appService "github.com/Light-ink-yht/ln-admin/internal/application/service"
	domainService "github.com/Light-ink-yht/ln-admin/internal/domain/service"
	infraRepo "github.com/Light-ink-yht/ln-admin/internal/infrastructure/repository"
	infraStorage "github.com/Light-ink-yht/ln-admin/internal/infrastructure/storage"
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
func SetupRouter(userAppService *appService.UserAppService, smsAppService *appService.SMSAppService, permissionService *appService.PermissionService) *gin.Engine {
	// 创建菜单仓库和服务
	menuRepo := infraRepo.NewMenuRepository()
	menuService := appService.NewMenuService(permissionService, menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)
	// 创建用户仓库（用于handler）
	userRepo := infraRepo.NewUserRepository()

	// 创建Gin引擎
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 创建处理器
	userHandler := handler.NewUserHandler(userAppService, smsAppService)
	roleHandler := handler.NewRoleHandler(permissionService, userRepo)
	permissionHandler := handler.NewPermissionHandler(permissionService, userRepo)

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
		menu.Use(middleware.Auth()) // JWT认证
		{
			menu.GET("/sidebar", menuHandler.GetSidebarMenus) // 获取侧边栏菜单（动态）
			menu.GET("/user", menuHandler.GetUserMenus)       // 获取用户下拉菜单（动态）
		}

		// 菜单管理相关路由（需要认证和权限验证）
		menuManagement := api.Group("/menu")
		menuManagement.Use(middleware.Auth())             // JWT认证
		menuManagement.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			menuManagement.GET("/list", menuHandler.ListMenus)    // 获取菜单列表
			menuManagement.GET("/tree", menuHandler.GetMenuTree)  // 获取菜单树
			menuManagement.POST("", menuHandler.CreateMenu)       // 创建菜单
			menuManagement.GET("/:id", menuHandler.GetMenu)       // 获取菜单详情（必须放在最后，避免与/list等冲突）
			menuManagement.PUT("/:id", menuHandler.UpdateMenu)    // 更新菜单
			menuManagement.DELETE("/:id", menuHandler.DeleteMenu) // 删除菜单
		}

		// 工作台相关路由（需要认证）
		workbenchRepo := infraRepo.NewWorkbenchRepository()
		workbenchService := appService.NewWorkbenchService(workbenchRepo, permissionService)
		workbenchHandler := handler.NewWorkbenchHandler(workbenchService)
		workbench := api.Group("/workbench")
		workbench.Use(middleware.Auth()) // 只需JWT认证
		{
			workbench.GET("/config", workbenchHandler.GetWorkbenchConfig)                  // 获取当前用户的工作台配置
			workbench.GET("/config/list", workbenchHandler.ListWorkbenchConfigs)           // 获取工作台配置列表
			workbench.GET("/config/:config_id", workbenchHandler.GetWorkbenchConfigByID)   // 获取工作台配置详情
			workbench.POST("/config", workbenchHandler.CreateWorkbenchConfig)              // 创建工作台配置
			workbench.PUT("/config", workbenchHandler.UpdateWorkbenchConfig)               // 更新工作台配置
			workbench.DELETE("/config/:config_id", workbenchHandler.DeleteWorkbenchConfig) // 删除工作台配置
		}

		// 短信管理相关路由（需要认证和权限验证）
		templateRepo := infraRepo.NewSMSTemplateRepository()
		templateService := appService.NewSMSTemplateService(templateRepo)
		codeRepo := infraRepo.NewSMSCodeRepository()
		codeService := appService.NewSMSCodeService(codeRepo)
		smsHandler := handler.NewSMSHandler(templateService, codeService)
		sms := api.Group("/sms")
		sms.Use(middleware.Auth())             // JWT认证
		sms.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			sms.GET("/template/list", smsHandler.GetTemplateList)          // 获取短信模板列表
			sms.GET("/template/:templateId", smsHandler.GetTemplateDetail) // 获取短信模板详情
			sms.POST("/template", smsHandler.CreateTemplate)               // 创建短信模板
			sms.PUT("/template/:templateId", smsHandler.UpdateTemplate)    // 更新短信模板
			sms.DELETE("/template/:templateId", smsHandler.DeleteTemplate) // 删除短信模板
			sms.GET("/code/list", smsHandler.GetCodeList)                  // 获取短信验证码列表
		}

		// 系统配置和日志相关路由（需要认证和权限验证）
		configRepo := infraRepo.NewSystemConfigRepository()
		configService := appService.NewSystemConfigService(configRepo)
		configHandler := handler.NewSystemConfigHandler(configService)
		logRepo := infraRepo.NewSystemLogRepository()
		logService := appService.NewSystemLogService(logRepo)
		logHandler := handler.NewSystemLogHandler(logService)
		system := api.Group("/system")
		system.Use(middleware.Auth())             // JWT认证
		system.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			system.GET("/config/list", configHandler.GetAllConfigs)             // 获取所有配置
			system.GET("/config/group/:group", configHandler.GetConfigsByGroup) // 根据分组获取配置
			system.GET("/config/:key", configHandler.GetConfigByKey)            // 根据配置键获取配置
			system.POST("/config", configHandler.CreateConfig)                  // 创建配置
			system.PUT("/config/:key", configHandler.UpdateConfig)              // 更新配置
			system.GET("/log/list", logHandler.GetLogList)                      // 获取系统日志列表
		}

		// 系统运维相关路由（需要认证和权限验证）
		monitorService := appService.NewSystemMonitorService()
		monitorHandler := handler.NewSystemMonitorHandler(monitorService)
		ops := api.Group("/ops")
		ops.Use(middleware.Auth())             // JWT认证
		ops.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			ops.GET("/monitor", monitorHandler.GetSystemMonitor) // 获取系统监控信息
		}

		// 文件管理相关路由（需要认证和权限验证）
		fileRepo := infraRepo.NewFileRepository()
		s3ConfigRepo := infraRepo.NewS3ConfigRepository()
		// 初始化本地文件存储
		localStorage := infraStorage.NewLocalStorage("./uploads", "/api/file/static")
		// 初始化S3存储（如果需要S3，需要先安装AWS SDK并配置）
		var s3Storage domainService.FileStorage = nil
		// 如果启用S3存储，可以在这里初始化：
		// s3Config, err := s3ConfigRepo.GetActiveConfig(context.Background())
		// if err == nil && s3Config != nil && s3Config.Status == "1" {
		// 	s3Storage, _ = infraStorage.NewS3Storage(infraStorage.S3Config{
		// 		Bucket:    s3Config.Bucket,
		// 		Region:    s3Config.Region,
		// 		Endpoint:  s3Config.Endpoint,
		// 		AccessKey: s3Config.AccessKeyID,
		// 		SecretKey: s3Config.SecretKey,
		// 		BaseURL:   s3Config.BaseURL,
		// 	})
		// }
		// 默认允许的文件扩展名（可以通过配置修改）
		allowedExts := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "svg",
			"doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "txt", "md", "csv",
			"mp4", "avi", "mov", "wmv", "flv", "mkv",
			"mp3", "wav", "flac", "aac",
			"zip", "rar", "7z", "tar", "gz"}
		// 最大文件大小：100MB
		maxFileSize := int64(100 * 1024 * 1024)
		fileService := appService.NewFileService(fileRepo, s3ConfigRepo, localStorage, s3Storage, maxFileSize, allowedExts)
		fileHandler := handler.NewFileHandler(fileService)
		s3ConfigService := appService.NewS3ConfigService(s3ConfigRepo)
		s3ConfigHandler := handler.NewS3ConfigHandler(s3ConfigService)
		file := api.Group("/file")
		file.Use(middleware.Auth())             // JWT认证
		file.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			file.POST("/upload", fileHandler.UploadFile)             // 上传文件
			file.GET("/list", fileHandler.ListFiles)                 // 获取文件列表
			file.GET("/:file_id", fileHandler.GetFile)               // 获取文件信息
			file.GET("/:file_id/download", fileHandler.DownloadFile) // 下载文件
			file.DELETE("/:file_id", fileHandler.DeleteFile)         // 删除文件
			// 存储配置管理
			file.GET("/storage/config", s3ConfigHandler.GetS3Config)                  // 获取存储配置
			file.POST("/storage/config", s3ConfigHandler.CreateOrUpdateS3Config)      // 创建或更新存储配置
			file.DELETE("/storage/config/:config_id", s3ConfigHandler.DeleteS3Config) // 删除存储配置
		}

		// 文件静态资源访问（需要认证）
		static := api.Group("/file/static")
		static.Use(middleware.Auth()) // JWT认证
		{
			// 静态文件服务由nginx或其他web服务器处理，这里只做路由占位
			// 实际文件访问应该通过GetFile接口获取URL
		}

		// 需要认证和权限验证的路由（管理功能）
		admin := api.Group("/user")
		admin.Use(middleware.Auth())             // JWT认证
		admin.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			admin.GET("/list", userHandler.ListUsers)                        // 获取用户列表
			admin.POST("", userHandler.CreateUser)                           // 创建用户
			admin.POST("/:userId/permissions", userHandler.GrantPermissions) // 给用户授权（必须放在/:userId之前）
			admin.GET("/:userId", userHandler.GetUserDetail)                 // 获取用户详情（必须放在最后）
			admin.PUT("/:userId", userHandler.UpdateUser)                    // 更新用户
			admin.DELETE("/:userId", userHandler.DeleteUser)                 // 删除用户
		}

		// 角色相关路由（需要认证和权限验证）
		role := api.Group("/role")
		role.Use(middleware.Auth())             // JWT认证
		role.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			role.GET("/list", roleHandler.ListRoles)        // 获取角色列表
			role.POST("", roleHandler.CreateRole)           // 创建角色
			role.GET("/:roleId", roleHandler.GetRoleDetail) // 获取角色详情（必须放在最后）
			role.PUT("/:roleId", roleHandler.UpdateRole)    // 更新角色
			role.DELETE("/:roleId", roleHandler.DeleteRole) // 删除角色
		}

		// 角色相关路由（只需要认证，不需要权限验证，用于下拉选择等）
		roleNoAuth := api.Group("/role")
		roleNoAuth.Use(middleware.Auth()) // 只需JWT认证
		{
			roleNoAuth.GET("/all", roleHandler.GetAllRoles) // 获取所有角色（不分页）
		}

		// 权限相关路由（需要认证和权限验证）
		permission := api.Group("/permission")
		permission.Use(middleware.Auth())             // JWT认证
		permission.Use(middleware.CasbinMiddleware()) // Casbin权限验证
		{
			permission.GET("/list", permissionHandler.ListPermissions)              // 获取权限列表
			permission.POST("", permissionHandler.CreatePermission)                 // 创建权限
			permission.GET("/:permissionId", permissionHandler.GetPermissionDetail) // 获取权限详情（必须放在最后）
			permission.PUT("/:permissionId", permissionHandler.UpdatePermission)    // 更新权限
			permission.DELETE("/:permissionId", permissionHandler.DeletePermission) // 删除权限
		}

		// 权限相关路由（只需要认证，不需要权限验证，用于下拉选择等）
		permissionNoAuth := api.Group("/permission")
		permissionNoAuth.Use(middleware.Auth()) // 只需JWT认证
		{
			permissionNoAuth.GET("/all", permissionHandler.GetAllPermissions) // 获取所有权限（不分页）
		}
	}

	// Swagger文档路由
	// 注意：Swagger UI的静态资源请求不需要认证，只有doc.json需要认证
	// 为了简化，Swagger UI完全开放访问（内部系统可以接受）
	// 如果需要保护，可以添加IP白名单或其他安全措施
	swagger := r.Group("/swagger")
	{
		url := ginSwagger.URL("/swagger/doc.json")
		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	// Swagger文档JSON需要认证（可选，如果需要保护API文档）
	// swaggerDoc := r.Group("/swagger")
	// swaggerDoc.Use(middleware.Auth())
	// {
	// 	swaggerDoc.GET("/doc.json", func(c *gin.Context) {
	// 		// 返回Swagger JSON文档
	// 	})
	// }

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

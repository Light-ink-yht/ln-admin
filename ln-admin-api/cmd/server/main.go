package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/redis"
	"github.com/Light-ink-yht/ln-admin/internal/interfaces/http/router"
	"github.com/Light-ink-yht/ln-admin/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title           LN Admin API
// @version         1.0
// @description     基于 Go + Gin + GORM 构建的现代化后台管理系统后端服务，采用 DDD（领域驱动设计）架构和 RBAC 权限控制
// @termsOfService  https://github.com/Light-ink-yht/ln-admin

// @contact.name   API Support
// @contact.url    https://github.com/Light-ink-yht/ln-admin
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT认证，格式：Bearer {token}。登录后获取 accessToken，在请求头中添加：Authorization: Bearer {accessToken}

func main() {
	bootstrap := NewBootstrap()

	// 加载配置
	if err := bootstrap.InitConfig(); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := bootstrap.InitLogger(); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 设置Gin模式
	if config.Cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	if err := bootstrap.InitDatabase(); err != nil {
		logger.Fatal("初始化数据库失败", logger.Fields(zap.Error(err))...)
	}

	// 初始化Redis
	if err := bootstrap.InitRedis(); err != nil {
		logger.Fatal("初始化Redis失败", logger.Fields(zap.Error(err))...)
	}
	defer redis.Close()

	// 初始化Casbin
	if err := bootstrap.InitCasbin(); err != nil {
		logger.Fatal("初始化Casbin失败", logger.Fields(zap.Error(err))...)
	}

	// 迁移数据库表结构
	if err := bootstrap.MigrateDatabase(); err != nil {
		logger.Fatal("数据库迁移失败", logger.Fields(zap.Error(err))...)
	}

	// 初始化仓库
	repos := bootstrap.InitRepositories()

	// 初始化服务
	services := bootstrap.InitServices(repos)

	// 初始化数据库日志服务
	bootstrap.InitDBLogger(repos)

	// 初始化默认数据（默认用户、角色、权限、短信模板）
	if err := bootstrap.InitDefaultData(repos, services); err != nil {
		logger.Warn("初始化默认数据失败", logger.Fields(zap.Error(err))...)
		// 不中断启动，只记录警告
	}

	// 设置路由
	r := router.SetupRouter(services.UserService, services.SMSAppService, services.PermissionService)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg.App.Port),
		Handler: r,
	}

	// 启动服务器（在goroutine中）
	go func() {
		logger.Info("服务器启动", logger.Fields(
			zap.String("address", srv.Addr),
			zap.String("mode", config.Cfg.App.Mode),
		)...)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", logger.Fields(zap.Error(err))...)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 设置5秒的超时时间用于关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器强制关闭", logger.Fields(zap.Error(err))...)
	}

	logger.Info("服务器已退出")
}

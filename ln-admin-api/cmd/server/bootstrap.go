package main

import (
	"fmt"
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/defaultdata"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/redis"
	infraRepo "github.com/Light-ink-yht/ln-admin/internal/infrastructure/repository"
	httpMiddleware "github.com/Light-ink-yht/ln-admin/internal/interfaces/http/middleware"
	"github.com/Light-ink-yht/ln-admin/pkg/config"
	logService "github.com/Light-ink-yht/ln-admin/pkg/service/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Bootstrap 应用启动初始化
type Bootstrap struct{}

// NewBootstrap 创建启动器
func NewBootstrap() *Bootstrap {
	return &Bootstrap{}
}

// InitConfig 初始化配置
func (b *Bootstrap) InitConfig() error {
	if err := config.Load("configs/config.yaml"); err != nil {
		return err
	}
	return nil
}

// InitLogger 初始化日志
func (b *Bootstrap) InitLogger() error {
	if err := logger.Init(); err != nil {
		return err
	}
	return nil
}

// InitDatabase 初始化数据库
func (b *Bootstrap) InitDatabase() error {
	if err := database.Init(); err != nil {
		return err
	}
	return nil
}

// InitRedis 初始化Redis
func (b *Bootstrap) InitRedis() error {
	if err := redis.Init(); err != nil {
		return err
	}
	return nil
}

// InitCasbin 初始化Casbin
func (b *Bootstrap) InitCasbin() error {
	if err := casbin.Init(); err != nil {
		return err
	}
	return nil
}

// MigrateDatabase 迁移数据库表结构
func (b *Bootstrap) MigrateDatabase() error {
	db := database.GetDB()

	// 定义表实体映射
	tables := []struct {
		table      interface{}
		getComment func() string
		tableName  string
	}{
		{&entity.AA01{}, (&entity.AA01{}).GetTableComment, "AA01"},
		{&entity.AA02{}, (&entity.AA02{}).GetTableComment, "AA02"},
		{&entity.AA03{}, (&entity.AA03{}).GetTableComment, "AA03"},
		{&entity.AA04{}, (&entity.AA04{}).GetTableComment, "AA04"},
		{&entity.AA05{}, (&entity.AA05{}).GetTableComment, "AA05"},
		{&entity.AA06{}, (&entity.AA06{}).GetTableComment, "AA06"},
		{&entity.AA07{}, (&entity.AA07{}).GetTableComment, "AA07"},
		{&entity.AA08{}, (&entity.AA08{}).GetTableComment, "AA08"},
		{&entity.AA09{}, (&entity.AA09{}).GetTableComment, "AA09"},
		{&entity.AA10{}, (&entity.AA10{}).GetTableComment, "AA10"},
	}

	// 执行 AutoMigrate
	var migrateTables []interface{}
	for _, t := range tables {
		migrateTables = append(migrateTables, t.table)
	}

	if err := db.AutoMigrate(migrateTables...); err != nil {
		return err
	}

	// 从实体结构体中读取注释并添加到数据库
	if err := b.addTableCommentsFromEntities(db, tables); err != nil {
		logger.Warn("添加表注释失败", logger.Fields(zap.Error(err))...)
		// 不中断启动，只记录警告
	}

	logger.Info("数据库表结构迁移成功")
	return nil
}

// addTableCommentsFromEntities 从实体结构体中读取注释并添加到数据库
func (b *Bootstrap) addTableCommentsFromEntities(db *gorm.DB, tables []struct {
	table      interface{}
	getComment func() string
	tableName  string
}) error {
	dbType := config.Cfg.Database.Type

	for _, t := range tables {
		comment := t.getComment()
		if comment == "" {
			continue
		}

		// 转义单引号（MySQL 和 Oracle 都需要）
		escapedComment := strings.ReplaceAll(comment, "'", "''")

		var sql string
		if dbType == "mysql" {
			sql = fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", t.tableName, escapedComment)
		} else if dbType == "oracle" {
			// Oracle 表注释使用 COMMENT ON TABLE
			sql = fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", t.tableName, escapedComment)
		} else {
			// 不支持的类型，跳过
			continue
		}

		if err := db.Exec(sql).Error; err != nil {
			logger.Warn(fmt.Sprintf("为表 %s 添加注释失败", t.tableName),
				logger.Fields(zap.Error(err))...)
			// 继续处理其他表，不中断
			continue
		}
		logger.Debug(fmt.Sprintf("为表 %s 添加注释成功", t.tableName))
	}

	return nil
}

// InitRepositories 初始化仓库
func (b *Bootstrap) InitRepositories() *Repositories {
	return &Repositories{
		UserRepository:         infraRepo.NewUserRepository(),
		SystemConfigRepository: infraRepo.NewSystemConfigRepository(),
		SystemLogRepository:    infraRepo.NewSystemLogRepository(),
		RoleRepository:         infraRepo.NewRoleRepository(),
		PermissionRepository:   infraRepo.NewPermissionRepository(),
		UserRoleRepository:     infraRepo.NewUserRoleRepository(),
		SMSTemplateRepository:  infraRepo.NewSMSTemplateRepository(),
	}
}

// InitServices 初始化服务
func (b *Bootstrap) InitServices(repos *Repositories) *Services {
	smsAppService := service.NewSMSAppService(repos.SystemConfigRepository, repos.SMSTemplateRepository)
	userService := service.NewUserAppService(repos.UserRepository, smsAppService)
	permissionService := service.NewPermissionService(
		repos.RoleRepository,
		repos.PermissionRepository,
		repos.UserRoleRepository,
	)

	return &Services{
		UserService:       userService,
		SMSAppService:     smsAppService,
		PermissionService: permissionService,
	}
}

// InitDBLogger 初始化数据库日志服务
func (b *Bootstrap) InitDBLogger(repos *Repositories) {
	dbLogger := logService.NewDBLogger(repos.SystemLogRepository)
	httpMiddleware.SetDBLogger(dbLogger)
	logger.Info("数据库日志服务初始化成功")
}

// InitDefaultData 初始化默认数据
func (b *Bootstrap) InitDefaultData(repos *Repositories, services *Services) error {
	if err := defaultdata.InitDefaultData(
		repos.UserRepository,
		repos.RoleRepository,
		repos.PermissionRepository,
		repos.UserRoleRepository,
		repos.SMSTemplateRepository,
		services.PermissionService,
	); err != nil {
		logger.Warn("初始化默认数据失败", logger.Fields(zap.Error(err))...)
		return err
	}
	return nil
}

// Repositories 仓库集合
type Repositories struct {
	UserRepository         repository.UserRepository
	SystemConfigRepository repository.SystemConfigRepository
	SystemLogRepository    repository.SystemLogRepository
	RoleRepository         repository.RoleRepository
	PermissionRepository   repository.PermissionRepository
	UserRoleRepository     repository.UserRoleRepository
	SMSTemplateRepository  repository.SMSTemplateRepository
}

// Services 服务集合
type Services struct {
	UserService       *service.UserAppService
	SMSAppService     *service.SMSAppService
	PermissionService *service.PermissionService
}

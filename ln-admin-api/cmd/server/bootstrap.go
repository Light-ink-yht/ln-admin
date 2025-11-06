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
		{&entity.AA10File{}, (&entity.AA10File{}).GetTableComment, "AA10"}, // 文件表
		{&entity.AA11{}, (&entity.AA11{}).GetTableComment, "AA11"},         // S3配置表
		{&entity.AA12{}, (&entity.AA12{}).GetTableComment, "AA12"},         // 角色菜单关联表
		{&entity.AA13{}, (&entity.AA13{}).GetTableComment, "AA13"},         // 角色权限授予关系表
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

	// 创建唯一联合索引（GORM AutoMigrate不会自动创建唯一联合索引）
	if err := b.createUniqueIndexes(db); err != nil {
		logger.Warn("创建唯一联合索引失败", logger.Fields(zap.Error(err))...)
		// 不中断启动，只记录警告
	}

	logger.Info("数据库表结构迁移成功")
	return nil
}

// createUniqueIndexes 创建唯一联合索引
func (b *Bootstrap) createUniqueIndexes(db *gorm.DB) error {
	dbType := config.Cfg.Database.Type

	// 定义需要创建的唯一联合索引
	indexes := []struct {
		tableName string
		indexName string
		columns   []string
	}{
		{"AA07", "uk_user_role", []string{"AAG001", "AAG002"}},                            // 用户角色关联表：防止重复分配
		{"AA12", "uk_role_menu", []string{"AAL001", "AAL002"}},                            // 角色菜单关联表：防止重复分配
		{"AA13", "uk_grantor_grantee_permission", []string{"AAM001", "AAM002", "AAM003"}}, // 角色权限授予关系表：防止重复授予
	}

	for _, idx := range indexes {
		// 先检查索引是否已存在
		var indexExists bool
		var checkSQL string

		if dbType == "mysql" {
			checkSQL = fmt.Sprintf("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = '%s' AND index_name = '%s'", idx.tableName, idx.indexName)
		} else if dbType == "oracle" {
			checkSQL = fmt.Sprintf("SELECT COUNT(*) FROM user_indexes WHERE index_name = '%s'", idx.indexName)
		} else {
			// 不支持的类型，跳过
			continue
		}

		// 检查索引是否存在
		var count int64
		if err := db.Raw(checkSQL).Scan(&count).Error; err == nil && count > 0 {
			indexExists = true
			logger.Debug(fmt.Sprintf("索引 %s 已存在，跳过创建", idx.indexName))
			continue
		}

		// 如果索引不存在，则创建
		if !indexExists {
			var sql string
			columns := strings.Join(idx.columns, ", ")
			if dbType == "mysql" {
				// MySQL: CREATE UNIQUE INDEX index_name ON table_name (col1, col2)
				sql = fmt.Sprintf("CREATE UNIQUE INDEX %s ON `%s` (%s)", idx.indexName, idx.tableName, columns)
			} else if dbType == "oracle" {
				// Oracle: CREATE UNIQUE INDEX index_name ON table_name (col1, col2)
				sql = fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s)", idx.indexName, idx.tableName, columns)
			}

			if err := db.Exec(sql).Error; err != nil {
				// 如果是索引已存在的错误，忽略（防止并发创建）
				errMsg := strings.ToLower(err.Error())
				if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "already exists") {
					logger.Debug(fmt.Sprintf("索引 %s 已存在，跳过创建", idx.indexName))
					continue
				}
				logger.Warn(fmt.Sprintf("创建索引 %s 失败", idx.indexName),
					logger.Fields(zap.Error(err))...)
				// 继续处理其他索引，不中断
				continue
			}
			logger.Info(fmt.Sprintf("创建唯一联合索引 %s 成功", idx.indexName))
		}
	}

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
		UserRepository:                infraRepo.NewUserRepository(),
		SystemConfigRepository:        infraRepo.NewSystemConfigRepository(),
		SystemLogRepository:           infraRepo.NewSystemLogRepository(),
		RoleRepository:                infraRepo.NewRoleRepository(),
		PermissionRepository:          infraRepo.NewPermissionRepository(),
		UserRoleRepository:            infraRepo.NewUserRoleRepository(),
		SMSTemplateRepository:         infraRepo.NewSMSTemplateRepository(),
		MenuRepository:                infraRepo.NewMenuRepository(),
		RolePermissionGrantRepository: infraRepo.NewRolePermissionGrantRepository(),
	}
}

// InitServices 初始化服务
func (b *Bootstrap) InitServices(repos *Repositories) *Services {
	smsAppService := service.NewSMSAppService(repos.SystemConfigRepository, repos.SMSTemplateRepository)
	permissionService := service.NewPermissionService(
		repos.RoleRepository,
		repos.PermissionRepository,
		repos.UserRoleRepository,
		repos.RolePermissionGrantRepository,
	)
	userService := service.NewUserAppService(repos.UserRepository, repos.RoleRepository, smsAppService, permissionService)

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
		repos.MenuRepository,
		services.PermissionService,
	); err != nil {
		logger.Warn("初始化默认数据失败", logger.Fields(zap.Error(err))...)
		return err
	}
	return nil
}

// Repositories 仓库集合
type Repositories struct {
	UserRepository                repository.UserRepository
	SystemConfigRepository        repository.SystemConfigRepository
	SystemLogRepository           repository.SystemLogRepository
	RoleRepository                repository.RoleRepository
	PermissionRepository          repository.PermissionRepository
	UserRoleRepository            repository.UserRoleRepository
	SMSTemplateRepository         repository.SMSTemplateRepository
	MenuRepository                repository.MenuRepository
	RolePermissionGrantRepository repository.RolePermissionGrantRepository
}

// Services 服务集合
type Services struct {
	UserService       *service.UserAppService
	SMSAppService     *service.SMSAppService
	PermissionService *service.PermissionService
}

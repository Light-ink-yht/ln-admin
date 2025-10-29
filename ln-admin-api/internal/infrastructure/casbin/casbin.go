package casbin

import (
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

// Init 初始化Casbin
func Init() error {
	// 使用GORM适配器，将权限规则存储在数据库中
	adapter, err := gormadapter.NewAdapterByDB(database.GetDB())
	if err != nil {
		return fmt.Errorf("创建Casbin适配器失败: %w", err)
	}

	// 加载RBAC模型配置
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return fmt.Errorf("创建Casbin Enforcer失败: %w", err)
	}

	// 启用自动保存
	enforcer.EnableAutoSave(true)

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %w", err)
	}

	Enforcer = enforcer

	logger.Info("Casbin权限控制初始化成功")
	return nil
}

// GetEnforcer 获取Casbin Enforcer实例
func GetEnforcer() *casbin.Enforcer {
	return Enforcer
}

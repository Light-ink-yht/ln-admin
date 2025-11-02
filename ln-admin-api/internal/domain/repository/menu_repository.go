package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// MenuRepository 菜单仓库接口
type MenuRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *entity.Menu) error

	// FindByID 根据ID查找菜单
	FindByID(ctx context.Context, menuID string) (*entity.Menu, error)

	// FindByKey 根据菜单标识查找菜单
	FindByKey(ctx context.Context, menuKey string) (*entity.Menu, error)

	// Update 更新菜单
	Update(ctx context.Context, menu *entity.Menu) error

	// Delete 删除菜单
	Delete(ctx context.Context, menuID string) error

	// FindAll 查找所有菜单
	FindAll(ctx context.Context) ([]*entity.Menu, error)

	// FindByType 根据菜单类型查找菜单
	FindByType(ctx context.Context, menuType string) ([]*entity.Menu, error)

	// FindByParentID 根据父菜单ID查找子菜单
	FindByParentID(ctx context.Context, parentID string) ([]*entity.Menu, error)

	// FindTreeByType 根据菜单类型查找菜单树（包含子菜单）
	FindTreeByType(ctx context.Context, menuType string) ([]*entity.Menu, error)
}

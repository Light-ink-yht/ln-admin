package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"
	"gorm.io/gorm"
)

var _ repository.MenuRepository = (*menuRepositoryImpl)(nil)

type menuRepositoryImpl struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓库实例
func NewMenuRepository() repository.MenuRepository {
	return &menuRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建菜单
func (r *menuRepositoryImpl) Create(ctx context.Context, menu *entity.Menu) error {
	aa09 := menu.ToAA09()
	if err := r.db.WithContext(ctx).Create(aa09).Error; err != nil {
		return fmt.Errorf("创建菜单失败: %w", err)
	}
	menu.FromAA09(aa09)
	return nil
}

// FindByID 根据ID查找菜单
func (r *menuRepositoryImpl) FindByID(ctx context.Context, menuID string) (*entity.Menu, error) {
	var aa09 entity.AA09
	if err := r.db.WithContext(ctx).Where("AAI001 = ?", menuID).First(&aa09).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	menu := &entity.Menu{}
	menu.FromAA09(&aa09)
	return menu, nil
}

// FindByKey 根据菜单标识查找菜单
func (r *menuRepositoryImpl) FindByKey(ctx context.Context, menuKey string) (*entity.Menu, error) {
	var aa09 entity.AA09
	if err := r.db.WithContext(ctx).Where("AAI002 = ?", menuKey).First(&aa09).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	menu := &entity.Menu{}
	menu.FromAA09(&aa09)
	return menu, nil
}

// Update 更新菜单
func (r *menuRepositoryImpl) Update(ctx context.Context, menu *entity.Menu) error {
	aa09 := menu.ToAA09()
	if err := r.db.WithContext(ctx).Where("AAI001 = ?", aa09.AAI001).Updates(aa09).Error; err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}
	return nil
}

// Delete 删除菜单
func (r *menuRepositoryImpl) Delete(ctx context.Context, menuID string) error {
	if err := r.db.WithContext(ctx).Where("AAI001 = ?", menuID).Delete(&entity.AA09{}).Error; err != nil {
		return fmt.Errorf("删除菜单失败: %w", err)
	}
	return nil
}

// FindAll 查找所有菜单
func (r *menuRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Menu, error) {
	var aa09List []entity.AA09
	if err := r.db.WithContext(ctx).Order("AAI007 ASC, AAI001 ASC").Find(&aa09List).Error; err != nil {
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}

	menus := make([]*entity.Menu, len(aa09List))
	for i, aa09 := range aa09List {
		menu := &entity.Menu{}
		menu.FromAA09(&aa09)
		menus[i] = menu
	}

	return menus, nil
}

// FindByType 根据菜单类型查找菜单
func (r *menuRepositoryImpl) FindByType(ctx context.Context, menuType string) ([]*entity.Menu, error) {
	var aa09List []entity.AA09
	if err := r.db.WithContext(ctx).Where("AAI009 = ?", menuType).Order("AAI007 ASC, AAI001 ASC").Find(&aa09List).Error; err != nil {
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}

	menus := make([]*entity.Menu, len(aa09List))
	for i, aa09 := range aa09List {
		menu := &entity.Menu{}
		menu.FromAA09(&aa09)
		menus[i] = menu
	}

	return menus, nil
}

// FindByParentID 根据父菜单ID查找子菜单
func (r *menuRepositoryImpl) FindByParentID(ctx context.Context, parentID string) ([]*entity.Menu, error) {
	var aa09List []entity.AA09
	if err := r.db.WithContext(ctx).Where("AAI006 = ?", parentID).Order("AAI007 ASC, AAI001 ASC").Find(&aa09List).Error; err != nil {
		return nil, fmt.Errorf("查询子菜单列表失败: %w", err)
	}

	menus := make([]*entity.Menu, len(aa09List))
	for i, aa09 := range aa09List {
		menu := &entity.Menu{}
		menu.FromAA09(&aa09)
		menus[i] = menu
	}

	return menus, nil
}

// FindTreeByType 根据菜单类型查找菜单树（包含子菜单）
func (r *menuRepositoryImpl) FindTreeByType(ctx context.Context, menuType string) ([]*entity.Menu, error) {
	// 查找所有该类型的菜单
	var aa09List []entity.AA09
	if err := r.db.WithContext(ctx).Where("AAI009 = ? AND AAI010 = '1'", menuType).Order("AAI007 ASC, AAI001 ASC").Find(&aa09List).Error; err != nil {
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}

	// 转换为实体
	menus := make([]*entity.Menu, len(aa09List))
	menuMap := make(map[string]*entity.Menu)

	for i, aa09 := range aa09List {
		menu := &entity.Menu{}
		menu.FromAA09(&aa09)
		menus[i] = menu
		menuMap[menu.MenuID] = menu
	}

	// 构建菜单树
	var rootMenus []*entity.Menu
	for _, menu := range menus {
		if menu.ParentID == "" || menu.ParentID == "0" {
			// 根菜单
			rootMenus = append(rootMenus, menu)
		} else {
			// 子菜单，添加到父菜单的Children中
			if parent, exists := menuMap[menu.ParentID]; exists {
				if parent.Children == nil {
					parent.Children = make([]*entity.Menu, 0)
				}
				parent.Children = append(parent.Children, menu)
			}
		}
	}

	return rootMenus, nil
}

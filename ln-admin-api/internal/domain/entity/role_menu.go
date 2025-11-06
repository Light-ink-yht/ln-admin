package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA12 表注释
const TableCommentAA12 = "角色菜单关联表：存储角色与菜单之间的多对多关联关系，实现基于角色的菜单访问控制。通过 AAL001（角色ID）和 AAL002（菜单ID）建立角色和菜单的关联，支持一个角色拥有多个菜单。超级管理员可以为不同角色分配不同的菜单权限，控制不同角色可以访问的页面。"

// AA12 角色菜单关联表
type AA12 struct {
	gorm.Model
	AAL001 string `gorm:"column:AAL001;type:varchar(50);index;index:idx_role_menu,priority:1;comment:角色ID"` // 角色ID
	AAL002 string `gorm:"column:AAL002;type:varchar(50);index;index:idx_role_menu,priority:2;comment:菜单ID"` // 菜单ID
	// 注意：唯一联合索引需要在数据库迁移时手动创建，或在GORM中通过Migrate方法创建
}

// TableName 指定表名
func (AA12) TableName() string {
	return "AA12"
}

// GetTableComment 获取表注释
func (AA12) GetTableComment() string {
	return TableCommentAA12
}

// RoleMenu 角色菜单关联领域实体
type RoleMenu struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	RoleID    string     `json:"role_id"` // AAL001
	MenuID    string     `json:"menu_id"` // AAL002
}

// ToAA12 转换为数据库实体
func (rm *RoleMenu) ToAA12() *AA12 {
	model := gorm.Model{
		ID:        rm.ID,
		CreatedAt: rm.CreatedAt,
		UpdatedAt: rm.UpdatedAt,
	}
	if rm.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *rm.DeletedAt}
	}
	return &AA12{
		Model:  model,
		AAL001: rm.RoleID,
		AAL002: rm.MenuID,
	}
}

// FromAA12 从数据库实体转换
func (rm *RoleMenu) FromAA12(aa12 *AA12) {
	rm.ID = aa12.ID
	rm.CreatedAt = aa12.CreatedAt
	rm.UpdatedAt = aa12.UpdatedAt
	if aa12.DeletedAt.Valid {
		rm.DeletedAt = &aa12.DeletedAt.Time
	}
	rm.RoleID = aa12.AAL001
	rm.MenuID = aa12.AAL002
}

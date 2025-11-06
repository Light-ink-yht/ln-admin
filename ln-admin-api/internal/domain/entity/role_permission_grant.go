package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA13 表注释
const TableCommentAA13 = "角色权限授予关系表：存储角色权限的授予关系，实现层级授权机制。通过 AAM001（授予者角色ID）、AAM002（被授予者角色ID）和 AAM003（权限ID）建立权限授予关系，记录谁授予给谁的权限。超级管理员拥有所有权限，可以为其他角色授予权限；管理员只能看到被授予的权限，并可以将这些权限授予给下一层角色。支持权限的层级传递，确保权限管理的安全性和可追溯性。"

// AA13 角色权限授予关系表
type AA13 struct {
	gorm.Model
	AAM001 string `gorm:"column:AAM001;type:varchar(50);index;index:idx_grantor_role_permission,priority:1;comment:授予者角色ID"`                                               // 授予者角色ID
	AAM002 string `gorm:"column:AAM002;type:varchar(50);index;index:idx_grantee_role_permission,priority:1;index:idx_grantor_role_permission,priority:2;comment:被授予者角色ID"` // 被授予者角色ID
	AAM003 string `gorm:"column:AAM003;type:varchar(50);index;index:idx_grantee_role_permission,priority:2;index:idx_grantor_role_permission,priority:3;comment:权限ID"`     // 权限ID
	// 注意：唯一联合索引需要在数据库迁移时手动创建，或在GORM中通过Migrate方法创建
}

// TableName 指定表名
func (AA13) TableName() string {
	return "AA13"
}

// GetTableComment 获取表注释
func (AA13) GetTableComment() string {
	return TableCommentAA13
}

// RolePermissionGrant 角色权限授予关系领域实体
type RolePermissionGrant struct {
	ID            uint       `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	GrantorRoleID string     `json:"grantor_role_id"` // AAM001 授予者角色ID
	GranteeRoleID string     `json:"grantee_role_id"` // AAM002 被授予者角色ID
	PermissionID  string     `json:"permission_id"`   // AAM003 权限ID
}

// ToAA13 转换为数据库实体
func (rpg *RolePermissionGrant) ToAA13() *AA13 {
	model := gorm.Model{
		ID:        rpg.ID,
		CreatedAt: rpg.CreatedAt,
		UpdatedAt: rpg.UpdatedAt,
	}
	if rpg.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *rpg.DeletedAt}
	}
	return &AA13{
		Model:  model,
		AAM001: rpg.GrantorRoleID,
		AAM002: rpg.GranteeRoleID,
		AAM003: rpg.PermissionID,
	}
}

// FromAA13 从数据库实体转换
func (rpg *RolePermissionGrant) FromAA13(aa13 *AA13) {
	rpg.ID = aa13.ID
	rpg.CreatedAt = aa13.CreatedAt
	rpg.UpdatedAt = aa13.UpdatedAt
	if aa13.DeletedAt.Valid {
		rpg.DeletedAt = &aa13.DeletedAt.Time
	}
	rpg.GrantorRoleID = aa13.AAM001
	rpg.GranteeRoleID = aa13.AAM002
	rpg.PermissionID = aa13.AAM003
}

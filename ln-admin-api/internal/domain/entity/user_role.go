package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA07 表注释
const TableCommentAA07 = "用户角色关联表：存储用户与角色之间的多对多关联关系，实现 RBAC 权限模型。通过 AAG001（用户ID）和 AAG002（角色ID）建立用户和角色的关联，支持一个用户拥有多个角色。在 Casbin 策略表中，会同步创建相应的角色分配策略（g策略），用于权限验证。支持按用户ID和角色ID建立索引，便于快速查询用户的角色列表或角色的用户列表。"

// AA07 用户角色关联表
type AA07 struct {
	gorm.Model
	AAG001 string `gorm:"column:AAG001;type:varchar(50);index;comment:用户ID"` // 用户ID
	AAG002 string `gorm:"column:AAG002;type:varchar(50);index;comment:角色ID"` // 角色ID
}

// TableName 指定表名
func (AA07) TableName() string {
	return "AA07"
}

// GetTableComment 获取表注释
func (AA07) GetTableComment() string {
	return TableCommentAA07
}

// UserRole 用户角色关联领域实体
type UserRole struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	UserID    string     `json:"user_id"` // AAG001
	RoleID    string     `json:"role_id"` // AAG002
}

// ToAA07 转换为数据库实体
func (ur *UserRole) ToAA07() *AA07 {
	model := gorm.Model{
		ID:        ur.ID,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
	}
	if ur.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *ur.DeletedAt}
	}
	return &AA07{
		Model:  model,
		AAG001: ur.UserID,
		AAG002: ur.RoleID,
	}
}

// FromAA07 从数据库实体转换
func (ur *UserRole) FromAA07(aa07 *AA07) {
	ur.ID = aa07.ID
	ur.CreatedAt = aa07.CreatedAt
	ur.UpdatedAt = aa07.UpdatedAt
	ur.DeletedAt = &aa07.DeletedAt.Time
	ur.UserID = aa07.AAG001
	ur.RoleID = aa07.AAG002
}

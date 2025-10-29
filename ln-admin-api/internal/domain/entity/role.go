package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA05 表注释
const TableCommentAA05 = "角色表：定义系统中的角色信息，如超级管理员、管理员、普通用户等。每个角色有唯一的角色标识（AAE002），通过角色标识进行权限判断。角色包含角色名称、描述、状态等字段。支持启用/禁用状态控制，主键为 AAE001（角色ID），AAE002（角色标识）为唯一索引。角色与权限通过 Casbin 策略表关联，实现 RBAC 权限模型。"

// AA05 角色表
type AA05 struct {
	gorm.Model
	AAE001 string `gorm:"column:AAE001;type:varchar(50);primaryKey;comment:角色ID"`          // 角色ID, 主键
	AAE002 string `gorm:"column:AAE002;type:varchar(100);uniqueIndex;comment:角色标识"`        // 角色标识（唯一）
	AAE003 string `gorm:"column:AAE003;type:varchar(100);comment:角色名称"`                    // 角色名称
	AAE004 string `gorm:"column:AAE004;type:varchar(500);comment:角色描述"`                    // 角色描述
	AAE005 string `gorm:"column:AAE005;type:varchar(10);default:'1';comment:状态 1 启用，2 禁用"` // 状态
	AAE006 string `gorm:"column:AAE006;type:varchar(50);comment:创建人"`                      // 创建人
	AAE007 string `gorm:"column:AAE007;type:varchar(50);comment:修改人"`                      // 修改人
}

// TableName 指定表名
func (AA05) TableName() string {
	return "AA05"
}

// GetTableComment 获取表注释
func (AA05) GetTableComment() string {
	return TableCommentAA05
}

// Role 角色领域实体
type Role struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	RoleID      string     `json:"role_id"`     // AAE001
	RoleKey     string     `json:"role_key"`    // AAE002
	RoleName    string     `json:"role_name"`   // AAE003
	Description string     `json:"description"` // AAE004
	Status      string     `json:"status"`      // AAE005
	CreatorID   string     `json:"creator_id"`  // AAE006
	ModifierID  string     `json:"modifier_id"` // AAE007
}

// ToAA05 转换为数据库实体
func (r *Role) ToAA05() *AA05 {
	model := gorm.Model{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	// 处理 DeletedAt（可能为 nil）
	if r.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *r.DeletedAt}
	}

	return &AA05{
		Model:  model,
		AAE001: r.RoleID,
		AAE002: r.RoleKey,
		AAE003: r.RoleName,
		AAE004: r.Description,
		AAE005: r.Status,
		AAE006: r.CreatorID,
		AAE007: r.ModifierID,
	}
}

// FromAA05 从数据库实体转换
func (r *Role) FromAA05(aa05 *AA05) {
	r.ID = aa05.ID
	r.CreatedAt = aa05.CreatedAt
	r.UpdatedAt = aa05.UpdatedAt
	r.DeletedAt = &aa05.DeletedAt.Time
	r.RoleID = aa05.AAE001
	r.RoleKey = aa05.AAE002
	r.RoleName = aa05.AAE003
	r.Description = aa05.AAE004
	r.Status = aa05.AAE005
	r.CreatorID = aa05.AAE006
	r.ModifierID = aa05.AAE007
}

// IsDisabled 判断角色是否禁用
func (r *Role) IsDisabled() bool {
	return r.Status == "2"
}

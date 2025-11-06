package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA06 表注释
const TableCommentAA06 = "权限表：定义系统的权限资源，包括资源路径、请求方法等权限信息。权限通过权限标识（AAF002）唯一标识，包含资源路径（API路径，如 /api/user/list）、请求方法（GET、POST、PUT、DELETE等）、权限名称、描述等字段。主键为 AAF001（权限ID），AAF002（权限标识）为唯一索引。权限与角色通过 Casbin 策略表关联，实现基于角色的访问控制（RBAC）。"

// AA06 权限表
type AA06 struct {
	gorm.Model
	AAF001 string `gorm:"column:AAF001;type:varchar(50);primaryKey;comment:权限ID"`                            // 权限ID, 主键
	AAF002 string `gorm:"column:AAF002;type:varchar(100);uniqueIndex;comment:权限标识"`                          // 权限标识（唯一）
	AAF003 string `gorm:"column:AAF003;type:varchar(100);comment:权限名称"`                                      // 权限名称
	AAF004 string `gorm:"column:AAF004;type:varchar(200);index:idx_resource_method,priority:1;comment:资源路径"` // 资源路径（API路径）
	AAF005 string `gorm:"column:AAF005;type:varchar(10);index:idx_resource_method,priority:2;comment:请求方法"`  // 请求方法（GET, POST, PUT, DELETE等）
	AAF006 string `gorm:"column:AAF006;type:varchar(500);comment:权限描述"`                                      // 权限描述
	AAF007 string `gorm:"column:AAF007;type:varchar(10);default:'1';index;comment:状态 1 启用，2 禁用"`             // 状态
	AAF008 string `gorm:"column:AAF008;type:varchar(50);comment:创建人"`                                        // 创建人
	AAF009 string `gorm:"column:AAF009;type:varchar(50);comment:修改人"`                                        // 修改人
	AAF010 string `gorm:"column:AAF010;type:varchar(50);index;comment:分类"`                                   // 分类（用户管理、角色管理等）
}

// TableName 指定表名
func (AA06) TableName() string {
	return "AA06"
}

// GetTableComment 获取表注释
func (AA06) GetTableComment() string {
	return TableCommentAA06
}

// Permission 权限领域实体
type Permission struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	PermissionID   string     `json:"permission_id"`   // AAF001
	PermissionKey  string     `json:"permission_key"`  // AAF002
	PermissionName string     `json:"permission_name"` // AAF003
	ResourcePath   string     `json:"resource_path"`   // AAF004
	Method         string     `json:"method"`          // AAF005
	Description    string     `json:"description"`     // AAF006
	Status         string     `json:"status"`          // AAF007
	CreatorID      string     `json:"creator_id"`      // AAF008
	ModifierID     string     `json:"modifier_id"`     // AAF009
	Category       string     `json:"category"`        // AAF010
}

// ToAA06 转换为数据库实体
func (p *Permission) ToAA06() *AA06 {
	model := gorm.Model{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if p.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt}
	}
	return &AA06{
		Model:  model,
		AAF001: p.PermissionID,
		AAF002: p.PermissionKey,
		AAF003: p.PermissionName,
		AAF004: p.ResourcePath,
		AAF005: p.Method,
		AAF006: p.Description,
		AAF007: p.Status,
		AAF008: p.CreatorID,
		AAF009: p.ModifierID,
		AAF010: p.Category,
	}
}

// FromAA06 从数据库实体转换
func (p *Permission) FromAA06(aa06 *AA06) {
	p.ID = aa06.ID
	p.CreatedAt = aa06.CreatedAt
	p.UpdatedAt = aa06.UpdatedAt
	p.DeletedAt = &aa06.DeletedAt.Time
	p.PermissionID = aa06.AAF001
	p.PermissionKey = aa06.AAF002
	p.PermissionName = aa06.AAF003
	p.ResourcePath = aa06.AAF004
	p.Method = aa06.AAF005
	p.Description = aa06.AAF006
	p.Status = aa06.AAF007
	p.CreatorID = aa06.AAF008
	p.ModifierID = aa06.AAF009
	p.Category = aa06.AAF010
}

// IsDisabled 判断权限是否禁用
func (p *Permission) IsDisabled() bool {
	return p.Status == "2"
}

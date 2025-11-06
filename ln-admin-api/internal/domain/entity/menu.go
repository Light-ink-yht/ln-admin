package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA09 表注释
const TableCommentAA09 = "菜单表：定义系统中的菜单项，包括侧边栏菜单和用户下拉菜单。菜单支持多级嵌套，通过父菜单ID（AAI006）实现层级关系。菜单包含菜单标识、标题、路径、图标、排序、权限标识等字段。主键为 AAI001（菜单ID），AAI002（菜单标识）为唯一索引。菜单可以通过权限标识关联权限，控制菜单的显示和访问。"

// AA09 菜单表
type AA09 struct {
	gorm.Model
	AAI001 string `gorm:"column:AAI001;type:varchar(50);primaryKey;comment:菜单ID"`                                                        // 菜单ID, 主键
	AAI002 string `gorm:"column:AAI002;type:varchar(100);uniqueIndex;comment:菜单标识"`                                                      // 菜单标识（唯一）
	AAI003 string `gorm:"column:AAI003;type:varchar(100);comment:菜单标题"`                                                                  // 菜单标题
	AAI004 string `gorm:"column:AAI004;type:varchar(200);comment:路由路径"`                                                                  // 路由路径
	AAI005 string `gorm:"column:AAI005;type:varchar(100);comment:图标名称"`                                                                  // 图标名称
	AAI006 string `gorm:"column:AAI006;type:varchar(50);index;index:idx_parent_type,priority:1;comment:父菜单ID"`                           // 父菜单ID（用于构建菜单树）
	AAI007 int    `gorm:"column:AAI007;type:int;default:0;comment:排序号"`                                                                  // 排序号
	AAI008 string `gorm:"column:AAI008;type:varchar(200);comment:权限标识"`                                                                  // 权限标识（关联权限表）
	AAI009 string `gorm:"column:AAI009;type:varchar(10);default:'1';index;index:idx_parent_type,priority:2;comment:菜单类型 1 侧边栏菜单 2 用户菜单"` // 菜单类型
	AAI010 string `gorm:"column:AAI010;type:varchar(10);default:'1';index;comment:状态 1 启用，2 禁用"`                                         // 状态
	AAI011 string `gorm:"column:AAI011;type:varchar(500);comment:菜单描述"`                                                                  // 菜单描述
	AAI012 string `gorm:"column:AAI012;type:varchar(50);comment:创建人"`                                                                    // 创建人
	AAI013 string `gorm:"column:AAI013;type:varchar(50);comment:修改人"`                                                                    // 修改人
}

// TableName 指定表名
func (AA09) TableName() string {
	return "AA09"
}

// GetTableComment 获取表注释
func (AA09) GetTableComment() string {
	return TableCommentAA09
}

// Menu 菜单领域实体
type Menu struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	MenuID      string     `json:"menu_id"`                     // AAI001
	MenuKey     string     `json:"menu_key"`                    // AAI002
	Title       string     `json:"title"`                       // AAI003
	Path        string     `json:"path"`                        // AAI004
	Icon        string     `json:"icon"`                        // AAI005
	ParentID    string     `json:"parent_id"`                   // AAI006
	Sort        int        `json:"sort"`                        // AAI007
	Permission  string     `json:"permission"`                  // AAI008
	MenuType    string     `json:"menu_type"`                   // AAI009: 1 侧边栏菜单 2 用户菜单
	Status      string     `json:"status"`                      // AAI010: 1 启用 2 禁用
	Description string     `json:"description"`                 // AAI011
	CreatorID   string     `json:"creator_id"`                  // AAI012
	ModifierID  string     `json:"modifier_id"`                 // AAI013
	Children    []*Menu    `json:"children,omitempty" gorm:"-"` // 子菜单（不存储到数据库）
}

// ToAA09 转换为数据库实体
func (m *Menu) ToAA09() *AA09 {
	model := gorm.Model{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *m.DeletedAt}
	}
	return &AA09{
		Model:  model,
		AAI001: m.MenuID,
		AAI002: m.MenuKey,
		AAI003: m.Title,
		AAI004: m.Path,
		AAI005: m.Icon,
		AAI006: m.ParentID,
		AAI007: m.Sort,
		AAI008: m.Permission,
		AAI009: m.MenuType,
		AAI010: m.Status,
		AAI011: m.Description,
		AAI012: m.CreatorID,
		AAI013: m.ModifierID,
	}
}

// FromAA09 从数据库实体转换
func (m *Menu) FromAA09(aa09 *AA09) {
	m.ID = aa09.ID
	m.CreatedAt = aa09.CreatedAt
	m.UpdatedAt = aa09.UpdatedAt
	if aa09.DeletedAt.Valid {
		m.DeletedAt = &aa09.DeletedAt.Time
	}
	m.MenuID = aa09.AAI001
	m.MenuKey = aa09.AAI002
	m.Title = aa09.AAI003
	m.Path = aa09.AAI004
	m.Icon = aa09.AAI005
	m.ParentID = aa09.AAI006
	m.Sort = aa09.AAI007
	m.Permission = aa09.AAI008
	m.MenuType = aa09.AAI009
	m.Status = aa09.AAI010
	m.Description = aa09.AAI011
	m.CreatorID = aa09.AAI012
	m.ModifierID = aa09.AAI013
}

// IsDisabled 判断菜单是否禁用
func (m *Menu) IsDisabled() bool {
	return m.Status == "2"
}

// IsSidebarMenu 判断是否为侧边栏菜单
func (m *Menu) IsSidebarMenu() bool {
	return m.MenuType == "1"
}

// IsUserMenu 判断是否为用户菜单
func (m *Menu) IsUserMenu() bool {
	return m.MenuType == "2"
}

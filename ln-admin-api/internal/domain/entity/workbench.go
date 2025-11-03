package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA10 表注释
const TableCommentAA10 = "工作台配置表：存储不同角色的工作台配置信息，包括统计卡片、快捷操作、公告等配置。每个角色对应一个工作台配置，配置内容以JSON格式存储，支持灵活配置工作台的显示内容和布局。主键为 AAJ001（配置ID），AAJ002（角色标识）为唯一索引，确保每个角色只有一个工作台配置。"

// AA10 工作台配置表
type AA10 struct {
	gorm.Model
	AAJ001 string `gorm:"column:AAJ001;type:varchar(50);primaryKey;comment:配置ID"`          // 配置ID, 主键
	AAJ002 string `gorm:"column:AAJ002;type:varchar(100);uniqueIndex;comment:角色标识"`        // 角色标识（唯一，关联AA05.AAE002）
	AAJ003 string `gorm:"column:AAJ003;type:varchar(100);comment:配置名称"`                    // 配置名称
	AAJ004 string `gorm:"column:AAJ004;type:text;comment:配置内容(JSON)"`                      // 配置内容（JSON格式）
	AAJ005 string `gorm:"column:AAJ005;type:varchar(10);default:'1';comment:状态 1 启用，2 禁用"` // 状态
	AAJ006 string `gorm:"column:AAJ006;type:varchar(500);comment:配置描述"`                    // 配置描述
	AAJ007 string `gorm:"column:AAJ007;type:varchar(50);comment:创建人"`                      // 创建人
	AAJ008 string `gorm:"column:AAJ008;type:varchar(50);comment:修改人"`                      // 修改人
}

// TableName 指定表名
func (AA10) TableName() string {
	return "AA10"
}

// GetTableComment 获取表注释
func (AA10) GetTableComment() string {
	return TableCommentAA10
}

// WorkbenchConfig 工作台配置领域实体
type WorkbenchConfig struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	ConfigID    string     `json:"config_id"`   // AAJ001
	RoleKey     string     `json:"role_key"`    // AAJ002
	ConfigName  string     `json:"config_name"` // AAJ003
	ConfigData  string     `json:"config_data"` // AAJ004 (JSON字符串)
	Status      string     `json:"status"`      // AAJ005
	Description string     `json:"description"` // AAJ006
	CreatorID   string     `json:"creator_id"`  // AAJ007
	ModifierID  string     `json:"modifier_id"` // AAJ008
}

// ToAA10 转换为数据库实体
func (wc *WorkbenchConfig) ToAA10() *AA10 {
	model := gorm.Model{
		ID:        wc.ID,
		CreatedAt: wc.CreatedAt,
		UpdatedAt: wc.UpdatedAt,
	}

	// 处理 DeletedAt（可能为 nil）
	if wc.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *wc.DeletedAt}
	}

	return &AA10{
		Model:  model,
		AAJ001: wc.ConfigID,
		AAJ002: wc.RoleKey,
		AAJ003: wc.ConfigName,
		AAJ004: wc.ConfigData,
		AAJ005: wc.Status,
		AAJ006: wc.Description,
		AAJ007: wc.CreatorID,
		AAJ008: wc.ModifierID,
	}
}

// FromAA10 从数据库实体转换
func (wc *WorkbenchConfig) FromAA10(aa10 *AA10) {
	wc.ID = aa10.ID
	wc.CreatedAt = aa10.CreatedAt
	wc.UpdatedAt = aa10.UpdatedAt
	wc.DeletedAt = &aa10.DeletedAt.Time
	wc.ConfigID = aa10.AAJ001
	wc.RoleKey = aa10.AAJ002
	wc.ConfigName = aa10.AAJ003
	wc.ConfigData = aa10.AAJ004
	wc.Status = aa10.AAJ005
	wc.Description = aa10.AAJ006
	wc.CreatorID = aa10.AAJ007
	wc.ModifierID = aa10.AAJ008
}

// IsDisabled 判断配置是否禁用
func (wc *WorkbenchConfig) IsDisabled() bool {
	return wc.Status == "2"
}

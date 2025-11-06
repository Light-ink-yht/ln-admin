package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA02 表注释
const TableCommentAA02 = "系统配置表：存储系统运行所需的各种配置项信息，支持动态配置管理。配置项按分组管理（如 sms、jwt、system 等），每个配置有唯一的配置键（AAB001），配置值（AAB002）支持 JSON 格式存储复杂配置。支持启用/禁用状态控制，便于配置的启用和停用。主键为 AAB001（配置键）。"

// AA02 系统配置表
type AA02 struct {
	gorm.Model
	AAB001 string     `gorm:"column:AAB001;type:varchar(50);primaryKey;comment:配置键"`                                                   // 配置键, 主键
	AAB002 string     `gorm:"column:AAB002;type:varchar(500);comment:配置值"`                                                             // 配置值
	AAB003 string     `gorm:"column:AAB003;type:varchar(100);comment:配置名称"`                                                            // 配置名称
	AAB004 string     `gorm:"column:AAB004;type:varchar(50);index;index:idx_group_status,priority:1;comment:配置分组"`                     // 配置分组（sms, jwt, system等）
	AAB005 string     `gorm:"column:AAB005;type:varchar(500);comment:配置描述"`                                                            // 配置描述
	AAB006 string     `gorm:"column:AAB006;type:varchar(10);default:'1';index;index:idx_group_status,priority:2;comment:状态 1 启用，2 禁用"` // 状态
	AAB007 *time.Time `gorm:"column:AAB007;type:datetime;comment:最后修改时间"`                                                              // 最后修改时间
	AAB008 string     `gorm:"column:AAB008;type:varchar(50);comment:创建人"`                                                              // 创建人
	AAB009 string     `gorm:"column:AAB009;type:varchar(50);comment:修改人"`                                                              // 修改人
}

// TableName 指定表名
func (AA02) TableName() string {
	return "AA02"
}

// GetTableComment 获取表注释
func (AA02) GetTableComment() string {
	return TableCommentAA02
}

// SystemConfig 系统配置领域实体
type SystemConfig struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ConfigKey      string     `json:"config_key"`                 // AAB001
	ConfigValue    string     `json:"config_value"`               // AAB002
	ConfigName     string     `json:"config_name"`                // AAB003
	ConfigGroup    string     `json:"config_group"`               // AAB004
	Description    string     `json:"description"`                // AAB005
	Status         string     `json:"status"`                     // AAB006
	LastModifyTime *time.Time `json:"last_modify_time,omitempty"` // AAB007
	CreatorID      string     `json:"creator_id"`                 // AAB008
	ModifierID     string     `json:"modifier_id"`                // AAB009
}

// ToAA02 转换为数据库实体
func (sc *SystemConfig) ToAA02() *AA02 {
	model := gorm.Model{
		ID:        sc.ID,
		CreatedAt: sc.CreatedAt,
		UpdatedAt: sc.UpdatedAt,
	}
	if sc.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *sc.DeletedAt}
	}
	return &AA02{
		Model:  model,
		AAB001: sc.ConfigKey,
		AAB002: sc.ConfigValue,
		AAB003: sc.ConfigName,
		AAB004: sc.ConfigGroup,
		AAB005: sc.Description,
		AAB006: sc.Status,
		AAB007: sc.LastModifyTime,
		AAB008: sc.CreatorID,
		AAB009: sc.ModifierID,
	}
}

// FromAA02 从数据库实体转换
func (sc *SystemConfig) FromAA02(aa02 *AA02) {
	sc.ID = aa02.ID
	sc.CreatedAt = aa02.CreatedAt
	sc.UpdatedAt = aa02.UpdatedAt
	sc.DeletedAt = &aa02.DeletedAt.Time
	sc.ConfigKey = aa02.AAB001
	sc.ConfigValue = aa02.AAB002
	sc.ConfigName = aa02.AAB003
	sc.ConfigGroup = aa02.AAB004
	sc.Description = aa02.AAB005
	sc.Status = aa02.AAB006
	sc.LastModifyTime = aa02.AAB007
	sc.CreatorID = aa02.AAB008
	sc.ModifierID = aa02.AAB009
}

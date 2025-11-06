package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA08 表注释
const TableCommentAA08 = "短信模板表：存储短信发送的模板信息，支持不同场景使用不同模板。每个模板有唯一的模板类型（AAH002），如 register（注册）、forgot（忘记密码）等。模板包含模板名称、腾讯云模板ID、短信签名/标题、模板内容（支持占位符如 {code}）、描述等字段。主键为 AAH001（模板ID），AAH002（模板类型）为唯一索引。短信服务根据业务场景（注册、忘记密码等）自动选择对应的模板进行发送。"

// AA08 短信模板表
type AA08 struct {
	gorm.Model
	AAH001 string `gorm:"column:AAH001;type:varchar(50);primaryKey;comment:模板ID"`                // 模板ID, 主键
	AAH002 string `gorm:"column:AAH002;type:varchar(50);uniqueIndex;comment:模板类型"`               // 模板类型（register, forgot等）
	AAH003 string `gorm:"column:AAH003;type:varchar(100);comment:模板名称"`                          // 模板名称
	AAH004 string `gorm:"column:AAH004;type:varchar(200);comment:腾讯云模板ID"`                       // 腾讯云模板ID
	AAH005 string `gorm:"column:AAH005;type:varchar(100);comment:短信签名/标题"`                       // 短信签名/标题
	AAH006 string `gorm:"column:AAH006;type:text;comment:模板内容"`                                  // 模板内容（支持占位符，如 {code}）
	AAH007 string `gorm:"column:AAH007;type:varchar(500);comment:模板描述"`                          // 模板描述
	AAH008 string `gorm:"column:AAH008;type:varchar(10);default:'1';index;comment:状态 1 启用，2 禁用"` // 状态
	AAH009 string `gorm:"column:AAH009;type:varchar(50);comment:创建人"`                            // 创建人
	AAH010 string `gorm:"column:AAH010;type:varchar(50);comment:修改人"`                            // 修改人
}

// TableName 指定表名
func (AA08) TableName() string {
	return "AA08"
}

// GetTableComment 获取表注释
func (AA08) GetTableComment() string {
	return TableCommentAA08
}

// SMSTemplate 短信模板领域实体
type SMSTemplate struct {
	ID                uint       `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	TemplateID        string     `json:"template_id"`         // AAH001
	Type              string     `json:"type"`                // AAH002
	TemplateName      string     `json:"template_name"`       // AAH003
	TencentTemplateID string     `json:"tencent_template_id"` // AAH004
	Title             string     `json:"title"`               // AAH005
	Content           string     `json:"content"`             // AAH006
	Description       string     `json:"description"`         // AAH007
	Status            string     `json:"status"`              // AAH008
	CreatorID         string     `json:"creator_id"`          // AAH009
	ModifierID        string     `json:"modifier_id"`         // AAH010
}

// ToAA08 转换为数据库实体
func (t *SMSTemplate) ToAA08() *AA08 {
	model := gorm.Model{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *t.DeletedAt}
	}
	return &AA08{
		Model:  model,
		AAH001: t.TemplateID,
		AAH002: t.Type,
		AAH003: t.TemplateName,
		AAH004: t.TencentTemplateID,
		AAH005: t.Title,
		AAH006: t.Content,
		AAH007: t.Description,
		AAH008: t.Status,
		AAH009: t.CreatorID,
		AAH010: t.ModifierID,
	}
}

// FromAA08 从数据库实体转换
func (t *SMSTemplate) FromAA08(aa08 *AA08) {
	t.ID = aa08.ID
	t.CreatedAt = aa08.CreatedAt
	t.UpdatedAt = aa08.UpdatedAt
	t.DeletedAt = &aa08.DeletedAt.Time
	t.TemplateID = aa08.AAH001
	t.Type = aa08.AAH002
	t.TemplateName = aa08.AAH003
	t.TencentTemplateID = aa08.AAH004
	t.Title = aa08.AAH005
	t.Content = aa08.AAH006
	t.Description = aa08.AAH007
	t.Status = aa08.AAH008
	t.CreatorID = aa08.AAH009
	t.ModifierID = aa08.AAH010
}

// IsDisabled 判断模板是否禁用
func (t *SMSTemplate) IsDisabled() bool {
	return t.Status == "2"
}

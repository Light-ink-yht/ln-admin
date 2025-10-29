package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA04 表注释
const TableCommentAA04 = "短信验证码表：存储短信验证码的发送记录，包括验证码、过期时间、使用状态等信息。支持多种验证码类型（register注册、forgot忘记密码等），记录发送的手机号、验证码、状态（未使用/已使用/已过期）、过期时间、使用时间、IP地址、发送次数等。通过验证码ID（AAD001）作为主键，支持按手机号建立索引，便于快速查询和验证。注意：实际验证码存储主要使用Redis，此表用于审计和记录。"

// AA04 短信验证码表
type AA04 struct {
	gorm.Model
	AAD001 string     `gorm:"column:AAD001;type:varchar(50);primaryKey;comment:验证码ID"`                 // 验证码ID, 主键
	AAD002 string     `gorm:"column:AAD002;type:varchar(20);index;comment:手机号"`                        // 手机号
	AAD003 string     `gorm:"column:AAD003;type:varchar(10);comment:验证码"`                              // 验证码
	AAD004 string     `gorm:"column:AAD004;type:varchar(20);comment:验证码类型"`                            // 验证码类型 (register, forgot, login等)
	AAD005 string     `gorm:"column:AAD005;type:varchar(10);default:'1';comment:状态 1 未使用，2 已使用，3 已过期"` // 状态
	AAD006 *time.Time `gorm:"column:AAD006;type:datetime;comment:过期时间"`                                // 过期时间
	AAD007 *time.Time `gorm:"column:AAD007;type:datetime;comment:使用时间"`                                // 使用时间
	AAD008 string     `gorm:"column:AAD008;type:varchar(50);comment:IP地址"`                             // IP地址
	AAD009 int        `gorm:"column:AAD009;type:int;default:0;comment:发送次数"`                           // 发送次数
}

// TableName 指定表名
func (AA04) TableName() string {
	return "AA04"
}

// GetTableComment 获取表注释
func (AA04) GetTableComment() string {
	return TableCommentAA04
}

// SMSCode 短信验证码领域实体
type SMSCode struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CodeID    string     `json:"code_id"`    // AAD001
	Phone     string     `json:"phone"`      // AAD002
	Code      string     `json:"code"`       // AAD003
	Type      string     `json:"type"`       // AAD004
	Status    string     `json:"status"`     // AAD005
	ExpireAt  *time.Time `json:"expire_at"`  // AAD006
	UsedAt    *time.Time `json:"used_at"`    // AAD007
	IP        string     `json:"ip"`         // AAD008
	SendCount int        `json:"send_count"` // AAD009
}

// ToAA04 转换为数据库实体
func (sc *SMSCode) ToAA04() *AA04 {
	model := gorm.Model{
		ID:        sc.ID,
		CreatedAt: sc.CreatedAt,
		UpdatedAt: sc.UpdatedAt,
	}
	if sc.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *sc.DeletedAt}
	}
	return &AA04{
		Model:  model,
		AAD001: sc.CodeID,
		AAD002: sc.Phone,
		AAD003: sc.Code,
		AAD004: sc.Type,
		AAD005: sc.Status,
		AAD006: sc.ExpireAt,
		AAD007: sc.UsedAt,
		AAD008: sc.IP,
		AAD009: sc.SendCount,
	}
}

// FromAA04 从数据库实体转换
func (sc *SMSCode) FromAA04(aa04 *AA04) {
	sc.ID = aa04.ID
	sc.CreatedAt = aa04.CreatedAt
	sc.UpdatedAt = aa04.UpdatedAt
	sc.DeletedAt = &aa04.DeletedAt.Time
	sc.CodeID = aa04.AAD001
	sc.Phone = aa04.AAD002
	sc.Code = aa04.AAD003
	sc.Type = aa04.AAD004
	sc.Status = aa04.AAD005
	sc.ExpireAt = aa04.AAD006
	sc.UsedAt = aa04.AAD007
	sc.IP = aa04.AAD008
	sc.SendCount = aa04.AAD009
}

// IsExpired 判断是否过期
func (sc *SMSCode) IsExpired() bool {
	if sc.ExpireAt == nil {
		return false
	}
	return time.Now().After(*sc.ExpireAt)
}

// IsUsed 判断是否已使用
func (sc *SMSCode) IsUsed() bool {
	return sc.Status == "2"
}

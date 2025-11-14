package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA01 表注释
const TableCommentAA01 = "用户表：存储系统用户的完整信息，包括登录凭证（手机号/邮箱、密码）、个人信息（昵称、姓名、头像、性别、生日）、状态管理（启用/禁用）、登录记录（登录次数、最后登录时间、最后登录IP）、密码管理（密码修改时间）等。支持软删除，主键为 AAA001（用户ID），AAA002（邮箱）和 AAA003（手机号）为唯一索引。"

// AA01 用户实体（与数据库表对应）
type AA01 struct {
	gorm.Model
	AAA001 string     `gorm:"column:AAA001;type:varchar(50);primaryKey;comment:用户id"`                // 用户id, 主键
	AAA002 *string    `gorm:"column:AAA002;type:varchar(100);uniqueIndex;comment:邮箱"`                // 邮箱, 全局唯一
	AAA003 *string    `gorm:"column:AAA003;type:varchar(20);uniqueIndex;comment:手机号"`                // 手机号, 全局唯一
	AAA004 string     `gorm:"column:AAA004;type:varchar(255);not null;comment:密码"`                   // 密码
	AAA005 string     `gorm:"column:AAA005;type:varchar(50);comment:昵称"`                             // 昵称
	AAA006 string     `gorm:"column:AAA006;type:varchar(50);comment:姓名"`                             // 姓名
	AAA007 string     `gorm:"column:AAA007;type:varchar(255);comment:头像"`                            // 头像
	AAA008 string     `gorm:"column:AAA008;type:varchar(10);default:'3';comment:性别 1 男，2 女，3 未知"`    // 性别
	AAA009 *time.Time `gorm:"column:AAA009;type:date;comment:生日"`                                    // 生日
	AAA010 string     `gorm:"column:AAA010;type:varchar(10);default:'1';index;comment:状态 1 启用，2 禁用"` // 状态
	AAA011 string     `gorm:"column:AAA011;type:varchar(500);comment:备注"`                            // 备注
	AAA012 int        `gorm:"column:AAA012;type:int;default:0;comment:登录次数"`                         // 登录次数
	AAA013 *time.Time `gorm:"column:AAA013;type:datetime;index;comment:最后登录时间"`                      // 最后登录时间
	AAA014 string     `gorm:"column:AAA014;type:varchar(50);comment:最后登录ip"`                         // 最后登录ip
	AAA015 *time.Time `gorm:"column:AAA015;type:datetime;comment:密码修改时间"`                            // 密码修改时间
	AAA016 string     `gorm:"column:AAA016;type:varchar(50);index;comment:创建人"`                      // 创建人
	AAA017 string     `gorm:"column:AAA017;type:varchar(50);index;comment:修改人"`                      // 修改人
}

// TableName 指定表名
func (AA01) TableName() string {
	return "AA01"
}

// GetTableComment 获取表注释
func (AA01) GetTableComment() string {
	return TableCommentAA01
}

// User 用户领域实体（业务层使用）
type User struct {
	ID                 uint       `json:"id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
	UserID             string     `json:"user_id"`                        // AAA001
	Email              *string    `json:"email,omitempty"`                // AAA002
	Phone              *string    `json:"phone,omitempty"`                // AAA003
	Password           string     `json:"-"`                              // AAA004 不序列化
	Nickname           string     `json:"nickname"`                       // AAA005
	FullName           string     `json:"full_name"`                      // AAA006
	Avatar             string     `json:"avatar"`                         // AAA007
	Gender             string     `json:"gender"`                         // AAA008
	Birthday           *time.Time `json:"birthday,omitempty"`             // AAA009
	Status             string     `json:"status"`                         // AAA010
	Remarks            string     `json:"remarks"`                        // AAA011
	LoginCount         int        `json:"login_count"`                    // AAA012
	LastLoginTime      *time.Time `json:"last_login_time,omitempty"`      // AAA013
	LastLoginIP        string     `json:"last_login_ip"`                  // AAA014
	PasswordChangeTime *time.Time `json:"password_change_time,omitempty"` // AAA015
	CreatorID          string     `json:"creator_id"`                     // AAA016
	ModifierID         string     `json:"modifier_id"`                    // AAA017
}

// ToAA01 转换为数据库实体
func (u *User) ToAA01() *AA01 {
	model := gorm.Model{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *u.DeletedAt, Valid: true}
	}
	return &AA01{
		Model:  model,
		AAA001: u.UserID,
		AAA002: u.Email,
		AAA003: u.Phone,
		AAA004: u.Password,
		AAA005: u.Nickname,
		AAA006: u.FullName,
		AAA007: u.Avatar,
		AAA008: u.Gender,
		AAA009: u.Birthday,
		AAA010: u.Status,
		AAA011: u.Remarks,
		AAA012: u.LoginCount,
		AAA013: u.LastLoginTime,
		AAA014: u.LastLoginIP,
		AAA015: u.PasswordChangeTime,
		AAA016: u.CreatorID,
		AAA017: u.ModifierID,
	}
}

// FromAA01 从数据库实体转换
func (u *User) FromAA01(aa01 *AA01) {
	u.ID = aa01.ID
	u.CreatedAt = aa01.CreatedAt
	u.UpdatedAt = aa01.UpdatedAt
	u.DeletedAt = &aa01.DeletedAt.Time
	u.UserID = aa01.AAA001
	u.Email = aa01.AAA002
	u.Phone = aa01.AAA003
	u.Password = aa01.AAA004
	u.Nickname = aa01.AAA005
	u.FullName = aa01.AAA006
	u.Avatar = aa01.AAA007
	u.Gender = aa01.AAA008
	u.Birthday = aa01.AAA009
	u.Status = aa01.AAA010
	u.Remarks = aa01.AAA011
	u.LoginCount = aa01.AAA012
	u.LastLoginTime = aa01.AAA013
	u.LastLoginIP = aa01.AAA014
	u.PasswordChangeTime = aa01.AAA015
	u.CreatorID = aa01.AAA016
	u.ModifierID = aa01.AAA017
}

// IsActive 判断用户是否启用
func (u *User) IsActive() bool {
	return u.Status == "1"
}

// IsDisabled 判断用户是否禁用
func (u *User) IsDisabled() bool {
	return u.Status == "2"
}

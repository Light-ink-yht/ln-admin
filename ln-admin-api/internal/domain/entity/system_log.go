package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA03 表注释
const TableCommentAA03 = "系统日志表：记录系统的操作日志、异常日志、访问日志等，用于审计和问题追踪。支持多级别日志（debug、info、warn、error），记录操作模块、操作类型、日志内容、用户信息、请求信息（IP、路径、方法、状态码、User-Agent）等。通过日志ID（AAC001）作为主键，支持按日志级别、日志时间等字段建立索引，便于快速查询和统计分析。"

// AA03 系统日志表
type AA03 struct {
	gorm.Model
	AAC001 string     `gorm:"column:AAC001;type:varchar(50);comment:日志ID, 主键"`                                   // 日志ID, 主键
	AAC002 string     `gorm:"column:AAC002;type:varchar(20);index;index:idx_level_time,priority:1;comment:日志级别"` // 日志级别 (debug, info, warn, error)
	AAC003 string     `gorm:"column:AAC003;type:varchar(200);comment:模块名称"`                                      // 模块名称
	AAC004 string     `gorm:"column:AAC004;type:varchar(100);index;comment:操作类型"`                                // 操作类型 (login, register, create, update, delete等)
	AAC005 string     `gorm:"column:AAC005;type:text;comment:日志内容"`                                              // 日志内容
	AAC006 string     `gorm:"column:AAC006;type:varchar(50);index;comment:用户ID"`                                 // 用户ID
	AAC007 string     `gorm:"column:AAC007;type:varchar(50);comment:IP地址"`                                       // IP地址
	AAC008 string     `gorm:"column:AAC008;type:varchar(500);comment:请求路径"`                                      // 请求路径
	AAC009 string     `gorm:"column:AAC009;type:varchar(20);comment:请求方法"`                                       // 请求方法 (GET, POST等)
	AAC010 int        `gorm:"column:AAC010;type:int;comment:响应状态码"`                                              // 响应状态码
	AAC011 string     `gorm:"column:AAC011;type:varchar(500);comment:用户代理"`                                      // 用户代理
	AAC012 string     `gorm:"column:AAC012;type:varchar(200);comment:错误信息"`                                      // 错误信息
	AAC013 *time.Time `gorm:"column:AAC013;type:datetime;index;index:idx_level_time,priority:2;comment:日志时间"`    // 日志时间
	AAC014 string     `gorm:"column:AAC014;type:varchar(50);comment:创建人"`                                        // 创建人
}

// TableName 指定表名
func (AA03) TableName() string {
	return "AA03"
}

// GetTableComment 获取表注释
func (AA03) GetTableComment() string {
	return TableCommentAA03
}

// SystemLog 系统日志领域实体
type SystemLog struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	LogID      string     `json:"log_id"`      // AAC001
	Level      string     `json:"level"`       // AAC002
	Module     string     `json:"module"`      // AAC003
	Action     string     `json:"action"`      // AAC004
	Content    string     `json:"content"`     // AAC005
	UserID     string     `json:"user_id"`     // AAC006
	IP         string     `json:"ip"`          // AAC007
	Path       string     `json:"path"`        // AAC008
	Method     string     `json:"method"`      // AAC009
	StatusCode int        `json:"status_code"` // AAC010
	UserAgent  string     `json:"user_agent"`  // AAC011
	ErrorMsg   string     `json:"error_msg"`   // AAC012
	LogTime    *time.Time `json:"log_time"`    // AAC013
	CreatorID  string     `json:"creator_id"`  // AAC014
}

// ToAA03 转换为数据库实体
func (sl *SystemLog) ToAA03() *AA03 {
	model := gorm.Model{
		ID:        sl.ID,
		CreatedAt: sl.CreatedAt,
		UpdatedAt: sl.UpdatedAt,
	}
	if sl.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *sl.DeletedAt}
	}
	return &AA03{
		Model:  model,
		AAC001: sl.LogID,
		AAC002: sl.Level,
		AAC003: sl.Module,
		AAC004: sl.Action,
		AAC005: sl.Content,
		AAC006: sl.UserID,
		AAC007: sl.IP,
		AAC008: sl.Path,
		AAC009: sl.Method,
		AAC010: sl.StatusCode,
		AAC011: sl.UserAgent,
		AAC012: sl.ErrorMsg,
		AAC013: sl.LogTime,
		AAC014: sl.CreatorID,
	}
}

// FromAA03 从数据库实体转换
func (sl *SystemLog) FromAA03(aa03 *AA03) {
	sl.ID = aa03.ID
	sl.CreatedAt = aa03.CreatedAt
	sl.UpdatedAt = aa03.UpdatedAt
	sl.DeletedAt = &aa03.DeletedAt.Time
	sl.LogID = aa03.AAC001
	sl.Level = aa03.AAC002
	sl.Module = aa03.AAC003
	sl.Action = aa03.AAC004
	sl.Content = aa03.AAC005
	sl.UserID = aa03.AAC006
	sl.IP = aa03.AAC007
	sl.Path = aa03.AAC008
	sl.Method = aa03.AAC009
	sl.StatusCode = aa03.AAC010
	sl.UserAgent = aa03.AAC011
	sl.ErrorMsg = aa03.AAC012
	sl.LogTime = aa03.AAC013
	sl.CreatorID = aa03.AAC014
}

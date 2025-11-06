package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA11 表注释（S3配置表）
const TableCommentAA11 = "S3存储配置表：存储S3云存储的配置信息，包括存储桶、区域、访问密钥等。系统支持本地存储和S3存储两种方式，通过配置开关控制使用哪种存储方式。主键为 AAK001（配置ID）。"

// AA11 S3配置数据库实体（与数据库表对应）
type AA11 struct {
	gorm.Model
	AAK001 string     `gorm:"column:AAK001;type:varchar(50);primaryKey;comment:配置ID"`                                                              // 配置ID, 主键
	AAK002 string     `gorm:"column:AAK002;type:varchar(10);default:'local';index;index:idx_storage_status,priority:1;comment:存储类型 local本地 s3云存储"` // 存储类型
	AAK003 string     `gorm:"column:AAK003;type:varchar(100);comment:存储桶名称"`                                                                       // 存储桶名称
	AAK004 string     `gorm:"column:AAK004;type:varchar(50);comment:区域"`                                                                           // 区域
	AAK005 string     `gorm:"column:AAK005;type:varchar(200);comment:端点地址"`                                                                        // 端点地址（可选，用于兼容S3兼容服务）
	AAK006 string     `gorm:"column:AAK006;type:varchar(100);comment:访问密钥ID"`                                                                      // 访问密钥ID
	AAK007 string     `gorm:"column:AAK007;type:varchar(200);comment:访问密钥"`                                                                        // 访问密钥
	AAK008 string     `gorm:"column:AAK008;type:varchar(500);comment:基础访问URL"`                                                                     // 基础访问URL（CDN或自定义域名，可选）
	AAK009 string     `gorm:"column:AAK009;type:varchar(10);default:'1';index;index:idx_storage_status,priority:2;comment:状态 1 启用，2 禁用"`           // 状态
	AAK010 string     `gorm:"column:AAK010;type:varchar(500);comment:备注"`                                                                          // 备注
	AAK011 string     `gorm:"column:AAK011;type:varchar(50);comment:创建人"`                                                                          // 创建人
	AAK012 string     `gorm:"column:AAK012;type:varchar(50);comment:修改人"`                                                                          // 修改人
	AAK013 *time.Time `gorm:"column:AAK013;type:datetime;comment:最后修改时间"`                                                                          // 最后修改时间
}

// TableName 指定表名
func (AA11) TableName() string {
	return "AA11"
}

// GetTableComment 获取表注释
func (AA11) GetTableComment() string {
	return TableCommentAA11
}

// S3Config S3配置领域实体（业务层使用）
type S3Config struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ConfigID       string     `json:"config_id"`                  // AAK001
	StorageType    string     `json:"storage_type"`               // AAK002
	Bucket         string     `json:"bucket"`                     // AAK003
	Region         string     `json:"region"`                     // AAK004
	Endpoint       string     `json:"endpoint"`                   // AAK005
	AccessKeyID    string     `json:"access_key_id"`              // AAK006
	SecretKey      string     `json:"secret_key"`                 // AAK007
	BaseURL        string     `json:"base_url"`                   // AAK008
	Status         string     `json:"status"`                     // AAK009
	Remarks        string     `json:"remarks"`                    // AAK010
	CreatorID      string     `json:"creator_id"`                 // AAK011
	ModifierID     string     `json:"modifier_id"`                // AAK012
	LastModifyTime *time.Time `json:"last_modify_time,omitempty"` // AAK013
}

// ToAA11 转换为数据库实体
func (s *S3Config) ToAA11() *AA11 {
	model := gorm.Model{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *s.DeletedAt}
	}
	return &AA11{
		Model:  model,
		AAK001: s.ConfigID,
		AAK002: s.StorageType,
		AAK003: s.Bucket,
		AAK004: s.Region,
		AAK005: s.Endpoint,
		AAK006: s.AccessKeyID,
		AAK007: s.SecretKey,
		AAK008: s.BaseURL,
		AAK009: s.Status,
		AAK010: s.Remarks,
		AAK011: s.CreatorID,
		AAK012: s.ModifierID,
		AAK013: s.LastModifyTime,
	}
}

// FromAA11 从数据库实体转换
func (s *S3Config) FromAA11(aa11 *AA11) {
	s.ID = aa11.ID
	s.CreatedAt = aa11.CreatedAt
	s.UpdatedAt = aa11.UpdatedAt
	if aa11.DeletedAt.Valid {
		s.DeletedAt = &aa11.DeletedAt.Time
	} else {
		s.DeletedAt = nil
	}
	s.ConfigID = aa11.AAK001
	s.StorageType = aa11.AAK002
	s.Bucket = aa11.AAK003
	s.Region = aa11.AAK004
	s.Endpoint = aa11.AAK005
	s.AccessKeyID = aa11.AAK006
	s.SecretKey = aa11.AAK007
	s.BaseURL = aa11.AAK008
	s.Status = aa11.AAK009
	s.Remarks = aa11.AAK010
	s.CreatorID = aa11.AAK011
	s.ModifierID = aa11.AAK012
	s.LastModifyTime = aa11.AAK013
}

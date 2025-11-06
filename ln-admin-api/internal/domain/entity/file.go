package entity

import (
	"time"

	"gorm.io/gorm"
)

// TableComment AA10 表注释（文件表）
const TableCommentAA10File = "文件表：存储系统上传的文件信息，包括文件名、存储路径、存储类型（本地/S3）、文件大小、MIME类型、文件扩展名等。主键为 AAJ001（文件ID）。"

// AA10File 文件数据库实体（与数据库表AA10对应）
type AA10File struct {
	gorm.Model
	AAJ001 string     `gorm:"column:AAJ001;type:varchar(50);primaryKey;comment:文件ID"`                                                                // 文件ID, 主键
	AAJ002 string     `gorm:"column:AAJ002;type:varchar(255);not null;comment:文件名"`                                                                  // 文件名
	AAJ003 string     `gorm:"column:AAJ003;type:varchar(255);not null;comment:原始文件名"`                                                                // 原始文件名
	AAJ004 string     `gorm:"column:AAJ004;type:varchar(500);not null;comment:文件路径"`                                                                 // 文件路径（本地路径或S3路径）
	AAJ005 string     `gorm:"column:AAJ005;type:varchar(20);default:'local';index;index:idx_storage_category,priority:1;comment:存储类型 local本地 s3云存储"` // 存储类型
	AAJ006 int64      `gorm:"column:AAJ006;type:bigint;comment:文件大小(字节)"`                                                                            // 文件大小
	AAJ007 string     `gorm:"column:AAJ007;type:varchar(100);comment:MIME类型"`                                                                        // MIME类型
	AAJ008 string     `gorm:"column:AAJ008;type:varchar(20);comment:文件扩展名"`                                                                          // 文件扩展名
	AAJ009 string     `gorm:"column:AAJ009;type:varchar(50);index;index:idx_storage_category,priority:2;comment:文件分类"`                               // 文件分类（image/document/video/audio/archive/other）
	AAJ010 string     `gorm:"column:AAJ010;type:varchar(50);index;comment:创建人"`                                                                      // 创建人
	AAJ011 string     `gorm:"column:AAJ011;type:varchar(50);comment:修改人"`                                                                            // 修改人
	AAJ012 *time.Time `gorm:"column:AAJ012;type:datetime;index;comment:上传时间"`                                                                        // 上传时间
}

// TableName 指定表名
func (AA10File) TableName() string {
	return "AA10"
}

// GetTableComment 获取表注释
func (AA10File) GetTableComment() string {
	return TableCommentAA10File
}

// File 文件领域实体（业务层使用）
type File struct {
	ID           uint       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	FileID       string     `json:"file_id"`       // AAJ001
	FileName     string     `json:"file_name"`     // AAJ002
	OriginalName string     `json:"original_name"` // AAJ003
	FilePath     string     `json:"file_path"`     // AAJ004
	StorageType  string     `json:"storage_type"`  // AAJ005
	FileSize     int64      `json:"file_size"`     // AAJ006
	MimeType     string     `json:"mime_type"`     // AAJ007
	Extension    string     `json:"extension"`     // AAJ008
	Category     string     `json:"category"`      // AAJ009
	CreatorID    string     `json:"creator_id"`    // AAJ010
	ModifierID   string     `json:"modifier_id"`   // AAJ011
	UploadTime   *time.Time `json:"upload_time"`   // AAJ012
}

// ToAA10 转换为数据库实体
func (f *File) ToAA10() *AA10File {
	model := gorm.Model{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
	if f.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *f.DeletedAt}
	}
	return &AA10File{
		Model:  model,
		AAJ001: f.FileID,
		AAJ002: f.FileName,
		AAJ003: f.OriginalName,
		AAJ004: f.FilePath,
		AAJ005: f.StorageType,
		AAJ006: f.FileSize,
		AAJ007: f.MimeType,
		AAJ008: f.Extension,
		AAJ009: f.Category,
		AAJ010: f.CreatorID,
		AAJ011: f.ModifierID,
		AAJ012: f.UploadTime,
	}
}

// FromAA10 从数据库实体转换
func (f *File) FromAA10(aa10 *AA10File) {
	f.ID = aa10.ID
	f.CreatedAt = aa10.CreatedAt
	f.UpdatedAt = aa10.UpdatedAt
	if aa10.DeletedAt.Valid {
		f.DeletedAt = &aa10.DeletedAt.Time
	} else {
		f.DeletedAt = nil
	}
	f.FileID = aa10.AAJ001
	f.FileName = aa10.AAJ002
	f.OriginalName = aa10.AAJ003
	f.FilePath = aa10.AAJ004
	f.StorageType = aa10.AAJ005
	f.FileSize = aa10.AAJ006
	f.MimeType = aa10.AAJ007
	f.Extension = aa10.AAJ008
	f.Category = aa10.AAJ009
	f.CreatorID = aa10.AAJ010
	f.ModifierID = aa10.AAJ011
	f.UploadTime = aa10.AAJ012
}

// GetFileCategory 根据扩展名获取文件分类
func GetFileCategory(ext string) string {
	extMap := map[string]string{
		// 图片
		"jpg": "image", "jpeg": "image", "png": "image", "gif": "image",
		"bmp": "image", "webp": "image", "svg": "image", "ico": "image",
		// 文档
		"doc": "document", "docx": "document", "xls": "document", "xlsx": "document",
		"ppt": "document", "pptx": "document", "pdf": "document", "txt": "document",
		"md": "document", "csv": "document",
		// 视频
		"mp4": "video", "avi": "video", "mov": "video", "wmv": "video",
		"flv": "video", "mkv": "video", "webm": "video",
		// 音频
		"mp3": "audio", "wav": "audio", "flac": "audio", "aac": "audio",
		"ogg": "audio", "wma": "audio",
		// 压缩包
		"zip": "archive", "rar": "archive", "7z": "archive", "tar": "archive",
		"gz": "archive", "bz2": "archive",
	}

	if category, ok := extMap[ext]; ok {
		return category
	}
	return "other"
}

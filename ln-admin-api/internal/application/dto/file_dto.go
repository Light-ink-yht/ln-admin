package dto

import "time"

// FileResponse 文件响应
type FileResponse struct {
	FileID       string     `json:"file_id"`
	FileName     string     `json:"file_name"`
	OriginalName string     `json:"original_name"`
	FilePath     string     `json:"file_path"`
	FileURL      string     `json:"file_url"` // 文件访问URL
	StorageType  string     `json:"storage_type"`
	FileSize     int64      `json:"file_size"`
	MimeType     string     `json:"mime_type"`
	Extension    string     `json:"extension"`
	Category     string     `json:"category"`
	UploadTime   *time.Time `json:"upload_time"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// FileListParams 文件列表查询参数
type FileListParams struct {
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
	FileName    string `form:"file_name"`
	StorageType string `form:"storage_type"`
	Category    string `form:"category"`
	Extension   string `form:"extension"`
	CreatorID   string `form:"creator_id"`
}

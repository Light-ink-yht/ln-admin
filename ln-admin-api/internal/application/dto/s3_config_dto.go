package dto

import "time"

// S3ConfigRequest S3配置请求
type S3ConfigRequest struct {
	StorageType string `json:"storage_type" binding:"required,oneof=local s3"` // 存储类型：local/s3
	Bucket      string `json:"bucket"`                                         // 存储桶名称（S3必填）
	Region      string `json:"region"`                                         // 区域（S3必填）
	Endpoint    string `json:"endpoint"`                                       // 端点地址（可选，用于兼容S3兼容服务）
	AccessKeyID string `json:"access_key_id"`                                  // 访问密钥ID（S3必填）
	SecretKey   string `json:"secret_key"`                                     // 访问密钥（S3必填）
	BaseURL     string `json:"base_url"`                                       // 基础访问URL（可选，CDN或自定义域名）
	Status      string `json:"status" binding:"required,oneof=1 2"`            // 状态：1启用，2禁用
	Remarks     string `json:"remarks"`                                        // 备注
}

// S3ConfigResponse S3配置响应
type S3ConfigResponse struct {
	ConfigID       string     `json:"config_id"`
	StorageType    string     `json:"storage_type"`
	Bucket         string     `json:"bucket"`
	Region         string     `json:"region"`
	Endpoint       string     `json:"endpoint"`
	AccessKeyID    string     `json:"access_key_id"`
	SecretKey      string     `json:"secret_key"` // 通常不返回，安全考虑
	BaseURL        string     `json:"base_url"`
	Status         string     `json:"status"`
	Remarks        string     `json:"remarks"`
	LastModifyTime *time.Time `json:"last_modify_time"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

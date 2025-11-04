// Package storage 文件存储实现
// 注意：S3存储功能需要安装AWS SDK依赖，如需使用请：
// 1. 执行：go get github.com/aws/aws-sdk-go-v2/aws github.com/aws/aws-sdk-go-v2/config github.com/aws/aws-sdk-go-v2/service/s3
// 2. 取消注释以下所有代码
package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/Light-ink-yht/ln-admin/internal/domain/service"
)

var _ service.FileStorage = (*S3Storage)(nil)

// S3Storage S3文件存储实现（需要安装AWS SDK后取消注释）
type S3Storage struct {
	bucket    string
	region    string
	endpoint  string
	accessKey string
	secretKey string
	baseURL   string
	// client    *s3.Client // 需要AWS SDK
}

// S3Config S3配置
type S3Config struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
	BaseURL   string
}

// NewS3Storage 创建S3存储实例
// 注意：此功能需要安装AWS SDK依赖
// 执行：go get github.com/aws/aws-sdk-go-v2/aws github.com/aws/aws-sdk-go-v2/config github.com/aws/aws-sdk-go-v2/service/s3
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	return nil, fmt.Errorf("S3存储功能需要安装AWS SDK依赖，请先执行：go get github.com/aws/aws-sdk-go-v2/aws github.com/aws/aws-sdk-go-v2/config github.com/aws/aws-sdk-go-v2/service/s3")
}

// Upload 上传文件到S3（需要AWS SDK）
func (s *S3Storage) Upload(ctx context.Context, file io.Reader, fileName string, fileSize int64) (string, error) {
	return "", fmt.Errorf("S3存储功能未启用")
}

// Delete 删除S3文件（需要AWS SDK）
func (s *S3Storage) Delete(ctx context.Context, filePath string) error {
	return fmt.Errorf("S3存储功能未启用")
}

// GetURL 获取文件访问URL（需要AWS SDK）
func (s *S3Storage) GetURL(ctx context.Context, filePath string) (string, error) {
	return "", fmt.Errorf("S3存储功能未启用")
}

// GetFile 获取文件内容（需要AWS SDK）
func (s *S3Storage) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("S3存储功能未启用")
}

package service

import (
	"context"
	"io"
)

// FileStorage 文件存储接口
type FileStorage interface {
	// Upload 上传文件
	Upload(ctx context.Context, file io.Reader, fileName string, fileSize int64) (string, error)

	// Delete 删除文件
	Delete(ctx context.Context, filePath string) error

	// GetURL 获取文件访问URL
	GetURL(ctx context.Context, filePath string) (string, error)

	// GetFile 获取文件内容
	GetFile(ctx context.Context, filePath string) (io.ReadCloser, error)
}

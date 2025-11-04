package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/domain/service"
	"github.com/google/uuid"
)

var _ service.FileStorage = (*LocalStorage)(nil)

// LocalStorage 本地文件存储实现
type LocalStorage struct {
	basePath string // 基础存储路径
	baseURL  string // 基础访问URL
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	// 确保基础路径存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(fmt.Sprintf("创建本地存储目录失败: %v", err))
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Upload 上传文件到本地
func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, fileName string, fileSize int64) (string, error) {
	// 生成唯一文件名
	ext := filepath.Ext(fileName)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期创建目录结构
	dateDir := time.Now().Format("2006/01/02")
	dirPath := filepath.Join(s.basePath, dateDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 文件完整路径
	filePath := filepath.Join(dirPath, newFileName)

	// 创建文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath) // 如果写入失败，删除已创建的文件
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 返回相对路径（相对于basePath）
	relativePath := filepath.Join(dateDir, newFileName)
	return relativePath, nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(ctx context.Context, filePath string) error {
	// filePath是相对路径，需要拼接basePath
	fullPath := filepath.Join(s.basePath, filePath)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，认为删除成功
		}
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetURL 获取文件访问URL
func (s *LocalStorage) GetURL(ctx context.Context, filePath string) (string, error) {
	// filePath是相对路径，拼接baseURL
	return fmt.Sprintf("%s/%s", s.baseURL, filePath), nil
}

// GetFile 获取文件内容
func (s *LocalStorage) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	// filePath是相对路径，需要拼接basePath
	fullPath := filepath.Join(s.basePath, filePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	return file, nil
}

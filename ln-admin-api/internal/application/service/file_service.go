package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/domain/service"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("文件不存在")
	ErrFileTooLarge = errors.New("文件过大")
)

// FileService 文件服务
type FileService struct {
	fileRepo           repository.FileRepository
	s3ConfigRepo       repository.S3ConfigRepository
	localStorage       service.FileStorage
	s3Storage          service.FileStorage
	maxFileSize        int64    // 最大文件大小（字节）
	allowedExts        []string // 允许的文件扩展名
	currentStorageType string   // 当前存储类型：local/s3
}

// NewFileService 创建文件服务
func NewFileService(
	fileRepo repository.FileRepository,
	s3ConfigRepo repository.S3ConfigRepository,
	localStorage service.FileStorage,
	s3Storage service.FileStorage,
	maxFileSize int64,
	allowedExts []string,
) *FileService {
	return &FileService{
		fileRepo:     fileRepo,
		s3ConfigRepo: s3ConfigRepo,
		localStorage: localStorage,
		s3Storage:    s3Storage,
		maxFileSize:  maxFileSize,
		allowedExts:  allowedExts,
	}
}

// getFileStorage 获取当前应该使用的文件存储
func (s *FileService) getFileStorage(ctx context.Context) (service.FileStorage, string, error) {
	// 从数据库获取S3配置，判断当前存储类型
	s3Config, err := s.s3ConfigRepo.GetActiveConfig(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("获取S3配置失败: %w", err)
	}

	// 如果有启用的S3配置，使用S3存储
	if s3Config != nil && s3Config.Status == "1" && s3Config.StorageType == "s3" {
		if s.s3Storage == nil {
			return nil, "", fmt.Errorf("S3存储未配置")
		}
		return s.s3Storage, "s3", nil
	}

	// 否则使用本地存储
	return s.localStorage, "local", nil
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, file io.Reader, fileName string, fileSize int64, creatorID string) (*dto.FileResponse, error) {
	// 验证文件大小
	if fileSize > s.maxFileSize {
		return nil, fmt.Errorf("%w: 文件大小不能超过 %d MB", ErrFileTooLarge, s.maxFileSize/(1024*1024))
	}

	// 验证文件扩展名
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
	if len(s.allowedExts) > 0 {
		allowed := false
		for _, allowedExt := range s.allowedExts {
			if strings.ToLower(allowedExt) == ext {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("不支持的文件类型: %s", ext)
		}
	}

	// 获取MIME类型
	mimeType := mime.TypeByExtension("." + ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 获取文件分类
	category := entity.GetFileCategory(ext)

	// 获取当前应该使用的存储
	fileStorage, storageType, err := s.getFileStorage(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取文件存储失败: %w", err)
	}

	// 上传文件到存储
	storagePath, err := fileStorage.Upload(ctx, file, fileName, fileSize)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 生成文件ID
	fileID := uuid.New().String()

	// 生成存储文件名
	storageFileName := filepath.Base(storagePath)

	// 创建文件记录
	now := time.Now()
	fileEntity := &entity.File{
		FileID:       fileID,
		FileName:     storageFileName,
		OriginalName: fileName,
		FilePath:     storagePath,
		StorageType:  storageType,
		FileSize:     fileSize,
		MimeType:     mimeType,
		Extension:    ext,
		Category:     category,
		CreatorID:    creatorID,
		ModifierID:   creatorID,
		UploadTime:   &now,
	}

	if err := s.fileRepo.Create(ctx, fileEntity); err != nil {
		// 如果数据库保存失败，尝试删除已上传的文件
		if fileStorage != nil {
			fileStorage.Delete(ctx, storagePath)
		}
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	// 获取文件访问URL
	var fileURL string
	if fileStorage != nil {
		fileURL, _ = fileStorage.GetURL(ctx, storagePath)
	}

	return &dto.FileResponse{
		FileID:       fileEntity.FileID,
		FileName:     fileEntity.FileName,
		OriginalName: fileEntity.OriginalName,
		FilePath:     fileEntity.FilePath,
		FileURL:      fileURL,
		StorageType:  fileEntity.StorageType,
		FileSize:     fileEntity.FileSize,
		MimeType:     fileEntity.MimeType,
		Extension:    fileEntity.Extension,
		Category:     fileEntity.Category,
		UploadTime:   fileEntity.UploadTime,
	}, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, fileID string) error {
	// 查找文件记录
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("查询文件失败: %w", err)
	}
	if file == nil {
		return ErrFileNotFound
	}

	// 根据存储类型选择对应的存储
	var fileStorage service.FileStorage
	if file.StorageType == "s3" {
		fileStorage = s.s3Storage
	} else {
		fileStorage = s.localStorage
	}

	// 删除存储中的文件
	if fileStorage != nil {
		if err := fileStorage.Delete(ctx, file.FilePath); err != nil {
			// 记录错误但不阻止数据库删除
			fmt.Printf("删除存储文件失败: %v\n", err)
		}
	}

	// 删除数据库记录
	if err := s.fileRepo.Delete(ctx, fileID); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	return nil
}

// GetFile 获取文件信息
func (s *FileService) GetFile(ctx context.Context, fileID string) (*dto.FileResponse, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}
	if file == nil {
		return nil, ErrFileNotFound
	}

	// 根据存储类型选择对应的存储
	var fileStorage service.FileStorage
	if file.StorageType == "s3" {
		fileStorage = s.s3Storage
	} else {
		fileStorage = s.localStorage
	}

	// 获取文件访问URL
	var fileURL string
	if fileStorage != nil {
		fileURL, _ = fileStorage.GetURL(ctx, file.FilePath)
	}

	return &dto.FileResponse{
		FileID:       file.FileID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileURL:      fileURL,
		StorageType:  file.StorageType,
		FileSize:     file.FileSize,
		MimeType:     file.MimeType,
		Extension:    file.Extension,
		Category:     file.Category,
		UploadTime:   file.UploadTime,
	}, nil
}

// GetFileContent 获取文件内容
func (s *FileService) GetFileContent(ctx context.Context, fileID string) (io.ReadCloser, string, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return nil, "", fmt.Errorf("查询文件失败: %w", err)
	}
	if file == nil {
		return nil, "", ErrFileNotFound
	}

	// 根据存储类型选择对应的存储
	var fileStorage service.FileStorage
	if file.StorageType == "s3" {
		fileStorage = s.s3Storage
	} else {
		fileStorage = s.localStorage
	}

	if fileStorage == nil {
		return nil, "", fmt.Errorf("文件存储未配置")
	}

	reader, err := fileStorage.GetFile(ctx, file.FilePath)
	if err != nil {
		return nil, "", fmt.Errorf("获取文件内容失败: %w", err)
	}

	return reader, file.MimeType, nil
}

// ListFiles 获取文件列表
func (s *FileService) ListFiles(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*dto.FileResponse, int64, error) {
	files, total, err := s.fileRepo.List(ctx, page, pageSize, conditions)
	if err != nil {
		return nil, 0, fmt.Errorf("查询文件列表失败: %w", err)
	}

	fileResponses := make([]*dto.FileResponse, len(files))
	for i, file := range files {
		// 根据存储类型选择对应的存储
		var fileStorage service.FileStorage
		if file.StorageType == "s3" {
			fileStorage = s.s3Storage
		} else {
			fileStorage = s.localStorage
		}

		// 获取文件访问URL
		var fileURL string
		if fileStorage != nil {
			fileURL, _ = fileStorage.GetURL(ctx, file.FilePath)
		}

		fileResponses[i] = &dto.FileResponse{
			FileID:       file.FileID,
			FileName:     file.FileName,
			OriginalName: file.OriginalName,
			FilePath:     file.FilePath,
			FileURL:      fileURL,
			StorageType:  file.StorageType,
			FileSize:     file.FileSize,
			MimeType:     file.MimeType,
			Extension:    file.Extension,
			Category:     file.Category,
			UploadTime:   file.UploadTime,
		}
	}

	return fileResponses, total, nil
}

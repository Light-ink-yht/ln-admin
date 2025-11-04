package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrS3ConfigNotFound = errors.New("S3配置不存在")
	ErrS3ConfigExists   = errors.New("S3配置已存在")
)

// S3ConfigService S3配置服务
type S3ConfigService struct {
	s3ConfigRepo repository.S3ConfigRepository
}

// NewS3ConfigService 创建S3配置服务
func NewS3ConfigService(s3ConfigRepo repository.S3ConfigRepository) *S3ConfigService {
	return &S3ConfigService{
		s3ConfigRepo: s3ConfigRepo,
	}
}

// GetS3Config 获取S3配置
func (s *S3ConfigService) GetS3Config(ctx context.Context) (*dto.S3ConfigResponse, error) {
	config, err := s.s3ConfigRepo.GetActiveConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取S3配置失败: %w", err)
	}
	if config == nil {
		return nil, nil
	}

	return &dto.S3ConfigResponse{
		ConfigID:       config.ConfigID,
		StorageType:    config.StorageType,
		Bucket:         config.Bucket,
		Region:         config.Region,
		Endpoint:       config.Endpoint,
		AccessKeyID:    config.AccessKeyID,
		SecretKey:      "", // 不返回密钥，安全考虑
		BaseURL:        config.BaseURL,
		Status:         config.Status,
		Remarks:        config.Remarks,
		LastModifyTime: config.LastModifyTime,
	}, nil
}

// CreateOrUpdateS3Config 创建或更新S3配置
func (s *S3ConfigService) CreateOrUpdateS3Config(ctx context.Context, req *dto.S3ConfigRequest, modifierID string) (*dto.S3ConfigResponse, error) {
	// 检查是否已存在配置
	existing, err := s.s3ConfigRepo.GetActiveConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询S3配置失败: %w", err)
	}

	now := time.Now()
	var config *entity.S3Config

	if existing != nil {
		// 更新现有配置
		existing.StorageType = req.StorageType
		existing.Bucket = req.Bucket
		existing.Region = req.Region
		existing.Endpoint = req.Endpoint
		existing.AccessKeyID = req.AccessKeyID
		existing.SecretKey = req.SecretKey
		existing.BaseURL = req.BaseURL
		existing.Status = req.Status
		existing.Remarks = req.Remarks
		existing.ModifierID = modifierID
		existing.LastModifyTime = &now

		if err := s.s3ConfigRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("更新S3配置失败: %w", err)
		}
		config = existing
	} else {
		// 创建新配置
		configID := uuid.New().String()
		config = &entity.S3Config{
			ConfigID:       configID,
			StorageType:    req.StorageType,
			Bucket:         req.Bucket,
			Region:         req.Region,
			Endpoint:       req.Endpoint,
			AccessKeyID:    req.AccessKeyID,
			SecretKey:      req.SecretKey,
			BaseURL:        req.BaseURL,
			Status:         req.Status,
			Remarks:        req.Remarks,
			CreatorID:      modifierID,
			ModifierID:     modifierID,
			LastModifyTime: &now,
		}

		if err := s.s3ConfigRepo.Create(ctx, config); err != nil {
			return nil, fmt.Errorf("创建S3配置失败: %w", err)
		}
	}

	return &dto.S3ConfigResponse{
		ConfigID:       config.ConfigID,
		StorageType:    config.StorageType,
		Bucket:         config.Bucket,
		Region:         config.Region,
		Endpoint:       config.Endpoint,
		AccessKeyID:    config.AccessKeyID,
		SecretKey:      "", // 不返回密钥
		BaseURL:        config.BaseURL,
		Status:         config.Status,
		Remarks:        config.Remarks,
		LastModifyTime: config.LastModifyTime,
	}, nil
}

// DeleteS3Config 删除S3配置
func (s *S3ConfigService) DeleteS3Config(ctx context.Context, configID string) error {
	return s.s3ConfigRepo.Delete(ctx, configID)
}

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"

	"gorm.io/gorm"
)

var _ repository.SMSCodeRepository = (*smsCodeRepositoryImpl)(nil)

type smsCodeRepositoryImpl struct {
	db *gorm.DB
}

// NewSMSCodeRepository 创建短信验证码仓库实例
func NewSMSCodeRepository() repository.SMSCodeRepository {
	return &smsCodeRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建验证码
func (r *smsCodeRepositoryImpl) Create(ctx context.Context, code *entity.SMSCode) error {
	aa04 := code.ToAA04()
	if err := r.db.WithContext(ctx).Create(aa04).Error; err != nil {
		return fmt.Errorf("创建验证码失败: %w", err)
	}
	code.FromAA04(aa04)
	return nil
}

// GetLatestByPhoneAndType 根据手机号和类型获取最新的验证码
func (r *smsCodeRepositoryImpl) GetLatestByPhoneAndType(ctx context.Context, phone, codeType string) (*entity.SMSCode, error) {
	var aa04 entity.AA04
	if err := r.db.WithContext(ctx).
		Where("AAD002 = ? AND AAD004 = ?", phone, codeType).
		Order("created_at DESC").
		First(&aa04).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询验证码失败: %w", err)
	}
	code := &entity.SMSCode{}
	code.FromAA04(&aa04)
	return code, nil
}

// UpdateStatus 更新验证码状态
func (r *smsCodeRepositoryImpl) UpdateStatus(ctx context.Context, codeID string, status string) error {
	updates := map[string]interface{}{
		"AAD005": status,
	}
	if status == "2" { // 已使用
		now := time.Now()
		updates["AAD007"] = &now
	}
	if err := r.db.WithContext(ctx).Model(&entity.AA04{}).
		Where("AAD001 = ?", codeID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("更新验证码状态失败: %w", err)
	}
	return nil
}

// MarkAsUsed 标记为已使用
func (r *smsCodeRepositoryImpl) MarkAsUsed(ctx context.Context, codeID string) error {
	return r.UpdateStatus(ctx, codeID, "2")
}

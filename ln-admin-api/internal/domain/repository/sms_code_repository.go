package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// SMSCodeRepository 短信验证码仓库接口
type SMSCodeRepository interface {
	// Create 创建验证码
	Create(ctx context.Context, code *entity.SMSCode) error

	// GetLatestByPhoneAndType 根据手机号和类型获取最新的验证码
	GetLatestByPhoneAndType(ctx context.Context, phone, codeType string) (*entity.SMSCode, error)

	// UpdateStatus 更新验证码状态
	UpdateStatus(ctx context.Context, codeID string, status string) error

	// MarkAsUsed 标记为已使用
	MarkAsUsed(ctx context.Context, codeID string) error
}

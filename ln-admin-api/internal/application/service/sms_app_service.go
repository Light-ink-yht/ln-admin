package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	redisStore "github.com/Light-ink-yht/ln-admin/internal/infrastructure/redis"
	"github.com/Light-ink-yht/ln-admin/pkg/service/sms"
)

var (
	ErrSmsCodeExpired = errors.New("短信验证码已过期")
	ErrSmsCodeUsed    = errors.New("短信验证码已使用")
	ErrSmsCodeInvalid = errors.New("短信验证码错误")
	ErrSmsSendTooFast = errors.New("发送验证码过于频繁，请稍后再试")
)

type SMSAppService struct {
	codeStore  *redisStore.SMSCodeStore
	configRepo repository.SystemConfigRepository
	smsService *sms.TencentSMSService
}

// NewSMSAppService 创建短信应用服务
func NewSMSAppService(configRepo repository.SystemConfigRepository, templateRepo repository.SMSTemplateRepository) *SMSAppService {
	smsService := sms.NewTencentSMSService(configRepo, templateRepo)
	codeStore := redisStore.NewSMSCodeStore()
	return &SMSAppService{
		codeStore:  codeStore,
		configRepo: configRepo,
		smsService: smsService,
	}
}

// SendSMS 发送短信验证码（使用Redis存储，10分钟过期）
func (s *SMSAppService) SendSMS(ctx context.Context, phone, codeType, clientIP string) error {
	// 检查发送频率（1分钟内只能发送一次）- 使用Redis Lua脚本
	minIntervalSeconds := 60
	canSend, waitSeconds, err := s.codeStore.CheckSendFrequency(ctx, phone, codeType, minIntervalSeconds)
	if err != nil {
		return fmt.Errorf("检查发送频率失败: %w", err)
	}
	if !canSend {
		return fmt.Errorf("%s，还需等待 %d 秒", ErrSmsSendTooFast.Error(), waitSeconds)
	}

	// 生成6位验证码
	code := sms.GenerateCode(6)

	// 存储验证码到Redis（10分钟过期）
	if err := s.codeStore.StoreCode(ctx, phone, code, codeType, clientIP); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 发送短信（使用对应类型的模板）
	if err := s.smsService.SendSMS(ctx, phone, code, codeType); err != nil {
		// 如果发送失败，删除Redis中的验证码
		_ = s.codeStore.DeleteCode(ctx, phone, codeType)
		return fmt.Errorf("发送短信失败: %w", err)
	}

	return nil
}

// VerifySMS 验证短信验证码（使用Redis Lua脚本保证原子性）
func (s *SMSAppService) VerifySMS(ctx context.Context, phone, code, codeType string) error {
	// 使用Redis Lua脚本验证验证码（原子操作）
	if err := s.codeStore.VerifyCode(ctx, phone, code, codeType); err != nil {
		// 转换Redis错误为应用层错误
		switch err {
		case redisStore.ErrCodeNotFound:
			return ErrSmsCodeInvalid
		case redisStore.ErrCodeExpired:
			return ErrSmsCodeExpired
		case redisStore.ErrCodeUsed:
			return ErrSmsCodeUsed
		case redisStore.ErrCodeInvalid:
			return ErrSmsCodeInvalid
		default:
			return err
		}
	}

	return nil
}

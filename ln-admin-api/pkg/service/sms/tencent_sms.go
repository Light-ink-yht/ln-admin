package sms

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"go.uber.org/zap"
)

type TencentSMSService struct {
	configRepo   repository.SystemConfigRepository
	templateRepo repository.SMSTemplateRepository
}

// NewTencentSMSService 创建腾讯云短信服务
func NewTencentSMSService(configRepo repository.SystemConfigRepository, templateRepo repository.SMSTemplateRepository) *TencentSMSService {
	return &TencentSMSService{
		configRepo:   configRepo,
		templateRepo: templateRepo,
	}
}

// GetConfig 获取短信基础配置（不包含模板ID）
func (s *TencentSMSService) GetConfig(ctx context.Context) (string, string, string, string, error) {
	configs, err := s.configRepo.GetByGroup(ctx, "sms")
	if err != nil {
		return "", "", "", "", fmt.Errorf("获取短信配置失败: %w", err)
	}

	var secretId, secretKey, appId, signName string
	for _, config := range configs {
		switch config.ConfigKey {
		case "sms_secret_id":
			secretId = config.ConfigValue
		case "sms_secret_key":
			secretKey = config.ConfigValue
		case "sms_app_id":
			appId = config.ConfigValue
		case "sms_sign_name":
			signName = config.ConfigValue
		}
	}

	if secretId == "" || secretKey == "" || appId == "" || signName == "" {
		return "", "", "", "", fmt.Errorf("短信配置不完整")
	}

	return secretId, secretKey, appId, signName, nil
}

// SendSMS 发送短信（根据codeType使用不同的模板）
func (s *TencentSMSService) SendSMS(ctx context.Context, phone, code, codeType string) error {
	// 检查短信功能是否启用
	if !config.Cfg.SMS.Enabled {
		// 短信功能未启用，只记录日志，不实际发送
		logger.Info("短信功能未启用（开发模式），仅记录日志",
			zap.String("phone", phone),
			zap.String("code", code),
			zap.String("code_type", codeType),
			zap.String("message", fmt.Sprintf("验证码：%s（开发环境，未实际发送）", code)))
		return nil
	}

	// 短信功能已启用，需要验证配置
	secretId, secretKey, appId, _, err := s.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("获取短信配置失败，请先在系统配置中填写短信相关配置: %w", err)
	}

	// 根据codeType获取完整的模板信息（从数据库模板表获取）
	template, err := s.GetTemplate(ctx, codeType)
	if err != nil {
		return fmt.Errorf("获取短信模板失败: %w", err)
	}

	// 格式化短信内容（替换占位符）
	content := FormatContent(template.Content, map[string]string{
		"code": code,
	})

	// 实例化认证对象
	credential := common.NewCredential(secretId, secretKey)

	// 实例化客户端配置对象
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	// 实例化sms client对象
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化请求对象
	request := sms.NewSendSmsRequest()

	// 设置请求参数
	request.SmsSdkAppId = common.StringPtr(appId)
	request.SignName = common.StringPtr(template.Title)               // 使用模板中的标题作为签名
	request.TemplateId = common.StringPtr(template.TencentTemplateID) // 使用腾讯云模板ID

	// 根据模板内容解析参数（简单的占位符替换，实际使用模板ID对应的参数）
	// 如果模板内容包含 {code}，则传入验证码作为参数
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{phone})

	// 发送请求
	response, err := client.SendSms(request)
	if err != nil {
		logger.Error("腾讯云短信发送失败",
			zap.Error(err),
			zap.String("phone", phone),
			zap.String("template_id", template.TencentTemplateID))
		return fmt.Errorf("短信发送失败: %w", err)
	}

	// 检查响应
	if response.Response == nil || len(response.Response.SendStatusSet) == 0 {
		return fmt.Errorf("短信发送响应异常")
	}

	sendStatus := response.Response.SendStatusSet[0]
	if *sendStatus.Code != "Ok" {
		logger.Error("腾讯云短信发送失败",
			zap.String("code", *sendStatus.Code),
			zap.String("message", *sendStatus.Message),
			zap.String("template_id", template.TencentTemplateID))
		return fmt.Errorf("短信发送失败: %s", *sendStatus.Message)
	}

	logger.Info("短信发送成功",
		zap.String("phone", phone),
		zap.String("code_type", codeType),
		zap.String("template_id", template.TencentTemplateID),
		zap.String("content", content))

	return nil
}

// GenerateCode 生成验证码
func GenerateCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}

// GetTemplate 根据类型获取完整的短信模板信息（从数据库模板表获取）
func (s *TencentSMSService) GetTemplate(ctx context.Context, codeType string) (*entity.SMSTemplate, error) {
	template, err := s.templateRepo.FindByType(ctx, codeType)
	if err != nil {
		return nil, fmt.Errorf("查询短信模板失败: %w", err)
	}
	if template == nil {
		return nil, fmt.Errorf("未找到类型为 %s 的短信模板", codeType)
	}

	// 验证模板状态
	if template.IsDisabled() {
		return nil, fmt.Errorf("短信模板 %s 已禁用", codeType)
	}

	return template, nil
}

// FormatContent 格式化短信内容，替换占位符
func FormatContent(template string, params map[string]string) string {
	content := template
	for key, value := range params {
		content = strings.ReplaceAll(content, fmt.Sprintf("{%s}", key), value)
	}
	return content
}

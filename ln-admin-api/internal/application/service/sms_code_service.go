package service

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
)

// SMSCodeService 短信验证码服务
type SMSCodeService struct {
	codeRepo repository.SMSCodeRepository
}

// NewSMSCodeService 创建短信验证码服务
func NewSMSCodeService(codeRepo repository.SMSCodeRepository) *SMSCodeService {
	return &SMSCodeService{
		codeRepo: codeRepo,
	}
}

// GetCodeList 获取验证码列表
func (s *SMSCodeService) GetCodeList(ctx context.Context, req *dto.SMSCodeListRequest) ([]*dto.SMSCodeResponse, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	conditions := make(map[string]interface{})
	if req.Phone != "" {
		conditions["phone"] = req.Phone
	}
	if req.Type != "" {
		conditions["type"] = req.Type
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}

	codes, total, err := s.codeRepo.List(ctx, page, pageSize, conditions)
	if err != nil {
		return nil, 0, fmt.Errorf("查询验证码列表失败: %w", err)
	}

	responses := make([]*dto.SMSCodeResponse, len(codes))
	for i, code := range codes {
		responses[i] = s.toCodeResponse(code)
	}

	return responses, total, nil
}

// toCodeResponse 转换为响应DTO
func (s *SMSCodeService) toCodeResponse(code *entity.SMSCode) *dto.SMSCodeResponse {
	expireAt := ""
	if code.ExpireAt != nil {
		expireAt = dto.FormatTime(*code.ExpireAt)
	}
	usedAt := ""
	if code.UsedAt != nil {
		usedAt = dto.FormatTime(*code.UsedAt)
	}

	return &dto.SMSCodeResponse{
		CodeID:    code.CodeID,
		Phone:     code.Phone,
		Code:      code.Code,
		Type:      code.Type,
		Status:    code.Status,
		ExpireAt:  expireAt,
		UsedAt:    usedAt,
		IP:        code.IP,
		SendCount: code.SendCount,
		CreatedAt: dto.FormatTime(code.CreatedAt),
	}
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/config"
	"github.com/Light-ink-yht/ln-admin/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserDisabled      = errors.New("用户已禁用")
	ErrPasswordIncorrect = errors.New("密码错误")
	ErrUserExists        = errors.New("用户已存在")
	ErrInvalidCaptcha    = errors.New("验证码错误")
	ErrInvalidSmsCode    = errors.New("短信验证码错误")
)

type UserAppService struct {
	userRepo      repository.UserRepository
	smsAppService *SMSAppService
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(userRepo repository.UserRepository, smsAppService *SMSAppService) *UserAppService {
	return &UserAppService{
		userRepo:      userRepo,
		smsAppService: smsAppService,
	}
}

// Login 用户登录（只使用图片验证码，不使用短信验证码）
func (s *UserAppService) Login(ctx context.Context, phone, password, loginIP string) (*dto.LoginResponse, error) {
	// 根据手机号查找用户
	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 检查用户状态
	if user.IsDisabled() {
		return nil, ErrUserDisabled
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return nil, ErrPasswordIncorrect
	}

	// 登录不需要验证短信验证码，只验证图片验证码（在handler层验证）

	// 更新登录信息
	if err := s.userRepo.UpdateLoginInfo(ctx, user.UserID, loginIP); err != nil {
		logger.Warn("更新登录信息失败", zap.Error(err))
	}

	// 重新查询用户以获取最新的登录信息
	user, err = s.userRepo.FindByID(ctx, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 生成双Token
	tokenPair, err := utils.GenerateTokenPair(user.UserID)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	// 转换为响应DTO
	userResp := s.toUserResponse(user)

	return &dto.LoginResponse{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresIn:        tokenPair.ExpiresIn,
		RefreshExpiresIn: tokenPair.RefreshExpiresIn,
		User:             *userResp,
	}, nil
}

// RefreshToken 刷新Token
func (s *UserAppService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	// 解析RefreshToken
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		if err == utils.ErrExpiredToken {
			return nil, fmt.Errorf("RefreshToken已过期，请重新登录")
		}
		return nil, fmt.Errorf("无效的RefreshToken")
	}

	// 生成新的双Token
	tokenPair, err := utils.GenerateTokenPair(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresIn:        tokenPair.ExpiresIn,
		RefreshExpiresIn: tokenPair.RefreshExpiresIn,
	}, nil
}

// Register 用户注册
func (s *UserAppService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// 验证短信验证码
	if err := s.smsAppService.VerifySMS(ctx, req.Phone, req.Code, "register"); err != nil {
		return nil, err
	}

	// 检查手机号是否已存在
	existingUser, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 生成用户ID
	userID := uuid.New().String()

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password, config.Cfg.Password.Cost)
	if err != nil {
		return nil, fmt.Errorf("加密密码失败: %w", err)
	}

	// 创建用户
	user := &entity.User{
		UserID:     userID,
		Phone:      &req.Phone,
		Password:   hashedPassword,
		Status:     "1", // 默认启用
		Gender:     "3", // 默认未知
		LoginCount: 0,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 重新查询用户
	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return s.toUserResponse(user), nil
}

// ForgotPassword 忘记密码
func (s *UserAppService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	// 验证短信验证码
	if err := s.smsAppService.VerifySMS(ctx, req.Phone, req.Code, "forgot"); err != nil {
		return err
	}

	// 查找用户
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.Password, config.Cfg.Password.Cost)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	// 更新密码
	now := time.Now()
	user.Password = hashedPassword
	user.PasswordChangeTime = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// GetUserInfo 获取用户信息
func (s *UserAppService) GetUserInfo(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.toUserResponse(user), nil
}

// ListUsers 获取用户列表
func (s *UserAppService) ListUsers(ctx context.Context, req *dto.UserListRequest) ([]*dto.UserResponse, int64, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.Phone != "" {
		conditions["phone"] = req.Phone
	}
	if req.Email != "" {
		conditions["email"] = req.Email
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}
	if req.Gender != "" {
		conditions["gender"] = req.Gender
	}
	if req.Nickname != "" {
		conditions["nickname"] = req.Nickname
	}
	if req.FullName != "" {
		conditions["full_name"] = req.FullName
	}

	// 查询用户列表
	users, total, err := s.userRepo.List(ctx, page, pageSize, conditions)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 转换为响应DTO
	userRespList := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userRespList[i] = s.toUserResponse(user)
	}

	return userRespList, total, nil
}

// toUserResponse 转换为用户响应DTO
func (s *UserAppService) toUserResponse(user *entity.User) *dto.UserResponse {
	var deletedAt *string
	if user.DeletedAt != nil {
		deletedAtStr := dto.FormatTime(*user.DeletedAt)
		deletedAt = &deletedAtStr
	}

	var lastLoginTime *string
	if user.LastLoginTime != nil {
		lastLoginTimeStr := dto.FormatTime(*user.LastLoginTime)
		lastLoginTime = &lastLoginTimeStr
	}

	var passwordChange *string
	if user.PasswordChangeTime != nil {
		passwordChangeStr := dto.FormatTime(*user.PasswordChangeTime)
		passwordChange = &passwordChangeStr
	}

	var birthday *string
	if user.Birthday != nil {
		birthdayStr := dto.FormatDate(user.Birthday)
		if birthdayStr != nil {
			birthday = birthdayStr
		}
	}

	return &dto.UserResponse{
		ID:             user.ID,
		CreatedAt:      dto.FormatTime(user.CreatedAt),
		UpdatedAt:      dto.FormatTime(user.UpdatedAt),
		DeletedAt:      deletedAt,
		UserId:         user.UserID,
		Email:          user.Email,
		Phone:          user.Phone,
		Nickname:       user.Nickname,
		FullName:       user.FullName,
		Avatar:         user.Avatar,
		Gender:         user.Gender,
		Birthday:       birthday,
		Status:         user.Status,
		Remarks:        user.Remarks,
		LoginCount:     user.LoginCount,
		LastLoginTime:  lastLoginTime,
		LastLoginIP:    user.LastLoginIP,
		PasswordChange: passwordChange,
		CreatorID:      user.CreatorID,
		ModifierID:     user.ModifierID,
	}
}

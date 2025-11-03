package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/casbin"
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
	userRepo          repository.UserRepository
	smsAppService     *SMSAppService
	permissionService *PermissionService
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(userRepo repository.UserRepository, smsAppService *SMSAppService, permissionService *PermissionService) *UserAppService {
	return &UserAppService{
		userRepo:          userRepo,
		smsAppService:     smsAppService,
		permissionService: permissionService,
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
// ListUsersWithConditions 使用条件映射获取用户列表（支持数组查询）
func (s *UserAppService) ListUsersWithConditions(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*dto.UserResponse, int64, error) {
	// 查询用户列表
	users, total, err := s.userRepo.List(ctx, page, pageSize, conditions)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 转换为响应DTO
	userResps := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResps[i] = s.toUserResponse(user)
	}

	return userResps, total, nil
}

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
	// 注意：这里接收的req可能包含数组类型（从handler层传入）
	// handler层会处理QueryArray，如果存在数组则传入数组，否则传入字符串
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

// CreateUser 创建用户
func (s *UserAppService) CreateUser(ctx context.Context, req *dto.CreateUserRequest, creatorID string) (*dto.UserResponse, error) {
	// 检查手机号是否已存在
	existing, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if req.Email != nil && *req.Email != "" {
		existing, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
		if existing != nil {
			return nil, fmt.Errorf("邮箱已被使用")
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password, config.Cfg.Password.Cost)
	if err != nil {
		return nil, fmt.Errorf("加密密码失败: %w", err)
	}

	// 生成用户ID
	userID := uuid.New().String()

	// 解析生日
	var birthday *time.Time
	if req.Birthday != nil && *req.Birthday != "" {
		parsed, err := time.Parse("2006-01-02", *req.Birthday)
		if err == nil {
			birthday = &parsed
		}
	}

	// 设置默认值
	gender := req.Gender
	if gender == "" {
		gender = "3"
	}
	status := req.Status
	if status == "" {
		status = "1"
	}

	// 创建用户实体
	user := &entity.User{
		UserID:     userID,
		Phone:      &req.Phone,
		Email:      req.Email,
		Password:   hashedPassword,
		Nickname:   req.Nickname,
		FullName:   req.FullName,
		Avatar:     req.Avatar,
		Gender:     gender,
		Birthday:   birthday,
		Status:     status,
		Remarks:    req.Remarks,
		CreatorID:  creatorID,
		ModifierID: creatorID,
	}

	// 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 分配角色
	if len(req.RoleIds) > 0 {
		for _, roleID := range req.RoleIds {
			if err := s.permissionService.AssignRoleToUser(ctx, userID, roleID); err != nil {
				logger.Warn("为用户分配角色失败",
					zap.String("user_id", userID),
					zap.String("role_id", roleID),
					zap.Error(err))
				// 继续执行，不因为角色分配失败而回滚用户创建
			}
		}
	}

	// 重新查询用户以获取完整信息
	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return s.toUserResponse(user), nil
}

// UpdateUser 更新用户
func (s *UserAppService) UpdateUser(ctx context.Context, userID string, req *dto.UpdateUserRequest, modifierID string) (*dto.UserResponse, error) {
	// 查询用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 检查邮箱是否已被其他用户使用（如果提供了邮箱）
	if req.Email != nil && *req.Email != "" {
		existing, err := s.userRepo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
		if existing != nil && existing.UserID != userID {
			return nil, fmt.Errorf("邮箱已被其他用户使用")
		}
		user.Email = req.Email
	}

	// 更新字段
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Remarks != nil {
		user.Remarks = *req.Remarks
	}
	if req.Birthday != nil && *req.Birthday != "" {
		parsed, err := time.Parse("2006-01-02", *req.Birthday)
		if err == nil {
			user.Birthday = &parsed
		}
	}

	user.ModifierID = modifierID

	// 更新用户
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 更新角色（如果提供了）
	if req.RoleIds != nil {
		// 获取用户当前角色
		currentRoles, err := s.permissionService.GetUserRoles(ctx, userID)
		if err != nil {
			logger.Warn("获取用户当前角色失败", zap.Error(err))
		} else {
			// 找出需要移除的角色
			currentRoleIDs := make(map[string]bool)
			for _, role := range currentRoles {
				currentRoleIDs[role.RoleID] = true
			}

			newRoleIDs := make(map[string]bool)
			for _, roleID := range req.RoleIds {
				newRoleIDs[roleID] = true
				// 如果角色不在当前角色中，则添加
				if !currentRoleIDs[roleID] {
					if err := s.permissionService.AssignRoleToUser(ctx, userID, roleID); err != nil {
						logger.Warn("为用户分配角色失败",
							zap.String("user_id", userID),
							zap.String("role_id", roleID),
							zap.Error(err))
					}
				}
			}

			// 移除不在新角色列表中的角色
			for _, role := range currentRoles {
				if !newRoleIDs[role.RoleID] {
					if err := s.permissionService.RemoveRoleFromUser(ctx, userID, role.RoleID); err != nil {
						logger.Warn("移除用户角色失败",
							zap.String("user_id", userID),
							zap.String("role_id", role.RoleID),
							zap.Error(err))
					}
				}
			}
		}
	}

	// 重新查询用户以获取完整信息
	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return s.toUserResponse(user), nil
}

// DeleteUser 删除用户
func (s *UserAppService) DeleteUser(ctx context.Context, userID string) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 获取用户角色并移除
	roles, err := s.permissionService.GetUserRoles(ctx, userID)
	if err == nil {
		for _, role := range roles {
			if err := s.permissionService.RemoveRoleFromUser(ctx, userID, role.RoleID); err != nil {
				logger.Warn("移除用户角色失败",
					zap.String("user_id", userID),
					zap.String("role_id", role.RoleID),
					zap.Error(err))
				// 继续执行，不因为角色移除失败而阻止用户删除
			}
		}
	}

	// 删除用户
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

// GetUserDetail 获取用户详情（包含角色和权限）
func (s *UserAppService) GetUserDetail(ctx context.Context, userID string) (*dto.UserDetailResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 转换为响应DTO
	userResp := s.toUserResponse(user)

	// 获取用户角色
	roles, err := s.permissionService.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Warn("获取用户角色失败", zap.Error(err))
		roles = []*entity.Role{}
	}

	// 转换为RoleInfo
	roleInfos := make([]dto.RoleInfo, len(roles))
	for i, role := range roles {
		roleInfos[i] = dto.RoleInfo{
			RoleID:   role.RoleID,
			RoleKey:  role.RoleKey,
			RoleName: role.RoleName,
			Status:   role.Status,
		}
	}

	// 获取用户权限（通过角色）
	permissionInfos := make([]dto.PermissionInfo, 0)
	roleKeys := make(map[string]bool)
	for _, role := range roles {
		roleKeys[role.RoleKey] = true
	}

	// 获取用户权限（包括角色权限和直接权限）
	enforcer := casbin.GetEnforcer()
	if enforcer != nil {
		permissionMap := make(map[string]dto.PermissionInfo)

		// 1. 获取角色权限
		for roleKey := range roleKeys {
			// 获取角色的所有权限策略
			policies, _ := enforcer.GetPermissionsForUser(roleKey)
			for _, policy := range policies {
				// policy格式：[roleKey, resourcePath, method]
				if len(policy) >= 3 {
					resourcePath := policy[1]
					method := policy[2]
					key := resourcePath + ":" + method
					if _, exists := permissionMap[key]; !exists {
						// 查询权限详情（通过resourcePath和method在数据库中查找）
						permission, _ := s.permissionService.GetPermissionByResourceAndMethod(ctx, resourcePath, method)
						if permission != nil {
							permissionMap[key] = dto.PermissionInfo{
								PermissionID:   permission.PermissionID,
								PermissionKey:  permission.PermissionKey,
								PermissionName: permission.PermissionName,
								ResourcePath:   permission.ResourcePath,
								Method:         permission.Method,
								Description:    permission.Description,
							}
						} else {
							// 如果数据库中找不到，使用简化信息
							permissionMap[key] = dto.PermissionInfo{
								PermissionID:   "",
								PermissionKey:  key,
								PermissionName: fmt.Sprintf("%s %s", method, resourcePath),
								ResourcePath:   resourcePath,
								Method:         method,
								Description:    "",
							}
						}
					}
				}
			}
		}
		// 2. 获取用户直接权限（p, user_id, resource_path, method）
		// 通过Casbin的GetFilteredPolicy获取用户直接权限策略
		// 策略格式：p, user_id, resource_path, method
		userDirectPolicies, _ := enforcer.GetFilteredPolicy(0, userID)
		for _, policy := range userDirectPolicies {
			// policy格式：[userID, resourcePath, method]
			if len(policy) >= 3 {
				resourcePath := policy[1]
				method := policy[2]
				key := resourcePath + ":" + method
				if _, exists := permissionMap[key]; !exists {
					// 查询权限详情
					permission, _ := s.permissionService.GetPermissionByResourceAndMethod(ctx, resourcePath, method)
					if permission != nil {
						permissionMap[key] = dto.PermissionInfo{
							PermissionID:   permission.PermissionID,
							PermissionKey:  permission.PermissionKey,
							PermissionName: permission.PermissionName,
							ResourcePath:   permission.ResourcePath,
							Method:         permission.Method,
							Description:    permission.Description,
						}
					} else {
						// 如果数据库中找不到，使用简化信息
						permissionMap[key] = dto.PermissionInfo{
							PermissionID:   "",
							PermissionKey:  key,
							PermissionName: fmt.Sprintf("%s %s", method, resourcePath),
							ResourcePath:   resourcePath,
							Method:         method,
							Description:    "",
						}
					}
				}
			}
		}

		// 转换为数组
		for _, perm := range permissionMap {
			permissionInfos = append(permissionInfos, perm)
		}
	}

	return &dto.UserDetailResponse{
		UserResponse: *userResp,
		Roles:        roleInfos,
		Permissions:  permissionInfos,
	}, nil
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

// GrantPermissions 给用户授权（直接分配权限）
func (s *UserAppService) GrantPermissions(ctx context.Context, userID string, permissionIDs []string) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 调用权限服务直接为用户分配权限
	if err := s.permissionService.AssignPermissionsToUser(ctx, userID, permissionIDs); err != nil {
		return fmt.Errorf("分配权限失败: %w", err)
	}

	return nil
}

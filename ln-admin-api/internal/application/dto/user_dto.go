package dto

import "time"

// LoginRequest 登录请求（只使用图片验证码）
type LoginRequest struct {
	Phone       string `json:"phone" binding:"required,min=11,max=11"`
	Password    string `json:"password" binding:"required,min=6,max=20"`
	CaptchaCode string `json:"captchaCode" binding:"required"` // 图片验证码
	CaptchaID   string `json:"captchaId" binding:"required"`   // 图片验证码ID
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone           string `json:"phone" binding:"required,min=11,max=11"`
	Code            string `json:"code" binding:"required"`
	Password        string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Phone           string `json:"phone" binding:"required,min=11,max=11"`
	Code            string `json:"code" binding:"required"`
	Password        string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

// SendSmsRequest 发送短信验证码请求（仅支持注册和忘记密码）
type SendSmsRequest struct {
	Phone string `json:"phone" binding:"required,min=11,max=11"`
	Type  string `json:"type" binding:"required,oneof=register forgot"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID             uint    `json:"ID"`
	CreatedAt      string  `json:"CreatedAt"`
	UpdatedAt      string  `json:"UpdatedAt"`
	DeletedAt      *string `json:"DeletedAt,omitempty"`
	UserId         string  `json:"userId"`
	Email          *string `json:"email,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Nickname       string  `json:"nickname"`
	FullName       string  `json:"fullName"`
	Avatar         string  `json:"avatar"`
	Gender         string  `json:"gender"`
	Birthday       *string `json:"birthday,omitempty"`
	Status         string  `json:"status"`
	Remarks        string  `json:"remarks"`
	LoginCount     int     `json:"loginCount"`
	LastLoginTime  *string `json:"lastLoginTime,omitempty"`
	LastLoginIP    string  `json:"lastLoginIp"`
	PasswordChange *string `json:"passwordChange,omitempty"`
	CreatorID      string  `json:"creatorId"`
	ModifierID     string  `json:"modifierId"`
}

// LoginResponse 登录响应（双Token）
type LoginResponse struct {
	AccessToken      string       `json:"accessToken"`
	RefreshToken     string       `json:"refreshToken"`
	ExpiresIn        int64        `json:"expiresIn"`
	RefreshExpiresIn int64        `json:"refreshExpiresIn"`
	User             UserResponse `json:"user"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresIn        int64  `json:"expiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID   string `json:"captchaId"`
	CaptchaCode string `json:"captcha_code"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Phone    string `form:"phone"`
	Email    string `form:"email"`
	Status   string `form:"status"`
	Gender   string `form:"gender"`
	Nickname string `form:"nickname"`
	FullName string `form:"full_name"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Phone    string   `json:"phone" binding:"required,min=11,max=11"`
	Email    *string  `json:"email,omitempty"`
	Password string   `json:"password" binding:"required,min=6,max=20"`
	Nickname string   `json:"nickname" binding:"required"`
	FullName string   `json:"fullName" binding:"required"`
	Gender   string   `json:"gender" binding:"omitempty,oneof=1 2 3"`
	Status   string   `json:"status" binding:"omitempty,oneof=1 2"`
	Avatar   string   `json:"avatar,omitempty"`
	Birthday *string  `json:"birthday,omitempty"`
	Remarks  string   `json:"remarks,omitempty"`
	RoleIds  []string `json:"roleIds,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    *string  `json:"email,omitempty"`
	Nickname *string  `json:"nickname,omitempty"`
	FullName *string  `json:"fullName,omitempty"`
	Gender   *string  `json:"gender,omitempty" binding:"omitempty,oneof=1 2 3"`
	Status   *string  `json:"status,omitempty" binding:"omitempty,oneof=1 2"`
	Avatar   *string  `json:"avatar,omitempty"`
	Birthday *string  `json:"birthday,omitempty"`
	Remarks  *string  `json:"remarks,omitempty"`
	RoleIds  []string `json:"roleIds,omitempty"`
}

// UserDetailResponse 用户详情响应（包含角色和权限）
type UserDetailResponse struct {
	UserResponse
	Roles       []RoleInfo       `json:"roles,omitempty"`
	Permissions []PermissionInfo `json:"permissions,omitempty"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	RoleID   string `json:"roleId"`
	RoleKey  string `json:"roleKey"`
	RoleName string `json:"roleName"`
	Status   string `json:"status,omitempty"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	PermissionID   string `json:"permissionId"`
	PermissionKey  string `json:"permissionKey"`
	PermissionName string `json:"permissionName"`
	ResourcePath   string `json:"resourcePath"`
	Method         string `json:"method"`
	Description    string `json:"description,omitempty"`
	Category       string `json:"category,omitempty"` // 分类
}

// FormatTime 格式化时间为字符串
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 格式化日期为字符串
func FormatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format("2006-01-02")
	return &formatted
}

// GrantPermissionsRequest 授权请求
type GrantPermissionsRequest struct {
	PermissionIds []string `json:"permissionIds" binding:"required"`
}

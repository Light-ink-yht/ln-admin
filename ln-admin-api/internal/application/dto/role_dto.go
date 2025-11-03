package dto

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	RoleKey  string `form:"role_key"`
	RoleName string `form:"role_name"`
	Status   string `form:"status"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID           uint    `json:"id"`
	CreatedAt    string  `json:"createdAt,omitempty"`
	UpdatedAt    string  `json:"updatedAt,omitempty"`
	DeletedAt    *string `json:"deletedAt,omitempty"`
	RoleID       string  `json:"roleId"`
	RoleKey      string  `json:"roleKey"`
	RoleName     string  `json:"roleName"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	CreatorID    string  `json:"creatorId,omitempty"`
	CreatorName  string  `json:"creatorName,omitempty"`
	ModifierID   string  `json:"modifierId,omitempty"`
	ModifierName string  `json:"modifierName,omitempty"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleKey     string `json:"roleKey" binding:"required,min=2,max=100"`
	RoleName    string `json:"roleName" binding:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=1 2"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleName    *string `json:"roleName,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty" binding:"omitempty,oneof=1 2"`
}

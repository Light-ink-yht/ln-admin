package dto

// PermissionListRequest 权限列表请求
type PermissionListRequest struct {
	Page           int    `form:"page" binding:"omitempty,min=1"`
	PageSize       int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	PermissionKey  string `form:"permission_key"`
	PermissionName string `form:"permission_name"`
	ResourcePath   string `form:"resource_path"`
	Method         string `form:"method"`
	Status         string `form:"status"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID             uint    `json:"id"`
	CreatedAt      string  `json:"createdAt,omitempty"`
	UpdatedAt      string  `json:"updatedAt,omitempty"`
	DeletedAt      *string `json:"deletedAt,omitempty"`
	PermissionID   string  `json:"permissionId"`
	PermissionKey  string  `json:"permissionKey"`
	PermissionName string  `json:"permissionName"`
	ResourcePath   string  `json:"resourcePath"`
	Method         string  `json:"method"`
	Description    string  `json:"description"`
	Status         string  `json:"status"`
	CreatorID      string  `json:"creatorId,omitempty"`
	CreatorName    string  `json:"creatorName,omitempty"`
	ModifierID     string  `json:"modifierId,omitempty"`
	ModifierName   string  `json:"modifierName,omitempty"`
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	PermissionKey  string `json:"permissionKey" binding:"required,min=2,max=100"`
	PermissionName string `json:"permissionName" binding:"required,min=1,max=100"`
	ResourcePath   string `json:"resourcePath" binding:"required"`
	Method         string `json:"method" binding:"required,oneof=GET POST PUT DELETE PATCH"`
	Description    string `json:"description,omitempty"`
	Status         string `json:"status" binding:"omitempty,oneof=1 2"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	PermissionName *string `json:"permissionName,omitempty" binding:"omitempty,min=1,max=100"`
	ResourcePath   *string `json:"resourcePath,omitempty"`
	Method         *string `json:"method,omitempty" binding:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Description    *string `json:"description,omitempty"`
	Status         *string `json:"status,omitempty" binding:"omitempty,oneof=1 2"`
}

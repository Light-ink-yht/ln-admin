package dto

// SystemConfigResponse 系统配置响应
type SystemConfigResponse struct {
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigName  string `json:"configName"`
	ConfigGroup string `json:"configGroup"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateSystemConfigRequest 创建系统配置请求
type CreateSystemConfigRequest struct {
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue" binding:"required"`
	ConfigName  string `json:"configName" binding:"required"`
	ConfigGroup string `json:"configGroup" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// UpdateSystemConfigRequest 更新系统配置请求
type UpdateSystemConfigRequest struct {
	ConfigValue string `json:"configValue"`
	ConfigName  string `json:"configName"`
	ConfigGroup string `json:"configGroup"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

package dto

// SMSTemplateResponse 短信模板响应
type SMSTemplateResponse struct {
	TemplateID        string `json:"templateId"`
	Type              string `json:"type"`
	TemplateName      string `json:"templateName"`
	TencentTemplateID string `json:"tencentTemplateId"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	CreatorID         string `json:"creatorId"`
	ModifierID        string `json:"modifierId"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

// SMSTemplateListRequest 短信模板列表请求
type SMSTemplateListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Type     string `form:"type"`
	Status   string `form:"status"`
}

// CreateSMSTemplateRequest 创建短信模板请求
type CreateSMSTemplateRequest struct {
	Type              string `json:"type" binding:"required"`
	TemplateName      string `json:"templateName" binding:"required"`
	TencentTemplateID string `json:"tencentTemplateId" binding:"required"`
	Title             string `json:"title" binding:"required"`
	Content           string `json:"content" binding:"required"`
	Description       string `json:"description"`
	Status            string `json:"status"`
}

// UpdateSMSTemplateRequest 更新短信模板请求
type UpdateSMSTemplateRequest struct {
	TemplateName      string `json:"templateName"`
	TencentTemplateID string `json:"tencentTemplateId"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Description       string `json:"description"`
	Status            string `json:"status"`
}

// SMSCodeResponse 短信验证码响应
type SMSCodeResponse struct {
	CodeID    string `json:"codeId"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	ExpireAt  string `json:"expireAt"`
	UsedAt    string `json:"usedAt"`
	IP        string `json:"ip"`
	SendCount int    `json:"sendCount"`
	CreatedAt string `json:"createdAt"`
}

// SMSCodeListRequest 短信验证码列表请求
type SMSCodeListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Phone    string `form:"phone"`
	Type     string `form:"type"`
	Status   string `form:"status"`
}

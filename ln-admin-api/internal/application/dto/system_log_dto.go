package dto

// SystemLogResponse 系统日志响应
type SystemLogResponse struct {
	LogID      string `json:"logId"`
	Level      string `json:"level"`
	Module     string `json:"module"`
	Action     string `json:"action"`
	Content    string `json:"content"`
	UserID     string `json:"userId"`
	IP         string `json:"ip"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	StatusCode int    `json:"statusCode"`
	UserAgent  string `json:"userAgent"`
	ErrorMsg   string `json:"errorMsg"`
	LogTime    string `json:"logTime"`
	CreatedAt  string `json:"createdAt"`
}

// SystemLogListRequest 系统日志列表请求
type SystemLogListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Level     string `form:"level"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	UserID    string `form:"userId"`
	IP        string `form:"ip"`
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

package dto

// WorkbenchStatItem 工作台统计项
type WorkbenchStatItem struct {
	Key    string `json:"key"`
	Title  string `json:"title"`
	Value  int    `json:"value"`
	Icon   string `json:"icon"`
	Color  string `json:"color"`
	Trend  string `json:"trend"` // "up" | "down"
	Change string `json:"change"`
	Desc   string `json:"desc"`
	Type   string `json:"type"` // "primary" | "success" | "warning" | "danger"
	Suffix string `json:"suffix,omitempty"`
}

// WorkbenchQuickAction 工作台快捷操作
type WorkbenchQuickAction struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Path  string `json:"path"`
}

// WorkbenchAnnouncement 工作台公告
type WorkbenchAnnouncement struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Author  string `json:"author"`
	Type    string `json:"type"` // "normal" | "important" | "warning"
}

// WorkbenchConfigData 工作台配置数据（JSON结构）
type WorkbenchConfigData struct {
	Stats          []WorkbenchStatItem     `json:"stats"`
	QuickActions   []WorkbenchQuickAction  `json:"quickActions"`
	Announcements  []WorkbenchAnnouncement `json:"announcements"`
	ShowSystemInfo bool                    `json:"showSystemInfo"`
}

// WorkbenchConfigResponse 工作台配置响应
type WorkbenchConfigResponse struct {
	RoleKey    string              `json:"role_key"`
	ConfigName string              `json:"config_name"`
	Config     WorkbenchConfigData `json:"config"`
}

// WorkbenchConfigCreateRequest 创建工作台配置请求
type WorkbenchConfigCreateRequest struct {
	RoleKey     string              `json:"role_key" binding:"required"`
	ConfigName  string              `json:"config_name" binding:"required"`
	Config      WorkbenchConfigData `json:"config" binding:"required"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
}

// WorkbenchConfigUpdateRequest 更新工作台配置请求
type WorkbenchConfigUpdateRequest struct {
	ConfigID    string              `json:"config_id" binding:"required"`
	ConfigName  string              `json:"config_name"`
	Config      WorkbenchConfigData `json:"config"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
}

// WorkbenchConfigListResponse 工作台配置列表响应
type WorkbenchConfigListResponse struct {
	ConfigID    string `json:"config_id"`
	RoleKey     string `json:"role_key"`
	ConfigName  string `json:"config_name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

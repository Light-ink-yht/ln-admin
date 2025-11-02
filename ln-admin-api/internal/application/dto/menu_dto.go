package dto

// MenuItem 菜单项DTO
type MenuItem struct {
	Key        string     `json:"key"`                  // 菜单唯一标识
	Title      string     `json:"title"`                // 菜单标题
	Path       string     `json:"path,omitempty"`       // 路由路径
	Icon       string     `json:"icon,omitempty"`       // 图标名称（如：UserOutlined）
	Children   []MenuItem `json:"children,omitempty"`   // 子菜单
	Permission string     `json:"permission,omitempty"` // 权限标识（用于权限过滤）
	Hidden     bool       `json:"hidden,omitempty"`     // 是否隐藏
}

// UserMenuItem 用户下拉菜单项DTO
type UserMenuItem struct {
	Key     string `json:"key"`               // 菜单唯一标识
	Label   string `json:"label"`             // 菜单标签
	Icon    string `json:"icon,omitempty"`    // 图标名称
	Divider bool   `json:"divider,omitempty"` // 是否在前方添加分隔线
}

// MenuCreateRequest 创建菜单请求DTO
type MenuCreateRequest struct {
	MenuKey     string `json:"menu_key" binding:"required"`  // 菜单标识
	Title       string `json:"title" binding:"required"`     // 菜单标题
	Path        string `json:"path"`                         // 路由路径
	Icon        string `json:"icon"`                         // 图标名称
	ParentID    string `json:"parent_id"`                    // 父菜单ID
	Sort        int    `json:"sort"`                         // 排序号
	Permission  string `json:"permission"`                   // 权限标识
	MenuType    string `json:"menu_type" binding:"required"` // 菜单类型：1 侧边栏菜单 2 用户菜单
	Description string `json:"description"`                  // 菜单描述
}

// MenuUpdateRequest 更新菜单请求DTO
type MenuUpdateRequest struct {
	Title       string `json:"title" binding:"required"` // 菜单标题
	Path        string `json:"path"`                     // 路由路径
	Icon        string `json:"icon"`                     // 图标名称
	ParentID    string `json:"parent_id"`                // 父菜单ID
	Sort        int    `json:"sort"`                     // 排序号
	Permission  string `json:"permission"`               // 权限标识
	Status      string `json:"status"`                   // 状态：1 启用 2 禁用
	Description string `json:"description"`              // 菜单描述
}

// MenuResponse 菜单响应DTO
type MenuResponse struct {
	ID          uint           `json:"id"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	MenuID      string         `json:"menu_id"`
	MenuKey     string         `json:"menu_key"`
	Title       string         `json:"title"`
	Path        string         `json:"path"`
	Icon        string         `json:"icon"`
	ParentID    string         `json:"parent_id"`
	Sort        int            `json:"sort"`
	Permission  string         `json:"permission"`
	MenuType    string         `json:"menu_type"`
	Status      string         `json:"status"`
	Description string         `json:"description"`
	CreatorID   string         `json:"creator_id"`
	ModifierID  string         `json:"modifier_id"`
	Children    []MenuResponse `json:"children,omitempty"` // 子菜单
}

// MenuListResponse 菜单列表响应DTO
type MenuListResponse struct {
	List  []MenuResponse `json:"list"`
	Total int64          `json:"total"`
}

package dto

// GrantRoleMenusRequest 为角色分配菜单请求
type GrantRoleMenusRequest struct {
	MenuIDs []string `json:"menu_ids" binding:"required"` // 菜单ID列表
}

// GetRoleMenusResponse 获取角色菜单响应
type GetRoleMenusResponse struct {
	MenuIDs []string `json:"menu_ids"` // 菜单ID列表
}

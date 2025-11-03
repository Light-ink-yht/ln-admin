package handler

import (
	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleHandler struct {
	permissionService *service.PermissionService
	userRepo          repository.UserRepository
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(permissionService *service.PermissionService, userRepo repository.UserRepository) *RoleHandler {
	return &RoleHandler{
		permissionService: permissionService,
		userRepo:          userRepo,
	}
}

// ListRoles 获取角色列表
// @Summary      获取角色列表
// @Description  分页获取角色列表，支持多条件查询
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page        query     int     false  "页码，从1开始"  default(1)
// @Param        page_size   query     int     false  "每页数量，最大100"  default(10)
// @Param        role_key    query     string  false  "角色标识（模糊查询）"
// @Param        role_name   query     string  false  "角色名称（模糊查询）"
// @Param        status      query     string  false  "状态（1启用，2禁用）"
// @Success      200         {object}  response.PageResponse{data=[]dto.RoleResponse}  "获取成功"
// @Failure      401         {object}  response.Response  "未授权"
// @Failure      403         {object}  response.Response  "权限不足"
// @Failure      500         {object}  response.Response  "服务器错误"
// @Router       /role/list [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	var req dto.RoleListRequest
	clientIP := c.ClientIP()

	// 从上下文获取用户ID（用于日志）
	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	logger.Info("开始获取角色列表",
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取角色列表失败",
			zap.String("操作", "获取角色列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	// 处理可能为数组的查询参数（多个搜索词）
	roleKeyValues := c.QueryArray("role_key")
	if len(roleKeyValues) > 0 {
		values := make([]interface{}, 0, len(roleKeyValues))
		for _, v := range roleKeyValues {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) > 0 {
			if len(values) == 1 {
				conditions["role_key"] = values[0]
			} else {
				conditions["role_key"] = values
			}
		}
	} else if req.RoleKey != "" {
		conditions["role_key"] = req.RoleKey
	}

	roleNameValues := c.QueryArray("role_name")
	if len(roleNameValues) > 0 {
		values := make([]interface{}, 0, len(roleNameValues))
		for _, v := range roleNameValues {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) > 0 {
			if len(values) == 1 {
				conditions["role_name"] = values[0]
			} else {
				conditions["role_name"] = values
			}
		}
	} else if req.RoleName != "" {
		conditions["role_name"] = req.RoleName
	}

	if req.Status != "" {
		conditions["status"] = req.Status
	}

	// 调用服务
	roles, total, err := h.permissionService.GetRoleList(c.Request.Context(), page, pageSize, conditions)
	if err != nil {
		logger.Error("获取角色列表失败",
			zap.String("操作", "获取角色列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取角色列表服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取角色列表失败")
		return
	}

	// 收集所有需要查询的用户ID（创建人和修改人）
	userIDSet := make(map[string]bool)
	for _, role := range roles {
		if role.CreatorID != "" {
			userIDSet[role.CreatorID] = true
		}
		if role.ModifierID != "" {
			userIDSet[role.ModifierID] = true
		}
	}

	// 批量查询用户信息
	userIDList := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDList = append(userIDList, id)
	}

	userMap := make(map[string]*entity.User)
	if len(userIDList) > 0 {
		users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDList)
		if err != nil {
			logger.Warn("批量查询用户信息失败", zap.Error(err))
			// 不中断流程，继续处理
		} else {
			for _, user := range users {
				userMap[user.UserID] = user
			}
		}
	}

	// 转换为响应DTO，填充用户姓名
	roleResps := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResp := h.toRoleResponse(role, userMap)
		roleResps[i] = roleResp
	}

	roleCount := len(roleResps)
	logger.Info("获取角色列表成功",
		zap.String("操作", "获取角色列表"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.Int64("total", total),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int("count", roleCount),
		zap.String("ip", clientIP))

	response.PageSuccess(c, roleResps, total, page, pageSize)
}

// GetAllRoles 获取所有角色（不分页）
// @Summary      获取所有角色
// @Description  获取所有角色列表，用于下拉选择等场景
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  response.Response{data=[]dto.RoleInfo}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /role/all [get]
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	clientIP := c.ClientIP()

	logger.Info("开始获取所有角色",
		zap.String("ip", clientIP))

	// 查询所有角色（使用大分页，通过PermissionService获取）
	roles, _, err := h.permissionService.GetAllRoles(c.Request.Context())
	if err != nil {
		logger.Error("获取所有角色失败",
			zap.String("操作", "获取所有角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "查询角色列表失败"),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取角色列表失败")
		return
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

	logger.Info("获取所有角色成功",
		zap.String("操作", "获取所有角色"),
		zap.String("结果", "成功"),
		zap.Int("count", len(roleInfos)),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取角色列表成功", roleInfos)
}

// GetRoleDetail 获取角色详情
// @Summary      获取角色详情
// @Description  根据角色ID获取角色详细信息
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        roleId  path      string  true  "角色ID"
// @Success      200     {object}  response.Response{data=dto.RoleResponse}  "获取成功"
// @Failure      401     {object}  response.Response  "未授权"
// @Failure      404     {object}  response.Response  "角色不存在"
// @Failure      500     {object}  response.Response  "服务器错误"
// @Router       /role/{roleId} [get]
func (h *RoleHandler) GetRoleDetail(c *gin.Context) {
	roleID := c.Param("roleId")
	clientIP := c.ClientIP()

	// 从上下文获取用户ID（用于日志）
	viewerID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			viewerID = uidStr
		}
	}

	logger.Info("开始获取角色详情",
		zap.String("role_id", roleID),
		zap.String("viewer_id", viewerID),
		zap.String("ip", clientIP))

	// 调用服务
	role, err := h.permissionService.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			logger.Warn("获取角色详情失败",
				zap.String("操作", "获取角色详情"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "角色不存在"),
				zap.String("role_id", roleID),
				zap.String("viewer_id", viewerID),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 404, "角色不存在")
			return
		}
		logger.Error("获取角色详情失败",
			zap.String("操作", "获取角色详情"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取角色详情服务内部错误"),
			zap.String("role_id", roleID),
			zap.String("viewer_id", viewerID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取角色详情失败")
		return
	}

	// 查询创建人和修改人信息
	userMap := make(map[string]*entity.User)
	if role.CreatorID != "" || role.ModifierID != "" {
		userIDs := make([]string, 0, 2)
		if role.CreatorID != "" {
			userIDs = append(userIDs, role.CreatorID)
		}
		if role.ModifierID != "" {
			userIDs = append(userIDs, role.ModifierID)
		}
		users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDs)
		if err == nil {
			for _, user := range users {
				userMap[user.UserID] = user
			}
		}
	}

	roleResp := h.toRoleResponse(role, userMap)

	logger.Info("获取角色详情成功",
		zap.String("操作", "获取角色详情"),
		zap.String("结果", "成功"),
		zap.String("role_id", roleID),
		zap.String("viewer_id", viewerID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取角色详情成功", roleResp)
}

// CreateRole 创建角色
// @Summary      创建角色
// @Description  创建新角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.CreateRoleRequest  true  "创建角色请求"
// @Success      200      {object}  response.Response{data=dto.RoleResponse}  "创建成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      403      {object}  response.Response  "权限不足"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /role [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	clientIP := c.ClientIP()

	// 从上下文获取当前用户ID
	creatorID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			creatorID = uidStr
		}
	}

	logger.Info("开始创建角色",
		zap.String("creator_id", creatorID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建角色失败",
			zap.String("操作", "创建角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("creator_id", creatorID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 设置默认状态
	status := req.Status
	if status == "" {
		status = "1" // 默认启用
	}

	// 转换为实体
	role := &entity.Role{
		RoleKey:     req.RoleKey,
		RoleName:    req.RoleName,
		Description: req.Description,
		Status:      status,
		CreatorID:   creatorID,
		ModifierID:  creatorID,
	}

	// 调用服务
	if err := h.permissionService.CreateRole(c.Request.Context(), role); err != nil {
		if err == service.ErrRoleExists {
			logger.Warn("创建角色失败",
				zap.String("操作", "创建角色"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "角色已存在"),
				zap.String("creator_id", creatorID),
				zap.String("role_key", req.RoleKey),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 400, "角色已存在")
			return
		}
		logger.Error("创建角色失败",
			zap.String("操作", "创建角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "创建角色服务内部错误"),
			zap.String("creator_id", creatorID),
			zap.String("role_key", req.RoleKey),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "创建角色失败")
		return
	}

	// 重新查询角色（获取完整信息，包括ID和时间戳）
	createdRole, err := h.permissionService.GetRoleByID(c.Request.Context(), role.RoleID)
	if err != nil {
		logger.Error("查询新创建的角色失败",
			zap.String("role_id", role.RoleID),
			zap.Error(err))
		// 仍然返回成功，但使用部分信息
		userMap := make(map[string]*entity.User)
		if creator, err := h.userRepo.FindByID(c.Request.Context(), creatorID); err == nil && creator != nil {
			userMap[creator.UserID] = creator
		}
		roleResp := h.toRoleResponse(role, userMap)
		response.SuccessWithMessage(c, "创建角色成功", roleResp)
		return
	}

	// 查询创建人信息
	userMap := make(map[string]*entity.User)
	if createdRole.CreatorID != "" {
		if creator, err := h.userRepo.FindByID(c.Request.Context(), createdRole.CreatorID); err == nil && creator != nil {
			userMap[creator.UserID] = creator
		}
	}
	if createdRole.ModifierID != "" {
		if modifier, err := h.userRepo.FindByID(c.Request.Context(), createdRole.ModifierID); err == nil && modifier != nil {
			userMap[modifier.UserID] = modifier
		}
	}

	roleResp := h.toRoleResponse(createdRole, userMap)

	logger.Info("创建角色成功",
		zap.String("操作", "创建角色"),
		zap.String("结果", "成功"),
		zap.String("creator_id", creatorID),
		zap.String("role_id", role.RoleID),
		zap.String("role_key", req.RoleKey),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "创建角色成功", roleResp)
}

// UpdateRole 更新角色
// @Summary      更新角色
// @Description  更新角色信息（角色标识不可修改）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        roleId   path      string               true  "角色ID"
// @Param        request  body      dto.UpdateRoleRequest  true  "更新角色请求"
// @Success      200      {object}  response.Response{data=dto.RoleResponse}  "更新成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      404      {object}  response.Response  "角色不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /role/{roleId} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("roleId")
	var req dto.UpdateRoleRequest
	clientIP := c.ClientIP()

	// 从上下文获取当前用户ID
	modifierID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			modifierID = uidStr
		}
	}

	logger.Info("开始更新角色",
		zap.String("role_id", roleID),
		zap.String("modifier_id", modifierID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新角色失败",
			zap.String("操作", "更新角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("role_id", roleID),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 先查询现有角色
	existingRole, err := h.permissionService.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			logger.Warn("更新角色失败",
				zap.String("操作", "更新角色"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "角色不存在"),
				zap.String("role_id", roleID),
				zap.String("modifier_id", modifierID),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 404, "角色不存在")
			return
		}
		logger.Error("更新角色失败",
			zap.String("操作", "更新角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "查询角色失败"),
			zap.String("role_id", roleID),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "更新角色失败")
		return
	}

	// 更新字段（只更新提供的字段）
	if req.RoleName != nil {
		existingRole.RoleName = *req.RoleName
	}
	if req.Description != nil {
		existingRole.Description = *req.Description
	}
	if req.Status != nil {
		existingRole.Status = *req.Status
	}
	existingRole.ModifierID = modifierID

	// 调用服务
	if err := h.permissionService.UpdateRole(c.Request.Context(), existingRole); err != nil {
		logger.Error("更新角色失败",
			zap.String("操作", "更新角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "更新角色服务内部错误"),
			zap.String("role_id", roleID),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "更新角色失败")
		return
	}

	// 重新查询角色（获取最新信息）
	updatedRole, err := h.permissionService.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		logger.Error("查询更新后的角色失败",
			zap.String("role_id", roleID),
			zap.Error(err))
		// 仍然返回成功，但使用部分信息
		userMap := make(map[string]*entity.User)
		if existingRole.CreatorID != "" || existingRole.ModifierID != "" {
			userIDs := make([]string, 0, 2)
			if existingRole.CreatorID != "" {
				userIDs = append(userIDs, existingRole.CreatorID)
			}
			if existingRole.ModifierID != "" {
				userIDs = append(userIDs, existingRole.ModifierID)
			}
			if users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDs); err == nil {
				for _, user := range users {
					userMap[user.UserID] = user
				}
			}
		}
		roleResp := h.toRoleResponse(existingRole, userMap)
		response.SuccessWithMessage(c, "更新角色成功", roleResp)
		return
	}

	// 查询创建人和修改人信息
	userMap := make(map[string]*entity.User)
	if updatedRole.CreatorID != "" || updatedRole.ModifierID != "" {
		userIDs := make([]string, 0, 2)
		if updatedRole.CreatorID != "" {
			userIDs = append(userIDs, updatedRole.CreatorID)
		}
		if updatedRole.ModifierID != "" {
			userIDs = append(userIDs, updatedRole.ModifierID)
		}
		if users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDs); err == nil {
			for _, user := range users {
				userMap[user.UserID] = user
			}
		}
	}

	roleResp := h.toRoleResponse(updatedRole, userMap)

	logger.Info("更新角色成功",
		zap.String("操作", "更新角色"),
		zap.String("结果", "成功"),
		zap.String("role_id", roleID),
		zap.String("modifier_id", modifierID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "更新角色成功", roleResp)
}

// DeleteRole 删除角色
// @Summary      删除角色
// @Description  根据角色ID删除角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        roleId  path      string  true  "角色ID"
// @Success      200     {object}  response.Response  "删除成功"
// @Failure      401     {object}  response.Response  "未授权"
// @Failure      404     {object}  response.Response  "角色不存在"
// @Failure      500     {object}  response.Response  "服务器错误"
// @Router       /role/{roleId} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("roleId")
	clientIP := c.ClientIP()

	// 从上下文获取当前用户ID
	deleterID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			deleterID = uidStr
		}
	}

	logger.Info("开始删除角色",
		zap.String("role_id", roleID),
		zap.String("deleter_id", deleterID),
		zap.String("ip", clientIP))

	// 调用服务
	if err := h.permissionService.DeleteRole(c.Request.Context(), roleID); err != nil {
		if err == service.ErrRoleNotFound {
			logger.Warn("删除角色失败",
				zap.String("操作", "删除角色"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "角色不存在"),
				zap.String("role_id", roleID),
				zap.String("deleter_id", deleterID),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 404, "角色不存在")
			return
		}
		logger.Error("删除角色失败",
			zap.String("操作", "删除角色"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "删除角色服务内部错误"),
			zap.String("role_id", roleID),
			zap.String("deleter_id", deleterID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "删除角色失败")
		return
	}

	logger.Info("删除角色成功",
		zap.String("操作", "删除角色"),
		zap.String("结果", "成功"),
		zap.String("role_id", roleID),
		zap.String("deleter_id", deleterID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "删除角色成功", nil)
}

// toRoleResponse 转换为角色响应DTO
func (h *RoleHandler) toRoleResponse(role *entity.Role, userMap map[string]*entity.User) dto.RoleResponse {
	var createdAt, updatedAt string
	var deletedAt *string

	if !role.CreatedAt.IsZero() {
		createdAt = role.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !role.UpdatedAt.IsZero() {
		updatedAt = role.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	if role.DeletedAt != nil && !role.DeletedAt.IsZero() {
		formatted := role.DeletedAt.Format("2006-01-02 15:04:05")
		deletedAt = &formatted
	}

	// 获取创建人姓名
	creatorName := ""
	if role.CreatorID != "" {
		if creator, exists := userMap[role.CreatorID]; exists && creator != nil {
			// 优先使用 fullName，其次使用 nickname，最后使用 userID
			if creator.FullName != "" {
				creatorName = creator.FullName
			} else if creator.Nickname != "" {
				creatorName = creator.Nickname
			} else {
				creatorName = role.CreatorID
			}
		} else {
			creatorName = role.CreatorID
		}
	}

	// 获取修改人姓名
	modifierName := ""
	if role.ModifierID != "" {
		if modifier, exists := userMap[role.ModifierID]; exists && modifier != nil {
			// 优先使用 fullName，其次使用 nickname，最后使用 userID
			if modifier.FullName != "" {
				modifierName = modifier.FullName
			} else if modifier.Nickname != "" {
				modifierName = modifier.Nickname
			} else {
				modifierName = role.ModifierID
			}
		} else {
			modifierName = role.ModifierID
		}
	}

	return dto.RoleResponse{
		ID:           role.ID,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		DeletedAt:    deletedAt,
		RoleID:       role.RoleID,
		RoleKey:      role.RoleKey,
		RoleName:     role.RoleName,
		Description:  role.Description,
		Status:       role.Status,
		CreatorID:    role.CreatorID,
		CreatorName:  creatorName,
		ModifierID:   role.ModifierID,
		ModifierName: modifierName,
	}
}

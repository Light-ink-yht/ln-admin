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

type PermissionHandler struct {
	permissionService *service.PermissionService
	userRepo          repository.UserRepository
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(permissionService *service.PermissionService, userRepo repository.UserRepository) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
		userRepo:          userRepo,
	}
}

// ListPermissions 获取权限列表
// @Summary      获取权限列表
// @Description  分页获取权限列表，支持多条件查询
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page           query     int     false  "页码，从1开始"  default(1)
// @Param        page_size      query     int     false  "每页数量，最大100"  default(10)
// @Param        permission_key query     string  false  "权限标识（模糊查询）"
// @Param        permission_name query   string  false  "权限名称（模糊查询）"
// @Param        resource_path  query     string  false  "资源路径（模糊查询）"
// @Param        method         query     string  false  "请求方法（GET, POST, PUT, DELETE）"
// @Param        status         query     string  false  "状态（1启用，2禁用）"
// @Success      200            {object}  response.PageResponse{data=[]dto.PermissionResponse}  "获取成功"
// @Failure      401            {object}  response.Response  "未授权"
// @Failure      403            {object}  response.Response  "权限不足"
// @Failure      500            {object}  response.Response  "服务器错误"
// @Router       /permission/list [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	var req dto.PermissionListRequest
	clientIP := c.ClientIP()

	// 从上下文获取用户ID（用于日志）
	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	logger.Info("开始获取权限列表",
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取权限列表失败",
			zap.String("操作", "获取权限列表"),
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
	permissionKeyValues := c.QueryArray("permission_key")
	if len(permissionKeyValues) > 0 {
		values := make([]interface{}, 0, len(permissionKeyValues))
		for _, v := range permissionKeyValues {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) > 0 {
			if len(values) == 1 {
				conditions["permission_key"] = values[0]
			} else {
				conditions["permission_key"] = values
			}
		}
	} else if req.PermissionKey != "" {
		conditions["permission_key"] = req.PermissionKey
	}

	permissionNameValues := c.QueryArray("permission_name")
	if len(permissionNameValues) > 0 {
		values := make([]interface{}, 0, len(permissionNameValues))
		for _, v := range permissionNameValues {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) > 0 {
			if len(values) == 1 {
				conditions["permission_name"] = values[0]
			} else {
				conditions["permission_name"] = values
			}
		}
	} else if req.PermissionName != "" {
		conditions["permission_name"] = req.PermissionName
	}

	resourcePathValues := c.QueryArray("resource_path")
	if len(resourcePathValues) > 0 {
		values := make([]interface{}, 0, len(resourcePathValues))
		for _, v := range resourcePathValues {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) > 0 {
			if len(values) == 1 {
				conditions["resource_path"] = values[0]
			} else {
				conditions["resource_path"] = values
			}
		}
	} else if req.ResourcePath != "" {
		conditions["resource_path"] = req.ResourcePath
	}

	if req.Method != "" {
		conditions["method"] = req.Method
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}
	if req.Category != "" {
		conditions["category"] = req.Category
	}

	// 调用服务
	permissions, total, err := h.permissionService.GetPermissionList(c.Request.Context(), page, pageSize, conditions)
	if err != nil {
		logger.Error("获取权限列表失败",
			zap.String("操作", "获取权限列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取权限列表服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取权限列表失败")
		return
	}

	// 收集所有需要查询的用户ID（创建人和修改人）
	userIDSet := make(map[string]bool)
	for _, permission := range permissions {
		if permission.CreatorID != "" {
			userIDSet[permission.CreatorID] = true
		}
		if permission.ModifierID != "" {
			userIDSet[permission.ModifierID] = true
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
	permissionResps := make([]dto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResps[i] = h.toPermissionResponse(permission, userMap)
	}

	permissionCount := len(permissionResps)
	logger.Info("获取权限列表成功",
		zap.String("操作", "获取权限列表"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.Int64("total", total),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int("count", permissionCount),
		zap.String("ip", clientIP))

	response.PageSuccess(c, permissionResps, total, page, pageSize)
}

// GetPermissionDetail 获取权限详情
// @Summary      获取权限详情
// @Description  根据权限ID获取权限详细信息
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        permissionId  path      string  true  "权限ID"
// @Success      200           {object}  response.Response{data=dto.PermissionResponse}  "获取成功"
// @Failure      401           {object}  response.Response  "未授权"
// @Failure      404           {object}  response.Response  "权限不存在"
// @Failure      500           {object}  response.Response  "服务器错误"
// @Router       /permission/{permissionId} [get]
func (h *PermissionHandler) GetPermissionDetail(c *gin.Context) {
	permissionID := c.Param("permissionId")
	clientIP := c.ClientIP()

	// 从上下文获取用户ID（用于日志）
	viewerID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			viewerID = uidStr
		}
	}

	logger.Info("开始获取权限详情",
		zap.String("permission_id", permissionID),
		zap.String("viewer_id", viewerID),
		zap.String("ip", clientIP))

	// 调用服务
	permission, err := h.permissionService.GetPermissionByID(c.Request.Context(), permissionID)
	if err != nil {
		if err == service.ErrPermissionNotFound {
			logger.Warn("获取权限详情失败",
				zap.String("操作", "获取权限详情"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "权限不存在"),
				zap.String("permission_id", permissionID),
				zap.String("viewer_id", viewerID),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 404, "权限不存在")
			return
		}
		logger.Error("获取权限详情失败",
			zap.String("操作", "获取权限详情"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取权限详情服务内部错误"),
			zap.String("permission_id", permissionID),
			zap.String("viewer_id", viewerID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取权限详情失败")
		return
	}

	// 查询创建人和修改人信息
	userMap := make(map[string]*entity.User)
	if permission.CreatorID != "" || permission.ModifierID != "" {
		userIDs := make([]string, 0, 2)
		if permission.CreatorID != "" {
			userIDs = append(userIDs, permission.CreatorID)
		}
		if permission.ModifierID != "" {
			userIDs = append(userIDs, permission.ModifierID)
		}
		users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDs)
		if err == nil {
			for _, user := range users {
				userMap[user.UserID] = user
			}
		}
	}

	permissionResp := h.toPermissionResponse(permission, userMap)

	logger.Info("获取权限详情成功",
		zap.String("操作", "获取权限详情"),
		zap.String("结果", "成功"),
		zap.String("permission_id", permissionID),
		zap.String("viewer_id", viewerID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取权限详情成功", permissionResp)
}

// CreatePermission 创建权限
// @Summary      创建权限
// @Description  创建新权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      dto.CreatePermissionRequest  true  "创建权限请求"
// @Success      200      {object}  response.Response{data=dto.PermissionResponse}  "创建成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      403      {object}  response.Response  "权限不足"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /permission [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	clientIP := c.ClientIP()

	// 从上下文获取当前用户ID
	creatorID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			creatorID = uidStr
		}
	}

	logger.Info("开始创建权限",
		zap.String("creator_id", creatorID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建权限失败",
			zap.String("操作", "创建权限"),
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
	permission := &entity.Permission{
		PermissionKey:  req.PermissionKey,
		PermissionName: req.PermissionName,
		ResourcePath:   req.ResourcePath,
		Method:         req.Method,
		Description:    req.Description,
		Status:         status,
		Category:       req.Category,
		CreatorID:      creatorID,
		ModifierID:     creatorID,
	}

	// 调用服务
	if err := h.permissionService.CreatePermission(c.Request.Context(), permission); err != nil {
		if err == service.ErrPermissionExists {
			logger.Warn("创建权限失败",
				zap.String("操作", "创建权限"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "权限已存在"),
				zap.String("creator_id", creatorID),
				zap.String("permission_key", req.PermissionKey),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 400, "权限已存在")
			return
		}
		logger.Error("创建权限失败",
			zap.String("操作", "创建权限"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "创建权限服务内部错误"),
			zap.String("creator_id", creatorID),
			zap.String("permission_key", req.PermissionKey),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "创建权限失败")
		return
	}

	// 重新查询权限（获取完整信息，包括ID和时间戳）
	createdPermission, err := h.permissionService.GetPermissionByID(c.Request.Context(), permission.PermissionID)
	if err != nil {
		logger.Error("查询新创建的权限失败",
			zap.String("permission_id", permission.PermissionID),
			zap.Error(err))
		// 仍然返回成功，但使用部分信息
		userMap := make(map[string]*entity.User)
		if creator, err := h.userRepo.FindByID(c.Request.Context(), creatorID); err == nil && creator != nil {
			userMap[creator.UserID] = creator
		}
		permissionResp := h.toPermissionResponse(permission, userMap)
		response.SuccessWithMessage(c, "创建权限成功", permissionResp)
		return
	}

	// 查询创建人信息
	userMap := make(map[string]*entity.User)
	if createdPermission.CreatorID != "" {
		if creator, err := h.userRepo.FindByID(c.Request.Context(), createdPermission.CreatorID); err == nil && creator != nil {
			userMap[creator.UserID] = creator
		}
	}
	if createdPermission.ModifierID != "" {
		if modifier, err := h.userRepo.FindByID(c.Request.Context(), createdPermission.ModifierID); err == nil && modifier != nil {
			userMap[modifier.UserID] = modifier
		}
	}

	permissionResp := h.toPermissionResponse(createdPermission, userMap)

	logger.Info("创建权限成功",
		zap.String("操作", "创建权限"),
		zap.String("结果", "成功"),
		zap.String("creator_id", creatorID),
		zap.String("permission_id", permission.PermissionID),
		zap.String("permission_key", req.PermissionKey),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "创建权限成功", permissionResp)
}

// UpdatePermission 更新权限
// @Summary      更新权限
// @Description  更新权限信息（权限标识不可修改）
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        permissionId  path      string                  true  "权限ID"
// @Param        request      body      dto.UpdatePermissionRequest  true  "更新权限请求"
// @Success      200          {object}  response.Response{data=dto.PermissionResponse}  "更新成功"
// @Failure      400          {object}  response.Response  "请求参数错误"
// @Failure      401          {object}  response.Response  "未授权"
// @Failure      404          {object}  response.Response  "权限不存在"
// @Failure      500          {object}  response.Response  "服务器错误"
// @Router       /permission/{permissionId} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	permissionID := c.Param("permissionId")
	var req dto.UpdatePermissionRequest
	clientIP := c.ClientIP()

	// 从上下文获取当前用户ID
	modifierID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			modifierID = uidStr
		}
	}

	logger.Info("开始更新权限",
		zap.String("permission_id", permissionID),
		zap.String("modifier_id", modifierID),
		zap.String("ip", clientIP))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("更新权限失败",
			zap.String("操作", "更新权限"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("permission_id", permissionID),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 先查询现有权限
	existingPermission, err := h.permissionService.GetPermissionByID(c.Request.Context(), permissionID)
	if err != nil {
		if err == service.ErrPermissionNotFound {
			logger.Warn("更新权限失败",
				zap.String("操作", "更新权限"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "权限不存在"),
				zap.String("permission_id", permissionID),
				zap.String("modifier_id", modifierID),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 404, "权限不存在")
			return
		}
		logger.Error("更新权限失败",
			zap.String("操作", "更新权限"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "查询权限失败"),
			zap.String("permission_id", permissionID),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "更新权限失败")
		return
	}

	// 更新字段（只更新提供的字段）
	if req.PermissionName != nil {
		existingPermission.PermissionName = *req.PermissionName
	}
	if req.ResourcePath != nil {
		existingPermission.ResourcePath = *req.ResourcePath
	}
	if req.Method != nil {
		existingPermission.Method = *req.Method
	}
	if req.Description != nil {
		existingPermission.Description = *req.Description
	}
	if req.Status != nil {
		existingPermission.Status = *req.Status
	}
	if req.Category != nil {
		existingPermission.Category = *req.Category
	}
	existingPermission.ModifierID = modifierID

	// 调用服务
	if err := h.permissionService.UpdatePermission(c.Request.Context(), existingPermission); err != nil {
		logger.Error("更新权限失败",
			zap.String("操作", "更新权限"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "更新权限服务内部错误"),
			zap.String("permission_id", permissionID),
			zap.String("modifier_id", modifierID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "更新权限失败")
		return
	}

	// 重新查询权限（获取最新信息）
	updatedPermission, err := h.permissionService.GetPermissionByID(c.Request.Context(), permissionID)
	if err != nil {
		logger.Error("查询更新后的权限失败",
			zap.String("permission_id", permissionID),
			zap.Error(err))
		// 仍然返回成功，但使用部分信息
		userMap := make(map[string]*entity.User)
		if existingPermission.CreatorID != "" || existingPermission.ModifierID != "" {
			userIDs := make([]string, 0, 2)
			if existingPermission.CreatorID != "" {
				userIDs = append(userIDs, existingPermission.CreatorID)
			}
			if existingPermission.ModifierID != "" {
				userIDs = append(userIDs, existingPermission.ModifierID)
			}
			if users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDs); err == nil {
				for _, user := range users {
					userMap[user.UserID] = user
				}
			}
		}
		permissionResp := h.toPermissionResponse(existingPermission, userMap)
		response.SuccessWithMessage(c, "更新权限成功", permissionResp)
		return
	}

	// 查询创建人和修改人信息
	userMap := make(map[string]*entity.User)
	if updatedPermission.CreatorID != "" || updatedPermission.ModifierID != "" {
		userIDs := make([]string, 0, 2)
		if updatedPermission.CreatorID != "" {
			userIDs = append(userIDs, updatedPermission.CreatorID)
		}
		if updatedPermission.ModifierID != "" {
			userIDs = append(userIDs, updatedPermission.ModifierID)
		}
		if users, err := h.userRepo.FindByIDs(c.Request.Context(), userIDs); err == nil {
			for _, user := range users {
				userMap[user.UserID] = user
			}
		}
	}

	permissionResp := h.toPermissionResponse(updatedPermission, userMap)

	logger.Info("更新权限成功",
		zap.String("操作", "更新权限"),
		zap.String("结果", "成功"),
		zap.String("permission_id", permissionID),
		zap.String("modifier_id", modifierID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "更新权限成功", permissionResp)
}

// DeletePermission 删除权限
// @Summary      删除权限
// @Description  根据权限ID删除权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        permissionId  path      string  true  "权限ID"
// @Success      200          {object}  response.Response  "删除成功"
// @Failure      401          {object}  response.Response  "未授权"
// @Failure      404          {object}  response.Response  "权限不存在"
// @Failure      500          {object}  response.Response  "服务器错误"
// @Router       /permission/{permissionId} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	permissionID := c.Param("permissionId")
	clientIP := c.ClientIP()

	// 从上下文获取当前用户ID
	deleterID := "system"
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			deleterID = uidStr
		}
	}

	logger.Info("开始删除权限",
		zap.String("permission_id", permissionID),
		zap.String("deleter_id", deleterID),
		zap.String("ip", clientIP))

	// 调用服务
	if err := h.permissionService.DeletePermission(c.Request.Context(), permissionID); err != nil {
		if err == service.ErrPermissionNotFound {
			logger.Warn("删除权限失败",
				zap.String("操作", "删除权限"),
				zap.String("结果", "失败"),
				zap.String("失败原因", "权限不存在"),
				zap.String("permission_id", permissionID),
				zap.String("deleter_id", deleterID),
				zap.String("ip", clientIP),
				zap.Error(err))
			response.Error(c, 404, "权限不存在")
			return
		}
		logger.Error("删除权限失败",
			zap.String("操作", "删除权限"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "删除权限服务内部错误"),
			zap.String("permission_id", permissionID),
			zap.String("deleter_id", deleterID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "删除权限失败")
		return
	}

	logger.Info("删除权限成功",
		zap.String("操作", "删除权限"),
		zap.String("结果", "成功"),
		zap.String("permission_id", permissionID),
		zap.String("deleter_id", deleterID),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "删除权限成功", nil)
}

// GetAllPermissions 获取所有权限（不分页，根据层级授权过滤）
// @Summary      获取所有权限
// @Description  获取当前用户可见的权限列表，不分页，用于下拉选择等场景。超级管理员可以看到所有权限，其他用户只能看到被授予的权限。
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200      {object}  response.Response{data=[]dto.PermissionResponse}  "获取成功"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /permission/all [get]
func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	clientIP := c.ClientIP()

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	logger.Info("开始获取所有权限",
		zap.String("user_id", userID),
		zap.String("ip", clientIP))

	// 调用服务，获取用户可见的权限（根据层级授权）
	permissions, err := h.permissionService.GetVisiblePermissions(c.Request.Context(), userID)
	if err != nil {
		logger.Error("获取所有权限失败",
			zap.String("操作", "获取所有权限"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取所有权限服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取所有权限失败")
		return
	}

	// 转换为响应DTO（不需要查询创建人和修改人，简化处理）
	permissionResps := make([]*dto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResps[i] = &dto.PermissionResponse{
			PermissionID:   permission.PermissionID,
			PermissionKey:  permission.PermissionKey,
			PermissionName: permission.PermissionName,
			ResourcePath:   permission.ResourcePath,
			Method:         permission.Method,
			Description:    permission.Description,
			Status:         permission.Status,
			CreatorID:      permission.CreatorID,
			ModifierID:     permission.ModifierID,
			CreatedAt:      dto.FormatTime(permission.CreatedAt),
			UpdatedAt:      dto.FormatTime(permission.UpdatedAt),
		}
	}

	permissionCount := 0
	if permissionResps != nil {
		permissionCount = len(permissionResps)
	}

	logger.Info("获取所有权限成功",
		zap.String("操作", "获取所有权限"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.Int("count", permissionCount),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取所有权限成功", permissionResps)
}

// toPermissionResponse 转换为权限响应DTO
func (h *PermissionHandler) toPermissionResponse(permission *entity.Permission, userMap map[string]*entity.User) dto.PermissionResponse {
	var createdAt, updatedAt string
	var deletedAt *string

	if !permission.CreatedAt.IsZero() {
		createdAt = permission.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !permission.UpdatedAt.IsZero() {
		updatedAt = permission.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	if permission.DeletedAt != nil && !permission.DeletedAt.IsZero() {
		formatted := permission.DeletedAt.Format("2006-01-02 15:04:05")
		deletedAt = &formatted
	}

	// 获取创建人姓名
	creatorName := ""
	if permission.CreatorID != "" {
		if creator, exists := userMap[permission.CreatorID]; exists && creator != nil {
			// 优先使用 fullName，其次使用 nickname，最后使用 userID
			if creator.FullName != "" {
				creatorName = creator.FullName
			} else if creator.Nickname != "" {
				creatorName = creator.Nickname
			} else {
				creatorName = permission.CreatorID
			}
		} else {
			creatorName = permission.CreatorID
		}
	}

	// 获取修改人姓名
	modifierName := ""
	if permission.ModifierID != "" {
		if modifier, exists := userMap[permission.ModifierID]; exists && modifier != nil {
			// 优先使用 fullName，其次使用 nickname，最后使用 userID
			if modifier.FullName != "" {
				modifierName = modifier.FullName
			} else if modifier.Nickname != "" {
				modifierName = modifier.Nickname
			} else {
				modifierName = permission.ModifierID
			}
		} else {
			modifierName = permission.ModifierID
		}
	}

	return dto.PermissionResponse{
		ID:             permission.ID,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DeletedAt:      deletedAt,
		PermissionID:   permission.PermissionID,
		PermissionKey:  permission.PermissionKey,
		PermissionName: permission.PermissionName,
		ResourcePath:   permission.ResourcePath,
		Method:         permission.Method,
		Description:    permission.Description,
		Status:         permission.Status,
		Category:       permission.Category,
		CreatorID:      permission.CreatorID,
		CreatorName:    creatorName,
		ModifierID:     permission.ModifierID,
		ModifierName:   modifierName,
	}
}

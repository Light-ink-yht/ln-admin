package handler

import (
	"io"
	"strconv"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// UploadFile 上传文件
// @Summary      上传文件
// @Description  支持本地上传和S3上传（根据配置）
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "文件"
// @Success      200   {object}  response.Response{data=dto.FileResponse}  "上传成功"
// @Failure      400   {object}  response.Response  "请求参数错误"
// @Failure      500   {object}  response.Response  "服务器错误"
// @Router       /file/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取用户ID（从JWT token中）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	creatorID := userID.(string)

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("上传文件失败",
			zap.String("操作", "上传文件"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取文件失败"),
			zap.String("user_id", creatorID),
			zap.Error(err))
		response.BadRequest(c, "获取文件失败: "+err.Error())
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		logger.Error("打开文件失败",
			zap.String("操作", "上传文件"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "打开文件失败"),
			zap.String("user_id", creatorID),
			zap.Error(err))
		response.BadRequest(c, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	// 上传文件
	fileResponse, err := h.fileService.UploadFile(c.Request.Context(), src, file.Filename, file.Size, creatorID)
	if err != nil {
		logger.Error("上传文件失败",
			zap.String("操作", "上传文件"),
			zap.String("结果", "失败"),
			zap.String("失败原因", err.Error()),
			zap.String("user_id", creatorID),
			zap.String("file_name", file.Filename),
			zap.Error(err))
		response.InternalError(c, "上传文件失败: "+err.Error())
		return
	}

	logger.Info("上传文件成功",
		zap.String("操作", "上传文件"),
		zap.String("结果", "成功"),
		zap.String("user_id", creatorID),
		zap.String("file_id", fileResponse.FileID),
		zap.String("file_name", fileResponse.OriginalName))

	response.SuccessWithMessage(c, "上传文件成功", fileResponse)
}

// DeleteFile 删除文件
// @Summary      删除文件
// @Description  删除文件记录和存储中的文件
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        file_id  path  string  true  "文件ID"
// @Success      200      {object}  response.Response  "删除成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      404      {object}  response.Response  "文件不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /file/:file_id [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		response.BadRequest(c, "文件ID不能为空")
		return
	}

	err := h.fileService.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		if err == service.ErrFileNotFound {
			response.NotFound(c, "文件不存在")
			return
		}
		logger.Error("删除文件失败",
			zap.String("操作", "删除文件"),
			zap.String("结果", "失败"),
			zap.String("file_id", fileID),
			zap.Error(err))
		response.InternalError(c, "删除文件失败: "+err.Error())
		return
	}

	logger.Info("删除文件成功",
		zap.String("操作", "删除文件"),
		zap.String("结果", "成功"),
		zap.String("file_id", fileID))

	response.SuccessWithMessage(c, "删除文件成功", nil)
}

// GetFile 获取文件信息
// @Summary      获取文件信息
// @Description  根据文件ID获取文件详细信息
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        file_id  path  string  true  "文件ID"
// @Success      200      {object}  response.Response{data=dto.FileResponse}  "成功"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      404      {object}  response.Response  "文件不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /file/:file_id [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		response.BadRequest(c, "文件ID不能为空")
		return
	}

	fileResponse, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if err == service.ErrFileNotFound {
			response.NotFound(c, "文件不存在")
			return
		}
		logger.Error("获取文件信息失败",
			zap.String("操作", "获取文件信息"),
			zap.String("结果", "失败"),
			zap.String("file_id", fileID),
			zap.Error(err))
		response.InternalError(c, "获取文件信息失败: "+err.Error())
		return
	}

	response.Success(c, fileResponse)
}

// DownloadFile 下载文件
// @Summary      下载文件
// @Description  根据文件ID下载文件
// @Tags         文件管理
// @Accept       json
// @Produce      application/octet-stream
// @Param        file_id  path  string  true  "文件ID"
// @Success      200      {file}  binary  "文件内容"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      404      {object}  response.Response  "文件不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /file/:file_id/download [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		response.BadRequest(c, "文件ID不能为空")
		return
	}

	// 获取文件信息
	fileResponse, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if err == service.ErrFileNotFound {
			response.NotFound(c, "文件不存在")
			return
		}
		response.InternalError(c, "获取文件信息失败: "+err.Error())
		return
	}

	// 获取文件内容
	reader, mimeType, err := h.fileService.GetFileContent(c.Request.Context(), fileID)
	if err != nil {
		response.InternalError(c, "获取文件内容失败: "+err.Error())
		return
	}
	defer reader.Close()

	// 设置响应头
	c.Header("Content-Type", mimeType)
	c.Header("Content-Disposition", `attachment; filename="`+fileResponse.OriginalName+`"`)
	c.Header("Content-Length", strconv.FormatInt(fileResponse.FileSize, 10))

	// 复制文件内容到响应
	if _, err := io.Copy(c.Writer, reader); err != nil {
		logger.Error("下载文件失败",
			zap.String("操作", "下载文件"),
			zap.String("结果", "失败"),
			zap.String("file_id", fileID),
			zap.Error(err))
		return
	}
}

// ListFiles 获取文件列表
// @Summary      获取文件列表
// @Description  分页获取文件列表，支持按文件名、存储类型、分类等筛选
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        page        query  int     false  "页码"        default(1)
// @Param        page_size   query  int     false  "每页数量"    default(10)
// @Param        file_name   query  string  false  "文件名（模糊搜索）"
// @Param        storage_type query  string  false  "存储类型（local/s3）"
// @Param        category    query  string  false  "文件分类（image/document/video/audio/archive/other）"
// @Param        extension   query  string  false  "文件扩展名"
// @Success      200         {object}  response.Response{data=[]dto.FileResponse}  "成功"
// @Failure      400         {object}  response.Response  "请求参数错误"
// @Failure      500         {object}  response.Response  "服务器错误"
// @Router       /file/list [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	var params dto.FileListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // 限制最大每页数量
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if params.FileName != "" {
		conditions["file_name"] = params.FileName
	}
	if params.StorageType != "" {
		conditions["storage_type"] = params.StorageType
	}
	if params.Category != "" {
		conditions["category"] = params.Category
	}
	if params.Extension != "" {
		conditions["extension"] = params.Extension
	}

	// 获取当前用户ID（可选，用于筛选用户上传的文件）
	if params.CreatorID != "" {
		conditions["creator_id"] = params.CreatorID
	}

	// 查询文件列表
	files, total, err := h.fileService.ListFiles(c.Request.Context(), params.Page, params.PageSize, conditions)
	if err != nil {
		logger.Error("获取文件列表失败",
			zap.String("操作", "获取文件列表"),
			zap.String("结果", "失败"),
			zap.Error(err))
		response.InternalError(c, "获取文件列表失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "获取文件列表成功", gin.H{
		"list":      files,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
	})
}

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"msg"`
	Data     interface{} `json:"data,omitempty"`
	Total    int64       `json:"total,omitempty"`
	Page     int         `json:"page,omitempty"`
	PageSize int         `json:"page_size,omitempty"`
}

// 响应状态码
const (
	CodeSuccess       = 200 // 成功：请求处理成功
	CodeBadRequest    = 400 // 错误请求：请求参数错误或格式不正确
	CodeUnauthorized  = 401 // 未授权：未提供有效的身份认证信息
	CodeForbidden     = 403 // 禁止访问：已认证但无权限执行该操作
	CodeNotFound      = 404 // 资源不存在：请求的资源未找到
	CodeInternalError = 500 // 服务器内部错误：服务器处理请求时发生意外错误
)

// 响应消息
const (
	MsgSuccess       = "success"
	MsgBadRequest    = "bad request"
	MsgUnauthorized  = "unauthorized"
	MsgForbidden     = "forbidden"
	MsgNotFound      = "not found"
	MsgInternalError = "internal server error"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// BadRequest 错误请求响应
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = MsgBadRequest
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MsgUnauthorized
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = MsgForbidden
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
	})
}

// NotFound 未找到响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MsgNotFound
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// InternalError 内部错误响应
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = MsgInternalError
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeInternalError,
		Message: message,
	})
}

// Error 自定义错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:     CodeSuccess,
		Message:  MsgSuccess,
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

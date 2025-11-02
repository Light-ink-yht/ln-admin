package handler

import (
	"fmt"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/Light-ink-yht/ln-admin/internal/application/service"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/logger"
	"github.com/Light-ink-yht/ln-admin/pkg/config"
	"github.com/Light-ink-yht/ln-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService   *service.UserAppService
	smsAppService *service.SMSAppService
	captchaStore  base64Captcha.Store
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserAppService, smsAppService *service.SMSAppService) *UserHandler {
	// 初始化验证码存储（使用内存存储，生产环境可改用Redis）
	store := base64Captcha.DefaultMemStore

	return &UserHandler{
		userService:   userService,
		smsAppService: smsAppService,
		captchaStore:  store,
	}
}

// GetCaptcha 获取图片验证码
// @Summary      获取图片验证码
// @Description  获取用于登录/注册的图片验证码，返回验证码ID和Base64图片
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.CaptchaResponse}  "成功"
// @Failure      500  {object}  response.Response  "服务器错误"
// @Router       /user/captcha [get]
func (h *UserHandler) GetCaptcha(c *gin.Context) {
	// 记录操作开始
	clientIP := c.ClientIP()
	logger.Info("开始获取图形验证码",
		zap.String("ip", clientIP),
		zap.String("user_agent", c.Request.UserAgent()))

	// 创建验证码驱动
	// 参数：height, width, length, maxSkew(最大倾斜角度0-1，越小越清晰), dotCount(噪声点数量，越小越清晰)
	driver := base64Captcha.NewDriverDigit(
		config.Cfg.Captcha.Height,
		config.Cfg.Captcha.Width,
		config.Cfg.Captcha.Length,
		0.3, // 降低倾斜角度提高清晰度（从0.7降到0.3）
		30,  // 降低噪声点数量提高清晰度（从80降到30）
	)

	// 生成验证码
	captcha := base64Captcha.NewCaptcha(driver, h.captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		logger.Error("获取图形验证码失败",
			zap.String("操作", "生成验证码"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "生成验证码时发生错误"),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "生成验证码失败")
		return
	}

	logger.Info("获取图形验证码成功",
		zap.String("操作", "生成验证码"),
		zap.String("结果", "成功"),
		zap.String("captcha_id", id),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "图形验证码获取成功", dto.CaptchaResponse{
		CaptchaID:   id,
		CaptchaCode: b64s,
	})

}

// Login 用户登录
// @Summary      用户登录
// @Description  使用手机号和密码登录，需要先获取图片验证码
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "登录请求"
// @Success      200      {object}  response.Response{data=dto.LoginResponse}  "登录成功，返回双Token和用户信息"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      401      {object}  response.Response  "认证失败"
// @Failure      403      {object}  response.Response  "用户已禁用"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	loginIP := c.ClientIP()

	// 记录登录尝试开始
	logger.Info("开始用户登录",
		zap.String("phone", req.Phone),
		zap.String("ip", loginIP),
		zap.String("user_agent", c.Request.UserAgent()))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("用户登录失败",
			zap.String("操作", "用户登录"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("phone", req.Phone),
			zap.String("ip", loginIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 验证图片验证码（用于防止机器人）
	if !h.captchaStore.Verify(req.CaptchaID, req.CaptchaCode, true) {
		logger.Warn("用户登录失败",
			zap.String("操作", "用户登录"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "图片验证码错误或已过期"),
			zap.String("phone", req.Phone),
			zap.String("captcha_id", req.CaptchaID),
			zap.String("ip", loginIP))
		response.Error(c, 400, "图片验证码错误")
		return
	}

	// 调用服务（登录不需要短信验证码）
	loginResp, err := h.userService.Login(c.Request.Context(), req.Phone, req.Password, loginIP)
	if err != nil {
		var failureReason string
		if err == service.ErrUserNotFound {
			failureReason = "用户不存在：手机号未注册"
		} else if err == service.ErrUserDisabled {
			failureReason = "账户已被禁用"
		} else if err == service.ErrPasswordIncorrect {
			failureReason = "密码错误"
		} else {
			failureReason = "登录服务内部错误"
		}

		logger.Warn("用户登录失败",
			zap.String("操作", "用户登录"),
			zap.String("结果", "失败"),
			zap.String("失败原因", failureReason),
			zap.String("phone", req.Phone),
			zap.String("ip", loginIP),
			zap.Error(err))

		if err == service.ErrUserNotFound {
			response.Error(c, 400, "用户不存在")
			return
		}
		if err == service.ErrUserDisabled {
			response.Error(c, 403, "用户已禁用")
			return
		}
		if err == service.ErrPasswordIncorrect {
			response.Error(c, 400, "密码错误")
			return
		}
		response.InternalError(c, "登录失败")
		return
	}

	logger.Info("用户登录成功",
		zap.String("操作", "用户登录"),
		zap.String("结果", "成功"),
		zap.String("phone", req.Phone),
		zap.String("user_id", loginResp.User.UserId),
		zap.String("ip", loginIP))

	response.SuccessWithMessage(c, "登录成功，欢迎回来", loginResp)
}

// Register 用户注册
// @Summary      用户注册
// @Description  使用手机号、密码和短信验证码注册新用户
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "注册请求"
// @Success      200      {object}  response.Response{data=dto.UserResponse}  "注册成功"
// @Failure      400      {object}  response.Response  "请求参数错误或验证码错误"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/signup [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	clientIP := c.ClientIP()

	// 记录注册尝试开始
	logger.Info("开始用户注册",
		zap.String("phone", req.Phone),
		zap.String("ip", clientIP),
		zap.String("user_agent", c.Request.UserAgent()))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("用户注册失败",
			zap.String("操作", "用户注册"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("phone", req.Phone),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	userResp, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		var failureReason string
		if err == service.ErrUserExists {
			failureReason = "用户已存在：手机号已被注册"
		} else if err == service.ErrSmsCodeInvalid {
			failureReason = "短信验证码错误"
		} else if err == service.ErrSmsCodeExpired {
			failureReason = "短信验证码已过期"
		} else if err == service.ErrSmsCodeUsed {
			failureReason = "短信验证码已使用"
		} else {
			failureReason = "注册服务内部错误"
		}

		logger.Warn("用户注册失败",
			zap.String("操作", "用户注册"),
			zap.String("结果", "失败"),
			zap.String("失败原因", failureReason),
			zap.String("phone", req.Phone),
			zap.String("ip", clientIP),
			zap.Error(err))

		if err == service.ErrUserExists {
			response.Error(c, 400, "用户已存在")
			return
		}
		if err == service.ErrSmsCodeInvalid || err == service.ErrSmsCodeExpired || err == service.ErrSmsCodeUsed {
			response.Error(c, 400, err.Error())
			return
		}
		response.InternalError(c, "注册失败")
		return
	}

	logger.Info("用户注册成功",
		zap.String("操作", "用户注册"),
		zap.String("结果", "成功"),
		zap.String("phone", req.Phone),
		zap.String("user_id", userResp.UserId),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "注册成功，欢迎加入", userResp)
}

// ForgotPassword 忘记密码
// @Summary      忘记密码
// @Description  通过手机号和短信验证码重置密码
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ForgotPasswordRequest  true  "忘记密码请求"
// @Success      200      {object}  response.Response  "密码重置成功"
// @Failure      400      {object}  response.Response  "请求参数错误或验证码错误"
// @Failure      404      {object}  response.Response  "用户不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	clientIP := c.ClientIP()

	// 记录重置密码尝试开始
	logger.Info("开始重置密码",
		zap.String("phone", req.Phone),
		zap.String("ip", clientIP),
		zap.String("user_agent", c.Request.UserAgent()))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("重置密码失败",
			zap.String("操作", "重置密码"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("phone", req.Phone),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	if err := h.userService.ForgotPassword(c.Request.Context(), &req); err != nil {
		var failureReason string
		if err == service.ErrUserNotFound {
			failureReason = "用户不存在：手机号未注册"
		} else if err == service.ErrSmsCodeInvalid {
			failureReason = "短信验证码错误"
		} else if err == service.ErrSmsCodeExpired {
			failureReason = "短信验证码已过期"
		} else if err == service.ErrSmsCodeUsed {
			failureReason = "短信验证码已使用"
		} else {
			failureReason = "重置密码服务内部错误"
		}

		logger.Warn("重置密码失败",
			zap.String("操作", "重置密码"),
			zap.String("结果", "失败"),
			zap.String("失败原因", failureReason),
			zap.String("phone", req.Phone),
			zap.String("ip", clientIP),
			zap.Error(err))

		if err == service.ErrUserNotFound {
			response.Error(c, 400, "用户不存在")
			return
		}
		if err == service.ErrSmsCodeInvalid || err == service.ErrSmsCodeExpired || err == service.ErrSmsCodeUsed {
			response.Error(c, 400, err.Error())
			return
		}
		response.InternalError(c, "重置密码失败")
		return
	}

	logger.Info("重置密码成功",
		zap.String("操作", "重置密码"),
		zap.String("结果", "成功"),
		zap.String("phone", req.Phone),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "密码重置成功，请使用新密码登录", nil)
}

// SendSms 发送短信验证码
// @Summary      发送短信验证码
// @Description  发送短信验证码到指定手机号，支持注册和忘记密码场景
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SendSmsRequest  true  "发送短信请求"
// @Success      200      {object}  response.Response  "验证码已发送"
// @Failure      400      {object}  response.Response  "请求参数错误"
// @Failure      429      {object}  response.Response  "发送过于频繁，请稍后再试"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/signup/code [post]
func (h *UserHandler) SendSms(c *gin.Context) {
	var req dto.SendSmsRequest
	clientIP := c.ClientIP()

	// 记录发送短信验证码尝试开始
	logger.Info("开始发送短信验证码",
		zap.String("phone", req.Phone),
		zap.String("type", req.Type),
		zap.String("ip", clientIP),
		zap.String("user_agent", c.Request.UserAgent()))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("发送短信验证码失败",
			zap.String("操作", "发送短信验证码"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("phone", req.Phone),
			zap.String("type", req.Type),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 调用短信服务
	if err := h.smsAppService.SendSMS(c.Request.Context(), req.Phone, req.Type, clientIP); err != nil {
		var failureReason string
		if err == service.ErrSmsSendTooFast {
			failureReason = "发送验证码过于频繁（1分钟内只能发送一次）"
		} else {
			failureReason = "短信服务内部错误"
		}

		logger.Warn("发送短信验证码失败",
			zap.String("操作", "发送短信验证码"),
			zap.String("结果", "失败"),
			zap.String("失败原因", failureReason),
			zap.String("phone", req.Phone),
			zap.String("type", req.Type),
			zap.String("ip", clientIP),
			zap.Error(err))

		if err == service.ErrSmsSendTooFast {
			response.Error(c, 429, "发送验证码过于频繁，请稍后再试")
			return
		}
		response.InternalError(c, "发送验证码失败")
		return
	}

	actionMsg := "注册"
	if req.Type == "forgot" {
		actionMsg = "重置密码"
	}

	logger.Info("发送短信验证码成功",
		zap.String("操作", "发送短信验证码"),
		zap.String("结果", "成功"),
		zap.String("phone", req.Phone),
		zap.String("type", req.Type),
		zap.String("action", actionMsg),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, fmt.Sprintf("%s验证码已发送至 %s，请注意查收（10分钟内有效）", actionMsg, req.Phone), nil)
}

// RefreshToken 刷新Token
// @Summary      刷新Token
// @Description  使用RefreshToken刷新AccessToken，获取新的双Token
// @Tags         认证相关
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RefreshTokenRequest  true  "刷新Token请求"
// @Success      200      {object}  response.Response{data=dto.RefreshTokenResponse}  "刷新成功"
// @Failure      401      {object}  response.Response  "Token无效或已过期"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	clientIP := c.ClientIP()

	// 记录刷新Token尝试开始
	logger.Info("开始刷新Token",
		zap.String("ip", clientIP),
		zap.String("user_agent", c.Request.UserAgent()))

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("刷新Token失败",
			zap.String("操作", "刷新Token"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	tokenResp, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Warn("刷新Token失败",
			zap.String("操作", "刷新Token"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "Token无效或已过期"),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.Error(c, 401, err.Error())
		return
	}

	logger.Info("刷新Token成功",
		zap.String("操作", "刷新Token"),
		zap.String("结果", "成功"),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "Token刷新成功，已更新访问凭证", tokenResp)
}

// GetUserInfo 获取当前用户信息
// @Summary      获取当前用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200      {object}  response.Response{data=dto.UserResponse}  "获取成功"
// @Failure      401      {object}  response.Response  "未授权"
// @Failure      404      {object}  response.Response  "用户不存在"
// @Failure      500      {object}  response.Response  "服务器错误"
// @Router       /user/userinfo [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	clientIP := c.ClientIP()

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Warn("获取用户信息失败",
			zap.String("操作", "获取用户信息"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "未授权：用户ID不存在"),
			zap.String("ip", clientIP))
		response.Unauthorized(c, "未授权")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Error("获取用户信息失败",
			zap.String("操作", "获取用户信息"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "用户ID类型错误"),
			zap.String("ip", clientIP))
		response.InternalError(c, "用户ID类型错误")
		return
	}

	// 记录获取用户信息尝试开始
	logger.Info("开始获取用户信息",
		zap.String("user_id", userIDStr),
		zap.String("ip", clientIP))

	// 调用服务
	userResp, err := h.userService.GetUserInfo(c.Request.Context(), userIDStr)
	if err != nil {
		var failureReason string
		if err == service.ErrUserNotFound {
			failureReason = "用户不存在"
		} else {
			failureReason = "获取用户信息服务内部错误"
		}

		logger.Warn("获取用户信息失败",
			zap.String("操作", "获取用户信息"),
			zap.String("结果", "失败"),
			zap.String("失败原因", failureReason),
			zap.String("user_id", userIDStr),
			zap.String("ip", clientIP),
			zap.Error(err))

		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
			return
		}
		response.InternalError(c, "获取用户信息失败")
		return
	}

	logger.Info("获取用户信息成功",
		zap.String("操作", "获取用户信息"),
		zap.String("结果", "成功"),
		zap.String("user_id", userIDStr),
		zap.String("phone", func() string {
			if userResp.Phone != nil {
				return *userResp.Phone
			}
			return ""
		}()),
		zap.String("ip", clientIP))

	response.SuccessWithMessage(c, "获取用户信息成功", userResp)
}

// ListUsers 获取用户列表
// @Summary      获取用户列表
// @Description  分页获取用户列表，支持多条件查询
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page        query     int     false  "页码，从1开始"  default(1)
// @Param        page_size   query     int     false  "每页数量，最大100"  default(10)
// @Param        phone       query     string  false  "手机号（模糊查询）"
// @Param        email       query     string  false  "邮箱（模糊查询）"
// @Param        status      query     string  false  "状态（1启用，2禁用）"
// @Param        gender      query     string  false  "性别（1男，2女，3未知）"
// @Param        nickname    query     string  false  "昵称（模糊查询）"
// @Param        full_name   query     string  false  "姓名（模糊查询）"
// @Success      200         {object}  response.PageResponse{data=[]dto.UserResponse}  "获取成功"
// @Failure      401         {object}  response.Response  "未授权"
// @Failure      403         {object}  response.Response  "权限不足"
// @Failure      500         {object}  response.Response  "服务器错误"
// @Router       /user/list [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	clientIP := c.ClientIP()

	// 从上下文获取用户ID（用于日志）
	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	// 记录获取用户列表尝试开始
	logger.Info("开始获取用户列表",
		zap.String("user_id", userID),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
		zap.String("ip", clientIP))

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("获取用户列表失败",
			zap.String("操作", "获取用户列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "请求参数验证失败"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	users, total, err := h.userService.ListUsers(c.Request.Context(), &req)
	if err != nil {
		logger.Error("获取用户列表失败",
			zap.String("操作", "获取用户列表"),
			zap.String("结果", "失败"),
			zap.String("失败原因", "获取用户列表服务内部错误"),
			zap.String("user_id", userID),
			zap.String("ip", clientIP),
			zap.Error(err))
		response.InternalError(c, "获取用户列表失败")
		return
	}

	userCount := 0
	if users != nil {
		userCount = len(users)
	}

	logger.Info("获取用户列表成功",
		zap.String("操作", "获取用户列表"),
		zap.String("结果", "成功"),
		zap.String("user_id", userID),
		zap.Int64("total", total),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
		zap.Int("count", userCount),
		zap.String("ip", clientIP))

	response.PageSuccess(c, users, total, req.Page, req.PageSize)
}

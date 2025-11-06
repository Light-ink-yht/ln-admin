import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 登录请求参数
 */
export interface LoginParams {
    phone: string
    password: string
    captchaCode: string
    captchaId: string
}

/**
 * 注册请求参数
 */
export interface RegisterParams {
    phone: string
    code: string
    password: string
    confirmPassword: string
}

/**
 * 忘记密码请求参数
 */
export interface ForgotPasswordParams {
    phone: string
    password: string
    confirmPassword: string
    code: string
}

/**
 * 发送短信验证码参数
 */
export interface SendSmsParams {
    phone: string
    type: 'register' | 'forgot'
}

/**
 * 验证码响应
 */
export interface CaptchaData {
    captchaId: string
    captcha_code: string
}

/**
 * 登录用户信息（与后端返回格式一致）
 */
export interface UserInfo {
    ID: number
    CreatedAt: string
    UpdatedAt: string
    DeletedAt: string | null
    userId: string
    email: string
    phone: string
    nickname: string
    fullName: string
    avatar: string
    gender: string
    birthday: string | null
    status: string
    remarks: string
    loginCount: number
    lastLoginTime: string
    lastLoginIp: string
    passwordChange: string | null
    creatorId: string
    modifierId: string
}

/**
 * 认证 API
 */
export const authApi = {
    /**
     * 获取图片验证码
     */
    getCaptcha(): Promise<ResponseData<CaptchaData>> {
        return request.get<CaptchaData>('/user/captcha', "", {
            showSuccessMessage: true,
        })
    },

    /**
     * 登录响应数据
     */
    login(params: LoginParams): Promise<ResponseData<{
        accessToken: string
        refreshToken: string
        expiresIn: number
        refreshExpiresIn: number
        user: UserInfo
    }>> {
        return request.post('/user/login', params, {
            showSuccessMessage: false, // 在组件中显示自定义消息，避免重复
        })
    },

    /**
     * 注册
     */
    register(params: RegisterParams): Promise<ResponseData<UserInfo>> {
        return request.post('/user/signup', params, {
            showSuccessMessage: true,
        })
    },

    /**
     * 忘记密码
     */
    forgotPassword(params: ForgotPasswordParams): Promise<ResponseData<void>> {
        return request.post('/user/password', params, {
            showSuccessMessage: true,
        })
    },

    /**
     * 发送短信验证码
     */
    sendSms(params: SendSmsParams): Promise<ResponseData<void>> {
        return request.post('/user/signup/code', params, {
            showSuccessMessage: true,
        })
    },

    /**
     * 退出登录
     */
    logout(): Promise<ResponseData<void>> {
        return request.post('/user/logout')
    },

    /**
     * 刷新令牌
     */
    refreshToken(): Promise<ResponseData<{ AccessToken: string; RefreshToken: string }>> {
        return request.refreshToken()
    },

    /**
     * 获取当前用户信息
     */
    getUserInfo(): Promise<ResponseData<UserInfo>> {
        return request.get<UserInfo>('/user/userinfo')
    },
}


import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import { message } from 'ant-design-vue'
import { envConfig } from '@/env.config'
import { parseJWTToken } from '@/utils/jwt'

/**
 * 响应数据接口
 */
export interface ResponseData<T = any> {
    code: number
    message?: string
    msg?: string
    data: T
    timestamp?: number
}

/**
 * 请求配置接口
 */
export interface RequestConfig extends AxiosRequestConfig {
    // 是否显示成功消息
    showSuccessMessage?: boolean
    // 是否显示错误消息
    showErrorMessage?: boolean
    // 重试次数
    retry?: number
    // 重试延迟
    retryDelay?: number
}

/**
 * 创建 axios 实例
 */
const service: AxiosInstance = axios.create({
    baseURL: envConfig.apiBaseUrl,
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json;charset=UTF-8',
    },
    // 配置参数序列化，支持数组参数（key[]=value1&key[]=value2）
    paramsSerializer: {
        indexes: null, // 使用 key[]=value 格式而不是 key[0]=value
    },
})

/**
 * 请求拦截器
 */
service.interceptors.request.use(
    (config: any) => {
        // 添加 token
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }

        // 添加时间戳，防止缓存
        if (config.method === 'get') {
            config.params = {
                ...config.params,
                _t: Date.now(),
            }
        }

        return config
    },
    (error: AxiosError) => {
        console.error('Request Error:', error)
        return Promise.reject(error)
    }
)

/**
 * 响应拦截器
 */
service.interceptors.response.use(
    (response: AxiosResponse<ResponseData>) => {
        const res = response.data
        const config = response.config as RequestConfig

        // 拦截并存储 token
        const accessToken = response.headers['x-access-token']
        const refreshToken = response.headers['x-refresh-token']

        if (accessToken) {
            localStorage.setItem('token', accessToken)
            // 从新token中更新角色信息
            try {
                const claims = parseJWTToken(accessToken)
                if (claims && claims.roles && Array.isArray(claims.roles)) {
                    localStorage.setItem('user_roles', JSON.stringify(claims.roles))
                }
            } catch (error) {
                console.error('更新角色信息失败:', error)
            }
        }
        if (refreshToken) {
            localStorage.setItem('refreshToken', refreshToken)
        }

        // 如果响应是文件流，直接返回
        if (response.config.responseType === 'blob') {
            return response as any
        }

        // 根据业务状态码处理（支持多种响应格式）
        // code规则：0=成功，1=警告，2=错误
        const responseMessage = res.message || res.msg || ''

        if (res.code === 0) {
            // 成功
            if (config.showSuccessMessage && responseMessage) {
                message.success(responseMessage)
            }
            return res as any
        } else if (res.code === 1) {
            // 警告（如：验证码发送太频繁、账户未激活等）
            if (config.showErrorMessage !== false) {
                message.warning(responseMessage || '操作警告')
            }
            return Promise.reject(new Error(responseMessage || '操作警告'))
        } else if (res.code === 2) {
            // 业务逻辑错误，使用警告提示
            const errorMsg = responseMessage || '操作失败'
            if (config.showErrorMessage !== false) {
                message.warning(errorMsg)
            }
            return Promise.reject(new Error(errorMsg))
        } else if (res.code === 401) {
            // 未授权，跳转到登录页（系统级别错误，使用error）
            message.error(responseMessage || '未授权，请重新登录')
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
            setTimeout(() => {
                window.location.href = '/login'
            }, 1500)
            return Promise.reject(new Error(responseMessage || '未授权'))
        } else if (res.code === 200) {
            // 兼容 HTTP 200 成功状态码
            if (config.showSuccessMessage && responseMessage) {
                message.success(responseMessage)
            }
            return res as any
        } else if (res.code === 403) {
            // 权限不足（业务逻辑错误，使用警告）
            const errorMsg = responseMessage || '没有权限访问该资源'
            if (config.showErrorMessage !== false) {
                message.warning(errorMsg)
            }
            return Promise.reject(new Error(errorMsg))
        } else if (res.code === 500) {
            // 服务器内部错误（系统级别错误，使用error）
            const errorMsg = responseMessage || '服务器内部错误'
            if (config.showErrorMessage !== false) {
                message.error(errorMsg)
            }
            return Promise.reject(new Error(errorMsg))
        } else {
            // 其他错误码（业务逻辑错误，使用警告）
            const errorMsg = responseMessage || '请求失败'
            if (config.showErrorMessage !== false) {
                message.warning(errorMsg)
            }
            return Promise.reject(new Error(errorMsg))
        }
    },
    (error: AxiosError<ResponseData>) => {
        const config = error.config as RequestConfig
        let errorMessage = '请求失败'

        if (error.response) {
            // 服务器返回了错误状态码
            const status = error.response.status
            const data = error.response.data

            switch (status) {
                case 400:
                    // 业务逻辑错误，使用警告
                    errorMessage = data?.message || data?.msg || '请求参数错误'
                    break
                case 401:
                    // 系统级别错误，使用error
                    errorMessage = '未授权，请重新登录'
                    localStorage.removeItem('token')
                    localStorage.removeItem('refreshToken')
                    setTimeout(() => {
                        window.location.href = '/login'
                    }, 1500)
                    break
                case 403:
                    // 业务逻辑错误，使用警告
                    errorMessage = data?.message || data?.msg || '拒绝访问'
                    break
                case 404:
                    // 业务逻辑错误，使用警告
                    errorMessage = data?.message || data?.msg || '请求资源不存在'
                    break
                case 500:
                case 502:
                case 503:
                case 504:
                    // 系统级别错误，使用error
                    errorMessage = status === 500 ? '服务器内部错误' : 
                                   status === 502 ? '网关错误' : 
                                   status === 503 ? '服务不可用' : '网关超时'
                    break
                default:
                    // 其他错误，根据状态码判断
                    if (status >= 500) {
                        // 服务器错误，使用error
                        errorMessage = data?.message || data?.msg || `服务器错误，错误码：${status}`
                    } else {
                        // 业务逻辑错误，使用警告
                        errorMessage = data?.message || data?.msg || `请求失败，错误码：${status}`
                    }
            }
        } else if (error.request) {
            // 网络错误，使用警告
            errorMessage = '网络连接失败，请检查网络'
        } else {
            // 其他错误，使用警告
            errorMessage = error.message || '请求失败'
        }

        if (config?.showErrorMessage !== false) {
            // 根据状态码判断使用error还是warning
            if (error.response?.status && error.response.status >= 500) {
                message.error(errorMessage)
            } else if (error.response?.status === 401) {
                message.error(errorMessage)
            } else {
                message.warning(errorMessage)
            }
        }

        return Promise.reject(error)
    }
)

/**
 * 请求方法封装
 */
export const request = {
    /**
     * GET 请求
     */
    get<T = any>(url: string, params?: any, config?: RequestConfig): Promise<ResponseData<T>> {
        return service.get(url, { ...config, params })
    },

    /**
     * POST 请求
     */
    post<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ResponseData<T>> {
        return service.post(url, data, config)
    },

    /**
     * PUT 请求
     */
    put<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ResponseData<T>> {
        return service.put(url, data, config)
    },

    /**
     * DELETE 请求
     */
    delete<T = any>(url: string, params?: any, config?: RequestConfig): Promise<ResponseData<T>> {
        return service.delete(url, { ...config, params })
    },

    /**
     * PATCH 请求
     */
    patch<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ResponseData<T>> {
        return service.patch(url, data, config)
    },

    /**
     * 上传文件
     */
    upload<T = any>(
        url: string,
        file: File | FormData,
        config?: RequestConfig
    ): Promise<ResponseData<T>> {
        const formData = file instanceof FormData ? file : new FormData()
        if (!(file instanceof FormData)) {
            formData.append('file', file)
        }

        return service.post(url, formData, {
            ...config,
            headers: {
                ...config?.headers,
                'Content-Type': 'multipart/form-data',
            },
        })
    },

    /**
     * 刷新令牌
     */
    async refreshToken(): Promise<ResponseData<{ AccessToken: string; RefreshToken: string }>> {
        const refreshToken = localStorage.getItem('refreshToken')

        if (!refreshToken) {
            const error = new Error('刷新令牌不存在')
            message.error('刷新令牌不存在，请重新登录')
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
            setTimeout(() => {
                window.location.href = '/login'
            }, 1500)
            return Promise.reject(error)
        }

        // 使用独立的 axios 请求，避免触发请求拦截器添加 Authorization 头
        const response = await axios.post<ResponseData<{ AccessToken: string; RefreshToken: string }>>(
            `${envConfig.apiBaseUrl}/refresh`,
            {},
            {
                headers: {
                    'X-Refresh-Token': refreshToken,
                    'Content-Type': 'application/json;charset=UTF-8',
                },
                timeout: 30000,
            }
        )

        // 手动处理响应头中的 token
        const accessToken = response.headers['x-access-token']
        const newRefreshToken = response.headers['x-refresh-token']

        if (accessToken) {
            localStorage.setItem('token', accessToken)
            // 从新token中更新角色信息
            try {
                const claims = parseJWTToken(accessToken)
                if (claims && claims.roles && Array.isArray(claims.roles)) {
                    localStorage.setItem('user_roles', JSON.stringify(claims.roles))
                }
            } catch (error) {
                console.error('更新角色信息失败:', error)
            }
        }
        if (newRefreshToken) {
            localStorage.setItem('refreshToken', newRefreshToken)
        }

        return response.data
    },

    /**
     * 下载文件
     */
    download(url: string, params?: any, filename?: string): Promise<void> {
        return service
            .get(url, {
                params,
                responseType: 'blob',
            })
            .then((response) => {
                const blob = new Blob([response.data])
                const downloadUrl = window.URL.createObjectURL(blob)
                const link = document.createElement('a')
                link.href = downloadUrl
                link.download = filename || 'download'
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
                window.URL.revokeObjectURL(downloadUrl)
            })
    },
}

/**
 * 导出默认实例
 */
export default service


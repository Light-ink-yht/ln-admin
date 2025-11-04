import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 系统配置信息
 */
export interface SystemConfig {
    configKey: string
    configValue: string
    configName: string
    configGroup: string
    description?: string
    status?: string
    createdAt?: string
    updatedAt?: string
}

/**
 * 创建系统配置请求
 */
export interface CreateSystemConfigRequest {
    configKey: string
    configValue: string
    configName: string
    configGroup: string
    description?: string
    status?: string
}

/**
 * 更新系统配置请求
 */
export interface UpdateSystemConfigRequest {
    configValue?: string
    configName?: string
    configGroup?: string
    description?: string
    status?: string
}

/**
 * 系统日志信息
 */
export interface SystemLog {
    logId: string
    level: string
    module: string
    action: string
    content: string
    userId: string
    ip: string
    path: string
    method: string
    statusCode: number
    userAgent: string
    errorMsg?: string
    logTime: string
    createdAt: string
}

/**
 * 系统日志列表请求参数
 */
export interface SystemLogListParams {
    page?: number
    pageSize?: number
    level?: string
    module?: string
    action?: string
    userId?: string
    ip?: string
    startTime?: string
    endTime?: string
}

/**
 * 系统配置 API
 */
export const systemApi = {
    /**
     * 获取所有配置
     */
    getAllConfigs(): Promise<ResponseData<SystemConfig[]>> {
        return request.get<SystemConfig[]>('/system/config/list')
    },

    /**
     * 根据分组获取配置
     */
    getConfigsByGroup(group: string): Promise<ResponseData<SystemConfig[]>> {
        return request.get<SystemConfig[]>(`/system/config/group/${group}`)
    },

    /**
     * 根据配置键获取配置
     */
    getConfigByKey(key: string): Promise<ResponseData<SystemConfig>> {
        return request.get<SystemConfig>(`/system/config/${key}`)
    },

    /**
     * 创建配置
     */
    createConfig(data: CreateSystemConfigRequest): Promise<ResponseData<void>> {
        return request.post<void>('/system/config', data)
    },

    /**
     * 更新配置
     */
    updateConfig(key: string, data: UpdateSystemConfigRequest): Promise<ResponseData<void>> {
        return request.put<void>(`/system/config/${key}`, data)
    },

    /**
     * 获取系统日志列表
     */
    getLogList(params: SystemLogListParams): Promise<ResponseData<{ list: SystemLog[]; total: number }>> {
        return request.get<{ list: SystemLog[]; total: number }>('/system/log/list', params)
    },
}



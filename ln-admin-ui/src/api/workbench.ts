import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 工作台统计项
 */
export interface WorkbenchStatItem {
    key: string
    title: string
    value: number
    icon: string
    color: string
    trend: 'up' | 'down'
    change: string
    desc: string
    type: 'primary' | 'success' | 'warning' | 'danger'
    suffix?: string
}

/**
 * 工作台快捷操作
 */
export interface WorkbenchQuickAction {
    key: string
    title: string
    desc: string
    icon: string
    color: string
    path: string
}

/**
 * 工作台公告
 */
export interface WorkbenchAnnouncement {
    key: string
    title: string
    content: string
    time: string
    author: string
    type: 'normal' | 'important' | 'warning'
}

/**
 * 工作台配置数据
 */
export interface WorkbenchConfigData {
    stats: WorkbenchStatItem[]
    quickActions: WorkbenchQuickAction[]
    announcements: WorkbenchAnnouncement[]
    showSystemInfo: boolean
}

/**
 * 工作台配置响应
 */
export interface WorkbenchConfigResponse {
    role_key: string
    config_name: string
    config: WorkbenchConfigData
}

/**
 * 工作台配置列表项
 */
export interface WorkbenchConfigListItem {
    config_id: string
    role_key: string
    config_name: string
    description: string
    status: string
    created_at: string
    updated_at: string
}

/**
 * 创建工作台配置请求
 */
export interface WorkbenchConfigCreateRequest {
    role_key: string
    config_name: string
    config: WorkbenchConfigData
    description?: string
    status?: string
}

/**
 * 更新工作台配置请求
 */
export interface WorkbenchConfigUpdateRequest {
    config_id: string
    config_name?: string
    config?: WorkbenchConfigData
    description?: string
    status?: string
}

/**
 * 工作台 API
 */
export const workbenchApi = {
    /**
     * 获取当前用户的工作台配置
     */
    getWorkbenchConfig(): Promise<ResponseData<WorkbenchConfigResponse>> {
        return request.get<WorkbenchConfigResponse>('/workbench/config')
    },

    /**
     * 获取工作台配置列表
     */
    getWorkbenchConfigList(): Promise<ResponseData<WorkbenchConfigListItem[]>> {
        return request.get<WorkbenchConfigListItem[]>('/workbench/config/list')
    },

    /**
     * 根据ID获取工作台配置详情
     */
    getWorkbenchConfigById(configId: string): Promise<ResponseData<WorkbenchConfigResponse>> {
        return request.get<WorkbenchConfigResponse>(`/workbench/config/${configId}`)
    },

    /**
     * 创建工作台配置
     */
    createWorkbenchConfig(req: WorkbenchConfigCreateRequest): Promise<ResponseData<void>> {
        return request.post('/workbench/config', req)
    },

    /**
     * 更新工作台配置
     */
    updateWorkbenchConfig(req: WorkbenchConfigUpdateRequest): Promise<ResponseData<void>> {
        return request.put('/workbench/config', req)
    },

    /**
     * 删除工作台配置
     */
    deleteWorkbenchConfig(configId: string): Promise<ResponseData<void>> {
        return request.delete(`/workbench/config/${configId}`)
    },
}


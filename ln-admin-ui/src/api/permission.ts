import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 权限信息
 */
export interface Permission {
    permissionId: string
    permissionKey: string
    permissionName: string
    resourcePath: string
    method: string
    description?: string
    status?: string
    category?: string
    creatorId?: string
    creatorName?: string
    modifierId?: string
    modifierName?: string
    createdAt?: string
    updatedAt?: string
    CreatedAt?: string
}

/**
 * 权限列表请求参数
 */
export interface PermissionListParams {
    page?: number
    page_size?: number
    permission_key?: string
    permission_name?: string
    resource_path?: string
    method?: string
    status?: string
    category?: string
}

/**
 * 分页响应数据
 */
export interface PageResponse<T> extends ResponseData<T> {
    total?: number
    page?: number
    page_size?: number
    pageSize?: number
}

/**
 * 创建权限请求
 */
export interface CreatePermissionRequest {
    permissionKey: string
    permissionName: string
    resourcePath: string
    method: string
    description?: string
    status?: string
    category?: string
}

/**
 * 更新权限请求
 */
export interface UpdatePermissionRequest {
    permissionName?: string
    resourcePath?: string
    method?: string
    description?: string
    status?: string
    category?: string
}

/**
 * 权限 API
 */
export const permissionApi = {
    /**
     * 获取权限列表
     */
    getPermissionList(params?: PermissionListParams): Promise<PageResponse<Permission[]>> {
        return request.get<Permission[]>('/permission/list', params) as Promise<PageResponse<Permission[]>>
    },

    /**
     * 获取所有权限（不分页）
     */
    getAllPermissions(): Promise<ResponseData<Permission[]>> {
        return request.get<Permission[]>('/permission/all') as Promise<ResponseData<Permission[]>>
    },

    /**
     * 获取权限详情
     */
    getPermissionDetail(permissionId: string): Promise<ResponseData<Permission>> {
        return request.get<Permission>(`/permission/${permissionId}`) as Promise<ResponseData<Permission>>
    },

    /**
     * 创建权限
     */
    createPermission(data: CreatePermissionRequest): Promise<ResponseData<Permission>> {
        return request.post<Permission>('/permission', data) as Promise<ResponseData<Permission>>
    },

    /**
     * 更新权限
     */
    updatePermission(permissionId: string, data: UpdatePermissionRequest): Promise<ResponseData<Permission>> {
        return request.put<Permission>(`/permission/${permissionId}`, data) as Promise<ResponseData<Permission>>
    },

    /**
     * 删除权限
     */
    deletePermission(permissionId: string): Promise<ResponseData<void>> {
        return request.delete<void>(`/permission/${permissionId}`) as Promise<ResponseData<void>>
    },
}


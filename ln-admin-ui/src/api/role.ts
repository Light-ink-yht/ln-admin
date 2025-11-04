import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 角色信息
 */
export interface Role {
    roleId: string
    roleKey: string
    roleName: string
    description?: string
    status?: string
    creatorId?: string
    creatorName?: string
    modifierId?: string
    modifierName?: string
    createdAt?: string
    updatedAt?: string
    CreatedAt?: string
}

/**
 * 创建角色请求
 */
export interface CreateRoleRequest {
    roleKey: string
    roleName: string
    description?: string
    status?: string
}

/**
 * 更新角色请求
 */
export interface UpdateRoleRequest {
    roleName?: string
    description?: string
    status?: string
}

/**
 * 角色列表请求参数
 */
export interface RoleListParams {
    page?: number
    page_size?: number
    role_key?: string
    role_name?: string
    status?: string
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
 * 角色 API
 */
export const roleApi = {
    /**
     * 获取角色列表
     */
    getRoleList(params?: RoleListParams): Promise<PageResponse<Role[]>> {
        return request.get<Role[]>('/role/list', params) as Promise<PageResponse<Role[]>>
    },

    /**
     * 获取所有角色（不分页）
     */
    getAllRoles(): Promise<ResponseData<Role[]>> {
        return request.get<Role[]>('/role/all', undefined, { showErrorMessage: false }) as Promise<ResponseData<Role[]>>
    },

    /**
     * 获取角色详情
     */
    getRoleDetail(roleId: string): Promise<ResponseData<Role>> {
        return request.get<Role>(`/role/${roleId}`) as Promise<ResponseData<Role>>
    },

    /**
     * 创建角色
     */
    createRole(data: CreateRoleRequest): Promise<ResponseData<Role>> {
        return request.post<Role>('/role', data) as Promise<ResponseData<Role>>
    },

    /**
     * 更新角色
     */
    updateRole(roleId: string, data: UpdateRoleRequest): Promise<ResponseData<Role>> {
        return request.put<Role>(`/role/${roleId}`, data) as Promise<ResponseData<Role>>
    },

    /**
     * 删除角色
     */
    deleteRole(roleId: string): Promise<ResponseData<void>> {
        return request.delete<void>(`/role/${roleId}`) as Promise<ResponseData<void>>
    },
}


import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 用户列表请求参数
 */
export interface UserListParams {
    page?: number
    page_size?: number
    phone?: string
    email?: string
    status?: string
    gender?: string
    nickname?: string
    full_name?: string
}

/**
 * 用户响应数据
 */
export interface UserResponse {
    ID: number
    CreatedAt: string
    UpdatedAt: string
    DeletedAt?: string | null
    userId: string
    email?: string | null
    phone?: string | null
    nickname: string
    fullName: string
    avatar: string
    gender: string
    birthday?: string | null
    status: string
    remarks: string
    loginCount: number
    lastLoginTime?: string | null
    lastLoginIp: string
    passwordChange?: string | null
    creatorId: string
    modifierId: string
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
 * 创建用户请求参数
 */
export interface CreateUserRequest {
    phone: string
    email?: string
    password: string
    nickname: string
    fullName: string
    gender?: string
    status?: string
    avatar?: string
    birthday?: string
    remarks?: string
    roleIds?: string[]
}

/**
 * 更新用户请求参数
 */
export interface UpdateUserRequest {
    email?: string
    nickname?: string
    fullName?: string
    gender?: string
    status?: string
    avatar?: string
    birthday?: string
    remarks?: string
    roleIds?: string[]
}

/**
 * 用户详情响应（包含角色和权限）
 */
export interface UserDetailResponse extends UserResponse {
    roles?: RoleInfo[]
    permissions?: PermissionInfo[]
}

/**
 * 角色信息
 */
export interface RoleInfo {
    roleId: string
    roleKey: string
    roleName: string
    status?: string
}

/**
 * 权限信息
 */
export interface PermissionInfo {
    permissionId: string
    permissionKey: string
    permissionName: string
    resourcePath: string
    method: string
    description?: string
}

/**
 * 用户 API
 */
export const userApi = {
    /**
     * 获取用户列表
     */
    getUserList(params?: UserListParams): Promise<PageResponse<UserResponse[]>> {
        return request.get<UserResponse[]>('/user/list', params) as Promise<PageResponse<UserResponse[]>>
    },

    /**
     * 获取用户详情
     */
    getUserDetail(userId: string): Promise<ResponseData<UserDetailResponse>> {
        return request.get<UserDetailResponse>(`/user/${userId}`) as Promise<ResponseData<UserDetailResponse>>
    },

    /**
     * 创建用户
     */
    createUser(data: CreateUserRequest): Promise<ResponseData<UserResponse>> {
        return request.post<UserResponse>('/user', data) as Promise<ResponseData<UserResponse>>
    },

    /**
     * 更新用户
     */
    updateUser(userId: string, data: UpdateUserRequest): Promise<ResponseData<UserResponse>> {
        return request.put<UserResponse>(`/user/${userId}`, data) as Promise<ResponseData<UserResponse>>
    },

    /**
     * 删除用户
     */
    deleteUser(userId: string): Promise<ResponseData<void>> {
        return request.delete<void>(`/user/${userId}`) as Promise<ResponseData<void>>
    },

    /**
     * 给用户授权（直接分配权限）
     */
    grantPermissions(userId: string, permissionIds: string[]): Promise<ResponseData<void>> {
        return request.post<void>(`/user/${userId}/permissions`, { permissionIds }) as Promise<ResponseData<void>>
    },
}


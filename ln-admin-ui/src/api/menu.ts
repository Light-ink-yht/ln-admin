import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'
import type { Component } from 'vue'

/**
 * 菜单项类型
 */
export interface MenuItem {
    key: string
    title: string
    path?: string
    icon?: Component | string
    children?: MenuItem[]
    permission?: string // 权限标识
    hidden?: boolean // 是否隐藏
}

/**
 * 用户下拉菜单项
 */
export interface UserMenuItem {
    key: string
    label: string
    icon?: Component | string
    divider?: boolean // 是否在此项前添加分隔线
    onClick?: () => void | Promise<void>
}

/**
 * 菜单响应数据
 */
export interface Menu {
    id: number
    created_at: string
    updated_at: string
    menu_id: string
    menu_key: string
    title: string
    path: string
    icon: string
    parent_id: string
    sort: number
    permission: string
    menu_type: string
    status: string
    description: string
    creator_id: string
    modifier_id: string
    children?: Menu[]
}

/**
 * 菜单列表请求参数
 */
export interface MenuListParams {
    page?: number
    pageSize?: number
    menu_type?: string
}

/**
 * 创建菜单请求
 */
export interface CreateMenuRequest {
    menu_key: string
    title: string
    path?: string
    icon?: string
    parent_id?: string
    sort?: number
    permission?: string
    menu_type: string
    description?: string
}

/**
 * 更新菜单请求
 */
export interface UpdateMenuRequest {
    title: string
    path?: string
    icon?: string
    parent_id?: string
    sort?: number
    permission?: string
    status?: string
    description?: string
}

/**
 * 菜单列表响应
 */
export interface MenuListResponse {
    list: Menu[]
    total: number
}

/**
 * 菜单 API
 */
export const menuApi = {
    /**
     * 获取侧边栏菜单列表
     */
    getSidebarMenus(): Promise<ResponseData<MenuItem[]>> {
        return request.get<MenuItem[]>('/menu/sidebar')
    },

    /**
     * 获取用户下拉菜单
     */
    getUserMenus(): Promise<ResponseData<UserMenuItem[]>> {
        return request.get<UserMenuItem[]>('/menu/user')
    },

    /**
     * 获取菜单列表
     */
    getMenuList(params?: MenuListParams): Promise<ResponseData<MenuListResponse>> {
        return request.get<MenuListResponse>('/menu/list', params)
    },

    /**
     * 获取菜单详情
     */
    getMenuDetail(menuId: string): Promise<ResponseData<Menu>> {
        return request.get<Menu>(`/menu/${menuId}`)
    },

    /**
     * 创建菜单
     */
    createMenu(data: CreateMenuRequest): Promise<ResponseData<Menu>> {
        return request.post<Menu>('/menu', data)
    },

    /**
     * 更新菜单
     */
    updateMenu(menuId: string, data: UpdateMenuRequest): Promise<ResponseData<Menu>> {
        return request.put<Menu>(`/menu/${menuId}`, data)
    },

    /**
     * 删除菜单
     */
    deleteMenu(menuId: string): Promise<ResponseData<void>> {
        return request.delete<void>(`/menu/${menuId}`)
    },

    /**
     * 获取菜单树
     */
    getMenuTree(menuType: string): Promise<ResponseData<Menu[]>> {
        return request.get<Menu[]>('/menu/tree', { menu_type: menuType })
    },

    /**
     * 为角色分配菜单
     */
    grantRoleMenus(roleId: string, menuIds: string[]): Promise<ResponseData<void>> {
        return request.post<void>(`/menu/role/${roleId}/menus`, { menu_ids: menuIds })
    },

    /**
     * 获取角色的菜单列表
     */
    getRoleMenus(roleId: string): Promise<ResponseData<{ menu_ids: string[] }>> {
        return request.get<{ menu_ids: string[] }>(`/menu/role/${roleId}/menus`)
    },
}


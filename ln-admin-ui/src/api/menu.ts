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
}


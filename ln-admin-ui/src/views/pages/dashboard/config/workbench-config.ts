/**
 * 工作台配置 - 根据不同角色显示不同的工作台内容
 */

import type { QuickAction, StatItem, Announcement } from '../components'

// 角色类型
export type RoleType = 'admin' | 'manager' | 'user' | 'guest'

// 工作台配置接口
export interface WorkbenchConfig {
    role: RoleType
    stats: StatItem[]
    quickActions: QuickAction[]
    announcements: Announcement[]
    showSystemInfo: boolean
    showRecentLogs: boolean
}

/**
 * 管理员工作台配置
 */
export const adminWorkbenchConfig: WorkbenchConfig = {
    role: 'admin',
    stats: [
        {
            key: 'users',
            title: '用户总数',
            value: 0,
            icon: 'TeamOutlined',
            color: '#1890ff',
            trend: 'up',
            change: '12%',
            desc: '较上月',
            type: 'primary',
            suffix: '人',
        },
        {
            key: 'roles',
            title: '角色数量',
            value: 0,
            icon: 'SafetyOutlined',
            color: '#52c41a',
            trend: 'up',
            change: '5%',
            desc: '较上月',
            type: 'success',
            suffix: '个',
        },
        {
            key: 'permissions',
            title: '权限数量',
            value: 0,
            icon: 'UnlockOutlined',
            color: '#faad14',
            trend: 'down',
            change: '2%',
            desc: '较上月',
            type: 'warning',
            suffix: '个',
        },
        {
            key: 'logs',
            title: '今日操作',
            value: 0,
            icon: 'FileTextOutlined',
            color: '#f5222d',
            trend: 'up',
            change: '23%',
            desc: '较昨日',
            type: 'danger',
            suffix: '次',
        },
    ],
    quickActions: [
        {
            key: 'user',
            title: '用户管理',
            desc: '管理系统用户',
            icon: 'UserAddOutlined',
            color: 'var(--color-primary)',
            path: '/user/list',
        },
        {
            key: 'role',
            title: '角色管理',
            desc: '配置用户角色',
            icon: 'SafetyOutlined',
            color: 'var(--color-primary)',
            path: '/role/list',
        },
        {
            key: 'permission',
            title: '权限管理',
            desc: '管理权限配置',
            icon: 'UnlockOutlined',
            color: 'var(--color-primary)',
            path: '/permission/list',
        },
        {
            key: 'log',
            title: '操作日志',
            desc: '查看系统日志',
            icon: 'AuditOutlined',
            color: 'var(--color-primary)',
            path: '/system/log',
        },
        {
            key: 'config',
            title: '系统配置',
            desc: '系统参数设置',
            icon: 'SettingOutlined',
            color: 'var(--color-primary)',
            path: '/system/config',
        },
        {
            key: 'database',
            title: '数据管理',
            desc: '数据库维护',
            icon: 'DatabaseOutlined',
            color: 'var(--color-primary)',
            path: '/system/config',
        },
    ],
    announcements: [
        {
            key: '1',
            title: '系统升级通知',
            content: '系统将于本周六晚上22:00-24:00进行升级维护，期间可能影响部分功能使用，请提前做好准备。',
            time: '2024-01-15 10:30',
            author: '系统管理员',
            type: 'important',
        },
    ],
    showSystemInfo: true,
    showRecentLogs: true,
}

/**
 * 管理者工作台配置
 */
export const managerWorkbenchConfig: WorkbenchConfig = {
    role: 'manager',
    stats: [
        {
            key: 'users',
            title: '团队成员',
            value: 0,
            icon: 'TeamOutlined',
            color: '#1890ff',
            trend: 'up',
            change: '8%',
            desc: '较上月',
            type: 'primary',
            suffix: '人',
        },
        {
            key: 'tasks',
            title: '待处理任务',
            value: 0,
            icon: 'FileTextOutlined',
            color: '#faad14',
            trend: 'down',
            change: '15%',
            desc: '较昨日',
            type: 'warning',
            suffix: '个',
        },
        {
            key: 'projects',
            title: '进行中项目',
            value: 0,
            icon: 'DatabaseOutlined',
            color: '#52c41a',
            trend: 'up',
            change: '3%',
            desc: '较上月',
            type: 'success',
            suffix: '个',
        },
        {
            key: 'reports',
            title: '本月报告',
            value: 0,
            icon: 'AuditOutlined',
            color: '#722ed1',
            trend: 'up',
            change: '20%',
            desc: '较上月',
            type: 'primary',
            suffix: '份',
        },
    ],
    quickActions: [
        {
            key: 'team',
            title: '团队管理',
            desc: '管理团队成员',
            icon: 'TeamOutlined',
            color: 'var(--color-primary)',
            path: '/user/list',
        },
        {
            key: 'task',
            title: '任务分配',
            desc: '分配和查看任务',
            icon: 'FileTextOutlined',
            color: 'var(--color-primary)',
            path: '/dashboard/analysis',
        },
        {
            key: 'report',
            title: '数据报表',
            desc: '查看业务报表',
            icon: 'AuditOutlined',
            color: 'var(--color-primary)',
            path: '/dashboard/analysis',
        },
        {
            key: 'log',
            title: '操作记录',
            desc: '查看操作日志',
            icon: 'AuditOutlined',
            color: 'var(--color-primary)',
            path: '/system/log',
        },
    ],
    announcements: [
        {
            key: '1',
            title: '团队会议通知',
            content: '本周五下午2点将举行团队会议，请准时参加。',
            time: '2024-01-14 09:00',
            author: '部门经理',
            type: 'warning',
        },
    ],
    showSystemInfo: false,
    showRecentLogs: true,
}

/**
 * 普通用户工作台配置
 */
export const userWorkbenchConfig: WorkbenchConfig = {
    role: 'user',
    stats: [
        {
            key: 'tasks',
            title: '我的任务',
            value: 0,
            icon: 'FileTextOutlined',
            color: '#1890ff',
            trend: 'down',
            change: '10%',
            desc: '较昨日',
            type: 'primary',
            suffix: '个',
        },
        {
            key: 'completed',
            title: '已完成',
            value: 0,
            icon: 'SafetyOutlined',
            color: '#52c41a',
            trend: 'up',
            change: '25%',
            desc: '较上月',
            type: 'success',
            suffix: '个',
        },
        {
            key: 'messages',
            title: '未读消息',
            value: 0,
            icon: 'BellOutlined',
            color: '#faad14',
            trend: 'down',
            change: '5%',
            desc: '较昨日',
            type: 'warning',
            suffix: '条',
        },
        {
            key: 'projects',
            title: '参与项目',
            value: 0,
            icon: 'DatabaseOutlined',
            color: '#722ed1',
            trend: 'up',
            change: '2%',
            desc: '较上月',
            type: 'primary',
            suffix: '个',
        },
    ],
    quickActions: [
        {
            key: 'profile',
            title: '个人资料',
            desc: '查看和编辑资料',
            icon: 'UserOutlined',
            color: 'var(--color-primary)',
            path: '/profile',
        },
        {
            key: 'settings',
            title: '账户设置',
            desc: '修改密码和设置',
            icon: 'SettingOutlined',
            color: 'var(--color-primary)',
            path: '/settings',
        },
        {
            key: 'tasks',
            title: '我的任务',
            desc: '查看我的任务',
            icon: 'FileTextOutlined',
            color: 'var(--color-primary)',
            path: '/dashboard/analysis',
        },
    ],
    announcements: [
        {
            key: '1',
            title: '系统使用指南',
            content: '欢迎使用 LN Admin 管理系统，如有任何问题请联系管理员。',
            time: '2024-01-10 08:00',
            author: '系统管理员',
            type: 'normal',
        },
    ],
    showSystemInfo: false,
    showRecentLogs: false,
}

/**
 * 访客工作台配置
 */
export const guestWorkbenchConfig: WorkbenchConfig = {
    role: 'guest',
    stats: [],
    quickActions: [
        {
            key: 'profile',
            title: '个人资料',
            desc: '查看个人资料',
            icon: 'UserOutlined',
            color: 'var(--color-primary)',
            path: '/profile',
        },
    ],
    announcements: [],
    showSystemInfo: false,
    showRecentLogs: false,
}

/**
 * 根据角色获取工作台配置
 */
export function getWorkbenchConfigByRole(roleKey: string | null | undefined): WorkbenchConfig {
    // 根据角色标识匹配配置
    if (!roleKey) {
        return guestWorkbenchConfig
    }

    const roleKeyLower = roleKey.toLowerCase()

    if (roleKeyLower.includes('admin') || roleKeyLower.includes('super')) {
        return adminWorkbenchConfig
    }

    if (roleKeyLower.includes('manager') || roleKeyLower.includes('manager')) {
        return managerWorkbenchConfig
    }

    if (roleKeyLower.includes('user') || roleKeyLower.includes('member')) {
        return userWorkbenchConfig
    }

    // 默认返回访客配置
    return guestWorkbenchConfig
}

/**
 * 角色配置映射表
 */
export const roleConfigMap: Record<string, WorkbenchConfig> = {
    admin: adminWorkbenchConfig,
    manager: managerWorkbenchConfig,
    user: userWorkbenchConfig,
    guest: guestWorkbenchConfig,
}


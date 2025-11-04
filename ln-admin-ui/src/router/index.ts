import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { message } from 'ant-design-vue'
import ThemeDemo from '@/views/ThemeDemo.vue'
import Login from '@/views/login/Login.vue'
import Layout from '@/views/layout/layout.vue'

// ==================== 路由元信息类型定义 ====================
declare module 'vue-router' {
    interface RouteMeta {
        title: string
        category?: 'dashboard' | 'user' | 'system' | 'profile'
        requiresAuth?: boolean
        icon?: string
        hidden?: boolean
        keepAlive?: boolean
    }
}

// ==================== 路由组件懒加载 ====================
// 仪表盘模块
const Workbench = () => import('@/views/pages/dashboard/Workbench.vue')
const DashboardAnalysis = () => import('@/views/pages/dashboard/DashboardAnalysis.vue')

// 用户管理模块
const UserList = () => import('@/views/pages/user/UserList.vue')
const RoleList = () => import('@/views/pages/user/RoleList.vue')
const PermissionList = () => import('@/views/pages/user/PermissionList.vue')

// 短信管理模块
const SMSTemplateList = () => import('@/views/pages/sms/SMSTemplateList.vue')
const SMSCodeList = () => import('@/views/pages/sms/SMSCodeList.vue')

// 菜单管理模块
const MenuList = () => import('@/views/pages/menu/MenuList.vue')

// 文件管理模块
const FileList = () => import('@/views/pages/file/FileList.vue')
const FileStorageConfig = () => import('@/views/pages/file/FileStorageConfig.vue')

// 系统运维模块
const SystemMonitor = () => import('@/views/pages/ops/SystemMonitor.vue')
const APIDoc = () => import('@/views/pages/ops/APIDoc.vue')

// 系统管理模块
const SystemConfig = () => import('@/views/pages/system/SystemConfig.vue')
const SystemLog = () => import('@/views/pages/system/SystemLog.vue')
const SettingsBasic = () => import('@/views/pages/system/SettingsBasic.vue')
const SettingsSecurity = () => import('@/views/pages/system/SettingsSecurity.vue')

// 用户信息模块
const UserProfile = () => import('@/views/pages/profile/UserProfile.vue')
const UserSettings = () => import('@/views/pages/profile/UserSettings.vue')

// ==================== 路由配置 ====================
const routes: RouteRecordRaw[] = [
    // 公开路由（无需认证）
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {
            title: '登录',
            requiresAuth: false,
        },
    },
    {
        path: '/theme-demo',
        name: 'ThemeDemo',
        component: ThemeDemo,
        meta: {
            title: '主题演示',
            requiresAuth: false,
        },
    },
    // 主布局路由（需要认证）
    {
        path: '/',
        component: Layout,
        redirect: '/dashboard/workbench',
        meta: {
            title: '首页',
            requiresAuth: true,
        },
        children: [
            // ========== 仪表盘模块 ==========
            {
                path: '',
                name: 'Workbench',
                component: Workbench,
                meta: {
                    title: '工作台',
                    category: 'dashboard',
                    requiresAuth: true,
                    keepAlive: true,
                },
            },
            {
                path: 'dashboard/workbench',
                name: 'DashboardWorkbench',
                component: Workbench,
                meta: {
                    title: '工作台',
                    category: 'dashboard',
                    requiresAuth: true,
                    keepAlive: true,
                },
            },
            {
                path: 'dashboard/analysis',
                name: 'DashboardAnalysis',
                component: DashboardAnalysis,
                meta: {
                    title: '分析页',
                    category: 'dashboard',
                    requiresAuth: true,
                },
            },
            // ========== 用户管理模块 ==========
            {
                path: 'user/list',
                name: 'UserList',
                component: UserList,
                meta: {
                    title: '用户列表',
                    category: 'user',
                    requiresAuth: true,
                },
            },
            {
                path: 'role/list',
                name: 'RoleList',
                component: RoleList,
                meta: {
                    title: '角色管理',
                    category: 'user',
                    requiresAuth: true,
                },
            },
            {
                path: 'permission/list',
                name: 'PermissionList',
                component: PermissionList,
                meta: {
                    title: '权限管理',
                    category: 'user',
                    requiresAuth: true,
                },
            },
            // ========== 短信管理模块 ==========
            {
                path: 'sms/template',
                name: 'SMSTemplateList',
                component: SMSTemplateList,
                meta: {
                    title: '短信模板',
                    category: 'sms',
                    requiresAuth: true,
                },
            },
            {
                path: 'sms/code',
                name: 'SMSCodeList',
                component: SMSCodeList,
                meta: {
                    title: '短信验证码',
                    category: 'sms',
                    requiresAuth: true,
                },
            },
            // ========== 菜单管理模块 ==========
            {
                path: 'menu/list',
                name: 'MenuList',
                component: MenuList,
                meta: {
                    title: '菜单列表',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            // ========== 文件管理模块 ==========
            {
                path: 'file/list',
                name: 'FileList',
                component: FileList,
                meta: {
                    title: '文件列表',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            {
                path: 'file/storage/config',
                name: 'FileStorageConfig',
                component: FileStorageConfig,
                meta: {
                    title: '存储配置',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            // ========== 系统运维模块 ==========
            {
                path: 'ops/api-doc',
                name: 'APIDoc',
                component: APIDoc,
                meta: {
                    title: '接口文档',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            {
                path: 'ops/monitor',
                name: 'SystemMonitor',
                component: SystemMonitor,
                meta: {
                    title: '系统监控',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            // ========== 系统管理模块 ==========
            {
                path: 'system/config',
                name: 'SystemConfig',
                component: SystemConfig,
                meta: {
                    title: '系统配置',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            {
                path: 'system/log',
                name: 'SystemLog',
                component: SystemLog,
                meta: {
                    title: '操作日志',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            {
                path: 'settings/basic',
                name: 'SettingsBasic',
                component: SettingsBasic,
                meta: {
                    title: '基础设置',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            {
                path: 'settings/security',
                name: 'SettingsSecurity',
                component: SettingsSecurity,
                meta: {
                    title: '安全设置',
                    category: 'system',
                    requiresAuth: true,
                },
            },
            // ========== 用户信息模块 ==========
            {
                path: 'profile',
                name: 'UserProfile',
                component: UserProfile,
                meta: {
                    title: '个人资料',
                    category: 'profile',
                    requiresAuth: true,
                },
            },
            {
                path: 'settings',
                name: 'UserSettings',
                component: UserSettings,
                meta: {
                    title: '设置',
                    category: 'profile',
                    requiresAuth: true,
                },
            },
        ],
    },
    // 404 页面
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        redirect: '/dashboard/workbench',
        meta: {
            title: '页面未找到',
            requiresAuth: false,
        },
    },
]

// ==================== 创建路由实例 ====================
const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
    scrollBehavior(to, from, savedPosition) {
        // 如果有保存的位置，恢复它
        if (savedPosition) {
            return savedPosition
        }
        // 对于相同路由（只是参数变化），保持当前位置
        if (to.path === from.path) {
            return false
        }
        // 对于不同路由，滚动到内容容器顶部（而不是window顶部）
        // 返回 false 可以禁用自动滚动，因为我们的布局使用了固定高度的容器
        return false
    },
})

// ==================== 路由守卫 ====================
/**
 * 检查用户是否已登录
 */
const isAuthenticated = (): boolean => {
    // 优先从 localStorage 检查（兼容性）
    const token = localStorage.getItem('token')
    return !!token
}

/**
 * 全局前置守卫
 */
router.beforeEach(
    async (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
        const requiresAuth = to.meta?.requiresAuth !== false

        // 如果路由需要认证但用户未登录
        if (requiresAuth && !isAuthenticated()) {
            // 保存目标路由，登录后可以重定向
            if (to.path !== '/login') {
                next({
                    path: '/login',
                    query: { redirect: to.fullPath },
                })
            } else {
                next()
            }
            return
        }

        // 如果已登录用户访问登录页，重定向到首页
        if (to.path === '/login' && isAuthenticated()) {
            next('/dashboard/workbench')
            return
        }

        // 设置页面标题
        if (to.meta?.title) {
            document.title = `${to.meta.title} - LN Admin`
        }

        next()
    }
)

/**
 * 全局后置守卫
 */
router.afterEach((to: RouteLocationNormalized, from: RouteLocationNormalized) => {
    // 可以在这里添加页面访问统计、埋点等逻辑
    if (import.meta.env.DEV) {
        console.log(`路由切换: ${from.path} -> ${to.path}`)
    }
})

/**
 * 全局错误处理
 */
router.onError((error) => {
    console.error('路由错误:', error)
    message.error('页面加载失败，请刷新重试')
})

// ==================== 路由工具函数 ====================
/**
 * 根据路由名称导航
 */
export const navigateByName = (name: string, params?: Record<string, any>) => {
    router.push({ name, params })
}

/**
 * 根据路由路径导航
 */
export const navigateByPath = (path: string, query?: Record<string, any>) => {
    router.push({ path, query })
}

/**
 * 返回上一页
 */
export const goBack = () => {
    router.back()
}

/**
 * 获取当前路由信息
 */
export const getCurrentRoute = () => {
    return router.currentRoute.value
}

export default router

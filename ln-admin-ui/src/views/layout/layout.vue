<template>
    <a-layout class="layout-container" :class="{ 'dark-layout': isDark }">
        <Header :user-info="userInfo" @toggle-sidebar="handleToggleSidebar" />
        <a-layout class="layout-body">
            <!-- 移动端抽屉 -->
            <a-drawer
                v-model:open="sidebarOpen"
                placement="left"
                :body-style="{ padding: 0 }"
                :width="200"
                class="mobile-drawer"
                @close="sidebarOpen = false"
            >
                <Sidebar @close="sidebarOpen = false" />
            </a-drawer>
            <!-- 桌面端侧边栏 -->
            <Sidebar class="desktop-sidebar" />
            <a-layout class="main-content" :style="{ background: contentBgColor }">
                <PageTabs
                    :tabs="pageTabs"
                    :active-key="activeTabKey"
                    @update:active-key="handleTabChange"
                    @tab-close="handleTabClose"
                    @tab-refresh="handleTabRefresh"
                    @close-other="handleCloseOther"
                    @close-all="handleCloseAll"
                />
                <div class="content-wrapper">
                    <a-layout-content :style="contentStyle">
                        <transition name="page-fade" mode="out-in">
                            <router-view v-slot="{ Component }">
                                <component :is="Component" />
                            </router-view>
                        </transition>
                    </a-layout-content>
                </div>
            </a-layout>
        </a-layout>
        <a-layout-footer class="footer-wrapper">
            <Footer />
        </a-layout-footer>
    </a-layout>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useThemeStore } from '@/stores/modules/theme'
import { useUserStore } from '@/stores/modules/user'
import { useSystemConfigStore } from '@/stores/modules/systemConfig'
import Header from './components/Header.vue'
import Sidebar from './components/Sidebar.vue'
import Footer from './components/Footer.vue'
import PageTabs from './components/PageTabs.vue'
import type { TabItem } from './components/PageTabs.vue'

const router = useRouter()
const route = useRoute()
const themeStore = useThemeStore()
const userStore = useUserStore()
const systemConfigStore = useSystemConfigStore()

// 用户信息（从store获取）
const userInfo = computed(() => userStore.userInfo)
const isDark = computed(() => themeStore.isDark)

// 移动端侧边栏控制
const sidebarOpen = ref(false)
const isMobile = ref(false)

const handleToggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
}

// 检测移动端
const checkMobile = () => {
    isMobile.value = window.innerWidth <= 768
}

onMounted(async () => {
    // 初始化用户信息（从存储恢复或从服务器获取）
    await userStore.initUserInfo()
    // 初始化系统配置
    await systemConfigStore.initConfigs()
    checkMobile()
    window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
    window.removeEventListener('resize', checkMobile)
})

// 路由标题映射
const routeTitleMap: Record<string, string> = {
    '/': '工作台',
    '/dashboard': '仪表盘',
    '/dashboard/workbench': '工作台',
    '/dashboard/analysis': '分析页',
    '/user': '用户管理',
    '/user/list': '用户列表',
    '/user/profile': '个人资料',
    '/role/list': '角色管理',
    '/permission/list': '权限管理',
    '/settings': '设置',
    '/settings/basic': '基础设置',
    '/settings/security': '安全设置',
    '/system/config': '系统配置',
    '/system/log': '操作日志',
}

// 根据路径获取标题
const getRouteTitle = (path: string, name?: string | symbol): string => {
    // 优先使用路由的 meta.title
    if (route.meta?.title && typeof route.meta.title === 'string') {
        return route.meta.title
    }
    // 优先使用映射表中的标题
    if (routeTitleMap[path]) {
        return routeTitleMap[path]
    }
    // 如果有路由名称且不是 symbol，使用路由名称
    if (name && typeof name === 'string') {
        return name
    }
    // 从路径推断标题
    const pathSegments = path.split('/').filter(Boolean)
    if (pathSegments.length > 0) {
        const lastSegment = pathSegments[pathSegments.length - 1]
        return lastSegment
            .replace(/-/g, ' ')
            .replace(/([A-Z])/g, ' $1')
            .trim()
            .split(' ')
            .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ') || '页面'
    }
    return '首页'
}

// 页面标签页
const pageTabs = ref<TabItem[]>([])
const activeTabKey = ref<string>('')

// 初始化首页标签
const initHomeTab = () => {
    const homeTab: TabItem = {
        key: 'home',
        title: '工作台',
        path: '/',
        closable: false,
    }
    pageTabs.value = [homeTab]
    activeTabKey.value = 'home'
}

// 初始化
initHomeTab()

const contentBgColor = computed(() => {
    return isDark.value ? '#141414' : '#f0f2f5'
})

const contentStyle = computed(() => {
    const padding = isMobile.value ? '16px' : '24px'
    
    if (isDark.value) {
        return {
            background: '#141414',
            padding: padding,
            margin: 0,
            minHeight: '280px',
            color: 'rgba(255, 255, 255, 0.85)',
            flex: 1,
            overflow: 'auto',
        }
    } else {
        return {
            background: '#fff',
            padding: padding,
            margin: 0,
            minHeight: '280px',
            flex: 1,
            overflow: 'auto',
        }
    }
})

// 获取用户信息（已迁移到store，保留此函数用于向后兼容或手动刷新）
const fetchUserInfo = async () => {
    await userStore.fetchUserInfo()
}

// 标签页管理
const handleTabChange = (key: string) => {
    activeTabKey.value = key
    const tab = pageTabs.value.find((t) => t.key === key)
    if (tab && tab.path) {
        router.push(tab.path)
    }
}

const handleTabClose = (key: string) => {
    const index = pageTabs.value.findIndex((t) => t.key === key)
    if (index > -1) {
        // 如果关闭的是当前标签，切换到前一个或后一个
        if (key === activeTabKey.value) {
            if (index > 0) {
                const prevTab = pageTabs.value[index - 1]
                if (prevTab) {
                    handleTabChange(prevTab.key)
                }
            } else if (pageTabs.value.length > 1) {
                const nextTab = pageTabs.value[index + 1]
                if (nextTab) {
                    handleTabChange(nextTab.key)
                }
            }
        }
        pageTabs.value.splice(index, 1)
        // 如果只剩下一个标签且不可关闭，确保它存在
        if (pageTabs.value.length === 0) {
            initHomeTab()
        }
    }
}

const handleTabRefresh = (key: string) => {
    // 刷新当前标签页的逻辑
    const tab = pageTabs.value.find((t) => t.key === key)
    if (tab) {
        router.push(tab.path)
    }
}

const handleCloseOther = (key: string) => {
    const keepTab = pageTabs.value.find((t) => t.key === key)
    if (keepTab) {
        pageTabs.value = [keepTab]
        activeTabKey.value = key
        if (keepTab.path) {
            router.push(keepTab.path)
        }
    }
}

const handleCloseAll = () => {
    initHomeTab()
    router.push('/')
}

// 监听路由变化，自动添加或更新标签页
watch(
    () => route.path,
    (newPath, oldPath) => {
        // 跳过登录页
        if (newPath === '/login') {
            return
        }

        // 查找是否已存在该路径的标签页
        const existingTab = pageTabs.value.find((t) => t.path === newPath)

        if (!existingTab) {
            // 不存在则添加新标签页
            const title = getRouteTitle(newPath, route.name)
            // 使用路径作为 key（避免重复路径创建多个标签）
            const pathKey = newPath.replace(/\//g, '-').replace(/^-|-$/g, '') || 'home'
            const newTab: TabItem = {
                key: pathKey,
                title: title,
                path: newPath,
                closable: newPath !== '/', // 首页不可关闭
            }
            
            // 如果首页不存在，确保首页存在
            const homeTab = pageTabs.value.find((t) => t.path === '/')
            if (!homeTab && newPath !== '/') {
                pageTabs.value.unshift({
                    key: 'home',
                    title: '工作台',
                    path: '/',
                    closable: false,
                })
            }
            
            pageTabs.value.push(newTab)
            activeTabKey.value = newTab.key
        } else {
            // 已存在则更新活动标签
            activeTabKey.value = existingTab.key
            // 更新标题（可能路由 meta 或名称有变化）
            const title = getRouteTitle(newPath, route.name)
            if (title !== existingTab.title) {
                existingTab.title = title
            }
        }
    },
    { immediate: true }
)

</script>

<style scoped>
.layout-container {
    min-height: 100vh;
    height: 100vh;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.layout-body {
    flex: 1;
    overflow: hidden;
    transition: background-color 0.3s;
    display: flex;
    min-height: 0;
    position: relative;
}

:deep(.ant-layout) {
    display: flex;
}

:deep(.ant-layout-sider) {
    flex-shrink: 0;
}

.main-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-height: 0;
}

.content-wrapper {
    margin: 20px 0;
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 0 24px;
    overflow: hidden;
    min-height: 0;
}

.desktop-sidebar {
    display: block;
}

.mobile-drawer {
    display: none;
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
    .desktop-sidebar {
        display: none !important;
    }

    .mobile-drawer {
        display: block;
    }

    .content-wrapper {
        padding: 0 16px;
    }

    .main-content {
        width: 100%;
    }
}

@media (max-width: 480px) {
    .content-wrapper {
        padding: 0 12px;
    }

    .layout-container {
        font-size: 14px;
    }
}

:deep(.ant-layout-content) {
    transition: background-color 0.3s, color 0.3s;
}

.dark-layout :deep(.ant-breadcrumb) {
    color: rgba(255, 255, 255, 0.65);
}

.dark-layout :deep(.ant-breadcrumb-link) {
    color: rgba(255, 255, 255, 0.85);
}

.layout-container:not(.dark-layout) .user-info:hover {
    background-color: rgba(0, 0, 0, 0.06);
}

.layout-container:not(.dark-layout) .theme-switch:hover {
    background-color: rgba(0, 0, 0, 0.06);
}

.footer-wrapper {
    padding: 0;
    flex-shrink: 0;
    margin: 0;
}

/* 页面切换动画 */
.page-fade-enter-active,
.page-fade-leave-active {
    transition: opacity 0.3s ease, transform 0.3s ease;
}

.page-fade-enter-from {
    opacity: 0;
    transform: translateX(20px);
}

.page-fade-leave-to {
    opacity: 0;
    transform: translateX(-20px);
}

.page-fade-enter-to,
.page-fade-leave-from {
    opacity: 1;
    transform: translateX(0);
}
</style>

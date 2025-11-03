<template>
    <a-layout-sider width="200" :style="siderStyle" class="sidebar-container">
        <a-spin :spinning="menuList.length === 0" style="width: 100%;">
            <a-menu
                v-model:selectedKeys="selectedKeys"
                v-model:openKeys="openKeys"
                mode="inline"
                :theme="isDark ? 'dark' : 'light'"
                :style="{ height: '100%', borderRight: 0 }"
                @click="handleMenuClick"
            >
                <!-- 有子菜单的情况 -->
                <a-sub-menu
                    v-for="(menu, index) in menusWithChildren"
                    :key="menu?.key ? `sub-${menu.key}` : `sub-menu-${index}`"
                >
                    <template #title>
                        <span>
                            <component v-if="menu?.icon" :is="getIcon(menu.icon)" />
                            {{ menu?.title || '' }}
                        </span>
                    </template>
                    <a-menu-item
                        v-for="(child, childIndex) in (menu.children || [])"
                        :key="child?.key || `child-${childIndex}`"
                        @click="handleMenuItemClick(child)"
                    >
                        {{ child?.title || '' }}
                    </a-menu-item>
                </a-sub-menu>
                <!-- 无子菜单的情况 -->
                <a-menu-item
                    v-for="(menu, index) in menusWithoutChildren"
                    :key="menu?.key ? `item-${menu.key}` : `item-menu-${index}`"
                    @click="handleMenuItemClick(menu)"
                >
                    <span v-if="menu?.icon">
                        <component :is="getIcon(menu.icon)" />
                    </span>
                    <span>{{ menu?.title || '' }}</span>
                </a-menu-item>
            </a-menu>
        </a-spin>
    </a-layout-sider>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
    UserOutlined,
    LaptopOutlined,
    NotificationOutlined,
    HomeOutlined,
    SettingOutlined,
    TeamOutlined,
    FileOutlined,
    DashboardOutlined,
} from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { menuApi, type MenuItem } from '@/api/menu'
import type { Component } from 'vue'

const themeStore = useThemeStore()
const router = useRouter()
const route = useRoute()
const isDark = computed(() => themeStore.isDark)

const emit = defineEmits<{
    close: []
}>()

// 初始化默认菜单，确保组件渲染时就有数据
const getInitialMenus = (): MenuItem[] => [
    {
        key: 'dashboard',
        title: '仪表盘',
        path: '/dashboard',
        icon: 'DashboardOutlined',
        children: [
            {
                key: 'workbench',
                title: '工作台',
                path: '/dashboard/workbench',
            },
            {
                key: 'analysis',
                title: '分析页',
                path: '/dashboard/analysis',
            },
        ],
    },
    {
        key: 'user',
        title: '用户管理',
        path: '/user/list',
        icon: 'TeamOutlined',
    },
    {
        key: 'system',
        title: '系统设置',
        path: '/system/config',
        icon: 'SettingOutlined',
    },
]

const menuList = ref<MenuItem[]>(getInitialMenus())
const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>(['dashboard'])

// 分离有子菜单和无子菜单的项，确保数据安全
const menusWithChildren = computed(() => {
    return menuList.value.filter((menu) => {
        return menu && menu.children && Array.isArray(menu.children) && menu.children.length > 0
    })
})

const menusWithoutChildren = computed(() => {
    return menuList.value.filter((menu) => {
        return menu && (!menu.children || !Array.isArray(menu.children) || menu.children.length === 0)
    })
})

// 图标映射
const iconMap: Record<string, Component> = {
    UserOutlined,
    LaptopOutlined,
    NotificationOutlined,
    HomeOutlined,
    SettingOutlined,
    TeamOutlined,
    FileOutlined,
    DashboardOutlined,
}

// 获取图标组件
const getIcon = (icon: Component | string): Component => {
    if (typeof icon === 'string') {
        return iconMap[icon] || HomeOutlined
    }
    return icon
}

// 获取菜单数据
const fetchMenus = async () => {
    try {
        const res = await menuApi.getSidebarMenus()
        if ((res.code === 200 || res.code === 0) && res.data && res.data.length > 0) {
            menuList.value = res.data
            // 默认展开第一个有子菜单的项
            const firstSubMenu = menuList.value.find((m) => m.children && m.children.length > 0)
            if (firstSubMenu) {
                openKeys.value = [firstSubMenu.key]
            }
            // 根据当前路由设置选中项
            updateSelectedKeys()
        }
        // 如果返回的数据为空，保持默认菜单
    } catch (error) {
        console.error('获取菜单失败，使用默认菜单:', error)
        // 失败时保持默认菜单（已在初始化时设置）
    }
}

// 设置默认菜单（API 失败时的降级方案）
const setDefaultMenus = () => {
    menuList.value = [
        {
            key: 'dashboard',
            title: '仪表盘',
            path: '/dashboard',
            icon: 'DashboardOutlined',
            children: [
                {
                    key: 'workbench',
                    title: '工作台',
                    path: '/dashboard/workbench',
                },
                {
                    key: 'analysis',
                    title: '分析页',
                    path: '/dashboard/analysis',
                },
            ],
        },
        {
            key: 'user',
            title: '用户管理',
            path: '/user/list',
            icon: 'TeamOutlined',
        },
        {
            key: 'system',
            title: '系统设置',
            path: '/system/config',
            icon: 'SettingOutlined',
        },
    ]
    openKeys.value = ['dashboard']
    updateSelectedKeys()
}

// 更新选中项（根据当前路由）
const updateSelectedKeys = () => {
    const findMenuKey = (menus: MenuItem[], path: string): string | null => {
        for (const menu of menus) {
            if (menu.path === path) {
                return menu.key
            }
            if (menu.children) {
                const found = findMenuKey(menu.children, path)
                if (found) return found
            }
        }
        return null
    }

    const currentPath = route.path
    const key = findMenuKey(menuList.value, currentPath)
    if (key) {
        selectedKeys.value = [key]
    }
}

const handleMenuClick = () => {
    // 移动端点击菜单项后自动关闭抽屉
    if (window.innerWidth <= 768) {
        emit('close')
    }
}

const handleMenuItemClick = (menu: MenuItem) => {
    if (menu.path && menu.path !== route.path) {
        router.push(menu.path)
    }
    handleMenuClick()
}

// 监听路由变化，更新选中项
watch(
    () => route.path,
    () => {
        updateSelectedKeys()
    },
    { immediate: true }
)

// 初始化时立即获取菜单
onMounted(async () => {
    // 尝试从 API 获取最新菜单（如果失败则保持默认菜单）
    await fetchMenus()
})

const siderStyle = computed(() => {
    if (isDark.value) {
        return {
            background: '#001529',
        }
    } else {
        return {
            background: '#fff',
        }
    }
})
</script>

<style scoped>
.sidebar-container {
    height: 100%;
}

:deep(.ant-layout-sider) {
    overflow-y: auto;
    overflow-x: hidden;
    transition: background-color 0.3s;
}


/* 移动端样式 */
@media (max-width: 768px) {
    .sidebar-container {
        height: 100vh;
    }
}
</style>


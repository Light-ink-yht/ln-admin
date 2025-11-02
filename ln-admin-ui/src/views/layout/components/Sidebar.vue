<template>
    <a-layout-sider width="200" :style="siderStyle" class="sidebar-container">
        <a-menu
            v-model:selectedKeys="selectedKeys"
            v-model:openKeys="openKeys"
            mode="inline"
            :theme="isDark ? 'dark' : 'light'"
            :style="{ height: '100%', borderRight: 0 }"
            @click="handleMenuClick"
        >
            <!-- 有子菜单的情况 -->
            <template v-for="menu in menuList" :key="menu.key">
                <a-sub-menu v-if="menu.children && menu.children.length > 0" :key="`sub-${menu.key}`">
                    <template #title>
                        <span>
                            <component v-if="menu.icon" :is="getIcon(menu.icon)" />
                            {{ menu.title }}
                        </span>
                    </template>
                    <a-menu-item
                        v-for="child in menu.children"
                        :key="child.key"
                        @click="handleMenuItemClick(child)"
                    >
                        {{ child.title }}
                    </a-menu-item>
                </a-sub-menu>
                <!-- 无子菜单的情况 -->
                <a-menu-item
                    v-else
                    :key="`item-${menu.key}`"
                    @click="handleMenuItemClick(menu)"
                >
                    <span v-if="menu.icon">
                        <component :is="getIcon(menu.icon)" />
                    </span>
                    <span>{{ menu.title }}</span>
                </a-menu-item>
            </template>
        </a-menu>
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

const menuList = ref<MenuItem[]>([])
const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>([])

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
        if (res.code === 200 || res.code === 0) {
            menuList.value = res.data || []
            // 默认展开第一个有子菜单的项
            const firstSubMenu = menuList.value.find((m) => m.children && m.children.length > 0)
            if (firstSubMenu) {
                openKeys.value = [firstSubMenu.key]
            }
            // 根据当前路由设置选中项
            updateSelectedKeys()
        }
    } catch (error) {
        console.error('获取菜单失败:', error)
        // 失败时使用默认菜单
        setDefaultMenus()
    }
}

// 设置默认菜单（API 失败时的降级方案）
const setDefaultMenus = () => {
    menuList.value = [
        {
            key: 'sub1',
            title: 'subnav 1',
            icon: 'UserOutlined',
            children: [
                { key: '1', title: 'option1', path: '/' },
                { key: '2', title: 'option2', path: '/' },
                { key: '3', title: 'option3', path: '/' },
                { key: '4', title: 'option4', path: '/' },
            ],
        },
        {
            key: 'sub2',
            title: 'subnav 2',
            icon: 'LaptopOutlined',
            children: [
                { key: '5', title: 'option5', path: '/' },
                { key: '6', title: 'option6', path: '/' },
                { key: '7', title: 'option7', path: '/' },
                { key: '8', title: 'option8', path: '/' },
            ],
        },
        {
            key: 'sub3',
            title: 'subnav 3',
            icon: 'NotificationOutlined',
            children: [
                { key: '9', title: 'option9', path: '/' },
                { key: '10', title: 'option10', path: '/' },
                { key: '11', title: 'option11', path: '/' },
                { key: '12', title: 'option12', path: '/' },
            ],
        },
    ]
    openKeys.value = ['sub1']
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

onMounted(() => {
    fetchMenus()
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


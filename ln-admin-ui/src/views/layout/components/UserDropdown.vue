<template>
    <a-dropdown v-if="currentUserInfo" placement="bottomRight" :trigger="['click']">
        <div class="user-info" :style="userInfoStyle">
            <a-avatar :src="currentUserInfo.avatar" :size="32">
                <template v-if="!currentUserInfo.avatar">
                    <UserOutlined />
                </template>
            </a-avatar>
            <span class="user-name">{{ currentUserInfo.nickname || currentUserInfo.fullName || '用户' }}</span>
            <DownOutlined />
        </div>
        <template #overlay>
            <a-menu @click="handleMenuClick">
                <template v-for="(item, index) in menuItems" :key="item.key">
                    <a-menu-divider v-if="item.divider && index > 0" />
                    <a-menu-item>
                        <component v-if="item.icon" :is="getIcon(item.icon)" />
                        <span>{{ item.label }}</span>
                    </a-menu-item>
                </template>
            </a-menu>
        </template>
    </a-dropdown>
    <div v-else class="user-info" :style="userInfoStyle">
        <a-avatar :size="32">
            <UserOutlined />
        </a-avatar>
        <span class="user-name">未登录</span>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, DownOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { useUserStore } from '@/stores/modules/user'
import { menuApi, type UserMenuItem } from '@/api/menu'
import type { UserInfo } from '@/api/auth'
import type { Component } from 'vue'

interface Props {
    userInfo: UserInfo | null
}

const props = defineProps<Props>()

const router = useRouter()
const themeStore = useThemeStore()
const userStore = useUserStore()
const isDark = computed(() => themeStore.isDark)

// 优先使用store中的用户信息
const currentUserInfo = computed(() => props.userInfo || userStore.userInfo)

const menuItems = ref<UserMenuItem[]>([])

// 图标映射
const iconMap: Record<string, Component> = {
    UserOutlined,
    SettingOutlined,
    LogoutOutlined,
}

// 获取图标组件
const getIcon = (icon: Component | string): Component => {
    if (typeof icon === 'string') {
        return iconMap[icon] || UserOutlined
    }
    return icon
}

// 获取用户菜单
const fetchUserMenus = async () => {
    try {
        const res = await menuApi.getUserMenus()
        if (res.code === 200 || res.code === 0) {
            menuItems.value = res.data || []
        }
    } catch (error) {
        console.error('获取用户菜单失败:', error)
        // 失败时使用默认菜单
        setDefaultMenus()
    }
}

// 设置默认菜单
const setDefaultMenus = () => {
    menuItems.value = [
        {
            key: 'profile',
            label: '个人资料',
            icon: 'UserOutlined',
            onClick: () => {
                router.push('/profile')
            },
        },
        {
            key: 'settings',
            label: '设置',
            icon: 'SettingOutlined',
            onClick: () => {
                router.push('/settings')
            },
        },
        {
            key: 'logout',
            label: '退出登录',
            icon: 'LogoutOutlined',
            divider: true,
            onClick: handleLogout,
        },
    ]
}

const handleLogout = async () => {
    try {
        // 调用退出登录API（可选）
        // await authApi.logout()
        
        // 清除store中的用户信息
        userStore.logout()
        
        // 跳转到登录页
        router.push('/login')
        message.success('已退出登录')
    } catch (error) {
        console.error('退出登录失败:', error)
        // 即使API调用失败，也清除本地信息
        userStore.logout()
        router.push('/login')
        message.error('退出登录失败')
    }
}

const handleMenuClick = ({ key }: { key: string }) => {
    const menuItem = menuItems.value.find((item) => item.key === key)
    if (menuItem && menuItem.onClick) {
        menuItem.onClick()
    }
}

const userInfoStyle = computed(() => {
    return {
        color: isDark.value ? '#fff' : 'rgba(0, 0, 0, 0.85)',
    }
})

onMounted(() => {
    fetchUserMenus()
})
</script>

<style scoped>
.user-info {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    padding: 6px 14px;
    border-radius: 8px;
    transition: all 0.3s;
    position: relative;
}

.user-info::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: 8px;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.1), transparent);
    opacity: 0;
    transition: opacity 0.3s;
}

.user-info:hover::before {
    opacity: 1;
}

.user-info:hover {
    background-color: rgba(255, 255, 255, 0.1);
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

:deep(.ant-avatar) {
    transition: transform 0.3s;
    border: 2px solid rgba(255, 255, 255, 0.2);
}

.user-info:hover :deep(.ant-avatar) {
    transform: scale(1.1);
}

.user-name {
    font-size: 14px;
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: color 0.3s;
}

/* 移动端优化 */
@media (max-width: 768px) {
    .user-name {
        max-width: 80px;
        font-size: 13px;
    }

    .user-info {
        padding: 2px 8px;
    }
}

@media (max-width: 480px) {
    .user-name {
        display: none;
    }

    .user-info {
        padding: 4px;
    }
}

:deep(.ant-dropdown-menu-item) {
    display: flex;
    align-items: center;
    gap: 8px;
}
</style>

<style>
/* 浅色模式下的悬停效果 */
.dark-layout .user-info:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.layout-container:not(.dark-layout) .user-info:hover {
    background-color: rgba(0, 0, 0, 0.06);
}
</style>


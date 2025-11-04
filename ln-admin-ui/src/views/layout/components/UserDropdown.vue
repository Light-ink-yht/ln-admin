<template>
    <a-dropdown v-if="currentUserInfo" placement="bottomRight" :trigger="['click']">
        <div class="user-info" :style="userInfoStyle">
            <a-avatar :src="currentUserInfo.avatar" :size="32">
                <template v-if="!currentUserInfo.avatar">
                    <UserOutlined />
                </template>
            </a-avatar>
            <span class="user-name">{{ currentUserInfo.nickname || currentUserInfo.fullName || '用户' }}</span>
            <DownOutlined class="dropdown-arrow" />
        </div>
        <template #overlay>
            <a-menu @click="handleMenuClick">
                <template v-for="processedItem in processedMenuItems" :key="processedItem.key">
                    <a-menu-divider v-if="processedItem.type === 'divider'" />
                    <a-menu-item v-else :key="processedItem.item!.key">
                        <component v-if="processedItem.item?.icon" :is="getIcon(processedItem.item.icon)" />
                        <span>{{ processedItem.item?.label }}</span>
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

// 处理菜单项，将 divider 和 menu-item 分开
const processedMenuItems = computed(() => {
    const result: Array<{ type: 'divider' | 'item'; item?: UserMenuItem; key: string }> = []
    menuItems.value.forEach((item, index) => {
        if (item.divider && index > 0) {
            result.push({ type: 'divider', key: `divider-${item.key}` })
        }
        result.push({ type: 'item', item, key: item.key })
    })
    return result
})

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
            // 过滤菜单，只保留个人资料和退出登录
            const allMenus = res.data || []
            const filteredMenus = allMenus.filter((item: UserMenuItem) => {
                const key = item.key || ''
                return key === 'profile' || key === 'logout'
            })
            
            // 为过滤后的菜单添加onClick处理
            menuItems.value = filteredMenus.map((item: UserMenuItem) => {
                if (!item.onClick) {
                    if (item.key === 'profile') {
                        item.onClick = () => {
                            router.push('/profile').catch(() => {})
                        }
                    } else if (item.key === 'logout') {
                        item.onClick = handleLogout
                    }
                }
                return item
            })
            
            // 如果过滤后没有菜单，使用默认菜单
            if (menuItems.value.length === 0) {
                setDefaultMenus()
            }
        } else {
            setDefaultMenus()
        }
    } catch (error) {
        console.error('获取用户菜单失败:', error)
        // 失败时使用默认菜单
        setDefaultMenus()
    }
}

// 设置默认菜单（只保留个人资料和退出登录）
const setDefaultMenus = () => {
    menuItems.value = [
        {
            key: 'profile',
            label: '个人资料',
            icon: 'UserOutlined',
            onClick: () => {
                router.push('/profile').catch(() => {})
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
    console.log('菜单点击:', key, menuItems.value)
    const menuItem = menuItems.value.find((item) => item.key === key)
    if (menuItem) {
        if (menuItem.onClick) {
            menuItem.onClick()
        } else {
            // 如果没有onClick，根据key执行默认操作
            if (key === 'profile') {
                router.push('/profile')
            } else if (key === 'logout') {
                handleLogout()
            }
        }
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

.dropdown-arrow {
    font-size: 12px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    opacity: 0.6;
    color: inherit;
    margin-left: 4px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    position: relative;
}

.user-info:hover .dropdown-arrow {
    opacity: 1;
    transform: translateY(2px) scale(1.1);
    color: var(--ant-primary-color, #1890ff);
}

/* 添加一个微妙的背景光晕效果 */
.dropdown-arrow::before {
    content: '';
    position: absolute;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(24, 144, 255, 0.15) 0%, transparent 70%);
    opacity: 0;
    transition: opacity 0.3s;
    transform: translate(-50%, -50%);
    top: 50%;
    left: 50%;
}

.user-info:hover .dropdown-arrow::before {
    opacity: 1;
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


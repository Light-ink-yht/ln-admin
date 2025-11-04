<template>
    <a-layout-header class="header" :style="headerStyle">
        <div class="header-left">
            <!-- 移动端菜单按钮 -->
            <a-button
                type="text"
                class="mobile-menu-btn"
                @click="$emit('toggle-sidebar')"
            >
                <template #icon>
                    <MenuOutlined />
                </template>
            </a-button>
            <div class="logo" :style="logoStyle">
                <img v-if="siteLogo" :src="siteLogo" alt="Logo" class="logo-img" />
                <span>{{ siteName }}</span>
            </div>
        </div>
        <!-- 顶部中间区域：面包屑或自定义菜单 -->
        <div class="header-center">
            <!-- 默认显示面包屑 -->
            <HeaderBreadcrumb v-if="!$slots.menu" />
            <!-- 插槽：允许某些页面插入自定义菜单 -->
            <slot name="menu"></slot>
        </div>
        <!-- 用户信息和主题切换（右上角） -->
        <div class="header-right">
            <ThemeSwitch />
            <UserDropdown :user-info="userInfo" />
        </div>
    </a-layout-header>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { MenuOutlined } from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { useSystemConfigStore } from '@/stores/modules/systemConfig'
import { useRoute } from 'vue-router'
import ThemeSwitch from './ThemeSwitch.vue'
import UserDropdown from './UserDropdown.vue'
// @ts-ignore: if type error for missing type, comment out for now and check path
import HeaderBreadcrumb from './HeaderBreadcrumb.vue'
import type { UserInfo } from '@/api/auth'

defineEmits<{
    'toggle-sidebar': []
}>()

interface Props {
    userInfo: UserInfo | null
}

const props = defineProps<Props>()

const themeStore = useThemeStore()
const systemConfigStore = useSystemConfigStore()
const route = useRoute()
const isDark = computed(() => themeStore.isDark)
const colorPrimary = computed(() => themeStore.colorPrimary)
const siteName = computed(() => systemConfigStore.siteName)
const siteLogo = computed(() => systemConfigStore.siteLogo)

const headerStyle = computed(() => {
    if (isDark.value) {
        return {
            background: 'linear-gradient(135deg, #001529 0%, #002140 100%)',
            borderBottom: '1px solid rgba(255, 255, 255, 0.08)',
            boxShadow: '0 2px 12px rgba(0, 0, 0, 0.3)',
        }
    } else {
        return {
            background: 'linear-gradient(135deg, #ffffff 0%, #fafafa 100%)',
            borderBottom: '1px solid rgba(0, 0, 0, 0.06)',
            boxShadow: '0 2px 12px rgba(0, 0, 0, 0.04)',
        }
    }
})

const logoStyle = computed(() => {
    return {
        color: isDark.value ? '#fff' : colorPrimary.value,
    }
})
</script>

<style scoped>
.header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0;
    padding-right: 24px; /* 只保留右侧 padding */
    height: 64px;
    position: relative;
    z-index: 10;
    transition: background-color 0.3s, border-color 0.3s;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 24px;
    flex: 0 0 auto;
    width: 200px; /* 固定宽度，与 sidebar 宽度一致，确保面包屑对齐 */
    padding-left: 24px; /* 添加左侧 padding */
    box-sizing: border-box;
}

.header-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 0 0 auto;
    padding-right: 8px;
}

.logo {
    font-size: 20px;
    font-weight: 600;
    white-space: nowrap;
    display: flex;
    align-items: center;
    gap: 10px;
    height: 64px;
    transition: all 0.3s;
    letter-spacing: 0.5px;
    position: relative;
}

.logo::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 0;
    height: 2px;
    background: linear-gradient(90deg, var(--ant-primary-color), transparent);
    transition: width 0.3s;
}

.logo:hover::after {
    width: 100%;
}

.logo-img {
    height: 36px;
    width: auto;
    object-fit: contain;
    border-radius: 4px;
    transition: transform 0.3s;
}

.logo:hover .logo-img {
    transform: scale(1.05);
}

.header-center {
    flex: 1;
    display: flex;
    align-items: center;
    margin: 0;
    min-width: 0; /* 允许 flex 项目缩小 */
    box-sizing: border-box;
}

.mobile-menu-btn {
    display: none;
    color: inherit;
    font-size: 18px;
    padding: 0 8px;
    margin-right: 8px;
}

.mobile-menu-btn:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.dark-layout .mobile-menu-btn:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.layout-container:not(.dark-layout) .mobile-menu-btn:hover {
    background-color: rgba(0, 0, 0, 0.06);
}

/* 移动端响应式 */
@media (max-width: 768px) {
    .mobile-menu-btn {
        display: flex;
        align-items: center;
    }

    .header-center {
        display: none; /* 移动端隐藏面包屑，HeaderBreadcrumb 组件内部也会隐藏 */
    }

    .header-left {
        margin-right: 0;
    }

    .header {
        padding: 0 16px;
    }

    .logo {
        font-size: 16px;
    }
}

@media (max-width: 480px) {
    .header {
        padding: 0 12px;
        height: 56px;
    }

    .logo {
        font-size: 14px;
        height: 56px;
    }
}
</style>


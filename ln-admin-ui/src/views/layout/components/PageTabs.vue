<template>
    <div class="page-tabs" :class="{ 'dark-tabs': isDark }" :style="tabsContainerStyle">
        <div class="tabs-wrapper" ref="tabsWrapperRef">
            <div
                v-for="tab in tabs"
                :key="tab.key"
                class="tab-item"
                :class="{ active: activeKey === tab.key }"
                :style="getTabStyle(tab)"
                @click="handleTabClick(tab)"
                @contextmenu.prevent="handleContextMenu($event, tab)"
            >
                <span class="tab-icon" v-if="tab.icon">
                    <component :is="tab.icon" />
                </span>
                <span class="tab-title">{{ tab.title }}</span>
                <span
                    class="tab-close"
                    v-if="tabs.length > 1 && tab.closable !== false"
                    @click.stop="handleClose(tab.key)"
                >
                    <CloseOutlined />
                </span>
            </div>
        </div>
        <div class="tabs-actions">
            <a-dropdown :trigger="['click']" placement="bottomRight">
                <a-button type="text" size="small" class="action-btn">
                    <template #icon>
                        <MoreOutlined />
                    </template>
                </a-button>
                <template #overlay>
                    <a-menu @click="handleMenuClick">
                        <a-menu-item key="closeOther">
                            <CloseCircleOutlined />
                            <span>关闭其他</span>
                        </a-menu-item>
                        <a-menu-item key="closeAll">
                            <DeleteOutlined />
                            <span>关闭全部</span>
                        </a-menu-item>
                        <a-menu-item key="refresh">
                            <ReloadOutlined />
                            <span>刷新当前</span>
                        </a-menu-item>
                    </a-menu>
                </template>
            </a-dropdown>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
    CloseOutlined,
    MoreOutlined,
    CloseCircleOutlined,
    DeleteOutlined,
    ReloadOutlined,
} from '@ant-design/icons-vue'
import { useThemeStore } from '@/stores/modules/theme'
import type { Component } from 'vue'

export interface TabItem {
    key: string
    title: string
    path: string
    icon?: Component
    closable?: boolean
}

interface Props {
    tabs?: TabItem[]
    activeKey?: string
}

const props = withDefaults(defineProps<Props>(), {
    tabs: () => [],
    activeKey: '',
})

const emit = defineEmits<{
    'update:activeKey': [key: string]
    'tab-close': [key: string]
    'tab-refresh': [key: string]
    'close-other': [key: string]
    'close-all': []
}>()

const router = useRouter()
const route = useRoute()
const themeStore = useThemeStore()
const tabsWrapperRef = ref<HTMLElement>()

const isDark = computed(() => themeStore.isDark)
const colorPrimary = computed(() => themeStore.colorPrimary)

const tabsContainerStyle = computed(() => {
    if (isDark.value) {
        return {
            background: 'linear-gradient(180deg, #1f1f1f 0%, #1a1a1a 100%)',
            borderBottom: '1px solid rgba(255, 255, 255, 0.08)',
        }
    } else {
        return {
            background: 'linear-gradient(180deg, #ffffff 0%, #fafafa 100%)',
            borderBottom: '1px solid rgba(0, 0, 0, 0.06)',
        }
    }
})

const getTabStyle = (tab: TabItem) => {
    const isActive = props.activeKey === tab.key
    if (isActive) {
        return {
            color: colorPrimary.value,
        }
    }
    return {}
}

const handleTabClick = (tab: TabItem) => {
    emit('update:activeKey', tab.key)
    if (tab.path && tab.path !== route.path) {
        router.push(tab.path)
    }
}

const handleClose = (key: string) => {
    const tab = props.tabs.find((t) => t.key === key)
    // 如果标签不可关闭，则不允许关闭
    if (tab && tab.closable === false) {
        return
    }
    emit('tab-close', key)
}

const handleContextMenu = (e: MouseEvent, tab: TabItem) => {
    // 可以在这里添加右键菜单功能
    console.log('右键菜单', tab)
}

const handleMenuClick = ({ key }: { key: string }) => {
    switch (key) {
        case 'closeOther':
            emit('close-other', props.activeKey)
            break
        case 'closeAll':
            emit('close-all')
            break
        case 'refresh':
            emit('tab-refresh', props.activeKey)
            break
    }
}

// 自动滚动到活动标签
const scrollToActiveTab = () => {
    if (tabsWrapperRef.value) {
        const activeTab = tabsWrapperRef.value.querySelector('.tab-item.active') as HTMLElement
        if (activeTab && tabsWrapperRef.value) {
            const container = tabsWrapperRef.value
            const containerRect = container.getBoundingClientRect()
            const tabRect = activeTab.getBoundingClientRect()
            
            // 检查标签是否在可视区域内
            const isVisible = 
                tabRect.left >= containerRect.left &&
                tabRect.right <= containerRect.right
            
            if (!isVisible) {
                // 如果标签不可见，滚动到中间位置
                activeTab.scrollIntoView({
                    behavior: 'smooth',
                    block: 'nearest',
                    inline: 'center',
                })
            }
        }
    }
}

watch(
    () => props.activeKey,
    async () => {
        // 使用 nextTick 确保 DOM 更新后再滚动
        await nextTick()
        scrollToActiveTab()
    },
    { immediate: true }
)

// 监听 tabs 变化，当 tabs 增加时也滚动
watch(
    () => props.tabs.length,
    async () => {
        await nextTick()
        scrollToActiveTab()
    }
)
</script>

<style scoped>
.page-tabs {
    display: flex;
    align-items: center;
    height: 44px;
    min-height: 44px;
    padding: 0 24px;
    overflow: hidden;
    transition: background-color 0.3s, border-color 0.3s;
    flex-shrink: 0;
    width: 100%;
    box-sizing: border-box;
    position: relative;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.dark-tabs {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

/* 移动端优化 */
@media (max-width: 768px) {
    .page-tabs {
        padding: 0 16px;
        height: 40px;
        min-height: 40px;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
    }

    .tab-item {
        min-width: 60px;
        max-width: 120px;
        padding: 6px 10px;
        font-size: 12px;
        height: 32px;
        margin-right: 2px;
    }

    .tab-title {
        font-size: 12px;
    }

    .tab-close {
        font-size: 10px;
        padding: 1px;
    }

    .tabs-actions {
        padding-left: 4px;
    }
}

@media (max-width: 480px) {
    .page-tabs {
        padding: 0 12px;
        height: 36px;
        min-height: 36px;
    }

    .tab-item {
        min-width: 50px;
        max-width: 100px;
        padding: 4px 8px;
        font-size: 11px;
        height: 28px;
        margin-right: 2px;
    }

    .tab-icon {
        display: none;
    }
}

.tabs-wrapper {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 2px;
    overflow-x: auto;
    overflow-y: hidden;
    scrollbar-width: thin;
    scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
    padding: 4px 0;
    min-width: 0; /* 允许 flex 项目缩小 */
    scroll-behavior: smooth; /* 平滑滚动 */
    -webkit-overflow-scrolling: touch; /* iOS 平滑滚动 */
}

.tabs-wrapper::-webkit-scrollbar {
    height: 8px; /* 增大滚动条高度，更容易看到和操作 */
}

.tabs-wrapper::-webkit-scrollbar-track {
    background: transparent;
    margin: 2px 0; /* 添加一些边距 */
}

.tabs-wrapper::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.25);
    border-radius: 4px;
    transition: background 0.2s;
}

.tabs-wrapper::-webkit-scrollbar-thumb:hover {
    background: rgba(0, 0, 0, 0.4);
}

.dark-tabs .tabs-wrapper::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.3);
}

.dark-tabs .tabs-wrapper::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.45);
}

/* 当需要滚动时显示滚动条 */
.tabs-wrapper:not(:hover)::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.15);
}

.dark-tabs .tabs-wrapper:not(:hover)::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
}

.tab-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    min-width: 80px;
    max-width: 200px;
    width: auto; /* 自动宽度，根据内容 */
    height: 36px;
    border-radius: 6px 6px 0 0;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    white-space: nowrap;
    font-size: 13px;
    color: rgba(0, 0, 0, 0.65);
    position: relative;
    user-select: none;
    margin-right: 4px;
    background: transparent;
    flex-shrink: 0; /* 防止标签被压缩 */
}

.tab-item::before {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: transparent;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 2px 2px 0 0;
}

.dark-tabs .tab-item {
    color: rgba(255, 255, 255, 0.65);
}

.tab-item:hover {
    background: rgba(0, 0, 0, 0.03);
    transform: translateY(-1px);
}

.dark-tabs .tab-item:hover {
    background: rgba(255, 255, 255, 0.06);
}

.tab-item.active {
    font-weight: 600;
    background: rgba(24, 144, 255, 0.08);
    box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

.tab-item.active::before {
    background: v-bind('colorPrimary');
    box-shadow: 0 0 8px v-bind('colorPrimary + "40"');
}

.dark-tabs .tab-item.active {
    background: rgba(255, 255, 255, 0.08);
    box-shadow: 0 2px 8px rgba(255, 255, 255, 0.1);
}

.dark-tabs .tab-item.active::before {
    box-shadow: 0 0 8px v-bind('colorPrimary + "60"');
}

.tab-item.active:hover {
    background: rgba(24, 144, 255, 0.12);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(24, 144, 255, 0.2);
}

.dark-tabs .tab-item.active:hover {
    background: rgba(255, 255, 255, 0.1);
    box-shadow: 0 4px 12px rgba(255, 255, 255, 0.15);
}

.tab-icon {
    display: flex;
    align-items: center;
    font-size: 14px;
    transition: transform 0.2s;
}

.tab-item:hover .tab-icon {
    transform: scale(1.1);
}

.tab-item.active .tab-icon {
    transform: scale(1.05);
}

.tab-title {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.tab-close {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 6px;
    padding: 2px 4px;
    border-radius: 4px;
    font-size: 12px;
    opacity: 0.5;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    width: 18px;
    height: 18px;
    flex-shrink: 0;
}

.tab-close:hover {
    opacity: 1;
    background: rgba(255, 77, 79, 0.15);
    color: #ff4d4f;
    transform: scale(1.1) rotate(90deg);
}

.dark-tabs .tab-close:hover {
    background: rgba(255, 77, 79, 0.25);
}

.tabs-actions {
    display: flex;
    align-items: center;
    padding-left: 12px;
    margin-left: 8px;
    border-left: 1px solid rgba(0, 0, 0, 0.08);
    flex-shrink: 0;
}

.dark-tabs .tabs-actions {
    border-left-color: rgba(255, 255, 255, 0.12);
}

.action-btn {
    color: inherit;
    transition: all 0.2s;
}

.action-btn:hover {
    background: rgba(0, 0, 0, 0.04);
    transform: scale(1.05);
}

.dark-tabs .action-btn:hover {
    background: rgba(255, 255, 255, 0.08);
}

:deep(.ant-dropdown-menu-item) {
    display: flex;
    align-items: center;
    gap: 8px;
}
</style>


<template>
    <a-card title="快捷链接" class="section-card" :bordered="false">
        <div class="quick-links-container">
            <a-spin :spinning="loading">
                <div class="links-list">
                    <transition-group name="link-list" tag="div">
                        <div
                            v-for="link in linksList"
                            :key="link.key"
                            class="link-item"
                            @click="handleLinkClick(link)"
                        >
                            <div class="link-icon-wrapper">
                                <div
                                    class="link-icon"
                                    :style="{ background: link.color || 'var(--color-primary)' }"
                                >
                                    <component :is="getIconComponent(link.icon)" />
                                </div>
                            </div>
                            <div class="link-content">
                                <div class="link-title">{{ link.title }}</div>
                                <div class="link-desc">{{ link.desc }}</div>
                            </div>
                            <div class="link-arrow">
                                <RightOutlined />
                            </div>
                        </div>
                    </transition-group>
                    <a-empty v-if="linksList.length === 0" description="暂无链接" :image="false" />
                </div>
            </a-spin>
        </div>
    </a-card>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
    RightOutlined,
    UserOutlined,
    FileTextOutlined,
    SettingOutlined,
    BarChartOutlined,
    DatabaseOutlined,
    SafetyOutlined,
    TeamOutlined,
    BellOutlined,
    AuditOutlined,
    UnlockOutlined,
    HomeOutlined,
} from '@ant-design/icons-vue'
import type { Component } from 'vue'

export interface QuickLink {
    key: string
    title: string
    desc: string
    icon: string
    color?: string
    path?: string
    url?: string
}

interface Props {
    links?: QuickLink[]
}

const props = withDefaults(defineProps<Props>(), {
    links: undefined,
})

const router = useRouter()
const loading = ref(false)

const iconMap: Record<string, Component> = {
    UserOutlined,
    FileTextOutlined,
    SettingOutlined,
    BarChartOutlined,
    DatabaseOutlined,
    SafetyOutlined,
    TeamOutlined,
    BellOutlined,
    AuditOutlined,
    UnlockOutlined,
    HomeOutlined,
}

const getIconComponent = (iconName: string): Component => {
    return iconMap[iconName] || HomeOutlined
}

const defaultLinks = ref<QuickLink[]>([
    {
        key: 'profile',
        title: '个人中心',
        desc: '查看个人信息',
        icon: 'UserOutlined',
        color: '#1890ff',
        path: '/profile',
    },
    {
        key: 'settings',
        title: '系统设置',
        desc: '系统参数配置',
        icon: 'SettingOutlined',
        color: '#722ed1',
        path: '/system/config',
    },
    {
        key: 'logs',
        title: '操作日志',
        desc: '查看操作记录',
        icon: 'AuditOutlined',
        color: '#faad14',
        path: '/system/log',
    },
    {
        key: 'analysis',
        title: '数据分析',
        desc: '查看数据报表',
        icon: 'BarChartOutlined',
        color: '#52c41a',
        path: '/dashboard/analysis',
    },
])

const linksList = computed(() => {
    return props.links && props.links.length > 0 ? props.links : defaultLinks.value
})

const handleLinkClick = (link: QuickLink) => {
    if (link.path) {
        router.push(link.path)
    } else if (link.url) {
        window.open(link.url, '_blank')
    }
}
</script>

<style scoped lang="less">
.section-card {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.3s;

    :deep(.ant-card-head) {
        border-bottom: 1px solid #f0f0f0;
    }

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }
}

.quick-links-container {
    min-height: 100px;
}

.links-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.link-item {
    display: flex;
    align-items: center;
    padding: 12px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    background: #fff;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(var(--color-primary-rgb, 24, 144, 255), 0.15);
        transform: translateX(4px);

        .link-arrow {
            opacity: 1;
            transform: translateX(0);
        }
    }
}

.link-icon-wrapper {
    flex-shrink: 0;
    margin-right: 12px;
}

.link-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    color: #fff;
    transition: all 0.3s;
}

.link-content {
    flex: 1;
    min-width: 0;
}

.link-title {
    font-size: 14px;
    font-weight: 500;
    color: #262626;
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.link-desc {
    font-size: 12px;
    color: #8c8c8c;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.link-arrow {
    flex-shrink: 0;
    margin-left: 8px;
    color: #bfbfbf;
    font-size: 14px;
    opacity: 0;
    transform: translateX(-8px);
    transition: all 0.3s;
}

// 列表动画
.link-list-enter-active,
.link-list-leave-active {
    transition: all 0.3s ease;
}

.link-list-enter-from {
    opacity: 0;
    transform: translateX(-10px);
}

.link-list-leave-to {
    opacity: 0;
    transform: translateX(10px);
}

// 响应式设计
@media (max-width: 768px) {
    .link-item {
        padding: 10px;
    }

    .link-icon {
        width: 32px;
        height: 32px;
        font-size: 16px;
    }

    .link-title {
        font-size: 13px;
    }

    .link-desc {
        font-size: 11px;
    }
}

// 暗色模式适配
body.dark-mode {
    .section-card {
        background: #1f1f1f;
        border-color: #303030;

        :deep(.ant-card-head) {
            border-color: #303030;
        }
    }

    .link-item {
        background: #1f1f1f;
        border-color: #303030;

        &:hover {
            border-color: var(--color-primary);
            background: #262626;
        }
    }

    .link-title {
        color: rgba(255, 255, 255, 0.85);
    }

    .link-desc {
        color: rgba(255, 255, 255, 0.65);
    }
}
</style>


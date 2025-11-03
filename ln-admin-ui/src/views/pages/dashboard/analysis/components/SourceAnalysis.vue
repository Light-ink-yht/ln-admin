<template>
    <a-card title="访问来源" class="source-analysis-card" :bordered="false">
        <div class="source-list">
            <div v-for="item in sourceData" :key="item.key" class="source-item">
                <div class="source-header">
                    <div class="source-info">
                        <span class="source-icon" :style="{ background: item.color }">
                            <component :is="getIconComponent(item.icon)" />
                        </span>
                        <div class="source-details">
                            <div class="source-name">{{ item.name }}</div>
                            <div class="source-value">{{ item.value }}</div>
                        </div>
                    </div>
                    <div class="source-percent">{{ item.percent }}%</div>
                </div>
                <a-progress
                    :percent="item.percent"
                    :stroke-color="item.color"
                    :show-info="false"
                    :stroke-width="8"
                />
            </div>
        </div>
    </a-card>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import {
    ChromeOutlined,
    GlobalOutlined,
    LinkOutlined,
    SearchOutlined,
} from '@ant-design/icons-vue'
import type { Component } from 'vue'

export interface SourceItem {
    key: string
    name: string
    value: string
    percent: number
    icon: string
    color: string
}

interface Props {
    sourceData?: SourceItem[]
}

const props = withDefaults(defineProps<Props>(), {
    sourceData: undefined,
})

const iconMap: Record<string, Component> = {
    ChromeOutlined,
    GlobalOutlined,
    LinkOutlined,
    SearchOutlined,
}

const getIconComponent = (iconName: string): Component => {
    return iconMap[iconName] || GlobalOutlined
}

const defaultSourceData: SourceItem[] = [
    {
        key: 'direct',
        name: '直接访问',
        value: '12,345',
        percent: 45,
        icon: 'LinkOutlined',
        color: 'var(--color-primary)',
    },
    {
        key: 'search',
        name: '搜索引擎',
        value: '8,901',
        percent: 32,
        icon: 'SearchOutlined',
        color: '#52c41a',
    },
    {
        key: 'referral',
        name: '外部链接',
        value: '4,567',
        percent: 17,
        icon: 'GlobalOutlined',
        color: '#faad14',
    },
    {
        key: 'social',
        name: '社交媒体',
        value: '2,034',
        percent: 6,
        icon: 'ChromeOutlined',
        color: '#f5222d',
    },
]

const sourceData = computed(() => {
    return props.sourceData && props.sourceData.length > 0 ? props.sourceData : defaultSourceData
})
</script>

<style scoped lang="less">
.source-analysis-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }
}

.source-list {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.source-item {
    .source-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
    }

    .source-info {
        display: flex;
        align-items: center;
        gap: 12px;
        flex: 1;
    }

    .source-icon {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 20px;
        color: #fff;
        flex-shrink: 0;
    }

    .source-details {
        flex: 1;
        min-width: 0;
    }

    .source-name {
        font-size: 14px;
        color: #262626;
        margin-bottom: 4px;
    }

    .source-value {
        font-size: 12px;
        color: #8c8c8c;
    }

    .source-percent {
        font-size: 18px;
        font-weight: 600;
        color: #262626;
    }
}

// 暗色模式
body.dark-mode {
    .source-analysis-card {
        background: #1f1f1f;
        border-color: #303030;
    }

    .source-name {
        color: rgba(255, 255, 255, 0.85);
    }

    .source-value {
        color: rgba(255, 255, 255, 0.65);
    }

    .source-percent {
        color: rgba(255, 255, 255, 0.85);
    }
}
</style>


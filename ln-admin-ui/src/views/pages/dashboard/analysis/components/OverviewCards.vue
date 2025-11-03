<template>
    <a-row :gutter="[16, 16]" class="overview-cards">
        <a-col :xs="24" :sm="12" :lg="6" v-for="card in cards" :key="card.key">
            <a-card class="overview-card" :class="card.type">
                <div class="card-icon" :style="{ background: card.color }">
                    <component :is="getIconComponent(card.icon)" />
                </div>
                <div class="card-content">
                    <div class="card-title">{{ card.title }}</div>
                    <div class="card-value">{{ card.value }}</div>
                    <div class="card-change" :class="card.trend">
                        <component :is="card.trend === 'up' ? ArrowUpOutlined : ArrowDownOutlined" />
                        {{ card.change }}
                    </div>
                </div>
            </a-card>
        </a-col>
    </a-row>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import {
    UserOutlined,
    EyeOutlined,
    ShoppingCartOutlined,
    DollarOutlined,
    ArrowUpOutlined,
    ArrowDownOutlined,
} from '@ant-design/icons-vue'
import type { Component } from 'vue'

export interface OverviewCard {
    key: string
    title: string
    value: string | number
    icon: string
    color: string
    trend: 'up' | 'down'
    change: string
    type: 'primary' | 'success' | 'warning' | 'danger'
}

interface Props {
    cards?: OverviewCard[]
}

const props = withDefaults(defineProps<Props>(), {
    cards: undefined,
})

const iconMap: Record<string, Component> = {
    UserOutlined,
    EyeOutlined,
    ShoppingCartOutlined,
    DollarOutlined,
}

const getIconComponent = (iconName: string): Component => {
    return iconMap[iconName] || UserOutlined
}

const defaultCards: OverviewCard[] = [
    {
        key: 'users',
        title: '总访问量',
        value: '12,345',
        icon: 'UserOutlined',
        color: 'var(--color-primary)',
        trend: 'up' as const,
        change: '12.5%',
        type: 'primary' as const,
    },
    {
        key: 'views',
        title: '页面浏览量',
        value: '45,678',
        icon: 'EyeOutlined',
        color: '#52c41a',
        trend: 'up' as const,
        change: '8.2%',
        type: 'success' as const,
    },
    {
        key: 'orders',
        title: '订单数量',
        value: '1,234',
        icon: 'ShoppingCartOutlined',
        color: '#faad14',
        trend: 'down' as const,
        change: '3.1%',
        type: 'warning' as const,
    },
    {
        key: 'revenue',
        title: '总收入',
        value: '¥89,012',
        icon: 'DollarOutlined',
        color: '#f5222d',
        trend: 'up' as const,
        change: '15.8%',
        type: 'danger' as const,
    },
]

const cards = computed(() => {
    return props.cards && props.cards.length > 0 ? props.cards : defaultCards
})
</script>

<style scoped lang="less">
.overview-cards {
    margin-bottom: 24px;
}

.overview-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.3s;
    position: relative;
    overflow: hidden;

    &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.12);
    }

    &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 4px;
        height: 100%;
        background: var(--color-primary);
        transition: width 0.3s;
    }

    &.primary::before {
        background: var(--color-primary);
    }

    &.success::before {
        background: #52c41a;
    }

    &.warning::before {
        background: #faad14;
    }

    &.danger::before {
        background: #f5222d;
    }

    :deep(.ant-card-body) {
        padding: 20px;
        display: flex;
        align-items: center;
        gap: 16px;
    }
}

.card-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
    color: #fff;
    flex-shrink: 0;
}

.card-content {
    flex: 1;
    min-width: 0;
}

.card-title {
    font-size: 14px;
    color: #8c8c8c;
    margin-bottom: 8px;
}

.card-value {
    font-size: 24px;
    font-weight: 600;
    color: #262626;
    margin-bottom: 8px;
}

.card-change {
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: 4px;

    &.up {
        color: #52c41a;
    }

    &.down {
        color: #f5222d;
    }
}

// 暗色模式
body.dark-mode {
    .overview-card {
        background: #1f1f1f;
        border-color: #303030;
    }

    .card-title {
        color: rgba(255, 255, 255, 0.65);
    }

    .card-value {
        color: rgba(255, 255, 255, 0.85);
    }
}
</style>


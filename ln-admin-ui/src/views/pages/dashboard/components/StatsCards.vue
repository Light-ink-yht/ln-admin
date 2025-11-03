<template>
    <a-row :gutter="[16, 16]" class="stats-row" v-if="stats.length > 0">
        <a-col :xs="24" :sm="12" :md="6" v-for="stat in stats" :key="stat.key">
            <a-card class="stat-card" :class="stat.type">
                <a-statistic
                    :title="stat.title"
                    :value="stat.value"
                    :value-style="{ color: stat.color }"
                >
                    <template #prefix>
                        <component :is="getIconComponent(stat.icon)" :style="{ fontSize: '24px', marginRight: '8px' }" />
                    </template>
                    <template #suffix>
                        <span v-if="stat.suffix">{{ stat.suffix }}</span>
                    </template>
                </a-statistic>
                <div class="stat-footer">
                    <span class="stat-change" :class="stat.trend">
                        <component :is="stat.trend === 'up' ? ArrowUpOutlined : ArrowDownOutlined" />
                        {{ stat.change }}
                    </span>
                    <span class="stat-desc">{{ stat.desc }}</span>
                </div>
            </a-card>
        </a-col>
    </a-row>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import {
    TeamOutlined,
    SafetyOutlined,
    FileTextOutlined,
    UnlockOutlined,
    ArrowUpOutlined,
    ArrowDownOutlined,
    BellOutlined,
    DatabaseOutlined,
    AuditOutlined,
} from '@ant-design/icons-vue'
import type { Component } from 'vue'

export interface StatItem {
    key: string
    title: string
    value: number
    icon: Component | string
    color: string
    trend: 'up' | 'down'
    change: string
    desc: string
    type: string
    suffix?: string
}

// 图标映射
const iconMap: Record<string, Component> = {
    TeamOutlined,
    SafetyOutlined,
    FileTextOutlined,
    UnlockOutlined,
    BellOutlined,
    DatabaseOutlined,
    AuditOutlined,
}

const getIconComponent = (icon: Component | string): Component => {
    if (typeof icon === 'string') {
        return iconMap[icon] || FileTextOutlined
    }
    return icon
}

const props = defineProps<{
    stats?: StatItem[]
}>()

// 默认统计数据（当没有传入数据时使用）
const defaultStats = ref<StatItem[]>([
    {
        key: 'users',
        title: '用户总数',
        value: 0,
        icon: 'TeamOutlined',
        color: '#1890ff',
        trend: 'up',
        change: '12%',
        desc: '较上月',
        type: 'primary',
        suffix: '人',
    },
    {
        key: 'roles',
        title: '角色数量',
        value: 0,
        icon: 'SafetyOutlined',
        color: '#52c41a',
        trend: 'up',
        change: '5%',
        desc: '较上月',
        type: 'success',
        suffix: '个',
    },
    {
        key: 'permissions',
        title: '权限数量',
        value: 0,
        icon: 'UnlockOutlined',
        color: '#faad14',
        trend: 'down',
        change: '2%',
        desc: '较上月',
        type: 'warning',
        suffix: '个',
    },
    {
        key: 'logs',
        title: '今日操作',
        value: 0,
        icon: 'FileTextOutlined',
        color: '#f5222d',
        trend: 'up',
        change: '23%',
        desc: '较昨日',
        type: 'danger',
        suffix: '次',
    },
])

// 使用传入的stats或默认stats
const stats = computed(() => {
    if (props.stats && props.stats.length > 0) {
        return props.stats
    }
    return defaultStats.value
})

// 模拟加载统计数据（仅在无外部数据时使用）
const loadStats = async () => {
    // 只有在使用默认数据时才加载
    if (!props.stats || props.stats.length === 0) {
        defaultStats.value.forEach((stat) => {
            switch (stat.key) {
                case 'users':
                    stat.value = 1024
                    break
                case 'roles':
                    stat.value = 32
                    break
                case 'permissions':
                    stat.value = 256
                    break
                case 'logs':
                    stat.value = 156
                    break
            }
        })
    }
}

onMounted(() => {
    loadStats()
})

defineExpose({
    loadStats,
    stats,
})
</script>

<style scoped lang="less">
.stats-row {
    margin-bottom: 24px;
}

.stat-card {
    border-radius: 8px;
    transition: all 0.3s;
    border: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    &:hover {
        transform: translateY(-4px);
        box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
    }

    &.primary {
        border-left: 4px solid var(--color-primary);
    }

    &.success {
        border-left: 4px solid #52c41a;
    }

    &.warning {
        border-left: 4px solid #faad14;
    }

    &.danger {
        border-left: 4px solid #f5222d;
    }

    .stat-footer {
        margin-top: 16px;
        padding-top: 16px;
        border-top: 1px solid #f0f0f0;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .stat-change {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-size: 14px;

        &.up {
            color: #52c41a;
        }

        &.down {
            color: #f5222d;
        }
    }

    .stat-desc {
        font-size: 12px;
        color: #8c8c8c;
    }
}

body.dark-mode {
    .stat-card {
        background: #1f1f1f;
        border-color: #303030;

        .stat-footer {
            border-color: #303030;
        }
    }
}
</style>


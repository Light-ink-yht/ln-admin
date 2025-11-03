<template>
    <a-card title="实时统计" class="realtime-stats-card" :bordered="false">
        <a-row :gutter="[16, 16]">
            <a-col :xs="24" :sm="12" :lg="6" v-for="stat in stats" :key="stat.key">
                <div class="stat-item">
                    <div class="stat-label">{{ stat.label }}</div>
                    <div class="stat-value" :style="{ color: stat.color || 'var(--color-primary)' }">
                        {{ stat.value }}
                    </div>
                    <div class="stat-desc">{{ stat.desc }}</div>
                </div>
            </a-col>
        </a-row>
    </a-card>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

export interface RealTimeStat {
    key: string
    label: string
    value: string | number
    desc: string
    color?: string
}

interface Props {
    stats?: RealTimeStat[]
}

const props = withDefaults(defineProps<Props>(), {
    stats: undefined,
})

const defaultStats: RealTimeStat[] = [
    {
        key: 'online',
        label: '在线用户',
        value: '1,234',
        desc: '当前在线',
        color: 'var(--color-primary)',
    },
    {
        key: 'today',
        label: '今日访问',
        value: '5,678',
        desc: '较昨日 +12%',
        color: '#52c41a',
    },
    {
        key: 'pv',
        label: '页面浏览',
        value: '12,345',
        desc: '今日总量',
        color: '#faad14',
    },
    {
        key: 'avg',
        label: '平均时长',
        value: '3分45秒',
        desc: '用户停留',
        color: '#f5222d',
    },
]

const stats = computed(() => {
    return props.stats && props.stats.length > 0 ? props.stats : defaultStats
})
</script>

<style scoped lang="less">
.realtime-stats-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }
}

.stat-item {
    padding: 16px;
    border-radius: 8px;
    background: #fafafa;
    text-align: center;
    transition: all 0.3s;

    &:hover {
        background: #f0f0f0;
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
    }
}

.stat-label {
    font-size: 14px;
    color: #8c8c8c;
    margin-bottom: 8px;
}

.stat-value {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 8px;
}

.stat-desc {
    font-size: 12px;
    color: #bfbfbf;
}

// 暗色模式
body.dark-mode {
    .realtime-stats-card {
        background: #1f1f1f;
        border-color: #303030;
    }

    .stat-item {
        background: #262626;
        border: 1px solid #303030;

        &:hover {
            background: #2d2d2d;
            border-color: var(--color-primary);
        }
    }

    .stat-label {
        color: rgba(255, 255, 255, 0.65);
    }

    .stat-desc {
        color: rgba(255, 255, 255, 0.45);
    }
}
</style>


<template>
    <a-card title="访问趋势" class="trend-chart-card" :bordered="false">
        <div class="chart-container">
            <div class="chart-legend">
                <div class="legend-item" v-for="item in legendData" :key="item.key">
                    <span class="legend-dot" :style="{ background: item.color }"></span>
                    <span class="legend-label">{{ item.label }}</span>
                    <span class="legend-value">{{ item.value }}</span>
                </div>
            </div>
            <div class="chart-content">
                <div class="chart-bars">
                    <div
                        v-for="item in chartData"
                        :key="item.date"
                        class="chart-item"
                        :title="`${item.date}: ${item.value}`"
                    >
                        <div class="bar-wrapper">
                            <div
                                class="bar"
                                :style="{
                                    height: `${(item.value / maxValue) * 100}%`,
                                    background: item.color || 'var(--color-primary)',
                                }"
                            >
                                <span class="bar-value">{{ item.value }}</span>
                            </div>
                        </div>
                        <div class="bar-label">{{ item.date }}</div>
                    </div>
                </div>
            </div>
        </div>
    </a-card>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

export interface TrendData {
    date: string
    value: number
    color?: string
}

interface Props {
    chartData?: TrendData[]
    legendData?: Array<{
        key: string
        label: string
        value: string
        color: string
    }>
}

const props = withDefaults(defineProps<Props>(), {
    chartData: undefined,
    legendData: undefined,
})

const defaultChartData: TrendData[] = [
    { date: '周一', value: 320, color: 'var(--color-primary)' },
    { date: '周二', value: 450, color: 'var(--color-primary)' },
    { date: '周三', value: 280, color: 'var(--color-primary)' },
    { date: '周四', value: 520, color: 'var(--color-primary)' },
    { date: '周五', value: 380, color: 'var(--color-primary)' },
    { date: '周六', value: 490, color: 'var(--color-primary)' },
    { date: '周日', value: 410, color: 'var(--color-primary)' },
]

const defaultLegendData = [
    {
        key: 'pv',
        label: '访问量',
        value: '2,456',
        color: 'var(--color-primary)',
    },
    {
        key: 'uv',
        label: '访客数',
        value: '1,234',
        color: '#52c41a',
    },
]

const chartData = computed(() => {
    return props.chartData && props.chartData.length > 0 ? props.chartData : defaultChartData
})

const legendData = computed(() => {
    return props.legendData && props.legendData.length > 0 ? props.legendData : defaultLegendData
})

const maxValue = computed(() => {
    if (!chartData.value || chartData.value.length === 0) {
        return 1
    }
    return Math.max(...chartData.value.map((item) => item.value), 0) || 1
})
</script>

<style scoped lang="less">
.trend-chart-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }
}

.chart-container {
    padding: 8px 0;
}

.chart-legend {
    display: flex;
    gap: 24px;
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #f0f0f0;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.legend-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
}

.legend-label {
    font-size: 14px;
    color: #8c8c8c;
}

.legend-value {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
}

.chart-content {
    height: 300px;
}

.chart-bars {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    height: 100%;
    gap: 8px;
}

.chart-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    height: 100%;
}

.bar-wrapper {
    flex: 1;
    width: 100%;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    position: relative;
}

.bar {
    width: 100%;
    min-height: 20px;
    border-radius: 4px 4px 0 0;
    position: relative;
    transition: all 0.3s;
    cursor: pointer;

    &:hover {
        opacity: 0.8;
        transform: scaleY(1.05);
    }
}

.bar-value {
    position: absolute;
    top: -24px;
    left: 50%;
    transform: translateX(-50%);
    font-size: 12px;
    color: #595959;
    white-space: nowrap;
    opacity: 0;
    transition: opacity 0.3s;
}

.chart-item:hover .bar-value {
    opacity: 1;
}

.bar-label {
    margin-top: 8px;
    font-size: 12px;
    color: #8c8c8c;
    text-align: center;
}

// 响应式
@media (max-width: 768px) {
    .chart-content {
        height: 200px;
    }

    .chart-legend {
        flex-wrap: wrap;
        gap: 16px;
    }
}

// 暗色模式
body.dark-mode {
    .trend-chart-card {
        background: #1f1f1f;
        border-color: #303030;
    }

    .chart-legend {
        border-color: #303030;
    }

    .legend-label {
        color: rgba(255, 255, 255, 0.65);
    }

    .legend-value {
        color: rgba(255, 255, 255, 0.85);
    }

    .bar-value {
        color: rgba(255, 255, 255, 0.65);
    }

    .bar-label {
        color: rgba(255, 255, 255, 0.65);
    }
}
</style>


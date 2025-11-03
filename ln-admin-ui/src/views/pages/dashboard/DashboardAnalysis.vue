<template>
    <div class="dashboard-analysis-container">
        <!-- 总览卡片 -->
        <OverviewCards :cards="overviewCards || undefined" />

        <!-- 实时统计 -->
        <RealTimeStats :stats="realTimeStats || undefined" style="margin-bottom: 24px;" />

        <!-- 主要数据区域 -->
        <a-row :gutter="[16, 16]" class="main-content-row">
            <!-- 左侧：趋势图表 -->
            <a-col :xs="24" :lg="16">
                <TrendChart :chart-data="trendData || undefined" :legend-data="legendData || undefined" />
            </a-col>

            <!-- 右侧：访问来源 -->
            <a-col :xs="24" :lg="8">
                <SourceAnalysis :source-data="sourceData || undefined" />
            </a-col>
        </a-row>

        <!-- 数据表格 -->
        <DataTable
            title="页面访问排行"
            :columns="tableColumns"
            :table-data="tableData || undefined"
            :pagination="{ pageSize: 10 }"
            style="margin-top: 24px;"
        />
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import {
    OverviewCards,
    TrendChart,
    DataTable,
    SourceAnalysis,
    RealTimeStats,
    type OverviewCard,
    type TrendData,
    type TableDataItem,
    type SourceItem,
    type RealTimeStat,
} from './analysis/components'
import type { TableColumnsType } from 'ant-design-vue'

const overviewCards = ref<OverviewCard[] | null>(null)
const realTimeStats = ref<RealTimeStat[] | null>(null)
const trendData = ref<TrendData[] | null>(null)
const legendData = ref<Array<{ key: string; label: string; value: string; color: string }> | null>(null)
const sourceData = ref<SourceItem[] | null>(null)
const tableData = ref<TableDataItem[] | null>(null)
const tableColumns = ref<TableColumnsType>([
    {
        title: '页面名称',
        dataIndex: 'name',
        key: 'name',
    },
    {
        title: '访问量',
        dataIndex: 'views',
        key: 'views',
        align: 'right',
    },
    {
        title: '变化趋势',
        key: 'trend',
        align: 'center',
    },
    {
        title: '状态',
        key: 'status',
        align: 'center',
    },
])

// 模拟数据加载
const loadData = async () => {
    try {
        // 这里可以从 API 获取数据
        // 示例：使用默认数据（组件内部会处理空值，使用默认数据）
        await new Promise((resolve) => setTimeout(resolve, 500))
        // 可以在这里设置数据，如果传 null，组件会使用默认数据
    } catch (error) {
        console.error('加载数据失败:', error)
    }
}

onMounted(() => {
    loadData()
})
</script>

<style scoped lang="less">
.dashboard-analysis-container {
    padding: 24px;
    min-height: 100%;
}

.main-content-row {
    margin-bottom: 0;
}

// 响应式设计
@media (max-width: 768px) {
    .dashboard-analysis-container {
        padding: 16px;
    }
}
</style>


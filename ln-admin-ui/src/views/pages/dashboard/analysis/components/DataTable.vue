<template>
    <a-card :title="title" class="data-table-card" :bordered="false">
        <a-table
            :columns="columns"
            :data-source="tableData"
            :pagination="pagination"
            :loading="loading"
            size="middle"
        >
            <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'trend'">
                    <a-tag :color="record.trend === 'up' ? 'success' : 'error'">
                        <template #icon>
                            <component
                                :is="record.trend === 'up' ? ArrowUpOutlined : ArrowDownOutlined"
                            />
                        </template>
                        {{ record.change }}
                    </a-tag>
                </template>
                <template v-else-if="column.key === 'status'">
                    <a-tag :color="getStatusColor(record.status)">
                        {{ record.status }}
                    </a-tag>
                </template>
            </template>
        </a-table>
    </a-card>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'

export interface TableDataItem {
    key: string
    [key: string]: any
}

interface Props {
    title?: string
    columns?: TableColumnsType
    tableData?: TableDataItem[]
    loading?: boolean
    pagination?: any
}

const props = withDefaults(defineProps<Props>(), {
    title: '数据表格',
    columns: undefined,
    tableData: undefined,
    loading: false,
    pagination: { pageSize: 10 },
})

const defaultColumns: TableColumnsType = [
    {
        title: '名称',
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
        title: '变化',
        key: 'trend',
        align: 'center',
    },
    {
        title: '状态',
        key: 'status',
        align: 'center',
    },
]

const defaultTableData = [
    {
        key: '1',
        name: '首页',
        views: 1234,
        change: '+12.5%',
        trend: 'up',
        status: '正常',
    },
    {
        key: '2',
        name: '用户中心',
        views: 856,
        change: '+8.2%',
        trend: 'up',
        status: '正常',
    },
    {
        key: '3',
        name: '设置页面',
        views: 432,
        change: '-3.1%',
        trend: 'down',
        status: '警告',
    },
]

const columns = computed(() => {
    return props.columns || defaultColumns
})

const tableData = computed(() => {
    return props.tableData || defaultTableData
})

const getStatusColor = (status: string) => {
    const colorMap: Record<string, string> = {
        正常: 'success',
        警告: 'warning',
        异常: 'error',
    }
    return colorMap[status] || 'default'
}
</script>

<style scoped lang="less">
.data-table-card {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }

    :deep(.ant-table) {
        .ant-table-thead > tr > th {
            background: #fafafa;
            font-weight: 600;
        }
    }
}

// 暗色模式
body.dark-mode {
    .data-table-card {
        background: #1f1f1f;
        border-color: #303030;

        :deep(.ant-table) {
            .ant-table-thead > tr > th {
                background: #262626;
            }
        }
    }
}
</style>


<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问操作日志的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="操作日志"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchLogList"
        :scroll="{ x: 1400 }"
        row-key="logId"
    >
        <!-- 自定义列插槽 -->
        <template #column-level="{ record }">
            <a-tag :color="getLevelColor(record.level)">
                {{ record.level.toUpperCase() }}
            </a-tag>
        </template>
        <template #column-method="{ record }">
            <a-tag :color="getMethodColor(record.method)">
                {{ record.method }}
            </a-tag>
        </template>
        <template #column-statusCode="{ record }">
            <a-tag :color="getStatusCodeColor(record.statusCode)">
                {{ record.statusCode }}
            </a-tag>
        </template>
        <template #column-actions="{ record }">
            <a-button type="link" size="small" @click="handleViewDetail(record)">
                查看详情
            </a-button>
        </template>
        <template #column-content="{ record }">
            <a-typography-paragraph
                :ellipsis="{ rows: 2, expandable: true, symbol: '展开' }"
                style="margin: 0; max-width: 300px"
            >
                {{ record.content }}
            </a-typography-paragraph>
        </template>
        <template #column-logTime="{ record }">
            <UserText :value="record.logTime" />
        </template>
    </DataList>

    <!-- 日志详情弹窗 -->
    <a-modal
        v-model:open="detailModalOpen"
        title="日志详情"
        width="800px"
        :footer="null"
    >
        <a-descriptions :column="2" bordered>
            <a-descriptions-item label="日志ID" :span="2">
                <code>{{ currentLog?.logId }}</code>
            </a-descriptions-item>
            <a-descriptions-item label="日志级别">
                <a-tag :color="getLevelColor(currentLog?.level || '')">
                    {{ currentLog?.level?.toUpperCase() }}
                </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="模块名称">
                {{ currentLog?.module }}
            </a-descriptions-item>
            <a-descriptions-item label="操作类型">
                <a-space>
                    <component :is="getActionIcon(currentLog?.action || '')" />
                    <span>{{ currentLog?.action }}</span>
                </a-space>
            </a-descriptions-item>
            <a-descriptions-item label="用户ID">
                {{ currentLog?.userId || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="IP地址">
                {{ currentLog?.ip }}
            </a-descriptions-item>
            <a-descriptions-item label="请求方法">
                <a-tag :color="getMethodColor(currentLog?.method || '')">
                    {{ currentLog?.method }}
                </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="状态码">
                <a-tag :color="getStatusCodeColor(currentLog?.statusCode || 0)">
                    {{ currentLog?.statusCode }}
                </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="请求路径" :span="2">
                <code>{{ currentLog?.path }}</code>
            </a-descriptions-item>
            <a-descriptions-item label="日志内容" :span="2">
                <a-typography-paragraph style="margin: 0; white-space: pre-wrap">
                    {{ currentLog?.content }}
                </a-typography-paragraph>
            </a-descriptions-item>
            <a-descriptions-item v-if="currentLog?.errorMsg" label="错误信息" :span="2">
                <a-typography-paragraph style="margin: 0; color: #ff4d4f; white-space: pre-wrap">
                    {{ currentLog.errorMsg }}
                </a-typography-paragraph>
            </a-descriptions-item>
            <a-descriptions-item label="用户代理" :span="2">
                <a-typography-paragraph style="margin: 0; font-size: 12px; color: #8c8c8c">
                    {{ currentLog?.userAgent }}
                </a-typography-paragraph>
            </a-descriptions-item>
            <a-descriptions-item label="日志时间">
                {{ currentLog?.logTime }}
            </a-descriptions-item>
            <a-descriptions-item label="创建时间">
                {{ currentLog?.createdAt }}
            </a-descriptions-item>
        </a-descriptions>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import {
    DataList,
    UserText,
    PermissionDeniedAlert,
    getSystemLogSearchFields,
    getSystemLogColumns,
    type PageResponse,
} from '@/views/pages/components'
import { systemApi, type SystemLog, type SystemLogListParams } from '@/api/system'
import {
    FileTextOutlined,
    LockOutlined,
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    EyeOutlined,
    SettingOutlined,
} from '@ant-design/icons-vue'

// 组件引用
const dataListRef = ref<InstanceType<typeof DataList>>()
// 权限提示
const permissionDenied = ref(false)

// 详情弹窗
const detailModalOpen = ref(false)
const currentLog = ref<SystemLog | null>(null)

// 搜索字段配置
const searchFields = getSystemLogSearchFields()

// 表格列定义
const columns = getSystemLogColumns()

// 获取日志列表
const fetchLogList = async (params: any): Promise<PageResponse<SystemLog[]>> => {
    try {
        // 处理时间范围
        const requestParams: SystemLogListParams = {
            page: params.page || params.page_size || 1,
            pageSize: params.pageSize || params.page_size || 10,
            level: params.level,
            module: params.module,
            action: params.action,
            userId: params.userId,
            ip: params.ip,
        }

        // 处理日期范围（dayjs对象）
        if (params.dateRange && Array.isArray(params.dateRange) && params.dateRange.length === 2) {
            const start = params.dateRange[0]
            const end = params.dateRange[1]
            if (start && end && typeof start.format === 'function' && typeof end.format === 'function') {
                requestParams.startTime = start.format('YYYY-MM-DD HH:mm:ss')
                requestParams.endTime = end.format('YYYY-MM-DD HH:mm:ss')
            }
        }

        const response = await systemApi.getLogList(requestParams)

        // 检查权限错误
        if (response.code === 403) {
            permissionDenied.value = true
            throw new Error(response.msg || response.message || '没有权限访问该资源')
        }

        // 权限验证通过，隐藏权限提示
        permissionDenied.value = false

        if (response.data && response.data.list) {
            return {
                ...response,
                data: response.data.list,
                total: response.data.total,
            }
        }

        return {
            ...response,
            data: [],
            total: 0,
        }
    } catch (error: any) {
        console.error('获取日志列表失败:', error)
        if (error.response?.status === 403 || error.message?.includes('没有权限')) {
            permissionDenied.value = true
        }
        throw error
    }
}

// 获取日志级别颜色
const getLevelColor = (level: string): string => {
    const colorMap: Record<string, string> = {
        info: 'blue',
        warn: 'orange',
        error: 'red',
        debug: 'default',
    }
    return colorMap[level.toLowerCase()] || 'default'
}

// 获取HTTP方法颜色
const getMethodColor = (method: string): string => {
    const colorMap: Record<string, string> = {
        GET: 'green',
        POST: 'blue',
        PUT: 'orange',
        DELETE: 'red',
        PATCH: 'purple',
    }
    return colorMap[method.toUpperCase()] || 'default'
}

// 获取状态码颜色
const getStatusCodeColor = (statusCode: number): string => {
    if (statusCode >= 200 && statusCode < 300) {
        return 'success'
    } else if (statusCode >= 300 && statusCode < 400) {
        return 'processing'
    } else if (statusCode >= 400 && statusCode < 500) {
        return 'warning'
    } else if (statusCode >= 500) {
        return 'error'
    }
    return 'default'
}

// 获取操作图标
const getActionIcon = (action: string) => {
    const iconMap: Record<string, any> = {
        login: LockOutlined,
        logout: LockOutlined,
        create: PlusOutlined,
        update: EditOutlined,
        delete: DeleteOutlined,
        view: EyeOutlined,
        config: SettingOutlined,
    }
    return iconMap[action.toLowerCase()] || FileTextOutlined
}

// 查看详情
const handleViewDetail = (record: SystemLog) => {
    currentLog.value = record
    detailModalOpen.value = true
}
</script>

<style scoped lang="less">
.system-log-container {
    padding: 0;
    background: transparent;
    min-height: 100%;
}
</style>

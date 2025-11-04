<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问短信验证码列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="短信验证码列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchCodeList"
        :scroll="{ x: 1600 }"
        row-key="codeId"
    >
        <!-- 自定义列插槽 -->
        <template #column-status="{ record }">
            <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
            </a-tag>
        </template>
        <template #column-expireAt="{ record }">
            <UserText :value="record.expireAt || '-'" />
        </template>
        <template #column-usedAt="{ record }">
            <UserText :value="record.usedAt || '-'" />
        </template>
        <template #column-createdAt="{ record }">
            <UserText :value="record.createdAt" />
        </template>
    </DataList>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import {
    DataList,
    UserText,
    PermissionDeniedAlert,
    getSMSCodeSearchFields,
    getSMSCodeColumns,
    type PageResponse,
} from '@/views/pages/components'
import { smsApi, type SMSCode, type SMSCodeListParams } from '@/api/sms'

// 组件引用
const dataListRef = ref<InstanceType<typeof DataList>>()
// 权限提示
const permissionDenied = ref(false)

// 搜索字段配置
const searchFields = getSMSCodeSearchFields()

// 表格列定义
const columns = getSMSCodeColumns()

// 获取状态颜色
const getStatusColor = (status: string): string => {
    const colorMap: Record<string, string> = {
        '1': 'blue',
        '2': 'green',
        '3': 'red',
    }
    return colorMap[status] || 'default'
}

// 获取状态文本
const getStatusText = (status: string): string => {
    const textMap: Record<string, string> = {
        '1': '未使用',
        '2': '已使用',
        '3': '已过期',
    }
    return textMap[status] || status
}

// 获取验证码列表
const fetchCodeList = async (params: any): Promise<PageResponse<SMSCode[]>> => {
    try {
        const requestParams: SMSCodeListParams = {
            page: params.page || params.page_size || 1,
            pageSize: params.pageSize || params.page_size || 10,
            phone: params.phone,
            type: params.type,
            status: params.status,
        }

        const response = await smsApi.getCodeList(requestParams)

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
        console.error('获取短信验证码列表失败:', error)
        if (error.response?.status === 403 || error.message?.includes('没有权限')) {
            permissionDenied.value = true
        }
        throw error
    }
}
</script>


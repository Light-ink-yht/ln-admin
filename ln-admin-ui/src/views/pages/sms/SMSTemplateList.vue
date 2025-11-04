<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问短信模板列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="短信模板列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchTemplateList"
        :actions="actions"
        :extra-actions="extraActions"
        :scroll="{ x: 1400 }"
        row-key="templateId"
    >
        <!-- 自定义列插槽 -->
        <template #column-status="{ record }">
            <UserStatus :status="record.status" />
        </template>
        <template #column-content="{ record }">
            <UserText :value="record.content" />
        </template>
        <template #column-createdAt="{ record }">
            <UserText :value="record.createdAt" />
        </template>
    </DataList>

    <!-- 模板表单弹窗 -->
    <SMSTemplateForm
        v-model:open="templateFormOpen"
        :template="currentTemplate"
        @success="handleFormSuccess"
    />

    <!-- 模板详情弹窗 -->
    <SMSTemplateDetail
        v-model:open="templateDetailOpen"
        :template-id="currentTemplateId"
    />
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import {
    DataList,
    UserStatus,
    UserText,
    PermissionDeniedAlert,
    getSMSTemplateSearchFields,
    getSMSTemplateColumns,
    type PageResponse,
    type ExtraAction,
    type ActionButton,
} from '@/views/pages/components'
import { smsApi, type SMSTemplate, type SMSTemplateListParams } from '@/api/sms'
import SMSTemplateForm from './components/SMSTemplateForm.vue'
import SMSTemplateDetail from './components/SMSTemplateDetail.vue'

// 组件引用
const dataListRef = ref<InstanceType<typeof DataList>>()
// 权限提示
const permissionDenied = ref(false)

// 组件状态
const templateFormOpen = ref(false)
const templateDetailOpen = ref(false)
const currentTemplate = ref<SMSTemplate | null>(null)
const currentTemplateId = ref<string | null>(null)

// 搜索字段配置
const searchFields = getSMSTemplateSearchFields()

// 表格列定义
const columns = getSMSTemplateColumns()

// 操作按钮
const actions: ActionButton[] = [
    {
        key: 'view',
        label: '查看',
        type: 'link',
        onClick: (record: SMSTemplate) => {
            handleViewTemplate(record)
        },
    },
    {
        key: 'edit',
        label: '编辑',
        type: 'link',
        onClick: (record: SMSTemplate) => {
            handleEditTemplate(record)
        },
    },
    {
        key: 'delete',
        label: '删除',
        type: 'link',
        danger: true,
        onClick: (record: SMSTemplate) => {
            handleDeleteTemplate(record)
        },
    },
]

// 额外操作按钮
const extraActions: ExtraAction[] = [
    {
        key: 'add',
        label: '添加模板',
        icon: PlusOutlined,
        type: 'primary',
        onClick: () => {
            handleAddTemplate()
        },
    },
]

// 获取模板列表
const fetchTemplateList = async (params: any): Promise<PageResponse<SMSTemplate[]>> => {
    try {
        const requestParams: SMSTemplateListParams = {
            page: params.page || params.page_size || 1,
            pageSize: params.pageSize || params.page_size || 10,
            type: params.type,
            status: params.status,
        }

        const response = await smsApi.getTemplateList(requestParams)

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
        console.error('获取短信模板列表失败:', error)
        if (error.response?.status === 403 || error.message?.includes('没有权限')) {
            permissionDenied.value = true
        }
        throw error
    }
}

// 添加模板
const handleAddTemplate = () => {
    currentTemplate.value = null
    templateFormOpen.value = true
}

// 查看模板
const handleViewTemplate = (record: SMSTemplate) => {
    currentTemplateId.value = record.templateId
    templateDetailOpen.value = true
}

// 编辑模板
const handleEditTemplate = (record: SMSTemplate) => {
    currentTemplate.value = record
    templateFormOpen.value = true
}

// 删除模板
const handleDeleteTemplate = (record: SMSTemplate) => {
    Modal.confirm({
        title: '确认删除',
        content: `确定要删除模板"${record.templateName}"吗？此操作不可恢复。`,
        okText: '确定',
        cancelText: '取消',
        okType: 'danger',
        onOk: async () => {
            try {
                await smsApi.deleteTemplate(record.templateId)
                message.success('删除成功')
                dataListRef.value?.refresh()
            } catch (error: any) {
                console.error('删除失败:', error)
                message.error(error.message || '删除失败')
            }
        },
    })
}

// 表单成功回调
const handleFormSuccess = () => {
    dataListRef.value?.refresh()
}
</script>


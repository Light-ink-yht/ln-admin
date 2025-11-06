<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问权限列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="权限列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchPermissionList"
        :actions="actions"
        :extra-actions="extraActions"
        :scroll="{ x: 1500 }"
        row-key="permissionId"
    >
        <!-- 自定义列插槽 -->
        <template #column-category="{ record }">
            <a-tag :color="getCategoryColor(record.category || getPermissionCategory(record.permissionKey))">
                {{ record.category || getPermissionCategory(record.permissionKey) }}
            </a-tag>
        </template>

        <template #column-status="{ record }">
            <UserStatus :status="record.status" />
        </template>

        <template #column-description="{ record }">
            <UserText :value="record.description" />
        </template>

        <template #column-createdAt="{ record }">
            <UserText :value="record.CreatedAt || record.createdAt" />
        </template>

        <template #column-updatedAt="{ record }">
            <UserText :value="record.updatedAt" />
        </template>

        <template #column-creatorId="{ record }">
            <UserText :value="record.creatorName || record.creatorId || '-'" />
        </template>

        <template #column-modifierId="{ record }">
            <UserText :value="record.modifierName || record.modifierId || '-'" />
        </template>
    </DataList>

    <!-- 权限表单弹窗 -->
    <PermissionForm
        v-model:open="permissionFormOpen"
        :permission="currentPermission"
        @success="handleFormSuccess"
    />

    <!-- 权限详情弹窗 -->
    <PermissionDetail
        v-model:open="permissionDetailOpen"
        :permission-id="currentPermissionId"
    />
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import {
    DataList,
    UserStatus,
    UserText,
    PermissionDeniedAlert,
    getPermissionSearchFields,
    getPermissionColumns,
    getPermissionCategory,
    mapPermissionList,
    type ExtraAction,
    type ActionButton,
    type PageResponse,
} from '@/views/pages/components'
import { permissionApi, type Permission } from '@/api/permission'
import PermissionForm from './components/PermissionForm.vue'
import PermissionDetail from './components/PermissionDetail.vue'

// 组件引用
const dataListRef = ref<InstanceType<typeof DataList>>()
// 权限提示
const permissionDenied = ref(false)

// 组件状态
const permissionFormOpen = ref(false)
const permissionDetailOpen = ref(false)
const currentPermission = ref<Permission | null>(null)
const currentPermissionId = ref<string | null>(null)

// 搜索字段配置（使用封装的函数）
const searchFields = getPermissionSearchFields()

// 表格列定义（使用封装的函数）
const columns = getPermissionColumns()

// 获取分类颜色
const getCategoryColor = (category: string): string => {
    const colorMap: Record<string, string> = {
        '用户管理': 'blue',
        '角色管理': 'green',
        '权限管理': 'orange',
        '系统配置': 'purple',
        '短信管理': 'cyan',
        '菜单管理': 'magenta',
        '系统运维': 'red',
        '文件管理': 'geekblue',
        '其他': 'default',
    }
    return colorMap[category] || 'default'
}

// 数据获取函数
const fetchPermissionList = async (params: any): Promise<PageResponse<Permission[]>> => {
    try {
        const response = await permissionApi.getPermissionList(params)
        console.log('权限列表API响应:', response)

        // 检查权限错误
        if (response.code === 403) {
            permissionDenied.value = true
            throw new Error(response.msg || response.message || '没有权限访问该资源')
        }

        // 权限验证通过，隐藏权限提示
        permissionDenied.value = false

        // 确保数据格式正确，处理可能的字段名差异
        if (response.data && Array.isArray(response.data)) {
            response.data = mapPermissionList(response.data)
        }

        return response
    } catch (error: any) {
        console.error('获取权限列表失败:', error)
        // 如果是权限错误，已经在上面设置了 permissionDenied
        if (error.message?.includes('权限') || error.message?.includes('没有权限')) {
            permissionDenied.value = true
        }
        throw error
    }
}

// 顶部操作按钮配置
const extraActions: ExtraAction[] = [
    {
        key: 'add',
        label: '添加权限',
        type: 'primary',
        icon: PlusOutlined,
        onClick: () => {
            currentPermission.value = null
            permissionFormOpen.value = true
        },
    },
]

// 操作按钮配置
const actions: ActionButton[] = [
    {
        key: 'view',
        label: '查看',
        type: 'link',
        onClick: (record: Permission) => {
            currentPermissionId.value = record.permissionId
            permissionDetailOpen.value = true
        },
    },
    {
        key: 'edit',
        label: '编辑',
        type: 'link',
        onClick: (record: Permission) => {
            currentPermission.value = record
            permissionFormOpen.value = true
        },
    },
    {
        key: 'delete',
        label: '删除',
        type: 'link',
        danger: true,
        confirm: '确定要删除该权限吗？删除后不可恢复！',
        onClick: async (record: Permission) => {
            try {
                await permissionApi.deletePermission(record.permissionId)
                message.success('删除成功')
                dataListRef.value?.refresh()
            } catch (error: any) {
                // 错误提示已在 request.ts 中统一处理，这里不再重复显示
            }
        },
    },
]

// 权限操作成功回调
const handleFormSuccess = () => {
    dataListRef.value?.refresh()
}
</script>

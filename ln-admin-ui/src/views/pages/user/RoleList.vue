<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问角色列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="角色列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchRoleList"
        :actions="actions"
        :extra-actions="extraActions"
        :scroll="{ x: 1200 }"
        row-key="roleId"
    >
        <!-- 自定义列插槽 -->
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

    <!-- 角色表单弹窗 -->
    <RoleForm
        v-model:open="roleFormOpen"
        :role="currentRole"
        @success="handleFormSuccess"
    />

    <!-- 角色详情弹窗 -->
    <RoleDetail
        v-model:open="roleDetailOpen"
        :role-id="currentRoleId"
    />

    <!-- 角色菜单授权弹窗 -->
    <RoleMenuGrant
        v-model:open="roleMenuGrantOpen"
        :role-id="currentRoleId"
        @success="handleMenuGrantSuccess"
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
    getRoleSearchFields,
    getRoleColumns,
    mapRoleList,
    type ExtraAction,
    type ActionButton,
    type PageResponse,
} from '@/views/pages/components'
import { roleApi, type Role } from '@/api/role'
import RoleForm from './components/RoleForm.vue'
import RoleDetail from './components/RoleDetail.vue'
import RoleMenuGrant from './components/RoleMenuGrant.vue'

// 组件引用
const dataListRef = ref<InstanceType<typeof DataList>>()
// 权限提示
const permissionDenied = ref(false)

// 组件状态
const roleFormOpen = ref(false)
const roleDetailOpen = ref(false)
const roleMenuGrantOpen = ref(false)
const currentRole = ref<Role | null>(null)
const currentRoleId = ref<string | null>(null)

// 搜索字段配置（使用封装的函数）
const searchFields = getRoleSearchFields()

// 表格列定义（使用封装的函数）
const columns = getRoleColumns()

// 数据获取函数
const fetchRoleList = async (params: any): Promise<PageResponse<Role[]>> => {
    try {
        const response = await roleApi.getRoleList(params)
        console.log('角色列表API响应:', response)
        
        // 检查权限错误
        if (response.code === 403) {
            permissionDenied.value = true
            throw new Error(response.msg || response.message || '没有权限访问该资源')
        }
        
        // 权限验证通过，隐藏权限提示
        permissionDenied.value = false
        
        // 确保数据格式正确，处理可能的字段名差异
        if (response.data && Array.isArray(response.data)) {
            response.data = mapRoleList(response.data)
        }
        
        return response
    } catch (error: any) {
        console.error('获取角色列表失败:', error)
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
        label: '添加角色',
        type: 'primary',
        icon: PlusOutlined,
        onClick: () => {
            currentRole.value = null
            roleFormOpen.value = true
        },
    },
]

// 操作按钮配置
const actions: ActionButton[] = [
    {
        key: 'view',
        label: '查看',
        type: 'link',
        onClick: (record: Role) => {
            currentRoleId.value = record.roleId
            roleDetailOpen.value = true
        },
    },
    {
        key: 'edit',
        label: '编辑',
        type: 'link',
        onClick: (record: Role) => {
            currentRole.value = record
            roleFormOpen.value = true
        },
    },
    {
        key: 'menu-grant',
        label: '菜单授权',
        type: 'link',
        onClick: (record: Role) => {
            currentRoleId.value = record.roleId
            roleMenuGrantOpen.value = true
        },
    },
    {
        key: 'delete',
        label: '删除',
        type: 'link',
        danger: true,
        confirm: '确定要删除该角色吗？删除后不可恢复！',
        onClick: async (record: Role) => {
            try {
                await roleApi.deleteRole(record.roleId)
                message.success('删除成功')
                dataListRef.value?.refresh()
            } catch (error: any) {
                message.error(error.message || '删除失败')
            }
        },
    },
]

// 角色操作成功回调
const handleFormSuccess = () => {
    dataListRef.value?.refresh()
}

// 菜单授权成功回调
const handleMenuGrantSuccess = () => {
    roleMenuGrantOpen.value = false
    message.success('菜单授权已更新')
}
</script>

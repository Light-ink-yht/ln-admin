<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问用户列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="用户列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchUserList"
        :actions="actions"
        :extra-actions="extraActions"
        :scroll="{ x: 1200 }"
        row-key="userId"
    >
        <!-- 自定义列插槽 -->
        <template #column-avatar="{ record }">
            <UserAvatar
                :avatar="record.avatar"
                :nickname="record.nickname"
                :full-name="record.fullName"
                :size="44"
            />
        </template>

        <template #column-status="{ record }">
            <UserStatus :status="record.status" />
        </template>

        <template #column-gender="{ record }">
            <UserGender :gender="record.gender" />
        </template>

        <template #column-phone="{ record }">
            <UserText :value="record.phone" />
        </template>

        <template #column-email="{ record }">
            <UserText :value="record.email" />
        </template>

        <template #column-lastLoginTime="{ record }">
            <UserText :value="record.lastLoginTime" />
        </template>

        <template #column-CreatedAt="{ record }">
            <UserText :value="record.CreatedAt" />
        </template>
    </DataList>

    <!-- 用户表单弹窗 -->
    <UserForm
        v-model:open="userFormOpen"
        :user="currentUser"
        @success="handleFormSuccess"
    />

    <!-- 用户详情弹窗 -->
    <UserDetail
        v-model:open="userDetailOpen"
        :user-id="currentUserId"
    />
</template>

<script lang="ts" setup>
import { PlusOutlined } from '@ant-design/icons-vue'
import {
    DataList,
    UserAvatar,
    UserStatus,
    UserGender,
    UserText,
    PermissionDeniedAlert,
    getUserSearchFields,
    getUserColumns,
    useUserList,
    type ExtraAction,
} from '@/views/pages/components'
import { userApi, type UserResponse, type PageResponse } from '@/api/user'
import UserForm from './components/UserForm.vue'
import UserDetail from './components/UserDetail.vue'

// 使用用户列表组合式函数
const {
    permissionDenied,
    userFormOpen,
    userDetailOpen,
    currentUser,
    currentUserId,
    dataListRef,
    fetchUserList,
    handleFormSuccess,
    actions,
} = useUserList(
    async (params: any) => {
        return await userApi.getUserList(params)
    },
    {
        onDelete: async (userId: string) => {
            await userApi.deleteUser(userId)
        },
    }
)

// 搜索字段配置（使用封装的函数）
const searchFields = getUserSearchFields()

// 表格列定义（使用封装的函数）
const columns = getUserColumns()

// 顶部操作按钮配置
const extraActions: ExtraAction[] = [
    {
        key: 'add',
        label: '添加用户',
        type: 'primary',
        icon: PlusOutlined,
        onClick: () => {
            currentUser.value = null
            userFormOpen.value = true
        },
    },
]
</script>


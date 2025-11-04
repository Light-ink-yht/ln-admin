<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问菜单列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="菜单列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchMenuList"
        :actions="actions"
        :extra-actions="extraActions"
        :scroll="{ x: 1400 }"
        row-key="menu_id"
    >
        <!-- 自定义列插槽 -->
        <template #column-menu_type="{ record }">
            <a-tag :color="record.menu_type === '1' ? 'blue' : 'green'">
                {{ record.menu_type === '1' ? '侧边栏菜单' : '用户菜单' }}
            </a-tag>
        </template>
        <template #column-status="{ record }">
            <UserStatus :status="record.status" />
        </template>
        <template #column-icon="{ record }">
            <UserText :value="record.icon || '-'" />
        </template>
        <template #column-created_at="{ record }">
            <UserText :value="record.created_at" />
        </template>
    </DataList>

    <!-- 菜单表单弹窗 -->
    <MenuForm
        v-model:open="menuFormOpen"
        :menu="currentMenu"
        @success="handleFormSuccess"
    />

    <!-- 菜单详情弹窗 -->
    <MenuDetail
        v-model:open="menuDetailOpen"
        :menu-id="currentMenuId"
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
    getMenuSearchFields,
    getMenuColumns,
    type PageResponse,
    type ExtraAction,
    type ActionButton,
} from '@/views/pages/components'
import { menuApi, type Menu, type MenuListParams } from '@/api/menu'
import MenuForm from './components/MenuForm.vue'
import MenuDetail from './components/MenuDetail.vue'

// 组件引用
const dataListRef = ref<InstanceType<typeof DataList>>()
// 权限提示
const permissionDenied = ref(false)

// 组件状态
const menuFormOpen = ref(false)
const menuDetailOpen = ref(false)
const currentMenu = ref<Menu | null>(null)
const currentMenuId = ref<string | null>(null)

// 搜索字段配置
const searchFields = getMenuSearchFields()

// 表格列定义
const columns = getMenuColumns()

// 操作按钮
const actions: ActionButton[] = [
    {
        key: 'view',
        label: '查看',
        type: 'link',
        onClick: (record: Menu) => {
            handleViewMenu(record)
        },
    },
    {
        key: 'edit',
        label: '编辑',
        type: 'link',
        onClick: (record: Menu) => {
            handleEditMenu(record)
        },
    },
    {
        key: 'delete',
        label: '删除',
        type: 'link',
        danger: true,
        onClick: (record: Menu) => {
            handleDeleteMenu(record)
        },
    },
]

// 额外操作按钮
const extraActions: ExtraAction[] = [
    {
        key: 'add',
        label: '添加菜单',
        icon: PlusOutlined,
        type: 'primary',
        onClick: () => {
            handleAddMenu()
        },
    },
]

// 获取菜单列表
const fetchMenuList = async (params: any): Promise<PageResponse<Menu[]>> => {
    try {
        const requestParams: MenuListParams = {
            page: params.page || params.page_size || 1,
            pageSize: params.pageSize || params.page_size || 10,
            menu_type: params.menu_type,
        }

        const response = await menuApi.getMenuList(requestParams)

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
        console.error('获取菜单列表失败:', error)
        if (error.response?.status === 403 || error.message?.includes('没有权限')) {
            permissionDenied.value = true
        }
        throw error
    }
}

// 添加菜单
const handleAddMenu = () => {
    currentMenu.value = null
    menuFormOpen.value = true
}

// 查看菜单
const handleViewMenu = (record: Menu) => {
    currentMenuId.value = record.menu_id
    menuDetailOpen.value = true
}

// 编辑菜单
const handleEditMenu = (record: Menu) => {
    currentMenu.value = record
    menuFormOpen.value = true
}

// 删除菜单
const handleDeleteMenu = (record: Menu) => {
    Modal.confirm({
        title: '确认删除',
        content: `确定要删除菜单"${record.title}"吗？此操作不可恢复。`,
        okText: '确定',
        cancelText: '取消',
        okType: 'danger',
        onOk: async () => {
            try {
                await menuApi.deleteMenu(record.menu_id)
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


<template>
    <a-modal
        v-model:open="visible"
        title="用户权限授权"
        :width="900"
        :confirm-loading="loading"
        @ok="handleSubmit"
        @cancel="handleCancel"
    >
        <a-spin :spinning="permissionsLoading">
            <div class="grant-permission-content">
                <a-alert
                    message="提示"
                    description="可以直接为用户分配权限，这些权限会独立于角色权限之外。如果用户同时拥有角色权限和直接权限，两者会合并生效。"
                    type="info"
                    show-icon
                    style="margin-bottom: 16px"
                />

                <a-input
                    v-model:value="searchKeyword"
                    placeholder="搜索权限（名称、标识、路径）"
                    allow-clear
                    style="margin-bottom: 16px"
                >
                    <template #prefix>
                        <SearchOutlined />
                    </template>
                </a-input>

                <div v-if="groupedFilteredPermissions && Object.keys(groupedFilteredPermissions).length > 0" class="permission-tree">
                    <a-collapse
                        v-for="(perms, category) in groupedFilteredPermissions"
                        :key="category"
                        :bordered="false"
                        ghost
                        class="permission-category"
                    >
                        <a-collapse-panel
                            :key="category"
                        >
                            <template #header>
                                <a-tag :color="getCategoryColor(category)" style="margin-right: 8px;">
                                    {{ category }}
                                </a-tag>
                            </template>
                            <template #extra>
                                <span style="color: #8c8c8c; font-size: 12px">
                                    {{ perms.length }} 个
                                </span>
                            </template>
                            <a-table
                                :columns="columns"
                                :data-source="perms"
                                :pagination="false"
                                :row-selection="getRowSelection(perms)"
                                size="small"
                                :row-key="(record) => record.permissionId"
                            >
                                <template #bodyCell="{ column, record }">
                                    <template v-if="column.key === 'method'">
                                        <a-tag :color="getMethodColor(record.method)">
                                            {{ record.method }}
                                        </a-tag>
                                    </template>
                                    <template v-else-if="column.key === 'status'">
                                        <a-badge
                                            :status="record.status === '1' ? 'success' : 'error'"
                                            :text="record.status === '1' ? '启用' : '禁用'"
                                        />
                                    </template>
                                </template>
                            </a-table>
                        </a-collapse-panel>
                    </a-collapse>
                </div>
                <a-empty v-else description="暂无权限" :image="Empty.PRESENTED_IMAGE_SIMPLE" />

                <div v-if="selectedPermissions.length > 0" class="selected-info">
                    <a-divider>已选择 {{ selectedPermissions.length }} 个权限</a-divider>
                </div>
            </div>
        </a-spin>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { message, Empty } from 'ant-design-vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import type { TableColumnsType } from 'ant-design-vue'
import { permissionApi, type Permission } from '@/api/permission'
import { userApi, type PermissionInfo } from '@/api/user'
import { getPermissionCategory } from '@/views/pages/components'

interface Props {
    open: boolean
    userId?: string | null
    currentPermissions?: PermissionInfo[]
}

const props = withDefaults(defineProps<Props>(), {
    open: false,
    userId: null,
    currentPermissions: () => [],
})

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void
    (e: 'success'): void
}>()

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

const loading = ref(false)
const permissionsLoading = ref(false)
const permissions = ref<Permission[]>([])
const selectedPermissionIds = ref<string[]>([])
const searchKeyword = ref('')

const getMethodColor = (method: string) => {
    const colorMap: Record<string, string> = {
        GET: 'blue',
        POST: 'green',
        PUT: 'orange',
        DELETE: 'red',
        PATCH: 'purple',
    }
    return colorMap[method] || 'default'
}

const columns: TableColumnsType = [
    {
        title: '权限名称',
        dataIndex: 'permissionName',
        key: 'permissionName',
        width: 150,
    },
    {
        title: '权限标识',
        dataIndex: 'permissionKey',
        key: 'permissionKey',
        width: 150,
    },
    {
        title: '资源路径',
        dataIndex: 'resourcePath',
        key: 'resourcePath',
        width: 200,
    },
    {
        title: '请求方法',
        key: 'method',
        width: 100,
    },
    {
        title: '状态',
        key: 'status',
        width: 80,
    },
    {
        title: '描述',
        dataIndex: 'description',
        key: 'description',
    },
]

// 获取分类颜色（与权限列表保持一致）
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

// 过滤权限
const filteredPermissions = computed(() => {
    if (!searchKeyword.value) {
        return permissions.value
    }
    const keyword = searchKeyword.value.toLowerCase()
    return permissions.value.filter(
        (p) =>
            p.permissionName?.toLowerCase().includes(keyword) ||
            p.permissionKey?.toLowerCase().includes(keyword) ||
            p.resourcePath?.toLowerCase().includes(keyword)
    )
})

// 按分类组织过滤后的权限（与权限列表保持一致）
const groupedFilteredPermissions = computed(() => {
    if (!filteredPermissions.value || filteredPermissions.value.length === 0) {
        return {}
    }

    const grouped: Record<string, Permission[]> = {}
    
    filteredPermissions.value.forEach((permission) => {
        // 优先使用后端返回的 category 字段，如果没有则从 permissionKey 提取
        const category = permission.category || getPermissionCategory(permission.permissionKey || '')
        
        if (!grouped[category]) {
            grouped[category] = []
        }
        grouped[category].push(permission)
    })

    return grouped
})

// 已选择的权限
const selectedPermissions = computed(() => {
    return permissions.value.filter((p) => selectedPermissionIds.value.includes(p.permissionId))
})

// 获取指定权限列表的行选择配置
const getRowSelection = (perms: Permission[]) => {
    return {
        selectedRowKeys: perms
            .filter((p) => selectedPermissionIds.value.includes(p.permissionId))
            .map((p) => p.permissionId),
        onChange: (selectedRowKeys: any[]) => {
            const permIds = selectedRowKeys as string[]
            // 先移除当前分类下已选择的所有权限
            perms.forEach((p) => {
                const index = selectedPermissionIds.value.indexOf(p.permissionId)
                if (index > -1) {
                    selectedPermissionIds.value.splice(index, 1)
                }
            })
            // 再添加新选择的权限
            selectedPermissionIds.value.push(...permIds)
        },
        getCheckboxProps: (record: Permission) => ({
            disabled: record.status !== '1', // 禁用状态的权限不能选择
        }),
    }
}

// 加载所有权限
const loadPermissions = async () => {
    permissionsLoading.value = true
    try {
        const response = await permissionApi.getAllPermissions()
        if (response.code === 200 || response.code === 0) {
            permissions.value = response.data || []
        }
    } catch (error: any) {
        console.error('加载权限列表失败:', error)
        message.error('加载权限列表失败')
    } finally {
        permissionsLoading.value = false
    }
}

// 初始化已选择的权限
const initSelectedPermissions = () => {
    if (props.currentPermissions && props.currentPermissions.length > 0) {
        selectedPermissionIds.value = props.currentPermissions.map((p) => p.permissionId)
    } else {
        selectedPermissionIds.value = []
    }
}

// 提交授权
const handleSubmit = async () => {
    if (!props.userId) {
        message.error('用户ID不能为空')
        return
    }

    loading.value = true
    try {
        // 调用后端API给用户分配权限
        await userApi.grantPermissions(props.userId, selectedPermissionIds.value)
        message.success('授权成功')
        emit('success')
        handleCancel()
    } catch (error: any) {
        console.error('授权失败:', error)
        message.error(error.message || '授权失败')
    } finally {
        loading.value = false
    }
}

// 取消
const handleCancel = () => {
    selectedPermissionIds.value = []
    searchKeyword.value = ''
    emit('update:open', false)
}

// 监听打开状态
watch(
    () => props.open,
    (val) => {
        if (val) {
            loadPermissions()
            initSelectedPermissions()
        }
    }
)
</script>

<style scoped lang="less">
.grant-permission-content {
    .selected-info {
        margin-top: 16px;
    }
}

.permission-tree {
    max-height: 500px;
    overflow-y: auto;
}

.permission-category {
    margin-bottom: 8px;

    &:last-child {
        margin-bottom: 0;
    }

    :deep(.ant-collapse-item) {
        border: 1px solid #e8e8e8;
        border-radius: 4px;
        margin-bottom: 8px;
        overflow: hidden;

        &:last-child {
            margin-bottom: 0;
        }
    }

    :deep(.ant-collapse-header) {
        background: #fafafa;
        padding: 12px 16px;
        font-weight: 500;
        color: #262626;
        display: flex;
        align-items: center;
    }

    :deep(.ant-collapse-content) {
        background: #fff;
        border-top: 1px solid #e8e8e8;
    }

    :deep(.ant-table) {
        margin: 0;
    }
}

/* 深色模式 */
:deep(.ant-layout-sider-dark),
.dark-mode {
    .permission-category {
        :deep(.ant-collapse-item) {
            border-color: #434343;
            background: #1f1f1f;
        }

        :deep(.ant-collapse-header) {
            background: #262626;
            color: rgba(255, 255, 255, 0.85);
        }

        :deep(.ant-collapse-content) {
            background: #1f1f1f;
            border-top-color: #434343;
        }
    }
}
</style>


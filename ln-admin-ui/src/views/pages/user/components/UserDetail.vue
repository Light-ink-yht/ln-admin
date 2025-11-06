<template>
    <a-modal
        v-model:open="visible"
        title="用户详情"
        :width="900"
        :footer="null"
    >
        <a-spin :spinning="loading">
            <a-descriptions :column="2" bordered>
                <a-descriptions-item label="用户ID">
                    {{ userDetail?.userId || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="手机号">
                    {{ userDetail?.phone || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="邮箱">
                    {{ userDetail?.email || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="昵称">
                    {{ userDetail?.nickname || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="姓名">
                    {{ userDetail?.fullName || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="性别">
                    {{ getGenderText(userDetail?.gender || '') }}
                </a-descriptions-item>
                <a-descriptions-item label="状态">
                    <a-badge
                        :status="userDetail?.status === '1' ? 'success' : 'error'"
                        :text="userDetail?.status === '1' ? '启用' : userDetail?.status === '2' ? '禁用' : '未知'"
                    />
                </a-descriptions-item>
                <a-descriptions-item label="头像">
                    <a-avatar
                        v-if="userDetail?.avatar"
                        :src="userDetail.avatar"
                        :size="64"
                    />
                    <span v-else>-</span>
                </a-descriptions-item>
                <a-descriptions-item label="生日">
                    {{ userDetail?.birthday || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="登录次数">
                    {{ userDetail?.loginCount || 0 }}
                </a-descriptions-item>
                <a-descriptions-item label="最后登录时间">
                    {{ userDetail?.lastLoginTime || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="最后登录IP">
                    {{ userDetail?.lastLoginIp || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="创建时间" :span="2">
                    {{ userDetail?.CreatedAt || '-' }}
                </a-descriptions-item>
                <a-descriptions-item label="备注" :span="2">
                    {{ userDetail?.remarks || '-' }}
                </a-descriptions-item>
            </a-descriptions>

            <!-- 角色信息 -->
            <a-divider>角色信息</a-divider>
            <a-table
                :columns="roleColumns"
                :data-source="userDetail?.roles || []"
                :pagination="false"
                size="small"
                :row-key="(record) => record.roleId"
            >
                <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'status'">
                        <a-badge
                            :status="record.status === '1' ? 'success' : 'error'"
                            :text="record.status === '1' ? '启用' : '禁用'"
                        />
                    </template>
                </template>
            </a-table>

            <!-- 权限信息（超级管理员不显示） -->
            <a-collapse v-if="!isSuperAdmin" :bordered="false" style="margin-top: 16px">
                <a-collapse-panel key="permissions" header="权限信息">
                    <template #extra>
                        <span style="color: #8c8c8c; font-size: 12px">
                            {{ userDetail?.permissions?.length || 0 }} 个权限
                        </span>
                    </template>
                    <div v-if="groupedPermissions && Object.keys(groupedPermissions).length > 0" class="permission-tree">
                        <a-collapse
                            v-for="(permissions, category) in groupedPermissions"
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
                                        {{ permissions.length }} 个
                                    </span>
                                </template>
                                <a-table
                                    :columns="permissionColumns"
                                    :data-source="permissions"
                                    :pagination="false"
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
                                                :status="(record as any).status === '1' ? 'success' : 'error'"
                                                :text="(record as any).status === '1' ? '启用' : '禁用'"
                                            />
                                        </template>
                                    </template>
                                </a-table>
                            </a-collapse-panel>
                        </a-collapse>
                    </div>
                    <a-empty v-else description="暂无权限" :image="Empty.PRESENTED_IMAGE_SIMPLE" />
                </a-collapse-panel>
            </a-collapse>
        </a-spin>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { message, Empty } from 'ant-design-vue'
import type { TableColumnsType } from 'ant-design-vue'
import { userApi, type UserDetailResponse, type RoleInfo, type PermissionInfo } from '@/api/user'
import { getPermissionCategory } from '@/views/pages/components'

interface Props {
    open: boolean
    userId?: string | null
}

const props = withDefaults(defineProps<Props>(), {
    open: false,
    userId: null,
})

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void
}>()

const loading = ref(false)
const userDetail = ref<UserDetailResponse | null>(null)

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

// 从资源路径提取分类
const getCategoryFromResourcePath = (resourcePath: string): string => {
    if (!resourcePath) {
        return '其他'
    }
    // 移除开头的 / 和尾部的 /，然后分割路径
    const path = resourcePath.replace(/^\/+|\/+$/g, '')
    const parts = path.split('/')
    if (parts.length === 0) {
        return '其他'
    }
    const firstPart = parts[0].toLowerCase()
    
    // 根据路径第一部分映射分类
    const categoryMap: Record<string, string> = {
        'user': '用户管理',
        'role': '角色管理',
        'permission': '权限管理',
        'system': '系统配置',
        'sms': '短信管理',
        'menu': '菜单管理',
        'file': '文件管理',
        'log': '系统运维',
        'swagger': '系统运维',
    }
    
    // 特殊处理：system/log 和 system/config
    if (firstPart === 'system') {
        if (parts.length > 1) {
            const secondPart = parts[1].toLowerCase()
            if (secondPart === 'log') {
                return '系统运维'
            }
            if (secondPart === 'config') {
                return '系统配置'
            }
        }
        return '系统配置'
    }
    
    return categoryMap[firstPart] || '其他'
}

// 按分类组织权限（与权限授权弹窗保持一致）
const groupedPermissions = computed(() => {
    if (!userDetail.value?.permissions || userDetail.value.permissions.length === 0) {
        return {}
    }

    const grouped: Record<string, PermissionInfo[]> = {}
    
    userDetail.value.permissions.forEach((permission) => {
        // 优先使用后端返回的 category 字段
        let category = (permission as any).category
        
        // 如果没有 category，尝试从 permissionKey 提取
        if (!category) {
            const permissionKey = (permission as any).permissionKey || ''
            category = getPermissionCategory(permissionKey)
        }
        
        // 如果还是没有，从 resourcePath 提取
        if (!category || category === '其他') {
            const resourcePath = (permission as any).resourcePath || ''
            category = getCategoryFromResourcePath(resourcePath)
        }
        
        // 确保分类不为空
        if (!category) {
            category = '其他'
        }
        
        if (!grouped[category]) {
            grouped[category] = []
        }
        grouped[category].push(permission)
    })

    return grouped
})

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

// 判断是否是超级管理员
const isSuperAdmin = computed(() => {
    if (!userDetail.value) {
        return false
    }
    // 通过手机号判断
    if (userDetail.value.phone === '18797131041') {
        return true
    }
    // 通过角色判断
    if (userDetail.value.roles && Array.isArray(userDetail.value.roles)) {
        return userDetail.value.roles.some((role: any) => role.roleKey === 'super_admin')
    }
    return false
})

const getGenderText = (gender: string) => {
    const genderMap: Record<string, string> = {
        '1': '男',
        '2': '女',
        '3': '未知',
    }
    return genderMap[gender] || '未知'
}

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

const roleColumns: TableColumnsType = [
    {
        title: '角色标识',
        dataIndex: 'roleKey',
        key: 'roleKey',
        width: 150,
    },
    {
        title: '角色名称',
        dataIndex: 'roleName',
        key: 'roleName',
        width: 150,
    },
    {
        title: '状态',
        key: 'status',
        width: 100,
    },
]

const permissionColumns: TableColumnsType = [
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

// 加载用户详情
const loadUserDetail = async () => {
    if (!props.userId) {
        return
    }

    loading.value = true
    try {
        const response = await userApi.getUserDetail(props.userId)
        if (response.code === 200 || response.code === 0) {
            userDetail.value = response.data
        } else {
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        }
    } catch (error: any) {
        console.error('获取用户详情失败:', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    } finally {
        loading.value = false
    }
}

// 监听打开状态和用户ID
watch(
    () => [props.open, props.userId],
    ([open, userId]) => {
        if (open && userId) {
            loadUserDetail()
        } else if (!open) {
            userDetail.value = null
        }
    },
    { immediate: true }
)
</script>

<style scoped lang="less">
:deep(.ant-descriptions-item-label) {
    font-weight: 500;
    background: #fafafa;
}

:deep(.ant-table) {
    margin-top: 16px;
}

:deep(.ant-divider) {
    margin: 24px 0 16px 0;
}

.permission-tree {
    margin-top: 8px;
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
    }

    :deep(.ant-collapse-content) {
        background: #fff;
        border-top: 1px solid #e8e8e8;
    }
}

.permission-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 12px;
    padding: 12px;
}

.permission-card {
    transition: all 0.3s;
    border: 1px solid #e8e8e8;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(var(--color-primary-rgb, 24, 144, 255), 0.15);
    }
}

.permission-content {
    .permission-header {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 12px;

        .method-tag {
            font-weight: 600;
            min-width: 60px;
            text-align: center;
        }

        .permission-name {
            font-weight: 500;
            font-size: 14px;
            color: #262626;
        }
    }

    .permission-info {
        .info-item {
            margin-bottom: 8px;
            font-size: 12px;
            line-height: 1.6;

            &:last-child {
                margin-bottom: 0;
            }

            .info-label {
                color: #8c8c8c;
                margin-right: 8px;
            }

            .info-value {
                color: #262626;

                code {
                    background: #f5f5f5;
                    padding: 2px 6px;
                    border-radius: 3px;
                    font-family: 'Courier New', monospace;
                    font-size: 11px;
                    color: #d46b08;
                }
            }
        }
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

    :deep(.ant-table) {
        margin: 0;
    }
}
</style>


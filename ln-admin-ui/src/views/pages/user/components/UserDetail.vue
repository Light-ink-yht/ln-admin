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

            <!-- 权限信息 -->
            <a-collapse :bordered="false" style="margin-top: 16px">
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
                                :header="getCategoryName(category)"
                            >
                                <template #extra>
                                    <span style="color: #8c8c8c; font-size: 12px">
                                        {{ permissions.length }} 个
                                    </span>
                                </template>
                                <div class="permission-list">
                                    <a-card
                                        v-for="permission in permissions"
                                        :key="permission.permissionId"
                                        size="small"
                                        class="permission-card"
                                    >
                                        <div class="permission-content">
                                            <div class="permission-header">
                                                <a-tag :color="getMethodColor(permission.method)" class="method-tag">
                                                    {{ permission.method }}
                                                </a-tag>
                                                <span class="permission-name">{{ permission.permissionName }}</span>
                                            </div>
                                            <div class="permission-info">
                                                <div class="info-item">
                                                    <span class="info-label">权限标识：</span>
                                                    <span class="info-value">{{ permission.permissionKey }}</span>
                                                </div>
                                                <div class="info-item">
                                                    <span class="info-label">资源路径：</span>
                                                    <code class="info-value">{{ permission.resourcePath }}</code>
                                                </div>
                                                <div v-if="permission.description" class="info-item">
                                                    <span class="info-label">描述：</span>
                                                    <span class="info-value">{{ permission.description }}</span>
                                                </div>
                                            </div>
                                        </div>
                                    </a-card>
                                </div>
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

// 权限分类映射
const categoryMap: Record<string, string> = {
    '/user': '用户管理',
    '/role': '角色管理',
    '/permission': '权限管理',
    '/menu': '菜单管理',
    '/workbench': '工作台管理',
    '/system': '系统管理',
}

// 按分类组织权限
const groupedPermissions = computed(() => {
    if (!userDetail.value?.permissions || userDetail.value.permissions.length === 0) {
        return {}
    }

    const grouped: Record<string, PermissionInfo[]> = {}
    
    userDetail.value.permissions.forEach((permission) => {
        // 从资源路径提取分类（例如：/api/user/list -> user）
        const path = permission.resourcePath || ''
        // 移除 /api 前缀
        const cleanPath = path.replace(/^\/api/, '')
        // 提取第一级路径作为分类
        const parts = cleanPath.split('/').filter((p) => p)
        const category = parts.length > 0 ? '/' + parts[0] : '/other'
        
        if (!grouped[category]) {
            grouped[category] = []
        }
        grouped[category].push(permission)
    })

    return grouped
})

// 获取分类名称
const getCategoryName = (category: string): string => {
    return categoryMap[category] || category || '其他'
}

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
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
        title: '权限标识',
        dataIndex: 'permissionKey',
        key: 'permissionKey',
        width: 150,
    },
    {
        title: '权限名称',
        dataIndex: 'permissionName',
        key: 'permissionName',
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
            message.error(response.msg || response.message || '获取用户详情失败')
        }
    } catch (error: any) {
        console.error('获取用户详情失败:', error)
        message.error(error.message || '获取用户详情失败')
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

    .permission-card {
        border-color: #434343;
        background: #1f1f1f;

        &:hover {
            border-color: var(--color-primary);
            background: #262626;
        }
    }

    .permission-content {
        .permission-header {
            .permission-name {
                color: rgba(255, 255, 255, 0.85);
            }
        }

        .permission-info {
            .info-item {
                .info-label {
                    color: rgba(255, 255, 255, 0.45);
                }

                .info-value {
                    color: rgba(255, 255, 255, 0.85);

                    code {
                        background: #2a2a2a;
                        color: #ffa940;
                    }
                }
            }
        }
    }
}
</style>


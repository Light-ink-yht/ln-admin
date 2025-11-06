<template>
    <a-modal
        title="角色详情"
        v-model:open="visible"
        :width="700"
        :footer="null"
        @cancel="handleCancel"
    >
        <a-spin :spinning="loading">
            <a-descriptions bordered :column="2" size="middle" class="role-detail-descriptions" v-if="roleDetail">
                <a-descriptions-item label="角色ID">{{ roleDetail.roleId || '-' }}</a-descriptions-item>
                <a-descriptions-item label="角色标识">{{ roleDetail.roleKey || '-' }}</a-descriptions-item>
                <a-descriptions-item label="角色名称">{{ roleDetail.roleName || '-' }}</a-descriptions-item>
                <a-descriptions-item label="状态">
                    <a-badge
                        :status="roleDetail.status === '1' ? 'success' : 'error'"
                        :text="roleDetail.status === '1' ? '启用' : roleDetail.status === '2' ? '禁用' : '未知'"
                    />
                </a-descriptions-item>
                <a-descriptions-item label="角色描述" :span="2">{{ roleDetail.description || '-' }}</a-descriptions-item>
                <a-descriptions-item label="创建时间">{{ roleDetail.CreatedAt || roleDetail.createdAt || '-' }}</a-descriptions-item>
                <a-descriptions-item label="更新时间">{{ roleDetail.updatedAt || '-' }}</a-descriptions-item>
            </a-descriptions>
        </a-spin>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { roleApi, type Role } from '@/api/role'

interface Props {
    open: boolean
    roleId: string | null
}

const props = defineProps<Props>()
const emit = defineEmits(['update:open'])

const loading = ref(false)
const roleDetail = ref<Role | null>(null)

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

// 加载角色详情
const loadRoleDetail = async (roleId: string) => {
    loading.value = true
    try {
        const response = await roleApi.getRoleDetail(roleId)
        if (response.code === 200 || response.code === 0) {
            roleDetail.value = response.data
        } else {
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
            roleDetail.value = null
        }
    } catch (error: any) {
        console.error('获取角色详情失败:', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        roleDetail.value = null
    } finally {
        loading.value = false
    }
}

// 监听弹窗打开状态和 roleId 变化
watch(
    () => props.open,
    (newVal) => {
        if (newVal && props.roleId) {
            loadRoleDetail(props.roleId)
        } else {
            roleDetail.value = null
        }
    },
    { immediate: true }
)

watch(
    () => props.roleId,
    (newVal) => {
        if (visible.value && newVal) {
            loadRoleDetail(newVal)
        }
    }
)

const handleCancel = () => {
    emit('update:open', false)
}
</script>

<style scoped lang="less">
.role-detail-descriptions {
    :deep(.ant-descriptions-row) {
        .ant-descriptions-item-label {
            width: 120px;
        }
        .ant-descriptions-item-content {
            width: calc(50% - 120px);
        }
    }
}
</style>


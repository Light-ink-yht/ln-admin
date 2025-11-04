<template>
    <a-modal
        title="权限详情"
        v-model:open="visible"
        :width="700"
        :footer="null"
        @cancel="handleCancel"
    >
        <a-spin :spinning="loading">
            <a-descriptions bordered :column="2" size="middle" class="permission-detail-descriptions" v-if="permissionDetail">
                <a-descriptions-item label="权限ID">{{ permissionDetail.permissionId || '-' }}</a-descriptions-item>
                <a-descriptions-item label="权限标识">{{ permissionDetail.permissionKey || '-' }}</a-descriptions-item>
                <a-descriptions-item label="权限名称">{{ permissionDetail.permissionName || '-' }}</a-descriptions-item>
                <a-descriptions-item label="状态">
                    <a-badge
                        :status="permissionDetail.status === '1' ? 'success' : 'error'"
                        :text="permissionDetail.status === '1' ? '启用' : permissionDetail.status === '2' ? '禁用' : '未知'"
                    />
                </a-descriptions-item>
                <a-descriptions-item label="资源路径" :span="2">{{ permissionDetail.resourcePath || '-' }}</a-descriptions-item>
                <a-descriptions-item label="请求方法">{{ permissionDetail.method || '-' }}</a-descriptions-item>
                <a-descriptions-item label="权限描述" :span="1">{{ permissionDetail.description || '-' }}</a-descriptions-item>
                <a-descriptions-item label="创建时间">{{ permissionDetail.CreatedAt || permissionDetail.createdAt || '-' }}</a-descriptions-item>
                <a-descriptions-item label="更新时间">{{ permissionDetail.updatedAt || '-' }}</a-descriptions-item>
                <a-descriptions-item label="创建人">{{ permissionDetail.creatorName || permissionDetail.creatorId || '-' }}</a-descriptions-item>
                <a-descriptions-item label="修改人">{{ permissionDetail.modifierName || permissionDetail.modifierId || '-' }}</a-descriptions-item>
            </a-descriptions>
        </a-spin>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { permissionApi, type Permission } from '@/api/permission'

interface Props {
    open: boolean
    permissionId: string | null
}

const props = defineProps<Props>()
const emit = defineEmits(['update:open'])

const loading = ref(false)
const permissionDetail = ref<Permission | null>(null)

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

// 加载权限详情
const loadPermissionDetail = async (permissionId: string) => {
    loading.value = true
    try {
        const response = await permissionApi.getPermissionDetail(permissionId)
        if (response.code === 200 || response.code === 0) {
            permissionDetail.value = response.data
        } else {
            message.error(response.message || '获取权限详情失败')
            permissionDetail.value = null
        }
    } catch (error: any) {
        console.error('获取权限详情失败:', error)
        if (error.response?.status !== 404) {
            message.error(error.message || '获取权限详情失败')
        }
        permissionDetail.value = null
    } finally {
        loading.value = false
    }
}

// 监听弹窗打开状态和 permissionId 变化
watch(
    () => props.open,
    (newVal) => {
        if (newVal && props.permissionId) {
            loadPermissionDetail(props.permissionId)
        } else {
            permissionDetail.value = null
        }
    },
    { immediate: true }
)

watch(
    () => props.permissionId,
    (newVal) => {
        if (visible.value && newVal) {
            loadPermissionDetail(newVal)
        }
    }
)

const handleCancel = () => {
    emit('update:open', false)
}
</script>

<style scoped lang="less">
.permission-detail-descriptions {
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


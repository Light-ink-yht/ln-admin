<template>
    <a-modal
        :title="isEdit ? '编辑权限' : '添加权限'"
        v-model:open="visible"
        :width="700"
        :confirm-loading="loading"
        @ok="handleSubmit"
        @cancel="handleCancel"
    >
        <a-form
            ref="formRef"
            :model="formData"
            :rules="rules"
            :label-col="{ span: 6 }"
            :wrapper-col="{ span: 18 }"
        >
            <a-form-item label="权限标识" name="permissionKey" v-if="!isEdit">
                <a-input v-model:value="formData.permissionKey" placeholder="请输入权限标识（英文，唯一，如：user:list）" />
            </a-form-item>

            <a-form-item label="权限名称" name="permissionName">
                <a-input v-model:value="formData.permissionName" placeholder="请输入权限名称" />
            </a-form-item>

            <a-form-item label="资源路径" name="resourcePath">
                <a-input v-model:value="formData.resourcePath" placeholder="请输入资源路径（如：/api/user/list）" />
            </a-form-item>

            <a-form-item label="请求方法" name="method">
                <a-select v-model:value="formData.method" placeholder="请选择请求方法">
                    <a-select-option value="GET">GET</a-select-option>
                    <a-select-option value="POST">POST</a-select-option>
                    <a-select-option value="PUT">PUT</a-select-option>
                    <a-select-option value="DELETE">DELETE</a-select-option>
                    <a-select-option value="PATCH">PATCH</a-select-option>
                </a-select>
            </a-form-item>

            <a-form-item label="权限描述" name="description">
                <a-textarea v-model:value="formData.description" placeholder="请输入权限描述" :rows="3" />
            </a-form-item>

            <a-form-item label="状态" name="status">
                <a-select v-model:value="formData.status" placeholder="请选择状态">
                    <a-select-option value="1">启用</a-select-option>
                    <a-select-option value="2">禁用</a-select-option>
                </a-select>
            </a-form-item>
        </a-form>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import { permissionApi, type CreatePermissionRequest, type UpdatePermissionRequest, type Permission } from '@/api/permission'

interface Props {
    open: boolean
    permission: Permission | null // null for add, Permission for edit
}

const props = defineProps<Props>()
const emit = defineEmits(['update:open', 'success'])

const formRef = ref<FormInstance>()
const loading = ref(false)

const isEdit = computed(() => !!props.permission)

const formData = reactive<CreatePermissionRequest & UpdatePermissionRequest>({
    permissionKey: '',
    permissionName: '',
    resourcePath: '',
    method: 'GET',
    description: '',
    status: '1', // Default to enabled
})

const rules = computed(() => ({
    permissionKey: [{ required: !isEdit.value, message: '请输入权限标识', trigger: 'blur' }],
    permissionName: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
    resourcePath: [{ required: true, message: '请输入资源路径', trigger: 'blur' }],
    method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
}))

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

// 初始化表单数据
const initFormData = () => {
    if (props.permission) {
        // 编辑模式
        formData.permissionName = props.permission.permissionName || ''
        formData.resourcePath = props.permission.resourcePath || ''
        formData.method = props.permission.method || 'GET'
        formData.description = props.permission.description || ''
        formData.status = props.permission.status || '1'
    } else {
        // 添加模式
        formData.permissionKey = ''
        formData.permissionName = ''
        formData.resourcePath = ''
        formData.method = 'GET'
        formData.description = ''
        formData.status = '1'
    }
}

// 提交表单
const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        loading.value = true

        const dataToSubmit: CreatePermissionRequest | UpdatePermissionRequest = {
            permissionName: formData.permissionName,
            resourcePath: formData.resourcePath,
            method: formData.method,
            description: formData.description || undefined,
            status: formData.status,
        }

        if (isEdit.value && props.permission) {
            // 编辑权限
            await permissionApi.updatePermission(props.permission.permissionId, dataToSubmit as UpdatePermissionRequest)
            message.success('更新权限成功')
        } else {
            // 添加权限
            await permissionApi.createPermission({
                ...dataToSubmit,
                permissionKey: formData.permissionKey,
            } as CreatePermissionRequest)
            message.success('添加权限成功')
        }
        emit('success')
        visible.value = false
    } catch (error: any) {
        console.error('提交权限表单失败:', error)
        if (error.response?.status !== 404) {
            message.error(error.message || '提交失败')
        }
    } finally {
        loading.value = false
    }
}

// 取消
const handleCancel = () => {
    visible.value = false
    formRef.value?.resetFields()
}

// 监听弹窗打开状态
watch(
    () => props.open,
    (newVal) => {
        if (newVal) {
            initFormData()
        } else {
            formRef.value?.resetFields()
        }
    },
    { immediate: true }
)
</script>


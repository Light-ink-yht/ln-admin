<template>
    <a-modal
        :title="isEdit ? '编辑角色' : '添加角色'"
        v-model:open="visible"
        :width="600"
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
            <a-form-item label="角色标识" name="roleKey" v-if="!isEdit">
                <a-input v-model:value="formData.roleKey" placeholder="请输入角色标识（英文，唯一）" />
            </a-form-item>

            <a-form-item label="角色名称" name="roleName">
                <a-input v-model:value="formData.roleName" placeholder="请输入角色名称" />
            </a-form-item>

            <a-form-item label="角色描述" name="description">
                <a-textarea v-model:value="formData.description" placeholder="请输入角色描述" :rows="3" />
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
import { roleApi, type CreateRoleRequest, type UpdateRoleRequest, type Role } from '@/api/role'

interface Props {
    open: boolean
    role: Role | null // null for add, Role for edit
}

const props = defineProps<Props>()
const emit = defineEmits(['update:open', 'success'])

const formRef = ref<FormInstance>()
const loading = ref(false)

const isEdit = computed(() => !!props.role)

const formData = reactive<CreateRoleRequest & UpdateRoleRequest>({
    roleKey: '',
    roleName: '',
    description: '',
    status: '1', // Default to enabled
})

const rules = computed(() => ({
    roleKey: [{ required: !isEdit.value, message: '请输入角色标识', trigger: 'blur' }],
    roleName: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
}))

const visible = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

// 初始化表单数据
const initFormData = () => {
    if (props.role) {
        // 编辑模式
        formData.roleName = props.role.roleName || ''
        formData.description = props.role.description || ''
        formData.status = props.role.status || '1'
    } else {
        // 添加模式
        formData.roleKey = ''
        formData.roleName = ''
        formData.description = ''
        formData.status = '1'
    }
}

// 提交表单
const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        loading.value = true

        const dataToSubmit: CreateRoleRequest | UpdateRoleRequest = {
            roleName: formData.roleName,
            description: formData.description || undefined,
            status: formData.status,
        }

        if (isEdit.value && props.role) {
            // 编辑角色
            await roleApi.updateRole(props.role.roleId, dataToSubmit as UpdateRoleRequest)
            message.success('更新角色成功')
        } else {
            // 添加角色
            await roleApi.createRole({
                ...dataToSubmit,
                roleKey: formData.roleKey,
            } as CreateRoleRequest)
            message.success('添加角色成功')
        }
        emit('success')
        visible.value = false
    } catch (error: any) {
        console.error('提交角色表单失败:', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
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


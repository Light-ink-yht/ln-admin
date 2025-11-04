<template>
    <a-modal
        v-model:open="modalOpen"
        :title="isEdit ? '编辑短信模板' : '添加短信模板'"
        width="800px"
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
            <a-form-item label="模板类型" name="type">
                <a-input
                    v-model:value="formData.type"
                    placeholder="请输入模板类型，如：register、forgot"
                    :disabled="isEdit"
                />
            </a-form-item>
            <a-form-item label="模板名称" name="templateName">
                <a-input
                    v-model:value="formData.templateName"
                    placeholder="请输入模板名称"
                />
            </a-form-item>
            <a-form-item label="腾讯云模板ID" name="tencentTemplateId">
                <a-input
                    v-model:value="formData.tencentTemplateId"
                    placeholder="请输入腾讯云模板ID"
                />
            </a-form-item>
            <a-form-item label="短信签名" name="title">
                <a-input
                    v-model:value="formData.title"
                    placeholder="请输入短信签名/标题"
                />
            </a-form-item>
            <a-form-item label="模板内容" name="content">
                <a-textarea
                    v-model:value="formData.content"
                    placeholder="请输入模板内容，支持占位符如：{code}"
                    :rows="4"
                />
            </a-form-item>
            <a-form-item label="模板描述" name="description">
                <a-textarea
                    v-model:value="formData.description"
                    placeholder="请输入模板描述"
                    :rows="3"
                />
            </a-form-item>
            <a-form-item label="状态" name="status">
                <a-radio-group v-model:value="formData.status">
                    <a-radio value="1">启用</a-radio>
                    <a-radio value="2">禁用</a-radio>
                </a-radio-group>
            </a-form-item>
        </a-form>
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from 'vue'
import { message } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import { smsApi, type SMSTemplate, type CreateSMSTemplateRequest, type UpdateSMSTemplateRequest } from '@/api/sms'

const props = defineProps<{
    open: boolean
    template?: SMSTemplate | null
}>()

const emit = defineEmits<{
    'update:open': [value: boolean]
    success: []
}>()

const modalOpen = ref(false)
const loading = ref(false)
const formRef = ref<FormInstance>()
const isEdit = ref(false)

const formData = reactive<CreateSMSTemplateRequest & { status?: string }>({
    type: '',
    templateName: '',
    tencentTemplateId: '',
    title: '',
    content: '',
    description: '',
    status: '1',
})

const rules = {
    type: [{ required: true, message: '请输入模板类型', trigger: 'blur' }],
    templateName: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
    tencentTemplateId: [{ required: true, message: '请输入腾讯云模板ID', trigger: 'blur' }],
    title: [{ required: true, message: '请输入短信签名', trigger: 'blur' }],
    content: [{ required: true, message: '请输入模板内容', trigger: 'blur' }],
}

// 监听 open 变化
watch(
    () => props.open,
    (newVal) => {
        modalOpen.value = newVal
        if (newVal) {
            initFormData()
        }
    },
    { immediate: true }
)

// 监听 modalOpen 变化，同步到父组件
watch(modalOpen, (newVal) => {
    emit('update:open', newVal)
})

// 初始化表单数据
const initFormData = () => {
    if (props.template) {
        isEdit.value = true
        Object.assign(formData, {
            type: props.template.type,
            templateName: props.template.templateName,
            tencentTemplateId: props.template.tencentTemplateId,
            title: props.template.title,
            content: props.template.content,
            description: props.template.description || '',
            status: props.template.status || '1',
        })
    } else {
        isEdit.value = false
        Object.assign(formData, {
            type: '',
            templateName: '',
            tencentTemplateId: '',
            title: '',
            content: '',
            description: '',
            status: '1',
        })
    }
    formRef.value?.resetFields()
}

// 提交表单
const handleSubmit = async () => {
    try {
        await formRef.value?.validate()
        loading.value = true

        if (isEdit.value && props.template) {
            const updateData: UpdateSMSTemplateRequest = {
                templateName: formData.templateName,
                tencentTemplateId: formData.tencentTemplateId,
                title: formData.title,
                content: formData.content,
                description: formData.description,
                status: formData.status,
            }
            await smsApi.updateTemplate(props.template.templateId, updateData)
            message.success('更新成功')
        } else {
            const createData: CreateSMSTemplateRequest = {
                type: formData.type,
                templateName: formData.templateName,
                tencentTemplateId: formData.tencentTemplateId,
                title: formData.title,
                content: formData.content,
                description: formData.description,
                status: formData.status,
            }
            await smsApi.createTemplate(createData)
            message.success('创建成功')
        }

        emit('success')
        handleCancel()
    } catch (error: any) {
        console.error('提交失败:', error)
        message.error(error.message || '提交失败')
    } finally {
        loading.value = false
    }
}

// 取消
const handleCancel = () => {
    modalOpen.value = false
    formRef.value?.resetFields()
}
</script>


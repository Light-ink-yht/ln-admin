<template>
    <a-modal
        v-model:open="modalOpen"
        title="短信模板详情"
        width="800px"
        :footer="null"
    >
        <a-descriptions :column="2" bordered v-if="template">
            <a-descriptions-item label="模板ID" :span="2">
                <code>{{ template.templateId }}</code>
            </a-descriptions-item>
            <a-descriptions-item label="模板类型">
                {{ template.type }}
            </a-descriptions-item>
            <a-descriptions-item label="模板名称">
                {{ template.templateName }}
            </a-descriptions-item>
            <a-descriptions-item label="腾讯云模板ID">
                {{ template.tencentTemplateId }}
            </a-descriptions-item>
            <a-descriptions-item label="短信签名">
                {{ template.title }}
            </a-descriptions-item>
            <a-descriptions-item label="状态">
                <UserStatus :status="template.status" />
            </a-descriptions-item>
            <a-descriptions-item label="模板内容" :span="2">
                <a-typography-paragraph style="margin: 0; white-space: pre-wrap">
                    {{ template.content }}
                </a-typography-paragraph>
            </a-descriptions-item>
            <a-descriptions-item label="模板描述" :span="2">
                {{ template.description || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="创建人">
                {{ template.creatorId || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="修改人">
                {{ template.modifierId || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="创建时间">
                {{ template.createdAt }}
            </a-descriptions-item>
            <a-descriptions-item label="更新时间">
                {{ template.updatedAt }}
            </a-descriptions-item>
        </a-descriptions>
        <a-empty v-else description="暂无数据" />
    </a-modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { UserStatus } from '@/views/pages/components'
import { smsApi, type SMSTemplate } from '@/api/sms'

const props = defineProps<{
    open: boolean
    templateId?: string | null
}>()

const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const modalOpen = ref(false)
const template = ref<SMSTemplate | null>(null)

// 监听 open 变化
watch(
    () => props.open,
    (newVal) => {
        modalOpen.value = newVal
        if (newVal && props.templateId) {
            loadTemplateDetail()
        }
    },
    { immediate: true }
)

// 监听 modalOpen 变化，同步到父组件
watch(modalOpen, (newVal) => {
    emit('update:open', newVal)
})

// 加载模板详情
const loadTemplateDetail = async () => {
    if (!props.templateId) return

    try {
        const response = await smsApi.getTemplateDetail(props.templateId)
        if (response.code === 200 || response.code === 0) {
            template.value = response.data
        }
    } catch (error: any) {
        console.error('加载模板详情失败:', error)
        template.value = null
    }
}
</script>


<template>
    <a-modal
        v-model:open="visible"
        :title="file?.original_name || '文件预览'"
        :width="800"
        :footer="null"
        @cancel="handleCancel"
    >
        <div class="file-preview-modal">
            <!-- 图片预览 -->
            <a-image
                v-if="isImage"
                :src="file?.file_url"
                :preview="{ mask: '预览' }"
                style="max-width: 100%; max-height: 600px;"
            />

            <!-- PDF预览 -->
            <iframe
                v-else-if="isPdf"
                :src="file?.file_url"
                style="width: 100%; height: 600px; border: none;"
            />

            <!-- 视频预览 -->
            <video
                v-else-if="isVideo"
                :src="file?.file_url"
                controls
                style="width: 100%; max-height: 600px;"
            />

            <!-- 音频预览 -->
            <audio
                v-else-if="isAudio"
                :src="file?.file_url"
                controls
                style="width: 100%;"
            />

            <!-- 其他文件类型 -->
            <div v-else class="file-info">
                <a-descriptions :column="2" bordered>
                    <a-descriptions-item label="文件名">
                        {{ file?.original_name }}
                    </a-descriptions-item>
                    <a-descriptions-item label="文件大小">
                        {{ formatFileSize(file?.file_size || 0) }}
                    </a-descriptions-item>
                    <a-descriptions-item label="文件类型">
                        {{ file?.mime_type }}
                    </a-descriptions-item>
                    <a-descriptions-item label="扩展名">
                        {{ file?.extension }}
                    </a-descriptions-item>
                    <a-descriptions-item label="存储类型" :span="2">
                        {{ file?.storage_type === 's3' ? 'S3存储' : '本地存储' }}
                    </a-descriptions-item>
                </a-descriptions>
                <div class="file-actions" style="margin-top: 24px; text-align: center;">
                    <a-button type="primary" @click="handleDownload">
                        <DownloadOutlined /> 下载文件
                    </a-button>
                </div>
            </div>
        </div>
    </a-modal>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { DownloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { FileResponse } from '@/api/file'
import { fileApi } from '@/api/file'
import { formatFileSize } from '../utils'

const props = defineProps<{
    open: boolean
    file: FileResponse | null
}>()

const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const visible = computed({
    get: () => props.open,
    set: (value) => emit('update:open', value),
})

const isImage = computed(() => props.file?.category === 'image')
const isPdf = computed(() => props.file?.extension?.toLowerCase() === 'pdf')
const isVideo = computed(() => props.file?.category === 'video')
const isAudio = computed(() => props.file?.category === 'audio')

const handleCancel = () => {
    visible.value = false
}

const handleDownload = async () => {
    if (!props.file) return

    try {
        const blob = await fileApi.downloadFile(props.file.file_id)
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = props.file.original_name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        message.success('下载成功')
    } catch (error: any) {
        console.error('下载文件失败:', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    }
}
</script>

<style scoped>
.file-preview-modal {
    min-height: 200px;
}

.file-info {
    padding: 20px;
}
</style>


<template>
    <div class="file-upload">
        <a-upload
            v-model:file-list="fileList"
            :multiple="multiple"
            :before-upload="beforeUpload"
            :custom-request="customRequest"
            :show-upload-list="showUploadList"
            :accept="accept"
            @remove="handleRemove"
        >
            <a-button v-if="!drag" :loading="uploading">
                <UploadOutlined />
                {{ buttonText }}
            </a-button>
            <a-upload-dragger v-else :disabled="uploading">
                <p class="ant-upload-drag-icon">
                    <InboxOutlined />
                </p>
                <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
                <p class="ant-upload-hint">
                    支持单个或批量上传
                </p>
            </a-upload-dragger>
        </a-upload>
    </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined, InboxOutlined } from '@ant-design/icons-vue'
import type { UploadFile, UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

interface Props {
    multiple?: boolean
    drag?: boolean
    accept?: string
    buttonText?: string
    showUploadList?: boolean
    maxSize?: number // MB
}

const props = withDefaults(defineProps<Props>(), {
    multiple: true,
    drag: false,
    accept: '*',
    buttonText: '上传文件',
    showUploadList: true,
    maxSize: 100,
})

const emit = defineEmits<{
    success: [file: FileResponse]
    error: [error: Error]
}>()

const fileList = ref<UploadFile[]>([])
const uploading = ref(false)

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        message.warning(`文件大小不能超过 ${props.maxSize}MB`)
        return false
    }
    return true
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError, onProgress } = options

    uploading.value = true

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            message.success('上传成功')
            onSuccess?.(response.data)
            emit('success', response.data as FileResponse)
        } else {
            const error = new Error(response.msg || '上传失败')
            onError?.(error)
            emit('error', error)
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        }
    } catch (error: any) {
        console.error('上传文件失败:', error)
        onError?.(error)
        emit('error', error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    } finally {
        uploading.value = false
    }
}

const handleRemove = (file: UploadFile) => {
    const index = fileList.value.findIndex((item) => item.uid === file.uid)
    if (index > -1) {
        fileList.value.splice(index, 1)
    }
}
</script>

<style scoped>
.file-upload {
    width: 100%;
}
</style>


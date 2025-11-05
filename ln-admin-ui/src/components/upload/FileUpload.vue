<template>
    <div class="file-upload">
        <a-upload
            v-model:file-list="fileList"
            :multiple="multiple"
            :accept="accept"
            :before-upload="beforeUpload"
            :custom-request="customRequest"
            :show-upload-list="showUploadList"
            :max-count="maxCount"
            :drag="drag"
        >
            <template v-if="drag">
                <a-upload-dragger :disabled="uploading">
                    <p class="ant-upload-drag-icon">
                        <InboxOutlined />
                    </p>
                    <p class="ant-upload-text">{{ placeholder }}</p>
                    <p class="ant-upload-hint">
                        {{ hint || `支持单个或批量上传，单个文件不超过 ${maxSize}MB` }}
                    </p>
                </a-upload-dragger>
            </template>
            <template v-else>
                <a-button :loading="uploading">
                    <template #icon><UploadOutlined /></template>
                    {{ buttonText }}
                </a-button>
            </template>
        </a-upload>

        <div v-if="showFileList && fileList.length > 0" class="file-list">
            <div
                v-for="(file, index) in fileList"
                :key="file.uid"
                class="file-item"
                :class="{ 'is-error': file.status === 'error' }"
            >
                <div class="file-info">
                    <FileOutlined class="file-icon" />
                    <span class="file-name" :title="file.name">{{ file.name }}</span>
                    <span class="file-size">{{ formatFileSize(file.size || 0) }}</span>
                </div>
                <div class="file-actions">
                    <a-button
                        v-if="file.status === 'done' && file.response"
                        type="link"
                        size="small"
                        @click="downloadFile(file.response)"
                    >
                        下载
                    </a-button>
                    <a-button
                        type="link"
                        danger
                        size="small"
                        @click="removeFile(index)"
                    >
                        删除
                    </a-button>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined, InboxOutlined, FileOutlined } from '@ant-design/icons-vue'
import type { UploadFile, UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

interface Props {
    modelValue?: UploadFile[]
    multiple?: boolean
    drag?: boolean
    accept?: string
    maxSize?: number // MB
    maxCount?: number
    buttonText?: string
    placeholder?: string
    hint?: string
    showUploadList?: boolean
    showFileList?: boolean
    autoUpload?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: () => [],
    multiple: true,
    drag: false,
    accept: '*',
    maxSize: 100,
    maxCount: 10,
    buttonText: '上传文件',
    placeholder: '点击或拖拽文件到此区域上传',
    showUploadList: true,
    showFileList: true,
    autoUpload: true,
})

const emit = defineEmits<{
    'update:modelValue': [value: UploadFile[]]
    'change': [value: UploadFile[]]
    'upload-success': [file: FileResponse, uploadFile: UploadFile]
    'upload-error': [error: Error, uploadFile: UploadFile]
}>()

const fileList = computed({
    get: () => props.modelValue || [],
    set: (val) => {
        emit('update:modelValue', val)
        emit('change', val)
    },
})

const uploading = ref(false)

const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        message.error(`文件大小不能超过 ${props.maxSize}MB`)
        return false
    }

    // 检查数量限制
    if (props.multiple && fileList.value.length >= props.maxCount) {
        message.error(`最多只能上传 ${props.maxCount} 个文件！`)
        return false
    }

    // 如果是自动上传，返回 true
    if (props.autoUpload) {
        return true
    }

    // 否则添加到列表但不自动上传
    const newFile: UploadFile = {
        uid: Date.now().toString(),
        name: file.name,
        size: file.size,
        type: file.type,
        status: 'ready',
        originFileObj: file,
    }
    fileList.value = [...fileList.value, newFile]
    return false
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError, onProgress } = options

    uploading.value = true

    // 更新文件状态
    const uploadFile = fileList.value.find((f) => f.uid === (file as any).uid)
    if (uploadFile) {
        uploadFile.status = 'uploading'
    }

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            const fileData = response.data as FileResponse
            
            if (uploadFile) {
                uploadFile.status = 'done'
                uploadFile.response = fileData
            }

            onSuccess?.(fileData)
            emit('upload-success', fileData, uploadFile as UploadFile)
            message.success('上传成功')
        } else {
            const error = new Error(response.msg || '上传失败')
            if (uploadFile) {
                uploadFile.status = 'error'
                uploadFile.response = { error: error.message }
            }
            onError?.(error)
            emit('upload-error', error, uploadFile as UploadFile)
            message.error(response.msg || '上传失败')
        }
    } catch (error: any) {
        console.error('上传文件失败:', error)
        if (uploadFile) {
            uploadFile.status = 'error'
            uploadFile.response = { error: error.message }
        }
        onError?.(error)
        emit('upload-error', error, uploadFile as UploadFile)
        message.error(error.message || '上传失败')
    } finally {
        uploading.value = false
    }
}

const removeFile = (index: number) => {
    const files = [...fileList.value]
    files.splice(index, 1)
    fileList.value = files
}

const downloadFile = async (fileResponse: FileResponse) => {
    try {
        const blob = await fileApi.downloadFile(fileResponse.file_id)
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = fileResponse.original_name || fileResponse.file_name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
    } catch (error: any) {
        message.error(error.message || '下载失败')
    }
}
</script>

<style scoped lang="less">
.file-upload {
    width: 100%;
}

.file-list {
    margin-top: 16px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.file-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    background: #fafafa;
    border: 1px solid #d9d9d9;
    border-radius: 4px;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        background: #f0f7ff;
    }

    &.is-error {
        border-color: #ff4d4f;
        background: #fff2f0;
    }
}

.file-info {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 0;
}

.file-icon {
    font-size: 16px;
    color: rgba(0, 0, 0, 0.45);
}

.file-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: rgba(0, 0, 0, 0.85);
}

.file-size {
    font-size: 12px;
    color: rgba(0, 0, 0, 0.45);
    white-space: nowrap;
}

.file-actions {
    display: flex;
    gap: 8px;
}

:deep(.ant-upload) {
    width: 100%;
}

:deep(.ant-upload-drag) {
    background: #fafafa;
    border: 1px dashed #d9d9d9;
    border-radius: 8px;

    &:hover {
        border-color: var(--color-primary);
    }
}
</style>


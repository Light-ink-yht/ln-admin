<template>
    <div class="archive-upload">
        <div v-if="modelValue && modelValue.length > 0" class="archive-list">
            <div
                v-for="(archive, index) in modelValue"
                :key="index"
                class="archive-item"
                :class="{ 'is-error': archive.error }"
            >
                <div class="archive-wrapper">
                    <div class="archive-icon">
                        <FileZipOutlined />
                    </div>
                    <div class="archive-info">
                        <div class="archive-name" :title="archive.name">{{ archive.name || '压缩包' }}</div>
                        <div class="archive-meta">
                            <span v-if="archive.size">{{ formatFileSize(archive.size) }}</span>
                            <span v-if="archive.type" class="archive-type">{{ archive.type.toUpperCase() }}</span>
                        </div>
                    </div>
                    <div class="archive-actions">
                        <a-button
                            v-if="archive.file_id && archive.url"
                            type="link"
                            size="small"
                            @click="downloadArchive(archive)"
                        >
                            <template #icon><DownloadOutlined /></template>
                            下载
                        </a-button>
                        <a-button
                            type="link"
                            danger
                            size="small"
                            @click="removeArchive(index)"
                        >
                            <template #icon><DeleteOutlined /></template>
                            删除
                        </a-button>
                    </div>
                    <div v-if="archive.uploading" class="upload-progress">
                        <a-progress
                            :percent="archive.progress || 0"
                            :status="archive.error ? 'exception' : 'active'"
                        />
                    </div>
                </div>
                <div v-if="archive.error" class="error-message">
                    {{ archive.error }}
                </div>
            </div>

            <div v-if="canAddMore" class="upload-trigger" @click="triggerUpload">
                <a-upload
                    ref="uploadRef"
                    :multiple="multiple"
                    :accept="accept"
                    :show-upload-list="false"
                    :before-upload="beforeUpload"
                    :custom-request="customRequest"
                    style="display: none;"
                />
                <div class="upload-placeholder">
                    <PlusOutlined />
                    <div class="upload-text">上传压缩包</div>
                </div>
            </div>
        </div>

        <div v-else class="upload-area" @click="triggerUpload">
            <a-upload
                ref="uploadRef"
                :multiple="multiple"
                :accept="accept"
                :show-upload-list="false"
                :before-upload="beforeUpload"
                :custom-request="customRequest"
                :drag="drag"
            >
                <div v-if="drag" class="upload-dragger">
                    <p class="ant-upload-drag-icon">
                        <InboxOutlined />
                    </p>
                    <p class="ant-upload-text">{{ placeholder }}</p>
                    <p class="ant-upload-hint">
                        支持 {{ acceptText }}，单个文件不超过 {{ maxSize }}MB
                    </p>
                </div>
                <div v-else class="upload-button">
                    <FileZipOutlined />
                    <div class="upload-text">{{ placeholder }}</div>
                </div>
            </a-upload>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import {
    PlusOutlined,
    DeleteOutlined,
    DownloadOutlined,
    FileZipOutlined,
    InboxOutlined,
} from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

/**
 * 压缩包项接口
 */
export interface ArchiveItem {
    url?: string
    file_id?: string
    name?: string
    size?: number
    type?: string
    uploading?: boolean
    progress?: number
    error?: string
    file?: File
}

interface Props {
    modelValue?: ArchiveItem[]
    multiple?: boolean
    drag?: boolean
    maxSize?: number // MB
    maxCount?: number
    accept?: string
    placeholder?: string
    autoUpload?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: () => [],
    multiple: false,
    drag: false,
    maxSize: 500,
    maxCount: 10,
    accept: '.zip,.rar,.7z,.tar,.gz',
    placeholder: '上传压缩包',
    autoUpload: true,
})

const emit = defineEmits<{
    'update:modelValue': [value: ArchiveItem[]]
    'change': [value: ArchiveItem[]]
    'upload-success': [file: FileResponse, index: number]
    'upload-error': [error: Error, index: number]
}>()

const uploadRef = ref()

const acceptText = computed(() => {
    if (props.accept.includes('.zip')) {
        return 'ZIP、RAR、7Z、TAR、GZ 等压缩格式'
    }
    return props.accept
})

const canAddMore = computed(() => {
    if (!props.multiple) return false
    return (props.modelValue?.length || 0) < props.maxCount
})

const currentArchives = computed({
    get: () => props.modelValue || [],
    set: (val) => {
        emit('update:modelValue', val)
        emit('change', val)
    },
})

const triggerUpload = () => {
    uploadRef.value?.$el.querySelector('input')?.click()
}

const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const getArchiveType = (fileName: string): string => {
    const ext = fileName.split('.').pop()?.toLowerCase() || ''
    return ext
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    // 检查文件类型
    const archiveTypes = ['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz']
    const fileExt = file.name.split('.').pop()?.toLowerCase() || ''
    
    if (!archiveTypes.includes(fileExt)) {
        message.warning('只能上传压缩包文件！支持 ZIP、RAR、7Z、TAR、GZ 等格式')
        return false
    }

    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        message.warning(`压缩包大小不能超过 ${props.maxSize}MB！`)
        return false
    }

    // 检查数量限制
    if (props.multiple && (currentArchives.value.length || 0) >= props.maxCount) {
        message.warning(`最多只能上传 ${props.maxCount} 个压缩包！`)
        return false
    }

    // 如果是自动上传，返回 true
    if (props.autoUpload) {
        return true
    }

    // 否则添加到列表但不自动上传
    const newArchive: ArchiveItem = {
        file,
        name: file.name,
        size: file.size,
        type: getArchiveType(file.name),
        uploading: false,
    }

    currentArchives.value = [...(currentArchives.value || []), newArchive]
    return false
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options

    const newArchive: ArchiveItem = {
        file: file as File,
        name: (file as File).name,
        size: (file as File).size,
        type: getArchiveType((file as File).name),
        uploading: true,
        progress: 0,
    }

    const index = currentArchives.value.length
    currentArchives.value = [...(currentArchives.value || []), newArchive]

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            const fileData = response.data as FileResponse
            newArchive.uploading = false
            newArchive.progress = 100
            newArchive.url = fileData.file_url || fileData.file_path
            newArchive.file_id = fileData.file_id
            newArchive.file = undefined

            currentArchives.value[index] = { ...newArchive }
            onSuccess?.(fileData)
            emit('upload-success', fileData, index)
            message.success('上传成功')
        } else {
            const error = new Error(response.msg || '上传失败')
            newArchive.error = error.message
            newArchive.uploading = false
            currentArchives.value[index] = { ...newArchive }
            onError?.(error)
            emit('upload-error', error, index)
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        }
    } catch (error: any) {
        console.error('上传压缩包失败:', error)
        newArchive.error = error.message || '上传失败'
        newArchive.uploading = false
        currentArchives.value[index] = { ...newArchive }
        onError?.(error)
        emit('upload-error', error, index)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    }
}

const removeArchive = (index: number) => {
    const archives = [...(currentArchives.value || [])]
    archives.splice(index, 1)
    currentArchives.value = archives
}

const downloadArchive = async (archive: ArchiveItem) => {
    if (!archive.file_id) return
    
    try {
        const blob = await fileApi.downloadFile(archive.file_id)
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = archive.name || 'archive.zip'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
    } catch (error: any) {
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    }
}
</script>

<style scoped lang="less">
.archive-upload {
    width: 100%;
}

.archive-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.archive-item {
    position: relative;
    background: #fafafa;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #d9d9d9;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    }

    &.is-error {
        border-color: #ff4d4f;
    }
}

.archive-wrapper {
    position: relative;
    padding: 16px;
    display: flex;
    align-items: center;
    gap: 16px;
}

.archive-icon {
    font-size: 32px;
    color: var(--color-primary);
    flex-shrink: 0;
}

.archive-info {
    flex: 1;
    min-width: 0;
}

.archive-name {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-bottom: 4px;
}

.archive-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.45);
}

.archive-type {
    padding: 2px 6px;
    background: rgba(var(--color-primary-rgb, 24, 144, 255), 0.1);
    color: var(--color-primary);
    border-radius: 2px;
    font-weight: 500;
}

.archive-actions {
    display: flex;
    gap: 8px;
    flex-shrink: 0;
}

.upload-progress {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 4px 16px;
    background: rgba(0, 0, 0, 0.05);
}

.error-message {
    padding: 0 16px 16px;
    font-size: 12px;
    color: #ff4d4f;
    text-align: center;
}

.upload-trigger {
    min-height: 120px;
    position: relative;
    cursor: pointer;
}

.upload-placeholder {
    width: 100%;
    height: 100%;
    min-height: 120px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px dashed #d9d9d9;
    border-radius: 8px;
    background: #fafafa;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        background: #f0f7ff;
    }
}

.upload-area {
    width: 100%;
}

.upload-dragger {
    padding: 40px;
}

.upload-button {
    width: 100%;
    min-height: 120px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px dashed #d9d9d9;
    border-radius: 8px;
    background: #fafafa;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        background: #f0f7ff;
    }
}

.upload-text {
    margin-top: 8px;
    color: rgba(0, 0, 0, 0.65);
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


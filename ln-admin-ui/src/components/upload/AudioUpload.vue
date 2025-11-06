<template>
    <div class="audio-upload">
        <div v-if="modelValue && modelValue.length > 0" class="audio-list">
            <div
                v-for="(audio, index) in modelValue"
                :key="index"
                class="audio-item"
                :class="{ 'is-error': audio.error }"
            >
                <div class="audio-wrapper">
                    <div class="audio-icon">
                        <SoundOutlined />
                    </div>
                    <div class="audio-player">
                        <audio
                            v-if="audio.url"
                            :src="audio.url"
                            controls
                            preload="metadata"
                            class="audio-element"
                        >
                            您的浏览器不支持音频播放
                        </audio>
                    </div>
                    <div class="audio-overlay">
                        <a-space>
                            <a-button
                                type="primary"
                                danger
                                shape="circle"
                                size="small"
                                @click="removeAudio(index)"
                            >
                                <template #icon><DeleteOutlined /></template>
                            </a-button>
                        </a-space>
                    </div>
                    <div v-if="audio.uploading" class="upload-progress">
                        <a-progress
                            :percent="audio.progress || 0"
                            :status="audio.error ? 'exception' : 'active'"
                            :show-info="false"
                        />
                    </div>
                </div>
                <div class="audio-info">
                    <div class="audio-name" :title="audio.name">{{ audio.name || '音频文件' }}</div>
                    <div class="audio-meta">
                        <span v-if="audio.size">{{ formatFileSize(audio.size) }}</span>
                        <span v-if="audio.duration" class="duration">{{ formatDuration(audio.duration) }}</span>
                    </div>
                </div>
                <div v-if="audio.error" class="error-message">
                    {{ audio.error }}
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
                    <div class="upload-text">上传音频</div>
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
                    <SoundOutlined />
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
    SoundOutlined,
    InboxOutlined,
} from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

/**
 * 音频项接口
 */
export interface AudioItem {
    url?: string
    file_id?: string
    name?: string
    size?: number
    duration?: number
    uploading?: boolean
    progress?: number
    error?: string
    file?: File
}

interface Props {
    modelValue?: AudioItem[]
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
    maxSize: 100,
    maxCount: 10,
    accept: 'audio/*',
    placeholder: '上传音频',
    autoUpload: true,
})

const emit = defineEmits<{
    'update:modelValue': [value: AudioItem[]]
    'change': [value: AudioItem[]]
    'upload-success': [file: FileResponse, index: number]
    'upload-error': [error: Error, index: number]
}>()

const uploadRef = ref()

const acceptText = computed(() => {
    if (props.accept === 'audio/*') return 'MP3、WAV、AAC、OGG 等音频格式'
    return props.accept
})

const canAddMore = computed(() => {
    if (!props.multiple) return false
    return (props.modelValue?.length || 0) < props.maxCount
})

const currentAudios = computed({
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

const formatDuration = (seconds: number): string => {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = Math.floor(seconds % 60)
    
    if (hours > 0) {
        return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
    }
    return `${minutes}:${secs.toString().padStart(2, '0')}`
}

const getAudioDuration = (file: File): Promise<number> => {
    return new Promise((resolve) => {
        const audio = document.createElement('audio')
        audio.preload = 'metadata'
        audio.onloadedmetadata = () => {
            window.URL.revokeObjectURL(audio.src)
            resolve(audio.duration)
        }
        audio.onerror = () => {
            resolve(0)
        }
        audio.src = URL.createObjectURL(file)
    })
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    // 检查文件类型
    if (!file.type.startsWith('audio/')) {
        message.warning('只能上传音频文件！')
        return false
    }

    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        message.warning(`音频大小不能超过 ${props.maxSize}MB！`)
        return false
    }

    // 检查数量限制
    if (props.multiple && (currentAudios.value.length || 0) >= props.maxCount) {
        message.warning(`最多只能上传 ${props.maxCount} 个音频！`)
        return false
    }

    // 如果是自动上传，返回 true
    if (props.autoUpload) {
        return true
    }

    // 否则添加到列表但不自动上传
    const newAudio: AudioItem = {
        file,
        name: file.name,
        size: file.size,
        uploading: false,
        url: URL.createObjectURL(file),
    }

    // 获取音频时长
    getAudioDuration(file).then((duration) => {
        newAudio.duration = duration
        const index = currentAudios.value.findIndex((a) => a.file === file)
        if (index > -1) {
            currentAudios.value[index] = { ...newAudio }
        }
    })

    currentAudios.value = [...(currentAudios.value || []), newAudio]
    return false
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options

    const newAudio: AudioItem = {
        file: file as File,
        name: (file as File).name,
        size: (file as File).size,
        uploading: true,
        progress: 0,
        url: URL.createObjectURL(file as File),
    }

    // 获取音频时长
    const duration = await getAudioDuration(file as File)
    newAudio.duration = duration

    const index = currentAudios.value.length
    currentAudios.value = [...(currentAudios.value || []), newAudio]

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            const fileData = response.data as FileResponse
            newAudio.uploading = false
            newAudio.progress = 100
            newAudio.url = fileData.file_url || fileData.file_path
            newAudio.file_id = fileData.file_id
            newAudio.file = undefined

            currentAudios.value[index] = { ...newAudio }
            onSuccess?.(fileData)
            emit('upload-success', fileData, index)
            message.success('上传成功')
        } else {
            const error = new Error(response.msg || '上传失败')
            newAudio.error = error.message
            newAudio.uploading = false
            currentAudios.value[index] = { ...newAudio }
            onError?.(error)
            emit('upload-error', error, index)
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        }
    } catch (error: any) {
        console.error('上传音频失败:', error)
        newAudio.error = error.message || '上传失败'
        newAudio.uploading = false
        currentAudios.value[index] = { ...newAudio }
        onError?.(error)
        emit('upload-error', error, index)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    }
}

const removeAudio = (index: number) => {
    const audios = [...(currentAudios.value || [])]
    const audio = audios[index]
    
    // 释放对象 URL
    if (audio.file && audio.url && audio.url.startsWith('blob:')) {
        URL.revokeObjectURL(audio.url)
    }

    audios.splice(index, 1)
    currentAudios.value = audios
}
</script>

<style scoped lang="less">
.audio-upload {
    width: 100%;
}

.audio-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.audio-item {
    position: relative;
    background: #fafafa;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #d9d9d9;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

        .audio-overlay {
            opacity: 1;
        }
    }

    &.is-error {
        border-color: #ff4d4f;
    }
}

.audio-wrapper {
    position: relative;
    padding: 16px;
    display: flex;
    align-items: center;
    gap: 16px;
}

.audio-icon {
    font-size: 32px;
    color: var(--color-primary);
    flex-shrink: 0;
}

.audio-player {
    flex: 1;
    min-width: 0;
}

.audio-element {
    width: 100%;
    height: 32px;
}

.audio-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s;
    pointer-events: none;
}

.audio-overlay :deep(.ant-btn) {
    pointer-events: auto;
}

.upload-progress {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 4px;
    background: rgba(0, 0, 0, 0.5);
}

.audio-info {
    padding: 0 16px 16px;
}

.audio-name {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-bottom: 4px;
}

.audio-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.45);
}

.duration {
    padding: 2px 6px;
    background: rgba(0, 0, 0, 0.05);
    border-radius: 2px;
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


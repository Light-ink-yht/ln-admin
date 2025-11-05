<template>
    <div class="video-upload">
        <div v-if="modelValue && modelValue.length > 0" class="video-list">
            <div
                v-for="(video, index) in modelValue"
                :key="index"
                class="video-item"
                :class="{ 'is-error': video.error }"
            >
                <div class="video-wrapper">
                    <video
                        v-if="video.url"
                        :src="video.url"
                        controls
                        preload="metadata"
                        class="video-preview"
                    >
                        您的浏览器不支持视频播放
                    </video>
                    <div v-else class="video-placeholder">
                        <PlayCircleOutlined />
                        <div class="placeholder-text">{{ video.name || '视频文件' }}</div>
                    </div>
                    <div class="video-overlay">
                        <a-space>
                            <a-button
                                type="primary"
                                shape="circle"
                                size="small"
                                @click="previewVideo(index)"
                            >
                                <template #icon><EyeOutlined /></template>
                            </a-button>
                            <a-button
                                type="primary"
                                danger
                                shape="circle"
                                size="small"
                                @click="removeVideo(index)"
                            >
                                <template #icon><DeleteOutlined /></template>
                            </a-button>
                        </a-space>
                    </div>
                    <div v-if="video.uploading" class="upload-progress">
                        <a-progress
                            :percent="video.progress || 0"
                            :status="video.error ? 'exception' : 'active'"
                            :show-info="false"
                        />
                    </div>
                </div>
                <div class="video-info">
                    <div class="video-name" :title="video.name">{{ video.name || '视频文件' }}</div>
                    <div class="video-meta">
                        <span v-if="video.size">{{ formatFileSize(video.size) }}</span>
                        <span v-if="video.duration" class="duration">{{ formatDuration(video.duration) }}</span>
                    </div>
                </div>
                <div v-if="video.error" class="error-message">
                    {{ video.error }}
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
                    <div class="upload-text">上传视频</div>
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
                    <PlayCircleOutlined />
                    <div class="upload-text">{{ placeholder }}</div>
                </div>
            </a-upload>
        </div>

        <!-- 视频预览 -->
        <a-modal
            v-model:open="previewVisible"
            :footer="null"
            centered
            width="900px"
            @cancel="previewVisible = false"
        >
            <video
                v-if="previewVideoUrl"
                :src="previewVideoUrl"
                controls
                autoplay
                style="width: 100%;"
            >
                您的浏览器不支持视频播放
            </video>
        </a-modal>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import {
    PlusOutlined,
    DeleteOutlined,
    EyeOutlined,
    PlayCircleOutlined,
    InboxOutlined,
} from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

/**
 * 视频项接口
 */
export interface VideoItem {
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
    modelValue?: VideoItem[]
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
    maxCount: 5,
    accept: 'video/*',
    placeholder: '上传视频',
    autoUpload: true,
})

const emit = defineEmits<{
    'update:modelValue': [value: VideoItem[]]
    'change': [value: VideoItem[]]
    'upload-success': [file: FileResponse, index: number]
    'upload-error': [error: Error, index: number]
}>()

const uploadRef = ref()
const previewVisible = ref(false)
const previewVideoUrl = ref('')

const acceptText = computed(() => {
    if (props.accept === 'video/*') return 'MP4、AVI、MOV、WMV 等视频格式'
    return props.accept
})

const canAddMore = computed(() => {
    if (!props.multiple) return false
    return (props.modelValue?.length || 0) < props.maxCount
})

const currentVideos = computed({
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

const getVideoDuration = (file: File): Promise<number> => {
    return new Promise((resolve) => {
        const video = document.createElement('video')
        video.preload = 'metadata'
        video.onloadedmetadata = () => {
            window.URL.revokeObjectURL(video.src)
            resolve(video.duration)
        }
        video.onerror = () => {
            resolve(0)
        }
        video.src = URL.createObjectURL(file)
    })
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    // 检查文件类型
    if (!file.type.startsWith('video/')) {
        message.error('只能上传视频文件！')
        return false
    }

    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        message.error(`视频大小不能超过 ${props.maxSize}MB！`)
        return false
    }

    // 检查数量限制
    if (props.multiple && (currentVideos.value.length || 0) >= props.maxCount) {
        message.error(`最多只能上传 ${props.maxCount} 个视频！`)
        return false
    }

    // 如果是自动上传，返回 true
    if (props.autoUpload) {
        return true
    }

    // 否则添加到列表但不自动上传
    const newVideo: VideoItem = {
        file,
        name: file.name,
        size: file.size,
        uploading: false,
        url: URL.createObjectURL(file),
    }

    // 获取视频时长
    getVideoDuration(file).then((duration) => {
        newVideo.duration = duration
        const index = currentVideos.value.findIndex((v) => v.file === file)
        if (index > -1) {
            currentVideos.value[index] = { ...newVideo }
        }
    })

    currentVideos.value = [...(currentVideos.value || []), newVideo]
    return false
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options

    const newVideo: VideoItem = {
        file: file as File,
        name: (file as File).name,
        size: (file as File).size,
        uploading: true,
        progress: 0,
        url: URL.createObjectURL(file as File),
    }

    // 获取视频时长
    const duration = await getVideoDuration(file as File)
    newVideo.duration = duration

    const index = currentVideos.value.length
    currentVideos.value = [...(currentVideos.value || []), newVideo]

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            const fileData = response.data as FileResponse
            newVideo.uploading = false
            newVideo.progress = 100
            newVideo.url = fileData.file_url || fileData.file_path
            newVideo.file_id = fileData.file_id
            newVideo.file = undefined

            currentVideos.value[index] = { ...newVideo }
            onSuccess?.(fileData)
            emit('upload-success', fileData, index)
            message.success('上传成功')
        } else {
            const error = new Error(response.msg || '上传失败')
            newVideo.error = error.message
            newVideo.uploading = false
            currentVideos.value[index] = { ...newVideo }
            onError?.(error)
            emit('upload-error', error, index)
            message.error(response.msg || '上传失败')
        }
    } catch (error: any) {
        console.error('上传视频失败:', error)
        newVideo.error = error.message || '上传失败'
        newVideo.uploading = false
        currentVideos.value[index] = { ...newVideo }
        onError?.(error)
        emit('upload-error', error, index)
        message.error(error.message || '上传失败')
    }
}

const removeVideo = (index: number) => {
    const videos = [...(currentVideos.value || [])]
    const video = videos[index]
    
    // 释放对象 URL
    if (video.file && video.url && video.url.startsWith('blob:')) {
        URL.revokeObjectURL(video.url)
    }

    videos.splice(index, 1)
    currentVideos.value = videos
}

const previewVideo = (index: number) => {
    const video = currentVideos.value[index]
    previewVideoUrl.value = video.url || ''
    previewVisible.value = true
}
</script>

<style scoped lang="less">
.video-upload {
    width: 100%;
}

.video-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
}

.video-item {
    position: relative;
    background: #fafafa;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #d9d9d9;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

        .video-overlay {
            opacity: 1;
        }
    }

    &.is-error {
        border-color: #ff4d4f;
    }
}

.video-wrapper {
    position: relative;
    width: 100%;
    padding-top: 56.25%; /* 16:9 比例 */
    background: #000;
    overflow: hidden;
}

.video-preview {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.video-placeholder {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: rgba(255, 255, 255, 0.65);
}

.video-placeholder .anticon {
    font-size: 48px;
    margin-bottom: 8px;
}

.placeholder-text {
    font-size: 14px;
}

.video-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s;
}

.upload-progress {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 4px;
    background: rgba(0, 0, 0, 0.5);
}

.video-info {
    padding: 12px;
}

.video-name {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-bottom: 4px;
}

.video-meta {
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
    padding: 0 12px 12px;
    font-size: 12px;
    color: #ff4d4f;
    text-align: center;
}

.upload-trigger {
    padding-top: 56.25%; /* 16:9 比例 */
    position: relative;
    cursor: pointer;
}

.upload-placeholder {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
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
    min-height: 200px;
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


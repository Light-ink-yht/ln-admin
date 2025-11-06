<template>
    <div class="image-upload">
        <div v-if="modelValue && modelValue.length > 0" class="image-list">
            <div
                v-for="(image, index) in modelValue"
                :key="index"
                class="image-item"
                :class="{ 'is-error': image.error }"
            >
                <div class="image-wrapper">
                    <img :src="getImageUrl(image)" :alt="image.name || '图片'" />
                    <div class="image-overlay">
                        <a-space>
                            <a-button
                                type="primary"
                                shape="circle"
                                size="small"
                                @click="previewImage(index)"
                            >
                                <template #icon><EyeOutlined /></template>
                            </a-button>
                            <a-button
                                type="primary"
                                danger
                                shape="circle"
                                size="small"
                                @click="removeImage(index)"
                            >
                                <template #icon><DeleteOutlined /></template>
                            </a-button>
                        </a-space>
                    </div>
                    <div v-if="image.uploading" class="upload-progress">
                        <a-progress
                            :percent="image.progress || 0"
                            :status="image.error ? 'exception' : 'active'"
                            :show-info="false"
                        />
                    </div>
                </div>
                <div v-if="image.error" class="error-message">
                    {{ image.error }}
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
                    <div class="upload-text">上传图片</div>
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
                    <p class="ant-upload-text">点击或拖拽图片到此区域上传</p>
                    <p class="ant-upload-hint">
                        支持 {{ acceptText }}，单个文件不超过 {{ maxSize }}MB
                    </p>
                </div>
                <div v-else class="upload-button">
                    <PlusOutlined />
                    <div class="upload-text">{{ placeholder }}</div>
                </div>
            </a-upload>
        </div>

        <!-- 图片预览 -->
        <a-modal
            v-model:open="previewVisible"
            :footer="null"
            centered
            width="800px"
            @cancel="previewVisible = false"
        >
            <img :src="previewImageUrl" style="width: 100%;" alt="预览" />
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
    InboxOutlined,
} from '@ant-design/icons-vue'
import type { UploadFile, UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

/**
 * 图片项接口
 */
export interface ImageItem {
    url?: string
    file_id?: string
    name?: string
    uploading?: boolean
    progress?: number
    error?: string
    file?: File
}

interface Props {
    modelValue?: ImageItem[]
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
    maxSize: 10,
    maxCount: 9,
    accept: 'image/*',
    placeholder: '上传图片',
    autoUpload: true,
})

const emit = defineEmits<{
    'update:modelValue': [value: ImageItem[]]
    'change': [value: ImageItem[]]
    'upload-success': [file: FileResponse, index: number]
    'upload-error': [error: Error, index: number]
}>()

const uploadRef = ref()
const previewVisible = ref(false)
const previewImageUrl = ref('')

const acceptText = computed(() => {
    if (props.accept === 'image/*') return 'JPG、PNG、GIF 等图片格式'
    return props.accept
})

const canAddMore = computed(() => {
    if (!props.multiple) return false
    return (props.modelValue?.length || 0) < props.maxCount
})

const currentImages = computed({
    get: () => props.modelValue || [],
    set: (val) => {
        emit('update:modelValue', val)
        emit('change', val)
    },
})

const triggerUpload = () => {
    uploadRef.value?.$el.querySelector('input')?.click()
}

const getImageUrl = (image: ImageItem): string => {
    if (image.url) return image.url
    if (image.file) {
        return URL.createObjectURL(image.file)
    }
    return ''
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    // 检查文件类型
    if (!file.type.startsWith('image/')) {
        message.warning('只能上传图片文件！')
        return false
    }

    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        message.warning(`图片大小不能超过 ${props.maxSize}MB！`)
        return false
    }

    // 检查数量限制
    if (props.multiple && (currentImages.value.length || 0) >= props.maxCount) {
        message.warning(`最多只能上传 ${props.maxCount} 张图片！`)
        return false
    }

    // 如果是自动上传，返回 true
    if (props.autoUpload) {
        return true
    }

    // 否则添加到列表但不自动上传
    const newImage: ImageItem = {
        file,
        name: file.name,
        uploading: false,
        url: URL.createObjectURL(file),
    }
    currentImages.value = [...(currentImages.value || []), newImage]
    return false
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError, onProgress } = options

    const newImage: ImageItem = {
        file: file as File,
        name: (file as File).name,
        uploading: true,
        progress: 0,
        url: URL.createObjectURL(file as File),
    }

    const index = currentImages.value.length
    currentImages.value = [...(currentImages.value || []), newImage]

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            const fileData = response.data as FileResponse
            newImage.uploading = false
            newImage.progress = 100
            newImage.url = fileData.file_url || fileData.file_path
            newImage.file_id = fileData.file_id
            newImage.file = undefined

            currentImages.value[index] = { ...newImage }
            onSuccess?.(fileData)
            emit('upload-success', fileData, index)
            message.success('上传成功')
        } else {
            const error = new Error(response.msg || '上传失败')
            newImage.error = error.message
            newImage.uploading = false
            currentImages.value[index] = { ...newImage }
            onError?.(error)
            emit('upload-error', error, index)
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        }
    } catch (error: any) {
        console.error('上传图片失败:', error)
        newImage.error = error.message || '上传失败'
        newImage.uploading = false
        currentImages.value[index] = { ...newImage }
        onError?.(error)
        emit('upload-error', error, index)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    }
}

const removeImage = (index: number) => {
    const images = [...(currentImages.value || [])]
    const image = images[index]
    
    // 释放对象 URL
    if (image.file && image.url && image.url.startsWith('blob:')) {
        URL.revokeObjectURL(image.url)
    }

    images.splice(index, 1)
    currentImages.value = images
}

const previewImage = (index: number) => {
    const image = currentImages.value[index]
    previewImageUrl.value = getImageUrl(image)
    previewVisible.value = true
}
</script>

<style scoped lang="less">
.image-upload {
    width: 100%;
}

.image-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 16px;
}

.image-item {
    position: relative;
    width: 100%;
    padding-top: 100%; /* 1:1 比例 */
    background: #fafafa;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #d9d9d9;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

        .image-overlay {
            opacity: 1;
        }
    }

    &.is-error {
        border-color: #ff4d4f;
    }
}

.image-wrapper {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.image-wrapper img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.image-overlay {
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

.error-message {
    position: absolute;
    bottom: -20px;
    left: 0;
    right: 0;
    font-size: 12px;
    color: #ff4d4f;
    text-align: center;
}

.upload-trigger {
    padding-top: 100%; /* 1:1 比例 */
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


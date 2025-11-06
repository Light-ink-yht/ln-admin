<template>
    <div class="image-upload-single">
        <div v-if="modelValue" class="image-preview">
            <div class="preview-wrapper">
                <img :src="modelValue" :alt="previewAlt" />
                <div class="preview-overlay">
                    <a-space>
                        <a-button
                            type="primary"
                            shape="circle"
                            size="small"
                            @click="previewImage"
                        >
                            <template #icon><EyeOutlined /></template>
                        </a-button>
                        <a-button
                            type="primary"
                            shape="circle"
                            size="small"
                            @click="triggerUpload"
                        >
                            <template #icon><EditOutlined /></template>
                        </a-button>
                        <a-button
                            type="primary"
                            danger
                            shape="circle"
                            size="small"
                            @click="removeImage"
                        >
                            <template #icon><DeleteOutlined /></template>
                        </a-button>
                    </a-space>
                </div>
                <div v-if="uploading" class="upload-progress">
                    <a-progress
                        :percent="uploadProgress"
                        :status="uploadError ? 'exception' : 'active'"
                        :show-info="false"
                    />
                </div>
            </div>
            <div v-if="uploadError" class="error-message">
                {{ uploadError }}
            </div>
        </div>

        <div v-else class="upload-area" @click="triggerUpload">
            <a-upload
                ref="uploadRef"
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
                        支持 {{ acceptText }}，文件不超过 {{ maxSize }}MB
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
            <img :src="modelValue" style="width: 100%;" alt="预览" />
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
    EditOutlined,
    InboxOutlined,
} from '@ant-design/icons-vue'
import type { UploadProps } from 'ant-design-vue'
import { fileApi, type FileResponse } from '@/api/file'

interface Props {
    modelValue?: string
    maxSize?: number // MB
    accept?: string
    placeholder?: string
    previewAlt?: string
    drag?: boolean
    autoUpload?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: '',
    maxSize: 10,
    accept: 'image/*',
    placeholder: '上传图片',
    previewAlt: '图片预览',
    drag: false,
    autoUpload: true,
})

const emit = defineEmits<{
    'update:modelValue': [value: string]
    'change': [value: string]
    'upload-success': [file: FileResponse]
    'upload-error': [error: Error]
}>()

const uploadRef = ref()
const previewVisible = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref('')

const acceptText = computed(() => {
    if (props.accept === 'image/*') return 'JPG、PNG、GIF 等图片格式'
    return props.accept
})

const triggerUpload = () => {
    uploadRef.value?.$el.querySelector('input')?.click()
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
    uploadError.value = ''

    // 检查文件类型
    if (!file.type.startsWith('image/')) {
        const error = '只能上传图片文件！'
        message.warning(error)
        uploadError.value = error
        return false
    }

    // 检查文件大小
    const maxSizeBytes = props.maxSize * 1024 * 1024
    if (file.size > maxSizeBytes) {
        const error = `图片大小不能超过 ${props.maxSize}MB！`
        message.warning(error)
        uploadError.value = error
        return false
    }

    // 如果是自动上传，返回 true
    if (props.autoUpload) {
        return true
    }

    // 否则预览但不自动上传
    const url = URL.createObjectURL(file)
    emit('update:modelValue', url)
    emit('change', url)
    return false
}

const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options

    uploading.value = true
    uploadProgress.value = 0
    uploadError.value = ''

    // 先显示预览
    const previewUrl = URL.createObjectURL(file as File)
    emit('update:modelValue', previewUrl)
    emit('change', previewUrl)

    try {
        const response = await fileApi.uploadFile(file as File)
        if (response.code === 200 || response.code === 0) {
            const fileData = response.data as FileResponse
            const finalUrl = fileData.file_url || fileData.file_path

            // 释放预览 URL
            URL.revokeObjectURL(previewUrl)

            emit('update:modelValue', finalUrl)
            emit('change', finalUrl)
            emit('upload-success', fileData)
            onSuccess?.(fileData)
            message.success('上传成功')
        } else {
            const error = new Error(response.msg || '上传失败')
            uploadError.value = error.message
            emit('upload-error', error)
            onError?.(error)
            // 错误提示已在 request.ts 中统一处理，这里不再重复显示
        }
    } catch (error: any) {
        console.error('上传图片失败:', error)
        uploadError.value = error.message || '上传失败'
        emit('upload-error', error)
        onError?.(error)
        // 错误提示已在 request.ts 中统一处理，这里不再重复显示
    } finally {
        uploading.value = false
        uploadProgress.value = 0
    }
}

const removeImage = () => {
    if (props.modelValue && props.modelValue.startsWith('blob:')) {
        URL.revokeObjectURL(props.modelValue)
    }
    emit('update:modelValue', '')
    emit('change', '')
    uploadError.value = ''
}

const previewImage = () => {
    previewVisible.value = true
}
</script>

<style scoped lang="less">
.image-upload-single {
    width: 100%;
}

.image-preview {
    width: 100%;
}

.preview-wrapper {
    position: relative;
    width: 100%;
    padding-top: 100%; /* 1:1 比例，可根据需要调整 */
    background: #fafafa;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid #d9d9d9;
    transition: all 0.3s;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

        .preview-overlay {
            opacity: 1;
        }
    }
}

.preview-wrapper img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.preview-overlay {
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
    margin-top: 8px;
    font-size: 12px;
    color: #ff4d4f;
    text-align: center;
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


<template>
    <div class="file-select">
        <div v-if="fileUrl && isImage" class="avatar-preview-wrapper">
            <div class="avatar-section">
                <a-avatar :size="120" :src="fileUrl">
                    <template v-if="!fileUrl">
                        <UserOutlined />
                    </template>
                </a-avatar>
                <div class="avatar-actions">
                    <a-button type="link" @click="openFileSelector">
                        <template #icon><UploadOutlined /></template>
                        更换图片
                    </a-button>
                    <a-button type="link" danger @click="clearFile">
                        <template #icon><DeleteOutlined /></template>
                        清除
                    </a-button>
                </div>
            </div>
            <a-input
                v-model:value="fileUrl"
                :placeholder="placeholder"
                readonly
                style="margin-top: 12px;"
                @click="openFileSelector"
            >
                <template #addonAfter>
                    <a-button type="link" size="small" @click="openFileSelector">选择文件</a-button>
                </template>
            </a-input>
        </div>
        <div v-else>
            <a-input
                v-model:value="fileUrl"
                :placeholder="placeholder"
                readonly
                @click="openFileSelector"
            >
                <template #addonAfter>
                    <a-button type="link" size="small" @click="openFileSelector">选择文件</a-button>
                </template>
            </a-input>
            <div v-if="fileUrl && !isImage" class="file-preview">
                <div class="preview-placeholder">
                    <FileOutlined />
                    <span>已选择文件</span>
                </div>
                <a-button type="link" danger size="small" @click="clearFile">清除</a-button>
            </div>
        </div>

        <!-- 文件选择弹窗 -->
        <a-modal
            v-model:open="fileSelectorVisible"
            title="选择文件"
            width="1200px"
            :footer="null"
            @cancel="handleCancel"
        >
            <FileListSelector
                :multiple="false"
                :accept-types="acceptTypes"
                @select="handleFileSelect"
            />
        </a-modal>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { FileOutlined, UserOutlined, UploadOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import FileListSelector from '@/views/pages/file/components/FileListSelector.vue'

interface Props {
    modelValue?: string
    placeholder?: string
    acceptTypes?: string[]
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: '',
    placeholder: '请选择文件',
    acceptTypes: () => [],
})

const emit = defineEmits<{
    'update:modelValue': [value: string]
}>()

const fileUrl = ref(props.modelValue)
const fileSelectorVisible = ref(false)

const isImage = computed(() => {
    if (!fileUrl.value) return false
    const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg']
    return imageExts.some(ext => fileUrl.value.toLowerCase().endsWith(ext))
})

watch(() => props.modelValue, (val) => {
    fileUrl.value = val
})

watch(fileUrl, (val) => {
    emit('update:modelValue', val)
})

const openFileSelector = () => {
    fileSelectorVisible.value = true
}

const handleFileSelect = (file: any) => {
    // 优先使用file_url，如果没有则使用file_path，最后使用url
    fileUrl.value = file.file_url || file.file_path || file.url || ''
    fileSelectorVisible.value = false
}

const handleCancel = () => {
    fileSelectorVisible.value = false
}

const clearFile = () => {
    fileUrl.value = ''
    emit('update:modelValue', '')
}
</script>

<style scoped>
.file-select {
    width: 100%;
}

.avatar-preview-wrapper {
    width: 100%;
}

.avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    margin-bottom: 12px;
}

.avatar-section :deep(.ant-avatar) {
    border: 4px solid rgba(255, 255, 255, 0.3);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.avatar-actions {
    display: flex;
    gap: 16px;
}

.avatar-actions :deep(.ant-btn) {
    color: #fff;
    padding: 0;
    height: auto;
}

.avatar-actions :deep(.ant-btn:hover) {
    color: rgba(255, 255, 255, 0.8);
}

.file-preview {
    margin-top: 8px;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px;
    background: #fafafa;
    border-radius: 4px;
}

.preview-placeholder {
    display: flex;
    align-items: center;
    gap: 8px;
    color: rgba(0, 0, 0, 0.45);
}
</style>


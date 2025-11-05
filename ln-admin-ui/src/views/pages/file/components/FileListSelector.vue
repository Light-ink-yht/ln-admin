<template>
    <div class="file-list-selector">
        <a-spin :spinning="loading">
            <div class="file-list-header">
                <a-space>
                    <a-input
                        v-model:value="searchParams.file_name"
                        placeholder="搜索文件名"
                        style="width: 200px"
                        allow-clear
                        @press-enter="handleSearch"
                    >
                        <template #prefix>
                            <SearchOutlined />
                        </template>
                    </a-input>
                    <a-select
                        v-model:value="searchParams.category"
                        placeholder="文件类型"
                        style="width: 150px"
                        allow-clear
                        @change="handleSearch"
                    >
                        <a-select-option value="image">图片</a-select-option>
                        <a-select-option value="document">文档</a-select-option>
                        <a-select-option value="video">视频</a-select-option>
                        <a-select-option value="audio">音频</a-select-option>
                        <a-select-option value="archive">压缩包</a-select-option>
                        <a-select-option value="other">其他</a-select-option>
                    </a-select>
                    <a-button type="primary" @click="handleSearch">
                        <template #icon><SearchOutlined /></template>
                        搜索
                    </a-button>
                    <a-button @click="handleReset">重置</a-button>
                    <a-button type="primary" @click="uploadVisible = true">
                        <template #icon><UploadOutlined /></template>
                        上传文件
                    </a-button>
                </a-space>
            </div>

            <div class="file-grid">
                <div
                    v-for="file in fileList"
                    :key="file.file_id"
                    class="file-item"
                    :class="{ selected: selectedFileId === file.file_id }"
                    @click="handleSelectFile(file)"
                >
                    <div class="file-thumbnail">
                        <img
                            v-if="file.category === 'image' && (file.file_url || file.file_path)"
                            :src="file.file_url || file.file_path"
                            alt=""
                            @error="handleImageError"
                        />
                        <FileOutlined v-else class="file-icon" />
                    </div>
                    <div class="file-info">
                        <div class="file-name" :title="file.original_name">{{ file.original_name }}</div>
                        <div class="file-meta">
                            <span>{{ formatFileSize(file.file_size) }}</span>
                            <span>{{ file.extension }}</span>
                        </div>
                    </div>
                    <CheckCircleOutlined v-if="selectedFileId === file.file_id" class="selected-icon" />
                </div>
            </div>

            <div class="file-pagination">
                <a-pagination
                    v-model:current="pagination.page"
                    v-model:page-size="pagination.pageSize"
                    :total="pagination.total"
                    :show-size-changer="true"
                    :show-total="(total) => `共 ${total} 条`"
                    @change="handlePageChange"
                />
            </div>
        </a-spin>

        <!-- 文件上传弹窗 -->
        <a-modal
            v-model:open="uploadVisible"
            title="上传文件"
            :footer="null"
            width="600px"
        >
            <FileUpload
                ref="fileUploadRef"
                :multiple="false"
                :drag="true"
                @success="handleUploadSuccess"
            />
        </a-modal>
    </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { SearchOutlined, FileOutlined, CheckCircleOutlined, UploadOutlined } from '@ant-design/icons-vue'
import { fileApi, type FileResponse } from '@/api/file'
import { formatFileSize } from '@/views/pages/file/utils'
import FileUpload from './FileUpload.vue'

interface Props {
    multiple?: boolean
    acceptTypes?: string[]
}

const props = withDefaults(defineProps<Props>(), {
    multiple: false,
    acceptTypes: () => [],
})

const emit = defineEmits<{
    select: [file: FileResponse]
}>()

const loading = ref(false)
const fileList = ref<FileResponse[]>([])
const selectedFileId = ref<string>('')
const uploadVisible = ref(false)
const fileUploadRef = ref()

const searchParams = reactive({
    file_name: '',
    category: '',
})

const pagination = reactive({
    page: 1,
    pageSize: 20,
    total: 0,
})

const fetchFileList = async () => {
    try {
        loading.value = true
        const response = await fileApi.getFileList({
            page: pagination.page,
            page_size: pagination.pageSize,
            file_name: searchParams.file_name || undefined,
            category: searchParams.category || undefined,
        })

        if (response.code === 200 || response.code === 0) {
            fileList.value = response.data?.list || []
            pagination.total = response.data?.total || 0
        } else {
            console.error('获取文件列表失败:', response.msg)
        }
    } catch (error) {
        console.error('获取文件列表失败:', error)
    } finally {
        loading.value = false
    }
}

const handleSelectFile = (file: FileResponse) => {
    selectedFileId.value = file.file_id
    emit('select', file)
}

const handleSearch = () => {
    pagination.page = 1
    fetchFileList()
}

const handleReset = () => {
    searchParams.file_name = ''
    searchParams.category = ''
    pagination.page = 1
    fetchFileList()
}

const handlePageChange = () => {
    fetchFileList()
}

const handleImageError = (e: Event) => {
    const img = e.target as HTMLImageElement
    img.style.display = 'none'
}

// 上传成功回调
const handleUploadSuccess = (file: FileResponse) => {
    uploadVisible.value = false
    // 刷新文件列表
    fetchFileList()
    // 自动选择刚上传的文件
    if (file) {
        handleSelectFile(file)
    }
}

onMounted(() => {
    fetchFileList()
})
</script>

<style scoped>
.file-list-selector {
    padding: 16px;
}

.file-list-header {
    margin-bottom: 16px;
}

.file-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
    margin-bottom: 16px;
    max-height: 500px;
    overflow-y: auto;
    padding: 8px;
}

.file-item {
    position: relative;
    border: 2px solid #f0f0f0;
    border-radius: 8px;
    padding: 12px;
    cursor: pointer;
    transition: all 0.3s;
    background: #fff;
}

.file-item:hover {
    border-color: #1890ff;
    box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

.file-item.selected {
    border-color: #1890ff;
    background: #e6f7ff;
}

.file-thumbnail {
    width: 100%;
    height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fafafa;
    border-radius: 4px;
    margin-bottom: 8px;
    overflow: hidden;
}

.file-thumbnail img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
}

.file-icon {
    font-size: 48px;
    color: rgba(0, 0, 0, 0.25);
}

.file-info {
    text-align: center;
}

.file-name {
    font-size: 12px;
    color: rgba(0, 0, 0, 0.85);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-bottom: 4px;
}

.file-meta {
    font-size: 11px;
    color: rgba(0, 0, 0, 0.45);
    display: flex;
    justify-content: space-between;
    gap: 4px;
}

.selected-icon {
    position: absolute;
    top: 8px;
    right: 8px;
    color: #1890ff;
    font-size: 20px;
}

.file-pagination {
    margin-top: 16px;
    text-align: right;
}
</style>


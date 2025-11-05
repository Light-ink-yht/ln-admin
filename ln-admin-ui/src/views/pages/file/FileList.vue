<template>
    <PermissionDeniedAlert
        v-model:visible="permissionDenied"
        description="您没有访问文件列表的权限，请联系管理员为您分配相应的权限。"
    />
    <DataList
        ref="dataListRef"
        title="文件列表"
        :search-fields="searchFields"
        :columns="columns"
        :fetch-data="fetchFileList"
        :actions="actions"
        :scroll="{ x: 1400 }"
        row-key="file_id"
    >
        <!-- 自定义列插槽 -->
        <template #column-preview="{ record }">
            <FilePreview :file="record" />
        </template>

        <template #column-category="{ record }">
            <FileCategory :category="record.category" />
        </template>

        <template #column-storage_type="{ record }">
            <a-tag :color="record.storage_type === 's3' ? 'blue' : 'green'">
                {{ record.storage_type === 's3' ? 'S3存储' : '本地存储' }}
            </a-tag>
        </template>

        <template #column-file_size="{ record }">
            {{ formatFileSize(record.file_size) }}
        </template>

        <template #column-upload_time="{ record }">
            {{ record.upload_time || '-' }}
        </template>

        <template #column-created_at="{ record }">
            {{ record.created_at }}
        </template>

        <template #column-action="{ record }">
            <a-space>
                <a-button type="link" size="small" @click="handlePreview(record)">预览</a-button>
                <a-button type="link" size="small" @click="handleDownload(record)">下载</a-button>
                <a-button type="link" size="small" danger @click="handleDelete(record)">删除</a-button>
            </a-space>
        </template>
    </DataList>

    <!-- 文件上传弹窗 -->
    <a-modal
        v-model:open="uploadVisible"
        title="上传文件"
        :footer="null"
        width="600px"
    >
        <FileUpload
            ref="fileUploadRef"
            :multiple="true"
            :drag="true"
            @success="handleUploadSuccess"
        />
    </a-modal>

    <!-- 文件预览弹窗 -->
    <FilePreviewModal
        v-model:open="previewVisible"
        :file="previewFile"
    />
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
    DataList,
    PermissionDeniedAlert,
    type ExtraAction,
} from '@/views/pages/components'
import { fileApi, type FileResponse, type FileListParams } from '@/api/file'
import FilePreview from './components/FilePreview.vue'
import FileCategory from './components/FileCategory.vue'
import FilePreviewModal from './components/FilePreviewModal.vue'
import FileUpload from './components/FileUpload.vue'
import { getFileColumns, getFileSearchFields, formatFileSize } from './utils'

const permissionDenied = ref(false)
const dataListRef = ref()
const previewVisible = ref(false)
const previewFile = ref<FileResponse | null>(null)
const uploadVisible = ref(false)
const fileUploadRef = ref()

// 搜索字段
const searchFields = getFileSearchFields()

// 列定义
const columns = getFileColumns()

// 获取文件列表
const fetchFileList = async (params: FileListParams) => {
    try {
        const response = await fileApi.getFileList(params)
        if (response.code === 200 || response.code === 0) {
            return {
                list: response.data?.list || [],
                total: response.data?.total || 0,
            }
        } else {
            if (response.code === 403) {
                permissionDenied.value = true
            }
            message.error(response.msg || '获取文件列表失败')
            return { list: [], total: 0 }
        }
    } catch (error: any) {
        console.error('获取文件列表失败:', error)
        if (error.response?.status === 403) {
            permissionDenied.value = true
        }
        message.error(error.message || '获取文件列表失败')
        return { list: [], total: 0 }
    }
}

// 操作按钮
const actions: ExtraAction[] = [
    {
        key: 'upload',
        label: '上传文件',
        type: 'primary',
        icon: 'UploadOutlined',
        onClick: () => {
            uploadVisible.value = true
        },
    },
]

// 预览文件
const handlePreview = (record: FileResponse) => {
    previewFile.value = record
    previewVisible.value = true
}

// 下载文件
const handleDownload = async (record: FileResponse) => {
    try {
        const blob = await fileApi.downloadFile(record.file_id)
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = record.original_name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        message.success('下载成功')
    } catch (error: any) {
        console.error('下载文件失败:', error)
        message.error(error.message || '下载文件失败')
    }
}

// 删除文件
const handleDelete = (record: FileResponse) => {
    Modal.confirm({
        title: '确认删除',
        content: `确定要删除文件 "${record.original_name}" 吗？`,
        onOk: async () => {
            try {
                const response = await fileApi.deleteFile(record.file_id)
                if (response.code === 200 || response.code === 0) {
                    message.success('删除成功')
                    dataListRef.value?.refresh()
                } else {
                    message.error(response.msg || '删除失败')
                }
            } catch (error: any) {
                console.error('删除文件失败:', error)
                message.error(error.message || '删除文件失败')
            }
        },
    })
}

// 上传成功回调
const handleUploadSuccess = () => {
    dataListRef.value?.refresh()
    uploadVisible.value = false
}

// 将操作添加到列定义
columns.push({
    title: '操作',
    key: 'action',
    fixed: 'right',
    width: 200,
    align: 'center',
    slots: { customRender: 'action' },
})
</script>



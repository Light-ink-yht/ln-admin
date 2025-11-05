import type { TableColumnsType } from 'ant-design-vue'
import type { SearchField } from '@/views/pages/components'

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

/**
 * 获取文件搜索字段
 */
export function getFileSearchFields(): SearchField[] {
    return [
        {
            key: 'file_name',
            label: '文件名',
            type: 'input',
            placeholder: '请输入文件名',
            width: '200px',
        },
        {
            key: 'storage_type',
            label: '存储类型',
            type: 'select',
            placeholder: '请选择存储类型',
            width: '120px',
            options: [
                { label: '本地存储', value: 'local' },
                { label: 'S3存储', value: 's3' },
            ],
        },
        {
            key: 'category',
            label: '文件分类',
            type: 'select',
            placeholder: '请选择分类',
            width: '120px',
            options: [
                { label: '图片', value: 'image' },
                { label: '文档', value: 'document' },
                { label: '视频', value: 'video' },
                { label: '音频', value: 'audio' },
                { label: '压缩包', value: 'archive' },
                { label: '其他', value: 'other' },
            ],
        },
        {
            key: 'extension',
            label: '扩展名',
            type: 'input',
            placeholder: '请输入扩展名',
            width: '120px',
        },
    ]
}

/**
 * 获取文件列定义
 */
export function getFileColumns(): TableColumnsType {
    return [
        {
            title: '预览',
            key: 'preview',
            dataIndex: 'preview',
            width: 100,
            align: 'center',
            slots: { customRender: 'preview' },
        },
        {
            title: '文件名',
            key: 'original_name',
            dataIndex: 'original_name',
            width: 250,
            ellipsis: true,
        },
        {
            title: '文件分类',
            key: 'category',
            dataIndex: 'category',
            width: 100,
            align: 'center',
            slots: { customRender: 'category' },
        },
        {
            title: '扩展名',
            key: 'extension',
            dataIndex: 'extension',
            width: 80,
            align: 'center',
        },
        {
            title: '文件大小',
            key: 'file_size',
            dataIndex: 'file_size',
            width: 120,
            align: 'right',
            slots: { customRender: 'file_size' },
        },
        {
            title: '存储类型',
            key: 'storage_type',
            dataIndex: 'storage_type',
            width: 100,
            align: 'center',
            slots: { customRender: 'storage_type' },
        },
        {
            title: 'MIME类型',
            key: 'mime_type',
            dataIndex: 'mime_type',
            width: 150,
            ellipsis: true,
        },
        {
            title: '上传时间',
            key: 'upload_time',
            dataIndex: 'upload_time',
            width: 180,
            slots: { customRender: 'upload_time' },
        },
        {
            title: '创建时间',
            key: 'created_at',
            dataIndex: 'created_at',
            width: 180,
            slots: { customRender: 'created_at' },
        },
    ]
}


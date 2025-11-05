import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 文件响应数据
 */
export interface FileResponse {
    file_id: string
    file_name: string
    original_name: string
    file_path: string
    file_url: string
    storage_type: 'local' | 's3'
    file_size: number
    mime_type: string
    extension: string
    category: 'image' | 'document' | 'video' | 'audio' | 'archive' | 'other'
    upload_time?: string
    created_at: string
    updated_at: string
}

/**
 * 文件列表请求参数
 */
export interface FileListParams {
    page?: number
    page_size?: number
    file_name?: string
    storage_type?: 'local' | 's3'
    category?: 'image' | 'document' | 'video' | 'audio' | 'archive' | 'other'
    extension?: string
    creator_id?: string
}

/**
 * 分页响应数据
 */
export interface FilePageResponse extends ResponseData<{
    list: FileResponse[]
    total: number
    page: number
    page_size: number
}> {}

/**
 * 文件 API
 */
export const fileApi = {
    /**
     * 上传文件
     */
    uploadFile(file: File): Promise<ResponseData<FileResponse>> {
        return request.upload<FileResponse>('/file/upload', file)
    },

    /**
     * 获取文件列表
     */
    getFileList(params?: FileListParams): Promise<FilePageResponse> {
        return request.get<{
            list: FileResponse[]
            total: number
            page: number
            page_size: number
        }>('/file/list', params) as Promise<FilePageResponse>
    },

    /**
     * 获取文件详情
     */
    getFileDetail(fileId: string): Promise<ResponseData<FileResponse>> {
        return request.get<FileResponse>(`/file/${fileId}`) as Promise<ResponseData<FileResponse>>
    },

    /**
     * 下载文件
     */
    async downloadFile(fileId: string): Promise<Blob> {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || ''}/api/file/${fileId}/download`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('accessToken') || ''}`,
            },
        })
        if (!response.ok) {
            throw new Error('下载失败')
        }
        return await response.blob()
    },

    /**
     * 删除文件
     */
    deleteFile(fileId: string): Promise<ResponseData<void>> {
        return request.delete<void>(`/file/${fileId}`) as Promise<ResponseData<void>>
    },

    /**
     * 获取存储配置
     */
    getStorageConfig(): Promise<ResponseData<StorageConfigResponse>> {
        return request.get<StorageConfigResponse>('/file/storage/config') as Promise<ResponseData<StorageConfigResponse>>
    },

    /**
     * 保存存储配置
     */
    saveStorageConfig(data: StorageConfigRequest): Promise<ResponseData<StorageConfigResponse>> {
        return request.post<StorageConfigResponse>('/file/storage/config', data) as Promise<ResponseData<StorageConfigResponse>>
    },
}

/**
 * 存储配置请求
 */
export interface StorageConfigRequest {
    storage_type: 'local' | 's3'
    bucket?: string
    region?: string
    endpoint?: string
    access_key_id?: string
    secret_key?: string
    base_url?: string
    status: '1' | '2'
    remarks?: string
}

/**
 * 存储配置响应
 */
export interface StorageConfigResponse {
    config_id?: string
    storage_type: 'local' | 's3'
    bucket?: string
    region?: string
    endpoint?: string
    access_key_id?: string
    secret_key?: string
    base_url?: string
    status: '1' | '2'
    remarks?: string
    last_modify_time?: string
    created_at?: string
    updated_at?: string
}


import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * CPU信息
 */
export interface CPUInfo {
    usage: number
    count: number
    model_name: string
    temperature?: number
}

/**
 * 内存信息
 */
export interface MemoryInfo {
    total: number
    used: number
    available: number
    usage: number
    total_gb: number
    used_gb: number
    available_gb: number
}

/**
 * 磁盘分区信息
 */
export interface DiskPartition {
    device: string
    mountpoint: string
    fstype: string
    total: number
    used: number
    available: number
    usage: number
    total_gb: number
    used_gb: number
    available_gb: number
}

/**
 * 磁盘信息
 */
export interface DiskInfo {
    total: number
    used: number
    available: number
    usage: number
    total_gb: number
    used_gb: number
    available_gb: number
    disks: DiskPartition[]
}

/**
 * 系统监控响应
 */
export interface SystemMonitorResponse {
    cpu: CPUInfo
    memory: MemoryInfo
    disk: DiskInfo
    uptime: string
}

/**
 * 系统运维 API
 */
export const opsApi = {
    /**
     * 获取系统监控信息
     */
    getSystemMonitor(): Promise<ResponseData<SystemMonitorResponse>> {
        return request.get<SystemMonitorResponse>('/ops/monitor')
    },
}


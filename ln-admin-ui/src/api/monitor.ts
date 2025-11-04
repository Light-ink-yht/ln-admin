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
 * 系统监控信息
 */
export interface SystemMonitorInfo {
    cpu: CPUInfo
    memory: MemoryInfo
    disk: DiskInfo
    uptime: string
}

/**
 * 系统监控 API
 */
export const monitorApi = {
    /**
     * 获取系统监控信息
     */
    getSystemMonitor(): Promise<ResponseData<SystemMonitorInfo>> {
        return request.get<SystemMonitorInfo>('/ops/monitor')
    },
}


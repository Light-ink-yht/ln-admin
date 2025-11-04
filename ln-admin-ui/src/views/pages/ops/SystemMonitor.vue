<template>
    <div class="system-monitor-container">
        <PermissionDeniedAlert
            v-model:visible="permissionDenied"
            description="您没有访问系统监控的权限，请联系管理员为您分配相应的权限。"
        />
        <a-spin :spinning="loading">
            <a-row :gutter="[16, 16]">
                <!-- 系统运行时间 -->
                <a-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
                    <a-card title="系统运行时间" :bordered="false" class="monitor-card">
                        <a-statistic :value="monitorData?.uptime || '未知'" />
                    </a-card>
                </a-col>

                <!-- CPU信息 -->
                <a-col :xs="24" :sm="12" :md="8" :lg="8" :xl="8">
                    <a-card title="CPU使用率" :bordered="false" class="monitor-card">
                        <a-progress
                            type="circle"
                            :percent="Math.round(monitorData?.cpu?.usage || 0)"
                            :status="getUsageStatus(monitorData?.cpu?.usage || 0)"
                            :format="(percent) => `${percent}%`"
                            :stroke-width="8"
                            :width="120"
                        />
                        <div class="monitor-info">
                            <div class="info-item">
                                <span class="label">CPU核心数：</span>
                                <span class="value">{{ monitorData?.cpu?.count || 0 }}</span>
                            </div>
                            <div class="info-item">
                                <span class="label">CPU型号：</span>
                                <span class="value">{{ monitorData?.cpu?.model_name || '未知' }}</span>
                            </div>
                        </div>
                    </a-card>
                </a-col>

                <!-- 内存信息 -->
                <a-col :xs="24" :sm="12" :md="8" :lg="8" :xl="8">
                    <a-card title="内存使用率" :bordered="false" class="monitor-card">
                        <a-progress
                            type="circle"
                            :percent="Math.round(monitorData?.memory?.usage || 0)"
                            :status="getUsageStatus(monitorData?.memory?.usage || 0)"
                            :format="(percent) => `${percent}%`"
                            :stroke-width="8"
                            :width="120"
                        />
                        <div class="monitor-info">
                            <div class="info-item">
                                <span class="label">总内存：</span>
                                <span class="value">{{ formatGB(monitorData?.memory?.total_gb || 0) }} GB</span>
                            </div>
                            <div class="info-item">
                                <span class="label">已使用：</span>
                                <span class="value">{{ formatGB(monitorData?.memory?.used_gb || 0) }} GB</span>
                            </div>
                            <div class="info-item">
                                <span class="label">可用：</span>
                                <span class="value">{{ formatGB(monitorData?.memory?.available_gb || 0) }} GB</span>
                            </div>
                        </div>
                    </a-card>
                </a-col>

                <!-- 磁盘信息 -->
                <a-col :xs="24" :sm="12" :md="8" :lg="8" :xl="8">
                    <a-card title="磁盘使用率" :bordered="false" class="monitor-card">
                        <a-progress
                            type="circle"
                            :percent="Math.round(monitorData?.disk?.usage || 0)"
                            :status="getUsageStatus(monitorData?.disk?.usage || 0)"
                            :format="(percent) => `${percent}%`"
                            :stroke-width="8"
                            :width="120"
                        />
                        <div class="monitor-info">
                            <div class="info-item">
                                <span class="label">总空间：</span>
                                <span class="value">{{ formatGB(monitorData?.disk?.total_gb || 0) }} GB</span>
                            </div>
                            <div class="info-item">
                                <span class="label">已使用：</span>
                                <span class="value">{{ formatGB(monitorData?.disk?.used_gb || 0) }} GB</span>
                            </div>
                            <div class="info-item">
                                <span class="label">可用：</span>
                                <span class="value">{{ formatGB(monitorData?.disk?.available_gb || 0) }} GB</span>
                            </div>
                        </div>
                    </a-card>
                </a-col>

                <!-- 磁盘分区列表 -->
                <a-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
                    <a-card title="磁盘分区详情" :bordered="false" class="monitor-card">
                        <a-table
                            :columns="diskColumns"
                            :data-source="monitorData?.disk?.disks || []"
                            :pagination="false"
                            row-key="device"
                        >
                            <template #bodyCell="{ column, record }">
                                <template v-if="column.key === 'usage'">
                                    <a-progress
                                        :percent="Math.round(record.usage)"
                                        :status="getUsageStatus(record.usage)"
                                        :show-info="true"
                                    />
                                </template>
                                <template v-else-if="column.key === 'total'">
                                    {{ formatGB(record.total_gb) }} GB
                                </template>
                                <template v-else-if="column.key === 'used'">
                                    {{ formatGB(record.used_gb) }} GB
                                </template>
                                <template v-else-if="column.key === 'available'">
                                    {{ formatGB(record.available_gb) }} GB
                                </template>
                            </template>
                        </a-table>
                    </a-card>
                </a-col>
            </a-row>
        </a-spin>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { opsApi, type SystemMonitorResponse } from '@/api/ops'
import type { TableColumnsType } from 'ant-design-vue'
import { PermissionDeniedAlert } from '@/views/pages/components'

const loading = ref(false)
const permissionDenied = ref(false)
const monitorData = ref<SystemMonitorResponse | null>(null)

// 磁盘列定义
const diskColumns: TableColumnsType = [
    {
        title: '设备',
        key: 'device',
        dataIndex: 'device',
        width: 150,
        align: 'center',
    },
    {
        title: '挂载点',
        key: 'mountpoint',
        dataIndex: 'mountpoint',
        width: 200,
        align: 'center',
    },
    {
        title: '文件系统',
        key: 'fstype',
        dataIndex: 'fstype',
        width: 120,
        align: 'center',
    },
    {
        title: '总空间',
        key: 'total',
        width: 120,
        align: 'center',
    },
    {
        title: '已使用',
        key: 'used',
        width: 120,
        align: 'center',
    },
    {
        title: '可用空间',
        key: 'available',
        width: 120,
        align: 'center',
    },
    {
        title: '使用率',
        key: 'usage',
        width: 200,
        align: 'center',
    },
]

// 获取监控数据
const fetchMonitorData = async () => {
    try {
        loading.value = true
        const response = await opsApi.getSystemMonitor()
        
        // 检查权限错误
        if (response.code === 403) {
            permissionDenied.value = true
            throw new Error(response.msg || response.message || '没有权限访问该资源')
        }
        
        // 权限验证通过，隐藏权限提示
        permissionDenied.value = false
        
        if (response.code === 200 || response.code === 0) {
            monitorData.value = response.data
        }
    } catch (error: any) {
        console.error('获取系统监控信息失败:', error)
        if (error.response?.status === 403 || error.message?.includes('没有权限')) {
            permissionDenied.value = true
        }
    } finally {
        loading.value = false
    }
}

// 格式化GB
const formatGB = (gb: number): string => {
    return gb.toFixed(2)
}

// 获取使用率状态
const getUsageStatus = (usage: number): 'success' | 'exception' | 'normal' | 'active' => {
    if (usage >= 90) return 'exception'
    if (usage >= 70) return 'active'
    return 'normal'
}

let refreshTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
    fetchMonitorData()
    // 每5秒刷新一次
    refreshTimer = setInterval(fetchMonitorData, 5000)
})

onUnmounted(() => {
    if (refreshTimer) {
        clearInterval(refreshTimer)
    }
})
</script>

<style scoped lang="less">
.system-monitor-container {
    padding: 0;
    background: transparent;
    min-height: 100%;

    .monitor-card {
        border-radius: 12px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
        transition: box-shadow 0.3s;

        &:hover {
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
        }

        :deep(.ant-card-head) {
            border-bottom: 2px solid #f0f0f0;
            padding: 12px 20px;
            background: linear-gradient(135deg, #ffffff 0%, #fafafa 100%);
        }

        :deep(.ant-card-body) {
            padding: 20px;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 16px;
        }
    }

    .monitor-info {
        width: 100%;
        margin-top: 16px;

        .info-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 0;
            border-bottom: 1px solid #f0f0f0;

            &:last-child {
                border-bottom: none;
            }

            .label {
                color: #8c8c8c;
                font-size: 14px;
            }

            .value {
                color: #262626;
                font-size: 14px;
                font-weight: 500;
            }
        }
    }
}

// 暗色模式
body.dark-mode {
    .system-monitor-container {
        .monitor-card {
            background: #1f1f1f;
            border-color: #303030;

            :deep(.ant-card-head) {
                border-bottom-color: #303030;
                background: linear-gradient(135deg, #1f1f1f 0%, #262626 100%);
            }

            .monitor-info {
                .info-item {
                    border-bottom-color: #434343;

                    .label {
                        color: rgba(255, 255, 255, 0.65);
                    }

                    .value {
                        color: rgba(255, 255, 255, 0.85);
                    }
                }
            }
        }
    }
}
</style>


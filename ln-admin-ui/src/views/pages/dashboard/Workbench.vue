<template>
    <div class="workbench-container">
        <!-- 欢迎横幅 -->
        <WelcomeBanner />

        <!-- 统计卡片 -->
        <StatsCards :stats="workbenchConfig?.config?.stats || []" />

        <!-- 主要功能区域 -->
        <a-row :gutter="[16, 16]" class="main-content-row">
            <!-- 左侧：快捷入口和公告 -->
            <a-col :xs="24" :lg="workbenchConfig?.config?.showSystemInfo ? 16 : 24">
                <!-- 快捷入口 -->
                <QuickActions 
                    :actions="workbenchConfig?.config?.quickActions || []" 
                    style="margin-bottom: 16px;" 
                />

                <!-- 系统公告 -->
                <Announcements :announcements="workbenchConfig?.config?.announcements || []" />
            </a-col>

            <!-- 右侧：系统信息 -->
            <a-col v-if="workbenchConfig?.config?.showSystemInfo" :xs="24" :lg="8">
                <SystemInfo />
            </a-col>
        </a-row>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { WelcomeBanner, StatsCards, QuickActions, Announcements, SystemInfo } from './components'
import { workbenchApi, type WorkbenchConfigResponse } from '@/api/workbench'

const workbenchConfig = ref<WorkbenchConfigResponse | null>(null)
const loading = ref(false)

// 从后端获取工作台配置
const fetchWorkbenchConfig = async () => {
    try {
        loading.value = true
        const res = await workbenchApi.getWorkbenchConfig()
        if (res.code === 200 || res.code === 0) {
            workbenchConfig.value = res.data
        }
    } catch (error) {
        console.error('获取工作台配置失败:', error)
        // 失败时使用空配置，组件会显示默认内容
        workbenchConfig.value = null
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchWorkbenchConfig()
})
</script>

<style scoped lang="less">
.workbench-container {
    padding: 24px;
    min-height: 100%;
}

.main-content-row {
    margin-top: 0;
}

// 响应式设计
@media (max-width: 768px) {
    .workbench-container {
        padding: 16px;
    }
}
</style>

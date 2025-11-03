<template>
    <a-card title="系统信息" class="section-card" :bordered="false">
        <a-descriptions :column="1" bordered size="small">
            <a-descriptions-item label="系统版本">v1.0.0</a-descriptions-item>
            <a-descriptions-item label="运行时间">{{ systemUptime }}</a-descriptions-item>
            <a-descriptions-item label="服务器时间">{{ serverTime }}</a-descriptions-item>
            <a-descriptions-item label="在线用户">{{ onlineUsers }}</a-descriptions-item>
        </a-descriptions>
    </a-card>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue'

const serverTime = ref('')
const systemUptime = ref('0 天 0 小时')
const onlineUsers = ref(128)

let timeInterval: number | null = null

// 计算运行时间（模拟）
const calculateUptime = () => {
    const days = Math.floor(Math.random() * 30)
    const hours = Math.floor(Math.random() * 24)
    systemUptime.value = `${days} 天 ${hours} 小时`
}

// 更新时间
const updateTime = () => {
    const now = new Date()
    serverTime.value = now.toLocaleTimeString('zh-CN')
}

onMounted(() => {
    updateTime()
    calculateUptime()
    timeInterval = window.setInterval(updateTime, 1000)
})

onUnmounted(() => {
    if (timeInterval) {
        clearInterval(timeInterval)
    }
})
</script>

<style scoped lang="less">
.section-card {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    :deep(.ant-card-head) {
        border-bottom: 1px solid #f0f0f0;
    }

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }

    :deep(.ant-descriptions-bordered .ant-descriptions-item-label) {
        background: #fafafa;
    }
}

body.dark-mode {
    .section-card {
        background: #1f1f1f;
        border-color: #303030;

        :deep(.ant-card-head) {
            border-color: #303030;
        }

        :deep(.ant-descriptions-bordered .ant-descriptions-item-label) {
            background: #262626;
        }

        :deep(.ant-descriptions-bordered .ant-descriptions-item-content) {
            background: #1f1f1f;
        }
    }
}
</style>


<template>
    <a-card title="最近操作" class="section-card" :bordered="false">
        <a-timeline>
            <a-timeline-item
                v-for="(log, index) in logs"
                :key="index"
                :color="log.color"
            >
                <div class="log-content">
                    <div class="log-title">{{ log.title }}</div>
                    <div class="log-time">{{ log.time }}</div>
                </div>
            </a-timeline-item>
        </a-timeline>
        <div v-if="logs.length === 0" class="empty-logs">
            <a-empty description="暂无操作记录" :image="false" />
        </div>
    </a-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

export interface LogItem {
    title: string
    time: string
    color: string
}

const props = defineProps<{
    logs?: LogItem[]
}>()

const logs = ref<LogItem[]>([
    { title: '用户 admin 登录系统', time: '10分钟前', color: 'blue' },
    { title: '创建了新角色：管理员', time: '1小时前', color: 'green' },
    { title: '修改了用户权限配置', time: '3小时前', color: 'orange' },
    { title: '备份了数据库', time: '1天前', color: 'gray' },
])

// 如果传入了外部数据，使用外部数据
if (props.logs) {
    logs.value = props.logs
}
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
}

.log-content {
    .log-title {
        font-size: 14px;
        color: #262626;
        margin-bottom: 4px;
    }

    .log-time {
        font-size: 12px;
        color: #8c8c8c;
    }
}

.empty-logs {
    padding: 20px 0;
}

body.dark-mode {
    .section-card {
        background: #1f1f1f;
        border-color: #303030;

        :deep(.ant-card-head) {
            border-color: #303030;
        }
    }

    .log-content .log-title {
        color: rgba(255, 255, 255, 0.85);
    }
}
</style>


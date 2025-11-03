<template>
    <a-card title="系统公告" class="section-card" :bordered="false">
        <a-spin :spinning="loading">
            <div class="announcements-container">
                <transition-group name="announcement-list" tag="div">
                    <div
                        v-for="item in announcementsList"
                        :key="item.key"
                        class="announcement-item"
                        :class="{
                            'announcement-important': item.type === 'important',
                            'announcement-warning': item.type === 'warning',
                        }"
                        @click="handleClick(item)"
                    >
                        <div class="announcement-icon-wrapper">
                            <div
                                class="announcement-icon"
                                :class="`icon-${item.type}`"
                            >
                                <BellOutlined />
                            </div>
                        </div>
                        <div class="announcement-content-wrapper">
                            <div class="announcement-header">
                                <span class="announcement-title">{{ item.title }}</span>
                                <div class="announcement-tags">
                                    <a-tag v-if="item.type === 'important'" color="red" size="small">
                                        <template #icon>
                                            <ExclamationCircleOutlined />
                                        </template>
                                        重要
                                    </a-tag>
                                    <a-tag v-else-if="item.type === 'warning'" color="orange" size="small">
                                        <template #icon>
                                            <WarningOutlined />
                                        </template>
                                        提醒
                                    </a-tag>
                                </div>
                            </div>
                            <div class="announcement-content">{{ item.content }}</div>
                            <div class="announcement-footer">
                                <span class="announcement-time">
                                    <ClockCircleOutlined />
                                    {{ item.time }}
                                </span>
                                <span class="announcement-author">
                                    <UserOutlined />
                                    {{ item.author }}
                                </span>
                            </div>
                        </div>
                        <div class="announcement-arrow">
                            <RightOutlined />
                        </div>
                    </div>
                </transition-group>
                <a-empty v-if="announcementsList.length === 0" description="暂无公告" :image="false" />
            </div>
        </a-spin>
        <div class="announcement-footer-actions" v-if="announcementsList.length > 0">
            <a-button type="link" @click="handleViewAll">
                查看全部
                <RightOutlined />
            </a-button>
        </div>
    </a-card>

    <!-- 公告详情模态框 -->
    <a-modal
        v-model:open="detailVisible"
        :title="currentAnnouncement?.title"
        :footer="null"
        width="600px"
        class="announcement-detail-modal"
        @cancel="handleCloseDetail"
    >
        <div v-if="currentAnnouncement" class="announcement-detail">
            <div class="detail-header">
                <div class="detail-title-wrapper">
                    <div
                        class="detail-icon"
                        :class="`icon-${currentAnnouncement.type}`"
                    >
                        <BellOutlined v-if="currentAnnouncement.type === 'normal'" />
                        <ExclamationCircleOutlined v-else-if="currentAnnouncement.type === 'important'" />
                        <WarningOutlined v-else />
                    </div>
                    <div class="detail-title-info">
                        <h3 class="detail-title">{{ currentAnnouncement.title }}</h3>
                        <div class="detail-tags">
                            <a-tag v-if="currentAnnouncement.type === 'important'" color="red">
                                <template #icon>
                                    <ExclamationCircleOutlined />
                                </template>
                                重要公告
                            </a-tag>
                            <a-tag v-else-if="currentAnnouncement.type === 'warning'" color="orange">
                                <template #icon>
                                    <WarningOutlined />
                                </template>
                                提醒公告
                            </a-tag>
                            <a-tag v-else color="blue">
                                <template #icon>
                                    <BellOutlined />
                                </template>
                                普通公告
                            </a-tag>
                        </div>
                    </div>
                </div>
            </div>
            <a-divider />
            <div class="detail-content">
                <div class="detail-content-text">{{ currentAnnouncement.content }}</div>
            </div>
            <a-divider />
            <div class="detail-footer">
                <div class="detail-meta">
                    <span class="detail-time">
                        <ClockCircleOutlined />
                        发布时间：{{ currentAnnouncement.time }}
                    </span>
                    <span class="detail-author">
                        <UserOutlined />
                        发布人：{{ currentAnnouncement.author }}
                    </span>
                </div>
            </div>
        </div>
    </a-modal>

    <!-- 全部公告抽屉 -->
    <a-drawer
        v-model:open="drawerVisible"
        title="全部公告"
        placement="right"
        width="600px"
        class="announcements-drawer"
    >
        <div class="drawer-content">
            <transition-group name="announcement-list" tag="div">
                <div
                    v-for="item in announcementsList"
                    :key="item.key"
                    class="announcement-item drawer-item"
                    :class="{
                        'announcement-important': item.type === 'important',
                        'announcement-warning': item.type === 'warning',
                    }"
                    @click="handleClick(item)"
                >
                    <div class="announcement-icon-wrapper">
                        <div
                            class="announcement-icon"
                            :class="`icon-${item.type}`"
                        >
                            <BellOutlined />
                        </div>
                    </div>
                    <div class="announcement-content-wrapper">
                        <div class="announcement-header">
                            <span class="announcement-title">{{ item.title }}</span>
                            <div class="announcement-tags">
                                <a-tag v-if="item.type === 'important'" color="red" size="small">
                                    <template #icon>
                                        <ExclamationCircleOutlined />
                                    </template>
                                    重要
                                </a-tag>
                                <a-tag v-else-if="item.type === 'warning'" color="orange" size="small">
                                    <template #icon>
                                        <WarningOutlined />
                                    </template>
                                    提醒
                                </a-tag>
                            </div>
                        </div>
                        <div class="announcement-content">{{ item.content }}</div>
                        <div class="announcement-footer">
                            <span class="announcement-time">
                                <ClockCircleOutlined />
                                {{ item.time }}
                            </span>
                            <span class="announcement-author">
                                <UserOutlined />
                                {{ item.author }}
                            </span>
                        </div>
                    </div>
                </div>
            </transition-group>
            <a-empty v-if="announcementsList.length === 0" description="暂无公告" />
        </div>
    </a-drawer>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import {
    BellOutlined,
    RightOutlined,
    ExclamationCircleOutlined,
    WarningOutlined,
    ClockCircleOutlined,
    UserOutlined,
} from '@ant-design/icons-vue'

export interface Announcement {
    key: string
    title: string
    content: string
    time: string
    author: string
    type: 'normal' | 'important' | 'warning'
}

interface Props {
    announcements?: Announcement[]
}

const props = withDefaults(defineProps<Props>(), {
    announcements: undefined,
})

const loading = ref(false)
const detailVisible = ref(false)
const drawerVisible = ref(false)
const currentAnnouncement = ref<Announcement | null>(null)

const defaultAnnouncements = ref<Announcement[]>([
    {
        key: '1',
        title: '系统升级通知',
        content: '系统将于本周六晚上22:00-24:00进行升级维护，期间可能影响部分功能使用，请提前做好准备。升级内容包括：\n1. 优化系统性能，提升响应速度\n2. 修复已知Bug，增强系统稳定性\n3. 新增数据备份功能，保障数据安全\n\n请各位用户提前保存重要数据，感谢配合！',
        time: '2024-01-15 10:30',
        author: '系统管理员',
        type: 'important',
    },
    {
        key: '2',
        title: '安全提醒',
        content: '请定期修改登录密码，建议使用强密码（包含大小写字母、数字和特殊字符）。为了您的账户安全，请注意以下几点：\n1. 密码长度至少8位\n2. 定期更换密码，建议每3个月更换一次\n3. 不要在多个平台使用相同密码\n4. 发现异常登录请及时联系管理员',
        time: '2024-01-14 14:20',
        author: '安全部门',
        type: 'warning',
    },
    {
        key: '3',
        title: '新功能上线',
        content: '数据导出功能已上线，现在可以导出用户列表、操作日志等数据。支持导出格式包括：\n1. Excel格式（.xlsx）\n2. CSV格式（.csv）\n3. PDF格式（.pdf）\n\n您可以在各个列表页面的右上角找到"导出"按钮，选择需要的格式进行导出。',
        time: '2024-01-13 09:15',
        author: '产品团队',
        type: 'normal',
    },
])

const announcementsList = computed(() => {
    return props.announcements && props.announcements.length > 0 
        ? props.announcements 
        : defaultAnnouncements.value
})

const handleClick = (item: Announcement) => {
    currentAnnouncement.value = item
    detailVisible.value = true
}

const handleCloseDetail = () => {
    detailVisible.value = false
    currentAnnouncement.value = null
}

const handleViewAll = () => {
    drawerVisible.value = true
}
</script>

<style scoped lang="less">
.section-card {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.3s;

    :deep(.ant-card-head) {
        border-bottom: 1px solid #f0f0f0;
    }

    :deep(.ant-card-head-title) {
        font-weight: 600;
    }
}

.announcements-container {
    min-height: 100px;
}

.announcement-item {
    display: flex;
    align-items: flex-start;
    padding: 16px;
    margin-bottom: 12px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    background: #fff;
    position: relative;
    overflow: hidden;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(var(--color-primary-rgb, 24, 144, 255), 0.15);
        transform: translateY(-2px);

        .announcement-arrow {
            opacity: 1;
            transform: translateX(0);
        }
    }

    &:last-child {
        margin-bottom: 0;
    }

    &.announcement-important {
        border-left: 4px solid #f5222d;
        background: linear-gradient(to right, rgba(245, 34, 45, 0.05), #fff);
    }

    &.announcement-warning {
        border-left: 4px solid #faad14;
        background: linear-gradient(to right, rgba(250, 173, 20, 0.05), #fff);
    }
}

.announcement-icon-wrapper {
    flex-shrink: 0;
    margin-right: 12px;
}

.announcement-icon {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    transition: all 0.3s;

    &.icon-important {
        background: rgba(245, 34, 45, 0.1);
        color: #f5222d;
    }

    &.icon-warning {
        background: rgba(250, 173, 20, 0.1);
        color: #faad14;
    }

    &.icon-normal {
        background: rgba(24, 144, 255, 0.1);
        color: var(--color-primary);
    }
}

.announcement-content-wrapper {
    flex: 1;
    min-width: 0;
}

.announcement-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
    gap: 8px;
}

.announcement-title {
    font-size: 15px;
    font-weight: 500;
    color: #262626;
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.announcement-tags {
    flex-shrink: 0;
}

.announcement-content {
    font-size: 14px;
    color: #595959;
    line-height: 1.6;
    margin-bottom: 12px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

.announcement-footer {
    display: flex;
    align-items: center;
    gap: 16px;
    font-size: 12px;
    color: #8c8c8c;
}

.announcement-time,
.announcement-author {
    display: inline-flex;
    align-items: center;
    gap: 4px;
}

.announcement-arrow {
    flex-shrink: 0;
    margin-left: 8px;
    color: #bfbfbf;
    font-size: 14px;
    opacity: 0;
    transform: translateX(-8px);
    transition: all 0.3s;
}

.announcement-footer-actions {
    text-align: center;
    padding-top: 16px;
    border-top: 1px solid #f0f0f0;
    margin-top: 16px;

    :deep(.ant-btn-link) {
        color: var(--color-primary);
        padding: 0;
        height: auto;

        &:hover {
            color: var(--color-primary-hover);
        }
    }
}

// 列表动画
.announcement-list-enter-active,
.announcement-list-leave-active {
    transition: all 0.3s ease;
}

.announcement-list-enter-from {
    opacity: 0;
    transform: translateY(-10px);
}

.announcement-list-leave-to {
    opacity: 0;
    transform: translateX(20px);
}

// 响应式设计
@media (max-width: 768px) {
    .announcement-item {
        padding: 12px;
    }

    .announcement-icon {
        width: 32px;
        height: 32px;
        font-size: 16px;
    }

    .announcement-title {
        font-size: 14px;
    }

    .announcement-content {
        font-size: 13px;
    }
}

// 详情模态框样式
.announcement-detail-modal {
    :deep(.ant-modal-header) {
        border-bottom: 1px solid #f0f0f0;
    }

    :deep(.ant-modal-body) {
        padding: 24px;
    }
}

.announcement-detail {
    .detail-header {
        margin-bottom: 16px;
    }

    .detail-title-wrapper {
        display: flex;
        align-items: flex-start;
        gap: 16px;
    }

    .detail-icon {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 24px;
        flex-shrink: 0;

        &.icon-important {
            background: rgba(245, 34, 45, 0.1);
            color: #f5222d;
        }

        &.icon-warning {
            background: rgba(250, 173, 20, 0.1);
            color: #faad14;
        }

        &.icon-normal {
            background: rgba(24, 144, 255, 0.1);
            color: var(--color-primary);
        }
    }

    .detail-title-info {
        flex: 1;
        min-width: 0;
    }

    .detail-title {
        font-size: 20px;
        font-weight: 600;
        color: #262626;
        margin: 0 0 12px 0;
        line-height: 1.4;
    }

    .detail-tags {
        margin-top: 8px;
    }

    .detail-content {
        padding: 16px 0;
        min-height: 100px;
    }

    .detail-content-text {
        font-size: 15px;
        color: #595959;
        line-height: 1.8;
        white-space: pre-wrap;
        word-break: break-word;
    }

    .detail-footer {
        margin-top: 16px;
    }

    .detail-meta {
        display: flex;
        align-items: center;
        gap: 24px;
        font-size: 14px;
        color: #8c8c8c;

        span {
            display: inline-flex;
            align-items: center;
            gap: 6px;
        }
    }
}

// 抽屉样式
.announcements-drawer {
    :deep(.ant-drawer-body) {
        padding: 24px;
    }
}

.drawer-content {
    .drawer-item {
        margin-bottom: 16px;
    }
}

// 暗色模式适配
body.dark-mode {
    .section-card {
        background: #1f1f1f;
        border-color: #303030;

        :deep(.ant-card-head) {
            border-color: #303030;
        }
    }

    .announcement-item {
        background: #1f1f1f;
        border-color: #303030;

        &:hover {
            border-color: var(--color-primary);
            background: #262626;
        }

        &.announcement-important {
            background: linear-gradient(to right, rgba(245, 34, 45, 0.1), #1f1f1f);
        }

        &.announcement-warning {
            background: linear-gradient(to right, rgba(250, 173, 20, 0.1), #1f1f1f);
        }
    }

    .announcement-title {
        color: rgba(255, 255, 255, 0.85);
    }

    .announcement-content {
        color: rgba(255, 255, 255, 0.65);
    }

    .announcement-footer-actions {
        border-color: #303030;
    }

    .announcement-detail {
        .detail-title {
            color: rgba(255, 255, 255, 0.85);
        }

        .detail-content-text {
            color: rgba(255, 255, 255, 0.65);
        }
    }

    .announcement-detail-modal {
        :deep(.ant-modal-header) {
            border-color: #303030;
            background: #1f1f1f;
        }

        :deep(.ant-modal-body) {
            background: #1f1f1f;
        }

        :deep(.ant-modal-title) {
            color: rgba(255, 255, 255, 0.85);
        }
    }

    .announcements-drawer {
        :deep(.ant-drawer-header) {
            border-color: #303030;
            background: #1f1f1f;
        }

        :deep(.ant-drawer-title) {
            color: rgba(255, 255, 255, 0.85);
        }

        :deep(.ant-drawer-body) {
            background: #1f1f1f;
        }
    }
}
</style>

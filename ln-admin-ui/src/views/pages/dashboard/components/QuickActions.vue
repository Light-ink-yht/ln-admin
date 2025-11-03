<template>
    <a-card title="快捷入口" class="section-card" :bordered="false">
                    <div class="quick-actions" v-if="actions.length > 0">
            <div
                v-for="action in actions"
                :key="action.key"
                class="quick-action-item"
                @click="handleAction(action)"
            >
                <div class="action-icon" :style="{ background: action.color }">
                    <component :is="getIconComponent(action.icon)" />
                </div>
                <div class="action-info">
                    <div class="action-title">{{ action.title }}</div>
                    <div class="action-desc">{{ action.desc }}</div>
                </div>
                <RightOutlined class="action-arrow" />
            </div>
        </div>
    </a-card>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
    UserAddOutlined,
    SafetyOutlined,
    UnlockOutlined,
    AuditOutlined,
    SettingOutlined,
    DatabaseOutlined,
    RightOutlined,
    TeamOutlined,
    FileTextOutlined,
    UserOutlined,
    BellOutlined,
} from '@ant-design/icons-vue'
import type { Component } from 'vue'

export interface QuickAction {
    key: string
    title: string
    desc: string
    icon: Component | string
    color: string
    path: string
}

// 图标映射
const iconMap: Record<string, Component> = {
    UserAddOutlined,
    SafetyOutlined,
    UnlockOutlined,
    AuditOutlined,
    SettingOutlined,
    DatabaseOutlined,
    TeamOutlined,
    FileTextOutlined,
    UserOutlined,
    BellOutlined,
}

const getIconComponent = (icon: Component | string): Component => {
    if (typeof icon === 'string') {
        return iconMap[icon] || UserOutlined
    }
    return icon
}

const props = defineProps<{
    actions?: QuickAction[]
}>()

const router = useRouter()

// 默认快捷操作（当没有传入数据时使用）
const defaultActions = ref<QuickAction[]>([
    {
        key: 'user',
        title: '用户管理',
        desc: '管理系统用户',
        icon: 'UserAddOutlined',
        color: 'var(--color-primary)',
        path: '/user/list',
    },
    {
        key: 'role',
        title: '角色管理',
        desc: '配置用户角色',
        icon: 'SafetyOutlined',
        color: 'var(--color-primary)',
        path: '/role/list',
    },
    {
        key: 'permission',
        title: '权限管理',
        desc: '管理权限配置',
        icon: 'UnlockOutlined',
        color: 'var(--color-primary)',
        path: '/permission/list',
    },
    {
        key: 'log',
        title: '操作日志',
        desc: '查看系统日志',
        icon: 'AuditOutlined',
        color: 'var(--color-primary)',
        path: '/system/log',
    },
    {
        key: 'config',
        title: '系统配置',
        desc: '系统参数设置',
        icon: 'SettingOutlined',
        color: 'var(--color-primary)',
        path: '/system/config',
    },
    {
        key: 'database',
        title: '数据管理',
        desc: '数据库维护',
        icon: 'DatabaseOutlined',
        color: 'var(--color-primary)',
        path: '/system/config',
    },
])

// 使用传入的actions或默认actions
const actions = computed(() => {
    if (props.actions && props.actions.length > 0) {
        return props.actions
    }
    return defaultActions.value
})

const handleAction = (action: QuickAction) => {
    router.push(action.path)
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

.quick-actions {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
}

.quick-action-item {
    display: flex;
    align-items: center;
    padding: 16px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    background: #fff;

    &:hover {
        border-color: var(--color-primary);
        box-shadow: 0 2px 8px rgba(var(--color-primary-rgb, 24, 144, 255), 0.15);
        transform: translateY(-2px);
    }

    .action-icon {
        width: 48px;
        height: 48px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        font-size: 24px;
        margin-right: 12px;
        flex-shrink: 0;
    }

    .action-info {
        flex: 1;
        min-width: 0;
    }

    .action-title {
        font-size: 16px;
        font-weight: 500;
        margin-bottom: 4px;
        color: #262626;
    }

    .action-desc {
        font-size: 12px;
        color: #8c8c8c;
    }

    .action-arrow {
        color: #bfbfbf;
        font-size: 14px;
        margin-left: 8px;
    }
}

@media (max-width: 768px) {
    .quick-actions {
        grid-template-columns: 1fr;
    }
}

body.dark-mode {
    .section-card {
        background: #1f1f1f;
        border-color: #303030;

        :deep(.ant-card-head) {
            border-color: #303030;
        }
    }

    .quick-action-item {
        background: #1f1f1f;
        border-color: #303030;

        .action-title {
            color: rgba(255, 255, 255, 0.85);
        }

        &:hover {
            border-color: var(--color-primary);
            background: #262626;
        }
    }
}
</style>


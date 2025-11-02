<template>
    <a-layout-footer class="footer" :style="footerStyle">
        <div class="footer-content">
            <div class="footer-info">
                <span>© 2024 LN Admin</span>
                <a-divider type="vertical" />
                <span>基于 Vue 3 + Ant Design Vue</span>
            </div>
            <div class="footer-links">
                <a href="#" @click.prevent>帮助文档</a>
                <a-divider type="vertical" />
                <a href="#" @click.prevent>问题反馈</a>
                <a-divider type="vertical" />
                <a href="#" @click.prevent>关于我们</a>
            </div>
        </div>
    </a-layout-footer>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useThemeStore } from '@/stores/modules/theme'

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.isDark)

const footerStyle = computed(() => {
    if (isDark.value) {
        return {
            background: '#000000',
            borderTop: '1px solid rgba(255, 255, 255, 0.1)',
            color: 'rgba(255, 255, 255, 0.65)',
        }
    } else {
        return {
            background: '#fafafa',
            borderTop: '1px solid #f0f0f0',
            color: 'rgba(0, 0, 0, 0.65)',
        }
    }
})
</script>

<style scoped>
.footer {
    padding: 16px 24px;
    text-align: center;
    transition: background-color 0.3s, border-color 0.3s, color 0.3s;
    flex-shrink: 0;
}

.footer-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
}

.footer-info {
    display: flex;
    align-items: center;
    gap: 8px;
}

.footer-links {
    display: flex;
    align-items: center;
    gap: 8px;
}

.footer-links a {
    color: inherit;
    text-decoration: none;
    transition: color 0.3s;
}

.footer-links a:hover {
    color: var(--color-primary, #1890ff);
}

:deep(.ant-divider-vertical) {
    margin: 0 8px;
}

@media (max-width: 768px) {
    .footer {
        padding: 12px 16px;
    }

    .footer-content {
        flex-direction: column;
        gap: 12px;
        text-align: center;
    }

    .footer-info,
    .footer-links {
        flex-wrap: wrap;
        justify-content: center;
    }
}

@media (max-width: 480px) {
    .footer {
        padding: 10px 12px;
        font-size: 12px;
    }

    .footer-content {
        gap: 8px;
    }

    .footer-info span,
    .footer-links a {
        font-size: 12px;
    }

    :deep(.ant-divider-vertical) {
        margin: 0 4px;
    }
}
</style>


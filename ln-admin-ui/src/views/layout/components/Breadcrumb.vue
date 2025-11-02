<template>
    <div class="breadcrumb-wrapper">
        <a-breadcrumb>
            <a-breadcrumb-item v-for="item in items" :key="item">
                {{ item }}
            </a-breadcrumb-item>
        </a-breadcrumb>
    </div>
</template>

<script lang="ts" setup>
interface Props {
    items?: string[]
}

withDefaults(defineProps<Props>(), {
    items: () => ['Home', 'List', 'App'],
})
</script>

<style scoped>
.breadcrumb-wrapper {
    padding: 16px 0;
    flex-shrink: 0;
}

/* 移动端优化 */
@media (max-width: 768px) {
    .breadcrumb-wrapper {
        padding: 12px 0;
    }

    :deep(.ant-breadcrumb) {
        font-size: 13px;
    }
}

@media (max-width: 480px) {
    .breadcrumb-wrapper {
        padding: 8px 0;
    }

    :deep(.ant-breadcrumb) {
        font-size: 12px;
    }
}
</style>


<template>
    <a-tag :color="getCategoryColor(category)">
        {{ getCategoryName(category) }}
    </a-tag>
</template>

<script lang="ts" setup>
defineProps<{
    category: string
}>()

function getCategoryName(category: string): string {
    const categoryMap: Record<string, string> = {
        image: '图片',
        document: '文档',
        video: '视频',
        audio: '音频',
        archive: '压缩包',
        other: '其他',
    }
    return categoryMap[category] || '未知'
}

function getCategoryColor(category: string): string {
    const colorMap: Record<string, string> = {
        image: 'green',
        document: 'blue',
        video: 'red',
        audio: 'orange',
        archive: 'purple',
        other: 'default',
    }
    return colorMap[category] || 'default'
}
</script>


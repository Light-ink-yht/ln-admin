<template>
    <a-badge
        :status="status === '1' ? 'success' : 'error'"
        :text="statusText"
    />
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
    status?: string
    enabledText?: string
    disabledText?: string
    unknownText?: string
}

const props = withDefaults(defineProps<Props>(), {
    status: '',
    enabledText: '启用',
    disabledText: '禁用',
    unknownText: '未知'
})

const statusText = computed(() => {
    if (props.status === '1') return props.enabledText
    if (props.status === '2') return props.disabledText
    return props.unknownText
})
</script>


<template>
    <span>{{ displayText }}</span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
    value?: string | null
    emptyText?: string
}

const props = withDefaults(defineProps<Props>(), {
    value: null,
    emptyText: '-'
})

const displayText = computed(() => {
    return props.value || props.emptyText
})
</script>


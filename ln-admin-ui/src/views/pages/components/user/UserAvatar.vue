<template>
    <a-avatar 
        :src="avatar || ''" 
        :size="size"
        :style="{ 
            backgroundColor: avatar ? 'transparent' : 'var(--color-primary)',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
        }"
    >
        {{ displayText }}
    </a-avatar>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
    avatar?: string
    nickname?: string
    fullName?: string
    size?: number
    defaultText?: string
}

const props = withDefaults(defineProps<Props>(), {
    avatar: '',
    nickname: '',
    fullName: '',
    size: 44,
    defaultText: 'U'
})

const displayText = computed(() => {
    const name = props.nickname || props.fullName || props.defaultText
    return name.charAt(0).toUpperCase()
})
</script>


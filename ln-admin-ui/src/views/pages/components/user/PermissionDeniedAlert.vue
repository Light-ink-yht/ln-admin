<template>
    <a-alert
        v-if="visible"
        message="权限不足"
        :description="description"
        type="warning"
        show-icon
        closable
        style="margin-bottom: 16px"
        @close="handleClose"
    />
</template>

<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
    visible?: boolean
    description?: string
}

const props = withDefaults(defineProps<Props>(), {
    visible: false,
    description: '您没有访问该资源的权限，请联系管理员为您分配相应的权限。',
})

const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
    (e: 'close'): void
}>()

const visible = computed({
    get: () => props.visible,
    set: (val) => emit('update:visible', val),
})

const handleClose = () => {
    visible.value = false
    emit('close')
}
</script>


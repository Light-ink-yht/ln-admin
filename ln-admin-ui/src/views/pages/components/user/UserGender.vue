<template>
    <span>{{ genderText }}</span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { getGenderText } from './userUtils'

interface Props {
    gender?: string
}

const props = withDefaults(defineProps<Props>(), {
    gender: ''
})

const genderText = computed(() => getGenderText(props.gender))
</script>


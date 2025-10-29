<template>
    <ConfigProvider :theme="antdTheme">
        <router-view />
    </ConfigProvider>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { ConfigProvider } from 'ant-design-vue'
import { useThemeStore } from '@/stores/modules/theme'
import { getAntdThemeConfig } from '@/theme/utils'
// type AntdThemeConfig = any

const themeStore = useThemeStore()

const antdTheme = computed(() => {
    return getAntdThemeConfig(themeStore.themeConfig)
})

onMounted(() => {
    themeStore.setAntdThemeConfigProvider(null)
    themeStore.initTheme()
})
</script>

<style>
:root {
    --color-primary: #1890ff;
    --color-primary-hover: #40a9ff;
    --color-primary-active: #096dd9;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
        'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
        'Noto Color Emoji';
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    transition: background-color 0.3s, color 0.3s;
}

body.dark-mode {
    background-color: #141414;
    color: rgba(255, 255, 255, 0.85);
}
</style>

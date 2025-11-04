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
    --color-primary-rgb: 24, 144, 255;
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

/* 全局滚动条样式 */
/* Webkit浏览器 (Chrome, Safari, Edge) */
::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

::-webkit-scrollbar-track {
    background: #f5f5f5;
    border-radius: 4px;
}

::-webkit-scrollbar-thumb {
    background: #d9d9d9;
    border-radius: 4px;
    transition: background 0.2s;
}

::-webkit-scrollbar-thumb:hover {
    background: #bfbfbf;
}

::-webkit-scrollbar-thumb:active {
    background: #999999;
}

::-webkit-scrollbar-corner {
    background: #f5f5f5;
}

/* 深色模式下的滚动条 */
body.dark-mode ::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

body.dark-mode ::-webkit-scrollbar-track {
    background: #262626;
    border-radius: 4px;
}

body.dark-mode ::-webkit-scrollbar-thumb {
    background: #434343;
    border-radius: 4px;
    transition: background 0.2s;
}

body.dark-mode ::-webkit-scrollbar-thumb:hover {
    background: #595959;
}

body.dark-mode ::-webkit-scrollbar-thumb:active {
    background: #737373;
}

body.dark-mode ::-webkit-scrollbar-corner {
    background: #262626;
}

/* Firefox浏览器 */
* {
    scrollbar-width: thin;
    scrollbar-color: #d9d9d9 #f5f5f5;
}

body.dark-mode * {
    scrollbar-width: thin;
    scrollbar-color: #434343 #262626;
}
</style>

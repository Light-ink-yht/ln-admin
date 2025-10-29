import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { DEFAULT_THEME, ThemeMode, type ThemeConfig } from '@/theme/config'
import { getAntdThemeConfig, generateCSSVariables, getSystemThemeMode } from '@/theme/utils'

export const useThemeStore = defineStore('theme', () => {
    // 状态
    const colorPrimary = ref<string>(DEFAULT_THEME.colorPrimary)
    const mode = ref<ThemeMode | string>(DEFAULT_THEME.mode)
    const antdConfigProvider = ref<any>(null)

    // 计算属性
    const themeConfig = computed<ThemeConfig>(() => ({
        colorPrimary: colorPrimary.value,
        mode: mode.value as ThemeMode,
    }))

    const isDark = computed(() => mode.value === ThemeMode.DARK || mode.value === 'dark')

    // 方法
    /**
     * 设置主题色
     */
    function setColorPrimary(color: string) {
        colorPrimary.value = color
        applyTheme()
    }

    /**
     * 设置主题模式
     */
    function setThemeMode(newMode: ThemeMode) {
        mode.value = newMode
        applyTheme()
    }

    /**
     * 切换主题模式
     */
    function toggleThemeMode() {
        mode.value = mode.value === ThemeMode.LIGHT ? ThemeMode.DARK : ThemeMode.LIGHT
        applyTheme()
    }

    /**
     * 跟随系统主题
     */
    function followSystemTheme() {
        mode.value = getSystemThemeMode()
        applyTheme()
    }

    /**
     * 应用主题
     */
    function applyTheme() {
        // 更新 Ant Design 主题
        const antdTheme = getAntdThemeConfig(themeConfig.value)

        // 更新 CSS 变量
        const cssVars = generateCSSVariables(colorPrimary.value)
        const root = document.documentElement

        Object.entries(cssVars).forEach(([key, value]) => {
            root.style.setProperty(key, value)
        })

        // 更新 body class
        if (isDark.value) {
            document.body.classList.add('dark-mode')
        } else {
            document.body.classList.remove('dark-mode')
        }

        // 设置 Ant Design ConfigProvider 主题
        setAntdTheme(antdTheme)
    }

    /**
     * 设置 Ant Design 主题（供外部调用）
     */
    function setAntdThemeConfigProvider(providerInstance: any) {
        antdConfigProvider.value = providerInstance
    }

    function setAntdTheme(theme: any) {
        if (antdConfigProvider.value) {
            antdConfigProvider.value.theme = theme
        }
    }

    /**
     * 初始化主题
     */
    function initTheme() {
        // 监听系统主题变化
        if (typeof window !== 'undefined') {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
            mediaQuery.addEventListener('change', () => {
                // 可以根据要自动跟随系统主题变化
                // followSystemTheme()
            })
        }

        applyTheme()
    }

    // 监听器
    watch([colorPrimary, mode], () => {
        applyTheme()
    })

    return {
        // 状态
        colorPrimary,
        mode,
        antdConfigProvider,
        // 计算属性
        themeConfig,
        isDark,
        // 方法
        setColorPrimary,
        setThemeMode,
        toggleThemeMode,
        followSystemTheme,
        applyTheme,
        setAntdThemeConfigProvider,
        setAntdTheme,
        initTheme,
    }
}, {
    persist: {
        key: 'theme-store',
        storage: localStorage,
    },
})

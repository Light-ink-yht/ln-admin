import { theme } from 'ant-design-vue'
import { ThemeMode, PRESET_COLORS } from './config'
import type { ThemeConfig } from './config'

/**
 * 获取 Ant Design 主题配置
 */
export function getAntdThemeConfig(config: ThemeConfig): any {
    const isDark = config.mode === ThemeMode.DARK

    return {
        algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: {
            colorPrimary: config.colorPrimary,
            borderRadius: 6,
        },
        components: {
            Button: {
                borderRadius: 6,
            },
            Card: {
                borderRadius: 8,
            },
        },
    }
}

/**
 * 从十六进制颜色转换为RGB值
 */
function hexToRgb(hex: string): string {
    const num = parseInt(hex.replace('#', ''), 16)
    const r = (num >> 16) & 255
    const g = (num >> 8) & 255
    const b = num & 255
    return `${r}, ${g}, ${b}`
}

/**
 * 生成 CSS 变量
 */
export function generateCSSVariables(colorPrimary: string): Record<string, string> {
    const cssVars: Record<string, string> = {
        '--color-primary': colorPrimary,
        '--color-primary-hover': adjustColor(colorPrimary, 10),
        '--color-primary-active': adjustColor(colorPrimary, -10),
        '--color-primary-rgb': hexToRgb(colorPrimary),
    }

    return cssVars
}

/**
 * 调整颜色亮度
 */
function adjustColor(color: string, percent: number): string {
    const num = parseInt(color.replace('#', ''), 16)
    const r = (num >> 16) + percent
    const g = (num >> 8 & 0x00FF) + percent
    const b = (num & 0x0000FF) + percent

    const newR = Math.max(0, Math.min(255, r))
    const newG = Math.max(0, Math.min(255, g))
    const newB = Math.max(0, Math.min(255, b))

    return '#' + (newR << 16 | newG << 8 | newB).toString(16).padStart(6, '0')
}

/**
 * 验证是否为有效的十六进制颜色
 */
export function isValidColor(color: string): boolean {
    return /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/.test(color)
}

/**
 * 获取系统主题模式偏好
 */
export function getSystemThemeMode(): ThemeMode {
    if (typeof window !== 'undefined' && window.matchMedia) {
        return window.matchMedia('(prefers-color-scheme: dark)').matches
            ? ThemeMode.DARK
            : ThemeMode.LIGHT
    }
    return ThemeMode.LIGHT
}

/**
 * 验证主题色是否在预设色中
 */
export function isPresetColor(color: string): boolean {
    return PRESET_COLORS.includes(color.toLowerCase())
}


/**
* 主题配置
*/

// 预设主题色
export const PRESET_COLORS = [
    '#1890ff', // 蓝色（默认）
    '#f5222d', // 红色
    '#52c41a', // 绿色
    '#faad14', // 橙色
    '#13c2c2', // 青色
    '#722ed1', // 紫色
    '#eb2f96', // 粉色
    '#fa8c16', // 橙黄
]

// 主题色名称映射
export const COLOR_NAMES: Record<string, string> = {
    '#1890ff': '蓝色',
    '#f5222d': '红色',
    '#52c41a': '绿色',
    '#faad14': '橙色',
    '#13c2c2': '青色',
    '#722ed1': '紫色',
    '#eb2f96': '粉色',
    '#fa8c16': '橙黄',
}

// 主题模式
export enum ThemeMode {
    LIGHT = 'light',
    DARK = 'dark',
}

// 默认主题配置
export const DEFAULT_THEME = {
    colorPrimary: '#1890ff',
    mode: ThemeMode.LIGHT,
}

// 主题配置类型
export interface ThemeConfig {
    colorPrimary: string
    mode: ThemeMode
}


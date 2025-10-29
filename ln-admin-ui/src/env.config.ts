/**
 * 环境变量配置
 * 注意：后端地址统一在 vite.config.ts 的 BACKEND_URL 中配置
 */

// 从环境变量中读取配置
export const envConfig = {
    // API 基础地址
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api',

    // 应用标题
    appTitle: import.meta.env.VITE_APP_TITLE || 'LN Admin UI',

    // 开发环境
    dev: import.meta.env.DEV,

    // 生产环境
    prod: import.meta.env.PROD,

    // 模式
    mode: import.meta.env.MODE,

    // 基础路径
    base: import.meta.env.BASE_URL,
}

/**
 * 是否是开发环境
 */
export const isDev = envConfig.dev

/**
 * 是否是生产环境
 */
export const isProd = envConfig.prod

/**
 * 获取 API 完整地址
 */
export function getApiUrl(path: string): string {
    if (path.startsWith('http')) {
        return path
    }
    return `${envConfig.apiBaseUrl}${path}`
}

/**
 * 输出环境配置（仅在开发环境）
 */
if (isDev) {
    console.log('环境配置:', envConfig)
}

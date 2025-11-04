import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { systemApi, type SystemConfig } from '@/api/system'

export interface SystemConfigMap {
    [key: string]: string
}

export const useSystemConfigStore = defineStore('systemConfig', () => {
    // 配置列表
    const configs = ref<SystemConfig[]>([])
    // 配置映射表（key -> value）
    const configMap = ref<SystemConfigMap>({})
    // 加载状态
    const loading = ref(false)

    // 获取配置值（带默认值）
    const getConfig = (key: string, defaultValue: string = ''): string => {
        return configMap.value[key] || defaultValue
    }

    // 计算属性：常用配置
    const siteName = computed(() => getConfig('site_name', 'LN Admin'))
    const siteLogo = computed(() => getConfig('site_logo', ''))
    const siteFavicon = computed(() => getConfig('site_favicon', ''))
    const siteCopyright = computed(() => getConfig('site_copyright', '© 2024 LN Admin 基于 Vue 3 + Ant Design Vue'))
    const siteDescription = computed(() => getConfig('site_description', 'LN Admin 管理系统'))
    const siteKeywords = computed(() => getConfig('site_keywords', 'LN Admin,管理系统,后台管理'))
    const siteBeian = computed(() => getConfig('site_beian', ''))
    const siteContactEmail = computed(() => getConfig('site_contact_email', ''))
    const siteContactPhone = computed(() => getConfig('site_contact_phone', ''))
    const siteAddress = computed(() => getConfig('site_address', ''))

    // 加载配置列表
    const loadConfigs = async () => {
        loading.value = true
        try {
            const response = await systemApi.getAllConfigs()
            if (response.code === 200 || response.code === 0) {
                configs.value = response.data || []
                // 更新配置映射
                const map: SystemConfigMap = {}
                configs.value.forEach((config) => {
                    map[config.configKey] = config.configValue
                })
                configMap.value = map
                // 更新favicon
                updateFavicon()
            }
        } catch (error) {
            console.error('加载系统配置失败:', error)
        } finally {
            loading.value = false
        }
    }

    // 更新配置值
    const updateConfig = (key: string, value: string) => {
        configMap.value[key] = value
        // 更新对应的config对象
        const config = configs.value.find((c) => c.configKey === key)
        if (config) {
            config.configValue = value
        }
        // 如果是特定配置，执行相应操作
        if (key === 'site_favicon') {
            updateFavicon()
        }
        if (key === 'site_name') {
            updateTitle()
        }
    }

    // 更新favicon
    const updateFavicon = () => {
        const favicon = getConfig('site_favicon', '')
        if (favicon) {
            const link = document.querySelector("link[rel*='icon']") as HTMLLinkElement
            if (link) {
                link.href = favicon
            } else {
                const newLink = document.createElement('link')
                newLink.rel = 'icon'
                newLink.href = favicon
                document.head.appendChild(newLink)
            }
        }
    }

    // 更新页面标题
    const updateTitle = () => {
        const name = getConfig('site_name', 'LN Admin')
        document.title = name
    }

    // 初始化配置
    const initConfigs = async () => {
        if (configs.value.length === 0) {
            await loadConfigs()
        }
        // 初始化标题和favicon
        updateTitle()
        updateFavicon()
    }

    return {
        configs,
        configMap,
        loading,
        getConfig,
        siteName,
        siteLogo,
        siteFavicon,
        siteCopyright,
        siteDescription,
        siteKeywords,
        siteBeian,
        siteContactEmail,
        siteContactPhone,
        siteAddress,
        loadConfigs,
        updateConfig,
        initConfigs,
    }
}, {
    persist: {
        key: 'system-config-store',
        storage: localStorage,
        paths: ['configMap'], // 只持久化配置映射表
    },
})


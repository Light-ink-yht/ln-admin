import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type UserInfo } from '@/api/auth'
import { getRolesFromToken, parseJWTToken } from '@/utils/jwt'

export const useUserStore = defineStore('user', () => {
    // ==================== 状态 ====================
    const userInfo = ref<UserInfo | null>(null)
    const token = ref<string>('')
    const refreshToken = ref<string>('')
    const roles = ref<string[]>([]) // 用户角色列表（从JWT token中解析）

    // ==================== 计算属性 ====================
    /**
     * 是否已登录
     */
    const isLoggedIn = computed(() => !!token.value && !!userInfo.value)

    /**
     * 用户ID
     */
    const userId = computed(() => userInfo.value?.userId || '')

    /**
     * 用户名（昵称或姓名）
     */
    const userName = computed(() => {
        return userInfo.value?.nickname || userInfo.value?.fullName || '用户'
    })

    /**
     * 用户头像
     */
    const userAvatar = computed(() => userInfo.value?.avatar || '')

    /**
     * 用户角色列表（从JWT token中解析）
     */
    const userRoles = computed(() => {
        return roles.value
    })

    /**
     * 是否是超级管理员
     */
    const isSuperAdmin = computed(() => {
        return roles.value.includes('super_admin')
    })

    /**
     * 主要角色标识
     */
    const primaryRole = computed(() => {
        if (roles.value && roles.value.length > 0) {
            return roles.value[0]
        }
        return null
    })

    // ==================== 方法 ====================
    /**
     * 设置用户信息
     */
    function setUserInfo(info: UserInfo | null) {
        userInfo.value = info
    }

    /**
     * 设置Token
     */
    function setToken(accessToken: string, refresh?: string) {
        token.value = accessToken
        if (refresh) {
            refreshToken.value = refresh
        }
        // 同步到 localStorage
        localStorage.setItem('token', accessToken)
        if (refresh) {
            localStorage.setItem('refreshToken', refresh)
        }
        
        // 从token中解析角色信息
        updateRolesFromToken(accessToken)
    }

    /**
     * 从Token中更新角色信息
     */
    function updateRolesFromToken(tokenString: string) {
        const claims = parseJWTToken(tokenString)
        if (claims && claims.roles && Array.isArray(claims.roles)) {
            roles.value = claims.roles
            // 持久化角色信息到 localStorage
            localStorage.setItem('user_roles', JSON.stringify(claims.roles))
        } else {
            roles.value = []
            localStorage.removeItem('user_roles')
        }
    }

    /**
     * 从存储中恢复Token
     */
    function restoreToken() {
        const storedToken = localStorage.getItem('token')
        const storedRefreshToken = localStorage.getItem('refreshToken')
        if (storedToken) {
            token.value = storedToken
            // 从token中恢复角色信息
            updateRolesFromToken(storedToken)
        } else {
            // 如果token不存在，尝试从localStorage恢复角色信息（向后兼容）
            const storedRoles = localStorage.getItem('user_roles')
            if (storedRoles) {
                try {
                    roles.value = JSON.parse(storedRoles)
                } catch (error) {
                    console.error('解析存储的角色信息失败:', error)
                    roles.value = []
                }
            }
        }
        if (storedRefreshToken) {
            refreshToken.value = storedRefreshToken
        }
    }

    /**
     * 获取用户信息（从服务器）
     */
    async function fetchUserInfo(): Promise<UserInfo | null> {
        try {
            const res = await authApi.getUserInfo()
            if (res.code === 200 || res.code === 0) {
                userInfo.value = res.data
                return res.data
            }
            return null
        } catch (error: any) {
            console.error('获取用户信息失败:', error)
            // 如果 token 失效，清除用户信息
            if (error?.response?.status === 401 || error?.message?.includes('未授权')) {
                clearUserInfo()
            }
            return null
        }
    }

    /**
     * 初始化用户信息（从存储恢复或从服务器获取）
     */
    async function initUserInfo(): Promise<void> {
        // 先恢复 token
        restoreToken()

        // 如果有 token 但没有用户信息，从服务器获取
        if (token.value && !userInfo.value) {
            await fetchUserInfo()
        }
    }

    /**
     * 登录（设置token和用户信息）
     */
    function login(user: UserInfo, accessToken: string, refresh?: string) {
        setUserInfo(user)
        setToken(accessToken, refresh)
    }

    /**
     * 退出登录（清除所有信息）
     */
    function logout() {
        userInfo.value = null
        token.value = ''
        refreshToken.value = ''
        roles.value = []
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('user_roles')
    }

    /**
     * 清除用户信息（但保留token，用于重试）
     */
    function clearUserInfo() {
        userInfo.value = null
    }

    /**
     * 更新用户信息（部分更新）
     */
    function updateUserInfo(partialInfo: Partial<UserInfo>) {
        if (userInfo.value) {
            userInfo.value = { ...userInfo.value, ...partialInfo }
        }
    }

    return {
        // 状态
        userInfo,
        token,
        refreshToken,
        roles,
        // 计算属性
        isLoggedIn,
        userId,
        userName,
        userAvatar,
        userRoles,
        isSuperAdmin,
        primaryRole,
        // 方法
        setUserInfo,
        setToken,
        restoreToken,
        fetchUserInfo,
        initUserInfo,
        login,
        logout,
        clearUserInfo,
        updateUserInfo,
        updateRolesFromToken,
    }
}, {
    persist: {
        key: 'user-store',
        storage: localStorage,
    },
})


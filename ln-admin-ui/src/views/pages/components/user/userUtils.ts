import type { UserResponse } from '@/api/user'

/**
 * 性别文本映射
 */
const genderMap: Record<string, string> = {
    '1': '男',
    '2': '女',
    '3': '未知',
}

/**
 * 获取性别文本
 * @param gender 性别代码
 * @returns 性别文本
 */
export function getGenderText(gender: string): string {
    return genderMap[gender] || '未知'
}

/**
 * 映射用户数据字段（处理 snake_case 到 camelCase 的转换）
 * @param user 原始用户数据
 * @returns 标准化后的用户数据
 */
export function mapUserData(user: any): UserResponse {
    return {
        ...user,
        userId: user.userId || user.user_id || '',
        fullName: user.fullName || user.full_name || '',
        nickname: user.nickname || '',
        phone: user.phone || null,
        email: user.email || null,
        avatar: user.avatar || '',
        gender: user.gender || '',
        status: user.status || '',
        lastLoginTime: user.lastLoginTime || user.last_login_time || null,
        lastLoginIp: user.lastLoginIp || user.last_login_ip || '',
        loginCount: user.loginCount || user.login_count || 0,
        creatorId: user.creatorId || user.creator_id || '',
        modifierId: user.modifierId || user.modifier_id || '',
        CreatedAt: user.CreatedAt || user.created_at || '',
    }
}

/**
 * 批量映射用户数据
 * @param users 用户数据数组
 * @returns 标准化后的用户数据数组
 */
export function mapUserList(users: any[]): UserResponse[] {
    return users.map(mapUserData)
}


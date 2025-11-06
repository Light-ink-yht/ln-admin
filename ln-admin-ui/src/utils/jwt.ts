/**
 * JWT Token 解析工具
 */

/**
 * JWT Token Claims 结构
 */
export interface JWTClaims {
    user_id: string
    roles?: string[]
    exp?: number
    iat?: number
    nbf?: number
    iss?: string
}

/**
 * 解析 JWT Token
 * @param token JWT Token 字符串
 * @returns 解析后的 Claims 对象，如果解析失败返回 null
 */
export function parseJWTToken(token: string): JWTClaims | null {
    try {
        // JWT Token 格式：header.payload.signature
        const parts = token.split('.')
        if (parts.length !== 3) {
            return null
        }

        // Base64 URL 解码 payload
        const payload = parts[1]
        
        // 添加 padding（如果需要）
        let base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
        const padding = base64.length % 4
        if (padding) {
            base64 += '='.repeat(4 - padding)
        }

        // 解码 Base64
        const decoded = atob(base64)
        
        // 解析 JSON
        const claims: JWTClaims = JSON.parse(decoded)
        
        return claims
    } catch (error) {
        console.error('解析 JWT Token 失败:', error)
        return null
    }
}

/**
 * 从 Token 中提取角色列表
 * @param token JWT Token 字符串
 * @returns 角色标识列表
 */
export function getRolesFromToken(token: string): string[] {
    const claims = parseJWTToken(token)
    if (claims && claims.roles && Array.isArray(claims.roles)) {
        return claims.roles
    }
    return []
}

/**
 * 从 Token 中提取用户ID
 * @param token JWT Token 字符串
 * @returns 用户ID
 */
export function getUserIdFromToken(token: string): string | null {
    const claims = parseJWTToken(token)
    if (claims && claims.user_id) {
        return claims.user_id
    }
    return null
}

/**
 * 检查 Token 是否过期
 * @param token JWT Token 字符串
 * @returns true 如果过期，false 如果未过期
 */
export function isTokenExpired(token: string): boolean {
    const claims = parseJWTToken(token)
    if (!claims || !claims.exp) {
        return true
    }
    
    // exp 是 Unix 时间戳（秒），Date.now() 是毫秒
    const expirationTime = claims.exp * 1000
    return Date.now() >= expirationTime
}


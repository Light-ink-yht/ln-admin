import type { Role } from '@/api/role'

/**
 * 映射角色数据字段（处理 snake_case 到 camelCase 的转换）
 * @param role 原始角色数据
 * @returns 标准化后的角色数据
 */
export function mapRoleData(role: any): Role {
    return {
        ...role,
        roleId: role.roleId || role.role_id || '',
        roleKey: role.roleKey || role.role_key || '',
        roleName: role.roleName || role.role_name || '',
        description: role.description || '',
        status: role.status || '',
        creatorId: role.creatorId || role.creator_id || '',
        creatorName: role.creatorName || role.creator_name || '',
        modifierId: role.modifierId || role.modifier_id || '',
        modifierName: role.modifierName || role.modifier_name || '',
        CreatedAt: role.CreatedAt || role.created_at || role.createdAt || '',
        updatedAt: role.updatedAt || role.updated_at || '',
        createdAt: role.createdAt || role.CreatedAt || role.created_at || '',
    }
}

/**
 * 批量映射角色数据
 * @param roles 角色数据数组
 * @returns 标准化后的角色数据数组
 */
export function mapRoleList(roles: any[]): Role[] {
    return roles.map(mapRoleData)
}


import type { Permission } from '@/api/permission'

/**
 * 映射权限数据字段（处理 snake_case 到 camelCase 的转换）
 * @param permission 原始权限数据
 * @returns 标准化后的权限数据
 */
export function mapPermissionData(permission: any): Permission {
    return {
        ...permission,
        permissionId: permission.permissionId || permission.permission_id || '',
        permissionKey: permission.permissionKey || permission.permission_key || '',
        permissionName: permission.permissionName || permission.permission_name || '',
        resourcePath: permission.resourcePath || permission.resource_path || '',
        method: permission.method || '',
        description: permission.description || '',
        status: permission.status || '',
        category: permission.category || '',
        creatorId: permission.creatorId || permission.creator_id || '',
        creatorName: permission.creatorName || permission.creator_name || '',
        modifierId: permission.modifierId || permission.modifier_id || '',
        modifierName: permission.modifierName || permission.modifier_name || '',
        CreatedAt: permission.CreatedAt || permission.created_at || permission.createdAt || '',
        updatedAt: permission.updatedAt || permission.updated_at || '',
        createdAt: permission.createdAt || permission.CreatedAt || permission.created_at || '',
    }
}

/**
 * 批量映射权限数据
 * @param permissions 权限数据数组
 * @returns 标准化后的权限数据数组
 */
export function mapPermissionList(permissions: any[]): Permission[] {
    return permissions.map(mapPermissionData)
}


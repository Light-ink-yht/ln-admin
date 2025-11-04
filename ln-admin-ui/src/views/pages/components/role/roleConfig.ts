import type { SearchField } from '../DataList.vue'
import type { TableColumnsType } from 'ant-design-vue'

/**
 * 获取角色搜索字段配置
 * @param customFields 自定义字段，可以覆盖或扩展默认字段
 * @returns 搜索字段配置数组
 */
export function getRoleSearchFields(customFields?: Partial<SearchField>[]): SearchField[] {
    const defaultFields: SearchField[] = [
        {
            key: 'role_key',
            label: '角色标识',
            type: 'input',
            placeholder: '请输入角色标识',
            width: '180px',
            multiple: true,
        },
        {
            key: 'role_name',
            label: '角色名称',
            type: 'input',
            placeholder: '请输入角色名称',
            width: '180px',
            multiple: true,
        },
        {
            key: 'status',
            label: '状态',
            type: 'select',
            placeholder: '请选择状态',
            width: '120px',
            options: [
                { label: '启用', value: '1' },
                { label: '禁用', value: '2' },
            ],
        },
    ]

    if (!customFields || customFields.length === 0) {
        return defaultFields
    }

    // 合并自定义字段
    const fieldMap = new Map<string, SearchField>()
    defaultFields.forEach((field) => fieldMap.set(field.key, field))

    customFields.forEach((customField) => {
        if (customField.key) {
            fieldMap.set(customField.key, { ...fieldMap.get(customField.key), ...customField } as SearchField)
        }
    })

    return Array.from(fieldMap.values())
}

/**
 * 获取角色表格列配置
 * @param options 配置选项
 * @returns 表格列配置数组
 */
export function getRoleColumns(options?: {
    showDescription?: boolean
    showCreatedAt?: boolean
    showUpdatedAt?: boolean
    showCreatorId?: boolean
    showModifierId?: boolean
    customColumns?: TableColumnsType
}): TableColumnsType {
    const {
        showDescription = true,
        showCreatedAt = true,
        showUpdatedAt = true,
        showCreatorId = true,
        showModifierId = true,
        customColumns = [],
    } = options || {}

    const columns: TableColumnsType = [
        {
            title: '角色标识',
            dataIndex: 'roleKey',
            key: 'roleKey',
            width: 150,
            align: 'center',
        },
        {
            title: '角色名称',
            dataIndex: 'roleName',
            key: 'roleName',
            width: 150,
            align: 'center',
        },
    ]

    if (showDescription) {
        columns.push({
            title: '角色描述',
            key: 'description',
            width: 200,
            align: 'center',
        })
    }

    columns.push({
        title: '状态',
        key: 'status',
        width: 100,
        align: 'center',
    })

    if (showCreatedAt) {
        columns.push({
            title: '创建时间',
            key: 'createdAt',
            width: 180,
            align: 'center',
        })
    }

    if (showUpdatedAt) {
        columns.push({
            title: '更新时间',
            key: 'updatedAt',
            width: 180,
            align: 'center',
        })
    }

    if (showCreatorId) {
        columns.push({
            title: '创建人',
            key: 'creatorId',
            width: 120,
            align: 'center',
        })
    }

    if (showModifierId) {
        columns.push({
            title: '修改人',
            key: 'modifierId',
            width: 120,
            align: 'center',
        })
    }

    // 添加自定义列
    if (customColumns.length > 0) {
        columns.push(...customColumns)
    }

    // 添加操作列（固定在最右侧）
    columns.push({
        title: '操作',
        key: 'action',
        width: 180,
        fixed: 'right',
        align: 'center',
    })

    return columns
}

/**
 * 角色状态选项（用于搜索和筛选）
 */
export const ROLE_STATUS_OPTIONS = [
    { label: '启用', value: '1' },
    { label: '禁用', value: '2' },
] as const


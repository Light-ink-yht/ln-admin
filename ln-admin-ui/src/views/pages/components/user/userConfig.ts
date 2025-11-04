import type { SearchField } from '../DataList.vue'
import type { TableColumnsType } from 'ant-design-vue'

/**
 * 获取用户搜索字段配置
 * @param customFields 自定义字段，可以覆盖或扩展默认字段
 * @returns 搜索字段配置数组
 */
export function getUserSearchFields(customFields?: Partial<SearchField>[]): SearchField[] {
    const defaultFields: SearchField[] = [
        {
            key: 'phone',
            label: '手机号',
            type: 'input',
            placeholder: '请输入手机号',
            width: '180px',
            multiple: true,
        },
        {
            key: 'email',
            label: '邮箱',
            type: 'input',
            placeholder: '请输入邮箱',
            width: '180px',
            multiple: true,
        },
        {
            key: 'nickname',
            label: '昵称',
            type: 'input',
            placeholder: '请输入昵称',
            width: '180px',
            multiple: true,
        },
        {
            key: 'full_name',
            label: '姓名',
            type: 'input',
            placeholder: '请输入姓名',
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
        {
            key: 'gender',
            label: '性别',
            type: 'select',
            placeholder: '请选择性别',
            width: '120px',
            options: [
                { label: '男', value: '1' },
                { label: '女', value: '2' },
                { label: '未知', value: '3' },
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
 * 获取用户表格列配置
 * @param options 配置选项
 * @returns 表格列配置数组
 */
export function getUserColumns(options?: {
    showAvatar?: boolean
    showGender?: boolean
    showStatus?: boolean
    showLoginCount?: boolean
    showLastLoginTime?: boolean
    showCreatedAt?: boolean
    customColumns?: TableColumnsType
}): TableColumnsType {
    const {
        showAvatar = true,
        showGender = true,
        showStatus = true,
        showLoginCount = true,
        showLastLoginTime = true,
        showCreatedAt = true,
        customColumns = [],
    } = options || {}

    const columns: TableColumnsType = [
        {
            title: '昵称',
            dataIndex: 'nickname',
            key: 'nickname',
            width: 120,
            align: 'center',
        },
        {
            title: '姓名',
            dataIndex: 'fullName',
            key: 'fullName',
            width: 120,
            align: 'center',
        },
        {
            title: '手机号',
            dataIndex: 'phone',
            key: 'phone',
            width: 130,
            align: 'center',
        },
        {
            title: '邮箱',
            dataIndex: 'email',
            key: 'email',
            width: 180,
            align: 'center',
        },
    ]

    if (showAvatar) {
        columns.push({
            title: '头像',
            key: 'avatar',
            width: 80,
            align: 'center',
        })
    }

    if (showGender) {
        columns.push({
            title: '性别',
            key: 'gender',
            width: 80,
            align: 'center',
        })
    }

    if (showStatus) {
        columns.push({
            title: '状态',
            key: 'status',
            width: 80,
            align: 'center',
        })
    }

    if (showLoginCount) {
        columns.push({
            title: '登录次数',
            dataIndex: 'loginCount',
            key: 'loginCount',
            width: 100,
            align: 'center',
        })
    }

    if (showLastLoginTime) {
        columns.push({
            title: '最后登录时间',
            dataIndex: 'lastLoginTime',
            key: 'lastLoginTime',
            width: 180,
            align: 'center',
        })
    }

    if (showCreatedAt) {
        columns.push({
            title: '创建时间',
            dataIndex: 'CreatedAt',
            key: 'CreatedAt',
            width: 180,
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
 * 用户状态选项（用于搜索和筛选）
 */
export const USER_STATUS_OPTIONS = [
    { label: '启用', value: '1' },
    { label: '禁用', value: '2' },
] as const

/**
 * 用户性别选项（用于搜索和筛选）
 */
export const USER_GENDER_OPTIONS = [
    { label: '男', value: '1' },
    { label: '女', value: '2' },
    { label: '未知', value: '3' },
] as const


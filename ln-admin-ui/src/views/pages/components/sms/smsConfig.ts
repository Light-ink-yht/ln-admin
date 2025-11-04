import type { SearchField } from '../DataList.vue'
import type { TableColumnsType } from 'ant-design-vue'

/**
 * 获取短信模板搜索字段配置
 */
export function getSMSTemplateSearchFields(): SearchField[] {
    return [
        {
            key: 'type',
            label: '模板类型',
            type: 'input',
            placeholder: '请输入模板类型',
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
}

/**
 * 获取短信模板表格列配置
 */
export function getSMSTemplateColumns(): TableColumnsType {
    return [
        {
            title: '模板ID',
            key: 'templateId',
            dataIndex: 'templateId',
            width: 200,
            align: 'center',
        },
        {
            title: '模板类型',
            key: 'type',
            dataIndex: 'type',
            width: 150,
            align: 'center',
        },
        {
            title: '模板名称',
            key: 'templateName',
            dataIndex: 'templateName',
            width: 180,
            align: 'center',
        },
        {
            title: '腾讯云模板ID',
            key: 'tencentTemplateId',
            dataIndex: 'tencentTemplateId',
            width: 200,
            align: 'center',
        },
        {
            title: '短信签名',
            key: 'title',
            dataIndex: 'title',
            width: 150,
            align: 'center',
        },
        {
            title: '模板内容',
            key: 'content',
            dataIndex: 'content',
            width: 300,
            ellipsis: true,
            align: 'center',
        },
        {
            title: '状态',
            key: 'status',
            dataIndex: 'status',
            width: 100,
            align: 'center',
        },
        {
            title: '创建时间',
            key: 'createdAt',
            dataIndex: 'createdAt',
            width: 180,
            align: 'center',
        },
        {
            title: '操作',
            key: 'action',
            width: 180,
            align: 'center',
            fixed: 'right',
        },
    ]
}

/**
 * 获取短信验证码搜索字段配置
 */
export function getSMSCodeSearchFields(): SearchField[] {
    return [
        {
            key: 'phone',
            label: '手机号',
            type: 'input',
            placeholder: '请输入手机号',
            width: '180px',
            multiple: true,
        },
        {
            key: 'type',
            label: '验证码类型',
            type: 'input',
            placeholder: '请输入验证码类型',
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
                { label: '未使用', value: '1' },
                { label: '已使用', value: '2' },
                { label: '已过期', value: '3' },
            ],
        },
    ]
}

/**
 * 获取短信验证码表格列配置
 */
export function getSMSCodeColumns(): TableColumnsType {
    return [
        {
            title: '验证码ID',
            key: 'codeId',
            dataIndex: 'codeId',
            width: 200,
            align: 'center',
        },
        {
            title: '手机号',
            key: 'phone',
            dataIndex: 'phone',
            width: 150,
            align: 'center',
        },
        {
            title: '验证码',
            key: 'code',
            dataIndex: 'code',
            width: 120,
            align: 'center',
        },
        {
            title: '验证码类型',
            key: 'type',
            dataIndex: 'type',
            width: 150,
            align: 'center',
        },
        {
            title: '状态',
            key: 'status',
            dataIndex: 'status',
            width: 100,
            align: 'center',
        },
        {
            title: '过期时间',
            key: 'expireAt',
            dataIndex: 'expireAt',
            width: 180,
            align: 'center',
        },
        {
            title: '使用时间',
            key: 'usedAt',
            dataIndex: 'usedAt',
            width: 180,
            align: 'center',
        },
        {
            title: 'IP地址',
            key: 'ip',
            dataIndex: 'ip',
            width: 150,
            align: 'center',
        },
        {
            title: '发送次数',
            key: 'sendCount',
            dataIndex: 'sendCount',
            width: 100,
            align: 'center',
        },
        {
            title: '创建时间',
            key: 'createdAt',
            dataIndex: 'createdAt',
            width: 180,
            align: 'center',
        },
    ]
}


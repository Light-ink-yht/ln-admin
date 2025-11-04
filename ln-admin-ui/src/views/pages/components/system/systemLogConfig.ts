import type { SearchField } from '../DataList.vue'
import type { TableColumnsType } from 'ant-design-vue'

/**
 * 获取系统日志搜索字段配置
 */
export function getSystemLogSearchFields(): SearchField[] {
    return [
        {
            key: 'level',
            label: '日志级别',
            type: 'select',
            placeholder: '请选择日志级别',
            width: '150px',
            options: [
                { label: 'Info', value: 'info' },
                { label: 'Warn', value: 'warn' },
                { label: 'Error', value: 'error' },
                { label: 'Debug', value: 'debug' },
            ],
        },
        {
            key: 'module',
            label: '模块名称',
            type: 'input',
            placeholder: '请输入模块名称',
            width: '180px',
            multiple: true,
        },
        {
            key: 'action',
            label: '操作类型',
            type: 'input',
            placeholder: '请输入操作类型',
            width: '180px',
            multiple: true,
        },
        {
            key: 'userId',
            label: '用户ID',
            type: 'input',
            placeholder: '请输入用户ID',
            width: '180px',
            multiple: true,
        },
        {
            key: 'ip',
            label: 'IP地址',
            type: 'input',
            placeholder: '请输入IP地址',
            width: '180px',
            multiple: true,
        },
        {
            key: 'dateRange',
            label: '时间范围',
            type: 'dateRange',
            width: '380px',
            placeholder: ['开始时间', '结束时间'],
        },
    ]
}

/**
 * 获取系统日志表格列配置
 */
export function getSystemLogColumns(): TableColumnsType {
    return [
        {
            title: '日志级别',
            key: 'level',
            dataIndex: 'level',
            width: 100,
            align: 'center',
            fixed: 'left',
        },
        {
            title: '模块',
            key: 'module',
            dataIndex: 'module',
            width: 150,
            ellipsis: true,
            align: 'center',
        },
        {
            title: '操作类型',
            key: 'action',
            dataIndex: 'action',
            width: 120,
            align: 'center',
        },
        {
            title: '日志内容',
            key: 'content',
            dataIndex: 'content',
            width: 300,
            ellipsis: true,
            align: 'center',
        },
        {
            title: '用户ID',
            key: 'userId',
            dataIndex: 'userId',
            width: 120,
            ellipsis: true,
            align: 'center',
        },
        {
            title: 'IP地址',
            key: 'ip',
            dataIndex: 'ip',
            width: 130,
            align: 'center',
        },
        {
            title: '请求方法',
            key: 'method',
            dataIndex: 'method',
            width: 100,
            align: 'center',
        },
        {
            title: '状态码',
            key: 'statusCode',
            dataIndex: 'statusCode',
            width: 100,
            align: 'center',
        },
        {
            title: '日志时间',
            key: 'logTime',
            dataIndex: 'logTime',
            width: 180,
            align: 'center',
        },
        {
            title: '操作',
            key: 'actions',
            width: 100,
            align: 'center',
            fixed: 'right',
        },
    ]
}

